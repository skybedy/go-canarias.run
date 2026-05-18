# Canary Islands Running Races Sources Assessment (2026-03-04)

## Scope
Assessment of these sources for building a unified Canary Islands race calendar:

1. https://www.runnea.com/carreras-populares/calendario/canarias/ano-2026/
2. https://www.finishers.com/en/destinations/europe/spain/canary-islands-2
3. https://www.ahotu.com/calendar/spain/canary-islands
4. https://cimarunning.com/blog/calendariodecarreras/
5. https://carreraspopularesgrancanaria.es/calendario/
6. https://atletismocanarias.es/actas-de-competicion/
7. https://atletismotenerife.es/wp-content/uploads/2025/12/CALENDARIO-26.pdf
8. https://atletismocanarias.es/wp-content/uploads/2026/01/Calendario-2026-actualizado-a-31-de-diciembre.pdf
9. https://www.cronolinecanarias.com/eventos/
10. https://procrono.com/eventos/
11. https://ascensotiming.es/eventos/
12. https://conchipcanarias.com/
13. https://bulltiming.es/proximos-eventos/

## Executive Summary
- Best recurring machine-readable sources: `runnea`, `carreraspopularesgrancanaria`, `cronolinecanarias`, `bulltiming`.
- Good secondary sources (partly manual cleanup needed): `cimarunning`, `finishers`.
- Sources to avoid as primary calendar feeds:
`actas-de-competicion` (results-oriented), `procrono`/`conchip` (iframes to external platform), `ahotu` (global mixed payload/noisy), PDF-only sources (good only as fallback/manual verification).
- Build periodic import from primary sources + dedup + optional PDF/manual reconciliation.

## Per-Source Assessment

### 1) Runnea (Canarias 2026)
- Type: SSR + `window.__NUXT__` state embedded in HTML.
- Canary filter: already hard-filtered by URL path (`/canarias/ano-2026/`).
- Data quality: strong. Includes race names, dates, province/region and JSON-LD `Event` blocks.
- Extraction effort: medium (parse embedded JSON state, avoid brittle CSS scraping).
- Suitability:
  - Recurring import: `YES` (high value).
  - One-off scrape: `YES`.
- Notes: NUXT payload contains references like `/json/races/...` and `api.runnea.com` markers.

### 2) Finishers (Canary destination)
- Type: Next.js (`__NEXT_DATA__` embedded).
- Canary filter: destination URL is region-specific; still verify per-row location before insert.
- Data quality: medium-high, but global platform and schema can evolve.
- Extraction effort: medium (parse `__NEXT_DATA__` JSON, not DOM text).
- Suitability:
  - Recurring import: `YES` (secondary feed).
  - One-off scrape: `YES`.

### 3) Ahotu (Canary calendar page)
- Type: Next.js-like heavy payload.
- Canary filter: page is Canary-specific, but payload contains noisy/global-looking entries.
- Data quality: medium-low for automatic trust (needs strict geographic post-filtering).
- Extraction effort: high (payload is noisy and heavy).
- Suitability:
  - Recurring import: `NO` as primary (too noisy).
  - One-off scrape: `MAYBE` for discovery only.

### 4) Cima Running calendar
- Type: WordPress page with visible HTML tables by month.
- Canary filter: mainly Gran Canaria (explicitly labeled).
- Data quality: medium. Many placeholder dates (`xx/MM/2026`) require cleanup.
- Extraction effort: low-medium (table parsing).
- Suitability:
  - Recurring import: `YES` (secondary/manual QC).
  - One-off scrape: `YES`.

### 5) Carreras Populares Gran Canaria
- Type: WordPress; structured page with direct race detail links and dated table rows.
- Canary filter: island-specific sections + direct race pages.
- Data quality: high for Gran Canaria segment.
- Extraction effort: low-medium.
- Suitability:
  - Recurring import: `YES` (primary for Gran Canaria).
  - One-off scrape: `YES`.

### 6) Atletismo Canarias - Actas de Competicion
- Type: WordPress page listing many PDF actas/results.
- Canary filter: federation scope, but content type is mostly results and official records.
- Data quality for race calendar: low (not a clean upcoming-race feed).
- Extraction effort: medium-high for little scheduling value.
- Suitability:
  - Recurring import: `NO` (not a calendar source).
  - One-off scrape: `NO` for calendar ingestion.

