# Research: إكمال تكامل React ↔ Go (014)

**Branch**: `014-complete-frontend-go-integration` | **Date**: 2026-05-21

## R1 — نطاق MVP vs موجات لاحقة

**Decision**: MVP = **US1–US3** (Settings tenant fields, Configuration MDM round-trip, Icons upload). US4–US6 في موجة 2 (sync status, stats, P3 polish).

**Rationale**: يغلق 100% من بنود UAT P1 في `FRONTEND-GO-BACKEND-INTEGRATION.md` §10 دون انتظار module `stats` أو منطق sync معقد.

**Alternatives considered**: إنجاز كل FR في PR واحد — مرفوض (تشتيت review؛ frontend + backend + module جديد).

---

## R2 — حفظ حقول Settings (`000015`)

**Decision**: توسيع `domain.Settings` + `settings_repo` SELECT/UPDATE لكل أعمدة `000015`؛ إدراجها في `POST /misc` body merge (نفس نمط Java — حقول tenant ضمن misc/lang وليس endpoint جديد).

**Rationale**: React [`settingsService.ts`](../../frontend/src/services/settingsService.ts) يحفظ عبر `misc` ثم `lang`؛ توسيع `normalizeSettings` + `miscBody` أقل كسراً من endpoint جديد.

**Alternatives considered**: `POST /settings/tenant` منفصل — مرفوض (كسر parity مسارات).

---

## R3 — Configuration `settingsjson` merge

**Decision**:

- **GET/PUT**: استمرار `json.Unmarshal(settingsjson, cfg)` على GET؛ على PUT استخراج كل حقول MDM من `Configuration` struct (بما فيها الحقول الموجودة في [`types.ts`](../../frontend/src/features/configurations/types.ts)) إلى `map[string]any` ثم `json.Marshal` → `settingsjson` مع دمج الأعمدة SQL الصريحة.
- **CAP**: `LEFT JOIN configurationapplicationparameters` في `ListConfigurationApplications`؛ `skipVersionCheck` في domain + scan.
- **remove/longTap**: إضافة إلى SELECT في `ListConfigurationApplications` (موجودة في INSERT فقط حالياً).

**Rationale**: React يرسل كائناً مسطحاً بمفاتيح camelCase؛ Go repo يدمج JSON عند الحفظ منذ 000007 لكن قراءة CAP ناقصة.

**Alternatives considered**: تخزين كل سياسة كعمود SQL — مرفوض (013 + constitution: `settingsjson` canonical).

---

## R4 — Icons upload

**Decision**: استخدام المسار الموجود `POST /rest/private/icon-files` ([`icon_file_handler.go`](../../serverBackendGo/internal/modules/icons/adapter/http/icon_file_handler.go)) من [`IconsPage`](../../frontend/src/features/icons/IconsPage.tsx) عبر `FormData` + `iconsService.uploadIconFile()` ثم `PUT /icons` بـ `{ name, fileId }`.

**Rationale**: Backend Phase 9 منجز؛ الفجوة frontend-only لـ FR-008.

**Alternatives considered**: رفع عبر `web-ui-files` — مرفوض (مسار مختلف لملفات UI وليس أيقونات مربعة).

---

## R5 — Module `stats` (P2)

**Decision**: وحدة جديدة `internal/modules/stats` — `PUT /rest/public/stats` بدون auth (مثل Java)، upsert على `(ts, instanceid)` في `usagestats`، `MODULE_STATS_ENABLED` في `config` + `modules.go`.

**Rationale**: جدول `000014` جاهز؛ Java `StatsResource` بسيط (PUT body = UsageStats).

**Alternatives considered**: دمج في `publicapi` — مرفوض (module-first constitution).

---

## R6 — تحديث `devicestatuses` من sync (P2)

**Decision**: في `sync/application` بعد `UpdateInfo`، استدعاء `DeviceStatusUpserter` (port) يحسب `applicationsstatus` / `configfilesstatus` من `info.applications` / `info.files` (نفس قيم Java: SUCCESS, FAILURE, VERSION_MISMATCH, OTHER) و `INSERT ... ON CONFLICT DO UPDATE`.

**Rationale**: 013 أنشأ الجدول؛ 014 يربطه بسلوك الوكيل دون انتظار خدمة Java `DeviceStatusService` كاملة.

**Alternatives considered**: Cron إعادة حساب — مرفوض (تأخير؛ فلاتر UI خاطئة حتى Cron).

---

## R7 — Device list enrichment (P2 optional)

**Decision**: اختياري — إضافة حقول مسطحة في `DeviceView` من `infojson` في `Search` SELECT (COALESCE) دون JOIN إضافي.

**Rationale**: يحسّن أعمدة الجدول؛ ليس حاجز MVP.

**Alternatives considered**: لا شيء — مقبول لـ P2.

---

## R8 — توثيق التكامل

**Decision**: كل story تُغلق بتحديث [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md) §4/§6/§10 + parity doc المعني.

**Rationale**: FR في spec + SC-005.

**Alternatives considered**: ملف جديد فقط — مرفوض (مصدر واحد للفرونت).
