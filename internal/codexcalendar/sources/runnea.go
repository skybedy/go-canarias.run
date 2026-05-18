package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const runneaURL = "https://www.runnea.com/carreras-populares/islas-canarias/"

type RunneaAdapter struct{}

func (RunneaAdapter) Name() string { return "runnea" }

func (RunneaAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, runneaURL)
	if err != nil {
		return nil, err
	}
	return ParseRunneaHTML(html), nil
}

func ParseRunneaHTML(html string) []calendar.Race {
	blob := extractAssignedObject(html, "window.__NUXT__")
	if blob == "" {
		return nil
	}

	var root any
	if err := json.Unmarshal([]byte(blob), &root); err != nil {
		return nil
	}

	var objs []map[string]any
	findObjects(root, &objs)

	out := make([]calendar.Race, 0)
	for _, o := range objs {
		name := firstNonEmpty(toString(o["name"]), toString(o["title"]), toString(o["race_name"]))
		dateRaw := firstNonEmpty(toString(o["date"]), toString(o["start_date"]), toString(o["event_date"]))
		dateLocal := parseDateLocal(dateRaw)
		if name == "" || dateLocal == "" {
			continue
		}

		city := firstNonEmpty(toString(o["city"]), toString(o["town"]), toString(o["municipality"]))
		island := firstNonEmpty(toString(o["island"]), toString(o["province"]))
		link := firstNonEmpty(toString(o["url"]), toString(o["link"]), toString(o["slug"]))
		if link != "" && !strings.HasPrefix(link, "http") {
			link = "/" + strings.TrimPrefix(link, "/")
		}

		conf := 0.90
		if city == "" && island == "" {
			conf = 0.82
		}

		race := buildRace("runnea", runneaURL, name, dateLocal, city, island, link, name, conf)
		out = append(out, race)
	}
	return calendar.MergeCalendar(out)
}
