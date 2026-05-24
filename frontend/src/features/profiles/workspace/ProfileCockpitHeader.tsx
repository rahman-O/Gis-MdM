import { X } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { ProfileHealthBadge } from '@/features/profiles/ProfileHealthBadge'
import { ProfileListBadges } from '@/features/profiles/ProfileListBadges'
import type { ProfileSummary } from '@/features/profiles/profileHubService'

interface Props {
  summary: ProfileSummary | null
  loading?: boolean
  onClose: () => void
  onEdit: () => void
  onPublish?: () => void
}

export function ProfileCockpitHeader({ summary, loading, onClose, onEdit, onPublish }: Props) {
  const lifecycle = summary?.lifecycle ?? '—'
  const version =
    summary?.publishedVersionNumber != null ? `v${summary.publishedVersionNumber}` : 'Unpublished'

  return (
    <header className="flex shrink-0 flex-col gap-3 border-b bg-background px-4 py-3 md:flex-row md:items-center md:justify-between">
      <div className="min-w-0 space-y-1">
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="truncate text-lg font-semibold">
            {loading ? 'Loading…' : (summary?.name ?? 'Profile')}
          </h2>
          {summary ? <ProfileHealthBadge health={summary.health} /> : null}
          <span className="rounded-md bg-muted px-2 py-0.5 text-xs capitalize text-muted-foreground">
            {lifecycle}
          </span>
          <span className="text-xs text-muted-foreground">{version}</span>
        </div>
        {summary ? (
          <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            <span>{summary.assignmentCount} folder assignment(s)</span>
            {(summary.rollout?.failed ?? 0) > 0 ? (
              <span className="text-destructive">
                {summary.rollout?.failed} rollout failure(s)
              </span>
            ) : null}
            <ProfileListBadges
              badges={
                summary.enabled === false
                  ? [...(summary.healthReasons ?? []), 'disabled']
                  : summary.healthReasons
              }
            />
          </div>
        ) : null}
      </div>
      <div className="flex shrink-0 flex-wrap items-center gap-2">
        <Button type="button" variant="outline" size="sm" onClick={onEdit}>
          Edit
        </Button>
        {onPublish ? (
          <Button
            type="button"
            variant="secondary"
            size="sm"
            disabled={!summary?.canPublish}
            onClick={onPublish}
          >
            Publish
          </Button>
        ) : null}
        <Button type="button" variant="ghost" size="icon" aria-label="Close workspace" onClick={onClose}>
          <X className="h-4 w-4" />
        </Button>
      </div>
    </header>
  )
}
