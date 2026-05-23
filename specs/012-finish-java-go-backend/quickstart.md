# Quickstart: إكمال نقل الباكند Java → Go (012)

**Branch**: `012-finish-java-go-backend`  
**Prerequisites**: Phases 1–8; Phase 9 P0 (push, cron, icon-files) working; Postgres seeded

**References**: [JAVA-GO-BACKEND-GAPS.md](../../JAVA-GO-BACKEND-GAPS.md), [FRONTEND-GO-BACKEND-INTEGRATION.md](../../FRONTEND-GO-BACKEND-INTEGRATION.md)

## 1. Start stack

```bash
cd serverBackendGo
docker compose down -v   # optional fresh DB
./scripts/db-up.sh
make migrate
make dev
```

Ensure `.env`:

```env
FILES_DIRECTORY=./data/files
MODULE_PUSH_NOTIFIER_ENABLED=true
PUSH_SCHEDULE_INTERVAL_SEC=60
MODULE_STATS_ENABLED=true
MODULE_VIDEOS_ENABLED=false
```

## 2. Obtain token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")
```

## 3. P0 regression (from 011)

```bash
# Push on config save
curl -s -X PUT http://localhost:8080/rest/private/configurations \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"id":1,"name":"Default","type":0}' | jq .

curl -s "http://localhost:8080/rest/notifications/device/hmdm-001" | jq .
```

## 4. P1 — Device search with filters

```bash
curl -s -X POST http://localhost:8080/rest/private/devices/search \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":20,"status":"green","sortBy":"LAST_UPDATE","sortDir":"desc"}' | jq .

curl -s http://localhost:8080/rest/private/devices/number/hmdm-001 \
  -H "Authorization: Bearer $TOKEN" | jq '.data.info'
```

## 5. P1 — Plugin exports

```bash
curl -s -X POST http://localhost:8080/rest/plugins/deviceinfo/deviceinfo/private/export \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"deviceNumber":"hmdm-001"}' -o /tmp/deviceinfo.csv

curl -s http://localhost:8080/rest/plugins/devicelog/log/rules/hmdm-001 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 6. P2 — Audit + static files

```bash
# Delete device → audit row
curl -s -X DELETE http://localhost:8080/rest/private/devices/999 \
  -H "Authorization: Bearer $TOKEN" | jq .

curl -s -X POST http://localhost:8080/rest/plugins/audit/private/log/search \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":10}' | jq .

# Static file (replace path with real uploaded file)
curl -s -o /tmp/agent.apk "http://localhost:8080/files/CUSTOMER_DIR/path/app.apk"
```

## 7. P3 — Stats & updates

```bash
curl -s -X PUT http://localhost:8080/rest/public/stats \
  -H 'Content-Type: application/json' \
  -d '{"customerId":1,"deviceId":1,"event":"test"}' | jq .

curl -s http://localhost:8080/rest/private/update/check \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 8. Full UAT (30 min checklist)

1. React login → dashboard counts  
2. Devices: apply filters, open detail panel (battery/model)  
3. Save configuration → agent notification  
4. Plugin settings page loads  
5. Optional: scheduled push task fires  
6. Stop Java WAR; repeat 1–4 on Go only  

## 9. Tests

```bash
cd serverBackendGo && go test ./...
```

Update gap trackers when all sections pass:

- `JAVA-GO-BACKEND-GAPS.md`
- `JAVA-GO-MIGRATION-STATUS.md`
- `serverBackendGo/docs/MIGRATION.md` Phase 9 → **done**
