# Tasks: Enrollment Routes — Controlled Onboarding Gateway

**Input**: Design documents from `/specs/021-enrollment-routes-ux/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Depends on**: [017-device-control-plane](../017-device-control-plane/tasks.md) enrollment routes + device tree; [019-profile-hub-ux](../019-profile-hub-ux/tasks.md) enrollment decoupling baseline (optional profile on create)

**Tests**: Unit tests for `intent_resolve` and validation per constitution IV (not full TDD).

**Organization**: Tasks grouped by user story (US1–US6). Foundational backend blocks all UI stories.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks in same phase)
- **[Story]**: US1–US6 from [spec.md](./spec.md)

## Path Conventions

- Backend: `serverBackendGo/internal/modules/enrollment_routes/`
- Migrations: `serverBackendGo/db/migrations/`
- Frontend: `frontend/src/features/enrollment-routes/`
- Parity: `serverBackendGo/docs/parity/enrollment-routes-ux.md`
- Contracts: `specs/021-enrollment-routes-ux/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify baseline and align contracts before implementation

- [X] T001 Confirm branch `021-enrollment-routes-ux` and modules enabled per [quickstart.md](./quickstart.md) prerequisites (`MODULE_ENROLLMENT_ROUTES_ENABLED`, device tree)
- [X] T002 [P] Review `specs/021-enrollment-routes-ux/contracts/` against `serverBackendGo/internal/modules/enrollment_routes/adapter/http/handler.go`
- [X] T003 [P] Create parity stub `serverBackendGo/docs/parity/enrollment-routes-ux.md` from [contracts/enrollment-routes-admin-api.md](./contracts/enrollment-routes-admin-api.md)
- [X] T004 [P] Add `ENROLLMENT_TREE_HEAVY_DEVICE_THRESHOLD` to `serverBackendGo/.env.example` and `serverBackendGo/internal/config/config.go` (default 500)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Migration `000030`, domain split, intent resolution, profile-free service layer, QR telemetry — MUST complete before user story UI work

**⚠️ CRITICAL**: No US1–US6 UI work until this phase checkpoint passes

- [X] T005 Create migration `serverBackendGo/db/migrations/000030_enrollment_routes_ux.up.sql` and `.down.sql` per [data-model.md](./data-model.md) (bootstrap columns, `is_recommended`, `container_placement_ack_at`, devices FK SET NULL, backfill)
- [X] T006 [P] Add `BootstrapIntent`, `EnrollmentRouteDefinition`, `EnrollmentRouteRuntimeState`, `EnrollmentRouteView` in `serverBackendGo/internal/modules/enrollment_routes/domain/route.go` and `domain/bootstrap_intent.go`
- [X] T007 [P] Add `EnrollmentDeleteImpact` in `serverBackendGo/internal/modules/enrollment_routes/domain/impact.go`
- [X] T008 Implement `ResolveBootstrapIntent` in `serverBackendGo/internal/modules/enrollment_routes/application/intent_resolve.go` (stable / specific / latest)
- [X] T009 [P] Unit tests for intent resolution in `serverBackendGo/internal/modules/enrollment_routes/application/intent_resolve_test.go`
- [X] T010 Remove `profileVersionId` requirement from `validateBinding` in `serverBackendGo/internal/modules/enrollment_routes/application/service.go`; map requests to `EnrollmentRouteView` without profile fields
- [X] T011 Extend `serverBackendGo/internal/modules/enrollment_routes/adapter/persistence/postgres/route_repo.go` for bootstrap columns, view mapping, and `ListBootstrapApps` / `ListTreeNodeOptions` queries
- [X] T012 Implement `Impact` in `serverBackendGo/internal/modules/enrollment_routes/application/impact.go` (three metrics per research R6)
- [X] T013 Register `GET /:id/impact`, `GET /options/tree-nodes`, `GET /options/bootstrap-apps` in `serverBackendGo/internal/modules/enrollment_routes/adapter/http/handler.go`
- [X] T014 Deprecate or return empty `GET /options/published-profile-versions` in `serverBackendGo/internal/modules/enrollment_routes/adapter/http/handler.go`
- [X] T015 Emit `enrollment_route.qr_viewed` `domain_events` on public QR/JSON hit when key resolves to enrollment route (extend qrcode/public handler per [research.md](./research.md) R6)
- [X] T016 [P] Strip `profileId`/`profileVersionId` from `frontend/src/features/enrollment-routes/enrollmentRouteService.ts` types; add `EnrollmentRouteView`, impact/options API methods per [contracts/enrollment-routes-admin-api.md](./contracts/enrollment-routes-admin-api.md)
- [X] T017 [P] Create `frontend/src/features/enrollment-routes/enrollmentRouteDialogState.ts` with state IDs from [contracts/enrollment-route-dialog-ux.md](./contracts/enrollment-route-dialog-ux.md)

