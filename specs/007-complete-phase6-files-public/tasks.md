---
description: "Task list for Phase 6 Files, Icons & Public API migration"
---

# Tasks: Phase 6 — Files, Icons & Public API

**Input**: `specs/007-complete-phase6-files-public/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–5 complete; Postgres via `./scripts/db-up.sh`; `FILES_DIRECTORY` writable

**Tests**: Included per FR-X05, FR-X06, and User Story 6.

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Files: `serverBackendGo/internal/modules/files/`
- Icons: `serverBackendGo/internal/modules/icons/`
- Public API: `serverBackendGo/internal/modules/publicapi/`
- Storage: `serverBackendGo/internal/platform/storage/`
- Platform: `serverBackendGo/internal/platform/auth/`, `serverBackendGo/internal/config/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 6 context and Java/React parity baseline.

- [X] T001 Verify feature context in `specs/007-complete-phase6-files-public/spec.md` against `serverBackendGo/docs/MIGRATION.md` Phase 6 pending row
- [X] T002 [P] Review Java `backend/server/src/main/java/com/hmdm/rest/resource/FilesResource.java`, `IconResource.java`, and `PublicResource.java` against `specs/007-complete-phase6-files-public/contracts/`
- [X] T003 [P] Review React `frontend/src/features/files/filesService.ts`, `frontend/src/features/applications/services/webUiFilesService.ts`, and `frontend/src/features/icons/iconsService.ts` for required JSON shapes
- [X] T004 Run baseline `cd serverBackendGo && go build ./...` and note `files`, `icons`, `publicapi` scaffolds in `serverBackendGo/internal/modules/*/module.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, shared storage, permissions, and module skeletons required by all user stories.

**⚠️ CRITICAL**: No endpoint work until migration `000008` applies and `platform/storage` exists.

- [X] T005 Create `serverBackendGo/db/migrations/000008_files_icons_core.up.sql` per `data-model.md` (`uploadedfiles`, `icons`, `configurationfiles.fileid`, `customers.sizelimit`, indexes)
- [X] T006 [P] Create `serverBackendGo/db/migrations/000008_files_icons_core.down.sql`
- [X] T007 Seed `permissions` rows `files`, `edit_files` and `userrolepermissions` for role 2 in `000008_files_icons_core.up.sql`
- [X] T008 [P] Add optional dev seed: sample `uploadedfiles` row + `icons` row for customer 1 in `000008_files_icons_core.up.sql`
- [X] T009 Extend `serverBackendGo/internal/platform/auth/permissions.go` with `PermFiles`, `PermEditFiles` and helper methods
- [X] T010 [P] Add tests in `serverBackendGo/internal/platform/auth/permissions_test.go` for files and edit_files
- [X] T011 Implement `serverBackendGo/internal/platform/storage/local.go` — `IsSafePath`, `CreateTemp`, `MoveToCustomer`, `DeleteFile`, `DirSizeMB`, `BuildPublicURL`
- [X] T012 [P] Add `serverBackendGo/internal/platform/storage/local_test.go` for path traversal and temp-prefix guards
- [X] T013 Extend `serverBackendGo/internal/config/config.go` with `HashSecret`, rebranding fields, `ModuleFilesEnabled`, `ModuleIconsEnabled`, `ModulePublicAPIEnabled`
- [X] T014 [P] Document new env vars in `serverBackendGo/.env.example` (`HASH_SECRET`, `REBRANDING_*`, `MODULE_*_ENABLED`)
- [X] T015 [P] Create `serverBackendGo/internal/modules/files/domain/file.go` — `UploadedFile`, `FileView`, `FileUploadResult`, `LimitResponse`, link DTOs per `contracts/files-api.md`
- [X] T016 [P] Create `serverBackendGo/internal/modules/icons/domain/icon.go` per `contracts/icons-api.md`
- [X] T017 [P] Create `serverBackendGo/internal/modules/publicapi/domain/public.go` — `NameResponse`, `UploadAppRequest` per `contracts/publicapi-api.md`
- [X] T018 Define `serverBackendGo/internal/modules/files/port/repository.go` — file CRUD, usage checks, configuration links; `CustomerReader` for filesDir/sizeLimit
- [X] T019 [P] Define `serverBackendGo/internal/modules/icons/port/repository.go`
- [X] T020 [P] Define `serverBackendGo/internal/modules/publicapi/port/device.go` and `application.go` — unsecure device lookup, insert application
- [X] T021 Define `serverBackendGo/internal/modules/files/port/push.go` no-op notifier stub
- [X] T022 Verify migration: `make dev` and confirm `uploadedfiles`, `icons` tables; `configurationfiles.fileid` column exists
- [X] T023 [P] Wire module flags in `serverBackendGo/internal/app/modules.go` for files/icons/publicapi when env enabled

**Checkpoint**: Schema + storage + ports ready — proceed to user stories.

---

## Phase 3: User Story 1 — Browse and manage the file library (Priority: P1) 🎯 MVP (part 1)

**Goal**: `GET /rest/private/web-ui-files/search`, `GET /search/{value}`, `POST /remove`.

**Independent Test**: JWT with `files` lists files; JWT with `edit_files` removes unused file; FILE_USED when linked.

### Tests for User Story 1

- [X] T024 [P] [US1] Add `serverBackendGo/internal/modules/files/application/service_test.go` for search/remove permission and FILE_USED
- [X] T025 [P] [US1] Add `serverBackendGo/internal/modules/files/adapter/http/handler_test.go` for GET `/search` and POST `/remove`

### Implementation for User Story 1

- [X] T026 [US1] Implement `serverBackendGo/internal/modules/files/adapter/persistence/postgres/file_repo.go` — `List`, `ListByValue`, `GetByID`, `Delete`, `IsUsedByConfiguration`, `IsUsedByIcon`
- [X] T027 [US1] Implement `Search`, `Remove` in `serverBackendGo/internal/modules/files/application/service.go` with `files`/`edit_files` and `platform/storage` delete
- [X] T028 [US1] Create `serverBackendGo/internal/modules/files/adapter/http/handler.go` — register `GET /search`, `GET /search/:value`, `POST /remove` with Swagger comments
- [X] T029 [US1] Wire `serverBackendGo/internal/modules/files/module.go` — repo → service → handler on `/web-ui-files`

**Checkpoint**: Files page list and delete work (`filesService.ts`).

---

## Phase 4: User Story 2 — Upload and commit files (Priority: P1) 🎯 MVP (part 2)

**Goal**: `POST /`, `POST /raw`, `POST /update` for APK/asset upload and commit.

**Independent Test**: Multipart upload returns `FileUploadResult`; POST update persists row and moves file from temp.

### Tests for User Story 2

- [X] T030 [P] [US2] Extend `serverBackendGo/internal/modules/files/application/service_test.go` for upload quota, duplicate path, unsafe tmp path
- [X] T031 [P] [US2] Extend `serverBackendGo/internal/modules/files/adapter/http/handler_test.go` for POST `/` and POST `/update`

### Implementation for User Story 2

- [X] T032 [P] [US2] Add `serverBackendGo/internal/modules/files/application/apkparse.go` — best-effort APK metadata (package, version, versionCode, arch) per `research.md` R3
- [X] T033 [US2] Extend `file_repo.go` with `Insert`, `Update`, `FindByPath`, duplicate checks
- [X] T034 [US2] Implement `Upload`, `UploadRaw`, `Create`, `CreateExternal`, `Update` in `serverBackendGo/internal/modules/files/application/service.go` using `platform/storage` and optional version-exists queries via applications port
- [X] T035 [US2] Register `POST /`, `POST /raw`, `POST /update` in `serverBackendGo/internal/modules/files/adapter/http/handler.go`

**Checkpoint**: Applications APK upload flow (`webUiFilesService.ts`) succeeds without Java.

---

## Phase 5: User Story 3 — File usage and configuration links (Priority: P2)

**Goal**: `GET /limit`, `GET /apps/{url}`, `GET /configurations/{id}`, `POST /configurations`.

**Independent Test**: GET limit returns MB used/limit; GET apps by URL returns applications; POST configurations updates links.

### Tests for User Story 3

- [X] T036 [P] [US3] Add tests in `serverBackendGo/internal/modules/files/application/service_test.go` for storage limit and configuration link scope filtering

### Implementation for User Story 3

- [X] T037 [US3] Extend `file_repo.go` — `GetFileConfigurations`, `UpdateFileConfigurations`, usage count queries for FileView enrichment
- [X] T038 [US3] Add applications URL lookup adapter or reuse `serverBackendGo/internal/modules/applications/port/repository.go` from files service via narrow port
- [X] T039 [US3] Implement `GetLimit`, `GetApplicationsByURL`, `GetFileConfigurations`, `UpdateFileConfigurations` in `serverBackendGo/internal/modules/files/application/service.go` (push stub via `port/push.go`)
- [X] T040 [US3] Register `GET /limit`, `GET /apps/*url`, `GET /configurations/:id`, `POST /configurations` in `serverBackendGo/internal/modules/files/adapter/http/handler.go`

**Checkpoint**: Storage limit and file-configuration endpoints match contracts.

---

## Phase 6: User Story 4 — Manage launcher icons (Priority: P2)

**Goal**: `GET /search`, `GET /search/{value}`, `PUT /`, `DELETE /{id}` on `/rest/private/icons`.

**Independent Test**: List icons; PUT save; DELETE requires `settings` permission.

### Tests for User Story 4

- [X] T041 [P] [US4] Add `serverBackendGo/internal/modules/icons/application/service_test.go` for search and delete permission
- [X] T042 [P] [US4] Add `serverBackendGo/internal/modules/icons/adapter/http/handler_test.go` for GET `/search` and DELETE `/:id`

### Implementation for User Story 4

- [X] T043 [US4] Implement `serverBackendGo/internal/modules/icons/adapter/persistence/postgres/icon_repo.go` — `List`, `ListByValue`, `Insert`, `Update`, `Delete` with tenant scope
- [X] T044 [US4] Implement `Search`, `Save`, `Delete` in `serverBackendGo/internal/modules/icons/application/service.go`
- [X] T045 [US4] Create `serverBackendGo/internal/modules/icons/adapter/http/handler.go` with routes and Swagger comments
- [X] T046 [US4] Wire `serverBackendGo/internal/modules/icons/module.go` — repo → service → handler on `/icons`

**Checkpoint**: Icons UI (`iconsService.ts`) works end-to-end.

---

## Phase 7: User Story 5 — Public rebranding and AppList upload (Priority: P2)

**Goal**: `GET /rest/public/name`, `GET /logo`, `POST /applications/upload`.

**Independent Test**: GET name returns rebranding JSON; invalid hash rejected; valid device hash creates application.

### Tests for User Story 5

- [X] T047 [P] [US5] Add `serverBackendGo/internal/modules/publicapi/application/service_test.go` for hash validation and device-not-found
- [X] T048 [P] [US5] Add `serverBackendGo/internal/modules/publicapi/adapter/http/handler_test.go` for GET `/name` and POST `/applications/upload` hash failure

### Implementation for User Story 5

- [X] T049 [US5] Implement `serverBackendGo/internal/modules/publicapi/adapter/persistence/postgres/device_repo.go` — unsecure device by number, duplicate pkg+version check, insert application
- [X] T050 [US5] Implement `GetRebranding`, `UploadApplication` in `serverBackendGo/internal/modules/publicapi/application/service.go` using `shared/crypto` MD5 and `platform/storage` for optional file write
- [X] T051 [US5] Create `serverBackendGo/internal/modules/publicapi/adapter/http/handler.go` — `GET /name`, `GET /logo` (stream/redirect), `POST /applications/upload` multipart
- [X] T052 [US5] Wire `serverBackendGo/internal/modules/publicapi/module.go` on `groups.Public` (no Bearer middleware)

**Checkpoint**: Public endpoints respond per `contracts/publicapi-api.md`.

---

## Phase 8: User Story 6 — Verifiable API and regression safety (Priority: P2)

**Goal**: Swagger tags, parity docs, module tests green.

**Independent Test**: `go test` passes; Swagger lists Files/Icons/Public; parity tables complete.

### Implementation for User Story 6

- [X] T053 [P] [US6] Add `@Router` Swagger comments on all handlers in `serverBackendGo/internal/modules/files/adapter/http/handler.go`, `icons/adapter/http/handler.go`, `publicapi/adapter/http/handler.go`
- [X] T054 [US6] Run `cd serverBackendGo && make swagger` and verify tags in `serverBackendGo/internal/platform/httpx/swagger/swagger.yaml`
- [X] T055 [P] [US6] Create `serverBackendGo/docs/parity/files.md` endpoint table (mark `GET /files/*` Partial)
- [X] T056 [P] [US6] Create `serverBackendGo/docs/parity/icons.md` endpoint table
- [X] T057 [P] [US6] Create `serverBackendGo/docs/parity/publicapi.md` endpoint table
- [X] T058 [US6] Run `go test ./internal/modules/files/... ./internal/modules/icons/... ./internal/modules/publicapi/... ./internal/platform/storage/...`

**Checkpoint**: CI-local test suite and Swagger aligned with contracts.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Refactors, docs, optional static files, migration status.

- [X] T059 [P] Refactor `serverBackendGo/internal/modules/configfiles/adapter/http/handler.go` to use `internal/platform/storage` for disk writes and URL building
- [X] T060 [P] Optional: add `GET /files/*` static or signed handler in `serverBackendGo/internal/app/` for dev agent URLs — document **Partial** in `docs/parity/files.md`
- [X] T061 Update `serverBackendGo/docs/MIGRATION.md` — Phase 6 row **done**; modules `files`, `icons`, `publicapi`
- [X] T062 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` Phase 6 section to **منجز** with parity links
- [X] T063 [P] Update `serverBackendGo/docs/parity/applications.md` — mark `web-ui-files` **Done** (remove Phase 6 deferral)
- [X] T064 Run validation steps in `specs/007-complete-phase6-files-public/quickstart.md` (curl + React Files/Applications smoke)
- [X] T065 Final `cd serverBackendGo && go build ./...` and fix any layer import violations

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **BLOCKS** all user stories
- **Phase 3 (US1)**: Depends on Phase 2
- **Phase 4 (US2)**: Depends on Phase 2; integrates with US1 repo (same module)
- **Phase 5 (US3)**: Depends on US2 file repo; optional applications port from Phase 5
- **Phase 6 (US4)**: Depends on Phase 2 (`uploadedfiles` for `fileId` FK); independent of US1–US3 for API surface
- **Phase 7 (US5)**: Depends on Phase 2 storage + applications insert port
- **Phase 8 (US6)**: Depends on handlers from US1–US5
- **Phase 9 (Polish)**: Depends on US6 minimum; T059–T060 can run after US2

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 | Phase 2 | File list/delete |
| US2 | Phase 2, US1 repo | Same `files` module; extends repo/service |
| US3 | US2 | Link endpoints need persisted files |
| US4 | Phase 2 | Icons need `uploadedfiles` table |
| US5 | Phase 2 | Public upload uses storage + device lookup |
| US6 | US1–US5 | Swagger + parity |

### Parallel Opportunities

- Phase 1: T002, T003
- Phase 2: T006–T010, T012, T014–T020, T023 [P] after T005
- US1: T024–T025 [P]
- US2: T030–T032 [P] (apkparse parallel to tests)
- US4: T041–T042, T043–T045 track parallel to US3 after US2
- US5: T047–T048 [P]
- US6: T053–T057 [P]
- Polish: T059–T063 [P]

### Parallel Example: MVP (US1 + US2)

```bash
# After Phase 2 complete:
# Track A: T024–T029 (US1 list/remove)
# Track B: T030–T035 (US2 upload/commit) — same module, sequence T029 before T035 handler merge
# Validate: quickstart §4–§5 + React Files + Applications upload
```

### Parallel Example: US4 + US5 after US2

```bash
# Developer A: T041–T046 icons module
# Developer B: T047–T052 publicapi module
# Both only need Phase 2 + uploadedfiles seed
```

---

## Implementation Strategy

### MVP First (US1 + US2)

1. Complete Phase 1–2: Setup + Foundational (migration + storage critical)
2. Complete Phase 3: US1 file list/delete
3. Complete Phase 4: US2 upload/commit
4. **STOP and VALIDATE**: React Files page + Applications APK upload (`quickstart.md` §4–§7)
5. Continue US3 → US4 → US5 → US6 → Polish

### Incremental Delivery

1. Foundation → US1 + US2 (admin file library + APK upload)
2. US3 (limits + configuration links)
3. US4 (icons) + US5 (public) in parallel
4. US6 + Polish (Swagger, parity, MIGRATION done)

### Suggested MVP Scope

**Minimum**: Phases 1–2 + **US1** + **US2** + T058 subset + T064 smoke.

Delivers Files page and Applications upload without Java; defers icons/public to next slice.

---

## Notes

- Table names follow Java casing in DB: `uploadedfiles` / `uploadedFiles` — use lowercase unquoted identifiers in Go migrations for Postgres consistency with existing `000007`.
- Icon DELETE requires `settings` permission only (Java parity); GET/PUT have no extra permission.
- Push on `POST /configurations` is stubbed until Phase 7.
- Branch: `007-complete-phase6-files-public`
