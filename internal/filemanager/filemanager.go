package filemanager

import (
	"os"
	"path/filepath"
)

type FileService interface {
	Create(dir string)
	Remove(dir string) error
	CreateTemp()
	RemoveTemp() error
	Save(path, content string) error
}

type FileManager struct {
	TempDir         string
	fallbackTempDir string
}

func NewService(fallbackTempDir string) (FileService, error) {
	return &FileManager{
		fallbackTempDir: fallbackTempDir,
	}, nil
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

func (f *FileManager) CreateTemp() {
	if f.TempDir != "" {
		return // Temp directory already exists
	}

	// Use the fallback temp directory if provided, otherwise use the system temp directory
	f.TempDir = filepath.Join(f.TempDir, "go-audiobook-temp")

	if _, err := os.Stat(f.TempDir); os.IsNotExist(err) {
		err = os.MkdirAll(f.TempDir, 0755)
		if err != nil {
			f.TempDir = f.fallbackTempDir
		}
	}
}

func (f *FileManager) RemoveTemp() error {
	if f.TempDir == "" {
		return nil // Nothing to remove
	}

	err := os.RemoveAll(f.TempDir)
	if err != nil {
		return err
	}

	f.TempDir = ""
	return nil
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

// func (f *FileManager) CreateTempDir(fallbackTempDir string) (string, error) {
// 	// Check if we can create a file in the system temp directory.
// 	sysTempDir := os.TempDir()
// 	sysAppTempDir := filepath.Join(sysTempDir, "go-audiobook")

// 	if _, err := os.Stat(sysAppTempDir); !os.IsNotExist(err) {
// 		return "", nil
// 	}

// 	err := os.MkdirAll(sysAppTempDir, 0755)
// 	if os.IsPermission(err) {
// 		fmt.Println("Warning: Permission denied for system temp directory:", sysAppTempDir)
// 	} else {
// 		fmt.Println("Error: Unable to create or access the system temp directory:", err)
// 	}

// 	if !f.CanCreate(sysAppTempDir) {
// 		return "", fmt.Errorf("Error: Unable to create or access the system temp directory.")
// 	}

// 	return sysAppTempDir, nil
// }

// func (f *FileManager) CanCreate(path string) bool {
// 	// Attempt to create a test file.
// 	testFile := filepath.Join(path, "test.tmp")
// 	file, err := os.Create(testFile)
// 	if err != nil {
// 		return false
// 	}

// 	file.Close()
// 	// Attempt to remove the test file.
// 	err = os.Remove(testFile)
// 	if err != nil {
// 		fmt.Println("Error: Unable to remove test file in system temp directory:", err)
// 		return false
// 	}
// 	return true
// }
