<!--
Sync Impact Report
- Version change: (template) → 1.0.0
- Modified principles: N/A (initial ratification; replaced all placeholders)
- Added sections: Technology & Structure Standards; Development Workflow & Quality Gates
- Removed sections: Generic template examples (Library-First, CLI, etc.)
- Templates: plan-template.md ✅ | spec-template.md ✅ | tasks-template.md ✅
- Commands: N/A (no .specify/templates/commands/)
- Deferred: none
-->

# Gis-MdM Constitution

## Core Principles

### I. Module-First Migration

Every backend capability MUST be delivered as one bounded context under
`serverBackendGo/internal/modules/<name>/`, migrated in phases documented in
`serverBackendGo/docs/MIGRATION.md` and `serverBackendGo/docs/NEXT_STEPS.md`.

- One module (or cohesive plugin subtree) per PR slice when possible.
- Java `*Resource.java` and DAOs are the source of truth until parity is signed off.
- Scaffold modules MUST register routes only when implementation exists; no fake 200s
  that hide missing parity.
- Cross-module reuse goes through `internal/platform/`, `internal/shared/`, or
  explicit `port` interfaces — never copy-paste handlers across modules.

**Rationale**: Headwind MDM is large; incremental modules keep the React frontend and
Android agents working while risk stays bounded.

### II. Layered Clean Architecture (NON-NEGOTIABLE)

Each module MUST follow the dependency rule:

```text
adapter/http, adapter/persistence → application → port, domain
```

- `domain/`: entities and DTOs; no imports from Gin, SQL, or other modules' adapters.
- `port/`: repository and external service interfaces consumed by `application/`.
- `application/`: use cases only; no HTTP or SQL literals.
- `adapter/`: HTTP (Gin), Postgres, or other IO; maps to/from domain at the edge.
- `module.go`: wires the module into `internal/app` via `module.Module`.

**Rationale**: Clear layers make modules easy to test, replace (e.g. persistence), and
extend without spaghetti.

### III. API Contract Parity

Public behavior MUST match the legacy Headwind MDM API unless a deliberate,
documented breaking change is approved.

- Paths: `/rest/public`, `/rest/private`, `/rest/plugins`, `/rest/plugin/main`.
- JSON envelope: `{ "status": "OK"|"ERROR", "message"?, "data"? }` via
  `internal/platform/httpx/response`.
- Password and auth flows MUST stay compatible with the React client (MD5 hex, session,
  JWT, optional RSA when `TRANSMIT_PASSWORD` is enabled).
- Each implemented module MUST have `serverBackendGo/docs/parity/<module>.md` listing
  endpoints and Java reference classes.

**Rationale**: The frontend and agents are consumers; parity avoids dual maintenance of
Java and Go in production.

### IV. Incremental, Testable Delivery

A module is "done" only when:

1. `go build ./...` passes in `serverBackendGo/`.
2. Application-layer unit tests exist for non-trivial logic (auth, crypto, parsing).
3. Manual or scripted smoke check is documented (curl or Swagger).
4. Parity doc is updated; `NEXT_STEPS.md` status row is updated when applicable.

Tests on adapters SHOULD use interfaces/stubs (`port` stubs or testcontainers only when
justified). Prefer testing `application/` over full HTTP unless contract tests are
required.

**Rationale**: Migration speed must not trade away regressions on critical paths (auth,
tenant settings, devices).

### V. Simplicity & Maintainability

- YAGNI: no frameworks inside modules, no generic "base handler" hierarchies unless
  three modules share identical code.
- Files SHOULD stay focused (one use-case group per `application/*.go` file).
- Migrations in `serverBackendGo/db/migrations/*.up.sql` MUST be idempotent where
  possible and named sequentially.
- Configuration via environment (`internal/config`); secrets MUST NOT be committed.
- Comments only for non-obvious business rules (legacy Headwind quirks, security).

**Rationale**: A small team maintaining a Java→Go migration needs readable code more
than clever abstractions.

### VI. Security & Multi-Tenancy

- All `/rest/private/*` routes MUST use JWT and/or session middleware unless explicitly
  public.
- Handlers MUST scope data by authenticated principal (`customerId`, permissions when
  implemented).
- Auth-sensitive flows (login, reset, signup) MUST avoid user enumeration leaks matching
  Java behavior.
- New endpoints MUST NOT log passwords, tokens, or PII at info level.

**Rationale**: MDM is a control plane; tenant isolation and auth correctness are
non-negotiable.

### VII. Observability & Operations

- Use structured logging via `internal/platform/logger` (`slog`).
- Errors returned to clients use stable `message` keys where Java does (e.g.
  `error.permission.denied`).
- Feature flags (`MODULE_*_ENABLED`) for optional modules; defaults documented in
  `.env.example`.
- Docker/scripts (`scripts/db-up.sh`, `make dev`) MUST remain the documented local path.

**Rationale**: Operators and developers need predictable runbooks during parallel Java/Go
operation.

## Technology & Structure Standards

| Area | Standard |
|------|----------|
| Go backend | `serverBackendGo/`, Go 1.22+, Gin, `lib/pq`, golang-migrate-style SQL |
| Legacy reference | `backend/` Java, `frontend/` React — do not break contracts |
| Module layout | `module.go`, `domain/`, `port/`, `application/`, `adapter/http/`, `adapter/persistence/postgres/` |
| Composition | `cmd/server` → `internal/app` registers all `module.Module` implementations |
| Docs | `docs/MIGRATION.md`, `docs/NEXT_STEPS.md`, `docs/parity/*.md`, `docs/AUTH_COMPLETE.md` |
| Frontend proxy | Vite → `http://localhost:8080` for `/rest` when using Go dev server |

Plugins live under `internal/modules/plugins/<name>/` with the same layer rules.

## Development Workflow & Quality Gates

1. **Specify** features against this constitution; plans MUST pass Constitution Check.
2. **Implement** one module slice; update parity + NEXT_STEPS.
3. **Verify** `go test ./...` and smoke paths from feature spec.
4. **Review** PRs for layer violations, path drift, and missing parity docs.
5. **Do not** expand scope into unrelated modules in the same change.

Complexity that violates Principles II or V MUST be recorded in the plan's
Complexity Tracking table with rejected simpler alternatives.

## Governance

- This constitution supersedes ad-hoc conventions for Spec Kit workflows and Go
  backend work in this repository.
- Amendments: update this file, bump version per semver below, refresh affected
  `.specify/templates/*`, and note impact in the Sync Impact Report comment.
- All implementation plans and specs MUST include a Constitution Check section;
  gates are defined in `.specify/templates/plan-template.md`.
- Runtime development guidance: `serverBackendGo/README.md`, `docs/MIGRATION.md`,
  `docs/NEXT_STEPS.md`.

**Version**: 1.0.0 | **Ratified**: 2026-05-20 | **Last Amended**: 2026-05-20
