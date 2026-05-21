---
description: "Task list for 012 — Finish Java→Go backend migration"
---

# Tasks: إكمال نقل الباكند Java → Go

**Input**: `specs/012-finish-java-go-backend/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–8 complete; [`specs/011-complete-migration-gaps/tasks.md`](../011-complete-migration-gaps/tasks.md) T001–T044 (push, cron, icon-files) **done**; [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md)

**Tests**: Application unit tests per constitution IV for search SQL builders, export streaming, quota, audit middleware, and schedule regression.

**Organization**: Tasks grouped by user story (US1–US5) for independent delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Platform: `serverBackendGo/internal/platform/audit/`, `platform/synchooks/`, `platform/push/` (regression)
- App: `serverBackendGo/internal/app/`
- Modules: `serverBackendGo/internal/modules/<name>/`
- Migrations: `serverBackendGo/db/migrations/`
- Parity: `serverBackendGo/docs/parity/`
- Trackers: `JAVA-GO-BACKEND-GAPS.md`, `JAVA-GO-MIGRATION-STATUS.md`

---

## Phase 1: Setup

**Purpose**: Confirm 012 context, Java references, and baseline build.

- [X] T001 Verify `specs/012-finish-java-go-backend/spec.md` against `JAVA-GO-BACKEND-GAPS.md` §4–§5 and `JAVA-GO-MIGRATION-STATUS.md` gap matrix
- [X] T002 [P] Review Java `DeviceResource.java`, `DeviceMapper.java` against `specs/012-finish-java-go-backend/contracts/devices-search-api.md`
- [X] T003 [P] Review Java `DeviceInfoResource.java`, `DeviceLogResource.java` against `specs/012-finish-java-go-backend/contracts/plugins-deviceinfo-gaps-api.md` and `plugins-devicelog-gaps-api.md`
- [X] T004 [P] Review Java `AuditFilter`, `SyncResource`, `CustomerResource`, `StatsResource`, `FilesResource` against `specs/012-finish-java-go-backend/contracts/platform-hardening-api.md`, `public-agent-gaps-api.md`, `files-static-api.md`
- [X] T005 Run `cd serverBackendGo && go build ./... && go test ./...` and record baseline pass/fail in `specs/012-finish-java-go-backend/quickstart.md` notes if needed

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Phase 9 P0 regression (FR-000) — ensure 011 deliverables still work before P1+.

**⚠️ CRITICAL**: No user story work until push/cron/icon-files regression passes.

- [X] T006 Verify `internal/platform/push/application/notifier.go` still wired from `internal/app/wiring.go` into configurations, devices, files modules
- [X] T007 [P] Run `specs/012-finish-java-go-backend/quickstart.md` §3 — configuration save enqueues `configUpdated` on `hmdm-001`
- [X] T008 [P] Run `specs/012-finish-java-go-backend/quickstart.md` §3 — `applicationSettings/notify` enqueues `appConfigUpdated`
- [X] T009 Verify `internal/app/scheduler.go` tick processes due `plugin_push_schedule` rows per `quickstart.md` §3 schedule snippet
- [X] T010 [P] Smoke `POST /rest/private/icon-files` per existing `serverBackendGo/docs/parity/icon-files.md`
- [X] T011 [P] Add regression test in `serverBackendGo/internal/platform/push/application/notifier_test.go` if any wiring drift found in T006
- [X] T012 Document Phase 2 checkpoint in `serverBackendGo/docs/MIGRATION.md` — Phase 9 remains **partial** until 012 complete

**Checkpoint**: P0 baseline green — US1+ implementation may start.

---

## Phase 3: User Story 1 — إدارة أجهزة (Priority: P1) 🎯 MVP

**Goal**: Full device search filters + `infojson` telemetry in detail view (FR-001, FR-002).

**Independent Test**: `quickstart.md` §4 — filtered search + `GET /private/devices/number/{n}` returns `data.info.batteryLevel`.

### Tests for User Story 1

- [X] T013 [P] [US1] Add `serverBackendGo/internal/modules/devices/application/service_test.go` — filter builder maps status to lastupdate bands
- [X] T014 [P] [US1] Add `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo_test.go` or SQL integration test for sort `LAST_UPDATE` desc

### Implementation for User Story 1

- [X] T015 [US1] Extend `serverBackendGo/internal/modules/devices/domain/device.go` `SearchRequest` with Java-aligned filter fields per `data-model.md`
- [X] T016 [US1] Add `DeviceInfoView` nested struct on `DeviceView` in `serverBackendGo/internal/modules/devices/domain/device.go`
- [X] T017 [US1] Implement `searchFilters` + `orderByClause` in `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go` for all React-sent filters
- [X] T018 [US1] Update `Count` query in `device_repo.go` to use same filters as `Search`
- [X] T019 [US1] Parse `devices.infojson` in `GetByNumber` — populate `info` in `serverBackendGo/internal/modules/devices/application/service.go`
- [X] T020 [US1] Wire handler `Search` in `serverBackendGo/internal/modules/devices/adapter/http/handler.go` — bind extended JSON body
- [X] T021 [P] [US1] Add handler test `serverBackendGo/internal/modules/devices/adapter/http/handler_test.go` for search with `status` filter
- [X] T022 [US1] Update `serverBackendGo/docs/parity/devices.md` — mark advanced filters and telemetry **Done**
- [X] T023 [P] [US1] Update `JAVA-GO-BACKEND-GAPS.md` §4.2 devices section — remove ⚠️ items closed by US1
- [X] T024 [P] [US1] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` devices table — note filters applied

