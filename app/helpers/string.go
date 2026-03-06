package helpers

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Slugify converts a string to a URL-friendly slug.
func Slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// Truncate shortens a string to max length, appending "..." if truncated.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Contains checks if needle exists in haystack (case-insensitive).
func Contains(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

// Title converts a string to title case: "hello world" → "Hello World".
func Title(s string) string {
	return strings.Title(strings.ToLower(s))
}

// Excerpt returns the first n words from a string.
func Excerpt(s string, words int) string {
	parts := strings.Fields(s)
	if len(parts) <= words {
		return s
	}
	return strings.Join(parts[:words], " ") + "..."
}

// StripHTML removes all HTML tags from a string.
func StripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return html.UnescapeString(re.ReplaceAllString(s, ""))
}

// Mask masks a string showing only first/last n chars: "secret123" → "se*****23".
func Mask(s string, showFirst, showLast int) string {
	if len(s) <= showFirst+showLast {
		return s
	}
	masked := s[:showFirst] + strings.Repeat("*", len(s)-showFirst-showLast) + s[len(s)-showLast:]
	return masked
}
