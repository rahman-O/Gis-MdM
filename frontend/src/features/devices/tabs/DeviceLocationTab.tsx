import { useCallback, useEffect, useRef, useState } from 'react'
import { Clock, History, MapPin, Radio } from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'
import { getDeviceLocations, type LocationPoint } from '@/features/devices/deviceService'

interface DeviceLocationTabProps {
  device: DeviceView
}

type Mode = 'normal' | 'live' | 'history'

function parseLocation(
  location: string | null | undefined
): { lat: number; lon: number; accuracy?: number; ts?: number } | null {
  if (!location?.trim()) return null
  try {
    const parsed = JSON.parse(location)
    if (parsed && typeof parsed.lat === 'number' && typeof parsed.lon === 'number') {
      return { lat: parsed.lat, lon: parsed.lon, accuracy: parsed.accuracy, ts: parsed.ts }
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

function formatTimeAgo(ts: number | undefined): string {
  if (!ts) return 'Unknown'
  const diff = Date.now() - ts
  if (diff < 60_000) return 'Just now'
  if (diff < 3_600_000) return `${Math.floor(diff / 60_000)} min ago`
  if (diff < 86_400_000) return `${Math.floor(diff / 3_600_000)} hours ago`
  return `${Math.floor(diff / 86_400_000)} days ago`
}

function formatTimestamp(ts: number): string {
  return new Date(ts).toLocaleString()
}

function buildMapUrl(lat: number, lon: number): string {
  const bbox = `${lon - 0.005},${lat - 0.005},${lon + 0.005},${lat + 0.005}`
  return `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}&layer=mapnik&marker=${lat},${lon}`
}

export function DeviceLocationTab({ device }: DeviceLocationTabProps) {
  const [mode, setMode] = useState<Mode>('normal')
  const [livePoint, setLivePoint] = useState<LocationPoint | null>(null)
  const [historyPoints, setHistoryPoints] = useState<LocationPoint[]>([])
  const [historyLoading, setHistoryLoading] = useState(false)
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  const loc = parseLocation(device.info?.location)

  // Live tracking polling
  const pollLocation = useCallback(async () => {
    try {
      const from = Date.now() - 120_000 // last 2 minutes
      const points = await getDeviceLocations(device.id, from, undefined, 1)
      if (points.length > 0) {
        setLivePoint(points[0])
      }
    } catch {
      // Silently ignore polling errors
    }
  }, [device.id])

  useEffect(() => {
    if (mode === 'live') {
      pollLocation()
      intervalRef.current = setInterval(pollLocation, 5000)
    }
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
        intervalRef.current = null
      }
    }
  }, [mode, pollLocation])

  const handleFetchHistory = async () => {
    setHistoryLoading(true)
    try {
      const from = dateFrom ? new Date(dateFrom).getTime() : undefined
      const to = dateTo ? new Date(dateTo).getTime() : undefined
      const points = await getDeviceLocations(device.id, from, to, 500)
      setHistoryPoints(points)
    } catch {
      setHistoryPoints([])
    } finally {
      setHistoryLoading(false)
    }
  }

  // Determine which location to show on map
  const displayLoc =
    mode === 'live' && livePoint
      ? { lat: livePoint.latitude, lon: livePoint.longitude, accuracy: livePoint.accuracy, ts: livePoint.timestamp }
      : mode === 'history' && historyPoints.length > 0
        ? {
            lat: historyPoints[0].latitude,
            lon: historyPoints[0].longitude,
            accuracy: historyPoints[0].accuracy,
            ts: historyPoints[0].timestamp,
          }
        : loc

  if (!displayLoc && mode === 'normal') {
    return (
      <div className="space-y-3">
        <ModeButtons mode={mode} setMode={setMode} />
        <div className="flex flex-col items-center justify-center gap-2 py-12">
          <MapPin className="text-muted-foreground h-8 w-8" />
          <p className="text-muted-foreground text-sm">
            Location not available — device hasn&apos;t reported GPS yet.
          </p>
        </div>
      </div>
    )
  }

  const mapUrl = displayLoc ? buildMapUrl(displayLoc.lat, displayLoc.lon) : null

  return (
    <div className="space-y-3">
      {/* Mode buttons */}
      <ModeButtons mode={mode} setMode={setMode} />

      {/* Coordinates */}
      {displayLoc && (
        <div className="flex items-center gap-4">
          <MapPin className="text-muted-foreground h-4 w-4 shrink-0" />
          <div className="flex flex-wrap gap-4">
            <div>
              <p className="text-muted-foreground text-xs">Latitude</p>
              <p className="text-xs font-medium">{displayLoc.lat.toFixed(6)}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Longitude</p>
              <p className="text-xs font-medium">{displayLoc.lon.toFixed(6)}</p>
            </div>
            {displayLoc.accuracy != null && (
              <div>
                <p className="text-muted-foreground text-xs">Accuracy</p>
                <p className="text-xs font-medium">{displayLoc.accuracy.toFixed(0)} m</p>
              </div>
            )}
            <div>
              <p className="text-muted-foreground text-xs">Last updated</p>
              <p className="text-xs font-medium">{formatTimeAgo(displayLoc.ts)}</p>
            </div>
          </div>
        </div>
      )}

      {/* Live tracking indicator */}
      {mode === 'live' && (
        <div className="flex items-center gap-2 rounded-md border border-green-200 bg-green-50 px-3 py-2 dark:border-green-800 dark:bg-green-950">
          <Radio className="h-4 w-4 animate-pulse text-green-600" />
          <span className="text-xs text-green-700 dark:text-green-300">
            Live tracking active — polling every 5 seconds
          </span>
        </div>
      )}

      {/* History controls */}
      {mode === 'history' && (
        <div className="flex flex-wrap items-end gap-3 rounded-md border p-3">
          <div>
            <label className="text-muted-foreground mb-1 block text-xs">From</label>
            <input
              type="datetime-local"
              value={dateFrom}
              onChange={(e) => setDateFrom(e.target.value)}
              className="rounded border px-2 py-1 text-xs"
            />
          </div>
          <div>
            <label className="text-muted-foreground mb-1 block text-xs">To</label>
            <input
              type="datetime-local"
              value={dateTo}
              onChange={(e) => setDateTo(e.target.value)}
              className="rounded border px-2 py-1 text-xs"
            />
          </div>
          <button
            onClick={handleFetchHistory}
            disabled={historyLoading}
            className="rounded bg-blue-600 px-3 py-1 text-xs text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {historyLoading ? 'Loading...' : 'Fetch History'}
          </button>
          {historyPoints.length > 0 && (
            <span className="text-muted-foreground text-xs">{historyPoints.length} points</span>
          )}
        </div>
      )}

      {/* Map embed */}
      {mapUrl && (
        <div className="overflow-hidden rounded-md border">
          <iframe
            title="Device location"
            src={mapUrl}
            className="h-[350px] w-full border-0"
            loading="lazy"
          />
        </div>
      )}

      {/* History points list */}
      {mode === 'history' && historyPoints.length > 0 && (
        <div className="max-h-[250px] overflow-y-auto rounded-md border">
          <table className="w-full text-xs">
            <thead className="bg-muted sticky top-0">
              <tr>
                <th className="px-2 py-1 text-left font-medium">Time</th>
                <th className="px-2 py-1 text-left font-medium">Latitude</th>
                <th className="px-2 py-1 text-left font-medium">Longitude</th>
                <th className="px-2 py-1 text-left font-medium">Accuracy</th>
                <th className="px-2 py-1 text-left font-medium">Speed</th>
              </tr>
            </thead>
            <tbody>
              {historyPoints.map((p, i) => (
                <tr key={i} className="border-t">
                  <td className="px-2 py-1">{formatTimestamp(p.timestamp)}</td>
                  <td className="px-2 py-1">{p.latitude.toFixed(6)}</td>
                  <td className="px-2 py-1">{p.longitude.toFixed(6)}</td>
                  <td className="px-2 py-1">{p.accuracy.toFixed(0)} m</td>
                  <td className="px-2 py-1">{p.speed.toFixed(1)} m/s</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function ModeButtons({ mode, setMode }: { mode: Mode; setMode: (m: Mode) => void }) {
  return (
    <div className="flex gap-2">
      <button
        onClick={() => setMode('normal')}
        className={`flex items-center gap-1 rounded px-3 py-1 text-xs ${
          mode === 'normal' ? 'bg-blue-600 text-white' : 'bg-muted text-muted-foreground hover:bg-muted/80'
        }`}
      >
        <MapPin className="h-3 w-3" />
        Normal
      </button>
      <button
        onClick={() => setMode('live')}
        className={`flex items-center gap-1 rounded px-3 py-1 text-xs ${
          mode === 'live' ? 'bg-green-600 text-white' : 'bg-muted text-muted-foreground hover:bg-muted/80'
        }`}
      >
        <Clock className="h-3 w-3" />
        Live Tracking
      </button>
      <button
        onClick={() => setMode('history')}
        className={`flex items-center gap-1 rounded px-3 py-1 text-xs ${
          mode === 'history' ? 'bg-purple-600 text-white' : 'bg-muted text-muted-foreground hover:bg-muted/80'
        }`}
      >
        <History className="h-3 w-3" />
        History
      </button>
    </div>
  )
}
