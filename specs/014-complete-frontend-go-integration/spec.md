# Feature Specification: إكمال تكامل React ↔ Go (فجوات الفرونت والباكند)

**Feature Branch**: `014-complete-frontend-go-integration`

**Created**: 2026-05-21

**Status**: Draft

**Input**: إكمال النقوصات والفجوات بين واجهة React والباكند Go الجديد، استناداً إلى [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md).

**Related**:

- تحليل التكامل: [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md)
- فجوات REST: [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md)
- Schema (منجز): [`specs/013-complete-database-gaps/spec.md`](../013-complete-database-gaps/spec.md)
- REST المتبقي: [`specs/012-finish-java-go-backend/spec.md`](../012-finish-java-go-backend/spec.md)

**Baseline (موجود):**

- ~96% من مسارات REST التي يستدعيها React مسجّلة في Go؛ المصادقة، أجهزة (فلاتر + `info`)، تكوينات CRUD، تطبيقات، إعدادات أساسية، ملخص لوحة التحكم بعد 013.
- Schema حتى migration `000017` (جداول `devicestatuses`, `userrolesettings`, `configurationapplicationparameters`, …).

**هدف المشروع:** أن تُكتمل تجربة المستخدم في React عند الاتصال بـ `serverBackendGo` دون Java — بحفظ/قراءة الحقول الناقصة، وإغلاق فجوات API/سلوك المذكورة في تحليل التكامل، مع تحديث وثائق التكامل والـ parity عند كل إنجاز.

---

## Gap Matrix (من FRONTEND-GO-BACKEND-INTEGRATION)

| الأولوية | الفجوة | جانب | التأثير |
|----------|--------|------|---------|
| **P1** | حقول إعدادات tenant (`000015`) غير معروضة/محفوظة في React | Frontend + Backend | صفحة Settings ناقصة؛ أسماء أعمدة مخصصة للأجهزة |
| **P1** | محرر التكوين: سياسات MDM + `skipVersionCheck` + `remove`/`longTap` round-trip | Frontend + Backend | حفظ تكوين كامل يفشل صامتاً أو يفقد خيارات |
| **P1** | رفع الأيقونات عبر `icon-files` بدل fileId يدوي | Frontend (+ Backend موجود) | صفحة Icons غير قابلة للاستخدام العملي |
| **P2** | تحديث `devicestatuses` من مزامنة الوكيل | Backend | فلاتر حالة التثبيت تعكس بيانات قديمة |
| **P2** | `PUT /rest/public/stats` (جدول `usagestats` جاهز) | Backend | أدوات/ترخيص خارج UI |
| **P2** | أعمدة إضافية في قائمة الأجهزة (`model`, `batteryLevel`, …) | Backend (اختياري Frontend) | أعمدة الجدول فارغة رغم وجود البيانات في `info` |
| **P2** | `devicesEnrolledMonthly` في الملخص | Backend | رسم بياني شهري ناقص |
| **P3** | `POST /private/update` في صفحة Updates | Frontend | لا يمكن تطبيق التحديث من UI |
| **P3** | `POST /private/hints/history` (mark shown) | Frontend | تلميحات تتكرر |
| **⊘** | WebSocket، صفحات plugins الفرعية، `videos` REST | — | خارج النطاق v1 (موثّق) |

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - إعدادات المستأجر الكاملة (Priority: P1)

كمسؤول MDM، أريد تعديل كل إعدادات المستأجر (مجموعة الأجهزة الافتراضية، تنسيق الهاتف، أسماء الحقول المخصصة، قوالب العرض) من صفحة Settings في React، وأن تُحفظ وتُعاد عند إعادة فتح الصفحة كما في Java.

**Why this priority**: Schema جاهز من 013 لكن الواجهة والـ API لا يعرضان الحقول — يكسر تخصيص قائمة الأجهزة وإعدادات التسجيل.

**Independent Test**: تغيير `phoneNumberFormat` و`customPropertyName1` ثم حفظ → إعادة تحميل Settings → القيم نفسها؛ تبويب أعمدة الأجهزة يعرض التسميات المخصصة.

**Acceptance Scenarios**:

1. **Given** مستخدم بصلاحية settings، **When** يفتح Settings ويعدّل حقول tenant الجديدة ويحفظ، **Then** القيم تظهر بعد إعادة التحميل.
2. **Given** أسماء custom1–3 مُعرَّفة، **When** يفتح تبويب أعمدة الأجهزة حسب الدور، **Then** تظهر التسميات الصحيحة في واجهة الاختيار.
3. **Given** قيم افتراضية في قاعدة جديدة، **When** يفتح Settings لأول مرة، **Then** لا تظهر أخطاء وحقول لها قيم منطقية.

