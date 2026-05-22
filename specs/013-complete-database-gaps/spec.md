# Feature Specification: إكمال فجوات قاعدة البيانات Java → Go

**Feature Branch**: `013-complete-database-gaps`

**Created**: 2026-05-21

**Status**: Draft

**Input**: إكمال النقوصات في مخطط قاعدة البيانات للباكند الجديد (`serverBackendGo`) بأفضل الممارسات، استناداً إلى [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md).

**Related**:

- تحليل الفجوات: [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md)
- فجوات REST/سلوك: [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md)
- مواصفة REST المتبقية: [`specs/012-finish-java-go-backend/spec.md`](../012-finish-java-go-backend/spec.md)
- Migrations الحالية: [`serverBackendGo/db/migrations/`](../../serverBackendGo/db/migrations/)

**Baseline (موجود في Go):**

- Migrations `000001`–`000010` + `000008_devices_search_extras` (جداول MDM أساسية، plugins أساسية، `settingsjson` على `configurations`).
- ~26 جدولاً أساسياً مطابقاً لـ Java؛ **~6 جداول حرجة** و**~13 جدول plugin اختياري** ناقصة (انظر Gap Matrix أدناه).

**هدف المشروع:** أن تصبح قاعدة Go **قابلة لاستبدال قاعدة Java** للمسارات التي يستخدمها React والوكلاء، مع migrations آمنة قابلة للتكرار، وفهارس وقيود متسقة، ومسار اختياري لاستيراد dump Java دون فقد بيانات التكوين.

---

## Gap Matrix (من التحليل)

| الأولوية | العنصر الناقص | التأثير على المستخدم/النظام |
|----------|---------------|------------------------------|
| **P1** | جدول `devicestatuses` | فلاتر حالة التثبيت، ترتيب أعمدة التطبيقات/الملفات، لوحة Summary |
| **P1** | جدول `userrolesettings` | إظهار/إخفاء أعمدة قائمة الأجهزة حسب الدور (React) |
| **P2** | `configurationapplicationparameters` | سياسة `skipVersionCheck` لكل تطبيق في التكوين |
| **P2** | `usagestats` | تتبع استخدام الخادم (إحصائيات تشغيلية) |
| **P2** | أعمدة: `applicationversions.apkhash`, `configurationapplications.remove`, `settings` (عرض أعمدة، `newdevicegroupid`, …) | تكامل تطبيقات/إعدادات متقدمة |
| **P2** | ترحيل بيانات `configurations` (أعمدة Java → `settingsjson`) | نشر من dump Java دون كسر التكوينات |
| **P3** | `trialkey`, جداول temp للرفع، جداول plugins اختيارية (GPS/WiFi، devicelocations، …) | ميزات اختيارية أو plugins غير مفعّلة في Go |

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - حالة تثبيت الأجهزة (Priority: P1)

كمسؤول MDM، أريد أن تُخزَّن حالة تثبيت التطبيقات وملفات التكوين لكل جهاز في قاعدة البيانات كما في Java، حتى تعمل فلاتر «حالة التثبيت» ولوحة الملخص بنفس المنطق بعد الانتقال إلى Go.

**Why this priority**: بدون `devicestatuses` تبقى فلاتر البحث وSummary مبسّطة أو غير دقيقة رغم تحسينات REST في 012.

**Independent Test**: بعد migration، لكل جهاز في بيئة الاختبار يوجد صف في `devicestatuses` (أو قيم افتراضية متسقة)؛ فلتر `installationStatus` في بحث الأجهزة يغيّر النتائج بشكل قابل للملاحظة.

**Acceptance Scenarios**:

1. **Given** جهازان بقيم `applicationsStatus` مختلفة، **When** أبحث بفلتر حالة التثبيت، **Then** أستلم فقط الأجهزة المطابقة.
2. **Given** تغيير حالة تثبيت بعد مزامنة وكيل (أو إعادة حساب)، **When** تُحدَّث السجلات، **Then** تنعكس القيم في البحث والملخص دون إعادة تشغيل الخادم.
3. **Given** قاعدة فارغة بعد migrate، **When** تُشغَّل seed/backfill، **Then** لا تفشل استعلامات JOIN على جدول الحالة.

---

### User Story 2 - تخصيص أعمدة قائمة الأجهزة (Priority: P1)

