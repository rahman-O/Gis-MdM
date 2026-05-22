# API Contract Delta: Configurations (014 — MDM round-trip)

**Base**: `/rest/private/configurations`  
**Auth**: JWT/session; permission `configurations` for write/search  
**Envelope**: Headwind standard

**Java**: `ConfigurationResource`  
**React**: `frontend/src/features/configurations/configurationService.ts`

---

## GET `/{id}` — editor payload (extended)

**Response `data`**: `Configuration` with:

1. **SQL-backed fields**: `name`, `description`, `type`, `password`, design colors, `qrCodeKey`, `baseUrl`, `mainAppId`, `contentAppId`, `defaultFilePath`, `permissive`, …

2. **`settingsjson` flattened** into top-level camelCase keys on the same object, e.g.:
   - `gps`, `wifi`, `bluetooth`, `mobileData`
   - `kioskMode`, `kioskHome`, `kioskLock`, `kioskExit`
   - `systemUpdateType`, `appUpdateMode`, `appPermissions`
   - (full set per React `Configuration` type / Java editor)

3. **`applications[]`** each item includes:

| Field | DB source |
|-------|-----------|
| `skipVersionCheck` | `configurationapplicationparameters.skipversioncheck` |
| `remove` | `configurationapplications.remove` |
| `longTap` | `configurationapplications.longtap` |
| (existing) | `action`, `showIcon`, `usedVersionId`, `screenOrder`, … |

4. **`files[]`**, **`applicationSettings[]`** — unchanged from Phase 5.

---

## PUT `/` — create/update (extended)

**Body**: same shape as GET response (React sends full editor state).

**Server behavior**:

1. Split: column fields → `UPDATE configurations SET ...`; policy keys → merge into `settingsjson` (marshal map, do not wipe unknown keys unless explicitly nulling in payload).

2. For each `applications[]` entry on save:
   - Upsert `configurationapplications` including `remove`, `longtap`
   - Upsert `configurationapplicationparameters` for `skipVersionCheck` when `applicationId` + `configurationId` known

3. Nested `files` / `applicationSettings` — existing logic unchanged.

**Response**: `{ "status": "OK", "data": { "id": <configurationId> } }` (parity existing).

---

## Regression endpoints (must not break)

- `GET /search`, `GET /list`, `POST /autocomplete`
- `DELETE /{id}`, `PUT /copy`
- `GET /applications`, `PUT /application/upgrade`

---

## Acceptance (014)

| Check | Expected |
|-------|----------|
| Save `kioskMode: true` | GET `/{id}` → `kioskMode: true` |
| `skipVersionCheck` on app 5 | GET → app 5 has `skipVersionCheck: true` |
| `remove` / `longTap` | Round-trip in `applications[]` |
