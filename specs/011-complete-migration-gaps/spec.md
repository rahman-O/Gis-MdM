# Feature Specification: Complete Java→Go Migration Gaps

**Feature Branch**: `011-complete-migration-gaps`

**Created**: 2026-05-21

**Status**: Draft

**Input**: العمل على إكمال غير المنقول والجزئي وغيرهم لإكمال فجوات النقل، استناداً إلى [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md).

**Related**: تحليل الفجوات (spec 010) — هذا المشروع هو **التنفيذ** لإغلاق الفجوات المُوثَّقة هناك.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - إشعارات الأجهزة الفورية (Priority: P1)

كمسؤول MDM، عند حفظ تكوين أو طلب إشعار جهاز، أريد أن يصل تنبيه فعلي للوكيل (مثل Java) حتى تُطبَّق التغييرات دون انتظار مزامنة يدوية.

**Why this priority**: `NoopPush` اليوم يعني أن React يعمل لكن الأجهزة لا تتلقى push — أكبر فجوة تشغيلية (P0 في تحليل الفجوات).

**Independent Test**: بعد حفظ configuration، الجهاز المسجّل يستقبل إشعاراً عبر قناة الإشعارات الحالية (نفس مسار Java) خلال دقيقة واحدة في بيئة اختبار.

**Acceptance Scenarios**:

1. **Given** تكوين مُحدَّث ومجموعة أجهزة مرتبطة، **When** أحفظ التكوين من لوحة التحكم، **Then** تُسجَّل رسالة push وتصل للجهاز (أو تظهر في طابور الإشعارات للوكيل).
2. **Given** طلب `applicationSettings/notify` لجهاز، **When** يُنفَّذ من الواجهة، **Then** لا يعود الرد نجاحاً صامتاً بدون إرسال.

---

### User Story 2 - الرسائل المجدولة تُرسل تلقائياً (Priority: P1)

كمسؤول يستخدم plugin push، أريد أن تُنفَّذ مهام `plugin_push_schedule` في وقتها كما في Java، دون تشغيل يدوي.

**Why this priority**: CRUD الجدولة موجود في Go لكن cron غير موجود — سلوك مكسور للمستخدمين الحاليين.

**Independent Test**: إنشاء مهمة مجدولة خلال دقيقتين من «الآن» → تظهر رسالة push أو سجل إرسال مطابق لـ Java.

**Acceptance Scenarios**:

1. **Given** مهمة مجدولة مستحقة، **When** يمر وقت التنفيذ، **Then** تُرسل الرسالة وتُحدَّث حالة المهمة.
2. **Given** مهمة ملغاة أو محذوفة، **When** يحين موعدها، **Then** لا تُرسل.

---

### User Story 3 - رفع أيقونات وملفات مكتمل (Priority: P2)

كمسؤول محتوى، أريد رفع ملف أيقونة (`icon-files`) وربطه بالتطبيقات كما في الباك اند القديم، مع احترام حدود التخزين حيث تُطبَّق في Java.

**Why this priority**: `IconFileResource` غير منقول؛ `icons` يدير metadata فقط.

**Independent Test**: `POST /rest/private/icon-files` يرجع نفس شكل الاستجابة ويُستخدم في تطبيق لاحقاً.

**Acceptance Scenarios**:

1. **Given** ملف صورة صالح، **When** أرفعه عبر icon-files، **Then** يُخزَّن ويُعاد معرف/مسار قابل للربط في `icons` أو التطبيقات.
2. **Given** تجاوز حصة التخزين (إن مُفعَّل)، **When** أرفع ملفاً، **Then** رسالة خطأ متوافقة مع Headwind envelope.

---

### User Story 4 - إضافات deviceinfo و devicelog كاملة (Priority: P2)

كمسؤول يراقب الأجهزة، أريد البحث المتقدم، التصدير، وقواعد السجل لكل جهاز كما في واجهة plugin القديمة.

**Why this priority**: مسارات export و rules ناقصة — الواجهة قد تفشل عند تفعيل هذه الشاشات.

**Independent Test**: كل endpoint المدرج في §8 من تحليل الفجوات (deviceinfo/devicelog) يُرجع 200/404 متوقعاً وليس 404 من Gin router.

**Acceptance Scenarios**:

1. **Given** جهاز ببيانات telemetry، **When** أطلب export deviceinfo، **Then** أستلم ملف/بيانات بنفس البنية المتوقعة من React أو Angular plugin.
2. **Given** جهاز بقواعد devicelog، **When** أطلب `rules/{deviceNumber}`، **Then** أستلم القواعد المطبّقة.

---

### User Story 5 - تدقيق وتزامن وعملاء (Priority: P3)

كمسؤول امتثال وتشغيل، أريد سجلات تدقيق تلقائية عند العمليات الحساسة، مزامنة وكيل تدعم hooks للإضافات، وإنشاء عميل جديد بنسخ افتراضيات معقولة.

