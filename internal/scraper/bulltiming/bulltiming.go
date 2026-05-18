package bulltiming

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

// BulltimingScraper is a scraper for bulltiming.es
type BulltimingScraper struct {
	Client *http.Client
}

// New returns a new instance of BulltimingScraper
func New() *BulltimingScraper {
	return &BulltimingScraper{
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the name of the scraper
func (s *BulltimingScraper) Name() string {
	return "Bulltiming"
}

// Scrape fetches events for Bulltiming
func (s *BulltimingScraper) Scrape(ctx context.Context) ([]models.Race, error) {
	url := "https://bulltiming.es/proximos-eventos/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

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

	doc.Find(".jet-listing-grid__item").Each(func(i int, sel *goquery.Selection) {
		titleSel := sel.Find(".elementor-heading-title a").First()
		title := strings.TrimSpace(titleSel.Text())
		if title == "" {
			return
		}

		link, _ := titleSel.Attr("href")
		link = strings.TrimSpace(link)

		dateSel := sel.Find(".jet-listing-dynamic-field__content").First()
		dateStr := strings.TrimSpace(dateSel.Text())
		parsedDate := parseSpanishDate(dateStr)

		raceType := "Running"
		tLower := strings.ToLower(title)
		if strings.Contains(tLower, "trail") {
			raceType = "Trail"
		} else if strings.Contains(tLower, "mtb") || strings.Contains(tLower, "bike") || strings.Contains(tLower, "ciclismo") {
			raceType = "MTB"
		} else if strings.Contains(tLower, "nado") || strings.Contains(tLower, "travesía") {
			raceType = "Plavání"
		}

		monthShort := ""
		if !parsedDate.IsZero() {
			monthShort = strings.ToUpper(parsedDate.Format("Jan"))
		}

		isl := utils.IdentifyIsland(title, link)
		if isl == "Canarias" && link != "" {
			detailText := utils.FetchText(ctx, link)
			isl = utils.IdentifyIsland(title, link, detailText)
		}

		race := models.Race{
			Name:       title,
			DateRaw:    dateStr,
			DateParsed: parsedDate.Format("2006-01-02"),
			Month:      monthShort,
			Island:     isl,
			Location:   "Canarias", // General
			Source:     "Bulltiming",
			URL:        link,
			Type:       raceType,
		}

		if link != "" {
			races = append(races, race)
		}
	})

	return races, nil
}

func parseSpanishDate(dateStr string) time.Time {
	// "15 marzo 2026"
	dateStr = strings.TrimSpace(dateStr)

	months := map[string]string{
		"enero":      "01",
		"febrero":    "02",
		"marzo":      "03",
		"abril":      "04",
		"mayo":       "05",
		"junio":      "06",
		"julio":      "07",
		"agosto":     "08",
		"septiembre": "09",
		"octubre":    "10",
		"noviembre":  "11",
		"diciembre":  "12",
	}

	for sp, num := range months {
		if strings.Contains(strings.ToLower(dateStr), sp) {
			dPart := strings.ToLower(dateStr)
			dPart = strings.ReplaceAll(dPart, " de ", " ")
			dPart = strings.ReplaceAll(dPart, sp, num)

			// Clean multiple spaces
			dPart = strings.Join(strings.Fields(dPart), " ")

			// Might look like "15 03 2026"
			finalLayout := "2 01 2006"

			parsed, err := time.Parse(finalLayout, dPart)
			if err == nil {
				return parsed
			}
		}
	}

	return time.Time{}
}
