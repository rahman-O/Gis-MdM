# Research: Phase 3 — Hints Module

**Date**: 2026-05-20

## R1 — SQL parity (source of truth)

**Decision**: Mirror MyBatis annotations in `UserMapper.java` exactly.

| Operation | SQL |
|-----------|-----|
| List history | `SELECT hintkey FROM userhints WHERE userid = $1` |
| Mark shown | `INSERT INTO userhints (userid, hintkey) VALUES ($1, $2)` |
| Enable | `DELETE FROM userhints WHERE userid = $1` |
| Disable | `DELETE ...` then `INSERT INTO userhints (userid, hintkey) SELECT $1, hintkey FROM userhinttypes` |

**Rationale**: Liquibase schema uses `userHints` / `userHintTypes`; PostgreSQL lowercases
unquoted identifiers to `userhints`, `userhinttypes`, `userid`, `hintkey`.

**Alternatives considered**: Storing history in JSON on `users` row — rejected (not Java parity).

## R2 — POST `/history` request body

**Decision**: Accept (1) raw JSON string body `"hint.step.1"`, (2) plain text body, and
(3) JSON object `{"hintKey":"..."}` if Angular sends wrapped form — primary path matches
Jersey `markHintAsShown(String hintKey)` and `$resource` POST parameter.

**Rationale**: Legacy `httpHintService.add(hintKey, ...)` passes the key as POST body;
React Hints page does not call mark-shown today but API must exist for tours.

**Alternatives considered**: Query param `?key=` — not used in Java resource.

## R3 — Duplicate hint key insert

**Decision**: Use `ON CONFLICT (userid, hintkey) DO NOTHING` and return `OK` anyway.

**Rationale**: Unique constraint `userHints_userId_hintKey_unique` in Liquibase; second
insert in Java may error — Go should be idempotent for better UX (spec acceptance).

**Alternatives considered**: Return `ERROR` on duplicate — worse for retries.

## R4 — Authentication

**Decision**: Rely on existing `JWTAuth` + `RequireAuth` + `EnrichPrincipal`; no extra
permission check (Java uses `SecurityContext.getCurrentUser()` only, not `settings`).

**Rationale**: Any logged-in user manages their own hints; React Hints nav uses
`permission: 'settings'` on UI only.

## R5 — Migration placement

**Decision**: New file `000004_hints_tables.up.sql` (not edit `000001` on deployed DBs).

**Rationale**: `000001_init` has no hint tables; additive migration is safest.

**Seed**: Insert four `userHintTypes` rows from Liquibase changeSet `19.09.19-18:48`.

## R6 — Module boundaries

**Decision**: Hints module owns `hint_repo.go`; does not import `users` or `auth` adapters.

**Rationale**: Constitution II — smallest module; only needs `principal.ID`.
