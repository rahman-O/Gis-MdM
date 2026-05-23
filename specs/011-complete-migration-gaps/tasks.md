---
description: "Task list for Phase 9 — Complete Java→Go migration gaps"
---

# Tasks: Phase 9 — Complete Migration Gaps

**Input**: `specs/011-complete-migration-gaps/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–8 complete; Postgres; seeded device `hmdm-001`; [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md)

**Tests**: Application unit tests per constitution IV for push notifier, schedule runner, and non-trivial export/quota logic.

**Organization**: Tasks grouped by user story (US1–US6) for independent delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Platform: `serverBackendGo/internal/platform/push/`, `platform/audit/`, `platform/synchooks/`
- App wiring: `serverBackendGo/internal/app/`
- Modules: `serverBackendGo/internal/modules/<name>/`
- Migrations: `serverBackendGo/db/migrations/`
- Parity: `serverBackendGo/docs/parity/`
- Tracker: `JAVA-GO-MIGRATION-GAP-ANALYSIS.md`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 9 context and Java/React baseline for gap closure.

- [X] T001 Verify feature context in `specs/011-complete-migration-gaps/spec.md` against `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §5–§12
- [X] T002 [P] Review Java `PushService.java`, `PushScheduleTaskModule.java`, `IconFileResource.java` against `specs/011-complete-migration-gaps/contracts/push-notifier-api.md` and `push-schedule-worker-api.md`
- [X] T003 [P] Review Java `DeviceInfoResource.java`, `DeviceLogResource.java` against `specs/011-complete-migration-gaps/contracts/plugins-deviceinfo-gaps-api.md` and `plugins-devicelog-gaps-api.md`
- [X] T004 [P] Review Java `StatsResource.java`, `VideosResource.java`, `CustomerResource.java` against `specs/011-complete-migration-gaps/contracts/public-agent-gaps-api.md` and `platform-hardening-api.md`
- [X] T005 Run baseline `cd serverBackendGo && go build ./...` and document current `NoopPush` wiring in `internal/modules/configurations/module.go`, `devices/module.go`, `files/module.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared push infrastructure, env flags, optional migrations, and app-level queue wiring.

**⚠️ CRITICAL**: No user story delivery until `platform/push` and `MessageQueue` injection exist.

- [X] T006 Add Phase 9 env vars to `serverBackendGo/internal/config/config.go`: `PushNotifierEnabled`, `PushScheduleIntervalSec`, `ModuleStatsEnabled`, `ModuleVideosEnabled`, `VideoDirectory`
- [X] T007 [P] Document Phase 9 env in `serverBackendGo/.env.example` and export in `serverBackendGo/scripts/dev.sh`
- [X] T008 Create `serverBackendGo/internal/platform/push/port/notifier.go` — `Notifier` interface per `contracts/push-notifier-api.md`
- [X] T009 Create `serverBackendGo/internal/platform/push/port/device_lookup.go` — `DeviceIDsByConfiguration`, `DeviceExists` interfaces
- [X] T010 Implement `serverBackendGo/internal/platform/push/adapter/postgres/device_lookup.go` — SQL `devices` by `configurationid`
- [X] T011 Implement `serverBackendGo/internal/platform/push/application/notifier.go` — `NotifyConfigurationChanged`, `NotifyDeviceApplicationSettings` using `notifications/port.MessageQueue`
- [X] T012 [P] Add `serverBackendGo/internal/platform/push/application/notifier_test.go` — mock queue asserts `configUpdated` / `appConfigUpdated` types
- [X] T013 Refactor `serverBackendGo/internal/app/modules.go` (or new `internal/app/wiring.go`) — construct single `notifications/adapter/persistence/postgres.QueueRepository` and `platform/push` notifier
- [X] T014 Update `serverBackendGo/internal/modules/configurations/module.go` — inject real `push.Notifier` when `PushNotifierEnabled` (replace `port.NoopPushNotifier{}`)
- [X] T015 [P] Update `serverBackendGo/internal/modules/devices/module.go` — inject real notifier (replace `port.NoopPush{}`)
- [X] T016 [P] Update `serverBackendGo/internal/modules/files/module.go` — inject real notifier (replace `fileport.NoopPush()`)
- [X] T017 Call `NotifyConfigurationChanged` from `serverBackendGo/internal/modules/configurations/application/service.go` after successful save
- [X] T018 [P] Call notifier from `serverBackendGo/internal/modules/devices/application/service.go` on `applicationSettings/notify` path
- [X] T019 [P] Call notifier from `serverBackendGo/internal/modules/files/application/service.go` when configuration links change
- [X] T020 Verify `cd serverBackendGo && go test ./internal/platform/push/... ./internal/modules/configurations/...` passes

**Checkpoint**: Push notifier wired — US1 implementation can complete integration testing.

---

## Phase 3: User Story 1 — إشعارات الأجهزة الفورية (Priority: P1) 🎯 MVP

**Goal**: Saving configuration or notifying a device enqueues real `pushmessages` rows (FR-001, FR-002, FR-004).

**Independent Test**: `quickstart.md` §3 — PUT configuration then `GET /rest/notifications/device/hmdm-001` shows pending message.

### Tests for User Story 1

- [X] T021 [P] [US1] Add `serverBackendGo/internal/modules/configurations/application/service_test.go` case — save triggers notifier mock once per device
- [X] T022 [P] [US1] Add smoke step to `specs/011-complete-migration-gaps/quickstart.md` §3 if notifier flags documented

### Implementation for User Story 1

- [X] T023 [US1] Ensure `serverBackendGo/internal/modules/configurations/application/service.go` logs enqueue errors without failing HTTP response (best-effort FR-004)
- [X] T024 [US1] Update `serverBackendGo/docs/parity/configurations.md` — remove NoopPush note; mark push notify **Done**
- [X] T025 [P] [US1] Update `serverBackendGo/docs/parity/devices.md` — `applicationSettings/notify` uses queue not stub
- [X] T026 [P] [US1] Update `serverBackendGo/docs/parity/files.md` — configuration link push **Done**
- [X] T027 [US1] Flip US1 rows in `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §6.3–§6.4 (configurations/devices push)