**Checkpoint**: Devices page fully usable on Go vs Java for search + detail panel.

---

## Phase 4: User Story 2 — Plugins export (Priority: P1)

**Goal**: Missing deviceinfo + devicelog REST endpoints (FR-003, FR-004).

**Independent Test**: `quickstart.md` §5 — export CSV + `GET .../rules/{deviceNumber}` return 200.

### Tests for User Story 2

- [ ] T025 [P] [US2] Add `serverBackendGo/internal/modules/plugins/deviceinfo/application/export_test.go` — CSV header row matches Java column order
- [ ] T026 [P] [US2] Add `serverBackendGo/internal/modules/plugins/devicelog/application/export_test.go` — stream does not load all rows in memory

### Implementation for User Story 2 — deviceinfo

- [ ] T027 [P] [US2] Add domain filters in `serverBackendGo/internal/modules/plugins/deviceinfo/domain/` for search/device and export
- [ ] T028 [US2] Extend `deviceinfo` repo in `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/persistence/postgres/` for export queries
- [ ] T029 [US2] Implement `SearchDevice` use case in `serverBackendGo/internal/modules/plugins/deviceinfo/application/service.go`
- [ ] T030 [US2] Implement `Export` use case — stream CSV in `application/export.go`
- [ ] T031 [US2] Implement `GetSettingsForDevice` in application layer
- [ ] T032 [US2] Register routes in `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/http/handler.go`: `POST .../private/search/device`, `POST .../private/export`, `GET ...-plugin-settings/device/:deviceNumber`
- [ ] T033 [P] [US2] Update `serverBackendGo/docs/parity/plugins-deviceinfo.md` — mark three endpoints **Done**

### Implementation for User Story 2 — devicelog

- [ ] T034 [P] [US2] Implement `ExportSearch` in `serverBackendGo/internal/modules/plugins/devicelog/application/service.go`
- [ ] T035 [US2] Implement `GetRulesForDevice` in `serverBackendGo/internal/modules/plugins/devicelog/application/service.go`
- [ ] T036 [US2] Register routes in `serverBackendGo/internal/modules/plugins/devicelog/adapter/http/handler.go`: `POST .../log/private/search/export`, `GET .../log/rules/:deviceNumber`
- [ ] T037 [P] [US2] Update `serverBackendGo/docs/parity/plugins-devicelog.md` — mark export + rules **Done**
- [ ] T038 [US2] Update `JAVA-GO-BACKEND-GAPS.md` §5 endpoints list — strike deviceinfo/devicelog ❌ rows

