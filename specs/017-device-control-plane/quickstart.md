# Quickstart: Device Control Plane

**Feature**: `017-device-control-plane` | **Branch**: `017-device-control-plane`

## Prerequisites

- Go server running (`make dev` or `serverBackendGo/cmd/server`)
- Postgres migrated through latest + feature migrations
- React admin on Vite proxy to `:8080`
- Headwind agent on test device (no app changes in v1)

## Sprint 1 — Tree + enrollment identity

1. Apply migrations (`device_tree_nodes`, `devices.agent_id`, `tree_node_id`, `enrollment_state`).
2. `curl -H "Authorization: Bearer $JWT" http://localhost:8080/rest/private/device-tree`
3. Create folder «Test site»; move a device via UI or API.
4. Create enrollment route (or migrated config) with `create=1` QR.
5. Enroll new device → appears under default tree folder within 60s (SC-002).

## Sprint 2 — Profiles publish

1. Open Profile → edit draft → Publish.
2. `GET /profiles/:id/impact` shows device counts.
3. With ≥50 devices, UI blocks publish without confirm.
4. Verify `profile_version_artifacts` row exists.

## Sprint 3 — Sync artifact

1. `curl` `GET /rest/public/sync/configuration/{deviceNumber}`
2. Compare JSON to pre-migration baseline (SC-005 staging sample).
3. Update route to newer profile version → device sync picks up v2 (route follow).

## Sprint 4 — Events (optional in v1 tail)

1. Publish profile → `domain_events` row → push worker batches notify.
2. Admin UI remains responsive during publish (FR-014).

## Migration smoke

Backfill is applied by migrations `000022`–`000024` (configurations → profiles + routes) and `000027` (device tree placement). No separate `0000xx_backfill_configurations` file.

```bash
cd serverBackendGo && make migrate
psql "$DATABASE_URL" -c "SELECT count(*) AS profiles FROM profiles;"
psql "$DATABASE_URL" -c "SELECT count(*) AS routes FROM enrollment_routes;"
psql "$DATABASE_URL" -c "SELECT count(*) AS orphan_devices FROM devices WHERE tree_node_id IS NULL;"  # expect 0
# Legacy link integrity (one profile per configuration row that was migrated):
psql "$DATABASE_URL" -c "SELECT count(*) FROM configurations c LEFT JOIN profiles p ON p.legacy_configuration_id = c.id WHERE p.id IS NULL;"
```

Expect `orphan_devices = 0` after `000027`. Unlinked configurations may remain if tenant had no rows to migrate.

## Go tests

```bash
cd serverBackendGo
go test ./internal/modules/device_tree/...
go test ./internal/modules/profiles/...
go test ./internal/modules/enrollment_routes/...
go test ./internal/modules/sync/...
```

## Blueprint gates (§20)

Before each sprint merge, complete Go/No-Go checklist in [DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md](../../DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md) §20.

## Related specs

- [016 quickstart](../016-config-sync-ux/quickstart.md) — policy locks & sync mapper
- [015 spec](../015-device-enrollment-sync/spec.md) — QR enrollment baseline
