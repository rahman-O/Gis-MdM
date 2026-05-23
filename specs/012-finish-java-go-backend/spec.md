# Feature Specification: إكمال نقل الباكند Java → Go

**Feature Branch**: `012-finish-java-go-backend`

**Created**: 2026-05-21

**Status**: Draft

**Input**: إكمال النقل من الباكند القديم (Java) إلى الجديد (Go)، استناداً إلى [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) و[`JAVA-GO-MIGRATION-STATUS.md`](../../JAVA-GO-MIGRATION-STATUS.md).

**Related**:

- تحليل الفجوات التفصيلي: [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md)
- حالة الهجرة: [`JAVA-GO-MIGRATION-STATUS.md`](../../JAVA-GO-MIGRATION-STATUS.md)
- تكامل الواجهة: [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md)
- مواصفة سابقة (Phase 9 جزئي): [`specs/011-complete-migration-gaps/spec.md`](../011-complete-migration-gaps/spec.md) — **يُكمّل** هذا المشروع ما تبقى بعد T046–T093

**Baseline (منجز مسبقاً — لا يُعاد تنفيذه):**

- Phases 1–8 في `serverBackendGo/docs/MIGRATION.md` (مسارات REST أساسية + React).
- Phase 9 جزئي: push notifier عند حفظ التكوين، cron جدولة plugin push، `POST /rest/private/icon-files`.

**هدف المشروع:** الوصول إلى **إمكانية إيقاف Java WAR** للعمليات اليومية (لوحة React + وكلاء MDM) مع **تطابق سلوكي ≥95%** مع Java للمسارات المستخدمة، وتوثيق صريح لما يبقى ⊘.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - إدارة أجهزة بنفس قدرات Java (Priority: P1)

كمسؤول MDM، أريد البحث عن الأجهزة بكل الفلاتر المتقدمة (الحالة، إصدار Android، وضع MDM، إلخ) ورؤية تفاصيل الجهاز (البطارية، النموذج، التطبيقات المثبتة) كما في الباكند القديم، حتى لا أفقد رؤية تشغيلية عند الانتقال إلى Go.

**Why this priority**: الفرونت يرسل فلاتر Java كاملة؛ Go يتجاهل معظمها اليوم — تجربة «أجهزة» ناقصة فوراً.

**Independent Test**: تنفيذ نفس طلب بحث من React مع فلاتر متعددة يُرجع مجموعة فرعية متسقة مع Java على نفس قاعدة البيانات؛ فتح تفاصيل جهاز يعرض حقول telemetry الأساسية.

**Acceptance Scenarios**:

1. **Given** قائمة أجهزة وفلتر `status` أو `launcherVersion`، **When** أبحث من لوحة الأجهزة، **Then** النتائج تُصفّى بنفس منطق Java (أو سلوك موثّق مكافئ).
2. **Given** جهاز مسجّل ببيانات telemetry، **When** أفتح تفاصيل الجهاز، **Then** أرى على الأقل مستوى البطارية والنموذج وإصدار Android إن وُجدت في قاعدة البيانات.
3. **Given** فرز حسب `lastUpdate`، **When** أطلب الصفحة الأولى، **Then** الترتيب يطابق توقع لوحة التحكم (الأحدث أولاً).

---

### User Story 2 - Plugins مراقبة وتصدير كامل (Priority: P1)

كمسؤول يستخدم إضافات deviceinfo و devicelog، أريد البحث المتقدم، التصدير، وقواعد السجل لكل جهاز كما في Java، دون الاعتماد على الباكند القديم.

**Why this priority**: 5+ endpoints REST ناقصة صراحة في Go — أي شاشة plugin قديمة أو تكامل مستقبلي يفشل.

**Independent Test**: كل مسار مذكور في §5 من `JAVA-GO-BACKEND-GAPS.md` (deviceinfo export/search/device/settings per device؛ devicelog export/rules) يُرجع استجابة envelope صحيحة وليس «route not found».

**Acceptance Scenarios**:

1. **Given** بيانات deviceinfo لجهاز، **When** أطلب تصدير deviceinfo، **Then** أستلم ملفاً أو payload قابلاً للتنزيل بنفس البنية التي يتوقعها العميل.
2. **Given** جهاز له قواعد devicelog، **When** أطلب `rules` لهذا الجهاز، **Then** أستلم قائمة القواعد المطبّقة.
3. **Given** بحث deviceinfo لجهاز واحد، **When** أنفّذ `search/device`، **Then** أستلم سجل الجهاز المطابق.

---

### User Story 3 - امتثال وتشغيل tenant (Priority: P2)

