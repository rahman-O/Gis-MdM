# تحليل اكتمال النقل: Java → Go

**تاريخ التحليل:** 2026-05-21  
**المصدر القديم:** [`backend/`](backend/) (Java / JAX-RS / MyBatis / Liquibase)  
**المصدر الجديد:** [`serverBackendGo/`](serverBackendGo/) (Go / Gin / Postgres migrations)  
**الواجهة:** [`frontend/`](frontend/) (React) — نفس مسارات `/rest/*`  
**خارطة الطريق:** [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md)  
**مواصفة إغلاق الفجوات:** [`specs/011-complete-migration-gaps/`](specs/011-complete-migration-gaps/)  
**تفصيل نواقص Go مقابل Java:** [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md)

---

## الإجابة المباشرة

| السؤال | الحكم |
|--------|--------|
| هل النقل **مكتمل 100%** مقارنة بـ Java؟ | **لا** |
| هل Go **جاهز تشغيلياً** للوحة React والوكلاء في المسارات الأساسية؟ | **نعم — تقريباً (~85–90%)** |
| هل يمكن إيقاف Java **بدون مخاطر** اليوم؟ | **لا — يُنصح بـ UAT وإغلاق Phase 9 أولاً** |

**الخلاصة:** المشروع الجديد **ليس نسخة 1:1** من القديم، لكنه **يغطي أغلب MDM** (مصادقة، عملاء، أجهزة، تكوينات، ملفات، وكلاء، plugins). النقل **مكتمل تشغيلياً للنواة** و**غير مكتمل رسمياً** من ناحية التطابق السلوكي والخلفية.

---

## نسب الاكتمال (تقدير)

| البُعد | التقدير | ملاحظة |
|--------|---------|--------|
| البنية (وحدات + مسارات REST أساسية) | **~92%** | 29 وحدة Go مقابل 35 مورد REST في Java |
| سلوك مطابق لـ React الحالي | **~88%** | مسارات الواجهة الرئيسية موجودة |
| سلوك مطابق للوكلاء (sync + notifications) | **~82%** | polling يعمل؛ MQTT/FCM غير منقول |
| Plugins (مسارات + cron + push queue) | **~80%** | exports وaudit تلقائي ناقص |
| تكاملات خلفية (filters, hooks, batch jobs) | **~35%** | AuditFilter، SyncResponseHook، … |
| **تطابق إجمالي مع Java WAR** | **~78%** | قبل إغلاق T045–T093 |

```
البنية REST أساسية     ████████████████████░  ~92%
React + لوحة التحكم      ████████████████████░  ~88%
وكلاء MDM                █████████████████░░░░  ~82%
Plugins                  ████████████████░░░░░  ~80%
خلفية / hooks / MQTT     ███████░░░░░░░░░░░░░░  ~35%
─────────────────────────────────────────────────
تطابق سلوكي إجمالي       ███████████████░░░░░░  ~78%
```

---

## مراحل الهجرة (Phases 1–9)

| Phase | الحالة في `MIGRATION.md` | الواقع التشغيلي |
|-------|--------------------------|------------------|
| **1** — auth, JWT | **done** | ✅ مكتمل |
| **1b** — signup, passwordreset, twofactor | **done** | ✅ أساسي؛ ⚠️ بدون Mailchimp |
| **2** — users, roles | **done** | ✅ لـ React؛ ⊘ impersonate / superadmin |
| **3** — customers, settings, hints, summary | **done** | ⚠️ bootstrap عميل؛ charts مبسّطة |
| **4** — devices, groups, configurations list | **done** | ⚠️ بحث أجهزة غني ناقص |
| **5** — applications, configurations, configfiles | **done** | ✅ CRUD؛ ⚠️ quota ملفات، plugin hooks |
| **6** — files, icons, publicapi | **done** | ✅؛ ⚠️ `/files/*` static للوكلاء |
| **7** — sync, push, notifications, updates, qrcode | **done** | ⚠️ sync hooks؛ updates APK؛ ليس MQTT |
| **8** — plugins | **done** | ⚠️ export؛ audit بحث فقط |
| **9** — gap closure | **partial** | ✅ push queue + cron + icon-files؛ باقي المهام معلّقة |

**تقدم Phase 9:** **45 / 93** مهمة في [`specs/011-complete-migration-gaps/tasks.md`](specs/011-complete-migration-gaps/tasks.md) (T001–T045 مكتملة تقريباً؛ T046–T093 متبقية).

### ما أُنجز في Phase 9 (جزئي)

