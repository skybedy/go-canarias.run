package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const bulltimingURL = "https://bulltiming.es/proximos-eventos/"

type BullTimingAdapter struct{}

func (BullTimingAdapter) Name() string { return "bulltiming" }

func (BullTimingAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, bulltimingURL)
	if err != nil {
		return nil, err
	}
	return ParseBullTimingHTML(html), nil
}

var bullItemStartRe = regexp.MustCompile(`(?is)<div\s+class="jet-listing-grid__item[^>]*>`)
var bullTitleLinkRe = regexp.MustCompile(`(?is)<div\s+class="elementor-heading-title[^>]*>\s*<a[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
var bullDateRe = regexp.MustCompile(`(?is)<div\s+class="jet-listing-dynamic-field__content">(.*?)</div>`)

func ParseBullTimingHTML(html string) []calendar.Race {
	starts := bullItemStartRe.FindAllStringIndex(html, -1)
	out := make([]calendar.Race, 0, len(starts))

	for i := 0; i < len(starts); i++ {
		start := starts[i][0]
		end := len(html)
		if i+1 < len(starts) {
			end = starts[i+1][0]
		}
		block := html[start:end]

		titleMatch := bullTitleLinkRe.FindStringSubmatch(block)
		if len(titleMatch) < 3 {
			continue
		}
		link := strings.TrimSpace(titleMatch[1])
		name := stripHTML(titleMatch[2])

		dateRaw := ""
		if dm := bullDateRe.FindStringSubmatch(block); len(dm) > 1 {
			dateRaw = stripHTML(dm[1])
		}
		dateLocal := parseDateLocal(dateRaw)
		if name == "" || dateLocal == "" {
			continue
		}

		island := inferIslandFromEventText(strings.ToLower(name + " " + link))
		if island == "" {
			island = "Canarias"
		}

		race := buildRace("bulltiming", bulltimingURL, name, dateLocal, "", island, link, name, 0.86)
		out = append(out, race)
	}

	return calendar.MergeCalendar(out)
}
