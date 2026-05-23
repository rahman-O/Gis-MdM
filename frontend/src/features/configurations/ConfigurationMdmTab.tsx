import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Textarea } from '@/shared/ui/textarea'

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
    <div className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <Label>Main app</Label>
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
            <SelectTrigger>
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
            <p className="text-xs text-muted-foreground">
              No applications were returned by backend for this customer/session.
            </p>
          ) : null}
        </div>
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <Label>Content app</Label>
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
            <SelectTrigger>
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
      <div className="space-y-2">
        <Label>Event receiving component</Label>
        <Input
          placeholder="com.example/.AdminReceiver"
          value={toText(configuration.eventReceivingComponent)}
          onChange={(event) => onChange({ ...configuration, eventReceivingComponent: event.target.value })}
        />
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <Label>Launcher URL override</Label>
          <Input
            placeholder="https://..."
            value={toText(configuration.launcherUrl)}
            onChange={(event) => onChange({ ...configuration, launcherUrl: event.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label>Wi-Fi SSID</Label>
          <Input
            value={toText(configuration.wifiSSID)}
            onChange={(event) => onChange({ ...configuration, wifiSSID: event.target.value })}
          />
        </div>
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <Label>Wi-Fi password</Label>
          <Input
            value={toText(configuration.wifiPassword)}
            onChange={(event) => onChange({ ...configuration, wifiPassword: event.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label>Wi-Fi security type</Label>
          <Input
            value={toText(configuration.wifiSecurityType)}
            onChange={(event) => onChange({ ...configuration, wifiSecurityType: event.target.value })}
          />
        </div>
      </div>
      <div className="space-y-2">
        <Label>QR parameters</Label>
        <Textarea
          rows={3}
          value={toText(configuration.qrParameters)}
          onChange={(event) => onChange({ ...configuration, qrParameters: event.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Admin extras</Label>
        <Textarea
          rows={3}
          value={toText(configuration.adminExtras)}
          onChange={(event) => onChange({ ...configuration, adminExtras: event.target.value })}
        />
      </div>
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2">
          <Checkbox
            checked={Boolean(configuration.mobileEnrollment)}
            onCheckedChange={(checked) =>
              onChange({ ...configuration, mobileEnrollment: checked === true })
            }
          />
          <Label>Mobile enrollment</Label>
        </div>
        <div className="flex items-center gap-2">
          <Checkbox
            checked={Boolean(configuration.encryptDevice)}
            onCheckedChange={(checked) =>
              onChange({ ...configuration, encryptDevice: checked === true })
            }
          />
          <Label>Encrypt device</Label>
        </div>
      </div>
      <div className="flex flex-wrap items-center gap-6">
        <div className="flex items-center gap-2">
          <Checkbox
            checked={Boolean(configuration.permissive)}
            disabled={Boolean(configuration.kioskMode)}
            onCheckedChange={(checked) => onChange({ ...configuration, permissive: checked === true })}
          />
          <Label>Permissive mode</Label>
        </div>
        <div className="flex items-center gap-2">
          <Checkbox
            checked={Boolean(configuration.lockSafeSettings)}
            disabled={Boolean(configuration.permissive)}
            onCheckedChange={(checked) =>
              onChange({ ...configuration, lockSafeSettings: checked === true })
            }
          />
          <Label>Lock safe settings</Label>
        </div>
        {Boolean(configuration.kioskMode) ? (
          <>
            {(
              [
                ['kioskHome', 'Kiosk: Home'],
                ['kioskRecents', 'Kiosk: Recents'],
                ['kioskNotifications', 'Kiosk: Notifications'],
                ['kioskSystemInfo', 'Kiosk: System info'],
                ['kioskKeyguard', 'Kiosk: Keyguard'],
                ['kioskLockButtons', 'Kiosk: Lock buttons'],
                ['kioskScreenOn', 'Kiosk: Keep screen on'],
                ['kioskExit', 'Kiosk: Allow exit'],
              ] as const
            ).map(([key, label]) => (
              <div key={key} className="flex items-center gap-2">
                <Checkbox
                  checked={Boolean(configuration[key])}
                  onCheckedChange={(checked) =>
                    onChange({ ...configuration, [key]: checked === true })
                  }
                />
                <Label>{label}</Label>
              </div>
            ))}
          </>
        ) : null}
      </div>
      <div className="space-y-2">
        <Label>Allowed classes</Label>
        <Textarea
          rows={3}
          value={toText(configuration.allowedClasses)}
          disabled={Boolean(configuration.permissive)}
          onChange={(event) => onChange({ ...configuration, allowedClasses: event.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>New server URL</Label>
        <Input
          placeholder="http://server:8080"
          value={toText(configuration.newServerUrl)}
          onChange={(event) => onChange({ ...configuration, newServerUrl: event.target.value })}
        />
      </div>
    </div>
  )
}
