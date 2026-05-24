# Tasks: Profile Rollout & Operations

**Input**: Design documents from `/specs/018-profile-rollout-ops/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Depends on**: [017-device-control-plane](../017-device-control-plane/tasks.md) complete (tree, profiles, publish, sync artifacts, enrollment routes)

**Tests**: Unit tests per constitution IV for resolver, assignment impact, rollout recompute (not full TDD).

**Organization**: Tasks grouped by user story (US1–US5). Sync integration phase bridges US1 and US3.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US5 from [spec.md](./spec.md)

## Path Conventions

- Backend: extend `serverBackendGo/internal/modules/profiles/`; hooks in `sync/`, `devices/`
- Migrations: `serverBackendGo/db/migrations/000028_profile_rollout_ops.up.sql`
- Frontend: `frontend/src/features/profiles/`
- Parity: `serverBackendGo/docs/parity/profile-rollout-ops.md`
- Contracts: `specs/018-profile-rollout-ops/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify 017 baseline, parity stub, contract alignment

- [ ] T001 Confirm branch `018-profile-rollout-ops` and 017 migrations applied (`000019`–`000027`) per [quickstart.md](./quickstart.md) prerequisites
- [X] T002 [P] Review `specs/018-profile-rollout-ops/contracts/` against existing `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [ ] T003 [P] Audit Java device list install-status columns vs `serverBackendGo/internal/modules/devices/domain/device.go` `DeviceApplication.Status` for rollout mapping
- [X] T004 Create parity stub `serverBackendGo/docs/parity/profile-rollout-ops.md` from [contracts/profile-rollout-api.md](./contracts/profile-rollout-api.md)
- [ ] T005 [P] Document recommended execution order (Foundational → US2 ∥ US1 → Sync → US3 → US4 → US5) in `specs/018-profile-rollout-ops/plan.md` appendix if not present

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, domain types, ports, feature flag — MUST complete before user stories

**⚠️ CRITICAL**: No user story work until this phase is complete

- [X] T006 Add migration `serverBackendGo/db/migrations/000028_profile_rollout_ops.up.sql` and `.down.sql` per [data-model.md](./data-model.md) (`profiles.enabled`, `profile_tree_assignments`, device rollout columns + indexes)
- [X] T007 [P] Add domain types `ProfileTreeAssignment`, `DeviceRolloutState`, `EffectiveProfileResolution` in `serverBackendGo/internal/modules/profiles/domain/assignment.go` and `rollout.go`
- [X] T008 [P] Define `EffectiveProfileResolver` port in `serverBackendGo/internal/modules/profiles/port/resolver.go`
- [X] T009 Add `MODULE_PROFILE_ROLLOUT_ENABLED` to `serverBackendGo/internal/config/config.go` and `serverBackendGo/.env.example` (default true when profiles enabled)
- [X] T010 Implement assignment + rollout list queries skeleton in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/assignment_repo.go` and `rollout_repo.go`
- [X] T011 Extend `serverBackendGo/internal/modules/profiles/module.go` to register rollout routes when flag enabled
- [X] T012 [P] Add stub methods on `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/profile_repo.go` for `ListVersions(profileId)` if missing
- [X] T013 [P] Create `frontend/src/features/profiles/profileRolloutService.ts` with typed API placeholders matching [contracts/profile-rollout-api.md](./contracts/profile-rollout-api.md)

**Checkpoint**: `go build ./...` passes; migration applies cleanly; modules wired behind flag

---

## Phase 3: User Story 1 — إسناد Profile لمجلد الشجرة (Priority: P1) 🎯 MVP

**Goal**: Assign published profile version to tree folder; impact preview; nearest-node-wins stored in DB

**Independent Test**: Assign v2 to folder with 3 devices → assignment listed; devices marked `pending` (spec US1)

### Tests for User Story 1

- [ ] T014 [P] [US1] Unit tests for subtree device count and nearest assignment path walk in `serverBackendGo/internal/modules/profiles/application/assignment_test.go`

### Implementation for User Story 1

