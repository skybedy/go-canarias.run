# PROJECT_CONTEXT.md

## Stručný popis projektu

`canarias.run` je Go webová aplikace pro agregaci a zobrazení běžeckých závodů (zejména Kanárské ostrovy) z více zdrojů. Aplikace pravidelně scrapuje zdroje, deduplikuje data a zobrazuje je přes server-rendered HTML.

## Aktuální stav

- Projekt je aktuálně na větvi `main`.
- Dokončen první foundation krok bez zásahu do core logiky scraperů.
- Přidána centralizovaná ENV konfigurace v `internal/config/config.go` a napojení v `main.go`.
- Přidán health endpoint `GET /healthz`.
- Přidány základní provozní soubory `README.md` a `Makefile`.
- Ověření prochází přes:
  - `GOCACHE=/tmp/go-build go test ./...`
  - `GOCACHE=/tmp/go-build go build ./...`

## Používaný stack

- Go 1.24.0
- Standardní `net/http` + `html/template`
- `goquery` pro HTML parsing
- JSON persistence (`data.json`)
- Vanilla JS + vlastní CSS

## Hlavní adresáře a soubory

- `main.go`: start aplikace, registrace scraperů, HTTP routing, napojení konfigurace
- `internal/config/config.go`: načítání a validace runtime konfigurace z ENV
- `internal/scraper/`: scrapery jednotlivých zdrojů + manager
- `internal/models/`: datové modely
- `internal/storage/`: persistence vrstva (JSON storage)
- `templates/`: HTML šablony
- `static/`: statické assety (CSS/JS)
- `data.json`: uložená data závodů
- `docs/`: projektová dokumentace
- `README.md`: základní run/test/build instrukce a ENV proměnné
- `Makefile`: sjednocené tasky `run`, `test`, `build`

## Jak projekt spustit

- Lokálně: `go run .` nebo `make run`
- Volitelný port: env `PORT` (default `8080`)
- Data soubor: env `DATA_FILE` (default `data.json`)
- Background scraping daemon:
  - `ENABLE_SCRAPER_DAEMON=true|false` (default `true`)
  - `SCRAPER_INTERVAL_HOURS` (default `24`)
- Health check endpoint: `GET /healthz`

## Jak projekt testovat

- `make test`
- nebo `GOCACHE=/tmp/go-build go test ./...`

## Jak projekt buildit

- `make build`
- nebo `GOCACHE=/tmp/go-build go build ./...`

## Známá omezení/problémy

- `Dockerfile`: nezjištěno (v repozitáři aktuálně není).
- `docker-compose.yml`: nezjištěno (v repozitáři aktuálně není).
- `package.json`: nezjištěno (v repozitáři aktuálně není).
- Některé provozní workflow (deploy, CI/CD) zatím není definováno v nalezených souborech.

## Poznámky pro další navázání

- Prioritní další krok refactoru: rozdělit `main.go` do `internal/web` a `internal/app` bez změny scraper funkcionality.
- Držet malé inkrementální změny a po každé ověřovat `go test ./...` a `go build ./...` s `GOCACHE=/tmp/go-build`.
- Pokud se mění behavior serveru nebo datový model, aktualizovat `DECISIONS.md` a `TODO.md`.
