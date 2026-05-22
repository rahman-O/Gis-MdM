# Implementation Plan: إكمال تكامل React ↔ Go

**Branch**: `014-complete-frontend-go-integration` | **Date**: 2026-05-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/014-complete-frontend-go-integration/spec.md`  
**Integration baseline**: [FRONTEND-GO-BACKEND-INTEGRATION.md](../../FRONTEND-GO-BACKEND-INTEGRATION.md)

## Summary

إغلاق فجوات **P1** بين React و `serverBackendGo` بعد اكتمال schema **013** (`000011`–`000017`): (1) حقول إعدادات tenant في API + Settings UI، (2) round-trip كامل لسياسات MDM وتطبيقات التكوين (`settingsjson`, CAP, `remove`/`longTap`)، (3) رفع الأيقونات عبر `icon-files`. موجة **P2**: ربط `devicestatuses` بـ `sync/info`، module `stats` لـ `PUT /public/stats`، تحسينات اختيارية للأجهزة/الملخص. موجة **P3**: Updates apply + hints history من الواجهة.

**Technical approach**: توسيع modules موجودة (`settings`, `configurations`, `icons`) + frontend services/pages؛ module جديد `stats`؛ port في `sync` لـ device status upsert — بدون migrations جديدة في 014.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 (`frontend`)

**Primary Dependencies**: Gin, `pgx`/`database/sql`, React, existing `apiClient` / feature services

**Storage**: PostgreSQL (schema من 013؛ لا DDL في 014)

**Testing**: `go test ./internal/modules/...`; manual UAT من [quickstart.md](./quickstart.md); تحديث parity docs

**Target Platform**: Linux/macOS dev (`make dev`); React dev proxy → Go `:8080`

**Project Type**: Full-stack integration (web admin + REST API)

**Performance Goals**: لا تغيير SLA؛ sync upsert خفيف (صف واحد لكل جهاز)

**Constraints**: Headwind JSON envelope؛ مسارات `/rest/private` و `/rest/public`؛ MD5/JWT auth كما هو

**Scale/Scope**: ~6 modules لمسات؛ 3–5 صفحات/خدمات React؛ 5 عقود API؛ خارج النطاق: WebSocket, plugins UI, `videos`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `settings`, `configurations`, `icons`, `sync`, new `stats`; phase في MIGRATION/NEXT_STEPS |
| **II. Layered Clean** | ✅ | تغييرات في domain/application/port/adapter فقط |
| **III. API Parity** | ✅ | Java Resources مرجع؛ contracts + `docs/parity/*.md` |
| **IV. Testable Delivery** | ✅ | service tests + quickstart smoke |
| **V. Simplicity** | ✅ | لا migrations؛ لا abstractions زائدة |
| **VI. Security** | ✅ | private routes + permissions؛ stats public كـ Java |
| **VII. Observability** | ✅ | `MODULE_STATS_ENABLED`; أخطاء httpx |

**Post-design**: جميع البوابات ✅ — لا استثناءات في Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/014-complete-frontend-go-integration/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/           # Phase 1
└── tasks.md             # /speckit-tasks (not created here)
```

### Source Code (repository root)

```text
serverBackendGo/
├── internal/modules/settings/          # US1: repo + handler misc/lang
├── internal/modules/configurations/      # US2: settingsjson merge, CAP, remove/longTap
├── internal/modules/icons/               # US3: (verify) icon-files
├── internal/modules/sync/                # US4: device status upsert after info
├── internal/modules/stats/               # US5: new module
├── internal/modules/devices/             # US4 optional: infojson columns
├── internal/modules/summary/             # US4 optional: monthly series
├── docs/parity/settings.md
├── docs/parity/configurations.md
├── docs/parity/icons.md
├── docs/parity/sync.md
└── docs/parity/stats.md                  # new

frontend/src/
├── features/settings/                    # US1: types, SettingsPage, settingsService
├── features/configurations/              # US2: types, editor tabs, configurationService
├── features/icons/                       # US3: IconsPage, iconsService
├── features/updates/                     # US6 P3
└── features/hints/                       # US6 P3

FRONTEND-GO-BACKEND-INTEGRATION.md        # تحديث عند كل story
```

**Structure Decision**: Monorepo موجود — تعديلات متوازية frontend + `serverBackendGo` دون هيكل جديد.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (none) | — | — |

---

## Phase 0: Research

**Output**: [research.md](./research.md) — قرارات R1–R8 (MVP US1–3، settings merge، settingsjson، icons، stats module، sync status، optional device columns، doc updates).

**NEEDS CLARIFICATION**: none (كلها محلولة).

---

## Phase 1: Design & Contracts

**Outputs**:

| Artifact | Path |
|----------|------|
| Data model | [data-model.md](./data-model.md) |
| Settings API delta | [contracts/settings-api.md](./contracts/settings-api.md) |
| Configurations round-trip | [contracts/configurations-api.md](./contracts/configurations-api.md) |
| Icons upload | [contracts/icons-api.md](./contracts/icons-api.md) |
| Stats public | [contracts/stats-api.md](./contracts/stats-api.md) |
| Sync device status | [contracts/sync-device-status.md](./contracts/sync-device-status.md) |
| UAT / smoke | [quickstart.md](./quickstart.md) |

---

## Phase 2: Implementation Waves (for `/speckit-tasks`)

### Wave A — MVP P1 (US1–US3)

| ID | Work | Backend | Frontend |
|----|------|---------|----------|
| A1 | Tenant settings fields | `settings` domain/repo/handler | `types.ts`, `SettingsPage`, `settingsService` |
| A2 | Configuration MDM + CAP + remove/longTap | `configurations` domain, `config_repo`, handler | `types.ts`, editor, `configurationService` |
| A3 | Icon upload UX | verify `icons` handler | `iconsService`, `IconsPage` |
| A4 | Docs | parity `settings`, `configurations`, `icons` | `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 |

**Exit**: UAT P1 checklist في quickstart — كل بنود `[x]`.

### Wave B — P2 (US4–US5 + optional)

| ID | Work |
|----|------|
| B1 | `DeviceStatusUpserter` في `sync` بعد `UpdateInfo` |
| B2 | Module `stats`: `PUT /rest/public/stats` + migration N/A (table exists) |
| B3 | Optional: device search columns from `infojson` |
| B4 | Optional: `devicesEnrolledMonthly` in summary |
| B5 | parity `sync.md`, `stats.md` |

### Wave C — P3 (US6)

| ID | Work |
|----|------|
| C1 | `updatesService` → `POST /private/update` on apply |
| C2 | `hintsService` → `POST /private/hints/history` mark shown |

### Dependencies

```text
013 migrations (000011–000017) → Wave A → Wave B → Wave C
```

**Prerequisite**: `make migrate` ≥ `000017` before any 014 smoke test.

---

## Phase 2 Planning Stop

`/speckit-plan` ends here. Next: `/speckit-tasks` → `tasks.md`, then `/speckit-implement`.

**Generated artifacts**:

- `specs/014-complete-frontend-go-integration/plan.md` (this file)
- `specs/014-complete-frontend-go-integration/research.md`
- `specs/014-complete-frontend-go-integration/data-model.md`
- `specs/014-complete-frontend-go-integration/quickstart.md`
- `specs/014-complete-frontend-go-integration/contracts/*.md`
