# API Contract: Product Updates (`/rest/private/update`)

**Base path**: `/rest/private/update`  
**Auth**: Bearer JWT  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.UpdateResource`  
**React reference**: `frontend/src/features/updates/updatesService.ts`

---

### GET `/check`

Fetch update manifest and compute outdated flags for web, launcher, and secondary APKs.

**Permissions**: Super-admin when multi-tenant (`!isSingleCustomer && !superAdmin` → denied)

**Response data**: `UpdateEntry[]` (`pkg`, `version`, `currentVersion`, `url`, `outdated`, `downloaded`, `updateDisabled`, …)

**Errors**: `ERROR` when manifest URL fetch fails

---

### POST `/`

Download and/or apply selected updates.

**Body**: `UpdateRequest`

```json
{
  "updates": [ { "pkg": "web", "version": "...", "outdated": true, ... } ],
  "update": true,
  "sendStats": false
}
```

**Behavior**:
- Downloads APK/web artifacts to `FILES_DIRECTORY` when `outdated && !updateDisabled`
- When `update=true`, upgrades configuration application versions for mobile packages
- `sendStats` → vendor stats URL (**partial** stub acceptable)

**Response data**: `UpdateEntry[]` with updated flags

**Permissions**: Same super-admin guard as check for multi-tenant
