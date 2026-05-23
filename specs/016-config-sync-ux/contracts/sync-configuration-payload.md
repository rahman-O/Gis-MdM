# API Contract: Sync configuration payload (016)

**Base**: `/rest/public/sync`  
**Auth**: `X-Request-Signature` (when `SECURE_ENROLLMENT` / Java parity)  
**Envelope**: raw JSON body (agent contract — not Headwind admin envelope)

**Java**: `SyncResource`, `SyncResponse.java`  
**Go**: `internal/modules/sync`

---

## GET `/configuration/{deviceId}`

Returns assembled policy for enrolled device.

## POST `/configuration/{deviceId}`

Creates device on demand when body includes `configuration` (qrcode key), then same response shape.

---

## Response body — required field parity (subset)

Go MUST populate at minimum (when set in DB):

| Field | Source |
|-------|--------|
| `deviceId`, `configurationId` | device row |
| `password`, `backgroundColor`, `textColor` | configuration columns |
| `permissive` | configuration |
| `applications[]` | `configurationapplications` + versions |
| `files[]` | `configurationfiles` + `BASE_URL` |
| `applicationSettings[]` | merge `configurationapplicationsettings` + `deviceapplicationsettings`; config `readonly` wins |
| `kioskMode`, `kioskHome`, `kioskLock`, … | `settingsjson` / columns per Java map |
| `restrictions` | configuration column or policy JSON |
| `gps`, `wifi`, `bluetooth`, `mobileData`, … | policy JSON |
| `pushOptions`, `requestUpdates` | policy JSON (existing) |
| `appName`, `vendor` | rebranding env (existing) |

**Implementation note**: add fields to `domain.SyncResponse` and map in `BuildSyncResponse` via dedicated mapper (avoid ad-hoc `map[string]any` in handler).

---

## POST `/applicationSettings/{deviceId}`

**Body**: array of `{ packageId, name, type, value, readonly?, lastUpdate? }`

**Rules**:

1. Resolve device → `configurationid`.
2. Load configuration default settings where `readonly = true`.
3. For each incoming item matching `(packageId, name)` on a readonly config default: **ignore update** (keep config value).
4. Other items: upsert `deviceapplicationsettings` (replace `ON CONFLICT DO NOTHING` with proper upsert where Java does).

---

## Acceptance

| # | Check |
|---|--------|
| S1 | After configuration save with new app, sync `applications` includes new entry with HTTPS `url` |
| S2 | `kioskMode: true` in configuration → sync JSON `kioskMode: true` |
| S3 | `restrictions` non-empty → sync JSON includes same string |
| S4 | Config app setting `readonly: true` → sync shows `readonly: true`; device POST cannot change value |
| S5 | Golden test: Go sync JSON ⊇ Java sync JSON for same DB fixture (critical keys) |
