# Feature Specification: Enrollment Routes — Controlled Onboarding Gateway

**Feature Branch**: `021-enrollment-routes-ux`

**Created**: 2026-05-24

**Status**: Draft (clarified 2026-05-24 — 5 decisions recorded in Clarifications)

**Input**: تحسين enrollment-routes — بوابة دخول مستقلة للأجهزة (ليس UI feature ولا امتداد للبرفايل): QR contract، منتقي شجرة بسياق، bootstrap app intent، dialogs احترافية، حذف آمن متعدد الأبعاد.

**Related** (خارج نطاق مفردات هذه المواصفة):

- [`specs/019-profile-hub-ux/spec.md`](../019-profile-hub-ux/spec.md) — Profiles / Assignments (نظام منفصل)
- [`specs/017-device-control-plane/spec.md`](../017-device-control-plane/spec.md) — شجرة الأجهزة، QR عام
- [`specs/018-profile-rollout-ops/spec.md`](../018-profile-rollout-ops/spec.md) — rollout على الشجرة (لا يُذكر داخل enrollment)

---

## Clarifications

### Session 2026-05-24

- Q: عند اختيار عقدة **Inheritable (حاوية)** كـ target node — ماذا يفعل المنتج؟ → A: **Allow with warning** — يُسمح بالحفظ؛ تحذير واضح في المنتقي وفي Overview يبقى ظاهراً حتى يختار المسؤول عقدة leaf (Locked) أو يؤكد بقاء الحاوية صراحة.
- Q: كيف يُنتَج **Pending QR** قبل أول حفظ أو أثناء تعديل غير محفوظ؟ → A: **Client preview** — المعاينة محلياً من قيم النموذج بنفس شكل Enrollment Contract (بدون `qrcodeKey` فعلي حتى الحفظ)؛ **Active QR** من الخادم بعد الحفظ فقط.
- Q: حذف مسار مع **أجهزة مسجّلة تاريخياً** (`historicalEnrolledCount > 0`) — ماذا يحدث؟ → A: **Allow with strengthened confirm** — الحذف مسموح؛ يعرض الأبعاد الثلاثة + **كتابة اسم المسار** للتأكيد النهائي؛ الأجهزة التاريخية تبقى في السجل ولا تُحذف مع المسار.
- Q: كيف يُحلّ **Stable channel** intent إلى إصدار فعلي؟ → A: **Catalog stable flag** — الإصدار المعلَّم stable/recommended في كتالوج التطبيقات؛ يُعاد الحل عند كل Save / توليد Active QR.
- Q: شارة الحالة في الـ header عند تعديل مسار **محفوظ** مع تغييرات غير محفوظة؟ → A: **Active + Unsaved** — تبقى Active؛ شارة ثانوية «Unsaved changes»؛ **Draft** فقط لمسار لم يُحفَظ أبداً.

---

## Core Concept *(mandatory)*

### What an Enrollment Route is

**Enrollment Route** = **controlled onboarding gateway** for device population.

| Domain | Responsibility |
|--------|----------------|
| **Enrollment Route** | Onboarding pipeline only: أين يدخل الجهاز، كيف يُعرَّف، ماذا يُثبَّت للـ bootstrap، حمولة QR |
| **Profile** (منفصل تماماً) | نظام السياسة على الشجرة — **خارج** نطاق enrollment |

### Strict vocabulary rule

داخل **Enrollment Routes** (واجهة، DTO للعميل، وثائق الميزة، رسائل المستخدم):

- **MUST NOT** استخدام أو الإشارة إلى: policy، profile، profile version، سياسة، برفايل، rollout، assignment bump.
- **MUST** الاقتصار على: target node، bootstrap app، device identity mode، QR payload، route status.

أي سلوك يتعلق بما يحدث **بعد** التسجيل (مثل تطبيق إعدادات على الجهاز) يُملكها وحدات **sync / profiles** — **ليست** جزءاً من عقد مسار التسجيل ولا تُوثَّق هنا.

### Enrollment Route knows only

| Concern | Description |
|---------|-------------|
| **Target node** | عقدة شجرة واحدة لاستقبال الأجهزة الجديدة |
| **Bootstrap app** | تطبيق/قناة التثبيت الأولي (intent + resolution) |
| **Device identity mode** | IMEI / Serial / Request |
| **QR payload** | Enrollment Contract (انظر أدناه) |

