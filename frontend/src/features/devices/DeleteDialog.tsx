import { useState } from 'react'
import { Loader2 } from 'lucide-react'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/ui/alert-dialog'
import { Button } from '@/shared/ui/button'
import type { DeviceView } from '@/features/devices/types'

export interface DeleteDialogProps {
  device: DeviceView | null
  onConfirm: () => Promise<void>
  onCancel: () => void
}

export function DeleteDialog({ device, onConfirm, onCancel }: DeleteDialogProps) {
  const [deleting, setDeleting] = useState(false)
  const [deleteError, setDeleteError] = useState<string | null>(null)
  const open = device != null

  const handleConfirm = async () => {
    setDeleteError(null)
    setDeleting(true)
    try {
      await onConfirm()
      onCancel()
    } catch (e) {
      setDeleteError(e instanceof Error ? e.message : 'Delete failed.')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <AlertDialog
      open={open}
      onOpenChange={(next) => {
        if (!next && !deleting) {
          setDeleteError(null)
          onCancel()
        }
      }}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete device?</AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently remove{' '}
            <span className="font-medium text-foreground">{device?.number ?? ''}</span> from the
            system. This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        {deleteError ? (
          <p className="text-sm text-destructive" role="alert">
            {deleteError}
          </p>
        ) : null}
        <AlertDialogFooter>
          <AlertDialogCancel disabled={deleting}>Cancel</AlertDialogCancel>
          <Button variant="destructive" disabled={deleting} onClick={() => void handleConfirm()}>
            {deleting ? <Loader2 className="h-4 w-4 animate-spin" aria-hidden /> : null}
            Delete
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
