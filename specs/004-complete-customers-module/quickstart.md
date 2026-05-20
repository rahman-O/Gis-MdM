# Quickstart: Phase 3 — Customers Module

**Branch**: `004-complete-customers-module`  
**Prerequisites**: Postgres up, migrations through `000005`, server on `:8080`, super-admin seed user

## 1. Start stack

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
```

Verify `customers` has extended columns after `000005` migration.

## 2. Obtain super-admin JWT

**Swagger**: `POST /public/jwt/login` with seed admin → **Authorize** with Bearer token.

**curl**:

```bash
TOKEN=$(curl -s -D - -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | grep -i '^Authorization:' | cut -d' ' -f2-)
```

## 3. Smoke — customer search (P1)

```bash
curl -s -X POST http://localhost:8080/rest/private/customers/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"currentPage":1,"pageSize":100,"searchValue":""}' | jq .
```

Expect: `status: OK`, `data.items` array, `data.totalItemsCount` number.

## 4. Smoke — impersonate (P1)

Pick a customer id from search (e.g. `2`):

```bash
curl -s http://localhost:8080/rest/private/customers/impersonate/2 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Expect: `status: OK`, `data.authToken` present, no `password` field.

## 5. Smoke — prefix check (P3)

```bash
curl -s http://localhost:8080/rest/private/customers/prefix/hmdm-/used \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Expect: `data: true` or `false`.

## 6. React control panel

1. Open frontend with Go proxy (`:8080`).
2. Log in as super admin.
3. Navigate to Control Panel.
4. Confirm customer list loads and **Impersonate** switches session.

## 7. Tests

```bash
cd serverBackendGo
go test ./internal/modules/customers/...
```

## 8. Parity sign-off

Update `serverBackendGo/docs/parity/customers.md` — mark P1 endpoints **Done**, note Partial on create side effects until Phase 4.
