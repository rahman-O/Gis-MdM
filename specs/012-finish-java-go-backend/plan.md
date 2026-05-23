# Implementation Plan: إكمال نقل الباكند Java → Go

**Branch**: `012-finish-java-go-backend` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/012-finish-java-go-backend/spec.md`  
**Gap sources**: [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md), [`JAVA-GO-MIGRATION-STATUS.md`](../../JAVA-GO-MIGRATION-STATUS.md)  
**Predecessor**: [`specs/011-complete-migration-gaps/`](../011-complete-migration-gaps/) — P0 (push notifier, cron, icon-files) **done**; this plan starts at **P1+**

## Summary

Complete remaining **Java→Go backend parity** so React and MDM agents can run on `serverBackendGo` without the Java WAR: enrich **devices** search/telemetry, implement **missing plugin REST** (deviceinfo/devicelog exports), add **platform hardening** (audit middleware, sync hooks, customer bootstrap, file quota/static serving), and deliver **public modules** (`stats`, optional `videos`, updates APK, summary charts).

**Technical approach**: Extend existing modules in place; add `stats` (and `videos` if in scope) as bounded modules; reuse **`internal/platform/push`** (Phase 9 baseline); introduce **`platform/audit`** and **`platform/synchooks`** only where cross-cutting; align SQL with Java `DeviceMapper` / `infojson` patterns; no MQTT in v1.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth`, `platform/httpx/response`, `platform/push`, `platform/storage`, `notifications/port.MessageQueue`, existing module repos

**Storage**: Postgres legacy schema; migrations `000011+` only for `usagestats`, missing plugin columns, or indexes for search filters

**Testing**: `go test ./...`; per-module `application/` tests; [quickstart.md](./quickstart.md); regression via `FRONTEND-GO-BACKEND-INTEGRATION.md` UAT paths

**Target Platform**: Linux/macOS dev (`:8080`); Vite proxy `/rest`; Android agents (sync, notifications, `/files/*`)

**Project Type**: Web service (Go) + React frontend + MDM agents

**Performance Goals**: Device search with filters &lt; 2s p95 for 10k devices (indexed columns); export stream start &lt; 5s; audit middleware overhead &lt; 5ms p95; static file serve from disk

**Constraints**: Headwind envelope; identical `/rest/...` paths; tenant scope; best-effort push (already implemented); `MODULE_*` / `ENABLED_PLUGINS` flags