**Checkpoint**: Plugin gap REST complete for monitoring plugins.

---

## Phase 5: User Story 3 — امتثال وتشغيل tenant (Priority: P2)

**Goal**: Audit middleware, sync hooks, customer bootstrap (FR-005, FR-006, FR-007).

**Independent Test**: Delete device → audit search finds row; new customer has default config; sync response includes plugin hook field when registered.

### Tests for User Story 3

- [ ] T039 [P] [US3] Add `serverBackendGo/internal/platform/audit/application/writer_test.go` — truncates body, skips excluded paths
- [ ] T040 [P] [US3] Add `serverBackendGo/internal/platform/synchooks/registry_test.go` — multiple hooks merge keys

### Implementation for User Story 3 — audit

- [ ] T041 [US3] Create `serverBackendGo/internal/platform/audit/port/repository.go` — `InsertAuditLog`
- [ ] T042 [US3] Implement `serverBackendGo/internal/platform/audit/adapter/postgres/audit_repo.go` → `plugin_audit_log`
- [ ] T043 [US3] Implement `serverBackendGo/internal/platform/audit/application/writer.go` — async-safe insert
- [ ] T044 [US3] Add Gin middleware in `serverBackendGo/internal/platform/audit/adapter/http/middleware.go`
- [ ] T045 [US3] Register audit middleware on private route group in `serverBackendGo/internal/app/modules.go` or `internal/platform/httpx/router.go`

### Implementation for User Story 3 — sync hooks

- [ ] T046 [P] [US3] Create `serverBackendGo/internal/platform/synchooks/registry.go` — `Register`, `Extend(map[string]any)`
- [ ] T047 [US3] Call `synchooks.Extend` from `serverBackendGo/internal/modules/sync/application/service.go` after building `SyncResponse`
- [ ] T048 [P] [US3] Register stub hook from `serverBackendGo/internal/modules/plugins/deviceinfo/module.go` (or platform) when plugin enabled — document extension key in parity

### Implementation for User Story 3 — customers bootstrap

- [ ] T049 [US3] Read Java `CustomerResource` create path and document template IDs in `serverBackendGo/docs/parity/customers.md`
- [ ] T050 [US3] Implement bootstrap in `serverBackendGo/internal/modules/customers/application/service.go` on create — copy default configuration (+ optional device)
- [ ] T051 [US3] Add transaction boundaries in `serverBackendGo/internal/modules/customers/adapter/persistence/postgres/` if multi-table insert required
- [ ] T052 [P] [US3] Update `JAVA-GO-BACKEND-GAPS.md` customers + audit + sync sections
- [ ] T053 [P] [US3] Update `serverBackendGo/docs/parity/plugins-audit.md` — AuditFilter **Done** (middleware)
- [ ] T054 [P] [US3] Update `serverBackendGo/docs/parity/sync.md` — SyncResponseHook **Done**

**Checkpoint**: Compliance and tenant onboarding paths match Java behavior.

---

## Phase 6: User Story 4 — ملفات وتخزين (Priority: P2)

**Goal**: Quota, `uploadedfiles`, agent static `/files/*` (FR-008, FR-009).

**Independent Test**: `quickstart.md` §6 — over-quota upload fails; `GET /files/...` serves bytes.

### Tests for User Story 4

- [ ] T055 [P] [US4] Add `serverBackendGo/internal/modules/configfiles/application/quota_test.go` — rejects when sum exceeds `sizeLimit`
- [ ] T056 [P] [US4] Add test for static handler path traversal rejection in `serverBackendGo/internal/app/files_static_test.go`

### Implementation for User Story 4

