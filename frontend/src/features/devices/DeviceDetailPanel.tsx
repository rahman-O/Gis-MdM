import { useEffect, useState } from 'react'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/shared/ui/sheet'
import { Skeleton } from '@/shared/ui/skeleton'
import * as deviceService from '@/features/devices/deviceService'
import type { DeviceView } from '@/features/devices/types'
import { StatusBadge } from '@/features/devices/StatusBadge'
import { formatLastSeen } from '@/features/devices/deviceFormat'

export interface DeviceDetailPanelProps {
  deviceNumber: string | null
  onClose: () => void
}

export function DeviceDetailPanel({ deviceNumber, onClose }: DeviceDetailPanelProps) {
  const [device, setDevice] = useState<DeviceView | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!deviceNumber) {
      setDevice(null)
      setError(null)
      return
    }
    let cancelled = false
    setLoading(true)
    setError(null)
    setDevice(null)
    deviceService
      .getDevice(deviceNumber)
      .then((d) => {
        if (!cancelled) setDevice(d)
      })
      .catch((e: unknown) => {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load device.')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [deviceNumber])

  const open = deviceNumber != null

  const groupsLabel =
    device?.groups?.length ?
      device.groups.map((g) => g.name ?? `#${g.id}`).join(', ')
    : '—'

  const battery = device?.info?.batteryLevel != null ? `${device.info.batteryLevel}%` : 'N/A'

  return (
    <Sheet open={open} onOpenChange={(o) => !o && onClose()}>
      <SheetContent side="right" className="flex w-full flex-col gap-4 overflow-y-auto sm:max-w-lg">
        <SheetHeader>
          <SheetTitle>Device details</SheetTitle>
        </SheetHeader>
        {loading ?
          <div className="space-y-3">
            <Skeleton className="h-8 w-3/4" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-2/3" />
          </div>
        : error ?
          <p className="text-sm text-destructive">{error}</p>
        : device ?
          <dl className="grid gap-3 text-sm">
            <div>
              <dt className="text-muted-foreground">Device number</dt>
              <dd className="font-medium">{device.number ?? '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Status</dt>
              <dd>
                <StatusBadge statusCode={device.statusCode} lastUpdate={device.lastUpdate} />
              </dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Configuration</dt>
              <dd>{device.configurationId != null ? `Configuration #${device.configurationId}` : '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Groups</dt>
              <dd>{groupsLabel}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Battery</dt>
              <dd>{battery}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Last update</dt>
              <dd>{formatLastSeen(device.lastUpdate)}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Location</dt>
              <dd>{device.info?.location?.trim() || 'Location unavailable'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Model</dt>
              <dd>{device.info?.model || device.model || '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">IMEI</dt>
              <dd>{device.info?.imei || device.imei || '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Phone</dt>
              <dd>{device.info?.phone || device.phone || '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Android version</dt>
              <dd>{device.info?.androidVersion || device.androidVersion || '—'}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Launcher version</dt>
              <dd>{device.info?.launcherVersion || device.launcherVersion || '—'}</dd>
            </div>
          </dl>
        : null}
      </SheetContent>
    </Sheet>
  )
}
