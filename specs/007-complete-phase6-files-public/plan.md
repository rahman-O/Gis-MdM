# Implementation Plan: Phase 6 — Files, Icons & Public API

**Branch**: `007-complete-phase6-files-public` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/007-complete-phase6-files-public/spec.md`

## Summary

Deliver **Phase 6** of the Go migration: replace scaffolds for **`files`**, **`icons`**, and **`publicapi`**
with full layered modules. Add migration `000008_files_icons_core` for `uploadedfiles`, `icons`,
`configurationfiles.fileid`, permissions (`files`, `edit_files`), and optional `customers.sizelimit`.
Introduce shared **local file storage** (`platform/storage` or equivalent) for temp upload, tenant-scoped
paths, and URL generation—reused by `configfiles` and `publicapi`. React **Files**, **Applications**
(APK upload), and **Icons** flows must work without Java; public **`/rest/public/name`**, **`/logo`**,
and **`/applications/upload`** match legacy behavior. APK metadata parsing is **best-effort** (parity
notes where Java `APKFileAnalyzer` exceeds Go MVP). **`GET /files/*`** device download servlet is
**partial** (platform static route in dev or documented deferral to Phase 7).

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `platform/auth` (`HasPermission`), `platform/httpx/response`,
`internal/shared/crypto` (MD5 for public upload hash), existing Phase 5 `applications` repo (URL lookup),
Phase 3 `customers` (filesDir, sizeLimit)

**Storage**: PostgreSQL migration `000008_files_icons_core.up.sql`; on-disk tree under
`FILES_DIRECTORY` + `customers.filesdir` (same as Java `files.directory`)

**Testing**: `go test ./internal/modules/files/... ./internal/modules/icons/... ./internal/modules/publicapi/...`;
handler tests with principal + permission flags; storage path-safety unit tests

**Target Platform**: Linux/macOS dev (`:8080`); Vite → `/rest`

**Project Type**: Web service + React (`frontend/src/features/files/`, `applications/services/webUiFilesService.ts`, `icons/`)

**Performance Goals**: File list load &lt; 5s p95 on seeded DB (spec SC-001); multipart uploads bounded by Gin max body

**Constraints**: Permissions `files`, `edit_files`, `settings` (icon delete); tenant scope; tmp-path must stay
under OS temp; `FileUtil.isSafePath` parity; push on file-configuration update **stub**; deprecated
`PublicFilesResource` **out of scope**

**Scale/Scope**: ~12 `web-ui-files` endpoints + 4 icons + 3 public; 3 modules + shared storage; ~50–65 Go files;
parity docs ×3

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `files`, `icons`, `publicapi`; Phase 6 row in `MIGRATION.md` |
| **II. Layered Clean** | ✅ | Full domain/port/application/adapter per module; shared storage in `platform/` |
| **III. API Parity** | ✅ | `contracts/*.md` + `docs/parity/files.md`, `icons.md`, `publicapi.md` |
| **IV. Testable Delivery** | ✅ | `go test` + `quickstart.md` + Swagger |
| **V. Simplicity** | ✅ | One storage helper; APK parse MVP; push/servlet partial documented |
| **VI. Security** | ✅ | Path traversal guards; public hash validation; tenant isolation |
| **VII. Observability** | ✅ | Legacy keys (`error.file.used`, `error.size.limit.exceeded`, etc.) |

**Post-design**: All gates ✅. `DownloadFilesServlet` parity deferred/partial—not a layer violation.

## Project Structure

### Documentation (this feature)

```text
specs/007-complete-phase6-files-public/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── files-api.md
│   ├── icons-api.md
│   └── publicapi-api.md
└── tasks.md                    # (/speckit-tasks)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   ├── 000008_files_icons_core.up.sql
│   └── 000008_files_icons_core.down.sql
├── docs/parity/
│   ├── files.md                                 # NEW
│   ├── icons.md                                 # NEW
│   └── publicapi.md                             # NEW
├── internal/config/config.go                    # + HashSecret, Rebranding*, Module* flags
├── internal/platform/
│   └── storage/                                 # NEW: LocalStore (safe paths, move, URL, quota)
├── internal/platform/auth/
│   └── permissions.go                           # + PermFiles, PermEditFiles
├── internal/modules/files/                      # REPLACE scaffold
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http + postgres
├── internal/modules/icons/                      # REPLACE scaffold
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http + postgres
├── internal/modules/publicapi/                  # REPLACE scaffold
│   ├── module.go
│   ├── domain/
│   ├── port/                                    # UnsecureDeviceLookup, ApplicationInsert
│   ├── application/
│   └── adapter/http
└── internal/modules/configfiles/                # REFACTOR: use platform/storage
```

**Structure Decision**: **files** owns `web-ui-files` and uploaded-file persistence; **icons** is a
small CRUD module; **publicapi** stays thin (no private auth). **platform/storage** avoids duplicating
`configfiles` disk logic. **publicapi** reuses `applications`/`devices` persistence via ports, not
direct SQL from handlers.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Database migrations & permissions

1. `000008`: `uploadedfiles` (Liquibase column set), `icons`, `configurationfiles.fileid` FK,
   `customers.sizelimit` if missing, indexes on `uploadedfiles(customerid)`.
2. Seed permissions `files`, `edit_files` for org-admin role (role id 2).
3. Sample seed: one `uploadedfiles` row + icon for customer 1 (optional smoke).

### Phase B — Platform storage & config

1. `internal/platform/storage`: `IsSafePath`, `CreateTemp`, `MoveToCustomer`, `DeleteFile`,
   `DirSizeMB`, `BuildPublicURL(baseURL, filesDir, customerDir, relativePath)`.
2. Extend `config.Config`: `HashSecret`, `RebrandingName`, `RebrandingLogo`, vendor/signup/terms links,
   `ModuleFilesEnabled`, `ModuleIconsEnabled`, `ModulePublicAPIEnabled`.
3. Update `.env.example` with new vars (mirror Java `context.xml` names).

### Phase C — Files module P1 (read + delete)

| Endpoint | Notes |
|----------|-------|
| `GET /search`, `GET /search/{value}` | `files` permission; FileView + usage flags |
| `POST /remove` | `edit_files`; FILE_USED guard |
| `GET /limit` | Multi-tenant quota |

### Phase D — Files module P2 (upload + commit)

| Endpoint | Notes |
|----------|-------|
| `POST /` | Multipart; APK parse branch |
| `POST /raw` | No APK parse |
| `POST /update` | Create / external / update |

### Phase E — Files module P3 (links)

| Endpoint | Notes |
|----------|-------|
| `GET /apps/{url}` | Delegate to applications repo by URL |
| `GET /configurations/{id}`, `POST /configurations` | Junction update; push stub |

### Phase F — Icons module

| Endpoint | Notes |
|----------|-------|
| `GET /search`, `GET /search/{value}` | Tenant scope |
| `PUT /` | Insert/update |
| `DELETE /{id}` | `settings` permission |

### Phase G — Public API module

| Endpoint | Notes |
|----------|-------|
| `GET /name` | NameResponse JSON |
| `GET /logo` | File stream or 302 to default |
| `POST /applications/upload` | MD5(deviceId+secret); UnsecureDAO device lookup |

### Phase H — Polish (optional static files)

| Item | Notes |
|------|-------|
| `GET /files/*` | Gin `StaticFS` or handler mirroring `DownloadFilesServlet` for dev URLs — **Partial** in parity |
| Refactor `configfiles` handler | Use `platform/storage` |
| Swagger, parity docs, `MIGRATION.md` Phase 6 **done** |
| `go test` + quickstart |

## Complexity Tracking

> No unjustified constitution violations. Shared `platform/storage` is justified because three modules
> (`files`, `configfiles`, `publicapi`) perform identical tenant path operations—copy-paste would violate
> Principle V. APK analyzer MVP is simpler than porting Java `APKFileAnalyzer` wholesale; parity doc
> records partial fields until a dedicated parser pass.
