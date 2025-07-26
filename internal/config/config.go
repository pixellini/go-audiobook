package config

import (
	"fmt"

	"github.com/pixellini/go-coqui/model"
	"github.com/pixellini/go-coqui/models/vocoder"
	"github.com/spf13/viper"
)

type Config struct {
	VerboseLogs bool    `mapstructure:"verbose_logs"`
	TestMode    bool    `mapstructure:"test_mode"`
	Epub        Epub    `mapstructure:"epub"`
	Output      Output  `mapstructure:"output"`
	Model       Model   `mapstructure:"model"`
	Vocoder     Vocoder `mapstructure:"vocoder"`
}

type Epub struct {
	Path        string `mapstructure:"path"`
	CoverImage  string `mapstructure:"cover_image"`
	Title       string `mapstructure:"title"`
	Author      string `mapstructure:"author"`
	Language    string `mapstructure:"language"`
	Publisher   string `mapstructure:"publisher"`
	Description string `mapstructure:"description"`
}

type Output struct {
	Path string `mapstructure:"path"`
	// Format   string `mapstructure:"format"`
	Filename string `mapstructure:"filename"`
}

type Model struct {
	Name        string         `mapstructure:"name"`
	Language    model.Language `mapstructure:"language"`
	SpeakerWav  string         `mapstructure:"speaker_wav"`
	SpeakerIdx  string         `mapstructure:"speaker_idx"`
	Concurrency uint8          `mapstructure:"concurrency"`
	MaxRetries  uint8          `mapstructure:"max_retries"`
	Device      model.Device   `mapstructure:"device"`
}

type Vocoder struct {
	Name     string         `mapstructure:"name"`
	Language model.Language `mapstructure:"language"`
}

const defaultConcurrency = 4

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	setDefaults()
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// check if the end of output doesn't contain a slash, otherwise add it
	if config.Output.Path != "" && config.Output.Path[len(config.Output.Path)-1] != '/' {
		config.Output.Path += "/"
	}

	// we can only allow positive numbers for concurrency
	if config.Model.Concurrency < 1 {
		config.Model.Concurrency = defaultConcurrency
	}

	return config, nil
}

func setDefaults() {
	// Set default values if not present in config
	viper.SetDefault("output.path", "./.dist/")
	viper.SetDefault("output.format", "m4b")

	// Model Defaults
	viper.SetDefault("model.name", "tts_models/multilingual/multi-dataset/xtts_v2")
	viper.SetDefault("model.speaker_idx", "p286")
	viper.SetDefault("model.language", model.English)
	viper.SetDefault("model.concurrency", defaultConcurrency)
	viper.SetDefault("model.max_retries", 5)
	viper.SetDefault("model.device", model.DeviceCPU)

	viper.SetDefault("vocoder.name", vocoder.PresetHifiganV2Blizzard2013.Name())
	viper.SetDefault("vocoder.language", model.English)
}

func (o Output) OutputFileName() string {
	return o.Filename + ".m4b"
}

func (o Output) FullPath() string {
	return o.Path + o.OutputFileName()
}
