package coqui

// ConfigOption defines an interface for TTS configuration options.
type ConfigOption interface {
	apply(*TTS)
}

type configOptionFunc func(*TTS)

// apply will set the configuration option on the TTS instance.
func (c configOptionFunc) apply(tts *TTS) {
	c(tts)
}

// WithLanguage sets the target language for TTS synthesis.
// Note: Language support varies by model.
func WithLanguage(l Language) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.Language = l
	})
}

// WithModel sets the TTS model to use.
func WithModel(m Model) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.Model = m
	})
}

// WithSpeaker sets the speaker for TTS synthesis.
// Automatically selects the appropriate configuration based on the model type.
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

// WithSpeakerWav sets the speaker sample file path for XTTS.
func WithSpeakerWav(path string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.SpeakerWavFile = path
	})
}

// WithSpeakerIdx sets the speaker index identifier for VITS.
func WithSpeakerIdx(idx string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.SpeakerIdx = idx
	})
}

// WithMaxRetries sets the maximum number of synthesis attempts on failure.
func WithMaxRetries(mr int) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.MaxRetries = mr
	})
}

// WithDistDir sets the output directory for generated audio files.
func WithDistDir(outputPath string) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		t.config.DistDir = outputPath
	})
}

// WithDevice sets the compute device for TTS synthesis.
// If Auto is specified, the best available device will be detected automatically.
func WithDevice(device Device) ConfigOption {
	return configOptionFunc(func(t *TTS) {
		// If user explicitly chooses Auto, detect the best device
		if device == Auto {
			device = DetectDevice(Auto)
		}
		t.config.Device = device
	})
}
