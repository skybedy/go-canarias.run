# TODO.md

## Teď

- [ ] Rozdělit `main.go` na menší celky (`internal/web`, `internal/app`) při zachování stávající funkčnosti scraperů.
- [ ] Přesunout HTTP routing a handlery z `main.go` do `internal/web`.
- [ ] Přesunout bootstrap aplikace (config, storage, registrace scraper manageru) do `internal/app`.
- [ ] Po refactoru ověřit beze změny chování endpointů `/` a `/healthz`.

## Další kroky

- [ ] Přidat/rozšířit testy pro edge cases v deduplikaci a normalizaci dat.
- [ ] Ověřit build příkaz v cílovém prostředí (Ubuntu VPS).
- [ ] Rozšířit `Makefile` o další provozní tasky podle potřeby (`migrate-up/down/status`, až bude DB vrstva).
- [ ] Připravit cílový DB návrh + migrační adresář `migrations/` pro odchod od `data.json`.
- [x] Přidat `.gitignore` pro build artefakty a logy (`bin/`, `*.log`).

## Později

- [ ] Zvážit přechod z `data.json` na robustnější persistentní úložiště, pokud poroste objem dat.
- [ ] Zavést CI workflow pro automatický `go test ./...` a `go build ./...`.
- [ ] Zvážit sjednocení HTTP vrstvy na Echo (`github.com/labstack/echo/v4`) pro konzistentní styl napříč projekty.
