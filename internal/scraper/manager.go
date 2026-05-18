package scraper

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strings"
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
		existing, err := m.store.GetAllRaces(ctx)
		if err != nil {
			log.Printf("[Manager] ⚠️ Nelze načíst stávající data pro merge: %v", err)
		}

		merged := MergeRaces(append(existing, allRaces...))
		err = m.store.SaveRaces(ctx, merged)
		if err != nil {
			log.Printf("[Manager] ❌ Nelze uložit získaná data: %v", err)
		} else {
			log.Printf("[Manager] 🎉 Úspěšně uloženo %d unikátních závodů z %d zdrojů (%d stažených položek před sloučením, %d sloučeno/odfiltrováno)", len(merged), len(m.scrapers), len(allRaces), len(allRaces)-len(merged))
		}
	} else {
		log.Println("[Manager] ⚠️ Pozor: Žádné závody ke stažení nebyly nalezeny.")
	}
}

// MergeRaces combines races from multiple scraper strategies and keeps the
// richest record for each date/name pair.
func MergeRaces(in []models.Race) []models.Race {
	merged := make(map[string]models.Race, len(in))
	for _, race := range in {
		key := raceKey(race)
		if key == "" {
			key = fallbackKey(race)
		}
		if key == "" {
			continue
		}

		if existing, ok := merged[key]; ok {
			merged[key] = mergeRace(existing, race)
			continue
		}
		merged[key] = race
	}

	out := make([]models.Race, 0, len(merged))
	for _, race := range merged {
		out = append(out, race)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DateParsed == out[j].DateParsed {
			return out[i].Name < out[j].Name
		}
		return out[i].DateParsed < out[j].DateParsed
	})
	return out
}

func mergeRace(existing, incoming models.Race) models.Race {
	out := existing
	if richerName(incoming.Name, out.Name) {
		out.Name = incoming.Name
	}
	if out.ID == "" {
		out.ID = incoming.ID
	}
	if out.DateRaw == "" {
		out.DateRaw = incoming.DateRaw
	}
	if out.DateParsed == "" {
		out.DateParsed = incoming.DateParsed
	}
	if out.Month == "" {
		out.Month = incoming.Month
	}
	if out.Island == "" || (out.Island == "Canarias" && incoming.Island != "") {
		out.Island = incoming.Island
	}
	if out.Location == "" || (out.Location == "Canarias" && incoming.Location != "") {
		out.Location = incoming.Location
	}
	if len(out.Distances) == 0 {
		out.Distances = incoming.Distances
	}
	out.Source = mergeSourceLabels(out.Source, incoming.Source)
	if out.Status == "" {
		out.Status = incoming.Status
	}
	if out.URL == "" || (strings.HasPrefix(out.URL, "#") && incoming.URL != "") {
		out.URL = incoming.URL
	}
	if out.Type == "" || out.Type == "Running" || out.Type == "running" {
		out.Type = incoming.Type
	}
	if out.Description == "" {
		out.Description = incoming.Description
	}
	return out
}

func raceKey(r models.Race) string {
	date := firstNonEmpty(r.DateParsed, r.DateRaw)
	name := normalizeKey(r.Name)
	if date == "" || name == "" {
		return ""
	}
	return date + "|" + name
}

func fallbackKey(r models.Race) string {
	url := strings.TrimSpace(strings.ToLower(r.URL))
	if url == "" || url == "#" {
		return ""
	}
	return "url|" + url
}

var nonKeyChars = regexp.MustCompile(`[^a-z0-9]+`)

func normalizeKey(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u",
		"ñ", "n", "ç", "c",
	).Replace(v)
	return strings.Trim(nonKeyChars.ReplaceAllString(v, ""), "_")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func richerName(a, b string) bool {
	return len(strings.TrimSpace(a)) > len(strings.TrimSpace(b))
}

func mergeSourceLabels(a, b string) string {
	labels := make([]string, 0, 4)
	seen := make(map[string]bool)
	for _, source := range []string{a, b} {
		for _, part := range strings.Split(source, "+") {
			part = strings.TrimSpace(part)
			if part == "" || seen[part] {
				continue
			}
			seen[part] = true
			labels = append(labels, part)
		}
	}
	return strings.Join(labels, "+")
}