**Checkpoint**: `make migrate` through 000030; `go test ./internal/modules/enrollment_routes/...`; curl POST route without `profileVersionId` → 201

---

## Phase 3: User Story 1 — Create gateway without profile concepts (Priority: P1) 🎯 MVP

**Goal**: Create enrollment route with target node + bootstrap app + identity mode; no profile/policy in API or UI

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 1 — POST without profile; create from UI without profile fields

### Implementation for User Story 1

- [X] T018 [P] [US1] Create `frontend/src/features/enrollment-routes/EnrollmentRouteForm.tsx` (name, description, identity mode only — no profile copy)
- [X] T019 [US1] Wire Create/Update payloads with `bootstrapIntent` + `bootstrapApplicationId` in `frontend/src/features/enrollment-routes/enrollmentRouteService.ts`
- [X] T020 [US1] Integrate `EnrollmentRouteForm` into `EnrollmentRouteDialog.tsx` CREATE mode in `frontend/src/features/enrollment-routes/EnrollmentRouteDialog.tsx`
- [X] T021 [US1] Remove published-profile gate from `handleNewRoute` in `frontend/src/features/enrollment-routes/EnrollmentRouteListPage.tsx`
- [X] T022 [P] [US1] Audit enrollment feature for forbidden strings (`profile`, `policy`, `سياسة`, `برفايل`) in `frontend/src/features/enrollment-routes/` and `frontend/src/i18n/locales/en.json` / `ar.json` keys under `enrollmentRoute.*`
- [X] T023 [US1] Handle `error.enrollment_route.stable_version_missing` and validation errors in `EnrollmentRouteForm.tsx`

**Checkpoint**: SC-004 — no profile fields in network responses or form; create succeeds without published profile

---

## Phase 4: User Story 2 — Dialog with live dual-column QR (Priority: P1)

**Goal**: Pending QR (client) updates live; Active QR (server) after save; Draft/Active/Unsaved badges

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 3 — Pending updates on field change; Active after save

### Implementation for User Story 2

- [X] T024 [P] [US2] Implement `buildEnrollmentContractPreview` in `frontend/src/features/enrollment-routes/buildEnrollmentContractPreview.ts` per [contracts/enrollment-contract-payload.md](./contracts/enrollment-contract-payload.md)
- [X] T025 [US2] Create `frontend/src/features/enrollment-routes/EnrollmentRouteQrColumn.tsx` (Pending watermark, disable copy when preview; Active via `enrollmentQrQuery` + `loadQrImageObjectUrl`)
- [X] T026 [US2] Create `frontend/src/features/enrollment-routes/EnrollmentRouteDialogHeader.tsx` (Draft / Active / Unsaved badges per clarifications)
- [X] T027 [US2] Wire two-column layout (form left, QR right) in `frontend/src/features/enrollment-routes/EnrollmentRouteDialog.tsx`
- [X] T028 [US2] Extend `GET .../qr` response mapping for resolved package/version in `serverBackendGo/internal/modules/enrollment_routes/application/service.go` and handler
- [X] T029 [US2] Enforce FR-006b: document in dialog footer that last Active QR remains scannable during unsaved Edit in `EnrollmentRouteDialog.tsx`

**Checkpoint**: SC-007 — Pending before save, Active after; no server preview API calls while editing

---

## Phase 5: User Story 3 — Target node with kind + context (Priority: P1)

**Goal**: Hierarchical tree picker with placement kind, breadcrumb, device count, heavily loaded + container warnings

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 4 — inheritable warning; invalid node blocked

### Implementation for User Story 3

