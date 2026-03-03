package cronocanarias

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"canarias.run/internal/models"
	"github.com/PuerkitoBio/goquery"
)

// Scraper implementuje rozhraní scraper.Scraper pro CronoCanarias
type Scraper struct {
	URL string
}

// NewScraper vytvoří novou instanci CronoCanarias scraperu
func NewScraper() *Scraper {
	return &Scraper{
		URL: "https://cronolinecanarias.com/eventos/",
	}
}

func (s *Scraper) Name() string {
	return "CronoCanarias"
}

// Scrape navštíví cílovou stránku, stáhne HTML a parsuje závody.
func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	if err != nil {
		return nil, err
	}

	// Dobrá praxe u scrapování - tvářit se jako validní prohlížeč
	req.Header.Set("User-Agent", "Mozilla/5.0 (canarias.run bot)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("chyba HTTP: status kód %d", resp.StatusCode)
	}

	// Parsování struktury DOM dokumentu
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var races []models.Race

	// Hledání DOM elementů na stránce.
	// Očekáváme např. karty nebo tabulky událostí.
	doc.Find(".event-card, .course-item, .post-item").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find("h3, .title, .event-title").Text())
		dateRaw := strings.TrimSpace(sel.Find(".date, .event-date").Text())
		link, _ := sel.Find("a").Attr("href")
		location := strings.TrimSpace(sel.Find(".location, .event-location").Text())

		if title == "" {
			return // Pokud nenajdeme ani nadpis, pravděpodobně nejde o závod
		}

		race := models.Race{
			ID:        fmt.Sprintf("cronocanarias_%d", i),
			Name:      title,
			DateRaw:   dateRaw,
			Month:     "TBD", // Bude později doplněno detailnějším parsováním datumu
			Island:    "Rozpoznáno z textu",
			Location:  location,
			Distances: []string{}, // Získáme později například z popisu závodu
			Source:    s.Name(),
			Status:    "open",
			URL:       link,
			Type:      "trail", // Výchozí odhad, lze zpřesnit na základě tagů
		}

		races = append(races, race)
	})

	return races, nil
}
