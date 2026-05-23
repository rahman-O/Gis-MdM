# Tasks: Configuration–Device Sync & Admin UX

**Input**: Design documents from `/specs/016-config-sync-ux/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Golden/unit tests included per plan.md (constitution IV); not full TDD — implement mapper then add failing golden test where practical.

**Organization**: Tasks grouped by user story (US1–US4) for independent delivery and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US4 from spec.md

## Path Conventions

- Backend: `serverBackendGo/internal/modules/{configurations,sync,push}/`
- Frontend: `frontend/src/features/configurations/`
- Parity: `serverBackendGo/docs/parity/`
- Java reference: `backend/server/src/main/java/com/hmdm/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Baseline audit and contract alignment before code changes

- [x] T001 Audit Java `SyncResponse.java` and sync assembly vs Go `BuildSyncResponse` in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go`; document field gap list in `specs/016-config-sync-ux/research.md` (appendix) or inline notes for implementer
- [x] T002 [P] Review contracts in `specs/016-config-sync-ux/contracts/` against current handlers in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` and `serverBackendGo/internal/modules/sync/adapter/http/handler.go`
- [x] T003 [P] Confirm dev prerequisites in `specs/016-config-sync-ux/quickstart.md` (DB, `BASE_URL`, JWT) run on local machine

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared types and repository hooks ALL user stories depend on

**⚠️ CRITICAL**: No user story implementation until this phase is complete

- [x] T004 Extend `SyncResponse` and related DTOs in `serverBackendGo/internal/modules/sync/domain/sync.go` with MDM fields from Java parity list (kiosk*, connectivity, `restrictions`, update modes, etc.)
- [x] T005 [P] Add configuration load port method for sync mapper in `serverBackendGo/internal/modules/sync/port/repository.go` (full configuration row + `settingsjson` + scalar columns) implemented in `device_sync_repo.go`
- [x] T006 [P] Create stub `serverBackendGo/internal/modules/sync/application/sync_configuration_mapper.go` with `MapConfigurationToSync` signature and wire call site from `BuildSyncResponse` (delegate gradually)
- [x] T007 [P] Add `PolicyLocks` type and allowlist helpers in `serverBackendGo/internal/modules/configurations/domain/policy_locks.go`; integrate read/write with `configuration_json.go`
- [x] T008 Add `policyLocks` to `frontend/src/features/configurations/types.ts` and normalize helpers in `frontend/src/features/configurations/configurationNormalize.ts`

**Checkpoint**: Extended types and mapper skeleton exist; `go build ./...` passes

---

## Phase 3: User Story 1 — تطبيق سياسة التكوين على الجهاز (Priority: P1) 🎯 MVP

**Goal**: Saved configuration policy (apps, MDM flags, files) reaches enrolled devices via sync payload and optional push notify

**Independent Test**: Save configuration with `kioskMode: true` and one app → `GET /rest/public/sync/configuration/{deviceNumber}` returns matching fields; device or curl shows policy applied per quickstart §2

### Tests for User Story 1

- [x] T009 [P] [US1] Add golden test fixture and expected JSON subset in `serverBackendGo/internal/modules/sync/application/sync_configuration_mapper_test.go`
- [ ] T010 [P] [US1] Add handler-level test for sync GET returning extended fields in `serverBackendGo/internal/modules/sync/application/service_test.go`

### Implementation for User Story 1

