package coqui

import "fmt"

// Default configuration values
const (
	defaultVitsVoice = "p287"
	defaultRetries   = 1
)

type Config struct {
	Language       Language
	Model          Model
	SpeakerWavFile string // Used by XTTS
	SpeakerIdx     string // Used by VITS
	MaxRetries     int
	DistDir        string
	Device         Device
}

type ConfigOption interface {
	apply(*TTS)
}

type configOptionFunc func(*TTS)

func (c configOptionFunc) apply(tts *TTS) {
	c(tts)
}

// Validate checks if the TTS configuration is valid
func (c *Config) Validate() error {
	// Use the Validator interface for cleaner validation
	if err := ValidateWithContext(c.Language, "language"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	if err := ValidateWithContext(c.Model, "model"); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	if !c.Model.SupportsLanguage(c.Language) {
		return fmt.Errorf("%w: model %s does not support language %s", ErrInvalidConfig, c.Model, c.Language)
	}

	return nil
}

func (c *Config) ToArgs() []string {
	args := []string{
		"--model_name", c.Model.Name(),
		"--device", c.Device.String(),
	}

	// Explicitly set CUDA usage based on device
	if c.Device == CUDA {
		args = append(args, "--use_cuda", "true")
	}

	if c.Model == XTTS {
		args = append(args,
			"--speaker_wav", c.SpeakerWavFile,
			"--language_idx", c.Language.String(),
		)
	}
	if c.Model == VITS {
		speakerIdx := c.SpeakerIdx
		if speakerIdx == "" {
			speakerIdx = defaultVitsVoice
		}
		args = append(args,
			"--speaker_idx", speakerIdx,
		)
	}

	return args
}

// WithLanguage sets the language for TTS
func WithLanguage(l Language) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.Language = l
	})
}

// WithModel sets the model for TTS
func WithModel(m Model) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.Model = m
	})
}

// WithSpeaker sets the speaker for TTS (automatically uses the right field based on model)
// This can be either a file path or an ID depending on the model that's chosen.
func WithSpeaker(speaker string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		if t.config.Model == XTTS {
			WithSpeakerWav(speaker).apply(t)
		}
		if t.config.Model == VITS {
			WithSpeakerIdx(speaker).apply(t)
		}
	})
}

// WithSpeakerWav sets the speaker wav file path for XTTS
func WithSpeakerWav(path string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.SpeakerWavFile = path
	})
}

// WithSpeakerIdx sets the speaker index for VITS
func WithSpeakerIdx(idx string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.SpeakerIdx = idx
	})
}

// WithMaxRetries sets the maximum number of retries that the TTS will attempt
func WithMaxRetries(mr int) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.MaxRetries = mr
	})
}

// WithDistDir sets the output path for the completed audio
func WithDistDir(outputPath string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.DistDir = outputPath
	})
}

func WithDevice(device Device) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		// If user explicitly chooses Auto, detect the best device
		if device == Auto {
			device = DetectDevice(Auto)
		}
		t.config.Device = device
	})
}
