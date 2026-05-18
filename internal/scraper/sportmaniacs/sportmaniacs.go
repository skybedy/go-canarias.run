package sportmaniacs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"canarias.run/internal/models"
	"canarias.run/internal/utils"
	"github.com/PuerkitoBio/goquery"
)

// Scraper implementuje rozhraní scraper.Scraper pro Sportmaniacs
type Scraper struct {
	URL string
}

func NewScraper() *Scraper {
	return &Scraper{
		URL: "https://sportmaniacs.com/es/races",
	}
}

func (s *Scraper) Name() string {
	return "Sportmaniacs"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	if err != nil {
		return nil, err
	}
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var races []models.Race

	// Zde by přišel selektor specifický pro Sportmaniacs, většinou ".race-card"
	doc.Find(".race-card, .list-group-item.race").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find(".race-name, h4").Text())
		dateRaw := strings.TrimSpace(sel.Find(".race-date, .date").Text())
		link, _ := sel.Find("a.btn, a.race-link").Attr("href")
		location := strings.TrimSpace(sel.Find(".race-location, .location").Text())

		if title == "" {
			return
		}

		// Doplníme úplnou URL, pokud je relativní
		if link != "" && strings.HasPrefix(link, "/") {
			link = "https://sportmaniacs.com" + link
		}

		race := models.Race{
			ID:        fmt.Sprintf("sportmaniacs_%d", i),
			Name:      title,
			DateRaw:   dateRaw,
			Month:     "TBD",
			Island:    utils.IdentifyIsland(title, location),
			Location:  location,
			Distances: []string{},
			Source:    s.Name(),
			Status:    "open",
			URL:       link,
			Type:      "asphalt",
		}

		// Only add if it's actually in Canarias
		if race.Island != "Canarias" || strings.Contains(strings.ToLower(title+location), "canaria") {
			races = append(races, race)
		}
	})

	return races, nil
}
