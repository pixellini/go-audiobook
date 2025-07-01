package epub

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pixellini/go-audiobook/internal/formatter"
	epubReader "github.com/taylorskalyo/goreader/epub"
)

type Epub struct {
	Title        string
	Author       string
	Description  string
	Language     string
	Introduction string // This is for the audiobook speaker as a welcome.
	Chapters     []Chapter
	dir          string
}

const (
	chapterHeaderPattern = `(?i)^(chapter|part|section|page)\s*\d*\.?\s*`
	titleTagPattern      = `<title>(.*?)</title>`
)

var (
	titleTagRegex = regexp.MustCompile(titleTagPattern)
)

func New(dir string) (*Epub, error) {
	epub := &Epub{
		dir: dir,
	}

	if err := epub.setMetadata(); err != nil {
		return epub, err
	}

	if err := epub.setChapters(); err != nil {
		return epub, err
	}

	return epub, nil
}

func (epub *Epub) setMetadata() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	epub.Title = book.Title
	epub.Description = book.Description
	epub.Language = book.Language
	epub.Author = book.Creator
	epub.Introduction = fmt.Sprintf("%s written by %s.", book.Title, book.Creator)

	return nil
}

func (epub *Epub) setChapters() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	// Filter only the chapters to read. Ignore contents, images, indexes etc.
	var chapters []epubReader.Item
	for _, item := range book.Manifest.Items {
		if isAcceptedChapterItem(item.ID) {
			chapters = append(chapters, item)
		}
	}

	for _, item := range chapters {
		r, err := item.Open()
		if err != nil {
			return fmt.Errorf("failed to open item %s: %w", item.ID, err)
		}

		content, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return fmt.Errorf("failed to read item %s: %w", item.ID, err)
		}

		rawHtmlContent := string(content)
		processedChapter := createChapter(rawHtmlContent)
		epub.Chapters = append(epub.Chapters, processedChapter)
	}

	introduction := Chapter{
		Title:      "Introduction",
		Paragraphs: formatter.SplitText(epub.Introduction),
		RawText:    epub.Introduction,
	}

	// Add the introduction to the beginning of the chapters
	epub.Chapters = append([]Chapter{introduction}, epub.Chapters...)

	// Remove chapters titled 'contents' or 'table of contents' (case-insensitive, trimmed)
	filtered := make([]Chapter, 0, len(epub.Chapters))
	for _, ch := range epub.Chapters {
		t := strings.TrimSpace(strings.ToLower(ch.Title))
		if t == "contents" || t == "table of contents" {
			continue
		}
		filtered = append(filtered, ch)
	}
	epub.Chapters = filtered

	return nil
}
