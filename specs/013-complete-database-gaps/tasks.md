---
description: "Task list for 013 — Complete database schema gaps (Java → Go)"
---

# Tasks: إكمال فجوات قاعدة البيانات Java → Go

**Input**: `specs/013-complete-database-gaps/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Migrations `000001`–`000010` + `000008_devices_search_extras` applied; [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md); optional coordination with [`specs/012-finish-java-go-backend/`](../012-finish-java-go-backend/) (REST consumes schema)

**Tests**: Unit tests for SQL filter builders and settings repo per constitution IV; migration smoke via `quickstart.md`.

**Organization**: Tasks grouped by user story (US1–US5) for independent delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Migrations: `serverBackendGo/db/migrations/`
- Modules: `serverBackendGo/internal/modules/<name>/`
- Parity: `serverBackendGo/docs/parity/`
- Contracts: `specs/013-complete-database-gaps/contracts/`
- Trackers: `JAVA-GO-DATABASE-GAPS.md`, `serverBackendGo/docs/MIGRATION.md`

---

## Phase 1: Setup

**Purpose**: Confirm gap analysis, Java Liquibase references, and migration baseline.

- [x] T001 Verify `specs/013-complete-database-gaps/spec.md` against `JAVA-GO-DATABASE-GAPS.md` §3–§7 and plan migration sequence `000011`–`000017`
- [x] T002 [P] Review Java `db.changelog.xml` `deviceStatuses`, `userRoleSettings`, `configurationApplicationParameters`, `usageStats` changelogs against `specs/013-complete-database-gaps/data-model.md`
- [x] T003 [P] Review `specs/013-complete-database-gaps/contracts/migrations-schema.md` and `repository-integration.md` for SQL/repository alignment
- [x] T004 List current `serverBackendGo/db/migrations/*.up.sql` and confirm next free prefix is `000011` (note duplicate `000008_*` names)
- [x] T005 Run `cd serverBackendGo && go build ./... && go test ./...` and record baseline in `specs/013-complete-database-gaps/quickstart.md` notes if needed

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Ensure existing schema is migratable before new tables; document 012 coordination.

**⚠️ CRITICAL**: No US1+ work until fresh `make migrate` through `000010` succeeds.

- [x] T006 Run `cd serverBackendGo && ./scripts/db-up.sh && make migrate` on empty volume — confirm through `000010_plugins_core` + `000008_devices_search_extras` without error
- [x] T007 [P] Verify `devices` and `settings` modules build against current schema (`go test ./internal/modules/devices/... ./internal/modules/settings/...`)
- [x] T008 Document in `specs/013-complete-database-gaps/plan.md` or `quickstart.md` — **012** `installationStatus` filter switches after T015–T018 (post `000011`)
- [x] T009 [P] Add checklist row to `serverBackendGo/docs/MIGRATION.md` — schema gap closure (013) **in progress**

**Checkpoint**: Baseline DB green — P1 migrations may start.

---

## Phase 3: User Story 1 — حالة تثبيت الأجهزة (Priority: P1) 🎯 MVP

**Goal**: Table `devicestatuses` + device search/summary SQL uses it for `installationStatus` (FR-001, FR-011 partial).

**Independent Test**: `quickstart.md` §2 — different `applicationsstatus` on two devices → search filter returns correct subset.

### Tests for User Story 1

- [x] T010 [P] [US1] Add `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_filters_test.go` — `installationStatus` builds `devicestatuses` JOIN clause
- [x] T011 [P] [US1] Extend `device_repo_test.go` or filter tests for `INSTALLATIONS`/`FILES` sort keys using `devicestatuses`

### Implementation for User Story 1

- [x] T012 [US1] Create `serverBackendGo/db/migrations/000011_devicestatuses_core.up.sql` per `contracts/migrations-schema.md` (table + index + backfill)
- [x] T013 [US1] Create `serverBackendGo/db/migrations/000011_devicestatuses_core.down.sql` — `DROP TABLE devicestatuses`
- [x] T014 [US1] Update `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_filters.go` — filter `installationStatus` on `devicestatuses.applicationsstatus`; remove `infojson` applications-only path
- [x] T015 [US1] Update `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_filters.go` — `orderExprInner` cases `INSTALLATIONS` and `FILES` via `devicestatuses`
- [x] T016 [US1] Ensure `Count` and `Search` in `device_repo.go` include `LEFT JOIN devicestatuses` when filter or sort requires it
- [x] T017 [P] [US1] Add optional `serverBackendGo/internal/modules/devices/domain/device.go` `DeviceStatus` struct if used by port layer
- [x] T018 [US1] Run `make migrate` and execute `quickstart.md` §2 device status filter smoke
- [x] T019 [P] [US1] Update `serverBackendGo/docs/parity/devices.md` — `devicestatuses` and `installationStatus` **Done**
- [x] T020 [P] [US1] Update `JAVA-GO-DATABASE-GAPS.md` §3.2 `devicestatuses` row → ✅

**Checkpoint**: Device install-status filtering works from DB table, not `infojson` workaround.

---

## Phase 4: User Story 2 — أعمدة قائمة الأجهزة (Priority: P1)

**Goal**: Table `userrolesettings` + settings API returns all `columnDisplayed*` flags (FR-002).

**Independent Test**: `quickstart.md` §3 — `GET /rest/private/settings/userRole/2` returns boolean column flags; PUT persists.

### Tests for User Story 2

- [x] T021 [P] [US2] Add `serverBackendGo/internal/modules/settings/application/service_test.go` — `GetUserRoleSettings` returns defaults when no row
- [x] T022 [P] [US2] Add `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo_test.go` — upsert `(roleid, customerid)` round-trip (sqlmock or integration)

### Implementation for User Story 2

- [x] T023 [US2] Create `serverBackendGo/db/migrations/000012_userrolesettings_core.up.sql` — full `columndisplayed*` columns per `data-model.md` + UNIQUE + seed roles 1–3
- [x] T024 [US2] Create `serverBackendGo/db/migrations/000012_userrolesettings_core.down.sql`
- [x] T025 [US2] Extend `serverBackendGo/internal/modules/settings/domain/settings.go` `UserRoleSettings` with all `columnDisplayed*` JSON fields
- [x] T026 [US2] Add `GetUserRoleSettings` / `SaveUserRoleSettings` to `serverBackendGo/internal/modules/settings/port/repository.go`
- [x] T027 [US2] Implement methods in `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo.go`
- [x] T028 [US2] Wire `serverBackendGo/internal/modules/settings/application/service.go` for get/save user role settings
- [x] T029 [US2] Update `serverBackendGo/internal/modules/settings/adapter/http/handler.go` `GetUserRole` and `SaveUserRolesCommon` to use service (not stub `{roleId}` only)
- [x] T030 [US2] Run `quickstart.md` §3 user role column settings smoke
- [x] T031 [P] [US2] Update `serverBackendGo/docs/parity/settings.md` — `userrolesettings` **Done**
- [x] T032 [P] [US2] Update `JAVA-GO-DATABASE-GAPS.md` §3.2 and §4.4 `userrolesettings` / `columnDisplayed*` → ✅

**Checkpoint**: React device list column prefs can load/save via Go API.

---

## Phase 5: User Story 3 — تكامل تطبيقات التكوين والإحصائيات (Priority: P2)

**Goal**: Migrations `000013`–`000016` for CAP, usagestats, settings columns, app version columns (FR-003–FR-005).

**Independent Test**: `quickstart.md` §4 — `\d configurationapplicationparameters`, `\d usagestats`, `apkhash` column exists.

### Implementation for User Story 3

- [x] T033 [P] [US3] Create `serverBackendGo/db/migrations/000013_configuration_application_parameters.up.sql` and `.down.sql` per `contracts/migrations-schema.md`
- [x] T034 [P] [US3] Create `serverBackendGo/db/migrations/000014_usagestats_core.up.sql` and `.down.sql`
- [x] T035 [P] [US3] Create `serverBackendGo/db/migrations/000015_settings_columns_extend.up.sql` and `.down.sql` per `data-model.md`
- [x] T036 [P] [US3] Create `serverBackendGo/db/migrations/000016_applications_columns_extend.up.sql` and `.down.sql` (`apkhash`, `remove`, `longtap`)
- [x] T037 [US3] Extend `serverBackendGo/internal/modules/settings/adapter/persistence/postgres/settings_repo.go` SELECT/UPDATE for new `settings` columns (`newdevicegroupid`, `phonenumberformat`, custom property fields)
- [x] T038 [P] [US3] Extend `serverBackendGo/internal/modules/settings/domain/settings.go` with new tenant settings fields
- [x] T039 [P] [US3] Add `configurationapplicationparameters` read/write in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` when `skipVersionCheck` present in payload
- [x] T040 [P] [US3] Add `apkhash` column scan/update in `serverBackendGo/internal/modules/applications/adapter/persistence/postgres/` (application versions repo)
- [x] T041 [US3] Add `remove` / `longtap` to `configurationapplications` persistence in `config_repo.go` or dedicated repo file
- [x] T042 [US3] Document in `serverBackendGo/docs/parity/configurations.md` — CAP table and extended columns
- [x] T043 [P] [US3] Add note in `specs/012-finish-java-go-backend/quickstart.md` or `JAVA-GO-BACKEND-GAPS.md` — `stats` module requires `000014` before INSERT
- [x] T044 [US3] Run `make migrate` and `quickstart.md` §4 schema checks

**Checkpoint**: P2 tables/columns exist; repos ready for 012 stats/config work.

---

## Phase 6: User Story 4 — ترحيل من قاعدة Java (Priority: P2)

**Goal**: Optional migration `000017` copies legacy `configurations` columns → `settingsjson` (FR-008).

**Independent Test**: On Java-column DB clone, `SELECT settingsjson FROM configurations LIMIT 3` shows merged policy keys; on greenfield Go DB migration is no-op.

### Implementation for User Story 4

- [x] T045 [US4] Create `serverBackendGo/db/migrations/000017_configurations_legacy_import.up.sql` per `contracts/legacy-config-import.md` (conditional `information_schema` + `jsonb` merge)
- [x] T046 [US4] Create `serverBackendGo/db/migrations/000017_configurations_legacy_import.down.sql` (document no-op or documented non-reversible)
- [x] T047 [US4] Add `serverBackendGo/docs/parity/configurations.md` section — legacy import path and greenfield ⊘
- [x] T048 [P] [US4] Update `JAVA-GO-DATABASE-GAPS.md` §4.3 — `settingsjson` import strategy **Done** with pointer to `000017`
- [x] T049 [US4] Manual verification checklist in `specs/013-complete-database-gaps/quickstart.md` §5 (≥3 configs, SC-005)

**Checkpoint**: Java dump migration path documented and executable.

---

## Phase 7: User Story 5 — Plugins اختيارية (Priority: P3)

**Goal**: Explicitly defer optional plugin tables; no schema unless future spec (FR-010).

**Independent Test**: Review docs — §3.3 plugin tables marked ⊘ with no migrations `000018+` for WiFi/GPS/devicelocations.

### Implementation for User Story 5

- [x] T050 [US5] Update `JAVA-GO-DATABASE-GAPS.md` §3.3 — mark plugin-only tables as **⊘ deferred** (reference 013 Out of Scope)
- [x] T051 [P] [US5] Add `specs/013-complete-database-gaps/research.md` note R7 — confirm no `000018` in v1
- [x] T052 [US5] Add one-line scope guard in `serverBackendGo/docs/MIGRATION.md` — optional plugin schema per-plugin specs only

**Checkpoint**: P3 documented; no accidental scope creep.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Summary charts, gap tracker closure, full validation (FR-011, FR-012, SC-001–SC-006).

- [x] T053 Update `serverBackendGo/internal/modules/summary/adapter/persistence/postgres/` queries to use `devicestatuses` for install-by-config charts per `contracts/repository-integration.md`
- [x] T054 [P] Update `serverBackendGo/docs/parity/summary.md` — remove simplified chart caveat when `devicestatuses` wired
- [x] T055 [P] Refresh `JAVA-GO-DATABASE-GAPS.md` §7 migration map — mark `000011`–`000017` ✅ with dates
- [x] T056 Update `serverBackendGo/docs/MIGRATION.md` — schema gap milestone (013) and dependency on 012 for REST
- [x] T057 Run full `specs/013-complete-database-gaps/quickstart.md` (§1–§7)
- [x] T058 Run `cd serverBackendGo && go test ./...`
- [x] T059 [P] Verify `000011_devicestatuses_core.down.sql` on dev DB (rollback one version) then re-migrate
- [x] T060 Confirm every migration `000011`–`000017` has paired `.down.sql` (SC-006)

**Checkpoint**: P1+P2 schema complete; trackers updated; tests green.

---

## Dependencies & Execution Order

### Phase Dependencies

| Phase | Depends on | Blocks |
|-------|------------|--------|
| 1 Setup | — | Phase 2 |
| 2 Foundational | Phase 1 | US1–US5 |
| 3 US1 (P1) | Phase 2 | US polish summary (partial) |
| 4 US2 (P1) | Phase 2 | — (parallel with US1 after Phase 2) |
| 5 US3 (P2) | Phase 2 | 012 stats INSERT |
| 6 US4 (P2) | US3 `000016` recommended | — |
| 7 US5 (P3) | Phase 2 | — (docs only) |
| 8 Polish | US1, US3 minimum | — |

### User Story Dependencies

- **US1** and **US2** can run **in parallel** after Phase 2 (different migrations `000011` vs `000012`).
- **US3** migrations `000013`–`000016` parallelizable as four file pairs.
- **US4** should run after core tables exist (`000007` configurations); best after US3.
- **US5** independent (documentation).
- **012 coordination**: device `installationStatus` repo change (T014–T016) should merge with or follow 012 US1 device work.

### Within Each User Story

1. Migration `.up.sql` + `.down.sql`
2. Domain/port extensions
3. Repository SQL
4. Application/handler wire
5. Parity + gap doc
6. Quickstart smoke

### Parallel Opportunities

- **Phase 1**: T002, T003 parallel
- **Phase 2**: T007, T009 parallel
- **US1 tests**: T010, T011 parallel
- **US2 tests**: T021, T022 parallel
- **US3 migrations**: T033–T036 all [P] parallel (different files)
- **US3 repos**: T037–T041 parallel after migrations applied
- **After Phase 2**: US1 + US2 in parallel on separate branches/files

---

## Parallel Example: User Story 1 + User Story 2

```bash
# Developer A — devicestatuses
Task T012–T016 (migrations + device_filters + device_repo)

# Developer B — userrolesettings  
Task T023–T029 (migrations + settings domain/repo/handler)

# Shared: run make migrate once both migration files merged (order 000011 then 000012)
```

---

## Parallel Example: User Story 3 migrations

```bash
Task T033  # 000013 CAP
Task T034  # 000014 usagestats
Task T035  # 000015 settings columns
Task T036  # 000016 applications columns
# Then: make migrate && T037–T041 repository updates
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1 + Phase 2
2. Complete Phase 3 (US1): `000011` + device repo JOIN
3. **STOP and VALIDATE** `quickstart.md` §2
4. Unblocks accurate `installationStatus` for 012 devices work

### Recommended order (single developer)

1. Setup → Foundational
2. US1 → US2 (P1 schema)
3. US3 (P2 columns/tables)
4. Polish T053–T054 (summary)
5. US4 (only if Java dump import needed)
6. US5 + Phase 8 validation

### Incremental Delivery

| Increment | Delivers |
|-----------|----------|
| US1 | Install status filtering + summary data source |
| US2 | Device table column preferences |
| US3 | Stats table + config params + APK hash |
| US4 | Legacy DB import path |
| Polish | Docs + full regression |

---

## Notes

- Do not renumber existing `000008_*` migrations; start at `000011`.
- Keep SQL seeds idempotent; business recalc of `devicestatuses` from agent sync is follow-up (not blocking schema).
- Coordinate with `012-finish-java-go-backend` to avoid conflicting edits to `device_filters.go`.
- Commit after each migration pair or story checkpoint.
