#!/usr/bin/env sh
# Free HTTPS tunnel for local MDM testing (Android QR enrollment requires HTTPS APK download).
#
# Usage (two terminals):
#   Terminal 1: make dev          # Go on :8080
#   Terminal 2: ./scripts/https-tunnel.sh
#
# Copy the printed https:// URL into serverBackendGo/.env as BASE_URL (no trailing slash),
# restart make dev, re-save the application (APK URL), then open enrollment QR again.

set -e
cd "$(dirname "$0")/.."

PORT="${1:-${SERVER_PORT:-8080}}"
TARGET="http://127.0.0.1:${PORT}"

echo "==> HTTPS tunnel for Headwind MDM dev"
echo "    Local backend: ${TARGET}"
echo ""
echo "After the tunnel starts:"
echo "  1. Copy the https://.... URL into .env → BASE_URL=https://...."
echo "  2. Restart: make dev"
echo "  3. Applications → re-save app with APK (refreshes download URL)"
echo "  4. Configuration → Save → Enrollment QR"
echo "  5. Scan QR on the phone (same Wi‑Fi or mobile data)"
echo ""

if command -v cloudflared >/dev/null 2>&1; then
  echo "Using Cloudflare Tunnel (free, no account for quick tunnels)."
  exec cloudflared tunnel --url "${TARGET}"
fi

if command -v ngrok >/dev/null 2>&1; then
  echo "Using ngrok."
  exec ngrok http "${PORT}"
fi

echo "No tunnel tool found. Install one of:"
echo "  brew install cloudflared    # recommended"
echo "  brew install ngrok/ngrok/ngrok"
echo ""
echo "Or one-shot without install:"
echo "  npx --yes localtunnel --port ${PORT}"
exit 1
