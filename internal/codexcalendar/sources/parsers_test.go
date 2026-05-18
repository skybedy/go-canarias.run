package sources

import "testing"

func TestParseRunneaHTML(t *testing.T) {
	html := `
	<html><body>
	<script>
	window.__NUXT__={"data":[{"name":"10K Arona","date":"2026-05-11","city":"Arona","island":"Tenerife","url":"/carreras/10k-arona"}]};
	</script>
	</body></html>`

	races := ParseRunneaHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].DateLocal != "2026-05-11" {
		t.Fatalf("unexpected date: %s", races[0].DateLocal)
	}
	if races[0].Sources[0].Source != "runnea" {
		t.Fatalf("missing source attribution")
	}
}

func TestParseCarrerasPopularesHTML(t *testing.T) {
	html := `
	<figure class="wp-block-table"><table><tbody>
	<tr>
	  <td class="has-text-align-center">01/06/2026</td>
	  <td class="has-text-align-center">Asfalto</td>
	  <td class="has-text-align-center"><a href="https://example.com/race-a">Carrera Popular LPGC</a></td>
	  <td class="has-text-align-center">Confirmada</td>
	</tr>
	</tbody></table></figure>`

	races := ParseCarrerasPopularesHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].Island == "" {
		t.Fatalf("expected island to be set")
	}
	if races[0].URL != "https://example.com/race-a" {
		t.Fatalf("expected parsed detail URL, got %s", races[0].URL)
	}
}

func TestParseCronolineHTML(t *testing.T) {
	html := `
	<article class="type-event">
	  <h3 class="entry-title"><a href="/event/trail-adeje">Trail Adeje Tenerife</a></h3>
	  <span class="start_date">2026-07-09</span>
	</article>`

	races := ParseCronolineHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].Island != "Tenerife" {
		t.Fatalf("expected Tenerife island, got %s", races[0].Island)
	}
}

func TestParseProcronoEventsHTML(t *testing.T) {
	html := `
	<div class="event-row">
	  <p class="event-celebration-date">14 de marzo de 2026, 16:40:00 CET</p>
	  <p class="event-title">XXIII MILLA BREÑA BAJA MÁGICA</p>
	  <a class="go-to-event registrations-enabled" href="https://inscripciones.procrono.com/inscripcion/xxiii-milla-brena-baja-magica/"></a>
	</div>`
	races := ParseProcronoEventsHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].DateLocal != "2026-03-14" {
		t.Fatalf("expected parsed date, got %s", races[0].DateLocal)
	}
}

func TestParseAscensoTimingHTML(t *testing.T) {
	html := `
	<div class="wp-block-uagb-container">
	  <img alt="Asomadero Trail - 25 de Abril de 2026 - Ascenso Timing" />
	  <a href="https://inscripciones.ascensotiming.es/inscripcion/asomadero-trail-2026/">Inscripciones</a>
	</div>`
	races := ParseAscensoTimingHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].Name != "Asomadero Trail" {
		t.Fatalf("expected parsed name, got %s", races[0].Name)
	}
	if races[0].DateLocal != "2026-04-25" {
		t.Fatalf("expected parsed date, got %s", races[0].DateLocal)
	}
}

func TestParseBullTimingHTML(t *testing.T) {
	html := `
	<div class="jet-listing-grid__item">
	  <div class="elementor-heading-title"><a href="https://bulltiming.es/eventos/2026-lpa-trail/">LPA Trail</a></div>
	  <div class="jet-listing-dynamic-field__content">15 marzo 2026</div>
	</div>`

	races := ParseBullTimingHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].DateLocal != "2026-03-15" {
		t.Fatalf("expected parsed spanish short date, got %s", races[0].DateLocal)
	}
}

func TestParseConchipEventsHTML(t *testing.T) {
	html := `
	<div class="event-row">
	  <p class="event-celebration-date">1 de enero de 2026, 9:00:00 CET</p>
	  <p class="event-title">III CIRCUITO DE CARRERAS X MONTAÑA ISLA DE TENERIFE 2026 - TROFEO MIZUNO</p>
	  <a class="go-to-event registrations-enabled" href="https://inscripciones.conchipcanarias.com/inscripcion/iii-circuito/"></a>
	</div>`

	races := ParseConchipEventsHTML(html)
	if len(races) != 1 {
		t.Fatalf("expected 1 race, got %d", len(races))
	}
	if races[0].DateLocal != "2026-01-01" {
		t.Fatalf("expected parsed Spanish long date, got %s", races[0].DateLocal)
	}
	if races[0].Island != "Tenerife" {
		t.Fatalf("expected Tenerife inference, got %s", races[0].Island)
	}
}