**كل ما عدا ذلك = خارج النطاق.**

---

## Enrollment Contract (QR Payload) *(mandatory)*

### Principle

QR يمثل **Enrollment Contract** — عقد onboarding مستقل **بدون أي مرجع** لبرفايل أو سياسة أو configuration policy.

### Logical payload fields

| Field | Required | Purpose |
|-------|----------|---------|
| `routeId` | yes | يحدد بوابة التسجيل |
| `targetNodeId` | yes | عقدة الشجرة المستهدفة |
| `mainAppPackage` | yes | حزمة التطبيق للـ bootstrap |
| `mainAppVersion` | yes* | إصدار محدد (*يُحلّ عند الحفظ حسب intent — انظر Bootstrap App) |
| `deviceIdentityMode` | yes | imei / serial / request |
| `bootstrapFlags` | optional | أعلام provisioning اختيارية (مثلاً `create=1`، إعدادات عامة legacy-compatible) |
| `qrcodeKey` | runtime | مفتاح عام للمسح (يُولَّد مع الحفظ الأول) |

**MUST NOT appear in contract**: `profileId`, `profileVersionId`, `policy*`, أي حقل سياسة.

### QR preview states (UI)

| State | When | Production | Source |
|-------|------|------------|--------|
| **Pending QR** | قبل أول حفظ أو أثناء تعديل غير محفوظ | **غير قابل للمسح للإنتاج** — معاينة بصرية فقط | **Client-side**: يُبنى من قيم النموذج الحالية بنفس حقول Enrollment Contract؛ `routeId`/`qrcodeKey` placeholder أو فارغ حتى الحفظ |
| **Active QR** | بعد حفظ تعريف المسار بنجاح | **قابل للمسح** — مفتاح عام فعلي | **Server-authoritative**: `GET .../qr` أو ما يعادله بعد الحفظ |

Pending MUST يتحدث **فوراً** عند تغيير target node أو bootstrap app أو identity mode (بدون round-trip للخادم).

Active MUST يُحدَّث من الخادم بعد كل Save ناجح؛ نسخ/تحميل QR متاح فقط في حالة Active.

معاينة QR **دائماً ظاهرة** في عمود الـ dialog — Pending أو Active حسب الحالة.

---

## UX Architecture *(mandatory)*

### Design goal

تحويل Enrollment Routes من «إعدادات مرتبطة بالبرفايل» إلى:

> **Controlled onboarding gateway for device population**

### Anti-patterns (MUST NOT)

| Anti-pattern | Why |
|--------------|-----|
| أي حقل أو نص «سياسة / برفايل» في شاشة المسار | يكسر الفصل المعماري |
| تمرير `profileVersionId` في DTO طبقة الواجهة | leakage مستقبلي |
| QR بدون عقد حقول موثّق | coupling عشوائي لاحقاً |
| قائمة مسطحة للشجرة بدون سياق عقدة | misplacement |
| حذف بتأكيد يعرض «رقم أجهزة» واحداً فقط | قرار غير واضح |
| إخفاء QR حتى الحفظ فقط | يفقد إحساس النظام الاحترافي |

### Target pattern: Enrollment Routes Hub

```text
List (خفيف)
  → Enrollment Route Dialog (واحد)
       Header: اسم المسار + status (Draft | Active)
       Body 2-column:
         Left:  configuration
         Right: live QR preview (Pending | Active)
       Footer: Save | Delete | Cancel
```

**NOT**: صفحة editor كاملة · nested dialogs · عمود QR يظهر فقط بعد الحفظ

### Dialog layout (required)

```
┌─────────────────────────────────────────────────────────┐
│ Header: Route name + status badge (Draft / Active)      │
├──────────────────────────┬──────────────────────────────┤
│ Left: Configuration      │ Right: Live QR preview       │
│  - Name, description       │  - Pending or Active label │
│  - Target node picker      │  - QR image + copy/download│
│  - Device identity mode    │                            │
│  - Bootstrap app intent    │                            │
├──────────────────────────┴──────────────────────────────┤
│ Footer: [Save] [Delete] [Cancel]                        │
└─────────────────────────────────────────────────────────┘
```

