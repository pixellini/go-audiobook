package epub

import (
	"regexp"
	"strings"

	"github.com/pixellini/go-audiobook/formatter"
)

// Chapter represents a single chapter in the EPUB book.
type Chapter struct {
	Title      string
	Paragraphs []string // Split-up content
	RawText    string   // Original, unsplit content
}

var chapterHeaderRegex = regexp.MustCompile(`(?i)^(chapter|part|section|page)\s*\d*\.?\s*`)

func createChapter(rawHtmlContent string) Chapter {
	cleanedContent := formatter.SplitText(formatter.CleanContent(rawHtmlContent))

	for i, line := range cleanedContent {
		if chapterHeaderRegex.MatchString(line) {
			cleanedContent[i] = chapterHeaderRegex.ReplaceAllString(line, "")
		}

		contentWithoutUrls := formatter.ReplaceUrlPattern(cleanedContent[i])
		cleanedContent[i] = formatter.ReplaceCodePattern(contentWithoutUrls)
	}

	return Chapter{
		Title:      getChapterTitle(rawHtmlContent),
		Paragraphs: cleanedContent,
		RawText:    rawHtmlContent,
	}
}

func getChapterTitle(htmlContent string) string {
	titleTagRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := titleTagRegex.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// isAcceptedChapterItem determines if an item ID is a section to be included as a chapter.
// TODO: There may be other acceptable sections and it will need to be investigated.
func isAcceptedChapterItem(id string) bool {
	id = strings.ToLower(id)
	return strings.Contains(id, "title") ||
		strings.Contains(id, "dedication") ||
		strings.Contains(id, "foreword") ||
		strings.Contains(id, "preface") ||
		strings.HasPrefix(id, "ch") ||
		strings.HasPrefix(id, "id")
}
