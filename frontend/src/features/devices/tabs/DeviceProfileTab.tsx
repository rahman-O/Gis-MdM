import { CheckCircle2, XCircle, RefreshCw, RotateCcw } from 'lucide-react'
import { Badge } from '@/shared/ui/badge'
import { Button } from '@/shared/ui/button'
import type { DeviceView } from '@/features/devices/types'

interface DeviceProfileTabProps {
  device: DeviceView
  configurationName: string | null
}

export function DeviceProfileTab({ device, configurationName }: DeviceProfileTabProps) {
  const info = device.info
  const applications = info?.applications ?? []

  return (
    <div className="space-y-4">
      {/* Profile info */}
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs">Configuration:</span>
          <span className="text-xs font-medium">{configurationName || '—'}</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs">Enrollment State:</span>
          <Badge variant="secondary" className="text-xs">
            {device.enrollmentState || '—'}
          </Badge>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs">MDM Mode:</span>
          <span className="text-xs font-medium">{info?.mdmMode ? 'Yes' : 'No'}</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs">Kiosk Mode:</span>
          <span className="text-xs font-medium">{info?.kioskMode ? 'Yes' : 'No'}</span>
        </div>
      </div>

      {/* Applications list */}
      <div className="space-y-2">
        <h4 className="text-sm font-medium">Required Applications</h4>
        {applications.length === 0 ? (
          <p className="text-muted-foreground text-xs">No applications configured.</p>
        ) : (
          <div className="space-y-1">
            {applications.map((app, index) => (
              <div key={`${app.pkg}-${index}`} className="flex items-center gap-2 rounded border p-1.5">
                {app.status === 'installed' ? (
                  <CheckCircle2 className="h-3.5 w-3.5 shrink-0 text-green-500" />
                ) : (
                  <XCircle className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                )}
                <span className="min-w-0 flex-1 truncate text-xs">{app.pkg || '—'}</span>
                <Badge variant="secondary" className="text-[10px]">
                  {app.status || 'unknown'}
                </Badge>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Action buttons */}
      <div className="flex gap-2 pt-2">
        <Button variant="outline" size="sm" className="text-xs">
          <RefreshCw className="mr-1.5 h-3 w-3" />
          Force Sync
        </Button>
        <Button variant="outline" size="sm" className="text-xs">
          <RotateCcw className="mr-1.5 h-3 w-3" />
          Reinstall Profile
        </Button>
      </div>
    </div>
  )
}
