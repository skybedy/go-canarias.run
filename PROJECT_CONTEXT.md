# PROJECT_CONTEXT.md

## Stručný popis projektu

`canarias.run` je Go webová aplikace pro agregaci a zobrazení běžeckých závodů (zejména Kanárské ostrovy) z více zdrojů. Aplikace pravidelně scrapuje zdroje, deduplikuje data a zobrazuje je přes server-rendered HTML.

## Aktuální stav

- Projekt je aktuálně na větvi `feature/functional-web-app`.
- Proběhl refactor bez změny core scraper logiky:
  - bootstrap aplikace přesunut do `internal/app`
  - HTTP vrstva přesunuta do `internal/web`
  - `main.go` je teď tenký entrypoint
- Přidána SQLite persistence (`internal/storage/sqlite_storage.go`) s ENV přepínačem storage driveru.
- Web stále renderuje tabulku závodů a data se periodicky aktualizují přes scraper daemon.
- Ověření prochází přes:
  - `GOCACHE=/tmp/go-build go test ./...`
  - `GOCACHE=/tmp/go-build go build ./...`

## Používaný stack

- Go 1.24.0
- Standardní `net/http` + `html/template`
- `goquery` pro HTML parsing
- SQLite nebo JSON persistence (`STORAGE_DRIVER=sqlite|json`)
- Vanilla JS + vlastní CSS

## Hlavní adresáře a soubory

- `main.go`: tenký entrypoint aplikace
- `internal/app/app.go`: bootstrap aplikace, výběr storage, registrace scraperů, spuštění daemonu
- `internal/web/server.go`: HTTP routing + handlery `/` a `/healthz`
- `internal/config/config.go`: načítání a validace runtime konfigurace z ENV
- `internal/scraper/`: scrapery jednotlivých zdrojů + manager
- `internal/models/`: datové modely
- `internal/storage/`: persistence vrstva (JSON + SQLite storage)
- `templates/`: HTML šablony
- `static/`: statické assety (CSS/JS)
- `data.json`: uložená data závodů
- `docs/`: projektová dokumentace
- `README.md`: základní run/test/build instrukce a ENV proměnné
- `Makefile`: sjednocené tasky `run`, `test`, `build`

## Jak projekt spustit

- Lokálně: `go run .` nebo `make run`
- Volitelný port: env `PORT` (default `8080`)
- Storage driver:
  - `STORAGE_DRIVER=sqlite` (default), `SQLITE_DSN`
  - `STORAGE_DRIVER=json`, `DATA_FILE`
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

- Další praktický krok je doplnit uživatelské server-side filtrování tabulky (např. ostrov, měsíc, typ) přes query parametry.
- Pokud se bude dál rozšiřovat DB vrstva, přidat migrační adresář `migrations/` a tasky v `Makefile`.
- Držet malé inkrementální změny a po každé ověřovat `go test ./...` a `go build ./...` s `GOCACHE=/tmp/go-build`.