- [X] T015 [US1] Implement `CountSubtreeDevices` and assignment CRUD in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/assignment_repo.go`
- [X] T016 [US1] Implement `AssignmentService` (validate published version, unique node, impact ≥50) in `serverBackendGo/internal/modules/profiles/application/assignment.go`
- [X] T017 [US1] On assign: set `devices.target_profile_version_id` and `profile_rollout_status=pending` for subtree in `assignment.go`
- [X] T018 [US1] Implement HTTP handlers `GET/PUT/DELETE /profiles/:id/assignments` and `GET .../assignments/impact` in `serverBackendGo/internal/modules/profiles/adapter/http/assignment_handler.go`
- [X] T019 [US1] Register assignment routes on profiles router in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T020 [P] [US1] Create `frontend/src/features/profiles/ProfileTreeAssignmentPanel.tsx` (tree picker, version picker, impact confirm)
- [X] T021 [US1] Integrate assignment panel tab in `frontend/src/features/profiles/ProfileEditorPage.tsx`
- [X] T022 [US1] Wire `profileRolloutService.ts` assignment API calls
- [X] T023 [US1] Reuse ≥50 confirm dialog pattern from `frontend/src/features/profiles/ProfilePublishDialog.tsx` for assignment confirm
- [X] T024 [US1] Document assignment endpoints in `serverBackendGo/docs/parity/profile-rollout-ops.md` § Tree assignments

**Checkpoint**: Admin can assign published version to folder and see affected device count

---

## Phase 4: User Story 2 — التنقل بين إصدارات Profile (Priority: P1)

**Goal**: Version dropdown, read-only published view, fork draft from version, unsaved guard

**Independent Test**: Switch between published v1 and draft; fork from v1 creates new draft (spec US2)

**Note**: Can start in parallel with US1 after Phase 2 (no assignment dependency).

### Tests for User Story 2

- [ ] T025 [P] [US2] Unit test `ListVersions` ordering and `ForkDraftFromVersion` rejects draft source in `serverBackendGo/internal/modules/profiles/application/draft_test.go`

### Implementation for User Story 2

- [X] T026 [US2] Implement `ListVersions` in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/profile_repo.go`
- [X] T027 [US2] Expose `GET /profiles/:id/versions` and `POST /profiles/:id/versions/:versionId/fork-draft` in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T028 [P] [US2] Create `frontend/src/features/profiles/ProfileVersionSelect.tsx` per [contracts/profile-editor-versions-ux.md](./contracts/profile-editor-versions-ux.md)
- [X] T029 [US2] Add unsaved-draft guard modal when switching versions in `frontend/src/features/profiles/ProfileEditorPage.tsx`
- [X] T030 [US2] Read-only mode for published versions (disable inputs, hide Save/Publish) in `ProfileEditorPage.tsx`
- [X] T031 [US2] Add route `/profiles/:profileId/versions/:versionId/edit` in `frontend/src/app/App.tsx`
- [ ] T032 [US2] Extend `frontend/src/features/profiles/profileService.ts` with `listProfileVersions` and `forkDraftFromVersion`
- [X] T033 [US2] Document version endpoints in `serverBackendGo/docs/parity/profile-rollout-ops.md` § Version navigation

**Checkpoint**: Version switcher and fork-draft work without assignment panel

---

## Phase 5: Sync & Effective Profile (Integration — blocks US3)

**Purpose**: Tree-over-route resolver; sync loads correct artifact; push on assignment

**Independent Test**: Device in assigned folder receives target version on sync; route-only device unchanged when no assignment

- [ ] T034 [P] Unit tests for resolver precedence (child over parent, tree over route) in `serverBackendGo/internal/modules/profiles/application/resolver_test.go`
- [X] T035 Implement `ResolveEffectiveProfile` in `serverBackendGo/internal/modules/profiles/application/resolver.go` per [research.md](./research.md) R1
- [X] T036 Integrate resolver into sync artifact load in `serverBackendGo/internal/modules/sync/application/service.go` (or `device_sync_repo.go` caller) per [contracts/sync-rollout.md](./contracts/sync-rollout.md)
- [X] T037 Skip artifact push when `profiles.enabled=false` in resolver + sync path
- [X] T038 On assignment change: insert `domain_events` row `ProfileAssignmentChanged` and notify via existing worker in `serverBackendGo/internal/modules/profiles/application/assignment.go`
- [X] T039 Document resolver order in `serverBackendGo/docs/parity/profile-rollout-ops.md` § Effective profile

**Checkpoint**: Sync uses tree assignment when present; disabled profile not pushed

---

## Phase 6: User Story 3 — مراقبة حالة التطبيق (Priority: P1)

**Goal**: Per-device rollout status grid; recompute from sync + device info; auto-refresh UI

**Independent Test**: Two devices → one `installed`, one `failed` with reason after sync (spec US3)

