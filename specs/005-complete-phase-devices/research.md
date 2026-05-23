# Research: Phase 4 — Devices & Groups

**Date**: 2026-05-20

## R1 — Migration strategy

**Decision**: Add `000006_devices_groups_core.up.sql` with Liquibase-aligned subset:
`groups` (customerid, name), `devices` (customerid, number, configurationid, lastupdate,
description, info/infojson columns as TEXT), `devicegroups` (deviceid, groupid),
`configurations` (id, name, customerid, permissive, mainappid nullable).

**Rationale**: `000001_init` has no devices table; `summary` repo already probes `devices`
existence. Full Liquibase import is too large for one PR.

**Alternatives considered**: Edit `000001_init` — rejected for deployed DBs.

## R2 — Device search request field names

**Decision**: Accept React `pageNum`, `pageSize`, `value`, `groupId`, `configurationId`, `status`,
`sortBy`, `sortDir`, and optional filters from `DeviceSearchRequest` in
`frontend/src/features/devices/types.ts`.

**Rationale**: Java `DeviceSearchRequest` uses `pageNum` not `currentPage` (unlike customers).

## R3 — Device list response shape

**Decision**: Return JSON matching `DeviceListResponse`:

```json
{
  "configurations": { "<id>": { "id", "name", "permissiveMode", ... } },
  "devices": { "items": [...], "totalItemsCount": N }
}
```

Wrapped in Headwind envelope `data`.

**Rationale**: React `getDevices()` expects this structure, not flat `PaginatedData`.

## R4 — Group–device relationship

**Decision**: Use `devicegroups` junction (current schema), not deprecated `devices.groupId`
column.

**Rationale**: Liquibase moved to M:N; `groupBulk` mutates `device.groups` list in Java.

## R5 — User group access on search

**Decision**: Mirror `allowedDevicesSelect` / `userDeviceGroupsAccess` — users with
`allDevicesAvailable=false` only see devices in assigned groups.

**Rationale**: `DeviceMapper.xml` INNER JOIN access; security requirement.

## R6 — Configuration enrichment in search (partial)

**Decision**: v1 returns configuration id/name/permissive in list; omit embedded applications
and files in configuration objects unless cheap join; document **Partial** in parity.

**Rationale**: Java loads apps/files per config in loop — heavy; unblock list UI first.

## R7 — Push notifications

**Decision**: `port.PushNotifier` no-op implementation; `notify` endpoints return `OK`.

**Rationale**: Phase 7 owns real push; Phase 4 needs API parity only.

## R8 — Permissions

**Decision**: Add `platformauth` helpers `HasEditDevices`, `HasEditDeviceDesc`; reuse
`HasPermission("settings")` for groups mutations.

**Rationale**: Java checks `edit_devices`, `edit_device_desc`, `settings` by name.

## R9 — Configurations list

**Decision**: Implement `GET /rest/private/configurations/list` in `configurations` module
returning `{ id, name }[]` for `principal.CustomerID` only.

**Rationale**: React devices page blocks without it; full ConfigurationResource is Phase 5.

## R10 — Summary statistics

**Decision**: Replace `EmptyDeviceStats()` stub in `summary_repo.go` with SQL ported from
`DeviceDAO` / `DeviceMapper.countAllDevicesForSummary` (online/offline/enrollment buckets).

**Rationale**: Spec US8; table exists after migration.

## R11 — Device limit on create

**Decision**: Check `settings.deviceLimit` / device count before insert (mirror Java
`settings.getDeviceLimit()`).

**Rationale**: `CustomerResource.insertCustomer` deferred devices; tenant admins hit limit on create.

## R12 — Error envelopes

**Decision**: Use existing httpx helpers; device duplicate → legacy device-exists message;
group not empty → `error.notempty.group`.

**Rationale**: React asserts on `status` + `message` keys.
