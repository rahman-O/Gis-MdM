#!/usr/bin/env sh
# Verifies HTTPS BASE_URL, QR JSON, and APK download URL for Android enrollment.
# Usage: ./scripts/verify-enrollment-https.sh [BASE_URL]
# Default BASE_URL from .env or first configuration qrcodekey from DB.

set -eu
cd "$(dirname "$0")/.."

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi

BASE="${1:-${BASE_URL:-}}"
BASE=$(echo "$BASE" | sed 's|/$||')

if [ -z "$BASE" ]; then
  echo "ERROR: pass BASE_URL or set it in .env"
  exit 1
fi

case "$BASE" in
  http://*)
    echo "ERROR: BASE_URL must be HTTPS for Android QR enrollment: $BASE"
    exit 1
    ;;
esac

KEY="${QR_CODE_KEY:-}"
if [ -z "$KEY" ]; then
  if command -v psql >/dev/null 2>&1 && [ -n "${DATABASE_URL:-}" ]; then
    KEY=$(psql "$DATABASE_URL" -tAc \
      "SELECT qrcodekey FROM configurations WHERE qrcodekey IS NOT NULL AND trim(qrcodekey) <> '' AND mainappid IS NOT NULL ORDER BY id DESC LIMIT 1" 2>/dev/null | tr -d ' \n' || true)
  elif docker compose ps -q db 2>/dev/null | grep -q .; then
    KEY=$(docker compose exec -T db psql -U hmdm -d hmdm -tAc \
      "SELECT qrcodekey FROM configurations WHERE qrcodekey IS NOT NULL AND trim(qrcodekey) <> '' AND mainappid IS NOT NULL ORDER BY id DESC LIMIT 1" 2>/dev/null | tr -d ' \n\r' || true)
  fi
fi
KEY="${KEY:-default-qr}"

echo "BASE_URL: $BASE"
echo "QR key:   $KEY"

printf '%s' "Public name endpoint... "
code=$(curl -sk -o /dev/null -w "%{http_code}" --max-time 20 "${BASE}/rest/public/name" 2>/dev/null) || code=000
code=$(echo "$code" | tr -cd '0-9')
if [ "$code" = "200" ]; then
  printf 'OK (%s)\n' "$code"
else
  printf 'FAIL (HTTP %s)\n' "$code"
  echo "  - Is 'make dev' running? (only one instance — port 8080 must not be 'address already in use')"
  echo "  - Is the tunnel up? In another terminal: make dev-https (leave it running; do not Ctrl+C)"
  if [ "$code" = "530" ] || [ "$code" = "000" ]; then
    echo "  - HTTP 530/000 usually means BASE_URL in .env is stale (tunnel stopped). Run: make dev-https"
  fi
  echo "  - If you are already in serverBackendGo/, do not run 'cd serverBackendGo' again."
  exit 1
fi

json_path="${BASE}/rest/public/qr/json/${KEY}?create=1&deviceId=mdm-verify-001"
echo "Fetching QR JSON..."
json=$(curl -sk --max-time 30 "$json_path" || true)
if [ -z "$json" ]; then
  echo "ERROR: empty QR JSON response"
  exit 1
fi

if ! echo "$json" | grep -q 'PROVISIONING_DEVICE_ADMIN_PACKAGE_DOWNLOAD_LOCATION'; then
  echo "ERROR: QR JSON missing APK download field. Server log may show main app URL missing."
  echo "$json" | head -c 400
  echo ""
  exit 1
fi

apk_url=$(echo "$json" | sed -n 's/.*PROVISIONING_DEVICE_ADMIN_PACKAGE_DOWNLOAD_LOCATION.:\"\([^\"]*\)\".*/\1/p' | head -1)
if [ -z "$apk_url" ]; then
  apk_url=$(echo "$json" | grep -oE 'https://[^\"]+' | head -1)
fi

echo "APK URL: $apk_url"

case "$apk_url" in
  https://*) echo "APK URL scheme: OK" ;;
  *)
    echo "ERROR: APK URL must be https:// — got: $apk_url"
    echo "Re-save application after restarting server with updated BASE_URL."
    exit 1
    ;;
esac

printf '%s' "APK download... "
# Range probe: full APK (~7MB) often exceeds 30s through trycloudflare; 206/200 on first byte is enough.
apk_code=$(curl -sk -o /dev/null -w "%{http_code}" --max-time 45 -r 0-0 "$apk_url" 2>/dev/null) || apk_code=000
apk_code=$(echo "$apk_code" | tr -cd '0-9')
case "$apk_code" in
  200|206) printf 'OK (%s)\n' "$apk_code" ;;
  *)
    printf 'FAIL (HTTP %s)\n' "$apk_code"
    echo "  - Large APK via tunnel may need: make repair-apk-urls after make dev-https"
    echo "  - Local check: curl -r 0-0 -o /dev/null -w '%{http_code}' http://127.0.0.1:8080/files/<apk>"
    exit 1
    ;;
esac

echo ""
echo "Enrollment HTTPS checks passed. Safe to scan QR on device."
