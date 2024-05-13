package tajweed

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"database/sql"

	_ "modernc.org/sqlite"
)

type Tajweed struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Surahs  []Surah `json:"surahs"`
	Edition Edition `json:"edition"`
}

type Edition struct {
	Identifier  string `json:"identifier"`
	Language    string `json:"language"`
	Name        string `json:"name"`
	EnglishName string `json:"englishName"`
	Format      string `json:"format"`
	Type        string `json:"type"`
}

type Surah struct {
	Number                 int    `json:"number"`
	Name                   string `json:"name"`
	EnglishName            string `json:"englishName"`
	EnglishNameTranslation string `json:"englishNameTranslation"`
	RevelationType         string `json:"revelationType"`
	Ayahs                  []Ayah `json:"ayahs"`
}

type Ayah struct {
	Number        int         `json:"number"`
	Text          string      `json:"text"`
	NumberInSurah int         `json:"numberInSurah"`
	Juz           int         `json:"juz"`
	Manzil        int         `json:"manzil"`
	Page          int         `json:"page"`
	Ruku          int         `json:"ruku"`
	HizbQuarter   int         `json:"hizbQuarter"`
	Sajda         interface{} `json:"sajda"`
}

type TajweedRes struct {
	Ayah    int
	Surah   int
	Tajweed string
}

// Structs as defined previously

func FetchAndInsertTajweed() {
	db, err := sql.Open("sqlite", "./db/quran.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Fetch Tajweed data
	tajweedData, err := fetchTajweed()
	if err != nil {
		log.Fatal("Error fetching tajweed data:", err)
	}

	// Begin a transaction for database updates
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error beginning transaction:", err)
	}

	// Prepare the SQL statement for updating data
	stmt, err := tx.Prepare("UPDATE ayah SET tajweed = ? WHERE surah_number = ? AND ayah_number = ?")
	if err != nil {
		log.Fatal("Error preparing statement:", err)
	}
	defer stmt.Close()

	// Update each ayah with the fetched Tajweed text
	for _, ayah := range tajweedData {
		_, err = stmt.Exec(ayah.Tajweed, ayah.Surah, ayah.Ayah)
		if err != nil {
			tx.Rollback()
			log.Fatal("Error executing update:", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Tajweed data inserted successfully")
}

func fetchTajweed() ([]TajweedRes, error) {
	url := "https://api.alquran.cloud/v1/quran/quran-tajweed"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var tajweed Tajweed
	err = json.Unmarshal(body, &tajweed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	var finalRes []TajweedRes
	for _, surah := range tajweed.Data.Surahs {
		for _, ayah := range surah.Ayahs {
			res := TajweedRes{
				Surah:   surah.Number,
				Ayah:    ayah.NumberInSurah,
				Tajweed: ayah.Text,
			}
			finalRes = append(finalRes, res)
		}
	}
	return finalRes, nil
}
