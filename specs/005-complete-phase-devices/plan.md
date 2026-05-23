# Implementation Plan: Phase 4 — Devices & Groups Module Migration

**Branch**: `005-complete-phase-devices` | **Date**: 2026-05-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/005-complete-phase-devices/spec.md`

## Summary

Deliver **Phase 4** of the Go migration: replace `devices` and `groups` scaffolds with full
`DeviceResource` / `GroupResource` parity for the React console; add Postgres migrations
(`devices`, `groups`, `devicegroups`, minimal `configurations`); implement read-only
`GET /configurations/list`; upgrade `summary` device stats SQL; stub push notify endpoints.
Largest effort is `POST /devices/search` returning `DeviceListResponse` shape
(`configurations` map + paginated `devices.items`).

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth` (permissions, principal, group access),
`platform/httpx/response`, existing `summary` module

**Storage**: PostgreSQL; new migrations `000006_devices_groups_core.up.sql` (+ optional
`000007_devices_search_columns.up.sql`) mirroring Liquibase subset; seed default group per tenant

**Testing**: `go test ./internal/modules/devices/... ./internal/modules/groups/...`; handler tests
with JWT principal + permission flags; extend summary tests

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`

**Project Type**: Web service + React (`frontend/src/features/devices/`, `groups/`)

**Performance Goals**: Device search &lt; 5s p95 with 500 devices (spec SC-001)

**Constraints**: Tenant-scoped (`customerId`); `pageNum`/`pageSize` (not `currentPage`);
permissions `edit_devices`, `edit_device_desc`, `settings`; push notify no-op; config list
read-only; search enrichment **partial** v1 (config name + status color, defer full apps/files)

**Scale/Scope**: ~11 device endpoints + 5 group endpoints + 1 config list + summary upgrade;
2 main modules + thin configurations + migration; ~25–35 Go files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `devices`, `groups` modules; Phase 4 closes `MIGRATION.md` row |
| **II. Layered Clean** | ✅ | domain/port/application/adapter per module |
| **III. API Parity** | ✅ | `contracts/*.md` + `docs/parity/devices.md`, `groups.md` |
| **IV. Testable Delivery** | ✅ | `go test` + `quickstart.md` + Swagger |
| **V. Simplicity** | ✅ | SQL from `DeviceMapper.xml` / `GroupDAO` incrementally; push stub |
| **VI. Security** | ✅ | Tenant + user group access + permission checks |
| **VII. Observability** | ✅ | Legacy error keys (`error.notempty.group`, device exists) |

**Post-design**: All gates ✅. Partial search enrichment and push documented in parity, not violations.

## Project Structure

### Documentation (this feature)

```text
specs/005-complete-phase-devices/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── devices-api.md
│   ├── groups-api.md
│   └── configurations-list-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000006_devices_groups_core.up.sql      # groups, devices, devicegroups, configurations
│   └── 000006_devices_groups_core.down.sql
├── docs/parity/
│   ├── devices.md                             # NEW
│   ├── groups.md                              # NEW
│   └── summary.md                             # UPDATE — real stats
├── internal/platform/auth/
│   └── permissions.go                         # + edit_devices, edit_device_desc constants
├── internal/modules/devices/
│   ├── module.go
│   ├── domain/          # Device, SearchRequest, DeviceListView, payloads
│   ├── port/
│   ├── application/     # Search, CRUD, bulk, autocomplete, app settings, description
│   └── adapter/
│       ├── http/handler.go
│       └── persistence/postgres/device_repo.go
├── internal/modules/groups/
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http + postgres
├── internal/modules/configurations/
│   └── adapter/http/handler.go                # GET /list only (minimal)
└── internal/modules/summary/
    └── adapter/persistence/postgres/summary_repo.go   # real GetDeviceStats SQL
```

**Structure Decision**: Two primary bounded modules (`devices`, `groups`). Configurations
list and summary upgrade are thin cross-cuts within Phase 4 scope per spec FR-X02/X03.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Database migrations & seed

1. `000006_devices_groups_core.up.sql`: `groups`, `devices`, `devicegroups`, `configurations`
   (id, name, customerid, mainappid nullable), optional `deviceapplicationsettings`.
2. Per-tenant default group + sample configuration + optional seed devices for dev smoke.
3. Verify `summary` `HasDevicesTable` returns true after migrate.

### Phase B — Groups module (P1 dependency for device filters)

| Endpoint | Handler |
|----------|---------|
| `GET /search` | List all groups for customer |
| `GET /search/{value}` | Filter by name (parity) |
| `POST /autocomplete` | Lookup items |
| `PUT /` | Create/update (`settings`) |
| `DELETE /{id}` | Delete if empty |

### Phase C — Devices module P1 (search + read)

| Endpoint | Notes |
|----------|-------|
| `POST /search` | Build `DeviceListResponse`; map `pageNum`; user group access join |
| `GET /number/{number}` | Single `DeviceView` |

### Phase D — Devices module P2 (mutations)

| Endpoint | Notes |
|----------|-------|
| `PUT /` | Create/update/bulk config ids; device limit check |
| `DELETE /{id}` | |
| `POST /deleteBulk`, `POST /groupBulk` | |
| `POST /autocomplete` | |
| `POST /{id}/description` | `edit_device_desc` |

### Phase E — Devices P3 + cross-cutting

| Item | Notes |
|------|-------|
| App settings GET/POST/notify | notify → no-op push port |
| `GET /configurations/list` | configurations module |
| `summary` repo | Port Java summary count queries |
| Swagger + parity docs + `MIGRATION.md` Phase 4 **done** |

### Phase F — Tests

- `devices/application` search + permission denied
- `groups/application` duplicate + not empty delete
- HTTP handler smoke with super-admin/org-admin principal

## Complexity Tracking

> No unjustified constitution violations. Device search SQL complexity justified by React
> contract; delivered incrementally (core list first, enrichment partial).
