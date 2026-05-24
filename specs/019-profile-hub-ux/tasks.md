# Tasks: Profile Hub & Enrollment UX

**Input**: Design documents from `/specs/019-profile-hub-ux/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Depends on**: [017-device-control-plane](../017-device-control-plane/tasks.md) complete; [018-profile-rollout-ops](../018-profile-rollout-ops/tasks.md) complete (assignments, rollout, publish, `000028`)

**Tests**: Unit tests for `ProfileHealthService` per constitution IV (not full TDD).

**Organization**: Tasks grouped by user story (US1–US9). Backend foundation blocks workspace UI.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US9 from [spec.md](./spec.md)

## Path Conventions

- Backend: `serverBackendGo/internal/modules/profiles/`, `enrollment_routes/`
- Migrations: `serverBackendGo/db/migrations/000029_profile_hub_activity_idx.*.sql`
- Frontend: `frontend/src/features/profiles/workspace/`, `frontend/src/features/enrollment-routes/`
- Parity: `serverBackendGo/docs/parity/profile-hub-ux.md`
- Contracts: `specs/019-profile-hub-ux/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify 017/018 baseline and align contracts before implementation

- [X] T001 Confirm branch `019-profile-hub-ux` and migrations through `000028` applied per [quickstart.md](./quickstart.md) prerequisites
- [X] T002 [P] Review `specs/019-profile-hub-ux/contracts/` against existing `profiles` and `enrollment_routes` HTTP handlers
- [X] T003 [P] Create parity stub `serverBackendGo/docs/parity/profile-hub-ux.md` from [contracts/profile-hub-api.md](./contracts/profile-hub-api.md) and [contracts/enrollment-routes-policy-decoupling.md](./contracts/enrollment-routes-policy-decoupling.md)
- [X] T004 [P] Add `PROFILE_STALE_PUBLISH_DAYS` (default 30) to `serverBackendGo/internal/config/config.go` and `serverBackendGo/.env.example`
- [X] T005 Document execution order (Foundational → US2 shell → US3/US4 → US5–US7 → US8/US9) in `specs/019-profile-hub-ux/plan.md` appendix if missing

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Hub APIs (health, summary, activity), optional migration, shared frontend services — MUST complete before user story UI work

**⚠️ CRITICAL**: No workspace or list-radar work until this phase is complete

