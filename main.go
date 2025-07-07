package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pixellini/go-audiobook/internal/audiobook"
	"github.com/pixellini/go-audiobook/internal/audioprocessor"
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
	tempDir := viper.GetString("temp_dir")
	distDir := viper.GetString("dist_dir")
	fsutils.CreateDirIfNotExist(tempDir)
	fsutils.CreateDirIfNotExist(distDir)
	return tempDir, distDir
}

// Get epub file.
func getEpubFileAndImage() (*epub.Epub, string) {
	// Get epub file path from config.
	epupPath := viper.GetString("epub_path")
	if epupPath == "" {
		panic("Missing required config value: 'epub_path' in config.json")
	}

	// Get image path from config.
	imgPath := viper.GetString("image_path")
	if imgPath == "" {
		// no image set so skip
	}

	book, err := epub.New(epupPath)
	if err != nil {
		panic(err)
	}
	return book, imgPath
}

// eg: ./.dist/Book Title-1.wav
func chapterFilePath(distDir, bookTitle string, chapterIdx int) string {
	return fmt.Sprintf("%s/%s-%d.%s", distDir, bookTitle, chapterIdx, processingFileType)
}

func generateChapterAudioFiles(epubBook *epub.Epub, Audiobook *audiobook.Audiobook, tempDir string, distDir string) {
	if len(epubBook.Chapters) == 0 {
		panic("No chapters found in epub file")
	}

	// Loop through chapters.
	for i, chapter := range epubBook.Chapters {
		// Skip chapter if already created.
		if len(chapter.Paragraphs) == 0 {
			continue
		}

		// eg: ./.dist/Book Title-1.wav
		chapterFile := chapterFilePath(distDir, epubBook.Title, i)

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
		tts.SynthesizeTextList(chapter.Paragraphs, epubBook.Language)

		// Output segments as .wav files
		chapterAudioSegmentFiles, err := fsutils.GetFilesFrom(tempDir, processingFileType)
		if err != nil {
			panic(err)
		}

		fsutils.SortFilesNumerically(chapterAudioSegmentFiles)

		// Combine all .wav paragraph files into a singular .wav file
		// This will become our Chapter wav file.
		err = audioprocessor.ConcatFiles(epubBook.Title, chapterAudioSegmentFiles, chapterFile)
		if err != nil {
			panic(err)
		}

		Audiobook.AddChapter(audiobook.AudiobookChapter{
			Title: chapter.Title,
			File:  chapterFile,
		})

		// Remove .wav paragraph files
		fsutils.RemoveAllFilesInDir(tempDir)
		fmt.Println("Processing: Chapter complete ✅")
	}
}

func main() {
	start := time.Now()

	loadConfig()

	book, image := getEpubFileAndImage()

	audiobook := audiobook.NewFromEpub(book, image)

	tempDir, distDir := setupDirs()

	ttsModel := tts.ModelXTTS
	if viper.GetBool("tts.use_vits") {
		ttsModel = tts.ModelVITS
	}
	fmt.Printf("Using %s model for TTS.\n", ttsModel)

	generateChapterAudioFiles(book, audiobook, tempDir, distDir)

	fmt.Println("\n\n--------------------------------------------------")

	// Combine all chapter .wav files FFmpeg.
	err := audiobook.Generate(distDir)
	if err != nil {
		panic(err)
	}

	// Delete temp files.
	err = fsutils.RemoveDirIfEmpty(tempDir)
	if err != nil {
		fmt.Println("Unable to remove temp directory.")
	}

	fmt.Println("Audiobook created!", time.Since(start))
}
