package coqui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type TTS struct {
	config *Config
}

// Synthesizer defines the interface for TTS synthesis
type Synthesizer interface {
	Synthesize(text string) ([]byte, error)
	SynthesizeContext(ctx context.Context, text string) ([]byte, error)
}

// ErrInvalidConfig is returned when configuration is invalid
var ErrInvalidConfig = errors.New("invalid configuration")

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

// NewWithXtts creates a new TTS instance configured for XTTS model
func NewWithXtts(speakerWav string, options ...ConfigOption) (*TTS, error) {
	opts := append([]ConfigOption{
		WithModel(XTTS),
		WithSpeakerWav(speakerWav),
	}, options...)
	return New(opts...)
}

// NewWithVits creates a new TTS instance configured for VITS model
// If speakerIdx is empty, the default speaker will be used
func NewWithVits(speakerIdx string, options ...ConfigOption) (*TTS, error) {
	opts := []ConfigOption{
		WithModel(VITS),
		WithLanguage(English), // VITS only supports English
	}

	if speakerIdx != "" {
		opts = append(opts, WithSpeakerIdx(speakerIdx))
	}

	opts = append(opts, options...)
	return New(opts...)
}

// Synthesize converts text to speech
func (t *TTS) Synthesize(text, output string) ([]byte, error) {
	return t.SynthesizeContext(context.Background(), text, output)
}

// SynthesizeContext converts text to speech with context support
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

// Config returns a copy of the current configuration
func (t *TTS) Config() Config {
	return *t.config
}

func (t *TTS) Configure(options ...ConfigOption) {
	for _, option := range options {
		option.apply(t)
	}
}

// exec executes the TTS command with the given text and output path
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
