import { FileText } from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'
import { formatLastSeen } from '@/features/devices/deviceFormat'

interface DeviceLogsTabProps {
  device: DeviceView
}

export function DeviceLogsTab({ device }: DeviceLogsTabProps) {
  const info = device.info

  return (
    <div className="space-y-4">
      {/* Timestamps */}
      <div className="space-y-2 rounded-md border p-2">
        <div className="flex items-center gap-2">
          <span className="text-muted-foreground text-xs">Last Sync:</span>
          <span className="text-xs font-medium">{formatLastSeen(device.lastUpdate)}</span>
        </div>
        {info?.enrollTime != null && info.enrollTime > 0 && (
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground text-xs">Enrollment Time:</span>
            <span className="text-xs font-medium">
              {new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
                new Date(info.enrollTime)
              )}
            </span>
          </div>
        )}
      </div>

      {/* Placeholder */}
      <div className="flex flex-col items-center justify-center gap-2 py-8">
        <FileText className="h-8 w-8 text-muted-foreground" />
        <p className="text-muted-foreground text-sm">Device logs will appear here</p>
      </div>
    </div>
  )
}
