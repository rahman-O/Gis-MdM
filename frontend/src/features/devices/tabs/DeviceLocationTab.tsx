import { MapPin } from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'

interface DeviceLocationTabProps {
  device: DeviceView
}

function parseLocation(
  location: string | null | undefined
): { lat: number; lon: number; accuracy?: number } | null {
  if (!location?.trim()) return null
  // Try JSON format first
  try {
    const parsed = JSON.parse(location)
    if (parsed && typeof parsed.lat === 'number' && typeof parsed.lon === 'number') {
      return { lat: parsed.lat, lon: parsed.lon, accuracy: parsed.accuracy }
    }
  } catch {
    // Not JSON, try comma-separated
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
      <div className="flex flex-col items-center justify-center gap-2 py-12">
        <MapPin className="text-muted-foreground h-8 w-8" />
        <p className="text-muted-foreground text-sm">
          Location not available — device hasn&apos;t reported GPS yet.
        </p>
      </div>
    )
  }

  const bbox = `${loc.lon - 0.01},${loc.lat - 0.01},${loc.lon + 0.01},${loc.lat + 0.01}`
  const mapUrl = `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}&layer=mapnik&marker=${loc.lat},${loc.lon}`

  return (
    <div className="space-y-3">
      {/* Coordinates */}
      <div className="flex items-center gap-4">
        <MapPin className="text-muted-foreground h-4 w-4 shrink-0" />
        <div className="flex gap-4">
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

      {/* Map embed */}
      <div className="overflow-hidden rounded-md border">
        <iframe
          title="Device location"
          src={mapUrl}
          className="h-[350px] w-full border-0"
          loading="lazy"
        />
      </div>
    </div>
  )
}
