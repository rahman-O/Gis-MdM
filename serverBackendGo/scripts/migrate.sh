#!/usr/bin/env sh
set -e
cd "$(dirname "$0")/.."

apply_file() {
  f="$1"
  echo "  -> $f"
  docker compose exec -T db psql -U hmdm -d hmdm -v ON_ERROR_STOP=1 -f - < "$f"
}

if command -v docker >/dev/null 2>&1 && docker compose ps -q db 2>/dev/null | grep -q .; then
  echo "Applying migrations via Docker..."
  for f in db/migrations/*.up.sql; do
    [ -f "$f" ] || continue
    apply_file "$f"
  done
  docker compose exec -T db psql -U hmdm -d hmdm -c \
    "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
     INSERT INTO schema_migrations (version) VALUES ('000001_init') ON CONFLICT DO NOTHING;" 2>/dev/null || true
else
  echo "Start Postgres first: ./scripts/db-up.sh"
  exit 1
fi

echo "Done. Login: admin / admin (MD5 uppercase hex in JSON body for API)."
