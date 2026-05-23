# Quickstart: إكمال تكامل React ↔ Go (014)

**Branch**: `014-complete-frontend-go-integration`  
**Prerequisite**: Feature **013** applied — DB at migration `000017` or later.

---

## 1. Environment

```bash
cd serverBackendGo
docker compose up -d   # if not running
./scripts/db-up.sh
make migrate           # must include 000011–000017
make dev               # API :8080
```

```bash
cd frontend
npm install
npm run dev            # proxy to Go (see vite config)
```

Login: org-admin (seed from dev) — JWT or session per `.env`.

---

## 2. Smoke — Wave A (MVP P1)

### US1 — Settings tenant fields

1. Open **Settings** → tab tenant / misc fields.
2. Set `phoneNumberFormat` to a distinct pattern; set `customPropertyName1` to `Asset Tag`.
3. Save → reload page → values persist.
4. Open **User role device columns** tab → labels show `Asset Tag` for custom1.

**API** (optional):

```bash
TOKEN=... # from login
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/rest/private/settings | jq .
```

Expect JSON includes `phoneNumberFormat`, `customPropertyName1`, …

### US2 — Configuration MDM round-trip

1. Open existing configuration in editor.
2. Change `kioskMode` / `wifi` in policy tabs; for one app enable **skip version check**, **remove**, **long tap**.
3. Save → reopen same id → all values unchanged.

**API**:

```bash
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/rest/private/configurations/1 | jq '.data.kioskMode, .data.applications[0].skipVersionCheck'
```

### US3 — Icons upload

1. **Icons** page → upload PNG → create icon with name.
2. Search list shows icon; preview URL loads.

**API**:

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/icon.png" \
  http://localhost:8080/rest/private/icon-files
```

---

## 3. Smoke — Wave B (P2)

### US4 — Device install status from sync

1. Note device id; call `POST /rest/public/sync/info` with payload including app install success/failure (see [contracts/sync-device-status.md](./contracts/sync-device-status.md)).
2. Query `devicestatuses` or filter devices by `installationStatus` → status updated.

### US5 — Stats

```bash
curl -s -X PUT http://localhost:8080/rest/public/stats \
  -H "Content-Type: application/json" \
  -d '{"instanceId":"dev-1","devicesTotal":10,"devicesOnline":3,"community":true}'
```

Verify row in `usagestats` for today's `ts`.

---

## 4. Smoke — Wave C (P3)

- **Updates**: trigger apply from Updates page → `POST /private/update` returns OK.
- **Hints**: dismiss hint → `POST /private/hints/history` → hint does not reappear.

---

## 5. Regression

```bash
cd serverBackendGo && go test ./...
```

Re-run UAT rows in [FRONTEND-GO-BACKEND-INTEGRATION.md](../../FRONTEND-GO-BACKEND-INTEGRATION.md) §10 after each wave.

---

## 6. Parity docs to update

| Module | Path |
|--------|------|
| settings | `serverBackendGo/docs/parity/settings.md` |
| configurations | `serverBackendGo/docs/parity/configurations.md` |
| icons | `serverBackendGo/docs/parity/icons.md` |
| sync | `serverBackendGo/docs/parity/sync.md` |
| stats | `serverBackendGo/docs/parity/stats.md` |