**Checkpoint**: MVP push path works for admin save + agent poll.

---

## Phase 4: User Story 2 — الرسائل المجدولة تُرسل تلقائياً (Priority: P1)

**Goal**: Background worker processes `plugin_push_schedule` every 60s (FR-003).

**Independent Test**: `quickstart.md` §4 — scheduled task fires within 2 minutes.

### Tests for User Story 2

- [X] T028 [P] [US2] Add `serverBackendGo/internal/modules/plugins/push/application/schedule_runner_test.go` — due task enqueues N messages and marks processed
- [X] T029 [P] [US2] Add `serverBackendGo/internal/modules/plugins/push/adapter/persistence/postgres/schedule_repo_test.go` for `findMatchingTime` SQL

### Implementation for User Story 2

- [X] T030 [US2] Add `serverBackendGo/internal/modules/plugins/push/port/schedule_repository.go` — `FindDue`, `MarkProcessed` per Java `PushScheduleDAO`
- [X] T031 [US2] Implement `serverBackendGo/internal/modules/plugins/push/adapter/persistence/postgres/schedule_repo.go`
- [X] T032 [US2] Implement `serverBackendGo/internal/modules/plugins/push/application/schedule_runner.go` — resolve scope via `internal/modules/plugins/shared/targets`
- [X] T033 [US2] Create `serverBackendGo/internal/app/scheduler.go` — ticker + graceful shutdown on context cancel
- [X] T034 [US2] Start scheduler from `serverBackendGo/cmd/server/main.go` when `MODULE_PLUGINS_PUSH_ENABLED` and plugin `push` enabled
- [X] T035 [US2] Update `serverBackendGo/docs/parity/push.md` — remove schedule cron **Partial** note
- [X] T036 [US2] Update `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §5 item 4 and §6.5 push cron to ✅

**Checkpoint**: Scheduled push tasks send without manual intervention.

---

## Phase 5: User Story 3 — رفع أيقونات (Priority: P2)

**Goal**: `POST /rest/private/icon-files` uploads square PNG icons (FR-005).

**Independent Test**: `quickstart.md` §5 — multipart upload returns `uploadedfiles` row.

### Tests for User Story 3

- [X] T037 [P] [US3] Add `serverBackendGo/internal/modules/icons/application/icon_file_test.go` — rejects non-square image with `error.icon.dimension.invalid`

### Implementation for User Story 3

- [X] T038 [US3] Add domain type `serverBackendGo/internal/modules/icons/domain/uploaded_file.go` matching Java `UploadedFile` JSON
- [X] T039 [US3] Extend `serverBackendGo/internal/modules/icons/port/repository.go` — `InsertUploadedFile` if missing
- [X] T040 [US3] Implement `serverBackendGo/internal/modules/icons/application/icon_file.go` — validate square, resize 144px, write PNG under `FILES_DIRECTORY`
- [X] T041 [US3] Add `serverBackendGo/internal/modules/icons/adapter/http/icon_file_handler.go` — register `POST /rest/private/icon-files`
- [X] T042 [US3] Wire icon-files route in `serverBackendGo/internal/modules/icons/module.go` (or dedicated `iconfiles` subgroup)
- [X] T043 [P] [US3] Create `serverBackendGo/docs/parity/icon-files.md` per `contracts/icon-files-api.md`
- [X] T044 [US3] Update `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §4 `IconFileResource` row to ✅

