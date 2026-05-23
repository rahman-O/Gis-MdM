# Implementation Plan: Configuration–Device Sync & Admin UX

**Branch**: `016-config-sync-ux` | **Date**: 2026-05-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/016-config-sync-ux/spec.md`

## Summary

Administrators save **Configurations** but enrolled **Android agents** often do not receive the full policy (kiosk, restrictions, new apps) because Go **`SyncResponse` is incomplete** vs Java, device app-settings upsert is weak, and push notify on save may be a no-op. The React editor mixes tab components with a duplicate inline MDM block and excess help copy.

**Approach**: (A) map full configuration → sync payload; (B) enforce configuration-level locks + readonly app settings; (C) wire push on save; (D) refactor configuration editor tabs and `policyLocks` UX; (E) golden/contract tests + parity docs.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 (`frontend`)

**Primary Dependencies**: Gin, `lib/pq`, existing `configurations`, `sync`, `devices`, `applications`, `push` modules; React + existing shadcn UI tabs

**Storage**: PostgreSQL (`configurations.settingsjson`, link tables); `FILES_DIRECTORY` for app/file URLs in sync

**Testing**: `go test ./internal/modules/configurations/... ./internal/modules/sync/...`; frontend unit tests for normalize/locks; [quickstart.md](./quickstart.md) device UAT

**Target Platform**: Go API `:8080`; Vite admin; Headwind launcher agent

**Project Type**: Web application (backend + frontend)

**Performance Goals**: Configuration save &lt; 2s p95; sync build &lt; 3s p95 on typical tenant configuration (&lt; 50 apps)

**Constraints**: REST path parity; Headwind admin envelope on private routes; layered modules per constitution; `policyLocks` v1 in JSON only

**Scale/Scope**: ~25–40 Go files, ~8–12 frontend files; no new module; optional migration deferred

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Work in `configurations`, `sync`, `devices` (notify), `push`; editor in `frontend/features/configurations` |
| **II. Layered Clean** | ✅ | `SyncConfigurationMapper` in `sync/application`; locks in `configurations/domain` |
| **III. API Parity** | ✅ | Contracts document Java field parity; paths unchanged |
| **IV. Testable Delivery** | ✅ | Golden sync tests + application unit tests + quickstart |
| **V. Simplicity** | ✅ | JSON `policyLocks`; no new table v1 |
| **VI. Security** | ✅ | Tenant-scoped repos; sync signature unchanged |
| **VII. Observability** | ✅ | Log skipped readonly overrides; push errors non-fatal |

**Post-design**: All gates ✅. Research resolved all unknowns ([research.md](./research.md)).

## Project Structure

### Documentation (this feature)

```text
specs/016-config-sync-ux/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── configurations-api.md
│   ├── sync-configuration-payload.md
│   └── configuration-editor-ux.md
└── tasks.md              # (/speckit-tasks — not created here)
```

### Source Code

```text
serverBackendGo/
├── internal/modules/configurations/
│   ├── domain/
│   │   ├── configuration.go
│   │   ├── configuration_json.go      # policyLocks merge
│   │   └── policy_locks.go            # NEW allowlist + helpers
│   ├── application/service.go         # save locks; push notify
│   └── adapter/persistence/postgres/config_repo.go
├── internal/modules/sync/
│   ├── domain/sync.go                 # extend SyncResponse fields
│   ├── application/
│   │   ├── sync_configuration_mapper.go  # NEW Java parity mapper
│   │   └── service.go                 # readonly enforcement on POST settings
│   └── adapter/persistence/postgres/device_sync_repo.go
├── internal/modules/configurations/module.go   # wire real PushNotifier
└── docs/parity/
    ├── configurations.md
    └── sync.md

frontend/src/features/configurations/
├── ConfigurationEditorPage.tsx        # tabs only; remove duplicate MDM
├── ConfigurationMdmTab.tsx            # NEW (extract from page)
├── ConfigurationRestrictionsTab.tsx   # NEW
├── configurationNormalize.ts          # policyLocks round-trip
├── configurationService.ts
└── types.ts                           # policyLocks type

خطط مستقبليه/
└── Configurations.md                    # backlog (already created)
```

**Structure Decision**: Cross-cutting sync mapper stays inside `sync` module; policy locks owned by `configurations` domain. Frontend tab split matches [configuration-editor-ux.md](./contracts/configuration-editor-ux.md).

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Sync payload parity (P1) — FR-001, FR-007, FR-008

1. Audit Java `SyncResponse` + assembly code vs Go `BuildSyncResponse`.
2. Extend `domain.SyncResponse` with MDM fields (kiosk*, connectivity, `restrictions`, update modes, …).
3. Implement `sync/application/sync_configuration_mapper.go` loading configuration row + unmarshaling `settingsjson`.
4. Merge `configurationapplicationsettings` + `deviceapplicationsettings` with readonly precedence.
5. Golden test: DB fixture → JSON compared to expected subset.
6. Update `docs/parity/sync.md`.

**Java reference**: `SyncResponse.java`, sync assembly in `ConfigurationService` / `UnsecureDAO`.

### Phase B — Configuration locks (P1) — FR-002, FR-003, FR-004

1. Add `policyLocks` to domain (`PolicyLocks` map, allowlisted keys).
2. GET/PUT configurations: round-trip `policyLocks` in API JSON.
3. On `POST /sync/applicationSettings`, skip updates for readonly config defaults; upsert others.
4. When building sync, apply locked scalar policy from configuration (ignore device-side drift for locked keys).
5. Unit tests: readonly + policyLocks scenarios.

### Phase C — Push on save (P1) — FR-001

1. Wire `configurations/module.go` to `push` module notifier (same pattern as `devices`).
2. Verify `configUpdated` message reaches notification channel for devices on configuration.

### Phase D — Configuration editor UX (P1) — FR-005, FR-006

1. Extract `ConfigurationMdmTab`, `ConfigurationRestrictionsTab`.
2. Remove duplicate MDM block and trim descriptions in `ConfigurationEditorPage`.
3. Lock toggles bound to `policyLocks`; save validation with tab/field errors.
4. Align tab order with [configuration-editor-ux.md](./contracts/configuration-editor-ux.md).

### Phase E — Save round-trip hardening (P2) — FR-009, FR-010

1. Ensure `Save` returns `ConfigurationResponseMap` with all MDM fields (extend if gaps remain).
2. Log/warn when application version missing URL (sync/build).
3. Frontend: reload after save; verify applications catalog linkage.

### Phase F — Verification

1. Run [quickstart.md](./quickstart.md).
2. Manual Android: change kiosk + app list → sync → confirm on device.
3. Update `JAVA-GO-MIGRATION-STATUS.md` row for configurations/sync if applicable.

## Complexity Tracking

> No constitution violations requiring justification.

| Item | Notes |
|------|-------|
| Large `SyncResponse` struct | Required for Java agent parity; alternative JSON blob rejected (agent expects typed fields) |

## Phase 0 & Phase 1 Outputs

- ✅ [research.md](./research.md) — decisions R1–R7  
- ✅ [data-model.md](./data-model.md)  
- ✅ [contracts/](./contracts/) — configurations, sync payload, editor UX  
- ✅ [quickstart.md](./quickstart.md)  

**Next command**: `/speckit-tasks` to generate `tasks.md`.
