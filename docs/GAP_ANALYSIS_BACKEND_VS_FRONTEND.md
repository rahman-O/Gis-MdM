# تحليل الفجوات: الـ backend (Headwind MDM) مقابل الـ frontend الحديث (React)

**تاريخ المسودة:** 2026-05-10  
**نطاق الوثيقة:** مقارنة بين واجهة **AngularJS المدمجة** تحت `backend/server/src/main/webapp/` والمشروع الجديد **`frontend/`** (React + Vite + TypeScript)، مع إشارة لمواصفات `.kiro/specs/` حيث وُجدت.

---

## 1. ملخص تنفيذي

- الـ **backend** يظل المصدر الرسمي للـ REST API (`/rest/...`): مصادقة، أجهزة، تطبيقات، تكوينات، مجموعات، مستخدمين، أدوار، ملخص، إعدادات، ملفات، أيقونات، عملاء (للمسؤول)، إشعارات دفع، تحديثات، تلميحات، تسجيل، استعادة كلمة المرور، إلخ.
- الـ **frontend** الحديث يغطي **جزءاً كبيراً** من مسارات التشغيل اليومي (لوحة، أجهزة، تطبيقات، تكوينات، مجموعات، مستخدمين، أدوار، إعدادات مختصرة، تسجيل الدخول، QR للتسجيل) لكنه **لا يعادل** بعد التجربة الكاملة لـ SPA القديمة: صفحات وجداول مسارات مفقودة أو مدمجة بشكل مختلف، **غياب دعم الإضافات (plugins)** في الواجهة، **غياب التعدد الكامل للغات في الواجهة** (I18n)، وعدة **واجهات API لم تُربط** بعد.

---

## 2. طبقة الـ backend (مختصر معماري)

### 2.1 الوحدة الرئيسية `backend/server`

- موارد JAX-RS تحت حزم مثل `com.hmdm.rest.resource` مع بادئة REST اعتيادية **`/rest`** (الـ frontend يستخدم `apiClient` بـ `baseURL` يشير إلى `/api` أو `/rest` حسب الإعداد).
- أمثلة موارد أولية:

| المورد | المسار الأساسي | ملاحظات |
|--------|------------------|----------|
| `AuthResource` | `/public/auth` | تسجيل الدخول/الخروج، خيارات المصادقة |
| `UserResource` | `/private/users` | مستخدمين، حالي، أدوار، استعادة، تصيير super-admin |
| `UserRoleResource` | `/private/roles` | صلاحيات، أدوار، CRUD |
| `DeviceResource` | `/private/devices` | بحث، حذف جماعي، مجموعات، إعدادات تطبيقات الجهاز |
| `ConfigurationResource` | `/private/configurations` | بحث، نسخ، ترقية تطبيقات |
| `ApplicationResource` | `/private/applications` | أندرويد/ويب، إصدارات، ربط بالتكوينات، admin مشترك |
| `GroupResource` | `/private/groups` | مجموعات |
| `SummaryResource` | **`/private/summary`** مع **`GET /devices`** | إحصاءات غنية (مخططات، تكوينات، حالات) |
| `SettingsResource` | `/private/settings` | تصميم افتراضي، لغة، متفرقات، أدوار مستخدم مشتركة |
| `FilesResource` | `/private/web-ui-files` | رفع/حدود/ملفات تطبيقات الويب |
| `IconResource` | `/private/icons` | أيقونات |
| `CustomerResource` | `/private/customers` | عملاء (لوحة تحكم متعددة المستأجرين) |
| `HintResource` | `/private/hints` | تفعيل/تعطيل التلميحات، سجل |
| `UpdateResource` | `/private/update/check` | فحص التحديثات |
| `PasswordResetResource` | `/public/passwordReset` | استعادة/إعادة تعيين |
| `PushApiResource` | `/private/push` | إشعارات دفع |
| `QRCodeResource` | `/public/qr` | QR للتسجيل (PNG/JSON) |

### 2.2 الإضافات (Plugins) تحت `backend/plugins/`

- `PluginResource` (`/plugin/main/...`): إضافات متاحة/مفعّلة/قوالب واجهات.
- إضافات مثل `devicelog` و`deviceinfo` تضيف موارد REST خاصة بها.

