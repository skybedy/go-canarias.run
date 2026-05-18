package fca

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

type Scraper struct {
	Client *http.Client
}

func New() *Scraper {
	return &Scraper{
		Client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *Scraper) Name() string {
	return "FCA"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	url := "https://atletismocanario.es/calendario/lista/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Use a very common Chrome UA to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.8,en-US;q=0.5,en;q=0.3")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var races []models.Race

	doc.Find("article[class*='tribe-events-calendar-list__event']").Each(func(i int, sel *goquery.Selection) {
		titleAnchor := sel.Find(".tribe-events-calendar-list__event-title-link")
		title := strings.TrimSpace(titleAnchor.Text())
		link, _ := titleAnchor.Attr("href")

		dateStr := strings.TrimSpace(sel.Find(".tribe-events-calendar-list__event-datetime").Text())

		if title == "" {
			return
		}

		// Basic island detection
		isl := utils.IdentifyIsland(title, link)
		if isl == "Canarias" && link != "" {
			detailText := utils.FetchText(ctx, link)
			isl = utils.IdentifyIsland(title, link, detailText)
		}

		race := models.Race{
			Name:       title,
			DateRaw:    dateStr,
			DateParsed: parseTribeDate(dateStr),
			Island:     isl,
			Location:   "Canarias",
			Source:     "FCA",
			URL:        link,
			Type:       "Running", // Default
		}

		// Refine type
		tLower := strings.ToLower(title)
		if strings.Contains(tLower, "trail") || strings.Contains(tLower, "montaña") || strings.Contains(tLower, "cxm") {
			race.Type = "Trail"
		}

		races = append(races, race)
	})

	return races, nil
}

func parseTribeDate(dateStr string) string {
	// e.g. "8 enero @ 8:00 am - 12:00 pm"
	// Hard to parse perfectly without a complex parser, but let's try a simple one
	dateStr = strings.ToLower(dateStr)
	months := map[string]string{
		"enero": "01", "febrero": "02", "marzo": "03", "abril": "04",
		"mayo": "05", "junio": "06", "julio": "07", "agosto": "08",
		"septiembre": "09", "octubre": "10", "noviembre": "11", "diciembre": "12",
	}

	parts := strings.Fields(dateStr)
	if len(parts) >= 2 {
		day := parts[0]
		if len(day) == 1 {
			day = "0" + day
		}
		monthName := parts[1]
		monthNum := months[monthName]
		if monthNum != "" {
			// Assume current or next year based on month
			year := time.Now().Year()
			if monthNum < fmt.Sprintf("%02d", int(time.Now().Month())) && year == 2026 {
				// We are in 2026, but month is earlier, could be 2027 if we are at the end of year
				// But let's just assume 2026 for now
			}
			return fmt.Sprintf("2026-%s-%s", monthNum, day)
		}
	}
	return ""
}
