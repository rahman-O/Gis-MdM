# تحليل الفجوات: الباكند القديم (Java) ↔ الجديد (Go)

**تاريخ التحليل:** 2026-05-21  
**القديم:** [`backend/`](backend/) — JAX-RS, MyBatis, Liquibase, Tomcat WAR  
**الجديد:** [`serverBackendGo/`](serverBackendGo/) — Gin, Postgres, migrations SQL  
**هذا الملف يجيب:** ما الذي **ينقص في الباكند الجديد (Go)** مقارنة بالقديم (Java)؟

**مراجع مرتبطة:**

- [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md) — حالة الهجرة العامة  
- [`FRONTEND-GO-BACKEND-INTEGRATION.md`](FRONTEND-GO-BACKEND-INTEGRATION.md) — تكامل React مع Go  
- [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md) — مراحل 1–9  
- [`serverBackendGo/docs/parity/`](serverBackendGo/docs/parity/) — تفصيل endpoint لكل وحدة  

---

## 1. الخلاصة التنفيذية

| المؤشر | القيمة التقريبية |
|--------|------------------|
| موارد REST في Java (`*Resource.java`) | **35** |
| وحدات Go مسجّلة (`module.go`) | **29** + `platform/push` |
| **مسارات REST غير موجودة أصلاً في Go** | **2** (`Stats`, `Videos`) |
| **مسارات موجودة لكن سلوك ناقص (⚠️)** | **~15 مجالاً** |
| **مكوّنات خلفية غير منقولة** | **~8** (AuditFilter, MQTT, hooks, …) |
| **تطابق سلوكي إجمالي مع Java** | **~78%** (Phase 9 جزئي) |

### أهم النواقص في Go (قائمة سريعة)

1. **وحدات REST كاملة مفقودة:** `StatsResource`, `VideosResource`  
2. **Plugins:** deviceinfo (export + search/device + إعدادات per-device), devicelog (export + rules)  
3. **أجهزة:** فلاتر بحث متقدمة، إثراء apps/files، `DeviceInfoView` كامل  
4. **عملاء:** نسخ أجهزة/تكوينات افتراضية عند إنشاء tenant، Mailchimp  
5. **ملفات:** تحميل static `/files/*` للوكلاء، quota كامل، `uploadedfiles`  
6. **مزامنة:** `SyncResponseHook` من plugins  
7. **تدقيق:** `AuditFilter` تلقائي (ليس فقط API بحث)  
8. **إشعارات:** MQTT/ActiveMQ (Go يستخدم طابور DB + polling)  
9. **تحديثات:** تنزيل APK بعيد، `sendStats`  
10. **تطبيقات:** plugin hooks عند الحفظ  
11. **Signup:** تكامل Mailchimp ونسخ إعدادات افتراضية كاملة  
12. **مستخدمون:** impersonate مستخدم، مسارات `superadmin/*` (⊘ لـ React غالباً)  
13. **PublicFilesResource** — ⊘ في Java أيضاً لكن servlet ملفات عامة قد يختلف  

---

## 2. منهجية التحليل

تمت المراجعة عبر:

| المصدر | الاستخدام |
|--------|-----------|
| جميع `*Resource.java` تحت `backend/` | قائمة REST الرسمية في Java |
| `handler.go` / `module.go` في `serverBackendGo/internal/modules/` | المسارات المسجّلة في Go |
| `serverBackendGo/docs/parity/*.md` (26 ملفاً) | حالة Done / Partial لكل endpoint |
| `specs/011-complete-migration-gaps/tasks.md` | مهام الإغلاق (45/93 منجزة) |
| مقارنة حقول domain (مثلاً `devices`, `configurations`) | فجوات payload |

**رموز الحالة (من منظور Go):**

| الرمز | المعنى |
|-------|--------|
| ✅ | منقول ويعادل Java للاستخدام اليومي |
| ⚠️ | مسار موجود؛ **سلوك أو حقول ناقصة** — هذا هو «النقص» |
| ❌ | **غير موجود في Go** |
| ⊘ | متعمد خارج النطاق / غير مستخدم في React |