كمسؤول يريد ضبط ما يراه المستخدمون في جدول الأجهزة، أريد إعدادات عرض الأعمدة مرتبطة بالدور والعميل (مثل Java)، حتى تتطابق لوحة React مع سياسة المؤسسة.

**Why this priority**: React يعتمد على `userrolesettings` / `columnDisplayed*` في Java؛ غياب الجدول يعني أعمدة افتراضية ثابتة أو غير محفوظة.

**Independent Test**: تغيير إعداد عمود (مثلاً إخفاء IMEI) يُحفظ ويُسترجع بعد إعادة تسجيل الدخول لنفس الدور/العميل.

**Acceptance Scenarios**:

1. **Given** دور «Organization Admin» وعميل نشط، **When** أحدّث أعمدة العرض المرئية، **Then** تُحفظ القيم في `userrolesettings` وتُطبَّق في واجهة الأجهزة.
2. **Given** أدوار متعددة لنفس العميل، **When** أضبط كل دوراً، **Then** الإعدادات مستقلة (قيد فريد role + customer).
3. **Given** عميل جديد، **When** يُنشأ أول مستخدم/دور، **Then** تُنشأ صفوف إعدادات افتراضية معقولة (مكافئة لـ seed Java).

---

### User Story 3 - تكامل تطبيقات التكوين والإحصائيات (Priority: P2)

كمسؤول محتوى ومسؤول تشغيل، أريد جداول معاملات التطبيقات (`skipVersionCheck`) وجدول إحصائيات الاستخدام، مع أعمدة ناقصة على إصدارات التطبيقات، حتى تكتمل دورة ترقية APK والتقارير التشغيلية.

**Why this priority**: يدعم محرر التكوينات ووحدة stats المخططة في 012 دون schema ad-hoc.

**Independent Test**: حفظ معامل تطبيق لزوج (configuration, application) يُ persist؛ إرسال إحصائية استخدام يُخزَّن في `usagestats`؛ عمود `apkhash` يقبل قيمة بعد رفع APK.

**Acceptance Scenarios**:

1. **Given** تكوين وتطبيق مرتبطان، **When** أفعّل/أعطّل تخطي فحص الإصدار، **Then** يُحفظ في `configurationapplicationparameters`.
2. **Given** نشر Go يرسل heartbeat إحصائي، **When** تُستقبل البيانات، **Then** تُسجَّل أو تُحدَّث صفاً لكل (يوم، instance) دون تكرار غير مقصود.
3. **Given** إصدار تطبيق مرفوع، **When** يُحسب hash الملف، **Then** يُخزَّن في `applicationversions` للتحقق لاحقاً.

---

### User Story 4 - ترحيل من قاعدة Java موجودة (Priority: P2)

كمسؤول نشر ينقل بيئة قائمة من Java إلى Go، أريد مسار ترحيل بيانات آمن لحقول التكوين المخزّنة كأعمدة في Java وكمفتاح JSON في Go، دون فقد سياسات MDM.

**Why this priority**: فجوة تصميمية (`settingsjson` vs أعمدة منفصلة) — بدون ترحيل، dump Java لا يعمل مع Go كما هو.

**Independent Test**: على نسخة DB تحتوي أعمدة `configurations` من Java، تشغيل خطوة ترحيل (migration أو script موثّق) ينتج `settingsjson` غير فارغ للتكوينات النشطة؛ التحقق اليدوي لعينة ≥3 تكوينات.

**Acceptance Scenarios**:

1. **Given** dump Java بأعمدة kiosk/network مملوءة، **When** أُشغّل ترحيل البيانات، **Then** تظهر نفس القيم في `settingsjson` بمفاتيح يتوقعها محرر التكوين.
2. **Given** تكوين بدون أعمدة legacy (بيئة Go خالصة)، **When** أُشغّل migrate، **Then** لا يحدث خطأ والقيم الافتراضية تبقى سليمة.
3. **Given** ترحيل فاشل جزئياً، **When** أُراجع السجل، **Then** أعرف أي معرفات تكوين تحتاج تدخل يدوي.

---

### User Story 5 - Plugins و جداول اختيارية (Priority: P3)

كمسؤول يفعّل plugins متقدمة (مثل telemetry WiFi/GPS)، أريد أن تُضاف جداول plugin فقط عند تفعيل الوحدة المقابلة في Go، دون تحميل schema كامل لكل التثبيتات.

