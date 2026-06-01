import { useEffect, useState } from 'react'
import {
  Smartphone,
  Terminal,
  Monitor,
  MoreHorizontal,
  Lock,
  Trash2,
  RotateCcw,
  RefreshCw,
  MapPin,
  Sliders,
  AppWindow,
  ScrollText,
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
import { Skeleton } from '@/shared/ui/skeleton'
import { cn } from '@/shared/utils/cn'
import type { DeviceView } from '@/features/devices/types'
import * as deviceService from '@/features/devices/deviceService'
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
  const [fullDevice, setFullDevice] = useState<DeviceView | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!device) {
      setFullDevice(null)
      return
    }
    let cancelled = false
    setLoading(true)
    deviceService
      .getDevice(device.number)
      .then((d) => {
        if (!cancelled) setFullDevice(d)
      })
      .catch(() => {
        // Fall back to the list-level device data if full fetch fails
        if (!cancelled) setFullDevice(device)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [device])

  if (!device) return null

  const displayDevice = fullDevice ?? device
  const online = isOnline(displayDevice.lastUpdate)

  return (
    <Dialog open={device != null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="flex max-h-[92vh] min-h-[70vh] w-[95vw] max-w-5xl flex-col overflow-hidden p-0 rounded-2xl border border-border/50 shadow-lg bg-background">
        {/* Header Telemetry Banner */}
        <div className="bg-muted/30 border-b border-border/55 p-5 flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-3.5">
            <div className="relative flex items-center justify-center h-11 w-11 rounded-xl bg-background border border-border/70 shadow-sm">
              <Smartphone className={cn("h-5 w-5 transition-colors", online ? 'text-emerald-500' : 'text-muted-foreground/60')} />
              {online ? (
                <span className="absolute -top-0.5 -right-0.5 flex h-2.5 w-2.5">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
                </span>
              ) : null}
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="text-lg font-bold text-foreground/90">{displayDevice.number || '—'}</span>
                <span className={cn(
                  "inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-semibold border transition-all",
                  online 
                    ? "bg-emerald-500/10 text-emerald-700 border-emerald-500/20" 
                    : "bg-muted text-muted-foreground border-border/60"
                )}>
                  {online ? 'Online' : 'Offline'}
                </span>
              </div>
              <p className="text-xs text-muted-foreground mt-0.5 flex flex-wrap items-center gap-1.5">
                <span>Model: {displayDevice.info?.model || displayDevice.model || '—'}</span>
                {configurationName ? (
                  <>
                    <span className="text-border">•</span>
                    <span className="font-medium text-foreground/75">{configurationName}</span>
                  </>
                ) : null}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" className="h-9 rounded-lg gap-1.5 text-xs font-medium border-border/80 shadow-xs hover:bg-accent/60 focus:outline-none" title="Send Script">
              <Terminal className="h-3.5 w-3.5 text-muted-foreground" />
              <span>Send Script</span>
            </Button>
            <Button variant="outline" size="sm" className="h-9 rounded-lg gap-1.5 text-xs font-medium border-border/80 shadow-xs hover:bg-accent/60 focus:outline-none" title="Remote View">
              <Monitor className="h-3.5 w-3.5 text-muted-foreground" />
              <span>Remote View</span>
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="icon" className="h-9 w-9 rounded-lg border-border/80 shadow-xs hover:bg-accent/60 focus:outline-none">
                  <MoreHorizontal className="h-4 w-4 text-muted-foreground" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48 p-1 border border-border/80 shadow-md rounded-lg">
                <DropdownMenuItem className="text-xs gap-2">
                  <Lock className="h-3.5 w-3.5 text-muted-foreground" />
                  Lock
                </DropdownMenuItem>
                <DropdownMenuItem className="text-xs gap-2">
                  <Trash2 className="h-3.5 w-3.5 text-muted-foreground" />
                  Wipe
                </DropdownMenuItem>
                <DropdownMenuItem className="text-xs gap-2">
                  <RotateCcw className="h-3.5 w-3.5 text-muted-foreground" />
                  Reboot
                </DropdownMenuItem>
                <DropdownMenuItem className="text-xs gap-2">
                  <RefreshCw className="h-3.5 w-3.5 text-muted-foreground" />
                  Reinstall Profile
                </DropdownMenuItem>
                <DropdownMenuItem className="text-xs gap-2 text-destructive focus:text-destructive focus:bg-destructive/10">
                  <Trash2 className="h-3.5 w-3.5" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {/* Tabs container */}
        <Tabs defaultValue="info" className="flex flex-1 flex-col overflow-hidden p-4">
          <TabsList className="mb-4 shrink-0 bg-muted/50 border border-border/45 p-1 rounded-xl flex gap-1 self-start">
            <TabsTrigger value="info" className="text-xs font-medium rounded-lg px-3 py-1.5 text-muted-foreground data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-xs transition-all duration-150">
              <Smartphone className="mr-1.5 h-3.5 w-3.5 text-muted-foreground/80" />
              Info
            </TabsTrigger>
            <TabsTrigger value="profile" className="text-xs font-medium rounded-lg px-3 py-1.5 text-muted-foreground data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-xs transition-all duration-150">
              <Sliders className="mr-1.5 h-3.5 w-3.5 text-muted-foreground/80" />
              Profile
            </TabsTrigger>
            <TabsTrigger value="location" className="text-xs font-medium rounded-lg px-3 py-1.5 text-muted-foreground data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-xs transition-all duration-150">
              <MapPin className="mr-1.5 h-3.5 w-3.5 text-muted-foreground/80" />
              Location
            </TabsTrigger>
            <TabsTrigger value="apps" className="text-xs font-medium rounded-lg px-3 py-1.5 text-muted-foreground data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-xs transition-all duration-150">
              <AppWindow className="mr-1.5 h-3.5 w-3.5 text-muted-foreground/80" />
              Apps
            </TabsTrigger>
            <TabsTrigger value="logs" className="text-xs font-medium rounded-lg px-3 py-1.5 text-muted-foreground data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-xs transition-all duration-150">
              <ScrollText className="mr-1.5 h-3.5 w-3.5 text-muted-foreground/80" />
              Logs
            </TabsTrigger>
          </TabsList>

          {loading ? (
            <div className="flex-1 space-y-4 p-4">
              <Skeleton className="h-8 w-3/4 rounded-lg animate-pulse" />
              <Skeleton className="h-4 w-full rounded-md animate-pulse" />
              <Skeleton className="h-4 w-full rounded-md animate-pulse" />
              <Skeleton className="h-4 w-2/3 rounded-md animate-pulse" />
              <Skeleton className="h-4 w-1/2 rounded-md animate-pulse" />
            </div>
          ) : (
            <div className="flex-1 overflow-y-auto min-h-0 bg-muted/10 border border-border/40 rounded-xl p-4">
              <TabsContent value="info" className="h-full focus-visible:outline-none">
                <DeviceInfoTab device={displayDevice} />
              </TabsContent>
              <TabsContent value="profile" className="h-full focus-visible:outline-none">
                <DeviceProfileTab device={displayDevice} configurationName={configurationName} />
              </TabsContent>
              <TabsContent value="location" className="h-full focus-visible:outline-none">
                <DeviceLocationTab device={displayDevice} />
              </TabsContent>
              <TabsContent value="apps" className="h-full focus-visible:outline-none">
                <DeviceAppsTab device={displayDevice} />
              </TabsContent>
              <TabsContent value="logs" className="h-full focus-visible:outline-none">
                <DeviceLogsTab device={displayDevice} />
              </TabsContent>
            </div>
          )}
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
