# Contract: Sync from Profile Artifact

**Feature**: `017-device-control-plane` | **Audience**: Headwind Android agent (unchanged client in v1)

## Public endpoints (unchanged paths)

| Method | Path |
|--------|------|
| GET/POST | `/rest/public/sync/configuration/:deviceId` |
| POST | `/rest/public/sync/applicationSettings/:deviceId` |
| POST | `/rest/public/sync/info` |

## Resolution order

1. Resolve device by `deviceId` (`number`).
2. Load `enrollment_route_id` → `enrollment_routes.profile_version_id`.
3. Load `profile_version_artifacts.artifact_json`.
4. Merge device-specific overrides (`deviceapplicationsettings`) with readonly rules from artifact.
5. Return `SyncResponse`.

**Route follow**: If admin updates route to newer `profile_version_id`, step 2 returns new version on next sync without per-device migration rows.

## Response fields (additive)

Existing Java-parity fields retained. Optional additions (agent may ignore):

```json
{
  "configurationId": 5,
  "profileId": 10,
  "profileVersionId": 46,
  "profileRevision": "sha256:abc...",
  "id": "DEVICE-NUMBER",
  "applications": [],
  "files": []
}
```

- `configurationId`: maps to `enrollment_routes.id` for agent compatibility (FR-009).
- `profileRevision`: `artifact_hash` for debugging/support.

## Artifact shape

`artifact_json` is the compiled policy document produced at publish time:

- Scalar MDM fields from profile version row
- `settings` map from `settingsjson`
- `applications[]`, `files[]`, `applicationSettings[]` frozen at publish
- Design fields (colors, background)

**Must not** re-query 10 junction tables on every sync GET (NFR-004).

## Enrollment

`POST /rest/public/sync/info` with valid route QR:

- `create=1` honored → insert device if absent
- Set `tree_node_id` = route `default_tree_node_id`
- Set `enrollment_route_id`, `enrollment_state=enrolled`
- Assign `agent_id` if new row

## Parity doc

Update `serverBackendGo/docs/parity/sync.md` and add `device-control-plane.md` when implemented.

## Reference

- Java: `SyncResource.java`, `QRCodeResource.java`
- Go today: `internal/modules/sync/application`, `internal/modules/qrcode`
