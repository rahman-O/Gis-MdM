import { useState } from 'react'
import { Smartphone, ChevronDown, ChevronRight, AlertTriangle } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Badge } from '@/shared/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface MdmKioskCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const KIOSK_TOGGLES: { key: keyof Configuration & string; label: string }[] = [
  { key: 'kioskHome', label: 'Home Button' },
  { key: 'kioskRecents', label: 'Recents Button' },
  { key: 'kioskNotifications', label: 'Notifications' },
  { key: 'kioskSystemInfo', label: 'System Info' },
  { key: 'kioskKeyguard', label: 'Keyguard' },
  { key: 'kioskLockButtons', label: 'Lock Buttons' },
  { key: 'kioskScreenOn', label: 'Keep Screen On' },
]

export function MdmKioskCard({ configuration, onChange }: MdmKioskCardProps) {
  const [open, setOpen] = useState(true)

  if (!configuration.kioskMode) return null

  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  const activeCount = KIOSK_TOGGLES.filter((t) => Boolean(configuration[t.key])).length

  const kioskExitValue = typeof configuration.kioskExit === 'string'
    ? configuration.kioskExit
    : configuration.kioskExit
      ? 'password'
      : 'none'

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Smartphone className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Kiosk Mode
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
            <div className="grid gap-2 sm:grid-cols-2">
              {KIOSK_TOGGLES.map(({ key, label }) => (
                <div
                  key={key}
                  className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all"
                >
                  <div className="flex items-center space-x-2.5">
                    <Checkbox
                      id={`kiosk-${key}`}
                      checked={Boolean(configuration[key])}
                      disabled={isPolicyLocked(configuration, key)}
                      onCheckedChange={(checked) =>
                        onChange({ ...configuration, [key]: checked === true })
                      }
                    />
                    <Label htmlFor={`kiosk-${key}`} className="cursor-pointer text-xs">
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

            {/* Kiosk Exit */}
            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <Label htmlFor="kiosk-exit" className="text-xs font-semibold text-muted-foreground">
                  Kiosk Exit Method
                </Label>
                <FieldLockToggle fieldKey="kioskExit" locked={isPolicyLocked(configuration, 'kioskExit')} onToggle={setLock} />
              </div>
              <Select
                value={kioskExitValue}
                disabled={isPolicyLocked(configuration, 'kioskExit')}
                onValueChange={(v) => onChange({ ...configuration, kioskExit: v as unknown as boolean })}
              >
                <SelectTrigger id="kiosk-exit">
                  <SelectValue placeholder="Select exit method" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="password">Password</SelectItem>
                  <SelectItem value="back">Back Button</SelectItem>
                  <SelectItem value="none">None</SelectItem>
                </SelectContent>
              </Select>
              {kioskExitValue === 'none' && (
                <div className="flex items-center gap-2 rounded-md bg-muted border p-2 text-muted-foreground text-xs">
                  <AlertTriangle className="h-4 w-4 shrink-0" />
                  <span>Warning: No exit method means the device cannot leave kiosk mode without MDM intervention.</span>
                </div>
              )}
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}