- [x] T011 [US1] Implement full field mapping in `serverBackendGo/internal/modules/sync/application/sync_configuration_mapper.go` from configuration columns + `settingsjson` per `specs/016-config-sync-ux/contracts/sync-configuration-payload.md`
- [x] T012 [US1] Refactor `BuildSyncResponse` in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go` to use mapper for policy fields; keep applications/files queries
- [x] T013 [US1] Merge `configurationapplicationsettings` with `deviceapplicationsettings` in mapper/repo with readonly precedence per `specs/016-config-sync-ux/data-model.md`
- [x] T014 [US1] Wire real `push.PushNotifier` in `serverBackendGo/internal/modules/configurations/module.go` (replace `NoopPushNotifier`) using same pattern as `devices` module
- [x] T015 [US1] Verify `NotifyConfigurationChanged` in `serverBackendGo/internal/modules/configurations/application/service.go` enqueues `configUpdated` for devices on configuration
- [x] T016 [US1] Update `serverBackendGo/docs/parity/sync.md` with extended `SyncResponse` field list and status

**Checkpoint**: Sync JSON includes kiosk/restrictions/apps; save triggers push when module enabled

---

## Phase 4: User Story 2 — تقييد الإعدادات على مستوى التكوين (Priority: P1)

**Goal**: `policyLocks` and readonly configuration app settings prevent device from overriding locked policy

**Independent Test**: Lock `mainAppId` in UI → save → device POST app settings with different value ignored → next sync still shows configuration value per quickstart §3

### Tests for User Story 2

- [x] T017 [P] [US2] Unit tests for `policy_locks.go` allowlist and merge in `serverBackendGo/internal/modules/configurations/domain/policy_locks_test.go`
- [ ] T018 [P] [US2] Unit tests for readonly skip logic in `serverBackendGo/internal/modules/sync/application/service_test.go` (`SaveApplicationSettings`)

### Implementation for User Story 2

- [x] T019 [US2] Persist `policyLocks` in `settingsjson` on save in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` and expose on GET via `configuration_json.go`
- [x] T020 [US2] Include `policyLocks` in save/get API responses in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` per `specs/016-config-sync-ux/contracts/configurations-api.md`
- [x] T021 [US2] Apply locked scalar policy when building sync in `sync_configuration_mapper.go` (configuration wins over device drift for locked keys)
- [x] T022 [US2] Implement readonly enforcement on `POST /rest/public/sync/applicationSettings/{deviceId}` in `serverBackendGo/internal/modules/sync/application/service.go` and fix upsert in `device_sync_repo.go` (`ON CONFLICT` → proper upsert for non-readonly rows)
- [x] T023 [US2] Load configuration default readonly settings in `device_sync_repo.go` for merge into sync `applicationSettings`
- [x] T024 [US2] Update `serverBackendGo/docs/parity/configurations.md` with `policyLocks` and readonly behavior

**Checkpoint**: Locked fields and readonly app settings cannot be overridden by device POST

---

## Phase 5: User Story 3 — محرر تكوينات منظّم بتابات (Priority: P1)

**Goal**: Professional tabbed configuration editor with lock toggles and minimal redundant copy

**Independent Test**: Open `/configurations/:id` → seven tabs render → lock on main app persists after reload → save error cites tab name per `specs/016-config-sync-ux/contracts/configuration-editor-ux.md`

### Implementation for User Story 3

- [x] T025 [P] [US3] Create `frontend/src/features/configurations/ConfigurationMdmTab.tsx` extracted from inline MDM fields in `ConfigurationEditorPage.tsx`
- [x] T026 [P] [US3] Create `frontend/src/features/configurations/ConfigurationRestrictionsTab.tsx` with `restrictions` field and related lock toggles
- [x] T027 [US3] Refactor `frontend/src/features/configurations/ConfigurationEditorPage.tsx` to tab-only layout (remove duplicate MDM block ~line 285, trim `CardDescription` paragraphs)
- [x] T028 [US3] Add per-field lock toggle component bound to `configuration.policyLocks` in MDM and Restrictions tabs
- [x] T029 [US3] Update `frontend/src/features/configurations/configurationService.ts` to send/receive `policyLocks` on PUT/GET
- [x] T030 [US3] Improve save validation messages to include tab + field labels in `ConfigurationEditorPage.tsx`
- [x] T031 [P] [US3] Add unit test for `policyLocks` normalize round-trip in `frontend/src/features/configurations/configurationNormalize.test.ts` (extend or create)

**Checkpoint**: Editor matches UX contract; no placeholder “phase 1” strings

---

## Phase 6: User Story 4 — مزامنة التطبيقات والقيود بشكل موثوق (Priority: P2)

**Goal**: Application list and restriction keys in sync match Java; configuration save round-trips all MDM fields in API response

**Independent Test**: Add app to configuration → save → sync `applications` includes entry with HTTPS URL; reload editor shows same `kioskMode`/`restrictions`/`applications` as saved

### Tests for User Story 4

- [ ] T032 [P] [US4] Extend `serverBackendGo/internal/modules/configurations/domain/configuration_json_test.go` for `restrictions` + `applications[]` round-trip cases
- [ ] T033 [P] [US4] Add integration-style test comparing sync applications count/URLs after configuration app link change in `device_sync_repo_test.go`

### Implementation for User Story 4

- [ ] T034 [US4] Ensure `Save` returns full `ConfigurationResponseMap` in `serverBackendGo/internal/modules/configurations/application/service.go` and `handler.go` (no dropped MDM keys — FR-009)
- [ ] T035 [US4] Log warning when application version missing URL during sync build in `sync_configuration_mapper.go` or `device_sync_repo.go` (FR-010)
- [ ] T036 [US4] Frontend: reload configuration after save and refresh applications tab state in `ConfigurationEditorPage.tsx` / `ConfigurationApplicationsTab.tsx`
- [ ] T037 [US4] Verify `PUT /private/configurations/application/upgrade` still updates linked version used by sync in `serverBackendGo/internal/modules/configurations/application/service.go`

**Checkpoint**: Editor reload matches DB; sync applications/restrictions align with Java sample for fixture

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, migration status, end-to-end validation

- [x] T038 [P] Run full `go test ./internal/modules/configurations/... ./internal/modules/sync/...` and fix failures
- [ ] T039 [P] Run `specs/016-config-sync-ux/quickstart.md` curl steps; document any env gaps in quickstart if needed
- [ ] T040 Update `JAVA-GO-MIGRATION-STATUS.md` configurations/sync rows for 016 scope
- [ ] T041 [P] Cross-link `خطط مستقبليه/Configurations.md` from `specs/016-config-sync-ux/spec.md` if not already linked
- [ ] T042 Manual Android UAT: change kiosk + add app → sync → confirm on device (record result in PR description)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: Start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **blocks all user stories**
- **Phase 3 (US1)**: Depends on Phase 2 — **MVP**
- **Phase 4 (US2)**: Depends on Phase 2; integrates with US1 mapper (T021 after T011)
- **Phase 5 (US3)**: Depends on Phase 2 (T008 types); can parallel US1/US2 after T008
- **Phase 6 (US4)**: Depends on US1 mapper (T011–T012) and configurations save path
- **Phase 7 (Polish)**: After desired stories complete

### User Story Dependencies

| Story | Depends on | Can start after |
|-------|------------|-----------------|
| US1 | Foundational | T004–T008 done |
| US2 | Foundational + US1 mapper stub | T006+T007; full lock apply after T011 |
| US3 | Foundational (T008) | T008 done (parallel with backend) |
| US4 | US1 + configurations save | T012, T019+ |

### Within Each User Story

- Tests before or alongside mapper implementation (golden tests may need fixture DB or sqlmock)
- Domain before HTTP
- Backend policy before frontend locks display (T028 after T020 recommended)

### Parallel Opportunities

- **Phase 1**: T002 ∥ T003
- **Phase 2**: T005 ∥ T006 ∥ T007 ∥ T008 (after T004)
- **US1**: T009 ∥ T010; T014 ∥ T011 (different modules)
- **US2**: T017 ∥ T018
- **US3**: T025 ∥ T026; T031 ∥ T027
- **US4**: T032 ∥ T033
- **Polish**: T038 ∥ T039 ∥ T041

---

## Parallel Example: User Story 1

```bash
# Tests in parallel:
# T009 golden test file
# T010 service_test.go extensions