### Tests for User Story 3

- [ ] T040 [P] [US3] Unit tests for status transitions (pending→installed, partial, failed) in `serverBackendGo/internal/modules/profiles/application/rollout_status_test.go`

### Implementation for User Story 3

- [X] T041 [US3] Implement `RolloutStatusService.Recompute` in `serverBackendGo/internal/modules/profiles/application/rollout_status.go` using target/applied version + `DeviceApplication.Status`
- [X] T042 [US3] Map agent status strings to admin reasons in `rollout_status.go` (document constants in parity doc)
- [ ] T043 [US3] Implement paginated `ListRolloutDevices` in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/rollout_repo.go`
- [X] T044 [US3] Implement `GET /profiles/:id/rollout/devices` and `POST .../rollout/recompute` in `serverBackendGo/internal/modules/profiles/adapter/http/rollout_handler.go`
- [X] T045 [US3] After sync config delivery: set `applied_profile_version_id` and call `Recompute` in `serverBackendGo/internal/modules/sync/application/service.go`
- [ ] T046 [US3] On device info ingest: hook `Recompute` in `serverBackendGo/internal/modules/devices/application/service.go` (or info update path)
- [X] T047 [P] [US3] Create `frontend/src/features/profiles/ProfileRolloutStatusPanel.tsx` (filters, status chips, 60s poll + manual refresh)
- [X] T048 [US3] Add Rollout tab to `frontend/src/features/profiles/ProfileEditorPage.tsx`
- [X] T049 [US3] Wire rollout APIs in `frontend/src/features/profiles/profileRolloutService.ts`
- [X] T050 [US3] Document rollout status in `serverBackendGo/docs/parity/profile-rollout-ops.md` § Rollout status

**Checkpoint**: Rollout grid updates within 2 minutes of sync; partial/failed show reasons

---

## Phase 7: User Story 4 — تعطيل وتفعيل Profile (Priority: P2)

**Goal**: Profile-level enabled flag; block assignment; warn on enrollment route; re-pending on enable

**Independent Test**: Disable → sync does not advance target; enable → devices pending again (spec US4)

### Implementation for User Story 4

- [X] T051 [US4] Implement `EnableService` (disable/enable, mark devices pending) in `serverBackendGo/internal/modules/profiles/application/enable.go`
- [X] T052 [US4] Add `POST /profiles/:id/disable` and `POST /profiles/:id/enable` in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T053 [US4] Reject new assignments when disabled in `assignment.go` with `error.profile.disabled`
- [X] T054 [P] [US4] Create `frontend/src/features/profiles/ProfileDisableBanner.tsx` and disabled badge on `frontend/src/features/profiles/ProfilesPage.tsx`
- [X] T055 [US4] Show warning when saving enrollment route bound to disabled profile in `frontend/src/features/enrollment-routes/EnrollmentRouteEditorPage.tsx`
- [X] T056 [US4] Document enable/disable in `serverBackendGo/docs/parity/profile-rollout-ops.md` § Profile lifecycle

**Checkpoint**: Disable/enable flows match spec without deleting assignments

---

## Phase 8: User Story 5 — أفضل الممارسات (Priority: P3)

**Goal**: Inline hints for draft→publish→assign→monitor workflow

**Independent Test**: Editor shows hint; assignment ≥50 shows confirm (spec US5)

### Implementation for User Story 5

- [X] T057 [P] [US5] Add best-practice hint banner in `frontend/src/features/profiles/ProfileEditorPage.tsx` per [contracts/profile-editor-versions-ux.md](./contracts/profile-editor-versions-ux.md)
- [ ] T058 [US5] Add i18n keys in `frontend/src/i18n/locales/en.json` and `frontend/src/i18n/locales/ar.json` (`profile.versions.*`, `profile.rollout.*`, `profile.assignment.*`)
- [X] T059 [US5] Ensure assignment impact dialog copy matches publish impact tone in `ProfileTreeAssignmentPanel.tsx`

**Checkpoint**: New admin sees guided copy in editor and assignment flow

---

## Phase 9: Polish & Cross-Cutting

**Purpose**: Docs, env, validation, regression tests

- [X] T060 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` with 018 profile rollout ops row
- [X] T061 [P] Update `serverBackendGo/docs/MIGRATION.md` with migration `000028` summary
- [X] T062 Run `go test ./internal/modules/profiles/application/...` and fix failures; record command in parity doc
- [ ] T063 Execute [quickstart.md](./quickstart.md) Sprints 1–4 manually and note gaps in parity doc
- [ ] T064 [P] Verify overlap edge case (parent/child assignments) with SQL or integration note in `specs/018-profile-rollout-ops/quickstart.md`
- [ ] T065 Final review: all FR-001–FR-014 traceable to tasks above in `specs/018-profile-rollout-ops/spec.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup + **017 complete** — blocks all user stories
- **US1 (Phase 3)**: After Foundational
- **US2 (Phase 4)**: After Foundational — **parallel with US1** (different files)
- **Sync integration (Phase 5)**: After US1 assignment service (T016–T017); **blocks US3**
- **US3 (Phase 6)**: After Phase 5
- **US4 (Phase 7)**: After Foundational; integrates with US1/US3 (enable checks)
- **US5 (Phase 8)**: After US1 + US2 UI exists
- **Polish (Phase 9)**: After P1 stories (US1–US3) at minimum

### User Story Dependencies

| Story | Depends on | Independent test when |
|-------|------------|------------------------|
| US1 | Foundational | Assignment CRUD + pending devices |
| US2 | Foundational | Version list + fork + read-only |
| US3 | US1 + Sync phase | Rollout grid statuses |
| US4 | Foundational (+ US1 for assignment block) | Disable blocks push |
| US5 | US1, US2 UI | Hints + confirm dialogs |

### Recommended implementation order (plan.md)

`Foundational → US2 ∥ US1 → Sync (Phase 5) → US3 → US4 → US5 → Polish`

### Parallel Opportunities

- **Phase 1**: T002, T003, T005 in parallel
- **Phase 2**: T007, T008, T012, T013 in parallel
- **US1 + US2**: Entire phases 3 and 4 in parallel after Phase 2
- **US3**: T047 parallel with T044–T046 after T041
- **Polish**: T060, T061, T064 in parallel

---

## Parallel Example: User Story 1 + 2

```bash
# After Phase 2 complete, two developers:
# Dev A — US1:
Task T015: assignment_repo.go
Task T016: assignment.go
Task T020: ProfileTreeAssignmentPanel.tsx

