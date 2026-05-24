# Feature Specification: Profile Hub & Enrollment UX

**Feature Branch**: `019-profile-hub-ux`

**Created**: 2026-05-23

**Status**: Draft (UX architecture refined 2026-05-23)

**Input**: تحسين تجربة إدارة Profiles بأفضل الممارسات: إلغاء إسناد Profile من مسارات التسجيل، فصل العرض عن التعديل، **Profile Workspace** (ليس modal تقليدي)، onboarding إسناد الشجرة، ومستوى enterprise خفيف (Intune / Workspace ONE–like) دون ازدحام بصري.

**Related**:

- [`specs/018-profile-rollout-ops/spec.md`](../018-profile-rollout-ops/spec.md) — إسناد الشجرة، حالة التطبيق، تعطيل/تفعيل
- [`specs/017-device-control-plane/spec.md`](../017-device-control-plane/spec.md) — Profiles، مسارات التسجيل، الشجرة
- [`DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md`](../../DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md) — Control Plane

---

## Clarifications

### Session 2026-05-23 (UX review)

- Q: هل «Dialog» يعني popup وسط الشاشة؟ → A: **لا** — **Profile Workspace**: شبه ملء الشاشة، backdrop خفيف، header ثابت، تنقل جانبي؛ ليس popup تقليدي ولا scroll واحد طويل لكل الميزات.
- Q: كيف نتجنب «صفحة كاملة متنكرة داخل dialog»؟ → A: **طبقات إدراكية ثلاث** (Cockpit header → Overview cards → Deep sections) + **sidebar navigation** بدل tabs أفقية كثيرة.
- Q: أين يكون Publish؟ → A: **إجراء عالمي في الـ header** (معطّل عند لا تغييرات أو فشل التحقق)؛ ليس مخفياً داخل Editor فقط.
- Q: هل nested dialogs مسموحة؟ → A: **ممنوعة** — استخدام side panels، inline drawers، أو step transitions فقط.

---

## UX Architecture *(mandatory for this feature)*

### Design goal

تحويل Profiles من «صفحات إعدادات متفرقة» إلى **Control Plane تشغيلي** بمستوى منتجات enterprise (Microsoft Intune، VMware Workspace ONE، Google AMC) **بطابع أخف وأسرع** — دون أن يصبح الواجهة cockpit نووي أو childish.

### Anti-patterns (MUST NOT)

| Anti-pattern | Why |
|--------------|-----|
| Popup مركزي صغير + scroll طويل واحد | يحس كصفحة متنكرة ومزدحمة |
| Tabs أفقية أكثر من ~4 عناصر | «قطار شحن» — صعب المسح |
| Overview + Versions + Rollout + Assignment + Editor + Publish في نفس الطبقة البصرية | وحش بصري |
| Nested dialogs (publish داخل assignment داخل workspace) | متاهة أبواب |
| Publish مخفي داخل Editor فقط | يفوّت إجراء production-critical |
| Form fields في وضع العرض | يخلط القراءة بالتحرير |

### Target pattern: Profile Workspace

```text
Profiles List (Control Radar)
    ↓ click
Profile Workspace
    ├── Layer 1: Fixed Cockpit Header (always visible)
    ├── Layer 2: Overview (cards, read-only, 5-second rule)
    ├── Sidebar: Overview | Assignments | Rollout | Versions | Editor | Activity
    └── Layer 3: Deep section content (one at a time; no nested modals)
```

**NOT**: `Page → Modal → Modal → Tab → Dialog`

### Layer 1 — Cockpit header (fixed)

يبقى ظاهراً دائماً أثناء التنقل داخل الـ Workspace:

- اسم Profile
- **Health** badge (Healthy / Warning / Error / Draft Only)
- **Lifecycle** badge: Draft (amber) | Published (green) | Disabled (gray/red)
- Published version number
- Assigned / Not assigned indicator
- Actions: **Edit** | **Publish** (global; disabled when no draft changes or validation fails) | **Close**

### Layer 2 — Overview (default section)

بطاقات (cards) فقط — **لا حقول إدخال**:

| Card | Content (example) |
|------|-------------------|
| Status | Active / Disabled |
| Assignment | 3 folders · Not assigned |
| Rollout | 124 installed · 3 failed |
| Apps | 7 linked apps |
| Kiosk | Enabled / Off |
| Last publish | 2h ago |

