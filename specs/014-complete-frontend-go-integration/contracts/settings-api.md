# API Contract Delta: Settings (014 — tenant fields)

**Base**: `/rest/private/settings`  
**Auth**: Session / JWT + permission `settings` where applicable  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?, "data"? }`

**Java**: `com.hmdm.rest.resource.SettingsResource`  
**React**: `frontend/src/features/settings/settingsService.ts`

**Schema**: columns from migration `000015` (prerequisite 013).

---

## GET `/` (load settings)

**Response `data`**: extended `Settings` object including:

| Field | Type | Notes |
|-------|------|--------|
| `newDeviceGroupId` | number \| null | default group for new devices |
| `phoneNumberFormat` | string | mask pattern |
| `customPropertyName1` … `3` | string | device list column labels |
| `customMultiline1` … `3` | boolean | |
| `customSend1` … `3` | boolean | |
| `desktopHeaderTemplate` | string | |
| `sendDescription` | boolean | |

Existing fields unchanged (`language`, `createNewDevices`, design colors, `idleLogout`, …).

---

## POST `/misc` (save tenant + misc)

**Body**: JSON object merging existing misc keys **and** all tenant fields above.

**Behavior**:

- Upsert `settings` row for `customerid` from auth context.
- Persist each `000015` column; do not drop unspecified keys on partial save (merge with current row).

**Response**: `{ "status": "OK" }` or ERROR with message.

---

## POST `/lang` (unchanged path)

Language-only save remains; tenant fields MUST NOT be cleared when only language is posted.

---

## User role device columns (related)

**GET/POST** `/userRole` — unchanged paths; column labels in UI read from GET settings `customPropertyName*`.

---

## Errors

| Case | HTTP / envelope |
|------|-----------------|
| Invalid `newDeviceGroupId` | ERROR, foreign key |
| Missing permission | 403 / ERROR |
| No settings row | create on first save (parity Java) |
