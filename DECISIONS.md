# DECISIONS.md

## Datum

- 2026-05-18

## Důležitá technická rozhodnutí

- Zatím nejsou zaznamenána samostatná zásadní rozhodnutí nad rámec aktuální implementace v kódu.
- Pro nové změny sem zapisovat rozhodnutí, která ovlivňují architekturu, datový model, runtime behavior nebo workflow.

## Použité technologie

- Go 1.24.0
- `net/http`, `html/template`
- `github.com/PuerkitoBio/goquery`
- JSON storage (`data.json`)

## Důvody důležitých voleb

- Go + stdlib: jednoduché nasazení, nízké runtime nároky, snadný build.
- JSON storage: rychlý start projektu bez externí DB závislosti.
- Scraper modulární po zdrojích: snadnější údržba a rozšiřování.

## Otevřené otázky

- Zda a kdy přejít z JSON persistence na databázi.
- Jak standardizovat a automatizovat validaci kvality scraped dat.
- Jak bude vypadat cílový CI/CD proces (zatím není definováno).