---

## 3. جرد Java — ماذا يوفّر المشروع القديم؟

### 3.1 Core server (`backend/server/` + `jwt/` + `notification/`)

| المورد Java | البادئة | Go module | حالة Go |
|-------------|---------|-----------|---------|
| `AuthResource` | `/rest/public/auth` | `auth` | ✅ |
| `JWTAuthResource` | `/rest/public/jwt` | `auth` | ✅ |
| `SignupResource` | `/rest/public/signup` | `signup` | ⚠️ |
| `PasswordResetResource` | `/rest/public/passwordReset` | `passwordreset` | ✅ |
| `UserResource` | `/rest/private/users` | `users` | ⚠️ |
| `UserRoleResource` | `/rest/private/roles` | `roles` | ✅ |
| `CustomerResource` | `/rest/private/customers` | `customers` | ⚠️ |
| `SettingsResource` | `/rest/private/settings` | `settings` | ✅ |
| `HintResource` | `/rest/private/hints` | `hints` | ✅ |
| `SummaryResource` | `/rest/private/summary` | `summary` | ⚠️ |
| `DeviceResource` | `/rest/private/devices` | `devices` | ⚠️ |
| `GroupResource` | `/rest/private/groups` | `groups` | ✅ |
| `ConfigurationResource` | `/rest/private/configurations` | `configurations` | ⚠️ |
| `ConfigurationFileResource` | `/rest/private/config-files` | `configfiles` | ⚠️ |
| `ApplicationResource` | `/rest/private/applications` | `applications` | ⚠️ |
| `FilesResource` | `/rest/private/web-ui-files` | `files` | ⚠️ |
| `IconResource` | `/rest/private/icons` | `icons` | ✅ |
| `IconFileResource` | `/rest/private/icon-files` | `icons` | ✅ (Phase 9) |
| `PublicResource` | `/rest/public` | `publicapi` | ✅ |
| `PublicFilesResource` | `/rest/public/files` | — | ⊘ |
| `SyncResource` | `/rest/public/sync` | `sync` | ⚠️ |
| `PushApiResource` | `/rest/private/push` | `push` | ✅ |
| `UpdateResource` | `/rest/private/update` | `updates` | ⚠️ |
| `QRCodeResource` | `/rest/public/qr` | `qrcode` | ✅ |
| `StatsResource` | `/rest/public/stats` | — | ❌ |
| `VideosResource` | `/rest/public/videos` | — | ❌ |
| `NotificationResource` | `/rest/notifications` | `notifications` | ✅ |

### 3.2 Plugins (`backend/plugins/`)

| المورد Java | Go | حالة Go |
|-------------|-----|---------|
| `PluginResource` | `plugins/platform` | ✅ |
| `AuditResource` | `plugins/audit` | ⚠️ |
| `MessagingResource` | `plugins/messaging` | ✅ |
| `PushResource` | `plugins/push` | ✅ (+ cron Phase 9) |
| `DeviceInfoResource` + settings | `plugins/deviceinfo` | ⚠️ |
| `DeviceLogResource` + settings | `plugins/devicelog` | ⚠️ |
| **xtra** (Angular فقط) | — | ⊘ |

### 3.3 Two-factor

| المسار | Go |
|--------|-----|
| `/rest/private/twofactor/*` | `twofactor` ✅ |

### 3.4 مكوّنات Java **ليست** REST (لكنها جزء من «الباكند»)

