package storage

import (
	"context"

	"canarias.run/internal/models"
)

// Storage definuje rozhraní pro ukládání a načítání závodů.
// Můžeme díky tomu snadno přepínat mezi ukládáním do JSONu, SQLite nebo např. PostgreSQL.
type Storage interface {
	// SaveRaces uloží seznam závodů (např. po úspěšném scrapování).
	SaveRaces(ctx context.Context, races []models.Race) error

	// GetAllRaces vrátí všechny uložené závody.
	GetAllRaces(ctx context.Context) ([]models.Race, error)

	// GetRacesFilter vrátí závody podle specifických filtrů (lze rozšířit).
	// GetRacesFilter(ctx context.Context, island, surface, month string) ([]models.Race, error)
}
