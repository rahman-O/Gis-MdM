# Contract: Legacy Java `configurations` → `settingsjson`

**When**: Database restored from Java Liquibase dump with extra columns on `configurations`.  
**When not**: Pure Go migrations only — `000017` is no-op.

## JSON key convention

- camelCase keys matching React configuration editor / Java bean property names.
- Merge: `settingsjson = COALESCE(settingsjson, '{}'::jsonb) || patch` (patch wins on conflict).

## Column → JSON mapping (minimum set for 013)

| Java column (if exists) | JSON key |
|-------------------------|----------|
| gps | gps |
| bluetooth | bluetooth |
| wifi | wifi |
| mobiledata | mobileData |
| kioskmode | kioskMode |
| blockstatusbar | blockStatusBar |
| systemupdatetype | systemUpdateType |
| systemupdatefrom | systemUpdateFrom |
| systemupdateto | systemUpdateTo |
| autoupdate | autoUpdate |
| usbstorage | usbStorage |
| requestupdates | requestUpdates |
| pushoptions | pushOptions |
| autobrightness | autoBrightness |
| brightness | brightness |
| managetimeout | manageTimeout |
| timeout | timeout |
| lockvolume | lockVolume |
| wifissid | wifiSSID |
| wifipassword | wifiPassword |
| wifisecuritytype | wifiSecurityType |
| passwordmode | passwordMode |
| kioskhome | kioskHome |
| kioskrecents | kioskRecents |
| kiosknotifications | kioskNotifications |
| kiosksysteminfo | kioskSystemInfo |
| kioskkeyguard | kioskKeyguard |
| orientation | orientation |
| rundefaultlauncher | runDefaultLauncher |
| timezone | timeZone |
| allowedclasses | allowedClasses |
| newserverurl | newServerUrl |
| locksafesettings | lockSafeSettings |
| disablescreenshots | disableScreenshots |
| restrictions | restrictions |
| keepalivetime | keepaliveTime |
| managevolume | manageVolume |
| volume | volume |
| showwifi | showWifi |
| mobileenrollment | mobileEnrollment |
| disablelocation | disableLocation |
| apppermissions | appPermissions |
| kioskexit | kioskExit |
| qrparameters | qrParameters |
| autostartforeground | autostartForeground |
| displaystatus | displayStatus |
| encryptdevice | encryptDevice |
| downloadupdates | downloadUpdates |
| adminextras | adminExtras |
| eventreceivingcomponent | eventReceivingComponent |
| usedefaultdesignsettings | useDefaultDesignSettings |
| iconsize | iconSize |
| desktopheader | desktopHeader |

Columns already explicit in Go (`backgroundcolor`, `textcolor`, `mainappid`, …) remain in SQL columns; also mirror into JSON if editor reads only `settingsjson`.

## Acceptance

- Sample 3 configurations: ≥90% visible editor toggles match pre-migration Java UI (manual UAT, SC-005).
- `SELECT settingsjson FROM configurations WHERE id IN (...)` non-empty `{}` after import.