| البند | Java | Go |
|-------|------|-----|
| إشعار بعد حفظ configuration | `PushService` | `internal/platform/push` → `configUpdated` في طابور |
| إشعار `applicationSettings/notify` | نفس الأسلوب | `appConfigUpdated` |
| cron جدولة push | `PushScheduleTaskModule` | `internal/app/scheduler.go` (~60s) |
| رفع أيقونة ملف | `IconFileResource` | `POST /rest/private/icon-files` — [parity/icon-files.md](serverBackendGo/docs/parity/icon-files.md) |

**ما لم يُغلق في Phase 9:** exports (deviceinfo/devicelog)، audit middleware، stats، videos، customers bootstrap، quota ملفات، sync hooks، إلخ.

---

## رموز الحالة

| الرمز | المعنى |
|-------|--------|
| ✅ | منقول ومناسب للاستخدام مع React/الوكلاء |
| ⚠️ | مسار موجود؛ سلوك أو تكامل ناقص (راجع parity) |
| ❌ | غير منقول — لا وحدة Go |
| ⊘ | خارج النطاق المتفق / غير مستخدم في React |

---

## Java: موارد REST (35 ملف `*Resource.java`)

### Core server (`backend/server/` + `jwt/` + `notification/`)

| المورد Java | المسار التقريبي | Go module | الحالة |
|-------------|-----------------|-----------|--------|
| `AuthResource` | `/rest/public/auth` | `auth` | ✅ |
| `JWTAuthResource` | `/rest/public/jwt` | `auth` | ✅ |
| `SignupResource` | `/rest/public/signup` | `signup` | ⚠️ |
| `PasswordResetResource` | `/rest/public/passwordReset` | `passwordreset` | ✅ |
| `UserResource` | `/rest/private/users` | `users` | ⚠️ / ⊘ |
| `UserRoleResource` | `/rest/private/roles` | `roles` | ✅ |
| `CustomerResource` | `/rest/private/customers` | `customers` | ⚠️ |
| `SettingsResource` | `/rest/private/settings` | `settings` | ✅ |
| `HintResource` | `/rest/private/hints` | `hints` | ✅ |
| `SummaryResource` | `/rest/private/summary` | `summary` | ⚠️ |
| `DeviceResource` | `/rest/private/devices` | `devices` | ⚠️ |
| `GroupResource` | `/rest/private/groups` | `groups` | ✅ |
| `ConfigurationResource` | `/rest/private/configurations` | `configurations` | ✅ |
| `ConfigurationFileResource` | `/rest/private/config-files` | `configfiles` | ⚠️ |
| `ApplicationResource` | `/rest/private/applications` | `applications` | ⚠️ |
| `FilesResource` | `/rest/private/files` | `files` | ⚠️ |
| `IconResource` | `/rest/private/icons` | `icons` | ✅ |
| `IconFileResource` | `/rest/private/icon-files` | `icons` | ✅ (Phase 9) |
| `PublicResource` | `/rest/public/...` | `publicapi` | ✅ |
| `PublicFilesResource` | — | — | ⊘ |
| `SyncResource` | `/rest/public/sync` | `sync` | ⚠️ |
| `PushApiResource` | `/rest/private/push` | `push` | ✅ |
| `UpdateResource` | `/rest/public/update` | `updates` | ⚠️ |
| `QRCodeResource` | `/rest/public/qrcode` | `qrcode` | ✅ |
| `StatsResource` | `PUT /rest/public/stats` | — | ❌ |
| `VideosResource` | `/rest/public/videos` | — | ❌ |
| `NotificationResource` | `/rest/notification` | `notifications` | ✅ |

### Plugins (`backend/plugins/`)

| المورد Java | Go module | الحالة |
|-------------|-----------|--------|
| `PluginResource` | `plugins/platform` | ✅ |
| `AuditResource` | `plugins/audit` | ⚠️ (بحث فقط؛ بدون `AuditFilter`) |
| `MessagingResource` | `plugins/messaging` | ✅ |
| `PushResource` (plugin) | `plugins/push` | ✅ (+ cron Phase 9) |
| `DeviceInfoResource` + settings | `plugins/deviceinfo` | ⚠️ |
| `DeviceLogResource` + settings | `plugins/devicelog` | ⚠️ |

### Two-factor (ليست Resource منفصلة في server)

| المسار | Go | الحالة |
|--------|-----|--------|
| `/rest/private/twofactor/*` | `twofactor` | ✅ — [parity/auth.md](serverBackendGo/docs/parity/auth.md) |

### ⊘ خارج REST

