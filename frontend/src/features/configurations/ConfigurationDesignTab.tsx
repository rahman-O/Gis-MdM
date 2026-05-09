import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import type { Configuration } from '@/features/configurations/types'

interface Props {
  configuration: Configuration
  onChange: (next: Configuration) => void
}

export function ConfigurationDesignTab({ configuration, onChange }: Props) {
  return (
    <div className="grid gap-4 md:grid-cols-2">
      <div className="flex items-center gap-2 md:col-span-2">
        <Checkbox
          checked={Boolean(configuration.useDefaultDesignSettings)}
          onCheckedChange={(checked) =>
            onChange({ ...configuration, useDefaultDesignSettings: checked === true })
          }
        />
        <Label>Use default design settings</Label>
      </div>
      <div className="space-y-2">
        <Label>Background color</Label>
        <Input
          placeholder="#FFFFFF"
          value={String(configuration.backgroundColor ?? '')}
          onChange={(e) => onChange({ ...configuration, backgroundColor: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Text color</Label>
        <Input
          placeholder="#000000"
          value={String(configuration.textColor ?? '')}
          onChange={(e) => onChange({ ...configuration, textColor: e.target.value })}
        />
      </div>
      <div className="space-y-2 md:col-span-2">
        <Label>Background image URL</Label>
        <Input
          placeholder="https://..."
          value={String(configuration.backgroundImageUrl ?? '')}
          onChange={(e) => onChange({ ...configuration, backgroundImageUrl: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Icon size</Label>
        <Select
          value={String(configuration.iconSize ?? 'SMALL')}
          onValueChange={(value) => onChange({ ...configuration, iconSize: value })}
        >
          <SelectTrigger><SelectValue placeholder="Icon size" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="SMALL">Small</SelectItem>
            <SelectItem value="MEDIUM">Medium</SelectItem>
            <SelectItem value="LARGE">Large</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label>Desktop header</Label>
        <Select
          value={String(configuration.desktopHeader ?? 'NO_HEADER')}
          onValueChange={(value) => onChange({ ...configuration, desktopHeader: value })}
        >
          <SelectTrigger><SelectValue placeholder="Header mode" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="NO_HEADER">No header</SelectItem>
            <SelectItem value="DEVICE_ID">Device ID</SelectItem>
            <SelectItem value="DESCRIPTION">Description</SelectItem>
            <SelectItem value="TEMPLATE">Template</SelectItem>
            <SelectItem value="CUSTOM1">Custom1</SelectItem>
            <SelectItem value="CUSTOM2">Custom2</SelectItem>
            <SelectItem value="CUSTOM3">Custom3</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2 md:col-span-2">
        <Label>Header template</Label>
        <Input
          value={String(configuration.desktopHeaderText ?? '')}
          onChange={(e) => onChange({ ...configuration, desktopHeaderText: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <Label>Orientation</Label>
        <Select
          value={String(configuration.orientation ?? 'AUTO')}
          onValueChange={(value) => onChange({ ...configuration, orientation: value })}
        >
          <SelectTrigger><SelectValue placeholder="Orientation" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="AUTO">Auto</SelectItem>
            <SelectItem value="PORTRAIT">Portrait</SelectItem>
            <SelectItem value="LANDSCAPE">Landscape</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="flex items-center gap-2">
        <Checkbox
          checked={Boolean(configuration.displayStatus)}
          onCheckedChange={(checked) => onChange({ ...configuration, displayStatus: checked === true })}
        />
        <Label>Display status</Label>
      </div>
    </div>
  )
}
