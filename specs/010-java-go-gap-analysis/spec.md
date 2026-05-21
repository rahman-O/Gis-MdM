# Feature Specification: Java–Go Migration Gap Analysis

**Feature Branch**: `010-java-go-gap-analysis`

**Created**: 2026-05-21

**Status**: Draft

**Input**: تحليل كامل ومفصل للمشروع القديم `backend/` (Java) ومقارنته بما تم نقله إلى `serverBackendGo/` (Go)، مع ملف MD في جذر المشروع يوضح الفروقات وما لم يُنقل لضمان هجرة ناجحة.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - رؤية شاملة لحالة الهجرة (Priority: P1)

كمسؤول تقني أو مالك المنتج، أريد مستنداً واحداً يوضح ما اكتمل نقله من الباك اند القديم وما تبقى، حتى أخطط للعمل المتبقي دون مفاجآت عند تشغيل الواجهة أو الوكلاء.

**Why this priority**: بدون خريطة فجوات واضحة لا يمكن إعلان «اكتمال الهجرة» بثقة.

**Independent Test**: يمكن التحقق بقراءة المستند ومطابقة 10 عينات عشوائية من endpoints Java مع حالة الجدول (منقول/جزئي/غير منقول).

**Acceptance Scenarios**:

1. **Given** المستند في جذر المشروع، **When** أبحث عن وحدة مثل `devices` أو `StatsResource`، **Then** أجد حالتها وحالة الفجوات المرتبطة.
2. **Given** قائمة مراحل MIGRATION.md، **When** أقارنها بالمستند، **Then** أفهم أن «done» لا يعني بالضرورة تطابقاً كاملاً.

---

### User Story 2 - ترتيب أولويات الإكمال (Priority: P2)

كفريق تطوير، أريد قائمة أولويات (P0–P3) للفجوات حسب التأثير على المستخدمين (إشعارات، جدولة، رفع ملفات).

**Why this priority**: يحوّل التحليل إلى خطة عمل قابلة للتنفيذ.

**Independent Test**: كل بند P0 في المستند مرتبط بسلوك يستخدمه React أو الوكلاء اليوم.

**Acceptance Scenarios**:

1. **Given** فجوة FCM/push، **When** أراجع §12 في المستند، **Then** تظهر ضمن P0 مع مرجع Java.

---

### User Story 3 - تتبع التحديثات المستقبلية (Priority: P3)

عند إغلاق فجوة في Go، أحدّث المستند وملف parity للوحدة حتى يبقى المرجع صالحاً.

**Why this priority**: يمنع تكرار العمل أو افتراضات قديمة.

**Acceptance Scenarios**:

1. **Given** اكتمال نقل `IconFileResource`، **When** أُحدّث parity و§4 في المستند، **Then** تتغير الحالة من ❌ إلى ✅.

---

### Edge Cases

- Endpoint مُهمل في Java (`PublicFilesResource`) يُوسَم ⊘ وليس ❌.
- مسارات React غير المستخدمة (`impersonate` users) تُوسَم out of scope مع ذكر السبب.
- وحدات Angular في plugins (`xtra`) خارج نطاق Go — مذكورة صراحة.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: يجب أن يوفّر المشروع مستنداً رئيسياً في الجذر: `JAVA-GO-MIGRATION-GAP-ANALYSIS.md`.
- **FR-002**: يجب أن يغطي المستند كل `*Resource.java` في `backend/` مع حالة: منقول / جزئي / غير منقول / خارج النطاق.
- **FR-003**: يجب أن يربط كل مجال بوحدة Go وملف parity عند الوجود.
- **FR-004**: يجب أن يفصّل الفجوات الجزئية (push، cron، audit، sync hooks، deviceinfo/devicelog export، إلخ).
- **FR-005**: يجب أن يتضمن قائمة endpoints غير المنقولة صراحة.
- **FR-006**: يجب أن يتضمن خطة أولويات P0–P3 لإغلاق الهجرة.
- **FR-007**: يجب أن يوضح تأثير الفجوات على `frontend/` والوكلاء حيث ينطبق.
- **FR-008**: يجب أن يُحدَّث المستند عند إغلاق فجوات (إجراء صيانة، ليس أتمتة إلزامية في هذه المرحلة).

### Key Entities

- **Migration module**: مجال وظيفي (auth, devices, plugin push, …) بربط Java ↔ Go.
- **Gap**: فجوة سلوكية أو endpoint أو مكوّن خلفية غير منقول.
- **Parity doc**: `serverBackendGo/docs/parity/<name>.md` — مصدر تفصيلي لكل وحدة.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: يمكن لمراجع تقني التحقق من حالة 100% من موارد REST Java المدرجة في §4 خلال أقل من 30 دقيقة باستخدام المستند فقط.
- **SC-002**: لا يبقى أكثر من 3 عناصر حرجة (P0) غير موثّقة في المستند بعد مراجعة الفريق.
- **SC-003**: عند إغلاق فجوة P0، يُحدَّث المستند خلال نفس دورة التطوير (نفس PR أو commit doc).
- **SC-004**: الفريق يستطيع ترتيب sprint واحد على الأقل من §12 دون الرجوع لقراءة كامل `backend/` يدوياً.

## Assumptions

- Java `backend/` يبقى مرجع الحقيقة حتى اكتمال الهجرة.
- `MIGRATION.md` Phases 1–8 «done» تعني تغطية مسارات أساسية وليس parity كامل.
- React الحالي هو العميل الأساسي؛ مسارات superadmin/impersonate قد تبقى خارج النطاق ما لم يُطلب خلاف ذلك.

## Deliverables

| Deliverable | Path |
|-------------|------|
| Gap analysis (primary) | [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../JAVA-GO-MIGRATION-GAP-ANALYSIS.md) |
| Feature spec (this file) | `specs/010-java-go-gap-analysis/spec.md` |
