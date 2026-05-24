import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import type { ProfileSummary } from '@/features/profiles/profileHubService'
import type { OverviewPolicyPreview } from '@/features/profiles/workspace/profileOverviewPreview'
import { useProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceState'

interface Props {
  summary: ProfileSummary
  draftPreview?: OverviewPolicyPreview | null
}

function cardsFromPreview(
  summary: ProfileSummary,
  preview: OverviewPolicyPreview,
  label: string
): { title: string; body: string; sub?: string }[] {
  const rollout = summary.rollout ?? {
    pending: 0,
    installed: 0,
    partial: 0,
    failed: 0,
    total: 0,
  }
  const folders = summary.assignedFolders ?? []
  return [
    {
      title: 'Status',
      body: summary.enabled ? 'Enabled' : 'Disabled',
      sub: `${label} · Lifecycle: ${summary.lifecycle}`,
    },
    {
      title: 'Assignment',
      body: `${summary.assignmentCount} folder(s)`,
      sub: folders.length ? folders.join(', ') : 'No tree assignment',
    },
    {
      title: 'Rollout',
      body: `${rollout.total} device(s) tracked`,
      sub: `Pending ${rollout.pending} · Failed ${rollout.failed}`,
    },
    {
      title: 'Apps',
      body: `${preview.appCount} application(s)`,
      sub: preview.mainAppName ? `Main: ${preview.mainAppName}` : undefined,
    },
    {
      title: 'Kiosk',
      body: preview.kioskMode ? 'Kiosk enabled' : 'Standard mode',
    },
    {
      title: 'Version',
      body: preview.versionNumber != null ? `v${preview.versionNumber}` : '—',
      sub: preview.status,
    },
  ]
}

export function ProfileOverviewSection({ summary, draftPreview }: Props) {
  const { setSection } = useProfileWorkspace()
  const published = summary.publishedContext
  const hasPublished = published != null

  const publishedPreview: OverviewPolicyPreview | null = hasPublished
    ? {
        versionId: published.versionId,
        versionNumber: published.versionNumber,
        status: published.status,
        kioskMode: published.pinnedSettings.kioskMode,
        appCount: published.pinnedSettings.appCount,
        mainAppName: published.pinnedSettings.mainAppName,
      }
    : null

  const activePreview =
    draftPreview && summary.hasUnpublishedDraft ? draftPreview : publishedPreview

  const cards =
    activePreview != null
      ? cardsFromPreview(
          summary,
          activePreview,
          draftPreview && summary.hasUnpublishedDraft
            ? `Saved draft${activePreview.versionNumber != null ? ` v${activePreview.versionNumber}` : ''}`
            : `Published${activePreview.versionNumber != null ? ` v${activePreview.versionNumber}` : ''}`
        )
      : []

  return (
    <div className="space-y-4">
      {summary.hasUnpublishedDraft ? (
        <div className="flex flex-wrap items-center justify-between gap-2 rounded-md border border-amber-500/40 bg-amber-50/80 px-3 py-2 text-sm dark:bg-amber-950/30">
          <span className="text-amber-900 dark:text-amber-100">
            Saved draft
            {draftPreview?.versionNumber != null ? ` v${draftPreview.versionNumber}` : ''}
            {hasPublished && summary.publishedVersionNumber != null
              ? ` · Published: v${summary.publishedVersionNumber}`
              : ''}
          </span>
          <Button type="button" size="sm" variant="outline" onClick={() => setSection('editor')}>
            Open editor
          </Button>
        </div>
      ) : null}

      {draftPreview && summary.hasUnpublishedDraft ? (
        <p className="text-xs text-muted-foreground">
          Cards below reflect your <strong>saved draft</strong> (last saved in editor). Publish to
          push these settings to devices and folder assignments.
        </p>
      ) : hasPublished ? (
        <p className="text-xs text-muted-foreground">
          Cards reflect published policy v{published.versionNumber}.
        </p>
      ) : (
        <p className="text-xs text-muted-foreground">
          No published version yet. Save a draft in the editor, then publish.
        </p>
      )}

      {hasPublished && draftPreview && summary.hasUnpublishedDraft ? (
        <div className="grid gap-4 lg:grid-cols-2">
          <div>
            <p className="mb-2 text-xs font-medium text-muted-foreground">Published (live)</p>
            <div className="grid gap-3 sm:grid-cols-2">
              {cardsFromPreview(summary, publishedPreview!, `Published v${published.versionNumber}`).map(
                (c) => (
                  <Card key={`pub-${c.title}`} className="opacity-90">
                    <CardHeader className="pb-2">
                      <CardTitle className="text-sm font-medium">{c.title}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-sm font-medium">{c.body}</p>
                      {c.sub ? <p className="mt-1 text-xs text-muted-foreground">{c.sub}</p> : null}
                    </CardContent>
                  </Card>
                )
              )}
            </div>
          </div>
          <div>
            <p className="mb-2 text-xs font-medium text-amber-800 dark:text-amber-200">Saved draft</p>
            <div className="grid gap-3 sm:grid-cols-2">
              {cards.map((c) => (
                <Card key={`draft-${c.title}`} className="border-amber-500/30">
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm font-medium">{c.title}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm font-medium">{c.body}</p>
                    {c.sub ? <p className="mt-1 text-xs text-muted-foreground">{c.sub}</p> : null}
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {cards.map((c) => (
            <Card key={c.title}>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">{c.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm font-medium">{c.body}</p>
                {c.sub ? <p className="mt-1 text-xs text-muted-foreground">{c.sub}</p> : null}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
