package epub

import (
	"html"
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

// unwantedHTMLTags defines HTML elements that should be removed for TTS processing
var unwantedHTMLTags = []string{
	"code", "pre", "kbd", "samp", "var", "figure", "script", "style",
	"nav", "aside", "footer", "form", "svg", "object", "embed",
	"iframe", "table", "math", "audio", "video",
}

// HTMLCleaner handles cleaning HTML content for EPUB text extraction
type HTMLCleaner struct {
	// Pre-compiled regexes for better performance
	doctypeRegex    *regexp.Regexp
	xmlDeclRegex    *regexp.Regexp
	allTagsRegex    *regexp.Regexp
	whitespaceRegex *regexp.Regexp

	// Content-specific regexes
	headRegex      *regexp.Regexp
	htmlTagRegex   *regexp.Regexp
	htmlCloseRegex *regexp.Regexp
	bodyTagRegex   *regexp.Regexp
	bodyCloseRegex *regexp.Regexp

	// Unwanted element regexes
	unwantedElements map[string]*regexp.Regexp
}

// NewHTMLCleaner creates a new HTML cleaner with pre-compiled regexes
func NewHTMLCleaner() *HTMLCleaner {
	cleaner := &HTMLCleaner{
		// Document structure
		doctypeRegex:    regexp.MustCompile(`(?s)<!DOCTYPE[^>]*>`),
		xmlDeclRegex:    regexp.MustCompile(`(?s)<\?xml[^>]*\?>`),
		allTagsRegex:    regexp.MustCompile(`<[^>]+>`),
		whitespaceRegex: regexp.MustCompile(`\s+`),

		// HTML structure
		headRegex:      regexp.MustCompile(`(?is)<head.*?>.*?</head>`),
		htmlTagRegex:   regexp.MustCompile(`(?s)<html[^>]*>`),
		htmlCloseRegex: regexp.MustCompile(`(?s)</html>`),
		bodyTagRegex:   regexp.MustCompile(`(?s)<body[^>]*>`),
		bodyCloseRegex: regexp.MustCompile(`(?s)</body>`),

		unwantedElements: make(map[string]*regexp.Regexp),
	}

	// Pre-compile unwanted element regexes
	for _, tag := range unwantedHTMLTags {
		pattern := `(?is)<` + tag + `.*?>.*?</` + tag + `>`
		cleaner.unwantedElements[tag] = regexp.MustCompile(pattern)
	}

	return cleaner
}

// CleanHTML removes all HTML tags and normalizes whitespace
func (c *HTMLCleaner) CleanHTML(content string) string {
	// Remove document declarations
	content = c.doctypeRegex.ReplaceAllString(content, "")
	content = c.xmlDeclRegex.ReplaceAllString(content, "")
	// Remove all HTML tags
	content = c.allTagsRegex.ReplaceAllString(content, "")
	// Unescape HTML entities
	content = html.UnescapeString(content)
	// Normalize whitespace
	content = c.whitespaceRegex.ReplaceAllString(content, " ")

	return strings.TrimSpace(content)
}

// CleanContent removes unwanted HTML elements while preserving text content
func (c *HTMLCleaner) CleanContent(htmlContent string) string {
	// Remove document structure
	htmlContent = c.headRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.xmlDeclRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.doctypeRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.htmlTagRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.htmlCloseRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.bodyTagRegex.ReplaceAllString(htmlContent, "")
	htmlContent = c.bodyCloseRegex.ReplaceAllString(htmlContent, "")

	// Remove unwanted elements
	for _, regex := range c.unwantedElements {
		htmlContent = regex.ReplaceAllString(htmlContent, "")
	}

	htmlContent = strings.TrimSpace(htmlContent)
	return strip.StripTags(htmlContent)
}

// ContentProcessor handles text pattern replacements for TTS optimization
type ContentProcessor struct {
	codePatternRegex *regexp.Regexp
	urlRegex         *regexp.Regexp
}

// NewContentProcessor creates a new content processor
func NewContentProcessor() *ContentProcessor {
	return &ContentProcessor{
		codePatternRegex: regexp.MustCompile(`(?m)^(Figure|Fig|Listing|Example)\b.*$`),
		urlRegex:         regexp.MustCompile(`(https?://\S+|www\.\S+)`),
	}
}

// ReplaceCodePatterns replaces code-related patterns with TTS-friendly alternatives
func (p *ContentProcessor) ReplaceCodePatterns(content string) string {
	return p.codePatternRegex.ReplaceAllString(content, "...")
}

// ReplaceURLs replaces URLs with TTS-friendly text
func (p *ContentProcessor) ReplaceURLs(content string) string {
	if p.urlRegex.MatchString(content) {
		return p.urlRegex.ReplaceAllString(content, "URL Link")
	}
	return content
}

// ProcessForTTS applies all TTS-friendly transformations
func (p *ContentProcessor) ProcessForTTS(content string) string {
	content = p.ReplaceCodePatterns(content)
	content = p.ReplaceURLs(content)
	return content
}
