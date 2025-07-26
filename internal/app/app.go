package app

import (
	"context"
	"fmt"
	"log"
	"os"
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
}

func New() (*Application, error) {
	// Load configuration
	c, err := config.Load()
	if err != nil {
		return nil, err
	}

	fm, err := filemanager.NewService("./.temp")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file manager")
	}

	// Initialize TTS service
	tts, err := ttsservice.NewCoquiService(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TTS service: %w", err)
	}

	ffmpeg := audioservice.NewFFMpegService()

	return &Application{
		config:      c,
		fileManager: fm,
		tts:         tts,
		audio:       ffmpeg,
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
	app.fileManager.Remove("./.temp/")
	// Load EPUB file
	epubPath := app.config.Epub.Path
	r, err := epubreader.NewGoEpubReaderService(epubPath)
	if err != nil {
		return err
	}

	defer r.Close()

	book, err := epub.NewFromFile(r)
	if err != nil {
		return err
	}

	fmt.Println(book.Metadata)

	rawChapters, err := r.GetChapters()
	if err != nil {
		return err
	}

	chapters, err := app.ProcessChapters(ctx, rawChapters)
	if err != nil {
		panic(err)
	}

	tempAudiobookFile := app.config.TempDir + "audiobook.wav"

	err = app.CombineChapters(chapters, tempAudiobookFile)
	if err != nil {
		panic(err)
	}

	for _, file := range chapters {
		err := os.Remove(file.Path)
		if err != nil {
			fmt.Printf("\nunable to remove file: %s", err)
		}
	}

	// Generate audiobook
	metadata, err := app.BuildMetadataFile(chapters)
	if err != nil {
		panic(err)
	}
	defer metadata.Close()

	err = app.audio.CreateAudiobook(
		tempAudiobookFile,
		app.config.Epub.CoverImage,
		metadata.Name(),
		app.config.Output.FullPath(),
	)
	if err != nil {
		panic(err)
	}

	_ = os.Remove(tempAudiobookFile)

	return nil
}

func (app *Application) ProcessChapters(ctx context.Context, chapters []*epubreader.EpubReaderChapter) ([]*epub.EpubChapter, error) {
	var processedChapters []*epub.EpubChapter

	chapterNumber := 1
	for _, chapter := range chapters[:6] {
		ch, err := epub.NewChapter(chapter.Id, chapter.Title, chapter.Content)
		if err != nil || !ch.IsValid() {
			continue
		}
		files, err := app.CreateChapterAudio(ctx, ch, chapterNumber)
		if err != nil {
			// TODO: fix this error
			return nil, fmt.Errorf("unable to create chapter audio")
		}

		fsutils.SortNumerically(files)

		chapter.Path = fmt.Sprintf("%s/chapter-%d.wav", app.config.TempDir, chapterNumber)
		err = app.audio.CombineFiles(files, chapter.Path)
		if err != nil {
			fmt.Printf("\nerror combining files for chapter %d: %v", chapterNumber, err)
		}

		processedChapters = append(processedChapters, ch)

		for _, file := range files {
			if err = os.Remove(file); err != nil {
				fmt.Printf("\nunable to remove file: %s", err)
			}
		}

		chapterNumber++
	}

	fmt.Println("processed", len(processedChapters))

	return processedChapters, nil
}

func (app *Application) CreateChapterAudio(ctx context.Context, chapter *epub.EpubChapter, chapterNumber int) ([]string, error) {
	// app.fileManager.CreateTemp()
	// defer app.fileManager.RemoveTemp()

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

	var files []string
	for i, paragraph := range text {
		if paragraph == "" {
			continue
		}
		sem <- struct{}{} // acquire a slot
		wg.Add(1)

		go func(p string, idx int) {
			fileName := fmt.Sprintf("part-%d.wav", idx+1)
			defer wg.Done()
			defer func() { <-sem }() // release the slot

			_, err := app.tts.SynthesizeContext(ctx, text[idx], fileName)
			if err != nil {
				log.Printf("error synthesizing chapter %d part %d: %v", chapterNumber, idx, err)
			}
			files = append(files, fmt.Sprintf("%s/%s", app.config.TempDir, fileName))
		}(paragraph, i)
	}

	wg.Wait()

	return files, nil
}

func (app *Application) BuildMetadataFile(chapters []*epub.EpubChapter) (*os.File, error) {
	file, err := os.CreateTemp(app.config.TempDir, "metadata-*.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file.WriteString(";FFMETADATA1\n")

	startTime := 0

	for _, chapter := range chapters {
		// Get the chapter duration
		// Then append it to the startTime so that we can calculate the endTime for the next chapter.
		duration, err := app.audio.GetDuration(chapter.Path)
		if err != nil {
			return nil, err
		}

		durationMs := int(duration * 1000)
		endTime := startTime + durationMs

		file.WriteString("[CHAPTER]\n")
		file.WriteString("TIMEBASE=1/1000\n")
		file.WriteString(fmt.Sprintf("START=%d\n", startTime))
		file.WriteString(fmt.Sprintf("END=%d\n", endTime))
		file.WriteString(fmt.Sprintf("title=%s\n\n", chapter.Title))

		startTime = endTime
	}

	return file, nil
}

func (app *Application) CombineChapters(chapters []*epub.EpubChapter, output string) error {
	var files []string
	for _, c := range chapters {
		fmt.Println("Path", c.Path)
		files = append(files, c.Path)
	}

	fmt.Println("-------FILES", files)
	err := app.audio.CombineFiles(files, output)
	if err != nil {
		return err
	}
	return nil
}
