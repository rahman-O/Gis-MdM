import { MapPin } from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'

interface DeviceLocationTabProps {
  device: DeviceView
}

function parseLocation(location: string | null | undefined): { lat: number; lon: number; accuracy?: number } | null {
  if (!location?.trim()) return null
  // Expected format: "lat,lon" or "lat,lon,accuracy" or JSON
  try {
    const parsed = JSON.parse(location)
    if (parsed && typeof parsed.lat === 'number' && typeof parsed.lon === 'number') {
      return { lat: parsed.lat, lon: parsed.lon, accuracy: parsed.accuracy }
    }
  } catch {
    // Try comma-separated
  }
  const parts = location.split(',').map((s) => s.trim())
  if (parts.length >= 2) {
    const lat = parseFloat(parts[0])
    const lon = parseFloat(parts[1])
    const accuracy = parts[2] ? parseFloat(parts[2]) : undefined
    if (!isNaN(lat) && !isNaN(lon)) {
      return { lat, lon, accuracy: accuracy && !isNaN(accuracy) ? accuracy : undefined }
    }
  }
  return null
}

export function DeviceLocationTab({ device }: DeviceLocationTabProps) {
  const loc = parseLocation(device.info?.location)

  if (!loc) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 py-8">
        <MapPin className="h-8 w-8 text-muted-foreground" />
        <p className="text-muted-foreground text-sm">Location not available</p>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      <div className="rounded-md border p-3">
        <div className="flex items-center gap-2">
          <MapPin className="h-4 w-4 text-muted-foreground" />
          <span className="text-sm font-medium">Last Known Location</span>
        </div>
        <div className="mt-2 grid grid-cols-3 gap-2">
          <div>
            <p className="text-muted-foreground text-xs">Latitude</p>
            <p className="text-xs font-medium">{loc.lat.toFixed(6)}</p>
          </div>
          <div>
            <p className="text-muted-foreground text-xs">Longitude</p>
            <p className="text-xs font-medium">{loc.lon.toFixed(6)}</p>
          </div>
          {loc.accuracy != null && (
            <div>
              <p className="text-muted-foreground text-xs">Accuracy</p>
              <p className="text-xs font-medium">{loc.accuracy.toFixed(0)} m</p>
            </div>
          )}
        </div>
      </div>
      <p className="text-muted-foreground text-xs">
        Map view will be available in a future update.
      </p>
    </div>
  )
}
