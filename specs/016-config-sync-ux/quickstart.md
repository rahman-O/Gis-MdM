# Quickstart: Configuration–Device Sync & Admin UX (016)

## Prerequisites

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate
make dev
```

- Stable `BASE_URL` (e.g. `https://mdm.studhub.app`) in `.env`
- `cloudflared` tunnel → `http://127.0.0.1:8080` if testing from physical device
- JWT login (admin with `configurations` permission)

```bash
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/rest/public/jwt/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' | jq -r '.data.token // .data.accessToken')
```

---

## 1. Configuration save round-trip

```bash
curl -s http://127.0.0.1:8080/rest/private/configurations/2 \
  -H "Authorization: Bearer $TOKEN" | jq '.data | {mainAppId, kioskMode, restrictions, policyLocks}'
```

Edit in UI: enable **kiosk**, set **restrictions**, lock **main app** → Save.

Re-run GET — values must match.

---

## 2. Sync payload includes policy

Replace `DEVICE` with enrolled device number:

```bash
curl -s "http://127.0.0.1:8080/rest/public/sync/configuration/DEVICE" \
  -H "X-CPU-Arch: arm64-v8a" | jq '{kioskMode, restrictions, applications: (.applications|length)}'
```

Expect `kioskMode` / `restrictions` present when set in configuration.

---

## 3. Readonly app setting

1. Configuration → App settings tab → set a value, check **Readonly**.
2. Save configuration.
3. Sync GET — setting appears with `"readonly": true`.
4. Simulate device POST with changed value — server keeps config value on next sync.

---

## 4. Push notify (optional)

With `MODULE_PUSH_ENABLED=true`, after configuration save check notifications endpoint or device log for `configUpdated`.

---

## 5. Frontend UX smoke

1. Open `/configurations/2`.
2. Walk all tabs — no duplicate MDM section on General tab.
3. Save with invalid name — error shows tab + field.

---

## 6. Java parity spot-check (developers)

Compare Go vs Java sync JSON for same DB seed (see `contracts/sync-configuration-payload.md` S5).

---

## Done when

- SC-001..SC-005 from [spec.md](./spec.md) pass on staging
- `docs/parity/configurations.md` and `sync.md` updated
