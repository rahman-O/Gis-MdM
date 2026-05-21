# Implementation Plan: Complete Java→Go Migration Gaps (Phase 9)

**Branch**: `011-complete-migration-gaps` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/011-complete-migration-gaps/spec.md`  
**Gap source**: [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md)

## Summary

Close **post–Phase 8 parity gaps** between Java `backend/` and Go `serverBackendGo/`: replace **`NoopPush`** notifiers with a shared **configuration/device push** service backed by the existing **`notifications` MessageQueue** (`pushmessages` + `pendingpushes`), add a **push schedule cron worker**, implement **missing REST** (`icon-files`, plugin export/rules), and deliver **P2/P3** hardening (audit middleware, sync hooks registry, customer bootstrap, stats/videos, devices search enrichment, agent file serving).

**Technical approach**: Prefer **reuse** of Phase 7 queue + Phase 8 targets over new FCM/MQTT clients in v1 — agents already poll `/rest/notifications`; Java `PushService` dual-sends MQTT + polling; Go v1 matches **polling path** (document MQTT as optional Phase 9b if env requires).

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth`, `platform/httpx/response`, `notifications/port.MessageQueue`, `plugins/shared/targets`, existing module repos

**Storage**: Postgres legacy schema; new migrations only where missing (`usagestats`, optional deviceinfo child tables); reuse `uploadedfiles` (migration `000008`) for icon-files

