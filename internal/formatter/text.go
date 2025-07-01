package formatter

import (
	"regexp"
	"strings"
)

func SplitText(content string) []string {
	// Step 1: Normalize line breaks (unwrap lines)
	// Replace single line breaks (not double) with a space
	reSingleLineBreak := regexp.MustCompile(`([^\n])\n([^\n])`)
	for reSingleLineBreak.MatchString(content) {
		content = reSingleLineBreak.ReplaceAllString(content, "$1 $2")
	}
	// Step 2: Split into sentences as before
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
