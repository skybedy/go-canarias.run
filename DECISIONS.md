# DECISIONS.md

## Datum

- 2026-05-18

## Důležitá technická rozhodnutí

- Zatím nejsou zaznamenána samostatná zásadní rozhodnutí nad rámec aktuální implementace v kódu.
- Pro nové změny sem zapisovat rozhodnutí, která ovlivňují architekturu, datový model, runtime behavior nebo workflow.
- Směr sjednocení stylu: projekt `go-canarias.run` se bude postupně přibližovat technickému stylu `go-tene.life` (vrstvy aplikace, env konfigurace, provozní příkazy, migrace, provozní standard).
- Přidána centralizovaná konfigurace (`internal/config`) a standardní běhové ENV (`PORT`, `DATA_FILE`, `ENABLE_SCRAPER_DAEMON`, `SCRAPER_INTERVAL_HOURS`) bez změny stávající scrape logiky.
- Přidán endpoint `GET /healthz` pro provozní kontrolu dostupnosti.
- Přidány `README.md` a `Makefile` (`run`, `test`, `build`) jako sjednocený vývojový workflow.
- Přidán `.gitignore` pro build/log artefakty (`bin/`, `*.log`), aby se provozní výstupy nemíchaly do verzovaného kódu.
- Refactor bez změny scraper logiky: bootstrap přesunut do `internal/app`, HTTP vrstva do `internal/web`, `main.go` ponechán jako tenký entrypoint.
- Zavedena SQLite persistence přes `internal/storage/sqlite_storage.go` a přepínání storage přes ENV (`STORAGE_DRIVER=sqlite|json`), s výchozím driverem `sqlite`.

## Použité technologie

- Go 1.24.0
- `net/http`, `html/template`
- `github.com/PuerkitoBio/goquery`
- SQLite (`modernc.org/sqlite`) + JSON storage (`data.json`)

## Důvody důležitých voleb

- Go + stdlib: jednoduché nasazení, nízké runtime nároky, snadný build.
- SQLite jako default: jednoduchá lokální databáze bez externího DB serveru, vhodná pro MVP produkční provoz.
- Scraper modulární po zdrojích: snadnější údržba a rozšiřování.

## Otevřené otázky

- Jak standardizovat a automatizovat validaci kvality scraped dat.
- Jak bude vypadat cílový CI/CD proces (zatím není definováno).
- Zda sjednotit i HTTP framework na Echo (`github.com/labstack/echo/v4`) nebo ponechat `net/http` a sjednotit pouze architekturu/workflow.
