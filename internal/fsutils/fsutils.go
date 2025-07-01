package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

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