| المكوّن | الموقع | في Go |
|---------|--------|-------|
| `AuditFilter` | `plugins/audit` | ❌ |
| `AuthFilter` / `JWTFilter` | common / jwt | ✅ middleware Gin |
| `LongPollingServlet` | notification | ✅ HTTP handler |
| `NotificationMqttTaskModule` / `MqttThrottledSender` | notification | ❌ |
| `PushScheduleTaskModule` | plugins/push | ✅ `app/scheduler.go` (Phase 9) |
| `SyncResponseHook` (Guice multibind) | SyncResource | ❌ |
| `BackgroundTaskRunnerService` | common | جزئي |
| `CustomerCreatedEventListener` | plugins | ❌ |
| `DeviceLogTaskModule` / `InsertDeviceLogRecordsTask` | devicelog | ❌ |
| `DeviceInfoTaskModule` | deviceinfo | ❌ |
| `MailchimpService` | server signup/customers | ❌ |
| `FileCheckTask` | server | ❌ |
| Liquibase (10+ changelogs) | `backend/**/liquibase/` | 10 migrations Go — **مراجعة يدوية** |

---

## 4. النواقص حسب الأولوية

### 4.1 ❌ غير منقول بالكامل (لا وحدة Go)

| # | Java | المسار | ماذا يفعل في القديم | تأثير النقص في Go |
|---|------|--------|---------------------|-------------------|
| 1 | `StatsResource` | `PUT /rest/public/stats` | يستقبل إحصائيات استخدام/أحداث | الوكلاء أو أدوات ترسل stats → **404** |
| 2 | `VideosResource` | `GET/POST /rest/public/videos/{fileName}` | فيديوهات مساعدة/تدريب | أي رابط قديم للفيديو → **404** |

**إجراء مقترح:** إنشاء `internal/modules/stats` و `internal/modules/videos` أو ⊘ موثّق إن لم تُستخدم.

---

### 4.2 ⚠️ REST موجود — **سلوك أو endpoints ناقصة**

#### أ) أجهزة (`DeviceResource` → `devices`)

| النقص في Go | Java يوفّره | التفاصيل |
|-------------|-------------|----------|
| ~~فلاتر بحث~~ | `DeviceSearchRequest` | **012 US1:** فلاتر React + `sortBy`/`sortDir` — جزئي: `installationStatus` عبر `infojson` فقط؛ `deviceStatuses` غير موجود في schema Go |
| إثراء القائمة | بحث غني | لا `applications` / `files` في صفوف قائمة البحث (موجودة في `GET /number/{n}` → `info`) |
| ~~`DeviceInfoView`~~ | `info` + `infojson` | **012 US1:** `GET /number/{number}` يُرجع `info` |
| `launcherVersion` في القائمة | أعمدة إضافية | غير مُرجَعة في صفوف البحث؛ فلتر `launcherVersion` عبر `infojson` |

**Parity:** [`serverBackendGo/docs/parity/devices.md`](serverBackendGo/docs/parity/devices.md)

---

#### ب) عملاء (`CustomerResource` → `customers`)

| النقص في Go | Java |
|-------------|------|
| Bootstrap tenant | عند `PUT` إنشاء عميل: نسخ أجهزة/تكوينات/إعدادات افتراضية |
| Mailchimp | اشتراك تسويقي اختياري |
| `GET /search/{value}` | موجود في Java (قديم) — Go يعتمد POST `/search` فقط ✅ للـ React |

**Parity:** [`parity/customers.md`](serverBackendGo/docs/parity/customers.md)

---

#### ج) تكوينات (`ConfigurationResource` → `configurations`)

| النقص في Go | Java |
|-------------|------|
| حقول تصميم/launcher كاملة | قد لا تُحفظ كل حقول `Configuration` الكبيرة |
| Push عند الحفظ | ✅ **أُغلق Phase 9** (`platform/push`) — كان Noop سابقاً |
| Plugin hooks | تعديلات plugins عند حفظ التكوين |

---

#### د) ملفات تكوين (`ConfigurationFileResource` → `configfiles`)

| النقص في Go | Java |
|-------------|------|
| `uploadedfiles` + checksum في DB | غالباً قرص فقط في v1 |
| Quota `sizeLimit` | غير مفروض بالكامل |
| ربط كامل مع configurations | مسار POST موجود؛ تكامل DB جزئي |

