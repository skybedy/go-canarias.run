package web

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"canarias.run/internal/app"
	"canarias.run/internal/i18n"
	"canarias.run/internal/models"
)

type Server struct {
	templates *template.Template
	app       *app.App
}

func New(a *app.App) (*Server, error) {
	funcs := template.FuncMap{
		"typeLabel": func(msgs map[string]string, raceType string) string {
			key := "type_" + raceType
			if v, ok := msgs[key]; ok {
				return v
			}
			return raceType
		},
		"formatDate": func(locale, dateParsed, dateRaw string) string {
			if dateParsed == "" || dateParsed == "0001-01-01" {
				return dateRaw
			}
			t, err := time.Parse("2006-01-02", dateParsed)
			if err != nil {
				return dateRaw
			}
			d, m, y := t.Day(), int(t.Month()), t.Year()%100
			switch locale {
			case "en":
				return fmt.Sprintf("%d/%d/%02d", m, d, y)
			case "cs":
				return fmt.Sprintf("%d.%d.%02d", d, m, y)
			default: // es
				return fmt.Sprintf("%d/%d/%02d", d, m, y)
			}
		},
	}
	tmpl, err := template.New("").Funcs(funcs).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{templates: tmpl, app: a}, nil
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/", s.handleIndex)

	addr := ":" + s.app.Config.Port
	log.Printf("Server posloucha na http://localhost:%s", s.app.Config.Port)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Zjisti locale z URL prefixu (/en, /cs) nebo použij výchozí (es)
	locale := i18n.DefaultLocale
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && i18n.IsSupportedLocale(parts[0]) {
		locale = parts[0]
		// Pouze kořenová stránka locale je platná (/en nebo /en/)
		if len(parts) > 1 && parts[1] != "" {
			http.NotFound(w, r)
			return
		}
	} else if path != "/" {
		http.NotFound(w, r)
		return
	}

	races, err := s.app.Store.GetAllRaces(context.Background())
	if err != nil {
		log.Printf("Chyba pri ziskavani zavodu: %v", err)
		races = []models.Race{}
	}

	sort.Slice(races, func(i, j int) bool {
		return races[i].DateParsed < races[j].DateParsed
	})

	seen := make(map[string]bool)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	filtered := races[:0]
	for _, race := range races {
		if race.DateParsed < "2026-01-01" {
			continue
		}
		nameKey := strings.ToLower(race.Name)
		nameKey = re.ReplaceAllString(nameKey, "")
		key := race.DateParsed + "_" + nameKey
		if seen[key] {
			continue
		}
		seen[key] = true
		filtered = append(filtered, race)
	}

	msgs := i18n.Messages(locale)
	err = s.templates.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title":     msgs["site_title"],
		"Races":     filtered,
		"Messages":  msgs,
		"Locale":    locale,
		"Languages": i18n.SupportedLanguages(),
	})
	if err != nil {
		log.Printf("Chyba pri renderovani sablony: %v", err)
		http.Error(w, "Vnitrni chyba serveru", http.StatusInternalServerError)
	}
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
