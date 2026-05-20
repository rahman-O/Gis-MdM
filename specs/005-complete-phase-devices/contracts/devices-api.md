# API Contract: Devices (`/rest/private/devices`)

**Base path**: `/rest/private/devices`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?: string, "data"?: T }`

**Java reference**: `com.hmdm.rest.resource.DeviceResource`

---

### POST `/search`

Primary devices grid.

**Body** (`DeviceSearchRequest` — React uses `pageNum`):

```json
{
  "pageNum": 1,
  "pageSize": 50,
  "value": "",
  "groupId": null,
  "configurationId": null,
  "status": null,
  "sortBy": null,
  "sortDir": "asc"
}
```

**Response data** (`DeviceListResponse`):

```json
{
  "configurations": {
    "1": { "id": 1, "name": "Default", "permissiveMode": false }
  },
  "devices": {
    "items": [
      {
        "id": 10,
        "number": "hmdm-001",
        "configurationId": 1,
        "lastUpdate": 1710000000000,
        "statusCode": "green",
        "groups": [{ "id": 1, "name": "General" }]
      }
    ],
    "totalItemsCount": 1
  }
}
```

**Permissions**: authenticated; respects user group access.

---

### GET `/number/{number}`

Single device for detail/edit.

**Response data**: `DeviceView`

**Errors**: device not found envelope

---

### PUT `/`

Create (`id` omitted), update (`id` set), bulk config update (`ids` + `configurationId`).

**Permissions**: `edit_devices`

**Errors**: permission denied, device exists

---

### DELETE `/{id}`

**Permissions**: `edit_devices`

---

### POST `/deleteBulk`

**Body**: `{ "ids": [1, 2, 3] }`

---

### POST `/groupBulk`

**Body**:

```json
{
  "ids": [1, 2],
  "action": "set",
  "groups": [{ "id": 2, "name": "Sales" }]
}
```

---

### POST `/autocomplete`

**Body**: JSON string filter value

**Response data**: lookup items (device numbers)

---

### GET `/{id}/applicationSettings`

**Response data**: array of application settings

---

### POST `/{id}/applicationSettings`

**Body**: settings array

---

### POST `/{id}/applicationSettings/notify`

**Success**: OK (push delivery stubbed in Phase 4)

---

### POST `/{id}/description`

**Body**: description string

**Permissions**: `edit_device_desc`

---

## React consumers

`frontend/src/features/devices/deviceService.ts` — all paths above.

## Partial parity (v1)

| Area | Note |
|------|------|
| Configuration apps/files in search | May omit nested apps/files |
| Push notify | No FCM |
| Full `infojson` telemetry | Minimal columns first |