على **mobile**: عمودان → تكديس (config ثم QR)؛ نفس الحالات.

### Modes

| Mode | Left column | Right column | Footer |
|------|-------------|--------------|--------|
| **Create / Edit** | حقول قابلة للتعديل | Pending QR (يتحدث live مع التغييرات) | Save, Cancel |
| **Overview** | قراءة فقط | Active QR | Edit, Delete, Close |

التبديل Overview ↔ Edit **داخل نفس الـ dialog** — لا تنقل URL لصفحة editor.

### Route status badges (header)

| Badge | Meaning |
|-------|---------|
| **Draft** | مسار **لم يُحفَظ أبداً** (Create فقط — لا `routeId` بعد) |
| **Active** | مسار محفوظ مرة على الأقل — QR عام فعّال من آخر Save |
| **Unsaved changes** | شارة **ثانوية** بجانب Active عندما النموذج ≠ آخر حالة محفوظة (Edit على مسار موجود) |

قواعد:

- Edit على مسار محفوظ: Header **Active** + **Unsaved changes** (إن وُجدت تغييرات)؛ عمود QR = **Pending** (معاينة التغييرات)؛ **Active QR السابق يبقى ساري المسح** حتى Save ناجح يحدّثه.
- Create (قبل أول Save): Header **Draft** فقط؛ لا Unsaved؛ عمود QR = Pending.
- Overview (بدون وضع Edit): **Active** فقط؛ عمود QR = Active من الخادم.

---

## Target Node Semantics *(tree)*

### Single node rule

يُختار **عقدة واحدة فقط** (فرع أو أب) كـ `targetNodeId`.

### Node kinds (placement eligibility)

| Kind | Label (UX) | Meaning |
|------|------------|---------|
| **Locked (leaf placement)** | «استقبال أجهزة» | عقدة صالحة كهدف نهائي لتسجيل أجهزة جديدة |
| **Inheritable (container)** | «حاوية توزيع» | عقدة أب — يُسمح باختيارها كهدف مع **تحذير غير مانع** (انظر أدناه) |

المنتقي MUST:

- تمييز نوع العقدة بصرياً (أيقونة/شارة).
- **السماح** باختيار عقد Inheritable (container) مع **تحذير واضح** في المنتقي قبل التأكيد وفي Overview بعد الحفظ.
- التحذير MUST يوضح أن العقدة حاوية وليست ورقة «استقبال أجهزة»؛ يبقى ظاهراً في Overview حتى يغيّر المسؤول إلى عقدة Locked أو يؤكد صراحة الإبقاء على الحاوية (checkbox «أفهم، أبقِ هذا المجلد»).
- **منع** فقط العقد غير الصالحة تقنياً (محذوفة، خارج نطاق العميل، بلا صلاحية) — ليس لأنها container.
- الإبقاء على **اختيار واحد** — لا multi-select.

### Context preview (before confirm)

قبل تأكيد الاختيار، يعرض المنتقي:

| Context | Purpose |
|---------|---------|
| **Full path breadcrumb** | وضوح الموقع في الشجرة |
| **Device count under node** | فهم الحمل الحالي |
| **Heavily loaded warning** | إن تجاوز عتبة قابلة للضبط (افتراض: تنبيه معلوماتي، لا منع إلا إن فرضت الشجرة ذلك) |

---

## Bootstrap App Intent *(not just version picker)*

### Problem

اختيار «رقم إصدار» فقط يسبب أعطال تشغيل عند تغيّر القنوات.

### Intent modes (required)

| Intent | Label (UX) | Resolution rule |
|--------|------------|-----------------|
| **Stable channel** | Recommended / Stable | الإصدار الذي يحمل علامة **stable / recommended** في كتالوج التطبيقات لنفس الحزمة؛ يُعاد الحل عند كل Save وتوليد Active QR |
| **Specific version** | Pinned version | الإصدار المختار صراحة — ثابت حتى يغيّره المسؤول |
| **Latest available** | Latest | أحدث إصدار **منشور/متاح** في الكتالوج (ليس بالضرورة نفس stable)؛ يُعاد الحل عند كل Save |

إن لم يوجد إصدار stable للحزمة: MUST خطأ تحقق عند الحفظ مع توجيه لاختيار Specific أو Latest.

