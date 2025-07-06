package tts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/spf13/viper"
)

func SynthesizeText(text, language, outputFile string) error {
	_, err := os.Stat(outputFile)
	if err == nil {
		fmt.Println("Skipping audio file:", outputFile)
		return nil
	}

	maxRetries := viper.GetInt("tts.max_retries")
	if maxRetries < 1 {
		maxRetries = 1
	}

	verbose := viper.GetBool("verbose_logs")
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := coquiTextToSpeechXTTS(text, language, outputFile)
		if err == nil {
			return nil
		}

		lastErr = fmt.Errorf("error generating audiobook for %s: %v", outputFile, err)
		fmt.Printf("TTS failed — (attempt %d/%d)\n", attempt, maxRetries)
	}

	if verbose && lastErr != nil {
		fmt.Println(lastErr)
	}

	return lastErr
}

func SynthesizeTextList(paragraphs []string, language string) {
	tempDir := viper.GetString("temp_dir")

	for i, content := range paragraphs {
		index := i + 1
		filename := fmt.Sprintf("part-%d.wav", index)
		outputAudioFile := filepath.Join(tempDir, filename)

		if fsutils.FileExists(outputAudioFile) {
			fmt.Println("Skipping audio file:", filename)
			continue
		}

		if err := SynthesizeText(content, language, outputAudioFile); err != nil {
			fmt.Printf("Error processing part %d: %v\n", index, err)
		}
	}
}
