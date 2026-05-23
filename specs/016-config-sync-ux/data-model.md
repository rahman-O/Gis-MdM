# Data Model: Configuration–Device Sync & Admin UX

**Feature**: `016-config-sync-ux` | **Date**: 2026-05-23

## Existing tables (unchanged schema v1)

| Table | Role |
|-------|------|
| `configurations` | Policy profile: scalars + `settingsjson` (MDM policy map) |
| `configurationapplications` | Apps linked to configuration (`action`, `remove`, `screenorder`, …) |
| `configurationapplicationparameters` | e.g. `skipversioncheck` per app |
| `configurationapplicationsettings` | Default app settings for configuration (`readonly`, `variable`, `value`) |
| `configurationfiles` | Files pushed with configuration |
| `devices` | `configurationid` links device to profile |
| `deviceapplicationsettings` | Per-device overrides from agent |

## Logical entities

### Configuration (aggregate)

- **Identity**: `id`, `customerid`, `name`, `description`, `type`
- **Enrollment**: `qrcodekey`, `baseurl`, `mainappid`, `contentappid`, `eventreceivingcomponent` (in policy JSON)
- **Design**: `password`, `backgroundcolor`, `textcolor`, `backgroundimageurl`
- **Policy map** (`settingsjson`): kiosk, connectivity, restrictions string, update modes, `pushOptions`, `requestUpdates`, …
- **Policy locks** (`settingsjson.policyLocks`): `Record<fieldKey, boolean>` — when `true`, device cannot override that policy field
- **Children**: `applications[]`, `files[]`, `applicationSettings[]`

### ConfigurationApplication (link)

- `applicationid`, `applicationversionid` (used version), `action`, UI flags (`showicon`, `screenorder`, `remove`, `longtap`, …)

### ConfigurationApplicationSetting

- `applicationpkg` / app reference, `name`, `type`, `value`, **`readonly`**, `comment`, `variable`
- Sync rule: rows with `readonly=true` are authoritative on device; agent updates rejected or merged back on read

### Device

- `number`, `configurationid`, `lastupdate`, `infojson`, custom fields
- Sync reads configuration by `configurationid`

### SyncResponse (agent DTO)

- Superset of current Go struct aligned with Java `SyncResponse`: device identity, design, **full MDM policy fields**, `applications[]`, `files[]`, `applicationSettings[]` (merged config defaults + device overrides respecting readonly)

## State transitions

```text
[Admin edits Configuration]
    → PUT /private/configurations
    → DB updated (columns + settingsjson + children)
    → Push notify configUpdated (devices on configuration)

[Device polls/sync]
    → GET|POST /public/sync/configuration/{deviceId}
    → BuildSyncResponse(configuration + device overrides)
    → Agent applies policy

[Device posts app settings]
    → POST /public/sync/applicationSettings/{deviceId}
    → Reject/skip keys matching readonly configuration defaults
    → Upsert allowed keys to deviceapplicationsettings
```

## Validation rules

| Rule | Enforcement |
|------|-------------|
| `name` required | configurations handler |
| `mainAppId` required for QR eligibility | frontend + QR service (existing) |
| `policyLocks` keys must be known MDM fields | domain allowlist in configurations module |
| Application version must exist for linked app | save configuration / upgrade endpoint |
| Readonly config setting cannot be changed by device POST | sync application service |
| Tenant scope | all queries filter by `customerid` from principal |

## Optional migration (Phase B — if JSON locks grow large)

- `000019_configuration_policy_locks.up.sql` — only if `policyLocks` in JSON proves insufficient; **not required for v1**.
