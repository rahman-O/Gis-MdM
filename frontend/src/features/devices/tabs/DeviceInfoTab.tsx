import {
  Smartphone,
  Battery,
  Hash,
  Clock,
  Phone,
  Globe,
  Shield,
  Monitor,
  Layers,
  FileText,
  Fingerprint,
} from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'
import { formatLastSeen } from '@/features/devices/deviceFormat'

interface DeviceInfoTabProps {
  device: DeviceView
}

interface InfoItem {
  icon: React.ReactNode
  label: string
  value: string | null | undefined
}

export function DeviceInfoTab({ device }: DeviceInfoTabProps) {
  const items: InfoItem[] = [
    { icon: <Hash className="h-3.5 w-3.5" />, label: 'Device Number', value: device.number },
    { icon: <Smartphone className="h-3.5 w-3.5" />, label: 'Model', value: device.model },
    { icon: <Layers className="h-3.5 w-3.5" />, label: 'Android', value: device.androidVersion },
    {
      icon: <Battery className="h-3.5 w-3.5" />,
      label: 'Battery',
      value: device.batteryLevel != null ? `${device.batteryLevel}%` : null,
    },
    { icon: <Fingerprint className="h-3.5 w-3.5" />, label: 'IMEI', value: device.imei },
    { icon: <Hash className="h-3.5 w-3.5" />, label: 'Serial', value: device.serial },
    { icon: <Phone className="h-3.5 w-3.5" />, label: 'Phone', value: device.phone },
    { icon: <Monitor className="h-3.5 w-3.5" />, label: 'Launcher', value: device.launcherVersion },
    {
      icon: <Clock className="h-3.5 w-3.5" />,
      label: 'Last Seen',
      value: device.lastUpdate ? formatLastSeen(device.lastUpdate) : null,
    },
    { icon: <FileText className="h-3.5 w-3.5" />, label: 'Description', value: device.description },
    { icon: <Globe className="h-3.5 w-3.5" />, label: 'Public IP', value: device.info?.publicIp },
    {
      icon: <Shield className="h-3.5 w-3.5" />,
      label: 'MDM Mode',
      value: device.info?.mdmMode != null ? (device.info.mdmMode ? 'Active' : 'Inactive') : null,
    },
    {
      icon: <Monitor className="h-3.5 w-3.5" />,
      label: 'Kiosk Mode',
      value: device.info?.kioskMode != null ? (device.info.kioskMode ? 'Active' : 'Inactive') : null,
    },
    {
      icon: <Clock className="h-3.5 w-3.5" />,
      label: 'Enrolled',
      value: device.info?.enrollTime ? new Date(device.info.enrollTime).toLocaleDateString() : null,
    },
  ]

  const visibleItems = items.filter(
    (item) => item.value != null && item.value !== '' && item.value !== '—'
  )

  if (visibleItems.length === 0) {
    return (
      <p className="text-muted-foreground py-8 text-center text-xs">
        No device information available yet.
      </p>
    )
  }

  return (
    <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
      {visibleItems.map((item) => (
        <div key={item.label} className="flex items-start gap-2 rounded-md p-2">
          <div className="text-muted-foreground mt-0.5">{item.icon}</div>
          <div className="min-w-0 flex-1">
            <p className="text-muted-foreground text-xs">{item.label}</p>
            <p className="truncate text-xs font-medium">{item.value}</p>
          </div>
        </div>
      ))}
    </div>
  )
}
