package epub

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/pixellini/go-audiobook/internal/formatter"
)

type Chapter struct {
	Title      string
	Paragraphs []string
	RawText    string
}

var chapterHeaderRegex = regexp.MustCompile(`(?i)^(chapter|part|section|page)\s*\d*\.?\s*`)

func createChapter(rawHtmlContent string) Chapter {
	cleanedContent := formatter.SplitText(cleanContent(rawHtmlContent))

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

		contentWithoutUrls := formatter.ReplaceUrlPattern(cleanedContent[i])
		cleanedContent[i] = formatter.ReplaceCodePattern(contentWithoutUrls)
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
		return matches[1]
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

func cleanContent(htmlContent string) string {
	htmlContent = regexp.MustCompile(`(?is)<head.*?>.*?</head>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<\?xml.*?\?>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<\?!DOCTYPE.*?\?>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<code.*?>.*?</code>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<pre.*?>.*?</pre>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<kbd.*?>.*?</kbd>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<samp.*?>.*?</samp>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<var.*?>.*?</var>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<figure.*?>.*?</figure>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<script.*?>.*?</script>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<style.*?>.*?</style>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<nav.*?>.*?</nav>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<aside.*?>.*?</aside>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<footer.*?>.*?</footer>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<form.*?>.*?</form>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<svg.*?>.*?</svg>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<object.*?>.*?</object>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<embed.*?>.*?</embed>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<iframe.*?>.*?</iframe>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<table.*?>.*?</table>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<math.*?>.*?</math>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<audio.*?>.*?</audio>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<video.*?>.*?</video>`).ReplaceAllString(htmlContent, "")
	htmlContent = strings.TrimSpace(htmlContent)
	return strip.StripTags(htmlContent)
}