---

### User Story 2 - محرر التكوين MDM كامل (Priority: P1)

كمسؤول، أريد فتح تكوين موجود، تعديل سياسات MDM (WiFi، GPS، وضع Kiosk، …) وتطبيقات التكوين (بما فيها تخطي فحص الإصدار وإزالة/ضغطة طويلة)، ثم الحفظ وإعادة الفتح دون فقدان الخيارات.

**Why this priority**: UAT التكامل يظهر `[ ] Configurations: حفظ سياسات MDM كاملة` — أعلى خطر على نشر الأجهزة.

**Independent Test**: تعديل 3 سياسات في تبويبات التكوين + تفعيل skipVersionCheck لتطبيق واحد → حفظ → GET بالمعرّف → نفس القيم.

**Acceptance Scenarios**:

1. **Given** تكوين بسياسات MDM محفوظة مسبقاً، **When** أفتح المحرر، **Then** أرى نفس القيم في النماذج.
2. **Given** أغيّر `kioskMode` و`wifi` وأحفظ، **When** أعيد فتح التكوين، **Then** التغييرات محفوظة.
3. **Given** تطبيق مرتبط بالتكوين، **When** أفعّل skipVersionCheck وأحفظ، **Then** عند إعادة الفتح يبقى الخيار مفعّلاً.
4. **Given** تطبيق بخيارات remove/longTap، **When** أحفظ التكوين، **Then** القيم تُعاد في قائمة تطبيقات التكوين.

---

### User Story 3 - رفع الأيقونات من الواجهة (Priority: P1)

كمسؤول، أريد رفع ملف صورة أيقونة من صفحة Icons وربطه بسجل أيقونة دون إدخال معرّف ملف يدوياً.

**Why this priority**: مسار Go `POST /private/icon-files` موجود؛ React فقط لا يستخدمه — يمنع سير عمل عملي.

**Independent Test**: رفع PNG → إنشاء/تحديث أيقونة → تظهر في قائمة البحث مع معاينة URL صالحة.

**Acceptance Scenarios**:

1. **Given** ملف صورة صالح، **When** أرفعه من صفحة Icons، **Then** يُنشأ سجل أيقونة مرتبط بالملف.
2. **Given** رفع فاشل (نوع غير مدعوم)، **When** أحاول الرفع، **Then** رسالة خطأ واضحة دون كسر الصفحة.

---

### User Story 4 - حالة التثبيت تعكس الوكيل (Priority: P2)

كنظام، عندما يرسل الوكيل معلومات التثبيت، يجب تحديث حالة الجهاز في قاعدة البيانات بحيث تعكس فلاتر «حالة التثبيت» ولوحة الملخص الواقع.

**Why this priority**: بعد 013 الجدول موجود لكن القيم الافتراضية ثابتة دون ربط sync.

**Independent Test**: محاكاة `sync/info` لجهاز → `devicestatuses.applicationsstatus` يتغير → فلتر البحث يعكس التغيير.

**Acceptance Scenarios**:

1. **Given** جهاز بنجاح تثبيت تطبيقات في payload الوكيل، **When** تتم المزامنة، **Then** `applicationsstatus` يصبح SUCCESS (أو القيمة المكافئة المتفق عليها).
2. **Given** فشل تثبيت، **When** تتم المزامنة، **Then** الحالة FAILURE أو VERSION_MISMATCH حسب البيانات الواردة.

---

### User Story 5 - إحصائيات الخادم (Priority: P2)

كمنصة، أريد استقبال تقارير استخدام الخادم (heartbeat) على المسار العام كما في Java، لتخزينها في `usagestats` دون كسر الوكلاء أو أدوات المراقبة.

**Why this priority**: جدول `usagestats` من 013؛ لا REST بعد — 404 للمستهلكين الخارجيين.

**Independent Test**: `PUT /rest/public/stats` بpayload نموذجي → صف في `usagestats` أو upsert يومي.

**Acceptance Scenarios**:

1. **Given** payload إحصائيات صالح، **When** يُرسل PUT عام، **Then** يُخزَّن السجل مع تاريخ اليوم ومعرّف المثيل.
2. **Given** إرسال مكرر لنفس اليوم والمثيل، **When** يُعاد الإرسال، **Then** يُحدَّث السجل دون تكرار مخالف للقيود.

---

