package scraper

import (
	"testing"

	"canarias.run/internal/models"
)

func TestMergeRacesDedupesAndKeepsRicherFields(t *testing.T) {
	races := []models.Race{
		{
			ID:         "a",
			Name:       "Trail Anaga",
			DateParsed: "2026-03-15",
			Island:     "Canarias",
			Source:     "CronoCanarias",
			Type:       "running",
		},
		{
			ID:         "b",
			Name:       "Trail Anaga",
			DateParsed: "2026-03-15",
			Island:     "Tenerife",
			Location:   "La Laguna",
			Source:     "CodexUnified:cronolinecanarias",
			URL:        "https://example.com/anaga",
			Type:       "trail",
		},
		{
			ID:         "c",
			Name:       "Carrera Nocturna",
			DateParsed: "2026-03-16",
			Source:     "TopTime",
		},
	}

	merged := MergeRaces(races)

	if len(merged) != 2 {
		t.Fatalf("expected same date/name races to merge, got %d races", len(merged))
	}

	got := findRace(t, merged, "Trail Anaga")
	if got.Island != "Tenerife" {
		t.Fatalf("expected richer island, got %q", got.Island)
	}
	if got.Location != "La Laguna" {
		t.Fatalf("expected richer location, got %q", got.Location)
	}
	if got.URL != "https://example.com/anaga" {
		t.Fatalf("expected URL from incoming race, got %q", got.URL)
	}
	if got.Type != "trail" {
		t.Fatalf("expected specific race type, got %q", got.Type)
	}
	if got.Source != "CronoCanarias+CodexUnified:cronolinecanarias" {
		t.Fatalf("expected merged source labels, got %q", got.Source)
	}
}

func findRace(t *testing.T, races []models.Race, name string) models.Race {
	t.Helper()
	for _, race := range races {
		if race.Name == name {
			return race
		}
	}
	t.Fatalf("race %q not found", name)
	return models.Race{}
}
