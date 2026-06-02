package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"canarias.run/internal/app"
	"canarias.run/internal/models"
)

type Server struct {
	templates *template.Template
	app       *app.App
}

func New(a *app.App) (*Server, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{templates: tmpl, app: a}, nil
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/healthz", s.handleHealthz)

	addr := ":" + s.app.Config.Port
	log.Printf("Server posloucha na http://localhost:%s", s.app.Config.Port)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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

	err = s.templates.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title": "canarias.run - Tvuj kalendar bezeckych zavodu",
		"Races": filtered,
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