**Checkpoint**: Icon file upload API matches Java contract.

---

## Phase 6: User Story 4 — deviceinfo & devicelog gaps (Priority: P2)

**Goal**: Missing plugin export/search/rules endpoints (FR-006, FR-007).

**Independent Test**: `quickstart.md` §6 — export endpoints return files, not Gin 404.

### Tests for User Story 4

- [ ] T045 [P] [US4] Add `serverBackendGo/internal/modules/plugins/deviceinfo/application/export_test.go` for filter → row count
- [ ] T046 [P] [US4] Add `serverBackendGo/internal/modules/plugins/devicelog/application/export_test.go` for CSV header row

### Implementation for User Story 4 — deviceinfo

- [ ] T047 [US4] Audit Java `DeviceInfoResource.java` export/search/device SQL; note missing tables in `specs/011-complete-migration-gaps/data-model.md`
- [ ] T048 [US4] Create `serverBackendGo/db/migrations/000011_deviceinfo_export.up.sql` and `.down.sql` **only if** GPS/WiFi tables absent in dev DB
- [ ] T049 [US4] Extend `serverBackendGo/internal/modules/plugins/deviceinfo/port/repository.go` — search by device, export query
- [ ] T050 [US4] Implement repos in `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/persistence/postgres/`
- [ ] T051 [US4] Add handlers in `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/http/handler.go` — `POST .../private/search/device`, `POST .../private/export`, `GET .../deviceinfo-plugin-settings/device/:deviceNumber`
- [ ] T052 [US4] Update `serverBackendGo/docs/parity/plugins-deviceinfo.md` — mark gap endpoints **Done**

### Implementation for User Story 4 — devicelog

