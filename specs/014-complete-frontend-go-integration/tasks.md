---
description: "Task list for 014 — Complete React ↔ Go frontend integration gaps"
---

# Tasks: إكمال تكامل React ↔ Go

**Input**: `specs/014-complete-frontend-go-integration/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Feature **013** applied — `make migrate` through `000017`; [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md); Go dev server + React proxy

**Tests**: Application/repo unit tests per constitution IV where logic is non-trivial; manual UAT via `quickstart.md` (no separate contract-test phase unless noted).

**Organization**: Tasks grouped by user story (US1–US6) for independent delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Go modules: `serverBackendGo/internal/modules/<name>/`
- Frontend: `frontend/src/features/<name>/`
- Parity: `serverBackendGo/docs/parity/`
- Contracts: `specs/014-complete-frontend-go-integration/contracts/`
- Tracker: `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 UAT checklist

---

## Phase 1: Setup

**Purpose**: Confirm 013 schema, contracts, and integration baseline before code changes.

- [x] T001 Verify `specs/014-complete-frontend-go-integration/spec.md` gap matrix against `FRONTEND-GO-BACKEND-INTEGRATION.md` §5–§7 and plan waves A/B/C
- [x] T002 [P] Read `specs/014-complete-frontend-go-integration/contracts/settings-api.md`, `configurations-api.md`, `icons-api.md` and note delta vs current handlers
- [x] T003 [P] Read `specs/014-complete-frontend-go-integration/contracts/sync-device-status.md` and `stats-api.md` for P2 scope
- [x] T004 Run `cd serverBackendGo && go build ./... && go test ./...` and record baseline in `specs/014-complete-frontend-go-integration/quickstart.md` notes if failures exist
- [x] T005 [P] Confirm DB at migration `000017` (`./scripts/db-up.sh && make migrate`) — prerequisite for all smoke tests

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Environment and docs ready; no new migrations in 014.

**⚠️ CRITICAL**: No user story work until T006–T008 complete.

- [x] T006 Document prerequisite in `serverBackendGo/docs/MIGRATION.md` — feature 014 (frontend integration) depends on 013 schema `000011`–`000017`
- [x] T007 [P] Add `MODULE_STATS_ENABLED` to `serverBackendGo/.env.example` with comment (default true in dev) per `contracts/stats-api.md`
- [x] T008 [P] Add row to `serverBackendGo/docs/NEXT_STEPS.md` — 014 in progress (settings UI, config MDM, icons upload, sync status, stats)
- [x] T009 Verify `serverBackendGo/internal/modules/settings/domain/settings.go` already lists `000015` fields — gap is repo/handler/frontend only (note in plan if needed)

**Checkpoint**: DB green at `000017`; team aligned on no DDL in 014.

---

## Phase 3: User Story 1 — إعدادات المستأجر الكاملة (Priority: P1) 🎯 MVP start

**Goal**: Tenant settings (`000015`) round-trip in Go API + Settings UI (FR-001–FR-003).

**Independent Test**: `quickstart.md` §2 US1 — change `phoneNumberFormat` and `customPropertyName1`, save, reload; role columns tab shows custom labels.

### Tests for User Story 1

- [x] T010 [P] [US1] Extend `serverBackendGo/internal/modules/settings/application/service_test.go` — GET settings returns `phoneNumberFormat` and `customPropertyName1` from repo stub
- [x] T011 [P] [US1] Add `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo_test.go` — SELECT/UPDATE includes `newdevicegroupid`, `phonenumberformat`, `custompropertyname1` columns (sqlmock)

### Implementation for User Story 1 (Backend)

- [x] T012 [US1] Update `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo.go` — SELECT all `000015` columns into `domain.Settings` on GetByCustomerID
- [x] T013 [US1] Update `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo.go` — UPDATE/INSERT misc save persists `000015` columns (merge, do not null-out omitted keys)
- [x] T014 [US1] Update `serverBackendGo/internal/modules/settings/adapter/http/handler.go` — `SaveMisc` decodes and passes tenant fields; `SaveLang` must not clear tenant columns
- [x] T015 [US1] Update `serverBackendGo/internal/modules/settings/application/service.go` if normalization/defaults needed for empty `phoneNumberFormat`

