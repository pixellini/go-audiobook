package epub

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	epubReader "github.com/taylorskalyo/goreader/epub"
)

type Epub struct {
	Title        string
	Author       string
	Description  string
	Language     string
	Introduction string // This is for the audiobook speaker as a welcome.
	Chapters     []Chapter
	dir          string
}

func New(dir string) (*Epub, error) {
	epub := &Epub{
		dir: dir,
	}

	if err := epub.setMetadata(); err != nil {
		return epub, err
	}

	if err := epub.setChapters(); err != nil {
		return epub, err
	}

	return epub, nil
}

func (epub *Epub) setMetadata() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	epub.Title = book.Title
	epub.Description = book.Description
	epub.Language = book.Language
	epub.Author = book.Creator
	epub.Introduction = fmt.Sprintf("%s written by %s.", book.Title, book.Creator)

	return nil
}

func (epub *Epub) setChapters() error {
	epubFile, err := epubReader.OpenReader(epub.dir)

	if err != nil {
		return fmt.Errorf("failed to open EPUB: %w", err)
	}

	defer epubFile.Close()

	if len(epubFile.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found in EPUB")
	}

	book := epubFile.Rootfiles[0]

	// Filter only the chapters to read. Ignore contents, images, indexes etc.
	var chapters []epubReader.Item
	for _, item := range book.Manifest.Items {
		if isAcceptedChapterItem(item.ID) {
			chapters = append(chapters, item)
		}
	}

	for _, item := range chapters {
		r, err := item.Open()
		if err != nil {
			return fmt.Errorf("failed to open item %s: %w", item.ID, err)
		}

		content, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return fmt.Errorf("failed to read item %s: %w", item.ID, err)
		}

		rawHtmlContent := string(content)
		processedChapter := createChapter(rawHtmlContent)
		epub.Chapters = append(epub.Chapters, processedChapter)
	}

	introduction := Chapter{
		Title:      "Introduction",
		Paragraphs: SplitText(epub.Introduction),
		RawText:    epub.Introduction,
	}

	// Add the introduction to the beginning of the chapters
	epub.Chapters = append([]Chapter{introduction}, epub.Chapters...)

	// Remove chapters titled 'contents', 'table of contents', or same as book title (case-insensitive, trimmed)
	filtered := make([]Chapter, 0, len(epub.Chapters))
	bookTitleLower := strings.TrimSpace(strings.ToLower(epub.Title))

	for _, ch := range epub.Chapters {
		t := strings.TrimSpace(strings.ToLower(ch.Title))
		if t == "contents" || t == "table of contents" || t == bookTitleLower {
			continue
		}
		filtered = append(filtered, ch)
	}
	epub.Chapters = filtered

	return nil
}

func SplitText(content string) []string {
	content = CleanHTML(content)

	reSingleLineBreak := regexp.MustCompile(`([^\n])\n([^\n])`)
	for reSingleLineBreak.MatchString(content) {
		content = reSingleLineBreak.ReplaceAllString(content, "$1 $2")
	}

	re := regexp.MustCompile(`(?m)([^\n]+?[\.!?](?:\s+|$)|[^\n]+$)`)
	matches := re.FindAllString(content, -1)
	var result []string
	for _, match := range matches {
		trimmed := strings.TrimSpace(match)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func CleanHTML(content string) string {
	content = regexp.MustCompile(`(?s)<!DOCTYPE[^>]*>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`(?s)<\?xml[^>]*\?>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")
	content = html.UnescapeString(content)
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	return strings.TrimSpace(content)
}

func CleanContent(htmlContent string) string {
	htmlContent = regexp.MustCompile(`(?is)<head.*?>.*?</head>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<\?xml.*?\?>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<!DOCTYPE[^>]*>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<html[^>]*>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)</html>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)<body[^>]*>`).ReplaceAllString(htmlContent, "")
	htmlContent = regexp.MustCompile(`(?s)</body>`).ReplaceAllString(htmlContent, "")
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

func ReplaceCodePattern(content string) string {
	re := regexp.MustCompile(`(?m)^(Figure|Fig|Listing|Example)\b.*$`)
	return re.ReplaceAllString(content, "...")
}

func ReplaceUrlPattern(content string) string {
	re := regexp.MustCompile(`(https?://\S+|www\.\S+)`)
	if re.MatchString(content) {
		return re.ReplaceAllString(content, "URL Link")
	}
	return content
}