| العنصر | ملاحظة |
|--------|--------|
| Plugin **xtra** | Angular assets فقط؛ لا REST |
| **Mailchimp** | اختياري في Java أصلاً |

---

## Go: الوحدات المسجّلة (29 `module.go`)

```
auth, signup, passwordreset, twofactor,
users, roles, customers, settings, hints, summary,
devices, groups, applications, configurations, configfiles,
files, icons, publicapi,
sync, push, notifications, updates, qrcode,
plugins/platform, plugins/audit, plugins/push,
plugins/messaging, plugins/deviceinfo, plugins/devicelog
```

**منصة مشتركة (ليست module مستقل):** `internal/platform/push` — notifier للطابور بعد Phase 9.

**وحدات Java بلا مكافئ Go:** `stats`, `videos`.

---

## تفصيل الفجوات

### ❌ غير منقول

| # | Java | المسار / السلوك | التأثير |
|---|------|-----------------|---------|
| 1 | `StatsResource` | `PUT /rest/public/stats` | إحصائيات usage من الوكلاء/الواجهة |
| 2 | `VideosResource` | `POST/GET /rest/public/videos/{fileName}` | فيديوهات مساعدة (إن وُجدت في UI) |

### ⚠️ منقول جزئياً — ملخص

| المجال | الفجوة الرئيسية |
|--------|------------------|
| **Signup** | بدون Mailchimp؛ نسخ إعدادات افتراضية غير كاملة |
| **Customers** | `PUT` إنشاء tenant بدون نسخ أجهزة/تكوينات افتراضية |
| **Devices** | فلاتر `mdmMode`, `launcherVersion`, `deviceStatuses`؛ إثراء apps/files في البحث |
| **Summary** | `devicestatuses` / charts — مبسّط |
| **Config files** | `uploadedfiles` + quota `sizeLimit` |
| **Applications** | plugin hooks عند الحفظ |
| **Files** | تحميل static `/files/*` للوكلاء |
| **Configurations** | push عند الحفظ — **✅ Phase 9** (queue) |
| **Sync** | `SyncResponseHook` من plugins |
| **Updates** | تنزيل APK بعيد؛ `sendStats` |
| **Push transport** | Java: MQTT + polling؛ Go: **DB queue + polling** (يعمل؛ بنية مختلفة) |
| **deviceinfo** | `POST .../private/search/device`, `POST .../private/export`, `GET settings/device/{deviceNumber}` |
| **devicelog** | `POST .../private/search/export`, `GET .../rules/{deviceNumber}` |
| **audit** | لا `AuditFilter` تلقائي |

### Endpoints Java الناقصة في Go (مرجع سريع)

```
POST .../deviceinfo/private/search/device
POST .../deviceinfo/private/export
GET  .../deviceinfo-plugin-settings/device/{deviceNumber}
POST .../devicelog/log/private/search/export
GET  .../devicelog/log/rules/{deviceNumber}
PUT  /rest/public/stats
POST /rest/public/videos/{fileName}
GET  /rest/public/videos/{fileName}
```

---

## مكونات Java غير REST (خلفية)

| المكوّن | الموقع التقريبي | في Go |
|---------|-----------------|-------|
| `AuditFilter` | `plugins/audit` | ❌ |
| `LongPollingServlet` | `notification` | ✅ (HTTP handler) |
| `BackgroundTaskRunnerService` | `common` | جزئي — cron push ✅ Phase 9 |
| `SyncResponseHook` + Guice | `SyncResource` | ❌ |
| `CustomerCreatedEventListener` | devicelog / deviceinfo | ❌ |
| `InsertDeviceLogRecordsTask` | devicelog | ❌ |
| `PushSenderMqtt` / ActiveMQ | push | ❌ (استُبدل بـ queue + poll) |
| MyBatis DAOs | `backend/common/...` | Postgres repos — نفس الجداول غالباً |
| Liquibase | `backend/**/liquibase/` | `serverBackendGo/db/migrations/` (10 ملفات) |

---

## وثائق التطابق (parity)

26 ملفاً تحت [`serverBackendGo/docs/parity/`](serverBackendGo/docs/parity/).  
ملفات تحتوي على **Partial** — راجعها قبل إعلان اكتمال:

