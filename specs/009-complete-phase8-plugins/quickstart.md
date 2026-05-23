# Quickstart: Phase 8 — Plugins

**Branch**: `009-complete-phase8-plugins`  
**Prerequisites**: Phases 1–7; Postgres; `make migrate` includes `000010_plugins_core`

## 1. Start stack

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate
make dev
```

## 2. Environment

```bash
ENABLED_PLUGINS=audit,push,messaging,deviceinfo,devicelog
MODULE_PLUGINS_ENABLED=true
MODULE_PLUGINS_PLATFORM_ENABLED=true
MODULE_PLUGINS_AUDIT_ENABLED=true
MODULE_PLUGINS_MESSAGING_ENABLED=true
MODULE_PLUGINS_DEVICEINFO_ENABLED=true
MODULE_PLUGINS_DEVICELOG_ENABLED=true
```

## 3. JWT token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")
```

## 4. Platform (React Plugin settings)

```bash
curl -s http://localhost:8080/rest/plugin/main/private/active \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

curl -s http://localhost:8080/rest/plugin/main/private/available \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

# Disable plugin id 3 for tenant (example)
curl -s -X POST http://localhost:8080/rest/plugin/main/private/disabled \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '[3]' | python3 -m json.tool

curl -s http://localhost:8080/rest/plugin/main/public/registered | python3 -m json.tool
```

## 5. Audit search

```bash
curl -s -X POST http://localhost:8080/rest/plugins/audit/private/log/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":20}' | python3 -m json.tool
```

## 6. Messaging

```bash
curl -s -X POST http://localhost:8080/rest/plugins/messaging/private/send \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"message":"hello","deviceNumbers":["hmdm-001"]}' | python3 -m json.tool

curl -s -X POST http://localhost:8080/rest/plugins/messaging/private/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":10}' | python3 -m json.tool

curl -s "http://localhost:8080/rest/notifications/device/hmdm-001" | python3 -m json.tool
```

## 7. Device info

```bash
curl -s -X PUT http://localhost:8080/rest/plugins/deviceinfo/deviceinfo/public/hmdm-001 \
  -H 'Content-Type: application/json' \
  -d '[{"attribute":"battery","value":"90"}]' | python3 -m json.tool

curl -s http://localhost:8080/rest/plugins/deviceinfo/deviceinfo/private/hmdm-001 \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
```

## 8. Device log

```bash
curl -s http://localhost:8080/rest/plugins/devicelog/devicelog-plugin-settings/private \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

curl -s -X POST http://localhost:8080/rest/plugins/devicelog/log/private/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":20}' | python3 -m json.tool
```

## 9. Push schedule tasks (Phase 8 completion)

```bash
curl -s -X POST http://localhost:8080/rest/plugins/push/private/searchTasks \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":10}' | python3 -m json.tool
```

## 10. Unit tests

```bash
cd serverBackendGo
go test ./internal/modules/plugins/... -count=1
go build ./...
```

## 11. React check

```bash
cd ../frontend && npm run dev
# Login → Settings → Plugins tab — no 404 on plugin/main endpoints
```
