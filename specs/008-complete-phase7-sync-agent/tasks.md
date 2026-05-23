---
description: "Task list for Phase 7 Agent Sync, Push, Notifications, Updates & QR migration"
---

# Tasks: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Input**: `specs/008-complete-phase7-sync-agent/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–6 complete; Postgres via `./scripts/db-up.sh`; seeded device `hmdm-001`

**Tests**: Included per User Story 8, FR-X05, and application-layer coverage for crypto/push/sync.

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Sync: `serverBackendGo/internal/modules/sync/`
- Notifications: `serverBackendGo/internal/modules/notifications/`
- Push API: `serverBackendGo/internal/modules/push/`
- Push plugin: `serverBackendGo/internal/modules/plugins/push/`
- Updates: `serverBackendGo/internal/modules/updates/`
- QR: `serverBackendGo/internal/modules/qrcode/`
- Crypto: `serverBackendGo/internal/shared/crypto/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 7 context and Java/React/agent parity baseline.

- [X] T001 Verify feature context in `specs/008-complete-phase7-sync-agent/spec.md` against `serverBackendGo/docs/MIGRATION.md` Phase 7 pending row
- [X] T002 [P] Review Java `SyncResource.java`, `NotificationResource.java`, `LongPollingServlet.java`, `PushApiResource.java`, `PushResource.java`, `UpdateResource.java`, `QRCodeResource.java` against `specs/008-complete-phase7-sync-agent/contracts/`
- [X] T003 [P] Review React `frontend/src/features/push/pushService.ts`, `frontend/src/features/updates/updatesService.ts`, `frontend/src/features/devices/enrollmentQrQuery.ts` for required paths and shapes
- [X] T004 Run baseline `cd serverBackendGo && go build ./...` and note scaffolds in `sync`, `push`, `notifications`, `updates`, `qrcode`, `plugins/push` `module.go` files

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, enrollment crypto, shared push queue, permissions, and domain/port skeletons.

**⚠️ CRITICAL**: No user story endpoints until migration `000009` applies and `MessageQueue` port exists.

- [X] T005 Create `serverBackendGo/db/migrations/000009_agent_push_notifications.up.sql` per `data-model.md` (`pushmessages`, `pendingpushes`, `plugin_push_messages`, `plugin_push_schedule`)
- [X] T006 [P] Create `serverBackendGo/db/migrations/000009_agent_push_notifications.down.sql`
- [X] T007 Seed permissions `push_api`, `plugin_push_send`, `plugin_push_delete` and `userrolepermissions` for role 2 in `000009_agent_push_notifications.up.sql`
- [X] T008 [P] Add optional dev seed: sample `pushmessages` + `pendingpushes` row for device `hmdm-001` in `000009_agent_push_notifications.up.sql`
- [X] T009 Implement `serverBackendGo/internal/shared/crypto/enrollment_signature.go` — `CheckRequestSignature`, `SignSyncResponse` per `research.md` R1
- [X] T010 [P] Add `serverBackendGo/internal/shared/crypto/enrollment_signature_test.go` for SHA1 request/response vectors
- [X] T011 Extend `serverBackendGo/internal/platform/auth/permissions.go` with `PermPushAPI`, `PermPluginPushSend`, `PermPluginPushDelete`
- [X] T012 [P] Extend `serverBackendGo/internal/platform/auth/permissions_test.go` for Phase 7 permissions
- [X] T013 Extend `serverBackendGo/internal/config/config.go` with `SecureEnrollment`, `PreventDuplicateEnrollment`, `PollingTimeoutMs`, rebranding mobile/vendor, `UpdateManifestURL`, `ModuleSyncEnabled`, `ModulePushEnabled`, `ModuleNotificationsEnabled`, `ModuleUpdatesEnabled`, `ModuleQRCodeEnabled`
- [X] T014 [P] Document Phase 7 env vars in `serverBackendGo/.env.example` and export in `serverBackendGo/scripts/dev.sh`
- [X] T015 Define `serverBackendGo/internal/modules/notifications/port/message_queue.go` — `Enqueue`, `ListPending`, `MarkDelivered` per `research.md` R3
- [X] T016 Implement `serverBackendGo/internal/modules/notifications/adapter/persistence/postgres/queue_repo.go` for `pushmessages` and `pendingpushes`
- [X] T017 [P] Create `serverBackendGo/internal/modules/sync/domain/sync.go` — `DeviceCreateOptions`, `SyncResponse`, `DeviceInfo`, `SyncApplicationSetting` per `contracts/sync-api.md`
- [X] T018 [P] Create `serverBackendGo/internal/modules/notifications/domain/message.go` — `PlainPushMessage` per `contracts/notifications-api.md`
- [X] T019 [P] Create `serverBackendGo/internal/modules/push/domain/push.go` — `PushRequest` per `contracts/push-api.md`
- [X] T020 [P] Create `serverBackendGo/internal/modules/updates/domain/update.go` — `UpdateEntry`, `UpdateRequest` per `contracts/updates-api.md`
- [X] T021 [P] Create `serverBackendGo/internal/modules/qrcode/domain/qr.go` — query params and provisioning bundle types per `contracts/qrcode-api.md`
- [X] T022 [P] Create `serverBackendGo/internal/modules/plugins/push/domain/plugin_push.go` — filter/send/schedule DTOs per `contracts/push-api.md`
- [X] T023 Define `serverBackendGo/internal/modules/sync/port/repository.go` — device lookup/create, configuration bundle, application settings persistence
- [X] T024 Verify migration: `cd serverBackendGo && make migrate` and confirm `pushmessages`, `pendingpushes`, `plugin_push_messages` tables exist
- [X] T025 Remove scaffold-only route registration from `serverBackendGo/internal/modules/sync/module.go`, `push/module.go`, `notifications/module.go`, `updates/module.go`, `qrcode/module.go`, `plugins/push/module.go` until handlers ship (per FR-X08)

**Checkpoint**: Schema + crypto + queue + ports ready — proceed to user stories.

---

## Phase 3: User Story 1 — Device enrollment and configuration sync (Priority: P1) 🎯 MVP

**Goal**: `POST`/`GET /rest/public/sync/configuration/{deviceId}` with `SyncResponse`, signatures, enrollment rules.

**Independent Test**: `curl` GET configuration for `hmdm-001` returns full payload; secure enrollment rejects bad signature.

### Tests for User Story 1

- [X] T026 [P] [US1] Add `serverBackendGo/internal/modules/sync/application/service_test.go` for enrollment, duplicate device, signature gate
- [X] T027 [P] [US1] Add `serverBackendGo/internal/modules/sync/adapter/http/handler_test.go` for GET `/configuration/:deviceId`

### Implementation for User Story 1

- [X] T028 [US1] Implement `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go` — device by number/old/imei, create on demand, configuration/apps/files aggregation
- [X] T029 [US1] Implement `serverBackendGo/internal/modules/sync/application/build_response.go` and `GetConfiguration` in `serverBackendGo/internal/modules/sync/application/service.go`
- [X] T030 [US1] Create `serverBackendGo/internal/modules/sync/adapter/http/handler.go` — register POST/GET `/configuration/:deviceId`, set `X-Response-Signature` and `X-IP-Address` headers
- [X] T031 [US1] Wire `serverBackendGo/internal/modules/sync/module.go` on `/rest/public/sync` with `MODULE_SYNC_ENABLED`

**Checkpoint**: Agent configuration sync works for seeded device.

---

## Phase 4: User Story 2 — Device telemetry and per-app settings (Priority: P1)

**Goal**: `POST /rest/public/sync/info` and `POST /rest/public/sync/applicationSettings/{deviceId}`.

**Independent Test**: POST info updates device row; POST applicationSettings persists settings; unknown device returns not found.

### Tests for User Story 2

- [X] T032 [P] [US2] Extend `serverBackendGo/internal/modules/sync/application/service_test.go` for info update and multi-tenant create rejection
- [X] T033 [P] [US2] Extend `serverBackendGo/internal/modules/sync/adapter/http/handler_test.go` for POST `/info` and `/applicationSettings/:deviceId`

### Implementation for User Story 2

- [X] T034 [US2] Extend `device_sync_repo.go` with `UpdateDeviceInfo`, `CompleteMigration`, `SaveApplicationSettings`
- [X] T035 [US2] Implement `UpdateInfo` and `SaveApplicationSettings` in `serverBackendGo/internal/modules/sync/application/service.go`
- [X] T036 [US2] Register POST `/info` and POST `/applicationSettings/:deviceId` in `serverBackendGo/internal/modules/sync/adapter/http/handler.go` with Swagger comments

**Checkpoint**: Agent heartbeat and settings endpoints operational.

---

## Phase 5: User Story 3 — Agent notification delivery (Priority: P1)

**Goal**: `GET /rest/notifications/device/{deviceNumber}` and long-poll `GET /rest/notification/polling/{deviceNumber}`.

**Independent Test**: After queue insert, GET returns messages; long-poll returns within timeout.

### Tests for User Story 3

- [X] T037 [P] [US3] Add `serverBackendGo/internal/modules/notifications/application/service_test.go` for pending list and mark delivered
- [X] T038 [P] [US3] Add `serverBackendGo/internal/modules/notifications/adapter/http/handler_test.go` for GET `/device/:deviceNumber`

### Implementation for User Story 3

- [X] T039 [US3] Implement `serverBackendGo/internal/modules/notifications/application/service.go` — resolve device, list/mark via `MessageQueue`
- [X] T040 [US3] Create `serverBackendGo/internal/modules/notifications/adapter/http/handler.go` — GET `/rest/notifications/device/:deviceNumber`
- [X] T041 [US3] Create `serverBackendGo/internal/modules/notifications/adapter/http/polling.go` — GET `/rest/notification/polling/:deviceNumber` with timeout and signature check
- [X] T042 [US3] Wire `serverBackendGo/internal/modules/notifications/module.go` — register JAX-RS group and polling route on engine root per `research.md` R4

**Checkpoint**: Agent can pull and long-poll pending push messages.

---

## Phase 6: User Story 4 — Administrator sends push from console (Priority: P2)

**Goal**: `POST /rest/private/push` with `push_api` permission and device/group/broadcast targeting.

**Independent Test**: JWT POST push queues row; GET notifications returns it for target device.

### Tests for User Story 4

- [X] T043 [P] [US4] Add `serverBackendGo/internal/modules/push/application/service_test.go` for permission, targeting expansion, invalid device
- [X] T044 [P] [US4] Add `serverBackendGo/internal/modules/push/adapter/http/handler_test.go` for POST `/` with Bearer auth

### Implementation for User Story 4

- [X] T045 [US4] Define `serverBackendGo/internal/modules/push/port/device_resolver.go` — tenant-scoped device list by numbers/groups/broadcast
- [X] T046 [US4] Implement `serverBackendGo/internal/modules/push/application/service.go` — enqueue via `notifications/port.MessageQueue`
- [X] T047 [US4] Create `serverBackendGo/internal/modules/push/adapter/http/handler.go` — POST `/rest/private/push` per `PushApiResource`
- [X] T048 [US4] Wire `serverBackendGo/internal/modules/push/module.go` with `MODULE_PUSH_ENABLED`

**Checkpoint**: React `pushService.ts` send path works end-to-end with notifications GET.

---

## Phase 7: User Story 5 — QR enrollment for configurations (Priority: P2)

**Goal**: `GET /rest/public/qr/{key}` PNG and `GET /rest/public/qr/json/{key}` provisioning JSON.

**Independent Test**: Valid `qrCodeKey` returns PNG content-type; json endpoint returns extras string.

### Tests for User Story 5

- [X] T049 [P] [US5] Add `serverBackendGo/internal/modules/qrcode/application/service_test.go` for unknown key and extras bundle fields
- [X] T050 [P] [US5] Add `serverBackendGo/internal/modules/qrcode/adapter/http/handler_test.go` for GET `/:id` response headers

### Implementation for User Story 5

- [X] T051 [US5] Add `github.com/skip2/go-qrcode` to `serverBackendGo/go.mod` per `research.md` R5
- [X] T052 [US5] Implement `serverBackendGo/internal/modules/qrcode/adapter/persistence/postgres/config_repo.go` — configuration by `qrCodeKey`, main app version URL/hash
- [X] T053 [US5] Implement `serverBackendGo/internal/modules/qrcode/application/extras_bundle.go` and `service.go` — loopback URL rewrite, SHA-256 helper using `platform/storage`
- [X] T054 [US5] Create `serverBackendGo/internal/modules/qrcode/adapter/http/handler.go` — GET `/:id` and GET `/json/:id`
- [X] T055 [US5] Wire `serverBackendGo/internal/modules/qrcode/module.go` on `/rest/public/qr` with `MODULE_QRCODE_ENABLED`

**Checkpoint**: React `EnrollmentQrPage` loads QR image without Java.

---

## Phase 8: User Story 6 — Check and apply product updates (Priority: P2)

**Goal**: `GET /rest/private/update/check` and `POST /rest/private/update`.

**Independent Test**: Super-admin GET check returns entries; POST with download flags updates filesystem rows.

### Tests for User Story 6

- [X] T056 [P] [US6] Add `serverBackendGo/internal/modules/updates/application/service_test.go` for multi-tenant non-super-admin denial and manifest parse
- [X] T057 [P] [US6] Add `serverBackendGo/internal/modules/updates/adapter/http/handler_test.go` for GET `/check`

### Implementation for User Story 6

- [X] T058 [US6] Implement `serverBackendGo/internal/modules/updates/application/manifest.go` — fetch `UpdateManifestURL`, parse web/launcher/mobile entries per Java `UpdateResource`
- [X] T059 [US6] Implement `serverBackendGo/internal/modules/updates/application/service.go` — `CheckUpdates`, `ApplyUpdates` with `platform/storage` download helper
- [X] T060 [US6] Create `serverBackendGo/internal/modules/updates/adapter/http/handler.go` — GET `/check`, POST `/` on `/rest/private/update`
- [X] T061 [US6] Wire `serverBackendGo/internal/modules/updates/module.go` with `MODULE_UPDATES_ENABLED`; document stats POST **partial** stub in service

**Checkpoint**: Updates page `checkUpdates()` succeeds against Go backend.

---

## Phase 9: User Story 7 — Push plugin administration (Priority: P3)

**Goal**: `/rest/plugins/push/private/*` search, send, delete, purge, schedule CRUD.

**Independent Test**: POST search returns paginated rows; send requires `plugin_push_send`; delete requires `plugin_push_delete`.

### Tests for User Story 7

- [X] T062 [P] [US7] Add `serverBackendGo/internal/modules/plugins/push/application/service_test.go` for plugin permissions and tenant scope
- [X] T063 [P] [US7] Add `serverBackendGo/internal/modules/plugins/push/adapter/http/handler_test.go` for POST `/private/search`

### Implementation for User Story 7

- [X] T064 [US7] Implement `serverBackendGo/internal/modules/plugins/push/adapter/persistence/postgres/message_repo.go` and `schedule_repo.go`
- [X] T065 [US7] Implement `serverBackendGo/internal/modules/plugins/push/application/service.go` — search/send/delete/purge/tasks; enqueue via shared `MessageQueue` + write `plugin_push_messages`
- [X] T066 [US7] Create `serverBackendGo/internal/modules/plugins/push/adapter/http/handler.go` — all routes in `contracts/push-api.md` plugin section
- [X] T067 [US7] Wire `serverBackendGo/internal/modules/plugins/push/module.go` on `/rest/plugins/push`; document schedule cron execution as **partial** in parity

**Checkpoint**: Legacy push plugin API surface available; schedule firing deferred.

---

## Phase 10: User Story 8 — Verifiable API & regression safety (Priority: P2)

**Goal**: Swagger tags, parity docs, cross-module smoke.

**Independent Test**: `make swagger` lists Phase 7 routes; `go test` passes; `quickstart.md` curl blocks succeed.

### Implementation for User Story 8

- [X] T068 [P] [US8] Create `serverBackendGo/docs/parity/sync.md` with endpoint table and Java reference
- [X] T069 [P] [US8] Create `serverBackendGo/docs/parity/notifications.md` including long-poll path note
- [X] T070 [P] [US8] Create `serverBackendGo/docs/parity/push.md` covering `/private/push` and plugin routes
- [X] T071 [P] [US8] Create `serverBackendGo/docs/parity/updates.md` and `serverBackendGo/docs/parity/qrcode.md`
- [X] T072 [US8] Run `cd serverBackendGo && make swagger` and verify Sync, Notifications, Push, Updates, QR tags in `internal/platform/httpx/swagger/swagger.yaml`
- [X] T073 [US8] Run `go test ./internal/modules/sync/... ./internal/modules/notifications/... ./internal/modules/push/... ./internal/modules/updates/... ./internal/modules/qrcode/... ./internal/modules/plugins/push/... ./internal/shared/crypto/...`

**Checkpoint**: Phase 7 modules documented and test-green.

---

## Phase 11: Polish & Cross-Cutting Concerns

**Purpose**: Migration status, docs, validation, no-op hooks.

- [X] T074 [P] Add no-op `SyncResponseHook` registry stub in `serverBackendGo/internal/modules/sync/application/hooks.go` with **partial** note in `docs/parity/sync.md`
- [X] T075 Update `serverBackendGo/docs/MIGRATION.md` — Phase 7 row **done**; modules `sync`, `push`, `notifications`, `updates`, `qrcode`
- [X] T076 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` Phase 7 section to **منجز** with parity links
- [X] T077 Run validation steps in `specs/008-complete-phase7-sync-agent/quickstart.md` (sync, push+notifications, QR, updates)
- [X] T078 [P] Optional: `cd ../frontend && npm run dev` — verify Enrollment QR page and Updates check
- [X] T079 Final `cd serverBackendGo && go build ./...` and fix layer import violations

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **BLOCKS** all user stories
- **Phase 3 (US1)**: Depends on Phase 2
- **Phase 4 (US2)**: Depends on Phase 3 (same `sync` module; extends repo/service)
- **Phase 5 (US3)**: Depends on Phase 2 queue repo
- **Phase 6 (US4)**: Depends on Phase 5 (`MessageQueue` delivery path)
- **Phase 7 (US5)**: Depends on Phase 2; independent of US1–US4
- **Phase 8 (US6)**: Depends on Phase 2; independent of US1–US5
- **Phase 9 (US7)**: Depends on Phase 5 queue + Phase 2 plugin tables
- **Phase 10 (US8)**: Depends on US1–US7 handlers
- **Phase 11 (Polish)**: Depends on Phase 10 minimum

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 | Phase 2 | Sync configuration |
| US2 | US1 | Same sync module |
| US3 | Phase 2 | Notifications queue |
| US4 | US3 | Push enqueue + agent pull |
| US5 | Phase 2 | QR public |
| US6 | Phase 2 | Updates private |
| US7 | US3, Phase 2 | Plugin + shared queue |
| US8 | US1–US7 | Docs/tests |

### Parallel Opportunities

- Phase 1: T002, T003
- Phase 2: T006–T008, T010–T012, T014, T017–T022 [P] after T005
- After Phase 2: **US5** (T049–T055) and **US6** (T056–T061) parallel with **US1** (T026–T031)
- After US3: **US4** sequential; **US7** after US4 optional
- Phase 10: T068–T071 [P]
- Phase 11: T074, T076, T078 [P]

### Parallel Example: P1 agent path (US1 + US2 + US3)

```bash
# After Phase 2:
# Track A: T026–T031 (US1 sync configuration)
# Track B: T037–T042 (US3 notifications) — can start once T016 queue repo done
# Then T034–T036 (US2) extends sync
# Validate: quickstart §3–§4
```

### Parallel Example: Admin UI (US5 + US6)

```bash
# After Phase 2:
# Developer A: T049–T055 qrcode
# Developer B: T056–T061 updates
# No dependency on sync/notifications
```

---

## Implementation Strategy

### MVP First (US1 + US2 + US3)

1. Complete Phase 1–2: Setup + Foundational (migration `000009` + queue + crypto)
2. Complete Phase 3–4: US1 + US2 sync (agent enroll + heartbeat)
3. Complete Phase 5: US3 notifications (agent message delivery)
4. **STOP and VALIDATE**: `quickstart.md` §3–§4 without Android emulator optional
5. Continue US4 → US5 → US6 → US7 → US8 → Polish

### Incremental Delivery

1. Foundation → US1 + US2 + US3 (agent operational on HTTP)
2. US4 (console push) — proves full push loop
3. US5 + US6 in parallel (QR + updates admin)
4. US7 (plugin API) → US8 + Polish (parity, MIGRATION done)

### Suggested MVP Scope

**Minimum**: Phases 1–2 + **US1** + **US2** + **US3** + T068–T069 subset + T077 sync/notifications smoke.

Delivers agent sync and push delivery without Java; defers QR, updates, and plugin UI to next slice.

---

## Notes

- Postgres table names: use lowercase `pushmessages`, `pendingpushes` (Java `pushMessages` folds).
- Long-poll path is **`/rest/notification/polling`** (singular), not `notifications`.
- Private push is **`POST /rest/private/push`** (React); plugin is **`/rest/plugins/push/private/*`** (Angular).
- `SyncResponseHook` and update `sendStats` are **partial** per spec — document in parity, do not fake success.
- Branch: `008-complete-phase7-sync-agent`
