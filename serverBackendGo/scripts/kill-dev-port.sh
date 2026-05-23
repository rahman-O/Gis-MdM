#!/usr/bin/env sh
# Stops whatever is listening on SERVER_PORT (default 8080) — usually a leftover `make dev`.
set -eu
cd "$(dirname "$0")/.."
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi
PORT="${SERVER_PORT:-8080}"

if ! command -v lsof >/dev/null 2>&1; then
  echo "lsof not found; cannot detect process on :${PORT}"
  exit 1
fi

pids=$(lsof -tiTCP:"${PORT}" -sTCP:LISTEN 2>/dev/null || true)
if [ -z "$pids" ]; then
  echo "Nothing listening on :${PORT}"
  exit 0
fi

echo "Stopping process(es) on :${PORT}: $pids"
for pid in $pids; do
  kill "$pid" 2>/dev/null || true
done
sleep 1
still=$(lsof -tiTCP:"${PORT}" -sTCP:LISTEN 2>/dev/null || true)
if [ -n "$still" ]; then
  echo "Force kill: $still"
  for pid in $still; do
    kill -9 "$pid" 2>/dev/null || true
  done
fi
echo "Port :${PORT} is free."
