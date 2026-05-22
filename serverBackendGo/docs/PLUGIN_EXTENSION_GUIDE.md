# دليل توسيع Plugins في Go — هيكلية وسياق موحّد

**الهدف:** إضافة plugin جديد (مدمج أو منفصل من [h-mdm](https://github.com/h-mdm?tab=repositories)) بأقل احتكاك، مع الحفاظ على parity مع Java وقواعد [`constitution.md`](../../.specify/memory/constitution.md).

**الوضع الحالي:** Plugins ليست JARs ديناميكية — كل plugin = **وحدة Go مُجمَّعة** + صف في `plugins` + flags في `.env`.

---

## 1. المبدأ: «عقد واحد، مجلد واحد، تسجيل واحد»

```text
internal/modules/plugins/<identifier>/
├── module.go              # PluginModule: Identifier + Register
├── domain/
├── port/                  # اختياري إن كان persistence معقّداً
├── application/
└── adapter/
    ├── http/              # handler + routes تحت /rest/plugins/<identifier>/...
    └── persistence/postgres/

db/migrations/0000NN_<identifier>_plugin.up.sql   # جداول + seed plugins + permissions
docs/parity/plugins-<identifier>.md
specs/0xx-plugin-<identifier>/                      # spec + contracts (اختياري للكبير)
```

**لا تُنشئ** مسارات خارج `/rest/plugins/...` أو `/rest/plugin/main/...` إلا إذا Java يفعل ذلك حرفياً (مثلاً `deviceinfo` public على `engine` مباشرة).

---

## 2. واجهة PluginModule (مقترح — طبقة رفيعة فوق `module.Module`)

اليوم كل plugin يطبّق `module.Module` فقط. للتوسيع السهل، عرّف واجهة اختيارية في `internal/platform/plugin/`:

```go
// internal/platform/plugin/plugin.go
package plugin

import "github.com/gis-mdm/server-backend-go/internal/module"

// Module is a Headwind MDM plugin (compile-time, not dynamic JAR).
type Module interface {
    module.Module
    // Identifier matches plugins.identifier and ENABLED_PLUGINS (lowercase).
    Identifier() string
    // Permissions seeded in migration (for docs / validation).
    PermissionNames() []string
}

// SyncContributor optional — merged into agent sync response.
type SyncContributor interface {
    ContributeSync(ctx context.Context, deviceID int64) (json.RawMessage, error)
}

// CustomerBootstrap optional — run after new tenant created.
type CustomerBootstrap interface {
    OnCustomerCreated(ctx context.Context, customerID int64) error
}
```

**التسجيل المركزي** في `internal/app/plugins_registry.go`:

```go
func pluginModules(platformCache *status.Cache) []module.Module {
    return []module.Module{
        pluginplatform.New(platformCache),
        pluginaudit.New(),
        pluginpush.New(),
        // ...
        pluginwifimanager.New(), // عند الإضافة
    }
}
```

`modules.go` يستدعي `registerModules(pluginModules(cache)...)` بدلاً من قائمة طويلة مكررة.

---

## 3. خطوات إضافة plugin جديد (checklist)

### المرحلة A — تحليل Java (مصدر الحقيقة)

1. استنساخ/قراءة مستودع h-mdm (مثلاً `hmdm-plugin-wifimanager`) أو `backend/plugins/<name>/`.
2. استخراج:
   - `*Resource.java` → جدول Method | Path | Auth | Permission
   - Liquibase → جداول `plugin_*`
   - `permissions` في SQL
   - هل يوجد public routes للوكيل؟
   - هل يشارك في `SyncResponse`؟

### المرحلة B — قاعدة البيانات

```sql
-- migration 0000NN_<identifier>_plugin.up.sql
CREATE TABLE IF NOT EXISTS plugin_<identifier>_... (...);

INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice, ...)
SELECT '<identifier>', '...', '...', 'plugin.<identifier>.localization.key.name', TRUE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = '<identifier>');

INSERT INTO permissions (name, description, superadmin) ...
```

### المرحلة C — كود Go

| ملف | المسؤولية |
|-----|-----------|
| `module.go` | `Identifier()`, flags `MODULE_PLUGINS_<ID>_ENABLED`, `Register` |
| `adapter/http/handler.go` | Gin handlers + `httpx.Envelope` |
| `application/service.go` | منطق الأعمال |
| `adapter/persistence/postgres/` | SQL |

**قالب `module.go`:**

```go
func (m *Module) Identifier() string { return "wifimanager" }

func (m *Module) Register(groups module.RouteGroups, deps module.Dependencies) error {
    if !deps.Config.ModulePluginsEnabled || !deps.Config.IsPluginEnabled(m.Identifier()) {
        deps.Log.Info("module disabled", "module", m.Name())
        return nil
    }
    if deps.DB == nil {
        return fmt.Errorf("plugins/%s requires DATABASE_URL", m.Identifier())
    }
    // repo → svc → handler
    base := groups.Plugins.Group("/" + m.Identifier()) // أو المسار الدقيق من Java
    h.Register(base, groups.Engine) // Engine فقط إن وُجد public path خاص
    return nil
}
```

### المرحلة D — تكوين

```env
ENABLED_PLUGINS=audit,push,messaging,deviceinfo,devicelog,wifimanager
MODULE_PLUGINS_ENABLED=true
MODULE_PLUGINS_WIFIMANAGER_ENABLED=true   # أضف في config.go
```

### المرحلة E — منصة + React

1. `docs/parity/plugins-<identifier>.md`
2. صفحة React تحت `frontend/src/features/plugins/<identifier>/` (لا Angular templates)
3. مفاتيح i18n مطابقة لـ `namelocalizationkey` إن لزم

---

## 4. طبقة platform المشتركة (أغلقها قبل plugins كثيرة)

| المكوّن | المسار المقترح | يعادل Java |
|---------|----------------|------------|
| **Plugin access middleware** | `internal/platform/plugin/middleware/access.go` | `PluginAccessFilter` |
| **Sync hook registry** | `internal/platform/plugin/sync/registry.go` | `SyncResponseHook` |
| **Audit recorder** | `internal/platform/audit/` | `AuditFilter` |
| **Customer bootstrap bus** | `internal/platform/plugin/lifecycle/bus.go` | `CustomerCreatedEventListener` |

### 4.1 PluginAccess middleware (مثال سلوك)

```text
على groups.Plugins:
  1. استخرج identifier من المسار (/rest/plugins/{id}/...)
  2. إن !IsPluginEnabled(env) → 404
  3. إن principal.customer + pluginsdisabled → 403
  4. c.Next()
```

يُسجَّل مرة واحدة في `httpx/router.go` بعد `RequireAuth`.

### 4.2 Sync registry

```go
type Registry struct { contributors []SyncContributor }
func (r *Registry) Apply(deviceID int64, resp *sync.Response) { ... }
```

`sync` module يستدعي `registry.Apply` قبل إرجاع JSON — كل plugin جديد يُسجّل نفسه في `app/wiring` دون تعديل `sync/service.go` لكل مرة.

---

## 5. هيكلية المجلدات الموصى بها (نهائي)

```text
serverBackendGo/
├── internal/
│   ├── app/
│   │   ├── modules.go           # core modules فقط
│   │   ├── plugins_registry.go  # كل plugins/*
│   │   └── wiring.go            # SyncRegistry, Audit, PushNotifier
│   ├── platform/
│   │   ├── plugin/
│   │   │   ├── plugin.go        # واجهات Module, SyncContributor, ...
│   │   │   ├── middleware/
│   │   │   ├── sync/
│   │   │   └── lifecycle/
│   │   ├── audit/
│   │   └── push/                # موجود
│   └── modules/
│       └── plugins/
│           ├── platform/        # catalog فقط
│           ├── shared/          # status.Cache, targets
│           ├── audit/
│           ├── push/
│           ├── messaging/
│           ├── deviceinfo/
│           ├── devicelog/
│           └── <new>/           # wifimanager, apn, ...
├── db/migrations/
│   └── 0000NN_<plugin>_plugin.up.sql
└── docs/
    ├── PLUGIN_EXTENSION_GUIDE.md  # هذا الملف
    └── parity/plugins-*.md
```

---

## 6. سياق Spec Kit (موصى به لكل plugin كبير)

```text
specs/015-plugin-wifimanager/
├── spec.md           # قصص المستخدم + FR
├── contracts/api.md  # جدول REST من Java
├── data-model.md     # جداول plugin_*
├── plan.md
└── tasks.md
```

يربط العمل بـ Speckit ويمنع نسيان permissions أو public routes.

---

## 7. ما لا تفعله

| خطأ | البديل |
|-----|--------|
| تحميل `.so` / plugin ديناميكي | compile-time module + `ENABLED_PLUGINS` |
| نسخ handler من Java حرفياً في HTTP | طبقة `application` + `port` |
| مسارات جديدة للواجهة | نفس `/rest/plugins/...` كما في Java |
| UI Angular (`javascriptmodulefile`) | صفحة React جديدة |
| plugin بدون صف في `plugins` | migration seed إلزامي |
| تعديل `sync/service.go` لكل plugin | `SyncContributor` registry |

---

## 8. ترتيب التنفيذ الموصى به (للمشروع ككل)

```text
1. platform/plugin/middleware (PluginAccess)     ← يحمي كل plugins الحالية والجديدة
2. platform/plugin/sync (Registry)               ← يسهل deviceinfo وغيره على الوكيل
3. platform/audit middleware                     ← امتثال
4. إكمال deviceinfo/devicelog endpoints الناقصة
5. plugins_registry.go + قالب scaffold
6. أول plugin خارجي pilot (مثلاً wifimanager) ب spec 015
```

---

## 9. مقارنة: Java منفصل vs Go

| Java (h-mdm) | Go (Gis-MdM) |
|--------------|--------------|
| Maven module → WAR | `internal/modules/plugins/<id>/` |
| Guice `PluginModule` | `module.go` + `Register` |
| `PluginAccessFilter` | `platform/plugin/middleware` |
| `SyncResponseHook` multibind | `SyncContributor` registry |
| Liquibase في plugin | `db/migrations/0000NN_*.sql` |
| Angular `*.module.js` | React feature folder |

---

## 10. مراجع

- [`MIGRATION.md`](MIGRATION.md) Phase 8
- [`parity/plugins-platform.md`](parity/plugins-platform.md)
- [`H-MDM-GITHUB-TO-GO-BACKEND-ANALYSIS.md`](../../H-MDM-GITHUB-TO-GO-BACKEND-ANALYSIS.md)
- [`JAVA-GO-BACKEND-GAPS.md`](../../JAVA-GO-BACKEND-GAPS.md) § plugins
- [h-mdm repositories](https://github.com/h-mdm?tab=repositories)

---

*عند إضافة أول plugin خارجي، أنشئ `docs/parity/plugins-<identifier>.md` وحدّث `ENABLED_PLUGINS` في `.env.example`.*
