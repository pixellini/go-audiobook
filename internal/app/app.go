package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pixellini/go-audiobook/internal/audioservice"
	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/epub"
	"github.com/pixellini/go-audiobook/internal/epubreader"
	"github.com/pixellini/go-audiobook/internal/filemanager"
	"github.com/pixellini/go-audiobook/internal/flags"
	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/textutils"
	"github.com/pixellini/go-audiobook/internal/ttsservice"
)

type Application struct {
	config      *config.Config
	fileManager filemanager.FileService
	tts         ttsservice.TTSservice
	audio       audioservice.AudioService

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
	// cacheDir := "./.cache/"
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file manager")
	}

	// Initialize TTS service
	tts, err := ttsservice.NewCoquiService(c, cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS service: %w", err)
	}

	ffmpeg := audioservice.NewFFMpegService(cacheDir)

	return &Application{
		config:      c,
		fileManager: fm,
		tts:         tts,
		audio:       ffmpeg,
		cacheDir:    cacheDir,
	}, nil
}

func (app *Application) Run() error {
	return app.RunContext(context.Background())
}

func (app *Application) RunContext(ctx context.Context) error {
	return app.run(ctx, nil)
}

func (app *Application) RunWithFlags(fl *flags.Flags) error {
	return app.RunWithFlagsContext(context.Background(), fl)
}

func (app *Application) RunWithFlagsContext(ctx context.Context, fl *flags.Flags) error {
	return app.run(ctx, nil)
}

func (app *Application) run(ctx context.Context, _ *flags.Flags) error {
	app.fileManager.Remove(app.config.Output.Path)
	defer app.fileManager.Remove(app.cacheDir)

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

	fmt.Println(book.Metadata)

	rawChapters, err := r.GetChapters()
	if err != nil {
		return err
	}

	chapters, err := app.ProcessChapters(ctx, rawChapters)
	if err != nil {
		return fmt.Errorf("failed to process chapters: %w", err)
	}

	tempAudiobookFile := filepath.Join(app.cacheDir, "audiobook.wav")

	metadata, err := app.BuildMetadataFile(chapters)
	if err != nil {
		return fmt.Errorf("failed to create metadata: %w", err)
	}
	defer metadata.Close()
	defer os.Remove(metadata.Name())

	// Turn all the chapter wav files into a singular wav file.
	err = app.CombineChapters(chapters, tempAudiobookFile)
	if err != nil {
		return fmt.Errorf("failed to combine chapter files: %w", err)
	}

	app.fileManager.Create(app.config.Output.Path)

	// Create the M4B audiobook file.
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

	return nil
}

func (app *Application) ProcessChapters(ctx context.Context, chapters []*epubreader.EpubReaderChapter) ([]*epub.EpubChapter, error) {
	processedChapters := make([]*epub.EpubChapter, 0, len(chapters))

	chapterNumber := 1
	for _, chapter := range chapters[:6] {
		ch, err := epub.NewChapter(chapter.Id, chapter.Title, chapter.Content)
		if err != nil || !ch.IsValid() {
			continue
		}
		files, err := app.CreateChapterAudio(ctx, ch, chapterNumber)
		if err != nil {
			return nil, fmt.Errorf("unable to create chapter audio for chapter %d: %w", chapterNumber, err)
		}

		if len(files) == 0 {
			log.Printf("No audio files created for chapter %d, skipping", chapterNumber)
			continue
		}

		fsutils.SortNumerically(files)

		ch.Path = filepath.Join(app.cacheDir, fmt.Sprintf("chapter-%d.wav", chapterNumber))

		// Check if chapter file already exists
		if _, err := os.Stat(ch.Path); err == nil {
			log.Printf("Chapter %d audio file already exists, skipping combination", chapterNumber)
		} else {
			if len(files) == 0 {
				return nil, fmt.Errorf("no audio files available for chapter %d", chapterNumber)
			}
			err = app.audio.CombineFiles(files, ch.Path)
			if err != nil {
				return nil, fmt.Errorf("error combining files for chapter %d: %w", chapterNumber, err)
			}
		}

		processedChapters = append(processedChapters, ch)

		for _, file := range files {
			if err = os.Remove(file); err != nil {
				fmt.Printf("\nunable to remove file: %s", err)
			}
		}

		chapterNumber++
	}

	return processedChapters, nil
}

