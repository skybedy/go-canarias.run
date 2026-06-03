SHELL := /bin/bash

MIGRATE := go run -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: run test build migrate-up migrate-down migrate-status dump-schema

run:
	GOCACHE=/tmp/go-build go run .

test:
	GOCACHE=/tmp/go-build go test ./...

build:
	GOCACHE=/tmp/go-build go build ./...

migrate-up:
	@set -a; source .env; set +a; \
	DB_URL="mysql://$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_DATABASE"; \
	$(MIGRATE) -path migrations -database "$$DB_URL" up

migrate-down:
	@set -a; source .env; set +a; \
	DB_URL="mysql://$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_DATABASE"; \
	$(MIGRATE) -path migrations -database "$$DB_URL" down 1

migrate-status:
	@set -a; source .env; set +a; \
	DB_URL="mysql://$$DB_USERNAME:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_DATABASE"; \
	$(MIGRATE) -path migrations -database "$$DB_URL" version || echo "No migrations applied yet"

dump-schema:
	@mkdir -p db
	@set -o pipefail; set -a; source .env; set +a; \
	MYSQL_PWD="$$DB_PASSWORD" mysqldump \
		-h "$$DB_HOST" -P "$$DB_PORT" -u "$$DB_USERNAME" \
		--no-data --skip-comments --skip-dump-date --single-transaction \
		"$$DB_DATABASE" races \
		| sed -E 's/AUTO_INCREMENT=[0-9]+ //g' > db/schema.sql
	@echo "Wrote db/schema.sql"
