# API Contract: Configurations list (Phase 4 read-only)

**Path**: `GET /rest/private/configurations/list`  
**Module**: `internal/modules/configurations/`  
**Auth**: Bearer / session  
**Phase**: 4 (minimal); full CRUD in Phase 5

**Java reference**: `ConfigurationResource` list endpoint used by Angular/React devices UI

---

### GET `/list`

Returns configuration id/name pairs for the authenticated user's customer.

**Response data**:

```json
[
  { "id": 1, "name": "Default" }
]
```

**Errors**: permission denied if unauthenticated

---

## React consumer

`deviceService.getConfigurations()` → populates configuration filter on Devices page.

## Out of scope

- PUT/POST configuration editor
- Configuration files, applications linkage
