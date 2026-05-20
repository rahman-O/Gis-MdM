# Implementation Plan: Phase 8 — Plugins Platform & Extension Modules

**Branch**: `009-complete-phase8-plugins` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/009-complete-phase8-plugins/spec.md`

## Summary

Deliver **Phase 8** of the Go migration: replace scaffolds for **`plugins/platform`**, **`plugins/audit`**, **`plugins/messaging`**, **`plugins/deviceinfo`**, **`plugins/devicelog`**, and complete **`plugins/push`** schedule-task CRUD. Add migration **`000010_plugins_core`** for `plugins`, `pluginsDisabled`, plugin-specific tables, permissions, and dev seeds. Wire **shared device targeting** and **notifications `MessageQueue`** for messaging send. Enable React **Plugin settings** (`pluginService.ts`) and legacy plugin REST consumers without Java.

**Partial**: audit servlet auto-capture; push schedule cron runner; deviceinfo export streaming polish; devicelog async upload pool.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth`, `platform/httpx/response`, Phase 4 `devices` resolution patterns, Phase 7 `notifications/port.MessageQueue`

**Storage**: PostgreSQL `000010_plugins_core.up.sql`; legacy `plugin_*` table names; reuse `plugin_push_*` from `000009`

**Testing**: `go test ./internal/modules/plugins/...`; `application/` unit tests per plugin; optional `adapter/http/handler_test.go`; HTTP smoke in [quickstart.md](./quickstart.md)

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`

**Project Type**: Web service + React (`PluginSettingsPage`) + Android agents (deviceinfo/devicelog public paths)

**Performance Goals**: Plugin list endpoints &lt; 500ms p95 seeded DB; messaging send for 100 devices &lt; 5s batch insert+enqueue

**Constraints**: Tenant scope on all private/plugin routes; `PluginMain` private routes need explicit JWT middleware; public plugin callbacks unauthenticated; Headwind JSON envelope; no fake scaffold 200s

**Scale/Scope**: ~35 REST endpoints across 6 plugin trees; ~100–120 new Go files; 6 parity docs; 1 migration

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Six plugin subtrees; Phase 8 in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | Shared `port` for targets + queue injection; no handler SQL |
| **III. API Parity** | ✅ | `contracts/*.md` + `docs/parity/plugins-*.md` |
| **IV. Testable Delivery** | ✅ | Unit + quickstart + build gate |
| **V. Simplicity** | ✅ | Reuse queue/targets; defer audit filter & cron |
| **VI. Security** | ✅ | Per-route permissions; tenant isolation tests |
| **VII. Observability** | ✅ | Legacy error keys; `MODULE_PLUGINS_*` in `.env.example` |

**Post-design**: All gates ✅. Partial items documented in parity—not hidden scaffolds.

## Project Structure

### Documentation (this feature)

```text
specs/009-complete-phase8-plugins/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── plugins-platform-api.md
│   ├── plugins-audit-api.md
│   ├── plugins-messaging-api.md
│   ├── plugins-deviceinfo-api.md
│   ├── plugins-devicelog-api.md
│   └── plugins-push-schedule-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000010_plugins_core.up.sql
│   └── 000010_plugins_core.down.sql
├── docs/parity/
│   ├── plugins-platform.md
│   ├── plugins-audit.md
│   ├── plugins-messaging.md
│   ├── plugins-deviceinfo.md
│   ├── plugins-devicelog.md
│   └── push.md                 # extend schedule section
├── internal/config/config.go   # ENABLED_PLUGINS, MODULE_PLUGINS_*
├── internal/modules/plugins/
│   ├── platform/               # REPLACE scaffold — PluginMain routes
│   ├── audit/
│   ├── messaging/
│   ├── deviceinfo/
│   ├── devicelog/
│   ├── push/                   # ADD schedule handlers to existing Phase 7 tree
│   └── shared/
│       ├── targets/            # port + postgres (from push targets_repo)
│       └── status/             # PluginStatusCache port
└── internal/app/modules.go     # Wire MessageQueue + module flags
```

**Structure Decision**: **platform** uses `groups.PluginMain` with JWT subgroup for `/private`. Other plugins use `groups.Plugins.Group("<name>")`. **messaging** depends on **notifications** queue via interface wired in `app`. **shared/targets** avoids duplicating `targets_repo.go`.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Migration `000010_plugins_core`

1. `CREATE TABLE plugins`, `pluginsDisabled` if not exists.
2. Plugin-specific tables: `plugin_audit_log`, `plugin_messaging_messages`, deviceinfo tables, devicelog postgres tables.
3. Seed `plugins` rows (audit, push, messaging, deviceinfo, devicelog).
4. Seed permissions + role 2 mappings; sample audit rows for smoke.

### Phase B — Config & shared infrastructure

1. `ENABLED_PLUGINS`, `MODULE_PLUGINS_*` flags in `config.go`, `.env.example`, `scripts/dev.sh`.
2. `plugins/shared/status` — disabled plugin cache.
3. `plugins/shared/targets` — extract from push `targets_repo`.

### Phase C — Platform (P1)

| Endpoint | Notes |
|----------|-------|
| GET `/private/available`, `/private/active` | Filter by env + `pluginsDisabled` |
| GET `/public/registered` | No auth |
| POST `/private/disabled` | `plugins_customer_access_management` |

Register JWT middleware on `PluginMain` private subgroup only.

### Phase D — Audit (P2)

| Endpoint | Notes |
|----------|-------|
| POST `/private/log/search` | `plugin_audit_access`, paginated |

### Phase E — Messaging (P2)

| Endpoint | Notes |
|----------|-------|
| POST `/private/search`, `/private/send` | Queue + DB |
| DELETE `/{id}`, GET `/private/purge/{days}` | |
| GET `/public/status/{id}/{status}` | On `groups.Public` or plugins public subgroup |

### Phase F — Device info (P2)

Settings + deviceinfo resources per [contracts/plugins-deviceinfo-api.md](./contracts/plugins-deviceinfo-api.md).

### Phase G — Device log (P2)

Settings + log resources; Postgres insert on upload.

### Phase H — Push schedule completion (P3)

Add `searchTasks`, PUT `task`, DELETE `task/{id}` to existing `plugins/push` module.

### Phase I — Docs, Swagger, verification

1. Parity docs ×5 + update `push.md`.
2. `make swagger` for new `@Router` comments.
3. Update `MIGRATION.md` / `NEXT_STEPS.md` Phase 8 → **done**.
4. Run `quickstart.md` + `go test ./internal/modules/plugins/...`.

## Complexity Tracking

| Item | Status | Notes |
|------|--------|-------|
| Shared targets package | Justified | Constitution V — 3 senders |
| Plugin status cache | Justified | Java parity for disabled checks |
| No audit servlet filter | Partial documented | FR-013 |
