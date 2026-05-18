package calendar

import "strings"

type locationHint struct {
	Needles []string
	City    string
	Island  string
}

var locationHints = []locationHint{
	{Needles: []string{"santa cruz"}, City: "Santa Cruz de Tenerife", Island: "Tenerife"},
	{Needles: []string{"la laguna", "san cristobal"}, City: "La Laguna", Island: "Tenerife"},
	{Needles: []string{"realejos"}, City: "Los Realejos", Island: "Tenerife"},
	{Needles: []string{"matanza"}, City: "La Matanza de Acentejo", Island: "Tenerife"},
	{Needles: []string{"abrigos"}, City: "Los Abrigos", Island: "Tenerife"},
	{Needles: []string{"granadilla"}, City: "Granadilla de Abona", Island: "Tenerife"},
	{Needles: []string{"fasnia"}, City: "Fasnia", Island: "Tenerife"},
	{Needles: []string{"sauzal", "ravelo"}, City: "El Sauzal", Island: "Tenerife"},
	{Needles: []string{"tamaimo"}, City: "Santiago del Teide", Island: "Tenerife"},
	{Needles: []string{"el rosario"}, City: "El Rosario", Island: "Tenerife"},
	{Needles: []string{"maretas"}, City: "Arico", Island: "Tenerife"},
	{Needles: []string{"carboneras"}, City: "Anaga", Island: "Tenerife"},
	{Needles: []string{"tenerife"}, City: "", Island: "Tenerife"},

	{Needles: []string{"las palmas", "lpa", "lpgc", "vegueta", "sebadal", "canteras"}, City: "Las Palmas", Island: "Gran Canaria"},
	{Needles: []string{"maspalomas"}, City: "Maspalomas", Island: "Gran Canaria"},
	{Needles: []string{"telde", "el goro"}, City: "Telde", Island: "Gran Canaria"},
	{Needles: []string{"arucas"}, City: "Arucas", Island: "Gran Canaria"},
	{Needles: []string{"agaete"}, City: "Agaete", Island: "Gran Canaria"},
	{Needles: []string{"galdar", "gáldar"}, City: "Galdar", Island: "Gran Canaria"},
	{Needles: []string{"tejeda"}, City: "Tejeda", Island: "Gran Canaria"},
	{Needles: []string{"teror"}, City: "Teror", Island: "Gran Canaria"},
	{Needles: []string{"mogan", "mogan"}, City: "Mogan", Island: "Gran Canaria"},
	{Needles: []string{"moya"}, City: "Moya", Island: "Gran Canaria"},
	{Needles: []string{"artenara"}, City: "Artenara", Island: "Gran Canaria"},
	{Needles: []string{"saucillo"}, City: "Galdar", Island: "Gran Canaria"},
	{Needles: []string{"transgrancanaria", "grancanaria", "gran canaria"}, City: "", Island: "Gran Canaria"},

	{Needles: []string{"lanzarote"}, City: "", Island: "Lanzarote"},
	{Needles: []string{"fuerteventura"}, City: "", Island: "Fuerteventura"},
	{Needles: []string{"breña", "brena", "tazacorte", "la palma"}, City: "", Island: "La Palma"},
}

// EnrichLocations heuristically fills missing city/island from race name and URL.
func EnrichLocations(in []Race) []Race {
	out := make([]Race, len(in))
	copy(out, in)

	for i := range out {
		r := out[i]
		text := strings.ToLower(strings.Join([]string{r.Name, r.URL, r.City, r.Island}, " "))
		for _, h := range locationHints {
			if !containsAny(text, h.Needles...) {
				continue
			}
			if (r.Island == "" || strings.EqualFold(r.Island, "canarias")) && h.Island != "" {
				r.Island = h.Island
			}
			if r.City == "" && h.City != "" {
				r.City = h.City
			}
			if r.City != "" && r.Island != "" && !strings.EqualFold(r.Island, "canarias") {
				break
			}
		}
		out[i] = r
	}
	return out
}

func containsAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, strings.ToLower(n)) {
			return true
		}
	}
	return false
}
