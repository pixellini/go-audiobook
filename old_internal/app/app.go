package app

import (
	"context"
	"fmt"

	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/filemanager"
	"github.com/pixellini/go-audiobook/internal/logger"
	"github.com/pixellini/go-audiobook/old_internal/audiobook"
	"github.com/pixellini/go-audiobook/old_internal/processor"
	"github.com/pixellini/go-audiobook/old_internal/ttsmanager"
	"github.com/pixellini/go-coqui/model"
)

// ApplicationService orchestrates the main application workflow
type ApplicationService interface {
	Run(ctx context.Context, resetProgress, finishAudiobook bool) error
}

// Application represents the main application
type Application struct {
	config      *config.Config
	fileManager *filemanager.FileManager
	// ttsInstance *coqui.TTS
	logger logger.Logger
}

// New creates a new Application instance
func New() (*Application, error) {
	appConfig, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	configService := config.NewViperConfigService(appConfig)
	fileManager := filemanager.NewFileManager(configService)
	// ttsInstance, err := ttsmanager.New(appConfig)
	loggerService := logger.New()

	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS manager: %w", err)
	}

	return &Application{
		config:      appConfig,
		fileManager: fileManager,
		// ttsInstance: ttsInstance,
		logger: loggerService,
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
	language, err := model.ParseLanguage(book.Language)
	if err != nil {
		return fmt.Errorf("failed to parse language: %w", err)
	}

	tm, err := ttsmanager.New(a.config, tempDir)
	if err != nil {
		return fmt.Errorf("failed to initialize TTS manager: %w", err)
	}

	fmt.Print(tm)

	// Create chapter processor
	p := processor.New(tm, a.config, tempDir, book.Title)

	// Generate chapter audio files
	if err := p.Chapters(ctx, book, audiobookInstance, finishAudiobook); err != nil {
		return fmt.Errorf("failed to generate chapter audio files: %w", err)
	}

	fmt.Println("\n--------------------------------------------------")

	// Generate final audiobook
	if err := audiobookInstance.Generate(distDir); err != nil {
		return fmt.Errorf("failed to generate final audiobook: %w", err)
	}

	// Cleanup
	if err := a.fileManager.Cleanup(tempDir); err != nil {
		fmt.Println("cleanup failed: %v", err)
	}

	return nil
}
