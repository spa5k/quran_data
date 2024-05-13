package main

import (
	"github.com/spa5k/quran_data/cmd/surah"
)

func main() {
	// Open the database
	// cmd/juzToSurah.json
	// Load Juz mappings
	// Fetch editions from the database
	// Process each edition
	// Assuming VerseText is defined correctly to match the JSON structure
	// Start transaction
	// Commit transaction
	// editions.InsertEditionsData()
	// translations.InsertTranslationsData()
	// juz.DownloadAndInsertJuz()
	// ayah.FetchAndInsertAyahInfo()
	// ayah.FetchAndInsertSajdah()
	// ayah.FetchAndInsertQuranText()
	// tajweed.FetchAndInsertTajweed()
	surah.FetchAndInsertSurahs()
}
