package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pixellini/go-audiobook/internal/audiobook"
	"github.com/pixellini/go-audiobook/internal/audioprocessor"
	"github.com/pixellini/go-audiobook/internal/device"
	"github.com/pixellini/go-audiobook/internal/epub"
	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/tts"
	"github.com/spf13/viper"
)

const processingFileType = "wav"

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Set default values if not present in config.
	viper.SetDefault("image_path", "")
	viper.SetDefault("temp_dir", "./.temp")
}

// Set up output directories.
func setupDirs() (string, string) {
	tempDir, err := fsutils.GetOrCreateTempDir()
	if err != nil {
		// We can't process the audiobook without a temp directory.
		panic(err)
	}

	distDir := viper.GetString("dist_dir")
	fsutils.CreateDirIfNotExist(distDir)
	return tempDir, distDir
}

func getEpubFile() *epub.Epub {
	// Get epub file path from config.
	epupPath := viper.GetString("epub_path")
	if epupPath == "" {
		panic("Missing required config value: 'epub_path' in config.json")
	}

	book, err := epub.New(epupPath)
	if err != nil {
		// We need a book to process the audiobook.
		panic(err)
	}
	return book
}

// eg: tempDir/Chapter Title-1.wav
func chapterFilePath(chaptersDir, bookTitle string, chapterIdx int) string {
	return fmt.Sprintf("%s/%s-%d.%s", chaptersDir, bookTitle, chapterIdx, processingFileType)
}

func generateChapterAudioFiles(epubBook *epub.Epub, Audiobook *audiobook.Audiobook, tempDir string) {
	if len(epubBook.Chapters) == 0 {
		panic("No chapters found in epub file")
	}

	chaptersDir := fmt.Sprintf("%s/chapters", tempDir)
	fsutils.CreateDirIfNotExist(chaptersDir)


	// Loop through chapters.
	for i, chapter := range epubBook.Chapters {
		// Skip chapter if already created.
		if len(chapter.Paragraphs) == 0 {
			continue
		}

		// eg: tempDir/Chapter Title-1.wav
		chapterFile := chapterFilePath(chaptersDir, epubBook.Title, i)

		// Check if the chaper was already in progress.
		// We can do this by checking if the .wav file has already been created.
		if fsutils.FileExists(chapterFile) {
			Audiobook.AddChapter(audiobook.AudiobookChapter{
				Title: chapter.Title,
				File:  chapterFile,
			})
			fmt.Println(chapterFile, "already exists. Skipping.")
			continue
		}

		// Split the chapter into audio segments
		fmt.Println("\n\nProcessing Chapter:", chapter.Title)
		fmt.Println("--------------------------------------------------")
		tts.SynthesizeTextList(tempDir, chapter.Paragraphs, epubBook.Language)

		// Output segments as .wav files
		chapterAudioSegmentFiles, err := fsutils.GetFilesFrom(tempDir, processingFileType)
		if err != nil {
			panic(err)
		}

		fsutils.SortFilesNumerically(chapterAudioSegmentFiles)

		// Combine all .wav paragraph files into a singular .wav file
		// This will become our Chapter wav file.
		fmt.Println("Chapter file:", chapterFile)
		err = audioprocessor.ConcatFiles(epubBook.Title, chapterAudioSegmentFiles, chapterFile)
		if err != nil {
			panic(err)
		}

		Audiobook.AddChapter(audiobook.AudiobookChapter{
			Title: chapter.Title,
			File:  chapterFile,
		})

		// Clean up tempDir (remove all part-*.wav files after chapter is created)
		fsutils.RemoveAllFilesInDir(tempDir)
		fmt.Println("Processing: Chapter complete ✅")
	}
}

func main() {
	start := time.Now()

	loadConfig()

	image := viper.GetString("image_path")
	book := getEpubFile()

	audiobook := audiobook.NewFromEpub(book, image)

	tempDir, distDir := setupDirs()

	device.Manager.Init()
	tts.BaseCoquiTTSConfig.Init(book.Language)

	generateChapterAudioFiles(book, audiobook, tempDir)

	fmt.Println("\n\n--------------------------------------------------")

	// Combine all chapter .wav files FFmpeg.
	err := audiobook.Generate(distDir)
	if err != nil {
		panic(err)
	}

	// Clean up chapters directory
	chaptersDir := fmt.Sprintf("%s/chapters", tempDir)
	err = fsutils.RemoveAllFilesInDir(chaptersDir)
	if err == nil {
		_ = fsutils.RemoveDirIfEmpty(chaptersDir)
	}

	// Clean up temp dir
	err = fsutils.RemoveDirIfEmpty(tempDir)
	if err != nil {
		fmt.Println("Unable to remove temp directory.")
	}

	fmt.Println("Audiobook created!", time.Since(start))
}
