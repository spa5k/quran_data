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
	Surah     int
	Ayah      int
	Text      string
	Tajweed   string
	EditionID int
}

var sources = []string{"ara-quranindopak", "ara-quranuthmanihaf", "ara-quranwarsh", "ara-quranqaloon", "ara-quransimple"}

func FetchAndInsertQuranText() {
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	for _, source := range sources {
		println("Fetching data for", source)

		// Start a new transaction for each source
		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Error beginning transaction:", err)
		}

		stmt, err := tx.Prepare("INSERT OR IGNORE INTO ayah (surah_number, ayah_number, edition_id, text) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal("Error preparing statement:", err)
		}
		defer stmt.Close()

		var sourceID int
		err = db.QueryRow("SELECT id FROM edition WHERE name = ?", source).Scan(&sourceID)
		if err != nil {
			log.Printf("Error fetching source ID for %s: %v", source, err)
			tx.Rollback()
			continue
		}

		url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/quran-api@1/editions/%s.json", source)
		data, err := fetchQuran(url)
		if err != nil {
			log.Fatal("Error fetching data:", err)
		}

		for _, ayah := range data.Quran {
			_, err = stmt.Exec(ayah.Chapter, ayah.Verse, sourceID, ayah.Text)
			if err != nil {
				log.Printf("Error executing insert for %d:%d: %v", ayah.Chapter, ayah.Verse, err)
				tx.Rollback()
				break
			}
		}

		stmt, err = tx.Prepare("UPDATE edition SET enabled = 1 WHERE name = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for _, edition := range sources {
			_, err := stmt.Exec(edition)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
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
