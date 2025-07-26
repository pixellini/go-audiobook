package epub

import (
	"fmt"
	"strings"

	"github.com/pixellini/go-audiobook/internal/epubreader"
)

type Epub struct {
	Metadata *EpubMetadata
	Chapters []EpubChapter
}

type EpubMetadata struct {
	Title       string
	Author      string
	Description string
	Language    string
	Publisher   string
	CoverImage  string
}

func New() *Epub {
	return &Epub{}
}

func NewFromFile(r epubreader.EpubReader) (*Epub, error) {
	if r == nil {
		return nil, fmt.Errorf("epub reader must be provided")
	}

	epub := New()

	epub.LoadMetadata(r)

	return epub, nil
}

func (e *Epub) LoadMetadata(r epubreader.EpubReader) {
	e.Metadata = &EpubMetadata{
		Title:       r.GetTitle(),
		Author:      r.GetAuthor(),
		Description: r.GetDescription(),
		Language:    r.GetLanguage(),
		Publisher:   r.GetPublisher(),
		CoverImage:  r.GetCoverImage(),
	}
}

type EpubChapter struct {
	Id string
	// Title is the title of the chapter.
	Title string
	// Content is the raw, unedited HTML chapter content.
	Content string

	Path string
}

func NewChapter(id, title, content string) (*EpubChapter, error) {
	if id == "" {
		return nil, fmt.Errorf("chapter id must be provided")
	}

	return &EpubChapter{
		Id:      id,
		Title:   title,
		Content: content,
	}, nil
}

func (c EpubChapter) IsValid() bool {
	// First check if it's a valid chapter type based on ID
	id := strings.ToLower(c.Id)
	isValidType :=
		strings.Contains(id, "dedication") ||
			// strings.Contains(id, "title") ||
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
