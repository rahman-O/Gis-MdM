# تحليل دعم الأنظمة المتعددة — Windows, macOS, Linux + Flutter

## المقدمة

هذا المستند يحلل إمكانية توسيع نظام MDM الحالي (المبني على Android) لدعم أنظمة Windows, macOS, Linux باستخدام Flutter كعميل موحّد لجميع المنصات (بما في ذلك استبدال تطبيق Android الأصلي).

---

## 1. تحليل البنية الحالية

### 1.1 الباك إند (serverBackendGo)

الباك إند مبني بهيكلية modular نظيفة (28 module) ويتواصل مع الأجهزة عبر REST API:

```
┌─────────────────────────────────────────────────────────────────┐
│ serverBackendGo (Go + Gin + PostgreSQL)                          │
├─────────────────────────────────────────────────────────────────┤
│ Public API (بدون auth):                                          │
│   /public/sync/configuration/:deviceId  ← بروتوكول المزامنة     │
│   /public/sync/info                     ← تقارير الجهاز         │
│   /public/qr/:key                       ← QR enrollment         │
│   /rest/notification/polling/:deviceId  ← push polling          │
├─────────────────────────────────────────────────────────────────┤
│ Private API (JWT auth):                                          │
│   /private/devices, /private/profiles, /private/applications    │
│   /private/enrollment-routes, /private/configurations           │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 تطبيق Android (hmdm-android)

تطبيق Android الأصلي يعمل كـ **Device Owner** ويقوم بـ:
- استقبال الإعدادات من السيرفر عبر `/public/sync/configuration`
- تثبيت/تحديث/حذف التطبيقات (APK)
- تطبيق سياسات MDM (kiosk, WiFi, GPS, Bluetooth, etc.)
- إرسال معلومات الجهاز (battery, IMEI, location)
- استقبال أوامر Push عبر polling
- العمل كـ Launcher (شاشة رئيسية بديلة)

### 1.3 بروتوكول المزامنة (Sync Protocol)

```json
// SyncResponse — ما يستقبله الجهاز
{
  "deviceId": "device-001",
  "configurationId": 5,
  "applications": [
    { "pkg": "com.app", "version": "1.0", "url": "https://..." }
  ],
  "files": [
    { "devicePath": "/data/config.json", "url": "https://..." }
  ],
  "applicationSettings": [...],
  // سياسات MDM:
  "kioskMode": false,
  "gps": true,
  "bluetooth": true,
  "brightness": 200,
  "volume": 80,
  "passwordMode": "strong",
  ...
}
```

---

## 2. تصنيف المكونات: خاص بـ Android vs عام

### 2.1 مكونات خاصة بـ Android (تحتاج بديل لكل منصة)

| المكون | السبب | البديل المطلوب |
|--------|-------|---------------|
| QR Enrollment (Device Owner Provisioning) | بروتوكول Android فقط | enrollment token / deep link / manual |
| APK Distribution | صيغة Android فقط | MSI/EXE (Win), DMG/PKG (Mac), DEB/AppImage (Linux), Flutter bundle |
| Kiosk Mode | Android Launcher API | OS-specific kiosk (Win: Assigned Access, Mac: Single App Mode) |
| Device Owner APIs | Android DevicePolicyManager | OS-specific MDM APIs (Win: MDM CSP, Mac: MDM profiles) |
| IMEI/Serial collection | Android TelephonyManager | OS-specific hardware info APIs |
| Split APK (arm64/armeabi) | Android ABI | لا حاجة — Flutter يبني لكل منصة |
| GPS/Bluetooth/WiFi control | Android system APIs | OS-specific APIs (محدودة على desktop) |
| System Update management | Android SystemUpdatePolicy | OS-specific (WSUS, softwareupdate, apt) |

### 2.2 مكونات عامة (تعمل لجميع المنصات بدون تغيير)

| المكون | السبب |
|--------|-------|
| Sync Protocol (REST API) | HTTP + JSON — أي منصة تدعمه |
| Push/Notification (polling) | HTTP polling — لا يعتمد على Firebase |
| Multi-tenant (customers, users, roles) | منطق أعمال بحت |
| Device Tree + Groups | هيكلية تنظيمية |
| Profile Versioning + Artifacts | سياسات مجمّعة كـ JSON |
| Application Settings (key-value) | بيانات عامة |
| File Distribution (URL-based) | تحميل ملفات عبر HTTP |
| Audit, Messaging, DeviceLog plugins | بيانات نصية |
| Enrollment Routes (structure) | الهيكل عام — فقط QR خاص بـ Android |

---

## 3. خطة التحويل إلى Flutter + دعم المنصات المتعددة

### 3.1 هيكلية Flutter Client الموحّد

```
flutter_mdm_agent/
├── lib/
│   ├── core/
│   │   ├── sync/              ← بروتوكول المزامنة (مشترك)
│   │   ├── push/              ← polling notifications (مشترك)
│   │   ├── enrollment/        ← تسجيل الجهاز (مشترك)
│   │   └── settings/          ← إعدادات التطبيق (مشترك)
│   ├── platform/
│   │   ├── android/           ← Device Owner, APK install, kiosk
│   │   ├── windows/           ← MDM CSP, MSI install, Assigned Access
│   │   ├── macos/             ← MDM profiles, PKG install, Single App
│   │   └── linux/             ← systemd, DEB install, kiosk via X11
│   ├── features/
│   │   ├── app_management/    ← تثبيت/تحديث التطبيقات
│   │   ├── policy_engine/     ← تطبيق السياسات
│   │   ├── device_info/       ← جمع معلومات الجهاز
│   │   └── ui/                ← واجهة المستخدم (launcher/kiosk)
│   └── main.dart
├── android/
├── windows/
├── macos/
├── linux/
└── pubspec.yaml
```

### 3.2 التغييرات المطلوبة في الباك إند

#### أ) إضافة حقل `platform` للأجهزة

```sql
-- Migration: add platform column
ALTER TABLE devices
  ADD COLUMN IF NOT EXISTS platform VARCHAR(20) NOT NULL DEFAULT 'android';
