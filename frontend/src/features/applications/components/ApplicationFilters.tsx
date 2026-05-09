import { Input } from '@/shared/ui/input'
import { Checkbox } from '@/shared/ui/checkbox'

interface Props {
  search: string
  onSearchChange: (value: string) => void
  showSystemApps: boolean
  onShowSystemAppsChange: (value: boolean) => void
  showMyAppsOnly: boolean
  onShowMyAppsOnlyChange: (value: boolean) => void
}

export function ApplicationFilters({
  search,
  onSearchChange,
  showSystemApps,
  onShowSystemAppsChange,
  showMyAppsOnly,
  onShowMyAppsOnlyChange,
}: Props) {
  return (
    <div className="space-y-3">
      <Input
        placeholder="Search applications..."
        value={search}
        onChange={(e) => onSearchChange(e.target.value)}
      />
      <div className="flex flex-wrap items-center gap-4 text-sm">
        <label className="flex items-center gap-2">
          <Checkbox checked={showSystemApps} onCheckedChange={(v) => onShowSystemAppsChange(v === true)} />
          Show system apps
        </label>
        <label className="flex items-center gap-2">
          <Checkbox checked={showMyAppsOnly} onCheckedChange={(v) => onShowMyAppsOnlyChange(v === true)} />
          Show my apps only
        </label>
      </div>
    </div>
  )
}
