import {
  Smartphone,
  Terminal,
  Monitor,
  MoreHorizontal,
  Lock,
  Trash2,
  RotateCcw,
  RefreshCw,
} from 'lucide-react'
import { Dialog, DialogContent } from '@/shared/ui/dialog'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/shared/ui/tabs'
import { Button } from '@/shared/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import { Badge } from '@/shared/ui/badge'
import type { DeviceView } from '@/features/devices/types'
import { DeviceInfoTab } from '@/features/devices/tabs/DeviceInfoTab'
import { DeviceProfileTab } from '@/features/devices/tabs/DeviceProfileTab'
import { DeviceLocationTab } from '@/features/devices/tabs/DeviceLocationTab'
import { DeviceAppsTab } from '@/features/devices/tabs/DeviceAppsTab'
import { DeviceLogsTab } from '@/features/devices/tabs/DeviceLogsTab'

export interface DeviceDetailsDialogProps {
  device: DeviceView | null
  configurationName: string | null
  onClose: () => void
}

function isOnline(lastUpdate: number | null | undefined): boolean {
  if (!lastUpdate) return false
  return Date.now() - lastUpdate < 5 * 60 * 1000
}

export function DeviceDetailsDialog({ device, configurationName, onClose }: DeviceDetailsDialogProps) {
  if (!device) return null

  const online = isOnline(device.lastUpdate)

  return (
    <Dialog open={device != null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="flex max-h-[90vh] max-w-4xl flex-col overflow-y-auto p-0">
        {/* Header */}
        <div className="flex items-center justify-between border-b p-3">
          <div className="flex items-center gap-2">
            <Smartphone className={`h-4 w-4 ${online ? 'text-green-500' : 'text-muted-foreground'}`} />
            <span className="text-sm font-medium">{device.number || '—'}</span>
            <Badge variant={online ? 'default' : 'secondary'} className="text-[10px]">
              {online ? 'Online' : 'Offline'}
            </Badge>
          </div>
          <div className="flex items-center gap-1">
            <Button variant="ghost" size="icon" className="h-7 w-7" title="Send Script">
              <Terminal className="h-3.5 w-3.5" />
            </Button>
            <Button variant="ghost" size="icon" className="h-7 w-7" title="Remote View">
              <Monitor className="h-3.5 w-3.5" />
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-7 w-7">
                  <MoreHorizontal className="h-3.5 w-3.5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>
                  <Lock className="mr-2 h-3.5 w-3.5" />
                  Lock
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Trash2 className="mr-2 h-3.5 w-3.5" />
                  Wipe
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <RotateCcw className="mr-2 h-3.5 w-3.5" />
                  Reboot
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <RefreshCw className="mr-2 h-3.5 w-3.5" />
                  Reinstall Profile
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive focus:text-destructive">
                  <Trash2 className="mr-2 h-3.5 w-3.5" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {/* Tabs */}
        <Tabs defaultValue="info" className="flex-1 px-3 pb-3">
          <TabsList className="mb-2">
            <TabsTrigger value="info" className="text-xs">
              <Smartphone className="mr-1 h-3 w-3" />
              Info
            </TabsTrigger>
            <TabsTrigger value="profile" className="text-xs">
              Profile
            </TabsTrigger>
            <TabsTrigger value="location" className="text-xs">
              Location
            </TabsTrigger>
            <TabsTrigger value="apps" className="text-xs">
              Apps
            </TabsTrigger>
            <TabsTrigger value="logs" className="text-xs">
              Logs
            </TabsTrigger>
          </TabsList>

          <TabsContent value="info">
            <DeviceInfoTab device={device} />
          </TabsContent>
          <TabsContent value="profile">
            <DeviceProfileTab device={device} configurationName={configurationName} />
          </TabsContent>
          <TabsContent value="location">
            <DeviceLocationTab device={device} />
          </TabsContent>
          <TabsContent value="apps">
            <DeviceAppsTab device={device} />
          </TabsContent>
          <TabsContent value="logs">
            <DeviceLogsTab device={device} />
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
