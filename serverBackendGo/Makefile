.PHONY: dev build test lint migrate tidy swagger db-up db-down

db-up:
	@./scripts/db-up.sh

db-down:
	docker compose down

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o internal/platform/httpx/swagger --parseDependency --parseInternal

dev:
	@./scripts/dev.sh

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

lint:
	@test -z "$$(go vet ./... 2>&1)" || go vet ./...

migrate:
	@./scripts/migrate.sh

tidy:
	go mod tidy
