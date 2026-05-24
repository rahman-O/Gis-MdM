const BADGE_LABELS: Record<string, string> = {
  no_assignment: 'No assignment',
  draft_changes: 'Draft changes',
  stale: 'Stale publish',
  rollout_issues: 'Rollout issues',
  disabled: 'Disabled',
  draft_only: 'Draft only',
}

interface Props {
  badges?: string[]
}

export function ProfileListBadges({ badges }: Props) {
  if (!badges?.length) return null
  return (
    <div className="flex flex-wrap gap-1">
      {badges.map((b) => (
        <span
          key={b}
          className="inline-flex rounded-md border border-border bg-muted/50 px-1.5 py-0.5 text-[10px] font-medium uppercase tracking-wide text-muted-foreground"
        >
          {BADGE_LABELS[b] ?? b.replace(/_/g, ' ')}
        </span>
      ))}
    </div>
  )
}
