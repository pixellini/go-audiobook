package filemanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileService interface {
	Create(dir string)
	Remove(dir string) error
	RemoveFiles(paths []string) error
	Save(path, content string) error
	CreateCacheDir() (cacheDir string, err error)
}

type FileManager struct{}

const defaultCacheDir = "./.cache/"

func New() FileService {
	return &FileManager{}
}

func (f *FileManager) Create(dir string) {
	if dir == "" {
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic("Error creating directory: " + err.Error())
		}
	}
}

func (f *FileManager) Remove(dir string) error {
	if dir == "" {
		return nil // Nothing to remove
	}

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil

}

func (f *FileManager) RemoveFiles(paths []string) error {
	var errs []error

	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			errs = append(errs, fmt.Errorf("\nunable to remove %s: %w", p, err))
		}
	}

	return errors.Join(errs...)
}

func (f *FileManager) Save(path, content string) error {
	if path == "" {
		return nil // Nothing to save
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileManager) CreateCacheDir() (cacheDir string, err error) {
	// Try system temp directory first
	sysCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	sysCachePath := filepath.Join(sysCacheDir, "go-audiobook") + "/"

	if err := f.createAndTestDir(sysCachePath); err == nil {
		return sysCachePath, nil
	}

	// Last resort: fallback to local temp directory
	if err := f.createAndTestDir(defaultCacheDir); err != nil {
		return "", fmt.Errorf("failed to create temp directory in system temp, unique temp, or local fallback: %w", err)
	}

	return defaultCacheDir, nil
}

func (f *FileManager) createAndTestDir(path string) error {
	// Try to create the directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// Test if we can actually write to it
	testFile := filepath.Join(path, "test.tmp")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()

	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		return err
	}

	return nil
}