كمسؤول نظام (super admin) ومسؤول امتثال، أريد إنشاء عميل جديد بموارد افتراضية، سجلات تدقيق تلقائية للعمليات الحساسة، ومزامنة وكيل تدعم حقول الإضافات المفعّلة.

**Why this priority**: يؤثر على onboarding عملاء جدد وامتثال — أقل من الأجهزة اليومية لكن حرج قبل إيقاف Java.

**Independent Test**: إنشاء customer عبر API يُنتج تكوينات/أجهزة قالب؛ حذف جهاز يُنشئ سجل audit؛ مزامنة وكيل مع plugin مفعّل تُرجع حقولاً إضافية في الاستجابة.

**Acceptance Scenarios**:

1. **Given** super admin ينشئ customer جديد، **When** يكتمل الحفظ، **Then** يوجد على الأقل تكوين/هيكل افتراضي يعادل Java (أو ما هو موثّق في parity customers).
2. **Given** عملية حذف جهاز ناجحة، **When** تُنفَّذ، **Then** يُسجَّل حدث audit تلقائياً (ليس فقط عبر API بحث يدوي).
3. **Given** plugin يوفّر توسيع مزامنة، **When** وكيل يزامن، **Then** استجابة المزامنة تحتوي الحقول الإضافية المتوقعة.

---

### User Story 4 - ملفات وتخزين موثوق (Priority: P2)

كمسؤول محتوى ووكيل جهاز، أريد رفع ملفات التكوين مع احترام الحصة، تتبع الملفات المرفوعة في قاعدة البيانات، وتحميل الملفات عبر مسار الوكيل `/files/*` كما في Java.

**Why this priority**: الوكلاء تعتمد على URLs ثابتة للملفات؛ quota يمنع امتلاء القرص.

**Independent Test**: رفع ملف يتجاوز الحصة يُرفض برسالة envelope؛ وكيل يحمّل ملفاً عبر URL عام؛ ربط ملف بتكوين يُحدّث `uploadedfiles`.

**Acceptance Scenarios**:

1. **Given** عميل عند حد التخزين، **When** أرفع ملف تكوين، **Then** أستلم خطأ quota متوافقاً مع Headwind.
2. **Given** ملف مرفوع مرتبط بتكوين، **When** يطلب الوكيل الملف عبر مسار الملفات العام، **Then** يحمّل المحتوى بنجاح.
3. **Given** رفع ناجح، **When** أحفظ التكوين، **Then** يُسجَّل الملف في سجل الملفات المرفوعة (إن كان Java يفعل ذلك).

---

### User Story 5 - وحدات عامة وتحديثات (Priority: P3)

كمسؤول ووكيل، أريد إرسال إحصائيات استخدام، الوصول لفيديوهات مساعدة (إن كانت مستخدمة)، وتطبيق تحديثات المنتج من مسار التحديثات دون Java.

**Why this priority**: وحدتان REST غير موجودتين (`stats`, `videos`)؛ تحديثات APK ناقصة.

**Independent Test**: `PUT` إحصائيات يُقبل ويُخزَّن؛ فيديو تدريبي يُرفع ويُحمّل (إن كان في النطاق)؛ فحص تحديث يؤدي إلى تنزيل/تطبيق عند التكوين.

**Acceptance Scenarios**:

1. **Given** عميل يرسل usage stats، **When** يُرسل إلى المسار العام للإحصائيات، **Then** تُحفظ دون خطأ 404.
2. **Given** سياسة تحديث مُفعّلة، **When** مسؤول يطبّق تحديثاً من لوحة التحديثات، **Then** يكتمل التنزيل/التطبيق أو يُوثَّق بديل مقبول للتشغيل.
3. **Given** فيديو مساعد مُرفوع (إن كان في النطاق)، **When** مستخدم يطلب التحميل، **Then** يستلم الملف.

---

### Edge Cases

- بيئة بدون MQTT: الإشعارات تعمل عبر طابور DB + polling (منجز في Phase 9) — لا يفشل حفظ التكوين.
- plugin معطّل في الإعدادات: لا تُسجَّل مساراته ولا مهامها الخلفية.
- Mailchimp عند التسجيل/إنشاء عميل: خارج النطاق — لا يُحجب باقي التدفق.
- واجهة Angular القديمة لـ plugin xtra: خارج النطاق.
- impersonate مستخدم ومسارات superadmin للمستخدمين: خارج النطاق ما لم يُطلب لاحقاً (React لا تستخدمها).
- قاعدة بيانات فارغة أو tenant جديد: bootstrap يجب ألا يفشل صامتاً.

