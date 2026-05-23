# Quickstart: Phase 2 Users & Roles verification

## Prerequisites

```bash
cd serverBackendGo
cp .env.example .env
./scripts/db-up.sh
make dev
```

Swagger: http://localhost:8080/swagger/index.html (`SWAGGER_ENABLED=true`)

Frontend (optional E2E):

```bash
cd ../frontend && npm run dev
# Vite proxies /rest → http://localhost:8080
```

## 1. Authenticate

**Session (browser / Swagger cookie)**:

```bash
curl -s -c /tmp/hmdm.jar -X POST http://localhost:8080/rest/public/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}'
```

**JWT**:

```bash
TOKEN=$(curl -s -D - -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | grep -i '^Authorization:' | awk '{print $2}')
export AUTH="Authorization: Bearer $TOKEN"
```

Use `-b /tmp/hmdm.jar` for session or `-H "$AUTH"` for JWT in steps below.

## 2. Users smoke

```bash
# Current user (Profile / session refresh)
curl -s -b /tmp/hmdm.jar http://localhost:8080/rest/private/users/current | jq .

# List users (admin)
curl -s -b /tmp/hmdm.jar http://localhost:8080/rest/private/users/all | jq .

# Update profile
curl -s -b /tmp/hmdm.jar -X PUT http://localhost:8080/rest/private/users/details \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"name":"Administrator","email":"admin@localhost"}' | jq .

# Change password (use MD5 hex in real clients; "admin" works on Go normalize)
curl -s -b /tmp/hmdm.jar -X PUT http://localhost:8080/rest/private/users/current \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"login":"admin","oldPassword":"21232F297A57A5A743894A0E4A801FC3","newPassword":"21232F297A57A5A743894A0E4A801FC3"}' | jq .

# Role dropdown
curl -s -b /tmp/hmdm.jar http://localhost:8080/rest/private/users/roles | jq .
```

## 3. Roles smoke

```bash
curl -s -b /tmp/hmdm.jar http://localhost:8080/rest/private/roles/permissions | jq .
curl -s -b /tmp/hmdm.jar http://localhost:8080/rest/private/roles/all | jq .
```

## 4. Automated tests

```bash
cd serverBackendGo
go test ./internal/modules/users/... ./internal/modules/roles/... ./internal/platform/auth/... -v
go build ./...
```

## 5. React E2E

1. Login as `admin` / `admin`
2. Open **Profile** — no failed `/private/users/current` or `/details`
3. Open **Users** — list loads; create/edit/delete (if permitted)
4. Open **Roles** (if in nav) — permissions and roles load

## 6. Regenerate Swagger after handler changes

```bash
make swagger
# Reload http://localhost:8080/swagger/index.html
```
