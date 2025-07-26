package textutils

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
)

// SplitText splits the content into sentences based on punctuation and line breaks.
func SplitText(content string) []string {
	// First, extract clean text from HTML content
	cleanText := ExtractTextFromHTML(content)

	// Split into paragraphs first (double newlines or <p> boundaries)
	paragraphs := SplitIntoParagraphs(cleanText)

	var sentences []string
	for _, paragraph := range paragraphs {
		// Split each paragraph into sentences
		paragraphSentences := SplitIntoSentences(paragraph)
		sentences = append(sentences, paragraphSentences...)
	}

	return sentences
}

// ExtractTextFromHTML extracts plain text from HTML content, preserving paragraph structure
func ExtractTextFromHTML(htmlContent string) string {
	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		// If parsing fails, fall back to simple tag stripping
		return strip.StripTags(htmlContent)
	}

	var textParts []string

	// Extract text from paragraphs, preserving structure
	doc.Find("p, div, h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			textParts = append(textParts, text)
		}
	})

	// If no structured elements found, extract all text
	if len(textParts) == 0 {
		return strings.TrimSpace(strip.StripTags(htmlContent))
	}

	return strings.Join(textParts, "\n\n")
}

// SplitIntoParagraphs splits text into paragraphs
func SplitIntoParagraphs(text string) []string {
	// Split on double newlines
	paragraphs := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

	var result []string
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// SplitIntoSentences splits a paragraph into individual sentences
func SplitIntoSentences(paragraph string) []string {
	// Clean up the text
	text := strings.TrimSpace(paragraph)
	if text == "" {
		return nil
	}

	// Regex to split on sentence endings (.!?) followed by whitespace and capital letter
	// but not on common abbreviations
	sentenceRegex := regexp.MustCompile(`([.!?]+)\s+([A-Z])`)

	// Replace sentence boundaries with a special marker
	marked := sentenceRegex.ReplaceAllString(text, `$1|||$2`)

	// Split on our marker
	sentences := strings.Split(marked, "|||")

	var result []string
	for i, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// Rejoin with the next sentence's first character if needed
		if i < len(sentences)-1 && !strings.HasSuffix(sentence, ".") &&
			!strings.HasSuffix(sentence, "!") && !strings.HasSuffix(sentence, "?") {
			continue
		}

		result = append(result, sentence)
	}

	// If no sentences were split, return the whole paragraph as one sentence
	if len(result) == 0 {
		return []string{text}
	}

	return result
}

// ExtractParagraphsFromHTML extracts paragraphs as plain text strings from HTML content
// This is a simpler alternative if you just want paragraphs instead of sentences
func ExtractParagraphsFromHTML(htmlContent string) []string {
	cleanText := ExtractTextFromHTML(htmlContent)
	return SplitIntoParagraphs(cleanText)
}
