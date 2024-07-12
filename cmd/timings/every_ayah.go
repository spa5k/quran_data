package timings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func FetchEveryAyahTimings() {
	source := "everyayah"
	jsonURL := "https://raw.githubusercontent.com/spa5k/quran_timings_api/master/everyayah/reciters.json"
	resp, err := http.Get(jsonURL)
	if err != nil {
		log.Fatalf("Failed to fetch JSON data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse JSON data
	var recitations Recitations
	err = json.Unmarshal(body, &recitations)
	if err != nil {
		log.Fatalf("Failed to parse JSON data: %v", err)
	}

	// Set up SQLite database
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Insert data into database
	insertSQL := `
	INSERT INTO reciters (source_id, reciter_name, style, slug, translated_name, language_name, source)
	VALUES (?, ?, ?, ?, ?, ?, ?);`
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	for _, recitation := range recitations.Recitations {
		modifiedID := recitation.ID + 1000
		_, err = stmt.Exec(
			modifiedID,
			recitation.ReciterName,
			recitation.Style,
			recitation.Slug,
			recitation.TranslatedName.Name,
			recitation.TranslatedName.LanguageName,
			source,
		)
		if err != nil {
			log.Printf("Failed to insert recitation (ID: %d): %v", recitation.ID, err)
		}
	}

	fmt.Println("Data successfully inserted into the database.")
}
