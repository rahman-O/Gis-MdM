# API Contract: Configurations (016 — sync, locks, UX)

**Base**: `/rest/private/configurations`  
**Auth**: permission `configurations`  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?, "data"? }`

**Java**: `ConfigurationResource`  
**Go**: `internal/modules/configurations/adapter/http`  
**React**: `configurationService.ts`, tab components under `features/configurations/`

---

## GET `/{id}` — editor + locks

**Response `data`**: existing `Configuration` shape plus:

| Field | Type | Notes |
|-------|------|-------|
| `policyLocks` | `object` | Map of field key → `true` when locked at configuration level |
| (existing) | | Flattened `settingsjson` keys (`kioskMode`, `restrictions`, …) |
| `applications[]` | | Full round-trip per 014 |
| `applicationSettings[]` | | Includes `readonly` per row |

---

## PUT `/` — save with locks

**Body**: full editor payload including optional `policyLocks`.

**Server**:

1. Persist scalars + merge policy into `settingsjson` (preserve unknown keys).
2. Persist `policyLocks` inside `settingsjson` (namespaced key).
3. Upsert `applications`, `files`, `applicationSettings` (including `readonly` on settings).
4. Call `NotifyConfigurationChanged(configurationId)` when push module enabled.

**Response**: `{ "status": "OK", "data": { ...ConfigurationResponseMap } }` — **must** include saved MDM fields and `policyLocks` (FR-009).

---

## Regression

- `GET /search`, `/list`, `POST /autocomplete`, `DELETE /{id}`, `PUT /copy`, `PUT /application/upgrade`, `GET /applications`, `GET /applications/{id}`

---

## Acceptance

| # | Check |
|---|--------|
| C1 | Save `policyLocks: { "mainAppId": true }` → GET returns same |
| C2 | Save `kioskMode: true` + `restrictions: "no_usb"` → GET round-trip |
| C3 | Save `applicationSettings` with `readonly: true` → GET preserves |
| C4 | After save, devices on configuration receive push message type `configUpdated` (when push enabled) |