**Scale/Scope**: ~12 modules extended; 2 new modules (`stats`, optional `videos`); ~20 endpoint/handler additions; 6 contract docs; update parity + root `JAVA-GO-*.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Per-gap modules; Phase 9 completion in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | platform/* = cross-cutting only; domain logic in modules |
| **III. API Parity** | ✅ | Contracts in `contracts/`; parity docs per FR-015 |
| **IV. Testable Delivery** | ✅ | Unit + quickstart smoke per user story |
| **V. Simplicity** | ✅ | Reuse push queue; no MQTT; videos ⊘ if unused |
| **VI. Security** | ✅ | Private routes scoped; audit excludes secrets |
| **VII. Observability** | ✅ | slog; stable error keys |

**Post-design**: All gates ✅. `platform/audit` and `platform/synchooks` justified as servlet-filter equivalents (see 011 research R6–R7).

## Project Structure

### Documentation (this feature)

```text
specs/012-finish-java-go-backend/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/           # Phase 1
│   ├── devices-search-api.md
│   ├── plugins-deviceinfo-gaps-api.md   # extends 011
│   ├── plugins-devicelog-gaps-api.md
│   ├── platform-hardening-api.md
│   ├── public-agent-gaps-api.md
│   └── files-static-api.md
└── tasks.md             # (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   └── 000011_*.up.sql              # usagestats; optional indexes
├── docs/
│   ├── MIGRATION.md                 # Phase 9 → done
│   └── parity/                      # update per module
├── internal/
│   ├── platform/
│   │   ├── push/                    # EXISTS — regression only
│   │   ├── audit/                   # NEW — middleware
│   │   └── synchooks/               # NEW — registry
│   ├── app/
│   │   ├── wiring.go                # static /files route
│   │   └── scheduler.go             # EXISTS — regression
│   └── modules/
│       ├── devices/                 # search filters + infojson parse
│       ├── customers/               # bootstrap
│       ├── configfiles/             # quota + uploadedfiles
│       ├── files/                   # static serve helper
│       ├── sync/                    # hook merge
│       ├── updates/                 # APK download
│       ├── summary/                 # charts data
│       ├── stats/                   # NEW
│       ├── videos/                  # NEW or ⊘
│       └── plugins/
│           ├── deviceinfo/          # 3 endpoints
│           ├── devicelog/           # 2 endpoints
│           └── audit/               # middleware wire
├── JAVA-GO-BACKEND-GAPS.md          # tracker
└── JAVA-GO-MIGRATION-STATUS.md
```

**Structure Decision**: Single Go backend (`serverBackendGo/`); frontend unchanged. Work ordered **P1 → P2 → P3** per spec user stories.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Regression baseline (FR-000)

1. Verify `platform/push`, `scheduler`, `icon-files` from 011 still pass `go test ./...` and quickstart §3–§5.
2. Lock behavior with handler tests if gaps found.

### Phase B — P1 Devices (FR-001, FR-002)

1. Extend `devices/domain.SearchRequest` with Java-aligned fields.
2. Update `device_repo.Search` / `Count` SQL (status from `lastupdate` bands, `infojson` filters, sort).
3. Parse `infojson` in `GetByNumber` → nested `info` object for React `DeviceDetailPanel`.
4. Update `docs/parity/devices.md`; smoke from [contracts/devices-search-api.md](./contracts/devices-search-api.md).

### Phase C — P1 Plugin gaps (FR-003, FR-004)

1. **deviceinfo**: handlers for `search/device`, `export`, `settings/device/{deviceNumber}` — Java reference in contracts.
2. **devicelog**: `search/export`, `rules/{deviceNumber}`.
3. Update parity docs; export smoke in quickstart §6.

### Phase D — P2 Platform hardening (FR-005–FR-009)

1. `platform/audit` middleware on private routes → `plugin_audit_log`.
2. `platform/synchooks` registry; `sync` merges hooks after core response.
3. `customers` bootstrap on create (template devices/configurations).
4. `configfiles` + `files`: `uploadedfiles` row, `sizeLimit` enforcement.
5. Gin route `GET /files/*` or static middleware using `FILES_DIRECTORY` + `BuildPublicURL` pattern.

### Phase E — P3 Public & polish (FR-010–FR-014)

1. New `stats` module: `PUT /rest/public/stats`.
2. `videos` module **or** document ⊘ in parity after product check.
3. `updates`: remote manifest download + `sendStats` call to stats.
4. `summary`: populate chart arrays from DB when `devicestatuses` available.
5. `configurations`: verify QR-critical fields persist (extend domain/repo if needed).

### Phase F — Governance (FR-015–FR-017)

1. Update `JAVA-GO-BACKEND-GAPS.md`, `JAVA-GO-MIGRATION-STATUS.md`, `MIGRATION.md` Phase 9 = **done**.
2. Full UAT script (30 min) in quickstart §8.
3. `.env.example` flags for new modules.

## Complexity Tracking

| Item | Why Needed | Simpler Alternative Rejected Because |
|------|------------|-------------------------------------|
| `platform/audit` | Servlet `AuditFilter` equivalent | Per-handler logging duplicates Java and misses routes |
| `platform/synchooks` | Guice `SyncResponseHook` set | Hard-coding plugin fields in sync breaks extensibility |
| `infojson` parsing in devices | React expects `device.info.*` | Flat columns duplicate Java JSON model |

## Dependencies & Risks

| Risk | Mitigation |
|------|------------|
| `infojson` schema drift | Mirror Java keys; test with seeded device |
| Export OOM | Stream CSV; limit row count |
| Videos unused | FR-011: product decision → ⊘ doc |
| 011 tasks overlap | Start tasks at T046 mapping in tasks.md |

## Artifacts Generated (Phase 0–1)

| Artifact | Path |
|----------|------|
| Research | [research.md](./research.md) |
| Data model | [data-model.md](./data-model.md) |
| Quickstart | [quickstart.md](./quickstart.md) |
| Contracts | [contracts/](./contracts/) |

**Next command**: `/speckit-tasks` to generate ordered `tasks.md`.
