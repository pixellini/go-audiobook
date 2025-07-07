package audioprocessor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pixellini/go-audiobook/internal/fsutils"
)

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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
		return CreateMP3File(fileOptions)
	case OutputFormatAAC:
		return CreateAACFile(fileOptions)
	case OutputFormatWAV:
		return CreateWAVFile(fileOptions)
	case OutputFormatM4B:
		fallthrough
	default:
		return CreateM4BFile(fileOptions)
	}
}