الواجهة MUST تعرض: اسم التطبيق، وضع الـ intent، والإصدار **المحلول** في Overview وفي معاينة QR.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create gateway without profile concepts (Priority: P1)

كمسؤول، أريد إنشاء بوابة تسجيل تحدد target node وbootstrap app وidentity mode وQR **دون** أي مفردات برفايل/سياسة في الواجهة.

**Why this priority**: يثبت نموذج البوابة المستقلة.

**Independent Test**: إنشاء مسار من القائمة → dialog → حفظ → Active QR؛ لا حقول برفايل في DOM/API client models.

**Acceptance Scenarios**:

1. **Given** dialog إنشاء، **When** يُحمَّل، **Then** لا حقول ولا نصوص policy/profile؛ فقط حقول الـ 4 concerns أعلاه.
2. **Given** نموذج مكتمل، **When** يحفظ، **Then** يُنشأ التعريف وتنتقل معاينة QR إلى **Active**.
3. **Given** لا برفايل في النظام، **When** ينشئ مساراً، **Then** الإنشاء ينجح (البرفايل ليس تبعية enrollment).

---

### User Story 2 - Dialog with live dual-column QR (Priority: P1)

كمسؤول، أريد رؤية **Pending QR** أثناء التعديل و**Active QR** بعد الحفظ في عمود ثابت.

**Acceptance Scenarios**:

1. **Given** Create/Edit، **When** يغيّر target node أو bootstrap app، **Then** يتحدث عمود Pending QR **فوراً على العميل** دون استدعاء preview API.
2. **Given** Create/Edit قبل أول حفظ، **When** يعرض Pending QR، **Then** لا يُعرض مفتاح عام قابل للمسح؛ تسمية «Pending — save to activate».
3. **Given** مسار محفوظ، **When** يفتح Overview، **Then** عمود يمين Active QR من الخادم + Header status Active.
4. **Given** تعديلات غير محفوظة على مسار Active، **When** يعرض Header، **Then** Active + شارة «Unsaved changes» (ليس Draft).
5. **Given** تعديلات غير محفوظة، **When** إغلاق، **Then** تأكيد فقدان التغييرات؛ Active QR السابق يبقى دون تغيير إن لم يحفظ.

---

### User Story 3 - Target node with kind + context (Priority: P1)

كمسؤول، أريد اختيار عقدة واحدة مع معرفة إنها «استقبال» أم «حاوية» ومعاينة المسار والحمل قبل التأكيد.

**Acceptance Scenarios**:

1. **Given** شجرة متعددة المستويات، **When** يفتح المنتقي، **Then** تظهر أنواع العقد (Locked / Inheritable) وbreadcrumb.
2. **Given** عقدة عليها عدد أجهزة كبير، **When** يحوم عليها، **Then** تحذير heavily loaded (عتبة قابلة للضبط).
3. **Given** عقدة Inheritable (container)، **When** يؤكد الاختيار ويحفظ، **Then** يُسمح بالحفظ مع تحذير مستمر في Overview حتى تغيير إلى Locked أو تأكيد صريح للإبقاء.
4. **Given** عقدة غير صالحة تقنياً (محذوفة/خارج النطاق)، **When** يحاول التأكيد، **Then** **منع** مع رسالة خطأ.

---

### User Story 4 - Bootstrap app intent (Priority: P2)

كمسؤول، أريد اختيار Stable أو إصدار محدد أو Latest بدل إدخال معرّف خام.

**Acceptance Scenarios**:

1. **Given** Stable intent وحزمة لها إصدار stable في الكتالوج، **When** يحفظ، **Then** QR contract يحتوي package + إصدار stable المعلَّم في الكتالوج (وليس مجرد «أحدث منشور»).
2. **Given** Specific version، **When** يغيّر الإصدار في الكتالوج، **Then** تحذير في Overview إن الإصدار لم يعد متاحاً.
3. **Given** Latest، **When** يحفظ لاحقاً مرة أخرى، **Then** قد يتغير الإصدار المحلول ويُحدَّث Active QR.

---

### User Story 5 - Safe delete with multi-dimensional impact (Priority: P2)

كمسؤول، أريد حذف مسار مع فهم أثر حقيقي قبل القرار.

