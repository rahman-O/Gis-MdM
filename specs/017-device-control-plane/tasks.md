# Tasks: Device Control Plane

**Input**: Design documents from `/specs/017-device-control-plane/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Unit/golden tests per constitution IV and plan.md (compile, tree path, publish impact); not full TDD.

**Organization**: Tasks grouped by user story (US1–US6). Cross-cutting sync, migration, and events in final phases.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US6 from spec.md

## Path Conventions

- Backend: `serverBackendGo/internal/modules/{device_tree,profiles,enrollment_routes,devices,sync,qrcode,push}/`
- Migrations: `serverBackendGo/db/migrations/`
- Frontend: `frontend/src/features/{device-tree,profiles,enrollment-routes,devices,onboarding}/`
- Parity: `serverBackendGo/docs/parity/device-control-plane.md`
- Blueprint gates: `DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md` §20

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Audit, parity stub, and contract alignment before module work

- [x] T001 Map plan phases to blueprint §18.10 delivery order; note execution caveat (US2 full test after US4) in `specs/017-device-control-plane/plan.md` appendix if needed
- [x] T002 [P] Audit Java `QRCodeResource.java` and `SyncResource.java` vs Go `serverBackendGo/internal/modules/qrcode/` and `serverBackendGo/internal/modules/sync/` for enroll + sync fields
- [x] T003 [P] Review all files in `specs/017-device-control-plane/contracts/` against existing `configurations` and `devices` handlers
- [x] T004 Create parity stub `serverBackendGo/docs/parity/device-control-plane.md` with endpoint checklist from contracts
- [x] T005 [P] Verify `specs/017-device-control-plane/quickstart.md` prerequisites (DB, JWT, agent) on local dev environment

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Module scaffolds, feature flags, and app wiring — MUST complete before user stories

**⚠️ CRITICAL**: No user story implementation until this phase is complete

- [x] T006 Scaffold `serverBackendGo/internal/modules/device_tree/` (`module.go`, `domain/`, `port/`, `application/`, `adapter/http/handler.go`, `adapter/persistence/postgres/`)
- [x] T007 [P] Scaffold `serverBackendGo/internal/modules/profiles/` with same layer layout
- [x] T008 [P] Scaffold `serverBackendGo/internal/modules/enrollment_routes/` with same layer layout
- [x] T009 Register `device_tree`, `profiles`, `enrollment_routes` modules in `serverBackendGo/internal/app/app.go` (or module registry) behind feature flags
- [x] T010 [P] Add `MODULE_DEVICE_TREE_ENABLED`, `MODULE_PROFILES_ENABLED`, `MODULE_ENROLLMENT_ROUTES_ENABLED` to `serverBackendGo/internal/config/config.go` and `serverBackendGo/.env.example`
- [x] T011 [P] Add frontend route placeholders in `frontend/src/App.tsx` for `/profiles`, `/enrollment-routes` (stub pages OK until story phases)
- [x] T012 Extend `serverBackendGo/internal/modules/devices/domain/device.go` with `AgentID`, `TreeNodeID`, `EnrollmentRouteID`, `EnrollmentState` fields (types only, no DB yet)

**Checkpoint**: `go build ./...` passes; modules registered but may return not-implemented until story phases

---

## Phase 3: User Story 1 — تنظيم الأجهزة في شجرة (Priority: P1) 🎯 MVP

**Goal**: Folder tree for device placement; filter table by node; move device; delete folder with mandatory relocation dialog

**Independent Test**: Create two folders, move a device between them → device appears only under selected folder in table (spec US1)

### Tests for User Story 1

- [x] T013 [P] [US1] Unit tests for cycle detection and path recompute in `serverBackendGo/internal/modules/device_tree/application/service_test.go`

### Implementation for User Story 1

- [x] T014 [US1] Add migration `serverBackendGo/db/migrations/0000xx_device_tree_nodes.up.sql` per `specs/017-device-control-plane/data-model.md` (`device_tree_nodes` + indexes)
- [x] T015 [US1] Implement `domain` tree node entity and validation errors in `serverBackendGo/internal/modules/device_tree/domain/tree_node.go`
- [x] T016 [US1] Implement Postgres repo CRUD + subtree device counts in `serverBackendGo/internal/modules/device_tree/adapter/persistence/postgres/tree_repo.go`
- [x] T017 [US1] Implement application service (create, rename, move, delete-with-relocation, path/depth recompute) in `serverBackendGo/internal/modules/device_tree/application/service.go`
- [x] T018 [US1] Implement HTTP handlers per `specs/017-device-control-plane/contracts/device-tree-api.md` in `serverBackendGo/internal/modules/device_tree/adapter/http/handler.go`
- [x] T019 [US1] Add migration `serverBackendGo/db/migrations/0000xx_devices_tree_node_id.up.sql` for `devices.tree_node_id` FK
- [x] T020 [US1] Seed default root folder per customer on first tree access in `serverBackendGo/internal/modules/device_tree/application/seed.go`
- [x] T021 [US1] Extend devices list query with `treeNodeId` + `includeDescendants` in `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go`
- [x] T022 [US1] Add `POST /devices/:id/move-tree` in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`
- [x] T023 [P] [US1] Create `frontend/src/features/device-tree/DeviceTreeSidebar.tsx` (select, expand, create folder)
- [x] T024 [P] [US1] Create `frontend/src/features/device-tree/DeleteTreeNodeDialog.tsx` (target folder picker per FR-001a)
- [x] T025 [US1] Integrate tree sidebar + filtered table in `frontend/src/features/devices/DevicesPage.tsx` (or equivalent devices page)
- [x] T026 [US1] Add `frontend/src/services/deviceTreeService.ts` for tree API calls
- [x] T027 [US1] Document tree endpoints in `serverBackendGo/docs/parity/device-control-plane.md` § Device tree

