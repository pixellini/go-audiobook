package formatter

import (
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

func CleanContent(htmlContent string) string {
	htmlContent = regexp.MustCompile(`(?is)<head.*?>.*?</head>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<\?xml.*?\?>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<\?!DOCTYPE.*?\?>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<code.*?>.*?</code>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<pre.*?>.*?</pre>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<kbd.*?>.*?</kbd>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<samp.*?>.*?</samp>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<var.*?>.*?</var>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?is)<figure.*?>.*?</figure>`).ReplaceAllString(htmlContent, "")
	htmlContent = strings.TrimSpace(htmlContent)
	return strip.StripTags(htmlContent)
}
