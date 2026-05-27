import {
  Smartphone,
  Battery,
  Wifi,
  HardDrive,
  MemoryStick,
  Hash,
  Clock,
  Eye,
  Cpu,
  Layers,
} from 'lucide-react'
import type { DeviceView } from '@/features/devices/types'
import { formatLastSeen } from '@/features/devices/deviceFormat'

interface DeviceInfoTabProps {
  device: DeviceView
}

interface InfoCardProps {
  icon: React.ReactNode
  label: string
  value: string
}

function InfoCard({ icon, label, value }: InfoCardProps) {
  return (
    <div className="flex items-start gap-2 rounded-md border p-2">
      <div className="text-muted-foreground mt-0.5">{icon}</div>
      <div className="min-w-0 flex-1">
        <p className="text-muted-foreground text-xs">{label}</p>
        <p className="truncate text-xs font-medium">{value}</p>
      </div>
    </div>
  )
}

export function DeviceInfoTab({ device }: DeviceInfoTabProps) {
  const info = device.info

  const cards: InfoCardProps[] = [
    {
      icon: <Smartphone className="h-3.5 w-3.5" />,
      label: 'Model',
      value: info?.model || device.model || '—',
    },
    {
      icon: <Layers className="h-3.5 w-3.5" />,
      label: 'Android',
      value: info?.androidVersion || device.androidVersion || '—',
    },
    {
      icon: <Battery className="h-3.5 w-3.5" />,
      label: 'Battery',
      value: info?.batteryLevel != null ? `${info.batteryLevel}%` : device.batteryLevel != null ? `${device.batteryLevel}%` : '—',
    },
    {
      icon: <Wifi className="h-3.5 w-3.5" />,
      label: 'Network / IP',
      value: info?.publicIp || '—',
    },
    {
      icon: <HardDrive className="h-3.5 w-3.5" />,
      label: 'Storage',
      value: '—',
    },
    {
      icon: <MemoryStick className="h-3.5 w-3.5" />,
      label: 'RAM',
      value: '—',
    },
    {
      icon: <Hash className="h-3.5 w-3.5" />,
      label: 'IMEI',
      value: info?.imei || device.imei || '—',
    },
    {
      icon: <Hash className="h-3.5 w-3.5" />,
      label: 'Serial',
      value: info?.serial || device.serial || '—',
    },
    {
      icon: <Eye className="h-3.5 w-3.5" />,
      label: 'Last Seen',
      value: formatLastSeen(device.lastUpdate),
    },
    {
      icon: <Clock className="h-3.5 w-3.5" />,
      label: 'Uptime',
      value: '—',
    },
    {
      icon: <Cpu className="h-3.5 w-3.5" />,
      label: 'Launcher Version',
      value: info?.launcherVersion || device.launcherVersion || '—',
    },
  ]

  return (
    <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-4">
      {cards.map((card) => (
        <InfoCard key={card.label} icon={card.icon} label={card.label} value={card.value} />
      ))}
    </div>
  )
}