### 7) Atletismo Tenerife PDF (CALENDARIO-26.pdf)
- Type: PDF only.
- Canary filter: Tenerife/Canarias federation context.
- Data quality: useful reference; text extraction works but tabular reconstruction is imperfect.
- Extraction effort: high for reliable automation.
- Suitability:
  - Recurring import: `NO` as primary automation.
  - One-off scrape/import: `YES` (bootstrap or reconciliation).

### 8) Atletismo Canarias PDF (Calendario-2026 actualizado...)
- Type: PDF only.
- Canary filter: federation calendar.
- Data quality: useful reference/checkpoint, but same PDF limitations.
- Extraction effort: high.
- Suitability:
  - Recurring import: `NO` as primary automation.
  - One-off scrape/import: `YES` (manual QA support).

### 9) Cronoline Canarias
- Type: WordPress with custom post type `event` rendered in HTML list.
- Canary filter: local operator; includes modality and event details.
- Data quality: high for their managed events.
- Extraction effort: low.
- Suitability:
  - Recurring import: `YES` (primary timing-provider feed).
  - One-off scrape: `YES`.

### 10) Procrono
- Type: WordPress shell embedding external iframe:
  - `https://www.avaibooksports.com/sport-entity/procrono/events/rankings`
- Canary filter: likely good inside provider platform.
- Data quality: unknown unless scraping iframe source directly.
- Extraction effort: medium-high (must integrate external platform).
- Suitability:
  - Recurring import: `MAYBE` only via direct `avaibooksports` endpoint.
  - One-off scrape: `MAYBE`.

### 11) Ascenso Timing
- Type: WordPress marketing page + links to registrations (`inscripciones.ascensotiming.es/...`).
- Canary filter: local, good context.
- Data quality: mixed (not always a clean full machine list on this page).
- Extraction effort: medium (likely better to target dedicated inscription domain endpoints).
- Suitability:
  - Recurring import: `MAYBE` (needs dedicated endpoint/source discovery).
  - One-off scrape: `YES`.

### 12) Conchip Canarias
- Type: WordPress page embedding iframe:
  - `https://inscripciones.conchipcanarias.com/sport-entity/conchipcanarias/events`
- Canary filter: local operator.
- Data quality: likely good in iframe source, weak in host page itself.
- Extraction effort: medium-high (integrate iframe source directly).
- Suitability:
  - Recurring import: `MAYBE` via iframe source.
  - One-off scrape: `MAYBE`.

### 13) Bulltiming
- Type: WordPress + JetEngine/JetSmartFilter configuration exposed in HTML.
- Canary filter: local operator.
- Data quality: good. Query config reveals post type (`glive_events`) and date meta key (`fecha-evento`).
- Extraction effort: medium (can parse HTML now; later optimize with WP/Jet endpoints if stable).
- Suitability:
  - Recurring import: `YES` (primary timing-provider feed).
  - One-off scrape: `YES`.

## Recommended Source Tiers

### Tier A (Primary recurring import)
- Runnea
- Carreras Populares Gran Canaria
- Cronoline Canarias
- Bulltiming

### Tier B (Secondary recurring import, stricter validation)
- Finishers
- Cima Running
- Ascenso Timing (pending better endpoint discovery)
- Procrono/Conchip only if direct external provider endpoints are stabilized

### Tier C (Reference/manual reconciliation)
- Atletismo Tenerife PDF
- Atletismo Canarias PDF
- Atletismo Canarias Actas (results verification, not calendar feed)
- Ahotu (discovery only, heavy cleanup needed)

## Practical Go Implementation Strategy

1. Add source adapters with stable extraction contracts:
   - `Fetch(ctx) ([]models.Race, error)`
2. For JS-heavy pages, parse embedded JSON state first (`__NUXT__`, `__NEXT_DATA__`) instead of CSS scraping.
3. Normalize and deduplicate by:
   - `date_local + normalized_name + island_or_city`
4. Store source attribution and confidence score per race:
   - `source`, `source_url`, `source_last_seen`, `confidence`.
5. Run imports daily + keep last successful snapshot for diffing.
6. Treat PDF ingestion as manual/ops task (not nightly primary pipeline).

## Decision: One-Off vs Recurring
- Recurring recommended: sources 1, 5, 9, 13 (+ optionally 2, 4).
- One-off or manual-only: sources 7, 8.
- Not recommended as primary recurring calendar feeds: 3, 6, 10, 12 (unless direct iframe-provider API/page integration is implemented and validated).