-- Values: 'android', 'windows', 'macos', 'linux', 'ios'
```

#### ب) توسيع Sync Response ليكون platform-aware

```go
// إضافة حقل platform في DeviceInfo
type DeviceInfo struct {
    DeviceID     string `json:"deviceId"`
    Platform     string `json:"platform"`     // NEW: android/windows/macos/linux
    OSVersion    string `json:"osVersion"`    // NEW: replaces androidVersion
    // ... existing fields
}

// إضافة حقل platform في SyncApplication
type SyncApplication struct {
    // ... existing fields
    Platform    string `json:"platform,omitempty"`    // NEW: target platform
    Installer   string `json:"installer,omitempty"`   // NEW: apk/msi/dmg/deb/flutter
}
```

#### ج) توسيع Application entity لدعم حزم متعددة

```sql
-- Migration: platform-specific application packages
CREATE TABLE application_packages (
    id SERIAL PRIMARY KEY,
    application_version_id INT REFERENCES applicationversions(id),
    platform VARCHAR(20) NOT NULL,  -- android, windows, macos, linux
    arch VARCHAR(20),               -- x64, arm64, universal
    url TEXT NOT NULL,
    file_path TEXT,
    checksum VARCHAR(128),
    installer_type VARCHAR(20)      -- apk, msi, exe, dmg, pkg, deb, appimage
);
```

#### د) Enrollment بدون QR (للمنصات غير Android)

```go
// New enrollment endpoint for non-Android platforms
// POST /public/enroll
type PlatformEnrollRequest struct {
    Platform     string `json:"platform"`      // windows/macos/linux
    DeviceID     string `json:"deviceId"`      // hostname or generated UUID
    EnrollToken  string `json:"enrollToken"`   // token from admin panel
    OSVersion    string `json:"osVersion"`
    Hostname     string `json:"hostname"`
    SerialNumber string `json:"serialNumber"`
}
```

#### هـ) سياسات خاصة بكل منصة

```go
// Platform-specific policy extensions in SyncResponse
type PlatformPolicy struct {
    // Windows-specific
    WindowsUpdatePolicy  *string `json:"windowsUpdatePolicy,omitempty"`
    BitLockerEnabled     *bool   `json:"bitLockerEnabled,omitempty"`
    FirewallEnabled      *bool   `json:"firewallEnabled,omitempty"`
    
    // macOS-specific
    FileVaultEnabled     *bool   `json:"fileVaultEnabled,omitempty"`
    GatekeeperEnabled    *bool   `json:"gatekeeperEnabled,omitempty"`
    
    // Linux-specific
    SELinuxMode          *string `json:"selinuxMode,omitempty"`
    FirewalldEnabled     *bool   `json:"firewalldEnabled,omitempty"`
    
    // Cross-platform
    ScreenLockTimeout    *int    `json:"screenLockTimeout,omitempty"`
    PasswordPolicy       *string `json:"passwordPolicy,omitempty"`
    AutoUpdateEnabled    *bool   `json:"autoUpdateEnabled,omitempty"`
}
```

---

## 4. التغييرات التفصيلية

### 4.1 الباك إند — ملخص التغييرات

| الملف/الوحدة | التغيير | الأولوية |
|-------------|---------|---------|
| `devices` domain | إضافة `Platform`, `OSVersion`, `Hostname` | عالية |
| `sync` module | دعم `X-Platform` header, platform-aware response | عالية |
| `applications` domain | إضافة `application_packages` table | عالية |
| `enrollment_routes` | إضافة enrollment token (بديل QR لـ desktop) | عالية |
| `qrcode` module | إبقاء كما هو (Android فقط) | — |
| `profiles` | إضافة platform filter للتطبيقات | متوسطة |
| `push/notifications` | لا تغيير (polling يعمل لكل المنصات) | — |
| `plugins/deviceinfo` | توسيع لدعم معلومات OS مختلفة | متوسطة |
| Admin UI (frontend) | إضافة فلتر platform في قائمة الأجهزة | متوسطة |

### 4.2 Flutter Client — المكونات الأساسية

| المكون | المسؤولية | المنصات |
|--------|----------|---------|
| `SyncService` | مزامنة الإعدادات مع السيرفر | الكل |
| `PushPoller` | استقبال أوامر Push | الكل |
| `AppInstaller` | تثبيت/تحديث التطبيقات | خاص بكل منصة |
| `PolicyEngine` | تطبيق السياسات | خاص بكل منصة |
| `DeviceInfoCollector` | جمع معلومات الجهاز | خاص بكل منصة |
| `EnrollmentManager` | تسجيل الجهاز | خاص بكل منصة |
| `KioskManager` | وضع القفل | خاص بكل منصة |
| `UIShell` | واجهة المستخدم/Launcher | مشترك مع تخصيص |

### 4.3 Enrollment لكل منصة

| المنصة | طريقة التسجيل | التفاصيل |
|--------|-------------|---------|
| **Android** | QR Code (Device Owner) | كما هو — بروتوكول Android الأصلي |
| **Windows** | Enrollment Token + MSI | المشرف يولّد token → المستخدم يثبّت MSI → التطبيق يسجّل تلقائياً |
| **macOS** | Enrollment Token + PKG | نفس المبدأ — أو عبر MDM profile (.mobileconfig) |
| **Linux** | Enrollment Token + Script | `curl install.sh \| bash` أو DEB/RPM package |
| **Flutter (all)** | Deep Link / Token | رابط `mdm://enroll?token=xxx` أو إدخال يدوي |