هذه الطبقة **موجودة وتعمل مع الـ Angular القديم**؛ الـ React **لا يستدعي** مسارات الإضافات أو يعرض قوالب الديناميكية.

---

## 3. طبقة الـ frontend الحديث (`frontend/src`)

### 3.1 المسارات المعرفة في `App.tsx`

| المسار | الصفحة |
|--------|--------|
| `/login` | تسجيل الدخول |
| `/dashboard` | لوحة (إحصاءات مبسطة) |
| `/devices` | أجهزة (بحث، ترقيم، لوحة تفاصيل، جماعي، QR) |
| `/groups` | مجموعات |
| `/applications` | تطبيقات |
| `/applications/admin` | تطبيقات مشتركة (مسؤول) |
| `/application/:id/versions` | إصدارات تطبيق |
| `/configurations` | قائمة التكوينات |
| `/configurations/:id/edit` | محرر التكوين (تبويبات) |
| `/users` | مستخدمون |
| `/roles` | أدوار وصلاحيات |
| `/settings` | إعدادات (نموذج موحّد) |
| `/qr/...` | تجربة QR للتسجيل |

### 3.2 استدعاءات API المستخدمة فعلياً (من مسح `src/`)

يقترب الاستخدام من:  
`/public/auth/*`, `/private/users/*`, `/private/roles/*`, `/private/devices/*`, `/private/groups/*`, `/private/configurations/*`, `/private/applications/*`, `/private/web-ui-files/*`, `/private/summary/devices`, `/private/settings`, `/public/qr/*`.

---

## 4. مقارنة المسارات: Angular (legacy) مقابل React

يستند التبويب القديم إلى `TabController` و`app.js`؛ التبويب الجديد يعتمد على React Router.

| ميزة / منطقة في SPA القديم | مسار Angular تقريبي | الحالة في React |
|-----------------------------|---------------------|------------------|
| ملخص/إحصاءات (charts) | `/summary` | **جزئي:** `/dashboard` يعرض 3 أرقام فقط من استجابة أوسع |
| الأجهزة | `/` (main) | **موجود** `/devices` بتجربة أغنى في بعض الجوانب |
| التطبيقات | `/applications` | **موجود** |
| إصدارات التطبيق | `/application/{id}/versions` | **موجود** |
| التكوينات | `/configurations` | **موجود**؛ المسار القديم للمحرر `/configuration/{id}` أصبح `/configurations/:id/edit` |
| الملفات (Web UI files) | `/files` | **غير موجود** كصفحة؛ الاستخدام **ضمن** تطبيقات الويب عبر `webUiFilesService` فقط |
| التصميم الافتراضي (Design) | `/designSettings` | **غير موجود** كصفحة؛ الـ API `POST /private/settings/design` متاح في الـ backend |
| إعدادات مشتركة للأدوار | `/commonSettings` | **غير موجود**؛ الـ API `POST /private/settings/userRoles/common` موجود |
| اللغة والإعدادات العامة | `/langSettings` (ومسارات مجزأة) | **مدمج جزئياً** في `/settings`؛ لا فصول كاملة مثل القديم |
| المستخدمون | `/users` | **موجود** |
| الأدوار | `/roles` | **موجود** |
| المجموعات | `/groups` | **موجود** |
| الأيقونات | `/icons` | **غير موجود** كصفحة (`/private/icons`) |
| التلميحات | `/hints` | **غير موجود** (`/private/hints`) |
| الإضافات | `/pluginSettings` | **غير موجود** — لا تكامل `plugin/main` |
| الملف الشخصي | `/profile` | **غير موجود** |
| التحديثات | `/updates` | **غير موجود** (`/private/update/check`) |
| لوحة التحكم (مسؤولو العملاء) | `/control-panel` | **غير موجود** |
| التسجيل العام | `/signup`, `/signupComplete` | **غير موجود** في React Router (قد يبقى على الـ WAR القديم) |
| استعادة كلمة المرور | `/passwordRecovery`, `/passwordReset/{token}` | **غير موجود** في React |
| المصادقة الثنائية | `/twofactor` | **غير موجود** في React |

