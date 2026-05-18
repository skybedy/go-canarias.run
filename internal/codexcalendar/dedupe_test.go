package calendar

import (
	"testing"
	"time"
)

func TestMergeCalendar_DedupesAndMergesSources(t *testing.T) {
	now := time.Date(2026, 3, 4, 12, 0, 0, 0, time.UTC)
	r1 := Race{
		Name:       "Media Maraton Las Palmas",
		DateLocal:  "2026-04-10",
		City:       "Las Palmas",
		Island:     "Gran Canaria",
		URL:        "https://a.example/r1",
		Confidence: 0.80,
		FirstSeen:  now,
		LastSeen:   now,
		Sources: []SourceAttribution{{
			Source: "runnea", URL: "https://a.example/r1", LastSeen: now, Confidence: 0.80,
		}},
	}
	r2 := Race{
		Name:       "Media-Maraton Las  Palmas",
		DateLocal:  "2026-04-10",
		City:       "Las Palmas",
		Island:     "Gran Canaria",
		URL:        "https://b.example/race",
		Confidence: 0.92,
		FirstSeen:  now.Add(1 * time.Hour),
		LastSeen:   now.Add(1 * time.Hour),
		Sources: []SourceAttribution{{
			Source: "cronolinecanarias", URL: "https://b.example/race", LastSeen: now.Add(1 * time.Hour), Confidence: 0.92,
		}},
	}

	merged := MergeCalendar([]Race{r1, r2})
	if len(merged) != 1 {
		t.Fatalf("expected 1 merged race, got %d", len(merged))
	}
	got := merged[0]
	if got.Confidence != 0.92 {
		t.Fatalf("expected confidence from stronger source, got %v", got.Confidence)
	}
	if got.URL != "https://b.example/race" {
		t.Fatalf("expected URL from stronger source, got %s", got.URL)
	}
	if len(got.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(got.Sources))
	}
}

func TestMergeRaces_StrongerSourceOverridesLocation(t *testing.T) {
	base := Race{
		Name:       "Race A",
		DateLocal:  "2026-05-01",
		City:       "Las Palmas",
		Island:     "Gran Canaria",
		Confidence: 0.80,
		URL:        "https://a.example/r",
	}
	stronger := Race{
		Name:       "Race A",
		DateLocal:  "2026-05-01",
		City:       "Santa Cruz de Tenerife",
		Island:     "Tenerife",
		Confidence: 0.95,
		URL:        "https://b.example/r",
	}

	got := MergeRaces(base, stronger)
	if got.Island != "Tenerife" {
		t.Fatalf("expected island from stronger source, got %s", got.Island)
	}
	if got.City != "Santa Cruz de Tenerife" {
		t.Fatalf("expected city from stronger source, got %s", got.City)
	}
}
