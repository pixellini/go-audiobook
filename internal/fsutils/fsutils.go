package fsutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false // some other I/O error occurred
}

// This should be a bit more generic
func SortNumerically(s []string) {
	sort.Slice(s, func(i, j int) bool {
		var a, b int

		_, err1 := fmt.Sscanf(filepath.Base(s[i]), "paragraph-%d.wav", &a)
		_, err2 := fmt.Sscanf(filepath.Base(s[j]), "paragraph-%d.wav", &b)

		if err1 != nil || err2 != nil {
			return s[i] < s[j]
		}

		return a < b
	})
}
