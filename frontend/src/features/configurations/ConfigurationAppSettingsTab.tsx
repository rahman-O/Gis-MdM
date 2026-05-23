import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import type { Configuration, ConfigurationApplicationSetting } from '@/features/configurations/types'

interface AppOption {
  id: number
  name: string
}

interface Props {
  configuration: Configuration
  applications: AppOption[]
  onChange: (next: Configuration) => void
}

export function ConfigurationAppSettingsTab({ configuration, applications, onChange }: Props) {
  const settings = Array.isArray(configuration.applicationSettings)
    ? configuration.applicationSettings
    : []

  const updateSetting = (index: number, patch: Partial<ConfigurationApplicationSetting>) => {
    const next = [...settings]
    next[index] = { ...next[index], ...patch }
    onChange({ ...configuration, applicationSettings: next })
  }

  const addSetting = () => {
    onChange({
      ...configuration,
      applicationSettings: [
        ...settings,
        { applicationId: null, name: '', type: 'STRING', value: '' },
      ],
    })
  }

  const removeSetting = (index: number) => {
    onChange({
      ...configuration,
      applicationSettings: settings.filter((_, idx) => idx !== index),
    })
  }

  return (
    <div className="space-y-4">
      <Button type="button" variant="outline" onClick={addSetting}>
        Add setting
      </Button>
      {settings.length === 0 ? (
        <p className="text-sm text-muted-foreground">No app settings configured.</p>
      ) : settings.map((setting, index) => (
        <div key={`setting-${index}`} className="rounded-md border p-3">
          <div className="grid gap-3 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Application</Label>
              <Select
                value={
                  setting.applicationId != null && setting.applicationId > 0
                    ? String(setting.applicationId)
                    : 'none'
                }
                onValueChange={(value) =>
                  updateSetting(index, {
                    applicationId: value === 'none' ? null : Number(value),
                  })
                }
              >
                <SelectTrigger><SelectValue placeholder="Select application" /></SelectTrigger>
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
            <div className="space-y-2">
              <Label>Name</Label>
              <Input
                value={String(setting.name ?? '')}
                onChange={(e) => updateSetting(index, { name: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Type</Label>
              <Input
                value={String(setting.type ?? '')}
                onChange={(e) => updateSetting(index, { type: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>Value</Label>
              <Input
                value={String(setting.value ?? '')}
                onChange={(e) => updateSetting(index, { value: e.target.value })}
              />
            </div>
          </div>
          <div className="mt-3">
            <Button type="button" size="sm" variant="destructive" onClick={() => removeSetting(index)}>
              Remove
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
