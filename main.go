package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func main() {
	start := time.Now()

	// Global variable configuration setup.
	viper.SetDefault("tempDir", "./.temp")
	viper.SetDefault("distDir", "./.dist")

	// Get text/epub file.

	// Set up output directories.

	// Loop through chapters.
	// Skip chapter if already created.
	// Check if the chaper was already in progress.
	// Split the chapter into audio segments
	// Output segments as .wav files
	// Concat .wav segments into a singular .wav file
	// Remove .wav segments

	// Combine all chapter .wav files FFmpeg.

	// Insert audiobook metadata and image.

	// Remove chapter .wav files

	// Delete temp files.

	fmt.Println("Audiobook created!", time.Since(start))
}
