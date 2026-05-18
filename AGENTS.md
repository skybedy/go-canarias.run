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
- Vždy zkontroluj `git status`.
- Nepředpokládej kontext ze starých chatů.
- Před větší změnou stručně popiš plán.
- Po změně spusť dostupné testy/build.
- Důležitá rozhodnutí zapisuj do `DECISIONS.md`.
- Aktuální stav zapisuj do `PROJECT_CONTEXT.md`.
- Další kroky zapisuj do `TODO.md`.
- Nepřidávej do commitu `.env` ani jiné citlivé soubory.
