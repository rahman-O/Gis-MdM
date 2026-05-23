# API Contract: Icons upload flow (014 — frontend wiring)

**Private base**: `/rest/private`

**Java**: `IconFileResource`, `IconResource`  
**React**: `frontend/src/features/icons/iconsService.ts`, `IconsPage.tsx`

Backend routes **already exist** (Phase 9); 014 completes React usage.

---

## Step 1 — Upload file

### POST `/icon-files`

**Content-Type**: `multipart/form-data`

| Part | Name | Required |
|------|------|----------|
| file | `file` | yes — image (PNG typical) |

**Auth**: permission `icons` (or parity equivalent)

**Response data** (example):

```json
{
  "fileId": 42,
  "url": "/files/..."
}
```

**Errors**: invalid type/size → `status: ERROR`, message for UI toast.

---

## Step 2 — Create or update icon record

### PUT `/icons`

**Body**:

```json
{
  "id": null,
  "name": "Home",
  "fileId": 42
}
```

**Response**: `{ "status": "OK", "data": { "id": 7, "name": "Home", "fileId": 42, ... } }`

---

## Step 3 — List / search (unchanged)

- `GET /icons/search`
- `GET /icons/search/{value}`

Preview in UI uses `url` from file record or constructed path from `fileId`.

---

## Frontend contract

`iconsService.uploadIconFile(file: File): Promise<{ fileId: number }>`  
then `saveIcon({ name, fileId })`.

No manual `fileId` input in UI after 014.
