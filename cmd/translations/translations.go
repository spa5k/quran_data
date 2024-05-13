package translations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type JuzMapping struct {
	Juz   int   `json:"juz"`
	Start Verse `json:"start"`
	End   Verse `json:"end"`
}

type Verse struct {
	Chapter int `json:"chapter"`
	Verse   int `json:"verse"`
}

type EditionData struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Source  string `json:"source"`
	Enabled int    `json:"enabled"`
}

type QuranText struct {
	Quran []Quran `json:"quran"`
}

type Quran struct {
	Chapter int64  `json:"chapter"`
	Verse   int64  `json:"verse"`
	Text    string `json:"text"`
}

type QuranVerses struct {
	Quran []struct {
		Chapter int    `json:"chapter"`
		Verse   int    `json:"verse"`
		Text    string `json:"text"`
	} `json:"quran"`
}

func loadJuzMappings(filename string) ([]JuzMapping, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mappings []JuzMapping
	err = json.NewDecoder(file).Decode(&mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

func findJuzNumber(mappings []JuzMapping, surah, ayah int) int {
	for _, mapping := range mappings {
		if (surah > mapping.Start.Chapter || (surah == mapping.Start.Chapter && ayah >= mapping.Start.Verse)) &&
			(surah < mapping.End.Chapter || (surah == mapping.End.Chapter && ayah <= mapping.End.Verse)) {
			return mapping.Juz
		}
	}
	return 0 // Default to 0 if no matching Juz is found
}

// Enable these editions before running this function
var enabledEditions = []string{
	"eng-mustafakhattaba",
	"eng-ummmuhammad",
}

func enableEditions(extraEdition *string) {
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("UPDATE edition SET enabled = 1 WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, edition := range enabledEditions {
		_, err := stmt.Exec(edition)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// Enable extra edition if provided
	if extraEdition != nil && *extraEdition != "" {
		println("Enabling extra edition " + *extraEdition)
		_, err := stmt.Exec(*extraEdition)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enabled editions successfully")
}

func InsertTranslationsData(extraEdition *string) {
	enableEditions(extraEdition)

	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	juzMappings, err := loadJuzMappings("./cmd/juzToSurah.json")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, name, source, enabled FROM edition WHERE enabled = 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var editions []EditionData
	for rows.Next() {
		var edition EditionData
		if err := rows.Scan(&edition.ID, &edition.Name, &edition.Source, &edition.Enabled); err != nil {
			log.Fatal(err)
		}
		editions = append(editions, edition)
	}

	for _, edition := range editions {
		editionName := strings.ReplaceAll(edition.Name, " ", "-")
		url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/fawazahmed0/quran-api@1/editions/%s.json", editionName)

		// if its not enabled, skip
		if edition.Enabled == 0 {
			continue
		}

		println("Fetching data for edition " + editionName + " from " + url)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var quranText QuranText
		if err := json.Unmarshal(body, &quranText); err != nil {
			log.Fatal(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := tx.Prepare("INSERT OR IGNORE INTO translation (surah_number, ayah_number, edition_id, text, juz_number) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		verses := quranText.Quran

		for _, verse := range verses {
			juzNumber := findJuzNumber(juzMappings, int(verse.Chapter), int(verse.Verse))
			_, err := stmt.Exec(verse.Chapter, verse.Verse, edition.ID, verse.Text, juzNumber)
			if err != nil {
				tx.Rollback()
				log.Fatal(err)
			}
		}

		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Translations inserted successfully")
}
