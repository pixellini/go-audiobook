package ttsmanager

import (
	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/logger"
	"github.com/pixellini/go-coqui"
	"github.com/pixellini/go-coqui/model"
	ttsModel "github.com/pixellini/go-coqui/models/tts"
	"github.com/pixellini/go-coqui/models/vocoder"
)

// NewTTSManager creates a new TTS manager
func New(c *config.Config, l logger.Logger) (*coqui.TTS, error) {
	// Convert string device to coqui.Device type
	device := model.Device(c.TTS.Device)

	return coqui.New(
		coqui.WithModelId(ttsModel.PresetVITSVCTK),
		coqui.WithVocoder(vocoder.PresetHifiganV2Blizzard2013),
		coqui.WithModelLanguage(model.English),
		coqui.WithMaxRetries(c.TTS.MaxRetries),
		coqui.WithDevice(device),
	)
}

// // InitializeTTS creates and initializes a TTS instance based on configuration
// func Init(config *config.Config, tempDir string, language model.Language) (*coqui.TTS, error) {
// 	var tts *coqui.TTS
// 	var err error

// 	// Convert string device to coqui.Device type
// 	device := model.Device(config.TTS.Device)

// 	if config.TTS.UseVits {
// 		// Use VITS model
// 		speakerIdx := config.TTS.VitsVoice
// 		if speakerIdx == "" {
// 			speakerIdx = defaultVitsSpeaker
// 		}
// 		tts, err = coqui.New(
// 			coqui.WithModelId(ttsModel.PresetVITSVCTK),
// 			coqui.WithVocoder(vocoder.PresetHifiganV2Blizzard2013),
// 			coqui.WithSpeakerIndex(speakerIdx),
// 			coqui.WithModelLanguage(language),
// 			coqui.WithMaxRetries(config.TTS.MaxRetries),
// 			coqui.WithDevice(device),
// 			coqui.WithOutputDir(tempDir),
// 		)
// 	} else {
// 		// Use XTTS model
// 		speakerWav := config.SpeakerWav
// 		if speakerWav == "" {
// 			// Coqui will panic anyway if speakerWav is empty
// 			panic("Speaker WAV file must be specified in the configuration")
// 		}
// 		tts, err = coqui.NewWithModelXttsV2(
// 			coqui.WithSpeakerSample(speakerWav),
// 			coqui.WithVocoder(vocoder.PresetHifiganV2Blizzard2013),
// 			coqui.WithModelLanguage(language),
// 			coqui.WithSpeaker(config.SpeakerWav),
// 			coqui.WithMaxRetries(config.TTS.MaxRetries),
// 			coqui.WithDevice(device),
// 			coqui.WithOutputDir(tempDir),
// 		)
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize TTS: %w", err)
// 	}

// 	return tts, nil
// }
