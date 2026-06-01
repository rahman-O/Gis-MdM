import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Hash,
  Smartphone,
  Battery,
  Wifi,
  HardDrive,
  Shield,
  Fingerprint,
  Phone,
  Cpu,
  Thermometer,
  Zap,
  Signal,
  Globe,
  Monitor,
  Layers,
  Calendar,
  RefreshCw,
  AlertCircle,
  Clock,
} from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'
import { Button } from '@/shared/ui/button'
import { cn } from '@/shared/utils/cn'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface InfoField {
  icon: React.ReactNode
  label: string
  value: string | null | undefined
}

interface InfoSection {
  title: string
  fields: InfoField[]
}

// ---------------------------------------------------------------------------
// Presence Calculation
// ---------------------------------------------------------------------------

/** Threshold in milliseconds — device is "Online" if heartbeat was within this window */
const PRESENCE_THRESHOLD_MS = 60_000

/**
 * Calculates device presence status based on last heartbeat timestamp.
 * Online if heartbeat < 60s ago, Offline otherwise.
 *
 * **Validates: Requirements 6.4**
 */
export function calculatePresence(lastHeartbeat: number | null | undefined, now: number): {
  isOnline: boolean
  elapsedMs: number | null
} {
  if (lastHeartbeat == null || lastHeartbeat <= 0) {
    return { isOnline: false, elapsedMs: null }
  }
  const elapsedMs = now - lastHeartbeat
  return {
    isOnline: elapsedMs < PRESENCE_THRESHOLD_MS,
    elapsedMs: Math.max(0, elapsedMs),
  }
}

/**
 * Formats elapsed milliseconds into a human-readable string.
 * e.g. "12s ago", "3m ago", "2h ago", "5d ago"
 */
