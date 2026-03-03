package scraper

import (
	"context"
	"log"
	"sync"
	"time"

	"canarias.run/internal/models"
	"canarias.run/internal/storage"
)

// Manager spravuje všechny registrované scrapery a spouští je.
type Manager struct {
	scrapers []Scraper
	store    storage.Storage
}

// NewManager vytvoří novou instanci manažeru.
func NewManager(store storage.Storage) *Manager {
	return &Manager{
		store: store,
	}
}

// Register přidá nový scraper do manažeru.
func (m *Manager) Register(s Scraper) {
	m.scrapers = append(m.scrapers, s)
}

// RunDaemon spustí sběr dat a bude ho periodicky opakovat.
func (m *Manager) RunDaemon(ctx context.Context, interval time.Duration) {
	// Pustíme to v samostatné goroutině
	go func() {
		m.ExecuteAll(ctx) // První spuštění hned

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.ExecuteAll(ctx)
			}
		}
	}()
}

// ExecuteAll spustí všechny scrapery paralelně a uloží výsledky do databáze.
func (m *Manager) ExecuteAll(ctx context.Context) {
	log.Println("[Manager] Začínám synchronizaci dat ze všech zdrojů...")

	var allRaces []models.Race
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, s := range m.scrapers {
		wg.Add(1)
		go func(sc Scraper) {
			defer wg.Done()

			// Scrapování dat pro konkrétní web
			races, err := sc.Scrape(ctx)
			if err != nil {
				log.Printf("[Scraper: %s] ❌ Chyba při scrapování: %v", sc.Name(), err)
				return
			}

			log.Printf("[Scraper: %s] ✅ Nalezeno %d závodů", sc.Name(), len(races))

			mu.Lock()
			allRaces = append(allRaces, races...)
			mu.Unlock()
		}(s)
	}

	// Počkáme, až se všechny weby projdou
	wg.Wait()

	if len(allRaces) > 0 {
		// V reálném nasazení bychom zde dělali mergování (aby se nesmazaly staré atd.),
		// pro MVP rovnou přepíšeme celou databázi novými daty.
		err := m.store.SaveRaces(ctx, allRaces)
		if err != nil {
			log.Printf("[Manager] ❌ Nelze uložit získaná data: %v", err)
		} else {
			log.Printf("[Manager] 🎉 Úspěšně uloženo %d závodů z %d zdrojů", len(allRaces), len(m.scrapers))
		}
	} else {
		log.Println("[Manager] ⚠️ Pozor: Žádné závody ke stažení nebyly nalezeny.")
	}
}
