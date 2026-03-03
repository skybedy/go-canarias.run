package models

import "strings"

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

func (r Race) DistancesStr() string {
	return strings.Join(r.Distances, ", ")
}
