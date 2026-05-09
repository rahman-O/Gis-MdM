import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Badge } from '@/shared/ui/badge'
import { Label } from '@/shared/ui/label'
import { Checkbox } from '@/shared/ui/checkbox'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { useDebounce } from '@/shared/hooks/useDebounce'
import type { ConfigurationOption, DeviceFilters, LookupItem } from '@/features/devices/types'

interface FilterPanelProps {
  filters: DeviceFilters
  groups: LookupItem[]
  configurations: ConfigurationOption[]
  onChange: (filters: DeviceFilters) => void
}

export function FilterPanel({ filters, groups, configurations, onChange }: FilterPanelProps) {
  const [expanded, setExpanded] = useState(false)
  const [androidVersionDraft, setAndroidVersionDraft] = useState(filters.androidVersion ?? '')
  const [launcherVersionDraft, setLauncherVersionDraft] = useState(filters.launcherVersion ?? '')
  const debouncedAndroid = useDebounce(androidVersionDraft, 300)
  const debouncedLauncher = useDebounce(launcherVersionDraft, 300)

  const activeCount = useMemo(() => {
    return Object.values(filters).filter((value) => {
      if (value == null) return false
      if (typeof value === 'string') return value.trim().length > 0
      return true
    }).length
  }, [filters])

  useEffect(() => {
    if (debouncedAndroid === filters.androidVersion && debouncedLauncher === filters.launcherVersion) return
    onChange({ ...filters, androidVersion: debouncedAndroid || null, launcherVersion: debouncedLauncher || null })
  }, [debouncedAndroid, debouncedLauncher, filters.androidVersion, filters.launcherVersion, onChange])

  return (
    <div className="rounded-md border p-3">
      <div className="flex items-center gap-2">
        <Button variant="outline" onClick={() => setExpanded((current) => !current)} type="button">
          {expanded ? 'Fewer filters' : 'More filters'}
        </Button>
        {activeCount > 0 ? <Badge variant="secondary">{activeCount}</Badge> : null}
      </div>

      {expanded ? (
        <div className="mt-3 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <div className="space-y-2">
            <Label>Group</Label>
            <Select
              value={filters.groupId == null ? 'none' : String(filters.groupId)}
              onValueChange={(value) => onChange({ ...filters, groupId: value === 'none' ? null : Number(value) })}
            >
              <SelectTrigger>
                <SelectValue placeholder="All groups" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">All groups</SelectItem>
                {groups.map((group) => (
                  <SelectItem key={group.id} value={String(group.id)}>
                    {group.name ?? `Group #${group.id}`}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Configuration</Label>
            <Select
              value={filters.configurationId == null ? 'none' : String(filters.configurationId)}
              onValueChange={(value) =>
                onChange({ ...filters, configurationId: value === 'none' ? null : Number(value) })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="All configurations" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">All configurations</SelectItem>
                {configurations.map((configuration) => (
                  <SelectItem key={configuration.id} value={String(configuration.id)}>
                    {configuration.name ?? `Configuration #${configuration.id}`}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Status</Label>
            <Select value={filters.status ?? 'all'} onValueChange={(value) => onChange({ ...filters, status: value === 'all' ? null : value })}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="green">Online</SelectItem>
                <SelectItem value="red">Offline</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Android version</Label>
            <Input value={androidVersionDraft} onChange={(event) => setAndroidVersionDraft(event.target.value)} />
          </div>

          <div className="space-y-2">
            <Label>Launcher version</Label>
            <Input value={launcherVersionDraft} onChange={(event) => setLauncherVersionDraft(event.target.value)} />
          </div>

          <div className="space-y-2">
            <Label>Sort by</Label>
            <Select value={filters.sortBy ?? 'NONE'} onValueChange={(value) => onChange({ ...filters, sortBy: value === 'NONE' ? null : value })}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="NONE">Default</SelectItem>
                <SelectItem value="LAST_UPDATE">Last seen</SelectItem>
                <SelectItem value="NUMBER">Number</SelectItem>
                <SelectItem value="ANDROID_VERSION">Android</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center gap-2 pt-6">
            <Checkbox
              checked={filters.fastSearch ?? false}
              onCheckedChange={(checked) => onChange({ ...filters, fastSearch: checked === true })}
            />
            <Label>Fast search</Label>
          </div>
        </div>
      ) : null}
    </div>
  )
}
