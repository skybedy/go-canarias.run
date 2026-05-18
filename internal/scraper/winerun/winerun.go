package winerun

import (
	"context"

	"canarias.run/internal/models"
)

type Scraper struct{}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) Name() string {
	return "Lanzarote Wine Run"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	return []models.Race{
		{
			Name:       "Lanzarote Wine Run 2026",
			DateRaw:    "13-14 Junio 2026",
			DateParsed: "2026-06-13",
			Island:     "Lanzarote",
			Location:   "La Geria",
			Source:     s.Name(),
			URL:        "http://www.lanzarotewinerun.com/",
			Type:       "Trail",
			Distances:  []string{"12k", "23k"},
		},
	}, nil
}
