# DuckDNS + free HTTPS (Let's Encrypt) for MDM

**DuckDNS does not issue SSL certificates.** It only maps a hostname (e.g. `studhub.duckdns.org`) to your public IP.

Certificates come from **Let's Encrypt** (free), usually via **Caddy** or **Certbot** on a machine that:

1. Owns the hostname (DNS A record → your public IP)
2. Answers HTTP challenge on port **80** (and serves HTTPS on **443**)

## Check what the domain points to now

```bash
curl -sI https://studhub.duckdns.org/ | head -5
curl -s https://studhub.duckdns.org/ | head -c 200
```

If you see another app (e.g. an embedding API), either:

- Point DuckDNS to the server where **Go MDM** runs, or
- Create a second hostname in DuckDNS, e.g. `mdm.studhub.duckdns.org` → same IP, and use that only for MDM.

## Option A — Caddy on the MDM host (recommended)

Install [Caddy](https://caddyserver.com/docs/install#mac) (`brew install caddy`).

`Caddyfile` example (MDM backend on `127.0.0.1:8080`):

```caddy
mdm.studhub.duckdns.org {
    reverse_proxy 127.0.0.1:8080
}
```

Run (needs ports 80/443 reachable from the internet — router port-forward to this Mac/VM):

```bash
sudo caddy run --config Caddyfile
```

Caddy obtains and renews the certificate automatically. Then in `serverBackendGo/.env`:

```env
BASE_URL=https://mdm.studhub.duckdns.org
```

Restart `make dev`, re-save application (APK URL), save configuration, regenerate QR.

**Router:** forward TCP **80** and **443** to this machine.

## Option B — Dev tunnel (no router / no DuckDNS change)

If you cannot open ports at home, use a temporary HTTPS URL:

```bash
make tunnel   # cloudflared → https://....trycloudflare.com
```

Set `BASE_URL` to that URL. No certificate management on your side.

## Option C — Cloudflare in front of DuckDNS

1. Add the domain in Cloudflare (DNS).
2. Orange-cloud proxy → free HTTPS edge certificate.
3. Tunnel or origin points to your Go server.

More setup; good for production-like tests.

## After HTTPS works

```bash
curl -sI "https://YOUR-HOST/rest/public/qr/json/YOUR_KEY?create=1&deviceId=test" | head -3
```

APK URL inside JSON must be `https://YOUR-HOST/files/...` with HTTP 200:

```bash
curl -sI "APK_URL_FROM_JSON" | head -3
```