**Checkpoint**: Admin can manage tree and move devices; §20 gate for tree phase passable

---

## Phase 4: User Story 2 — تسجيل QR يظهر في الشجرة (Priority: P1)

**Goal**: QR enrollment auto-creates device in route default tree folder; `create=1` always; IMEI default device id mode

**Independent Test**: Scan valid route QR → device in correct tree folder within 60s with clear enrollment state (spec US2)

**Note**: Until US4 completes, bind enroll flow to legacy `configurations` row as interim `enrollment_route_id` (same numeric id); replace resolver in US4.

### Tests for User Story 2

- [x] T028 [P] [US2] Unit test enroll handler sets `tree_node_id` and `enrollment_state` in `serverBackendGo/internal/modules/sync/application/enroll_test.go`

### Implementation for User Story 2

- [x] T029 [US2] Add migration `serverBackendGo/db/migrations/0000xx_devices_control_plane_identity.up.sql` (`agent_id` UUID, `enrollment_state`, `enrollment_route_id`)
- [x] T030 [US2] Backfill `agent_id` for existing devices in same migration or follow-up SQL script
- [x] T031 [US2] Assign `agent_id` on device create in `serverBackendGo/internal/modules/devices/application/service.go`
- [x] T032 [US2] Force `create=1` in QR JSON builder in `serverBackendGo/internal/modules/qrcode/application/provisioning.go` (or equivalent)
- [x] T033 [US2] Read `default_device_id_mode` (default `imei`) when building QR in qrcode module; interim: column on `configurations` or constant until `enrollment_routes` table (US4)
- [x] T034 [US2] On successful `POST /rest/public/sync/info`, auto-create device if missing, set `tree_node_id` from route default folder, set `enrollment_route_id` and `enrollment_state=enrolled` in `serverBackendGo/internal/modules/sync/application/service.go`
- [x] T035 [US2] Reject duplicate enroll when tenant policy requires in enroll path (match Java behavior)
- [x] T036 [US2] Show `enrollment_state` column/badge in `frontend/src/features/devices/DevicesTable.tsx` (or devices list component)
- [x] T037 [US2] Set `createOnDemand={true}` always in `frontend/src/features/configurations/EnrollmentQrExperience.tsx` (move to enrollment-routes feature in US4)
- [x] T038 [US2] Document enroll + identity in `serverBackendGo/docs/parity/device-control-plane.md` § Enrollment

**Checkpoint**: New device appears under default tree folder after QR enroll (with interim configuration binding if US4 not done)

---

## Phase 5: User Story 3 — Profile للقيود والإعدادات (Priority: P1)

**Goal**: Profile editor with same tabs as configurations (restrictions, MDM, apps, design, files); policy not edited on enrollment route

**Independent Test**: Create Profile with one restriction + one app → save draft → values persist on reopen (spec US3)

### Tests for User Story 3

- [X] T039 [P] [US3] Unit tests for draft save and version fork in `serverBackendGo/internal/modules/profiles/application/draft_test.go`

### Implementation for User Story 3

