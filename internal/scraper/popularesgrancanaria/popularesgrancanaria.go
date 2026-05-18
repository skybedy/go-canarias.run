package popularesgrancanaria

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"canarias.run/internal/models"
	"canarias.run/internal/utils"
	"github.com/PuerkitoBio/goquery"
)

// PopularesGranCanariaScraper is a scraper for carreraspopularesgrancanaria.es
type PopularesGranCanariaScraper struct {
	Client *http.Client
}

// New returns a new instance of PopularesGranCanariaScraper
func New() *PopularesGranCanariaScraper {
	return &PopularesGranCanariaScraper{
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the name of the scraper
func (s *PopularesGranCanariaScraper) Name() string {
	return "PopularesGranCanaria"
}

// Scrape fetches events for PopularesGranCanaria
func (s *PopularesGranCanariaScraper) Scrape(ctx context.Context) ([]models.Race, error) {
	url := "https://carreraspopularesgrancanaria.es/calendario-de-carreras-populares-gran-canaria/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var races []models.Race

	doc.Find("table tbody tr").Each(func(i int, sel *goquery.Selection) {
		var cols []string
		sel.Find("td").Each(func(j int, td *goquery.Selection) {
			cols = append(cols, strings.TrimSpace(td.Text()))
		})

		if len(cols) >= 3 {
			dateStr := cols[0]
			raceTypeRaw := cols[1]
			raceName := cols[2]

			// Try to get explicit link
			link, _ := sel.Find("td").Eq(2).Find("a").Attr("href")
			link = strings.TrimSpace(link)

			if raceName == "" || dateStr == "" {
				return
			}

			// e.g., dateStr is "05/11/2025" or "15/11/2026"
			parsedDate, err := time.Parse("02/01/2006", dateStr)

			var raceDate time.Time
			if err == nil {
				raceDate = parsedDate
			}

			// Parse race type somewhat consistently
			raceType := "Running"
			rLower := strings.ToLower(raceTypeRaw)
			if strings.Contains(rLower, "trail") {
				raceType = "Trail"
			} else if strings.Contains(rLower, "asfalto") {
				raceType = "Asfalto" // or Road Running
			} else if strings.Contains(rLower, "orientación") {
				raceType = "Orientace" // Orienteering
			} else if strings.Contains(rLower, "mtb") || strings.Contains(rLower, "bici") {
				raceType = "MTB"
			} else if strings.Contains(rLower, "marcha") {
				raceType = "Marcha Nórdica"
			} else if raceTypeRaw != "" {
				raceType = raceTypeRaw
			}

			monthShort := ""
			if !raceDate.IsZero() {
				monthShort = strings.ToUpper(raceDate.Format("Jan"))
			}

			race := models.Race{
				Name:       raceName,
				DateRaw:    dateStr,
				DateParsed: raceDate.Format("2006-01-02"),
				Month:      monthShort,
				Island:     utils.IdentifyIsland(raceName, "Gran Canaria"),
				Location:   "Gran Canaria",
				Source:     "PopularesGranCanaria",
				URL:        link,
				Type:       raceType,
			}

			races = append(races, race)
		}
	})

	return races, nil
}