**Parity:** [`parity/configfiles.md`](serverBackendGo/docs/parity/configfiles.md)

---

#### هـ) ملفات الواجهة (`FilesResource` → `files`)

| النقص في Go | Java |
|-------------|------|
| `GET /files/{filePath}` للوكلاء | servlet/static — **اختياري في dev** |
| Push عند `POST /configurations` (ربط ملف بتكوين) | Phase 9 يغطي configurations/devices؛ ملفات قد تبقى stub |
| `GET /web-ui-files/search/{value}` | ✅ موجود في Go |

**Parity:** [`parity/files.md`](serverBackendGo/docs/parity/files.md)

---

#### و) تطبيقات (`ApplicationResource` → `applications`)

| النقص في Go | Java |
|-------------|------|
| Plugin hooks عند حفظ تطبيق | تكامل plugins/platform |
| أيقونة من رفع مباشر | iconId عبر `icon-files` منفصل ✅ |

**Parity:** [`parity/applications.md`](serverBackendGo/docs/parity/applications.md)

---

#### ز) مزامنة وكلاء (`SyncResource` → `sync`)

| النقص في Go | Java |
|-------------|------|
| `SyncResponseHook` | plugins توسّع `SyncResponse` |
| دمج settings كامل | قد تُفقد حقول من plugins مفعّلة |
| توقيعات/تسجيل آمن | موجودة — راجع env `SECURE_ENROLLMENT` |

**Parity:** [`parity/sync.md`](serverBackendGo/docs/parity/sync.md)

---

#### ح) تحديثات (`UpdateResource` → `updates`)

| النقص في Go | Java |
|-------------|------|
| `POST /` تنزيل APK بعيد | Apply كامل |
| `sendStats` بعد التحديث | مرتبط بـ `StatsResource` ❌ |
| `GET /check` | ✅ |

**Parity:** [`parity/updates.md`](serverBackendGo/docs/parity/updates.md)

---

#### ط) ملخص لوحة (`SummaryResource` → `summary`)

| النقص في Go | Java |
|-------------|------|
| Charts `devicestatuses` حسب config | قد تُرجع قيم فارغة/مبسّطة |
| إحصائيات شهرية غنية | `EmptyDeviceStats` shape موجود؛ بيانات DB قد تكون ناقصة |

**Parity:** [`parity/summary.md`](serverBackendGo/docs/parity/summary.md)

---

#### ي) تسجيل (`SignupResource` → `signup`)

| النقص في Go | Java |
|-------------|------|
| Mailchimp | بعد التحقق من البريد |
| نسخ إعدادات/حدود أجهزة افتراضية | عند `complete` |
| المسارات | `verifyEmail`, `verifyToken`, `complete`, `canSignup` ✅ في Go |

**Parity:** [`parity/auth.md`](serverBackendGo/docs/parity/auth.md)

---

#### ك) مستخدمون (`UserResource` → `users`)

| النقص في Go | Java | ملاحظة |
|-------------|------|--------|
| `GET /{id}` | تفاصيل مستخدم | ⊘ React لا يستخدم |
| `GET /impersonate/{id}` | انتحال مستخدم | ⊘ |
| `GET/PUT /superadmin/*` | إدارة عبر tenants | ⊘ |

مسارات React الأساسية: **✅**

**Parity:** [`parity/users.md`](serverBackendGo/docs/parity/users.md)

---

#### ل) Plugins — **أكبر فجوات REST بعد Stats/Videos**

##### deviceinfo

| Endpoint Java | Go | الحالة |
|---------------|-----|--------|
| `POST .../private/search/dynamic` | ✅ | ✅ |
| `GET .../private/{deviceNumber}` | ✅ | ⚠️ بيانات مبسّطة |
| `PUT .../public/{deviceNumber}` | ✅ | ✅ |
| `GET/PUT ...-plugin-settings/private` | ✅ | ✅ |
| `POST .../private/search/device` | ❌ | **ناقص** |
| `POST .../private/export` | ❌ | **ناقص** |
| `GET ...-plugin-settings/device/{deviceNumber}` | ❌ | **ناقص** |

