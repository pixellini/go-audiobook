package metadata

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/pixellini/go-audiobook/internal/epub"
)

type Metadata struct {
	bw *bufio.Writer
	f  *os.File
}

func New(dir string) (*Metadata, error) {
	file, err := os.CreateTemp(dir, "metadata-*.txt")
	if err != nil {
		return nil, err
	}

	bw := bufio.NewWriter(file)

	if _, err := bw.WriteString(";FFMETADATA1\n"); err != nil {
		file.Close()
		return nil, err
	}

	return &Metadata{
		bw: bw,
		f:  file,
	}, nil
}

func (m *Metadata) AddDetails(md *epub.EpubMetadata) error {
	var errs []error

	if err := m.writeProperty("title", md.Title); err != nil {
		errs = append(errs, err)
	}

	if err := m.writeProperty("artist", md.Author); err != nil {
		errs = append(errs, err)
	}

	if err := m.writeProperty("description", md.Description); err != nil {
		errs = append(errs, err)
	}

	if err := m.writeProperty("publisher", md.Publisher); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (m *Metadata) AddChapter(title string, start, end int) error {
	_, err := fmt.Fprintf(m.bw,
		"[CHAPTER]\nTIMEBASE=1/1000\nSTART=%d\nEND=%d\ntitle=%s\n\n",
		start, end, title)
	return err
}

func (m *Metadata) writeProperty(key, val string) error {
	if val == "" {
		return nil
	}
	_, err := fmt.Fprintf(m.bw, "%s=%s\n", key, val)
	return err
}

func (m *Metadata) Name() string {
	return m.f.Name()
}

func (m *Metadata) Close() error {
	if err := m.bw.Flush(); err != nil {
		m.f.Close()
		return err
	}
	return m.f.Close()
}

func (m *Metadata) Remove() error {
	// return os.Remove(m.Name())
	return nil
}
