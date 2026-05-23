# API Contract: Hints (`/rest/private/hints`)

**Base path**: `/rest/private/hints`  
**Auth**: Session cookie and/or `Authorization: Bearer <jwt>`  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?: string, "data"?: T }`

**Java reference**: `com.hmdm.rest.resource.HintResource`

---

### GET `/history`

Returns hint keys already shown to the current user.

**Response data**: `string[]` (e.g. `["hint.step.1","hint.step.2"]`)

**Errors**: `403` if unauthenticated; `ERROR` + `error.internal.server` on failure

---

### POST `/history`

Marks a hint as shown for the current user.

**Body** (legacy-compatible):

- JSON string: `"hint.step.1"`
- Or plain text: `hint.step.1`

**Success**: `status: OK`, `data` optional

**Errors**: empty key → validation error; `403` if unauthenticated

---

### POST `/enable`

Clears all hint history for the current user (re-enable tutorials).

**Body**: none

**Success**: `status: OK`

---

### POST `/disable`

Clears history then inserts every `userHintTypes.hintKey` for the current user.

**Body**: none

**Success**: `status: OK`

**Response data** (optional): not required by React; Java returns OK only

---

## React consumers

| Client | Endpoints used |
|--------|----------------|
| `frontend/src/features/hints/hintsService.ts` | GET history, POST enable, POST disable |
| Legacy Angular `hint.service.js` | + POST history (mark shown) |

---

## Out of scope

- CRUD on `userHintTypes` catalog
- Per-customer or admin-managed hints for other users
- Hint content / localization (UI-only)