- [X] T040 [US3] Add migrations `serverBackendGo/db/migrations/0000xx_profiles.up.sql` and `0000xx_profile_versions.up.sql` per data-model
- [X] T041 [US3] Add junction migrations `profile_applications`, `profile_files`, `profile_application_settings`, `profile_application_parameters` (or version-scoped names per research R3)
- [X] T042 [US3] Implement domain entities `Profile`, `ProfileVersion` in `serverBackendGo/internal/modules/profiles/domain/`
- [X] T043 [US3] Implement Postgres repos in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/`
- [X] T044 [US3] Implement draft save/load application service in `serverBackendGo/internal/modules/profiles/application/draft.go` (port settingsjson + junctions from configurations patterns)
- [X] T045 [US3] Implement HTTP handlers per `specs/017-device-control-plane/contracts/profiles-api.md` in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T046 [P] [US3] Copy and adapt `frontend/src/features/configurations/` → `frontend/src/features/profiles/` (list, editor, tabs, normalize, types)
- [X] T047 [US3] Wire `frontend/src/features/profiles/ProfileEditorPage.tsx` with Restrictions, MDM, Apps, Design, Files tabs per `specs/017-device-control-plane/contracts/frontend-control-plane-ux.md`
- [X] T048 [US3] Carry `policyLocks` from 016 in `frontend/src/features/profiles/types.ts` and `configurationNormalize.ts` equivalent
- [X] T049 [US3] Reuse 016 readonly/lock behavior in profile save path via configurations domain helpers or shared `internal/shared/policylocks` if extracted
- [X] T050 [US3] Add `GET /profiles` list endpoint with basic metadata in profiles handler

**Checkpoint**: Admin can create/edit Profile draft with full tabs; no publish required yet for US3-only test

---

## Phase 6: User Story 4 — مسار تسجيل يربط QR والشجرة والProfile (Priority: P1)

**Goal**: Enrollment route binds published profile version + default tree folder + QR; UI label «مسار تسجيل»; no restriction tabs on route editor

**Independent Test**: Two routes, same Profile, different tree folders → two enrollments land in different folders (spec US4)

### Tests for User Story 4

- [X] T051 [P] [US4] Unit tests for route validation (published version required) in `serverBackendGo/internal/modules/enrollment_routes/application/service_test.go`

### Implementation for User Story 4

- [X] T052 [US4] Add migration `serverBackendGo/db/migrations/0000xx_enrollment_routes.up.sql` per data-model
- [X] T053 [US4] Implement domain `EnrollmentRoute` in `serverBackendGo/internal/modules/enrollment_routes/domain/route.go`
- [X] T054 [US4] Implement Postgres repo and application service in `serverBackendGo/internal/modules/enrollment_routes/`
- [X] T055 [US4] Implement HTTP handlers per `specs/017-device-control-plane/contracts/enrollment-routes-api.md`
- [X] T056 [US4] Wire QR resolution by `qrcodekey` to `enrollment_routes` in `serverBackendGo/internal/modules/qrcode/` (replace interim configurations binding from US2)
- [X] T057 [US4] Set `default_device_id_mode` default `imei` on route create (FR-002a)
- [X] T058 [US4] Implement route-follow resolver stub: device sync reads `enrollment_routes.profile_version_id` in `serverBackendGo/internal/modules/sync/application/route_resolver.go`
- [X] T059 [P] [US4] Create `frontend/src/features/enrollment-routes/EnrollmentRouteListPage.tsx` and `EnrollmentRouteEditorPage.tsx` (binding fields only)
- [X] T060 [US4] Move QR panel to `frontend/src/features/enrollment-routes/` with helper text per contract (restrictions in Profile)
- [X] T061 [US4] Add i18n keys `nav.enrollmentRoutes`, `enrollmentRoute.help.profileOnly` in frontend locale files
- [X] T062 [US4] Block route save without published `profileVersionId` and `defaultTreeNodeId` in editor validation

**Checkpoint**: Two routes + enrollments prove tree placement + shared profile version binding

---

## Phase 7: User Story 5 — نشر Profile بأمان (Priority: P1)

**Goal**: Draft/publish versions, compiled artifact, usage panel, mandatory confirm when ≥50 devices, async notify without UI freeze

**Independent Test**: Profile with 50+ devices → publish shows confirm dialog → next sync receives new policy (spec US5)

### Tests for User Story 5

- [X] T063 [P] [US5] Golden test compile output in `serverBackendGo/internal/modules/profiles/application/compile_test.go` vs legacy `sync_configuration_mapper` sample
- [X] T064 [P] [US5] Unit tests for impact counts and `requiresConfirmDialog` threshold in `serverBackendGo/internal/modules/profiles/application/impact_test.go`

### Implementation for User Story 5

- [X] T065 [US5] Add migration `serverBackendGo/db/migrations/0000xx_profile_version_artifacts.up.sql`
- [X] T066 [US5] Implement `compile.go` in `serverBackendGo/internal/modules/profiles/application/` reusing `serverBackendGo/internal/modules/sync/application/sync_configuration_mapper.go`
- [X] T067 [US5] Implement `publish.go` (draft → published, bump version, write artifact + hash) in profiles application
- [X] T068 [US5] Implement `GET /profiles/:id/impact` and `POST .../publish` with `confirmImpact` per `contracts/profiles-api.md`
- [X] T069 [US5] Add usage counts (devices, routes) to `GET /profiles/:id` for usage panel
- [X] T070 [US5] Add `ProfileUsagePanel.tsx` in `frontend/src/features/profiles/` showing device and route counts
- [X] T071 [US5] Implement publish confirm modal when `requiresConfirmDialog` (≥50 devices) in `frontend/src/features/profiles/ProfilePublishDialog.tsx`
- [X] T072 [US5] Wire sync to load `profile_version_artifacts.artifact_json` in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go` per `contracts/sync-artifact.md`
- [X] T073 [US5] Keep `configurationId` in sync response mapped to `enrollment_routes.id` in `serverBackendGo/internal/modules/sync/domain/sync.go`
- [X] T074 [US5] Add optional `profileRevision` / `profileVersionId` fields to sync response (agent may ignore)
- [X] T075 [US5] Add migration `serverBackendGo/db/migrations/0000xx_domain_events.up.sql` and outbox insert on publish in `serverBackendGo/internal/modules/profiles/application/publish.go`
- [X] T076 [US5] Implement batched push worker consuming `domain_events` in `serverBackendGo/internal/platform/push/` or `internal/modules/push/application/worker.go` (debounced, FR-014)
- [X] T077 [US5] Update `serverBackendGo/docs/parity/sync.md` for artifact-based sync path

