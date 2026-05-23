# Implementation Plan: إكمال فجوات قاعدة البيانات Java → Go

**Branch**: `013-complete-database-gaps` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/013-complete-database-gaps/spec.md`  
**Gap source**: [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md)  
**Parallel work**: [`specs/012-finish-java-go-backend/`](../012-finish-java-go-backend/) (REST/API — consumes schema from this plan)

## Summary

Close **PostgreSQL schema gaps** between legacy Java Liquibase and Go `golang-migrate` so React, agents, and Go repositories behave like production Headwind MDM: add **P1** tables `devicestatuses` and `userrolesettings`, **P2** tables/columns (`configurationapplicationparameters`, `usagestats`, `apkhash`, settings extensions), optional **legacy import** of Java `configurations` columns into `settingsjson`, then wire **devices**, **settings**, and **summary** repos to use the new schema.

**Technical approach**: Seven sequential migrations `000011`–`000017`; idempotent SQL; full `down` files; repository JOINs replacing interim `infojson`-only filters; extend stub `UserRoleSettings` in settings module; document parity in root gap tracker.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: `database/sql`, `lib/pq`, golang-migrate, existing module repos (devices, settings, configurations, summary, stats)

**Storage**: PostgreSQL 14; new tables per [data-model.md](./data-model.md); no new stored procedures

**Testing**: `make migrate` on fresh DB; `go test` on `devices`, `settings`, `summary`; [quickstart.md](./quickstart.md) SQL + curl smoke

**Target Platform**: Linux/macOS dev; optional Java dump import path

**Project Type**: Schema + repository layer (backend)

**Performance Goals**: Device search with `devicestatuses` JOIN &lt; 2s p95 at 10k devices (index on `applicationsstatus`); role settings read &lt; 100ms

**Constraints**: Migrations reversible; `IF NOT EXISTS`; lowercase identifiers; minimal seed/backfill in SQL only

**Scale/Scope**: 4 new tables; ~30 new columns across 3 tables; 4 modules touched; 4 contract docs; 7 migration pairs

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Schema per module consumers; `db/migrations` versioned |
| **II. Layered Clean** | ✅ | SQL in `adapter/persistence`; no domain imports SQL |
| **III. API Parity** | ✅ | Repos enable existing REST (012); no path changes in 013 |
| **IV. Testable Delivery** | ✅ | migrate smoke + repo unit tests in quickstart |
| **V. Simplicity** | ✅ | No PL/pgSQL except optional 000017; no 55-table port |
| **VI. Security** | ✅ | Tenant scope unchanged; migrations no secrets |
| **VII. Observability** | ✅ | migrate failures logged; gap doc updated |

**Post-design**: All gates ✅. Cross-cutting migrations justified (single DB, multiple module consumers).

## Project Structure

### Documentation (this feature)

```text
specs/013-complete-database-gaps/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── migrations-schema.md
│   ├── legacy-config-import.md
│   └── repository-integration.md
└── tasks.md             # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000011_devicestatuses_core.up.sql
│   ├── 000011_devicestatuses_core.down.sql
│   ├── 000012_userrolesettings_core.up.sql
│   ├── 000012_userrolesettings_core.down.sql
│   ├── 000013_configuration_application_parameters.up.sql
│   ├── 000013_configuration_application_parameters.down.sql
│   ├── 000014_usagestats_core.up.sql
│   ├── 000014_usagestats_core.down.sql
│   ├── 000015_settings_columns_extend.up.sql
│   ├── 000015_settings_columns_extend.down.sql
│   ├── 000016_applications_columns_extend.up.sql
│   ├── 000016_applications_columns_extend.down.sql
│   ├── 000017_configurations_legacy_import.up.sql
│   └── 000017_configurations_legacy_import.down.sql
├── internal/modules/
│   ├── devices/adapter/persistence/postgres/   # JOIN devicestatuses
│   ├── settings/domain/settings.go             # full UserRoleSettings
│   ├── settings/adapter/persistence/postgres/  # userrolesettings CRUD
│   ├── settings/adapter/http/handler.go        # return column flags
│   ├── summary/adapter/persistence/postgres/   # charts from devicestatuses
│   └── configurations/...                      # optional cap params (P2)
├── docs/
│   └── parity/devices.md, settings.md, summary.md  # schema notes
└── JAVA-GO-DATABASE-GAPS.md                    # ✅ updates FR-012
```

**Structure Decision**: Migrations centralized in `db/migrations/`; module repos updated only where FR-011 requires reads/writes.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — P1 schema: device + role UI (FR-001, FR-002)

1. `000011_devicestatuses_core` per [contracts/migrations-schema.md](./contracts/migrations-schema.md).
2. `000012_userrolesettings_core` with full column set + seed.
3. Wire `device_repo` per [contracts/repository-integration.md](./contracts/repository-integration.md).
4. Implement `settings` repo + handler for `UserRoleSettings`.
5. Update `docs/parity/devices.md`, `settings.md`.

### Phase B — P2 schema: params, stats, columns (FR-003–FR-005)

1. `000013`–`000016` migrations.
2. Optional `configurations` repo for `configurationapplicationparameters`.
3. Confirm `stats` module (012) targets `usagestats` table after `000014`.

### Phase C — Legacy import (FR-008)

1. `000017_configurations_legacy_import` per [contracts/legacy-config-import.md](./contracts/legacy-config-import.md).
2. Document ⊘ on greenfield Go DB in parity doc.

### Phase D — Summary + docs (FR-011, FR-012)

1. `summary` repo SQL uses `devicestatuses`.
2. Refresh [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) §3.2, §4, §7.
3. Note in [`serverBackendGo/docs/MIGRATION.md`](../../serverBackendGo/docs/MIGRATION.md) schema milestone.

### Phase E — Validation (SC-001–SC-006)

1. [quickstart.md](./quickstart.md) full pass.
2. `go test ./...`
3. Rollback one migration pair in dev to verify `down`.

## Complexity Tracking

| Item | Why | Simpler alternative rejected |
|------|-----|------------------------------|
| 7 migration files | Independent rollback slices | Single 2000-line migration — unmaintainable |
| 000017 conditional PL/pgSQL | Java dumps vary by version | Manual script only — not repeatable in CI |

## Dependencies & Coordination

| Dependency | Action |
|------------|--------|
| **012 devices US1** | After `000011`, switch `installationStatus` filter to `devicestatuses` |
| **012 stats** | Requires `000014` before INSERT |
| **012 summary charts** | Requires `000011` for accurate install-by-config |
| **Fresh dev** | `make migrate` runs 001→017 in order |

## Out of Scope (plan confirms spec)

- Plugin optional tables (WiFi/GPS, devicelocations, photo, …) — P3 ⊘
- `trialkey`, temp upload tables
- Stored procedures `mdm_*`
- REST route additions (012)

## Phase 0 & 1 Artifacts

| Artifact | Status |
|----------|--------|
| [research.md](./research.md) | ✅ Complete |
| [data-model.md](./data-model.md) | ✅ Complete |
| [contracts/](./contracts/) | ✅ Complete |
| [quickstart.md](./quickstart.md) | ✅ Complete |

**Next command**: `/speckit-tasks` to generate ordered `tasks.md`.
