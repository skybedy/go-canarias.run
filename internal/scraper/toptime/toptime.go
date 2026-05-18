package toptime

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"canarias.run/internal/models"
	"canarias.run/internal/utils"
	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	URL string
}

func New() *Scraper {
	return &Scraper{
		URL: "https://toptime.live/",
	}
}

func (s *Scraper) Name() string {
	return "TopTime"
}

func (s *Scraper) Scrape(ctx context.Context) ([]models.Race, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (canarias.run bot)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
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
	dateRegex := regexp.MustCompile(`(\d{2})[/.](\d{2})[/.](\d{4})`)

	doc.Find("article.elementor-post").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find(".elementor-heading-title").Text())
		if title == "" {
			return
		}

		link, _ := sel.Find("a.elementor-button-link").First().Attr("href")
		if link == "" {
			link, _ = sel.Find("a").First().Attr("href")
		}

		// Date detection in card text
		cardText := sel.Text()
		dateMatch := dateRegex.FindString(cardText)

		isl := utils.IdentifyIsland(title, link, cardText)

		race := models.Race{
			Name:       title,
			DateRaw:    dateMatch,
			DateParsed: utils.FormatDate(dateMatch),
			Island:     isl,
			Location:   "Canarias",
			Source:     s.Name(),
			URL:        link,
			Type:       "Running",
		}

		if strings.Contains(strings.ToLower(title), "trail") {
			race.Type = "Trail"
		}

		if race.DateParsed != "" {
			races = append(races, race)
		}
	})

	return races, nil
}