### Implementation for User Story 1 (Frontend)

- [x] T016 [P] [US1] Extend `frontend/src/features/settings/types.ts` — add tenant fields per `data-model.md` (`newDeviceGroupId`, `phoneNumberFormat`, `customPropertyName*`, `customMultiline*`, `customSend*`, `desktopHeaderTemplate`, `sendDescription`)
- [x] T017 [P] [US1] Extend `frontend/src/features/settings/settingsService.ts` — map tenant fields in GET/POST misc payload per `contracts/settings-api.md`
- [x] T018 [US1] Update `frontend/src/features/settings/SettingsPage.tsx` — form sections for tenant/misc fields (group picker, phone format, custom labels, desktop template, send flags)
- [x] T019 [US1] Update `frontend/src/features/settings/SettingsRoleColumnsTab.tsx` — display `customPropertyName1`–`3` as column labels when choosing device list columns

### Documentation for User Story 1

- [x] T020 [P] [US1] Update `serverBackendGo/docs/parity/settings.md` — tenant `000015` fields GET/POST **Done**
- [x] T021 [P] [US1] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 — mark Settings tenant UAT `[x]` when smoke passes

**Checkpoint**: US1 independently testable via quickstart §2.

---

## Phase 4: User Story 2 — محرر التكوين MDM كامل (Priority: P1)

**Goal**: Full `settingsjson` + CAP + `remove`/`longTap` round-trip (FR-004–FR-007).

**Independent Test**: `quickstart.md` §2 US2 — change `kioskMode`, `wifi`, `skipVersionCheck`, `remove`, `longTap`; save; GET `/{id}` matches.

### Tests for User Story 2

