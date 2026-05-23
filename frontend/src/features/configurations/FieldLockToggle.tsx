import { Lock, LockOpen } from 'lucide-react'
import { Button } from '@/shared/ui/button'

interface FieldLockToggleProps {
  fieldKey: string
  locked: boolean
  onToggle: (fieldKey: string, locked: boolean) => void
  disabled?: boolean
}

export function FieldLockToggle({ fieldKey, locked, onToggle, disabled }: FieldLockToggleProps) {
  return (
    <Button
      type="button"
      variant="ghost"
      size="icon"
      className="h-8 w-8 shrink-0"
      disabled={disabled}
      aria-label={locked ? `Unlock ${fieldKey}` : `Lock ${fieldKey}`}
      aria-pressed={locked}
      onClick={() => onToggle(fieldKey, !locked)}
    >
      {locked ? <Lock className="h-4 w-4" /> : <LockOpen className="h-4 w-4 text-muted-foreground" />}
    </Button>
  )
}
