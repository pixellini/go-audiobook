package epub

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Chapter struct {
	Title      string
	Paragraphs []string
	RawText    string
}

var nonNumberedChapterTitles = []string{
	"chapter",
	"part",
	"introduction",
	"acknowledgements",
	"prologue",
	"epilogue",
	"foreword",
	"preface",
	"afterword",
	"conclusion",
	"appendix",
	"glossary",
	"index",
}

var chapterHeaderRegex = regexp.MustCompile(`(?i)^(chapter|part|section|page)\s*\d*\.?\s*`)

func CreateChapterAnnouncement(chapterIndex int, title string) string {
	lowerTitle := strings.ToLower(strings.TrimSpace(title))

	for _, prefix := range nonNumberedChapterTitles {
		if strings.HasPrefix(lowerTitle, prefix) {
			return title
		}
	}

	// E.g. "Chapter 1: The Beginning"
	return fmt.Sprintf("Chapter %d: %s", chapterIndex, title)
}

func createChapter(rawHtmlContent string) Chapter {
	cleanedContent := SplitText(CleanContent(rawHtmlContent))

	title := getChapterTitle(rawHtmlContent)
	if title == "" || title == "Converted Ebook" {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHtmlContent))
		if err == nil {
			heading := ""
			doc.Find("h1, h2, h3, h4, h5, h6").EachWithBreak(func(i int, s *goquery.Selection) bool {
				heading = strings.TrimSpace(s.Text())
				return heading == ""
			})
			if heading != "" {
				title = heading
			}
		}
	}

	if title == "" || title == "Converted Ebook" {
		for _, line := range cleanedContent {
			if strings.TrimSpace(line) != "" {
				title = strings.TrimSpace(line)
				break
			}
		}
	}

	for i, line := range cleanedContent {
		if chapterHeaderRegex.MatchString(line) {
			cleanedContent[i] = chapterHeaderRegex.ReplaceAllString(line, "")
		}

		contentWithoutUrls := ReplaceUrlPattern(cleanedContent[i])
		cleanedContent[i] = ReplaceCodePattern(contentWithoutUrls)
	}

	return Chapter{
		Title:      title,
		Paragraphs: cleanedContent,
		RawText:    rawHtmlContent,
	}
}

func getChapterTitle(htmlContent string) string {
	titleTagRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := titleTagRegex.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := matches[1]
		// Reject titles that look like filenames (e.g., "ch001.xhtml")
		// This seems to happen a lot when the file is converted to EPUB.
		if strings.HasSuffix(title, ".xhtml") || strings.HasSuffix(title, ".html") {
			return ""
		}
		return title
	}
	return ""
}

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