- [x] T022 [P] [US2] Extend `serverBackendGo/internal/modules/configurations/application/service_test.go` — save/load preserves policy key from `settingsjson` merge
- [x] T023 [P] [US2] Add or extend repo test in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/` — `ListConfigurationApplications` returns `skipVersionCheck`, `remove`, `longTap`

### Implementation for User Story 2 (Backend)

- [x] T024 [US2] Update `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` `GetByID` — after load, `json.Unmarshal` `settingsjson` into flat policy fields on `domain.Configuration` (camelCase keys per Java/React)
- [x] T025 [US2] Update `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` `Save` — marshal non-column MDM fields from struct into `settingsjson` merge map; never replace with `{}` on partial body
- [x] T026 [US2] Update `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` `ListConfigurationApplications` — `LEFT JOIN configurationapplicationparameters` for `skipversioncheck`; SELECT `ca.remove`, `ca.longtap`
- [x] T027 [US2] Verify `serverBackendGo/internal/modules/configurations/domain/configuration.go` documents all policy + app link fields used by handler JSON tags
- [x] T028 [US2] Update `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` if request/response mapping omits nested policy keys on PUT/GET

### Implementation for User Story 2 (Frontend)

- [x] T029 [P] [US2] Audit `frontend/src/features/configurations/types.ts` — ensure MDM policy keys match Go/Java camelCase list in `contracts/configurations-api.md`
- [x] T030 [P] [US2] Update `frontend/src/features/configurations/configurationNormalize.ts` — round-trip `skipVersionCheck`, `remove`, `longTap` on applications array
- [x] T031 [US2] Update `frontend/src/features/configurations/configurationService.ts` — send full editor payload on PUT including policy tabs and app flags
- [x] T032 [US2] Update `frontend/src/features/configurations/ConfigurationApplicationsTab.tsx` — UI toggles for skipVersionCheck, remove, longTap bound to form state
- [x] T033 [US2] Spot-check `frontend/src/features/configurations/ConfigurationForm.tsx` and policy tabs — bound fields persist after save/reopen

### Documentation for User Story 2

- [x] T034 [P] [US2] Update `serverBackendGo/docs/parity/configurations.md` — `settingsjson` flatten, CAP, remove/longTap **Done**
- [x] T035 [P] [US2] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 — mark Configuration MDM UAT `[x]`

**Checkpoint**: US2 independently testable; P1 config editor no silent data loss.

---

## Phase 5: User Story 3 — رفع الأيقونات من الواجهة (Priority: P1)

**Goal**: Icons page uses `POST /icon-files` then `PUT /icons` (FR-008).

**Independent Test**: `quickstart.md` §2 US3 — upload PNG, create icon, list shows preview without manual fileId.

### Implementation for User Story 3

- [x] T036 [P] [US3] Verify `serverBackendGo/internal/modules/icons/adapter/http/icon_file_handler.go` matches `contracts/icons-api.md` response shape (`fileId`, optional `url`)
- [x] T037 [P] [US3] Add `uploadIconFile(file: File)` to `frontend/src/features/icons/iconsService.ts` — multipart `POST /rest/private/icon-files`
- [x] T038 [US3] Update `frontend/src/features/icons/IconsPage.tsx` — file input + upload flow; remove manual numeric `fileId` requirement
- [x] T039 [US3] Add error toasts for invalid file type/size on upload failure in `frontend/src/features/icons/IconsPage.tsx`

### Documentation for User Story 3

- [x] T040 [P] [US3] Update `serverBackendGo/docs/parity/icons.md` — frontend upload flow documented **Done**
- [x] T041 [P] [US3] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 — mark Icons upload UAT `[x]`

**Checkpoint**: **MVP (Wave A) complete** — US1 + US2 + US3; SC-001 achievable.

---

## Phase 6: User Story 4 — حالة التثبيت تعكس الوكيل (Priority: P2)

**Goal**: `devicestatuses` updated from `POST /public/sync/info` (FR-009).

**Independent Test**: `quickstart.md` §3 US4 — sync payload changes status → device search `installationStatus` filter matches.

### Tests for User Story 4

- [x] T042 [P] [US4] Add `serverBackendGo/internal/modules/sync/application/service_test.go` case — after `UpdateInfo`, `DeviceStatusUpserter` called with derived SUCCESS/FAILURE

### Implementation for User Story 4

- [x] T043 [US4] Add `DeviceStatusUpserter` interface to `serverBackendGo/internal/modules/sync/port/repository.go` (or new `port/device_status.go`)
- [x] T044 [US4] Implement upsert in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_status_repo.go` — `INSERT ... ON CONFLICT (deviceid) DO UPDATE` per `contracts/sync-device-status.md`
- [x] T045 [US4] Add status derivation helper in `serverBackendGo/internal/modules/sync/application/service.go` — aggregate applications/files from sync `info` payload (Java parity rules)
- [x] T046 [US4] Wire upsert call in `serverBackendGo/internal/modules/sync/application/service.go` after successful device info update
- [x] T047 [US4] Register repo in `serverBackendGo/internal/modules/sync/module.go`
- [x] T048 [US4] Run `quickstart.md` §3 US4 smoke — sync then filter devices by `installationStatus`

### Documentation for User Story 4

