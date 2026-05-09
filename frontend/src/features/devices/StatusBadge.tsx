import { Badge } from '@/shared/ui/badge'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/shared/ui/tooltip'

export interface StatusBadgeProps {
  statusCode: string | null | undefined
  lastUpdate?: number | null
}

function mapStatus(statusCode: string | null | undefined): { label: string; variant: 'default' | 'destructive' | 'secondary' } {
  switch (statusCode) {
    case 'green':
      return { label: 'Online', variant: 'default' }
    case 'red':
      return { label: 'Offline', variant: 'destructive' }
    case 'yellow':
      return { label: 'Warning', variant: 'secondary' }
    case 'brown':
      return { label: 'Inactive', variant: 'secondary' }
    default:
      return { label: 'Unknown', variant: 'secondary' }
  }
}

function humanizeSince(lastUpdate: number | null | undefined): string {
  if (!lastUpdate || lastUpdate <= 0) return 'Last seen: unavailable'
  const diffMs = Date.now() - lastUpdate
  if (diffMs < 60_000) return 'Last seen: just now'
  const minutes = Math.floor(diffMs / 60_000)
  if (minutes < 60) return `Last seen: ${minutes} minute${minutes === 1 ? '' : 's'} ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `Last seen: ${hours} hour${hours === 1 ? '' : 's'} ago`
  const days = Math.floor(hours / 24)
  return `Last seen: ${days} day${days === 1 ? '' : 's'} ago`
}

export function StatusBadge({ statusCode, lastUpdate }: StatusBadgeProps) {
  const status = mapStatus(statusCode)
  let exact = 'N/A'
  if (lastUpdate && lastUpdate > 0) {
    exact = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(lastUpdate))
  }

  return (
    <TooltipProvider delayDuration={150}>
      <Tooltip>
        <TooltipTrigger asChild>
          <span>
            <Badge variant={status.variant}>{status.label}</Badge>
          </span>
        </TooltipTrigger>
        <TooltipContent>
          <p>{humanizeSince(lastUpdate)}</p>
          <p className="text-xs text-muted-foreground">{exact}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
