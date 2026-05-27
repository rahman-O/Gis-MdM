import { useState } from 'react'
import { Search } from 'lucide-react'
import { Badge } from '@/shared/ui/badge'
import { Input } from '@/shared/ui/input'
import type { DeviceView } from '@/features/devices/types'

interface DeviceAppsTabProps {
  device: DeviceView
}

export function DeviceAppsTab({ device }: DeviceAppsTabProps) {
  const [filter, setFilter] = useState('')
  const applications = device.info?.applications ?? []

  const filtered = filter.trim()
    ? applications.filter(
        (app) =>
          app.pkg?.toLowerCase().includes(filter.toLowerCase()) ||
          app.version?.toLowerCase().includes(filter.toLowerCase())
      )
    : applications

  return (
    <div className="space-y-2">
      {/* Search */}
      <div className="relative">
        <Search className="text-muted-foreground absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2" />
        <Input
          className="h-8 pl-8 text-xs"
          placeholder="Filter apps..."
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
        />
      </div>

      {/* App list */}
      {filtered.length === 0 ? (
        <p className="text-muted-foreground py-4 text-center text-xs">
          {applications.length === 0 ? 'No applications reported.' : 'No apps match the filter.'}
        </p>
      ) : (
        <div className="max-h-[50vh] space-y-1 overflow-y-auto">
          {filtered.map((app, index) => (
            <div
              key={`${app.pkg}-${index}`}
              className="flex items-center gap-2 rounded border p-1.5"
            >
              <span className="min-w-0 flex-1 truncate text-xs font-medium">
                {app.pkg || '—'}
              </span>
              <span className="text-muted-foreground shrink-0 text-[10px]">
                {app.version || '—'}
              </span>
              <Badge variant="secondary" className="shrink-0 text-[10px]">
                {app.status || 'unknown'}
              </Badge>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
