# Quickstart: Phase 6 — Files, Icons & Public API

**Branch**: `007-complete-phase6-files-public`  
**Prerequisites**: Phases 1–5 complete; Postgres; `FILES_DIRECTORY` writable (`./data/files`)

## 1. Start stack with migrations

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
```

Confirm tables: `uploadedfiles`, `icons`; `configurationfiles.fileid` column.

## 2. Environment

Ensure `.env` includes:

```bash
FILES_DIRECTORY=./data/files
BASE_URL=http://localhost:8080
HASH_SECRET=changeme-C3z9vi54
MODULE_FILES_ENABLED=true
MODULE_ICONS_ENABLED=true
MODULE_PUBLICAPI_ENABLED=true
```

## 3. JWT login

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")
```

Grant org-admin `files` and `edit_files` via migration seed or role UI.

## 4. Files smoke

```bash
# List
curl -s http://localhost:8080/rest/private/web-ui-files/search \
  -H "Authorization: Bearer $TOKEN" | jq .

# Storage limit
curl -s http://localhost:8080/rest/private/web-ui-files/limit \
  -H "Authorization: Bearer $TOKEN" | jq .

# Multipart upload (small file)
curl -s -X POST http://localhost:8080/rest/private/web-ui-files/raw \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/etc/hosts" | jq .
```

Use returned `serverPath` / `name` in a follow-up `POST /update` body to commit (see Java flow).

## 5. Icons smoke

```bash
curl -s http://localhost:8080/rest/private/icons/search \
  -H "Authorization: Bearer $TOKEN" | jq .

curl -s -X PUT http://localhost:8080/rest/private/icons \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test Icon","fileId":1}' | jq .
```

## 6. Public API smoke

```bash
curl -s http://localhost:8080/rest/public/name | jq .

curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/rest/public/logo
```

## 7. React verification

```bash
cd ../frontend && npm run dev
```

Login → **Files** page loads without 404 on `/private/web-ui-files/search`.  
**Applications** → upload APK uses `/private/web-ui-files` endpoints.

## 8. Tests

```bash
cd serverBackendGo
go test ./internal/modules/files/... ./internal/modules/icons/... ./internal/modules/publicapi/... ./internal/platform/storage/...
```

## 9. Swagger

```bash
make swagger
# Open http://localhost:8080/swagger/index.html — tags Files, Icons, Public API
```
