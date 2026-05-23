# HTTPS for local enrollment testing

Android **Device Owner** provisioning (QR enrollment) usually **cannot download the admin APK over plain HTTP**.  
If the phone shows **"Couldn't download the admin app"**, the APK URL in the QR is often:

- `http://localhost:...` (phone cannot reach your PC)
- `http://192.168.x.x:8080/...` (blocked or unreliable on recent Android)
- missing / empty (no APK committed — see applications save fix)

Use a **public HTTPS URL** that tunnels to your local Go server.

## Quick setup (automated — recommended)

### Terminal 1 — backend

```bash
cd serverBackendGo
./scripts/db-up.sh
make migrate
make dev
```

### Terminal 2 — HTTPS tunnel + update `.env`

```bash
cd serverBackendGo
make dev-https
```

This installs `cloudflared` if needed, starts the tunnel, sets `BASE_URL=https://....trycloudflare.com` in `.env`, and runs checks.

**Then restart Terminal 1** (`Ctrl+C` → `make dev`) so QR/APK URLs use HTTPS.

```bash
make verify-https   # after restart
```

Stop tunnel: `make dev-https-stop`

## Manual tunnel (optional)

```bash
make tunnel
```

Install tunnel tool if needed:

```bash
brew install cloudflared
```

Copy the URL shown, e.g. `https://random-words.trycloudflare.com` (no trailing `/`).

### Update `.env`

```env
BASE_URL=https://random-words.trycloudflare.com
```

Restart `make dev` in terminal 1.

### Refresh APK + QR URLs

1. **Applications** → open your launcher app → upload APK again → **Save** (writes `https://.../files/...` URL).
2. **Configurations** → MDM → **Main app** = that version → **Save**.
3. **Enrollment QR** → verify JSON/APK:

```bash
KEY=your-qrcodekey
curl -s "https://YOUR-TUNNEL/rest/public/qr/json/${KEY}?create=1&deviceId=test-001" | head -c 600
```

APK line must start with `https://YOUR-TUNNEL/files/...`

```bash
# Replace APK_URL from JSON
curl -sI "APK_URL" | head -3
```

Expect `HTTP/2 200` or `HTTP/1.1 200`.

## Alternatives

| Tool | Install | Command |
|------|---------|---------|
| Cloudflare | `brew install cloudflared` | `./scripts/https-tunnel.sh` |
| ngrok | `brew install ngrok/ngrok/ngrok` | same script (auto-detect) |
| localtunnel | none (npx) | `npx localtunnel --port 8080` |

## Notes

- Tunnel URL **changes** each time you restart cloudflared (unless you use a named ngrok domain).
- Phone does **not** need to be on the same Wi‑Fi when using a public tunnel.
- Admin UI can stay on `http://localhost:5173` (Vite proxy); only **BASE_URL** must be HTTPS for QR/sync from the device.
- Production: use real TLS (reverse proxy + certificate), not a dev tunnel.