- [X] T030 [US3] Implement `placementKind` computation in `serverBackendGo/internal/modules/enrollment_routes/application/tree_options.go` (locked vs inheritable per [research.md](./research.md) R4)
- [X] T031 [US3] Wire `GET /options/tree-nodes` handler to tree options service in `serverBackendGo/internal/modules/enrollment_routes/adapter/http/handler.go`
- [x] T032 [P] [US3] Create `frontend/src/features/enrollment-routes/TargetNodePicker.tsx` (inline panel inside dialog, single select, context preview)
- [x] T033 [US3] Integrate `TargetNodePicker` into `EnrollmentRouteForm.tsx` with selected path display
- [X] T034 [US3] Show persistent container warning in Overview + `acknowledgeContainerPlacement` on save in `EnrollmentRouteDialog.tsx` / form
- [X] T035 [US3] Persist `container_placement_ack_at` on save in `serverBackendGo/internal/modules/enrollment_routes/adapter/persistence/postgres/route_repo.go`

**Checkpoint**: User Story 3 acceptance scenarios; SC-003 measurable in manual test

---

## Phase 6: User Story 4 — Bootstrap app intent (Priority: P2)

**Goal**: Stable / Specific / Latest picker with resolved version display

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 1 stable resolution + Sprint 4 specific pin

### Implementation for User Story 4

- [X] T036 [US4] Implement `ListBootstrapApps` with `isRecommended` / `isLatest` flags in `serverBackendGo/internal/modules/enrollment_routes/adapter/persistence/postgres/route_repo.go`
- [X] T037 [US4] Wire `GET /options/bootstrap-apps` in `serverBackendGo/internal/modules/enrollment_routes/adapter/http/handler.go`
- [x] T038 [P] [US4] Create `frontend/src/features/enrollment-routes/BootstrapAppPicker.tsx` (intent radio + version dropdown when specific)
- [x] T039 [US4] Integrate `BootstrapAppPicker` into `EnrollmentRouteForm.tsx`; show resolved version line in Overview
- [X] T040 [US4] Update `buildEnrollmentContractPreview.ts` when intent/app changes for Pending QR

**Checkpoint**: Stable uses `is_recommended` row, not latest; save fails clearly when stable missing

---

## Phase 7: User Story 5 — Safe delete with multi-dimensional impact (Priority: P2)

**Goal**: Delete shows three impact metrics; typed route name when any impact > 0; historical devices do not block delete

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 5

### Implementation for User Story 5

- [X] T041 [US5] Implement `DELETE` with impact pre-check in `serverBackendGo/internal/modules/enrollment_routes/application/service.go`
- [X] T042 [P] [US5] Create `frontend/src/features/enrollment-routes/DeleteRouteConfirm.tsx` (inline steps DELETE_STEP1/2 per contract — no nested dialog)
- [X] T043 [US5] Wire delete flow from `EnrollmentRouteDialog.tsx` footer through `DeleteRouteConfirm.tsx`
- [X] T044 [US5] Add `getEnrollmentRouteImpact` consumer in `frontend/src/features/enrollment-routes/enrollmentRouteService.ts`
- [x] T045 [US5] Unit test impact counts in `serverBackendGo/internal/modules/enrollment_routes/application/impact_test.go`

**Checkpoint**: SC-006 — all deletes show three dimensions; typed name when impact > 0

---

## Phase 8: User Story 6 — List hub + overview/edit in one shell (Priority: P2)

**Goal**: List opens Overview dialog; Edit in same shell; remove full-page editor routes

**Independent Test**: [quickstart.md](./quickstart.md) Sprint 3 — no `/enrollment-routes/:id` page navigation

### Implementation for User Story 6

- [X] T046 [US6] Refactor `frontend/src/features/enrollment-routes/EnrollmentRouteListPage.tsx` to open `EnrollmentRouteDialog` on row click (OVERVIEW mode)
- [X] T047 [US6] Implement Overview read-only column + Edit transition in `frontend/src/features/enrollment-routes/EnrollmentRouteDialog.tsx`
- [X] T048 [US6] Remove routes `/enrollment-routes/new` and `/enrollment-routes/:id` from `frontend/src/app/App.tsx`; redirect legacy URLs if needed
- [X] T049 [P] [US6] Remove or thin `frontend/src/features/enrollment-routes/EnrollmentRouteEditorPage.tsx` to re-export dialog-only flow
- [x] T050 [US6] Update `serverBackendGo/internal/modules/profiles/application/onboarding.go` checklist paths to enrollment dialog (no profile prerequisite text)