- [ ] T053 [P] [US4] Audit Java `DeviceLogResource.java` for `search/export` and `rules/{deviceNumber}`
- [ ] T054 [US4] Extend `serverBackendGo/internal/modules/plugins/devicelog/port/repository.go` and postgres adapter for export + rules
- [ ] T055 [US4] Add handlers in `serverBackendGo/internal/modules/plugins/devicelog/adapter/http/handler.go` — `POST .../private/search/export`, `GET .../log/rules/:deviceNumber`
- [ ] T056 [US4] Update `serverBackendGo/docs/parity/plugins-devicelog.md` — mark gap endpoints **Done**
- [ ] T057 [US4] Update `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §8 partial endpoints for deviceinfo/devicelog to ✅

**Checkpoint**: Plugin monitoring UIs can call export/rules APIs.

---

## Phase 7: User Story 5 — تدقيق وتزامن وعملاء (Priority: P3)

**Goal**: Audit middleware, sync hooks, customer bootstrap, storage quota (FR-009–FR-012).

**Independent Test**: DELETE device creates audit row; sync returns plugin fields; new customer has default config.

### Tests for User Story 5

- [ ] T058 [P] [US5] Add `serverBackendGo/internal/platform/audit/middleware_test.go` — records POST with principal
- [ ] T059 [P] [US5] Add `serverBackendGo/internal/platform/synchooks/registry_test.go` — merge order
- [ ] T060 [P] [US5] Add `serverBackendGo/internal/modules/customers/application/bootstrap_test.go` — create copies template config

### Implementation for User Story 5 — audit

- [ ] T061 [US5] Create `serverBackendGo/internal/platform/audit/middleware.go` — insert `plugin_audit_log` on private routes; skip `/swagger`, health
- [ ] T062 [US5] Register audit middleware on `RouteGroups.Private` in `serverBackendGo/internal/platform/httpx/router.go` or app bootstrap
- [ ] T063 [US5] Update `serverBackendGo/docs/parity/plugins-audit.md` — document auto-capture **Done**

### Implementation for User Story 5 — sync hooks

- [ ] T064 [US5] Create `serverBackendGo/internal/platform/synchooks/registry.go` — `Register`, `ApplyAll` interfaces
- [ ] T065 [US5] Integrate hook merge in `serverBackendGo/internal/modules/sync/application/build_response.go` after core payload built
- [ ] T066 [P] [US5] Register no-op or deviceinfo hook stub from `serverBackendGo/internal/modules/plugins/deviceinfo/module.go` when enabled

### Implementation for User Story 5 — customers & quota

- [ ] T067 [US5] Implement tenant bootstrap in `serverBackendGo/internal/modules/customers/application/service.go` on create — copy default configuration/devices per Java `CustomerResource`
- [ ] T068 [US5] Add `serverBackendGo/internal/shared/storage/quota.go` — sum `uploadedfiles` sizes vs `customers.sizeLimit`
- [ ] T069 [US5] Enforce quota in `serverBackendGo/internal/modules/configfiles/application/service.go` and `files/application/service.go` and `icons/application/icon_file.go`
- [ ] T070 [US5] Update `serverBackendGo/docs/parity/customers.md` and `configfiles.md` — bootstrap/quota **Done**
- [ ] T071 [US5] Update `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §6.2 customers and §7 audit/sync rows

**Checkpoint**: Compliance and tenant onboarding parity improved.

---

## Phase 8: User Story 6 — وحدات عامة ووكلاء (Priority: P4)

**Goal**: stats, videos, updates APK, summary charts, devices search enrichment, agent file serving (FR-013–FR-018).

**Independent Test**: `quickstart.md` §7 + public stats PUT; devices search accepts new filters.

### Tests for User Story 6

- [ ] T072 [P] [US6] Add `serverBackendGo/internal/modules/stats/application/service_test.go` — insert usage row
- [ ] T073 [P] [US6] Add `serverBackendGo/internal/modules/devices/application/search_test.go` — filter `mdmMode` applied in SQL

### Implementation for User Story 6 — stats & videos

- [ ] T074 [US6] Create `serverBackendGo/db/migrations/000011_usage_stats.up.sql` and `.down.sql` from Java `UsageStats` domain columns
- [ ] T075 [US6] Scaffold `serverBackendGo/internal/modules/stats/` — domain, port, application, adapter/http handler `PUT /rest/public/stats`
- [ ] T076 [US6] Register `stats` module in `serverBackendGo/internal/app/modules.go` with `MODULE_STATS_ENABLED`
- [ ] T077 [US6] Scaffold `serverBackendGo/internal/modules/videos/` — `POST` upload, `GET /{fileName}` stream per `contracts/public-agent-gaps-api.md`
- [ ] T078 [US6] Register `videos` module with `VIDEO_DIRECTORY` env

### Implementation for User Story 6 — updates, summary, devices, static files