**5-second rule (SC-006)**: خلال 5 ثوانٍ من فتح Overview يفهم المسؤول: الحالة، منشور أم لا، مُسنَد أم لا، حجم الأثر (أجهزة)، وجود مشاكل — دون فتح Editor.

### Layer 3 — Deep sections (on demand)

تُفتح من **sidebar** واحدة في كل مرة:

| Section | Purpose |
|---------|---------|
| Overview | Layer 2 cards (home) |
| Assignments | Tree assignment + **miniature tree preview** |
| Rollout | Device status grid (018) |
| Versions | Version list, fork draft |
| Editor | Full policy tabs — **edit mode only** |
| Activity | Timeline (publish, assign, disable, rollout events) |

### Read vs Edit — visual system

| Mode | Visual treatment |
|------|------------------|
| **Read (Overview + deep read-only views)** | خلفية هادئة، cards، labels، حدود قليلة، **no inputs** |
| **Edit (Editor section)** | accent لون مختلف، **top warning bar** («تغيير سياسة إنتاج»)، **sticky save bar**، inputs واضحة |

الدماغ يميّز فوراً: «أقرأ» vs «أغيّر سياسة production».

### Navigation

- **Sidebar** داخل Workspace (ليس tabs أفقية كثيرة): `Overview · Assignments · Rollout · Versions · Editor · Activity`
- على **mobile**: full-screen sheet، sidebar → drawer، ملخص قابل للطي، إجراءات سفلية ثابتة

### Density

- صفوف مدمجة (compact rows)، تباعد متوسط، cards خفيفة، typography واضحة — بين sparse وdense المفرط

### Nested UI rule

تأكيد النشر (Impact Preview)، إسناد مجلد، أو تفاصيل rollout: **side panel أو inline drawer** داخل Workspace — **never** dialog فوق dialog.

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — إزالة ربط Profile من مسار التسجيل (Priority: P1)

كمسؤول MDM، عند إنشاء أو تعديل **مسار تسجيل** أُعرّف فقط ما يخص التسجيل (الاسم، الوصف، مجلد الشجرة الافتراضي، معرّف الجهاز، تطبيق MDM للـ QR) **دون** اختيار Profile. السياسة من **Tree assignment** فقط.

**Why this priority**: مصدر حقيقة واحد للسياسة — يقلل chaos مستقبلياً.

**Independent Test**: مسار جديد بلا Profile → تسجيل جهاز → سياسة من إسناد المجلد فقط.

**Acceptance Scenarios**:

1. **Given** محرر مسار تسجيل، **When** يُفتح، **Then** لا يظهر محدّد Profile.
2. **Given** مسار قديم كان مربوطاً بProfile، **When** يُحفظ، **Then** يبقى صالحاً للتسجيل + ملاحظة: السياسة من Profiles → Assignments.
3. **Given** مجلد مُسنَد + مسار بنفس المجلد الافتراضي، **When** يزامن الجهاز، **Then** يستلم سياسة الإسناد.

---

### User Story 2 — Profile Workspace بدل صفحة أو modal تقليدي (Priority: P1)

كمسؤول، من قائمة Profiles أنقر فيُفتح **Profile Workspace** (شبه ملء الشاشة، backdrop خفيف) مع Cockpit header وOverview افتراضي — دون مغادرة سياق القائمة.

**Why this priority**: سرعة workflow + إحساس Control Plane (Linear / Notion / Vercel–like).

**Independent Test**: فتح Workspace → sidebar إلى Rollout → Close → العودة لنفس scroll القائمة.

**Acceptance Scenarios**:

1. **Given** قائمة Profiles، **When** ينقر صفاً، **Then** يُفتح Workspace بعرض ≥90% على desktop مع padding جانبي وليس popup 400px.
2. **Given** Workspace مفتوح، **When** يتنقل بين الأقسام، **Then** الـ header يبقى ثابتاً والمحتوى يتبدل في المنطقة الرئيسية فقط.
3. **Given** رابط `/profiles/:id/edit` قديم، **When** يُفتح، **Then** يُعاد توجيهه لفتح Workspace على القسم المناسب (Overview أو Editor).

---

### User Story 3 — Overview: قراءة في 5 ثوانٍ + فصل التحرير (Priority: P1)

