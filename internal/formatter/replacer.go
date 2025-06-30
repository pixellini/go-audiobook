package formatter

import "regexp"

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
