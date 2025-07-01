package tts

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pixellini/go-audiobook/internal/formatter"
)

/*
This is a temporary solution for text-to-speech using Coqui TTS via command line execution.
Proper integration of Coqui TTS into Go is pending and will replace this approach in the future.
*/
func coquiTextToSpeechXTTS(text, language, outputFile string) ([]byte, error) {
	fmt.Println("Processing:", text)
	language = formatter.FormatToStandardLanguage(language)

	speaker := os.Getenv("SPEAKER")
	if speaker == "" {
		panic("SPEAKER is not set")
	}

	cmd := exec.Command("tts",
		"--text", text,
		"--model_name", "tts_models/multilingual/multi-dataset/xtts_v2",
		"--speaker_wav", "./speakers/"+speaker,
		"--language_idx", language,
		"--out_path", outputFile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("error generating audiobook for %s: %v, output: %s", outputFile, err, string(output))
	}

	return output, nil
}
