#!/usr/bin/env sh
set -e
cd "$(dirname "$0")/.."
if docker compose version >/dev/null 2>&1; then
  docker compose up -d db
elif command -v docker-compose >/dev/null 2>&1; then
  docker-compose up -d db
else
  echo "Docker Compose not found. Install Docker Desktop or start PostgreSQL manually on port 5432."
  exit 1
fi
echo "Waiting for PostgreSQL..."
for i in 1 2 3 4 5 6 7 8 9 10; do
  if docker compose exec -T db pg_isready -U hmdm -d hmdm >/dev/null 2>&1; then
    echo "PostgreSQL is ready (postgres://hmdm:hmdm@localhost:5432/hmdm)"
    exit 0
  fi
  sleep 2
done
echo "Database container started but health check timed out. Check: docker compose logs db"
exit 1
