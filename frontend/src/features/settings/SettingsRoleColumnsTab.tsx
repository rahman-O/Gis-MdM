import { useCallback, useEffect, useMemo, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { useToast } from '@/shared/hooks/use-toast'
import * as settingsService from '@/features/settings/settingsService'

const KEYS = [
  'columnDisplayedDeviceStatus',
  'columnDisplayedDeviceDate',
  'columnDisplayedDeviceNumber',
  'columnDisplayedDeviceImei',
  'columnDisplayedDevicePhone',
  'columnDisplayedDeviceModel',
  'columnDisplayedDevicePermissionsStatus',
  'columnDisplayedDeviceAppInstallStatus',
  'columnDisplayedDeviceFilesStatus',
  'columnDisplayedDeviceConfiguration',
  'columnDisplayedDeviceDesc',
  'columnDisplayedDeviceGroup',
  'columnDisplayedLauncherVersion',
  'columnDisplayedBatteryLevel',
  'columnDisplayedMdmMode',
  'columnDisplayedKioskMode',
  'columnDisplayedDefaultLauncher',
  'columnDisplayedAndroidVersion',
  'columnDisplayedEnrollmentDate',
  'columnDisplayedSerial',
  'columnDisplayedPublicIp',
  'columnDisplayedCustom1',
  'columnDisplayedCustom2',
  'columnDisplayedCustom3',
] as const

const LABELS: Record<(typeof KEYS)[number], string> = {
  columnDisplayedDeviceStatus: 'Device status',
  columnDisplayedDeviceDate: 'Last update',
  columnDisplayedDeviceNumber: 'Device number',
  columnDisplayedDeviceImei: 'IMEI',
  columnDisplayedDevicePhone: 'Phone number',
  columnDisplayedDeviceModel: 'Device model',
  columnDisplayedDevicePermissionsStatus: 'Permissions status',
  columnDisplayedDeviceAppInstallStatus: 'App install status',
  columnDisplayedDeviceFilesStatus: 'Files status',
  columnDisplayedDeviceConfiguration: 'Configuration',
  columnDisplayedDeviceDesc: 'Description',
  columnDisplayedDeviceGroup: 'Group',
  columnDisplayedLauncherVersion: 'Launcher version',
  columnDisplayedBatteryLevel: 'Battery level',
  columnDisplayedMdmMode: 'MDM mode',
  columnDisplayedKioskMode: 'Kiosk mode',
  columnDisplayedDefaultLauncher: 'Default launcher',
  columnDisplayedAndroidVersion: 'Android version',
  columnDisplayedEnrollmentDate: 'Enrollment date',
  columnDisplayedSerial: 'Serial number',
  columnDisplayedPublicIp: 'Public IP',
  columnDisplayedCustom1: '',
  columnDisplayedCustom2: '',
  columnDisplayedCustom3: '',
}

const DEFAULT_VISIBILITY: Partial<Record<(typeof KEYS)[number], boolean>> = {
  columnDisplayedDeviceStatus: true,
  columnDisplayedDeviceDate: true,
  columnDisplayedDeviceNumber: true,
  columnDisplayedDevicePermissionsStatus: true,
  columnDisplayedDeviceAppInstallStatus: true,
  columnDisplayedDeviceFilesStatus: true,
  columnDisplayedDeviceConfiguration: true,
}

function coerceBool(raw: Record<string, unknown>, key: (typeof KEYS)[number]): boolean {
  const v = raw[key]
  if (typeof v === 'boolean') return v
  if (DEFAULT_VISIBILITY[key] !== undefined) return DEFAULT_VISIBILITY[key] as boolean
  return Boolean(v ?? false)
}

export function SettingsRoleColumnsTab() {
  const { toast } = useToast()
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [roles, setRoles] = useState<settingsService.UserRoleListRow[]>([])
  const [activeRoleId, setActiveRoleId] = useState<number | null>(null)
  const [snapshots, setSnapshots] = useState<Record<number, Record<string, unknown>>>({})
  const [customerNames, setCustomerNames] = useState<{ c1?: string; c2?: string; c3?: string }>({})

  const hydrate = useCallback(async () => {
    setLoading(true)
    try {
      const meta = await settingsService.fetchRawSettings()
      setCustomerNames({
        c1: typeof meta.customPropertyName1 === 'string' ? meta.customPropertyName1 : undefined,
        c2: typeof meta.customPropertyName2 === 'string' ? meta.customPropertyName2 : undefined,
        c3: typeof meta.customPropertyName3 === 'string' ? meta.customPropertyName3 : undefined,
      })

      const list = await settingsService.listAssignableUserRoles()
      setRoles(list)
      if (!list.length) {
        setActiveRoleId(null)
        setSnapshots({})
        return
      }

      const next: Record<number, Record<string, unknown>> = {}
      for (const r of list) {
        next[r.id] = await settingsService.getUserRoleColumns(r.id)
      }

      setSnapshots(next)
      setActiveRoleId(list[0].id)
    } catch {
      toast({ variant: 'destructive', title: 'Failed to load column defaults' })
    } finally {
      setLoading(false)
    }
  }, [toast])

  useEffect(() => {
    void hydrate()
  }, [hydrate])

  const toggle = useCallback((key: (typeof KEYS)[number], checked: boolean) => {
    if (activeRoleId == null) return
    setSnapshots((prev) => {
      const row = prev[activeRoleId] ?? {}
      return { ...prev, [activeRoleId]: { ...row, [key]: checked } }
    })
  }, [activeRoleId])

  const visibleKeys = useMemo(() => {
    return KEYS.filter((k) => {
      if (k === 'columnDisplayedCustom1') return Boolean(customerNames.c1)
      if (k === 'columnDisplayedCustom2') return Boolean(customerNames.c2)
      if (k === 'columnDisplayedCustom3') return Boolean(customerNames.c3)
      return true
    })
  }, [customerNames.c1, customerNames.c2, customerNames.c3])

  async function save() {
    setBusy(true)
    try {
      const rows = roles.map((r) => snapshots[r.id]).filter(Boolean) as Record<string, unknown>[]
      if (!rows.length) return
      await settingsService.saveUserRolesCommon(rows)
      toast({ title: 'Saved role column preferences' })
    } catch {
      toast({ variant: 'destructive', title: 'Could not save common settings' })
    } finally {
      setBusy(false)
    }
  }

  function label(k: (typeof KEYS)[number]): string {
    if (k === 'columnDisplayedCustom1') return customerNames.c1 ?? 'Custom 1'
    if (k === 'columnDisplayedCustom2') return customerNames.c2 ?? 'Custom 2'
    if (k === 'columnDisplayedCustom3') return customerNames.c3 ?? 'Custom 3'
    return LABELS[k]
  }

  if (loading) {
    return (
      <div className="text-muted-foreground flex items-center gap-2 text-sm">
        <Loader2 className="h-4 w-4 animate-spin" />
        Loading device table preferences…
      </div>
    )
  }

  if (!roles.length) {
    return <p className="text-muted-foreground text-sm">No roles are available.</p>
  }

  const model = activeRoleId != null ? snapshots[activeRoleId] ?? {} : {}

  return (
    <div className="max-w-3xl space-y-4">
      <p className="text-muted-foreground text-sm">
        Mirrors the legacy “Common settings” tab: persisted through POST `/private/settings/userRoles/common` with each
        role row that was fetched for this tenant.
      </p>
      <div>
        <label className="text-muted-foreground text-xs">Role</label>
        <Select
          value={activeRoleId == null ? '' : String(activeRoleId)}
          onValueChange={(id) => setActiveRoleId(Number(id))}
        >
          <SelectTrigger>
            <SelectValue placeholder="Pick role" />
          </SelectTrigger>
          <SelectContent>
            {roles.map((r) => (
              <SelectItem key={r.id} value={String(r.id)}>
                {r.name ?? `Role #${r.id}`}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="grid gap-3 sm:grid-cols-2">
        {visibleKeys.map((k) => (
          <label key={k} className="flex items-center gap-2 text-sm">
            <Checkbox checked={coerceBool(model, k)} onCheckedChange={(c) => toggle(k, c === true)} />
            <span>{label(k)}</span>
          </label>
        ))}
      </div>
      <Button type="button" variant="outline" disabled={busy} onClick={() => void save()}>
        {busy ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
        Save columns
      </Button>
    </div>
  )
}