func (app *Application) CreateChapterAudio(ctx context.Context, chapter *epub.EpubChapter, chapterNumber int) ([]string, error) {
	// Split up the text and clean the content
	text := textutils.ExtractParagraphsFromHTML(chapter.Content)
	if len(text) == 0 {
		return nil, fmt.Errorf("chapter does not have text")
	}

	// We can use the first paragraph as the title
	chapter.Title = fmt.Sprintf("Chapter %d: %s", chapterNumber, text[0])
	text[0] = chapter.Title // Use the first paragraph as the title

	sem := make(chan struct{}, app.config.Model.Concurrency)
	var wg sync.WaitGroup

	// Use a mutex to protect the files slice
	var mu sync.Mutex
	var files []string

	for i, paragraph := range text[:3] {
		if paragraph == "" {
			continue
		}

		fileName := fmt.Sprintf("part-%d.wav", i+1)
		filePath := filepath.Join(app.cacheDir, fileName)

		// Check if file already exists
		if _, err := os.Stat(filePath); err == nil {
			// File exists, add it to the list
			mu.Lock()
			files = append(files, filePath)
			mu.Unlock()
			continue
		}

		// File doesn't exist, need to synthesize
		sem <- struct{}{} // acquire a slot
		wg.Add(1)

		go func(p string, idx int, fPath string) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot

			_, err := app.tts.SynthesizeContext(ctx, text[idx], fileName)
			if err != nil {
				log.Printf("error synthesizing chapter %d part %d: %v", chapterNumber, idx, err)
				return // Don't add failed files
			}

			// Verify file was created
			if _, err := os.Stat(fPath); err != nil {
				log.Printf("synthesized file not found for chapter %d part %d: %v", chapterNumber, idx, err)
				return
			}

			// Thread-safe append
			mu.Lock()
			files = append(files, fPath)
			mu.Unlock()
		}(paragraph, i, filePath)
	}

	wg.Wait()

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
	for _, file := range chapters {
		err := os.Remove(file.Path)
		if err != nil {
			fmt.Printf("\nunable to remove file: %s", err)
		}
	}

	return nil
}

func (app *Application) BuildMetadataFile(chapters []*epub.EpubChapter) (*os.File, error) {
	file, err := os.CreateTemp(app.cacheDir, "metadata-*.txt")
	if err != nil {
		return nil, err
	}
	// Don't defer close here since we're returning the file

	_, err = file.WriteString(";FFMETADATA1\n")
	if err != nil {
		file.Close()
		return nil, err
	}

	startTime := 0

	for _, chapter := range chapters {
		// Verify chapter file exists before trying to get duration
		if _, err := os.Stat(chapter.Path); err != nil {
			file.Close()
			return nil, fmt.Errorf("chapter file does not exist: %s - %w", chapter.Path, err)
		}

		// Get the chapter duration
		// Then append it to the startTime so that we can calculate the endTime for the next chapter.
		duration, err := app.audio.GetDuration(chapter.Path)
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to get duration for chapter file %s: %w", chapter.Path, err)
		}

		durationMs := int(duration * 1000)
		endTime := startTime + durationMs

		_, err = file.WriteString("[CHAPTER]\n")
		if err != nil {
			file.Close()
			return nil, err
		}
		_, err = file.WriteString("TIMEBASE=1/1000\n")
		if err != nil {
			file.Close()
			return nil, err
		}
		_, err = file.WriteString(fmt.Sprintf("START=%d\n", startTime))
		if err != nil {
			file.Close()
			return nil, err
		}
		_, err = file.WriteString(fmt.Sprintf("END=%d\n", endTime))
		if err != nil {
			file.Close()
			return nil, err
		}
		_, err = file.WriteString(fmt.Sprintf("title=%s\n\n", chapter.Title))
		if err != nil {
			file.Close()
			return nil, err
		}

		startTime = endTime
	}

	// Flush and seek to beginning for reading
	err = file.Sync()
	if err != nil {
		file.Close()
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}
