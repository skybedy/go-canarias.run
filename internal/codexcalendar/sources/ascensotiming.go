package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const ascensoTimingURL = "https://ascensotiming.es/eventos/"

type AscensoTimingAdapter struct{}

func (AscensoTimingAdapter) Name() string { return "ascensotiming" }

func (AscensoTimingAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, ascensoTimingURL)
	if err != nil {
		return nil, err
	}
	return ParseAscensoTimingHTML(html), nil
}

var ascensoCardRe = regexp.MustCompile(`(?is)<img[^>]*alt="([^"]+)"[^>]*>.*?<a[^>]*href="(https://inscripciones\.ascensotiming\.es/inscripcion/[^"]+)"[^>]*>\s*Inscripciones\s*</a>`)

func ParseAscensoTimingHTML(html string) []calendar.Race {
	matches := ascensoCardRe.FindAllStringSubmatch(html, -1)
	out := make([]calendar.Race, 0, len(matches))

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		alt := stripHTML(m[1])
		link := strings.TrimSpace(m[2])
		if link == "" {
			continue
		}

		name := alt
		dateRaw := ""
		parts := strings.Split(alt, " - ")
		if len(parts) >= 2 {
			name = strings.TrimSpace(parts[0])
			dateRaw = strings.TrimSpace(parts[1])
		}
		dateLocal := parseDateLocal(dateRaw)
		if name == "" || dateLocal == "" {
			continue
		}

		island := inferIslandFromEventText(strings.ToLower(name + " " + link + " " + alt))
		if island == "" {
			island = "Tenerife"
		}

		race := buildRace("ascensotiming", ascensoTimingURL, name, dateLocal, "", island, link, name, 0.88)
		out = append(out, race)
	}

	return calendar.MergeCalendar(out)
}
