include .env
BIN = reader
NOW = $(shell date '+%Y_%m_%d_%H_%M_%S')

.PHONY: all tailwind run sqlite migrate migrate-down backup nlp nlp-debug

all: tailwind templ
	@go build -o bin/$(BIN) cmd/web/*

tailwind:
	@$(TAILWIND) -i app.css -o ./static/tailwind-output.css --minify

templ:
	@go tool templ generate

run: all
	@./bin/$(BIN)


# DB ##################################################
sqlite:
	@sqlite3 $(DB)

migrate:
	migrate -database "sqlite3://$(DB)" -path ./migrations up

migrate-down:
	migrate -database "sqlite3://$(DB)" -path ./migrations down 1

backup:
	sqlite3 $(DB) ".backup $(DB_BACKUP_DIR)/$(NOW).db"

sqlc:
	go tool sqlc generate

######################################################

nlp:
	cd nlp && uv run fastapi run

nlp-debug:
	cd nlp && uv run fastapi dev
