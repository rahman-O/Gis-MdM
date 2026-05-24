# Tasks: Profile Workspace Maturity

**Input**: Design documents from `/specs/020-profile-workspace-maturity/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Depends on**: [019-profile-hub-ux](../019-profile-hub-ux/tasks.md) Phases 1–4 minimum (workspace shell + hub APIs); [018-profile-rollout-ops](../018-profile-rollout-ops/tasks.md) complete

**Tests**: Unit tests for `VersionDeleteService` and publish assignment bump per constitution IV (not full TDD).

**Organization**: Tasks grouped by user story (US1–US7). Backend foundation blocks all UI stories.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US7 from [spec.md](./spec.md)

## Path Conventions

- Backend: `serverBackendGo/internal/modules/profiles/`
- Frontend: `frontend/src/features/profiles/workspace/`, `frontend/src/app/App.tsx`
- Parity: `serverBackendGo/docs/parity/profile-workspace-maturity.md`
- Contracts: `specs/020-profile-workspace-maturity/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify 019/018 baseline and align contracts before implementation

- [X] T001 Confirm branch `020-profile-workspace-maturity` and 019 workspace baseline applied per [quickstart.md](./quickstart.md) prerequisites
- [X] T002 [P] Review `specs/020-profile-workspace-maturity/contracts/` against existing `profiles` HTTP handlers and 018 rollout routes
- [X] T003 [P] Create parity stub `serverBackendGo/docs/parity/profile-workspace-maturity.md` from [contracts/profile-workspace-maturity-api.md](./contracts/profile-workspace-maturity-api.md)
- [X] T004 [P] Add `versionId` and `readOnly` query param support to `frontend/src/features/profiles/workspace/profileWorkspaceState.tsx` (plan R7)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Publish assignment bump, version delete, summary `publishedContext`, shared workspace refresh — MUST complete before user story UI work

**⚠️ CRITICAL**: No US2–US7 UI work until this phase is complete

- [X] T005 [P] Add `PublishImpactAssignment`, `PublishedContext` types in `serverBackendGo/internal/modules/profiles/domain/publish.go` and `domain/hub.go`
- [X] T006 [P] Add version delete domain types and error sentinels in `serverBackendGo/internal/modules/profiles/domain/version_delete.go`
- [X] T007 Implement `VersionDeleteService` in `serverBackendGo/internal/modules/profiles/application/version_delete.go` per [data-model.md](./data-model.md) guards
- [X] T008 [P] Unit tests for version delete guards in `serverBackendGo/internal/modules/profiles/application/version_delete_test.go`
- [X] T009 Extend `PublishService.Impact` to return `assignmentsToUpdate` in `serverBackendGo/internal/modules/profiles/application/publish.go`
- [X] T010 Implement transactional `BumpAllAssignmentsOnPublish` in `serverBackendGo/internal/modules/profiles/application/publish.go` (reuse `AssignmentService` / rollout recompute from 018)
- [X] T011 [P] Unit tests for publish assignment bump in `serverBackendGo/internal/modules/profiles/application/publish_assignments_test.go`
- [X] T012 Extend `HubService.Summary` with `publishedContext` block in `serverBackendGo/internal/modules/profiles/application/hub.go`
- [X] T013 Emit `ProfileVersionDeleted` domain event from `version_delete.go` via `ProfileRepository.InsertDomainEvent`
- [X] T014 Implement `DELETE /profiles/:id/versions/:versionId` in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T015 Extend `GET /profiles/:id/impact` and publish response `assignmentsUpdated` in `serverBackendGo/internal/modules/profiles/adapter/http/handler.go`
- [X] T016 Wire new services in `serverBackendGo/internal/modules/profiles/module.go`
- [X] T017 [P] Create `frontend/src/features/profiles/workspace/profileWorkspaceEvents.ts` refresh bus per [research.md](./research.md) R5
- [X] T018 [P] Extend `frontend/src/features/profiles/profileHubService.ts` with `publishedContext`, `deleteProfileVersion`, extended `PublishProfileResult`

**Checkpoint**: `go test ./internal/modules/profiles/application/...` passes; curl DELETE version + impact assignments preview documented in parity

