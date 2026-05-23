# Quickstart: Phase 5 — Applications, Configurations & Config Files

**Branch**: `006-complete-phase5-apps-config`  
**Prerequisites**: Phases 1–4 complete; Postgres; `make dev`

## 1. Start stack with migrations

```bash
cd serverBackendGo
docker compose down -v   # optional fresh DB
./scripts/db-up.sh
make dev
```

Confirm tables: `applications`, `applicationversions`, `configurationapplications`,
`configurationfiles` (and extended `configurations` columns).

## 2. JWT login

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")
```

Grant org-admin `applications` + `configurations` permissions via migration seed or role UI.

## 3. Configurations smoke

```bash
# List (Configurations page)
curl -s http://localhost:8080/rest/private/configurations/search \
  -H "Authorization: Bearer $TOKEN" | jq .

# Phase 4 regression — devices dropdown
curl -s http://localhost:8080/rest/private/configurations/list \
  -H "Authorization: Bearer $TOKEN" | jq .

# Detail
curl -s http://localhost:8080/rest/private/configurations/1 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 4. Applications smoke

```bash
curl -s http://localhost:8080/rest/private/applications/search \
  -H "Authorization: Bearer $TOKEN" | jq .

curl -s http://localhost:8080/rest/private/applications/1/versions \
  -H "Authorization: Bearer $TOKEN" | jq .
```

## 5. Create configuration (minimal)

```bash
curl -s -X PUT http://localhost:8080/rest/private/configurations \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test Policy","description":"smoke","type":0,"applications":[],"files":[]}' | jq .
```

## 6. Config file upload (optional)

```bash
curl -s -X POST http://localhost:8080/rest/private/config-files \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/small.txt" | jq .
```

Requires `FILES_DIRECTORY` (or equivalent) in `.env` matching Java `files.directory`.

## 7. Tests

```bash
go test ./internal/modules/applications/... ./internal/modules/configurations/... ./internal/modules/configfiles/...
```

## 8. Swagger

```bash
make swagger && make dev
```

Tags: **Applications**, **Configurations**, **ConfigFiles**.

## 9. React E2E

1. **Configurations** — list, create, edit tabs, copy, delete (if no devices assigned).
2. **Applications** — list, add Android/web app, versions page, link configurations dialog.
3. **Devices** — configuration dropdown still loads (`GET /list`).

## 10. Super-admin (optional)

Login as super-admin → Admin Applications page → `GET /private/applications/admin/search`.