**Why this priority**: ~13 جدولاً في Java ليست مطلوبة لـ MVP؛ يُدار كنطاق منفصل per-plugin.

**Independent Test**: تفعيل plugin deviceinfo الموسّع ينشئ جداول WiFi/GPS؛ تثبيت Go بدون plugin لا ينشئ تلك الجداول.

**Acceptance Scenarios**:

1. **Given** plugin غير مفعّل في التكوين، **When** أُشغّل migrations الأساسية فقط، **Then** لا تُنشأ جداول plugin الاختيارية.
2. **Given** قرار تفعيل plugin devicelocations لاحقاً، **When** تُضاف migration فرعية، **Then** الجداول الثلاثة للمواقع تُنشأ مع FK صحيحة دون كسر migrations سابقة.

---

### Edge Cases

- تشغيل `migrate` على قاعدة فيها جزء من الجداول يدوياً (idempotent `IF NOT EXISTS`).
- downgrade migration يُسقِط الجداول/الأعمدة بترتيب عكسي دون كسر FK.
- عميل بدون أجهزة: `devicestatuses` فارغ — البحث لا يفشل.
- backfill `devicestatuses` من `infojson` عند غياب حاسبة Java: قيم افتراضية `FAILURE` / `OTHER` موثّقة.
- استيراد dump بأسماء أعمدة camelCase من Liquibase — التطبيع إلى lowercase في PostgreSQL.
- حجم `settingsjson` كبير جداً لسياسة واحدة — حد معقول أو تقسيم مفاتيح (document in assumptions).

---

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **نطاق**: `serverBackendGo/db/migrations/` فقط للـ schema؛ تحديث repositories في `internal/modules/*` حيث تُستهلك الجداول الجديدة (devices, settings/roles, configurations, applications, stats).
- **مرجع Java**: Liquibase [`backend/server/src/main/resources/liquibase/db.changelog.xml`](../../backend/server/src/main/resources/liquibase/db.changelog.xml) وجداول plugins في [`backend/plugins/`](../../backend/plugins/).
- **أداة الهجرة**: golang-migrate (ملف `.up.sql` / `.down.sql` لكل إصدار، ترقيم متسلسل من `000011`).
- **Parity docs**: تحديث [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) وملخص في [`serverBackendGo/docs/MIGRATION.md`](../../serverBackendGo/docs/MIGRATION.md).
- **طبقات**: لا منطق أعمال في SQL إلا seeds/backfill محدودة؛ حساب الحالة يبقى في `application/` حيث أمكن.

### Functional Requirements

- **FR-001**: النظام MUST توفير migration `devicestatuses` مطابقة لـ Java (جهاز واحد = صف واحد، `applicationsStatus`, `configFilesStatus`) مع FK إلى `devices`.
- **FR-002**: النظام MUST توفير migration `userrolesettings` مع أعمدة `columnDisplayed*` المستخدمة في React وقيد فريد `(roleId, customerId)` وseed افتراضي للأدوار القياسية.
- **FR-003**: النظام MUST توفير migration `configurationapplicationparameters` مع قيد فريد `(configurationId, applicationId)`.
- **FR-004**: النظام MUST توفير migration `usagestats` مع قيد فريد `(ts, instanceId)` كما في Java.
- **FR-005**: النظام MUST إضافة الأعمدة الناقصة P2: `applicationversions.apkhash`, `configurationapplications.remove` (و`longtap` إن مطلوب للـ UI), وأعمدة `settings` الحرجة (`newdevicegroupid`, `phonenumberformat`, أسماء الخصائص المخصصة) أو توثيق ⊘ صريح إن بقيت في `settingsjson` فقط.
- **FR-006**: كل migration MUST تكون قابلة للإلغاء (`down`) وتستخدم `IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS` حيث ينطبق.
- **FR-007**: النظام MUST توفير backfill اختياري لـ `devicestatuses` للأجهزة الموجودة بعد إنشاء الجدول.
- **FR-008**: النظام MUST توثيق مسار ترحيل `configurations` (أعمدة Java → `settingsjson`) كخطوة migration منفصلة أو script مع قائمة مفاتيح JSON.
- **FR-009**: النظام MUST إنشاء فهارس على FK وأعمدة البحث الشائعة (`devicestatuses.deviceId`, `usagestats.ts`) وفق أفضل ممارسات PostgreSQL.
- **FR-010**: جداول plugins الاختيارية (§3.3 في Gap doc) MUST تبقى خارج النطاق v1 إلا إذا رُبطت بمواصفة plugin منفصلة — يُوثَّق ⊘ في spec/plan.
- **FR-011**: بعد اكتمال P1+P2، مستودعات Go التي تعتمد على `installationStatus` / Summary MUST تقرأ من `devicestatuses` بدلاً من الاعتماد فقط على `infojson`.
- **FR-012**: تحديث [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) MUST يعكس الحالة الجديدة (✅/⚠️/❌) لكل بند P1–P2.

