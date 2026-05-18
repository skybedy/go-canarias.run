package calendar

import "time"

// ProvisionalRaces returns a temporary manually seeded calendar
// suitable for local web preview before full live ingestion is wired.
func ProvisionalRaces() []Race {
	now := time.Now().UTC()
	mk := func(name, date, city, island, url, source string, conf float64) Race {
		r := Race{
			Name:           name,
			NormalizedName: NormalizeName(name),
			DateLocal:      date,
			City:           city,
			Island:         island,
			URL:            url,
			Confidence:     conf,
			FirstSeen:      now,
			LastSeen:       now,
			Sources: []SourceAttribution{{
				Source:     source,
				URL:        url,
				LastSeen:   now,
				RawTitle:   name,
				Confidence: conf,
			}},
		}
		r.ID = DedupeKey(r)
		return r
	}

	return MergeCalendar([]Race{
		mk("10K Arona", "2026-03-22", "Arona", "Tenerife", "https://www.runnea.com/carreras-populares/10k-arona/", "runnea", 0.90),
		mk("Media Maraton Las Palmas", "2026-04-10", "Las Palmas", "Gran Canaria", "https://www.carreraspopularesgrancanaria.com/", "carreraspopularesgrancanaria", 0.86),
		mk("Trail Adeje", "2026-04-19", "Adeje", "Tenerife", "https://cronolinecanarias.com/eventos/", "cronolinecanarias", 0.82),
		mk("Night Run La Laguna", "2026-05-02", "La Laguna", "Tenerife", "https://bulltiming.com/", "bulltiming", 0.84),
		mk("Vertical Los Realejos", "2026-05-16", "Los Realejos", "Tenerife", "https://cronolinecanarias.com/eventos/", "cronolinecanarias", 0.80),
		mk("Gran Fondo Maspalomas", "2026-05-24", "Maspalomas", "Gran Canaria", "https://www.carreraspopularesgrancanaria.com/", "carreraspopularesgrancanaria", 0.85),
		mk("Ruta del Vino Lanzarote", "2026-06-07", "La Geria", "Lanzarote", "https://www.runnea.com/carreras-populares/", "runnea", 0.83),
		mk("Cross Puerto del Rosario", "2026-06-21", "Puerto del Rosario", "Fuerteventura", "https://bulltiming.com/", "bulltiming", 0.78),
		mk("Subida Tejina", "2026-07-05", "Tejina", "Tenerife", "https://cronolinecanarias.com/eventos/", "cronolinecanarias", 0.79),
		mk("10K Las Canteras", "2026-09-13", "Las Palmas", "Gran Canaria", "https://www.carreraspopularesgrancanaria.com/", "carreraspopularesgrancanaria", 0.87),
		mk("Media Maraton Santa Cruz", "2026-10-04", "Santa Cruz de Tenerife", "Tenerife", "https://www.runnea.com/carreras-populares/", "runnea", 0.88),
		mk("San Silvestre Canaria", "2026-12-31", "Las Palmas", "Gran Canaria", "https://bulltiming.com/", "bulltiming", 0.76),
	})
}
