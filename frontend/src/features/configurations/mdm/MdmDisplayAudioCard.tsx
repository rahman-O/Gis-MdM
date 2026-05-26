import { useState } from 'react'
import { Monitor, ChevronDown, ChevronRight } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Input } from '@/shared/ui/input'
import { Badge } from '@/shared/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface MdmDisplayAudioCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const DISPLAY_TOGGLES: { key: keyof Configuration & string; label: string }[] = [
  { key: 'displayStatus', label: 'Lock Status Bar' },
  { key: 'autoBrightness', label: 'Auto Brightness' },
  { key: 'manageTimeout', label: 'Manage Timeout' },
  { key: 'lockVolume', label: 'Lock Volume' },
  { key: 'manageVolume', label: 'Manage Volume' },
]

function toText(value: unknown): string {
  return value == null ? '' : String(value)
}

export function MdmDisplayAudioCard({ configuration, onChange }: MdmDisplayAudioCardProps) {
  const [open, setOpen] = useState(true)

  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  const activeCount = DISPLAY_TOGGLES.filter((t) => Boolean(configuration[t.key])).length +
    (configuration.brightness != null ? 1 : 0) +
    (configuration.timeout != null ? 1 : 0) +
    (configuration.volume != null ? 1 : 0)

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Monitor className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Display & Audio
              </CardTitle>
            </div>
            {activeCount > 0 && (
              <Badge variant="secondary" className="text-xs">
                {activeCount} active
              </Badge>
            )}
          </CardHeader>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <CardContent className="px-3 pb-3 pt-0 space-y-4">
            {/* Toggles */}
            <div className="grid gap-2 sm:grid-cols-2">
              {DISPLAY_TOGGLES.map(({ key, label }) => (
                <div
                  key={key}
                  className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all"
                >
                  <div className="flex items-center space-x-2.5">
                    <Checkbox
                      id={`da-${key}`}
                      checked={Boolean(configuration[key])}
                      disabled={isPolicyLocked(configuration, key)}
                      onCheckedChange={(checked) =>
                        onChange({ ...configuration, [key]: checked === true })
                      }
                    />
                    <Label htmlFor={`da-${key}`} className="cursor-pointer text-xs">
                      {label}
                    </Label>
                  </div>
                  <FieldLockToggle
                    fieldKey={key}
                    locked={isPolicyLocked(configuration, key)}
                    onToggle={setLock}
                  />
                </div>
              ))}
            </div>

            {/* Numeric Inputs */}
            <div className="grid gap-2 sm:grid-cols-3">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-brightness" className="text-xs font-semibold text-muted-foreground">
                    Brightness (0-255)
                  </Label>
                  <FieldLockToggle fieldKey="brightness" locked={isPolicyLocked(configuration, 'brightness')} onToggle={setLock} />
                </div>
                <Input
                  id="da-brightness"
                  type="number"
                  min={0}
                  max={255}
                  placeholder="0-255"
                  value={toText(configuration.brightness)}
                  disabled={isPolicyLocked(configuration, 'brightness')}
                  onChange={(e) =>
                    onChange({ ...configuration, brightness: e.target.value ? Number(e.target.value) : null })
                  }
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-timeout" className="text-xs font-semibold text-muted-foreground">
                    Timeout (seconds)
                  </Label>
                  <FieldLockToggle fieldKey="timeout" locked={isPolicyLocked(configuration, 'timeout')} onToggle={setLock} />
                </div>
                <Input
                  id="da-timeout"
                  type="number"
                  min={0}
                  placeholder="Seconds"
                  value={toText(configuration.timeout)}
                  disabled={isPolicyLocked(configuration, 'timeout')}
                  onChange={(e) =>
                    onChange({ ...configuration, timeout: e.target.value ? Number(e.target.value) : null })
                  }
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-volume" className="text-xs font-semibold text-muted-foreground">
                    Volume (0-100)
                  </Label>
                  <FieldLockToggle fieldKey="volume" locked={isPolicyLocked(configuration, 'volume')} onToggle={setLock} />
                </div>
                <Input
                  id="da-volume"
                  type="number"
                  min={0}
                  max={100}
                  placeholder="0-100"
                  value={toText(configuration.volume)}
                  disabled={isPolicyLocked(configuration, 'volume')}
                  onChange={(e) =>
                    onChange({ ...configuration, volume: e.target.value ? Number(e.target.value) : null })
                  }
                />
              </div>
            </div>

            {/* Selects */}
            <div className="grid gap-2 sm:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-orientation" className="text-xs font-semibold text-muted-foreground">
                    Orientation
                  </Label>
                  <FieldLockToggle fieldKey="orientation" locked={isPolicyLocked(configuration, 'orientation')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.orientation ?? ''}
                  disabled={isPolicyLocked(configuration, 'orientation')}
                  onValueChange={(v) => onChange({ ...configuration, orientation: v || null })}
                >
                  <SelectTrigger id="da-orientation">
                    <SelectValue placeholder="Select orientation" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="free">Free</SelectItem>
                    <SelectItem value="portrait">Portrait</SelectItem>
                    <SelectItem value="landscape">Landscape</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-iconSize" className="text-xs font-semibold text-muted-foreground">
                    Icon Size
                  </Label>
                  <FieldLockToggle fieldKey="iconSize" locked={isPolicyLocked(configuration, 'iconSize')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.iconSize ?? ''}
                  disabled={isPolicyLocked(configuration, 'iconSize')}
                  onValueChange={(v) => onChange({ ...configuration, iconSize: v || null })}
                >
                  <SelectTrigger id="da-iconSize">
                    <SelectValue placeholder="Select size" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="small">Small</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="large">Large</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Color & Background */}
            <div className="grid gap-2 sm:grid-cols-3">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-bgColor" className="text-xs font-semibold text-muted-foreground">
                    Background Color
                  </Label>
                  <FieldLockToggle fieldKey="backgroundColor" locked={isPolicyLocked(configuration, 'backgroundColor')} onToggle={setLock} />
                </div>
                <Input
                  id="da-bgColor"
                  type="color"
                  value={configuration.backgroundColor ?? '#ffffff'}
                  disabled={isPolicyLocked(configuration, 'backgroundColor')}
                  onChange={(e) => onChange({ ...configuration, backgroundColor: e.target.value })}
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-textColor" className="text-xs font-semibold text-muted-foreground">
                    Text Color
                  </Label>
                  <FieldLockToggle fieldKey="textColor" locked={isPolicyLocked(configuration, 'textColor')} onToggle={setLock} />
                </div>
                <Input
                  id="da-textColor"
                  type="color"
                  value={configuration.textColor ?? '#000000'}
                  disabled={isPolicyLocked(configuration, 'textColor')}
                  onChange={(e) => onChange({ ...configuration, textColor: e.target.value })}
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="da-bgImage" className="text-xs font-semibold text-muted-foreground">
                    Background Image URL
                  </Label>
                  <FieldLockToggle fieldKey="backgroundImageUrl" locked={isPolicyLocked(configuration, 'backgroundImageUrl')} onToggle={setLock} />
                </div>
                <Input
                  id="da-bgImage"
                  type="url"
                  placeholder="https://..."
                  value={toText(configuration.backgroundImageUrl)}
                  disabled={isPolicyLocked(configuration, 'backgroundImageUrl')}
                  onChange={(e) => onChange({ ...configuration, backgroundImageUrl: e.target.value || null })}
                />
              </div>
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}
