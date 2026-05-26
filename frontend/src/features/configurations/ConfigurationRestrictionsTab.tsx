import { Cpu, Monitor, Smartphone, Shield, Package, Wifi, ShieldCheck } from 'lucide-react'
import type { Configuration } from '@/features/configurations/types'
import { MdmHardwareCard } from '@/features/configurations/mdm/MdmHardwareCard'
import { MdmDisplayAudioCard } from '@/features/configurations/mdm/MdmDisplayAudioCard'
import { MdmKioskCard } from '@/features/configurations/mdm/MdmKioskCard'
import { MdmSecurityCard } from '@/features/configurations/mdm/MdmSecurityCard'
import { MdmAppsUpdatesCard } from '@/features/configurations/mdm/MdmAppsUpdatesCard'
import { MdmNetworkCard } from '@/features/configurations/mdm/MdmNetworkCard'
import { RestrictionsPickerCard } from '@/features/configurations/mdm/RestrictionsPickerCard'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/shared/ui/tabs'

interface ConfigurationRestrictionsTabProps {
  configuration: Configuration
  onChange: (configuration: Configuration) => void
}

export function ConfigurationRestrictionsTab({ configuration, onChange }: ConfigurationRestrictionsTabProps) {
  return (
    <Tabs defaultValue="hardware" className="w-full">
      <TabsList className="flex flex-wrap h-auto gap-1 bg-muted/50 p-1.5 rounded-lg">
        <TabsTrigger value="hardware" className="flex items-center gap-1.5 text-xs">
          <Cpu className="h-3.5 w-3.5 text-muted-foreground" />
          Hardware
        </TabsTrigger>
        <TabsTrigger value="display" className="flex items-center gap-1.5 text-xs">
          <Monitor className="h-3.5 w-3.5 text-muted-foreground" />
          Display & Audio
        </TabsTrigger>
        {configuration.kioskMode ? (
          <TabsTrigger value="kiosk" className="flex items-center gap-1.5 text-xs">
            <Smartphone className="h-3.5 w-3.5 text-muted-foreground" />
            Kiosk
          </TabsTrigger>
        ) : null}
        <TabsTrigger value="security" className="flex items-center gap-1.5 text-xs">
          <Shield className="h-3.5 w-3.5 text-muted-foreground" />
          Security
        </TabsTrigger>
        <TabsTrigger value="apps" className="flex items-center gap-1.5 text-xs">
          <Package className="h-3.5 w-3.5 text-muted-foreground" />
          Apps & Updates
        </TabsTrigger>
        <TabsTrigger value="network" className="flex items-center gap-1.5 text-xs">
          <Wifi className="h-3.5 w-3.5 text-muted-foreground" />
          Network
        </TabsTrigger>
        <TabsTrigger value="all" className="flex items-center gap-1.5 text-xs">
          <ShieldCheck className="h-3.5 w-3.5 text-muted-foreground" />
          All Restrictions
        </TabsTrigger>
      </TabsList>

      <TabsContent value="hardware">
        <MdmHardwareCard configuration={configuration} onChange={onChange} />
      </TabsContent>

      <TabsContent value="display">
        <MdmDisplayAudioCard configuration={configuration} onChange={onChange} />
      </TabsContent>

      {configuration.kioskMode ? (
        <TabsContent value="kiosk">
          <MdmKioskCard configuration={configuration} onChange={onChange} />
        </TabsContent>
      ) : null}

      <TabsContent value="security">
        <MdmSecurityCard configuration={configuration} onChange={onChange} />
      </TabsContent>

      <TabsContent value="apps">
        <MdmAppsUpdatesCard configuration={configuration} onChange={onChange} />
      </TabsContent>

      <TabsContent value="network">
        <MdmNetworkCard configuration={configuration} onChange={onChange} />
      </TabsContent>

      <TabsContent value="all">
        <RestrictionsPickerCard configuration={configuration} onChange={onChange} />
      </TabsContent>
    </Tabs>
  )
}
