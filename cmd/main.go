package main

import (
	"fmt"
	"log"

	"github.com/spa5k/quran_data/cmd/ayah"
	"github.com/spa5k/quran_data/cmd/editions"
	"github.com/spa5k/quran_data/cmd/juz"
	"github.com/spa5k/quran_data/cmd/surah"
	"github.com/spa5k/quran_data/cmd/tajweed"
	"github.com/spa5k/quran_data/cmd/timings"
	"github.com/spa5k/quran_data/cmd/translations"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "quran_data"}
	var edition string

	translationsCmd := &cobra.Command{
		Use:   "translations",
		Short: "Fetch and insert translations data",
		Run: func(cmd *cobra.Command, args []string) {
			translations.InsertTranslationsData(&edition)
		},
	}
	translationsCmd.Flags().StringVarP(&edition, "edition", "e", "", "Specify the edition for which to fetch translations")

	rootCmd.AddCommand(translationsCmd)
	rootCmd.AddCommand(
		makeCmd("editions", "Insert editions data", editions.InsertEditionsData),
		makeCmd("juz", "Download and insert Juz data", juz.DownloadAndInsertJuz),
		makeCmd("ayahinfo", "Fetch and insert Ayah info", ayah.FetchAndInsertAyahInfo),
		makeCmd("sajdah", "Fetch and insert Sajdah info", ayah.FetchAndInsertSajdah),
		makeCmd("qurantext", "Fetch and insert Quran text", ayah.FetchAndInsertQuranText),
		makeCmd("tajweed", "Fetch and insert Tajweed data", tajweed.FetchAndInsertTajweed),
		makeCmd("surahs", "Fetch and insert Surahs data", surah.FetchAndInsertSurahs),
		makeCmd("timings_quran", "Fetch and insert timings data", timings.FetchQuranComAyahTimings),
		makeCmd("timings_every_ayah", "Fetch and insert timings data", timings.FetchEveryAyahTimings),
		makeCmd("recitations_every_ayah", "Fetch and insert timings data", timings.FetchAndSaveEveryAyahRecitations),
		makeCmd("recitations_quran_com", "Fetch and insert timings data", timings.FetchAndSaveDataFolderRecitations),
		makeCmd("all", "Run all data import functions sequentially", runAll),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func makeCmd(name, description string, action func()) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			action()
		},
	}
}

func runAll() {
	fmt.Println("Running all data import functions sequentially...")
	editions.InsertEditionsData()
	translations.InsertTranslationsData(nil)
	juz.DownloadAndInsertJuz()
	ayah.FetchAndInsertAyahInfo()
	ayah.FetchAndInsertSajdah()
	ayah.FetchAndInsertQuranText()
	tajweed.FetchAndInsertTajweed()
	surah.FetchAndInsertSurahs()
	timings.FetchQuranComAyahTimings()
	timings.FetchEveryAyahTimings()

	// Fetch and save timings data from Data Folder
	timings.FetchAndSaveDataFolderRecitations()
	timings.FetchAndSaveEveryAyahRecitations()

	fmt.Println("All data import functions completed.")
}
