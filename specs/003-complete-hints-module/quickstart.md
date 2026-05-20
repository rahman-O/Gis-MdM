# Quickstart: Phase 3 — Hints Module

**Branch**: `003-complete-hints-module`  
**Prerequisites**: Postgres up, migrations applied, server on `:8080`

## 1. Start stack

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
```

Verify tables: `userhints`, `userhinttypes` exist after `000004` migration.

## 2. Obtain JWT (Swagger or curl)

**Swagger**: `POST /public/jwt/login` → copy `Authorization` header → **Authorize**.

**curl**:

```bash
TOKEN=$(curl -s -D - -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | grep -i '^Authorization:' | cut -d' ' -f2-)
echo "Token: $TOKEN"
```

## 3. Smoke tests

### Empty history (new user or after enable)

```bash
curl -s http://localhost:8080/rest/private/hints/history \
  -H "Authorization: Bearer $TOKEN" | jq .
# Expect: {"status":"OK","data":[]}
```

### Mark hint shown

```bash
curl -s -X POST http://localhost:8080/rest/private/hints/history \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '"hint.step.1"' | jq .
```

```bash
curl -s http://localhost:8080/rest/private/hints/history \
  -H "Authorization: Bearer $TOKEN" | jq .
# Expect data contains "hint.step.1"
```

### Disable (all catalog keys)

```bash
curl -s -X POST http://localhost:8080/rest/private/hints/disable \
  -H "Authorization: Bearer $TOKEN" | jq .
```

```bash
curl -s http://localhost:8080/rest/private/hints/history \
  -H "Authorization: Bearer $TOKEN" | jq .
# Expect 4 keys: hint.step.1 .. hint.step.4
```

### Enable (clear history)

```bash
curl -s -X POST http://localhost:8080/rest/private/hints/enable \
  -H "Authorization: Bearer $TOKEN" | jq .
```

```bash
curl -s http://localhost:8080/rest/private/hints/history \
  -H "Authorization: Bearer $TOKEN" | jq .
# Expect: []
```

## 4. Automated tests

```bash
cd serverBackendGo
go test ./internal/modules/hints/... -v
go build ./...
```

## 5. React UI

1. Vite proxy to `http://localhost:8080`
2. Login as `admin`
3. Open **Hints** (`/hints`)
4. Enable / Disable / refresh list — no 403/404

## 6. Swagger

Open `http://localhost:8080/swagger/index.html` → Authorize → test all four `/private/hints/*` operations.

## Sign-off checklist

- [ ] All four endpoints return Headwind envelope
- [ ] Unauthenticated calls return 403
- [ ] `docs/parity/hints.md` all Done
- [ ] `go test ./internal/modules/hints/...` passes