**Checkpoint**: SC-002 — 100% create/edit from dialog; list context preserved

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: i18n, parity completion, onboarding, docs

- [X] T051 [P] Add `enrollmentRoute.*` i18n keys (status, qr, tree, delete, bootstrap) in `frontend/src/i18n/locales/en.json` and `frontend/src/i18n/locales/ar.json` — no profile/policy keys
- [x] T052 Complete `serverBackendGo/docs/parity/enrollment-routes-ux.md` with Java refs and breaking DTO notes
- [x] T053 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` row for enrollment routes UX if present
- [x] T054 Run [quickstart.md](./quickstart.md) Sprints 1–6 and fix gaps
- [x] T055 [P] Mobile sheet layout for `EnrollmentRouteDialog.tsx` (stack config then QR per contract)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Start immediately
- **Foundational (Phase 2)**: After Setup — **blocks all user stories**
- **US1 (Phase 3)**: After Foundational — **MVP** (backend + minimal create dialog)
- **US2 (Phase 4)**: After US1 form exists (shares dialog)
- **US3 (Phase 5)**: After Foundational; integrates into form (can parallel US2 after T018)
- **US4 (Phase 6)**: After Foundational; integrates into form (parallel with US3 after T018)
- **US5 (Phase 7)**: After Foundational impact API (T012–T013); UI after dialog shell (T020)
- **US6 (Phase 8)**: After US1–US2 dialog layout; completes list/overview routing
- **Polish (Phase 9)**: After desired stories complete

### User Story Dependencies

| Story | Depends on |
|-------|------------|
| US1 | Phase 2 |
| US2 | US1 form + dialog shell (T020) |
| US3 | Phase 2 tree options (T030–T031); form (T018) |
| US4 | Phase 2 intent (T008); form (T018) |
| US5 | Phase 2 impact (T012–T013); dialog (T020) |
| US6 | US1–US2 dialog (best after T027) |

### Parallel Opportunities

- Phase 1: T002–T004 parallel
- Phase 2: T006–T007, T009, T016–T017 parallel after T005
- US1: T018 ∥ T022
- US2: T024 ∥ T026
- US3: T032 ∥ T030 (backend) after T031
- US4: T038 ∥ T036
- US5: T042 ∥ T045
- US6: T049 ∥ T051

### Parallel Example: Foundational backend

```bash
# After T005 migration:
T006 domain types
T007 impact domain
T008 intent_resolve.go
T009 intent_resolve_test.go
T016 frontend service types
T017 dialog state types
# Then:
T010 → T011 → T012 → T013 → T014 → T015
```

### Parallel Example: P1 UI (after T020 dialog exists)

```bash
Developer A: US2 (T024–T029) QR column
Developer B: US3 (T032–T035) tree picker
Developer C: US4 (T038–T040) bootstrap picker
```

---

## Implementation Strategy

### MVP First (Recommended)

1. Phase 1 Setup  
2. Phase 2 Foundational (**required**)  
3. Phase 3 US1 — profile-free create  
4. Phase 4 US2 — Pending/Active QR column  
5. **STOP**: Validate quickstart Sprints 1 & 3  

### Incremental delivery

1. Foundation → US1 → US2 (P1 gateway + QR)  
2. US3 tree picker → US4 bootstrap intent  
3. US5 delete safety → US6 list/overview routing cleanup  
4. Phase 9 polish + full quickstart  

### Suggested MVP scope

**Phases 1–4** (through US2): admin can create a policy-free enrollment gateway with live Pending QR and Active QR after save.

---

## Notes

- Do **not** import `profiles` package from `enrollment_routes` application layer (compat only in sync resolver)  
- `qrcodekey` immutable v1 — no rotate task in this feature  
- Pending QR: client only ([contracts/enrollment-contract-payload.md](./contracts/enrollment-contract-payload.md))  
- Total tasks: **55**
