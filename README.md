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
- `DATA_FILE` (default `data.json`)
- `ENABLE_SCRAPER_DAEMON` (default `true`)
- `SCRAPER_INTERVAL_HOURS` (default `24`)

Příklad:

```bash
PORT=8090 DATA_FILE=data.json ENABLE_SCRAPER_DAEMON=true SCRAPER_INTERVAL_HOURS=12 go run .
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
