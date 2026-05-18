# AGENTS.md

## Role Codexu v projektu

Codex je technický spolupracovník pro tento repozitář. Jeho úkol je navrhovat a implementovat změny bezpečně, s ohledem na aktuální stav kódu a s průběžnou aktualizací projektové dokumentace (`PROJECT_CONTEXT.md`, `TODO.md`, `DECISIONS.md`).

## Stack projektu (zjištěný z aktuálního stavu)

- Backend: Go (`go 1.24.0`)
- HTTP server: standardní knihovna `net/http`
- Šablony: Go `html/template`
- Frontend: Vanilla JavaScript (`static/js/app.js`), vlastní CSS (`static/css/styles.css`)
- Data storage: JSON soubor (`data.json`) přes interní storage vrstvu
- Scraping/parsing: vlastní interní scrapers + `github.com/PuerkitoBio/goquery`

## Obecné preference

- Hlavní jazyk preferuj Go, pokud projekt neurčuje jinak.
- Frontend preferuj Vanilla JavaScript.
- UI navrhuj jednoduše, čistě a prakticky.
- Pokud je potřeba stylování, preferuj Tailwind.
- Nepřidávej zbytečně složité frameworky.
- Pracuji na Linux Mint.
- Server bývá Ubuntu VPS.

## Pravidla práce

- Nejdřív vždy čti aktuální stav projektu.
- Před změnami vždy ověř aktuální adresář (`pwd`).
- Vždy zkontroluj `git status`.
- Nepředpokládej kontext ze starých chatů.
- Před větší změnou stručně popiš plán.
- Po změně spusť dostupné testy/build.
- Důležitá rozhodnutí zapisuj do `DECISIONS.md`.
- Aktuální stav zapisuj do `PROJECT_CONTEXT.md`.
- Další kroky zapisuj do `TODO.md`.
- Nepřidávej do commitu `.env` ani jiné citlivé soubory.
- Nevracej (`revert`) nesouvisející rozpracované změny.
- Nedělej destruktivní Git operace bez výslovného pokynu.
- Dělej malé, cílené změny s jasným účelem.
- Nepřidávej nové závislosti bez jasného důvodu.
- Pokud ověření po změně nejde spustit, napiš to explicitně.

## Jednotný styl převzatý z go-tene.life

- Preferovaný aplikační styl pro nové části:
  - jasně oddělené vrstvy `web/handlers`, `store/repository`, `models`, `collectors/scrapers`
  - konfigurace přes env proměnné + včasná validace
  - provozní úlohy spouštět jako samostatné příkazy (CLI mode) i periodicky
  - připravenost na DB migrace (`migrations/*.up.sql`, `*.down.sql`) a `Makefile` tasky
- Při zavádění stylu zachovej kompatibilitu s aktuálním během aplikace.

## Oprávnění a autonomní běh Codexu

- Cíl: minimalizovat ruční schvalování a zabránit zaseknutí práce při nepřítomnosti uživatele.
- V UI schvalování používej ukládání pravidel (persist/always allow) pro opakované bezpečné příkazy.
- Preferované schválené prefixy:
  - `git checkout`
  - `git branch`
  - `git merge`
  - `git commit -m`
  - `git add`
  - `go test ./...`
  - `go build ./...`
- Pokud to UI umožní, používej méně restriktivní režim schvalování pro aktuální projekt/session.
- Codex nemůže měnit globální policy sám; může pouze žádat o schválení a navrhovat prefix pravidla.
