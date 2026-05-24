import { cn } from '@/shared/utils/cn'
import type { ProfileHealth } from '@/features/profiles/profileHubService'

const LABELS: Record<ProfileHealth, string> = {
  healthy: 'Healthy',
  warning: 'Warning',
  error: 'Error',
  draft_only: 'Draft only',
}

const STYLES: Record<ProfileHealth, string> = {
  healthy: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950 dark:text-emerald-200',
  warning: 'bg-amber-100 text-amber-900 dark:bg-amber-950 dark:text-amber-200',
  error: 'bg-red-100 text-red-800 dark:bg-red-950 dark:text-red-200',
  draft_only: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200',
}

interface Props {
  health: ProfileHealth | string
  className?: string
}

export function ProfileHealthBadge({ health, className }: Props) {
  const key = (health in LABELS ? health : 'warning') as ProfileHealth
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium',
        STYLES[key],
        className
      )}
    >
      {LABELS[key]}
    </span>
  )
}