# After T011 mapper done:
# T014 push wiring (configurations/module.go) in parallel with T016 docs
```

---

## Parallel Example: User Story 3 (Frontend-only track)

```bash
# After T008 types exist:
# T025 ConfigurationMdmTab.tsx
# T026 ConfigurationRestrictionsTab.tsx
# Then T027 integrate in ConfigurationEditorPage.tsx
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1–2
2. Complete Phase 3 (US1): sync payload + push
3. **STOP and VALIDATE**: `quickstart.md` §1–§2
4. Demo device receiving kiosk/app list

### Incremental Delivery

1. Foundation → **US1** (sync works) → **US2** (locks) → **US3** (editor UX) → **US4** (hardening) → Polish
2. US3 can start after T008 for parallel frontend work while backend finishes US1

### Suggested MVP Scope

- **Minimum**: Phase 1 + 2 + Phase 3 (US1) — 16 tasks (T001–T016)
- **Production-ready P1**: Add Phase 4 + 5 (US2 + US3) — through T031
- **Full feature**: Include Phase 6–7 (US4 + polish) — through T042

---

## Task Summary

| Phase | Story | Task IDs | Count |
|-------|-------|----------|-------|
| 1 Setup | — | T001–T003 | 3 |
| 2 Foundational | — | T004–T008 | 5 |
| 3 US1 | US1 | T009–T016 | 8 |
| 4 US2 | US2 | T017–T024 | 8 |
| 5 US3 | US3 | T025–T031 | 7 |
| 6 US4 | US4 | T032–T037 | 6 |
| 7 Polish | — | T038–T042 | 5 |
| **Total** | | **T001–T042** | **42** |

**Format validation**: All tasks use `- [ ]`, task ID, optional `[P]`, story label where required, and file paths.

---

## Notes

- No new DB migration for v1 (`policyLocks` in `settingsjson`)
- Java reference paths: `SyncResponse.java`, `ConfigurationResource`, `ConfigurationMapper.xml`
- Keep REST paths unchanged; update parity docs only
- Commit per task group or per user story checkpoint
