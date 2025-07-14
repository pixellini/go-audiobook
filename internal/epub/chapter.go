package epub

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Chapter struct {
	Id         string
	Index      int // This is the index of the chapter in the book
	Title      string
	Paragraphs []string
	Content    string // This will be HTML content
}

var nonNumberedChapterTitles = []string{
	"chapter",
	"part",
	"introduction",
	"acknowledgements",
	"prologue",
	"epilogue",
	"foreword",
	"preface",
	"afterword",
	"conclusion",
	"appendix",
	"glossary",
	"index",
}

var chapterHeaderRegex = regexp.MustCompile(`(?i)^(chapter|part|section|page)\s*\d*\.?\s*`)

// CreateChapterAnnouncement generates a chapter announcement string based on the chapter index and title.
// TODO: I want this to be the first line of the chapter
func CreateChapterAnnouncement(chapterIndex int, title string) string {
	lowerTitle := strings.ToLower(strings.TrimSpace(title))

	for _, prefix := range nonNumberedChapterTitles {
		if strings.HasPrefix(lowerTitle, prefix) {
			return title
		}
	}

	// E.g. "Chapter 1: The Beginning"
	return fmt.Sprintf("Chapter %d: %s", chapterIndex, title)
}

// NewChapter creates a new Chapter instance with the given ID and HTML content.
func NewChapter(id, htmlContent string) *Chapter {
	return &Chapter{
		Id:         id,
		Title:      "", // Title will be set later
		Paragraphs: []string{},
		Content:    htmlContent,
	}
}

// LoadTitle extracts the title from the chapter content.
func (c *Chapter) LoadTitle() {
	cleaner := NewHTMLCleaner()
	cleanedContent := SplitText(cleaner.CleanContent(c.Content))

	title := ""

	titleTagRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := titleTagRegex.FindStringSubmatch(c.Content)
	if len(matches) > 1 {
		title := matches[1]
		// Reject titles that look like filenames (e.g., "ch001.xhtml")
		// This seems to happen a lot when the file is converted to EPUB.
		if strings.HasSuffix(title, ".xhtml") || strings.HasSuffix(title, ".html") {
			title = ""
		}
	}

	if title == "" || title == "Converted Ebook" {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(c.Content))
		if err == nil {
			heading := ""
			doc.Find("h1, h2, h3, h4, h5, h6").EachWithBreak(func(i int, s *goquery.Selection) bool {
				heading = strings.TrimSpace(s.Text())
				return heading == ""
			})
			if heading != "" {
				title = heading
			}
		}
	}

	if title == "" || title == "Converted Ebook" {
		for _, line := range cleanedContent {
			if strings.TrimSpace(line) != "" {
				title = strings.TrimSpace(line)
				break
			}
		}
	}

	c.Title = title
}

// Clean processes the chapter content to remove HTML tags, clean up text, and split into paragraphs.
func (c *Chapter) Clean() {
	cleaner := NewHTMLCleaner()
	processor := NewContentProcessor()

	cleanedContent := SplitText(cleaner.CleanContent(c.Content))

	for i, line := range cleanedContent {
		if chapterHeaderRegex.MatchString(line) {
			cleanedContent[i] = chapterHeaderRegex.ReplaceAllString(line, "")
		}

		contentWithoutUrls := processor.ReplaceURLs(cleanedContent[i])
		cleanedContent[i] = processor.ReplaceCodePatterns(contentWithoutUrls)
	}

	c.Paragraphs = cleanedContent
}

// IsValid checks if the chapter is a valid chapter based on its ID and title.
// It filters out both invalid chapter types and unwanted content like table of contents.
// This method should be called after LoadTitle() has been called.
func (c *Chapter) IsValid() bool {
	// First check if it's a valid chapter type based on ID
	id := strings.ToLower(c.Id)
	isValidType := strings.Contains(id, "title") ||
		strings.Contains(id, "dedication") ||
		strings.Contains(id, "foreword") ||
		strings.Contains(id, "preface") ||
		strings.HasPrefix(id, "ch") ||
		strings.HasPrefix(id, "id")

	if !isValidType {
		return false
	}

	// Then check if the chapter title should be filtered out
	t := strings.TrimSpace(strings.ToLower(c.Title))
	shouldFilter := t == "contents" || t == "table of contents"

	return !shouldFilter
}

// AddIndex adds an index to the chapter ID for easier identification.
func (c *Chapter) AddIndex(index int) {
	c.Index = index
}

func (c *Chapter) GetAnnouncement() string {
	lowerTitle := strings.ToLower(strings.TrimSpace(c.Title))

	// Handle introduction and other special chapters
	for _, prefix := range nonNumberedChapterTitles {
		if strings.HasPrefix(lowerTitle, prefix) {
			return c.Title
		}
	}

	// E.g. "Chapter 1: The Beginning"
	return fmt.Sprintf("Chapter %d: %s", c.Index, c.Title)
}
