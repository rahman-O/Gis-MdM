---
description: "Task list for Phase 8 Plugins Platform & Extension Modules migration"
---

# Tasks: Phase 8 — Plugins Platform & Extension Modules

**Input**: `specs/009-complete-phase8-plugins/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–7 complete; Postgres via `./scripts/db-up.sh`; migration `000009` applied; seeded device `hmdm-001`

**Tests**: Included per spec FR-011, FR-012, and constitution IV (application-layer + selective handler tests).

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Platform: `serverBackendGo/internal/modules/plugins/platform/`
- Audit: `serverBackendGo/internal/modules/plugins/audit/`
- Messaging: `serverBackendGo/internal/modules/plugins/messaging/`
- Device info: `serverBackendGo/internal/modules/plugins/deviceinfo/`
- Device log: `serverBackendGo/internal/modules/plugins/devicelog/`
- Push plugin: `serverBackendGo/internal/modules/plugins/push/`
- Shared: `serverBackendGo/internal/modules/plugins/shared/`
- Migrations: `serverBackendGo/db/migrations/`
- Parity: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 8 context and Java/React parity baseline.

- [X] T001 Verify feature context in `specs/009-complete-phase8-plugins/spec.md` against `serverBackendGo/docs/MIGRATION.md` Phase 8 pending row
- [X] T002 [P] Review Java `PluginResource.java`, `AuditResource.java`, `MessagingResource.java`, `DeviceInfoResource.java`, `DeviceLogResource.java`, `PushResource.java` against `specs/009-complete-phase8-plugins/contracts/`
- [X] T003 [P] Review React `frontend/src/features/plugins/pluginService.ts` for `/plugin/main/private/*` paths
- [X] T004 Run baseline `cd serverBackendGo && go build ./...` and note scaffolds in `internal/modules/plugins/*/module.go` (platform, audit, messaging, deviceinfo, devicelog)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Migration `000010`, shared infrastructure, permissions, and ports before any plugin endpoints.

**⚠️ CRITICAL**: No user story HTTP handlers until migration applies and shared `targets` + `status` ports exist.

- [X] T005 Create `serverBackendGo/db/migrations/000010_plugins_core.up.sql` — `plugins`, `pluginsDisabled`, `plugin_audit_log`, `plugin_messaging_messages`, deviceinfo tables, devicelog postgres tables per `data-model.md`
- [X] T006 [P] Create `serverBackendGo/db/migrations/000010_plugins_core.down.sql`
- [X] T007 Seed `plugins` catalog rows (audit, push, messaging, deviceinfo, devicelog) and permissions (`plugins_customer_access_management`, `plugin_audit_access`, `plugin_messaging_send`, `plugin_messaging_delete`, `plugin_deviceinfo_access`, `plugin_devicelog_access`) with role 2 grants in `000010_plugins_core.up.sql`
- [X] T008 [P] Add dev smoke seeds: sample `plugin_audit_log` rows and disabled-plugin fixture in `000010_plugins_core.up.sql`
- [X] T009 Extend `serverBackendGo/internal/platform/auth/permissions.go` with Phase 8 plugin permission constants
- [X] T010 [P] Extend `serverBackendGo/internal/platform/auth/permissions_test.go` for new permissions
- [X] T011 Extend `serverBackendGo/internal/config/config.go` with `EnabledPlugins`, `ModulePluginsEnabled`, `ModulePluginsPlatformEnabled`, `ModulePluginsAuditEnabled`, `ModulePluginsMessagingEnabled`, `ModulePluginsDeviceinfoEnabled`, `ModulePluginsDevicelogEnabled`
- [X] T012 [P] Document Phase 8 env vars in `serverBackendGo/.env.example` and export in `serverBackendGo/scripts/dev.sh`
- [X] T013 Create `serverBackendGo/internal/modules/plugins/shared/status/port.go` and `cache.go` — customer disabled plugin cache per `research.md` R6
- [X] T014 [P] Create `serverBackendGo/internal/modules/plugins/shared/targets/port.go` — extract device/group/broadcast resolution interface from `plugins/push/adapter/persistence/postgres/targets_repo.go`
- [X] T015 Implement `serverBackendGo/internal/modules/plugins/shared/targets/postgres/resolver.go` — move/refactor SQL from `targets_repo.go` (use `pq.Array` for ANY clauses)
- [X] T016 Refactor `serverBackendGo/internal/modules/plugins/push/adapter/persistence/postgres/targets_repo.go` to delegate to `shared/targets` (no behavior regression)
- [X] T017 [P] Create `serverBackendGo/internal/modules/plugins/platform/domain/plugin.go` — Plugin, DisabledPlugin DTOs per `contracts/plugins-platform-api.md`
- [X] T018 [P] Create `serverBackendGo/internal/modules/plugins/audit/domain/audit.go` — AuditLogRecord, AuditLogFilter
- [X] T019 [P] Create `serverBackendGo/internal/modules/plugins/messaging/domain/message.go` — Message, MessageFilter, SendRequest
- [X] T020 [P] Create `serverBackendGo/internal/modules/plugins/deviceinfo/domain/deviceinfo.go` — settings and dynamic info DTOs per contract
- [X] T021 [P] Create `serverBackendGo/internal/modules/plugins/devicelog/domain/devicelog.go` — settings, rules, log record DTOs per contract
- [X] T022 Define `serverBackendGo/internal/modules/plugins/platform/port/repository.go` — available/active/registered/disabled plugin queries
- [X] T023 Wire `notifications/port.MessageQueue` into `serverBackendGo/internal/app/modules.go` for injection into messaging (and existing push) modules
- [X] T024 Verify migration: `cd serverBackendGo && make migrate` and confirm `plugins`, `plugin_audit_log`, `plugin_messaging_messages` exist
- [X] T025 Add `PermPluginsCustomerAccess` (or equivalent) usage pattern in platform service for `POST /private/disabled`

**Checkpoint**: Schema + config + shared ports ready — proceed to user stories.

---

## Phase 3: User Story 1 — Tenant plugin catalog and enablement (Priority: P1) 🎯 MVP

**Goal**: React Plugin settings works via `/rest/plugin/main/private/*` and public registered list.

**Independent Test**: `quickstart.md` §4 — active/available GET, disabled POST, registered GET without auth.

### Tests for User Story 1

- [X] T026 [P] [US1] Add `serverBackendGo/internal/modules/plugins/platform/application/service_test.go` — enabled list filter, disabled merge, permission gate
- [X] T027 [P] [US1] Add `serverBackendGo/internal/modules/plugins/platform/adapter/http/handler_test.go` for GET active/available and POST disabled

### Implementation for User Story 1

- [X] T028 [US1] Implement `serverBackendGo/internal/modules/plugins/platform/adapter/persistence/postgres/plugin_repo.go` — findAvailable, findActive, findRegistered, saveDisabled
- [X] T029 [US1] Implement `serverBackendGo/internal/modules/plugins/platform/application/service.go` — `ENABLED_PLUGINS` filter + `shared/status` cache invalidation on save
- [X] T030 [US1] Create `serverBackendGo/internal/modules/plugins/platform/adapter/http/handler.go` — register routes on `groups.PluginMain` with JWT middleware on `/private` subgroup only per `research.md` R1
- [X] T031 [US1] Wire `serverBackendGo/internal/modules/plugins/platform/module.go` — replace scaffold; `MODULE_PLUGINS_PLATFORM_ENABLED`
- [X] T032 [US1] Create `serverBackendGo/docs/parity/plugins-platform.md` — endpoint table vs `PluginResource.java`

**Checkpoint**: Plugin settings page loads in React without 404 on `plugin/main`.

---

## Phase 4: User Story 2 — Audit log search (Priority: P2)

**Goal**: `POST /rest/plugins/audit/private/log/search` with pagination and `plugin_audit_access`.

**Independent Test**: `quickstart.md` §5 — search returns seeded audit rows; denied without permission.

### Tests for User Story 2

- [X] T033 [P] [US2] Add `serverBackendGo/internal/modules/plugins/audit/application/service_test.go` — filter builder, tenant scope
- [X] T034 [P] [US2] Add `serverBackendGo/internal/modules/plugins/audit/adapter/http/handler_test.go` for POST log search

### Implementation for User Story 2

- [X] T035 [US2] Define `serverBackendGo/internal/modules/plugins/audit/port/repository.go`
- [X] T036 [US2] Implement `serverBackendGo/internal/modules/plugins/audit/adapter/persistence/postgres/audit_repo.go` — search + count with customer filter
- [X] T037 [US2] Implement `serverBackendGo/internal/modules/plugins/audit/application/service.go`
- [X] T038 [US2] Create `serverBackendGo/internal/modules/plugins/audit/adapter/http/handler.go` — POST `/private/log/search` on `groups.Plugins.Group("/audit")`
- [X] T039 [US2] Wire `serverBackendGo/internal/modules/plugins/audit/module.go` — replace scaffold; `MODULE_PLUGINS_AUDIT_ENABLED`
- [X] T040 [US2] Create `serverBackendGo/docs/parity/plugins-audit.md` — note servlet filter **partial** per FR-013

**Checkpoint**: Audit search smoke passes for admin JWT.

---

## Phase 5: User Story 3 — Messaging plugin (Priority: P2)

**Goal**: Messaging send/search/delete/purge + public status; agent delivery via `MessageQueue`.

**Independent Test**: `quickstart.md` §6 — send to `hmdm-001`, search history, notifications pull shows delivery path.

### Tests for User Story 3

- [X] T041 [P] [US3] Add `serverBackendGo/internal/modules/plugins/messaging/application/service_test.go` — send targets, plugin disabled guard
- [X] T042 [P] [US3] Add `serverBackendGo/internal/modules/plugins/messaging/adapter/http/handler_test.go` for POST send and search

### Implementation for User Story 3

- [X] T043 [US3] Define `serverBackendGo/internal/modules/plugins/messaging/port/repository.go` — message CRUD + purge
- [X] T044 [US3] Implement `serverBackendGo/internal/modules/plugins/messaging/adapter/persistence/postgres/message_repo.go`
- [X] T045 [US3] Implement `serverBackendGo/internal/modules/plugins/messaging/application/service.go` — use `shared/targets` + `MessageQueue.Enqueue` per `research.md` R3
- [X] T046 [US3] Create `serverBackendGo/internal/modules/plugins/messaging/adapter/http/handler.go` — private search/send/delete/purge on `groups.Plugins.Group("/messaging")`
- [X] T047 [US3] Register GET `/rest/plugins/messaging/public/status/:id/:status` without JWT (public subgroup or dedicated route) per contract
- [X] T048 [US3] Wire `serverBackendGo/internal/modules/plugins/messaging/module.go` — replace scaffold; inject queue + targets; `MODULE_PLUGINS_MESSAGING_ENABLED`
- [X] T049 [US3] Create `serverBackendGo/docs/parity/plugins-messaging.md`

**Checkpoint**: Messaging send enqueues agent notification; history searchable.

---

## Phase 6: User Story 4 — Device info plugin (Priority: P2)

**Goal**: Public dynamic upload; private detail/search/export; settings GET/PUT.

**Independent Test**: `quickstart.md` §7 — PUT public payload, GET private detail for `hmdm-001`.

### Tests for User Story 4

- [X] T050 [P] [US4] Add `serverBackendGo/internal/modules/plugins/deviceinfo/application/service_test.go` — preserve period, plugin disabled check
- [X] T051 [P] [US4] Add `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/http/handler_test.go` for settings GET/PUT

### Implementation for User Story 4

- [X] T052 [US4] Define `serverBackendGo/internal/modules/plugins/deviceinfo/port/repository.go` — settings, dynamic info, device search
- [X] T053 [US4] Implement `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/persistence/postgres/settings_repo.go` and `deviceinfo_repo.go`
- [X] T054 [US4] Implement `serverBackendGo/internal/modules/plugins/deviceinfo/application/service.go` — settings + upload + search/export (export **partial** note if CSV simplified)
- [X] T055 [US4] Create `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/http/settings_handler.go` — `/deviceinfo-plugin-settings/private`
- [X] T056 [US4] Create `serverBackendGo/internal/modules/plugins/deviceinfo/adapter/http/deviceinfo_handler.go` — public PUT + private routes per contract
- [X] T057 [US4] Wire `serverBackendGo/internal/modules/plugins/deviceinfo/module.go` — replace scaffold; `MODULE_PLUGINS_DEVICEINFO_ENABLED`
- [X] T058 [US4] Create `serverBackendGo/docs/parity/plugins-deviceinfo.md`

**Checkpoint**: Device info upload and admin read paths operational.

---

## Phase 7: User Story 5 — Device log plugin (Priority: P2)

**Goal**: Settings/rules CRUD, public rules + upload, private search/export (Postgres).

**Independent Test**: `quickstart.md` §8 — settings GET, log search after optional upload seed.

### Tests for User Story 5

- [X] T059 [P] [US5] Add `serverBackendGo/internal/modules/plugins/devicelog/application/service_test.go` — rule validation, tenant scope
- [X] T060 [P] [US5] Add `serverBackendGo/internal/modules/plugins/devicelog/adapter/http/handler_test.go` for settings GET

### Implementation for User Story 5

- [X] T061 [US5] Define `serverBackendGo/internal/modules/plugins/devicelog/port/repository.go` — settings, rules, log insert/search
- [X] T062 [US5] Implement `serverBackendGo/internal/modules/plugins/devicelog/adapter/persistence/postgres/settings_repo.go`, `rules_repo.go`, `log_repo.go`
- [X] T063 [US5] Implement `serverBackendGo/internal/modules/plugins/devicelog/application/service.go` — synchronous batch insert on upload per `research.md` R9
- [X] T064 [US5] Create `serverBackendGo/internal/modules/plugins/devicelog/adapter/http/settings_handler.go` — `/devicelog-plugin-settings/private` + rule PUT/DELETE
- [X] T065 [US5] Create `serverBackendGo/internal/modules/plugins/devicelog/adapter/http/log_handler.go` — public rules/upload + private search/export
- [X] T066 [US5] Wire `serverBackendGo/internal/modules/plugins/devicelog/module.go` — replace scaffold; `MODULE_PLUGINS_DEVICELOG_ENABLED`
- [X] T067 [US5] Create `serverBackendGo/docs/parity/plugins-devicelog.md`

**Checkpoint**: Devicelog admin search and public upload paths respond correctly.

---

## Phase 8: User Story 6 — Push plugin schedule tasks (Priority: P3)

**Goal**: Complete Phase 7 partial — `searchTasks`, PUT task, DELETE task on existing `plugins/push`.

**Independent Test**: `quickstart.md` §9 — searchTasks returns rows; PUT/DELETE round-trip on `plugin_push_schedule`.

### Tests for User Story 6

- [X] T068 [P] [US6] Add `serverBackendGo/internal/modules/plugins/push/application/schedule_test.go` — device scope resolution, permission names match Java (`plugin_push_delete` for save)
- [X] T069 [P] [US6] Extend `serverBackendGo/internal/modules/plugins/push/adapter/http/handler_test.go` for schedule endpoints

### Implementation for User Story 6

- [X] T070 [US6] Define schedule port methods in `serverBackendGo/internal/modules/plugins/push/port/` or extend existing repository interface
- [X] T071 [US6] Implement `serverBackendGo/internal/modules/plugins/push/adapter/persistence/postgres/schedule_repo.go` — CRUD on `plugin_push_schedule`
- [X] T072 [US6] Extend `serverBackendGo/internal/modules/plugins/push/application/service.go` with `SearchTasks`, `SaveTask`, `DeleteTask`
- [X] T073 [US6] Extend `serverBackendGo/internal/modules/plugins/push/adapter/http/handler.go` — POST `searchTasks`, PUT `task`, DELETE `task/:id` per `contracts/plugins-push-schedule-api.md`
- [X] T074 [US6] Update `serverBackendGo/docs/parity/push.md` — schedule section; cron execution **partial**

**Checkpoint**: Phase 7 push message endpoints still pass regression smoke; schedule CRUD works.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Swagger, docs, migration status, full verification.

- [X] T075 [P] Add Swagger `// @Router` comments on all new handlers; run `cd serverBackendGo && make swagger`
- [X] T076 [P] Update `serverBackendGo/docs/MIGRATION.md` — Phase 8 row **done** with parity links
- [X] T077 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — Phase 8 **منجز**
- [X] T078 Run `cd serverBackendGo && go test ./internal/modules/plugins/... -count=1` and fix failures
- [X] T079 Run full `specs/009-complete-phase8-plugins/quickstart.md` smoke (§3–§9) against `make dev`
- [X] T080 [P] Manual React check: login → Settings → Plugins tab — no errors in Network for `plugin/main`
- [X] T081 [P] Add negative tenant isolation test in `audit` or `messaging` `application/service_test.go` (SC-006)
- [X] T082 Verify `go build ./...` and remove any remaining `module scaffold registered` logs for implemented plugins in `internal/app/modules.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup — **blocks all user stories**
- **US1 (Phase 3)**: After Phase 2 — **MVP** (no dependency on US2–US6)
- **US2–US5 (Phases 4–7)**: After Phase 2; independent of each other (may run in parallel)
- **US6 (Phase 8)**: After Phase 2; extends existing `plugins/push` from Phase 7
- **Polish (Phase 9)**: After desired user stories complete

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 | Phase 2 | Platform catalog only |
| US2 | Phase 2 | Audit tables in 000010 |
| US3 | Phase 2 + Phase 7 queue | Uses `MessageQueue` + `shared/targets` |
| US4 | Phase 2 | Deviceinfo tables |
| US5 | Phase 2 | Devicelog tables |
| US6 | Phase 2 + Phase 7 push | Schedule on existing module |

### Parallel Opportunities

- Phase 1: T002, T003 [P]
- Phase 2: T006–T008, T010–T012, T017–T021 [P] after T005
- After Phase 2: **US2**, **US3**, **US4**, **US5** can run in parallel (different module trees)
- **US1** should complete first for React MVP
- Phase 9: T075–T077, T080–T081 [P]

### Parallel Example: MVP + one extension

```bash
# After Phase 2:
# Track A: T026–T032 (US1 platform) → validate quickstart §4
# Track B (parallel): T033–T040 (US2 audit)
```

### Parallel Example: Agent-facing plugins

```bash
# After Phase 2:
# Developer A: T050–T058 (US4 deviceinfo)
# Developer B: T059–T067 (US5 devicelog)
```

---

## Implementation Strategy

### MVP First (US1 only)

1. Complete Phase 1–2
2. Complete Phase 3 (US1 platform)
3. **STOP and VALIDATE**: `quickstart.md` §4 + React Plugin settings
4. Continue US2 → US3 → US4 → US5 → US6 → Polish

### Incremental Delivery

1. Foundation → US1 (React plugin settings)
2. US2 (audit) + US3 (messaging) — operator tooling
3. US4 + US5 in parallel — agent telemetry/logs
4. US6 — push schedule completion
5. Polish — MIGRATION done, full quickstart

### Suggested MVP Scope

**Minimum**: Phases 1–2 + **US1** + T079 subset (platform smoke only).

Delivers React Plugin settings without Java; defers audit, messaging, deviceinfo, devicelog, schedule to next slices.

---

## Notes

- `PluginMain` private routes need explicit JWT middleware (group is not auto-protected in `httpx/router.go`).
- Java uses `plugin_push_delete` for PUT schedule task — match in Go, do not rename permission.
- Table `pluginsDisabled` uses camelCase column names from legacy liquibase.
- Branch: `009-complete-phase8-plugins`
