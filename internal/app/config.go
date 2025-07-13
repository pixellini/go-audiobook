package app

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	EpubPath     string    `mapstructure:"epub_path"`
	ImagePath    string    `mapstructure:"image_path"`
	SpeakerWav   string    `mapstructure:"speaker_wav"`
	DistDir      string    `mapstructure:"dist_dir"`
	OutputFormat string    `mapstructure:"output_format"`
	VerboseLogs  bool      `mapstructure:"verbose_logs"`
	TestMode     bool      `mapstructure:"test_mode"`
	TTS          TTSConfig `mapstructure:"tts"`
}

// TTSConfig represents TTS-specific configuration
type TTSConfig struct {
	MaxRetries         int    `mapstructure:"max_retries"`
	ParallelAudioCount int    `mapstructure:"parallel_audio_count"`
	UseVits            bool   `mapstructure:"use_vits"`
	VitsVoice          string `mapstructure:"vits_voice"`
	Device             string `mapstructure:"device"`
}

// LoadConfig loads the configuration from file and environment
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Set default values if not present in config
	viper.SetDefault("image_path", "")
	viper.SetDefault("temp_dir", "./.temp")
	viper.SetDefault("output_format", "m4b")
	viper.SetDefault("verbose_logs", false)
	viper.SetDefault("test_mode", false)
	viper.SetDefault("tts.max_retries", 3)
	viper.SetDefault("tts.parallel_audio_count", 4)
	viper.SetDefault("tts.use_vits", false)
	viper.SetDefault("tts.vits_voice", "p287")
	viper.SetDefault("tts.device", "auto")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
