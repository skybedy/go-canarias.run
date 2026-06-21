package cronocanarias

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

var mesesES = map[string]string{
	"enero": "01", "febrero": "02", "marzo": "03", "abril": "04",
	"mayo": "05", "junio": "06", "julio": "07", "agosto": "08",
	"septiembre": "09", "octubre": "10", "noviembre": "11", "diciembre": "12",
}

// parseDateES parsuje španělský formát "31 diciembre, 2026" na "2026-12-31"
func parseDateES(raw string) (string, string) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	raw = strings.ReplaceAll(raw, ",", "")
	parts := strings.Fields(raw)
	if len(parts) != 3 {
		return "", ""
	}
	day := parts[0]
	if len(day) == 1 {
		day = "0" + day
	}
	month, ok := mesesES[parts[1]]
	if !ok {
		return "", ""
	}
	year := parts[2]
	dateParsed := fmt.Sprintf("%s-%s-%s", year, month, day)
	t, err := time.Parse("2006-01-02", dateParsed)
	if err != nil {
		return "", ""
	}
	monthShort := strings.ToUpper(t.Format("Jan"))
	return dateParsed, monthShort
}

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

		dateParsed, monthShort := parseDateES(dateRaw)

		race := models.Race{
			Name:       title,
			DateRaw:    dateRaw,
			DateParsed: dateParsed,
			Month:      monthShort,
			Island:     isl,
			Location:   "",
			Source:     s.Name(),
			Status:     "open",
			URL:        link,
			Type:       modalityDesc,
		}

		races = append(races, race)
	})

	return races, nil
}
