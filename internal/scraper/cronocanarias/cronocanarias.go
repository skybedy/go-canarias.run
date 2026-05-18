package cronocanarias

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"canarias.run/internal/models"
	"canarias.run/internal/utils"
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

	// Hledáme element <article>, který obsahuje třídu event
	doc.Find("article.type-event").Each(func(i int, sel *goquery.Selection) {
		titleSel := sel.Find(".entry-title a")
		title := strings.TrimSpace(titleSel.Text())
		link, _ := titleSel.Attr("href")
		dateRaw := strings.TrimSpace(sel.Find(".start_date").Text())

		if title == "" {
			return // Pokud nenajdeme nadpis, pravděpodobně nejde o platný závod
		}

		// Pokus o zjištění typu (running, mtb, atd.) z ikonek
		modalityDesc := "trail/asphalt" // výchozí
		modalityImg := sel.Find(".modality img.icon-madality")
		if modalityImg.Length() > 0 {
			modTitle, _ := modalityImg.Attr("title")
			if modTitle != "" {
				modalityDesc = strings.ToLower(modTitle)
			}
		}

		isl := utils.IdentifyIsland(title, link)
		if isl == "Canarias" && link != "" {
			detailText := utils.FetchText(ctx, link)
			isl = utils.IdentifyIsland(title, link, detailText)
		}

		race := models.Race{
			ID:        fmt.Sprintf("cronoline_%d", i),
			Name:      title,
			DateRaw:   dateRaw,
			Month:     "TBD", // Bude zpracováno společným datumových parserem
			Island:    isl,
			Location:  "",
			Distances: []string{}, // Nelze získat bez detailu
			Source:    s.Name(),
			Status:    "open",
			URL:       link,
			Type:      modalityDesc,
		}

		races = append(races, race)
	})

	return races, nil
}