| الملف | موضوع الفجوة |
|-------|----------------|
| [customers.md](serverBackendGo/docs/parity/customers.md) | bootstrap عند الإنشاء |
| [devices.md](serverBackendGo/docs/parity/devices.md) | بحث غني؛ push notify ✅ |
| [configfiles.md](serverBackendGo/docs/parity/configfiles.md) | quota |
| [files.md](serverBackendGo/docs/parity/files.md) | static agent؛ push |
| [applications.md](serverBackendGo/docs/parity/applications.md) | plugin hooks |
| [sync.md](serverBackendGo/docs/parity/sync.md) | hooks |
| [updates.md](serverBackendGo/docs/parity/updates.md) | APK؛ sendStats |
| [summary.md](serverBackendGo/docs/parity/summary.md) | charts |
| [plugins-audit.md](serverBackendGo/docs/parity/plugins-audit.md) | AuditFilter |
| [plugins-deviceinfo.md](serverBackendGo/docs/parity/plugins-deviceinfo.md) | export / search/device |
| [plugins-devicelog.md](serverBackendGo/docs/parity/plugins-devicelog.md) | export / rules |
| [auth.md](serverBackendGo/docs/parity/auth.md) | signup / 2FA |
| [configurations.md](serverBackendGo/docs/parity/configurations.md) | push على الحفظ ✅ |

---

## جاهزية React والوكلاء

### يعمل على Go (مسارات React الرئيسية)

- تسجيل الدخول (session + JWT)، 2FA، إعادة تعيين كلمة المرور
- مستخدمون، أدوار، عملاء (بحث + impersonate)
- أجهزة، مجموعات، تطبيقات، تكوينات، ملفات، أيقونات (+ رفع icon-files)
- إعدادات، تلميحات، ملخص لوحة
- push يدوي + plugin push + جدولة cron
- مزامنة وكيل، إشعارات polling، QR، تحديثات أساسية
- plugins: platform، messaging، audit (بحث)، deviceinfo/devicelog (أساسي)

### قد يفشل أو يكون ناقصاً

- تصدير deviceinfo / devicelog من الواجهة
- `/rest/public/videos` و `/rest/public/stats`
- توقع **MQTT فوري** بدل polling (الوكيل عادة يدعم polling)
- سيناريوهات tenant جديدة تحتاج نسخ أجهزة/تكوينات افتراضية كاملة

---

## معيار «اكتمال النقل» الرسمي

يُعتبر النقل **مكتملاً** عندما:

1. كل `*Resource.java` إما **✅** أو **⊘** موثّق في parity
2. لا توجد فجوات **⚠️ حرجة** في: push، cron، exports، audit تلقائي، agents
3. `cd serverBackendGo && go test ./...` + smoke React + agent enroll/sync
4. Phase 9 في [`MIGRATION.md`](serverBackendGo/docs/MIGRATION.md) → **done** (وليس partial)
5. إكمال **T046–T093** في [`tasks.md`](specs/011-complete-migration-gaps/tasks.md) أو ⊘ صريح لما لا يُستخدم

---

## الخطوات التالية (مقترحة)

| الأولوية | المهمة | Spec / tasks |
|----------|--------|----------------|
| P1 | deviceinfo + devicelog exports وrules | US4 — T046+ |
| P1 | `AuditFilter` + sync hooks + customers bootstrap | US5 |
| P2 | `stats`, `videos`, updates APK، devices search enrichment | US6 |
| P3 | تحديث parity القديمة؛ `make swagger`؛ UAT كامل | Polish |
| اختياري | MQTT/FCM | فقط إن كان polling غير كافٍ في الإنتاج |

**فرع العمل:** `011-complete-migration-gaps` — راجع [`specs/011-complete-migration-gaps/quickstart.md`](specs/011-complete-migration-gaps/quickstart.md).

---

## مراجع

| الوثيقة | الغرض |
|---------|--------|
| [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md) | خارطة مراحل 1–9 |
| [`serverBackendGo/docs/NEXT_STEPS.md`](serverBackendGo/docs/NEXT_STEPS.md) | أولويات التنفيذ |
| [`specs/011-complete-migration-gaps/spec.md`](specs/011-complete-migration-gaps/spec.md) | متطلبات إغلاق الفجوات |
| [`specs/011-complete-migration-gaps/tasks.md`](specs/011-complete-migration-gaps/tasks.md) | مهام قابلة للتنفيذ (93) |
| [`backend/server/.../resource/`](backend/server/src/main/java/com/hmdm/rest/resource/) | مصدر Java للـ REST |
| [`.specify/memory/constitution.md`](.specify/memory/constitution.md) | مبادئ الهندسة |

---

*هذا الملف يلخّص حالة النقل وقت كتابته. عند إغلاق مهام Phase 9 حدّث الأقسام §«Phase 9» و§«نسب الاكتمال» و§«Endpoints الناقصة».*
