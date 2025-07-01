package formatter

import (
	"strings"
)

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
