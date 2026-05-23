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

export interface ConfigurationDeleteTarget {
  id: number
  name: string | null | undefined
}

export interface ConfigurationDeleteDialogProps {
  configuration: ConfigurationDeleteTarget | null
  onConfirm: () => Promise<void>
  onCancel: () => void
}

export function ConfigurationDeleteDialog({
  configuration,
  onConfirm,
  onCancel,
}: ConfigurationDeleteDialogProps) {
  const [deleting, setDeleting] = useState(false)
  const [deleteError, setDeleteError] = useState<string | null>(null)
  const open = configuration != null

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
          <AlertDialogTitle>Delete configuration?</AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently remove{' '}
            <span className="text-foreground font-medium">{configuration?.name ?? ''}</span>.
            Devices referencing it may block deletion.
          </AlertDialogDescription>
        </AlertDialogHeader>
        {deleteError ?
          <p className="text-destructive text-sm" role="alert">
            {deleteError}
          </p>
        : null}
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