---

## 5. مراحل التنفيذ المقترحة

### المرحلة 1: تجهيز الباك إند (2-3 أسابيع)

1. إضافة `platform` column لجدول `devices`
2. إنشاء `application_packages` table
3. إضافة enrollment token endpoint (`POST /public/enroll`)
4. توسيع `SyncResponse` بحقل `platformPolicy`
5. إضافة `X-Platform` header handling في sync module
6. تعديل admin UI لعرض platform في قائمة الأجهزة

### المرحلة 2: Flutter Agent — Core (3-4 أسابيع)

1. إنشاء مشروع Flutter (android + windows + macos + linux)
2. بناء `SyncService` (REST client + periodic sync)
3. بناء `PushPoller` (HTTP long-polling)
4. بناء `EnrollmentManager` (token-based enrollment)
5. بناء `DeviceInfoCollector` (platform-specific via method channels)
6. بناء UI shell أساسي

### المرحلة 3: Flutter Agent — Android (2 أسابيع)

1. تحويل وظائف hmdm-android الأساسية:
   - Device Owner provisioning
   - APK silent install
   - Kiosk mode (launcher replacement)
   - System policy enforcement (GPS, BT, WiFi)
2. اختبار التوافق مع QR enrollment الحالي

### المرحلة 4: Flutter Agent — Windows (2-3 أسابيع)

1. MSI/EXE silent install via PowerShell
2. Windows Assigned Access (kiosk)
3. Registry-based policy enforcement
4. Windows service for background sync
5. Hardware info via WMI

### المرحلة 5: Flutter Agent — macOS (2 أسابيع)

1. PKG/DMG install via `installer` command
2. Single App Mode (kiosk)
3. System Preferences enforcement
4. LaunchDaemon for background sync
5. Hardware info via `system_profiler`

