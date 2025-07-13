package processor

import (
	"context"
	"fmt"

	"github.com/pixellini/go-audiobook/internal/audiobook"
	"github.com/pixellini/go-audiobook/internal/config"
	"github.com/pixellini/go-audiobook/internal/epub"
	"github.com/pixellini/go-audiobook/internal/fs"
	"github.com/pixellini/go-audiobook/internal/interfaces"
	"github.com/pixellini/go-audiobook/internal/utils"
)

const processingFileType = "wav"

// ChapterProcessor handles the processing of individual chapters
type ChapterProcessor struct {
	tts       interfaces.TTSService
	config    *config.Config
	logger    interfaces.Logger
	tempDir   string
	bookTitle string
}

// NewChapterProcessor creates a new chapter processor
func NewChapterProcessor(tts interfaces.TTSService, config *config.Config, logger interfaces.Logger, tempDir, bookTitle string) *ChapterProcessor {
	return &ChapterProcessor{
		tts:       tts,
		config:    config,
		logger:    logger,
		tempDir:   tempDir,
		bookTitle: bookTitle,
	}
}

// ProcessChapters processes all chapters in the EPUB
func (cp *ChapterProcessor) ProcessChapters(ctx context.Context, epubBook *epub.Epub, audiobookInstance *audiobook.Audiobook, finishAudiobook bool) error {
	if len(epubBook.Chapters) == 0 {
		return fmt.Errorf("no chapters found in epub file")
	}

	chaptersDir := fmt.Sprintf("%s/chapters", cp.tempDir)
	if err := fs.CreateDirIfNotExist(chaptersDir); err != nil {
		return fmt.Errorf("failed to create chapters directory: %w", err)
	}

	if finishAudiobook {
		cp.logger.Info("Finish audiobook generation flag is set. Already processed chapters will be used.")
	}

	for i, chapter := range epubBook.Chapters {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if len(chapter.Paragraphs) == 0 {
			continue
		}

		if err := cp.processChapter(ctx, chapter, i, chaptersDir, audiobookInstance, finishAudiobook); err != nil {
			return fmt.Errorf("failed to process chapter %d: %w", i+1, err)
		}
	}

	return nil
}

// processChapter processes a single chapter
func (cp *ChapterProcessor) processChapter(ctx context.Context, chapter epub.Chapter, index int, chaptersDir string, audiobookInstance *audiobook.Audiobook, finishAudiobook bool) error {
	chapterFile := cp.chapterFilePath(chaptersDir, index)

	// Check if the chapter was already processed
	if fs.FileExists(chapterFile) {
		audiobookInstance.AddChapter(audiobook.AudiobookChapter{
			Title: chapter.Title,
			File:  chapterFile,
		})
		cp.logger.Infof("Chapter %d already exists. Skipping.", index+1)
		return nil
	}

	if finishAudiobook {
		return nil
	}

	cp.logger.Infof("Processing Chapter: %s", chapter.Title)

	allContent := append([]string{
		epub.CreateChapterAnnouncement(index, chapter.Title),
	}, chapter.Paragraphs...)

	if err := cp.synthesizeContent(ctx, allContent); err != nil {
		return fmt.Errorf("failed to synthesize content: %w", err)
	}

	if err := cp.combineAudioSegments(chapterFile); err != nil {
		return fmt.Errorf("failed to combine audio segments: %w", err)
	}

	audiobookInstance.AddChapter(audiobook.AudiobookChapter{
		Title: chapter.Title,
		File:  chapterFile,
	})

	// Clean up temporary files
	if err := fs.RemoveAllFilesInDir(cp.tempDir); err != nil {
		cp.logger.Warnf("Failed to clean up temp files: %v", err)
	}

	cp.logger.Info("Chapter processing complete ✅")
	return nil
}

// synthesizeContent synthesizes all content using TTS
func (cp *ChapterProcessor) synthesizeContent(ctx context.Context, content []string) error {
	concurrency := cp.config.TTS.Concurrency
	if concurrency <= 0 {
		concurrency = 4 // default fallback
	}

	utils.ParallelForEach(content, concurrency, func(i int, text string) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if text == "" {
			return
		}

		outputFile := fmt.Sprintf("%s/%d.%s", cp.tempDir, i, processingFileType)

		_, err := cp.tts.Synthesize(text, outputFile)
		if err != nil {
			cp.logger.Errorf("Failed to synthesize text segment %d: %v", i, err)
		}
	})

	return nil
}

// combineAudioSegments combines all audio segments into a single chapter file
func (cp *ChapterProcessor) combineAudioSegments(chapterFile string) error {
	chapterAudioSegmentFiles, err := fs.GetFilesFrom(cp.tempDir, processingFileType)
	if err != nil {
		return fmt.Errorf("failed to get audio segment files: %w", err)
	}

	fs.SortFilesNumerically(chapterAudioSegmentFiles)

	if err := audiobook.ConcatFiles(cp.bookTitle, chapterAudioSegmentFiles, chapterFile); err != nil {
		return fmt.Errorf("failed to concatenate audio files: %w", err)
	}

	return nil
}

// chapterFilePath generates the file path for a chapter audio file
func (cp *ChapterProcessor) chapterFilePath(chaptersDir string, chapterIdx int) string {
	return fmt.Sprintf("%s/%s-%d.%s", chaptersDir, cp.bookTitle, chapterIdx, processingFileType)
}