**Why this priority**: P2 في خطة الفجوات — مهم لكن أقل من push الفوري.

**Independent Test**: طلب خاص مُسجَّل في audit؛ sync response يحتوي حقول plugin عند تفعيل plugin؛ `PUT customers` ينشئ بنية tenant أولية.

**Acceptance Scenarios**:

1. **Given** عملية حذف جهاز، **When** تُنفَّذ بنجاح، **Then** يُكتب سجل audit (ليس بحثاً فقط).
2. **Given** plugin مفعّل يوفّر SyncResponseHook، **When** وكيل يزامن، **Then** الحقول الإضافية موجودة في الاستجابة.
3. **Given** super admin ينشئ customer، **When** يكتمل PUT، **Then** توجد تكوينات/أجهزة افتراضية أو ما يعادل سلوك Java الموثّق.

---

### User Story 6 - وحدات عامة ووكلاء (Priority: P4)

كمسؤول ووكيل، أريد إحصائيات استخدام، فيديوهات مساعدة (إن لزم)، تحديثات APK كاملة، وبحث أجهزة غني كما في Java.

**Why this priority**: P3 في تحليل الفجوات — أقل استخداماً من React الحالي.

**Independent Test**: `PUT /rest/public/stats` يقبل payload؛ `GET/POST videos` يعملان؛ `updates` يحمّل APK عند التكوين؛ بحث أجهزة يدعم الفلاتر المذكورة في parity.

**Acceptance Scenarios**:

1. **Given** وكيل يرسل usage stats، **When** `PUT /rest/public/stats`، **Then** تُحفظ في قاعدة البيانات.
2. **Given** إعداد تحديث بعيد، **When** يطلب الوكيل التحديث، **Then** يمكن تنزيل/تطبيق APK حسب Java (أو documented equivalent).

---

### Edge Cases

- FCM/MQTT غير مُكوَّن في البيئة → تسجيل تحذير وعدم فشل حفظ التكوين (سلوك Java: best-effort push).
- plugin معطّل عبر `ENABLED_PLUGINS` → لا تُسجَّل routes ولا cron لذلك plugin.
- Mailchimp و `PublicFilesResource` و plugin **xtra** → خارج النطاق (لا تُنفَّذ إلا بطلب صريح لاحقاً).
- impersonate / superadmin user endpoints → خارج النطاق ما لم يُستخدمها React (كما في parity users).

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Phase**: Post-Phase 8 hardening — إكمال parity دون كسر Phases 1–8.
- **Modules touched** (متوقع): `devices`, `configurations`, `files`, `push`, `plugins/push`, `icons` أو `files`, `plugins/deviceinfo`, `plugins/devicelog`, `plugins/audit`, `sync`, `customers`, `updates`, `signup` (اختياري), وحدات جديدة `stats` / `videos` إن لزم.
- **Java references**: جدول §4–§7 في [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md).
- **REST**: الحفاظ على `/rest/...` الحالية؛ لا تغيير breaking دون توثيق.
- **Parity**: تحديث `serverBackendGo/docs/parity/<module>.md` و`JAVA-GO-MIGRATION-GAP-ANALYSIS.md` عند إغلاق كل فجوة.
- **Layers**: كل سلوك جديد في `application/` + `port/`؛ IO في `adapter/` فقط.

### Functional Requirements

#### P0 — Push وتكاملات حرجة

- **FR-001**: النظام MUST استبدال `NoopPush` في `devices` و`configurations` (و`files` عند ربط configurations) بمُرسِل push يعادل Java (`PushSender` / طابور notifications).
- **FR-002**: النظام MUST دعم إرسال push عبر `POST /rest/private/push` بنفس العقد الحالي للواجهة.
- **FR-003**: النظام MUST تشغيل worker دوري لـ `plugin_push_schedule` يعادل `PushScheduleTaskModule` (استعلام مهام مستحقة + إرسال + تحديث حالة).
- **FR-004**: النظام MUST تسجيل فشل الإرسال دون إرجاع خطأ للمستخدم عند حفظ تكوين (ما لم يُحدد خلاف ذلك في Java).

#### P1 — REST غير المنقول + plugins endpoints

- **FR-005**: النظام MUST تنفيذ `POST /rest/private/icon-files` (multipart، resize إن وُجد في Java) وربط النتيجة بـ `icons`/التطبيقات.
- **FR-006**: النظام MUST تنفيذ مسارات deviceinfo الناقصة: `private/search/device`, `private/export`, `deviceinfo-plugin-settings/device/{deviceNumber}`.
- **FR-007**: النظام MUST تنفيذ مسارات devicelog الناقصة: `private/search/export`, `log/rules/{deviceNumber}`.
- **FR-008**: النظام MUST تحديث parity docs لكل module أعلاه عند الاكتمال.

#### P2 — خلفية وسلوك جزئي

