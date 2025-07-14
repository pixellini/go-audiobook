package epub

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	epubReader "github.com/taylorskalyo/goreader/epub"
)

type Epub struct {
	Title        string
	Author       string
	Description  string
	Language     string
	Introduction string // This is for the audiobook speaker as a welcome.
	Chapters     []Chapter

	dir string

	// Cleaners for HTML processing
	htmlCleaner      *HTMLCleaner
	contentProcessor *ContentProcessor
}

// New creates a new EPUB instance with the specified directory
func New(dir string) (*Epub, error) {
	if dir == "" {
		return nil, fmt.Errorf("EPUB directory cannot be empty")
	}

	return &Epub{
		dir:              dir,
		htmlCleaner:      NewHTMLCleaner(),
		contentProcessor: NewContentProcessor(),
	}, nil
}

// NewWithMetadata creates a new EPUB with metadata
func NewWithMetadata(dir, title, author, description, language string) (*Epub, error) {
	epub, err := New(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create EPUB: %w", err)
	}

	epub.SetMetadata(dir, title, author, description, language)

	return epub, nil
}

// SetMetadata sets the metadata for the EPUB
func (epub *Epub) SetMetadata(dir, title, author, description, language string) {
	epub.Title = title
	epub.Author = author
	epub.Description = description
	epub.Language = language
}

// LoadMetadata loads the metadata from the EPUB file
func (epub *Epub) LoadMetadata() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	epub.SetMetadata(epub.dir, book.Title, book.Creator, book.Description, book.Language)
	epub.SetIntroduction()

	return nil
}

// LoadChapters loads all chapters from the EPUB file and processes them.
func (epub *Epub) LoadChapters() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	// We don't use the forloop index because some chapters may be skipped
	indexCounter := 1

	for _, item := range book.Manifest.Items {
		r, err := item.Open()
		if err != nil {
			return fmt.Errorf("failed to open item %s: %w", item.ID, err)
		}

		content, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return fmt.Errorf("failed to read item %s: %w", item.ID, err)
		}

		ch := NewChapter(item.ID, string(content))

		ch.Clean()
		ch.LoadTitle()

		// Sometimes the chapter title is the same as the book title, and can create a small chapter with no content.
		chapterHasSameTitleAsBook := strings.TrimSpace(strings.ToLower(ch.Title)) == strings.TrimSpace(strings.ToLower(epub.Title))

		// Check if chapter is valid (includes both ID-based and title-based filtering)
		if !ch.IsValid() || chapterHasSameTitleAsBook {
			continue
		}

		ch.AddIndex(indexCounter)
		// Prepend the chapter announcement
		announcement := ch.GetAnnouncement()
		if announcement != "" {
			ch.Paragraphs = append([]string{announcement}, ch.Paragraphs...)
		}
		indexCounter++

		epub.Chapters = append(epub.Chapters, *ch)
	}

	introduction := Chapter{
		Title:      "Introduction",
		Paragraphs: append([]string{"Introduction."}, SplitText(epub.Introduction)...),
		Content:    epub.Introduction,
		Index:      0, // Set introduction as chapter 0
	}

	// Add the introduction to the beginning of the chapters
	epub.Chapters = append([]Chapter{introduction}, epub.Chapters...)

	return nil
}

// GetDir returns the directory where the EPUB is located
func (epub *Epub) GetDir() string {
	return epub.dir
}

// SetIntroduction sets the introduction text for the audiobook
func (epub *Epub) SetIntroduction() {
	epub.Introduction = fmt.Sprintf("%s written by %s.", epub.Title, epub.Author)
}

// SplitText splits the content into sentences based on punctuation and line breaks.
func SplitText(content string) []string {
	// Use the cleaner from the package
	cleaner := NewHTMLCleaner()
	content = cleaner.CleanHTML(content)

	reSingleLineBreak := regexp.MustCompile(`([^\n])\n([^\n])`)
	for reSingleLineBreak.MatchString(content) {
		content = reSingleLineBreak.ReplaceAllString(content, "$1 $2")
	}

	re := regexp.MustCompile(`(?m)([^\n]+?[\.!?](?:\s+|$)|[^\n]+$)`)
	matches := re.FindAllString(content, -1)
	var result []string
	for _, match := range matches {
		trimmed := strings.TrimSpace(match)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
