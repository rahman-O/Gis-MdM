# Contract: Repository integration after schema (013)

**Purpose**: Define SQL/repository behavior once migrations `000011`–`000016` are applied. REST paths unchanged (012 owns API parity).

## Devices — `installationStatus` filter

**Before (012 interim)**: `EXISTS (… infojson applications …)` or ignored.

**After (013)**:

```sql
LEFT JOIN devicestatuses ds ON ds.deviceid = d.id
-- when installationStatus = $n:
AND ds.applicationsstatus = $n
```

**Sort `INSTALLATIONS` / `FILES`**: use `COALESCE(ds.applicationsstatus,'FAILURE')` and `COALESCE(ds.configfilesstatus,'OTHER')` in ORDER BY when `sortBy` requests.

**Recalc hook (future)**: `DeviceStatusService` equivalent in `devices/application` or `sync` after `info` update — out of 013 v1 unless trivial INSERT/UPDATE on sync.

---

## Settings — `UserRoleSettings`

**Java**: `UserRoleSettingsDAO.getUserRoleSettings(customerId, roleId)`.

**Go endpoints** (existing scaffold):

- `GET /rest/private/settings/user-role/{roleId}` → full `UserRoleSettings` JSON with all `columnDisplayed*` booleans.
- `PUT` same path → upsert on `(roleid, customerid)` from principal.

**Repository**:

- `GetUserRoleSettings(ctx, customerID, roleID)`
- `SaveUserRoleSettings(ctx, customerID, roleID, domain.UserRoleSettings)`

Default row: if missing, return all `true` (Java behavior) or insert on first save.

---

## Summary — install charts

**Java**: counts by `deviceStatuses.applicationsStatus` / config file status.

**Go** (`summary` module): replace simplified counts with:

```sql
SELECT c.id, c.name, ds.applicationsstatus, COUNT(DISTINCT d.id)
FROM configurations c
JOIN devices d ON d.configurationid = c.id
LEFT JOIN devicestatuses ds ON ds.deviceid = d.id
WHERE d.customerid = $1
GROUP BY c.id, c.name, ds.applicationsstatus
```

---

## Configuration application parameters

**Persistence** (when config editor saves app list):

- Upsert `configurationapplicationparameters` per (configurationId, applicationId) when `skipVersionCheck` flag present in payload (align with Java DTO field name).

---

## Usage stats

**Table only in 013**; **insert** contract owned by `stats` module (012):

- Upsert on conflict `(ts, instanceid)` for daily heartbeat.

---

## Tests required

| Area | Test type |
|------|-----------|
| `devicestatuses` filter | `device_filters_test.go` + integration with JOIN |
| `userrolesettings` | `settings/application` repo stub test |
| migrations | `make migrate` on empty DB + smoke `\d devicestatuses` |
| legacy 000017 | manual on Java dump clone only |
