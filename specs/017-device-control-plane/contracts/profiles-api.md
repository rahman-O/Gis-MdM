# Contract: Profiles API

**Feature**: `017-device-control-plane` | **Audience**: Admin React client

**Base**: `/rest/private` | **Envelope**: Headwind standard

## Endpoints

### `GET /profiles`

List profiles with usage summary.

**Response item**:

```json
{
  "id": 10,
  "name": "Kiosk retail",
  "description": "",
  "publishedVersion": 3,
  "draftVersionId": 45,
  "deviceCount": 120,
  "enrollmentRouteCount": 2
}
```

### `GET /profiles/:id`

Profile detail + current draft version id + published version id.

### `GET /profiles/:id/versions/:versionId`

Full editor payload (same shape as legacy configuration GET): scalars, `settings`, `applications`, `files`, `applicationSettings`, `policyLocks`.

### `PUT /profiles/:id/versions/:versionId`

Save **draft** only. Does not affect devices until publish.

### `POST /profiles/:id/versions/:versionId/publish`

Publish draft as new version.

**Body** (optional): `{ "confirmImpact": true }` — required when `affectedDeviceCount >= 50`.

**Response `data`**:

```json
{
  "publishedVersionId": 46,
  "versionNumber": 4,
  "artifactHash": "abc...",
  "affectedDevices": 120,
  "affectedRoutes": 2
}
```

**Pre-publish GET** `GET /profiles/:id/impact?versionId=45`:

```json
{
  "deviceCount": 120,
  "enrollmentRouteCount": 2,
  "requiresConfirmDialog": true
}
```

### `POST /profiles`

Create profile + initial draft version.

### `DELETE /profiles/:id`

Only if no enrollment routes and no devices (or soft-delete policy TBD in tasks).

## Alias (transition)

| Legacy | Maps to |
|--------|---------|
| `GET /configurations/:id` | Profile + published version OR enrollment route shim |
| `PUT /configurations/:id` | Save draft on linked profile |

Alias removed after one release; FR-010a UI does not expose Configurations.

## Permissions

`configurations` permission maps to `profiles` in v1.