**Checkpoint**: Publish compiles artifact; ≥50 devices requires confirm; sync serves artifact; publish does not block UI

---

## Phase 8: User Story 6 — إعداد أولي موجّه (Priority: P2)

**Goal**: Checklist and wizard: tree → Profile → publish → route → QR

**Independent Test**: New tenant completes wizard → test device in tree within 15 minutes (spec US6)

### Implementation for User Story 6

- [X] T078 [P] [US6] Create `frontend/src/features/onboarding/OnboardingChecklist.tsx` (dashboard widget for incomplete setup)
- [X] T079 [US6] Create `frontend/src/features/onboarding/OnboardingWizardPage.tsx` with steps per `contracts/frontend-control-plane-ux.md`
- [X] T080 [US6] Add route `/onboarding` in `frontend/src/App.tsx` and link from checklist
- [X] T081 [US6] Guard enrollment route create: redirect to Profile create when no published profile exists
- [X] T082 [US6] Add backend `GET /rest/private/onboarding/status` in `serverBackendGo/internal/modules/profiles/adapter/http/onboarding_handler.go` (or shared admin handler) returning checklist flags

**Checkpoint**: New admin guided through full control plane setup

---

## Phase 9: Migration, Nav Cutover & Polish (Cross-Cutting)

**Purpose**: Data backfill, replace Configurations nav, API alias, blueprint gates, quickstart validation

