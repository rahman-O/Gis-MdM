# Implementation Plan: Phase 3 — Hints Module Migration

**Branch**: `003-complete-hints-module` | **Date**: 2026-05-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/003-complete-hints-module/spec.md`

## Summary

Replace the **hints** scaffold in `serverBackendGo` with full parity for Java
`HintResource` and `UserDAO` hint persistence: load per-user hint history, mark keys
as shown, enable (clear history), and disable (mark all catalog keys). Add Postgres
migration for `userHints` / `userHintTypes`, clean-architecture module wiring,
Swagger `BearerAuth`, application tests, and `docs/parity/hints.md`. Reuse existing
JWT/session middleware and Headwind response envelope.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, existing `platform/auth`, `platform/httpx/response`,
`swaggo/gin-swagger` (Bearer auth already configured in `cmd/server/main.go`)

**Storage**: PostgreSQL (`DATABASE_URL`); new migration `000004_hints_tables.up.sql`
(tables missing from `000001_init.up.sql`)

**Testing**: `go test ./internal/modules/hints/...`; `application/*_test.go` with stub repo;
optional `adapter/http/handler_test.go` with JWT or session principal

**Target Platform**: Linux/macOS dev (`:8080`); Docker Postgres via `scripts/db-up.sh`

**Project Type**: Web service + React consumer (`frontend/src/features/hints/`)

**Performance Goals**: History GET &lt; 100ms p95 on seed DB; single-user row ops only

**Constraints**: Four endpoints only; current-user scoping via `Principal.ID`; no new
hint catalog management API; `POST /history` body compatible with Angular `$resource`

**Scale/Scope**: 4 endpoints; 2 tables; 1 module; ~6–8 implementation files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `hints` module; Phase 3 in `MIGRATION.md` / `NEXT_STEPS.md` #3 |
| **II. Layered Clean** | ✅ | domain/port/application/adapter; replace scaffold |
| **III. API Parity** | ✅ | `contracts/hints-api.md` + `docs/parity/hints.md` |
| **IV. Testable Delivery** | ✅ | `go test` + `quickstart.md` + Swagger |
| **V. Simplicity** | ✅ | Direct SQL mirroring MyBatis annotations; no shared user repo coupling |
| **VI. Security** | ✅ | `RequireAuth` + principal user id only |
| **VII. Observability** | ✅ | `error.internal.server` on failures; slog in handlers |

**Post-design**: All gates remain ✅. No Complexity Tracking entries required.

## Project Structure

### Documentation (this feature)

```text
specs/003-complete-hints-module/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1 smoke + Swagger
├── contracts/
│   └── hints-api.md
└── tasks.md             # (/speckit-tasks — not created by plan)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   └── 000004_hints_tables.up.sql       # userHints + userHintTypes + seed keys
├── docs/parity/
│   └── hints.md                         # NEW — all HintResource endpoints Done
├── internal/modules/hints/
│   ├── module.go                        # wire repo + service + http
│   ├── domain/
│   │   └── hint.go                      # HintKey string alias, validation
│   ├── port/
│   │   └── repository.go
│   ├── application/
│   │   ├── service.go                   # GetHistory, MarkShown, Enable, Disable
│   │   └── service_test.go
│   └── adapter/
│       ├── http/
│       │   ├── handler.go               # + Swagger @Security BearerAuth
│       │   └── handler_test.go          # optional
│       └── persistence/postgres/
│           └── hint_repo.go             # SQL from UserMapper annotations
└── docs/NEXT_STEPS.md                   # mark hints row منجز after implement
```

**Structure Decision**: Single bounded module under `internal/modules/hints/`. No
dependency on `users` module — only `platformauth.PrincipalFromContext` for `userId`.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Database migration

1. Add `000004_hints_tables.up.sql` with `userHints`, `userHintTypes`, unique
   `(userid, hintkey)`, FK cascade to `users`, seed `hint.step.1` … `hint.step.4`.
2. Verify `database.Migrate` applies on fresh `db-up`.

### Phase B — Persistence & application

| SQL (Java) | Use case |
|------------|----------|
| `SELECT hintKey FROM userHints WHERE userId = $1` | `GetHistory` |
| `INSERT INTO userHints ... ON CONFLICT DO NOTHING` | `MarkShown` (idempotent) |
| `DELETE FROM userHints WHERE userId = $1` | `Enable` |
| `DELETE` + `INSERT ... SELECT hintKey FROM userHintTypes` | `Disable` |

### Phase C — HTTP & module wiring

| Endpoint | Handler |
|----------|---------|
| `GET /history` | `GetHistory` → `OK` + `string[]` |
| `POST /history` | `MarkShown` — parse body (JSON string or raw) |
| `POST /enable` | `Enable` → `OK` |
| `POST /disable` | `Disable` → `OK` |

Replace `module.go` scaffold with real `Register` on `groups.Private.Group("/hints")`.

### Phase D — Swagger, parity, docs

1. `@Security BearerAuth` on all four handlers (match Phase 2 pattern).
2. `make swagger`.
3. Fill `docs/parity/hints.md`.
4. Update `NEXT_STEPS.md` row #3 to **منجز**.

### Phase E — Tests

- Stub repo tests: enable clears, disable inserts catalog count, duplicate mark-shown.
- Handler test: authenticated GET `/history` returns 200 envelope.

## Complexity Tracking

> No constitution violations requiring justification.
