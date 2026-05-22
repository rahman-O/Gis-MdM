#!/usr/bin/env sh
# Starts a free HTTPS tunnel, writes BASE_URL into .env, and verifies QR/APK endpoints.
#
# Prereq: backend running on SERVER_PORT (default 8080) — `make dev` in another terminal.
#
# Usage:
#   ./scripts/dev-https-auto.sh
#   ./scripts/dev-https-auto.sh stop   # stop background tunnel

set -eu
cd "$(dirname "$0")/.."
ROOT="$(pwd)"
ENV_FILE="${ROOT}/.env"
PORT="${SERVER_PORT:-8080}"
TARGET="http://127.0.0.1:${PORT}"
PID_FILE="${ROOT}/.https-tunnel.pid"
LOG_FILE="${ROOT}/.https-tunnel.log"

stop_tunnel() {
  if [ -f "$PID_FILE" ]; then
    pid=$(cat "$PID_FILE" 2>/dev/null || true)
    if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
      echo "Stopped tunnel (pid $pid)."
    fi
    rm -f "$PID_FILE"
  fi
}

if [ "${1:-}" = "stop" ]; then
  stop_tunnel
  exit 0
fi

if ! command -v cloudflared >/dev/null 2>&1; then
  echo "Installing cloudflared via Homebrew..."
  if ! command -v brew >/dev/null 2>&1; then
    echo "Install Homebrew first: https://brew.sh"
    exit 1
  fi
  brew install cloudflared
fi

if ! curl -sf --max-time 2 "${TARGET}/rest/public/name" >/dev/null 2>&1; then
  echo "ERROR: Go backend not reachable at ${TARGET}"
  echo "Start it first: make dev"
  exit 1
fi

stop_tunnel
: >"$LOG_FILE"

echo "Starting Cloudflare quick tunnel → ${TARGET}"
cloudflared tunnel --url "${TARGET}" >>"$LOG_FILE" 2>&1 &
echo $! >"$PID_FILE"
cf_pid=$(cat "$PID_FILE")

public_url=""
i=0
while [ "$i" -lt 45 ]; do
  public_url=$(grep -oE 'https://[a-zA-Z0-9.-]+\.trycloudflare\.com' "$LOG_FILE" 2>/dev/null | head -1 || true)
  if [ -n "$public_url" ]; then
    break
  fi
  if ! kill -0 "$cf_pid" 2>/dev/null; then
    echo "cloudflared exited early. Log:"
    tail -20 "$LOG_FILE" 2>/dev/null || true
    exit 1
  fi
  sleep 1
  i=$((i + 1))
done

if [ -z "$public_url" ]; then
  echo "Could not read tunnel URL from log. See ${LOG_FILE}"
  exit 1
fi

public_url=$(echo "$public_url" | sed 's|/$||')
echo "Tunnel URL: ${public_url}"

# Update or append BASE_URL in .env
if [ ! -f "$ENV_FILE" ]; then
  cp .env.example "$ENV_FILE" 2>/dev/null || touch "$ENV_FILE"
fi
if grep -q '^BASE_URL=' "$ENV_FILE" 2>/dev/null; then
  # BSD sed
  sed -i '' "s|^BASE_URL=.*|BASE_URL=${public_url}|" "$ENV_FILE"
else
  printf '\nBASE_URL=%s\n' "$public_url" >>"$ENV_FILE"
fi
echo "Updated ${ENV_FILE} → BASE_URL=${public_url}"
echo ""
echo "IMPORTANT: Restart the Go server (Ctrl+C in make dev, then make dev again)."
echo ""

# Wait until cloudflared registers (can take 30–90s on slow DNS)
wait_i=0
while [ "$wait_i" -lt 90 ]; do
  if grep -q 'Registered tunnel connection' "$LOG_FILE" 2>/dev/null; then
    break
  fi
  sleep 2
  wait_i=$((wait_i + 2))
done
sleep 5

echo "==> Verifying public HTTPS routes (via tunnel)"
./scripts/verify-enrollment-https.sh "$public_url" || verify_rc=$?
verify_rc=${verify_rc:-0}

echo ""
echo "Tunnel PID: $cf_pid (log: ${LOG_FILE})"
echo "Stop tunnel: ./scripts/dev-https-auto.sh stop"
echo ""
echo "Next: restart make dev → Applications (re-save APK) → Configuration (Save) → scan QR"

exit "$verify_rc"
