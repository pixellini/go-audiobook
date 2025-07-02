package tts

import (
	"fmt"
	"os/exec"

	"github.com/pixellini/go-audiobook/internal/formatter"
	"github.com/spf13/viper"
)

// Basically a Coqui TTS wrapper
func coquiTextToSpeechXTTS(text, language, outputFile string) ([]byte, error) {
	fmt.Println("Processing:", text)
	language = formatter.FormatToStandardLanguage(language)

	speakerWav := viper.GetString("speaker_wav")
	if speakerWav == "" {
		panic("Missing required config value: 'speaker_wav' in config.json")
	}

	cmd := exec.Command("tts",
		"--text", text,
		"--model_name", "tts_models/multilingual/multi-dataset/xtts_v2",
		"--speaker_wav", speakerWav,
		"--language_idx", language,
		"--out_path", outputFile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("error generating audiobook for %s: %v, output: %s", outputFile, err, string(output))
	}

	return output, nil
}
