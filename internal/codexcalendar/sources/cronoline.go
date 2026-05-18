package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const cronolineURL = "https://cronolinecanarias.com/eventos/"

type CronolineCanariasAdapter struct{}

func (CronolineCanariasAdapter) Name() string { return "cronolinecanarias" }

func (CronolineCanariasAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, cronolineURL)
	if err != nil {
		return nil, err
	}
	return ParseCronolineHTML(html), nil
}

var eventBlockRe = regexp.MustCompile(`(?is)<article[^>]*type-event[^>]*>(.*?)</article>`)
var titleLinkRe = regexp.MustCompile(`(?is)<[^>]*class=["'][^"']*entry-title[^"']*["'][^>]*>\\s*<a[^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
var cronolineStartDateRe = regexp.MustCompile(`(?is)<[^>]*class=["'][^"']*start_date[^"']*["'][^>]*>(.*?)</[^>]+>`)

func ParseCronolineHTML(html string) []calendar.Race {
	chunks := eventBlockRe.FindAllStringSubmatch(html, -1)
	out := make([]calendar.Race, 0, len(chunks))

	for _, m := range chunks {
		if len(m) < 2 {
			continue
		}
		body := m[1]
		name := ""
		link := ""
		if a := titleLinkRe.FindStringSubmatch(body); len(a) > 2 {
			link = a[1]
			name = stripHTML(a[2])
		}
		if name == "" {
			if h := headingRe.FindStringSubmatch(body); len(h) > 1 {
				name = stripHTML(h[1])
			}
		}
		if name == "" {
			continue
		}

		dateRaw := ""
		if t := cronolineStartDateRe.FindStringSubmatch(body); len(t) > 1 {
			dateRaw = stripHTML(t[1])
		} else if t := timeDateRe.FindStringSubmatch(body); len(t) > 1 {
			dateRaw = t[1]
		} else if d := freeDateRe.FindStringSubmatch(stripHTML(body)); len(d) > 1 {
			dateRaw = d[1]
		}
		dateLocal := parseDateLocal(dateRaw)
		if dateLocal == "" {
			continue
		}

		text := strings.ToLower(stripHTML(body))
		city := ""
		island := ""
		switch {
		case strings.Contains(text, "tenerife"):
			island = "Tenerife"
		case strings.Contains(text, "gran canaria"):
			island = "Gran Canaria"
		case strings.Contains(text, "lanzarote"):
			island = "Lanzarote"
		case strings.Contains(text, "fuerteventura"):
			island = "Fuerteventura"
		}
		if island == "" {
			island = inferIslandFromEventText(strings.ToLower(name + " " + link + " " + text))
		}

		race := buildRace("cronolinecanarias", cronolineURL, name, dateLocal, city, island, link, name, 0.82)
		out = append(out, race)
	}

	return calendar.MergeCalendar(out)
}
