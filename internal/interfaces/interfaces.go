package interfaces

import (
	"context"

	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/epub"
)

// TTSService defines the interface for text-to-speech operations
type TTSService interface {
	Synthesize(text, outputFile string) ([]byte, error)
}

// Logger defines the interface for logging operations
type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
}

// FileService defines the interface for file operations
type FileService interface {
	GetImageFile() string
	GetEpubFile() (*epub.Epub, error)
	SetupDirs() (tempDir, distDir string, err error)
	ResetProgress(tempDir, distDir string) error
	Cleanup(tempDir string) error
}

// ConfigService defines the interface for configuration access
type ConfigService interface {
	GetEpubPath() string
	GetImagePath() string
	GetDistDir() string
	GetSpeakerWav() string
	IsTestMode() bool
	GetTTSConfig() config.TTSConfig
	GetConcurrency() int
}

// ApplicationService orchestrates the main application workflow
type ApplicationService interface {
	Run(ctx context.Context, resetProgress, finishAudiobook bool) error
}