- [X] T083 Add migration `serverBackendGo/db/migrations/0000xx_backfill_configurations_to_profiles.up.sql` (configuration → profile v1 + route + device links) per FR-010
- [X] T084 Implement backfill verification script or SQL checks documented in `specs/017-device-control-plane/quickstart.md` § Migration smoke
- [X] T085 Replace `Configurations` with `Profiles` and `Enrollment routes` in `frontend/src/navItems.ts` per FR-010a
- [X] T086 Add router redirect `/configurations` → `/profiles` in `frontend/src/App.tsx`
- [X] T087 Implement temporary `GET/PUT /rest/private/configurations/*` alias delegating to profiles/routes in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go`
- [X] T088 [P] Assign devices without `tree_node_id` to customer root during backfill in migration or job
- [X] T089 Run §20 Go/No-Go checklist from `DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md` and record results in parity doc
- [X] T090 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` and `serverBackendGo/docs/MIGRATION.md` with device control plane phase status
- [X] T091 Execute full `specs/017-device-control-plane/quickstart.md` smoke (Sprints 1–4) and fix gaps
- [X] T092 [P] Run `go test ./internal/modules/device_tree/... ./internal/modules/profiles/... ./internal/modules/enrollment_routes/... ./internal/modules/sync/...` and document results in parity doc

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup — **blocks all user stories**
- **US1 (Phase 3)**: After Foundational — **MVP**
- **US2 (Phase 4)**: After US1 (tree columns); full QR test after **US4** (or interim configurations binding per T034 note)
- **US3 (Phase 5)**: After Foundational; independent of US2/US4 for draft editor
- **US4 (Phase 6)**: After US1 + US3 (published profile version to bind)
- **US5 (Phase 7)**: After US3 + US4 (artifact + route-follow + impact counts need routes/devices)
- **US6 (Phase 8)**: After US3, US4, US5 (wizard needs publish + routes)
- **Polish (Phase 9)**: After P1 stories (US1–US5) for migration and nav cutover

### User Story Dependencies

| Story | Depends on | Independent test when |
|-------|------------|------------------------|
| US1 | Foundational | Tree + move device works |
| US2 | US1; US4 for full route semantics | QR → device in tree folder |
| US3 | Foundational | Profile draft save/load |
| US4 | US1, US3 (published version) | Two routes, two tree folders |
| US5 | US3, US4 | Publish + sync artifact |
| US6 | US3–US5 | Wizard end-to-end |

### Recommended implementation order (plan.md)

`US1 → US3 → US4 → US2 (complete binding) → US5 → Phase 9 → US6`

### Parallel Opportunities

- **Phase 1**: T002, T003, T005 in parallel
- **Phase 2**: T007, T008, T010, T011 in parallel
- **US1**: T023, T024 in parallel after T018
- **US3**: T039, T046 in parallel
- **US4**: T051, T059 in parallel
- **US5**: T063, T064 in parallel
- **US6**: T078 parallel with T082
- **Polish**: T088, T090, T092 in parallel

---

## Parallel Example: User Story 1

```bash
# After T018 (HTTP handlers done):
# Frontend tree UI in parallel:
Task T023: DeviceTreeSidebar.tsx
Task T024: DeleteTreeNodeDialog.tsx
```

---

## Parallel Example: User Story 3

```bash
# Backend domain/repos (T042–T045) then parallel:
Task T046: Copy configurations → profiles frontend
Task T039: draft_test.go
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1–2
2. Complete Phase 3 (US1)
3. **STOP and VALIDATE**: Tree + device move + folder filter (quickstart Sprint 1)
4. Demo before profiles/routes

### Incremental Delivery (recommended)

1. Setup + Foundational
2. US1 (tree) → gate §20
3. US3 (profile draft editor)
4. US4 (enrollment routes)
5. US2 + US5 (enroll + publish + artifact sync)
6. Phase 9 migration + nav cutover
7. US6 onboarding (P2)

### Parallel Team Strategy

- **Dev A**: US1 device_tree + devices filter
- **Dev B**: US3 profiles backend (after Foundational)
- **Dev C**: US4 enrollment_routes (after US3 publish API exists)
- Merge before US5 sync artifact integration

---

## Notes

- v1: **no Android app changes** — test with existing Headwind agent only
- Groups remain separate from tree (FR-011) — do not merge into device_tree module
- All private routes: tenant scope via `customerId` from JWT principal
- Commit after each task group; run blueprint §20 gate before merging each sprint

---

## Task Summary

| Phase | Story | Task IDs | Count |
|-------|-------|----------|-------|
| 1 Setup | — | T001–T005 | 5 |
| 2 Foundational | — | T006–T012 | 7 |
| 3 US1 | US1 | T013–T027 | 15 |
| 4 US2 | US2 | T028–T038 | 11 |
| 5 US3 | US3 | T039–T050 | 12 |
| 6 US4 | US4 | T051–T062 | 12 |
| 7 US5 | US5 | T063–T077 | 15 |
| 8 US6 | US6 | T078–T082 | 5 |
| 9 Polish | — | T083–T092 | 10 |
| **Total** | | **T001–T092** | **92** |

**Format validation**: All tasks use `- [ ]`, sequential `T###` IDs, story labels on user-story phases only, and explicit file paths.