- [x] T049 [P] [US4] Create or update `serverBackendGo/docs/parity/sync.md` — device status side effect on `/sync/info` **Done**
- [x] T050 [P] [US4] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` §6 — sync → `devicestatuses` row ✅

**Checkpoint**: US4 independently testable; SC-004 path enabled in QA.

---

## Phase 7: User Story 5 — إحصائيات الخادم (Priority: P2)

**Goal**: New `stats` module — `PUT /rest/public/stats` → `usagestats` (FR-010).

**Independent Test**: `quickstart.md` §3 US5 — PUT sample payload → row in `usagestats`; repeat upserts same `(ts, instanceid)`.

### Tests for User Story 5

- [x] T051 [P] [US5] Add `serverBackendGo/internal/modules/stats/application/service_test.go` — upsert validates numeric fields ≥ 0

### Implementation for User Story 5

- [x] T052 [US5] Create `serverBackendGo/internal/modules/stats/domain/usage_stats.go` — DTO per `data-model.md`
- [x] T053 [US5] Create `serverBackendGo/internal/modules/stats/port/repository.go` — `UpsertUsageStats`
- [x] T054 [US5] Create `serverBackendGo/internal/modules/stats/adapter/persistence/postgres/stats_repo.go` — upsert on `(ts, instanceid)`
- [x] T055 [US5] Create `serverBackendGo/internal/modules/stats/application/service.go` and `adapter/http/handler.go` — `PUT /rest/public/stats`, public, no auth
- [x] T056 [US5] Create `serverBackendGo/internal/modules/stats/module.go` — register routes; gate with `MODULE_STATS_ENABLED` in `serverBackendGo/internal/platform/config`
- [x] T057 [US5] Wire module in `serverBackendGo/internal/app/modules.go` (or equivalent registry)
- [x] T058 [US5] Run `quickstart.md` §3 US5 curl smoke against running `make dev`

### Documentation for User Story 5

- [x] T059 [P] [US5] Create `serverBackendGo/docs/parity/stats.md` — `PUT /rest/public/stats` **Done**
- [x] T060 [P] [US5] Update `JAVA-GO-BACKEND-GAPS.md` — stats endpoint row ✅ if listed

**Checkpoint**: US5 independently testable; external heartbeat consumers unblocked.

---

## Phase 8: User Story 6 — تحسينات ثانوية (Priority: P3)

**Goal**: Updates apply, hints mark shown; optional device list columns + summary monthly series (FR-011–FR-014).

**Independent Test**: Each sub-item in `quickstart.md` §4 — separate pass per acceptance scenario.

### Implementation for User Story 6 (P3 — required)

- [x] T061 [P] [US6] Update `frontend/src/features/updates/updatesService.ts` — call `POST /rest/private/update` on user confirm apply
- [x] T062 [US6] Update `frontend/src/features/updates/UpdatesPage.tsx` — wire apply button to service; show success/error toast from envelope
- [x] T063 [P] [US6] Update `frontend/src/features/hints/hintsService.ts` — `markHintShown` → `POST /rest/private/hints/history` with hint id
- [x] T064 [US6] Update `frontend/src/features/hints/HintsPage.tsx` — call mark shown on dialog close/dismiss

### Implementation for User Story 6 (P2 optional)

- [ ] T065 [P] [US6] Update `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go` — expose `model`, `batteryLevel`, `androidVersion` from `infojson` in search row DTO (FR-011)
- [ ] T066 [P] [US6] Update `frontend/src/features/devices/types.ts` and `DevicesPage.tsx` columns if new fields returned
- [ ] T067 [US6] Update `serverBackendGo/internal/modules/summary/adapter/persistence/postgres/summary_repo.go` — real `devicesEnrolledMonthly` series from enrollment dates (FR-012)
- [ ] T068 [P] [US6] Update `serverBackendGo/docs/parity/devices.md` and `serverBackendGo/docs/parity/summary.md` if T065–T067 implemented

**Checkpoint**: US6 P3 items done; optional P2 enrichments marked in integration doc.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Regression, docs, and feature closure.

- [x] T069 Run `cd serverBackendGo && go test ./...` and fix regressions from all phases
- [x] T070 [P] Execute full `specs/014-complete-frontend-go-integration/quickstart.md` Waves A–C; fix gaps found
- [x] T071 [P] Update `FRONTEND-GO-BACKEND-INTEGRATION.md` §1 field parity table and §11 priorities — reflect 014 closure
- [x] T072 Update `serverBackendGo/docs/MIGRATION.md` and `serverBackendGo/docs/NEXT_STEPS.md` — mark 014 complete or note remaining optional items
- [ ] T073 [P] Mark completed items in `specs/014-complete-frontend-go-integration/spec.md` success criteria SC-001–SC-006 where verified
- [ ] T074 [P] Regenerate Swagger if handlers added: `serverBackendGo` swag script / `internal/platform/httpx/swagger/` per project convention

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — **blocks all user stories**
- **US1–US3 (Phases 3–5)**: Depend on Foundational — **MVP Wave A**; US2/US3 can parallelize after US1 backend if staffed
- **US4–US5 (Phases 6–7)**: Depend on Foundational; independent of US1–US3 (can run in parallel with Wave A)
- **US6 (Phase 8)**: Depends on Foundational; P3 frontend tasks independent of US4–US5
- **Polish (Phase 9)**: Depends on desired user stories complete

### User Story Dependencies

| Story | Depends on | Independent test |
|-------|------------|------------------|
| US1 | Phase 2 | quickstart §2 US1 |
| US2 | Phase 2 | quickstart §2 US2 |
| US3 | Phase 2 | quickstart §2 US3 |
| US4 | Phase 2, 013 `devicestatuses` | quickstart §3 US4 |
| US5 | Phase 2, 013 `usagestats` table | quickstart §3 US5 |
| US6 | Phase 2 | quickstart §4 per item |

### Within Each User Story

- Tests (when listed) before or alongside implementation
- Backend persistence before handler smoke
- Frontend services before pages
- Parity + integration doc last per story

### Parallel Opportunities

- **Phase 1**: T002, T003, T005 parallel
- **Phase 2**: T007, T008 parallel
- **US1**: T010–T011 parallel; T016–T017 parallel after T012–T014
- **US2**: T022–T023 parallel; T029–T030 parallel
- **US3**: T036–T037 parallel
- **US4**: T042 parallel with T043 start; T049–T050 parallel
- **US5**: T051–T059 backend files sequential per layer but repo/handler split [P] where noted
- **US6**: T061–T063 parallel; T065–T067 parallel
- **Cross-story**: After Phase 2, **US4+US5** can run while **US1–US3** run (different modules)

---

## Parallel Example: User Story 1

```bash
# Tests in parallel:
T010 service_test.go
T011 settings_repo_test.go

