---
description: "Task list for device enrollment and sync reliability (015)"
---

# Tasks: Device Enrollment & Sync Reliability

**Input**: `specs/015-device-enrollment-sync/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–7 Go modules exist; Postgres via `serverBackendGo/scripts/db-up.sh`; configuration with `qrcodekey`, `mainAppId`, and launcher APK on disk

**Tests**: Unit/golden tests for QR provisioning and sync create-on-demand per plan Phase A/B and constitution IV (application-layer coverage).

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from [spec.md](./spec.md)

## Path Conventions

- QR: `serverBackendGo/internal/modules/qrcode/`
- Sync: `serverBackendGo/internal/modules/sync/`
- Static files: `serverBackendGo/internal/platform/storage/`, `serverBackendGo/internal/app/`
- Devices: `serverBackendGo/internal/modules/devices/`
- Java reference: `backend/server/src/main/java/com/hmdm/rest/resource/QRCodeResource.java`, `SyncResource.java`
- Frontend: `frontend/src/features/devices/`, `frontend/src/features/configurations/`
- Docs: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm baseline, environment, and Java parity references before code changes.

- [X] T001 Review root-cause items R1–R3 in `specs/015-device-enrollment-sync/research.md` against current `serverBackendGo/internal/modules/qrcode/application/service.go`
- [X] T002 [P] Compare Java `backend/server/src/main/java/com/hmdm/rest/resource/QRCodeResource.java` and `UnsecureDAO.createNewDeviceOnDemand` to `specs/015-device-enrollment-sync/contracts/`
- [X] T003 [P] Document required env vars (`BASE_URL`, `FILES_DIRECTORY`, `HASH_SECRET`, `SECURE_ENROLLMENT`) in `serverBackendGo/.env.example` if any enrollment vars are missing
- [X] T004 Run baseline `cd serverBackendGo && go build ./...` and note current QR JSON output via `specs/015-device-enrollment-sync/quickstart.md` step 4 (fail expected before fixes)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Static file serving and correct configuration lookup on create-on-demand — **blocks US1 QR enrollment end-to-end**.

**⚠️ CRITICAL**: No QR UAT on a physical device until this phase completes.

- [X] T005 Implement safe static file handler with path traversal protection in `serverBackendGo/internal/platform/storage/static_files.go`
- [X] T006 Register `GET /files/*filepath` on Gin engine in `serverBackendGo/internal/app/app.go` (or `wiring.go`) using `FILES_DIRECTORY` from config
- [X] T007 [P] Add unit tests for static handler path sanitization in `serverBackendGo/internal/platform/storage/static_files_test.go`
- [X] T008 Fix `resolveConfigurationID` to lookup `configurations.qrcodekey` first, then fallback to `name`, in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go`
- [X] T009 [P] Add multi-tenant `customer` name resolution on `CreateOnDemand` in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go` matching Java `UnsecureDAO`
- [X] T010 Verify `storage.BuildPublicURL` in `serverBackendGo/internal/platform/storage/local.go` matches `/files/*` route shape from `specs/015-device-enrollment-sync/contracts/public-files-api.md`

**Checkpoint**: APK URLs from sync/QR return HTTP 200; POST enroll with `configuration=<qrCodeKey>` resolves correct `configurationid`.

---

## Phase 3: User Story 1 — تسجيل جهاز جديد عبر QR (Priority: P1) 🎯 MVP

**Goal**: Full Android Device Owner provisioning JSON from `GET /rest/public/qr/{key}` and `/json/{key}`; device enrolls on first `POST /rest/public/sync/configuration/{deviceId}` with `create=1`.

**Independent Test**: Quickstart steps 4–6 — QR JSON contains checksum and `com.hmdm.*` keys; new device appears in DB within 3 minutes of simulated enroll.

### Tests for User Story 1

- [X] T011 [P] [US1] Add golden/structure test for provisioning JSON in `serverBackendGo/internal/modules/qrcode/application/provisioning_test.go` per `specs/015-device-enrollment-sync/contracts/qrcode-api.md`
- [X] T012 [P] [US1] Add integration test for `CreateOnDemand` with `configuration` = qr key in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo_test.go`

### Implementation for User Story 1

- [X] T013 [P] [US1] Extend `QRConfig` fields and SQL SELECT in `serverBackendGo/internal/modules/qrcode/port/repository.go` and `serverBackendGo/internal/modules/qrcode/adapter/persistence/postgres/config_repo.go` (WiFi, encryption, `qrparameters`, `eventreceivingcomponent`, `apkhash`, `launcherurl`)
- [X] T014 [P] [US1] Implement SHA-256 base64 APK digest in `serverBackendGo/internal/modules/qrcode/application/apk_hash.go` (local `/files/` path and remote URL)
- [X] T015 [US1] Implement `ProvisioningBuilder` with query params `deviceId`, `create`, `useId`, `group` in `serverBackendGo/internal/modules/qrcode/application/provisioning.go` porting Java `generateExtrasBundle` / outer wrapper
- [X] T016 [US1] Refactor `serverBackendGo/internal/modules/qrcode/application/service.go` to use `ProvisioningBuilder` for both JSON and PNG endpoints
- [X] T017 [US1] Improve `serverBackendGo/internal/modules/qrcode/adapter/http/handler.go` — 500 with structured log when main app URL missing; preserve `application/json` for `/json/{id}`
- [X] T018 [US1] Pass `BASE_URL` and server context path into QR service from `serverBackendGo/internal/modules/qrcode/module.go`
- [ ] T019 [US1] Run `specs/015-device-enrollment-sync/quickstart.md` steps 4–6 and confirm Flow A in `specs/015-device-enrollment-sync/contracts/enrollment-e2e.md`

**Checkpoint**: QR scan path can download launcher APK and complete first sync enroll.

---

## Phase 4: User Story 2 — إضافة جهاز من لوحة التحكم والمزامنة الأولى (Priority: P1)

**Goal**: Admin pre-creates device via private API; agent syncs with same `number` without duplicate rows; create-on-demand still works with QR key.

**Independent Test**: `PUT /rest/private/devices` then `POST /rest/public/sync/configuration/{number}` — same `number`, updated `lastupdate`, valid `SyncResponse`.

### Tests for User Story 2

- [ ] T020 [P] [US2] Extend `serverBackendGo/internal/modules/sync/adapter/http/handler_test.go` for POST enroll with `DeviceCreateOptions.configuration` = qr key
- [ ] T021 [P] [US2] Verify devices PUT handler test coverage in `serverBackendGo/internal/modules/devices/adapter/http/handler_test.go` for create-by-number path

### Implementation for User Story 2

- [ ] T022 [US2] Audit `serverBackendGo/internal/modules/devices/application/service.go` create/update — ensure `number` uniqueness and `configurationId` assignment match Java for pre-register flow
- [ ] T023 [US2] Confirm `EnrollConfiguration` does not create duplicate when device pre-exists in `serverBackendGo/internal/modules/sync/application/service.go`
- [ ] T024 [US2] Document and execute Flow B validation in `specs/015-device-enrollment-sync/contracts/enrollment-e2e.md` via quickstart + manual UI device create

**Checkpoint**: Both QR create-on-demand and admin pre-register + agent sync paths work.

---

## Phase 5: User Story 3 — مزامنة مستمرة وحالة الجهاز في اللوحة (Priority: P1)

**Goal**: After enroll, `POST /sync/info` updates telemetry and `lastupdate`; console shows device active; config changes reach agent via sync or push.

**Independent Test**: Post-enroll `sync/info` then device search shows recent `lastUpdate`; optional `configUpdated` push smoke from quickstart.

### Implementation for User Story 3

- [X] T025 [US3] Audit `BuildSyncResponse` in `serverBackendGo/internal/modules/sync/adapter/persistence/postgres/device_sync_repo.go` against Java `SyncResource` for fields blocking launcher post-enroll
- [X] T026 [P] [US3] Load device `applicationSettings` into `SyncResponse` on configuration sync if Java includes them in `device_sync_repo.go`
- [X] T027 [US3] Ensure `TouchLastUpdate` runs on successful GET and POST configuration in `serverBackendGo/internal/modules/sync/application/service.go`
- [X] T028 [US3] Verify `DeviceStatusWriter` wiring for `devicestatuses` upsert in `serverBackendGo/internal/modules/sync/module.go` after `POST /info`
- [ ] T029 [US3] Smoke `POST /rest/public/sync/info` and admin device search per `specs/015-device-enrollment-sync/quickstart.md` steps 6–7
- [ ] T030 [P] [US3] Verify `POST /rest/private/devices/{id}/applicationSettings/notify` enqueues push in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`

**Checkpoint**: Enrolled device appears online in React device list after agent heartbeat.

---

## Phase 6: User Story 4 — تجربة مستخدم احترافية عند الفشل (Priority: P2)

**Goal**: Actionable admin/field errors for invalid QR, missing APK, wrong `BASE_URL`, secure enrollment signature failures.

**Independent Test**: Trigger known failures (bad qr key, 404 APK, `SECURE_ENROLLMENT=true` without signature) — UI or API surfaces clear messages without silent success.

### Implementation for User Story 4

- [ ] T031 [P] [US4] Show QR load failure message (HTTP 500 / network) in `frontend/src/features/devices/QrDialog.tsx`
- [ ] T032 [P] [US4] Surface QR eligibility failure reason in `frontend/src/features/configurations/ConfigurationEditorPage.tsx` when main app URL missing
- [X] T033 [US4] Improve enrollment page error states in `frontend/src/features/devices/EnrollmentQrPage.tsx` and `frontend/src/features/devices/EnrollmentQrExperience.tsx`
- [ ] T034 [US4] Add server-side slog hints for QR misconfiguration in `serverBackendGo/internal/modules/qrcode/application/service.go` (no secrets in logs)

**Checkpoint**: Operators can diagnose enrollment failures without reading server code.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Parity documentation, full UAT, and regression gates.

- [X] T035 [P] Update `serverBackendGo/docs/parity/qrcode.md` — document full provisioning bundle and verification date
- [X] T036 [P] Update `serverBackendGo/docs/parity/sync.md` — qrcodekey create-on-demand and enrollment notes
- [X] T037 [P] Add public static `/files/*` section to `serverBackendGo/docs/parity/files.md`
- [ ] T038 Run full `specs/015-device-enrollment-sync/quickstart.md` on LAN IP with Android device; record pass/fail in quickstart checklist section 9
- [X] T039 Run `cd serverBackendGo && go test ./... && go build ./...` and fix regressions
- [ ] T040 [P] Update `JAVA-GO-BACKEND-GAPS.md` § sync/QR/static files rows to reflect closed P0 gaps when UAT passes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS all user stories** (especially US1 APK download)
- **US1 (Phase 3)**: Depends on Foundational (T005–T010)
- **US2 (Phase 4)**: Depends on Foundational (T008–T009); benefits from US1 QR fixes but independently testable via pre-registered device
- **US3 (Phase 5)**: Depends on US1 or US2 having at least one enrolled device
- **US4 (Phase 6)**: Can start after US1 handler errors exist; full validation after US1–US3
- **Polish (Phase 7)**: Depends on desired user stories complete

### User Story Dependencies

| Story | Depends on | Can parallelize after |
|-------|------------|------------------------|
| US1 | Phase 2 | Phase 2 complete |
| US2 | Phase 2 | Phase 2 complete (parallel with US1 implementation) |
| US3 | US1 or US2 enroll path | US1 checkpoint |
| US4 | US1 QR errors | US1 T017+ |

### Within Each User Story

- Tests (T011–T012, T020–T021) SHOULD be written before or alongside implementation
- Repository/port before application service
- Application before HTTP handler
- Story checkpoint before next priority

### Parallel Opportunities

**Phase 2** (after T005):

```text
T007 static_files_test.go  ||  T009 customer resolution
```

**Phase 3 US1** (after T015):

```text
T013 config_repo.go  ||  T014 apk_hash.go  ||  T011 provisioning_test.go
```

**Phase 7**:

```text
T035 parity qrcode  ||  T036 parity sync  ||  T037 parity files
```

**Cross-story** (team of 2+ after Phase 2):

```text
Developer A: Phase 3 US1 (QR)
Developer B: Phase 4 US2 (devices + sync tests)
```

---

## Parallel Example: User Story 1

```bash
# After T015 ProvisioningBuilder exists, in parallel:
# T013 — extend config_repo.go / port QRConfig
# T014 — apk_hash.go
# T011 — provisioning_test.go golden assertions
```

---

## Implementation Strategy

### MVP First (User Story 1 + Foundational)

1. Complete Phase 1: Setup (T001–T004)
2. Complete Phase 2: Foundational (T005–T010) — **critical**
3. Complete Phase 3: User Story 1 (T011–T019)
4. **STOP and VALIDATE**: `quickstart.md` steps 4–6 + optional Android QR scan
5. Demo: device row in DB + APK download 200

### Incremental Delivery

1. Foundational → US1 = **MVP enrollment via QR**
2. Add US2 → pre-register + sync path
3. Add US3 → online status and policy refresh
4. Add US4 → operator-friendly errors
5. Polish → parity docs + production checklist

### Suggested MVP Scope

**T001–T019 only** (Setup + Foundational + US1) — satisfies SC-001 primary metric (QR enroll success).

---

## Notes

- Do not change REST paths (`/rest/public/qr`, `/rest/public/sync`, `/files`)
- Reuse existing Phase 7 modules; avoid duplicate static file middleware from `011-complete-migration-gaps` if already merged
- `SyncResponseHook` remains documented ⊘ unless a plugin blocks enroll
- Commit after each checkpoint; run `go test` on touched packages before PR
