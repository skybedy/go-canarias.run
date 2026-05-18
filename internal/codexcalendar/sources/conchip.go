package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const conchipEventsURL = "https://inscripciones.conchipcanarias.com/sport-entity/conchipcanarias/events"

type ConchipCanariasAdapter struct{}

func (ConchipCanariasAdapter) Name() string { return "conchipcanarias" }

func (ConchipCanariasAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, conchipEventsURL)
	if err != nil {
		return nil, err
	}
	return ParseConchipEventsHTML(html), nil
}

var conchipDateRe = regexp.MustCompile(`(?is)<p\s+class="event-celebration-date"[^>]*>(.*?)</p>`)
var conchipTitleRe = regexp.MustCompile(`(?is)<p\s+class="event-title"[^>]*>(.*?)</p>`)
var conchipLinkRe = regexp.MustCompile(`(?is)<a[^>]*class="[^"]*go-to-event[^"]*"[^>]*href="([^"]+)"`)

func ParseConchipEventsHTML(html string) []calendar.Race {
	dates := conchipDateRe.FindAllStringSubmatch(html, -1)
	titles := conchipTitleRe.FindAllStringSubmatch(html, -1)
	links := conchipLinkRe.FindAllStringSubmatch(html, -1)

	n := minInt(len(dates), len(titles), len(links))
	out := make([]calendar.Race, 0, n)
	for i := 0; i < n; i++ {
		dateLocal := parseDateLocal(stripHTML(dates[i][1]))
		name := stripHTML(titles[i][1])
		link := strings.TrimSpace(links[i][1])
		if dateLocal == "" || name == "" {
			continue
		}

		island := inferIslandFromEventText(strings.ToLower(name + " " + link))
		if strings.Contains(strings.ToLower(name), "canarias") && island == "" {
			island = "Canarias"
		}

		race := buildRace("conchipcanarias", conchipEventsURL, name, dateLocal, "", island, link, name, 0.92)
		out = append(out, race)
	}

	return calendar.MergeCalendar(out)
}

func minInt(v int, others ...int) int {
	m := v
	for _, x := range others {
		if x < m {
			m = x
		}
	}
	return m
}
