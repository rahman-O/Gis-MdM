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
export FILES_DIRECTORY="${FILES_DIRECTORY:-$(pwd)/data/files}"
mkdir -p "$FILES_DIRECTORY"
go run ./cmd/server
