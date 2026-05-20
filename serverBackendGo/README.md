# serverBackendGo

Go backend for Headwind MDM — gradual migration from the Java [`backend/`](../backend/) REST API.

## Architecture

- **Modular Clean**: each domain lives under `internal/modules/<name>` with `domain`, `port`, `application`, and `adapter` layers.
- **Composition root**: `cmd/server` → `internal/app` wires platform (DB, HTTP) and modules.
- **API parity**: routes keep legacy prefixes (`/rest/public`, `/rest/private`, `/rest/plugins`).

See [docs/MIGRATION.md](docs/MIGRATION.md) for the phased migration order (start with **auth**).

## Quick start

```bash
cp .env.example .env

# 1) Start PostgreSQL (required for auth/login)
./scripts/db-up.sh
# Or use an existing MDM database and set DATABASE_URL in .env

# 2) Run API server
make dev
```

Auth endpoints are fully implemented (session + JWT). Swagger UI: `http://localhost:8080/swagger/index.html` when `SWAGGER_ENABLED=true`.

### Troubleshooting: `connection refused` on port 5432

The server **requires** a running PostgreSQL instance when `DATABASE_URL` is set (default in `.env.example`).

```text
database: ping database: dial tcp [::1]:5432: connect: connection refused
```

**Fix (Docker):**

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
```

**If you use the repo root stack** (`./docker-start.sh`), Postgres is often on **port 5433**. Update `.env`:

```env
DATABASE_URL=postgres://hmdm:hmdm@localhost:5433/hmdm?sslmode=disable
```

On first `make dev`, migrations in `db/migrations/` create a minimal schema and user **`admin`** / **`admin`** (API expects password as MD5 uppercase hex, same as the React frontend).

If the DB volume was created **before** migrations existed, reset and re-apply:

```bash
docker compose down -v
./scripts/db-up.sh
make dev
```

Or: `make migrate` then restart the server.

## Adding a module

1. Create `internal/modules/<name>/` following the `auth` layout.
2. Implement `module.go` (`Name`, `Register`).
3. Register in `internal/app/modules.go`.
4. Document endpoints in `docs/parity/<name>.md`.

## Layout

```
cmd/server/          Entry point
internal/app/        Bootstrap and wiring
internal/module/     Module contract (avoids import cycles)
internal/config/     Environment configuration
internal/platform/   Database, HTTP, logging
internal/modules/    Domain modules (auth first)
internal/shared/     Cross-cutting utilities
migrations/          SQL migrations (future)
docs/                Migration guides and API parity checklists
```

## Status

| Phase | Module        | Status   |
|-------|---------------|----------|
| 1     | auth          | Done |
| 1b    | signup, passwordreset | Scaffold |
| 2+    | others        | Scaffold |

Do not copy logic from [`backend-go/`](../backend-go/) — implement module-by-module from Java sources.
