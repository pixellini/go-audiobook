package audiobook

import (
	"fmt"
	"os"

	"github.com/pixellini/go-audiobook/internal/audioprocessor"
	"github.com/pixellini/go-audiobook/internal/epub"
)

type Audiobook struct {
	Title    string
	Author   string
	Metadata string
	Image    string
	Chapters []AudiobookChapter
}

type AudiobookChapter struct {
	Title string
	File  string
}

func New(title, author, metadata, image string) *Audiobook {
	return &Audiobook{
		Title:    title,
		Author:   author,
		Metadata: metadata,
		Image:    image,
	}
}

func NewFromEpub(e *epub.Epub, image string) *Audiobook {
	return New(e.Title, e.Author, e.Description, image)
}

func (a *Audiobook) AddChapter(chapter AudiobookChapter) {
	a.Chapters = append(a.Chapters, chapter)
}

func (a *Audiobook) Generate(outputDir string) error {
	metadataFile, err := a.buildMetadataFile()
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	chapterWavFiles := make([]string, len(a.Chapters))
	for i, chapter := range a.Chapters {
		chapterWavFiles[i] = chapter.File
	}

	tempWavFileOutput := fmt.Sprintf("%s/concatenated.wav", outputDir)

	err = audioprocessor.ConcatFiles("concat", chapterWavFiles, tempWavFileOutput)
	if err != nil {
		os.Remove(tempWavFileOutput)
		fmt.Println("Error concatenating audio files:", err)
		return err
	}

	outputName := fmt.Sprintf("%s/%s - %s", outputDir, a.Title, a.Author)

	err = audioprocessor.CreateM4BFile(audioprocessor.FileOptions{
		Name:       tempWavFileOutput,
		Image:      a.Image,
		Metadata:   metadataFile.Name(),
		OutputName: outputName,
	})
	if err != nil {
		os.Remove(tempWavFileOutput)
		fmt.Println("Error creating m4b file:", err)
		return err
	}

	os.Remove(tempWavFileOutput)
	a.cleanUp()

	return nil
}

func (a *Audiobook) cleanUp() {
	// Remove the chapter files.
	for _, chapter := range a.Chapters {
		err := os.Remove(chapter.File)

		if err != nil {
			fmt.Println("Error removing chapter file:", err)
		}
	}
}

func (a *Audiobook) buildMetadataFile() (*os.File, error) {
	file, err := os.CreateTemp("", "chapters.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file.WriteString(";FFMETADATA1\n")

	startTime := 0
	for _, chapterFile := range a.Chapters {
		title := chapterFile.Title

		// Get the chapter duration
		// Then append it to the startTime so that we can calculate the endTime for the next chapter.
		duration, err := audioprocessor.GetDuration(chapterFile.File)
		if err != nil {
			return nil, err
		}

		durationMs := int(duration * 1000)
		endTime := startTime + durationMs

		file.WriteString("[CHAPTER]\n")
		file.WriteString("TIMEBASE=1/1000\n")
		file.WriteString(fmt.Sprintf("START=%d\n", startTime))
		file.WriteString(fmt.Sprintf("END=%d\n", endTime))
		file.WriteString(fmt.Sprintf("title=%s\n\n", title))

		startTime = endTime
	}

	return file, nil
}