- [ ] T057 [US4] Add `CheckQuota` + `RecordUploadedFile` to `serverBackendGo/internal/modules/configfiles/port/repository.go` and postgres adapter
- [ ] T058 [US4] Call quota check from `serverBackendGo/internal/modules/configfiles/application/service.go` on upload
- [ ] T059 [P] [US4] Mirror quota + `uploadedfiles` insert in `serverBackendGo/internal/modules/files/application/service.go` for web-ui-files upload paths
- [ ] T060 [US4] Add `serverBackendGo/db/migrations/000011_uploadedfiles_indexes.up.sql` if missing indexes for quota sum (optional)
- [ ] T061 [US4] Implement `GET /files/*filepath` handler in `serverBackendGo/internal/app/files_static.go` using `platform/storage.LocalStore` + `FILES_DIRECTORY`
- [ ] T062 [US4] Register static route on engine in `serverBackendGo/internal/app/wiring.go` or `cmd/server/main.go`
- [ ] T063 [US4] Update `serverBackendGo/docs/parity/configfiles.md` and `files.md` — quota + static **Done**
- [ ] T064 [P] [US4] Update `contracts/files-static-api.md` status and `JAVA-GO-BACKEND-GAPS.md` files section

**Checkpoint**: Agents can download files; uploads respect storage limits.

---

## Phase 7: User Story 5 — وحدات عامة وتحديثات (Priority: P3)

**Goal**: stats module, optional videos, updates APK, summary charts, configurations QR fields (FR-010–FR-014).

**Independent Test**: `quickstart.md` §7 — `PUT /rest/public/stats` OK; update check/apply smoke.

### Tests for User Story 5

- [ ] T065 [P] [US5] Add `serverBackendGo/internal/modules/stats/application/service_test.go` — persists payload
- [ ] T066 [P] [US5] Add `serverBackendGo/internal/modules/updates/application/apply_test.go` — mock HTTP download for manifest APK

### Implementation for User Story 5 — stats

- [ ] T067 [US5] Add migration `serverBackendGo/db/migrations/000011_usage_stats.up.sql` from Java `UsageStats` columns
- [ ] T068 [US5] Scaffold `serverBackendGo/internal/modules/stats/` (domain, port, application, adapter/http, module.go)
- [ ] T069 [US5] Implement `PUT /rest/public/stats` in `serverBackendGo/internal/modules/stats/adapter/http/handler.go`
- [ ] T070 [US5] Register module in `serverBackendGo/internal/app/modules.go` with `MODULE_STATS_ENABLED`
- [ ] T071 [P] [US5] Add `serverBackendGo/docs/parity/stats.md` — **Done**

### Implementation for User Story 5 — videos (conditional)

- [ ] T072 [US5] Product check: if videos unused, add `serverBackendGo/docs/parity/videos.md` with ⊘ and skip module; else scaffold `internal/modules/videos/` per `public-agent-gaps-api.md`
- [ ] T073 [P] [US5] Add `MODULE_VIDEOS_ENABLED` and `VIDEO_DIRECTORY` to `serverBackendGo/.env.example` (default false)

### Implementation for User Story 5 — updates + summary + configurations

- [ ] T074 [US5] Implement remote APK download in `serverBackendGo/internal/modules/updates/application/service.go` `Apply` path
- [ ] T075 [US5] Call stats port from updates after apply (`sendStats`) via interface in `serverBackendGo/internal/modules/updates/port/`
- [ ] T076 [US5] Extend `serverBackendGo/internal/modules/summary/adapter/persistence/postgres/summary_repo.go` — populate chart arrays from `devicestatuses` when present
- [ ] T077 [US5] Audit `PUT /configurations` save path — ensure QR fields (`qrCodeKey`, `baseUrl`, `mainAppId`, design colors) persist in `configurations` repo
- [ ] T078 [P] [US5] Update `serverBackendGo/docs/parity/updates.md`, `summary.md`, `configurations.md`
- [ ] T079 [P] [US5] Update `JAVA-GO-BACKEND-GAPS.md` — mark stats/videos/updates/summary items

**Checkpoint**: Public agent gaps closed or explicitly ⊘.

---

## Phase 8: Polish & Cross-Cutting

**Purpose**: Governance, docs, full UAT, swagger (FR-015–FR-017).

