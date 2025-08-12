package epub

import (
	"fmt"
	"slices"
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

// Validation lists
var excludedTypes = []string{"css", "ncx", "stylesheet", "style"}
var invalidTitlePrefixes = []string{"<?xml", "@page", ".", "#"}
var validIdTypes = []string{"dedication", "title", "foreword", "preface", "ch", "id"}
var filteredTitles = []string{"contents", "table of contents"}

func (c EpubChapter) IsValid() bool {
	id := strings.ToLower(c.Id)
	title := strings.TrimSpace(strings.ToLower(c.Title))

	return !hasExcludedIdType(id) &&
		!hasInvalidTitlePrefix(title) &&
		hasValidIdType(id) &&
		!isFilteredTitle(title)
}

// Validation methods
func matchesAny(slice []string, value string, matchFunc func(string, string) bool) bool {
	return slices.ContainsFunc(slice, func(item string) bool {
		return matchFunc(value, item)
	})
}

func hasExcludedIdType(id string) bool {
	return matchesAny(excludedTypes, id, strings.Contains)
}

func hasInvalidTitlePrefix(title string) bool {
	return matchesAny(invalidTitlePrefixes, title, strings.HasPrefix)
}

func hasValidIdType(id string) bool {
	return matchesAny(validIdTypes, id, func(value, pattern string) bool {
		return strings.Contains(value, pattern) || strings.HasPrefix(value, pattern)
	})
}

func isFilteredTitle(title string) bool {
	return matchesAny(filteredTitles, title, func(value, pattern string) bool {
		return value == pattern
	})
}
