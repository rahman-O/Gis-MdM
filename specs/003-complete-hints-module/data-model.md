# Data Model: Phase 3 — Hints

**Date**: 2026-05-20  
**Storage**: PostgreSQL (migration `000004_hints_tables.up.sql`)

## Entity: UserHint

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| userId | int | FK → `users(id)` ON DELETE CASCADE |
| hintKey | varchar(100) | NOT NULL |
| created | timestamp | DEFAULT NOW() (optional in API) |

**Constraints**: UNIQUE (`userId`, `hintKey`)

**API exposure**: Only `hintKey` strings in `GET /history` response array.

## Entity: UserHintType (catalog)

| Field | Type | Rules |
|-------|------|-------|
| hintKey | varchar(100) | PK / UNIQUE |

**Purpose**: Master list of all tutorial keys; `disable` copies all rows into `userHints`
for the current user.

**Seed data** (from Java Liquibase):

- `hint.step.1`
- `hint.step.2`
- `hint.step.3`
- `hint.step.4`

## Relationships

```text
users (1) ──< userHints (many)
userHintTypes (catalog) ── copied into userHints on disable
```

## Application operations

| Use case | DB effect |
|----------|-----------|
| GetHistory | SELECT hintkeys for principal user id |
| MarkShown | INSERT one row (ignore duplicate) |
| Enable | DELETE all rows for user |
| Disable | DELETE all, then INSERT SELECT from catalog |

## Validation

- `hintKey` non-empty, max 100 chars
- Operations require authenticated principal with valid `userId`
- No cross-user access

## State (per user)

```text
[]  -- no rows → tours may show
[key1, key2, ...]  -- shown hints → tours skip those keys
[all catalog keys]  -- after disable → no tours
```
