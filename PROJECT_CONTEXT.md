# PROJECT_CONTEXT.md

## Stručný popis projektu

`canarias.run` je Go webová aplikace pro agregaci a zobrazení běžeckých závodů (zejména Kanárské ostrovy) z více zdrojů. Aplikace pravidelně scrapuje zdroje, deduplikuje data a zobrazuje je přes server-rendered HTML.

## Aktuální stav

- Projekt je aktivně rozpracovaný na větvi `combine-race-sources`.
- V pracovním stromu je větší množství lokálních změn (tracked i untracked), včetně nových scraperů a modulu `internal/codexcalendar`.
- Aplikace má funkční testy (`go test ./...` prošlo).
- Build je funkční (`go build ./...` prošel při použití `GOCACHE=/tmp/go-build`).

## Používaný stack

- Go 1.24.0
- Standardní `net/http` + `html/template`
- `goquery` pro HTML parsing
- JSON persistence (`data.json`)
- Vanilla JS + vlastní CSS

## Hlavní adresáře a soubory

- `main.go`: start aplikace, registrace scraperů, HTTP routing
- `internal/scraper/`: scrapery jednotlivých zdrojů + manager
- `internal/codexcalendar/`: agregace/normalizace/deduplikace zdrojů
- `internal/models/`: datové modely
- `internal/storage/`: persistence vrstva (JSON storage)
- `templates/`: HTML šablony
- `static/`: statické assety (CSS/JS)
- `data.json`: uložená data závodů
- `docs/`: projektová dokumentace

## Jak projekt spustit

- Lokálně: `go run .`
- Volitelný port: env `PORT` (default `8080`)
- Aplikace po startu spustí background scraping daemon s intervalem 24h.

## Jak projekt testovat

- `go test ./...`

## Jak projekt buildit

- `go build ./...`
- V sandbox/restricted prostředí může být potřeba: `GOCACHE=/tmp/go-build go build ./...`

## Známá omezení/problémy

- `README.md`: nezjištěno (v repozitáři aktuálně není).
- `Dockerfile`: nezjištěno (v repozitáři aktuálně není).
- `docker-compose.yml`: nezjištěno (v repozitáři aktuálně není).
- `package.json`: nezjištěno (v repozitáři aktuálně není).
- Některé provozní workflow (deploy, CI/CD) zatím není definováno v nalezených souborech.

## Poznámky pro další navázání

- Před každou změnou nejdřív zkontrolovat `git status`, protože workspace obsahuje rozpracované změny napříč více moduly.
- Při práci se scrapery průběžně ověřovat deduplikaci a merge logiku (`internal/scraper`, `internal/codexcalendar`).
- Pokud se mění behavior serveru nebo datový model, aktualizovat `DECISIONS.md` a `TODO.md`.
