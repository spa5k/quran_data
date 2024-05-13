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

type Data struct {
	Sajdas Sajdas `json:"sajdas"`
}

type Sajdas struct {
	References []SajdasReference `json:"references"`
}

type SajdasReference struct {
	Sajda       int64 `json:"sajda"`
	Chapter     int64 `json:"chapter"`
	Verse       int64 `json:"verse"`
	Recommended bool  `json:"recommended"`
	Obligatory  bool  `json:"obligatory"`
}

func FetchAndInsertSajdah() {
	db, err := sql.Open("sqlite", "./db/quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Fetch the Sajdah data from the API
	resp, err := http.Get("https://cdn.jsdelivr.net/gh/fawazahmed0/quran-api@1/info.json")
	if err != nil {
		log.Fatal("Error fetching data:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}

	// Prepare the SQL statement for inserting data
	stmt, err := tx.Prepare("INSERT INTO sajdah (sajdah_number, surah_number, ayah_number, recommended, obligatory) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	// Insert each Sajdah into the database
	for _, sajdah := range data.Sajdas.References {
		_, err = stmt.Exec(sajdah.Sajda, sajdah.Chapter, sajdah.Verse, boolToInt(sajdah.Recommended), boolToInt(sajdah.Obligatory))
		if err != nil {
			tx.Rollback()
			log.Fatal("Error executing insert:", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Sajdah data inserted successfully")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
