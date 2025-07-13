package audiobook

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pixellini/go-audiobook/internal/fsutils"
	"github.com/pixellini/go-audiobook/internal/utils"
	"github.com/pixellini/go-audiobook/pkg/epub"
	"github.com/spf13/viper"
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

type FileOptions struct {
	Name       string
	Image      string
	Metadata   string
	OutputName string
}

type OutputFileFormat string

const (
	OutputFormatM4B     OutputFileFormat = "m4b"
	OutputFormatMP3     OutputFileFormat = "mp3"
	OutputFormatAAC     OutputFileFormat = "m4a"
	OutputFormatWAV     OutputFileFormat = "wav"
	DefaultOutputFormat OutputFileFormat = OutputFormatM4B
)

func New(title, author, metadata, image string) *Audiobook {
	return &Audiobook{
		Title:    title,
		Author:   author,
		Metadata: metadata,
		Image:    image,
	}
}

func NewWithEPUB(e *epub.Epub, image string) *Audiobook {
	return New(e.Title, e.Author, e.Description, image)
}

func (a *Audiobook) AddChapter(chapter AudiobookChapter) {
	a.Chapters = append(a.Chapters, chapter)
}

func (a *Audiobook) Generate(outputDir string) error {
	outputFormat := OutputFileFormat(viper.GetString("output_format"))
	if outputFormat == "" {
		fmt.Println("No format specified, defaulting to M4B")
		outputFormat = DefaultOutputFormat
	}

	chapterWavFiles := make([]string, len(a.Chapters))
	for i, chapter := range a.Chapters {
		chapterWavFiles[i] = chapter.File
	}
	tempWavFileOutput := fmt.Sprintf("%s/concatenated.wav", outputDir)

	err := ConcatFiles("concat", chapterWavFiles, tempWavFileOutput)
	if err != nil {
		os.Remove(tempWavFileOutput)
		fmt.Println("Error concatenating audio files:", err)
		return err
	}

	outputName := fmt.Sprintf("%s/%s - %s", outputDir, a.Title, a.Author)

	metadataFile, err := a.buildMetadataFile()
	if err != nil {
		os.Remove(tempWavFileOutput)
		return err
	}
	defer metadataFile.Close()

	err = CreateFileFromFormat(outputFormat, FileOptions{
		Name:       tempWavFileOutput,
		Image:      a.Image,
		Metadata:   metadataFile.Name(),
		OutputName: outputName,
	})

	if err != nil {
		os.Remove(tempWavFileOutput)
		fmt.Printf("Error creating %s file: %v\n", outputFormat, err)
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
		duration, err := GetDuration(chapterFile.File)
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

func ConcatFiles(title string, files []string, output string) error {
	// Create the temporary file list
	listFileName, err := fsutils.CreateTempFileListTextFile(files, "ffmpeg_concat_*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(listFileName)

	cmd := exec.Command(
		"ffmpeg", "-y",
		"-f", "concat", "-safe", "0",
		"-i", listFileName,
		"-metadata", fmt.Sprintf("title=%s", title),
		"-c", "copy",
		output,
	)

	utils.LogCommandIfVerbose(cmd)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to concatenate audio files: %v", err)
	}
	return nil
}

func GetDuration(audioFilePath string) (float64, error) {
	probeResult, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", audioFilePath).Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(probeResult))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func CreateFile(fileOptions FileOptions, extension OutputFileFormat, additionalOptions []string) error {
	outputFile := fmt.Sprintf("%s.%s", fileOptions.OutputName, extension)

	baseFlags := []string{
		"-y",
		"-i", fileOptions.Name,
	}

	// Cover art image for non-wav files
	if extension != OutputFormatWAV {
		baseFlags = append(baseFlags, "-f", "image2", "-i", fileOptions.Image)
	}

	// Combine base flags, extra flags, and output file
	cmdArgs := append(baseFlags, additionalOptions...)
	cmdArgs = append(cmdArgs, outputFile)

	cmd := exec.Command("ffmpeg", cmdArgs...)

	utils.LogCommandIfVerbose(cmd)

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to create %s file: %v", extension, err)
		return err
	}

	return nil
}

func CreateM4BFile(fileOptions FileOptions) error {
	return CreateFile(fileOptions, OutputFormatM4B, []string{
		"-f", "ffmetadata", "-i", fileOptions.Metadata,
		"-map", "0:a",
		"-map", "1",
		"-map_metadata", "2",
		"-id3v2_version", "3",
		"-write_id3v1", "1",
		"-c:a", "aac",
		"-b:a", "64k",
		"-c:v", "mjpeg",
		"-disposition:v:0", "attached_pic",
		"-f", "ipod",
	})
}

func CreateMP3File(fileOptions FileOptions) error {
	return CreateFile(fileOptions, OutputFormatMP3, []string{
		"-map", "0:a",
		"-map", "1",
		"-id3v2_version", "3",
		"-write_id3v1", "1",
		"-c:a", "libmp3lame",
		"-b:a", "128k",
		"-c:v", "mjpeg",
		"-disposition:v:0", "attached_pic",
	})
}

func CreateAACFile(fileOptions FileOptions) error {
	return CreateFile(fileOptions, OutputFormatAAC, []string{
		"-map", "0:a",
		"-map", "1",
		"-c:a", "aac",
		"-b:a", "96k",
		"-c:v", "mjpeg",
		"-disposition:v:0", "attached_pic",
	})
}

func CreateWAVFile(fileOptions FileOptions) error {
	return CreateFile(fileOptions, OutputFormatWAV, []string{
		"-c:a", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
	})
}

func CreateFileFromFormat(format OutputFileFormat, fileOptions FileOptions) error {
	switch format {
	case OutputFormatMP3:
		fmt.Println("Creating MP3 file...")
		return CreateMP3File(fileOptions)
	case OutputFormatAAC:
		fmt.Println("Creating AAC file...")
		return CreateAACFile(fileOptions)
	case OutputFormatWAV:
		fmt.Println("Creating WAV file...")
		return CreateWAVFile(fileOptions)
	case OutputFormatM4B:
		fmt.Println("Creating M4B file...")
		return CreateM4BFile(fileOptions)
	default:
		fmt.Printf("Format '%s' is not supported. Defaulting to M4B\n", format)
		return CreateM4BFile(fileOptions)
	}
}
