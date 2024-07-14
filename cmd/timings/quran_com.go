package timings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Recitation struct {
	ID             int    `json:"id"`
	ReciterName    string `json:"reciter_name"`
	Style          string `json:"style"`
	Slug           string `json:"slug"`
	TranslatedName struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"translated_name"`
}

type Recitations struct {
	Recitations []Recitation `json:"recitations"`
}

func FetchQuranComAyahTimings() {
	source := "quran_com"
	jsonURL := "https://raw.githubusercontent.com/spa5k/quran_timings_api/master/data/reciters.json"

	// Fetch JSON data
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
		_, err = stmt.Exec(
			recitation.ID,
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

// Define the structure of the API response for verses
type Verse struct {
	Chapter  int     `json:"chapter"`
	Segments [][]int `json:"segments"`
	URL      string  `json:"url"`
	Verse    int     `json:"verse"`
}

type VersesResponse struct {
	Verses []Verse `json:"verses"`
}

func FetchAndSaveDataFolderRecitations() {
	source := "quran_com"

	// Set up SQLite database
	db, err := sql.Open("sqlite", "quran.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Fetch reciters from the database
	recitersQuery := `SELECT id, slug FROM reciters WHERE source = ?`
	rows, err := db.Query(recitersQuery, source)
	if err != nil {
		log.Fatalf("Failed to fetch reciters from database: %v", err)
	}
	defer rows.Close()

	var reciters []Recitation
	for rows.Next() {
		var reciter Recitation
		err := rows.Scan(&reciter.ID, &reciter.Slug)
		if err != nil {
			log.Fatalf("Failed to scan reciter row: %v", err)
		}
		reciters = append(reciters, reciter)
	}

	// Channel to collect recitation data
	type RecitationData struct {
		ReciterID      int
		SurahNumber    int
		RecitationData []int
		URL            string
		Err            error
	}
	recitationChan := make(chan RecitationData, len(reciters)*114)

	var wg sync.WaitGroup

	// Base URLs for round-robin
	baseURLs := []string{
		"https://cdn.jsdelivr.net/gh/spa5k/quran_timings_api@master/data",
		"https://cdn.statically.io/gh/spa5k/quran_timings_api/master/data",
		"https://raw.githubusercontent.com/spa5k/quran_timings_api/master/data",
	}
	baseURLIndex := 0

	// Fetch recitations concurrently in batches
	for _, reciter := range reciters {
		for surahNumber := 1; surahNumber <= 114; surahNumber++ {
			wg.Add(1)
			go func(reciterID int, slug string, surahNumber int) {
				defer wg.Done()

				cacheKey := fmt.Sprintf("%d-%d", reciterID, surahNumber)
				cacheMutex.Lock()
				if cache[cacheKey] {
					cacheMutex.Unlock()
					fmt.Printf("Data for recitation %d chapter %d already fetched, skipping...\n", reciterID, surahNumber)
					return
				}
				cache[cacheKey] = true
				cacheMutex.Unlock()

				// Round-robin base URL selection
				baseURL := baseURLs[baseURLIndex%len(baseURLs)]
				baseURLIndex++

				recitationURL := fmt.Sprintf("%s/%s/%d.json", baseURL, slug, surahNumber)
				recitationResp, err := http.Get(recitationURL)
				if err != nil {
					recitationChan <- RecitationData{ReciterID: reciterID, SurahNumber: surahNumber, URL: recitationURL, Err: err}
					return
				}
				defer recitationResp.Body.Close()

				recitationBody, err := io.ReadAll(recitationResp.Body)
				if err != nil {
					recitationChan <- RecitationData{ReciterID: reciterID, SurahNumber: surahNumber, URL: recitationURL, Err: err}
					return
				}

				var verses []Verse
				err = json.Unmarshal(recitationBody, &verses)
				if err != nil {
					recitationChan <- RecitationData{ReciterID: reciterID, SurahNumber: surahNumber, URL: recitationURL, Err: err}
					return
				}

				// Extract the last element of all segments
				var recitationData []int
				for _, verse := range verses {
					for _, segment := range verse.Segments {
						if len(segment) > 0 {
							recitationData = append(recitationData, segment[len(segment)-1])
						}
					}
				}

				recitationChan <- RecitationData{ReciterID: reciterID, SurahNumber: surahNumber, RecitationData: recitationData, URL: recitationURL}
			}(reciter.ID, reciter.Slug, surahNumber)

			// Batch processing: wait for 5 seconds after every 10 requests
			if (surahNumber % 10) == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(recitationChan)
	}()

	// Insert recitations into the database
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	insertRecitationsSQL := `
	INSERT INTO recitations (reciter_id, surah_number, recitation_data)
	VALUES (?, ?, ?);`
	stmtRecitations, err := tx.Prepare(insertRecitationsSQL)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement for recitations: %v", err)
	}
	defer stmtRecitations.Close()

	for recitationData := range recitationChan {
		if recitationData.Err != nil {
			log.Printf("Failed to fetch recitation data from URL %s: %v", recitationData.URL, recitationData.Err)
			continue
		}
		recitationDataText, err := json.Marshal(recitationData.RecitationData)
		if err != nil {
			log.Printf("Failed to marshal recitation data for reciter (ID: %d, Surah: %d): %v", recitationData.ReciterID, recitationData.SurahNumber, err)
			continue
		}
		stringRecitationData := string(recitationDataText)
		_, err = stmtRecitations.Exec(recitationData.ReciterID, recitationData.SurahNumber, stringRecitationData)
		if err != nil {
			log.Printf("Failed to insert recitation for reciter (ID: %d, Surah: %d): %v", recitationData.ReciterID, recitationData.SurahNumber, err)
		}

	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Data successfully inserted into the database.")
}
