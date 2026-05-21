# API Contract: Shared push notifier (Phase 9 P0)

**Not a public REST surface** — internal service contract + verification via existing endpoints.

**Java reference**: `com.hmdm.notification.PushService`, `PushMessage` types

## Message types

| Type | When enqueued |
|------|----------------|
| `configUpdated` | Configuration saved; files linked to configuration; bulk notify |
| `appConfigUpdated` | Device application settings notify |

## Internal interface (`platform/push.Notifier`)

```text
NotifyConfigurationChanged(ctx, configurationID int) error
NotifyDeviceApplicationSettings(ctx, deviceID int64) error
```

**Behavior**:

1. Resolve all `devices.id` where `devices.configurationid = configurationID`.
2. For each device: `MessageQueue.Enqueue(ctx, deviceID, messageType, "")`.
3. Errors logged; **return nil** to HTTP save handlers (best-effort, FR-004).

## Verification endpoints (existing)

| Method | Path | Expectation after wiring |
|--------|------|--------------------------|
| PUT | `/rest/private/configurations` | Creates rows in `pushmessages` for affected devices |
| POST | `/rest/private/devices/{id}/applicationSettings/notify` | `appConfigUpdated` message |
| POST | `/rest/private/web-ui-files/configurations` | `configUpdated` for linked configs |
| POST | `/rest/private/push` | Unchanged — manual push API |

## Auth

Notifier runs in authenticated request context; enqueue uses device IDs only (no cross-tenant IDs).

## Parity updates

- `serverBackendGo/docs/parity/configurations.md` — remove NoopPush note
- `serverBackendGo/docs/parity/devices.md` — FCM stub → queue enqueue
- `serverBackendGo/docs/parity/files.md` — push on configuration link