**الوضع الافتراضي**: قسم Overview (cards) + زر **Edit** في الـ header ينقل لقسم Editor — لا تحرير مباشر في Overview.

**Acceptance Scenarios**:

1. **Given** Workspace جديد، **When** يُفتح، **Then** القسم الافتراضي Overview ببطاقات الحالة دون inputs.
2. **Given** Overview، **When** يُقاس زمن الفهم، **Then** 90% من المختبرين يجيبون صحيحاً على: منشور؟ مُسنَد؟ معطّل؟ مشاكل rollout؟ خلال **5 ثوانٍ** (اختبار قبول).
3. **Given** Edit من الـ header، **When** يُفعَّل، **Then** warning bar + accent + sticky save وتظهر تبويبات السياسة الكاملة في قسم Editor فقط.
4. **Given** تغييرات غير محفوظة في Editor، **When** يغادر القسم أو يغلق Workspace، **Then** حفظ / تجاهل / إلغاء.

---

### User Story 4 — Health state & قائمة Control Radar (Priority: P1)

كمسؤول، أرى **صحة Profile** في القائمة وفي الـ header دون فتح التفاصيل؛ القائمة تعمل كـ «رادار تشغيل».

**Health states**:

| State | Meaning |
|-------|---------|
| **Healthy** | منشور + مُسنَد + rollout stable (لا failures حرجة) |
| **Warning** | لا إسناد، أو مسودة غير منشورة، أو Stale (لم يُنشر منذ مدة — threshold في الخطة) |
| **Error** | rollout failures أو تعطيل مع أجهزة معلقة |
| **Draft Only** | لا إصدار منشور |

**List badges (visible without opening)**:

| Badge | Meaning |
|-------|---------|
| No Assignment | يحتاج إسناد |
| Disabled | policy off |
| Draft Changes | unpublished draft |
| Rollout Issues | failures |
| Stale | لم يُنشر منذ مدة |

**Acceptance Scenarios**:

1. **Given** Profile بلا إسناد، **When** يُعرض في القائمة، **Then** شارة No Assignment + Health Warning.
2. **Given** 3 أجهزة failed rollout، **When** القائمة، **Then** Rollout Issues + Health Error (أو Warning حسب السياسة).

---

### User Story 5 — Assignments بصري: miniature tree (Priority: P1)

كمسؤول، في قسم Assignments أرى **معاينة شجرة مصغّرة** مع تمييز المجلد المختار (subtle highlight) قبل التأكيد — لتقليل أخطاء الإسناد.

**Acceptance Scenarios**:

1. **Given** شجرة أجهزة، **When** أختار مجلداً للإسناد، **Then** يظهر مسار المجلد في معاينة شجرية (مثال: HQ → Baghdad → Sales) والمجلد المحدد مميز بصرياً.
2. **Given** تأكيد إسناد، **When** يُنفَّذ، **Then** يُحدَّث Overview card وHealth دون nested dialog (panel جانبي أو inline).

---

### User Story 6 — إسناد الشجرة عند أول إنشاء (Priority: P1)

بعد إنشاء Profile، Workspace يفتح على **معالج موجّه** (3 خطوات موصى بها): اسم/هوية → نشر أول إصدار → إسناد شجرة (مع tree preview).

**Acceptance Scenarios**:

1. **Given** Profile جديد، **When** يُنشأ، **Then** Workspace يفتح على خطوة الإسناد أو يوجّه للنشر أولاً ثم الإسناد — **بدون nested wizard dialog**.
2. **Given** تخطي إسناد، **When** يُغلق، **Then** Warning مستمر + CTA «Assign now» في Overview.

---

### User Story 7 — Publish عالمي + Impact Preview (Priority: P1)

كمسؤول، أضغط **Publish** من الـ header؛ قبل التأكيد أرى **Impact Preview**: عدد الأجهزة، المجلدات المتأثرة، وملاحظة عن مسارات التسجيل المرتبطة بالمجلدات (إن وُجدت) — دون ربط سياسة بالمسار نفسه.

**Acceptance Scenarios**:

