package coqui

// Model represents a TTS model type supported by Coqui TTS.
// Each model has different capabilities, speed, and language support.
type Model string

// VitsConfig holds configuration specific to VITS model.
// Currently stores voice selection for VITS synthesis.
type VitsConfig struct {
	// Voice specifies the speaker voice identifier for VITS.
	Voice string
}

// XttsConfig holds configuration specific to XTTS model.
// Currently stores speaker sample path for XTTS synthesis.
type XttsConfig struct {
	// Speaker specifies the path to the speaker sample file.
	Speaker string
}

const (
	// XTTS represents the XTTS model with multilingual support.
	// Supports custom voice cloning from audio samples.
	XTTS Model = "XTTS"
	// VITS represents the VITS model with English support.
	// Uses predefined speaker voices for synthesis.
	VITS Model = "VITS"

	// modelNameXTTS is the full Coqui model identifier for XTTS.
	modelNameXTTS = "tts_models/multilingual/multi-dataset/xtts_v2"
	// modelNameVITS is the full Coqui model identifier for VITS.
	modelNameVITS = "tts_models/en/vctk/vits"
)

// String returns the string representation of the Model.
func (m Model) String() string {
	return string(m)
}

// IsValid checks if the model is supported by this package.
// Returns true for all supported model types.
func (m Model) IsValid() bool {
	switch m {
	case XTTS, VITS:
		return true
	default:
		return false
	}
}

// Name returns the full Coqui TTS model identifier.
// Used internally when calling the Coqui TTS Python process.
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

// SupportsLanguage checks if the model supports the specified language.
// Different models have different language capabilities.
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
