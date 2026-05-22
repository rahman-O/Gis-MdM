.PHONY: dev dev-stop build test lint migrate tidy swagger db-up db-down tunnel dev-https dev-https-stop verify-https setup-studhub-https repair-apk-urls help

db-up:
	@./scripts/db-up.sh

db-down:
	docker compose down

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o internal/platform/httpx/swagger --parseDependency --parseInternal

help:
	@echo "Run from this directory (serverBackendGo/). Do not 'cd serverBackendGo' if you are already here."
	@echo ""
	@echo "HTTPS QR enrollment (two terminals):"
	@echo "  Terminal 1: make dev-stop && make dev"
	@echo "  Terminal 2: make dev-https          # keep running; updates .env BASE_URL"
	@echo "  Terminal 1: Ctrl+C, then make dev   # reload BASE_URL"
	@echo "  Terminal 2: make repair-apk-urls && make verify-https"
	@echo ""
	@echo "DuckDNS (stable URL): make setup-studhub-https HOST=mdm.studhub.duckdns.org"

dev:
	@./scripts/dev.sh

dev-stop:
	@chmod +x ./scripts/kill-dev-port.sh 2>/dev/null || true
	@./scripts/kill-dev-port.sh

tunnel:
	@./scripts/https-tunnel.sh

dev-https:
	@chmod +x ./scripts/dev-https-auto.sh ./scripts/verify-enrollment-https.sh 2>/dev/null || true
	@./scripts/dev-https-auto.sh

dev-https-stop:
	@./scripts/dev-https-auto.sh stop

verify-https:
	@chmod +x ./scripts/verify-enrollment-https.sh 2>/dev/null || true
	@./scripts/verify-enrollment-https.sh

setup-studhub-https:
	@chmod +x ./scripts/setup-studhub-https.sh 2>/dev/null || true
	@./scripts/setup-studhub-https.sh $(HOST)

repair-apk-urls:
	@chmod +x ./scripts/repair-apk-urls.sh 2>/dev/null || true
	@./scripts/repair-apk-urls.sh

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
