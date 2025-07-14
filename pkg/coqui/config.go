package coqui

import "fmt"

// Default configuration values
const (
	defaultVitsVoice = "p287"
	defaultRetries   = 1
)

// Config holds the configuration for TTS synthesis.
type Config struct {
	// Language specifies the target language for synthesis (XTTS only).
	// See language.go for all supported language codes.
	Language Language `json:"language" yaml:"language"`
	// Model specifies which TTS model to use.
	Model Model `json:"model" yaml:"model"`
	// SpeakerWavFile is the path to the speaker sample file (XTTS only).
	// Should be a clear audio sample of the desired voice (1-3 minutes recommended).
	SpeakerWavFile string `json:"speakerWavFile" yaml:"speakerWavFile"`
	// SpeakerIdx is the speaker index identifier (VITS only).
	// Use speaker IDs like "p225", "p287", etc. from the VCTK dataset.
	SpeakerIdx string `json:"speakerIdx" yaml:"speakerIdx"`
	// MaxRetries is the maximum number of synthesis attempts on failure.
	// Recommended range is 1-5; higher values increase reliability but slow down failure recovery.
	MaxRetries int `json:"maxRetries" yaml:"maxRetries"`
	// DistDir is the output directory for generated audio files.
	// If empty, files are saved to the current working directory.
	DistDir string `json:"distDir" yaml:"distDir"`
	// Device specifies the compute device (auto/cpu/cuda/mps).
	// Use "auto" for automatic detection, "cuda" for GPU acceleration if available.
	Device Device `json:"device" yaml:"device"`
}

// Validate checks if the TTS configuration is valid and returns an error
// if any configuration values are invalid or incompatible.
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

// ToArgs converts the TTS configuration to command-line arguments
// for the underlying Coqui TTS Python process.
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
