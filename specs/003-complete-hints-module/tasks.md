---
description: "Task list for Phase 3 Hints module migration"
---

# Tasks: Phase 3 — Hints Module Migration

**Input**: `specs/003-complete-hints-module/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Auth + JWT/Swagger Bearer (Phase 1–2); Postgres via `./scripts/db-up.sh`

**Tests**: Included per FR-007 and User Story 5.

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Go backend: `serverBackendGo/internal/modules/hints/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`, `specs/003-complete-hints-module/quickstart.md`

---

## Phase 1: Setup

**Purpose**: Confirm migration context and Java/React parity baseline.

- [x] T001 Verify feature context in `specs/003-complete-hints-module/spec.md` against `serverBackendGo/docs/NEXT_STEPS.md` #3 and `MIGRATION.md` Phase 3
- [x] T002 [P] Review Java `backend/server/src/main/java/com/hmdm/rest/resource/HintResource.java` and `UserMapper.java` hint SQL annotations
- [x] T003 [P] Review React `frontend/src/features/hints/hintsService.ts` and legacy `hint.service.js` endpoint usage
- [x] T004 Run baseline `cd serverBackendGo && go build ./...` and note current `internal/modules/hints/module.go` scaffold state

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Database tables and persistence layer required by all user stories.

**⚠️ CRITICAL**: No hint endpoint work until migration applies successfully.

- [x] T005 Create `serverBackendGo/db/migrations/000004_hints_tables.up.sql` with `userhints`, `userhinttypes`, unique `(userid, hintkey)`, FK to `users`, seed `hint.step.1`–`hint.step.4`
- [x] T006 [P] Add `serverBackendGo/db/migrations/000004_hints_tables.down.sql` dropping hint tables
- [x] T007 Create `serverBackendGo/internal/modules/hints/domain/hint.go` with `HintKey` validation (non-empty, max 100 chars)
- [x] T008 Define `serverBackendGo/internal/modules/hints/port/repository.go` with `GetHistory`, `MarkShown`, `Enable`, `Disable` methods
- [x] T009 Implement `serverBackendGo/internal/modules/hints/adapter/persistence/postgres/hint_repo.go` with all SQL per `research.md` (ON CONFLICT for mark-shown)
- [x] T010 Verify migration applies: restart `make dev` or run migrate and confirm tables exist in Postgres

**Checkpoint**: `hint_repo` compiles; tables present — proceed to user stories.

---

## Phase 3: User Story 1 — Restore hint history (Priority: P1) 🎯 MVP

**Goal**: `GET /rest/private/hints/history` returns `OK` + `string[]` for authenticated user.

**Independent Test**: JWT login → `GET /history` → `[]` or list of keys.

### Tests for User Story 1

- [x] T011 [P] [US1] Add `serverBackendGo/internal/modules/hints/application/service_test.go` stub test for `GetHistory` empty and populated lists
- [x] T012 [P] [US1] Add `serverBackendGo/internal/modules/hints/adapter/http/handler_test.go` for `GET /history` with 403 without auth and 200 with principal

### Implementation for User Story 1

- [x] T013 [US1] Implement `GetHistory` in `serverBackendGo/internal/modules/hints/application/service.go` using `platformauth.Principal`
- [x] T014 [US1] Create `serverBackendGo/internal/modules/hints/adapter/http/handler.go` with `GetHistory` handler and Headwind `response.OK`
- [x] T015 [US1] Register `GET /history` in handler `Register` on `groups.Private.Group("/hints")` in `serverBackendGo/internal/modules/hints/module.go`
- [x] T016 [US1] Wire `module.go`: repo → service → handler; require `deps.DB`; remove scaffold-only log

**Checkpoint**: `quickstart.md` §3 empty-history curl succeeds.

---

## Phase 4: User Story 2 — Record hint as shown (Priority: P1)

**Goal**: `POST /rest/private/hints/history` persists hint key (idempotent).

**Independent Test**: POST `"hint.step.1"` → GET includes key.

### Tests for User Story 2

- [x] T017 [P] [US2] Extend `service_test.go` for `MarkShown` duplicate key idempotent and empty key error
- [x] T018 [P] [US2] Extend `handler_test.go` for `POST /history` with JSON string body

### Implementation for User Story 2

- [x] T019 [US2] Implement `MarkShown` in `serverBackendGo/internal/modules/hints/application/service.go` with body validation
- [x] T020 [US2] Add `parseHintKeyBody` in `serverBackendGo/internal/modules/hints/adapter/http/handler.go` (JSON string, raw text per `research.md`)
- [x] T021 [US2] Add `MarkShown` handler and register `POST /history` in `serverBackendGo/internal/modules/hints/adapter/http/handler.go`

**Checkpoint**: Mark-shown smoke from `quickstart.md` §3 passes.

---

## Phase 5: User Story 3 — Re-enable tutorials (Priority: P2)

**Goal**: `POST /rest/private/hints/enable` clears user hint history.

**Independent Test**: After mark-shown → enable → GET returns `[]`.

### Tests for User Story 3

- [x] T022 [P] [US3] Add `service_test.go` case for `Enable` clearing history rows

### Implementation for User Story 3

- [x] T023 [US3] Implement `Enable` in `serverBackendGo/internal/modules/hints/application/service.go`
- [x] T024 [US3] Add `Enable` handler and register `POST /enable` in `serverBackendGo/internal/modules/hints/adapter/http/handler.go`

**Checkpoint**: Enable curl in `quickstart.md` §3 passes.

---

## Phase 6: User Story 4 — Disable all tutorials (Priority: P2)

**Goal**: `POST /rest/private/hints/disable` marks all catalog keys shown.

**Independent Test**: disable → GET returns four seed keys.

### Tests for User Story 4

- [x] T025 [P] [US4] Add `service_test.go` case for `Disable` inserting all catalog keys from stub repo

### Implementation for User Story 4

- [x] T026 [US4] Implement `Disable` in `serverBackendGo/internal/modules/hints/application/service.go`
- [x] T027 [US4] Add `Disable` handler and register `POST /disable` in `serverBackendGo/internal/modules/hints/adapter/http/handler.go`

**Checkpoint**: Disable curl in `quickstart.md` §3 returns four keys.

---

## Phase 7: User Story 5 — Verifiable API contract (Priority: P2)

**Goal**: Swagger, parity doc, full test pass, sign-off.

**Independent Test**: `make swagger` + `go test ./internal/modules/hints/...` + Swagger Authorize all four routes.

### Implementation for User Story 5

- [x] T028 [P] [US5] Add Swagger `// @Summary` and `// @Security BearerAuth` on all handlers in `serverBackendGo/internal/modules/hints/adapter/http/handler.go`
- [x] T029 [US5] Run `make swagger` in `serverBackendGo/` and commit regenerated `internal/platform/httpx/swagger/docs.go` if changed
- [x] T030 [P] [US5] Create `serverBackendGo/docs/parity/hints.md` marking all HintResource endpoints Done
- [x] T031 [US5] Update `serverBackendGo/docs/NEXT_STEPS.md` row #3 hints to **منجز**
- [x] T032 [US5] Update `serverBackendGo/docs/MIGRATION.md` Phase 3 hints row with link to `parity/hints.md`
- [x] T033 [US5] Run `specs/003-complete-hints-module/quickstart.md` full smoke and fix failures
- [x] T034 [US5] Run `cd serverBackendGo && go test ./internal/modules/hints/... -v` and `go build ./...`

