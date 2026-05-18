package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const procronoEventsURL = "https://www.avaibooksports.com/sport-entity/procrono/events"

type ProcronoAdapter struct{}

func (ProcronoAdapter) Name() string { return "procrono" }

func (ProcronoAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, procronoEventsURL)
	if err != nil {
		return nil, err
	}
	return ParseProcronoEventsHTML(html), nil
}

var procronoDateRe = regexp.MustCompile(`(?is)<p\s+class="event-celebration-date"[^>]*>(.*?)</p>`)
var procronoTitleRe = regexp.MustCompile(`(?is)<p\s+class="event-title"[^>]*>(.*?)</p>`)
var procronoLinkRe = regexp.MustCompile(`(?is)<a[^>]*class="[^"]*go-to-event[^"]*"[^>]*href="([^"]+)"`)

func ParseProcronoEventsHTML(html string) []calendar.Race {
	dates := procronoDateRe.FindAllStringSubmatch(html, -1)
	titles := procronoTitleRe.FindAllStringSubmatch(html, -1)
	links := procronoLinkRe.FindAllStringSubmatch(html, -1)

	n := minInt(len(dates), len(titles), len(links))
	out := make([]calendar.Race, 0, n)
	for i := 0; i < n; i++ {
		name := stripHTML(titles[i][1])
		dateLocal := parseDateLocal(stripHTML(dates[i][1]))
		link := strings.TrimSpace(links[i][1])
		if name == "" || dateLocal == "" {
			continue
		}

		island := inferIslandFromEventText(strings.ToLower(name + " " + link))
		if island == "" {
			island = "Canarias"
		}

		race := buildRace("procrono", procronoEventsURL, name, dateLocal, "", island, link, name, 0.91)
		out = append(out, race)
	}

	return calendar.MergeCalendar(out)
}
