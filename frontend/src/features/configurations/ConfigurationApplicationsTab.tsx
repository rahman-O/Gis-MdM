import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import type { Configuration, ConfigurationApplication } from '@/features/configurations/types'
import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmList } from '@/services/hmdmEnvelope'

interface AppVersion {
  id: number
  version: string
  versionCode: number
}

interface AppOption {
  id: number
  name: string
  latestVersionId?: number | null
  versions?: AppVersion[]
}

interface Props {
  configuration: Configuration
  applications: AppOption[]
  upgrading: boolean
  onChange: (next: Configuration) => void
  onUpgrade: (applicationId: number) => Promise<void>
}

export function ConfigurationApplicationsTab({
  configuration,
  applications,
  upgrading,
  onChange,
  onUpgrade,
}: Props) {
  const [newAppId, setNewAppId] = useState<number | null>(null)
  const [versionsMap, setVersionsMap] = useState<Map<number, AppVersion[]>>(new Map())
  const linked = Array.isArray(configuration.applications) ? configuration.applications : []

  // Stable key for the linked app IDs so the effect doesn't re-run on every render.
  const linkedAppIdsKey = useMemo(() => {
    const ids = linked.map(app => (app as Record<string, unknown>).applicationId as number ?? app.id).filter((id): id is number => id != null && id > 0)
    return [...new Set(ids)].sort((a, b) => a - b).join(',')
  }, [linked])

  useEffect(() => {
    if (!linkedAppIdsKey) return
    const uniqueIds = linkedAppIdsKey.split(',').map(Number)

    void Promise.all(
      uniqueIds.map(async (appId) => {
        try {
          const res = await apiClient.get<HmdmEnvelope<unknown[]>>(`/private/applications/${appId}/versions`)
          const versions = unwrapHmdmList(res.data, '').map((v: unknown) => {
            const rec = v as Record<string, unknown>
            return {
              id: Number(rec.id),
              version: String(rec.version ?? ''),
              versionCode: Number(rec.versionCode ?? 0),
            }
          })
          return [appId, versions] as const
        } catch {
          return [appId, []] as const
        }
      })
    ).then((results) => {
      const map = new Map<number, AppVersion[]>()
      for (const [id, versions] of results) {
        map.set(id, versions)
      }
      setVersionsMap(map)
    })
  }, [linkedAppIdsKey])

  const selectableApps = useMemo(() => {
    const linkedIds = new Set(linked.map((app) => Number(app.id ?? 0)))
    return applications.filter((app) => !linkedIds.has(app.id))
  }, [applications, linked])

  const updateApp = (index: number, patch: Partial<ConfigurationApplication>) => {
    const next = [...linked]
    next[index] = { ...next[index], ...patch }
    onChange({ ...configuration, applications: next })
  }

  const removeApp = (index: number) => {
    const next = linked.filter((_, idx) => idx !== index)
    onChange({ ...configuration, applications: next })
  }

  const addApp = () => {
    if (newAppId == null) return
    const option = applications.find((item) => item.id === newAppId)
    if (!option) return
    const latestVersionId = Number(option.latestVersionId ?? 0)
    onChange({
      ...configuration,
      applications: [
        ...linked,
        {
          id: option.id,
          applicationId: option.id,
          name: option.name,
          action: 1,
          showIcon: true,
          version: null,
          // Persist as a linked row; otherwise backend may store NULL version id and UI drops it on reload.
          ...(latestVersionId > 0
            ? { usedVersionId: latestVersionId, latestVersion: latestVersionId }
            : {}),
        },
      ],
    })
    setNewAppId(null)
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end gap-2">
        <div className="min-w-[260px] space-y-2">
          <Label>Add application</Label>
          <Select
            value={newAppId == null ? 'none' : String(newAppId)}
            onValueChange={(value) => setNewAppId(value === 'none' ? null : Number(value))}
          >
            <SelectTrigger><SelectValue placeholder="Choose app" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="none">Select application</SelectItem>
              {selectableApps.map((app) => (
                <SelectItem key={app.id} value={String(app.id)}>
                  {app.name || `Application #${app.id}`}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <Button type="button" variant="outline" onClick={addApp} disabled={newAppId == null}>
          Add
        </Button>
      </div>

      <div className="space-y-3">
        {linked.length === 0 ? (
          <p className="text-sm text-muted-foreground">No applications linked.</p>
        ) : linked.map((app, index) => (
          <div key={`${app.id ?? 'app'}-${index}`} className="rounded-md border p-3">
            <div className="mb-3 flex items-center justify-between">
              <p className="font-medium text-sm">{String(app.name ?? `Application #${app.id ?? '?'}`)}</p>
              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={upgrading || app.id == null}
                  onClick={() => app.id != null && void onUpgrade(app.id)}
                >
                  Upgrade
                </Button>
                <Button type="button" variant="destructive" size="sm" onClick={() => removeApp(index)}>
                  Remove
                </Button>
              </div>
            </div>
            <div className="grid gap-3 md:grid-cols-2">
              <div className="space-y-2">
                <Label>Action</Label>
                <Select
                  value={String(app.action ?? 1)}
                  onValueChange={(value) => updateApp(index, { action: Number(value) })}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">Install</SelectItem>
                    <SelectItem value="2">Do not install</SelectItem>
                    <SelectItem value="3">Remove</SelectItem>
                    <SelectItem value="4">Permit</SelectItem>
                    <SelectItem value="5">Prohibit</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Version</Label>
                <Select
                  value={String((app as Record<string, unknown>).usedVersionId ?? '')}
                  onValueChange={(value) => {
                    const versionId = Number(value)
                    const versions = versionsMap.get((app as Record<string, unknown>).applicationId as number ?? app.id ?? 0) ?? []
                    const selected = versions.find(v => v.id === versionId)
                    updateApp(index, {
                      usedVersionId: versionId > 0 ? versionId : null,
                      version: selected?.version ?? null,
                    } as Partial<ConfigurationApplication>)
                  }}
                >
                  <SelectTrigger><SelectValue placeholder="Select version" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="0">Latest</SelectItem>
                    {(versionsMap.get((app as Record<string, unknown>).applicationId as number ?? app.id ?? 0) ?? []).map((ver) => (
                      <SelectItem key={ver.id} value={String(ver.id)}>
                        {ver.version} (code: {ver.versionCode})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div className="mt-3 flex flex-wrap gap-4">
              <label className="flex items-center gap-2 text-sm">
                <Checkbox
                  checked={(app as Record<string, unknown>).showIcon !== false}
                  onCheckedChange={(c) => updateApp(index, { showIcon: c === true })}
                />
                Show on home screen
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox
                  checked={Boolean((app as Record<string, unknown>).skipVersionCheck)}
                  onCheckedChange={(c) => updateApp(index, { skipVersionCheck: c === true })}
                />
                Skip version check
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox
                  checked={Boolean(app.remove)}
                  onCheckedChange={(c) => updateApp(index, { remove: c === true })}
                />
                Remove
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox
                  checked={Boolean((app as Record<string, unknown>).longTap)}
                  onCheckedChange={(c) => updateApp(index, { longTap: c === true })}
                />
                Long tap
              </label>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
