package formatter

import (
	"regexp"
	"strings"
)

func SplitText(content string) []string {
	re := regexp.MustCompile(`(?m)([^\n]+?[\.!?](?:\s+|$)|[^\n]+$)`)
	matches := re.FindAllString(content, -1)
	var result []string
	for _, match := range matches {
		trimmed := strings.TrimSpace(match)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
