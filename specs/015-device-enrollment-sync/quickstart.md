# Quickstart: Device Enrollment & Sync UAT

**Branch**: `015-device-enrollment-sync`  
**Prerequisites**: Postgres migrated; configuration with valid `qrCodeKey`, `mainAppId` + APK URL; `BASE_URL` reachable from phone

## 1. Start stack

```bash
cd serverBackendGo
docker compose down -v   # optional clean DB
./scripts/db-up.sh
make migrate
make dev
```

## 2. Environment (critical)

```bash
# Use LAN IP reachable from Android — NOT localhost on device
export BASE_URL=http://192.168.1.10:8080
export FILES_DIRECTORY=./data/files
export HASH_SECRET=changeme-C3z9vi54
export SECURE_ENROLLMENT=false
export PREVENT_DUPLICATE_ENROLLMENT=false
export MODULE_SYNC_ENABLED=true
export MODULE_QRCODE_ENABLED=true
```

Login (JWT for admin UI): `admin` / `admin` (seed).

## 3. Verify configuration QR eligibility

```bash
psql "$DATABASE_URL" -c \
  "SELECT id, name, qrcodekey, mainappid, launcherurl FROM configurations LIMIT 5;"
```

Note `qrcodekey` (e.g. `default`) for steps below.

## 4. QR provisioning smoke

```bash
KEY=default   # replace with your qrcodekey
curl -s -o /tmp/qr.png -w "PNG HTTP %{http_code}\n" \
  "${BASE_URL}/rest/public/qr/${KEY}?size=250&create=1&deviceId=test-device-001"

curl -s "${BASE_URL}/rest/public/qr/json/${KEY}?create=1&deviceId=test-device-001" | head -c 800
echo
```

**Pass**: JSON contains `PROVISIONING_DEVICE_ADMIN_PACKAGE_CHECKSUM`, `com.hmdm.BASE_URL`, `com.hmdm.CONFIG`.

**Fail**: Only three lines / missing checksum → QR parity not fixed yet.

## 5. Static files smoke

Extract APK path from QR JSON (`PACKAGE_DOWNLOAD_LOCATION`), then:

```bash
curl -s -o /dev/null -w "APK HTTP %{http_code}\n" "<APK_URL_FROM_QR>"
```

**Pass**: HTTP 200. **Fail**: 404 → implement `/files/*` handler.

## 6. Enrollment sync smoke

```bash
DEVICE=test-device-001
curl -s -X POST "${BASE_URL}/rest/public/sync/configuration/${DEVICE}" \
  -H 'Content-Type: application/json' \
  -d "{\"configuration\":\"${KEY}\",\"groups\":[]}" | jq .

curl -s -X POST "${BASE_URL}/rest/public/sync/info" \
  -H 'Content-Type: application/json' \
  -d "{\"deviceId\":\"${DEVICE}\",\"batteryLevel\":88}" | jq .
```

**Pass**: First call `status=OK` with `data.applications` / `data.files`; device row exists:

```bash
psql "$DATABASE_URL" -c \
  "SELECT id, number, configurationid, lastupdate FROM devices WHERE lower(number)=lower('${DEVICE}');"
```

## 7. Admin UI check

1. Open React → Devices → search `test-device-001`.
2. Device appears with recent activity after `sync/info`.
3. Enrollment QR page: generate QR with **Create on demand** → scan with physical device on same Wi‑Fi as `BASE_URL`.

## 8. Secure enrollment (optional)

```bash
export SECURE_ENROLLMENT=true
SIG=$(python3 -c "import hashlib; print(hashlib.sha1(('changeme-C3z9vi54'+'${DEVICE}').encode()).hexdigest().upper())")
curl -s "${BASE_URL}/rest/public/sync/configuration/${DEVICE}" \
  -H "X-Request-Signature: $SIG" | jq .
```

## 9. Regression checklist (SC-002)

- [ ] Flow A: QR `create=1` → new device in DB  
- [ ] Flow B: Pre-create device in UI → sync same number  
- [ ] Invalid QR key → 500 + log (not silent success)  
- [ ] `PREVENT_DUPLICATE_ENROLLMENT=true` blocks second enroll  
- [ ] Change configuration → push or poll updates agent  

## 10. Sign-off

When all pass, update:

- `serverBackendGo/docs/parity/qrcode.md`
- `serverBackendGo/docs/parity/sync.md`
- New or updated static files section in `docs/parity/files.md`