- [ ] T079 [US6] Complete remote APK download in `serverBackendGo/internal/modules/updates/application/service.go` per Java `UpdateResource`
- [ ] T080 [US6] Implement `sendStats` persistence in updates handler using `stats` repo or shared DAO
- [ ] T081 [US6] Extend `serverBackendGo/internal/modules/summary/adapter/persistence/postgres/summary_repo.go` — use `devicestatuses` when table exists
- [ ] T082 [US6] Extend `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go` — filters `mdmMode`, `launcherVersion`, `deviceStatuses`
- [ ] T083 [US6] Enrich devices search response in `serverBackendGo/internal/modules/devices/application/service.go` — nested apps/files map per parity
- [ ] T084 [US6] Add agent file download route in `serverBackendGo/internal/platform/httpx/router.go` or `files` module — `GET` under customer files path per Java servlet
- [ ] T085 [P] [US6] Create `serverBackendGo/docs/parity/stats.md` and `videos.md`
- [ ] T086 [P] [US6] Update `serverBackendGo/docs/parity/updates.md`, `summary.md`, `devices.md`
- [ ] T087 [US6] Mark `StatsResource`, `VideosResource` ✅ in `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §4–§5

**Checkpoint**: P3 public/agent gaps closed or explicitly ⊘ in parity.

---

## Phase 9: Polish & Cross-Cutting

**Purpose**: Documentation, swagger, migration roadmap, regression.

- [ ] T088 Add **Phase 9** row to `serverBackendGo/docs/MIGRATION.md` with modules and parity links
- [ ] T089 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — post-gap hardening complete / remaining ⊘ items
- [ ] T090 Run `cd serverBackendGo && make swagger` and commit regenerated `internal/platform/httpx/swagger/*`
- [ ] T091 Run full `cd serverBackendGo && go test ./...` and document results in `specs/011-complete-migration-gaps/quickstart.md` §7
- [ ] T092 [P] Final pass `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` — SC-001 checklist: all P0/P1 items ✅ or ⊘ with reason
- [ ] T093 [P] Verify Phases 1–8 smoke still passes (`serverBackendGo/scripts/dev.sh` or manual curl list from Phase 7 quickstart)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup — **blocks all user stories**
- **US1 (Phase 3)**: Depends on Foundational (T008–T019)
- **US2 (Phase 4)**: Depends on Foundational (shared `MessageQueue`); can run parallel to US1 after T013
- **US3 (Phase 5)**: Depends on Foundational; independent of US1/US2
- **US4 (Phase 6)**: Depends on Foundational; independent of US1–US3
- **US5 (Phase 7)**: Depends on Foundational; sync hooks may reference US4 deviceinfo stub (T066)
- **US6 (Phase 8)**: Depends on Foundational; stats migration (T074) before updates sendStats (T080)
- **Polish (Phase 9)**: Depends on desired user stories complete

### User Story Dependencies

| Story | Depends on | Can parallel with |
|-------|------------|-------------------|
| US1 | Phase 2 | US2 after T013 |
| US2 | Phase 2, targets | US1 |
| US3 | Phase 2 | US4, US6 stats |
| US4 | Phase 2, optional migration | US3 |
| US5 | Phase 2 | US6 (partial) |
| US6 | Phase 2, T074 for T080 | US3, US4 |

### Within Each User Story

- Unit tests (T021, T028, etc.) alongside or immediately after implementation tasks
- Parity + gap tracker updates are last tasks per story

---

## Parallel Execution Examples

### After Phase 2 completes

```text
Developer A: US1 (T023–T027) + US2 (T030–T036)
Developer B: US3 (T038–T044)
Developer C: US4 deviceinfo (T047–T052) then devicelog (T053–T056)
```

### US4 internal parallel

```text
T047 audit Java + T053 audit Java [P]
T050 deviceinfo postgres + T054 devicelog postgres [P]
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

1. Complete Phase 1–2
2. Complete US1 (push notifier) — **demo value immediately**
3. Complete US2 (schedule worker)
4. **STOP and VALIDATE** via `quickstart.md` §3–§4
5. Deploy/demo before P2 plugin gaps

### Incremental Delivery

| Increment | Stories | Outcome |
|-----------|---------|---------|
| MVP | US1 + US2 | Operational push parity |
| P2 | US3 + US4 | Admin content + plugin monitoring |
| P3 | US5 | Audit, sync hooks, tenants |
| P4 | US6 | Stats/videos/updates/search polish |

### Suggested task counts

| Phase | Tasks | Story |
|-------|-------|-------|
| Setup | 5 | — |
| Foundational | 15 | — |
| US1 | 7 | P1 MVP |
| US2 | 9 | P1 |
| US3 | 8 | P2 |
| US4 | 13 | P2 |
| US5 | 14 | P3 |
| US6 | 16 | P4 |
| Polish | 6 | — |
| **Total** | **93** | |

---

## Notes

- Mailchimp, xtra plugin UI, `PublicFilesResource`, user impersonate/superadmin remain **out of scope** — do not add tasks unless spec amended.
- MQTT/FCM client deferred per `research.md` R1 — polling queue is sufficient for Phase 9 v1.
- Mark tasks `[X]` only after smoke + parity update for that slice.
