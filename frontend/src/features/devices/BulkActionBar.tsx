import { Button } from '@/shared/ui/button'

interface BulkActionBarProps {
  selectedCount: number
  onDeleteSelected: () => void
  onSetConfiguration: () => void
  onSetGroup: () => void
}

export function BulkActionBar({
  selectedCount,
  onDeleteSelected,
  onSetConfiguration,
  onSetGroup,
}: BulkActionBarProps) {
  return (
    <div className="flex flex-wrap items-center gap-2 rounded-md border bg-muted/30 p-3">
      <span className="text-sm font-medium">{selectedCount} selected</span>
      <Button type="button" variant="destructive" size="sm" onClick={onDeleteSelected}>
        Delete Selected
      </Button>
      <Button type="button" variant="outline" size="sm" onClick={onSetConfiguration}>
        Set Configuration
      </Button>
      <Button type="button" variant="outline" size="sm" onClick={onSetGroup}>
        Set Group
      </Button>
    </div>
  )
}