1. **Given** مسودة جاهزة، **When** Publish من الـ header، **Then** side panel يعرض: «يؤثر على N جهازاً، M مجلداً» + تأكيد ≥50 جهاز (018).
2. **Given** لا تغييرات في المسودة، **When** Publish، **Then** الزر معطّل مع تلميح «لا تغييرات للنشر».

---

### User Story 8 — Activity timeline (Priority: P2)

كمسؤول، في قسم Activity أرى خط زمني: نشر، إسناد، تعطيل، فشل rollout — حتى لو أحداث أساسية في v1.

**Acceptance Scenarios**:

1. **Given** نشر v12، **When** أفتح Activity، **Then** يظهر حدث «[User] published v12» مع وقت نسبي.
2. **Given** لا أحداث بعد، **Then** حالة فارغة توضيحية — ليس جدولاً فارغاً صامتاً.

---

### User Story 9 — Mobile Workspace (Priority: P2)

على شاشة ضيقة، Workspace يصبح full-screen sheet مع drawer للتنقل وملخص قابل للطي.

**Acceptance Scenarios**:

1. **Given** عرض جوال، **When** يُفتح Profile، **Then** full-screen + أزرار إجراء سفلية (Edit, Publish, Close).

---

### Edge Cases

- Profile بلا مسودة/منشور → Health Draft Only + CTA إنشاء مسودة.
- محاولة Publish بفشل validation → تعطيل الزر + رسائل عند الحقول في Editor.
- مسار تسجيل legacy بـ profile_version_id → تسجيل يعمل؛ UI لا تعرض الربط؛ Impact Preview قد يذكر «routes in folder» إعلامياً فقط.
- صلاحيات قراءة فقط → لا Edit/Publish؛ Overview + Rollout read-only.
- إغلاق Workspace أثناء Impact Preview → إلغاء النشر وليس nested close hell.

---

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Modules**: `enrollment_routes`, `profiles`, `sync` (كما في النسخة السابقة)؛ إضافة أحداث Activity إن لزم (`domain_events` موجود 017).
- **Parity**: `parity/profile-hub-ux.md` يصف Workspace sections + إزالة Profile من route UI.
- **Frontend**: مكوّن Workspace جديد؛ refactor `ProfilesPage`, `ProfileEditorPage` → Workspace؛ `EnrollmentRouteEditorPage` بدون Profile picker.

### Functional Requirements

**Enrollment & policy source**

- **FR-001**: مسار التسجيل MUST NOT يعرض اختيار Profile.
- **FR-002**: توجيه صريح: السياسة من Profiles → Assignments.

**Workspace shell**

- **FR-003**: فتح Profile من القائمة MUST يفتح **Profile Workspace** (ليس صفحة كاملة افتراضية ولا popup صغير).
- **FR-004**: Workspace MUST يطبّق **Layer 1 header ثابت** + **Layer 2 Overview cards** + **sidebar navigation** للأقسام العميقة.
- **FR-005**: Workspace MUST NOT يستخدم nested dialogs؛ التأكيدات عبر side panel / drawer / inline steps.
- **FR-006**: التنقل الجانبي MUST يتضمن على الأقل: Overview, Assignments, Rollout, Versions, Editor, Activity.

**Read / Edit**

- **FR-007**: Overview MUST read-only (cards فقط؛ لا form fields).
- **FR-008**: التحرير MUST في قسم Editor فقط مع تمييز بصري قوي (warning bar, accent, sticky save).
- **FR-009**: زر Edit في الـ header ينقل لقسم Editor؛ Overview لا يسمح بتعديل مباشر.

**Cockpit & lifecycle**

- **FR-010**: Header MUST يعرض: اسم، Health، lifecycle (Draft/Published/Disabled)، إصدار منشور، assigned/not assigned.
- **FR-011**: **Publish** MUST في الـ header (معطّل عند لا تغييرات أو validation fail).
- **FR-012**: قبل Publish MUST **Impact Preview** (أجهزة، مجلدات؛ تأكيد عتبة 018 عند اللزوم).

**Health & list**

- **FR-013**: كل Profile MUST له Health محسوب: Healthy | Warning | Error | Draft Only.
- **FR-014**: قائمة Profiles MUST تعرض شارات: No Assignment, Disabled, Draft Changes, Rollout Issues, Stale (حسب البيانات المتاحة).

**Assignments**

