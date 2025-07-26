package fsutils

import (
	"fmt"
	"path/filepath"
	"sort"
)

func SortNumerically(s []string) {
	sort.Slice(s, func(i, j int) bool {
		var a, b int

		_, err1 := fmt.Sscanf(filepath.Base(s[i]), "part-%d.wav", &a)
		_, err2 := fmt.Sscanf(filepath.Base(s[j]), "part-%d.wav", &b)

		if err1 != nil || err2 != nil {
			return s[i] < s[j]
		}

		return a < b
	})
}