---

## 5. فجوات API: endpoints في الـ backend دون تكامل واضح في React

| المنطقة | Endpoint تقريبي | الملاحظة |
|---------|-----------------|----------|
| العملاء (multi-tenant) | `/private/customers/*` | لا يوجد `grep` لاستخدام في `frontend/src` — إدارة العملاء/التنكر غير مبنية في SPA الجديد |
| الأيقونات | `/private/icons/*` | غير مستخدم في المسارات الحالية |
| التلميحات | `/private/hints/*` | غير مستخدم |
| فحص التحديث | `/private/update/check` | غير مستخدم |
| الإشعارات الدفعية | `/private/push` | غير مستخدم |
| إعدادات التصميم | `POST /private/settings/design` | غير مستدعى من صفحة مخصصة |
| إعدادات أدوار المستخدم الشائعة | `POST /private/settings/userRoles/common` | غير مستدعى |
| مستخدم حسب التفاصيل / مسارات super-admin إضافية | أجزاء من `/private/users/*` | قد تكون غير مستخدمة حسب الدور |
| إضافات | `/plugin/main/private/...` | غير مستخدم |
| Plugin APIs | `/plugins/devicelog/...`, `/plugins/deviceinfo/...` | غير مستخدم |

*(القائمة تقريبية للتخطيط؛ التحقق الدوري من `grep` على المشروع يحدّثها.)*

---

## 6. فجوات وظيفية ومنتجية تفصيلية

### 6.1 لوحة التحكم (Dashboard)

- **الـ backend** يعيد عبر `GET /private/summary/devices` كائناً غنياً: `statusSummary`, `installSummary`, قوائم حسب التكوين، تسجيل شهري، إلخ (`SummaryResource`).
- **الـ React** (`DashboardPage` + `summaryService`) يستغل **جزءاً صغيراً** (عناوين مثل إجمالي الأجهزة والمسجّلين).
- مواصفات **`.kiro/specs/dashboard-improvements/tasks.md`** تصف لوحة متقدمة (5 بطاقات، جدول أجهزة حديث، تحديث كل 60 ثانية) — **المهام ما زالت غير منفذة** (`[ ]`).

### 6.2 الإعدادات (Settings)

- الصفحة الحديثة تركز على حقول مختارة؛ الـ **legacy** يفصل التصميم الافتراضي، الإعدادات المشتركة للأدوار، اللغة، مع حقول مثل **idle logout** وغيرها.
- بعض حقول الـ UI في React **لا تُحفظ** في جدول `settings` (مثل أسماء قد تخص جدول العملاء) — تتطلب either **ربط API العملاء** أو إزالة/توضيح للمستخدم.

### 6.3 التعددية اللغوية (i18n)

- الـ Angular يحمّل حزماً من `localization/*.js` ولغات كثيرة.
- الـ React الحالي **إنجليزي أساساً** في الواجهة (نصوص ثابتة في المكوّنات).

### 6.4 الملف الشخصي وكلمة المرور والجلسة

- الـ legacy: `/profile`, استعادة كلمة المرور، 2FA.
- الـ React: **لا صفحات مطابقة** في المسارات الحالية (قد تُستكمل لاحقاً عبر نفس موارد `UserResource` و`PasswordResetResource`).

### 6.5 الإضافات (Plugins)

- الـ backend يدعم إضافات ديناميكية (`settingsViewTemplate`, `functionsViewTemplate`).
- الـ React **لا يحمّل** هذه القوالب — أي منطق إضافي معرّف في الإضافات **غير مرئي** في الواجهة الجديدة.

### 6.6 الخرائط والـ Leaflet

- الـ SPA القديم يضمّ مكتبات Leaflet (مثلاً في `app.js` `SUPPORTED_LIBS`).
- لا يوجد في المسح السريع مسار خرائط واضح في React الحالي؛ أي خريطة أجهزة في القديم قد تكون **غير منقولة**.

---

## 7. مواصفات Kiro — ما يشير إلى عمل غير مكتمل

