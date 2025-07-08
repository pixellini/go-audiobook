package tts

import (
	"fmt"
	"os/exec"

	"github.com/pixellini/go-audiobook/internal/device"
	"github.com/pixellini/go-audiobook/internal/formatter"
	"github.com/spf13/viper"
)

type TTSModel string

// Available TTS models for go-audiobook
const (
	ModelXTTS TTSModel = "XTTS"
	ModelVITS TTSModel = "VITS"
)

const (
	modelNameXTTS = "tts_models/multilingual/multi-dataset/xtts_v2"
	modelNameVITS = "tts_models/en/vctk/vits"
)

const defaultVitsVoice = "p287"

func coquiTextToSpeechXTTS(text, language, outputFile string) ([]byte, error) {
	fmt.Println("Processing:", text)

	speakerWav := viper.GetString("speaker_wav")
	if !viper.IsSet(speakerWav) {
		panic("Missing required config value: 'speaker_wav' in config.json")
	}

	args := []string{
		"--text", text,
		"--model_name", modelNameXTTS,
		"--speaker_wav", speakerWav,
		"--language_idx", language,
		"--out_path", outputFile,
		"--device", string(device.Manager.Device),
	}
	if device.Manager.Device == device.DeviceCUDA {
		args = append(args, "--use_cuda", "true")
	}

	cmd := exec.Command("tts", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error generating audiobook for %s: %v, output: %s", outputFile, err, string(output))
	}

	return output, nil
}

func coquiTextToSpeechVITS(text, outputFile, voice string) ([]byte, error) {
	fmt.Println("Processing:", text)

	args := []string{
		"--text", text,
		"--model_name", modelNameVITS,
		"--speaker_idx", voice,
		"--out_path", outputFile,
		"--device", string(device.Manager.Device),
	}
	if device.Manager.Device == device.DeviceCUDA {
		args = append(args, "--use_cuda", "true")
	}

	cmd := exec.Command("tts", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error generating audiobook (VITS) for %s: %v, output: %s", outputFile, err, string(output))
	}

	return output, nil
}

// Note: The "language" parameter will only used by XTTS (not VITS).
// For VITS, only English is supported, therefore the EPUB language must be in English.
func coquiTextToSpeech(text, language, outputFile string) ([]byte, error) {
	language = formatter.FormatToStandardLanguage(language)
	useVitsModel := viper.GetBool("tts.use_vits")

	if useVitsModel {
		// VITS only supports English. Panic if language is not English.
		if language != "en" {
			panic("The VITS model currently only supports English, please make sure your EPUB is in English.")
		}

		voice := viper.GetString("tts.vits_voice")
		if !viper.IsSet(voice) {
			voice = defaultVitsVoice
		}

		return coquiTextToSpeechVITS(text, outputFile, voice)
	}

	// Default to XTTS if VITS is not used
	return coquiTextToSpeechXTTS(text, language, outputFile)
}
