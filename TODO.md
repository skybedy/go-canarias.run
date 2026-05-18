# TODO.md

## Teď

- [ ] Dokončit rozpracované slučování více zdrojů závodů ve větvi `combine-race-sources`.
- [ ] Zkontrolovat konzistenci dat po deduplikaci (jméno, datum, ostrov, URL).
- [ ] Ověřit, že homepage korektně renderuje filtrované závody od roku 2026.

## Další kroky

- [ ] Doplnit `README.md` (lokální run, test, build, architektura scraperů).
- [ ] Přidat/rozšířit testy pro edge cases v deduplikaci a normalizaci dat.
- [ ] Ověřit build příkaz v cílovém prostředí (Ubuntu VPS).

## Později

- [ ] Zvážit přechod z `data.json` na robustnější persistentní úložiště, pokud poroste objem dat.
- [ ] Přidat jednoduchý health endpoint pro provozní monitoring.
- [ ] Zavést CI workflow pro automatický `go test ./...` a `go build ./...`.
