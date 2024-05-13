package ayah

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type Ayah struct {
	Verse   int    `json:"verse"`
	Chapter int    `json:"chapter"`
	Text    string `json:"text"`
}

type Quran struct {
	Quran []Ayah `json:"quran"`
}

type QuranText struct {
	Surah   int
	Ayah    int
	Text    string
	Tajweed string
}

var sources = []string{"ara-quranindopak", "ara-quranuthmanihaf", "ara-quranwarsh", "ara-quranqaloon", "ara-quransimple"}

func FetchAndInsertQuranText() {
	db, err := sql.Open("sqlite", "./db/quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}

	stmt, err := tx.Prepare("INSERT INTO ayah (surah_number, ayah_number, edition_id, text) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	finalRes := make([]QuranText, 0)

	// Fetch text for each source and update finalRes
	for _, source := range sources {
		url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/quran-api@1/editions/%s.json", source)
		data, err := fetchQuran(url)
		if err != nil {
			log.Fatal("Error fetching data:", err)
		}

		// sourceKey[source]
		for i, ayah := range data.Quran {
			if len(finalRes) <= i {
				finalRes = append(finalRes, QuranText{Surah: ayah.Chapter, Ayah: ayah.Verse, Text: ayah.Text})
			}
			finalRes[i].Text = ayah.Text // Assuming Uthmani text as default
		}
	}

	// Insert each Quran text into the database
	for _, ayah := range finalRes {
		_, err = stmt.Exec(ayah.Surah, ayah.Ayah, 1, ayah.Text) // Example: using Uthmani text
		if err != nil {
			tx.Rollback()
			log.Fatal("Error executing insert:", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Quran text added")
}

func fetchQuran(url string) (Quran, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Quran{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Quran{}, err
	}

	var quran Quran
	err = json.Unmarshal(body, &quran)
	if err != nil {
		return Quran{}, err
	}
	return quran, nil
}