**Acceptance Scenarios**:

1. **Given** حذف مسار، **When** يفتح تأكيد الأثر، **Then** يعرض **ثلاثة** أبعاد على الأقل:
   - أجهزة **قيد** التسجيل عبر QR (جلسات/محاولات حالية إن وُجدت)
   - أجهزة **مسجّلة تاريخياً** عبر هذا المسار
   - **مسحات QR نشطة** خلال آخر 7 أيام
2. **Given** الأبعاد الثلاثة كلها صفر، **When** يؤكد، **Then** حذف مباشر مع ملخص قصير (تأكيد واحد).
3. **Given** أي بُعد أثر > 0 (بما فيها أجهزة تاريخية فقط)، **When** يتابع الحذف، **Then** خطوة تأكيد ثانية داخل نفس الـ dialog: **كتابة اسم المسار** + عرض الأبعاد الثلاثة — **لا** منع الحذف بسبب `historicalEnrolledCount` وحده.
4. **Given** حذف ناجح مع أجهزة تاريخية، **When** يكتمل، **Then** يُزال المسار من القائمة ويُعطَّل QR؛ سجلات الأجهزة المسجّلة سابقاً **تبقى** دون حذف.

---

### User Story 6 - List hub + overview/edit in one shell (Priority: P2)

من القائمة: فتح مسار → Overview في dialog؛ Edit في نفس الإطار.

**Acceptance Scenarios**:

1. **Given** قائمة، **When** نقر صف، **Then** dialog Overview بالتخطيط الثنائي.
2. **Given** Overview، **When** Edit، **Then** يسار editable + يمين Pending حتى الحفظ.

---

### Edge Cases

- اسم مكرر — خطأ دون إغلاق dialog.
- target node محذوف — إلزام إعادة الاختيار قبل Save.
- bootstrap app أزيل من الكتالوج — تحذير + إعادة intent.
- انقطاع شبكة عند Save — إعادة محاولة مع الإبقاء على Pending QR (client preview لا يتأثر).
- محاولة نسخ/مسح Pending QR — معطّل أو رسالة «Save to activate».
- صلاحيات قراءة فقط — Overview + Active QR؛ لا Save/Delete.

---

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Module**: `internal/modules/enrollment_routes/` — **لا** استيراد أو منطق من `profiles` باستثناء migration/compat في طبقة منفصلة (انظر أدناه).
- **Conceptual split (required)**:
  - **`EnrollmentRouteDefinition`**: الاسم، الوصف، `targetNodeId`, `deviceIdentityMode`, bootstrap app intent + resolved ids، `type` — ما يحرّره المسؤول.
  - **`EnrollmentRouteRuntimeState`**: `qrcodeKey`, timestamps, QR generation metadata, إحصائيات مسح/تسجيل للأثر — ما يتغير بتشغيل النظام.
- **Java parity**: legacy QR enrollment في `backend/` — الحفاظ على مسارات `/rest/public/qr` و sync العامة؛ تغيير **محتوى** العقد فقط بما يتوافق مع Enrollment Contract.
- **REST admin**: `/rest/private/enrollment-routes` — request/response **للعميل** MUST NOT تتضمن `profileVersionId`.
- **Legacy `profile_version_id`**: يبقى في **migration/compat layer** فقط (DB + resolver خارج enrollment UI)؛ **ممنوع** في DTO الواجهة و OpenAPI للعميل الجديد.
- **Parity doc**: `serverBackendGo/docs/parity/enrollment-routes-ux.md`
- **Endpoints إضافية مقترحة**: `GET .../enrollment-routes/:id/impact` (أبعاد الحذف)، `GET .../options/tree-nodes` (أنواع + counts)، `GET .../options/bootstrap-apps`
- **Frontend**: `frontend/src/features/enrollment-routes/` — client types **EnrollmentRouteDefinitionView** بدون حقول profile.

### Functional Requirements

#### Domain boundary

