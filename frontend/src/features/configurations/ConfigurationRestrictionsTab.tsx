import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Textarea } from '@/shared/ui/textarea'

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

  const boolField = (fieldKey: keyof Configuration, label: string) => (
    <div className="flex items-center gap-2">
      <Checkbox
        checked={Boolean(configuration[fieldKey])}
        disabled={isPolicyLocked(configuration, fieldKey)}
        onCheckedChange={(checked) =>
          onChange({ ...configuration, [fieldKey]: checked === true })
        }
      />
      <Label>{label}</Label>
      <FieldLockToggle
        fieldKey={fieldKey}
        locked={isPolicyLocked(configuration, fieldKey)}
        onToggle={setLock}
      />
    </div>
  )

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <Label>Restrictions (UserManager)</Label>
          <FieldLockToggle
            fieldKey="restrictions"
            locked={isPolicyLocked(configuration, 'restrictions')}
            onToggle={setLock}
          />
        </div>
        <Textarea
          rows={4}
          value={toText(configuration.restrictions)}
          disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
          onChange={(event) => onChange({ ...configuration, restrictions: event.target.value })}
        />
      </div>
      <div className="grid gap-3 sm:grid-cols-2">
        {boolField('gps', 'GPS')}
        {boolField('bluetooth', 'Bluetooth')}
        {boolField('wifi', 'Wi-Fi')}
        {boolField('mobileData', 'Mobile data')}
        {boolField('usbStorage', 'USB storage')}
        {boolField('disableScreenshots', 'Disable screenshots')}
      </div>
    </div>
  )
}
