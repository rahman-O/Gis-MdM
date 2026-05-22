# Research: Device Enrollment & Sync Reliability

**Branch**: `015-device-enrollment-sync` | **Date**: 2026-05-21

## R1 — Root cause: QR provisioning payload (P0)

**Decision**: Re-port `QRCodeResource.generateQRCode` / `generateExtrasBundle` / `calculateApkHash` from Java into `qrcode/application` as a single `ProvisioningBuilder` that produces the **full** Android Device Owner JSON (outer wrapper + nested `PROVISIONING_ADMIN_EXTRAS_BUNDLE`).

**Findings (Go vs Java)**:

| Area | Java | Go today | Impact |
|------|------|----------|--------|
| PNG QR content | Full provisioning JSON (APK URL, SHA-256 checksum, WiFi, skip encryption, admin extras) | Minimal 3-field fragment only | Device Owner setup fails on scan |
| Query params `deviceId`, `create`, `useId`, `group` | Embedded in admin extras (`com.hmdm.*`) | **Ignored** | Create-on-demand enrollment broken |
| APK hash | SHA-256 base64 from file or URL | Missing | Android rejects APK download |
| `mainAppId` / `launcherUrl` | Resolved from DB | Partial (URL only, no hash) | Broken or insecure install |
| `eventReceivingComponent` | Configurable, default `AdminReceiver` | Hardcoded `AdminReceiver` | Wrong component for custom launchers |
| `/json/{key}` | Same bundle as PNG (minus image) | Returns inner fragment only | Launcher JSON path broken |

**Rationale**: User symptom «مسح QR لا يُكمل التسجيل» matches incomplete provisioning, not merely sync.

**Alternatives considered**:

- Document “use manual device ID only” — rejected; QR is primary MDM onboarding.
- Frontend-generated QR — rejected; agents and Device Owner expect server bundle.

---

## R2 — Root cause: create-on-demand configuration lookup (P0)

**Decision**: In `sync` `CreateOnDemand`, resolve `DeviceCreateOptions.configuration` by **`configurations.qrcodekey`** (case-insensitive), matching `UnsecureDAO.createNewDeviceOnDemand` → `getConfigurationByQRCodeKey`, with fallback to configuration **name** only when no QR key match (documented compat).

**Findings**: Go `resolveConfigurationID` uses `lower(name) = lower($2)` only. Agent sends QR **key** in `configuration` field after scanning (`com.hmdm.CONFIG` in admin extras).

**Rationale**: Even with fixed QR PNG, enrollment POST with `create=1` would assign wrong/default configuration.

**Alternatives considered**:

- Change React to send configuration display name — breaks parity with Java agent behavior.

---

## R3 — Static file serving for agent downloads (P0)

**Decision**: Add `GET /files/*` (or `/rest/public/files` if servlet path requires — verify Java `PublicFilesResource` deprecation) on Gin engine in `internal/app` or `platform/storage`, mapping to `FILES_DIRECTORY` + customer subdirs, matching URLs emitted by `storage.BuildPublicURL` and QR APK links.

**Findings**: `JAVA-GO-BACKEND-GAPS.md` lists `/files/*` as P1 gap; sync and QR URLs point to `{BASE_URL}/files/{filesDir}/...` but Go server has **no** static handler in `cmd/` / `internal/app`.

**Rationale**: Enrollment can “succeed” at API level while launcher cannot download APK or config files.

**Alternatives considered**:

- External nginx only — acceptable for prod but breaks local `make dev` quickstart; plan includes both dev middleware and ops note.

---

## R4 — Environment consistency (P1)

**Decision**: Document and validate in quickstart:

- `BASE_URL` must be reachable from the phone (not `localhost` unless rewritten).
- `HASH_SECRET` identical for JWT-unrelated signing and AppList MD5.
- `SECURE_ENROLLMENT` off for first UAT; enable only after signature smoke passes.
- Configuration must have `mainAppId` with `applicationversions.url` **or** `launcherUrl` for QR eligibility (mirror React `qrEligibility`).

**Rationale**: Common false negatives in dev (loopback, wrong secret).

---

## R5 — Sync payload completeness (P1)

**Decision**: Extend `BuildSyncResponse` only where agent-proven gaps exist:

- Load `applicationSettings` for device on configuration sync if Java does.
- Merge additional `settingsjson` keys agents expect (audit against Java `getDeviceSettingInternal`).
- Keep `SyncResponseHook` as documented no-op unless a plugin is required for enroll.

**Rationale**: Parity doc marks core sync Done; focus on fields that block launcher post-enroll.

**Alternatives considered**:

- Full Java line-by-line port in one PR — split into QR/static (P0) then payload polish (P1).

---

## R6 — Frontend enrollment UX (P2)

**Decision**: Minimal React changes:

- `ConfigurationEditorPage` / `EnrollmentQrPage`: surface server base URL hint and QR eligibility errors from API.
- `QrDialog`: show actionable message when PNG load returns 500.
- Optional: poll device search after QR session for “device appeared” feedback.

**Rationale**: FR-011; no full UI redesign.

---

## R7 — Testing strategy

**Decision**:

1. **Unit**: `qrcode/application` provisioning JSON golden tests vs Java fixture strings.
2. **Unit**: `sync` `CreateOnDemand` with configuration = qrCodeKey.
3. **HTTP integration**: quickstart curl chain + optional `handler_test` for QR JSON structure.
4. **Manual UAT**: real Android device on same LAN as `BASE_URL`.

**Rationale**: Constitution IV; enrollment is integration-heavy.
