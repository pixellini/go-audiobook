package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pixellini/go-audiobook/internal/audioservice"
	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/epub"
	"github.com/pixellini/go-audiobook/internal/epubreader"
	"github.com/pixellini/go-audiobook/internal/filemanager"
	"github.com/pixellini/go-audiobook/internal/flags"
	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/logger"
	"github.com/pixellini/go-audiobook/internal/metadata"
	"github.com/pixellini/go-audiobook/internal/textutils"
	"github.com/pixellini/go-audiobook/internal/ttsservice"
	"github.com/pixellini/go-audiobook/internal/tui"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	config      *config.Config
	fileManager filemanager.FileService
	tts         ttsservice.TTSservice
	audio       audioservice.AudioService
	flag        *flags.Flags
	tui         tui.TUIService
	logger      logger.Logger

	chapters []*epubreader.EpubReaderChapter
	cacheDir string
}

func New() (*Application, error) {
	// Load configuration
	c, err := config.Load()
	if err != nil {
		return nil, err
	}

	fm := filemanager.New()
	cacheDir, err := fm.CreateCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file manager")
	}

	// Initialize TTS service
	tts, err := ttsservice.NewCoquiService(c, cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS service: %w", err)
	}

	ffmpeg := audioservice.NewFFMpegService(cacheDir)

	// Create logger based on config
	var (
		l logger.Logger
		t tui.TUIService
	)
	if c.VerboseLogs {
		l = logger.NewLogger()
		t = tui.NewEmpty()
	} else {
		t = tui.NewBubbleTeaUI()
		l = logger.NewSilentLogger()
	}

	return &Application{
		config:      c,
		fileManager: fm,
		tts:         tts,
		audio:       ffmpeg,
		tui:         t,
		logger:      l,
		cacheDir:    cacheDir,
	}, nil
}

func NewWithFlags(fl *flags.Flags) (*Application, error) {
	app, err := New()
	if err != nil {
		return app, err
	}

	app.flag = fl

	return app, nil
}

func (app *Application) Run() error {
	return app.RunContext(context.Background())
}

func (app *Application) RunContext(ctx context.Context) error {
	return app.run(ctx)
}

func (app *Application) run(ctx context.Context) error {
	defer app.fileManager.Remove(app.cacheDir)

	if app.flag.ResetProgress {
		app.Reset()
	}

	// Use TUI unless verbose logging is enabled in config
	// Start the TUI for progress tracking
	if err := app.tui.Start(); err != nil {
		log.Printf("Failed to start TUI: %v", err)
	} else {
		defer app.tui.Stop()
	}

	start := time.Now()

	epubPath := app.config.Epub.Path
	r, err := epubreader.NewGoEpubReaderService(epubPath)
	if err != nil {
		return err
	}

	defer r.Close()

	book, err := epub.NewFromFile(r)
	if err != nil {
		return fmt.Errorf("failed to read epub file: %w", err)
	}
	if app.config.Output.Filename == "" {
		app.config.Output.Filename = book.Metadata.Title
	}

	if fsutils.FileExists(app.config.Output.FullPath()) {
		return fmt.Errorf("File '%s' has already been created.", app.config.Output.OutputFileName())
	}

	rawChapters, err := r.GetChapters()
	if err != nil {
		return err
	}

	chapters, err := app.ProcessChapters(ctx, rawChapters)
	if err != nil {
		return fmt.Errorf("failed to process chapters: %w", err)
	}

	tempAudiobookFile := filepath.Join(app.cacheDir, "audiobook.wav")

	metadata, err := app.BuildMetadataFile(book.Metadata, chapters)

	if err != nil {
		return fmt.Errorf("failed to create metadata: %w", err)
	}
	defer metadata.Remove()

	// Close the metadata file to ensure all data is written before FFmpeg uses it
	err = metadata.Close()
	if err != nil {
		return fmt.Errorf("failed to close metadata file: %w", err)
	}

	// Turn all the chapter wav files into a singular wav file.
	app.tui.UpdateProgress("Combining all chapters into audiobook...")
	err = app.CombineChapters(chapters, tempAudiobookFile)
	if err != nil {
		return fmt.Errorf("failed to combine chapter files: %w", err)
	}

	app.fileManager.Create(app.config.Output.Path)

	// Create the M4B audiobook file.
	app.tui.UpdateProgress("Creating final audiobook file...")
	err = app.audio.CreateAudiobook(
		tempAudiobookFile,
		app.config.Epub.CoverImage,
		metadata.Name(),
		app.config.Output.FullPath(),
	)
	if err != nil {
		return fmt.Errorf("failed to create audiobook file: %w", err)
	}

	_ = os.Remove(tempAudiobookFile)

	// Show completion message
	completionTime := time.Since(start).Truncate(time.Second)
	completionMsg := fmt.Sprintf("🎉 Audiobook \"%s\" created successfully! (%v)", app.config.Output.OutputFileName(), completionTime)
	if app.tui != nil {
		app.tui.Finish(completionMsg)
	} else {
		fmt.Print(completionMsg)
	}

	return nil
}

