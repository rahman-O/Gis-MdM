import { Button } from '@/shared/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import type { Application } from '@/features/applications/model/types'

interface Props {
  open: boolean
  duplicates: Application[]
  onCreateNew: () => void
  onAttachVersion: (appId: number) => void
  onClose: () => void
}

export function DuplicatePackageDialog({
  open,
  duplicates,
  onCreateNew,
  onAttachVersion,
  onClose,
}: Props) {
  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Duplicate package detected</DialogTitle>
          <DialogDescription>
            An application with the same package already exists. Choose how to proceed.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-2">
          {duplicates.map((app) => (
            <div key={app.id} className="flex items-center justify-between rounded border p-2">
              <span className="text-sm">{app.name ?? app.pkg ?? `Application #${app.id}`}</span>
              <Button size="sm" variant="outline" onClick={() => app.id && onAttachVersion(app.id)}>
                Add as version
              </Button>
            </div>
          ))}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cancel</Button>
          <Button onClick={onCreateNew}>Create new app anyway</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
