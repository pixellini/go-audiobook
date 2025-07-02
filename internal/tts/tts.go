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

	output, err := coquiTextToSpeechXTTS(text, language, outputFile)
	if err != nil {
		return fmt.Errorf("error generating audiobook for %s: %v, output: %s", outputFile, err, string(output))
	}

	return nil
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
