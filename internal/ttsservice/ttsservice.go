package ttsservice

import (
	"context"

	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-coqui"
	"github.com/pixellini/go-coqui/models/tts"
)

type TTSservice interface {
	Synthesize(text, output string) ([]byte, error)
	SynthesizeContext(ctx context.Context, text, output string) ([]byte, error)
}

type CoquiTTSService struct {
	tts *coqui.TTS
}

func NewCoquiService(config *config.Config, outputDir string) (*CoquiTTSService, error) {

	tts, err := coqui.New(
		// coqui.WithModelPath(config.Model.Name),
		// coqui.WithSpeaker(config.Model.SpeakerWav),
		coqui.WithModelId(tts.PresetVITSVCTK),
		coqui.WithSpeakerIndex(config.Model.SpeakerIdx),
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

	// if config.Vocoder.Name != "" {
	// 	// Don't have a vocoder path option...
	// 	tts.Configure(
	// 	// coqui.WithVocoder(),
	// 	// coqui.WithVocoderLanguage(config.Vocoder.Language),
	// 	)
	// }

	return &CoquiTTSService{tts: tts}, nil
}

func (c *CoquiTTSService) Synthesize(text, output string) ([]byte, error) {
	return c.tts.SynthesizeContext(context.Background(), text, output)
}

func (c *CoquiTTSService) SynthesizeContext(ctx context.Context, text, output string) ([]byte, error) {
	bytes, err := c.tts.SynthesizeContext(ctx, text, output)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
