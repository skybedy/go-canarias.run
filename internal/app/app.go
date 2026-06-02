package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"canarias.run/internal/config"
	"canarias.run/internal/models"
	"canarias.run/internal/scraper"
	"canarias.run/internal/scraper/ahotu"
	"canarias.run/internal/scraper/ascensotiming"
	"canarias.run/internal/scraper/bulltiming"
	"canarias.run/internal/scraper/codex"
	"canarias.run/internal/scraper/conchipcanarias"
	"canarias.run/internal/scraper/cronocanarias"
	fca "canarias.run/internal/scraper/fca"
	"canarias.run/internal/scraper/fuerteventura"
	"canarias.run/internal/scraper/gomeraparadise"
	"canarias.run/internal/scraper/lanzarote"
	"canarias.run/internal/scraper/popularesgrancanaria"
	"canarias.run/internal/scraper/procrono"
	"canarias.run/internal/scraper/sportmaniacs"
	"canarias.run/internal/scraper/toptime"
	"canarias.run/internal/scraper/winerun"
	"canarias.run/internal/storage"
)

type App struct {
	Config  config.Config
	Store   storage.Storage
	Manager *scraper.Manager
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	store, err := newStorage(cfg)
	if err != nil {
		return nil, err
	}

	m := scraper.NewManager(store)
	registerScrapers(m)

	return &App{Config: cfg, Store: store, Manager: m}, nil
}

func (a *App) StartScraperDaemon(ctx context.Context) {
	if a.Config.EnableScraperDaemon {
		log.Printf("Spoustim scraper daemon (interval: %s)", a.Config.ScraperInterval)
		a.Manager.RunDaemon(ctx, a.Config.ScraperInterval)
		return
	}
	log.Println("Background scraping daemon je vypnuty pres ENABLE_SCRAPER_DAEMON.")
}

func (a *App) EnsureSeedData(ctx context.Context) {
	races, err := a.Store.GetAllRaces(ctx)
	if err != nil {
		log.Printf("Chyba pri cteni dat: %v", err)
		return
	}
	if len(races) > 0 {
		return
	}

	seed := []models.Race{
		{
			ID:         "1",
			DateRaw:    "15. 3. 2026",
			DateParsed: "2026-03-15",
			Month:      "MAR",
			Name:       "Transgrancanaria Ultra",
			Island:     "Gran Canaria",
			Location:   "Maspalomas",
			Distances:  []string{"126km", "84km", "46km"},
			Source:     "CronoCanarias",
			Status:     "open",
			URL:        "#",
			Type:       "trail",
		},
		{
			ID:         "2",
			DateRaw:    "2. 4. 2026",
			DateParsed: "2026-04-02",
			Month:      "APR",
			Name:       "Media Maraton Las Galletas",
			Island:     "Tenerife",
			Location:   "Arona",
			Distances:  []string{"21km", "10km"},
			Source:     "Sportmaniacs",
			Status:     "closing",
			URL:        "#",
			Type:       "asphalt",
		},
	}

	if err := a.Store.SaveRaces(ctx, seed); err != nil {
		log.Printf("Nepodarilo se ulozit uvodni data: %v", err)
	}
}

func newStorage(cfg config.Config) (storage.Storage, error) {
	switch strings.ToLower(cfg.StorageDriver) {
	case "json":
		return storage.NewJSONStorage(cfg.DataFile), nil
	case "sqlite":
		return storage.NewSQLiteStorage(cfg.SQLiteDSN)
	default:
		return nil, fmt.Errorf("unsupported STORAGE_DRIVER: %s", cfg.StorageDriver)
	}
}

func registerScrapers(m *scraper.Manager) {
	m.Register(cronocanarias.NewScraper())
	m.Register(sportmaniacs.NewScraper())
	m.Register(procrono.New())
	m.Register(conchipcanarias.New())
	m.Register(bulltiming.New())
	m.Register(ascensotiming.New())
	m.Register(popularesgrancanaria.New())
	m.Register(fca.New())
	m.Register(ahotu.New())
	m.Register(toptime.New())
	m.Register(gomeraparadise.New())
	m.Register(winerun.New())
	m.Register(fuerteventura.New())
	m.Register(lanzarote.New())
	m.Register(codex.New())
}