### User Story 6 - تحسينات ثانوية للوحة والأجهزة (Priority: P3)

كمستخدم، أريد (اختياري v1) تطبيق تحديثات من صفحة Updates، وتمييز التلميحات كمُعرَضة، ورؤية أعمدة إضافية في جدول الأجهزة عند الحاجة.

**Why this priority**: تحسين تجربة دون حجب الإطلاق.

**Independent Test**: كل عنصر يُختبر منفصلاً حسب القبول أدناه.

**Acceptance Scenarios**:

1. **Given** تحديث متاح، **When** أؤكد التطبيق من Updates، **Then** يُستدعى مسار التطبيق ويُعرض نتيجة للمستخدم.
2. **Given** تلميحاً عُرض، **When** أغلق الحوار، **Then** لا يُعرض مرة أخرى في الجلسة التالية (mark shown).

---

### Edge Cases

- حفظ تكوين بحقول MDM فارغة أو غير صالحة → رسائل خطأ دون تلف `settingsjson` الحالي.
- تكوين قديم من dump Java بمفاتيح JSON مختلفة → القراءة تعرض قيم معقولة أو افتراضيات (بعد 000017).
- رفع أيقونة بحجم كبير أو مسار ملف غير موجود → رفض آمن.
- مستخدم بدون صلاحية `settings` أو `configurations` → رفض صلاحيات متسق مع باقي الوحدات.
- `MODULE_SIGNUP_ENABLED=false` → صفحة التسجيل تبقى مع رسالة واضح (لا ضمن هذا المشروع إلا توثيق).

---

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM)*

**Backend (`serverBackendGo/`):**

- Modules: `settings`, `configurations`, `applications`, `icons`, `sync`, `stats` (جديد), `summary`, `devices` (اختياري).
- Java references: `SettingsResource`, `ConfigurationResource`, `IconFileResource`, `SyncResource`, `StatsResource`, `SummaryResource`, `DeviceResource`.
- REST paths unchanged: `/rest/private/*`, `/rest/public/stats`.
- Parity docs: `serverBackendGo/docs/parity/settings.md`, `configurations.md`, `icons.md`, `sync.md`, `summary.md`, `devices.md`, + `stats.md` (new).
- Layered architecture per module (domain / application / port / adapter).

**Frontend (`frontend/`):**

- Services: `settingsService.ts`, `configurationService.ts`, `iconsService.ts`, `updatesService.ts`, `hintsService.ts` — تحديث العقود لتطابق Go.
- Envelope parsing via `hmdmEnvelope.ts` دون تغيير شكل الاستجابة.
- لا تغيير `baseURL` أو مسارات REST عن Java إلا موثّق.

**Documentation:**

- تحديث [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md) عند إغلاق كل فجوة (§4, §6, §11).

### Functional Requirements

**Settings (P1)**

- **FR-001**: النظام MUST إرجاع حقول tenant الموسّعة (مجموعة جهاز جديد، تنسيق هاتف، أسماء/أعلام الحقول المخصصة، قالب رأس سطح المكتب، إرسال الوصف) في GET إعدادات المستأجر.
- **FR-002**: النظام MUST قبول حفظ نفس الحقول عبر مسارات الإعدادات الحالية (misc/lang/design أو مسار موحّد متفق عليه).
- **FR-003**: واجهة React MUST عرض وتحرير هذه الحقول في صفحة Settings وربطها بالحفظ.

**Configurations (P1)**

- **FR-004**: النظام MUST دمج `settingsjson` مع حقول التكوين الصريحة عند GET/PUT تكوين بحيث تبقى سياسات MDM قابلة للقراءة والكتابة.
- **FR-005**: النظام MUST إرجاع `skipVersionCheck` لكل تطبيق في تكوين عند القراءة إن وُجد في `configurationapplicationparameters`.
- **FR-006**: النظام MUST حفظ `skipVersionCheck`, `remove`, `longTap` عند حفظ التكوين.
- **FR-007**: محرر التكوين في React MUST إرسال مفاتيح سياسة MDM بصيغة camelCase المتوافقة مع Java عند الحفظ.

**Icons (P1)**

- **FR-008**: صفحة Icons في React MUST دعم رفع ملف عبر مسار رفع الأيقونات وربط النتيجة بسجل الأيقونة.

**Sync / Stats (P2)**

- **FR-009**: عند استلام معلومات جهاز من مسار المزامنة العام، النظام MUST تحديث أو إنشاء صف `devicestatuses` متسقاً مع حالة التثبيت المبلّغ عنها.
- **FR-010**: النظام MUST توفير مسار عام لقبول إحصائيات الاستخدام وتخزينها في `usagestats` مع منع التكرار لنفس اليوم والمثيل.