---

## Phase 3: User Story 1 — Workspace كمسار وحيد (Priority: P1) 🎯 MVP

**Goal**: Profile Workspace is the only admin path; legacy editor routes redirect into workspace

**Independent Test**: Open `/profiles/2/edit` → redirects to workspace; no full-page editor ([quickstart.md](./quickstart.md) Sprint 1)

### Implementation for User Story 1

- [X] T019 [P] [US1] Extract `ProfileEditorCore.tsx` from `frontend/src/features/profiles/ProfileEditorPage.tsx` for embedded use
- [X] T020 [US1] Refactor `ProfileEditorPage.tsx` to thin wrapper or remove export used by routes
- [X] T021 [US1] Remove standalone routes `/profiles/:profileId/edit` and `/profiles/:profileId/versions/:versionId/edit` from `frontend/src/app/App.tsx`
- [X] T022 [US1] Extend `frontend/src/features/profiles/ProfileEditRedirect.tsx` for `versionId` query → `section=editor&versionId=`
- [X] T023 [US1] Update `frontend/src/features/profiles/ProfilesPage.tsx` create flow to open workspace (Overview or Assignments) not `/profiles/:id/edit`
- [X] T024 [P] [US1] Audit and remove internal `navigate('/profiles/.../edit')` calls in `frontend/src/features/profiles/`

**Checkpoint**: SC-001 — all profile admin flows stay inside workspace overlay

---

## Phase 4: User Story 2 — Assignments تعرض إصدار السياسة الفعلي (Priority: P1)

**Goal**: Assignments shows published version context, compact summary, link to read-only Editor; no blank screen

**Independent Test**: Published profile → Assignments shows vN strip + View full policy ([quickstart.md](./quickstart.md) Sprint 2)

### Implementation for User Story 2

- [X] T025 [P] [US2] Create `frontend/src/features/profiles/workspace/ProfileAssignmentsContext.tsx` (published bar + pinned strip)
- [X] T026 [US2] Integrate `ProfileAssignmentsContext` above `ProfileTreeAssignmentPanel` in `frontend/src/features/profiles/workspace/ProfileWorkspaceContent.tsx`
- [X] T027 [US2] Wire «View full policy» to `setSection('editor')` with `versionId=published` and `readOnly=1` in `profileWorkspaceState.tsx`
- [X] T028 [US2] Add empty state (no published) with CTA to Versions/Editor in `ProfileAssignmentsContext.tsx`
- [X] T029 [US2] Replace silent failure with error + Retry in `ProfileWorkspaceContent.tsx` Assignments section
- [X] T030 [US2] Fix workspace Dialog flex layout (`min-h-0`, scroll region) in `frontend/src/features/profiles/workspace/ProfileWorkspace.tsx` per 019 bugfix

**Checkpoint**: SC-002, SC-003 for Assignments; no empty content area

---

## Phase 5: User Story 3 — Editor: مسودة، منشور، وإصدارات سابقة (Priority: P1)

**Goal**: Editor section loads draft/published/historical versions with unsaved guards and sticky save

**Independent Test**: Switch versions with dirty warning; save draft without leaving workspace ([quickstart.md](./quickstart.md) Sprint 7)

### Implementation for User Story 3

- [X] T031 [P] [US3] Create `frontend/src/features/profiles/workspace/ProfileEditorSection.tsx` mounting `ProfileEditorCore`
- [X] T032 [US3] Replace embedded `ProfileEditorPage` usage in `ProfileWorkspaceContent.tsx` with `ProfileEditorSection.tsx`
- [X] T033 [US3] Add version switcher (018 `listProfileVersions`) to `ProfileEditorSection.tsx`
- [X] T034 [US3] Implement read-only mode when `readOnly=1` query in `ProfileEditorSection.tsx`
- [X] T035 [US3] Enforce unsaved-changes guard on section switch and workspace close in `profileWorkspaceState.tsx` + `ProfileWorkspace.tsx`
- [X] T036 [US3] Add sticky save bar and «last saved» indicator in `ProfileEditorSection.tsx`
- [X] T037 [P] [US3] Wire «Fork draft from published» via `profileRolloutService` in `ProfileEditorSection.tsx`

