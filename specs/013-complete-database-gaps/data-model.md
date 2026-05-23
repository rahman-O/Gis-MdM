# Data Model: إكمال فجوات قاعدة البيانات (013)

**Branch**: `013-complete-database-gaps` | **Date**: 2026-05-21  
**Reference**: [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md)

## Migration sequence

| Version | Name | Scope |
|---------|------|--------|
| `000011` | `devicestatuses_core` | Table + backfill + index |
| `000012` | `userrolesettings_core` | Table + seed + unique (roleid, customerid) |
| `000013` | `configuration_application_parameters` | Table + unique (configurationid, applicationid) |
| `000014` | `usagestats_core` | Table + unique (ts, instanceid) |
| `000015` | `settings_columns_extend` | `newdevicegroupid`, `phonenumberformat`, custom property names, multiline/send flags |
| `000016` | `applications_columns_extend` | `applicationversions.apkhash`, `configurationapplications.remove`, `longtap` |
| `000017` | `configurations_legacy_import` | Optional data migration Java columns → `settingsjson` (no-op on pure Go DB) |

---

## New tables

### `devicestatuses`

| Column | Type | Constraints |
|--------|------|-------------|
| deviceid | INT | PK, FK → devices(id) ON DELETE CASCADE |
| configfilesstatus | VARCHAR(100) | nullable |
| applicationsstatus | VARCHAR(100) | nullable |

**Relationship**: 1:1 with `devices`.

**Typical values**: `applicationsstatus` ∈ `SUCCESS`, `FAILURE`, `VERSION_MISMATCH`, … (Java `DeviceApplicationsStatus`); `configfilesstatus` ∈ `SUCCESS`, `OTHER`, …

---

### `userrolesettings`

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL | PK |
| roleid | INT | FK → userroles |
| customerid | INT | FK → customers |
| columndisplayeddevicestatus | BOOLEAN | |
| columndisplayeddevicedate | BOOLEAN | |
| columndisplayeddevicenumber | BOOLEAN | |
| columndisplayeddevicemodel | BOOLEAN | |
| columndisplayeddevicepermissionsstatus | BOOLEAN | |
| columndisplayeddeviceappinstallstatus | BOOLEAN | |
| columndisplayeddeviceconfiguration | BOOLEAN | |
| columndisplayeddeviceimei | BOOLEAN | |
| columndisplayeddevicephone | BOOLEAN | |
| columndisplayeddevicedesc | BOOLEAN | |
| columndisplayeddevicegroup | BOOLEAN | |
| columndisplayedlauncherversion | BOOLEAN | |
| columndisplayeddevicefilesstatus | BOOLEAN | Java 19.04.20 |
| columndisplayedbatterylevel | BOOLEAN | |
| columndisplayeddefaultlauncher | BOOLEAN | |
| columndisplayedcustom1 | BOOLEAN | |
| columndisplayedcustom2 | BOOLEAN | |
| columndisplayedcustom3 | BOOLEAN | |
| columndisplayedmdmmode | BOOLEAN | |
| columndisplayedkioskmode | BOOLEAN | |
| columndisplayedandroidversion | BOOLEAN | |
| columndisplayedenrollmentdate | BOOLEAN | |
| columndisplayedserial | BOOLEAN | |
| columndisplayedpublicip | BOOLEAN | |

**Unique**: `(roleid, customerid)`

---

### `configurationapplicationparameters`

| Column | Type | Constraints |
|--------|------|-------------|
| id | SERIAL | PK |
| configurationid | INT | FK → configurations |
| applicationid | INT | FK → applications |
| skipversioncheck | BOOLEAN | NOT NULL DEFAULT FALSE |

**Unique**: `(configurationid, applicationid)`

---

### `usagestats`

| Column | Type | Constraints |
|--------|------|-------------|
| id | SERIAL | PK |
| ts | DATE | NOT NULL DEFAULT CURRENT_DATE |
| instanceid | VARCHAR(255) | |
| webversion | VARCHAR(255) | |
| community | BOOLEAN | NOT NULL DEFAULT TRUE |
| devicestotal | INT | NOT NULL DEFAULT 0 |
| devicesonline | INT | NOT NULL DEFAULT 0 |
| cputotal | INT | NOT NULL DEFAULT 0 |
| cpuused | INT | NOT NULL DEFAULT 0 |
| ramtotal | INT | NOT NULL DEFAULT 0 |
| ramused | INT | NOT NULL DEFAULT 0 |
| scheme | VARCHAR(255) | |
| arch | VARCHAR(255) | |
| os | VARCHAR(255) | |

**Unique**: `(ts, instanceid)`

---

## Extended existing tables

### `settings` (000015)

| New column | Type | Default |
|------------|------|---------|
| newdevicegroupid | INT | nullable, FK groups optional |
| phonenumberformat | VARCHAR(50) | `'+9 (999) 999-99-99'` |
| custompropertyname1 | VARCHAR(200) | |
| custompropertyname2 | VARCHAR(200) | |
| custompropertyname3 | VARCHAR(200) | |
| custommultiline1 | BOOLEAN | FALSE |
| custommultiline2 | BOOLEAN | FALSE |
| custommultiline3 | BOOLEAN | FALSE |
| customsend1 | BOOLEAN | FALSE |
| customsend2 | BOOLEAN | FALSE |
| customsend3 | BOOLEAN | FALSE |
| desktopheadertemplate | TEXT | |
| senddescription | BOOLEAN | FALSE |

### `applicationversions` (000016)

| Column | Type |
|--------|------|
| apkhash | VARCHAR(100) |

### `configurationapplications` (000016)

| Column | Type | Default |
|--------|------|---------|
| remove | BOOLEAN | FALSE |
| longtap | BOOLEAN | FALSE |

### `configurations` (unchanged structure)

- **`settingsjson` JSONB**: canonical store for MDM policy keys not promoted to dedicated columns.
- **000017**: merges legacy SQL columns into JSON when present (see contract).

---

## Repository touchpoints (post-migration)

| Module | Change |
|--------|--------|
| `devices/adapter/persistence/postgres` | JOIN `devicestatuses`; filter `installationStatus` on `applicationsstatus` |
| `settings/adapter/persistence/postgres` | CRUD `userrolesettings`; extend `settings` SELECT/UPDATE |
| `configurations/adapter/persistence/postgres` | Optional read/write `configurationapplicationparameters` |
| `summary/adapter/persistence/postgres` | Use `devicestatuses` for install-by-config charts |
| `stats` (012) | INSERT/UPSERT `usagestats` |

---

## Entities (domain)

- **DeviceStatus** — maps 1:1 to `devicestatuses` row.
- **UserRoleSettings** — full column flags; maps to `userrolesettings`.
- **ConfigurationApplicationParameter** — `SkipVersionCheck bool`.
- **UsageStats** — server telemetry snapshot.
- **Settings** — extended tenant settings fields.