**Enhancements (P2–P3)**

- **FR-011**: (اختياري P2) بحث الأجهزة MAY إرجاع حقول شائعة (`model`, `batteryLevel`, `androidVersion`) في صفوف القائمة دون فتح التفاصيل.
- **FR-012**: (اختياري P2) ملخص الأجهزة MAY يوفّر سلسلة تسجيل شهرية أدق من الأصفار الثابتة.
- **FR-013**: (P3) صفحة Updates MUST استدعاء مسار تطبيق التحديث عند تأكيد المستخدم.
- **FR-014**: (P3) واجهة التلميحات MUST تسجيل التلميح كمُعرَض عبر API عند الإغلاق.

### Key Entities

- **Tenant Settings**: إعدادات المستأجر بما فيها تخصيص الحقول المخصصة وتنسيق الهاتف.
- **Configuration Policy**: سياسات MDM مجمّعة في JSON مع بيانات وصفية SQL.
- **Configuration Application Link**: ربط تطبيق بتكوين مع إصدار وخيارات remove/longTap/skipVersionCheck.
- **Device Install Status**: حالة تثبيت التطبيقات/الملفات لكل جهاز (جدول موجود).
- **Usage Stats Snapshot**: لقطة استخدام خادم يومية لكل مثيل.
- **Icon Asset**: أيقونة مرتبطة بملف مرفوع.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% من بنود UAT P1 في [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md) §10 (Settings tenant، Configuration MDM، Icons upload) تُعلَم منجزة بعد الاختبار اليدوي.
- **SC-002**: مسؤول MDM يكمل سيناريو «تعديل تكوين + حفظ + إعادة فتح» لـ 5 سياسات MDM مختلفة دون فقدان قيم في أول محاولة.
- **SC-003**: مسؤول MDM يرفع أيقونة ويراها في القائمة خلال دقيقة واحدة دون إدخال معرّف ملف يدوي.
- **SC-004**: بعد مزامنة وكيل اختبارية، فلتر «حالة التثبيت» يعكس الحالة المحدّثة لـ ≥90% من الأجهزة المختبرة (≥10 أجهزة في بيئة QA).
- **SC-005**: `FRONTEND-GO-BACKEND-INTEGRATION.md` يُحدَّث بحيث تصبح نسب «توافق الحقول» ≥95% للمسارات المستخدمة من React (مذكور في §1).
- **SC-006**: لا تظهر انحدارات على مسارات UAT المعلّمة منجزة سابقاً (login، بحث أجهزة، userRole columns، summary charts).

---

## Assumptions

- Schema 013 (`000011`–`000017`) مُطبَّق في بيئة التطوير والاختبار.
- React يبقى العميل الوحيد للوحة الإدارة في هذا المشروع؛ Angular/plugins الفرعية خارج النطاق v1.
- WebSocket، `videos` REST، وصفحات React لـ plugin push/messaging/audit **خارج النطاق** — تُوثَّق كـ ⊘ في تحليل التكامل.
- `MODULE_STATS_ENABLED` يُفعَّل في dev عند اختبار FR-010.
- bootstrap عميل جديد (نسخ قوالب عند إنشاء customer) وMailchimp **خارج النطاق** — تبقى في JAVA-GO-BACKEND-GAPS.
- مفاتيح JSON للتكوين تتبع اصطلاح Java/React الحالي (camelCase) كما في [`legacy-config-import.md`](../013-complete-database-gaps/contracts/legacy-config-import.md).

---

## Out of Scope (v1)

- WebSocket `/rest/ws/connect`
- REST `videos` وواجهة فيديوهات التدريب
- صفحات React كاملة لـ plugins (push, messaging, audit, deviceinfo, devicelog)
- CRUD عملاء من React (لوحة التحكم تستخدم search + impersonate فقط)
- نسخ bootstrap كامل عند إنشاء عميل (Java customers)
- دعم MQTT / Long polling للواجهة الإدارية (وكلاء فقط عبر notifications الحالي)

---

## Dependencies

- **013-complete-database-gaps**: migrations وAPI أساسية لـ `userrolesettings`, `devicestatuses`, CAP, `usagestats`.
- **012-finish-java-go-backend**: يمكن التنسيق على أجزاء متداخلة (devices enrichment) دون تعارض.
- [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) للعناصر التي تبقى backend-only بعد 014.
