package ascensotiming

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

// AscensotimingScraper is a scraper for ascensotiming.es
type AscensotimingScraper struct {
	Client *http.Client
}

// New returns a new instance of AscensotimingScraper
func New() *AscensotimingScraper {
	return &AscensotimingScraper{
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Name returns the name of the scraper
func (s *AscensotimingScraper) Name() string {
	return "Ascensotiming"
}

// Scrape fetches events for Ascensotiming
func (s *AscensotimingScraper) Scrape(ctx context.Context) ([]models.Race, error) {
	url := "https://ascensotiming.es/eventos/"
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

	// Ascensotiming places event cards inside UAGB containers
	doc.Find(".wp-block-uagb-container").Each(func(i int, sel *goquery.Selection) {
		link := ""
		title := ""
		dateStr := ""

		sel.Find("a").Each(func(j int, asel *goquery.Selection) {
			text := strings.TrimSpace(asel.Text())
			if strings.Contains(strings.ToLower(text), "inscripcione") || strings.Contains(strings.ToLower(text), "inscripción") {
				link, _ = asel.Attr("href")
			}
		})

		// They usually put the title and date in the image alt attribute
		imgAlt, _ := sel.Find("img").Attr("alt")
		imgAlt = strings.TrimSpace(imgAlt)

		if imgAlt != "" {
			parts := strings.Split(imgAlt, " - ")
			title = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				dateStr = strings.TrimSpace(parts[1])
			}
		}

		if title == "" || link == "" {
			return
		}

		parsedDate := parseSpanishDate(dateStr)

		raceType := "Running"
		tLower := strings.ToLower(title)
		if strings.Contains(tLower, "trail") {
			raceType = "Trail"
		} else if strings.Contains(tLower, "dh ") || strings.Contains(tLower, "descenso") || strings.Contains(tLower, "mtb") || strings.Contains(tLower, "bici") {
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
			Location:   "Tenerife", // Default logic
			Source:     "Ascensotiming",
			URL:        link,
			Type:       raceType,
		}

		races = append(races, race)
	})

	return races, nil
}

func parseSpanishDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	dateStr = strings.ToLower(dateStr)

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
		if strings.Contains(dateStr, sp) {
			dPart := strings.ReplaceAll(dateStr, " de ", " ")
			dPart = strings.ReplaceAll(dPart, " del ", " ")
			dPart = strings.ReplaceAll(dPart, sp, num)

			// Clean multiple spaces
			dPart = strings.Join(strings.Fields(dPart), " ")

			// e.g. "25 04 2026"
			finalLayout := "2 01 2006"

			parsed, err := time.Parse(finalLayout, dPart)
			if err == nil {
				return parsed
			}
		}
	}

	return time.Time{}
}
