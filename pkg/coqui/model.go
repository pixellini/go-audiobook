package coqui

type Model string

type VitsConfig struct {
	Voice string
}

type XttsConfig struct {
	Speaker string
}

// Available Models for Coqui TTS
const (
	XTTS Model = "XTTS"
	VITS Model = "VITS"

	modelNameXTTS = "tts_models/multilingual/multi-dataset/xtts_v2"
	modelNameVITS = "tts_models/en/vctk/vits"
)

func (m Model) String() string {
	return string(m)
}

// IsValid checks if the model is supported
func (m Model) IsValid() bool {
	switch m {
	case XTTS, VITS:
		return true
	default:
		return false
	}
}

// Name returns the full model name for the TTS engine
func (m Model) Name() string {
	switch m {
	case XTTS:
		return modelNameXTTS
	case VITS:
		return modelNameVITS
	default:
		return modelNameXTTS
	}
}

// SupportsLanguage checks if the model supports the given language
func (m Model) SupportsLanguage(lang Language) bool {
	switch m {
	case VITS:
		// VITS only supports English
		return lang == English
	case XTTS:
		// XTTS supports all valid languages from Coqui
		return lang.IsValid()
	default:
		return false
	}
}
