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
# .env may copy production paths; remap /var/lib/hmdm for local dev (no sudo on macOS).
case "${FILES_DIRECTORY:-}" in
  ""|/var/lib/hmdm|/var/lib/hmdm/*)
    export FILES_DIRECTORY="$(pwd)/data/files"
    ;;
esac
export FILES_DIRECTORY="${FILES_DIRECTORY:-$(pwd)/data/files}"
mkdir -p "$FILES_DIRECTORY"
go run ./cmd/server
