package coqui

import "strings"

type Language string

const (
	English    Language = "en"
	Spanish    Language = "es"
	French     Language = "fr"
	German     Language = "de"
	Italian    Language = "it"
	Portuguese Language = "pt"
	Polish     Language = "pl"
	Turkish    Language = "tr"
	Russian    Language = "ru"
	Dutch      Language = "nl"
	Czech      Language = "cs"
	Arabic     Language = "ar"
	Chinese    Language = "zh"
	Japanese   Language = "ja"
	Hungarian  Language = "hu"
	Korean     Language = "ko"
)

// All supported languages
var supportedLanguages = []Language{
	English,
	Spanish,
	French,
	German,
	Italian,
	Portuguese,
	Polish,
	Turkish,
	Russian,
	Dutch,
	Czech,
	Arabic,
	Chinese,
	Japanese,
	Hungarian,
	Korean,
}

// String returns the string representation of the language
func (l Language) String() string {
	return string(l)
}

// IsValid checks if the language is supported
func (l Language) IsValid() bool {
	for _, lang := range supportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// ParseLanguage takes a language string (e.g., "en-US", "en", "es-ES")
// and returns the corresponding Language. Returns English for unknown languages.
func ParseLanguage(s string) Language {
	// Extract language code (before the "-")
	if idx := strings.Index(s, "-"); idx != -1 {
		s = s[:idx]
	}

	lang := Language(strings.ToLower(s))

	// Validate it's a supported language
	if lang.IsValid() {
		return lang
	}

	return English
}

// MustParseLanguage parses a language string and panics if invalid
func MustParseLanguage(s string) Language {
	lang := ParseLanguage(s)
	if !lang.IsValid() {
		panic("Coqui TTS does not support language: " + s)
	}
	return lang
}

// All returns all supported languages
func All() []Language {
	return append([]Language(nil), supportedLanguages...)
}