- **FR-001**: Enrollment Route MUST يُعرَّف حصرياً بـ: target node، bootstrap app (intent)، device identity mode، QR contract metadata — **بدون** مفاهيم policy/profile.
- **FR-002**: وثائق الميزة وواجهة المستخدم ونماذج العميل MUST NOT تستخدم كلمة policy أو profile أو مرادفاتها العربية.
- **FR-003**: DTO طبقة UI/API consumer MUST NOT يتضمن `profileVersionId` أو `profileId` — حتى للقراءة؛ legacy يُقرأ فقط عبر أدوات migration/admin منفصلة إن لزم.

#### QR contract

- **FR-004**: حمولة QR MUST تطابق Enrollment Contract (routeId, targetNodeId, mainAppPackage, mainAppVersion, deviceIdentityMode, optional bootstrapFlags) — **بدون** حقول سياسة.
- **FR-005**: الواجهة MUST تعرض معاينة QR **دائماً** في عمود يمين: **Pending** (client-built من النموذج، غير قابل للمسح للإنتاج) قبل/أثناء التعديل غير المحفوظ؛ **Active** (server-generated، قابل للمسح ونسخ/تحميل) بعد الحفظ.
- **FR-005a**: Pending QR MUST NOT يستدعي endpoint preview على الخادم؛ Active QR MUST يأتي من الخادم بعد Save.

#### Dialog UX

- **FR-006**: Create/Edit/Overview MUST تستخدم التخطيط: Header (name + status badges) · Body 2-column · Footer (Save/Delete/Cancel).
- **FR-006a**: **Draft** badge فقط قبل أول Save؛ **Active** بعد أول Save؛ **Unsaved changes** شارة ثانوية عند Edit مع تغييرات غير محفوظة — MUST NOT استخدام Draft لمسار محفوظ يُعدَّل.
- **FR-006b**: أثناء Edit مع Unsaved، Active QR السابق MUST يبقى قابلاً للمسح حتى Save يحدّث التعريف والـ QR.
- **FR-007**: **100%** happy path للإنشاء/التعديل من dialog — لا route لصفحة editor كاملة في التنقل الافتراضي.

#### Target node

- **FR-008**: منتقي العقدة MUST يدعم تمييز Locked vs Inheritable واختيار **عقدة واحدة**؛ اختيار Inheritable **مسموح** مع تحذير غير مانع.
- **FR-009**: قبل التأكيد MUST يعرض: breadcrumb كامل، عدد أجهزة تحت العقدة، تحذير heavily loaded عند تجاوز العتبة، وتحذير container إن لم تكن العقدة Locked.
- **FR-009a**: بعد الحفظ على target Inheritable، Overview MUST يعرض تحذيراً مستمراً حتى تغيير إلى Locked أو تأكيد صريح «أبقِ هذا المجلد».

#### Bootstrap app

- **FR-010**: المسؤول MUST يختار intent: Stable (recommended) | Specific version | Latest available.
- **FR-011**: عند الحفظ/توليد Active QR، النظام MUST يحلّ intent إلى package + version فعليين في العقد: Stable → إصدار **stable/recommended** في الكتالوج؛ Specific → المختار؛ Latest → أحدث متاح منشور.
- **FR-011a**: Stable MUST NOT يُعادل Latest تلقائياً؛ إن غاب stable للحزمة، Save MUST يفشل مع رسالة واضحة.

#### Delete safety

- **FR-012**: تأكيد الحذف MUST يعرض: (1) أجهزة قيد التسجيل عبر QR، (2) أجهزة مسجّلة تاريخياً، (3) مسحات QR نشطة آخر 7 أيام.
- **FR-013**: عند **أي** بُعد أثر > 0، MUST خطوة تأكيد ثانية داخل نفس الـ dialog: **كتابة اسم المسار** (typed route name) مطابقة لاسم المسار — لا nested modal.
- **FR-013a**: الحذف MUST **لا يُمنع** بسبب `historicalEnrolledCount > 0` وحده؛ الأجهزة التاريخية لا تُحذف مع المسار — يُعطَّل QR والمسار فقط.
- **FR-013b**: عند الأبعاد الثلاثة = 0، تأكيد واحد كافٍ (بدون typed name).

#### Backend model

- **FR-014**: التخزين والخدمات MUST تفصل Definition عن RuntimeState (جداول أو aggregates منفصلة — قرار التنفيذ في plan).
- **FR-015**: تعديل Definition MUST NOT يكسر QR keys القائمة دون سياسة versioning صريحة (runtime يحتفظ بـ key حتى rotate متعمد).

