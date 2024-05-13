package ayah

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

type Verse struct {
	VerseKey        string `json:"verse_key"`
	VerseNumber     int    `json:"verse_number"`
	HizbNumber      int    `json:"hizb_number"`
	ManzilNumber    int    `json:"manzil_number"`
	PageNumber      int    `json:"page_number"`
	RubElHizbNumber int    `json:"rub_el_hizb_number"`
	RukuNumber      int    `json:"ruku_number"`
	ChapterNumber   int    // This will be set manually
}

type AyahAPIResponse struct {
	Verses []Verse `json:"verses"`
}

func FetchAndInsertAyahInfo() {
	db, err := sql.Open("sqlite", "./db/quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}

	stmt, err := tx.Prepare("INSERT INTO ayah_info (surah_number, ayah_number, ayah_key, hizb, rub_el_hizb, ruku, manzil, page, juz) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	for chapter := 1; chapter <= 114; chapter++ {
		url := fmt.Sprintf("https://api.quran.com/api/v4/verses/by_chapter/%d?language=en&words=false&page=1&per_page=350", chapter)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Error fetching data:", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
		}

		var apiResponse AyahAPIResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Fatal("Error unmarshaling JSON:", err)
		}

		for _, verse := range apiResponse.Verses {
			verse.ChapterNumber = chapter
			_, err = stmt.Exec(verse.ChapterNumber, verse.VerseNumber, verse.VerseKey, verse.HizbNumber, verse.RubElHizbNumber, verse.RukuNumber, verse.ManzilNumber, verse.PageNumber, strconv.Itoa((verse.PageNumber-1)/20+1))
			if err != nil {
				tx.Rollback()
				log.Fatal("Error executing insert:", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Ayah data inserted successfully")
}
