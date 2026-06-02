# TODO.md

## Teď

- [ ] Doplnit server-side filtrování tabulky závodů (min. ostrov, měsíc, typ) přes query parametry.
- [ ] Přidat endpoint/akci pro manuální sync scraperu (např. admin-only trigger).
- [ ] Přidat základní stránkování tabulky, pokud bude dataset růst.

## Další kroky

- [ ] Přidat/rozšířit testy pro edge cases v deduplikaci a normalizaci dat.
- [ ] Ověřit build příkaz v cílovém prostředí (Ubuntu VPS).
- [ ] Rozšířit `Makefile` o další provozní tasky podle potřeby (`migrate-up/down/status`, až bude DB vrstva).
- [ ] Připravit migrační adresář `migrations/` a první SQL migraci pro tabulku `races`.
- [x] Přidat `.gitignore` pro build artefakty a logy (`bin/`, `*.log`).
- [x] Rozdělit `main.go` na `internal/app` + `internal/web` bez změny scraper funkcionality.
- [x] Přidat DB persistence (SQLite) přes storage driver (`STORAGE_DRIVER`).

## Později

- [ ] Zavést CI workflow pro automatický `go test ./...` a `go build ./...`.
- [ ] Zvážit sjednocení HTTP vrstvy na Echo (`github.com/labstack/echo/v4`) pro konzistentní styl napříč projekty.