# Dev B — US2:
Task T026: ListVersions in profile_repo.go
Task T028: ProfileVersionSelect.tsx
Task T029: unsaved guard in ProfileEditorPage.tsx
```

---

## Parallel Example: User Story 3

```bash
# After Phase 5 (sync resolver):
Task T041: rollout_status.go
Task T043: rollout_repo.go
Task T047: ProfileRolloutStatusPanel.tsx  # parallel once API contract stable
```

---

## Implementation Strategy

### MVP First (User Story 1 + Sync)

1. Complete Phase 1–2
2. Complete US1 (assignment)
3. Complete Phase 5 (resolver + sync)
4. **STOP and VALIDATE**: Assign to folder → device syncs target version ([quickstart.md](./quickstart.md) Sprint 2)

### Incremental Delivery

1. Setup + Foundational
2. US2 (version UX) — quick win for editors
3. US1 + Sync — operational assignment
4. US3 — rollout visibility
5. US4 — disable/enable
6. US5 + Polish

### Parallel Team Strategy

- **Dev A**: US1 assignment backend + panel
- **Dev B**: US2 version navigation
- **Merge** before Phase 5 sync integration
- **Dev C**: US3 rollout after sync phase

---

## Notes

- v1: **no Android agent changes** — infer status from existing device info + sync
- Do not auto-bump tree assignments on publish; admin re-assigns to new version id
- `installing` state optional in v1; may map to `pending` (research R3)
- Commit after each task group; run quickstart sprint before merging

---

## Task Summary

| Phase | Story | Task IDs | Count |
|-------|-------|----------|-------|
| 1 Setup | — | T001–T005 | 5 |
| 2 Foundational | — | T006–T013 | 8 |
| 3 US1 | US1 | T014–T024 | 11 |
| 4 US2 | US2 | T025–T033 | 9 |
| 5 Sync | — | T034–T039 | 6 |
| 6 US3 | US3 | T040–T050 | 11 |
| 7 US4 | US4 | T051–T056 | 6 |
| 8 US5 | US5 | T057–T059 | 3 |
| 9 Polish | — | T060–T065 | 6 |
| **Total** | | **T001–T065** | **65** |

**Format validation**: All tasks use `- [ ]`, sequential `T###` IDs, story labels on user-story phases only, and explicit file paths.
