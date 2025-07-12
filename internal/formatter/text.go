package formatter

import (
	"html"
	"regexp"
	"strings"
)

func SplitText(content string) []string {
	content = CleanHTML(content)

	reSingleLineBreak := regexp.MustCompile(`([^\n])\n([^\n])`)
	for reSingleLineBreak.MatchString(content) {
		content = reSingleLineBreak.ReplaceAllString(content, "$1 $2")
	}

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

func CleanHTML(content string) string {
	content = regexp.MustCompile(`(?s)<!DOCTYPE[^>]*>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`(?s)<\?xml[^>]*\?>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")
	content = html.UnescapeString(content)
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	return strings.TrimSpace(content)
}
