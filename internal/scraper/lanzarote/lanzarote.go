package lanzarote

import (
	"context"

	"canarias.run/internal/models"
)

type Scraper struct{}

func New() *Scraper {
	return &Scraper{}
}

func (s *Scraper) Name() string {
	return "Lanzarote Deportes"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	return []models.Race{
		{
			Name:       "Carrera Mascaritas",
			DateRaw:    "28/02/2026",
			DateParsed: "2026-02-28",
			Island:     "Lanzarote",
			Location:   "Tinajo",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Trail/Carnival",
		},
		{
			Name:       "Volcano Triatlón Lanzarote",
			DateRaw:    "25/04/2026",
			DateParsed: "2026-04-25",
			Island:     "Lanzarote",
			Location:   "Tinajo",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Triathlon",
		},
		{
			Name:       "Famara Total Trail 2026",
			DateRaw:    "13-15/08/2026",
			DateParsed: "2026-08-13",
			Island:     "Lanzarote",
			Location:   "Teguise",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Trail",
		},
		{
			Name:       "Media Marathon de Playa Blanca",
			DateRaw:    "02/11/2026",
			DateParsed: "2026-11-02",
			Island:     "Lanzarote",
			Location:   "Playa Blanca",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Road",
		},
		{
			Name:       "Caminata del Vino",
			DateRaw:    "14/12/2026",
			DateParsed: "2026-12-14",
			Island:     "Lanzarote",
			Location:   "Yaiza",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Hiking",
		},
		{
			Name:       "Haria Titan 2026",
			DateRaw:    "Octubre 2026",
			DateParsed: "2026-10-10", // Estimates
			Island:     "Lanzarote",
			Location:   "Haría",
			Source:     s.Name(),
			URL:        "https://www.lanzarotedeportes.com/",
			Type:       "Trail",
		},
	}, nil
}
