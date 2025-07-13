package tts

import (
	"fmt"

	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/pkg/coqui"
)

const (
	defaultVitsSpeaker = "p287"
)

// TTSManager implements the TTSService interface
type TTSManager struct {
	tts *coqui.TTS
}

// NewTTSManager creates a new TTS manager
func NewManager(tts *coqui.TTS) *TTSManager {
	return &TTSManager{
		tts: tts,
	}
}

// Synthesize synthesizes text to audio file
func (tm *TTSManager) Synthesize(text, outputFile string) ([]byte, error) {
	return tm.tts.Synthesize(text, outputFile)
}

// InitializeTTS creates and initializes a TTS instance based on configuration
func Init(config *config.Config, language coqui.Language) (*coqui.TTS, error) {
	var tts *coqui.TTS
	var err error

	if config.TTS.UseVits {
		// Use VITS model
		speakerIdx := config.TTS.VitsVoice
		if speakerIdx == "" {
			speakerIdx = defaultVitsSpeaker
		}
		tts, err = coqui.NewWithVits(speakerIdx,
			coqui.WithLanguage(language),
			coqui.WithMaxRetries(config.TTS.MaxRetries),
		)
	} else {
		// Use XTTS model
		speakerWav := config.SpeakerWav
		if speakerWav == "" {
			panic("Speaker WAV file must be specified in the configuration")
		}
		tts, err = coqui.NewWithXtts(speakerWav,
			coqui.WithLanguage(language),
			coqui.WithMaxRetries(config.TTS.MaxRetries),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS: %w", err)
	}

	return tts, nil
}
