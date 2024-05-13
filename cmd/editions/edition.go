package editions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

type Edition struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	Direction string `json:"direction"`
	Source    string `json:"source"`
	Type      string `json:"type"`
	Enabled   int    `json:"enabled"`
}
type Type string

const (
	QuranTransliteration Type = "QURAN_TRANSLITERATION"
	Quran                Type = "QURAN"
	Transliteration      Type = "TRANSLITERATION"
	Translation          Type = "TRANSLATION"
)

func GetEditionType(name string) Type {
	quranTransliteration := regexp.MustCompile(`^ara_quran.*la$`)
	quranRegex := regexp.MustCompile(`^ara_quran`)

	switch {
	case quranTransliteration.MatchString(name) || (strings.HasPrefix(name, "ara_quran") && strings.Contains(name, "la")):
		return QuranTransliteration
	case quranRegex.MatchString(name) && !strings.Contains(name, "la"):
		return Quran
	case !strings.HasPrefix(name, "ara_quran") && strings.HasSuffix(name, "la"):
		return Transliteration
	default:
		return Translation
	}
}

func InsertEditionsData() {
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	resp, err := http.Get("https://cdn.jsdelivr.net/gh/fawazahmed0/quran-api@1/editions.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]Edition
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO edition (name, author, language, direction, source, type, enabled) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for key, edition := range data {
		editionType := GetEditionType(key)
		_, err := stmt.Exec(edition.Name, edition.Author, edition.Language, edition.Direction, edition.Source, editionType, 0)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Editions inserted successfully")
}
