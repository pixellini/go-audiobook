package ttsservice

import (
	"context"
	"io"
	"log"

	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-coqui"
	"github.com/pixellini/go-coqui/models/tts"
	"github.com/pixellini/go-coqui/models/vocoder"
)

type TTSservice interface {
	Synthesize(text, output string) ([]byte, error)
	SynthesizeContext(ctx context.Context, text, output string) ([]byte, error)
}

type CoquiTTSService struct {
	tts          *coqui.TTS
	suppressLogs bool
}

func NewCoquiService(config *config.Config, outputDir string) (*CoquiTTSService, error) {
	tts, err := coqui.New(
		// coqui.WithModelPath(config.Model.Name),
		coqui.WithModelId(tts.PresetXTTSv2),
		coqui.WithSpeaker(config.Model.SpeakerWav),
		// coqui.WithModelId(tts.PresetVITSVCTK),
		// coqui.WithSpeakerIndex(config.Model.SpeakerIdx),
		coqui.WithDevice(config.Model.Device),
		coqui.WithMaxRetries(int(config.Model.MaxRetries)),
		coqui.WithOutputDir(outputDir),
	)

	if err != nil {
		return nil, err
	}

	if config.Model.SpeakerWav != "" {
		tts.Configure(
			coqui.WithSpeakerSample(config.Model.SpeakerWav),
		)
	} else {
		tts.Configure(
			coqui.WithSpeakerIndex(config.Model.SpeakerIdx),
		)
	}

	if config.Vocoder.Name != "" {
		tts.Configure(
			coqui.WithVocoder(vocoder.PresetHifiganV2Blizzard2013),
			coqui.WithVocoderLanguage(config.Vocoder.Language),
		)
	}

	return &CoquiTTSService{
		tts:          tts,
		suppressLogs: !config.VerboseLogs,
	}, nil
}

// We silence the standard library logger when verbose logs are disabled.
// We no longer redirect global Stdout/Stderr to avoid serialising concurrent calls.
// This enables true concurrency even when verbose_logs is false.
func (c *CoquiTTSService) maybeSilenceLogs() func() {
	if !c.suppressLogs {
		return func() {}
	}
	origLogOut := log.Writer()
	log.SetOutput(io.Discard)
	return func() { log.SetOutput(origLogOut) }
}

func (c *CoquiTTSService) Synthesize(text, output string) ([]byte, error) {
	return c.tts.SynthesizeContext(context.Background(), text, output)
}

func (c *CoquiTTSService) SynthesizeContext(ctx context.Context, text, output string) ([]byte, error) {
	restore := c.maybeSilenceLogs()
	defer restore()

	bytes, err := c.tts.SynthesizeContext(ctx, text, output)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
