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
