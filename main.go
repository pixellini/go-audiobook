package main

import (
	"fmt"
	"time"

	"github.com/pixellini/go-audiobook/epub"
	"github.com/pixellini/go-audiobook/fsutils"
	"github.com/spf13/viper"
)

const processingFileType = "wav"

func main() {
	start := time.Now()

	// Global variable configuration setup.
	viper.SetDefault("tempDir", "./.temp")
	viper.SetDefault("distDir", "./.dist")

	// Get text/epub file.
	book, err := epub.New("test/Writing a Compiler in Go - Thorsten Ball.epub")
	if err != nil {
		panic(err)
	}

	fmt.Println(book.Title)
	fmt.Println(book.Author)
	fmt.Println(book.Description)
	fmt.Println(book.Language)
	fmt.Println(book.Introduction)

	// Set up output directories.
	tempDir := viper.GetString("tempDir")
	distDir := viper.GetString("distDir")

	fsutils.CreateDirIfNotExist(tempDir)
	fsutils.CreateDirIfNotExist(distDir)

	// Loop through chapters.
	for i, chapter := range book.Chapters {
		// Skip chapter if already created.
		if len(chapter.Paragraphs) == 0 {
			continue
		}

		// eg: ./.dist/Chapter Name-1.wav
		chapterFile := fmt.Sprintf("%s/%s-%d.wav", distDir, book.Title, i)

		// Check if the chaper was already in progress.
		if _, err := fsutils.FileExists(chapterFile); err != nil {
			fmt.Println(chapterFile, "already exists. Skipping.")
			continue
		}

		// Split the chapter into audio segments
		fmt.Println("Processing Chapter:", chapter.Title)
		// process here...

		// Output segments as .wav files
		chapterAudioSegmentFiles, err := fsutils.GetFilesFrom(tempDir, processingFileType)
		if err != nil {
			panic(err)
		}

		fsutils.SortFilesNumerically(chapterAudioSegmentFiles)

		// Concat .wav segments into a singular .wav file
		// chapter creation...

		// Remove .wav segments
		fsutils.RemoveAllFilesInDir(tempDir)
	}

	chapterAudioFiles, err := fsutils.GetFilesFrom(distDir, processingFileType)
	if err != nil {
		panic(err)
	}
	// Combine all chapter .wav files FFmpeg.

	// Insert audiobook metadata and image.

	// Remove chapter .wav files
	err = fsutils.RemoveFilesFrom(distDir, chapterAudioFiles)
	if err != nil {
		fmt.Println("Unable to remove chapter audio files from dist directory.")
	}

	// Delete temp files.
	err = fsutils.RemoveDirIfEmpty(tempDir)
	if err != nil {
		fmt.Println("Unable to remove temp directory.")
	}

	fmt.Println("Audiobook created!", time.Since(start))
}
