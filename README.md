# go-canarias.run

Kalendář běžeckých závodů na Kanárských ostrovech postavený v Go.  
Aplikace sbírá data z více zdrojů, agreguje je a zobrazuje přes server-rendered HTML.

## Požadavky

- Go 1.24+

## Spuštění lokálně

```bash
go run .
```

Výchozí adresa: `http://localhost:8080`

## Konfigurace (ENV)

- `PORT` (default `8080`)
- `STORAGE_DRIVER` (default `sqlite`, možnosti: `sqlite|mariadb|json`)
- `SQLITE_DSN` (default `file:canarias.db?_pragma=busy_timeout(5000)`, při `STORAGE_DRIVER=sqlite`)
- `MARIADB_DSN` (při `STORAGE_DRIVER=mariadb`, např. `user:pass@tcp(host:3306)/dbname?parseTime=true`)
- `DATA_FILE` (default `data.json`, při `STORAGE_DRIVER=json`)
- `ENABLE_SCRAPER_DAEMON` (default `true`)
- `SCRAPER_INTERVAL_HOURS` (default `24`)

Příklad s MariaDB:

```bash
PORT=8090 STORAGE_DRIVER=mariadb MARIADB_DSN='user:pass@tcp(localhost:3306)/canarias?parseTime=true' go run .
```

## Migrace (MariaDB)

Migrace používají [golang-migrate](https://github.com/golang-migrate/migrate). Připrav `.env` soubor s `DB_USERNAME`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_DATABASE`.

```bash
make migrate-up      # aplikuje nové migrace
make migrate-down    # vrátí poslední migraci
make migrate-status  # zobrazí aktuální verzi migrací
make dump-schema     # exportuje schéma do db/schema.sql
```

## Health check

- `GET /healthz` vrací `200 OK` a tělo `ok`

## Testy

```bash
make test
```

nebo:

```bash
GOCACHE=/tmp/go-build go test ./...
```

## Build

```bash
make build
```

nebo:

```bash
GOCACHE=/tmp/go-build go build ./...
```
