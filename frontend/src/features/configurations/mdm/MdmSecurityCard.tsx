import { useState } from 'react'
import { Shield, ChevronDown, ChevronRight } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { restrictionsToSet, setToRestrictions } from '@/features/configurations/mdm/restrictionsRegistry'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Input } from '@/shared/ui/input'
import { Badge } from '@/shared/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface MdmSecurityCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const SECURITY_RESTRICTIONS = [
  { key: 'no_camera', label: 'Disable Camera' },
  { key: 'no_factory_reset', label: 'Disable Factory Reset' },
  { key: 'no_safe_boot', label: 'Disable Safe Boot' },
  { key: 'no_debugging_features', label: 'Disable Debugging' },
  { key: 'no_shutdown', label: 'Disable Shutdown' },
] as const

function toText(value: unknown): string {
  return value == null ? '' : String(value)
}

export function MdmSecurityCard({ configuration, onChange }: MdmSecurityCardProps) {
  const [open, setOpen] = useState(true)

  const setLock = (fieldKey: string, locked: boolean) => {
    onChange(togglePolicyLock(configuration, fieldKey, locked))
  }

  const restrictionSet = restrictionsToSet(configuration.restrictions)

  const toggleRestriction = (key: string, enabled: boolean) => {
    const updated = new Set(restrictionSet)
    if (enabled) {
      updated.add(key)
    } else {
      updated.delete(key)
    }
    onChange({ ...configuration, restrictions: setToRestrictions(updated) })
  }

  const activeCount =
    (configuration.passwordMode && configuration.passwordMode !== 'none' ? 1 : 0) +
    (configuration.password ? 1 : 0) +
    (configuration.lockSafeSettings ? 1 : 0) +
    SECURITY_RESTRICTIONS.filter((r) => restrictionSet.has(r.key)).length

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Shield className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Security
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
            {/* Password Mode & Admin Password */}
            <div className="grid gap-2 sm:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="sec-passwordMode" className="text-xs font-semibold text-muted-foreground">
                    Password Mode
                  </Label>
                  <FieldLockToggle fieldKey="passwordMode" locked={isPolicyLocked(configuration, 'passwordMode')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.passwordMode ?? ''}
                  disabled={isPolicyLocked(configuration, 'passwordMode')}
                  onValueChange={(v) => onChange({ ...configuration, passwordMode: v || null })}
                >
                  <SelectTrigger id="sec-passwordMode">
                    <SelectValue placeholder="Select mode" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="strong">Strong</SelectItem>
                    <SelectItem value="numeric">Numeric</SelectItem>
                    <SelectItem value="any">Any</SelectItem>
                    <SelectItem value="none">None</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="sec-password" className="text-xs font-semibold text-muted-foreground">
                    Admin Password
                  </Label>
                  <FieldLockToggle fieldKey="password" locked={isPolicyLocked(configuration, 'password')} onToggle={setLock} />
                </div>
                <Input
                  id="sec-password"
                  type="password"
                  placeholder="Admin password"
                  value={toText(configuration.password)}
                  disabled={isPolicyLocked(configuration, 'password')}
                  onChange={(e) => onChange({ ...configuration, password: e.target.value || null })}
                />
              </div>
            </div>

            {/* Lock Safe Settings */}
            <div className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all">
              <div className="flex items-center space-x-2.5">
                <Checkbox
                  id="sec-lockSafeSettings"
                  checked={Boolean(configuration.lockSafeSettings)}
                  disabled={isPolicyLocked(configuration, 'lockSafeSettings')}
                  onCheckedChange={(checked) =>
                    onChange({ ...configuration, lockSafeSettings: checked === true })
                  }
                />
                <Label htmlFor="sec-lockSafeSettings" className="cursor-pointer text-xs">
                  Lock Safe Settings
                </Label>
              </div>
              <FieldLockToggle
                fieldKey="lockSafeSettings"
                locked={isPolicyLocked(configuration, 'lockSafeSettings')}
                onToggle={setLock}
              />
            </div>

            {/* Security Restrictions */}
            <div className="space-y-2">
              <Label className="text-xs font-semibold text-muted-foreground">
                Security Restrictions
              </Label>
              <div className="grid gap-2 sm:grid-cols-2">
                {SECURITY_RESTRICTIONS.map(({ key, label }) => (
                  <div
                    key={key}
                    className="flex items-center space-x-2.5 rounded-md border p-2 hover:bg-muted/20 transition-all"
                  >
                    <Checkbox
                      id={`sec-${key}`}
                      checked={restrictionSet.has(key)}
                      disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
                      onCheckedChange={(checked) => toggleRestriction(key, checked === true)}
                    />
                    <Label htmlFor={`sec-${key}`} className="cursor-pointer text-xs">
                      {label}
                    </Label>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}
