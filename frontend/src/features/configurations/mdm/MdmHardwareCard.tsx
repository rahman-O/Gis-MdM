import { useState } from 'react'
import { Cpu, ChevronDown, ChevronRight } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Badge } from '@/shared/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface MdmHardwareCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const HARDWARE_TOGGLES: { key: keyof Configuration & string; label: string }[] = [
  { key: 'gps', label: 'GPS' },
  { key: 'bluetooth', label: 'Bluetooth' },
  { key: 'wifi', label: 'Wi-Fi' },
  { key: 'mobileData', label: 'Mobile Data' },
  { key: 'usbStorage', label: 'USB Storage' },
  { key: 'disableScreenshots', label: 'Disable Screenshots' },
]

export function MdmHardwareCard({ configuration, onChange }: MdmHardwareCardProps) {
  const [open, setOpen] = useState(true)

  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  const activeCount = HARDWARE_TOGGLES.filter(
    (t) => Boolean(configuration[t.key]),
  ).length

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Cpu className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Hardware
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
          <CardContent className="px-3 pb-3 pt-0">
            <div className="grid gap-2 sm:grid-cols-2">
              {HARDWARE_TOGGLES.map(({ key, label }) => (
                <div
                  key={key}
                  className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all"
                >
                  <div className="flex items-center space-x-2.5">
                    <Checkbox
                      id={`hw-${key}`}
                      checked={Boolean(configuration[key])}
                      disabled={isPolicyLocked(configuration, key)}
                      onCheckedChange={(checked) =>
                        onChange({ ...configuration, [key]: checked === true })
                      }
                    />
                    <Label htmlFor={`hw-${key}`} className="cursor-pointer text-xs">
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
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}