**Testing**: `go test ./internal/platform/push/...` + per-module `application/` tests; HTTP smoke in [quickstart.md](./quickstart.md); regression on Phases 1–8 scripts

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`; Android agents on sync/notification paths

**Project Type**: Web service (Go) + React frontend + MDM agents

**Performance Goals**: Push enqueue for 500 devices &lt; 3s; schedule tick every 60s processes due rows without blocking HTTP; export endpoints stream &lt; 30s for 10k rows (paginated/chunked)

**Constraints**: Headwind JSON envelope; identical `/rest/...` paths; tenant scope; best-effort push (log errors, don’t fail saves); `ENABLED_PLUGINS` / `MODULE_*` flags

**Scale/Scope**: ~15 modules touched; ~4–6 new migrations max; ~25 new/changed endpoints; 8 contract docs; update 15+ parity files + root gap analysis

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Phase 9 slice; modules listed below; no monolith package |
| **II. Layered Clean** | ✅ | `platform/push` = port + application; modules inject interface |
| **III. API Parity** | ✅ | Contracts per gap area; parity docs updated per FR-008/FR-019 |
| **IV. Testable Delivery** | ✅ | Unit tests for notifier + scheduler; quickstart smoke |
| **V. Simplicity** | ✅ | Reuse MessageQueue; single cron goroutine; no MQTT until justified |
| **VI. Security** | ✅ | Private routes + audit middleware; stats public documented |
| **VII. Observability** | ✅ | slog on push failures; env flags in `.env.example` |

**Post-design**: All gates ✅. Shared `platform/push` is cross-module infrastructure (like `httpx/response`), not a violation of module boundaries.

## Project Structure

### Documentation (this feature)

```text
specs/011-complete-migration-gaps/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── push-notifier-api.md
│   ├── push-schedule-worker-api.md
│   ├── icon-files-api.md
│   ├── plugins-deviceinfo-gaps-api.md
│   ├── plugins-devicelog-gaps-api.md
│   ├── platform-hardening-api.md
│   └── public-agent-gaps-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000011_usage_stats.up.sql          # if absent in DB
│   └── 000011_deviceinfo_export.up.sql    # only if export needs missing tables
├── docs/
│   ├── MIGRATION.md                       # Phase 9 row
│   ├── NEXT_STEPS.md
│   └── parity/                            # update per module
├── internal/
│   ├── platform/push/                     # NEW — ConfigurationNotifier
│   │   ├── port/notifier.go
│   │   ├── application/notifier.go
│   │   └── adapter/postgres/devices_by_config.go
│   ├── platform/audit/                   # NEW — optional HTTP audit middleware
│   ├── platform/synchooks/                # NEW — SyncResponseHook registry
│   ├── app/
│   │   ├── modules.go                    # wire QueueRepo → modules
│   │   └── scheduler.go                  # NEW — push schedule ticker
│   └── modules/
│       ├── configurations/               # inject real notifier
│       ├── devices/
│       ├── files/
│       ├── icons/                        # or sub-route icon-files
│       ├── customers/                    # bootstrap on create
│       ├── devices/                      # search enrichment
│       ├── sync/                         # hook merge
│       ├── updates/
│       ├── summary/
│       ├── stats/                        # NEW module
│       ├── videos/                       # NEW module (or publicapi extension)
│       └── plugins/
│           ├── push/                     # schedule worker hooks
│           ├── deviceinfo/               # missing handlers
│           ├── devicelog/
│           └── audit/                    # middleware registration
└── JAVA-GO-MIGRATION-GAP-ANALYSIS.md      # tracker updates
```

**Structure Decision**: **Web application** — Go backend modules + parity docs; frontend unchanged except benefiting from fixed APIs. New code favors **`internal/platform/*`** for cross-cutting push/audit/sync hooks; bounded modules for stats/videos.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — P0 Shared push notifier (FR-001, FR-002, FR-004)

1. Add `internal/platform/push` with:
   - `NotifyConfigurationChanged(ctx, configurationID)` → devices by config → `Enqueue(deviceID, "configUpdated", "")`
   - `NotifyDeviceApplicationSettings(ctx, deviceID)` → `appConfigUpdated`
2. Postgres helper: device IDs by `configurationId` (mirror `DeviceDAO.getDeviceIdsByConfigurationId`).
3. Wire in `app/modules.go`: construct `QueueRepository` once; pass `push.Notifier` into `configurations`, `devices`, `files` module registration (replace `NoopPush*`).
4. Verify `POST /rest/private/push` still works (unchanged).
5. Update parity: `configurations.md`, `devices.md`, `files.md`.

**Java refs**: `com.hmdm.notification.PushService`, `ConfigurationResource`, `DeviceResource`, `FilesResource`

### Phase B — P0 Push schedule worker (FR-003)

1. Add `internal/app/scheduler.go` — `time.Ticker` 60s (configurable `PUSH_SCHEDULE_INTERVAL_SEC`).
2. Service in `plugins/push/application/schedule_runner.go`:
   - `findMatchingTime()` SQL (port from `PushScheduleDAO`)
   - Resolve devices by scope device|group|configuration (reuse `shared/targets`)
   - Enqueue plugin message types via `MessageQueue` + mark task sent/processed
3. Start scheduler from `cmd/server` when `MODULE_PLUGINS_PUSH_ENABLED` && plugin enabled.
4. Update `docs/parity/push.md` + gap analysis §5 item 4.

**Java ref**: `PushScheduleTaskModule`

### Phase C — P1 Icon files (FR-005)

1. Extend `icons` module OR add `iconfiles` handler at `POST /rest/private/icon-files`.
2. Multipart field `file`; validate square image; resize 144px PNG; write under `FILES_DIRECTORY/{filesDir}/`.
3. Insert `uploadedfiles` row; return `UploadedFile` JSON (Java shape).
4. Errors: `error.icon.dimension.invalid`.
5. Parity doc `icons.md` or new `icon-files.md`.

**Java ref**: `IconFileResource`

### Phase D — P1 Plugin endpoint gaps (FR-006, FR-007)

| Module | Endpoints |
|--------|-----------|
| deviceinfo | `POST .../private/search/device`, `POST .../private/export`, `GET .../deviceinfo-plugin-settings/device/{deviceNumber}` |
| devicelog | `POST .../private/search/export`, `GET .../log/rules/{deviceNumber}` |

1. Read Java handlers for SQL/export format (CSV/stream).
2. Add handlers + application methods + repo queries.
3. Migration `000011_deviceinfo_export` only if tables missing.
4. Update `plugins-deviceinfo.md`, `plugins-devicelog.md`.

### Phase E — P2 Platform hardening (FR-009–FR-013)

| Item | Approach |
|------|----------|
| Audit middleware | `platform/audit` Gin middleware on `Private` group; insert `plugin_audit_log`; skip health/swagger |
| Sync hooks | `platform/synchooks.Registry`; plugins register at init; `sync` merges into response JSON |
| Customers bootstrap | On `PUT /` create: copy default configuration/devices from template customer or SQL seed |
| Files/configfiles quota | Enforce `customers.sizeLimit` + sum `uploadedfiles` sizes |
| Devices search | Add filter fields + optional apps/files enrichment in search result map |

### Phase F — P3 Public/agent gaps (FR-014–FR-018)

| Module | Work |
|--------|------|
| `stats` | `PUT /rest/public/stats` → `usagestats` table migration + DAO |
| `videos` | `POST/GET /rest/public/videos/{fileName}` under `FILES_DIRECTORY` or `VIDEO_DIRECTORY` env |
| `updates` | Remote APK fetch + persist; `sendStats` writes usage row |
| `summary` | Use `devicestatuses` if table exists |
| Static files | Gin route or middleware for agent downloads matching Java servlet path |

### Phase G — Governance (FR-019, FR-020)

1. After each phase: flip rows in `JAVA-GO-MIGRATION-GAP-ANALYSIS.md`.
2. Add **Phase 9** section to `MIGRATION.md`.
3. `make swagger` for new routes.

## Complexity Tracking

| Item | Why Needed | Simpler Alternative Rejected Because |
|------|------------|-------------------------------------|
| `platform/push` shared package | 3 modules need identical notify semantics | Duplicating enqueue logic in configurations/devices/files violates DRY and drifts from Java `PushService` |
| App-level scheduler | Cron spans plugins/push + notifications queue | Per-request cron impossible; inline goroutine in plugin module hides lifecycle |
| Sync hook registry | Multiple plugins extend sync without editing core | Hard-coding deviceinfo fields in sync module breaks module-first principle |

*No constitution gate failures requiring exemption.*

## Risk & Rollback

- **Risk**: Push storm on bulk configuration save → throttle batch enqueue (optional, match Java MQTT throttling later).
- **Rollback**: Feature flags `MODULE_PUSH_NOTIFIER_ENABLED=false` keeps Noop wiring; disable scheduler via env.
