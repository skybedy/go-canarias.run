package gomeraparadise

import (
	"context"

	"canarias.run/internal/models"
)

type Scraper struct{}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) Name() string {
	return "Gomera Paradise"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	// Hardcoded because it's a major event that we want to ENSURE is there
	// and its site is hard to scrape dynamically for date/distances if not using browser
	return []models.Race{
		{
			Name:       "Gomera Paradise Trail 2026",
			DateRaw:    "18-20 Septiembre 2026",
			DateParsed: "2026-09-18",
			Island:     "La Gomera",
			Location:   "San Sebastián de La Gomera",
			Source:     s.Name(),
			URL:        "https://gomeraparadise.com/",
			Type:       "Trail",
			Distances:  []string{"10k", "16k", "30k", "43k", "75k"},
		},
	}, nil
}
