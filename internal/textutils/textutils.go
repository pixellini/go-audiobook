package textutils

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
)

const textDelimiter = "|||"
const maxCharSizeLimitTTS = 250

var (
	sentenceRegex  = regexp.MustCompile(`([.!?])\s+([A-Z])`)
	paragraphRegex = regexp.MustCompile(`\n\s*\n`)
)

// SplitText splits the content into sentences based on punctuation and line breaks.
func SplitText(content string) []string {
	cleanText := ExtractTextFromHTML(content)
	return SplitIntoParagraphs(cleanText)
}

// ExtractTextFromHTML extracts plain text from HTML content, preserving paragraph structure
func ExtractTextFromHTML(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return cleanText(strip.StripTags(htmlContent))
	}

	var textParts []string
	doc.Find("p, div, h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && isValidContent(text) {
			textParts = append(textParts, text)
		}
	})

	if len(textParts) == 0 {
		return cleanText(strip.StripTags(htmlContent))
	}

	return strings.Join(textParts, "\n\n")
}

// isValidContent filters out XML/CSS declarations and other non-content
func isValidContent(text string) bool {
	text = strings.TrimSpace(text)

	return text != "" &&
		!strings.HasPrefix(text, "<?xml") &&
		!strings.HasPrefix(text, "<!DOCTYPE") &&
		!strings.HasPrefix(text, "@page") &&
		!(strings.IndexByte(text, '.') == 0 && strings.IndexByte(text, '{') != -1) &&
		!(strings.IndexByte(text, '#') == 0 && strings.IndexByte(text, '{') != -1)
}

// cleanText removes unwanted content lines
func cleanText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && isValidContent(line) {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// SplitIntoParagraphs splits text into manageable chunks for TTS
func SplitIntoParagraphs(text string) []string {
	paragraphs := paragraphRegex.Split(text, -1)
	var result []string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		// Split the long paragraphs
		if len(p) > maxCharSizeLimitTTS+50 {
			result = append(result, splitLongText(p)...)
		} else {
			result = append(result, p)
		}
	}

	return result
}

// splitLongText splits long text on sentences, and if still too long, on commas
func splitLongText(text string) []string {
	// Try splitting on sentence boundaries first.
	// The triple pipe ||| is functioning as a custom delimiter that's unlikely to exist in normal text.
	formattedText := sentenceRegex.ReplaceAllString(text, "$1"+textDelimiter+"$2")
	sentences := strings.Split(formattedText, textDelimiter)

	var result []string
	textLengthLimit := maxCharSizeLimitTTS - 50

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		// If sentence is still too long, split on commas
		if len(sentence) > textLengthLimit {
			result = append(result, splitOnCommas(sentence)...)
		} else {
			result = append(result, sentence)
		}
	}

	// Fallback if no splitting occurred
	if len(result) == 0 {
		if len(text) > textLengthLimit {
			return splitOnCommas(text)
		}
		return []string{text}
	}

	return result
}

// splitOnCommas splits text on commas
// There's additional logic here because strings.Split() will only give us tiny fragments that waste the 250 char TTS limit.
// Instead, we pack as many comma parts as possible into each chunk until we hit the limit, then start fresh.
// This way we get fewer API calls and more natural sounding speech.
func splitOnCommas(text string) []string {
	// The size limit for TTS is usually 250 characters, so if the paragraph is shorter than that, then it should be fine
	if len(text) < maxCharSizeLimitTTS || strings.IndexByte(text, ',') == -1 {
		return []string{text}
	}

	parts := strings.Split(text, ",")
	if len(parts) <= 1 {
		return []string{text}
	}

	var result []string
	var builder strings.Builder
	builder.Grow(maxCharSizeLimitTTS + 50)

	// Here is where we loop through and pack as many comma parts into the TTS limit
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Calculate the length we'd need if we add this part
		needLen := builder.Len()
		if needLen > 0 {
			needLen++
		}
		needLen += len(part)
		if i < len(parts)-1 {
			needLen++
		}

		// If adding this part would exceed TTS limit and we have content, flush current chunk (sentence)
		if needLen > maxCharSizeLimitTTS && builder.Len() > 0 {
			result = append(result, builder.String())
			builder.Reset()
		}

		// Add a separator
		if builder.Len() > 0 {
			builder.WriteByte(' ')
		}

		builder.WriteString(part)

		// Add comma back unless it's the end of the sentence
		if i < len(parts)-1 {
			builder.WriteByte(',')
		}
	}

	// Add remaining content
	if builder.Len() > 0 {
		result = append(result, builder.String())
	}

	return result
}

// ExtractParagraphsFromHTML extracts paragraphs as plain text strings from HTML content
func ExtractParagraphsFromHTML(htmlContent string) []string {
	cleanText := ExtractTextFromHTML(htmlContent)
	return SplitIntoParagraphs(cleanText)
}

// ExtractTitleFromHTML extracts the chapter title from HTML content by looking for heading tags
func ExtractTitleFromHTML(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	// Look for title in heading tags, prioritising h1, then h2, etc.
	headingSelectors := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	
	for _, selector := range headingSelectors {
		title := strings.TrimSpace(doc.Find(selector).First().Text())
		if title != "" && isValidContent(title) {
			return title
		}
	}

	return ""
}