**Checkpoint**: SC-001–SC-005 from spec.md satisfied.

---

## Phase 8: Polish & Cross-Cutting Concerns

- [x] T035 [P] Remove `serverBackendGo/internal/modules/hints/domain/doc.go` scaffold if still present
- [x] T036 Delete unused `serverBackendGo/internal/modules/hints/adapter/http/routes.go` scaffold or merge into `handler.go` Register only
- [x] T037 Manual E2E: React `/hints` page — Enable, Disable, list refresh with Go-only backend (`frontend` Vite → `:8080`)
- [x] T038 Verify unauthenticated `GET /rest/private/hints/history` returns 403 (not 404)

---

## Dependencies & Execution Order

### Phase Dependencies

```text
Phase 1 Setup
    ↓
Phase 2 Foundational (migration + repo) — BLOCKS all stories
    ↓
Phase 3 US1 (GET history) — MVP
    ↓
Phase 4 US2 (POST history)
    ↓
Phase 5 US3 (enable) — can parallel with US4 after US2
    ↓
Phase 6 US4 (disable)
    ↓
Phase 7 US5 (Swagger, parity, sign-off)
    ↓
Phase 8 Polish
```

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 | Phase 2 | MVP |
| US2 | US1 handler/module pattern | Same service file |
| US3 | US2 | Tests assume mark-shown works |
| US4 | US3 optional | Independent SQL path |
| US5 | US1–US4 handlers | Sign-off |

### Within Each User Story

- Tests alongside or immediately after application layer
- Handlers after application
- `module.go` wired in US1; extended routes in US2–US4

---

## Parallel Example: User Story 1

```bash
# After Phase 2:
T011 service_test.go
T012 handler_test.go
T007 domain/hint.go   # if not done in Phase 2
```

---

## Parallel Example: User Story 5

```bash
T028 Swagger comments
T030 parity/hints.md
T032 MIGRATION.md update
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1 + Phase 2.
2. Complete Phase 3 (US1): `GET /history` only.
3. **STOP and VALIDATE**: `quickstart.md` first curl block.
4. Demo React Hints page list load (may be empty).

### Incremental Delivery

1. US1 → US2 → US3 → US4 → US5 → Polish.
2. Each story adds one endpoint family independently testable.

### Suggested MVP Scope

- **Minimum**: Through T016 (Phase 2 + US1) — 16 tasks.
- **Full Phase 3**: Through T034 — 34 tasks.

---

## Task Summary

| Phase | Task IDs | Count |
|-------|----------|-------|
| Setup | T001–T004 | 4 |
| Foundational | T005–T010 | 6 |
| US1 History | T011–T016 | 6 |
| US2 Mark shown | T017–T021 | 5 |
| US3 Enable | T022–T024 | 3 |
| US4 Disable | T025–T027 | 3 |
| US5 Verify | T028–T034 | 7 |
| Polish | T035–T038 | 4 |
| **Total** | T001–T038 | **38** |

**Parallel opportunities**: 14 tasks marked `[P]`

**Independent test criteria**: See each phase **Checkpoint** and spec.md User Story sections.