#### General

- **FR-016**: التحقق: name، target node صالح، identity mode، bootstrap app قابل للحل — **لا** برفايل.
- **FR-017**: صلاحيات عرض/إضافة/تعديل/حذف حسب نموذج المسارات الحالي.
- **FR-018**: دعم AR/EN لكل النصوص الجديدة.

### Key Entities

- **EnrollmentRouteDefinition**: تعريف البوابة — editable بواسطة المسؤول؛ بدون profile.
- **EnrollmentRouteRuntimeState**: مفتاح QR، حالة Active/Draft، إحصائيات المسح والتسجيل للأثر.
- **EnrollmentContractPayload**: عقد QR المنطقي — الحقول في قسم Enrollment Contract.
- **TargetNodeSelection**: عقدة واحدة + kind (locked/inheritable) + path breadcrumb.
- **BootstrapAppIntent**: stable | specific | latest + application reference.
- **EnrollmentDeleteImpact**: enrollingNowCount, historicalEnrolledCount, activeQrScans7d — للعرض والتأكيد فقط؛ لا يمنع الحذف تلقائياً عند historical > 0.

### Out of scope (explicit)

- أي UI أو API يذكر policy/profile داخل enrollment.
- حل «ماذا يُطبَّق على الجهاز بعد التسجيل» — مملوك لـ sync/profiles (لا يُذكر في عقد enrollment).
- إدارة Profiles، Assignments، Rollout.
- wireflow تفصيلي state machine — يُنتج في `/speckit-plan` كـ contract منفصل إن لزم.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: مسؤول ينشئ بوابة كاملة (node + bootstrap + QR) في **< 3 دقائق** دون شاشات خارج enrollment hub.
- **SC-002**: **100%** create/edit من dialog ثنائي الأعمدة (happy path).
- **SC-003** (tree picker quality):
  - **Median time-to-selection** لاختيار target node: **< 10 ثوانٍ** في اختبار قبول (n ≥ 10 مستخدمين).
  - **Mis-selection rate** (اختيار عقدة ثم تغييرها قبل الحفظ): **< 5%** من الجلسات.
- **SC-004**: **0** حقول profile/policy في نماذج العميل أو استجابات enrollment admin API الموجهة للواجهة.
- **SC-005**: **0** أخطاء «معرّف تطبيق غير صالح» من إدخال خام في الاختبار المعياري (بعد intent picker).
- **SC-006**: **100%** عمليات حذف تعرض الأبعاد الثلاثة للأثر قبل التأكيد النهائي.
- **SC-007**: **100%** جلسات Create/Edit تعرض Pending QR قبل الحفظ وActive QR بعده.

---

## Assumptions

- عتبة «heavily loaded» افتراضية (مثلاً > 500 جهاز) قابلة للضبط per customer لاحقاً.
- «أجهزة قيد التسجيل» = محاولات enroll لم تكتمل أو devices في حالة provisioning عبر route — يُعرّف بدقة في plan مع الـ backend.
- «مسحات QR نشطة 7d» = أحداث domain أو telemetry موجودة أو تُضاف في migration 021+.
- Stable channel = علامة **stable/recommended** صريحة في كتالوج التطبيقات (Clarification 2026-05-24).
- Legacy DB column `profile_version_id` يبقى حتى migration لاحقة؛ **لا** يظهر للعميل.

---

## Dependencies

- **device_tree**: أنواع العقد (locked/inheritable)، counts تحت العقدة.
- **applications**: كتالوج + قنوات stable/latest.
- **domain_events** (أو ما يعادلها): لمسحات QR 7d إن لم تكن موجودة.
- **sync/profiles modules**: أي legacy resolution — **خارج** نطاق تنفيذ enrollment_routes في هذه الميزة.

---

## Planning artifacts (next phase)

يُوصى بإنشاء في `/speckit-plan`:

| Artifact | Content |
|----------|---------|
| `contracts/enrollment-contract-payload.md` | JSON schema للـ QR + أمثلة Pending/Active |
| `contracts/enrollment-route-dialog-ux.md` | State machine: List → Dialog(Create\|Overview\|Edit) → Delete confirm |
| `data-model.md` | Definition vs RuntimeState tables/aggregates |