# Frontend types + service in parallel (after backend T012–T014):
T016 types.ts
T017 settingsService.ts
```

---

## Parallel Example: Wave A (MVP)

```bash
# Three developers after Phase 2:
Dev A: US1 (T012–T021)
Dev B: US2 (T024–T035)
Dev C: US3 (T036–T041)
# Coordinate only on shared go test ./... at end of day
```

---

## Implementation Strategy

### MVP First (Wave A — US1 + US2 + US3)

1. Complete Phase 1 + Phase 2
2. Complete Phase 3 (US1) → validate quickstart §2 US1
3. Complete Phase 4 (US2) → validate §2 US2
4. Complete Phase 5 (US3) → validate §2 US3
5. **STOP and VALIDATE**: SC-001 (all P1 UAT in integration doc)
6. Demo to stakeholders before P2

### Incremental Delivery

1. Wave A (US1–US3) → P1 integration closed
2. Wave B (US4–US5) → agent status + stats heartbeat
3. Wave C (US6) → UX polish + optional columns/charts
4. Phase 9 → docs and regression

### Suggested MVP scope note

Skill default “US1 only” is **narrower** than this feature’s MVP: per plan **Wave A = US1–US3** (all P1 gaps). Minimum shippable increment for production UI on Go: **all three P1 stories**.

---

## Notes

- `domain/settings.go` already declares `000015` fields — do not duplicate; wire repo/handler/UI (T009, T012–T019).
- No new migrations in 014; if schema missing, fix environment (013), not new DDL.
- Commit after each user story checkpoint; optional `speckit.git.commit` after tasks generation.
- Avoid editing Java backend except as reference (`backend/server/.../SettingsResource.java`, etc.).
