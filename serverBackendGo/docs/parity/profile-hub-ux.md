# Parity: Profile Hub UX (019)

**Feature**: `019-profile-hub-ux`  
**Contracts**: `specs/019-profile-hub-ux/contracts/`

## Hub API

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /private/profiles` | Extended | Adds `health`, `healthReasons`, `badges`, `assignmentCount`, `rolloutFailureCount` via `HubService.List` |
| `GET /private/profiles/:id/summary` | New | Cockpit + overview cards |
| `GET /private/profiles/:id/activity` | New | Timeline from `domain_events` (`limit` query, default 50) |

### Health rules

- `draft_only`: no published version
- `error`: rollout failures while enabled
- `warning`: no assignment, stale publish (`PROFILE_STALE_PUBLISH_DAYS`, default 30), or disabled with issues
- `healthy`: otherwise when published

### Activity events

Emitted on enable/disable (`ProfileEnabled`, `ProfileDisabled`). Additional events from publish/assignment depend on 017/018 emitters.

## Enrollment routes (policy decoupling)

| Change | Status |
|--------|--------|
| `profileVersionId` optional on create/update | Done |
| `validateBinding` skips profile when `profileVersionId <= 0` | Done |
| Admin UI profile version picker | Removed |
| `GET .../options/published-profile-versions` | Deprecated for new UI; legacy clients may still call |

## Frontend workspace

- List opens `ProfileWorkspace` via `?open={id}&section=`
- Legacy `/profiles/:id/edit` redirects to `?open={id}&section=editor`
- Desktop: Dialog ~96vw×94vh; mobile: full Sheet

## Manual verification

```bash
# List with health
curl -s -H "Cookie: ..." "$BASE_URL/rest/private/profiles" | jq '.data[0] | {id, health, badges}'

# Summary
curl -s -H "Cookie: ..." "$BASE_URL/rest/private/profiles/1/summary" | jq '.data.health'

# Activity
curl -s -H "Cookie: ..." "$BASE_URL/rest/private/profiles/1/activity?limit=10" | jq '.data.items | length'

# Route without profile
curl -s -X POST -H "Content-Type: application/json" -d '{"name":"Test","defaultTreeNodeId":1,"mainAppId":1}' \
  "$BASE_URL/rest/private/enrollment-routes"
```
