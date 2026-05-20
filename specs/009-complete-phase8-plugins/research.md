# Research: Phase 8 — Plugins Platform & Extension Modules

**Branch**: `009-complete-phase8-plugins` | **Date**: 2026-05-21

## R1 — Plugin platform routing and auth

**Decision**: Register `plugins/platform` on `groups.PluginMain` (`/rest/plugin/main`). Apply JWT + `RequireAuth` + `EnrichPrincipal` only to `/private/*` sub-routes; expose `/public/registered` without auth (matches Java `getRegisteredPlugins`).

**Rationale**: `httpx.BuildRouteGroups` creates `PluginMain` without global JWT (unlike `/rest/plugins`). React `pluginService.ts` uses Bearer on private paths only.

**Alternatives considered**:
- Move platform to `/rest/private/plugin/main` — rejected (breaks frontend paths).
- JWT on entire `PluginMain` group — rejected (blocks public registered).

## R2 — Build-enabled plugin catalog

**Decision**: Config var `ENABLED_PLUGINS` (comma-separated identifiers, default `audit,push,messaging,deviceinfo,devicelog`) filters rows from `plugins` table, mirroring Java `PluginList.isPluginEnabled`.

**Rationale**: Java compile-time classpath controls which plugins load; Go has no classpath—env list is the operational equivalent.

**Alternatives considered**:
- Hard-code in code only — rejected (operators cannot disable plugins without rebuild).
- Read from DB flag only — insufficient (Java still filters by build).

## R3 — Device message delivery for messaging plugin

**Decision**: Messaging `send` persists `plugin_messaging_messages` **and** enqueues agent delivery via `notifications/port.MessageQueue` (same as Phase 7 private push). Message type/payload mapped to queue fields per Java `PushService` usage.

**Rationale**: Java `MessagingResource` injects `PushService`; duplicating queue SQL in messaging violates constitution reuse.

**Alternatives considered**:
- Messaging-only DB without agent delivery — rejected (agents would never see messages).
- Direct HTTP to devices — not in legacy model.

## R4 — Target resolution for plugin send endpoints

**Decision**: Extract or share device/group/broadcast resolution from `plugins/push` / `push/adapter/persistence/postgres/targets_repo.go` into `internal/modules/plugins/shared/targets` **port** consumed by messaging and push plugin send.

**Rationale**: Identical SQL and permission semantics; avoids third copy in messaging.

**Alternatives considered**:
- Duplicate targets repo per plugin — rejected (constitution V / maintenance).

## R5 — Database bootstrap for greenfield Go dev

**Decision**: Migration `000010_plugins_core` creates `plugins`, `pluginsDisabled`, all `plugin_*` tables missing from Go migrations 001–009, seeds plugin catalog rows + permissions (`plugins_customer_access_management`, `plugin_audit_access`, `plugin_messaging_*`, `plugin_deviceinfo_access`, `plugin_devicelog_access`), and sample audit/messaging rows for smoke.

**Rationale**: Go migrations never created `plugins` table; fresh `make migrate` must support Phase 8 smoke without Java liquibase.

**Alternatives considered**:
- Require Java DB import only — rejected (violates SC-003 quickstart on fresh DB).

## R6 — Plugin disabled cache

**Decision**: In-process map `customerID → disabled plugin IDs`, invalidated on `POST /private/disabled` (Java `PluginStatusCache` equivalent). Checked by messaging/deviceinfo/devicelog before mutating.

**Rationale**: Java plugins call `pluginStatusCache.isPluginDisabled`; behavior parity for send/upload guards.

**Alternatives considered**:
- Query `pluginsDisabled` on every request — acceptable fallback if cache omitted; cache preferred for parity.

## R7 — Audit servlet filter (deferred)

**Decision**: **Partial** — implement only `POST /private/log/search`. Do not port `AuditFilter` servlet wrapper in Phase 8.

**Rationale**: Spec FR-013; filter is cross-cutting and not required for React plugin settings MVP.

**Alternatives considered**:
- Full filter port — large scope, no React consumer in Go phase.

## R8 — Deviceinfo export format

**Decision**: `POST .../private/export` returns downloadable CSV stream (or JSON fallback documented in parity) matching Java `DeviceInfoExportService` column set; simplify streaming if needed with **partial** note.

**Rationale**: Export is admin-facing; core search/detail/settings required first.

## R9 — Devicelog upload concurrency

**Decision**: Synchronous insert in handler for Phase 8 (batch insert in transaction); optional background goroutine pool **deferred** unless perf smoke fails.

**Rationale**: Java uses `ExecutorService` with 5 threads; Go can start synchronous for parity correctness, optimize later.

## R10 — Push schedule task permissions

**Decision**: Match Java exactly: `searchTasks` authenticated; `PUT /private/task` and `DELETE /private/task/{id}` require `plugin_push_delete` (Java source uses delete permission for save—document in parity, do not “fix”).

**Rationale**: API parity constitution III.

## R11 — Testing strategy

**Decision**:
- Unit tests: plugin list filtering, disabled merge, audit/messaging filter builders, schedule validation.
- HTTP: `handler_test.go` for platform + one plugin (pattern from Phase 5–7).
- `quickstart.md`: curl smoke for all families.

**Rationale**: Constitution IV; user requested best practices with testing.

## R12 — Feature flags

**Decision**: Add `MODULE_PLUGINS_PLATFORM_ENABLED`, `MODULE_PLUGINS_AUDIT_ENABLED`, `MODULE_PLUGINS_MESSAGING_ENABLED`, `MODULE_PLUGINS_DEVICEINFO_ENABLED`, `MODULE_PLUGINS_DEVICELOG_ENABLED` (push plugin flag inherits Phase 7 module registration). Master `MODULE_PLUGINS_ENABLED=true` enables all.

**Rationale**: Aligns with constitution VII module toggles.
