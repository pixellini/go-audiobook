package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// CreateDirIfNotExist creates a directory if it does not exist.
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// FileExists returns true if the file exists.
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// RemoveDirIfEmpty removes the directory if it is empty.
func RemoveDirIfEmpty(dir string) error {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return os.Remove(dir)
	}

	return nil
}

// RemoveFilesInDir removes all files and subdirectories in the specified directory.
func RemoveAllFilesInDir(dir string) error {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dir, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			return err
		}
	}

	return nil
}

func RemoveFilesFrom(dir string, files []string) error {
	for _, chapterFile := range files {
		fullPath := filepath.Join(dir, chapterFile)
		fmt.Println(fullPath)
		err := os.Remove(fullPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// SortFilesNumerically sorts file paths by the number in their base name (e.g., part-1.wav).
func SortFilesNumerically(filePaths []string) {
	sort.Slice(filePaths, func(i, j int) bool {
		var a, b int

		_, err1 := fmt.Sscanf(filepath.Base(filePaths[i]), "part-%d.wav", &a)
		_, err2 := fmt.Sscanf(filepath.Base(filePaths[j]), "part-%d.wav", &b)

		if err1 != nil || err2 != nil {
			return filePaths[i] < filePaths[j]
		}

		return a < b
	})
}

// GetFilesFrom returns a list of files with the given extension in the specified directory.
func GetFilesFrom(dir, fileType string) ([]string, error) {
	if fileType == "" {
		return nil, fmt.Errorf("file type required")
	}

	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == "."+fileType {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

// CreateTempFileListTextFile creates a temporary file listing the given files, one per line, for use with ffmpeg.
func CreateTempFileListTextFile(files []string, fileName string) (string, error) {
	listFile, err := os.CreateTemp("", fileName)

	if err != nil {
		return "", fmt.Errorf("failed to create temporary list file: %w", err)
	}

	defer func() {
		if err != nil {
			listFile.Close()
			os.Remove(listFile.Name())
		}
	}()

	for _, file := range files {
		absPath, err := filepath.Abs(file)

		if err != nil {
			return "", fmt.Errorf("path resolution failed: %w", err)
		}

		if _, err := fmt.Fprintf(listFile, "file '%s'\n", absPath); err != nil {
			return "", fmt.Errorf("list file write failed: %w", err)
		}
	}

	if err := listFile.Sync(); err != nil {
		return "", fmt.Errorf("file sync failed: %w", err)
	}

	if err := listFile.Close(); err != nil {
		return "", fmt.Errorf("file close failed: %w", err)
	}

	return listFile.Name(), nil
}
