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

var languageExceptions = map[string]string{
	"zh-Hans": "zh-Hans", // Simplified Chinese
	"zh-Hant": "zh-Hant", // Traditional Chinese
}

var languagePrefixMap = map[string]string{
	"en-": "en", // English
	"es-": "es", // Spanish
	"fr-": "fr", // French
	"de-": "de", // German
	"it-": "it", // Italian
	"pt-": "pt", // Portuguese
	"pl-": "pl", // Polish
	"tr-": "tr", // Turkish
	"ru-": "ru", // Russian
	"nl-": "nl", // Dutch
	"cs-": "cs", // Czech
	"ar-": "ar", // Arabic
	"zh-": "zh", // Chinese
	"ja-": "ja", // Japanese
	"hu-": "hu", // Hungarian
	"ko-": "ko", // Korean
}

func FormatToStandardLanguage(language string) string {
	if std, ok := languageExceptions[language]; ok {
		return std
	}
	for prefix, std := range languagePrefixMap {
		if strings.HasPrefix(language, prefix) {
			return std
		}
	}
	return language
}