**Parity:** [`parity/plugins-deviceinfo.md`](serverBackendGo/docs/parity/plugins-deviceinfo.md)

##### devicelog

| Endpoint Java | Go | الحالة |
|---------------|-----|--------|
| `POST .../log/private/search` | ✅ | ✅ |
| `POST .../log/list/{deviceNumber}` (public) | ✅ | ✅ |
| GET/PUT settings + rules CRUD | ✅ | ✅ |
| `POST .../private/search/export` | ❌ | **ناقص** |
| `GET .../rules/{deviceNumber}` | ❌ | **ناقص** |
| معالجة دفعات `InsertDeviceLogRecordsTask` | ❌ | **خلفية ناقصة** |

**Parity:** [`parity/plugins-devicelog.md`](serverBackendGo/docs/parity/plugins-devicelog.md)

##### audit

| النقص | Java | Go |
|-------|------|-----|
| `AuditFilter` تلقائي على كل طلب private | servlet filter | ❌ — **بحث فقط** `POST .../log/search` |

**Parity:** [`parity/plugins-audit.md`](serverBackendGo/docs/parity/plugins-audit.md)

##### push (plugin)

| العنصر | Java | Go |
|--------|------|-----|
| CRUD + searchTasks + task | ✅ | ✅ |
| `PushScheduleTaskModule` cron | ✅ | ✅ Phase 9 scheduler |
| إرسال FCM/MQTT مباشر | ActiveMQ/MQTT | ❌ — **queue + polling** |

**Parity:** [`parity/push.md`](serverBackendGo/docs/parity/push.md)

---

### 4.3 ✅ منقول بشكل كافٍ (لا يُعدّ «نقصاً» للعمل اليومي)

`auth`, `passwordreset`, `twofactor`, `roles`, `groups`, `hints`, `settings` (أساسي), `push` API, `notifications` + polling, `qrcode`, `plugins/platform`, `plugins/messaging`, `icons` metadata, `icon-files`, `publicapi` (name/logo/APK upload), `applications` CRUD, `configurations` CRUD, `customers` search/impersonate, `users` CRUD للـ React.

---

## 5. جدول Endpoints — Java موجود وGo **ناقص**

```
❌ PUT  /rest/public/stats
❌ POST /rest/public/videos/{fileName}
❌ GET  /rest/public/videos/{fileName}

❌ POST /rest/plugins/deviceinfo/deviceinfo/private/search/device
❌ POST /rest/plugins/deviceinfo/deviceinfo/private/export
❌ GET  /rest/plugins/deviceinfo/deviceinfo-plugin-settings/device/{deviceNumber}

❌ POST /rest/plugins/devicelog/log/private/search/export
❌ GET  /rest/plugins/devicelog/log/rules/{deviceNumber}

⊘ GET  /rest/public/files/{filePath}          (PublicFilesResource — قديم)
⊘ GET  /rest/private/users/{id}
⊘ GET  /rest/private/users/impersonate/{id}
⊘ GET  /rest/private/users/superadmin/*
```

---

## 6. فجوات البنية التحتية (ليست endpoints فقط)

| المجال | Java | النقص في Go |
|--------|------|-------------|
| **Push للأجهزة** | MQTT + `pushmessages` + polling | طابور Postgres + long poll ✅ يعمل؛ **ليس نفس مسار MQTT** |
| **جدولة push** | `PushScheduleTaskModule` | ✅ cron ~60s |
| **تدقيق تلقائي** | `AuditFilter` | ❌ |
| **Sync plugins** | `SyncResponseHook` | ❌ |
| **أحداث إنشاء عميل** | listeners في plugins | ❌ |
| **سجلات جهاز دفعات** | `InsertDeviceLogRecordsTask` | ❌ |
| **تسويق** | Mailchimp | ❌ |
| **Schema** | Liquibase متعدد plugins | 10 migrations — تحقق من جداول plugins عند مشاكل DB |
| **ملفات على القرص** | `/var/lib/hmdm/files` + servlet | `FILES_DIRECTORY` محلي؛ **لا Gin static `/files`** افتراضياً |

