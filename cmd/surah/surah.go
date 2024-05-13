package surah

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type SurahsResponse struct {
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	ID              int64           `json:"id"`
	RevelationPlace RevelationPlace `json:"revelation_place"`
	RevelationOrder int64           `json:"revelation_order"`
	BismillahPre    bool            `json:"bismillah_pre"`
	NameSimple      string          `json:"name_simple"`
	NameComplex     string          `json:"name_complex"`
	NameArabic      string          `json:"name_arabic"`
	VersesCount     int64           `json:"verses_count"`
	Pages           []int64         `json:"pages"`
	TranslatedName  TranslatedName  `json:"translated_name"`
}

type TranslatedName struct {
	LanguageName LanguageName `json:"language_name"`
	Name         string       `json:"name"`
}

type (
	RevelationPlace string
	LanguageName    string
)

const (
	Madinah RevelationPlace = "madinah"
	Makkah  RevelationPlace = "makkah"
	English LanguageName    = "english"
)

func FetchAndInsertSurahs() {
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO surah (surah_number, name_simple, name_complex, name_arabic, ayah_start, ayah_end, revelation_place, page_start, page_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	resp, err := http.Get("https://api.quran.com/api/v4/chapters")
	if err != nil {
		log.Fatal("Error fetching data:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var surahsResponse SurahsResponse
	err = json.Unmarshal(body, &surahsResponse)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	for _, chapter := range surahsResponse.Chapters {
		_, err = stmt.Exec(chapter.ID, chapter.NameSimple, chapter.NameComplex, chapter.NameArabic, 1, chapter.VersesCount, chapter.RevelationPlace, chapter.Pages[0], chapter.Pages[len(chapter.Pages)-1]) // Juz number is not provided
		if err != nil {
			tx.Rollback()
			log.Fatal("Error executing insert:", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Surah data inserted successfully")
}
