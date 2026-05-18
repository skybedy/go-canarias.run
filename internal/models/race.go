package models

import (
	"strings"
	"time"
)

type Race struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DateRaw     string   `json:"date_raw"`    // Originální datum ze zdroje
	DateParsed  string   `json:"date_parsed"` // Naparsové datum do ISO, např. 2026-03-15
	Month       string   `json:"month"`       // Pro rychlé filtry, např. MAR
	Island      string   `json:"island"`
	Location    string   `json:"location"`
	Distances   []string `json:"distances"`   // Pole vzdáleností, např. ["126km", "84km"]
	Source      string   `json:"source"`      // CronoCanarias, Trackingsport atd.
	Status      string   `json:"status"`      // open, closing, past, cancelled
	URL         string   `json:"url"`         // Odkaz na registraci nebo detail
	Type        string   `json:"type"`        // trail, asphalt
	Description string   `json:"description"` // Volitelný popis
}

func (r Race) Day() string {
	parts := strings.Split(r.DateRaw, ".")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return "?"
}

// DateFormatted returns the date as dd.mm.yyyy, parsed from DateParsed (ISO 2006-01-02).
// Falls back to DateRaw if parsing fails.
func (r Race) DateFormatted() string {
	if r.DateParsed == "" || r.DateParsed == "0001-01-01" {
		return r.DateRaw
	}
	t, err := time.Parse("2006-01-02", r.DateParsed)
	if err != nil {
		return r.DateRaw
	}
	return t.Format("02.01.2006")
}

func (r Race) DistancesStr() string {
	return strings.Join(r.Distances, ", ")
}

func (r Race) TypeCZ() string {
	t := strings.ToLower(r.Type)
	switch {
	case strings.Contains(t, "trail") || strings.Contains(t, "hory") || strings.Contains(t, "montaña") || strings.Contains(t, "montana"):
		return "Horský běh / Trail"
	case strings.Contains(t, "asphalt") || strings.Contains(t, "asfalto") || strings.Contains(t, "silnice"):
		return "Silniční běh"
	case strings.Contains(t, "mtb") || strings.Contains(t, "bici"):
		return "Horská kola (MTB)"
	case strings.Contains(t, "plavání") || strings.Contains(t, "swim") || strings.Contains(t, "nado") || strings.Contains(t, "travesía"):
		return "Plavání"
	case strings.Contains(t, "marcha") || strings.Contains(t, "nordic"):
		return "Nordic Walking"
	case strings.Contains(t, "orient"):
		return "Orientační běh"
	default:
		if r.Type != "" {
			return r.Type
		}
		return "Závod"
	}
}
