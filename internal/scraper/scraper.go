package scraper

import (
	"context"

	"canarias.run/internal/models"
)

// Scraper je společné rozhraní pro všechny zdroje dat.
// Každý zdroj (CronoCanarias, Trackingsport atd.) bude implementovat toto rozhraní,
// čímž sjednotíme způsob, jakým získáváme data pro naši aplikaci.
type Scraper interface {
	// Name vrací jméno zdroje, např. "CronoCanarias"
	Name() string

	// Scrape stáhne závody ze zdroje.
	// Přijímá context pro možnost zrušení operace.
	Scrape(ctx context.Context) ([]models.Race, error)
}
