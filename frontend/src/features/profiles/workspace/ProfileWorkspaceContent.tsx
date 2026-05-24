import { useEffect, useState } from 'react'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import { ProfileEditorSection } from '@/features/profiles/workspace/ProfileEditorSection'
import { ProfileOverviewSection } from '@/features/profiles/workspace/ProfileOverviewSection'
import { ProfileVersionsSection } from '@/features/profiles/workspace/ProfileVersionsSection'
import { notifyProfileWorkspace, subscribeProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceEvents'
import { ProfileRolloutStatusPanel } from '@/features/profiles/ProfileRolloutStatusPanel'
import { ProfileTreeAssignmentPanel } from '@/features/profiles/ProfileTreeAssignmentPanel'
import {
  getProfileActivity,
  getProfileSummary,
  type ProfileSummary,
} from '@/features/profiles/profileHubService'
import {
  enrichPublishedPreviewFromVersion,
  loadDraftOverviewPreview,
} from '@/features/profiles/workspace/profileOverviewPreview'
import type { OverviewPolicyPreview } from '@/features/profiles/workspace/profileOverviewPreview'
import { useProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceState'

function ActivitySection({ profileId }: { profileId: number }) {
  const [items, setItems] = useState<
    { id: number; summaryKey: string; occurredAt: string; eventType: string }[]
  >([])
  const [loading, setLoading] = useState(true)
  const [refreshGen, setRefreshGen] = useState(0)

  useEffect(() => {
    return subscribeProfileWorkspace(profileId, () => setRefreshGen((g) => g + 1))
  }, [profileId])

  useEffect(() => {
    setLoading(true)
    void getProfileActivity(profileId)
      .then((page) => setItems(page.items ?? []))
      .finally(() => setLoading(false))
  }, [profileId, refreshGen])

  if (loading) return <Skeleton className="h-32 w-full" />

  return (
    <div className="space-y-3 rounded-lg border p-4">
      <h3 className="text-sm font-semibold">Activity</h3>
      {items.length === 0 ? (
        <p className="text-sm text-muted-foreground">No activity recorded yet.</p>
      ) : (
        items.map((e) => (
          <div
            key={e.id}
            className={`border-l-2 pl-3 text-sm ${
              e.eventType === 'ProfileVersionDeleted' ? 'border-destructive' : 'border-muted'
            }`}
          >
            <p className="font-medium">
              {e.eventType === 'ProfileVersionDeleted' ? 'Version deleted' : e.eventType}
            </p>
            <p className="text-xs text-muted-foreground">
              {new Date(e.occurredAt).toLocaleString()} · {e.summaryKey}
            </p>
          </div>
        ))
      )}
    </div>
  )
}

export function ProfileWorkspaceContent() {
  const { profileId, section } = useProfileWorkspace()
  const [refreshGen, setRefreshGen] = useState(0)
  const [summary, setSummary] = useState<ProfileSummary | null>(null)
  const [draftPreview, setDraftPreview] = useState<OverviewPolicyPreview | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)

  const loadSummary = () => {
    if (profileId == null) return
    setLoading(true)
    setLoadError(null)
    void getProfileSummary(profileId)
      .then(async (s) => {
        let nextSummary = s
        const publishedId = s.publishedContext?.versionId ?? s.publishedVersionId
        if (publishedId != null && s.publishedContext) {
          const pinned = s.publishedContext.pinnedSettings
          if (pinned.appCount === 0 || !pinned.mainAppName) {
            const enriched = await enrichPublishedPreviewFromVersion(profileId, publishedId, {
              kioskMode: pinned.kioskMode,
              appCount: pinned.appCount,
              mainAppName: pinned.mainAppName,
            })
            if (enriched.appCount > 0 || enriched.mainAppName) {
              nextSummary = {
                ...s,
                publishedContext: {
                  ...s.publishedContext,
                  pinnedSettings: {
                    ...pinned,
                    appCount: enriched.appCount || pinned.appCount,
                    mainAppName: enriched.mainAppName ?? pinned.mainAppName,
                    kioskMode: enriched.kioskMode || pinned.kioskMode,
                  },
                },
              }
            }
          }
        }
        setSummary(nextSummary)
        if (nextSummary.hasUnpublishedDraft && nextSummary.draftVersionId) {
          const preview = await loadDraftOverviewPreview(profileId, nextSummary.draftVersionId)
          setDraftPreview(preview)
        } else {
          setDraftPreview(null)
        }
      })
      .catch((e: unknown) => {
        setSummary(null)
        setDraftPreview(null)
        setLoadError(e instanceof Error ? e.message : 'Failed to load profile summary.')
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    loadSummary()
    // eslint-disable-next-line react-hooks/exhaustive-deps -- reload when profile or refresh bus changes
  }, [profileId, refreshGen])

  useEffect(() => {
    if (profileId == null) return
    return subscribeProfileWorkspace(profileId, () => setRefreshGen((g) => g + 1))
  }, [profileId])

  if (profileId == null) {
    return <p className="p-4 text-sm text-muted-foreground">No profile selected.</p>
  }

  if (section === 'editor') {
    return <ProfileEditorSection profileId={profileId} />
  }

  const needsSummary =
    section === 'overview' || section === 'assignments' || section === 'versions'

  return (
    <div className="bg-muted/30 p-4">
      {needsSummary && loading ? <Skeleton className="h-48 w-full" /> : null}
      {needsSummary && !loading && loadError ? (
        <div className="flex flex-col items-start gap-3 rounded-md border border-destructive/40 bg-destructive/10 p-4">
          <div className="flex items-center gap-2 text-sm text-destructive">
            <AlertCircle className="h-4 w-4 shrink-0" />
            <span>{loadError}</span>
          </div>
          <Button type="button" variant="outline" size="sm" onClick={loadSummary}>
            Retry
          </Button>
        </div>
      ) : null}
      {section === 'overview' && !loading && summary ? (
        <ProfileOverviewSection summary={summary} draftPreview={draftPreview} />
      ) : null}
      {section === 'assignments' && (!needsSummary || (!loading && !loadError)) ? (
        <>
          <ProfileTreeAssignmentPanel
            profileId={profileId}
            publishedVersionId={summary?.publishedVersionId ?? summary?.publishedContext?.versionId}
            publishedVersionNumber={
              summary?.publishedVersionNumber ?? summary?.publishedContext?.versionNumber
            }
            hasUnpublishedDraft={summary?.hasUnpublishedDraft}
            refreshKey={refreshGen}
            onAssignmentChange={() => notifyProfileWorkspace(profileId)}
          />
        </>
      ) : null}
      {section === 'rollout' ? <ProfileRolloutStatusPanel profileId={profileId} /> : null}
      {section === 'versions' ? (
        <ProfileVersionsSection
          profileId={profileId}
          publishedVersionId={summary?.publishedVersionId ?? summary?.publishedContext?.versionId}
        />
      ) : null}
      {section === 'activity' ? <ActivitySection profileId={profileId} /> : null}
    </div>
  )
}