---

## 7. Phase 9 — ما أُغلق وما بقي

| البند | قبل Java-only | Go الآن |
|-------|---------------|---------|
| Push عند حفظ configuration | Java `PushService` | ✅ `platform/push` |
| Notify `applicationSettings` | Java | ✅ |
| cron `plugin_push_schedule` | Java `PushScheduleTaskModule` | ✅ `scheduler.go` |
| `POST /icon-files` | Java `IconFileResource` | ✅ |

**ما زال ناقصاً (مهام T046–T093):** exports plugins، audit middleware، stats، videos، customers bootstrap، device search enrichment، static files، sync hooks، إلخ — راجع [`specs/011-complete-migration-gaps/tasks.md`](specs/011-complete-migration-gaps/tasks.md).

---

## 8. خطة إغلاق النواقص (مقترحة)

| الأولوية | النقص | Spec / مهمة |
|----------|-------|-------------|
| **P0** | deviceinfo + devicelog exports/endpoints | US4 |
| **P0** | AuditFilter + sync hooks + customers bootstrap | US5 |
| **P1** | devices search filters + DeviceInfo في الاستجابة | US6 / devices |
| **P1** | configfiles quota + `uploadedfiles` | US5 |
| **P1** | `/files/*` static للوكلاء | US6 |
| **P2** | `stats`, `videos` modules | US6 |
| **P2** | updates APK download + sendStats | US6 |
| **P2** | applications plugin hooks | Polish |
| **P3** | MQTT (فقط إن polling غير كافٍ) | قرار تشغيل |
| **⊘** | Mailchimp, superadmin users, xtra | توثيق خارج النطاق |

---

## 9. قائمة تحقق: «هل Go يغني عن Java؟»

```text
استبدال Java في الإنتاج — آمن إذا:
  [ ] كل سيناريوهات React UAT تمر على Go
  [ ] وكلاء: enroll + sync + notifications + تحديث APK (إن مستخدم)
  [ ] لا تعتمد على: stats, videos, deviceinfo export, devicelog export, MQTT إلزامي
  [ ] فلاتر أجهزة المتقدمة مقبولة مؤقتاً أو مُنفَّذة
  [ ] Phase 9 = done في MIGRATION.md

غير آمن إذا:
  [ ] تعتمد على واجهة Angular القديمة لـ plugins export
  [ ] تعتمد على AuditFilter لامتثال كامل
  [ ] تعتمد على إنشاء tenant بنسخ كامل افتراضي دون عمل يدوي
```

---

## 10. ملخص: «ماذا ينقص الباكند الجديد؟»

**باختصار:** Go يغطي **أغلب REST** الذي تحتاجه React والوكلاء (~92% من المسارات)، لكن **ينقصه:**

1. **وحدتان كاملتان:** إحصائيات عامة + فيديوهات.  
2. **5+ endpoints plugins** (تصدير، قواعد devicelog، إعدادات deviceinfo per-device).  
3. **سلوك خلفي:** تدقيق تلقائي، hooks مزامنة، MQTT، مهام دفعات devicelog، Mailchimp.  
4. **عمق بيانات:** بحث أجهزة غني، telemetry جهاز، charts ملخص، quota ملفات، bootstrap عملاء.  
5. **بنية ملفات الوكلاء:** خدمة static `/files/*` (اختياري لكن Java يوفّرها).

للتفاصيل التشغيلية مع الواجهة: [`FRONTEND-GO-BACKEND-INTEGRATION.md`](FRONTEND-GO-BACKEND-INTEGRATION.md).  
لحالة المشروع الكاملة: [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md).

---

*حدّث هذا الملف عند إغلاق مهام `011-complete-migration-gaps` أو إضافة وحدات `stats` / `videos`.*
