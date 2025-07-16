package app

import (
	"context"
	"fmt"

	"github.com/pixellini/go-audiobook/internal/audiobook"
	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/filemanager"
	"github.com/pixellini/go-audiobook/internal/interfaces"
	"github.com/pixellini/go-audiobook/internal/logger"
	"github.com/pixellini/go-audiobook/internal/processor"
	"github.com/pixellini/go-audiobook/internal/tts"
	"github.com/pixellini/go-coqui"
)

// Application represents the main application
type Application struct {
	config      *config.Config
	fileManager *filemanager.FileManager
	logger      interfaces.Logger
}

// New creates a new Application instance
func New() (*Application, error) {
	appConfig, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	logger := logger.NewSimpleLogger()
	configService := config.NewViperConfigService(appConfig)
	fileManager := filemanager.NewFileManager(configService, logger)

	return &Application{
		config:      appConfig,
		fileManager: fileManager,
		logger:      logger,
	}, nil
}

// Run executes the main application logic
func (a *Application) Run(resetProgress, finishAudiobook bool) error {
	return a.RunWithContext(context.Background(), resetProgress, finishAudiobook)
}

// RunWithContext executes the main application logic with context support
func (a *Application) RunWithContext(ctx context.Context, resetProgress, finishAudiobook bool) error {
	// Get image and book files
	image := a.fileManager.GetImageFile()
	book, err := a.fileManager.GetEpubFile()
	if err != nil {
		return fmt.Errorf("failed to load EPUB file: %w", err)
	}

	// Create audiobook instance
	audiobookInstance := audiobook.NewWithEPUB(book, image)

	// Setup directories
	tempDir, distDir, err := a.fileManager.SetupDirs()
	if err != nil {
		return fmt.Errorf("failed to setup directories: %w", err)
	}

	// Reset progress if requested
	if resetProgress {
		if err := a.fileManager.ResetProgress(tempDir, distDir); err != nil {
			return fmt.Errorf("failed to reset progress: %w", err)
		}
	}

	// Parse the epub language and initialize TTS
	language := coqui.ParseLanguage(book.Language)
	coquiTTS, err := tts.Init(a.config, language)
	if err != nil {
		return err
	}

	ttsService := tts.NewManager(coquiTTS)

	// Create chapter processor
	p := processor.NewChapterProcessor(ttsService, a.config, a.logger, tempDir, book.Title)

	// Generate chapter audio files
	if err := p.ProcessChapters(ctx, book, audiobookInstance, finishAudiobook); err != nil {
		return fmt.Errorf("failed to generate chapter audio files: %w", err)
	}

	a.logger.Info("\n\n--------------------------------------------------")

	// Generate final audiobook
	if err := audiobookInstance.Generate(distDir); err != nil {
		return fmt.Errorf("failed to generate final audiobook: %w", err)
	}

	// Cleanup
	if err := a.fileManager.Cleanup(tempDir); err != nil {
		a.logger.Warnf("cleanup failed: %v", err)
	}

	return nil
}