**Checkpoint**: US3 acceptance scenarios 1–5 pass inside workspace

---

## Phase 6: User Story 4 — إسناد لعدة مجلدات/فروع (Priority: P1)

**Goal**: Multiple folder assignments visible; parent/child overlap explained in UI

**Independent Test**: Assign two folders; both listed; parent+child shows override hint ([quickstart.md](./quickstart.md) Sprint 3–4)

### Implementation for User Story 4

- [X] T038 [US4] Add overlap hint helper using `treePath` from assignment list in `frontend/src/features/profiles/ProfileTreeAssignmentPanel.tsx`
- [X] T039 [US4] Show informational banner when parent and child both assigned in `ProfileAssignmentsContext.tsx`
- [X] T040 [P] [US4] Refresh assignment list after PUT/DELETE via `profileWorkspaceEvents.ts` in `ProfileTreeAssignmentPanel.tsx`
- [X] T041 [US4] Document parent/child nearest-wins in `serverBackendGo/docs/parity/profile-workspace-maturity.md` § Assignments

**Checkpoint**: SC-006; clarification A visual check

---

## Phase 7: User Story 5 — حذف إصدار بأمان (Priority: P2)

**Goal**: Delete draft or unused historical published version with guards and Versions UI refresh

**Independent Test**: Delete extra draft; block delete of active published ([quickstart.md](./quickstart.md) Sprint 6)

### Implementation for User Story 5

- [X] T042 [P] [US5] Add `deleteProfileVersion` to `frontend/src/features/profiles/profileRolloutService.ts` or `profileHubService.ts`
- [X] T043 [US5] Enhance `VersionsSection` in `ProfileWorkspaceContent.tsx` with per-row Delete + confirm copy
- [X] T044 [US5] Map API error keys to user messages for delete failures in Versions UI
- [X] T045 [US5] Emit workspace refresh after successful delete in `profileWorkspaceEvents.ts`
- [X] T046 [P] [US5] Show `ProfileVersionDeleted` in Activity section when event type present

**Checkpoint**: SC-004 delete flow; FR-005 rules enforced server-side

---

## Phase 8: User Story 6 — Overview متزامن مع الإعدادات (Priority: P2)

**Goal**: Overview cards reflect published settings; draft badge only; auto-refresh after mutations

**Independent Test**: Publish → Overview updates within 5s without browser refresh ([quickstart.md](./quickstart.md) Sprint 5)

### Implementation for User Story 6

- [X] T047 [P] [US6] Extract `ProfileOverviewSection.tsx` from `ProfileWorkspaceContent.tsx` with published-only card logic
- [X] T048 [US6] Add `hasUnpublishedDraft` banner with link to Editor in `ProfileOverviewSection.tsx`
- [X] T049 [US6] When no published version, show draft-based cards with «Not published yet» label per spec
- [X] T050 [US6] Subscribe Overview to `profileWorkspaceEvents.ts` and refetch summary on bump
- [X] T051 [P] [US6] Invalidate summary in cockpit header after refresh in `ProfileWorkspace.tsx`

**Checkpoint**: SC-005; clarification A (published cards + draft badge)

---

## Phase 9: User Story 7 — تجربة مسودة / حفظ / نشر (Priority: P2)

**Goal**: Publish from header via side sheet with assignment bump preview; no nested publish dialog

**Independent Test**: Edit → Save → Publish sheet shows folders → confirm → Assignments/Overview update ([quickstart.md](./quickstart.md) Sprint 5)

### Implementation for User Story 7

- [X] T052 [P] [US7] Create `frontend/src/features/profiles/workspace/ProfilePublishImpactSheet.tsx` per [contracts/profile-workspace-sections-ux.md](./contracts/profile-workspace-sections-ux.md)
- [X] T053 [US7] Wire cockpit header Publish to open impact sheet (not `ProfilePublishDialog`) in `ProfileWorkspace.tsx`
- [X] T054 [US7] Render `assignmentsToUpdate` table in `ProfilePublishImpactSheet.tsx` from extended impact API
- [X] T055 [US7] On publish success: close sheet, toast, bump workspace refresh, optional navigate to Overview
- [X] T056 [US7] Disable duplicate Publish in embedded editor when workspace header owns publish in `ProfileEditorSection.tsx`
- [X] T057 [P] [US7] Extend `getProfileImpact` / publish client to send `confirmImpact` and read `assignmentsUpdated`

