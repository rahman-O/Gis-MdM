# Research: إكمال فجوات قاعدة البيانات (013)

**Branch**: `013-complete-database-gaps` | **Date**: 2026-05-21

## R1 — ترقيم migrations بعد `000010`

**Decision**: استخدام `000011` … `000017` (ملف `.up.sql` / `.down.sql` لكل إصدار). عدم إعادة ترقيم `000008_devices_search_extras` (موجود على الفرع 012).

**Rationale**: golang-migrate يعتمد على الترتيب الرقمي؛ `000008_files_icons_core` و`000008_devices_search_extras` يتشاركان البادئة — يعملان لأن الأسماء مختلفة، لكن الإصدارات الجديدة تبدأ من `000011` كما في [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) §7.

**Alternatives considered**: دمج كل التغييرات في migration واحدة — مرفوض (صعوبة rollback، مراجعة PR).

---

## R2 — `devicestatuses` وقيم الحالة

**Decision**:

- جدول مطابق Java: PK = `deviceid` FK → `devices(id) ON DELETE CASCADE`.
- أعمدة: `configfilesstatus VARCHAR(100)`, `applicationsstatus VARCHAR(100)`.
- Backfill في نفس migration `000011`: `INSERT … SELECT deviceid, 'OTHER', 'FAILURE' FROM devices` حيث لا صف.
- فلتر `installationStatus` في `device_repo` يستخدم `LEFT JOIN devicestatuses ds ON ds.deviceid = d.id` بدلاً من `infojson` فقط.

**Rationale**: يطابق [`DeviceMapper.xml`](../../backend/common/src/main/java/com/hmdm/persistence/mapper/DeviceMapper.xml) و[`DeviceStatusService`](../../backend/common/src/main/java/com/hmdm/service/DeviceStatusService.java). قيم enum Java: `SUCCESS`, `FAILURE`, `VERSION_MISMATCH`, … — نخزّن كنص كما في Java.

**Alternatives considered**: حساب الحالة من `infojson` فقط — مرفوض لـ FR-011 وSummary.

---

## R3 — `userrolesettings` وأعمدة العرض

**Decision**:

- إنشاء جدول بكل أعمدة `columndisplayed*` من Liquibase (15+ عموداً إضافياً من changelogs لاحقة: battery, files, mdm, kiosk, android, enroll, serial, publicip, custom1–3, defaultlauncher).
- توسيع `settings/domain.UserRoleSettings` و`settings_repo` لـ `GetUserRoleSettings` / `Save` من الجدول.
- Seed: نسخ منطق Java `04.10.19-13:50` — كل الأعمدة `TRUE` للأدوار 1–3 لكل `customerid` موجود.

**Rationale**: [`settings/adapter/http/handler.go`](../../serverBackendGo/internal/modules/settings/adapter/http/handler.go) يُرجع اليوم `{roleId}` فقط؛ React [`devices.controller.js`](../../backend/server/src/main/webapp/app/components/main/controller/devices.controller.js) يستدعي `getUserRoleSettings`.

**Alternatives considered**: تخزين التفضيلات في `settings` — يخالف نموذج Java multi-role.

---

## R4 — `configurations.settingsjson` مقابل أعمدة Java

**Decision**:

- **مصدر جديد (Go-native)**: الكتابة تستمر عبر `config_repo` إلى `settingsjson` + الأعمدة الصريحة الحالية.
- **استيراد legacy**: migration `000017` **اختيارية** — `DO $$ … IF column exists` ينسخ أعمدة Java المعروفة إلى مفاتيح camelCase في JSON (مثلاً `kioskMode`, `gps`, `blockStatusBar`) باستخدام `jsonb_build_object` / `||`.
- توثيق قائمة المفاتيح في [`contracts/legacy-config-import.md`](./contracts/legacy-config-import.md).

**Rationale**: [`Configuration.Extra`](../../serverBackendGo/internal/modules/configurations/domain/configuration.go) و`settingsjson` موجودان؛ dump Java يحتوي أعمدة غير موجودة في schema Go.

**Alternatives considered**: إضافة 60+ عموداً لـ `configurations` — مرفوض (تعارض مع قرار 007).

---

## R5 — دوال SQL Java

**Decision**: لا stored procedures في PostgreSQL لهذه الميزة. `mdm_config_app_upgrade` → منطق في `applications/application` عند الحاجة (خارج 013 v1 إن لم يكن مستدعى). فلاتر launcher تبقى `infojson` + join `applicationversions` لاحقاً.

**Rationale**: spec Assumptions + constitution V (بساطة).

---

## R6 — `usagestats` ووحدة stats

**Decision**: migration `000014` فقط في 013؛ module `stats` REST يبقى في 012 (أو يُفعّل بعد 014). جدول مطابق Java مع UNIQUE `(ts, instanceid)`.

**Rationale**: فصل schema (013) عن API (012) يسمح بالتنفيذ المتوازي.

---

## R7 — جداول plugins اختيارية (P3)

**Decision**: ⊘ في v1 — لا migrations لـ `plugin_deviceinfo_deviceparams_wifi` وغيرها. توثيق في plan Out of Scope.

**Rationale**: spec FR-010؛ لا modules Go لـ devicelocations/photo.

---

## R8 — Idempotency و down migrations

**Decision**: كل `up` يستخدم `CREATE TABLE IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`; كل `down` يعكس بالترتيب (DROP TABLE / DROP COLUMN). لا `CASCADE` على بيانات production في down إلا حيث الجدول جديد.

**Rationale**: FR-006, SC-006؛ constitution V.
