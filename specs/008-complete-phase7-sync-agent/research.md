# Research: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Branch**: `008-complete-phase7-sync-agent` | **Date**: 2026-05-21

## R1 — Enrollment request/response signatures

**Decision**: Extend `internal/shared/crypto` with `EnrollmentRequestSignature(header, secret+deviceId)` and
`SyncResponseSignature(secret, syncResponseJSON)` matching Java `CryptoUtil`:
`SHA1(value)` uppercase hex for requests; response signature = `SHA1(hashSecret + compactJSON)` where
compact JSON is marshaled struct with whitespace stripped.

**Rationale**: Agents and `secureEnrollment` depend on exact algorithm; MD5 is only used for public AppList
upload (Phase 6), not sync.

**Alternatives considered**:
- HMAC-SHA256 — breaks agent compatibility.
- Skip signatures in dev — rejected; must be env-flagged like Java (`SECURE_ENROLLMENT`).

---

## R2 — SyncResponse assembly

**Decision**: Implement `sync/application/build_response.go` that loads device, configuration, applications
(with versions/URLs from Phase 5 repos), configuration files (Phase 5/6), and settings in one transaction
scope; port interfaces `DeviceSyncRepository`, `ConfigurationSyncRepository` in `sync/port` backed by
postgres adapters reusing SQL patterns from `devices` and `configurations` modules (no cross-module SQL
imports).

**Rationale**: `SyncResource.getDeviceSettingInternal` is ~200 lines of aggregation; isolating in
application keeps handler thin.

**Alternatives considered**:
- Single mega `UnsecureDAO` Go file — violates module boundaries.
- Call HTTP to other modules — rejected.

**Partial**: `SyncResponseHook` Guice extensions → no-op hook registry; documented in parity.

---

## R3 — Push message persistence (shared queue)

**Decision**: Migration `000009` creates `pushmessages` and `pendingpushes` (lowercase identifiers, matching
PostgreSQL fold of Java `pushMessages`). Single `internal/platform/pushstore` package OR
`notifications/adapter/persistence/postgres/queue_repo.go` implementing `port.MessageQueue` consumed by:
- `push` module (`PushApiResource`)
- `plugins/push` module (plugin history + queue)
- `notifications` module (delivery GET + long-poll)

**Rationale**: Java uses one `NotificationDAO` + `PushService`; duplicating insert logic in three modules
violates constitution V.

**Alternatives considered**:
- Separate tables per module — diverges from Java.
- In-memory queue — breaks multi-instance and agent polling.

---

## R4 — Long-polling in Gin

**Decision**: Register `GET /rest/notification/polling/*deviceNumber` on root engine (not only `/rest/public`
group) via `notifications` module; handler uses `context.WithTimeout` for `POLLING_TIMEOUT_MS` (default 60s),
polls DB every 1s for pending message, returns JSON array on hit or empty body on timeout (match servlet).

**Rationale**: Java servlet path is singular `notification` not `notifications`; agents hard-code URL.

**Alternatives considered**:
- WebSockets — not in legacy API.
- SSE — not used by Headwind agent.

---

## R5 — QR code generation

**Decision**: Use `github.com/skip2/go-qrcode` for PNG; JSON endpoint returns raw provisioning string from
`generateExtrasBundle` logic ported from `QRCodeResource` (Android Device Admin extras format).

**Rationale**: Java uses `qrgen`; Go library is stable and avoids cgo.

**Alternatives considered**:
- Return redirect to external QR API — breaks offline enrollment.
- Embed only URL string without PNG — React enrollment page needs image.

---

## R6 — Updates manifest

**Decision**: Port `UpdateSettings.MANIFEST_URL` pattern: fetch remote manifest with `CUSTOMER_DOMAIN`
substituted from `BASE_URL` host; compare versions using installed files + `applicationversions` rows;
`POST /private/update` downloads via `platform/storage` HTTP client; stats POST to vendor URL is **stub**
(log-only) unless `UPDATE_STATS_ENABLED=true`.

**Rationale**: Super-admin check and manifest parsing are required for React `checkUpdates`; stats are
non-blocking for MDM core.

**Alternatives considered**:
- Disable updates in Go — breaks maintenance UI.
- Bundle static manifest in repo — diverges from vendor update channel.

---

## R7 — Push plugin vs private push API

**Decision**: Two modules:
- `internal/modules/push` — `POST /rest/private/push` only (`PushApiResource`).
- `internal/modules/plugins/push` — `/rest/plugins/push/private/*` (`PushResource`).

Both call shared `MessageQueue` port; plugin also writes `plugin_push_messages` history table.

**Rationale**: Paths and permissions differ; React uses `/private/push`; Angular plugin uses `/plugins/push`.

---

## R8 — Feature flags & config

**Decision**: Add to `config.Config`:
`SecureEnrollment`, `PreventDuplicateEnrollment`, `PollingTimeoutMs`, `RebrandingMobileName`,
`RebrandingVendorName` (for sync branding), `ModuleSyncEnabled`, `ModulePushEnabled`,
`ModuleNotificationsEnabled`, `ModuleUpdatesEnabled`, `ModuleQRCodeEnabled`, `UpdateManifestURL` (override).

**Rationale**: Mirrors Java `context.xml` / named bindings in `SyncResource` constructor.

---

## R9 — Event bus (battery/location)

**Decision**: Persist telemetry in `devices` row only; fire no external events in Phase 7. Parity doc notes
**partial** vs Java `EventService`.

**Rationale**: No Go consumer for events yet; scope control.

---

## R10 — Testing strategy

**Decision**:
- Unit: signature crypto, push targeting expansion, QR key not found, update guard (non-super-admin).
- HTTP: sync configuration GET for seeded `hmdm-001`, notifications GET after push POST, QR PNG content-type.
- Manual: `quickstart.md` curl blocks.

**Rationale**: Constitution IV; agent E2E requires Android emulator (out of CI scope).
