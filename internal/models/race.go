package models

import (
	"fmt"
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

// DateFormatted returns the date as DD/MM/YYYY.
// Falls back to DateRaw if DateParsed is missing or invalid.
func (r Race) DateFormatted() string {
	if r.DateParsed == "" || r.DateParsed == "0001-01-01" {
		return r.DateRaw
	}
	t, err := time.Parse("2006-01-02", r.DateParsed)
	if err != nil {
		return r.DateRaw
	}
	return fmt.Sprintf("%d/%d/%02d", t.Day(), t.Month(), t.Year()%100)
}

func (r Race) DistancesStr() string {
	return strings.Join(r.Distances, ", ")
}

func (r Race) TypeLabel() string {
	switch r.Type {
	case "trail":
		return "Trail"
	case "road":
		return "Road"
	case "cross":
		return "Cross"
	case "orienteering":
		return "Orienteering"
	case "ocr":
		return "OCR"
	default:
		return r.Type
	}
}
