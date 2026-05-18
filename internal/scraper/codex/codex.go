package codex

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	codexcalendar "canarias.run/internal/codexcalendar"
	"canarias.run/internal/codexcalendar/sources"
	"canarias.run/internal/models"
)

type Scraper struct {
	Client *http.Client
}

func New() *Scraper {
	return &Scraper{
		Client: &http.Client{Timeout: 25 * time.Second},
	}
}

func (s *Scraper) Name() string {
	return "CodexUnified"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	client := s.Client
	if client == nil {
		client = http.DefaultClient
	}

	races, err := codexcalendar.Aggregate(ctx, client, sources.PriorityAdapters()...)
	races = codexcalendar.EnrichLocations(races)

	out := make([]models.Race, 0, len(races))
	for _, race := range races {
		if race.DateLocal == "" {
			continue
		}
		out = append(out, convertRace(race))
	}

	if err != nil && len(out) > 0 {
		log.Printf("[Scraper: %s] partial source error: %v", s.Name(), err)
		return out, nil
	}
	return out, err
}

func convertRace(r codexcalendar.Race) models.Race {
	return models.Race{
		ID:          "codex|" + r.ID,
		Name:        r.Name,
		DateRaw:     r.DateLocal,
		DateParsed:  r.DateLocal,
		Month:       monthAbbrev(r.DateLocal),
		Island:      r.Island,
		Location:    r.City,
		Source:      sourceLabel(r),
		Status:      "open",
		URL:         r.URL,
		Type:        inferType(r.Name),
		Description: "Imported from the Codex unified calendar pipeline.",
	}
}

func sourceLabel(r codexcalendar.Race) string {
	if len(r.Sources) == 0 {
		return "CodexUnified"
	}
	seen := make(map[string]bool, len(r.Sources))
	var names []string
	for _, source := range r.Sources {
		if source.Source == "" || seen[source.Source] {
			continue
		}
		seen[source.Source] = true
		names = append(names, source.Source)
	}
	if len(names) == 0 {
		return "CodexUnified"
	}
	return "CodexUnified:" + strings.Join(names, "+")
}

func monthAbbrev(dateLocal string) string {
	if len(dateLocal) < 7 {
		return ""
	}
	switch dateLocal[5:7] {
	case "01":
		return "JAN"
	case "02":
		return "FEB"
	case "03":
		return "MAR"
	case "04":
		return "APR"
	case "05":
		return "MAY"
	case "06":
		return "JUN"
	case "07":
		return "JUL"
	case "08":
		return "AUG"
	case "09":
		return "SEP"
	case "10":
		return "OCT"
	case "11":
		return "NOV"
	case "12":
		return "DEC"
	default:
		return ""
	}
}

func inferType(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "trail"), strings.Contains(lower, "montaña"), strings.Contains(lower, "montana"), strings.Contains(lower, "vertical"), strings.Contains(lower, "cxm"):
		return "trail"
	case strings.Contains(lower, "mtb"), strings.Contains(lower, "btt"), strings.Contains(lower, "bike"), strings.Contains(lower, "bici"):
		return "mtb"
	case strings.Contains(lower, "swim"), strings.Contains(lower, "nado"), strings.Contains(lower, "travesia"), strings.Contains(lower, "travesía"):
		return "swim"
	default:
		return "running"
	}
}
