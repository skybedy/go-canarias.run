package procrono

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

// ProcronoScraper is a scraper for procrono.com (via avaibooksports)
type ProcronoScraper struct {
	Client *http.Client
}

// New returns a new instance of ProcronoScraper
func New() *ProcronoScraper {
	return &ProcronoScraper{
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the name of the scraper
func (s *ProcronoScraper) Name() string {
	return "Procrono"
}

// Scrape fetches events for Procrono
func (s *ProcronoScraper) Scrape(ctx context.Context) ([]models.Race, error) {
	// fetching from the avaibooksports iframe URL since it contains the actual upcoming events list
	url := "https://www.avaibooksports.com/sport-entity/procrono/events"
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

	doc.Find(".event-row").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find(".event-title").Text())
		if title == "" {
			return
		}

		// Example: "14 de marzo de 2026, 16:40:00 CET"
		dateStr := strings.TrimSpace(sel.Find(".event-celebration-date").Text())
		parsedDate := parseSpanishDate(dateStr)

		link, _ := sel.Find(".go-to-event").Attr("href")
		link = strings.TrimSpace(link)

		raceType := "Running"
		tLower := strings.ToLower(title)
		if strings.Contains(tLower, "trail") {
			raceType = "Trail"
		} else if strings.Contains(tLower, "mtb") || strings.Contains(tLower, "motocross") || strings.Contains(tLower, "bicicleta") {
			raceType = "MTB"
		} else if strings.Contains(tLower, "nado") || strings.Contains(tLower, "travesía") || strings.Contains(tLower, "travesia") {
			raceType = "Plavání" // Swim
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
			Location:   "Canarias",
			Source:     "Procrono",
			URL:        link,
			Type:       raceType,
		}

		// Optional rule: only keep valid dates and non-empty URLs
		if link != "" {
			races = append(races, race)
		}
	})

	return races, nil
}

func parseSpanishDate(dateStr string) time.Time {
	// Standardize to easier parsing: "14 de marzo de 2026, 16:40:00 CET"
	// Remove known timezone suffixes just in case, we'll parse locally
	dateStr = strings.Replace(dateStr, " CET", "", 1)
	dateStr = strings.Replace(dateStr, " CEST", "", 1)
	dateStr = strings.TrimSpace(dateStr)

	// Map of Spanish month names to English or numbers to assist parsing
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
			// Replace "14 de marzo de 2026, 16:40:00" -> "14 03 2026 16:40:00"
			// Split by commas first
			parts := strings.Split(dateStr, ",")
			if len(parts) >= 1 {
				// Handle date part: "14 de marzo de 2026"
				dPart := strings.TrimSpace(parts[0])
				dPart = strings.ToLower(dPart)
				dPart = strings.ReplaceAll(dPart, " de ", " ")
				dPart = strings.ReplaceAll(dPart, sp, num)

				// dPart should now look like "14 03 2026"
				tPart := "00:00:00"
				if len(parts) == 2 {
					tPart = strings.TrimSpace(parts[1])
				}

				finalLayout := "2 01 2006 15:04:05" // "14 03 2026 16:40:00"
				finalStr := fmt.Sprintf("%s %s", dPart, tPart)

				parsed, err := time.Parse(finalLayout, finalStr)
				if err == nil {
					return parsed
				}
			}
		}
	}

	return time.Time{}
}
