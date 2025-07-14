package coqui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// TTS represents a text-to-speech synthesis engine.
// Configured with a specific model, language, and device settings.
type TTS struct {
	config *Config
}

// Synthesizer defines the interface for TTS synthesis operations.
// Implementations should support both blocking and context-aware synthesis.
type Synthesizer interface {
	Synthesize(text string) ([]byte, error)
	SynthesizeContext(ctx context.Context, text string) ([]byte, error)
}

// ErrInvalidConfig is returned when TTS configuration validation fails.
var ErrInvalidConfig = errors.New("invalid configuration")

// New creates a new TTS instance with the specified configuration options.
func New(options ...ConfigOption) (*TTS, error) {
	// Build the config, apply the defaults
	tts := &TTS{
		config: &Config{
			Language:   English,
			Model:      XTTS,
			MaxRetries: defaultRetries,
		},
	}

	for _, option := range options {
		option.apply(tts)
	}

	// Validate configuration
	if err := tts.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid TTS configuration: %w", err)
	}

	return tts, nil
}

// NewWithXtts creates a new TTS instance configured for the XTTS model.
// Requires a speaker sample file path for voice cloning.
func NewWithXtts(speakerWav string, options ...ConfigOption) (*TTS, error) {
	opts := append([]ConfigOption{
		WithModel(XTTS),
		WithSpeakerWav(speakerWav),
	}, options...)
	return New(opts...)
}

// NewWithVits creates a new TTS instance configured for the VITS model.
// If speakerIdx is empty, the default speaker (p287) will be used.
func NewWithVits(speakerIdx string, options ...ConfigOption) (*TTS, error) {
	opts := []ConfigOption{
		WithModel(VITS),
		WithLanguage(English), // VITS requires English
	}

	if speakerIdx != "" {
		opts = append(opts, WithSpeakerIdx(speakerIdx))
	}

	opts = append(opts, options...)
	return New(opts...)
}

// NewFromConfig creates a new TTS instance from a Config struct.
// Allows loading configuration from JSON/YAML files with optional overrides.
// Additional ConfigOption parameters can override config file settings.
func NewFromConfig(config *Config, options ...ConfigOption) (*TTS, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	opts := []ConfigOption{
		WithModel(config.Model),
		WithLanguage(config.Language),
		WithDevice(config.Device),
		WithMaxRetries(config.MaxRetries),
	}

	if config.SpeakerWavFile != "" {
		opts = append(opts, WithSpeakerWav(config.SpeakerWavFile))
	}
	if config.SpeakerIdx != "" {
		opts = append(opts, WithSpeakerIdx(config.SpeakerIdx))
	}
	if config.DistDir != "" {
		opts = append(opts, WithDistDir(config.DistDir))
	}

	// Allow additional options to override config file settings
	opts = append(opts, options...)

	return New(opts...)
}

// Synthesize converts text to speech and saves it to the specified output file.
// This is a convenience method that uses context.Background().
func (t *TTS) Synthesize(text, output string) ([]byte, error) {
	return t.SynthesizeContext(context.Background(), text, output)
}

// SynthesizeContext converts text to speech with context support for cancellation.
// Supports automatic retries on failure and returns the command output on success.
// Returns an error if the output file already exists.
func (t *TTS) SynthesizeContext(ctx context.Context, text, output string) ([]byte, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	_, err := os.Stat(output)
	if err == nil {
		return nil, fmt.Errorf("audio file already created")
	}

	var lastErr error
	for attempt := 1; attempt <= t.config.MaxRetries; attempt++ {
		cmdOutput, err := t.exec(ctx, text, output)
		if err == nil {
			return cmdOutput, nil
		}

		lastErr = err
		log.Print(err)
		log.Printf("TTS failed — (attempt %d/%d)\n", attempt, t.config.MaxRetries)
	}

	return nil, lastErr
}

// Config returns a copy of the current TTS configuration.
// The returned Config can be safely modified without affecting the TTS instance.
func (t *TTS) Config() Config {
	return *t.config
}

// Configure applies additional configuration options to the TTS instance.
// Use this to modify settings after the TTS instance has been created.
func (t *TTS) Configure(options ...ConfigOption) {
	for _, option := range options {
		option.apply(t)
	}
}

// exec executes the Coqui TTS command with the specified text and output path.
// This is an internal method that handles the actual subprocess execution.
func (t *TTS) exec(ctx context.Context, text, output string) ([]byte, error) {
	args := t.config.ToArgs()
	args = append(args,
		"--text", text,
		"--out_path", output,
	)

	cmd := exec.CommandContext(ctx, "tts", args...)

	fmt.Printf("\nProcessing text: %q", text)

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return cmdOutput, fmt.Errorf("TTS command failed: %w", err)
	}

	return cmdOutput, nil
}
