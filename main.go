package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
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

// Global variable pro šablony a úložiště
var (
	tmpl  *template.Template
	store storage.Storage
)

func init() {
	// Parsování šablon
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	log.Println("=== Inicializace canarias.run ===")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Chyba konfigurace: %v", err)
	}

	// Inicializace JSON úložiště
	store = storage.NewJSONStorage(cfg.DataFile)

	// Přípojka pro systém stahování (Scraper Orchestrator)
	m := scraper.NewManager(store)
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

	// Spustit Scrapovacího robota na pozadí (synch jednou za 24 hodin)
	if cfg.EnableScraperDaemon {
		log.Printf("Spouštím robota na pozadí pro scrapování dat (interval: %s)...", cfg.ScraperInterval)
		m.RunDaemon(context.Background(), cfg.ScraperInterval)
	} else {
		log.Println("Background scraping daemon je vypnutý přes ENABLE_SCRAPER_DAEMON.")
	}

	// Pokusíme se načíst data, pokud neexistují, vložíme úvodní
	// initMockDataIfEmpty() (Mockdata už teď nepotřebujeme plnit, ale pro sichr si to zde necháme jako fall-back)
	initMockDataIfEmpty()

	// Servírování statických souborů
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Hlavní route
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/healthz", handleHealthz)

	log.Printf("Server poslouchá na http://localhost:%s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Chyba při spouštění serveru: %v", err)
	}
}

func initMockDataIfEmpty() {
	ctx := context.Background()
	races, err := store.GetAllRaces(ctx)
	if err != nil {
		log.Printf("Chyba při čtení dat: %v", err)
		return
	}

	// Pokud už nějaká data máme, nic neděláme
	if len(races) > 0 {
		return
	}

	mockRaces := []models.Race{
		{
			ID:        "1",
			DateRaw:   "15. 3. 2026",
			Month:     "MAR",
			Name:      "Transgrancanaria Ultra",
			Island:    "Gran Canaria",
			Location:  "Maspalomas",
			Distances: []string{"126km", "84km", "46km"},
			Source:    "CronoCanarias",
			Status:    "open",
			URL:       "#",
			Type:      "trail",
		},
		{
			ID:        "2",
			DateRaw:   "2. 4. 2026",
			Month:     "APR",
			Name:      "Media Maratón Las Galletas",
			Island:    "Tenerife",
			Location:  "Arona",
			Distances: []string{"21km", "10km"},
			Source:    "Sportmaniacs",
			Status:    "closing",
			URL:       "#",
			Type:      "asphalt",
		},
	}

	err = store.SaveRaces(ctx, mockRaces)
	if err != nil {
		log.Printf("Nepodařilo se uložit úvodní data: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Získáme všechny závody z našeho úložiště
	races, err := store.GetAllRaces(context.Background())
	if err != nil {
		log.Printf("Chyba při získávání závodů: %v", err)
		races = []models.Race{} // Fallback na prázdné pole
	}

	// Seřadit závody chronologicky podle DateParsed
	sort.Slice(races, func(i, j int) bool {
		return races[i].DateParsed < races[j].DateParsed
	})

	// Vyfiltrovat pouze závody roku 2026 a novější a odstranit duplicity
	seen := make(map[string]bool)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	filtered := races[:0]

	for _, r := range races {
		// Rok 2026+
		if r.DateParsed < "2026-01-01" {
			continue
		}

		// Klíč pro duplicity: normalizované jméno + datum
		nameKey := strings.ToLower(r.Name)
		nameKey = re.ReplaceAllString(nameKey, "")
		key := r.DateParsed + "_" + nameKey

		if !seen[key] {
			seen[key] = true
			filtered = append(filtered, r)
		}
	}
	races = filtered

	// Renderování indexové šablony s daty
	err = tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Title": "canarias.run - Tvůj ultimátní kalendář běžeckých závodů",
		"Races": races,
	})
	if err != nil {
		log.Printf("Chyba při renderování šablony: %v", err)
		http.Error(w, "Vnitřní chyba serveru", http.StatusInternalServerError)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