- [ ] T080 Run `cd serverBackendGo && go test ./...` — fix failures
- [ ] T081 [P] Run full `specs/012-finish-java-go-backend/quickstart.md` §8 UAT checklist on Go-only stack
- [ ] T082 [P] Run React smoke against Go — login, devices filters, configurations, push per `FRONTEND-GO-BACKEND-INTEGRATION.md`
- [ ] T083 Update `JAVA-GO-MIGRATION-STATUS.md` — Phase 9 **done**, ≥95% behavioral parity note
- [ ] T084 [P] Update `JAVA-GO-BACKEND-GAPS.md` — resolve all ❌/critical ⚠️ or document ⊘
- [ ] T085 Update `serverBackendGo/docs/MIGRATION.md` — Phase 9 status **done**
- [ ] T086 [P] Run `cd serverBackendGo && make swagger` for new handlers; restart dev and verify `/swagger/index.html`
- [ ] T087 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — next work beyond Java parity
- [ ] T088 Mark completed tasks in this file and link PR / branch `012-finish-java-go-backend`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Setup — **BLOCKS** all user stories
- **Phase 3 (US1)**: After Phase 2 — MVP for React devices page
- **Phase 4 (US2)**: After Phase 2 — independent of US1 (parallel OK)
- **Phase 5 (US3)**: After Phase 2 — audit middleware should register before heavy private testing
- **Phase 6 (US4)**: After Phase 2 — independent; benefits from US3 customer bootstrap optional order
- **Phase 7 (US5)**: After Phase 2 — updates `sendStats` depends on T069 stats module
- **Phase 8 (Polish)**: After desired user stories complete

### User Story Dependencies

| Story | Depends on | Can parallel with |
|-------|------------|-------------------|
| US1 | Phase 2 | US2, US4 |
| US2 | Phase 2 | US1, US4 |
| US3 | Phase 2 | US1, US2 (audit middleware global) |
| US4 | Phase 2 | US1, US2 |
| US5 | Phase 2; US5 updates needs stats (T069 before T075) | US1–US4 except T075 |

### Within Each User Story

- Domain/port before application before HTTP routes
- Parity doc update last in story
- Tests after implementation unless noted for SQL builders

### Parallel Opportunities

- T002–T004 (Java reviews) in Phase 1
- T007–T010 smoke in Phase 2
- US1 and US2 entire phases in parallel after Phase 2
- US3 audit vs synchooks vs customers (T041–T054) partially parallel across subfolders
- T083–T087 doc updates in Polish

---

## Parallel Example: User Story 2

```bash
# deviceinfo track (developer A):
T027 → T028 → T029 → T030 → T031 → T032 → T033

# devicelog track (developer B):
T034 → T035 → T036 → T037

# merge: T038 gap tracker
```

---

## Parallel Example: User Story 3

```bash
# Track A — audit:
T041 → T042 → T043 → T044 → T045

# Track B — sync hooks:
T046 → T047 → T048

# Track C — customers:
T049 → T050 → T051
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1 + Phase 2 (regression)
2. Complete Phase 3 (US1 — devices)
3. **STOP and VALIDATE** — React devices page with filters + detail panel
4. Demo without waiting for plugins/public modules

### Incremental Delivery

1. Phase 2 → US1 (devices) → deploy/demo
2. Add US2 (plugins export)
3. Add US3 (audit/sync/customers) + US4 (files) in parallel
4. Add US5 (stats/updates/summary)
5. Phase 8 governance → declare Java WAR optional

### Mapping to 011 tasks

| 011 (done) | 012 continues |
|------------|----------------|
| T001–T044 push, cron, icon-files | Phase 2 regression only |
| T045–T093 (never done in 011) | **This file T013–T088** |

---

## Notes

- Do not re-implement `platform/push` or schedule cron except regression (Phase 2).
- Videos: default ⊘ unless T072 confirms product need.
- Commit per task group; stop at each **Checkpoint**.
- `[P]` tasks must not edit the same file concurrently.
