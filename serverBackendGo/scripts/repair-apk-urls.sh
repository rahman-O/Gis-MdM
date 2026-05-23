#!/usr/bin/env sh
# Rebuilds empty applicationversions.url from filepath + BASE_URL (after HTTPS tunnel / BASE_URL change).
# Requires: Postgres (docker compose db) and committed APK files under FILES_DIRECTORY.
set -eu
cd "$(dirname "$0")/.."

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi

BASE=$(echo "${BASE_URL:-}" | sed 's|/$||')
if [ -z "$BASE" ]; then
  echo "ERROR: set BASE_URL in .env first (e.g. make dev-https)"
  exit 1
fi

case "$BASE" in
  https://*) ;;
  *)
    echo "ERROR: BASE_URL must be HTTPS: $BASE"
    exit 1
    ;;
esac

sql_base=$(echo "$BASE" | sed "s/'/''/g")

run_sql() {
  if command -v psql >/dev/null 2>&1 && [ -n "${DATABASE_URL:-}" ]; then
    psql "$DATABASE_URL" -c "$1"
    return
  fi
  if docker compose ps -q db 2>/dev/null | grep -q .; then
    docker compose exec -T db psql -U hmdm -d hmdm -c "$1"
    return
  fi
  echo "ERROR: need psql or running 'docker compose' db service"
  exit 1
}

echo "Repairing APK URLs and configuration baseUrl with BASE_URL=$BASE"

run_sql "
UPDATE configurations
SET baseurl = '${sql_base}'
WHERE baseurl IS NULL OR trim(baseurl) = ''
   OR baseurl ~ '^https?://[^/]*\\.trycloudflare\\.com'
   OR baseurl ~ '^https?://(localhost|127\\.0\\.0\\.1)';
"

run_sql "
UPDATE applicationversions av
SET url = '${sql_base}/files/' || cu.filesdir || '/' || replace(av.filepath, E'\\\\', '/')
FROM applications a
JOIN customers cu ON cu.id = a.customerid
WHERE av.applicationid = a.id
  AND (av.url IS NULL OR trim(av.url) = '')
  AND av.filepath IS NOT NULL AND trim(av.filepath) <> ''
  AND av.filepath NOT LIKE '%1111111%';
"

run_sql "
UPDATE applicationversions
SET url = regexp_replace(
  regexp_replace(url, '^http://localhost:8080', '${sql_base}'),
  '^http://127\\.0\\.0\\.1:8080', '${sql_base}')
WHERE url ~ '^http://(localhost|127\\.0\\.0\\.1)(:8080)?/';
"

# Previous dev-https tunnels (trycloudflare host changes every session)
run_sql "
UPDATE applicationversions
SET url = '${sql_base}' || regexp_replace(url, '^https://[^/]+', '')
WHERE url ~ '^https://';
"

run_sql "
UPDATE applications
SET url = '${sql_base}' || regexp_replace(url, '^https://[^/]+', '')
WHERE url ~ '^https://';
"

run_sql "
UPDATE applications a
SET url = sub.url
FROM (
  SELECT av.applicationid, av.url
  FROM applicationversions av
  WHERE av.url IS NOT NULL AND trim(av.url) <> ''
  ORDER BY av.versioncode DESC, av.id DESC
) sub
WHERE a.id = sub.applicationid
  AND (a.url IS NULL OR trim(a.url) = '' OR a.url ~ '^http://(localhost|127\\.0\\.0\\.1)');
"

echo "Done. Restart make dev if running, then: make verify-https"
