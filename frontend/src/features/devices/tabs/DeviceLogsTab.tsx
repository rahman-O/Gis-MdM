import { Clock, FileText } from 'lucide-react'
import { Badge } from '@/shared/ui/badge'
import { Separator } from '@/shared/ui/separator'
import type { DeviceView } from '@/features/devices/types'
import { formatLastSeen } from '@/features/devices/deviceFormat'

interface DeviceLogsTabProps {
  device: DeviceView
}

export function DeviceLogsTab({ device }: DeviceLogsTabProps) {
  const info = device.info

  const entries: { label: string; value: string | null }[] = [
    { label: 'Last Sync', value: device.lastUpdate ? formatLastSeen(device.lastUpdate) : null },
    {
      label: 'Enrollment Time',
      value:
        info?.enrollTime != null && info.enrollTime > 0
          ? new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
              new Date(info.enrollTime)
            )
          : null,
    },
  ]

  const visibleEntries = entries.filter((e) => e.value != null)

  return (
    <div className="space-y-4">
      {/* Timestamps & state */}
      {(visibleEntries.length > 0 || device.enrollmentState || device.statusCode) && (
        <div className="space-y-2">
          {visibleEntries.map((entry) => (
            <div key={entry.label} className="flex items-center gap-2">
              <Clock className="text-muted-foreground h-3.5 w-3.5" />
              <span className="text-muted-foreground text-xs">{entry.label}:</span>
              <span className="text-xs font-medium">{entry.value}</span>
            </div>
          ))}
          {device.enrollmentState && (
            <div className="flex items-center gap-2">
              <span className="text-muted-foreground text-xs">Enrollment State:</span>
              <Badge variant="secondary" className="text-[10px]">
                {device.enrollmentState}
              </Badge>
            </div>
          )}
          {device.statusCode && (
            <div className="flex items-center gap-2">
              <span className="text-muted-foreground text-xs">Status Code:</span>
              <Badge variant="outline" className="text-[10px]">
                {device.statusCode}
              </Badge>
            </div>
          )}
        </div>
      )}

      <Separator />

      {/* Placeholder */}
      <div className="flex flex-col items-center justify-center gap-2 py-8">
        <FileText className="text-muted-foreground h-8 w-8" />
        <p className="text-muted-foreground text-xs">
          Detailed logs will be available when the device agent reports events.
        </p>
      </div>
    </div>
  )
}
