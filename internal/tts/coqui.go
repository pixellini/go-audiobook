package tts

import (
	"fmt"
	"os/exec"

	"github.com/pixellini/go-audiobook/internal/formatter"
	"github.com/spf13/viper"
)

const defaultVitsVoice = "p287"
const xttsModelName = "tts_models/multilingual/multi-dataset/xtts_v2"
const vitsModelName = "tts_models/en/vctk/vits"

// Coqui TTS wrapper for XTTS
func coquiTextToSpeechXTTS(text, language, outputFile string) ([]byte, error) {
	fmt.Println("Processing (XTTS):", text)
	language = formatter.FormatToStandardLanguage(language)

	speakerWav := viper.GetString("speaker_wav")
	if speakerWav == "" {
		panic("Missing required config value: 'speaker_wav' in config.json")
	}

	cmd := exec.Command("tts",
		"--text", text,
		"--model_name", xttsModelName,
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

// Coqui TTS wrapper for VITS
func coquiTextToSpeechVITS(text, outputFile string, voice string) ([]byte, error) {
	fmt.Println("Processing (VITS):", text)

	cmd := exec.Command("tts",
		"--text", text,
		"--model_name", vitsModelName,
		"--speaker_idx", voice,
		"--out_path", outputFile,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error generating audiobook (VITS) for %s: %v, output: %s", outputFile, err, string(output))
	}

	return output, nil
}

// Unified Coqui TTS function
func coquiTextToSpeech(text, language, outputFile string) ([]byte, error) {
	useVitsModel := viper.GetBool("tts.use_vits")

	if useVitsModel {
		voice := viper.GetString("tts.vits_voice")

		if voice == "" {
			voice = defaultVitsVoice
		}

		return coquiTextToSpeechVITS(text, outputFile, voice)
	}

	// Default to XTTS if VITS is not used
	return coquiTextToSpeechXTTS(text, language, outputFile)
}
