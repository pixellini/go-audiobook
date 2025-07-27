package ttsservice

import (
	"context"
	"io"
	"log"
	"os"
	"sync"

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

var (
	devNullOnce sync.Once
	devNullFile *os.File
	suppressMu  sync.Mutex
)

// suppressOutput redirects Stdout, Stderr and the default logger to
// /dev/null (or io.Discard for the logger).
// We can call the returned restore function to put everything back.
// TODO: It might be better to move this to go-coqui as an option instead of doing it here.
func (c *CoquiTTSService) suppressOutput() (restore func(), err error) {
	if !c.suppressLogs {
		return func() {}, nil
	}

	var openErr error
	// only open /dev/null once
	devNullOnce.Do(func() {
		devNullFile, openErr = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	})
	if openErr != nil {
		return nil, openErr
	}

	// allow only one active suppression at a time
	suppressMu.Lock()

	origStdout := os.Stdout
	origStderr := os.Stderr
	origLogOut := log.Writer()

	os.Stdout = devNullFile
	os.Stderr = devNullFile
	log.SetOutput(io.Discard)

	return func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
		log.SetOutput(origLogOut)
		suppressMu.Unlock()
	}, nil
}

func (c *CoquiTTSService) Synthesize(text, output string) ([]byte, error) {
	restore, err := c.suppressOutput()
	if err != nil {
		return nil, err
	}
	defer restore()

	return c.tts.SynthesizeContext(context.Background(), text, output)
}

func (c *CoquiTTSService) SynthesizeContext(ctx context.Context, text, output string) ([]byte, error) {
	restore, err := c.suppressOutput()
	if err != nil {
		return nil, err
	}
	defer restore()

	bytes, err := c.tts.SynthesizeContext(ctx, text, output)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
