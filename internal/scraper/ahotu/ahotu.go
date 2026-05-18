package ahotu

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Scraper) Name() string {
	return "Ahotu"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	url := "https://www.ahotu.com/calendar/running/spain/canary-islands"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (canarias.run bot)")

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

	// Ahotu stores data in __NEXT_DATA__ script tag
	scriptTag := doc.Find("#__NEXT_DATA__")
	if scriptTag.Length() == 0 {
		log.Println("[Scraper: Ahotu] ⚠️ __NEXT_DATA__ not found, falling back to manual scrape")
		return s.scrapeManual(doc), nil
	}
	script := scriptTag.Text()

	var data struct {
		Props struct {
			PageProps struct {
				InitialState struct {
					Calendar struct {
						Events []struct {
							Name     string `json:"name"`
							Slug     string `json:"slug"`
							Date     string `json:"date"`
							Location struct {
								City    string `json:"city"`
								Country string `json:"country"`
							} `json:"location"`
							Categories []struct {
								Name string `json:"name"`
							} `json:"categories"`
						} `json:"events"`
					} `json:"calendar"`
				} `json:"initialState"`
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal([]byte(script), &data); err != nil {
		// Fallback to manual parsing if structure is different
		return s.scrapeManual(doc), nil
	}

	events := data.Props.PageProps.InitialState.Calendar.Events
	for _, e := range events {
		if e.Name == "" {
			continue
		}

		link := "https://www.ahotu.com/event/" + e.Slug
		isl := utils.IdentifyIsland(e.Name, e.Location.City, link)

		race := models.Race{
			Name:       e.Name,
			DateRaw:    e.Date,
			DateParsed: e.Date, // Ahotu uses ISO
			Island:     isl,
			Location:   e.Location.City,
			Source:     "Ahotu",
			URL:        link,
			Type:       "Running",
		}

		if len(e.Categories) > 0 {
			race.Type = e.Categories[0].Name
		}

		// For Ahotu, we MUST be certain it's in the islands.
		// If IdentifyIsland returns "Canarias", it means no specific island keyword was found.
		if isl != "Canarias" {
			races = append(races, race)
		}
	}

	if len(races) == 0 {
		log.Println("[Scraper: Ahotu] Manual fallback")
		manual := s.scrapeManual(doc)
		for _, m := range manual {
			if m.Island != "Canarias" {
				races = append(races, m)
			}
		}
	}

	return races, nil
}

func (s *Scraper) scrapeManual(doc *goquery.Document) []models.Race {
	var races []models.Race
	doc.Find("a[href^='/event/']").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find("h3").Text())
		if title == "" {
			title = strings.TrimSpace(sel.Text())
		}
		link, _ := sel.Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://www.ahotu.com" + link
		}

		if title == "" || strings.Contains(strings.ToLower(title), "ahotu") {
			return
		}

		races = append(races, models.Race{
			Name:     title,
			Island:   utils.IdentifyIsland(title, link),
			Location: "Canarias",
			Source:   "Ahotu",
			URL:      link,
			Type:     "Running",
		})
	})
	return races
}
