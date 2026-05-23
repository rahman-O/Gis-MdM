import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Textarea } from '@/shared/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

export interface MdmAppOption {
  applicationId: number
  versionId: number
  action: number
  name: string
}

interface ConfigurationMdmTabProps {
  configuration: Configuration
  selectableMdmApps: MdmAppOption[]
  onChange: (configuration: Configuration) => void
}

function toText(value: unknown): string {
  return value == null ? '' : String(value)
}

export function ConfigurationMdmTab({
  configuration,
  selectableMdmApps,
  onChange,
}: ConfigurationMdmTabProps) {
  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  return (
    <div className="grid gap-6 md:grid-cols-2">
      {/* COLUMN 1: Apps & Enrollment */}
      <div className="space-y-6">
        {/* CARD 1: Core Applications */}
        <Card className="shadow-sm border">
          <CardHeader className="bg-muted/15 border-b py-3 px-4">
            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Applications & Provisioning</CardTitle>
            <CardDescription className="text-xs mt-0.5">Primary packages running on the device in MDM mode.</CardDescription>
          </CardHeader>
          <CardContent className="p-4 grid gap-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="mainApp" className="text-xs font-semibold text-muted-foreground uppercase">Main App</Label>
                  <FieldLockToggle
                    fieldKey="mainAppId"
                    locked={isPolicyLocked(configuration, 'mainAppId')}
                    onToggle={setLock}
                  />
                </div>
                <Select
                  disabled={isPolicyLocked(configuration, 'mainAppId')}
                  value={
                    configuration.mainAppId != null && configuration.mainAppId > 0
                      ? String(configuration.mainAppId)
                      : 'none'
                  }
                  onValueChange={(value) =>
                    onChange({ ...configuration, mainAppId: value === 'none' ? null : Number(value) })
                  }
                >
                  <SelectTrigger id="mainApp">
                    <SelectValue placeholder="Select main app" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {selectableMdmApps.map((app) => (
                      <SelectItem key={`m-${app.applicationId}-${app.versionId}`} value={String(app.versionId)}>
                        {app.name || `Application #${app.applicationId}`}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {selectableMdmApps.length === 0 ? (
                  <p className="text-[10px] text-muted-foreground">
                    No applications returned by backend.
                  </p>
                ) : null}
              </div>

              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="contentApp" className="text-xs font-semibold text-muted-foreground uppercase">Content App</Label>
                  <FieldLockToggle
                    fieldKey="contentAppId"
                    locked={isPolicyLocked(configuration, 'contentAppId')}
                    onToggle={setLock}
                  />
                </div>
                <Select
                  disabled={isPolicyLocked(configuration, 'contentAppId')}
                  value={
                    configuration.contentAppId != null && configuration.contentAppId > 0
                      ? String(configuration.contentAppId)
                      : 'none'
                  }
                  onValueChange={(value) =>
                    onChange({ ...configuration, contentAppId: value === 'none' ? null : Number(value) })
                  }
                >
                  <SelectTrigger id="contentApp">
                    <SelectValue placeholder="Select content app" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {selectableMdmApps.map((app) => (
                      <SelectItem key={`c-${app.applicationId}-${app.versionId}`} value={String(app.versionId)}>
                        {app.name || `Application #${app.applicationId}`}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="eventReceivingComponent" className="text-xs font-semibold text-muted-foreground uppercase">Event Receiving Component</Label>
              <Input
                id="eventReceivingComponent"
                placeholder="com.example/.AdminReceiver"
                value={toText(configuration.eventReceivingComponent)}
                onChange={(event) => onChange({ ...configuration, eventReceivingComponent: event.target.value })}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="launcherUrl" className="text-xs font-semibold text-muted-foreground uppercase">Launcher URL Override</Label>
              <Input
                id="launcherUrl"
                placeholder="https://example.com/launcher"
                value={toText(configuration.launcherUrl)}
                onChange={(event) => onChange({ ...configuration, launcherUrl: event.target.value })}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="newServerUrl" className="text-xs font-semibold text-muted-foreground uppercase">New Server URL Override</Label>
              <Input
                id="newServerUrl"
                placeholder="http://server-address:8080"
                value={toText(configuration.newServerUrl)}
                onChange={(event) => onChange({ ...configuration, newServerUrl: event.target.value })}
              />
            </div>
          </CardContent>
        </Card>

        {/* CARD 2: Enrollment Policies & System Toggles */}
        <Card className="shadow-sm border">
          <CardHeader className="bg-muted/15 border-b py-3 px-4">
            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Enrollment & System Modes</CardTitle>
            <CardDescription className="text-xs mt-0.5">Control enrollment security and overall device restrictions.</CardDescription>
          </CardHeader>
          <CardContent className="p-4 grid gap-4">
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="mobileEnrollment"
                  checked={Boolean(configuration.mobileEnrollment)}
                  onCheckedChange={(checked) => onChange({ ...configuration, mobileEnrollment: checked === true })}
                />
                <Label htmlFor="mobileEnrollment" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Mobile enrollment</Label>
              </div>

              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="encryptDevice"
                  checked={Boolean(configuration.encryptDevice)}
                  onCheckedChange={(checked) => onChange({ ...configuration, encryptDevice: checked === true })}
                />
                <Label htmlFor="encryptDevice" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Encrypt device</Label>
              </div>

              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="permissive"
                  checked={Boolean(configuration.permissive)}
                  disabled={Boolean(configuration.kioskMode)}
                  onCheckedChange={(checked) => onChange({ ...configuration, permissive: checked === true })}
                />
                <Label htmlFor="permissive" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Permissive mode</Label>
              </div>

              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="lockSafeSettings"
                  checked={Boolean(configuration.lockSafeSettings)}
                  disabled={Boolean(configuration.permissive)}
                  onCheckedChange={(checked) => onChange({ ...configuration, lockSafeSettings: checked === true })}
                />
                <Label htmlFor="lockSafeSettings" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">Lock safe settings</Label>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="allowedClasses" className="text-xs font-semibold text-muted-foreground uppercase">Allowed Classes (Whitelisting)</Label>
              <Textarea
                id="allowedClasses"
                rows={3}
                placeholder="com.android.settings.Settings"
                value={toText(configuration.allowedClasses)}
                disabled={Boolean(configuration.permissive)}
                onChange={(event) => onChange({ ...configuration, allowedClasses: event.target.value })}
              />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* COLUMN 2: Network & Advanced Customizations */}
      <div className="space-y-6">
        {/* CARD 3: Network Provisioning (Wi-Fi) */}
        <Card className="shadow-sm border">
          <CardHeader className="bg-muted/15 border-b py-3 px-4">
            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Network Provisioning</CardTitle>
            <CardDescription className="text-xs mt-0.5">Preset Wi-Fi network credentials used during device setup.</CardDescription>
          </CardHeader>
          <CardContent className="p-4 grid gap-4">
            <div className="space-y-1.5">
              <Label htmlFor="wifiSSID" className="text-xs font-semibold text-muted-foreground uppercase">Wi-Fi SSID</Label>
              <Input
                id="wifiSSID"
                placeholder="MyCorporateWiFi"
                value={toText(configuration.wifiSSID)}
                onChange={(event) => onChange({ ...configuration, wifiSSID: event.target.value })}
              />
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="wifiPassword" className="text-xs font-semibold text-muted-foreground uppercase">Wi-Fi Password</Label>
                <Input
                  id="wifiPassword"
                  placeholder="Enter Wi-Fi password"
                  value={toText(configuration.wifiPassword)}
                  onChange={(event) => onChange({ ...configuration, wifiPassword: event.target.value })}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wifiSecurityType" className="text-xs font-semibold text-muted-foreground uppercase">Wi-Fi Security Type</Label>
                <Input
                  id="wifiSecurityType"
                  placeholder="WPA/WPA2"
                  value={toText(configuration.wifiSecurityType)}
                  onChange={(event) => onChange({ ...configuration, wifiSecurityType: event.target.value })}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* CARD 4: Advanced Parameters */}
        <Card className="shadow-sm border">
          <CardHeader className="bg-muted/15 border-b py-3 px-4">
            <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Advanced Parameters</CardTitle>
            <CardDescription className="text-xs mt-0.5">Parameters encoded in the provisioning QR code and intent payloads.</CardDescription>
          </CardHeader>
          <CardContent className="p-4 grid gap-4">
            <div className="space-y-1.5">
              <Label htmlFor="qrParameters" className="text-xs font-semibold text-muted-foreground uppercase">QR Code Custom Parameters</Label>
              <Textarea
                id="qrParameters"
                rows={3}
                placeholder="android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED=true"
                value={toText(configuration.qrParameters)}
                onChange={(event) => onChange({ ...configuration, qrParameters: event.target.value })}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="adminExtras" className="text-xs font-semibold text-muted-foreground uppercase">Admin Extras (JSON/Key-Value)</Label>
              <Textarea
                id="adminExtras"
                rows={3}
                placeholder='{"serverUrl": "https://example.com"}'
                value={toText(configuration.adminExtras)}
                onChange={(event) => onChange({ ...configuration, adminExtras: event.target.value })}
              />
            </div>
          </CardContent>
        </Card>

        {/* CARD 5: Kiosk Mode Options (Only if Kiosk is enabled in Common tab) */}
        {Boolean(configuration.kioskMode) ? (
          <Card className="shadow-sm border border-primary/20 bg-primary/5">
            <CardHeader className="border-b border-primary/10 py-3 px-4 bg-primary/10">
              <CardTitle className="text-sm font-bold uppercase tracking-wider text-primary">Kiosk Customizations</CardTitle>
              <CardDescription className="text-xs mt-0.5 text-primary/80">Configure allowed system layout controls when in kiosk mode.</CardDescription>
            </CardHeader>
            <CardContent className="p-4 grid gap-3 sm:grid-cols-2">
              {(
                [
                  ['kioskHome', 'Kiosk: Home Button'],
                  ['kioskRecents', 'Kiosk: Recents Menu'],
                  ['kioskNotifications', 'Kiosk: Notifications Panel'],
                  ['kioskSystemInfo', 'Kiosk: System Info Area'],
                  ['kioskKeyguard', 'Kiosk: Keyguard Lock'],
                  ['kioskLockButtons', 'Kiosk: Lock physical buttons'],
                  ['kioskScreenOn', 'Kiosk: Keep screen alive'],
                  ['kioskExit', 'Kiosk: Allow app exit'],
                ] as const
              ).map(([key, label]) => (
                <div key={key} className="flex items-center space-x-2.5 rounded-md border border-primary/10 bg-background/50 p-2.5 hover:bg-background/80 transition-all duration-200">
                  <Checkbox
                    id={key}
                    checked={Boolean(configuration[key])}
                    onCheckedChange={(checked) => onChange({ ...configuration, [key]: checked === true })}
                  />
                  <Label htmlFor={key} className="cursor-pointer text-xs font-bold uppercase tracking-wider text-foreground flex-1">
                    {label.replace('Kiosk: ', '')}
                  </Label>
                </div>
              ))}
            </CardContent>
          </Card>
        ) : null}
      </div>
    </div>
  )
}
