package calendar

import (
	"regexp"
	"strings"
	"unicode"
)

var multiSpace = regexp.MustCompile(`\s+`)

func NormalizeName(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
			continue
		}
		b.WriteRune(' ')
	}
	out := multiSpace.ReplaceAllString(b.String(), " ")
	return strings.TrimSpace(out)
}

func NormalizeLocation(city, island string) string {
	c := NormalizeName(city)
	i := NormalizeName(island)
	if c != "" && i != "" {
		return c + "|" + i
	}
	if i != "" {
		return i
	}
	return c
}