| الموضوع | الملف | حالة تقريبية |
|---------|-------|----------------|
| تحسينات لوحة التحكم | `dashboard-improvements/tasks.md` | المسار الصحيح للملخص `GET /private/summary/devices`؛ بقية البنود الاختيارية قد تبقى `[ ]` |
| بنية حديثة | `hmdm-modern-architecture/tasks.md` | أجزاء اختيارية (اختبارات خاصية) لا تزال مفتوحة |
| أجهزة / إضافة-تعديل | `devices-*/tasks.md` | بنود اختيارية غير منجزة |
| إعدادات | `settings-management/tasks.md` | توسّعت بـ General + تبويب Design + أعمدة الأدوار + `idleLogout` في `misc` |

استخدم هذه الملفات كـ **قائمة تحقق** وليس كمصدر وحيد للحقيقة — يجب مطابقة الكود فعلياً.

---

## 8. خلاصة الأولويات المقترحة

1. **لوحة تحكم كاملة** توافق `SummaryResponse` (أو تصميم `.kiro`) بدل الثلاثة بطاقات فقط.
2. **التكامل مع الإضافات** أو توثيق صريح بأن الإضافات تبقى على واجهة الـ WAR القديم فقط حتى إشعار لاحق.
3. **صفحات إدارية مفقودة:** ملفات، أيقونات، تلميحات، تحديث، عملاء (super-admin)، ملف شخصي، استعادة كلمة المرور، 2FA، تسجيل — حسب أولوية المنتج.
4. **i18n** في React إن كان المنتج متعدد اللغات.
5. **توحيد إعدادات Headwind:** فصول `/design`, `/userRoles/common` أو دمجها بوضوح في `/settings`.
6. **مراجعة أذونات `navItems`** مقابل إخفاء عناصر Angular القديمة حسب الدور.

---

## 9. كيفية تحديث هذا المستند

- عند إضافة مسار جديد في `App.tsx` أو استدعاء API جديد: حدّث الأقسام 3 و 5 و 6.
- عند إزالة واجهة Angular من الإنتاج: صحّح القسم 4.
- اربط كل فجوة كبيرة بـ issue/ticket في أداة التتبع إن وُجدت.

---

*هذا التحليل مبني على هيكلة المستودع ومسح ثابت للملفات؛ لا يغني عن اختبار يدوي/تشغيلي لكل ميزة في بيئة staging.*

---

## 10. مصفوفة تتبع التنفيذ (تحديث تلقائي مع تقدم SPA)

| المسار / المنطقة | الحالة |
|------------------|--------|
| لوحة التحكم (إحصاء + مخططات + آخر الأجهزة + تحديث 60 ث) | **مُنفَّذ** (`dashboardService`, محاذاة `GET /private/summary/devices`) |
| تصميم / إعدادات مشتركة / لغة + idle / عميل ضمن الإعدادات | **مُنفَّذ** (`/settings`: General + Design + Device columns، `idleLogout` → `misc`) |
| الملفات | **مُنفَّذ** (`/files`, `filesService`) |
| الأيقونات / التلميحات / التحديثات / الدفع | **مُنفَّذ** |
| عملاء (`/control-panel`) | **مُنفَّذ** (بحث أساسي، تنكر ضمن SPA) |
| ملف شخصي / استعادة كلمة المرور / إعادة التعيين | **مُنفَّذ** |
| 2FA | **مُنفَّذ في الواجهة** (`/twofactor` + استدعاءات `verify/set` + QR blob) — يتطلّب أن يوفّر الـ WAR فعلياً `TwoFactorAuthResource` |
| التسجيل العام | **مُنفَّذ** (مسارات React + `SignupResource`) |
| i18n | **أساسيات** (`i18next` — `en`/`ar`, تبديل EN/AR في الهيدر، سلاسل تسجيل الدخول) |
| الخريطة (`hmdmMap`) | **مسار توثيقي** `/maps` — الأداة لم تُستدعَ من قوالب legacy |
| الإضافات | **إدارة تشغيل/تعطيل** (`PluginSettingsPage` يطابق المنطق القديم) |
| أذونات القائمة | **مُنفَّذ** (`navItems` + `canManageRoles` + `singleCustomer` في الجلسة) |
