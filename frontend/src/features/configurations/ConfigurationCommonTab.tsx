import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/shared/ui/card'
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
    <div className="grid gap-6 md:grid-cols-2">
      {/* CARD 1: Profile & Sync Settings */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Profile & Connectivity</CardTitle>
          <CardDescription className="text-xs mt-0.5">Basic identifiers and push synchronization details.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="name" className="text-xs font-semibold text-muted-foreground uppercase">Configuration Name</Label>
              <Input
                id="name"
                value={String(configuration.name ?? '')}
                placeholder="Enter config name"
                onChange={(e) => onChange({ ...configuration, name: e.target.value })}
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="password" className="text-xs font-semibold text-muted-foreground uppercase">Admin Password</Label>
              <Input
                id="password"
                type="password"
                value={String(configuration.password ?? '')}
                placeholder="••••••••"
                onChange={(e) => onChange({ ...configuration, password: e.target.value })}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="description" className="text-xs font-semibold text-muted-foreground uppercase">Description</Label>
            <Input
              id="description"
              value={String(configuration.description ?? '')}
              placeholder="Provide a description for this configuration profile"
              onChange={(e) => onChange({ ...configuration, description: e.target.value })}
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="pushOptions" className="text-xs font-semibold text-muted-foreground uppercase">Push Connection Mode</Label>
              <Select
                value={String(configuration.pushOptions ?? 'mqttWorker')}
                onValueChange={(value) => onChange({ ...configuration, pushOptions: value })}
              >
                <SelectTrigger id="pushOptions"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="mqttWorker">MQTT Worker</SelectItem>
                  <SelectItem value="mqttAlarm">MQTT Alarm</SelectItem>
                  <SelectItem value="polling">Polling</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {String(configuration.pushOptions ?? '') === 'mqttAlarm' ? (
              <div className="space-y-1.5">
                <Label htmlFor="keepaliveTime" className="text-xs font-semibold text-muted-foreground uppercase">MQTT Keepalive (seconds)</Label>
                <Select
                  value={String(configuration.keepaliveTime ?? 300)}
                  onValueChange={(value) => onChange({ ...configuration, keepaliveTime: Number(value) })}
                >
                  <SelectTrigger id="keepaliveTime"><SelectValue /></SelectTrigger>
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

            <div className="space-y-1.5">
              <Label htmlFor="requestUpdates" className="text-xs font-semibold text-muted-foreground uppercase">Request Location Updates</Label>
              <Select
                value={String(configuration.requestUpdates ?? 'DONOTTRACK')}
                onValueChange={(value) => onChange({ ...configuration, requestUpdates: value })}
              >
                <SelectTrigger id="requestUpdates"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="DONOTTRACK">Do Not Track</SelectItem>
                  <SelectItem value="GPS">GPS</SelectItem>
                  <SelectItem value="WIFI">Wi-Fi</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* CARD 2: Device Hardware Toggles */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Device Hardware & Connectivity</CardTitle>
          <CardDescription className="text-xs mt-0.5">Toggle hardware state restrictions and mobile behaviors.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="gps" className="text-xs font-semibold text-muted-foreground uppercase">GPS State</Label>
              <Select
                value={triStateToValue(configuration.gps)}
                onValueChange={(value) => onChange({ ...configuration, gps: valueToTriState(value) })}
              >
                <SelectTrigger id="gps"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Any (Do not manage)</SelectItem>
                  <SelectItem value="off">Enforce Off</SelectItem>
                  <SelectItem value="on">Enforce On</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="bluetooth" className="text-xs font-semibold text-muted-foreground uppercase">Bluetooth State</Label>
              <Select
                value={triStateToValue(configuration.bluetooth)}
                onValueChange={(value) => onChange({ ...configuration, bluetooth: valueToTriState(value) })}
              >
                <SelectTrigger id="bluetooth"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Any (Do not manage)</SelectItem>
                  <SelectItem value="off">Enforce Off</SelectItem>
                  <SelectItem value="on">Enforce On</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="wifi" className="text-xs font-semibold text-muted-foreground uppercase">Wi-Fi Enforcer</Label>
              <Select
                value={triStateToValue(configuration.wifi)}
                onValueChange={(value) => onChange({ ...configuration, wifi: valueToTriState(value) })}
              >
                <SelectTrigger id="wifi"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Any (Do not manage)</SelectItem>
                  <SelectItem value="off">Enforce Off</SelectItem>
                  <SelectItem value="on">Enforce On</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="mobileData" className="text-xs font-semibold text-muted-foreground uppercase">Mobile Data State</Label>
              <Select
                value={triStateToValue(configuration.mobileData)}
                onValueChange={(value) => onChange({ ...configuration, mobileData: valueToTriState(value) })}
              >
                <SelectTrigger id="mobileData"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Any (Do not manage)</SelectItem>
                  <SelectItem value="off">Enforce Off</SelectItem>
                  <SelectItem value="on">Enforce On</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="downloadUpdates" className="text-xs font-semibold text-muted-foreground uppercase">Download App Updates</Label>
              <Select
                value={String(configuration.downloadUpdates ?? 'UNLIMITED')}
                onValueChange={(value) => onChange({ ...configuration, downloadUpdates: value })}
              >
                <SelectTrigger id="downloadUpdates"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="UNLIMITED">Unlimited Networks</SelectItem>
                  <SelectItem value="LIMITED">Limited Networks</SelectItem>
                  <SelectItem value="WIFI">Wi-Fi Only</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="passwordMode" className="text-xs font-semibold text-muted-foreground uppercase">Screen Password Mode</Label>
              <Select
                value={String(configuration.passwordMode ?? 'any')}
                onValueChange={(value) => onChange({ ...configuration, passwordMode: value === 'any' ? null : value })}
              >
                <SelectTrigger id="passwordMode"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="any">Any / No Enforce</SelectItem>
                  <SelectItem value="present">Enforce Present</SelectItem>
                  <SelectItem value="easy">Easy (Alphanumeric)</SelectItem>
                  <SelectItem value="moderate">Moderate (Numbers/Letters)</SelectItem>
                  <SelectItem value="strong">Strong (Complex)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid gap-3 sm:grid-cols-2 mt-2">
            <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="usbStorage"
                checked={Boolean(configuration.usbStorage)}
                onCheckedChange={(checked) => onChange({ ...configuration, usbStorage: checked === true })}
              />
              <Label htmlFor="usbStorage" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Allow USB storage</Label>
            </div>

            <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="lockVolume"
                checked={Boolean(configuration.lockVolume)}
                onCheckedChange={(checked) => onChange({ ...configuration, lockVolume: checked === true })}
              />
              <Label htmlFor="lockVolume" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Lock volume keys</Label>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* CARD 3: Display, Sound & Timezone Controls */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Display, Audio & Timezones</CardTitle>
          <CardDescription className="text-xs mt-0.5">Control system values like timeout, brightness, volume, and clock settings.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            {/* Screen Timeout Block */}
            <div className="space-y-3">
              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="manageTimeout"
                  checked={Boolean(configuration.manageTimeout)}
                  onCheckedChange={(checked) => onChange({ ...configuration, manageTimeout: checked === true })}
                />
                <Label htmlFor="manageTimeout" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Manage screen timeout</Label>
              </div>
              {Boolean(configuration.manageTimeout) && (
                <div className="pl-2 space-y-1">
                  <Label htmlFor="timeout" className="text-xs font-semibold text-muted-foreground uppercase">Timeout (seconds)</Label>
                  <Input
                    id="timeout"
                    type="number"
                    value={String(configuration.timeout ?? 60)}
                    onChange={(e) => onChange({ ...configuration, timeout: Number(e.target.value) || 0 })}
                  />
                </div>
              )}
            </div>

            {/* Volume Control Block */}
            <div className="space-y-3">
              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="manageVolume"
                  checked={Boolean(configuration.manageVolume)}
                  onCheckedChange={(checked) => onChange({ ...configuration, manageVolume: checked === true })}
                />
                <Label htmlFor="manageVolume" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Manage device volume</Label>
              </div>
              {Boolean(configuration.manageVolume) && (
                <div className="pl-2 space-y-1">
                  <Label htmlFor="volume" className="text-xs font-semibold text-muted-foreground uppercase">Volume level (0-100)</Label>
                  <Input
                    id="volume"
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
              )}
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2 border-t pt-4">
            {/* Brightness Mode Block */}
            <div className="space-y-2">
              <Label htmlFor="brightnessMode" className="text-xs font-semibold text-muted-foreground uppercase">Brightness Mode</Label>
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
                <SelectTrigger id="brightnessMode"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Do not manage</SelectItem>
                  <SelectItem value="manual">Manual Level</SelectItem>
                  <SelectItem value="auto">Automatic (Ambient)</SelectItem>
                </SelectContent>
              </Select>
              {configuration.autoBrightness === false && (
                <div className="space-y-1 mt-2">
                  <Label htmlFor="brightness" className="text-xs font-semibold text-muted-foreground uppercase">Brightness (0-255)</Label>
                  <Input
                    id="brightness"
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
              )}
            </div>

            {/* Timezone Mode Block */}
            <div className="space-y-2">
              <Label htmlFor="timezoneMode" className="text-xs font-semibold text-muted-foreground uppercase">Time Zone Mode</Label>
              <Select
                value={String(configuration.timeZoneMode ?? 'default')}
                onValueChange={(value) => onChange({ ...configuration, timeZoneMode: value })}
              >
                <SelectTrigger id="timezoneMode"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="default">Default Behavior</SelectItem>
                  <SelectItem value="auto">Auto Sync Network</SelectItem>
                  <SelectItem value="manual">Manual Timezone</SelectItem>
                </SelectContent>
              </Select>
              {String(configuration.timeZoneMode ?? 'default') === 'manual' && (
                <div className="space-y-1 mt-2">
                  <Label htmlFor="timezone" className="text-xs font-semibold text-muted-foreground uppercase">Manual Time Zone</Label>
                  <Input
                    id="timezone"
                    placeholder="Europe/Moscow"
                    value={String(configuration.timeZone ?? '')}
                    onChange={(e) => onChange({ ...configuration, timeZone: e.target.value })}
                  />
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* CARD 4: System Policies & OS Updates */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">System Policies & Kiosk</CardTitle>
          <CardDescription className="text-xs mt-0.5">Control Android update rules, system layout switches, and kiosk behaviors.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            {/* System Update Policy */}
            <div className="space-y-2">
              <Label htmlFor="systemUpdateType" className="text-xs font-semibold text-muted-foreground uppercase">System OTA Updates Policy</Label>
              <Select
                value={String(configuration.systemUpdateType ?? 0)}
                onValueChange={(value) => onChange({ ...configuration, systemUpdateType: Number(value) })}
              >
                <SelectTrigger id="systemUpdateType"><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">Default (System Controlled)</SelectItem>
                  <SelectItem value="1">Immediate (As available)</SelectItem>
                  <SelectItem value="2">Scheduled Window</SelectItem>
                  <SelectItem value="3">Postponed (Up to 30 days)</SelectItem>
                </SelectContent>
              </Select>
              {Number(configuration.systemUpdateType ?? 0) === 2 && (
                <div className="grid grid-cols-2 gap-2 mt-2">
                  <div className="space-y-1">
                    <Label htmlFor="systemUpdateFrom" className="text-[10px] font-semibold text-muted-foreground uppercase">From</Label>
                    <Input
                      id="systemUpdateFrom"
                      type="time"
                      value={String(configuration.systemUpdateFrom ?? '')}
                      onChange={(e) => onChange({ ...configuration, systemUpdateFrom: e.target.value })}
                    />
                  </div>
                  <div className="space-y-1">
                    <Label htmlFor="systemUpdateTo" className="text-[10px] font-semibold text-muted-foreground uppercase">To</Label>
                    <Input
                      id="systemUpdateTo"
                      type="time"
                      value={String(configuration.systemUpdateTo ?? '')}
                      onChange={(e) => onChange({ ...configuration, systemUpdateTo: e.target.value })}
                    />
                  </div>
                </div>
              )}
            </div>

            {/* App updates block */}
            <div className="space-y-2">
              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="scheduleAppUpdate"
                  checked={Boolean(configuration.scheduleAppUpdate)}
                  onCheckedChange={(checked) => onChange({ ...configuration, scheduleAppUpdate: checked === true })}
                />
                <Label htmlFor="scheduleAppUpdate" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Schedule App Updates</Label>
              </div>
              {Boolean(configuration.scheduleAppUpdate) && (
                <div className="grid grid-cols-2 gap-2 mt-2 pl-2">
                  <div className="space-y-1">
                    <Label htmlFor="appUpdateFrom" className="text-[10px] font-semibold text-muted-foreground uppercase">From</Label>
                    <Input
                      id="appUpdateFrom"
                      type="time"
                      value={String(configuration.appUpdateFrom ?? '')}
                      onChange={(e) => onChange({ ...configuration, appUpdateFrom: e.target.value })}
                    />
                  </div>
                  <div className="space-y-1">
                    <Label htmlFor="appUpdateTo" className="text-[10px] font-semibold text-muted-foreground uppercase">To</Label>
                    <Input
                      id="appUpdateTo"
                      type="time"
                      value={String(configuration.appUpdateTo ?? '')}
                      onChange={(e) => onChange({ ...configuration, appUpdateTo: e.target.value })}
                    />
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="grid gap-3 sm:grid-cols-2 border-t pt-4">
            <div className="flex items-center space-x-2.5 rounded-md border p-2.5 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="runDefaultLauncher"
                checked={Boolean(configuration.runDefaultLauncher)}
                onCheckedChange={(checked) => onChange({ ...configuration, runDefaultLauncher: checked === true })}
              />
              <Label htmlFor="runDefaultLauncher" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Use default launcher</Label>
            </div>

            <div className="flex items-center space-x-2.5 rounded-md border p-2.5 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="showWifi"
                checked={Boolean(configuration.showWifi)}
                onCheckedChange={(checked) => onChange({ ...configuration, showWifi: checked === true })}
              />
              <Label htmlFor="showWifi" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Show Wi-Fi Icon</Label>
            </div>

            <div className="flex items-center space-x-2.5 rounded-md border p-2.5 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="disableScreenshots"
                checked={Boolean(configuration.disableScreenshots)}
                onCheckedChange={(checked) => onChange({ ...configuration, disableScreenshots: checked === true })}
              />
              <Label htmlFor="disableScreenshots" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Disable screenshots</Label>
            </div>

            <div className="flex items-center space-x-2.5 rounded-md border p-2.5 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="autostartForeground"
                checked={Boolean(configuration.autostartForeground)}
                onCheckedChange={(checked) => onChange({ ...configuration, autostartForeground: checked === true })}
              />
              <Label htmlFor="autostartForeground" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Autostart in foreground</Label>
            </div>
          </div>

          <div className="space-y-3 border-t pt-4">
            <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
              <Checkbox
                id="kioskMode"
                checked={Boolean(configuration.kioskMode)}
                onCheckedChange={(checked) => onChange({ ...configuration, kioskMode: checked === true })}
              />
              <Label htmlFor="kioskMode" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Enforce Kiosk Mode</Label>
            </div>
            {Boolean(configuration.kioskMode) && (
              <div className="pl-2 space-y-1.5">
                <Label htmlFor="contentAppId" className="text-xs font-semibold text-muted-foreground uppercase">Content App (Kiosk Target)</Label>
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
                  <SelectTrigger id="contentAppId"><SelectValue placeholder="Select app..." /></SelectTrigger>
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
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
