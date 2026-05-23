#!/usr/bin/env sh
# Configures BASE_URL for DuckDNS + Caddy (Let's Encrypt). Does not start Caddy (needs sudo + open ports).
set -eu
cd "$(dirname "$0")/.."

HOST="${1:-mdm.studhub.duckdns.org}"
BASE="https://${HOST}"
ENV_FILE=".env"

if ! command -v caddy >/dev/null 2>&1; then
  echo "Note: Caddy not installed yet (needed to terminate TLS on this machine)."
  echo "  brew install caddy"
  echo "Continuing — only updating BASE_URL in .env."
  echo ""
fi

if grep -q '^BASE_URL=' "$ENV_FILE" 2>/dev/null; then
  sed -i '' "s|^BASE_URL=.*|BASE_URL=${BASE}|" "$ENV_FILE"
else
  echo "BASE_URL=${BASE}" >>"$ENV_FILE"
fi

echo "Set BASE_URL=${BASE} in .env"
echo ""
echo "1. DuckDNS: point ${HOST} → your public IP"
echo "2. Router: forward ports 80, 443 → this machine"
echo "3. Terminal A: make dev"
echo "4. Terminal B: sudo caddy run --config deploy/caddy/Caddyfile.studhub"
echo "5. ./scripts/verify-enrollment-https.sh ${BASE}"
echo ""
echo "Note: https://studhub.duckdns.org/ may still serve another app until DNS/IP is updated."
