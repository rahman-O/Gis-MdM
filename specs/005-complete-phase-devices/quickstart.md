# Quickstart: Phase 4 — Devices & Groups

**Branch**: `005-complete-phase-devices`  
**Prerequisites**: Phases 1–3 complete; Postgres; `make dev`

## 1. Start stack with migrations

```bash
cd serverBackendGo
docker compose down -v   # optional fresh DB
./scripts/db-up.sh
make dev
```

Confirm tables: `devices`, `groups`, `devicegroups`, `configurations`.

## 2. JWT login

```bash
TOKEN=$(curl -s -D - -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | grep -i '^Authorization:' | cut -d' ' -f2-)
```

Swagger: Authorize with `Bearer <token>`.

## 3. Groups smoke

```bash
curl -s http://localhost:8080/rest/private/groups/search \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 4. Configurations list (Devices UI dependency)

```bash
curl -s http://localhost:8080/rest/private/configurations/list \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 5. Device search smoke

```bash
curl -s -X POST http://localhost:8080/rest/private/devices/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"pageNum":1,"pageSize":50,"value":""}' | jq .
```

Expect `data.devices.items` and `data.configurations`.

## 6. Dashboard summary

```bash
curl -s http://localhost:8080/rest/private/summary/devices \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Expect non-empty counters when seed devices exist.

## 7. Create group (settings permission)

```bash
curl -s -X PUT http://localhost:8080/rest/private/groups \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"id":null,"name":"Test Group"}' | jq .
```

## 8. Tests

```bash
go test ./internal/modules/devices/... ./internal/modules/groups/...
```

## 9. Swagger

```bash
make swagger && make dev
```

Tags: **Devices**, **Groups**, **Configurations** (list).

## 10. React E2E

1. Open Devices page — list loads, configuration dropdown populated.
2. Create/edit device, bulk delete, groups admin page.
