# Data Model: Phase 6 — Files, Icons & Public API

**Date**: 2026-05-21  
**Storage**: PostgreSQL (`000008_files_icons_core.up.sql`) + filesystem `FILES_DIRECTORY`

## Entity: UploadedFile (`uploadedfiles`)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| customerid | int | FK → customers; NOT NULL |
| filepath | text | Relative path under tenant dir; nullable per Liquibase 24.08.25 |
| description | text | optional |
| uploadtime | bigint | ms since epoch |
| devicepath | text | path on device |
| external | boolean | default false |
| externalurl | text | when external=true |
| replacevariables | boolean | default false |

**Relationships**: Referenced by `icons.fileid`, `configurationfiles.fileid`; logical link to
applications via matching `url` string (not FK).

## Entity: Icon (`icons`)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| customerid | int | FK → customers |
| name | varchar(64) | NOT NULL |
| fileid | int | FK → uploadedfiles ON DELETE CASCADE |

**Relationships**: Optional FK from `applications.iconid` (Phase 5 schema).

## Entity: ConfigurationFile (extended)

| Field | Type | Rules |
|-------|------|-------|
| fileid | int | NEW in 000008; FK → uploadedfiles |

Used for `isFileUsed` / `getFileConfigurations` parity.

## Entity: Customer (existing, extended)

| Field | Type | Rules |
|-------|------|-------|
| filesdir | varchar | subdirectory under `FILES_DIRECTORY` |
| sizelimit | int | MB quota; 0 = unlimited; master tenant may skip enforcement |

## DTO: FileView (API list)

| Field | Notes |
|-------|-------|
| id, filePath, description, url | From UploadedFile + computed url |
| external, externalUrl | |
| usedByConfigurations | int count or id list per Java |
| usedByIcons | int count or id list |
| fileName | display name derived from path |

## DTO: FileUploadResult

| Field | Notes |
|-------|-------|
| name | original filename (UTF-8 corrected) |
| serverPath | absolute temp path after multipart |
| fileDetails | APKFileDetails subset when parsed |
| exists, complete | version conflict hints |
| application | Application copy when pkg exists |

## DTO: LimitResponse

| Field | Notes |
|-------|-------|
| sizeUsed | MB used under tenant dir |
| sizeLimit | MB limit from customer |

## DTO: FileConfigurationLink

| Field | Notes |
|-------|-------|
| fileId, configurationId | junction keys |
| notify | bool; triggers push stub |

## DTO: Icon (API)

| Field | Notes |
|-------|-------|
| id, name, fileId, fileName | `fileName` joined from uploadedfiles |

## DTO: NameResponse (public)

| Field | Notes |
|-------|-------|
| appName, vendorName, vendorLink, signupLink, termsLink | from env |

## DTO: UploadAppRequest (public multipart `app` JSON)

| Field | Notes |
|-------|-------|
| deviceId, hash | required; MD5 validation |
| name, pkg, version | application metadata |
| localPath, fileName | when binary included |
| showIcon, useKiosk, runAfterInstall, runAtBoot, system | booleans |

## Validation rules

| Rule | Error / response |
|------|------------------|
| Missing `files` permission | `error.permission.denied` |
| Missing `edit_files` on mutate | `error.permission.denied` |
| Unsafe path / tmp outside temp | permission denied |
| File referenced by config or icon | `FILE_USED` / `error.file.used` |
| Duplicate filepath in tenant | `FILE_EXISTS` |
| Storage over quota | `error.size.limit.exceeded` |
| Invalid public hash | `Invalid hash` / ERROR |
| Unknown device on public upload | device not found envelope |
| Duplicate pkg+version on public upload | duplicate application |

## State: upload commit flow

```text
multipart POST → temp file on disk (serverPath)
       ↓
POST /update (create) → move temp → tenant dir → INSERT uploadedfiles → return row + url
```

External create skips move; sets `external=true` and `externalUrl`.

## Filesystem layout

```text
{FILES_DIRECTORY}/
  {customer.filesdir}/          # e.g. customer-1/
    {relative filePath}         # APKs, assets
```

Public URL: `{BASE_URL}/files/{filesdir}/{encoded path segments}` per Java `UploadedFile.getUrl`.