- [X] T006 Add migration `serverBackendGo/db/migrations/000029_profile_hub_activity_idx.up.sql` and `.down.sql` per [data-model.md](./data-model.md)
- [X] T007 [P] Add `ProfileHealth`, `ProfileSummary`, `ProfileActivityEvent` types in `serverBackendGo/internal/modules/profiles/domain/hub.go`
- [X] T008 [P] Implement `ProfileHealthService` rules in `serverBackendGo/internal/modules/profiles/application/health.go` per [research.md](./research.md) R4
- [X] T009 [P] Unit tests for health transitions in `serverBackendGo/internal/modules/profiles/application/health_test.go`
- [X] T010 Extend `List` query in `serverBackendGo/internal/modules/profiles/adapter/persistence/postgres/profile_repo.go` with `health`, `badges`, `assignmentCount`, `rolloutFailureCount`
- [X] T011 Implement `SummaryService` aggregator in `serverBackendGo/internal/modules/profiles/application/summary.go`
- [X] T012 Implement `ActivityService` reading `domain_events` in `serverBackendGo/internal/modules/profiles/application/activity.go`
- [X] T013 Emit `ProfileEnabled` / `ProfileDisabled` domain events from `serverBackendGo/internal/modules/profiles/application/enable.go`
- [X] T014 Implement `GET /profiles/:id/summary` and `GET /profiles/:id/activity` in `serverBackendGo/internal/modules/profiles/adapter/http/hub_handlers.go`
- [X] T015 Register hub routes and extend list response in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go` and `module.go`
- [X] T016 [P] Create `frontend/src/features/profiles/profileHubService.ts` (`getProfileSummary`, `getProfileActivity`, extended list types)
- [X] T017 [P] Create `frontend/src/features/profiles/workspace/profileWorkspaceState.ts` (section, dirty, secondaryPanel types)

**Checkpoint**: `go test ./internal/modules/profiles/application/...` passes; summary/activity endpoints return data in Swagger/curl

---

## Phase 3: User Story 1 — إزالة ربط Profile من مسار التسجيل (Priority: P1) 🎯 MVP (policy decoupling)

**Goal**: Enrollment routes are enrollment-only; no profile version picker in UI or required API field

**Independent Test**: Create enrollment route without `profileVersionId` → saves → UI shows policy helper text only ([quickstart.md](./quickstart.md) Sprint 1)

### Implementation for User Story 1

- [X] T018 [US1] Make `profileVersionId` optional in `serverBackendGo/internal/modules/enrollment_routes/domain/route.go` create/update DTOs
- [X] T019 [US1] Update `validateBinding` in `serverBackendGo/internal/modules/enrollment_routes/application/service.go` to not require published profile version
- [X] T020 [US1] Update `serverBackendGo/internal/modules/enrollment_routes/application/service_test.go` for routes without profile
- [X] T021 [P] [US1] Remove profile version `<Select>` and `listPublishedProfileVersions` usage from `frontend/src/features/enrollment-routes/EnrollmentRouteEditorPage.tsx`
- [X] T022 [P] [US1] Add policy helper callout linking to Profiles → Assignments in `EnrollmentRouteEditorPage.tsx`
- [X] T023 [US1] Document decoupling in `serverBackendGo/docs/parity/profile-hub-ux.md` § Enrollment routes and mark `published-profile-versions` endpoint deprecated in UI

**Checkpoint**: New routes enroll devices; policy comes only from tree assignment (018)

---

## Phase 4: User Story 2 — Profile Workspace shell (Priority: P1)

**Goal**: Click profile in list opens near-fullscreen Workspace with fixed header and sidebar (not classic popup)

**Independent Test**: Open profile → sidebar navigate → Close → list state preserved ([quickstart.md](./quickstart.md) Sprint 2)

### Implementation for User Story 2

- [X] T024 [P] [US2] Create `frontend/src/features/profiles/workspace/ProfileWorkspace.tsx` (Dialog desktop / Sheet mobile per [research.md](./research.md) R1)
- [X] T025 [P] [US2] Create `frontend/src/features/profiles/workspace/ProfileWorkspaceSidebar.tsx` with sections per [contracts/profile-workspace-ux.md](./contracts/profile-workspace-ux.md)
- [X] T026 [US2] Create `frontend/src/features/profiles/workspace/ProfileCockpitHeader.tsx` (name, badges placeholders, Close)
- [X] T027 [US2] Add `ProfileWorkspaceProvider` and URL sync `?open=&section=` in `profileWorkspaceState.ts`
- [X] T028 [US2] Wire `ProfilesPage.tsx` row click to open workspace instead of `navigate(/profiles/:id/edit)`
- [X] T029 [US2] Add legacy redirect `/profiles/:profileId/edit` → `/profiles?open={id}&section=editor` in `frontend/src/app/App.tsx`

**Checkpoint**: Workspace opens/closes; sidebar switches sections; zero nested dialogs

---

## Phase 5: User Story 4 — Health state & Control Radar list (Priority: P1)

**Goal**: List shows health chip and operational badges without opening profile

**Independent Test**: Profile with no assignment shows `No Assignment` + Warning health in list ([spec.md](./spec.md) US4)

**Note**: Backend in Phase 2; this phase is list UI (can start after T015–T016).

### Tests for User Story 4

- [X] T030 [P] [US4] Extend list handler response test or handler smoke in `serverBackendGo/internal/modules/profiles/adapter/http/handler_test.go` if file exists, else document curl check in parity

### Implementation for User Story 4

- [X] T031 [P] [US4] Create `frontend/src/features/profiles/ProfileHealthBadge.tsx` and `ProfileListBadges.tsx`
- [X] T032 [US4] Update `frontend/src/features/profiles/ProfilesPage.tsx` to render health + badges from extended list API
- [X] T033 [US4] Wire health/lifecycle badges into `ProfileCockpitHeader.tsx` from summary on workspace open

**Checkpoint**: List acts as control radar; header matches list health after open

---

## Phase 6: User Story 3 — Overview & read/edit separation (Priority: P1)

**Goal**: Default Overview cards (no inputs); Edit moves to Editor section with distinct visuals

**Independent Test**: Overview has no inputs; Edit shows warning bar + sticky save ([quickstart.md](./quickstart.md) Sprint 3 partial)

### Implementation for User Story 3

- [ ] T034 [P] [US3] Create `frontend/src/features/profiles/workspace/ProfileOverviewCards.tsx` loading `getProfileSummary` per [contracts/profile-workspace-ux.md](./contracts/profile-workspace-ux.md)
- [ ] T035 [US3] Mount Overview as default section in `ProfileWorkspace.tsx`
- [ ] T036 [P] [US3] Create `frontend/src/features/profiles/workspace/sections/EditorSection.tsx` wrapping existing editor tabs with read/edit chrome
- [ ] T037 [US3] Implement header **Edit** → `section=editor` and read-mode styling (`bg-muted/30` vs editor accent) in `ProfileWorkspace.tsx`
- [ ] T038 [US3] Add unsaved-changes guard when leaving Editor or closing workspace in `profileWorkspaceState.ts`

**Checkpoint**: 5-second rule testable on Overview cards (SC-006)

---

## Phase 7: User Story 7 — Publish from header + Impact Preview (Priority: P1)

**Goal**: Publish in cockpit header opens side sheet impact panel — not nested dialog

**Independent Test**: Publish from header → right sheet with device/folder counts → confirm ([quickstart.md](./quickstart.md) Sprint 3)

### Implementation for User Story 7

- [ ] T039 [P] [US7] Create `frontend/src/features/profiles/workspace/ProfilePublishImpactSheet.tsx` (Sheet side=right, reuse `profileService.getProfileImpact` / publish)
- [ ] T040 [US7] Wire header **Publish** button with `canPublish` from summary and open `secondaryPanel=publish-impact` in `ProfileCockpitHeader.tsx`
- [ ] T041 [US7] Disable Publish when `!hasUnpublishedDraft` or validation errors; tooltip reason

**Checkpoint**: Publish never requires navigating into Editor tabs only

---

## Phase 8: User Story 5 — Assignments + miniature tree (Priority: P1)

**Goal**: Assignments section shows tree preview with selected folder highlight

**Independent Test**: Select folder in assignment → preview path highlights before save ([spec.md](./spec.md) US5)

### Implementation for User Story 5

- [ ] T042 [P] [US5] Create `frontend/src/features/profiles/workspace/TreePreview.tsx` (compact tree from `getDeviceTree`, max depth 4)
- [ ] T043 [US5] Create `frontend/src/features/profiles/workspace/sections/AssignmentsSection.tsx` wrapping `ProfileTreeAssignmentPanel.tsx` + `TreePreview`
- [ ] T044 [US5] Use inline confirm or right Sheet for ≥50 device assignment impact (no nested Dialog) in `AssignmentsSection.tsx`

**Checkpoint**: Assignment UX shows visual folder path

---

## Phase 9: User Story 6 — Create profile → assign wizard (Priority: P1)

**Goal**: After create, workspace opens on Assignments with publish-first guidance if needed

**Independent Test**: New profile → wizard on Assignments → assign → Overview shows assignment card ([quickstart.md](./quickstart.md) Sprint 4)

### Implementation for User Story 6

- [ ] T045 [US6] Update `ProfilesPage.tsx` / `ProfileForm` onCreate success to open workspace `?open={id}&wizard=assign`
- [ ] T046 [US6] Add wizard banner in `AssignmentsSection.tsx` (publish-first vs assign step) using summary flags
- [ ] T047 [US6] On skip assign, show persistent Warning on Overview card with CTA «Assign now» in `ProfileOverviewCards.tsx`

**Checkpoint**: SC-001 create+publish+assign under 3 minutes in manual test

---

## Phase 10: Deep sections — Rollout, Versions, Editor integration (Priority: P1, supports US2–US3)

**Goal**: Mount 018 panels inside workspace sections (not standalone page)

**Independent Test**: Sidebar → Rollout shows grid; Versions lists; Editor saves draft

### Implementation

- [ ] T048 [P] Create `frontend/src/features/profiles/workspace/sections/RolloutSection.tsx` wrapping `ProfileRolloutStatusPanel.tsx`
- [ ] T049 [P] Create `frontend/src/features/profiles/workspace/sections/VersionsSection.tsx` wrapping `ProfileVersionSelect.tsx` + version list
- [ ] T050 Register all sections in `ProfileWorkspace.tsx` switch by `section` state
- [ ] T051 Deprecate default navigation to `ProfileEditorPage.tsx`; keep file as thin redirect or embed only via `EditorSection`

**Checkpoint**: Full workspace navigation without `/profiles/:id/edit` as primary path

---

## Phase 11: User Story 8 — Activity timeline (Priority: P2)

**Goal**: Activity section shows publish/assign/enable events from API

**Independent Test**: After publish + assign, Activity lists both events ([quickstart.md](./quickstart.md) Sprint 5)

### Implementation for User Story 8

- [ ] T052 [P] [US8] Create `frontend/src/features/profiles/workspace/sections/ActivitySection.tsx` consuming `getProfileActivity`
- [ ] T053 [US8] Add i18n keys for `profile.activity.*` summary templates in `frontend/src/i18n/locales/en.json` and `ar.json`
- [ ] T054 [US8] Document activity event types in `serverBackendGo/docs/parity/profile-hub-ux.md` § Activity

**Checkpoint**: Activity section non-empty after typical admin workflow

---

## Phase 12: User Story 9 — Mobile Workspace (Priority: P2)

**Goal**: Narrow viewport uses full-screen sheet, drawer nav, bottom actions

**Independent Test**: Viewport &lt;768px → full screen + bottom Edit/Publish/Close ([quickstart.md](./quickstart.md) Sprint 6)

### Implementation for User Story 9

- [ ] T055 [US9] Add responsive breakpoint switch Dialog→Sheet in `ProfileWorkspace.tsx`
- [ ] T056 [US9] Add mobile drawer toggle for sidebar in `ProfileWorkspaceSidebar.tsx`
- [ ] T057 [US9] Add bottom action bar on mobile in `ProfileCockpitHeader.tsx` or `ProfileWorkspace.tsx`

**Checkpoint**: Mobile smoke passes without horizontal tab overflow

---

## Phase 13: Polish & Cross-Cutting

**Purpose**: i18n, docs, regression, accessibility

- [ ] T058 [P] Add `profile.workspace.*`, `profile.health.*`, `profile.badges.*` keys in `frontend/src/i18n/locales/en.json` and `ar.json`
- [ ] T059 [P] Add lifecycle color tokens (draft amber, published green, disabled muted) in workspace components per [research.md](./research.md) R10
- [ ] T060 [P] Optional P2 keyboard shortcuts (E, Esc, Ctrl+S) in `ProfileWorkspace.tsx`
- [ ] T061 Update `serverBackendGo/docs/NEXT_STEPS.md` with 019 profile hub UX row
- [ ] T062 Update `serverBackendGo/docs/MIGRATION.md` with `000029` summary if migration added
- [ ] T063 Run `go test ./internal/modules/profiles/application/...` and `go build ./...`; record in parity doc
- [ ] T064 Execute [quickstart.md](./quickstart.md) Sprints 1–6 and note gaps in parity doc
- [ ] T065 UX audit: confirm zero nested dialogs in workspace flows (SC-007)
- [ ] T066 Final FR traceability review against [spec.md](./spec.md) FR-001–FR-020

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup + **018 complete** — blocks all UI stories
- **US1 (Phase 3)**: After Foundational — can parallel with US2 backend done
- **US2 (Phase 4)**: After T016–T017 (Foundational frontend stubs)
- **US4 (Phase 5)**: After T015–T016; parallel with US2 once shell exists
- **US3, US7 (Phases 6–7)**: After US2 shell + summary API
- **US5, US6 (Phases 8–9)**: After US2; US5 needs 018 assignment APIs
- **Deep sections (Phase 10)**: After US2; parallel with US5–US7
- **US8, US9 (Phases 11–12)**: After workspace sections mounted
- **Polish (Phase 13)**: After P1 stories

### User Story Dependencies

| Story | Depends on | Independent test when |
|-------|------------|------------------------|
| US1 | Foundational | Route saves without profile field |
| US2 | Foundational | Workspace open/close + sidebar |
| US3 | US2 | Overview read-only + Editor section |
| US4 | Foundational (API) | List badges visible |
| US5 | US2 + 018 | Tree preview on assign |
| US6 | US2 + US5 section | Create → assign wizard |
| US7 | US2 + summary | Header publish + impact sheet |
| US8 | US2 + T012–T014 | Activity timeline |
| US9 | US2 | Mobile layout |

### Recommended order (from plan.md)

```text
Setup → Foundational → US1 ∥ US2 → US4 → US3 + US7 → US5 + US6 → Phase 10 → US8 + US9 → Polish
```

### Parallel Opportunities

- **Phase 1**: T002, T003, T004 in parallel
- **Phase 2**: T007–T009, T016–T017 in parallel after T006
- **US1**: T021, T022 in parallel after T019
- **US2**: T024, T025, T026 in parallel
- **Phase 10**: T048, T049 in parallel
- **Polish**: T058, T059, T060 in parallel

---

## Parallel Example: User Story 2

```bash
# After Phase 2 complete:
Task T024: ProfileWorkspace.tsx
Task T025: ProfileWorkspaceSidebar.tsx
Task T026: ProfileCockpitHeader.tsx
# Then sequentially: T027 URL state → T028 ProfilesPage → T029 App redirect
```

---

## Parallel Example: Foundational backend

```bash
Task T008: health.go
Task T009: health_test.go
Task T012: activity.go
# Then T011 summary.go → T014 hub_handlers.go → T015 register routes
```

---

## Implementation Strategy

### MVP First (US1 + US2 + US4)

1. Complete Phase 1–2 (hub APIs + health list)
2. Complete US1 (enrollment decoupling)
3. Complete US2 (workspace shell)
4. Complete US4 (list radar)
5. **STOP and VALIDATE**: [quickstart.md](./quickstart.md) Sprints 1–2

### Incremental Delivery

1. Foundational → US1 + US2 + US4 (operational shell + policy model)
2. US3 + US7 (read/edit + publish header)
3. US5 + US6 + Phase 10 (assignments + full sections)
4. US8 + US9 + Polish

### Parallel Team Strategy

- **Dev A**: Phase 2 backend + US1 enrollment
- **Dev B**: US2 workspace shell + US4 list
- **Merge** before US3/US7
- **Dev C**: US5/US6 assignments + Phase 10 sections

---

## Task Summary

| Phase | Story | Task IDs | Count |
|-------|-------|----------|-------|
| 1 Setup | — | T001–T005 | 5 |
| 2 Foundational | — | T006–T017 | 12 |
| 3 US1 | US1 | T018–T023 | 6 |
| 4 US2 | US2 | T024–T029 | 6 |
| 5 US4 | US4 | T030–T033 | 4 |
| 6 US3 | US3 | T034–T038 | 5 |
| 7 US7 | US7 | T039–T041 | 3 |
| 8 US5 | US5 | T042–T044 | 3 |
| 9 US6 | US6 | T045–T047 | 3 |
| 10 Sections | — | T048–T051 | 4 |
| 11 US8 | US8 | T052–T054 | 3 |
| 12 US9 | US9 | T055–T057 | 3 |
| 13 Polish | — | T058–T066 | 9 |
| **Total** | | **T001–T066** | **66** |

**Format validation**: All tasks use `- [ ]`, sequential `T###` IDs, `[USn]` on user-story phases, explicit file paths.

**Suggested MVP scope**: Phase 1–2 + US1 + US2 + US4 (T001–T033) — enrollment decoupling, workspace shell, control radar list.

---

## Notes

- Reuses 018 components inside workspace sections; refactor `ProfileEditorPage` rather than duplicate policy tabs.
- **No nested dialogs** — enforce in code review (SC-007).
- `ProfileDisableBanner` logic moves into Overview/Header (enabled toggle in cockpit P2 if not in v1 scope).
- Commit after each phase checkpoint; run quickstart sprint before merge.
