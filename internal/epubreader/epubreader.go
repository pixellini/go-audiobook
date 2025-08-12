package epubreader

import (
	"fmt"
	"io"

	"github.com/pixellini/go-audiobook/internal/textutils"
	epubReader "github.com/taylorskalyo/goreader/epub"
)

type EpubReader interface {
	GetTitle() string
	GetAuthor() string
	GetDescription() string
	GetLanguage() string
	GetPublisher() string
	GetCoverImage() string
	GetChapter(index int) (*EpubReaderChapter, error)
	GetChapters() ([]*EpubReaderChapter, error)
	Close() error
}

type EpubReaderChapter struct {
	Id string
	// Title is the title of the chapter.
	Title string
	// Content is the raw, unedited HTML chapter content.
	Content string

	Path string
}

type GoEpubReaderService struct {
	epubFile *epubReader.ReadCloser
	r        *epubReader.Rootfile
	path     string
}

func NewGoEpubReaderService(path string) (EpubReader, error) {
	epubFile, err := epubReader.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open EPUB: %w", err)
	}

	if len(epubFile.Rootfiles) == 0 {
		epubFile.Close()
		return nil, fmt.Errorf("no rootfiles found in EPUB")
	}

	return &GoEpubReaderService{
		epubFile: epubFile,
		r:        epubFile.Rootfiles[0],
		path:     path,
	}, nil
}

func (g *GoEpubReaderService) GetTitle() string {
	return g.r.Title
}

func (g *GoEpubReaderService) GetAuthor() string {
	return g.r.Creator
}

func (g *GoEpubReaderService) GetDescription() string {
	return g.r.Description
}

func (g *GoEpubReaderService) GetLanguage() string {
	return g.r.Language
}

func (g *GoEpubReaderService) GetPublisher() string {
	return g.r.Publisher
}

func (g *GoEpubReaderService) GetCoverImage() string {
	return ""
}

func (g *GoEpubReaderService) GetChapter(index int) (*EpubReaderChapter, error) {
	item := g.r.Manifest.Items[index]
	r, err := item.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open item %s: %w", item.ID, err)
	}

	content, err := io.ReadAll(r)

	r.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read item %s: %w", item.ID, err)
	}

	contentStr := string(content)
	title := textutils.ExtractTitleFromHTML(contentStr)

	return &EpubReaderChapter{
		Id:      item.ID,
		Title:   title,
		Content: contentStr,
	}, nil
}

func (g *GoEpubReaderService) GetChapters() ([]*EpubReaderChapter, error) {
	chapterCount := len(g.r.Manifest.Items)
	chapters := make([]*EpubReaderChapter, 0, chapterCount)

	for i := range chapterCount {
		chapter, err := g.GetChapter(i)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (g *GoEpubReaderService) Close() error {
	if g.epubFile != nil {
		g.epubFile.Close()
	}
	return nil
}
