# Implementation Plan: Device Enrollment & Sync Reliability

**Branch**: `015-device-enrollment-sync` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/015-device-enrollment-sync/spec.md`

## Summary

End-to-end **QR → APK download → POST sync enrollment → ongoing sync/info** must match Java Headwind MDM behavior. Code review shows Phase 7 marked Done but three **P0 operational gaps** block real devices:

1. **QR module** returns an incomplete provisioning bundle (missing checksum, WiFi, admin extras, query params).
2. **Sync create-on-demand** resolves `configuration` by **name** instead of **`qrcodekey`**.
3. **No `/files/*` static serving** in Go — agent cannot download APKs/files linked from sync/QR.

Plan: fix P0 in `qrcode`, `sync`, and `internal/app`; polish sync payload + React UX; document UAT in quickstart and parity docs.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `github.com/skip2/go-qrcode`, `internal/shared/crypto`, existing Phase 5–7 modules (`sync`, `qrcode`, `devices`, `configurations`, `files`, `push`, `notifications`)

**Storage**: PostgreSQL (existing schema); filesystem `FILES_DIRECTORY` per customer `filesdir`

**Testing**: `go test ./internal/modules/qrcode/... ./internal/modules/sync/...`; golden JSON tests; [quickstart.md](./quickstart.md) manual/Android UAT

**Target Platform**: Linux/macOS dev server `:8080`; React admin; Android Headwind launcher

**Project Type**: Web service (Go) + React admin + MDM agent (external APK)

**Performance Goals**: QR generation &lt; 2s with local APK hash; sync p95 &lt; 3s on seeded DB (per spec SC-001/008)

**Constraints**: No REST path changes; Headwind `{status,data,message}` envelope; layer boundaries per constitution

**Scale/Scope**: ~15–25 Go files touched across 3 areas + optional small React edits; no new migration unless index on `qrcodekey` needed

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Fixes within `qrcode`, `sync`, `app` static; no new monolith |
| **II. Layered Clean** | ✅ | Provisioning logic in `qrcode/application`; HTTP thin; sync repo fix in adapter |
| **III. API Parity** | ✅ | Contracts document Java-equivalent behavior; paths unchanged |
| **IV. Testable Delivery** | ✅ | Unit + quickstart + parity doc updates required for done |
| **V. Simplicity** | ✅ | Single `ProvisioningBuilder`; one static file mount |
| **VI. Security** | ✅ | Public routes intentional; path traversal blocked on `/files` |
| **VII. Observability** | ✅ | Structured errors on QR misconfiguration; existing error keys on sync |

**Post-design**: All gates ✅. Research resolved unknowns (no NEEDS CLARIFICATION).

## Project Structure

### Documentation (this feature)

```text
specs/015-device-enrollment-sync/
├── plan.md              # This file
├── research.md          # Root-cause analysis (R1–R7)
├── data-model.md
├── quickstart.md        # UAT script
├── contracts/
│   ├── qrcode-api.md
│   ├── sync-api.md
│   ├── public-files-api.md
│   └── enrollment-e2e.md
└── tasks.md             # (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
serverBackendGo/
├── internal/app/
│   └── router.go (or modules.go)     # Mount GET /files/*
├── internal/platform/storage/
│   └── static_files.go               # Safe file server helper
├── internal/modules/qrcode/
│   ├── application/
│   │   ├── provisioning.go           # NEW: port Java QRCodeResource bundle
│   │   └── apk_hash.go               # SHA-256 base64 from local/remote APK
│   ├── port/repository.go            # Extended QRConfig fields
│   └── adapter/persistence/postgres/config_repo.go
├── internal/modules/sync/
│   └── adapter/persistence/postgres/device_sync_repo.go  # qrcodekey lookup
├── docs/parity/
│   ├── qrcode.md
│   ├── sync.md
│   └── files.md                      # static /files section
frontend/src/features/
├── devices/enrollmentQrQuery.ts      # unchanged contract
├── configurations/ConfigurationEditorPage.tsx  # optional UX hints
└── devices/QrDialog.tsx
```

**Structure Decision**: Cross-cutting static files in `internal/app` + `platform/storage` (not a new module). QR and sync remain bounded contexts.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — QR provisioning parity (P0)

1. Extend `QRConfig` SQL load: WiFi, encryption, mobile enrollment, `qrparameters`, `eventreceivingcomponent`, `apkhash`, `launcherurl`.
2. Implement `ProvisioningBuilder` with query `deviceId`, `create`, `useId`, `group`, `contextPath` (`/rest` or `""` per deployment).
3. APK SHA-256: read from `FILES_DIRECTORY` when URL under `/files/`, else HTTP digest (match Java).
4. Wire PNG + JSON handlers to same builder; return 500 with logged reason when main app URL missing.
5. Golden tests: compare output structure to Java sample for seeded configuration.

**Java reference**: `QRCodeResource.java` lines 179–447.

### Phase B — Sync enrollment fix (P0)

1. `resolveConfigurationID`: lookup `WHERE lower(qrcodekey)=lower($1)` first; fallback `name` if no row.
2. Multi-tenant: resolve `customer` name on create when body includes `customer` (match `UnsecureDAO`).
3. Integration test: POST with `configuration` = qr key creates device with correct `configurationid`.

**Java reference**: `UnsecureDAO.createNewDeviceOnDemand` line 613.

### Phase C — Public static files (P0)

1. Register `GET /files/*filepath` on Gin engine before or after `/rest`.
2. Use `FILES_DIRECTORY` + sanitized path; set content types for `.apk`.
3. Verify `BuildPublicURL` matches served paths in quickstart step 5.

**Reference**: `JAVA-GO-BACKEND-GAPS.md` §8 P1 `/files/*`.

### Phase D — Sync payload & ops polish (P1)

1. Audit `BuildSyncResponse` vs Java `SyncResource` for missing settings/apps fields blocking launcher.
2. Load device `applicationSettings` on sync if Java does.
3. Ensure `TouchLastUpdate` on successful GET/POST configuration sync.
4. Update parity docs with Verified date + known ⊘ (`SyncResponseHook`).

### Phase E — Frontend UX & documentation (P2)

1. QR load errors: show API failure message in `QrDialog` / enrollment page.
2. Configuration editor: link to quickstart when QR not eligible.
3. Run full [quickstart.md](./quickstart.md) on LAN + one Android device; record results in PR.

## Risk & Dependency Notes

- **Depends on**: Valid seed/configuration with uploaded launcher APK under customer files dir.
- **Does not fix**: MQTT-only push, full plugin `SyncResponseHook`, Mailchimp — out of scope per spec.
- **Parallel work**: `011-complete-migration-gaps` static files task merges here; avoid duplicate implementations.

## Complexity Tracking

> No constitution violations requiring justification.

| Item | Status |
|------|--------|
| (none) | — |
