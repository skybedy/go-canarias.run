package sources

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"canarias.run/internal/codexcalendar"
)

const carrerasPopularesURL = "https://carreraspopularesgrancanaria.es/calendario/"

type CarrerasPopularesGCAdapter struct{}

func (CarrerasPopularesGCAdapter) Name() string { return "carreraspopularesgrancanaria" }

func (CarrerasPopularesGCAdapter) Fetch(ctx context.Context, client *http.Client) ([]calendar.Race, error) {
	html, err := fetchHTML(ctx, client, carrerasPopularesURL)
	if err != nil {
		return nil, err
	}
	return ParseCarrerasPopularesHTML(html), nil
}

var articleRe = regexp.MustCompile(`(?is)<article[^>]*>(.*?)</article>`)
var hrefRe = regexp.MustCompile(`(?is)href=["']([^"']+)["']`)
var headingRe = regexp.MustCompile(`(?is)<h[1-4][^>]*>(.*?)</h[1-4]>`)
var timeDateRe = regexp.MustCompile(`(?is)<time[^>]*datetime=["']([^"']+)["'][^>]*>`)
var freeDateRe = regexp.MustCompile(`\b(\d{1,2}[/-]\d{1,2}[/-]\d{4}|\d{4}-\d{1,2}-\d{1,2})\b`)
var trRe = regexp.MustCompile(`(?is)<tr[^>]*>(.*?)</tr>`)
var tdRe = regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)

func ParseCarrerasPopularesHTML(html string) []calendar.Race {
	out := parseCarrerasTable(html)
	out = append(out, parseCarrerasArticles(html)...)
	return calendar.MergeCalendar(out)
}

func parseCarrerasTable(html string) []calendar.Race {
	rows := trRe.FindAllStringSubmatch(html, -1)
	out := make([]calendar.Race, 0, len(rows))
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		cells := tdRe.FindAllStringSubmatch(row[1], -1)
		if len(cells) < 3 {
			continue
		}

		dateRaw := stripHTML(cells[0][1])
		dateLocal := parseDateLocal(dateRaw)
		if dateLocal == "" {
			continue
		}

		nameCell := cells[2][1]
		name := stripHTML(nameCell)
		if name == "" {
			continue
		}

		link := ""
		if href := hrefRe.FindStringSubmatch(nameCell); len(href) > 1 {
			link = href[1]
		}

		rowText := strings.ToLower(stripHTML(row[1]))
		island := inferIslandFromRowText(rowText)
		if island == "" {
			island = inferIslandFromEventText(strings.ToLower(name + " " + link + " " + rowText))
		}

		race := buildRace("carreraspopularesgrancanaria", carrerasPopularesURL, name, dateLocal, "", island, link, name, 0.90)
		out = append(out, race)
	}
	return out
}

func parseCarrerasArticles(html string) []calendar.Race {
	chunks := articleRe.FindAllStringSubmatch(html, -1)
	out := make([]calendar.Race, 0, len(chunks))

	for _, m := range chunks {
		if len(m) < 2 {
			continue
		}
		body := m[1]

		name := ""
		if h := headingRe.FindStringSubmatch(body); len(h) > 1 {
			name = stripHTML(h[1])
		}
		if name == "" {
			continue
		}

		dateRaw := ""
		if t := timeDateRe.FindStringSubmatch(body); len(t) > 1 {
			dateRaw = t[1]
		} else if d := freeDateRe.FindStringSubmatch(stripHTML(body)); len(d) > 1 {
			dateRaw = d[1]
		}
		dateLocal := parseDateLocal(dateRaw)
		if dateLocal == "" {
			continue
		}

		link := ""
		if href := hrefRe.FindStringSubmatch(body); len(href) > 1 {
			link = href[1]
		}

		city := ""
		text := strings.ToLower(stripHTML(body))
		island := inferIslandFromRowText(text)
		if island == "" {
			island = inferIslandFromEventText(strings.ToLower(name + " " + link + " " + text))
		}
		if strings.Contains(text, "las palmas") {
			city = "Las Palmas"
		}

		race := buildRace("carreraspopularesgrancanaria", carrerasPopularesURL, name, dateLocal, city, island, link, name, 0.86)
		out = append(out, race)
	}
	return out
}

func inferIslandFromRowText(text string) string {
	switch {
	case strings.Contains(text, "tenerife"):
		return "Tenerife"
	case strings.Contains(text, "gran canaria"), strings.Contains(text, "las palmas"), strings.Contains(text, "lpgc"):
		return "Gran Canaria"
	case strings.Contains(text, "lanzarote"):
		return "Lanzarote"
	case strings.Contains(text, "fuerteventura"):
		return "Fuerteventura"
	case strings.Contains(text, "la palma"):
		return "La Palma"
	case strings.Contains(text, "la gomera"):
		return "La Gomera"
	case strings.Contains(text, "el hierro"):
		return "El Hierro"
	case strings.Contains(text, "canarias"):
		return "Canarias"
	default:
		return ""
	}
}

func inferIslandFromEventText(text string) string {
	matchAny := func(needles ...string) bool {
		for _, n := range needles {
			if strings.Contains(text, n) {
				return true
			}
		}
		return false
	}

	switch {
	case matchAny(
		"tenerife", "arona", "adeje", "laguna", "realejos", "acentejo", "los abrigos",
		"granadilla", "fasnia", "tejina", "la esperanza", "sauzal", "ravelo", "tamaimo", "santa cruz",
	):
		return "Tenerife"
	case matchAny(
		"gran canaria", "las palmas", "lpgc", "teror", "tejeda", "mogan", "agaete", "arucas", "galdar",
		"telde", "agüimes", "aguimes", "santa lucia", "santa lucía", "maspalomas", "vegueta", "isleta",
	):
		return "Gran Canaria"
	case matchAny("lanzarote", "tinajo", "yaiza", "tazacorte"):
		// tazacorte is La Palma; keep specific check below.
		if strings.Contains(text, "tazacorte") {
			return "La Palma"
		}
		return "Lanzarote"
	case matchAny("fuerteventura", "puerto del rosario", "corralejo", "antigua", "morro jable", "la lajita"):
		return "Fuerteventura"
	case matchAny("la palma", "breña alta", "brena alta", "breña baja", "brena baja", "tazacorte"):
		return "La Palma"
	case matchAny("la gomera", "san sebastian de la gomera", "san sebastián de la gomera"):
		return "La Gomera"
	case matchAny("el hierro"):
		return "El Hierro"
	default:
		return ""
	}
}
