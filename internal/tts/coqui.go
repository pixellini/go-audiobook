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

type CoquiTTSConfig struct {
	ModelName  TTSModel
	Language   string
	SpeakerWav string
	Voice      string
	args       []string
}

// This is a singleton instance of CoquiTTSConfig.
// The config won't change, so we only need to initialise the configuration once.
var BaseCoquiTTSConfig = &CoquiTTSConfig{}

func (c *CoquiTTSConfig) Init(language string) {
	language = formatter.FormatToStandardLanguage(language)

	c.ModelName = ModelXTTS
	c.Language = language

	if viper.GetBool("tts.use_vits") {
		// VITS only supports English. Panic if language is not English.
		if language != "en" {
			panic("The VITS model currently only supports English, please make sure your EPUB is in English.")
		}

		c.ModelName = ModelVITS
		c.setVoice()

	} else {
		c.setSpeaker()
	}

	fmt.Printf("Using %s model for TTS.\n", c.ModelName)

	c.buildBaseCommandArgs()
}

// These are the default arguments for the TTS command.
// They won't change based on the model.
// Therefore, the "--text" and "--out_path" arguments are not needed here,
func (c *CoquiTTSConfig) buildBaseCommandArgs() {
	// Reset args to avoid duplicate arguments on repeated calls
	c.args = []string{}
	if c.ModelName == ModelXTTS {
		c.args = append(c.args,
			"--model_name", modelNameXTTS,
			"--speaker_wav", c.SpeakerWav,
			"--language_idx", c.Language,
		)
	} else {
		c.args = append(c.args,
			"--model_name", modelNameVITS,
			"--speaker_idx", c.Voice,
		)
	}

	if device.Manager.Device == device.DeviceCUDA {
		c.args = append(c.args, "--use_cuda", "true")
	}
}

func (c *CoquiTTSConfig) setSpeaker() {
	speakerWav := viper.GetString("speaker_wav")
	if speakerWav == "" {
		panic("Missing required config value: 'speaker_wav'")
	}
	c.SpeakerWav = speakerWav
}

func (c *CoquiTTSConfig) setVoice() {
	voice := viper.GetString("tts.vits_voice")
	if voice == "" {
		c.Voice = defaultVitsVoice
	} else {
		c.Voice = voice
	}
}

// Note: The "language" parameter will only used by XTTS (not VITS).
// For VITS, only English is supported, therefore the EPUB language must be in English.
func coquiTextToSpeech(text, outputFile string) ([]byte, error) {
	// Sudden idea: We could use this empty check to implement a speaker pause.
	if text == "" {
		return nil, fmt.Errorf("paragraph is empty")
	}

	args := []string{
		"--text", text,
		"--out_path", outputFile,
	}
	args = append(args, BaseCoquiTTSConfig.args...)

	cmd := exec.Command("tts", args...)

	fmt.Println("Processing:", text)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, fmt.Errorf("error generating audiobook (%s) for %s: %v, output: %s", BaseCoquiTTSConfig.ModelName, outputFile, err, string(output))
	}

	return output, nil
}
