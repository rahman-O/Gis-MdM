# ============================================================================
# Gis-MdM — Root Makefile
# ============================================================================
# Orchestrates all sub-projects: backend, frontend, flutter agent.

BACKEND_DIR  := serverBackendGo
FRONTEND_DIR := frontend
AGENT_DIR    := flutter_mdm_agent
ADB          := $(HOME)/Library/Android/sdk/platform-tools/adb

.PHONY: help dev stop db backend frontend agent install logs status clean

# ─── Help ────────────────────────────────────────────────────────────────────

help: ## Show this help
	@echo ""
	@echo "  Gis-MdM — Mobile Device Management Platform"
	@echo "  ════════════════════════════════════════════"
	@echo ""
	@echo "  Development:"
	@echo "    make dev          Start all services (db + backend + frontend)"
	@echo "    make stop         Stop all services"
	@echo "    make db           Start database only"
	@echo "    make backend      Start Go backend only"
	@echo "    make frontend     Start React frontend only"
	@echo ""
	@echo "  Flutter Agent:"
	@echo "    make agent        Build debug APK"
	@echo "    make agent-release Build release APK"
	@echo "    make install      Install APK on connected device"
	@echo "    make logs         Show device logs (MDM agent)"
	@echo "    make status       Check device connection & agent status"
	@echo ""
	@echo "  Database:"
	@echo "    make migrate      Run database migrations"
	@echo "    make db-reset     Reset database (WARNING: deletes all data)"
	@echo "    make db-query     Interactive psql session"
	@echo ""
	@echo "  Build & Deploy:"
	@echo "    make build        Build all (backend + frontend + agent)"
	@echo "    make lint         Lint all projects"
	@echo "    make test         Run all tests"
	@echo "    make clean        Clean build artifacts"
	@echo ""

# ─── Development ─────────────────────────────────────────────────────────────

dev: ## Start all services for development
	@echo "Starting database..."
	@cd $(BACKEND_DIR) && docker compose up -d
	@echo "Starting backend..."
	@cd $(BACKEND_DIR) && (go run ./cmd/server &)
	@echo "Starting frontend..."
	@cd $(FRONTEND_DIR) && (npm run dev &)
	@echo ""
	@echo "✓ All services starting"
	@echo "  Frontend: http://localhost:5173"
	@echo "  Backend:  http://localhost:8081"
	@echo "  Swagger:  http://localhost:8081/swagger/index.html"

stop: ## Stop all services
	@-pkill -f "go run ./cmd/server" 2>/dev/null
	@-pkill -f "vite" 2>/dev/null
	@cd $(BACKEND_DIR) && docker compose stop
	@echo "✓ All services stopped"

db: ## Start database only
	@cd $(BACKEND_DIR) && docker compose up -d

backend: ## Start Go backend only
	@cd $(BACKEND_DIR) && go run ./cmd/server

frontend: ## Start React frontend only
	@cd $(FRONTEND_DIR) && npm run dev

# ─── Flutter Agent ───────────────────────────────────────────────────────────

agent: ## Build Flutter agent (debug APK)
	@cd $(AGENT_DIR) && flutter build apk --debug
	@echo ""
	@echo "✓ APK built: $(AGENT_DIR)/build/app/outputs/flutter-apk/app-debug.apk"

agent-release: ## Build Flutter agent (release APK)
	@cd $(AGENT_DIR) && flutter build apk --release
	@echo ""
	@echo "✓ APK built: $(AGENT_DIR)/build/app/outputs/flutter-apk/app-release.apk"

install: ## Install debug APK on connected device
	@$(ADB) install -r $(AGENT_DIR)/build/app/outputs/flutter-apk/app-debug.apk
	@$(ADB) shell am start -n com.gismdm.mdm_agent/.MainActivity
	@echo "✓ Agent installed and launched"

logs: ## Show MDM agent logs from device
	@$(ADB) logcat --pid=$$($(ADB) shell pidof com.gismdm.mdm_agent) | grep -v "WifiHAL\|WifiVendorHal"

status: ## Check device & agent status
	@echo "─── Connected Devices ───"
	@$(ADB) devices
	@echo ""
	@echo "─── Agent Process ───"
	@$(ADB) shell pidof com.gismdm.mdm_agent && echo "✓ Agent is running" || echo "✗ Agent is NOT running"
	@echo ""
	@echo "─── Database Check ───"
	@docker exec serverbackendgo-db psql -U hmdm -d hmdm -c \
		"SELECT number, substring(info, 1, 100) as info FROM devices LIMIT 5;" 2>/dev/null || echo "(database not accessible)"

# ─── Database ────────────────────────────────────────────────────────────────

migrate: ## Run database migrations
	@cd $(BACKEND_DIR) && $(MAKE) migrate

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "⚠️  This will DELETE all data. Press Ctrl+C to cancel."
	@sleep 3
	@cd $(BACKEND_DIR) && docker compose down -v && docker compose up -d
	@sleep 3
	@cd $(BACKEND_DIR) && $(MAKE) migrate
	@echo "✓ Database reset complete"

db-query: ## Open interactive psql session
	@docker exec -it serverbackendgo-db psql -U hmdm -d hmdm

# ─── Build & Quality ─────────────────────────────────────────────────────────

build: ## Build all projects
	@echo "Building backend..."
	@cd $(BACKEND_DIR) && go build ./...
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm run build
	@echo "Building agent..."
	@cd $(AGENT_DIR) && flutter build apk --release
	@echo "✓ All projects built"

lint: ## Lint all projects
	@echo "Linting backend..."
	@cd $(BACKEND_DIR) && go vet ./...
	@echo "Linting frontend..."
	@cd $(FRONTEND_DIR) && npx tsc --noEmit
	@echo "Linting agent..."
	@cd $(AGENT_DIR) && dart analyze lib/
	@echo "✓ All lints passed"

test: ## Run all tests
	@echo "Testing backend..."
	@cd $(BACKEND_DIR) && go test ./... 2>&1 | tail -5
	@echo "Testing frontend..."
	@cd $(FRONTEND_DIR) && npx vitest --run 2>&1 | tail -5
	@echo "Testing agent..."
	@cd $(AGENT_DIR) && flutter test 2>&1 | tail -5

clean: ## Clean build artifacts
	@cd $(BACKEND_DIR) && rm -f server
	@cd $(FRONTEND_DIR) && rm -rf dist/
	@cd $(AGENT_DIR) && flutter clean
	@echo "✓ Cleaned"
