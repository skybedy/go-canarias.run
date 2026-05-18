package calendar

import "testing"

func TestEnrichLocations_FromNameAndURL(t *testing.T) {
	in := []Race{
		{Name: "XII 5KM NOCTURNA EL ROSARIO 2026", URL: "https://inscripciones.example/inscripcion/xii-5km-nocturna-el-rosario-2026/"},
		{Name: "Transgrancanaria", URL: "https://example.com/transgrancanaria"},
		{Name: "Santa Cruz Night Run", URL: "https://example.com/santa-cruz-night-run"},
	}
	out := EnrichLocations(in)

	if out[0].Island != "Tenerife" {
		t.Fatalf("expected Tenerife, got %s", out[0].Island)
	}
	if out[1].Island != "Gran Canaria" {
		t.Fatalf("expected Gran Canaria, got %s", out[1].Island)
	}
	if out[2].City == "" || out[2].Island != "Tenerife" {
		t.Fatalf("expected Santa Cruz de Tenerife inference, got city=%q island=%q", out[2].City, out[2].Island)
	}
}

func TestEnrichLocations_DoesNotOverwriteSpecificIsland(t *testing.T) {
	in := []Race{{Name: "LPA Trail", City: "Las Palmas", Island: "Gran Canaria"}}
	out := EnrichLocations(in)
	if out[0].Island != "Gran Canaria" {
		t.Fatalf("expected existing island to stay, got %s", out[0].Island)
	}
}