- **FR-009**: النظام MUST تسجيل audit تلقائياً لطلبات private الحساسة عبر middleware يعادل `AuditFilter` (مع استثناءات health/static موثّقة).
- **FR-010**: النظام MUST دعم registry لـ `SyncResponseHook` بحيث plugins توسّع استجابة `/rest/public/sync` دون تعديل core sync لكل plugin.
- **FR-011**: النظام MUST إكمال `customers` PUT لإنشاء tenant بنسخ افتراضيات (أجهزة/تكوينات) كما في `CustomerResource` Java.
- **FR-012**: النظام MUST فرض `sizeLimit` / `uploadedfiles` في `configfiles` و`files` حيث تطبّق Java quota.
- **FR-013**: النظام MUST إثراء `devices` POST `/search` بالفلاتر `mdmMode`, `launcherVersion`, `deviceStatuses` وإثراء النتائج حيث يعتمد React على Java.

#### P3 — وحدات عامة ووكلاء

- **FR-014**: النظام MUST تنفيذ `PUT /rest/public/stats` مع persistence في جدول usage stats (Java `UsageStatsDAO`).
- **FR-015**: النظام MUST تنفيذ `POST` و`GET /rest/public/videos/{fileName}` إن كانت الواجهة أو التوثيق يشيران لاستخدامها؛ وإلا توثيق ⊘ صريح في parity.
- **FR-016**: النظام MUST إكمال `updates`: تنزيل APK بعيد و`sendStats` غير stub عند تفعيل الإعداد.
- **FR-017**: النظام MUST تحسين `summary` charts عند توفر جداول `devicestatuses` أو ما يعادلها في schema.
- **FR-018**: النظام MUST توفير تحميل ملفات الوكيل `/files/*` (static أو handler) متوافق مع Java servlet.

#### حوكمة وتوثيق

- **FR-019**: عند إغلاق فجوة، الفريق MUST تحديث [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md) من ⚠️/❌ إلى ✅.
- **FR-020**: النظام MUST عدم تسجيل routes لوحدات معطّلة بـ env flags (نفس نمط Phase 8).

### Key Entities

- **Push notification**: رسالة موجهة لجهاز/مجموعة؛ ترتبط بطابور notifications و/أو FCM.
- **Scheduled push task**: صف في `plugin_push_schedule` مع وقت تنفيذ وحالة.
- **Icon file upload**: ملف مرفوع مرتبط بـ customer ومسار تخزين.
- **Device telemetry / log export**: مجموعة بيانات للتصدير من جداول plugin deviceinfo/devicelog.
- **Audit log record**: حدث HTTP/عملية مع user وcustomer وtimestamp.
- **Usage stats**: payload عام من عميل لإحصائيات الاستخدام.
- **Tenant bootstrap**: customer جديد مع موارد افتراضية منسوخة من قالب Java.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% من بنود P0 في §12 من تحليل الفجوات مُغلقة ومُعلَّمة ✅ في `JAVA-GO-MIGRATION-GAP-ANALYSIS.md` خلال دورة التطوير هذه.
- **SC-002**: 100% من endpoints §8 (غير المنقول + الجزئي الناقص) إما منفَّذة أو ⊘ موثّقة بسبب عدم الاستخدام.
- **SC-003**: React smoke الحالي (`frontend` + `serverBackendGo/scripts/dev.sh`) يمر بدون regressions على مسارات Phases 1–8.
- **SC-004**: كل module مُعدَّل له `docs/parity/<name>.md` محدّث بحالة **Done** للمسارات المُنفَّذة.
- **SC-005**: مسؤول MDM يمكنه إكمال سيناريو: حفظ تكوين → push → مزامنة وكيل → (اختياري) رسالة مجدولة — دون الرجوع لباك اند Java.

## Assumptions

- Postgres legacy schema هو المصدر؛ migrations جديدة فقط عند الحاجة (usagestats، GPS tables، إلخ).
- Java `backend/` يبقى مرجع السلوك عند الاختلاف.
- FCM يُفعَّل عبر متغيرات بيئة مثل Java (`FCM_*` أو ما يعادلها في المشروع).
- Mailchimp، xtra Angular UI، `PublicFilesResource` خارج النطاق الافتراضي.
- الأولوية: P0 → P1 → P2 → P3 كما في تحليل الفجوات؛ يمكن تقسيم التنفيذ على عدة PRs ضمن نفس feature branch.

## Out of Scope (v1)

- واجهات Angular القديمة لـ plugins (`webapp` assets) — تبقى على React/proxy.
- `UserResource` impersonate و superadmin (ما لم يُطلب صراحة).
- `MailchimpService` إلا بمتطلب عمل منفصل.
- plugin **xtra** بالكامل.

## Deliverables

| Deliverable | Path |
|-------------|------|
| Gap tracker (يُحدَّث مستمراً) | [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md) |
| This specification | `specs/011-complete-migration-gaps/spec.md` |
| Implementation (later) | `serverBackendGo/internal/modules/...` + parity docs |
