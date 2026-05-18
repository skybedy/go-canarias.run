package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"canarias.run/internal/codexcalendar"
)

var tagStripper = regexp.MustCompile(`(?s)<[^>]*>`)
var ws = regexp.MustCompile(`\s+`)

func fetchHTML(ctx context.Context, client *http.Client, target string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func stripHTML(s string) string {
	s = tagStripper.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = ws.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

func absURL(baseURL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	u, err := url.Parse(href)
	if err == nil && u.IsAbs() {
		return href
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}

func parseDateLocal(s string) string {
	s = strings.ReplaceAll(s, ",", " ")
	s = strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	if s == "" {
		return ""
	}
	layouts := []string{
		"2006-01-02",
		"2006-1-2",
		"02-01-2006",
		"2-1-2006",
		"02.01.2006",
		"2.1.2006",
		"02/01/2006",
		"2/1/2006",
		"2006/01/02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"02 Jan 2006",
		"2 Jan 2006",
		"02 January 2006",
		"2 January 2006",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// common Spanish month forms (marzo, mayo, etc.)
	es := strings.NewReplacer(
		"enero", "january",
		"febrero", "february",
		"marzo", "march",
		"abril", "april",
		"mayo", "may",
		"junio", "june",
		"julio", "july",
		"agosto", "august",
		"septiembre", "september",
		"setiembre", "september",
		"octubre", "october",
		"noviembre", "november",
		"diciembre", "december",
	)
	lower := strings.ToLower(s)
	translated := es.Replace(lower)
	for _, l := range []string{"2 January 2006", "02 January 2006"} {
		if t, err := time.Parse(l, translated); err == nil {
			return t.Format("2006-01-02")
		}
	}
	// "1 de enero de 2026, 9:00:00 CET"
	spanishLongRe := regexp.MustCompile(`(?i)\b(\d{1,2})\s+de\s+([a-záéíóúñ]+)\s+de\s+(\d{4})\b`)
	if m := spanishLongRe.FindStringSubmatch(lower); len(m) == 4 {
		candidate := strings.TrimSpace(m[1] + " " + es.Replace(m[2]) + " " + m[3])
		for _, l := range []string{"2 January 2006", "02 January 2006"} {
			if t, err := time.Parse(l, candidate); err == nil {
				return t.Format("2006-01-02")
			}
		}
	}
	// "15 marzo 2026"
	spanishShortRe := regexp.MustCompile(`(?i)\b(\d{1,2})\s+([a-záéíóúñ]+)\s+(\d{4})\b`)
	if m := spanishShortRe.FindStringSubmatch(lower); len(m) == 4 {
		candidate := strings.TrimSpace(m[1] + " " + es.Replace(m[2]) + " " + m[3])
		for _, l := range []string{"2 January 2006", "02 January 2006"} {
			if t, err := time.Parse(l, candidate); err == nil {
				return t.Format("2006-01-02")
			}
		}
	}

	return ""
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case json.Number:
		return t.String()
	case float64:
		return strings.TrimSpace(fmt.Sprintf("%.0f", t))
	case int:
		return fmt.Sprintf("%d", t)
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}

func buildRace(source, baseURL, name, dateLocal, city, island, link, rawTitle string, confidence float64) calendar.Race {
	now := time.Now().UTC()
	dateLocal = parseDateLocal(dateLocal)
	link = absURL(baseURL, link)
	return calendar.Race{
		Name:           strings.TrimSpace(name),
		NormalizedName: calendar.NormalizeName(name),
		DateLocal:      dateLocal,
		City:           strings.TrimSpace(city),
		Island:         strings.TrimSpace(island),
		URL:            link,
		Confidence:     confidence,
		FirstSeen:      now,
		LastSeen:       now,
		Sources: []calendar.SourceAttribution{{
			Source:     source,
			URL:        link,
			LastSeen:   now,
			RawTitle:   strings.TrimSpace(rawTitle),
			Confidence: confidence,
		}},
	}
}

func containsAnyKey(m map[string]any, keys ...string) bool {
	for _, k := range keys {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

func findObjects(v any, out *[]map[string]any) {
	switch t := v.(type) {
	case map[string]any:
		*out = append(*out, t)
		for _, child := range t {
			findObjects(child, out)
		}
	case []any:
		for _, child := range t {
			findObjects(child, out)
		}
	}
}

func extractAssignedObject(html, variable string) string {
	needle := variable + "="
	i := strings.Index(html, needle)
	if i < 0 {
		needle = variable + " ="
		i = strings.Index(html, needle)
	}
	if i < 0 {
		return ""
	}
	start := strings.IndexAny(html[i:], "[{")
	if start < 0 {
		return ""
	}
	start += i

	open := html[start]
	var close byte = '}'
	if open == '[' {
		close = ']'
	}
	depth := 0
	inString := false
	escaped := false
	for j := start; j < len(html); j++ {
		c := html[j]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inString = false
			}
			continue
		}
		if c == '"' {
			inString = true
			continue
		}
		if c == open {
			depth++
		}
		if c == close {
			depth--
			if depth == 0 {
				return html[start : j+1]
			}
		}
	}
	return ""
}

func extractJSONScripts(html string) []string {
	re := regexp.MustCompile(`(?is)<script[^>]*>(.*?)</script>`)
	ms := re.FindAllStringSubmatch(html, -1)
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		if len(m) < 2 {
			continue
		}
		body := strings.TrimSpace(m[1])
		if body == "" {
			continue
		}
		if strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[") {
			out = append(out, body)
		}
	}
	return out
}