### Key Entities

- **DeviceStatus**: حالة تثبيت ملفات التكوين والتطبيقات لجهاز واحد؛ يرتبط 1:1 بـ Device.
- **UserRoleSettings**: تفضيلات عرض أعمدة قائمة الأجهزة لكل (Role, Customer).
- **ConfigurationApplicationParameter**: علامة `skipVersionCheck` لزوج (Configuration, Application).
- **UsageStats**: لقطة استخدام يومية/مثيل للخادم (أجهزة، موارد، إصدار).
- **Configuration (extended)**: سياسات MDM في `settingsjson` + أعمدة صريحة للحقول الحرجة في API.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% من بنود P1 في Gap Matrix (جداول `devicestatuses`, `userrolesettings`) تُعلَم ✅ في `JAVA-GO-DATABASE-GAPS.md` بعد التنفيذ.
- **SC-002**: تشغيل `make migrate` من صفر حتى `000014` (أو آخر migration معرّف) ينجح بدون أخطاء على PostgreSQL 14.
- **SC-003**: في بيئة اختبار بجهازين على الأقل، فلتر `installationStatus` يُنتج نتائج مختلفة عند اختلاف `devicestatuses` (قابل للتحقق في ≤5 دقائق UAT).
- **SC-004**: حفظ إعدادات أعمدة الدور يبقى بعد إعادة تشغيل الخادم وإعادة تسجيل الدخول (اختبار واحد لكل دور admin).
- **SC-005**: ترحيل عينة من 3 تكوينات من نموذج أعمدة Java إلى `settingsjson` يحافظ على ≥90% من المفاتيح المرئية في محرر التكوين (مراجعة يدوية).
- **SC-006**: لا توجد migrations بدون ملف `down` مقابل لنفس الإصدار.

---

## Assumptions

- PostgreSQL 14+ كما في `docker-compose.yml` الحالي؛ أسماء الجداول lowercase غير مقتبسة.
- `configurations.settingsjson` يبقى المصدر الرئيسي لسياسات MDM الجديدة؛ الأعمدة المنفصلة في Java تُرحَّل إلى JSON وليس بالعكس (إلا الحقول التي يقرأها Go صراحة اليوم).
- Feature `012-finish-java-go-backend` قد يُنفَّذ بالتوازي؛ 013 يركز على **schema** ويُحدّث repos عند الحاجة لتفعيل الجداول.
- Plugins اختيارية (devicelocations, photo, knox, …) **خارج نطاق v1** ما لم يُطلب صراحة لاحقاً.
- `trialkey` و جداول temp للرفع **P3** — لا تمنع إيقاف Java للعمليات اليومية.
- دوال SQL Java (`mdm_device_launcher_version`, …) **لا تُنقل** كإجراءات مخزنة؛ المنطق يبقى في Go/SQL queries (موثّق في plan).

---

## Out of Scope

- إعادة كتابة كامل Liquibase Java إلى Go دفعة واحدة (~55 جدول).
- دعم MySQL أو قواعد أخرى غير PostgreSQL.
- Mailchimp / WebSocket / MQTT schema.
- تغيير مسارات REST (يغطيها 012).
- إنشاء وحدات Go جديدة بالكامل لكل plugin اختياري — فقط schema عند الحاجة.

---

## Dependencies

- [`JAVA-GO-DATABASE-GAPS.md`](../../JAVA-GO-DATABASE-GAPS.md) كمصدر حقيقة للفجوات.
- Migrations موجودة حتى `000010` + `000008_devices_search_extras`.
- Constitution: [`.specify/memory/constitution.md`](../../.specify/memory/constitution.md) — migrations versioned، tests للـ repos المتأثرة.