**Checkpoint**: FR-011; no nested dialogs during publish (SC-007 regression)

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Parity, quickstart validation, redirects, layout hardening

- [X] T058 [P] Add redirect component or route for bookmarked version URLs in `frontend/src/app/App.tsx` if not covered by T022
- [X] T059 Complete `serverBackendGo/docs/parity/profile-workspace-maturity.md` with curl examples from [quickstart.md](./quickstart.md)
- [X] T060 Run quickstart Sprints 1–7 and fix gaps documented in parity
- [X] T061 [P] Add i18n keys for new strings (Assignments strip, publish sheet, delete errors) in `frontend/src/locales/` if project uses them
- [X] T062 Verify `profileWorkspaceState.ts` re-export shim remains valid after `.tsx` changes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Start immediately
- **Foundational (Phase 2)**: After Setup — **blocks all user stories**
- **US1 (Phase 3)**: After Foundational — recommended first UI slice (MVP)
- **US2–US3 (Phases 4–5)**: After Foundational; US2 ∥ US1 partially; US3 benefits from US1 redirects
- **US4 (Phase 6)**: After US2 (assignment panel)
- **US5 (Phase 7)**: After Foundational (DELETE API); UI after US3 Versions section exists
- **US6 (Phase 8)**: After Foundational + workspace events (T017); best after US7 publish refresh
- **US7 (Phase 9)**: After Foundational publish bump (T009–T011)
- **Polish (Phase 10)**: After desired stories complete

### User Story Dependencies

| Story | Depends on |
|-------|------------|
| US1 | Phase 2 |
| US2 | Phase 2, workspace shell (019) |
| US3 | US1 (embedded editor routes) |
| US4 | US2 |
| US5 | Phase 2 DELETE API |
| US6 | T017 refresh bus; summary API |
| US7 | Phase 2 publish bump |

### Parallel Opportunities

- Phase 1: T002–T004 parallel
- Phase 2: T005–T006, T008, T011, T017–T018 parallel after T007/T009 started
- US1: T019 ∥ T024
- US2: T025 ∥ T030
- US5: T042 ∥ T046
- US6: T047 ∥ T051
- US7: T052 ∥ T057

### Parallel Example: Foundational backend

```bash
# Parallel domain + tests:
T005 domain publish types
T006 domain version_delete types
T008 version_delete_test.go
T011 publish_assignments_test.go

# Then sequential integration:
T009 → T010 → T014 → T015 → T016
```

### Parallel Example: P1 UI after foundation

```bash
Developer A: US1 (T019–T024)
Developer B: US2 (T025–T030)
Developer C: US3 (T031–T037) after T019 extracts ProfileEditorCore
```

---

## Implementation Strategy

### MVP First (Recommended)

1. Phase 1 Setup  
2. Phase 2 Foundational (**required**)  
3. Phase 3 US1 — workspace-only routes  
4. Phase 4 US2 — fix Assignments blank screen  
5. **STOP**: Validate quickstart Sprints 1–2  

### Incremental delivery

1. Foundation → US1 → US2 → US3 (P1 complete)  
2. US7 publish sheet → US6 overview sync  
3. US5 version delete → US4 overlap polish  
4. Phase 10 parity + quickstart  

### Suggested MVP scope

**Phases 1–4** (through US2): backend `publishedContext` + workspace routes + Assignments context — addresses reported production pain fastest.

---

## Notes

- Reuse 018 `ProfileTreeAssignmentPanel`, `profileRolloutService` — do not reimplement assignment APIs  
- Publish bump MUST be same transaction as `PublishVersion` ([contracts/publish-assignment-bump.md](./contracts/publish-assignment-bump.md))  
- Overview MUST NOT show draft body in cards when published exists (clarification A)  
- Total tasks: **62**
