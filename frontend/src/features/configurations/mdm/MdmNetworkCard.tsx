import { useState } from 'react'
import { Wifi, ChevronDown, ChevronRight } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { FieldLockToggle } from '@/features/configurations/FieldLockToggle'
import { isPolicyLocked, togglePolicyLock } from '@/features/configurations/configurationPolicyLocks'
import { restrictionsToSet, setToRestrictions } from '@/features/configurations/mdm/restrictionsRegistry'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Badge } from '@/shared/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/shared/ui/collapsible'
import { cn } from '@/shared/utils/cn'

interface MdmNetworkCardProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

const NETWORK_RESTRICTIONS = [
  { key: 'no_config_wifi', label: 'Disable Wi-Fi Configuration' },
  { key: 'no_config_vpn', label: 'Disable VPN Configuration' },
  { key: 'no_config_tethering', label: 'Disable Tethering Configuration' },
  { key: 'no_bluetooth_sharing', label: 'Disable Bluetooth Sharing' },
  { key: 'no_bluetooth', label: 'Disable Bluetooth' },
  { key: 'no_sms', label: 'Disable SMS' },
  { key: 'no_outgoing_calls', label: 'Disable Outgoing Calls' },
  { key: 'no_config_mobile_networks', label: 'Disable Mobile Network Config' },
  { key: 'no_outgoing_beam', label: 'Disable NFC Beam' },
  { key: 'no_airplane_mode', label: 'Disable Airplane Mode' },
  { key: 'no_config_cell_broadcasts', label: 'Disable Cell Broadcast Config' },
  { key: 'no_data_roaming', label: 'Disable Data Roaming' },
] as const

export function MdmNetworkCard({ configuration, onChange }: MdmNetworkCardProps) {
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
    NETWORK_RESTRICTIONS.filter((r) => restrictionSet.has(r.key)).length +
    (configuration.showWifi ? 1 : 0)

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className={cn('shadow-sm border')}>
        <CollapsibleTrigger className="w-full text-left">
          <CardHeader className="py-2 px-3 flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              {open ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
              <Wifi className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-bold text-muted-foreground">
                Network
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
            {/* Network Settings */}
            <div className="grid gap-2 sm:grid-cols-2">
              <div className="space-y-1.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="net-pushOptions" className="text-xs font-semibold text-muted-foreground">
                    Push Options
                  </Label>
                  <FieldLockToggle fieldKey="pushOptions" locked={isPolicyLocked(configuration, 'pushOptions')} onToggle={setLock} />
                </div>
                <Select
                  value={configuration.pushOptions ?? ''}
                  disabled={isPolicyLocked(configuration, 'pushOptions')}
                  onValueChange={(v) => onChange({ ...configuration, pushOptions: v || null })}
                >
                  <SelectTrigger id="net-pushOptions">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="mqttWorker">MQTT Worker</SelectItem>
                    <SelectItem value="mqttAlarm">MQTT Alarm</SelectItem>
                    <SelectItem value="polling">Polling</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex items-center justify-between rounded-md border p-2 hover:bg-muted/20 transition-all">
                <div className="flex items-center space-x-2.5">
                  <Checkbox
                    id="net-showWifi"
                    checked={Boolean(configuration.showWifi)}
                    disabled={isPolicyLocked(configuration, 'showWifi')}
                    onCheckedChange={(checked) =>
                      onChange({ ...configuration, showWifi: checked === true })
                    }
                  />
                  <Label htmlFor="net-showWifi" className="cursor-pointer text-xs">
                    Show Wi-Fi
                  </Label>
                </div>
                <FieldLockToggle
                  fieldKey="showWifi"
                  locked={isPolicyLocked(configuration, 'showWifi')}
                  onToggle={setLock}
                />
              </div>
            </div>

            {/* Network Restrictions */}
            <div className="space-y-2">
              <Label className="text-xs font-semibold text-muted-foreground">
                Network Restrictions
              </Label>
              <div className="grid gap-2 sm:grid-cols-2">
                {NETWORK_RESTRICTIONS.map(({ key, label }) => (
                  <div
                    key={key}
                    className="flex items-center space-x-2.5 rounded-md border p-2 hover:bg-muted/20 transition-all"
                  >
                    <Checkbox
                      id={`net-${key}`}
                      checked={restrictionSet.has(key)}
                      disabled={Boolean(configuration.permissive) || isPolicyLocked(configuration, 'restrictions')}
                      onCheckedChange={(checked) => toggleRestriction(key, checked === true)}
                    />
                    <Label htmlFor={`net-${key}`} className="cursor-pointer text-xs">
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