### المرحلة 6: Flutter Agent — Linux (2 أسابيع)

1. DEB/AppImage install via `dpkg`/`apt`
2. X11/Wayland kiosk mode
3. systemd service for background sync
4. Hardware info via `/sys/` and `lshw`

---

## 6. ما لا يحتاج تغيير

| المكون | السبب |
|--------|-------|
| Push/Notification system | Polling-based — يعمل لأي HTTP client |
| Multi-tenant architecture | منطق أعمال بحت |
| Profile versioning | JSON artifacts — platform-neutral |
| Device tree + groups | هيكلية تنظيمية |
| Admin web UI (React) | يبقى كما هو — يدير كل المنصات |
| Audit, messaging, devicelog | بيانات نصية |
| File distribution | URL-based download |
| Application settings | Key-value pairs |

---

## 7. تحديات وملاحظات

### 7.1 تحديات تقنية

| التحدي | الحل المقترح |
|--------|-------------|
| Silent app install على desktop | Windows: MSI + PowerShell, macOS: PKG + installer, Linux: apt/dpkg |
| Kiosk mode على desktop | Windows: Assigned Access, macOS: Single App Mode, Linux: custom X session |
| Background service | Windows: Windows Service, macOS: LaunchDaemon, Linux: systemd unit |
| Device Owner equivalent | لا يوجد مكافئ مباشر — نعتمد على صلاحيات admin/root |
| Hardware info | Platform channels في Flutter → native code لكل OS |

### 7.2 قيود Flutter على Desktop

| القيد | التأثير | الحل |
|-------|--------|------|
| لا يوجد Device Owner API | لا يمكن فرض سياسات بنفس قوة Android | نعتمد على OS-level MDM APIs |
| صلاحيات محدودة | لا يمكن تثبيت تطبيقات بصمت بدون admin | نطلب صلاحيات admin عند التثبيت |
| لا يوجد launcher replacement | لا يمكن استبدال shell | نستخدم kiosk mode الخاص بكل OS |
| Platform channels مطلوبة | كود native لكل منصة | Dart + native plugins |

### 7.3 ميزة Flutter

| الميزة | الفائدة |
|--------|--------|
| كود مشترك 60-70% | SyncService, PushPoller, UI, enrollment logic |
| UI موحّد | نفس الواجهة على كل المنصات |
| Hot reload | تطوير سريع |
| Plugin ecosystem | مكتبات جاهزة لـ hardware info, file system, etc. |
| Single codebase | صيانة أسهل |

---

## 8. ملخص التغييرات المطلوبة

### الباك إند (serverBackendGo)

```
التغييرات:
├── DB: إضافة devices.platform, application_packages table
├── sync module: platform-aware response, X-Platform header
├── enrollment: token-based enrollment endpoint (بديل QR)
├── applications: multi-platform package support
├── profiles: platform filter for apps
├── admin UI: platform column in device list
└── deviceinfo plugin: OS-agnostic telemetry schema
```

### Flutter Client (جديد)

```
flutter_mdm_agent/
├── Core (مشترك 70%):
│   ├── Sync protocol client
│   ├── Push polling
│   ├── Enrollment (token-based)
│   ├── App settings sync
│   └── UI shell
├── Platform-specific (30%):
│   ├── Android: Device Owner, APK install, kiosk
│   ├── Windows: MSI install, Assigned Access, WMI
│   ├── macOS: PKG install, Single App Mode
│   └── Linux: DEB install, systemd, X11 kiosk
└── Plugins:
    ├── Device info collector
    ├── Location tracker
    └── Log reporter
```

---

## 9. الخلاصة

**هل يمكن دعم Windows/macOS/Linux؟** — **نعم**، الباك إند مصمم بشكل modular والبروتوكول (REST + JSON + polling) لا يعتمد على Android. التغييرات المطلوبة في الباك إند محدودة (إضافة platform field + enrollment token + multi-platform packages).

**هل Flutter مناسب؟** — **نعم**، Flutter يدعم Android + Windows + macOS + Linux من codebase واحد. 60-70% من كود العميل سيكون مشتركاً (sync, push, UI, enrollment). الـ 30% المتبقي هو platform-specific (تثبيت تطبيقات, kiosk, hardware info) ويُنفّذ عبر platform channels.

**الجهد المقدّر:** 12-16 أسبوع لدعم كامل لجميع المنصات (بما في ذلك تحويل Android من Java إلى Flutter).
