import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import type { Configuration } from '@/features/configurations/types'

interface AppOption {
  id: number
  name: string
}

interface Props {
  configuration: Configuration
  applications: AppOption[]
  onChange: (next: Configuration) => void
}

export function ConfigurationCommonTab({ configuration, applications, onChange }: Props) {
  const triStateToValue = (v: boolean | null | undefined): 'any' | 'on' | 'off' =>
    v === true ? 'on' : v === false ? 'off' : 'any'
  const valueToTriState = (v: string): boolean | null => (v === 'on' ? true : v === 'off' ? false : null)

  return (
    <div className="grid gap-4 md:grid-cols-2">
      <div className="space-y-2">
        <Label>Name</Label>
        <Input
          value={String(configuration.name ?? '')}
          onChange={(e) => onChange({ ...configuration, name: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Admin password</Label>
        <Input
          type="password"
          value={String(configuration.password ?? '')}
          onChange={(e) => onChange({ ...configuration, password: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Push options</Label>
        <Select
          value={String(configuration.pushOptions ?? 'mqttWorker')}
          onValueChange={(value) => onChange({ ...configuration, pushOptions: value })}
        >
          <SelectTrigger><SelectValue placeholder="Select push mode" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="mqttWorker">MQTT Worker</SelectItem>
            <SelectItem value="mqttAlarm">MQTT Alarm</SelectItem>
            <SelectItem value="polling">Polling</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {String(configuration.pushOptions ?? '') === 'mqttAlarm' ? (
        <div className="space-y-2">
          <Label>Keepalive (seconds)</Label>
          <Select
            value={String(configuration.keepaliveTime ?? 300)}
            onValueChange={(value) => onChange({ ...configuration, keepaliveTime: Number(value) })}
          >
            <SelectTrigger><SelectValue placeholder="Keepalive" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="60">60</SelectItem>
              <SelectItem value="120">120</SelectItem>
              <SelectItem value="180">180</SelectItem>
              <SelectItem value="300">300</SelectItem>
              <SelectItem value="600">600</SelectItem>
              <SelectItem value="900">900</SelectItem>
            </SelectContent>
          </Select>
        </div>
      ) : null}
      <div className="space-y-2">
        <Label>Request updates</Label>
        <Select
          value={String(configuration.requestUpdates ?? 'DONOTTRACK')}
          onValueChange={(value) => onChange({ ...configuration, requestUpdates: value })}
        >
          <SelectTrigger><SelectValue placeholder="Select tracking mode" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="DONOTTRACK">Do Not Track</SelectItem>
            <SelectItem value="GPS">GPS</SelectItem>
            <SelectItem value="WIFI">Wi-Fi</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2 md:col-span-2">
        <Label>Description</Label>
        <Input
          value={String(configuration.description ?? '')}
          onChange={(e) => onChange({ ...configuration, description: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>GPS</Label>
        <Select
          value={triStateToValue(configuration.gps)}
          onValueChange={(value) => onChange({ ...configuration, gps: valueToTriState(value) })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="any">Any</SelectItem>
            <SelectItem value="off">Off</SelectItem>
            <SelectItem value="on">On</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Bluetooth</Label>
        <Select
          value={triStateToValue(configuration.bluetooth)}
          onValueChange={(value) => onChange({ ...configuration, bluetooth: valueToTriState(value) })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="any">Any</SelectItem>
            <SelectItem value="off">Off</SelectItem>
            <SelectItem value="on">On</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Wi-Fi</Label>
        <Select
          value={triStateToValue(configuration.wifi)}
          onValueChange={(value) => onChange({ ...configuration, wifi: valueToTriState(value) })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="any">Any</SelectItem>
            <SelectItem value="off">Off</SelectItem>
            <SelectItem value="on">On</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Mobile data</Label>
        <Select
          value={triStateToValue(configuration.mobileData)}
          onValueChange={(value) => onChange({ ...configuration, mobileData: valueToTriState(value) })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="any">Any</SelectItem>
            <SelectItem value="off">Off</SelectItem>
            <SelectItem value="on">On</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Password mode</Label>
        <Select
          value={String(configuration.passwordMode ?? 'any')}
          onValueChange={(value) => onChange({ ...configuration, passwordMode: value === 'any' ? null : value })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="any">Any</SelectItem>
            <SelectItem value="present">Present</SelectItem>
            <SelectItem value="easy">Easy</SelectItem>
            <SelectItem value="moderate">Moderate</SelectItem>
            <SelectItem value="strong">Strong</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Download updates</Label>
        <Select
          value={String(configuration.downloadUpdates ?? 'UNLIMITED')}
          onValueChange={(value) => onChange({ ...configuration, downloadUpdates: value })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="UNLIMITED">Unlimited</SelectItem>
            <SelectItem value="LIMITED">Limited</SelectItem>
            <SelectItem value="WIFI">Wi-Fi only</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.manageTimeout)}
          onCheckedChange={(checked) => onChange({ ...configuration, manageTimeout: checked === true })}
        />
        <Label>Manage screen timeout</Label>
      </div>
      <div className="space-y-2">
        <Label>Timeout (seconds)</Label>
        <Input
          type="number"
          value={String(configuration.timeout ?? 60)}
          onChange={(e) => onChange({ ...configuration, timeout: Number(e.target.value) || 0 })}
          disabled={!Boolean(configuration.manageTimeout)}
        />
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.usbStorage)}
          onCheckedChange={(checked) => onChange({ ...configuration, usbStorage: checked === true })}
        />
        <Label>Allow USB storage</Label>
      </div>
      <div className="space-y-2">
        <Label>Brightness mode</Label>
        <Select
          value={
            configuration.autoBrightness === true
              ? 'auto'
              : configuration.autoBrightness === false
                ? 'manual'
                : 'none'
          }
          onValueChange={(value) =>
            onChange({
              ...configuration,
              autoBrightness: value === 'auto' ? true : value === 'manual' ? false : null,
            })
          }
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="none">Do not manage</SelectItem>
            <SelectItem value="manual">Manual</SelectItem>
            <SelectItem value="auto">Auto</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {configuration.autoBrightness === false ? (
        <div className="space-y-2">
          <Label>Brightness (0-255)</Label>
          <Input
            type="number"
            min={0}
            max={255}
            value={String(configuration.brightness ?? 180)}
            onChange={(e) => {
              const n = Number(e.target.value)
              onChange({ ...configuration, brightness: Number.isFinite(n) ? Math.min(255, Math.max(0, n)) : 0 })
            }}
          />
        </div>
      ) : null}
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.lockVolume)}
          onCheckedChange={(checked) => onChange({ ...configuration, lockVolume: checked === true })}
        />
        <Label>Lock volume keys</Label>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.manageVolume)}
          onCheckedChange={(checked) => onChange({ ...configuration, manageVolume: checked === true })}
        />
        <Label>Manage volume</Label>
      </div>
      {Boolean(configuration.manageVolume) ? (
        <div className="space-y-2">
          <Label>Volume (0-100)</Label>
          <Input
            type="number"
            min={0}
            max={100}
            value={String(configuration.volume ?? 50)}
            onChange={(e) => {
              const n = Number(e.target.value)
              onChange({ ...configuration, volume: Number.isFinite(n) ? Math.min(100, Math.max(0, n)) : 0 })
            }}
          />
        </div>
      ) : null}
      <div className="space-y-2">
        <Label>Time zone mode</Label>
        <Select
          value={String(configuration.timeZoneMode ?? 'default')}
          onValueChange={(value) =>
            onChange({ ...configuration, timeZoneMode: value })
          }
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="default">Default</SelectItem>
            <SelectItem value="auto">Auto</SelectItem>
            <SelectItem value="manual">Manual</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {String(configuration.timeZoneMode ?? 'default') === 'manual' ? (
        <div className="space-y-2">
          <Label>Time zone</Label>
          <Input
            placeholder="Europe/Moscow"
            value={String(configuration.timeZone ?? '')}
            onChange={(e) => onChange({ ...configuration, timeZone: e.target.value })}
          />
        </div>
      ) : null}
      <div className="space-y-2">
        <Label>System update policy</Label>
        <Select
          value={String(configuration.systemUpdateType ?? 0)}
          onValueChange={(value) => onChange({ ...configuration, systemUpdateType: Number(value) })}
        >
          <SelectTrigger><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="0">Default</SelectItem>
            <SelectItem value="1">Immediate</SelectItem>
            <SelectItem value="2">Scheduled</SelectItem>
            <SelectItem value="3">Postponed</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {Number(configuration.systemUpdateType ?? 0) === 2 ? (
        <>
          <div className="space-y-2">
            <Label>System update from</Label>
            <Input
              type="time"
              value={String(configuration.systemUpdateFrom ?? '')}
              onChange={(e) => onChange({ ...configuration, systemUpdateFrom: e.target.value })}
            />
          </div>
          <div className="space-y-2">
            <Label>System update to</Label>
            <Input
              type="time"
              value={String(configuration.systemUpdateTo ?? '')}
              onChange={(e) => onChange({ ...configuration, systemUpdateTo: e.target.value })}
            />
          </div>
        </>
      ) : null}
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.scheduleAppUpdate)}
          onCheckedChange={(checked) => onChange({ ...configuration, scheduleAppUpdate: checked === true })}
        />
        <Label>Schedule app updates</Label>
      </div>
      {Boolean(configuration.scheduleAppUpdate) ? (
        <>
          <div className="space-y-2">
            <Label>App updates from</Label>
            <Input
              type="time"
              value={String(configuration.appUpdateFrom ?? '')}
              onChange={(e) => onChange({ ...configuration, appUpdateFrom: e.target.value })}
            />
          </div>
          <div className="space-y-2">
            <Label>App updates to</Label>
            <Input
              type="time"
              value={String(configuration.appUpdateTo ?? '')}
              onChange={(e) => onChange({ ...configuration, appUpdateTo: e.target.value })}
            />
          </div>
        </>
      ) : null}
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.runDefaultLauncher)}
          onCheckedChange={(checked) => onChange({ ...configuration, runDefaultLauncher: checked === true })}
        />
        <Label>Use default launcher</Label>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.showWifi)}
          onCheckedChange={(checked) => onChange({ ...configuration, showWifi: checked === true })}
        />
        <Label>Show Wi-Fi icon</Label>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.disableScreenshots)}
          onCheckedChange={(checked) => onChange({ ...configuration, disableScreenshots: checked === true })}
        />
        <Label>Disable screenshots</Label>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.autostartForeground)}
          onCheckedChange={(checked) => onChange({ ...configuration, autostartForeground: checked === true })}
        />
        <Label>Autostart in foreground</Label>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.kioskMode)}
          onCheckedChange={(checked) => onChange({ ...configuration, kioskMode: checked === true })}
        />
        <Label>Kiosk mode</Label>
      </div>
      <div className="space-y-2">
        <Label>Content app (required in kiosk mode)</Label>
        <Select
          value={
            configuration.contentAppId != null && configuration.contentAppId > 0
              ? String(configuration.contentAppId)
              : 'none'
          }
          onValueChange={(value) =>
            onChange({
              ...configuration,
              contentAppId: value === 'none' ? null : Number(value),
            })
          }
        >
          <SelectTrigger><SelectValue placeholder="Select content app" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="none">None</SelectItem>
            {applications.map((app) => (
              <SelectItem key={app.id} value={String(app.id)}>
                {app.name || `Application #${app.id}`}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
