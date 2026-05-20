#!/usr/bin/env sh
set -e
cd "$(dirname "$0")/.."
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi
export SERVER_PORT="${SERVER_PORT:-8080}"
go run ./cmd/server
