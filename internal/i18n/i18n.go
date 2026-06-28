package i18n

import "strings"

const DefaultLocale = "es"

type LanguageOption struct {
	Code string
	Name string
}

var languages = []LanguageOption{
	{Code: "es", Name: "Español"},
	{Code: "en", Name: "English"},
	{Code: "cs", Name: "Čeština"},
}

var supportedLocales = func() map[string]struct{} {
	m := make(map[string]struct{}, len(languages))
	for _, lang := range languages {
		m[lang.Code] = struct{}{}
	}
	return m
}()

var dictionary = map[string]map[string]string{
	"es": {
		"site_title":   "Canarias.run – Carreras populares en Canarias",
		"page_heading": "Carreras populares en Canarias",
		"col_date":     "Fecha",
		"col_name":     "Carrera",
		"col_location": "Lugar",
		"col_type":     "Tipo",
		"no_races":     "No hay carreras disponibles.",
		"lang_switch":  "Idioma",
		"type_trail":        "Trail",
		"type_road":         "Asfalto",
		"type_cross":        "Cross",
		"type_orienteering": "Orientación",
		"type_ocr":          "OCR",
	},
	"en": {
		"site_title":   "Canarias.run – Running races in the Canary Islands",
		"page_heading": "Running races in the Canary Islands",
		"col_date":     "Date",
		"col_name":     "Race",
		"col_location": "Location",
		"col_type":     "Type",
		"no_races":     "No races available.",
		"lang_switch":  "Language",
		"type_trail":        "Trail",
		"type_road":         "Road",
		"type_cross":        "Cross",
		"type_orienteering": "Orienteering",
		"type_ocr":          "OCR",
	},
	"cs": {
		"site_title":   "Canarias.run – Běžecké závody na Kanárských ostrovech",
		"page_heading": "Běžecké závody na Kanárských ostrovech",
		"col_date":     "Datum",
		"col_name":     "Závod",
		"col_location": "Místo",
		"col_type":     "Typ",
		"no_races":     "Žádné závody k dispozici.",
		"lang_switch":  "Jazyk",
		"type_trail":        "Trail",
		"type_road":         "Silniční běh",
		"type_cross":        "Cross",
		"type_orienteering": "Orientační běh",
		"type_ocr":          "OCR",
	},
}

func SupportedLanguages() []LanguageOption {
	out := make([]LanguageOption, len(languages))
	copy(out, languages)
	return out
}

func IsSupportedLocale(locale string) bool {
	_, ok := supportedLocales[locale]
	return ok
}

func NormalizeLocale(locale string) string {
	if IsSupportedLocale(locale) {
		return locale
	}
	return DefaultLocale
}

func Messages(locale string) map[string]string {
	loc := NormalizeLocale(locale)
	fallback := dictionary[DefaultLocale]
	chosen := dictionary[loc]

	merged := make(map[string]string, len(fallback))
	for k, v := range fallback {
		merged[k] = v
	}
	for k, v := range chosen {
		merged[k] = v
	}
	return merged
}

func T(locale, key string) string {
	msgs := Messages(locale)
	if val, ok := msgs[key]; ok {
		return val
	}
	return key
}

func LocalePrefix(locale string) string {
	loc := NormalizeLocale(locale)
	if loc == DefaultLocale {
		return ""
	}
	return "/" + loc
}

func LocaleURL(locale, path string) string {
	cleanPath := path
	if cleanPath == "" {
		cleanPath = "/"
	}
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}
	if cleanPath != "/" {
		cleanPath = strings.TrimRight(cleanPath, "/")
	}
	prefix := LocalePrefix(locale)
	if cleanPath == "/" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}
	return prefix + cleanPath
}

func StripLocalePrefix(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "/"
	}
	if IsSupportedLocale(parts[0]) {
		if len(parts) == 1 {
			return "/"
		}
		return "/" + strings.Join(parts[1:], "/")
	}
	return path
}
