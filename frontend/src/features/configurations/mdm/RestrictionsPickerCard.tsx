import { useState, useMemo } from 'react'
import { ShieldCheck, ChevronDown, ChevronRight, Search } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import {
  RESTRICTIONS_REGISTRY,
  CATEGORY_LABELS,
  restrictionsToSet,
  setToRestrictions,
  type RestrictionCategory,
  type RestrictionDefinition,
} from '@/features/configurations/mdm/restrictionsRegistry'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Input } from '@/shared/ui/input'
import { Badge } from '@/shared/ui/badge'
import { Textarea } from '@/shared/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface RestrictionsPickerCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

/** Restrictions already shown in Security, Apps, and Network cards */
const SHOWN_ELSEWHERE = new Set([
  // Security (5)
  'no_camera',
  'no_factory_reset',
  'no_safe_boot',
  'no_debugging_features',
  'no_shutdown',
  // Apps (4)
  'no_install_apps',
  'no_uninstall_apps',
  'no_install_unknown_sources',
  'no_install_unknown_sources_globally',
  // Network (12)
  'no_config_wifi',
  'no_config_vpn',
  'no_config_tethering',
  'no_bluetooth_sharing',
  'no_bluetooth',
  'no_sms',
  'no_outgoing_calls',
  'no_config_mobile_networks',
  'no_outgoing_beam',
  'no_airplane_mode',
  'no_config_cell_broadcasts',
  'no_data_roaming',
])

/** Only show restrictions NOT already displayed in other cards */
const PICKER_RESTRICTIONS = RESTRICTIONS_REGISTRY.filter(
  (r) => !SHOWN_ELSEWHERE.has(r.key),
)

function groupByCategory(
  restrictions: RestrictionDefinition[],
): Map<RestrictionCategory, RestrictionDefinition[]> {
  const map = new Map<RestrictionCategory, RestrictionDefinition[]>()
  for (const r of restrictions) {
    const list = map.get(r.category) ?? []
    list.push(r)
    map.set(r.category, list)
  }
  return map
}

export function RestrictionsPickerCard({ configuration, onChange }: RestrictionsPickerCardProps) {
  const [open, setOpen] = useState(false)
  const [search, setSearch] = useState('')
  const [showAdvanced, setShowAdvanced] = useState(false)

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

  const filteredRestrictions = useMemo(() => {
    if (!search.trim()) return PICKER_RESTRICTIONS
    const lower = search.toLowerCase()
    return PICKER_RESTRICTIONS.filter(
      (r) =>
        r.labelEn.toLowerCase().includes(lower) ||
        r.key.toLowerCase().includes(lower) ||
        r.descriptionEn.toLowerCase().includes(lower),
    )
  }, [search])

  const grouped = useMemo(() => groupByCategory(filteredRestrictions), [filteredRestrictions])

  const activeCount = PICKER_RESTRICTIONS.filter((r) => restrictionSet.has(r.key)).length

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <ShieldCheck className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                All Restrictions
              </CardTitle>
            </div>
            <div className="flex items-center gap-2">
              {activeCount > 0 && (
                <Badge variant="secondary" className="text-xs">
                  {activeCount} active
                </Badge>
              )}
              <Badge variant="outline" className="text-xs">
                {PICKER_RESTRICTIONS.length} restrictions
              </Badge>
            </div>
          </CardHeader>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <CardContent className="px-3 pb-3 pt-0 space-y-4">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search restrictions..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>

            {/* Grouped Restrictions */}
            {Array.from(grouped.entries()).map(([category, restrictions]) => (
              <div key={category} className="space-y-2">
                <Label className="text-xs font-semibold text-muted-foreground">
                  {CATEGORY_LABELS[category].en}
                </Label>
                <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                  {restrictions.map((r) => (
                    <div
                      key={r.key}
                      className="flex items-center space-x-2.5 rounded-md border p-2 hover:bg-muted/20 transition-all"
                    >
                      <Checkbox
                        id={`rp-${r.key}`}
                        checked={restrictionSet.has(r.key)}
                        disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
                        onCheckedChange={(checked) => toggleRestriction(r.key, checked === true)}
                      />
                      <Label htmlFor={`rp-${r.key}`} className="cursor-pointer text-xs flex-1">
                        {r.labelEn}
                      </Label>
                      {r.minAndroid !== '5.0' && (
                        <Badge variant="outline" className="text-[10px] shrink-0">
                          {r.minAndroid}+
                        </Badge>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            ))}

            {filteredRestrictions.length === 0 && (
              <p className="text-sm text-muted-foreground text-center py-4">
                No restrictions match your search.
              </p>
            )}

            {/* Advanced: Raw Textarea */}
            <div className="border-t pt-4 space-y-2">
              <button
                type="button"
                className="text-xs font-semibold text-muted-foreground uppercase hover:text-foreground transition-colors"
                onClick={() => setShowAdvanced(!showAdvanced)}
              >
                {showAdvanced ? '▾ Hide Advanced' : '▸ Show Advanced'}
              </button>
              {showAdvanced && (
                <div className="space-y-1.5">
                  <div className="flex items-center justify-between">
                    <Label htmlFor="rp-raw" className="text-xs font-semibold text-muted-foreground">
                      Raw Restrictions (comma-separated)
                    </Label>
                    <FieldLockToggle
                      fieldKey="restrictions"
                      locked={isPolicyLocked(configuration, 'restrictions')}
                      onToggle={setLock}
                    />
                  </div>
                  <Textarea
                    id="rp-raw"
                    rows={4}
                    placeholder="e.g. no_camera,no_sms,no_fun"
                    value={configuration.restrictions ?? ''}
                    disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
                    onChange={(e) => onChange({ ...configuration, restrictions: e.target.value })}
                  />
                  <p className="text-xs text-muted-foreground">
                    Editing this field directly will sync with the checkboxes above and in other cards.
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}
