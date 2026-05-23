# Quickstart: إكمال فجوات قاعدة البيانات (013)

**Branch**: `013-complete-database-gaps`  
**Prerequisites**: Postgres 14, migrations through `000010`, Go backend buildable

**013 baseline (2026-05-21):** `go test ./...` green; `make migrate` applies `000011`–`000017`.

## 1. Apply migrations

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate
```

Expected new objects: `devicestatuses`, `userrolesettings`, `configurationapplicationparameters`, `usagestats`, extended columns on `settings` / `applicationversions` / `configurationapplications`.

Verify:

```bash
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "\dt devicestatuses"
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "\d userrolesettings"
```

## 2. Device status filter (US1)

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")

# Set different application status on two devices
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "
UPDATE devicestatuses SET applicationsstatus = 'SUCCESS'
WHERE deviceid = (SELECT id FROM devices WHERE number = 'hmdm-001');
UPDATE devicestatuses SET applicationsstatus = 'FAILURE'
WHERE deviceid = (SELECT id FROM devices WHERE number = 'hmdm-002');
"

curl -s -X POST http://localhost:8080/rest/private/devices/search \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":50,"installationStatus":"SUCCESS"}' | jq '.data.devices.items[].number'
# Expect: hmdm-001 only (after repo wired)
```

## 3. User role column settings (US2)

```bash
curl -s http://localhost:8080/rest/private/settings/user-role/2 \
  -H "Authorization: Bearer $TOKEN" | jq .
# Expect: columnDisplayed* booleans (not only roleId)

curl -s -X PUT http://localhost:8080/rest/private/settings/user-role/2 \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"roleId":2,"columnDisplayedDeviceImei":false}' | jq .
```

## 4. Schema-only checks (US3)

```bash
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "\d configurationapplicationparameters"
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "\d usagestats"
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c "\d+ applicationversions" | grep apkhash
```

## 5. Legacy import (US4, optional)

Only on DB with Java `configurations` extra columns:

```bash
make migrate   # includes 000017
docker exec serverbackendgo-db psql -U hmdm -d hmdm -c \
  "SELECT id, name, jsonb_object_keys(settingsjson) FROM configurations LIMIT 3;"
```

## 6. Rollback test

```bash
# Down one version at a time using migrate CLI if configured, or restore volume:
docker compose down -v && ./scripts/db-up.sh && make migrate
```

## 7. Tests

```bash
cd serverBackendGo && go test ./internal/modules/devices/... ./internal/modules/settings/...
```

## 8. Update trackers

After UAT pass, mark rows in [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) §3.2 and §7 as ✅.
