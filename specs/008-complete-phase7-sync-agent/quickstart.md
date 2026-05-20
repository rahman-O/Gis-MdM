# Quickstart: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Branch**: `008-complete-phase7-sync-agent`  
**Prerequisites**: Phases 1–6 complete; Postgres; seeded device `hmdm-001`

## 1. Start stack

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate    # applies 000009_agent_push_notifications
make dev
```

## 2. Environment

```bash
# .env or exports
HASH_SECRET=changeme-C3z9vi54
SECURE_ENROLLMENT=false          # true requires X-Request-Signature on sync/poll
PREVENT_DUPLICATE_ENROLLMENT=false
POLLING_TIMEOUT_MS=60000
MODULE_SYNC_ENABLED=true
MODULE_PUSH_ENABLED=true
MODULE_NOTIFICATIONS_ENABLED=true
MODULE_UPDATES_ENABLED=true
MODULE_QRCODE_ENABLED=true
```

## 3. Sync smoke (public)

```bash
# Configuration sync for seeded device
curl -s "http://localhost:8080/rest/public/sync/configuration/hmdm-001" | jq .

# Device info heartbeat
curl -s -X POST http://localhost:8080/rest/public/sync/info \
  -H 'Content-Type: application/json' \
  -d '{"deviceId":"hmdm-001","batteryLevel":85}' | jq .
```

With `SECURE_ENROLLMENT=true`, compute signature:

```bash
SIG=$(python3 -c "import hashlib; print(hashlib.sha1(('changeme-C3z9vi54'+'hmdm-001').encode()).hexdigest().upper())")
curl -s "http://localhost:8080/rest/public/sync/configuration/hmdm-001" \
  -H "X-Request-Signature: $SIG" | jq .
```

## 4. Push + notifications smoke

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id_token'])")

curl -s -X POST http://localhost:8080/rest/private/push \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"messageType":"configUpdated","payload":"","deviceNumbers":["hmdm-001"]}' | jq .

curl -s "http://localhost:8080/rest/notifications/device/hmdm-001" | jq .
```

## 5. QR smoke

Replace `YOUR_QR_KEY` with `qrCodeKey` from configurations table.

```bash
curl -s -o /tmp/qr.png -w "%{http_code}\n" \
  "http://localhost:8080/rest/public/qr/YOUR_QR_KEY?size=250"

curl -s "http://localhost:8080/rest/public/qr/json/YOUR_QR_KEY" | head -c 400
```

## 6. Updates smoke

```bash
curl -s http://localhost:8080/rest/private/update/check \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Requires super-admin in multi-tenant mode (seed admin qualifies).

## 7. React verification

```bash
cd ../frontend && npm run dev
```

- Open **Configurations** → QR enrollment for a configuration with `qrCodeKey`.
- **Devices** → send push (if UI wired).
- **Updates** page → check for updates (network to vendor manifest).

## 8. Tests

```bash
cd serverBackendGo
go test ./internal/modules/sync/... ./internal/modules/notifications/... \
  ./internal/modules/push/... ./internal/modules/updates/... \
  ./internal/modules/qrcode/... ./internal/shared/crypto/...
make swagger
```
