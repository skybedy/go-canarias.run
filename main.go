package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"canarias.run/internal/models"
	"canarias.run/internal/scraper"
	"canarias.run/internal/scraper/cronocanarias"
	"canarias.run/internal/scraper/sportmaniacs"
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

	// Inicializace JSON úložiště
	store = storage.NewJSONStorage("data.json")

	// Přípojka pro systém stahování (Scraper Orchestrator)
	m := scraper.NewManager(store)
	m.Register(cronocanarias.NewScraper())
	m.Register(sportmaniacs.NewScraper())

	// Spustit Scrapovacího robota na pozadí (synch jednou za 24 hodin)
	log.Println("Spouštím robota na pozadí pro scrapování reálných dat...")
	m.RunDaemon(context.Background(), 24*time.Hour)

	// Pokusíme se načíst data, pokud neexistují, vložíme úvodní
	// initMockDataIfEmpty() (Mockdata už teď nepotřebujeme plnit, ale pro sichr si to zde necháme jako fall-back)
	initMockDataIfEmpty()

	// Servírování statických souborů
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Hlavní route
	http.HandleFunc("/", handleIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server poslouchá na http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, nil)
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
