package fuerteventura

import (
	"context"

	"canarias.run/internal/models"
)

type Scraper struct{}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) Name() string {
	return "Federación Fuerteventura"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	return []models.Race{
		{
			Name:       "2ª Maratón Dany Sport Corralejo Grandes Playas",
			DateRaw:    "17/01/2026",
			DateParsed: "2026-01-17",
			Island:     "Fuerteventura",
			Location:   "Corralejo",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Road",
		},
		{
			Name:       "XXIV Cross Insular Municipio de Antigua",
			DateRaw:    "07/02/2026",
			DateParsed: "2026-02-07",
			Island:     "Fuerteventura",
			Location:   "Antigua",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Cross",
		},
		{
			Name:       "9ª Carrera Carnavalera",
			DateRaw:    "07/03/2026",
			DateParsed: "2026-03-07",
			Island:     "Fuerteventura",
			Location:   "Corralejo",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Road",
		},
		{
			Name:       "XI Carrera Trotadunas",
			DateRaw:    "25/04/2026",
			DateParsed: "2026-04-25",
			Island:     "Fuerteventura",
			Location:   "Puerto del Rosario",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Trail",
		},
		{
			Name:       "Paladea La Oliva",
			DateRaw:    "10/05/2026",
			DateParsed: "2026-05-10",
			Island:     "Fuerteventura",
			Location:   "Villaverde",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Trail/MTB",
		},
		{
			Name:       "XIV Milla Barceló Caleta de Fuste",
			DateRaw:    "22/05/2026",
			DateParsed: "2026-05-22",
			Island:     "Fuerteventura",
			Location:   "Caleta de Fuste",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Road",
		},
		{
			Name:       "XI Baifo Extreme",
			DateRaw:    "13/09/2026",
			DateParsed: "2026-09-13",
			Island:     "Fuerteventura",
			Location:   "Tindaya",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "OCR/Trail",
		},
		{
			Name:       "XI Subida a Betancuria",
			DateRaw:    "17/10/2026",
			DateParsed: "2026-10-17",
			Island:     "Fuerteventura",
			Location:   "Betancuria",
			Source:     s.Name(),
			URL:        "https://atletismomajorero.es/",
			Type:       "Trail",
		},
	}, nil
}
