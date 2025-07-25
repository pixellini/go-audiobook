package audioservice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type AudioService interface {
	CombineFiles(inputFiles []string, outputFile string) error
	ConvertFile(inputFile, outputFile string) error
	GetDuration(audioFilePath string) (float64, error)
	CreateAudiobook(file, image, metadataPath, output string) error
}

type FFMpegService struct {
	outputDir string
}

func NewFFMpegService(outputDir string) *FFMpegService {
	return &FFMpegService{
		outputDir: outputDir,
	}
}

func (f *FFMpegService) CombineFiles(inputFiles []string, outputFile string) error {
	// fmt.Println("Combining audio files:", len(inputFiles), "files into", outputFile)
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided for combination")
	}

	// Verify all input files exist
	for _, inputFile := range inputFiles {
		if _, err := os.Stat(inputFile); err != nil {
			return fmt.Errorf("input file does not exist: %s - %w", inputFile, err)
		}
	}

	// Create a temporary file list for FFmpeg concat
	file, err := os.CreateTemp(f.outputDir, "filelist-*.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name()) // Clean up the temporary file

	// Write file list to temporary file
	err = f.writeFileList(inputFiles, file.Name())
	if err != nil {
		return fmt.Errorf("failed to create file list: %v", err)
	}

	// Use FFmpeg concat demuxer
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", file.Name(),
		"-c", "copy",
		"-y", // Overwrite output file if it exists
		outputFile,
	}

	err = f.ffmpeg(args...)
	if err != nil {
		return fmt.Errorf("failed to concatenate audio files: %v", err)
	}

	return nil
}

// writeFileList creates a file list for FFmpeg concat demuxer
func (f *FFMpegService) writeFileList(inputFiles []string, listFile string) error {
	var fileListContent strings.Builder
	for _, fpath := range inputFiles {
		// Convert to absolute path to avoid issues
		absPath, err := filepath.Abs(fpath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %v", fpath, err)
		}
		// Escape single quotes in the path and wrap in single quotes
		escapedPath := strings.ReplaceAll(absPath, "'", "'\"'\"'")
		fileListContent.WriteString(fmt.Sprintf("file '%s'\n", escapedPath))
	}

	return os.WriteFile(listFile, []byte(fileListContent.String()), 0644)
}

func (f *FFMpegService) ConvertFile(inputFile, outputFile string) error {
	return nil
}

func (f *FFMpegService) GetDuration(audioFilePath string) (float64, error) {
	probeResult, err := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		audioFilePath,
	).Output()
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

func (f *FFMpegService) CreateAudiobook(file, image, metadataPath, output string) error {
	return f.ffmpeg(
		"-i", file,
		"-i", metadataPath,
		"-i", image,
		"-map", "0:a",
		"-map", "2:v",
		"-map_metadata", "1",
		"-c:a", "aac",
		"-b:a", "64k",
		"-c:v", "png",
		"-disposition:v:0", "attached_pic",
		"-f", "ipod",
		output,
	)
}

func (f *FFMpegService) ffmpeg(args ...string) error {
	return f.ffmpegContext(context.Background(), args...)
}

func (f *FFMpegService) ffmpegContext(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v\nStderr: %s", err, stderr.String())
	}
	return nil
}
