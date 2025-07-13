package filemanager

import (
	"fmt"

	"github.com/pixellini/go-audiobook/internal/epub"
	"github.com/pixellini/go-audiobook/internal/fs"
	"github.com/pixellini/go-audiobook/internal/interfaces"
)

const (
	testBookPath  = "./examples/test/book.epub"
	testImagePath = "./examples/test/cover.png"
)

// FileManager implements the FileService interface
type FileManager struct {
	config interfaces.ConfigService
	logger interfaces.Logger
}

// NewFileManager creates a new file manager
func NewFileManager(config interfaces.ConfigService, logger interfaces.Logger) *FileManager {
	return &FileManager{
		config: config,
		logger: logger,
	}
}

// GetImageFile returns the path to the cover image
func (fm *FileManager) GetImageFile() string {
	if fm.config.IsTestMode() {
		fm.logger.Info("TEST MODE ENABLED: Using mock cover image.")
		return testImagePath
	}

	image := fm.config.GetImagePath()
	if image == "" {
		fm.logger.Warn("No image path provided in config.")
	}
	return image
}

// GetEpubFile loads and returns the EPUB file
func (fm *FileManager) GetEpubFile() (*epub.Epub, error) {
	epubPath := fm.config.GetEpubPath()

	if fm.config.IsTestMode() {
		fm.logger.Info("TEST MODE ENABLED: Using mock epub book.")
		epubPath = testBookPath
	}

	if epubPath == "" {
		return nil, fmt.Errorf("missing required config value: 'epub_path' in config.json")
	}

	book, err := epub.New(epubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load EPUB file: %w", err)
	}
	return book, nil
}

// SetupDirs creates and returns the temp and dist directories
func (fm *FileManager) SetupDirs() (string, string, error) {
	tempDir, err := fs.GetOrCreateTempDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	distDir := fm.config.GetDistDir()
	if err := fs.CreateDirIfNotExist(distDir); err != nil {
		return "", "", fmt.Errorf("failed to create dist directory: %w", err)
	}

	return tempDir, distDir, nil
}

// ResetProgress resets the audiobook generation process
func (fm *FileManager) ResetProgress(tempDir, distDir string) error {
	fm.logger.Info("Resetting audiobook generation process...")

	if err := fs.EmptyDir(tempDir); err != nil {
		return fmt.Errorf("failed to empty temp directory: %w", err)
	}

	if err := fs.EmptyDir(distDir); err != nil {
		return fmt.Errorf("failed to empty dist directory: %w", err)
	}

	return nil
}

// Cleanup performs the removal of temporary directories and files
func (fm *FileManager) Cleanup(tempDir string) error {
	// Clean and remove chapters directory
	chaptersDir := fmt.Sprintf("%s/chapters", tempDir)
	if err := fs.RemoveAllFilesInDir(chaptersDir); err == nil {
		_ = fs.RemoveDirIfEmpty(chaptersDir)
	}

	// Clean and remove temp dir
	if err := fs.RemoveDirIfEmpty(tempDir); err != nil {
		return fmt.Errorf("unable to remove temp directory: %w", err)
	}

	return nil
}
