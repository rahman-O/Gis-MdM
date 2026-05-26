import { useState } from 'react'
import { Package, ChevronDown, ChevronRight } from 'lucide-react'
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

interface MdmAppsUpdatesCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const APP_TOGGLES: { key: keyof Configuration & string; label: string }[] = [
  { key: 'permissive', label: 'Permissive' },
  { key: 'runDefaultLauncher', label: 'Run Default Launcher' },
  { key: 'autostartForeground', label: 'Autostart Foreground' },
  { key: 'scheduleAppUpdate', label: 'Schedule App Update' },
]

const APP_RESTRICTIONS = [
  { key: 'no_install_apps', label: 'Disable App Installation' },
  { key: 'no_uninstall_apps', label: 'Disable App Uninstallation' },
  { key: 'no_install_unknown_sources', label: 'Disable Unknown Sources' },
  { key: 'no_install_unknown_sources_globally', label: 'Disable Unknown Sources Globally' },
] as const

function toText(value: unknown): string {
  return value == null ? '' : String(value)
}

export function MdmAppsUpdatesCard({ configuration, onChange }: MdmAppsUpdatesCardProps) {
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
    APP_TOGGLES.filter((t) => Boolean(configuration[t.key])).length +
    APP_RESTRICTIONS.filter((r) => restrictionSet.has(r.key)).length

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Package className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Apps & Updates
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
              {APP_TOGGLES.map(({ key, label }) => (
                <div
                  key={key}
                  className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all"
                >
                  <div className="flex items-center space-x-2.5">
                    <Checkbox
                      id={`app-${key}`}
                      checked={Boolean(configuration[key])}
                      disabled={isPolicyLocked(configuration, key)}
                      onCheckedChange={(checked) =>
                        onChange({ ...configuration, [key]: checked === true })
                      }
                    />
                    <Label htmlFor={`app-${key}`} className="cursor-pointer text-xs">
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

            {/* Selects */}
            <div className="grid gap-2 sm:grid-cols-3">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-permissions" className="text-xs font-semibold text-muted-foreground">
                    App Permissions
                  </Label>
                  <FieldLockToggle fieldKey="appPermissions" locked={isPolicyLocked(configuration, 'appPermissions')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.appPermissions ?? ''}
                  disabled={isPolicyLocked(configuration, 'appPermissions')}
                  onValueChange={(v) => onChange({ ...configuration, appPermissions: v || null })}
                >
                  <SelectTrigger id="app-permissions">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="grant">Grant</SelectItem>
                    <SelectItem value="deny">Deny</SelectItem>
                    <SelectItem value="default">Default</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-updateType" className="text-xs font-semibold text-muted-foreground">
                    System Update Type
                  </Label>
                  <FieldLockToggle fieldKey="systemUpdateType" locked={isPolicyLocked(configuration, 'systemUpdateType')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.systemUpdateType != null ? String(configuration.systemUpdateType) : ''}
                  disabled={isPolicyLocked(configuration, 'systemUpdateType')}
                  onValueChange={(v) => onChange({ ...configuration, systemUpdateType: v ? Number(v) : null })}
                >
                  <SelectTrigger id="app-updateType">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="0">Auto</SelectItem>
                    <SelectItem value="1">Deferred</SelectItem>
                    <SelectItem value="2">Window</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-downloadUpdates" className="text-xs font-semibold text-muted-foreground">
                    Download Updates
                  </Label>
                  <FieldLockToggle fieldKey="downloadUpdates" locked={isPolicyLocked(configuration, 'downloadUpdates')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.downloadUpdates ?? ''}
                  disabled={isPolicyLocked(configuration, 'downloadUpdates')}
                  onValueChange={(v) => onChange({ ...configuration, downloadUpdates: v || null })}
                >
                  <SelectTrigger id="app-downloadUpdates">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="wifi">Wi-Fi Only</SelectItem>
                    <SelectItem value="any">Any Network</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Time Inputs */}
            <div className="grid gap-2 sm:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-sysUpdateFrom" className="text-xs font-semibold text-muted-foreground">
                    System Update From
                  </Label>
                  <FieldLockToggle fieldKey="systemUpdateFrom" locked={isPolicyLocked(configuration, 'systemUpdateFrom')} onToggle={setLock} />
                </div>
                <Input
                  id="app-sysUpdateFrom"
                  type="time"
                  value={toText(configuration.systemUpdateFrom)}
                  disabled={isPolicyLocked(configuration, 'systemUpdateFrom')}
                  onChange={(e) => onChange({ ...configuration, systemUpdateFrom: e.target.value || null })}
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-sysUpdateTo" className="text-xs font-semibold text-muted-foreground">
                    System Update To
                  </Label>
                  <FieldLockToggle fieldKey="systemUpdateTo" locked={isPolicyLocked(configuration, 'systemUpdateTo')} onToggle={setLock} />
                </div>
                <Input
                  id="app-sysUpdateTo"
                  type="time"
                  value={toText(configuration.systemUpdateTo)}
                  disabled={isPolicyLocked(configuration, 'systemUpdateTo')}
                  onChange={(e) => onChange({ ...configuration, systemUpdateTo: e.target.value || null })}
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-appUpdateFrom" className="text-xs font-semibold text-muted-foreground">
                    App Update From
                  </Label>
                  <FieldLockToggle fieldKey="appUpdateFrom" locked={isPolicyLocked(configuration, 'appUpdateFrom')} onToggle={setLock} />
                </div>
                <Input
                  id="app-appUpdateFrom"
                  type="time"
                  value={toText(configuration.appUpdateFrom)}
                  disabled={isPolicyLocked(configuration, 'appUpdateFrom')}
                  onChange={(e) => onChange({ ...configuration, appUpdateFrom: e.target.value || null })}
                />
              </div>
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="app-appUpdateTo" className="text-xs font-semibold text-muted-foreground">
                    App Update To
                  </Label>
                  <FieldLockToggle fieldKey="appUpdateTo" locked={isPolicyLocked(configuration, 'appUpdateTo')} onToggle={setLock} />
                </div>
                <Input
                  id="app-appUpdateTo"
                  type="time"
                  value={toText(configuration.appUpdateTo)}
                  disabled={isPolicyLocked(configuration, 'appUpdateTo')}
                  onChange={(e) => onChange({ ...configuration, appUpdateTo: e.target.value || null })}
                />
              </div>
            </div>

            {/* App Restrictions */}
            <div className="space-y-2">
              <Label className="text-xs font-semibold text-muted-foreground">
                App Restrictions
              </Label>
              <div className="grid gap-2 sm:grid-cols-2">
                {APP_RESTRICTIONS.map(({ key, label }) => (
                  <div
                    key={key}
                    className="flex items-center space-x-2.5 rounded-md border p-2 hover:bg-muted/20 transition-all"
                  >
                    <Checkbox
                      id={`app-r-${key}`}
                      checked={restrictionSet.has(key)}
                      disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
                      onCheckedChange={(checked) => toggleRestriction(key, checked === true)}
                    />
                    <Label htmlFor={`app-r-${key}`} className="cursor-pointer text-xs">
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
