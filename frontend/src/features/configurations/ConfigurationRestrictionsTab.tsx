import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Textarea } from '@/shared/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

interface ConfigurationRestrictionsTabProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

function toText(value: unknown): string {
  return value == null ? '' : String(value)
}

export function ConfigurationRestrictionsTab({ configuration, onChange }: ConfigurationRestrictionsTabProps) {
  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  const boolField = (fieldKey: keyof Configuration & string, label: string) => (
    <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200 justify-between">
      <div className="flex items-center space-x-2.5 flex-1">
        <Checkbox
          id={fieldKey}
          checked={Boolean(configuration[fieldKey])}
          disabled={isPolicyLocked(configuration, fieldKey)}
          onCheckedChange={(checked) =>
            onChange({ ...configuration, [fieldKey]: checked === true })
          }
        />
        <Label htmlFor={fieldKey} className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">
          {label}
        </Label>
      </div>
      <FieldLockToggle
        fieldKey={fieldKey}
        locked={isPolicyLocked(configuration, fieldKey)}
        onToggle={setLock}
      />
    </div>
  )

  return (
    <Card className="shadow-sm border">
      <CardHeader className="bg-muted/15 border-b py-3 px-4">
        <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Device Restrictions Policy</CardTitle>
        <CardDescription className="text-xs mt-0.5">Enforce standard Android device-level policies via UserManager.</CardDescription>
      </CardHeader>
      <CardContent className="p-4 space-y-4">
        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <Label htmlFor="restrictions" className="text-xs font-semibold text-muted-foreground uppercase">Custom UserManager Restrictions</Label>
            <FieldLockToggle
              fieldKey="restrictions"
              locked={isPolicyLocked(configuration, 'restrictions')}
              onToggle={setLock}
            />
          </div>
          <Textarea
            id="restrictions"
            rows={4}
            placeholder="e.g. no_config_bluetooth, no_safe_boot, no_sms"
            value={toText(configuration.restrictions)}
            disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
            onChange={(event) => onChange({ ...configuration, restrictions: event.target.value })}
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          {boolField('gps', 'GPS')}
          {boolField('bluetooth', 'Bluetooth')}
          {boolField('wifi', 'Wi-Fi')}
          {boolField('mobileData', 'Mobile data')}
          {boolField('usbStorage', 'USB storage')}
          {boolField('disableScreenshots', 'Disable screenshots')}
        </div>
      </CardContent>
    </Card>
  )
}