---

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Phase**: إكمال Phase 9 وما بعده — إغلاق الفجوات دون كسر Phases 1–8.
- **Target**: `serverBackendGo/internal/modules/` — وحدات جديدة `stats`, `videos` عند الحاجة؛ توسيع `devices`, `customers`, `configfiles`, `files`, `sync`, `plugins/*`, `updates`, `summary`, `signup` (اختياري).
- **Java references**: جميع `*Resource.java` في [`backend/`](../../backend/) — جدول التطابق في [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) §3–§5.
- **REST paths**: الحفاظ على `/rest/...` الحالية؛ لا تغيير breaking للفرونت React.
- **Parity**: تحديث `serverBackendGo/docs/parity/<module>.md` وملفات الجذر `JAVA-GO-*.md` عند إغلاق كل فجوة.
- **Layers**: سلوك في `application/` + `port/`؛ HTTP/DB في `adapter/` فقط.

### Functional Requirements

#### مكتمل مسبقاً (تحقق regression فقط)

- **FR-000**: النظام MUST الإبقاء على push notifier عند حفظ التكوين وإشعار `applicationSettings/notify` وcron جدولة plugin push و`icon-files` كما في Phase 9.

#### P1 — أجهزة و plugins

- **FR-001**: النظام MUST تطبيق فلاتر بحث الأجهزة التي يرسلها الفرونت (`status`, `androidVersion`, `mdmMode`, `kioskMode`, `launcherVersion`, `installationStatus`, نطاقات التواريخ, `sortBy`, `sortDir`, وغيرها الموثّقة في Java `DeviceSearchRequest`).
- **FR-002**: النظام MUST إرجاع تفاصيل جهاد غنية (على الأقل telemetry الأساسية: battery, model, androidVersion) عند `GET /private/devices/number/{number}`.
- **FR-003**: النظام MUST تنفيذ `POST .../deviceinfo/private/search/device`, `POST .../deviceinfo/private/export`, `GET .../deviceinfo-plugin-settings/device/{deviceNumber}`.
- **FR-004**: النظام MUST تنفيذ `POST .../devicelog/log/private/search/export`, `GET .../devicelog/log/rules/{deviceNumber}`.

#### P2 — خلفية وملفات وعملاء

- **FR-005**: النظام MUST تسجيل audit تلقائي لعمليات private الحساسة (middleware يعادل `AuditFilter`) مع استثناءات موثّقة.
- **FR-006**: النظام MUST دعم توسيع استجابة المزامنة عبر آلية hooks للإضافات (يعادل `SyncResponseHook` في Java).
- **FR-007**: النظام MUST إكمال إنشاء customer بـ bootstrap (نسخ تكوينات/أجهزة افتراضية) كما في `CustomerResource` Java.
- **FR-008**: النظام MUST فرض quota تخزين (`sizeLimit`) وتسجيل `uploadedfiles` في مسار رفع ملفات التكوين والملفات.
- **FR-009**: النظام MUST توفير تحميل ملفات للوكلاء عبر مسار يعادل `/files/*` في Java.

#### P3 — وحدات عامة وتحسينات

- **FR-010**: النظام MUST تنفيذ `PUT /rest/public/stats` مع حفظ البيانات.
- **FR-011**: النظام MUST تنفيذ `POST` و`GET /rest/public/videos/{fileName}` **أو** توثيق ⊘ صريح في parity إن لم تُستخدم في الإنتاج.
- **FR-012**: النظام MUST إكمال `updates`: تنزيل APK بعيد و`sendStats` غير stub.
- **FR-013**: النظام MUST تحسين إحصائيات لوحة `summary/devices` (charts) عند توفر البيانات في schema.
- **FR-014**: النظام MUST إكمال حقول حفظ `configurations` و`applications` الحرجة للـ QR والتصميم (أو توثيق الحقول المدعومة).

#### حوكمة وتوثيق

- **FR-015**: عند إغلاق فجوة، الفريق MUST تحديث [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) و[`JAVA-GO-MIGRATION-STATUS.md`](../../JAVA-GO-MIGRATION-STATUS.md) وPhase 9 في `MIGRATION.md` إلى **done** عند اكتمال النطاق.
- **FR-016**: النظام MUST تمرير `go test ./...` وsmoke React + agent (enroll, sync, notification poll) قبل إعلان اكتمال النقل.
- **FR-017**: النظام MUST عدم تسجيل routes لوحدات معطّلة بمتغيرات البيئة (نفس نمط المشروع الحالي).

### Key Entities

