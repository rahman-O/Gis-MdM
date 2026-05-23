# Implementation Plan: Phase 3 — Customers Module Migration

**Branch**: `004-complete-customers-module` | **Date**: 2026-05-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/004-complete-customers-module/spec.md`

## Summary

Replace the **customers** scaffold in `serverBackendGo` with super-admin parity for Java
`CustomerResource`: paginated search and impersonation (React control panel P1), plus
CRUD/prefix-check endpoints for full resource coverage. Add migration `000005_customers_extend`,
clean-architecture module (`domain` / `port` / `application` / `adapter`), `RequireSuperAdmin`
checks, Swagger `BearerAuth`, tests, and `docs/parity/customers.md`. Customer **create** omits
default devices and configuration copy until Phase 4/5 schema exists (documented Partial).

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth`, `platform/httpx/response`,
`internal/platform/crypto` (tokens/passwords), `swaggo/gin-swagger`

**Storage**: PostgreSQL; migration `000005_customers_extend.up.sql` for Liquibase-aligned
customer columns; reuse existing `customers` / `users` from `000001_init`

**Testing**: `go test ./internal/modules/customers/...`; application stubs + handler tests with
super-admin principal

**Target Platform**: Linux/macOS dev (`:8080`); Docker Postgres via `scripts/db-up.sh`

**Project Type**: Web service + React consumer (`frontend/src/features/customers/`)

**Performance Goals**: Search POST &lt; 500ms p95 on seed DB with 100+ customers

**Constraints**: Super-admin only; Headwind `PaginatedData` shape; no Mailchimp; deprecated GET
search not implemented; impersonation returns JWT-friendly payload

**Scale/Scope**: 6 active endpoints; 1 module; ~12–15 Go files; completes Phase 3 roadmap row

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `customers` module; closes Phase 3 in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | domain/port/application/adapter; replace scaffold |
| **III. API Parity** | ✅ | `contracts/customers-api.md` + `docs/parity/customers.md` |
| **IV. Testable Delivery** | ✅ | `go test` + `quickstart.md` + Swagger |
| **V. Simplicity** | ✅ | SQL mirrors `CustomerMapper.xml`; localized user lookup in adapter |
| **VI. Security** | ✅ | Super-admin gate; no password in impersonate response |
| **VII. Observability** | ✅ | Stable `message` keys; slog on errors |

**Post-design**: All gates ✅. Partial create parity (devices/config) documented in research R5 —
not a layer violation, scope deferral per migration phases.

## Project Structure

### Documentation (this feature)

```text
specs/004-complete-customers-module/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1 smoke
├── contracts/
│   └── customers-api.md
└── tasks.md             # (/speckit-tasks — not created by plan)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   └── 000005_customers_extend.up.sql    # email, accounttype, customerstatus, times, etc.
├── docs/parity/
│   └── customers.md                      # NEW — endpoint matrix
├── internal/platform/auth/
│   └── superadmin.go                     # NEW — RequireSuperAdmin helper (optional file)
├── internal/modules/customers/
│   ├── module.go                         # wire repo + service + http
│   ├── domain/
│   │   ├── customer.go                   # Customer, SearchRequest, Paginated
│   │   └── transliterate.go              # port of CustomerDAO.transliterate
│   ├── port/
│   │   ├── customer_repository.go
│   │   └── user_lookup.go                # org admin + token update for impersonate
│   ├── application/
│   │   ├── service.go                    # Search, Impersonate, Save, Delete, Edit, Prefix
│   │   └── service_test.go
│   └── adapter/
│       ├── http/
│       │   ├── handler.go                # + Swagger @Security BearerAuth
│       │   └── handler_test.go
│       └── persistence/postgres/
│           ├── customer_repo.go          # search/count/CRUD/prefix
│           └── user_lookup.go            # findOrgAdmin, ensureAuthToken
└── docs/
    ├── MIGRATION.md                      # Phase 3 customers → done
    └── NEXT_STEPS.md                     # customers row منجز; next devices
```

**Structure Decision**: Single bounded module. No import of `users`/`auth` application layers;
duplicate minimal SQL in customers postgres adapter per research R7.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Database migration

1. `000005_customers_extend.up.sql` — add columns from Liquibase subset (see `data-model.md`).
2. `.down.sql` for rollback.
3. Seed: ensure at least one non-master customer + org admin for impersonate smoke.

### Phase B — P1: Search + impersonate

| Endpoint | Application | Notes |
|----------|-------------|-------|
| `POST /search` | `Search` | Build dynamic SQL from `CustomerMapper.xml` |
| `GET /impersonate/:id` | `Impersonate` | Org admin role 2; token generation; block reset token |

HTTP: `RequireAuth` + `EnrichPrincipal` + super-admin check.

### Phase C — P2: CRUD + prefix

| Endpoint | Application |
|----------|-------------|
| `GET /:id/edit` | `GetForEdit` |
| `PUT /` | `Save` (create/update, duplicate checks, org admin sync) |
| `DELETE /:id` | `Delete` |
| `GET /prefix/:prefix/used` | `PrefixUsed` |

Create path: customer insert + org admin + `adminCredentials`; **no** default devices (Partial).

### Phase D — Platform + docs

1. `platformauth.RequireSuperAdmin` (or inline in handler).
2. Replace `module.go` scaffold — register on `groups.Private.Group("/customers")`.
3. `make swagger` — Bearer on all handlers.
4. `docs/parity/customers.md` + update `MIGRATION.md` / `NEXT_STEPS.md`.

### Phase E — Tests

- `application/service_test.go`: search pagination, duplicate name, impersonate no admin.
- `handler_test.go`: super admin POST search 200; non–super-admin 403.

## Complexity Tracking

> No constitution violations. Create side-effect deferral (devices/config) is phased migration,
> recorded in parity doc as **Partial** until Phase 4/5.
