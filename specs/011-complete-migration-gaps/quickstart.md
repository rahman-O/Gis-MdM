# Quickstart: Phase 9 — Complete Migration Gaps

**Branch**: `011-complete-migration-gaps`  
**Prerequisites**: Phases 1–8 complete; Postgres; seeded device `hmdm-001`; [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md)

## 1. Start stack

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate    # applies 000011_* when added
make dev
```

## 2. Environment (Phase 9)

```bash
MODULE_PUSH_NOTIFIER_ENABLED=true
PUSH_SCHEDULE_INTERVAL_SEC=60
MODULE_PLUGINS_PUSH_ENABLED=true
MODULE_STATS_ENABLED=true
MODULE_VIDEOS_ENABLED=true
VIDEO_DIRECTORY=./data/videos
HASH_SECRET=changeme-C3z9vi54
```

## 3. P0 — Push notifier smoke

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")

# Save configuration (triggers configUpdated enqueue)
curl -s -X PUT http://localhost:8080/rest/private/configurations \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d @fixtures/config-minimal.json | jq .

# Agent poll sees pending message
curl -s "http://localhost:8080/rest/notifications/device/hmdm-001" | jq .
```

Verify DB:

```sql
SELECT id, messagetype, deviceid FROM pushmessages ORDER BY id DESC LIMIT 5;
```

## 4. P0 — Schedule worker smoke

```bash
# Create task due in ~90 seconds
curl -s -X PUT http://localhost:8080/rest/plugins/push/private/task \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"messageType":"textMessage","payload":"scheduled","scope":"device","deviceNumber":"hmdm-001","scheduledTime":'$(($(date +%s)*1000+90000))'}' | jq .

# Wait 2 minutes, then poll notifications again
sleep 120
curl -s "http://localhost:8080/rest/notifications/device/hmdm-001" | jq .
```

## 5. P1 — Icon files

```bash
curl -s -X POST http://localhost:8080/rest/private/icon-files \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@fixtures/icon-square-256.png" | jq .
```

## 6. P1 — Plugin exports (after implementation)

```bash
curl -s -X POST http://localhost:8080/rest/plugins/deviceinfo/deviceinfo/private/export \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"deviceNumber":"hmdm-001"}' -o /tmp/deviceinfo-export.csv

curl -s -X POST http://localhost:8080/rest/plugins/devicelog/log/private/search/export \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"deviceNumber":"hmdm-001"}' -o /tmp/devicelog-export.csv
```

## 7. Regression

```bash
cd serverBackendGo
go test ./...
# Prior phase smoke (auth, plugins platform, sync)
bash scripts/dev.sh smoke 2>/dev/null || true
```

## 8. Update gap tracker

After each subsection passes, edit `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` §4–§6 and matching `docs/parity/*.md`.
