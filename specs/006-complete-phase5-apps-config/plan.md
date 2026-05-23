# Implementation Plan: Phase 5 — Applications, Configurations & Config Files

**Branch**: `006-complete-phase5-apps-config` | **Date**: 2026-05-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/006-complete-phase5-apps-config/spec.md`

## Summary

Deliver **Phase 5** of the Go migration: replace `applications` and `configfiles` scaffolds and
**extend** `configurations` from list-only (Phase 4) to full `ConfigurationResource` parity.
Add Postgres migration `000007_applications_configurations_core` for `applications`,
`applicationversions`, `configurationapplications`, `configurationfiles`, and extended
`configurations` columns. React **Configurations** and **Applications** pages must work
without Java; `/private/web-ui-files` (APK upload) remains **Phase 6**.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth` (`HasPermission`, super-admin helpers),
`platform/httpx/response`, existing Phase 4 `configurations` list handler (refactor into full module)

**Storage**: PostgreSQL; migration `000007_applications_configurations_core.up.sql` (Liquibase
subset); tenant `filesdir` for config-files disk writes under configured `files.directory`

**Testing**: `go test ./internal/modules/applications/... ./internal/modules/configurations/... ./internal/modules/configfiles/...`; handler tests with principal + permission flags

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`

**Project Type**: Web service + React (`frontend/src/features/configurations/`, `applications/`)

**Performance Goals**: Configuration list + detail load &lt; 5s p95 on seeded DB (spec SC-001)

**Constraints**: Permissions `applications`, `configurations`; tenant scope on all private routes;
super-admin only on `/applications/admin/*`; push on config upgrade **stub**; preserve Phase 4
`GET /configurations/list`; large binary APK via `/private/web-ui-files` **deferred Phase 6**

**Scale/Scope**: ~15 configuration endpoints + ~18 application endpoints + 1 config-files upload;
3 modules; ~45–55 Go files; parity docs ×3

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `applications`, `configurations`, `configfiles`; Phase 5 row in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | Refactor configurations from thin handler → full module layout |
| **III. API Parity** | ✅ | `contracts/*.md` + `docs/parity/applications.md`, `configurations.md`, `configfiles.md` |
| **IV. Testable Delivery** | ✅ | `go test` + `quickstart.md` + Swagger |
| **V. Simplicity** | ✅ | SQL from Java DAOs incrementally; push/web-ui-files stubbed/deferred |
| **VI. Security** | ✅ | Tenant + permission checks; super-admin guard on admin routes |
| **VII. Observability** | ✅ | Legacy error keys (duplicate name, notempty configuration, etc.) |

**Post-design**: All gates ✅. Phase 6 file upload and push delivery documented as out-of-scope, not violations.

## Project Structure

### Documentation (this feature)

```text
specs/006-complete-phase5-apps-config/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── applications-api.md
│   ├── configurations-api.md
│   └── configfiles-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000007_applications_configurations_core.up.sql
│   └── 000007_applications_configurations_core.down.sql
├── docs/parity/
│   ├── applications.md                          # NEW
│   ├── configurations.md                        # NEW (full; list was Phase 4)
│   └── configfiles.md                             # NEW
├── internal/platform/auth/
│   └── permissions.go                           # + PermApplications, PermConfigurations
├── internal/modules/configurations/             # EXTEND (was list-only)
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http + postgres
├── internal/modules/applications/               # REPLACE scaffold
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http + postgres
└── internal/modules/configfiles/                # REPLACE scaffold
    ├── module.go
    └── adapter/http/handler.go                  # multipart POST + disk write
```

**Structure Decision**: **Configurations** and **applications** are peer bounded contexts;
**configfiles** is a thin upload adapter (no heavy domain). Configuration save aggregates
nested `applications[]` and `files[]` from React `PUT` body in one transaction like Java.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Database migrations & permissions

1. `000007`: `applications`, `applicationversions`, `configurationapplications`,
   `configurationfiles`, `configurationapplicationsettings` (subset); extend `configurations`
   (type, password, design columns, qrcodekey, baseurl fields per React `Configuration` type).
2. Seed permissions `applications`, `configurations` for org-admin role; sample app + version +
   link to default configuration.
3. Verify Phase 4 `GET /list` still works.

### Phase B — Configurations module P1 (read)

| Endpoint | Notes |
|----------|-------|
| `GET /search`, `GET /search/{value}` | `configurations` permission |
| `GET /{id}` | Full editor payload |
| `GET /list` | Keep Phase 4 behavior |
| `POST /autocomplete` | Lookup items |

### Phase C — Configurations module P2 (write + apps on config)

| Endpoint | Notes |
|----------|-------|
| `PUT /` | Create/update; persist nested apps/files/settings |
| `DELETE /{id}` | Block if devices assigned |
| `PUT /copy` | Clone configuration |
| `GET /applications`, `GET /applications/{id}` | Picker / tab data |
| `PUT /application/upgrade` | Version bump; push stub |

### Phase D — Applications module P1–P2

| Endpoint | Notes |
|----------|-------|
| `GET /search`, `GET /search/{value}`, `POST /autocomplete` | Tenant catalog |
| `GET /{id}`, `GET /{id}/versions` | Detail + versions page |
| `PUT /android`, `PUT /web`, `PUT /versions` | CRUD |
| `DELETE /{id}`, `DELETE /versions/{id}` | |
| `PUT /validatePkg` | Duplicate package warning |
| `GET/POST /configurations`, `GET/POST /version/.../configurations` | Link dialogs |

### Phase E — Applications P3 (super-admin)

| Endpoint | Notes |
|----------|-------|
| `GET /admin/search`, `GET /admin/search/{value}` | Super-admin only |
| `GET /admin/common/{id}` | Merge to shared app |

### Phase F — Config files

| Endpoint | Notes |
|----------|-------|
| `POST /private/config-files` | Multipart `file`; write under `customer.filesdir` |

### Phase G — Polish

- Swagger `@Router` on all handlers; `make swagger`
- Parity docs; `MIGRATION.md` Phase 5 **done**
- `go test` + `quickstart.md` validation
- Optional: enrich Phase 4 device search `configurations` map (partial)

## Complexity Tracking

> No unjustified constitution violations. Configuration PUT aggregates many child tables —
justified by Java `ConfigurationDAO` single-save semantics and React editor payload shape.
