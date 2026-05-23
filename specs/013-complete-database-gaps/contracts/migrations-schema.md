# Contract: Database migrations (013)

**Tool**: golang-migrate  
**Path**: `serverBackendGo/db/migrations/`  
**Baseline**: through `000010_plugins_core` + `000008_devices_search_extras`

## Rules (all migrations)

1. Pair every `NNNNNN_name.up.sql` with `NNNNNN_name.down.sql`.
2. Use lowercase unquoted identifiers (PostgreSQL folds to lowercase).
3. Prefer `IF NOT EXISTS` / `IF EXISTS` for idempotent deploys.
4. Seeds/backfill only in `up`; reversible in `down` where safe (DROP TABLE for new tables).
5. No business logic beyond INSERT defaults/backfill.

## 000011 — `devicestatuses_core`

**Up**:

```sql
CREATE TABLE IF NOT EXISTS devicestatuses (
    deviceid INT NOT NULL PRIMARY KEY REFERENCES devices (id) ON DELETE CASCADE,
    configfilesstatus VARCHAR(100),
    applicationsstatus VARCHAR(100)
);
CREATE INDEX IF NOT EXISTS devicestatuses_apps_status_idx ON devicestatuses (applicationsstatus);
INSERT INTO devicestatuses (deviceid, configfilesstatus, applicationsstatus)
SELECT id, 'OTHER', 'FAILURE' FROM devices d
WHERE NOT EXISTS (SELECT 1 FROM devicestatuses ds WHERE ds.deviceid = d.id);
```

**Down**: `DROP TABLE IF EXISTS devicestatuses;`

**Consumers**: `devices` search SQL, `summary` aggregates.

---

## 000012 — `userrolesettings_core`

**Up**: `CREATE TABLE userrolesettings (...)` — all `columndisplayed*` columns per [data-model.md](../data-model.md); `UNIQUE (roleid, customerid)`; seed rows for roles 1–3 × each customer (columns default TRUE).

**Down**: `DROP TABLE IF EXISTS userrolesettings;`

**Consumers**: `GET/PUT` settings role endpoints (Java `UserRoleSettingsResource` parity).

---

## 000013 — `configuration_application_parameters`

**Up**:

```sql
CREATE TABLE IF NOT EXISTS configurationapplicationparameters (
    id SERIAL PRIMARY KEY,
    configurationid INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    applicationid INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    skipversioncheck BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE UNIQUE INDEX IF NOT EXISTS cap_config_app_uidx
    ON configurationapplicationparameters (configurationid, applicationid);
```

**Down**: `DROP TABLE IF EXISTS configurationapplicationparameters;`

---

## 000014 — `usagestats_core`

**Up**: table per Java `07.03.23-11:23`; `UNIQUE (ts, instanceid)`.

**Down**: `DROP TABLE IF EXISTS usagestats;`

**Consumers**: `stats` module (012), `PUT /rest/public/stats`.

---

## 000015 — `settings_columns_extend`

**Up**: `ALTER TABLE settings ADD COLUMN IF NOT EXISTS …` (see data-model).

**Down**: `ALTER TABLE settings DROP COLUMN IF EXISTS …` (reverse order).

---

## 000016 — `applications_columns_extend`

**Up**:

- `ALTER TABLE applicationversions ADD COLUMN IF NOT EXISTS apkhash VARCHAR(100);`
- `ALTER TABLE configurationapplications ADD COLUMN IF NOT EXISTS remove BOOLEAN NOT NULL DEFAULT FALSE;`
- `ALTER TABLE configurationapplications ADD COLUMN IF NOT EXISTS longtap BOOLEAN NOT NULL DEFAULT FALSE;`

**Down**: drop columns.

---

## 000017 — `configurations_legacy_import`

**Up**: PL/pgSQL block — **only if** legacy columns exist on `configurations` (from Java dump):

- Detect via `information_schema.columns`.
- `UPDATE configurations SET settingsjson = settingsjson || jsonb_build_object(...)` mapping documented in [legacy-config-import.md](./legacy-config-import.md).

**Down**: no-op or remove merged keys (document: non-reversible without backup).

---

## Verification

```bash
cd serverBackendGo
make migrate
psql "$DATABASE_URL" -c "\dt devicestatuses"
psql "$DATABASE_URL" -c "\d userrolesettings"
```

See [quickstart.md](../quickstart.md).
