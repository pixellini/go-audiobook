package tts

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/utils"
	"github.com/spf13/viper"
)

const defaultMaxRetries = 1
const defaultParallelAudioCount = 1

func SynthesizeText(text, language, outputFile string) error {
	_, err := os.Stat(outputFile)
	if err == nil {
		fmt.Println("Skipping audio file:", outputFile)
		return nil
	}

	maxRetries := viper.GetInt("tts.max_retries")
	if maxRetries < defaultMaxRetries {
		maxRetries = defaultMaxRetries
	}

	verbose := viper.GetBool("verbose_logs")
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		output, err := coquiTextToSpeech(text, outputFile)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("TTS failed — (attempt %d/%d)\n", attempt, maxRetries)

		if verbose {
			fmt.Printf("Output: %s\n", string(output))
		}
	}

	if verbose && lastErr != nil {
		fmt.Println(lastErr)
	}

	return lastErr
}

func SynthesizeTextList(paragraphs []string, language string) {
	tempDir := viper.GetString("temp_dir")

	parallelAudioCount := viper.GetInt("tts.parallel_audio_count")
	if parallelAudioCount < defaultParallelAudioCount {
		parallelAudioCount = defaultParallelAudioCount
	}

	utils.ParallelForEach(paragraphs, parallelAudioCount, func(i int, content string) {
		index := i + 1
		filename := fmt.Sprintf("part-%d.wav", index)
		outputAudioFile := filepath.Join(tempDir, filename)

		if fsutils.FileExists(outputAudioFile) {
			fmt.Println("Skipping audio file:", filename)
			return
		}

		if err := SynthesizeText(content, language, outputAudioFile); err != nil {
			fmt.Printf("Error processing part %d: %v\n", index, err)
		}
	})
}
