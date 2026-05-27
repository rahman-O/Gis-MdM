import { MoreHorizontal, RefreshCw, RotateCcw, Package } from 'lucide-react'
import { Badge } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import { Separator } from '@/shared/ui/separator'
import type { DeviceView } from '@/features/devices/types'

interface DeviceProfileTabProps {
  device: DeviceView
  configurationName: string | null
}

export function DeviceProfileTab({ device, configurationName }: DeviceProfileTabProps) {
  const applications = device.info?.applications ?? []
  const groups = device.groups ?? []

  return (
    <div className="space-y-4">
      {/* Profile details */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-xs">Configuration</span>
          <span className="text-xs font-medium">{configurationName || 'None'}</span>
        </div>
        <Separator />
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-xs">Enrollment State</span>
          {device.enrollmentState ? (
            <Badge variant="secondary" className="text-[10px]">
              {device.enrollmentState}
            </Badge>
          ) : (
            <span className="text-muted-foreground text-xs">—</span>
          )}
        </div>
        <Separator />
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-xs">MDM Mode</span>
          <Badge
            variant={device.info?.mdmMode ? 'default' : 'secondary'}
            className="text-[10px]"
          >
            {device.info?.mdmMode ? 'Active' : 'Inactive'}
          </Badge>
        </div>
        <Separator />
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground text-xs">Kiosk Mode</span>
          <Badge
            variant={device.info?.kioskMode ? 'default' : 'secondary'}
            className="text-[10px]"
          >
            {device.info?.kioskMode ? 'Active' : 'Inactive'}
          </Badge>
        </div>
        {groups.length > 0 && (
          <>
            <Separator />
            <div className="flex items-start justify-between">
              <span className="text-muted-foreground text-xs">Groups</span>
              <div className="flex flex-wrap justify-end gap-1">
                {groups.map((g) => (
                  <Badge key={g.id} variant="outline" className="text-[10px]">
                    {g.name || `#${g.id}`}
                  </Badge>
                ))}
              </div>
            </div>
          </>
        )}
      </div>

      {/* Actions dropdown */}
      <div className="flex justify-end">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" className="text-xs">
              <MoreHorizontal className="mr-1.5 h-3 w-3" />
              Actions
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>
              <RefreshCw className="mr-2 h-3.5 w-3.5" />
              Force Sync
            </DropdownMenuItem>
            <DropdownMenuItem>
              <RotateCcw className="mr-2 h-3.5 w-3.5" />
              Reinstall Profile
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <Separator />

      {/* Applications */}
      <div className="space-y-2">
        <div className="flex items-center gap-1.5">
          <Package className="h-3.5 w-3.5 text-muted-foreground" />
          <span className="text-xs font-medium">Applications</span>
        </div>
        {applications.length === 0 ? (
          <p className="text-muted-foreground text-xs">
            No app data — device hasn&apos;t synced yet.
          </p>
        ) : (
          <div className="space-y-1">
            {applications.map((app, index) => (
              <div
                key={`${app.pkg}-${index}`}
                className="flex items-center justify-between py-1"
              >
                <span className="min-w-0 flex-1 truncate text-xs">{app.pkg}</span>
                <div className="flex items-center gap-2">
                  <span className="text-muted-foreground text-[10px]">{app.version}</span>
                  <Badge
                    variant={app.status === 'installed' ? 'default' : 'secondary'}
                    className="text-[10px]"
                  >
                    {app.status || 'unknown'}
                  </Badge>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