- **Device search criteria**: معايير تصفية وترتيب لقائمة الأجهزة.
- **Device telemetry**: بيانات حالة الجهاز (بطارية، طراز، إصدار، …).
- **Plugin export bundle**: بيانات مجمّعة للتصدير من deviceinfo/devicelog.
- **Audit event**: سجل عملية مع مستخدم وعميل وطابع زمني.
- **Tenant bootstrap template**: موارد افتراضية عند إنشاء customer.
- **Uploaded file record**: ملف مرفوع مع checksum وحصة تخزين.
- **Usage statistics**: إحصائيات استخدام من عميل أو وكيل.
- **Product update manifest**: إدخال تحديث مع حالة outdated.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: تطابق سلوكي إجمالي مع Java يصل إلى **≥95%** للمسارات المستخدمة من React والوكلاء (مقيساً بعدد بنود ❌/⚠️ في `JAVA-GO-BACKEND-GAPS.md` → 0 حرجة).
- **SC-002**: **100%** من endpoints §5 في `JAVA-GO-BACKEND-GAPS.md` (القائمة الصريحة للناقص) منفَّذة أو ⊘ موثّقة.
- **SC-003**: مسؤول MDM يكمل خلال **30 دقيقة** UAT: login → أجهزة (فلاتر) → تكوين → push → مزامنة وكيل → (اختياري) push مجدول — **بدون Java**.
- **SC-004**: **0** regressions على مسارات Phases 1–8 في smoke الفرونت الحالي.
- **SC-005**: كل وحدة مُعدَّلة لها `docs/parity/<name>.md` بحالة **Done** للمسارات المُنفَّذة.
- **SC-006**: `MIGRATION.md` Phase 9 = **done** وملفات الجذر `JAVA-GO-*.md` محدّثة بتاريخ الإغلاق.

---

## Assumptions

- Postgres schema الحالي (legacy) هو المصدر؛ migrations إضافية فقط عند الحاجة (usage stats, جداول plugin إضافية).
- Java [`backend/`](../../backend/) يبقى مرجع السلوك عند أي اختلاف.
- إشعارات الأجهزة عبر **طابور DB + long polling** كافية للإنتاج؛ MQTT اختياري وخارج النطاق الافتراضي.
- Mailchimp، plugin xtra، `PublicFilesResource`، impersonate مستخدم — خارج النطاق v1.
- يُكمّل هذا المشروع مهام `specs/011-complete-migration-gaps/tasks.md` من T046 فصاعداً (لا يُعاد P0 المنفّذ).
- الفرونت React يبقى العميل الرئيسي؛ لا إلزام بدعم Angular plugin UI القديمة.

---

## Out of Scope (v1)

- واجهة Angular / plugin **xtra**.
- `MailchimpService` واشتراكات التسويق.
- `UserResource` impersonate و `/superadmin/*` (ما لم يُطلب صراحة لاحقاً).
- `PublicFilesResource` legacy.
- استبدال MQTT إلزامياً — فقط إن قررت العمليات لاحقاً.
- WebSocket `/rest/ws/connect` (غير مستخدم في React الحالي).

---

## Deliverables

| المخرج | المسار |
|--------|--------|
| مواصفة هذا المشروع | `specs/012-finish-java-go-backend/spec.md` |
| تتبع الفجوات (يُحدَّث) | [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) |
| حالة الهجرة | [`JAVA-GO-MIGRATION-STATUS.md`](../../JAVA-GO-MIGRATION-STATUS.md) |
| تكامل الفرونت | [`FRONTEND-GO-BACKEND-INTEGRATION.md`](../../FRONTEND-GO-BACKEND-INTEGRATION.md) |
| تنفيذ (لاحقاً عبر `/speckit-plan`) | `serverBackendGo/` + parity docs |

---

## Gap Closure Matrix (مرجع سريع من التحليل)

| الأولوية | النقص في Go | مرجع Java |
|----------|-------------|-----------|
| P1 | فلاتر + telemetry أجهزة | `DeviceResource` |
| P1 | deviceinfo export/search/device/settings | `DeviceInfoResource` |
| P1 | devicelog export/rules | `DeviceLogResource` |
| P2 | AuditFilter | `plugins/audit` |
| P2 | SyncResponseHook | `SyncResource` |
| P2 | customers bootstrap | `CustomerResource` |
| P2 | configfiles quota + uploadedfiles | `ConfigurationFileResource` |
| P2 | `/files/*` static | `FilesResource` servlet |
| P3 | stats | `StatsResource` |
| P3 | videos | `VideosResource` |
| P3 | updates APK + sendStats | `UpdateResource` |
| P3 | summary charts | `SummaryResource` |
