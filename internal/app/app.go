package app

import (
	"fmt"

	"github.com/pixellini/go-audiobook/internal/audiobook"
	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/utils"
	"github.com/pixellini/go-audiobook/pkg/coqui"
	"github.com/pixellini/go-audiobook/pkg/epub"
	"github.com/spf13/viper"
)

const processingFileType = "wav"

const (
	testBookPath  = "./examples/test/book.epub"
	testImagePath = "./examples/test/cover.png"
)

// Application represents the main application
type Application struct {
	config *Config
	tts    *coqui.TTS
}

// New creates a new Application instance
func New() (*Application, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &Application{
		config: config,
		tts:    nil, // Will be initialized in Run() when we know the language
	}, nil
}

// Run executes the main application logic
func (a *Application) Run(resetProgress, finishAudiobook bool) error {
	// Get image and book files
	image := a.getImageFile()
	book, err := a.getEpubFile()
	if err != nil {
		return fmt.Errorf("failed to load EPUB file: %w", err)
	}

	// Create audiobook instance
	audiobookInstance := audiobook.NewWithEPUB(book, image)

	// Setup directories
	tempDir, distDir, err := a.setupDirs()
	if err != nil {
		return fmt.Errorf("failed to setup directories: %w", err)
	}

	// Reset progress if requested
	if resetProgress {
		if err := a.resetAudiobookGeneration(tempDir, distDir); err != nil {
			return fmt.Errorf("failed to reset progress: %w", err)
		}
	}

	// Parse the epub language and initialize TTS
	language := coqui.ParseLanguage(book.Language)

	// Configure TTS based on app config
	var tts *coqui.TTS

	if a.config.TTS.UseVits {
		// Use VITS model
		speakerIdx := a.config.TTS.VitsVoice
		if speakerIdx == "" {
			speakerIdx = "" // Will use default
		}
		tts, err = coqui.NewWithVits(speakerIdx,
			coqui.WithLanguage(language),
			coqui.WithMaxRetries(a.config.TTS.MaxRetries),
		)
	} else {
		// Use XTTS model
		speakerWav := a.config.SpeakerWav
		if speakerWav == "" {
			speakerWav = "speakers/speaker.wav" // default fallback
		}
		tts, err = coqui.NewWithXtts(speakerWav,
			coqui.WithLanguage(language),
			coqui.WithMaxRetries(a.config.TTS.MaxRetries),
		)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize TTS: %w", err)
	}
	a.tts = tts

	// Generate chapter audio files
	if err := a.generateChapterAudioFiles(book, audiobookInstance, tempDir, finishAudiobook); err != nil {
		return fmt.Errorf("failed to generate chapter audio files: %w", err)
	}

	fmt.Println("\n\n--------------------------------------------------")

	// Generate final audiobook
	if err := audiobookInstance.Generate(distDir); err != nil {
		return fmt.Errorf("failed to generate final audiobook: %w", err)
	}

	// Cleanup
	if err := a.cleanup(tempDir); err != nil {
		fmt.Printf("Warning: cleanup failed: %v\n", err)
	}

	return nil
}

// getImageFile returns the path to the cover image
func (a *Application) getImageFile() string {
	if viper.GetBool("test_mode") {
		fmt.Println("TEST MODE ENABLED: Using mock cover image.")
		return testImagePath
	}

	image := viper.GetString("image_path")
	if image == "" {
		fmt.Println("WARNING: No image path provided in config.")
	}
	return image
}