- **FR-015**: قسم Assignments MUST يعرض **miniature tree preview** مع تمييز المجلد المختار.
- **FR-016**: إنشاء Profile جديد MUST يوجّه لخطوة إسناد (معالج داخل Workspace بدون wizard modal منفصل).

**Activity & mobile**

- **FR-017**: قسم Activity MUST يعرض timeline أحداث أساسية (publish, assign, enable/disable, rollout failures) — v1 يمكن أن يبدأ بأحداث من `domain_events` إن وُجدت.
- **FR-018**: على mobile MUST full-screen sheet + drawer navigation + bottom actions.

**Compatibility**

- **FR-019**: روابط محرر قديمة MUST تفتح Workspace على القسم المناسب.
- **FR-020**: مسارات API قد تبقي `profileVersionId` deprecated؛ UI تعامل الشجرة كمصدر وحيد.

### Key Entities

- **Profile Workspace**: غلاف UX — header + sidebar + content area.
- **Profile Health**: حالة تشغيل مجمّعة للقائمة والـ header.
- **Overview cards**: لقطات قراءة فقط.
- **Impact Preview**: تأكيد قبل النشر.
- **Activity event**: حدث خط زمني للمسؤول.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: إنشاء Profile + نشر + إسناد شجرة في **< 3 دقائق** داخل Workspace دون صفحات منفصلة.
- **SC-002**: **100%** مسارات تسجيل جديدة بلا حقل Profile في UI.
- **SC-003**: فتح Profile والعودة للقائمة في **نقرتين** (فتح Workspace، إغلاق).
- **SC-004**: **90%** مهام «مراجعة الحالة» تُنجز من Overview دون Editor.
- **SC-005**: تقليل تذاكر «أين أربط السياسة؟» — مسار واحد: Assignments.
- **SC-006 (5-second rule)**: **90%** من المختبرين يحددون خلال 5 ثوانٍ: منشور؟ مُسنَد؟ معطّل؟ يوجد مشكلة rollout؟ — من Overview فقط.
- **SC-007**: **صفر** nested dialogs في مسار قبول Workspace (audit UX checklist).

---

## Assumptions

- 018 يوفّر بيانات rollout وassignment وenable/disable.
- Activity v1 قد تقتصر على أحداث الخادم المسجّلة؛ توسيع لاحقاً.
- **Stale** threshold (مثلاً 30 يوماً بدون نشر) يُحدد في plan.
- ألوان lifecycle: Draft amber, Published green, Disabled gray/red — في plan/design tokens.

---

## Out of Scope (v1)

- تغيير agent Android.
- إعادة تصميم صفحة الأجهزة بالكامل.
- حذف `profile_version_id` من DB.
- Analytics متقدمة على Activity (فلترة، تصدير).

---

## UX Reference — v1 Recommended Checklist

_عناصر كانت «اقتراحات» وأصبحت جزءاً من اتجاه v1 بعد مراجعة المنتج:_

| # | Item | v1 |
|---|------|-----|
| 1 | Profile Workspace (not classic modal) | **Required** |
| 2 | 3-layer layout + sidebar nav | **Required** |
| 3 | Overview cards + 5-second rule | **Required** |
| 4 | Health + list badges (Control Radar) | **Required** |
| 5 | Read/Edit visual system | **Required** |
| 6 | Tree miniature in Assignments | **Required** |
| 7 | Publish in header + Impact Preview | **Required** |
| 8 | No nested dialogs | **Required** |
| 9 | Activity timeline (basic) | **Required** |
| 10 | Mobile full-screen sheet | **Required** |
| 11 | Create wizard inside Workspace | **Required** |
| 12 | Keyboard shortcuts (E, Esc, Ctrl+S) | Optional P2 |

---

## Dependencies

| Dependency | Reason |
|------------|--------|
| 017 device control plane | Tree, profiles, routes |
| 018 profile rollout ops | Assignment, rollout, publish impact data |

---

## Product Evaluation (internal)

| Dimension | Target |
|-----------|--------|
| Product architecture | Control Plane واضح — policy من الشجرة فقط |
| Scalability | Sidebar يستوعب أقساماً جديدة دون tabs أفقية |
| Enterprise readiness | Health, Impact Preview, Activity |
| Primary risk | UI complexity — **mitigated by layered workspace** |
| Critical success factor | Workspace layers + 5-second Overview |

---
