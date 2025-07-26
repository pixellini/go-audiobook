package interfaces

import (
	"github.com/pixellini/go-audiobook/internal/epub"
)

// TTSService defines the interface for text-to-speech operations
type TTSService interface {
	Synthesize(text, outputFile string) ([]byte, error)
}

// Logger defines the interface for logging operations

// FileService defines the interface for file operations
type FileService interface {
	GetImageFile() string
	GetEpubFile() (*epub.Epub, error)
	SetupDirs() (tempDir, distDir string, err error)
	ResetProgress(tempDir, distDir string) error
	Cleanup(tempDir string) error
}
