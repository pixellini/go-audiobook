package config

import (
	"github.com/spf13/viper"
)

// ViperConfigService implements ConfigService using Viper
type ViperConfigService struct {
	config *Config
}

// NewViperConfigService creates a new config service
func NewViperConfigService(config *Config) *ViperConfigService {
	return &ViperConfigService{
		config: config,
	}
}

// GetEpubPath returns the EPUB file path
func (c *ViperConfigService) GetEpubPath() string {
	return viper.GetString("epub_path")
}

// GetImagePath returns the image file path
func (c *ViperConfigService) GetImagePath() string {
	return viper.GetString("image_path")
}

// GetDistDir returns the distribution directory
func (c *ViperConfigService) GetDistDir() string {
	return viper.GetString("dist_dir")
}

// GetSpeakerWav returns the speaker WAV file path
func (c *ViperConfigService) GetSpeakerWav() string {
	return c.config.SpeakerWav
}

// IsTestMode returns whether test mode is enabled
func (c *ViperConfigService) IsTestMode() bool {
	return viper.GetBool("test_mode")
}

// GetTTSConfig returns the TTS configuration
func (c *ViperConfigService) GetTTSConfig() TTSConfig {
	return c.config.TTS
}

// GetConcurrency returns the TTS concurrency setting
func (c *ViperConfigService) GetConcurrency() int {
	return viper.GetInt("tts.concurrency")
}