// getEpubFile loads and returns the EPUB file
func (a *Application) getEpubFile() (*epub.Epub, error) {
	epubPath := viper.GetString("epub_path")

	if viper.GetBool("test_mode") {
		fmt.Println("TEST MODE ENABLED: Using mock epub book.")
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

// setupDirs creates and returns the temp and dist directories
func (a *Application) setupDirs() (string, string, error) {
	tempDir, err := fsutils.GetOrCreateTempDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	distDir := viper.GetString("dist_dir")
	if err := fsutils.CreateDirIfNotExist(distDir); err != nil {
		return "", "", fmt.Errorf("failed to create dist directory: %w", err)
	}

	return tempDir, distDir, nil
}

// chapterFilePath generates the file path for a chapter audio file
func (a *Application) chapterFilePath(chaptersDir, bookTitle string, chapterIdx int) string {
	return fmt.Sprintf("%s/%s-%d.%s", chaptersDir, bookTitle, chapterIdx, processingFileType)
}

// generateChapterAudioFiles processes all chapters and generates audio files
func (a *Application) generateChapterAudioFiles(epubBook *epub.Epub, audiobookInstance *audiobook.Audiobook, tempDir string, finishAudiobook bool) error {
	if len(epubBook.Chapters) == 0 {
		return fmt.Errorf("no chapters found in epub file")
	}

	chaptersDir := fmt.Sprintf("%s/chapters", tempDir)
	if err := fsutils.CreateDirIfNotExist(chaptersDir); err != nil {
		return fmt.Errorf("failed to create chapters directory: %w", err)
	}

	if finishAudiobook {
		fmt.Println("The finish audiobook generation flag is set. Already processed chapters will be used to create the final audiobook.")
	}

	// Loop through chapters
	for i, chapter := range epubBook.Chapters {
		if len(chapter.Paragraphs) == 0 {
			continue
		}

		chapterFile := a.chapterFilePath(chaptersDir, epubBook.Title, i)

		// Check if the chapter was already processed
		if fsutils.FileExists(chapterFile) {
			audiobookInstance.AddChapter(audiobook.AudiobookChapter{
				Title: chapter.Title,
				File:  chapterFile,
			})
			fmt.Printf("Chapter %d already exists. Skipping.\n", i+1)
			continue
		}

		if finishAudiobook {
			continue
		}

		// Process chapter
		if err := a.processChapter(chapter, i, epubBook, tempDir, chapterFile, audiobookInstance); err != nil {
			return fmt.Errorf("failed to process chapter %d: %w", i+1, err)
		}
	}

	return nil
}

// processChapter processes a single chapter
func (a *Application) processChapter(chapter epub.Chapter, index int, epubBook *epub.Epub, tempDir, chapterFile string, audiobookInstance *audiobook.Audiobook) error {
	fmt.Printf("\n\nProcessing Chapter: %s\n", chapter.Title)
	fmt.Println("--------------------------------------------------")

	allContent := append([]string{
		epub.CreateChapterAnnouncement(index, chapter.Title),
	}, chapter.Paragraphs...)

	utils.ParallelForEach(allContent, viper.GetInt("tts.concurrency"), func(i int, text string) {
		if text == "" {
			return
		}

		outputFile := fmt.Sprintf("%s/%d.%s", tempDir, i, processingFileType)

		_, err := a.tts.Synthesize(text, outputFile)
		if err != nil {
			fmt.Println("failed to synthesize text segment", i, err)
		}
	})

	// Get generated audio segment files
	chapterAudioSegmentFiles, err := fsutils.GetFilesFrom(tempDir, processingFileType)
	if err != nil {
		return fmt.Errorf("failed to get audio segment files: %w", err)
	}

	fsutils.SortFilesNumerically(chapterAudioSegmentFiles)

	// Combine all audio segments into a single chapter file
	if err := audiobook.ConcatFiles(epubBook.Title, chapterAudioSegmentFiles, chapterFile); err != nil {
		return fmt.Errorf("failed to concatenate audio files: %w", err)
	}

	audiobookInstance.AddChapter(audiobook.AudiobookChapter{
		Title: chapter.Title,
		File:  chapterFile,
	})

	// Clean up temporary files
	if err := fsutils.RemoveAllFilesInDir(tempDir); err != nil {
		fmt.Printf("Warning: failed to clean up temp files: %v\n", err)
	}

	fmt.Println("Processing: Chapter complete ✅")
	return nil
}

// resetAudiobookGeneration resets the audiobook generation process
func (a *Application) resetAudiobookGeneration(tempDir, distDir string) error {
	fmt.Println("Resetting audiobook generation process...")

	if err := fsutils.EmptyDir(tempDir); err != nil {
		return fmt.Errorf("failed to empty temp directory: %w", err)
	}

	if err := fsutils.EmptyDir(distDir); err != nil {
		return fmt.Errorf("failed to empty dist directory: %w", err)
	}

	return nil
}

// cleanup performs the removal of temporary directories and files
func (a *Application) cleanup(tempDir string) error {
	// Clean and remove chapters directory
	chaptersDir := fmt.Sprintf("%s/chapters", tempDir)
	if err := fsutils.RemoveAllFilesInDir(chaptersDir); err == nil {
		_ = fsutils.RemoveDirIfEmpty(chaptersDir)
	}

	// Clean and remove temp dir
	if err := fsutils.RemoveDirIfEmpty(tempDir); err != nil {
		return fmt.Errorf("unable to remove temp directory: %w", err)
	}

	return nil
}
