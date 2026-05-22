# Contract: Device install status from sync (014)

**Trigger**: `POST /rest/public/sync/info` (existing sync module)  
**Side effect**: upsert `devicestatuses` for `deviceid`

**Java reference**: `DeviceStatusService`, sync info handler  
**Table**: `devicestatuses` (`000011`, prerequisite 013)

---

## When status updates

After successful `UpdateInfo` for a device (same transaction or immediately after commit):

1. Parse agent payload sections for applications and config files install results.
2. Compute aggregate status per dimension.
3. `INSERT INTO devicestatuses (deviceid, applicationsstatus, configfilesstatus, ...) ON CONFLICT (deviceid) DO UPDATE`.

---

## Status enum values (applications)

Align with Java / device filter `installationStatus`:

| Value | Meaning |
|-------|---------|
| `SUCCESS` | All reported apps OK |
| `FAILURE` | At least one hard failure |
| `VERSION_MISMATCH` | Version mismatch without hard failure |
| (default / empty) | No agent report yet |

**Config files** (`configfilesstatus`): `SUCCESS`, `OTHER`, `FAILURE` — simplified parity with Java aggregates.

---

## Derivation rules (simplified parity)

**Applications** (from `info.applications` or equivalent JSON in sync body):

```
if any status == FAILURE → applicationsstatus = FAILURE
else if any status == VERSION_MISMATCH → applicationsstatus = VERSION_MISMATCH
else if any reported → applicationsstatus = SUCCESS
else → leave unchanged or NULL
```

**Config files**: analogous on file install entries.

Exact field paths in JSON MUST match Java agent payload (verify against `SyncResource` / agent protocol during implement).

---

## Consumers

| Consumer | Use |
|----------|-----|
| `GET /private/devices/search` | filter `installationStatus` (013) |
| Summary `installSummary` | per-config breakdown (013) |

---

## Testing

1. POST sync info with mock failure for one app.
2. `SELECT * FROM devicestatuses WHERE deviceid = ?`
3. Device search with `installationStatus=FAILURE` includes device.

---

## Module design

- **Port**: `DeviceStatusUpserter` in `sync/port`
- **Adapter**: postgres implementation in `sync/adapter/persistence` or reuse `devices` repo if shared
- **No new HTTP routes** in 014
