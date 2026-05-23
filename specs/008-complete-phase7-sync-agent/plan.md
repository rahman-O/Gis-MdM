# Implementation Plan: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Branch**: `008-complete-phase7-sync-agent` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/008-complete-phase7-sync-agent/spec.md`

## Summary

Deliver **Phase 7** of the Go migration: replace scaffolds for **`sync`**, **`notifications`**, **`push`**,
**`updates`**, **`qrcode`**, and implement **`plugins/push`** with full layered modules. Add migration
`000009_agent_push_notifications` for `pushmessages`, `pendingpushes`, `plugin_push_messages`,
`plugin_push_schedule`, and permissions (`push_api`, `plugin_push_send`, `plugin_push_delete`). Introduce
**shared push queue** persistence (notifications-owned port used by push modules) and **enrollment crypto**
(SHA1 signatures) in `internal/shared/crypto`. Android agents must enroll/sync and receive push commands
via HTTP; React must send push, show QR enrollment, and check updates without Java. **Partial** items:
`SyncResponseHook`, `EventService` side effects, update stats POST, plugin schedule cron execution, FCM.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth`, `platform/httpx/response`, `platform/storage`,
`internal/shared/crypto`, Phase 4 `devices`, Phase 5 `configurations`/`applications`, Phase 6 `files`;
QR: `github.com/skip2/go-qrcode`

**Storage**: PostgreSQL `000009_agent_push_notifications.up.sql`; existing device/configuration schema;
HTTP download for updates manifest/APKs

**Testing**: `go test` on `sync/application`, `notifications/application`, `push/application`, `shared/crypto`;
HTTP smoke in `quickstart.md`; optional Android agent verify manual

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`

**Project Type**: Web service + React (`pushService.ts`, `updatesService.ts`, `EnrollmentQrPage`) + Android agents

**Performance Goals**: Sync response &lt; 3s p95 seeded DB (SC-001); long-poll returns within configured timeout;
push queue insert for 1k devices &lt; 10s batch (best-effort, Java parity)

**Constraints**: Public routes unauthenticated but signature-gated when `SECURE_ENROLLMENT`; tenant scope on
private/plugin routes; no FCM; singular `/rest/notification/polling` path; Headwind JSON envelope on REST JSON

**Scale/Scope**: ~20 agent/admin endpoints across 6 module trees; ~70–90 Go files; parity docs ×5

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Six bounded contexts; Phase 7 in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | Shared queue via `port` only; no handler SQL |
| **III. API Parity** | ✅ | `contracts/*.md` + `docs/parity/*.md` |
| **IV. Testable Delivery** | ✅ | Unit + quickstart smoke |
| **V. Simplicity** | ✅ | One queue repo; hooks/events/stats partial |
| **VI. Security** | ✅ | Signatures, `push_api`, super-admin updates guard |
| **VII. Observability** | ✅ | Legacy error keys; `MODULE_*` flags in `.env.example` |

**Post-design**: All gates ✅. Partial features documented—not scaffold fake 200s.

## Project Structure

### Documentation (this feature)

```text
specs/008-complete-phase7-sync-agent/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── sync-api.md
│   ├── notifications-api.md
│   ├── push-api.md
│   ├── updates-api.md
│   └── qrcode-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000009_agent_push_notifications.up.sql
│   └── 000009_agent_push_notifications.down.sql
├── docs/parity/
│   ├── sync.md
│   ├── notifications.md
│   ├── push.md
│   ├── updates.md
│   └── qrcode.md
├── internal/config/config.go              # + SecureEnrollment, PreventDuplicateEnrollment,
│                                          #   PollingTimeoutMs, Module* Phase 7, UpdateManifestURL
├── internal/shared/crypto/
│   └── enrollment_signature.go          # SHA1 request/response signatures
├── internal/modules/sync/               # REPLACE scaffold
│   ├── module.go, domain/, port/, application/, adapter/http, adapter/persistence/postgres
├── internal/modules/notifications/
│   ├── port/message_queue.go            # Shared queue interface
│   └── adapter/http/handler.go + polling.go
├── internal/modules/push/               # POST /rest/private/push only
├── internal/modules/plugins/push/       # /rest/plugins/push/private/*
├── internal/modules/updates/
├── internal/modules/qrcode/
└── internal/app/modules.go              # Wire deps; remove scaffold-only route groups
```

**Structure Decision**: **notifications** owns `MessageQueue` persistence; **push** and **plugins/push** depend
on it via interfaces wired in `module.go`. **sync** aggregates read models from device/configuration repos
without importing other modules' adapters. **qrcode** is read-only public module. **updates** uses
`platform/storage` HTTP download helper.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Migration & permissions

1. `000009`: `pushmessages`, `pendingpushes`, `plugin_push_messages`, `plugin_push_schedule`.
2. Seed permissions `push_api`, `plugin_push_send`, `plugin_push_delete` for role 2.
3. Optional seed: pending message for device `hmdm-001` smoke.

### Phase B — Platform crypto & config

1. `enrollment_signature.go`: `CheckRequestSignature`, `SignSyncResponse`.
2. Extend `config.Config` + `.env.example` + `scripts/dev.sh` exports.

### Phase C — Notifications + queue (P1)

| Component | Notes |
|-----------|-------|
| `MessageQueue` port | Insert, list pending, mark delivered |
| GET `/rest/notifications/device/{n}` | PlainPushMessage list |
| GET `/rest/notification/polling/{n}` | Long-poll handler |

### Phase D — Sync module (P1)

| Endpoint | Notes |
|----------|-------|
| POST/GET `/configuration/{id}` | SyncResponse build, signatures, enrollment |
| POST `/info` | Telemetry |
| POST `/applicationSettings/{id}` | Settings persist |

### Phase E — Private push API (P2)

| Endpoint | Notes |
|----------|-------|
| POST `/rest/private/push` | `push_api`; device/group/broadcast expansion |

### Phase F — QR module (P2)

| Endpoint | Notes |
|----------|-------|
| GET `/public/qr/{key}` | PNG |
| GET `/public/qr/json/{key}` | JSON extras |

### Phase G — Updates module (P2)

| Endpoint | Notes |
|----------|-------|
| GET `/private/update/check` | Manifest fetch + compare |
| POST `/private/update` | Download/apply |

### Phase H — Push plugin (P3)

| Endpoint | Notes |
|----------|-------|
| `/rest/plugins/push/private/*` | search, send, delete, purge, tasks |

### Phase I — Docs & verification

1. Parity docs ×5; `MIGRATION.md` Phase 7 → done.
2. `make swagger`; `go test ./...`; `quickstart.md` smoke.

## Complexity Tracking

| Item | Why Needed | Simpler Alternative Rejected |
|------|------------|------------------------------|
| Shared `MessageQueue` in notifications | Single delivery pipeline for agent + admin push | Duplicate DAOs in push/plugins break delivery |
| Long-poll separate route prefix | Java servlet path `notification` vs `notifications` | Renaming would break agents |
| `plugins/push` subtree | Distinct `/rest/plugins` mount and permissions | Merging into `push` module blurs plugin contract |

## Dependencies & Risks

| Risk | Mitigation |
|------|------------|
| Large SyncResponse JSON drift | Golden-file test from seeded DB; compare fields to Java sample |
| Schedule cron not ported | Mark plugin schedule execution **partial**; CRUD still works |
| Manifest URL offline in dev | Env override or empty list + parity note |
| QR APK hash CPU cost | Cache hash per version URL in request scope |

## Next Command

Run **`/speckit-tasks`** to generate ordered `tasks.md`, then **`/speckit-implement`**.
