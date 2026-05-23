# Research: Phase 6 — Files, Icons & Public API

**Date**: 2026-05-21

## R1 — Schema: `uploadedfiles` and `icons` in Go migrations

**Decision**: Add migration `000008_files_icons_core` creating `uploadedfiles` and `icons` with
Liquibase-equivalent columns; add `configurationfiles.fileid INT REFERENCES uploadedfiles(id)`.

**Rationale**: Go migrations through `000007` do not include these tables; Java MyBatis mappers
assume they exist. Phase 5 `configurationfiles` lacks `fileid`, blocking “file in use” checks.

**Alternatives considered**:
- Rely on full Liquibase import only — rejected; local `make dev` uses golang-migrate SQL only.
- Skip DB and disk-only — rejected; React list/delete requires `uploadedfiles` rows.

---

## R2 — Shared file storage location

**Decision**: New package `internal/platform/storage` with path-safety, temp upload, move into
`FILES_DIRECTORY/{customer.filesdir}/`, URL builder, directory size for quota.

**Rationale**: `configfiles` handler already writes disks; `FilesResource` duplicates `FileUtil`
logic. One helper keeps modules thin (Principle V).

**Alternatives considered**:
- All logic inside `files` module only — rejected; `configfiles` and `publicapi` would import
  files module adapter (layer violation).
- External object store (S3) — rejected; Java uses local disk; out of scope.

---

## R3 — APK metadata on multipart upload

**Decision**: **MVP**: parse APK as ZIP, read `AndroidManifest.xml` via lightweight decoder
(e.g. `github.com/shogo82148/androidbinary/apk` or minimal custom XML) for `package`,
`versionName`, `versionCode`, and split `arch` when present. If parse fails, return
`FileUploadResult` with `name` + `serverPath` only (no `fileDetails`).

**Rationale**: React upload flow uses `exists` / `complete` / `application` hints from Java
`APKFileAnalyzer`; partial metadata is better than blocking upload. Full split-APK rules
mirrored from Java when version rows exist in DB (query `applicationversions`).

**Alternatives considered**:
- Shell out to `aapt`/`aapt2` — rejected for dev portability.
- No parsing (Phase 5 deferral) — rejected; spec SC-002 requires APK upload UX.

---

## R4 — Public upload hash algorithm

**Decision**: `MD5(deviceId + hashSecret)` uppercase compare, matching `PublicResource` and
`CryptoUtil.getMD5String` usage; reuse `internal/shared/crypto.MD5UpperHex` on concatenation
string (verify byte-exact match with Java in tests).

**Rationale**: AppList utility and agents depend on stable hash contract.

**Alternatives considered**:
- SHA-256 upgrade — rejected (breaking parity).

---

## R5 — Device/agent file download (`DownloadFilesServlet`)

**Decision**: **Partial in Phase 6**: document in `parity/files.md`; implement optional Gin route
`GET /files/*` under `internal/app` using `platform/storage` + `PUBLIC_IP_ALLOWLIST` when time
permits in Phase H. React admin flows do not require it; application `url` fields reference
`{baseUrl}/files/...` URLs used by agents in Phase 7.

**Rationale**: Servlet includes secure-enrollment signature checks tied to `ApplicationDAO.isMainApp`;
full parity is a cross-cutting concern beyond admin REST.

**Alternatives considered**:
- Full servlet port in Phase 6 — deferred to limit scope; not blocking SC-001–SC-003.

---

## R6 — Icon permissions

**Decision**: Match Java: GET/PUT icons without explicit permission check in `IconResource`;
DELETE requires `settings` permission.

**Rationale**: Parity over tightening security in migration slice.

**Alternatives considered**:
- Require `files` on all icon routes — rejected without product approval.

---

## R7 — Push on file-configuration update

**Decision**: Stub `port.PushNotifier` no-op; return `OK` after DB update.

**Rationale**: Same approach as Phase 5 configuration upgrade; Phase 7 owns real push.

---

## R8 — Rebranding configuration

**Decision**: Env vars: `REBRANDING_NAME`, `REBRANDING_LOGO`, `REBRANDING_VENDOR_NAME`,
`REBRANDING_VENDOR_LINK`, `REBRANDING_SIGNUP_LINK`, `REBRANDING_TERMS_LINK`, plus
`HASH_SECRET` for public upload.

**Rationale**: Mirrors Java `@Named("rebranding.*")` and `hash.secret` in `context.xml`.

**Alternatives considered**:
- Database-driven rebranding — rejected; Java uses deployment parameters.