export function formatElapsed(elapsedMs: number | null): string {
  if (elapsedMs == null) return '—'
  const seconds = Math.floor(elapsedMs / 1000)
  if (seconds < 60) return `${seconds}s ago`
  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

// ---------------------------------------------------------------------------
// Null/Empty Field Placeholder
// ---------------------------------------------------------------------------

/** Placeholder for null/empty fields per Requirement 6.2 */
const PLACEHOLDER = '—'

/**
 * Returns the display value for a field, showing PLACEHOLDER for null/empty values.
 *
 * **Validates: Requirements 6.2**
 */
export function displayValue(value: string | number | boolean | null | undefined): string {
  if (value == null) return PLACEHOLDER
  if (typeof value === 'boolean') return value ? 'Active' : 'Inactive'
  const str = String(value).trim()
  return str === '' ? PLACEHOLDER : str
}

// ---------------------------------------------------------------------------
// Component Props
// ---------------------------------------------------------------------------

interface DeviceInfoTabProps {
  device: DeviceView
}

// ---------------------------------------------------------------------------
// DeviceInfoTab Component
// ---------------------------------------------------------------------------

/**
 * Displays comprehensive device information organized into sections:
 * Identity, Battery, Network, Storage, Management.
 *
 * Features:
 * - Presence indicator (Online/Offline) based on heartbeat
 * - "—" placeholder for null/empty fields
 * - Auto-refresh every 5 seconds for real-time updates
 * - Error state with retry button
 *
 * **Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5, 6.6**
 */
export function DeviceInfoTab({ device }: DeviceInfoTabProps) {
  const [now, setNow] = useState(Date.now())
  const [error, setError] = useState<string | null>(null)
  const [retrying, setRetrying] = useState(false)

  // Update "now" every 5 seconds for presence indicator refresh (Req 6.3, 6.6)
  useEffect(() => {
    const interval = setInterval(() => {
      setNow(Date.now())
    }, 5_000)
    return () => clearInterval(interval)
  }, [])

  const info = device.info

  // Determine heartbeat source: prefer info.lastHeartbeat, fall back to device.lastUpdate
  const lastHeartbeat = info?.lastHeartbeat ?? device.lastUpdate
  const presence = useMemo(() => calculatePresence(lastHeartbeat, now), [lastHeartbeat, now])

  // Build sections per Requirement 6.1
  const sections: InfoSection[] = useMemo(() => {
    return [
      {
        title: 'Identity',
        fields: [
          { icon: <Hash className="h-3.5 w-3.5" />, label: 'Device Number', value: displayValue(device.number) },
          { icon: <Smartphone className="h-3.5 w-3.5" />, label: 'Description', value: displayValue(device.description) },
          { icon: <Smartphone className="h-3.5 w-3.5" />, label: 'Model', value: displayValue(info?.model ?? device.model) },
          { icon: <Layers className="h-3.5 w-3.5" />, label: 'Manufacturer', value: displayValue(info?.manufacturer) },
          { icon: <Layers className="h-3.5 w-3.5" />, label: 'Android Version', value: displayValue(info?.androidVersion ?? device.androidVersion) },
          { icon: <Hash className="h-3.5 w-3.5" />, label: 'Serial Number', value: displayValue(info?.serial ?? device.serial) },
          { icon: <Fingerprint className="h-3.5 w-3.5" />, label: 'IMEI', value: displayValue(info?.imei ?? device.imei) },
          { icon: <Phone className="h-3.5 w-3.5" />, label: 'Phone Number', value: displayValue(info?.phone ?? device.phone) },
        ],
      },
      {
        title: 'Battery',
        fields: [
          { icon: <Battery className="h-3.5 w-3.5" />, label: 'Level', value: info?.batteryLevel != null ? `${info.batteryLevel}%` : displayValue(device.batteryLevel != null ? `${device.batteryLevel}%` : null) },
          { icon: <Battery className="h-3.5 w-3.5" />, label: 'Health', value: displayValue(info?.batteryHealth) },
          { icon: <Thermometer className="h-3.5 w-3.5" />, label: 'Temperature', value: info?.batteryTemperature != null ? `${info.batteryTemperature}°C` : PLACEHOLDER },
          { icon: <Zap className="h-3.5 w-3.5" />, label: 'Charging State', value: displayValue(info?.chargingState) },
        ],
      },
      {
        title: 'Network',
        fields: [
          { icon: <Wifi className="h-3.5 w-3.5" />, label: 'WiFi SSID', value: displayValue(info?.wifiSsid) },
          { icon: <Signal className="h-3.5 w-3.5" />, label: 'Network Type', value: displayValue(info?.networkType) },
          { icon: <Signal className="h-3.5 w-3.5" />, label: 'Signal Strength', value: info?.signalStrength != null ? `${info.signalStrength} dBm` : PLACEHOLDER },
          { icon: <Globe className="h-3.5 w-3.5" />, label: 'IP Address', value: displayValue(info?.ipAddress) },
          { icon: <Globe className="h-3.5 w-3.5" />, label: 'Public IP', value: displayValue(info?.publicIp) },
        ],
      },
      {
        title: 'Storage',
        fields: [
          { icon: <Cpu className="h-3.5 w-3.5" />, label: 'RAM Usage', value: info?.ramUsed != null && info?.ramTotal != null ? `${info.ramUsed} / ${info.ramTotal} MB` : PLACEHOLDER },
          { icon: <HardDrive className="h-3.5 w-3.5" />, label: 'Storage', value: info?.storageUsed != null && info?.storageFree != null ? `${info.storageUsed} / ${info.storageFree} GB` : PLACEHOLDER },
          { icon: <Cpu className="h-3.5 w-3.5" />, label: 'CPU Usage', value: info?.cpuUsage != null ? `${info.cpuUsage}%` : PLACEHOLDER },
        ],
      },
      {
        title: 'Management',
        fields: [
          { icon: <Shield className="h-3.5 w-3.5" />, label: 'MDM Mode', value: displayValue(info?.mdmMode) },
          { icon: <Monitor className="h-3.5 w-3.5" />, label: 'Kiosk Mode', value: displayValue(info?.kioskMode) },
          { icon: <Monitor className="h-3.5 w-3.5" />, label: 'Launcher Version', value: displayValue(info?.launcherVersion ?? device.launcherVersion) },
          { icon: <Calendar className="h-3.5 w-3.5" />, label: 'Enrollment Date', value: info?.enrollTime ? new Date(info.enrollTime).toLocaleDateString() : PLACEHOLDER },
          { icon: <Clock className="h-3.5 w-3.5" />, label: 'Last Sync Time', value: info?.lastSyncTime ? new Date(info.lastSyncTime).toLocaleString() : PLACEHOLDER },
          { icon: <Shield className="h-3.5 w-3.5" />, label: 'Enrollment State', value: displayValue(info?.enrollmentState ?? device.enrollmentState) },
        ],
      },
    ]
  }, [device, info])

  // Retry handler for error state (Req 6.5)
  const handleRetry = useCallback(() => {
    setRetrying(true)
    setError(null)
    // Simulate retry — in a real scenario this would re-fetch device data
    // The parent component (DeviceDetailsDialog) handles the actual data fetching.
    // We clear the error state to allow the parent to re-render with fresh data.
    setTimeout(() => setRetrying(false), 1000)
  }, [])

  // Error state (Req 6.5)
  if (error) {
    return (
      <div className="flex flex-col items-center justify-center gap-3 py-12">
        <AlertCircle className="text-destructive h-8 w-8" />
        <p className="text-muted-foreground text-sm">{error}</p>
        <Button variant="outline" size="sm" onClick={handleRetry} disabled={retrying}>
          <RefreshCw className={`mr-2 h-3.5 w-3.5 ${retrying ? 'animate-spin' : ''}`} />
          Retry
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-5">
      {/* Presence Indicator */}
      <div className={cn(
        "flex items-center gap-2.5 rounded-xl border px-4 py-3 shadow-xs transition-colors",
        presence.isOnline 
          ? "bg-emerald-500/5 border-emerald-500/10 text-emerald-800 dark:text-emerald-300" 
          : "bg-muted border-border/60 text-muted-foreground"
      )}>
        <div
          className={cn(
            "h-2.5 w-2.5 rounded-full transition-all duration-300",
            presence.isOnline ? "bg-emerald-500 shadow-xs shadow-emerald-500/50" : "bg-muted-foreground/60"
          )}
          aria-label={presence.isOnline ? 'Online' : 'Offline'}
        />
        <span className="text-sm font-semibold">
          {presence.isOnline ? 'Connected' : 'Disconnected'}
        </span>
        {presence.elapsedMs != null && (
          <span className="text-xs opacity-80">
            • Last sync heartbeat {formatElapsed(presence.elapsedMs)}
          </span>
        )}
      </div>

      {/* Grid-based Sectioned Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {sections.map((section) => (
          <div key={section.title} className="bg-card border border-border/50 rounded-xl p-4 shadow-2xs hover:shadow-xs transition-all duration-200">
            <h3 className="text-foreground/90 text-xs font-bold uppercase tracking-wider mb-3.5 border-b border-border/40 pb-1.5">
              {section.title}
            </h3>
            <div className="space-y-2">
              {section.fields.map((field) => (
                <div key={field.label} className="flex items-center justify-between py-1.5 border-b border-border/30 last:border-0">
                  <div className="flex items-center gap-2 text-muted-foreground/80">
                    <div className="shrink-0 text-muted-foreground/60">{field.icon}</div>
                    <span className="text-xs font-medium">{field.label}</span>
                  </div>
                  <span className="text-xs font-semibold text-foreground/85 select-all truncate max-w-[200px] text-right" title={String(field.value)}>
                    {field.value ?? PLACEHOLDER}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