func (app *Application) ProcessChapters(ctx context.Context, chapters []*epubreader.EpubReaderChapter) ([]*epub.EpubChapter, error) {
	processedChapters := make([]*epub.EpubChapter, 0, len(chapters))

	chapterNumber := 1
	for _, chapter := range chapters {
		ch, err := epub.NewChapter(chapter.Id, chapter.Title, chapter.Content)
		if err != nil || !ch.IsValid() {
			continue
		}
		ch.Path = filepath.Join(app.cacheDir, fmt.Sprintf("chapter-%d.wav", chapterNumber))

		// Update TUI with current chapter being processed
		app.tui.UpdateProgress(fmt.Sprintf("Processing — %s", ch.Title))

		if fsutils.FileExists(ch.Path) {
			processedChapters = append(processedChapters, ch)
			// Mark chapter as complete (cached)
			app.tui.CompleteCurrentTask(fmt.Sprintf("Chapter %d (cached): %s", chapterNumber, ch.Title))
			chapterNumber++
			continue
		}

		if app.flag.FinishAudiobook {
			// Skip this chapter if we're finishing audiobook and the file doesn't exist
			chapterNumber++
			continue
		}

		files, err := app.CreateChapterAudio(ctx, ch, chapterNumber)
		if err != nil {
			return nil, fmt.Errorf("unable to create chapter audio for chapter %d: %w", chapterNumber, err)
		}

		if len(files) == 0 {
			app.logger.Printf("No audio files created for chapter %d, skipping", chapterNumber)
			continue
		}

		// Sort the paragraph files numerically, because they might be out of order due to concurrency.
		fsutils.SortNumerically(files)

		err = app.audio.CombineFiles(files, ch.Path)
		if err != nil {
			return nil, fmt.Errorf("error combining files for chapter %d: %w", chapterNumber, err)
		}

		processedChapters = append(processedChapters, ch)

		// Mark chapter as complete
		app.tui.CompleteCurrentTask(fmt.Sprintf("Chapter %d completed: %s", chapterNumber, ch.Title))

		app.fileManager.RemoveFiles(files)

		chapterNumber++
	}

	return processedChapters, nil
}

func (app *Application) CreateChapterAudio(ctx context.Context, chapter *epub.EpubChapter, chapterNumber int) ([]string, error) {
	text := textutils.ExtractParagraphsFromHTML(chapter.Content)
	if len(text) == 0 {
		return nil, fmt.Errorf("chapter does not have text")
	}

	chapter.Title = fmt.Sprintf("Chapter %d: %s", chapterNumber, text[0])
	text[0] = chapter.Title

	var files []string
	var mu sync.Mutex
	var completed int
	totalParagraphs := len(text)

	// Initialize progress
	app.tui.UpdateProgressWithBar(
		fmt.Sprintf("Processing — %s", chapter.Title), 0, totalParagraphs,
	)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(int(app.config.Model.Concurrency))

	add := func(p string) error {
		mu.Lock()
		files = append(files, p)
		completed++
		currentCompleted := completed
		mu.Unlock()

		// Update progress bar
		app.tui.UpdateProgressWithBar(
			fmt.Sprintf("Processing — %s", chapter.Title), currentCompleted, totalParagraphs,
		)

		return nil
	}

	for i, p := range text {
		if p == "" {
			continue
		}
		i, p := i, p // capture
		eg.Go(func() error {
			name := fmt.Sprintf("paragraph-%d.wav", i+1)
			path := filepath.Join(app.cacheDir, name)

			if fsutils.FileExists(path) {
				return add(path)
			}

			if app.flag.FinishAudiobook {
				return nil
			}

			if _, err := app.tts.SynthesizeContext(ctx, p, name); err != nil {
				return fmt.Errorf("chapter %d part %d: %w", chapterNumber, i+1, err)
			}
			if !fsutils.FileExists(path) {
				return fmt.Errorf("synthesised file missing for chapter %d part %d", chapterNumber, i+1)
			}

			return add(path)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return files, nil
}

func (app *Application) CombineChapters(chapters []*epub.EpubChapter, output string) error {
	var files []string
	for _, c := range chapters {
		files = append(files, c.Path)
	}

	err := app.audio.CombineFiles(files, output)
	if err != nil {
		return err
	}

	// No longer need the chapter audio files
	app.fileManager.RemoveFiles(files)

	return nil
}

func (app *Application) BuildMetadataFile(book *epub.EpubMetadata, chapters []*epub.EpubChapter) (*metadata.Metadata, error) {
	metaFile, err := metadata.New(app.cacheDir)
	if err != nil {
		return nil, err
	}
	metaFile.AddDetails(book)

	startTime := 0

	for _, chapter := range chapters {
		// Verify chapter file exists before trying to get duration
		if _, err := os.Stat(chapter.Path); err != nil {
			metaFile.Close()
			return nil, fmt.Errorf("chapter file does not exist: %s - %w", chapter.Path, err)
		}

		// Get the chapter duration
		// Then append it to the startTime so that we can calculate the endTime for the next chapter.
		duration, err := app.audio.GetDuration(chapter.Path)
		if err != nil {
			metaFile.Close()
			return nil, fmt.Errorf("failed to get duration for chapter: %s: %w", chapter.Title, err)
		}

		durationMs := int(duration * 1000)
		endTime := startTime + durationMs

		err = metaFile.AddChapter(chapter.Title, startTime, endTime)
		if err != nil {
			metaFile.Close()
			return nil, fmt.Errorf("failed to create metadata for chapter: %s: %w", chapter.Title, err)
		}

		startTime = endTime
	}

	return metaFile, nil
}

func (app *Application) Reset() {
	app.fileManager.Remove(app.config.Output.Path)
	app.fileManager.Remove(app.cacheDir)
}

func (app *Application) ToggleTUI() {

}
