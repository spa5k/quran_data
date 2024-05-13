package juz

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
	Verses   Verses    `json:"verses"`
	Chapters []Chapter `json:"chapters"`
	Sajdas   Sajdas    `json:"sajdas"`
	Rukus    Rukus     `json:"rukus"`
	Pages    Pages     `json:"pages"`
	Manzils  Manzils   `json:"manzils"`
	Maqras   Maqras    `json:"maqras"`
	Juzs     Juzs      `json:"juzs"`
}

type Chapter struct {
	Chapter     int64      `json:"chapter"`
	Name        string     `json:"name"`
	Englishname string     `json:"englishname"`
	Arabicname  string     `json:"arabicname"`
	Revelation  Revelation `json:"revelation"`
	Verses      []Verse    `json:"verses"`
}

type Verse struct {
	Verse  int64       `json:"verse"`
	Line   int64       `json:"line"`
	Juz    int64       `json:"juz"`
	Manzil int64       `json:"manzil"`
	Page   int64       `json:"page"`
	Ruku   int64       `json:"ruku"`
	Maqra  int64       `json:"maqra"`
	Sajda  *SajdaUnion `json:"sajda"`
}

type SajdaClass struct {
	No          int64 `json:"no"`
	Recommended bool  `json:"recommended"`
	Obligatory  bool  `json:"obligatory"`
}

type Juzs struct {
	Count      int64           `json:"count"`
	References []JuzsReference `json:"references"`
}

type JuzsReference struct {
	Juz   int64 `json:"juz"`
	Start End   `json:"start"`
	End   End   `json:"end"`
}

type End struct {
	Chapter int64 `json:"chapter"`
	Verse   int64 `json:"verse"`
}

type Manzils struct {
	Count      int64              `json:"count"`
	References []ManzilsReference `json:"references"`
}

type ManzilsReference struct {
	Manzil int64 `json:"manzil"`
	Start  End   `json:"start"`
	End    End   `json:"end"`
}

type Maqras struct {
	Count      int64             `json:"count"`
	References []MaqrasReference `json:"references"`
}

type MaqrasReference struct {
	Maqra int64 `json:"maqra"`
	Start End   `json:"start"`
	End   End   `json:"end"`
}

type Pages struct {
	Count      int64            `json:"count"`
	References []PagesReference `json:"references"`
}

type PagesReference struct {
	Page  int64 `json:"page"`
	Start End   `json:"start"`
	End   End   `json:"end"`
}

type Rukus struct {
	Count      int64            `json:"count"`
	References []RukusReference `json:"references"`
}

type RukusReference struct {
	Ruku  int64 `json:"ruku"`
	Start End   `json:"start"`
	End   End   `json:"end"`
}

type Sajdas struct {
	Count      int64             `json:"count"`
	References []SajdasReference `json:"references"`
}

type SajdasReference struct {
	Sajda       int64 `json:"sajda"`
	Chapter     int64 `json:"chapter"`
	Verse       int64 `json:"verse"`
	Recommended bool  `json:"recommended"`
	Obligatory  bool  `json:"obligatory"`
}

type Verses struct {
	Count int64 `json:"count"`
}

type Revelation string

const (
	Madina Revelation = "Madina"
	Mecca  Revelation = "Mecca"
)

type SajdaUnion struct {
	Bool       *bool
	SajdaClass *SajdaClass
}

func (su *SajdaUnion) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a boolean
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		su.Bool = &b
		return nil
	}

	// Try to unmarshal as a SajdaClass
	var sc SajdaClass
	if err := json.Unmarshal(data, &sc); err == nil {
		su.SajdaClass = &sc
		return nil
	}

	return fmt.Errorf("data could not be unmarshaled as either bool or SajdaClass: %s", data)
}

func DownloadAndInsertJuz() {
	// Open the database connection
	db, err := sql.Open("sqlite", "./db/quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Fetch the data from the API
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

	// 	CREATE TABLE "juz" (
	//     id INTEGER not null primary key autoincrement,
	//     juz_number INTEGER not null,
	//     start_surah INTEGER not null,
	//     start_ayah INTEGER not null,
	//     end_surah INTEGER not null,
	//     end_ayah INTEGER not null
	// );
	// Example of inserting chapters
	stmt, err := tx.Prepare("INSERT INTO juz (juz_number, start_surah, start_ayah, end_surah, end_ayah) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	for _, chapter := range data.Chapters {
		_, err = stmt.Exec(chapter.Chapter, chapter.Verses[0].Juz, chapter.Verses[0].Verse, chapter.Verses[len(chapter.Verses)-1].Juz, chapter.Verses[len(chapter.Verses)-1].Verse)
		if err != nil {
			tx.Rollback()
			log.Fatal("Error executing insert:", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Data inserted successfully")
}
