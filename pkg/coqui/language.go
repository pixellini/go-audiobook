package coqui

import "strings"

// Language represents a supported language for TTS synthesis.
// Uses ISO 639-1 two-letter language codes (e.g., "en", "es", "fr").
// Coqui TTS doesn't take into account regional variations (e.g., "en-US" vs "en-GB").
// That will come from the speaker file supplied to the model.
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

// supportedLanguages contains all languages supported by Coqui TTS.
// Note: Language support varies by model; check model documentation.
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

// String returns the ISO 639-1 language code as a string.
func (l Language) String() string {
	return string(l)
}

// IsValid checks if the language is supported by Coqui TTS.
// Returns true for all languages in the supportedLanguages list.
func (l Language) IsValid() bool {
	for _, lang := range supportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// ParseLanguage parses a language string and returns the corresponding Language.
// Accepts formats like "en-US", "en", "es-ES" and extracts the language code.
// Returns English as the default for unsupported or invalid languages.
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

// MustParseLanguage parses a language string and panics if invalid.
// Use this when you need to ensure the language is valid at initialization time.
func MustParseLanguage(s string) Language {
	lang := ParseLanguage(s)
	if !lang.IsValid() {
		panic("Coqui TTS does not support language: " + s)
	}
	return lang
}

// All returns a copy of all supported languages.
// Useful for displaying available language options to users.
func All() []Language {
	return append([]Language(nil), supportedLanguages...)
}
