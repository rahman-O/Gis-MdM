import { ExternalLink } from 'lucide-react'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Button } from '@/shared/ui/button'

interface QrDialogProps {
  enrollmentUrl: string | null
  onClose: () => void
}

export function QrDialog({ enrollmentUrl, onClose }: QrDialogProps) {
  return (
    <Dialog open={enrollmentUrl != null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Enrollment link</DialogTitle>
        </DialogHeader>
        {enrollmentUrl ? <p className="break-all rounded border bg-muted/30 p-3 text-sm">{enrollmentUrl}</p> : null}
        {!enrollmentUrl ? <p className="text-sm text-muted-foreground">QR is not available for this device configuration.</p> : null}

        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            Close
          </Button>
          <Button
            type="button"
            onClick={() => {
              if (!enrollmentUrl) return
              window.open(enrollmentUrl, '_blank', 'noopener,noreferrer')
            }}
            disabled={!enrollmentUrl}
          >
            <ExternalLink className="mr-2 h-4 w-4" />
            Open
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
