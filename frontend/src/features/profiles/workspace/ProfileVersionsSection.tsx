import { useCallback, useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Skeleton } from '@/shared/ui/skeleton'
import * as profileService from '@/features/profiles/profileService'
import { deleteProfileVersion } from '@/features/profiles/profileHubService'
import {
  forkDraftFromVersion,
  listProfileVersions,
  type ProfileVersionListItem,
} from '@/features/profiles/profileRolloutService'
import { notifyProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceEvents'
import { useProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceState'

const DELETE_ERROR_MESSAGES: Record<string, string> = {
  'error.profile.version.delete.activePublished':
    'Cannot delete the version currently published for this profile.',
  'error.profile.version.delete.assigned':
    'This version is still assigned to a folder. Remove the assignment first.',
  'error.profile.version.delete.devicesTarget':
    'Devices are still targeting this version. Wait for rollout to finish or reassign.',
  'error.notfound.profile': 'Version not found.',
}

function mapDeleteError(message: string): string {
  return DELETE_ERROR_MESSAGES[message] ?? message
}

interface Props {
  profileId: number
  publishedVersionId?: number | null
}

export function ProfileVersionsSection({ profileId, publishedVersionId }: Props) {
  const { setEditorTarget } = useProfileWorkspace()
  const [versions, setVersions] = useState<ProfileVersionListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const v = await listProfileVersions(profileId)
      setVersions(v)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load versions.')
    } finally {
      setLoading(false)
    }
  }, [profileId])

  useEffect(() => {
    void load()
  }, [load])

  const handleDelete = async (v: ProfileVersionListItem) => {
    const isActivePublished =
      publishedVersionId != null && v.versionId === publishedVersionId && v.status === 'published'
    if (isActivePublished) {
      window.alert('The active published version cannot be deleted.')
      return
    }
    const label =
      v.status === 'draft'
        ? `Delete draft v${v.versionNumber}? This cannot be undone.`
        : `Delete historical published v${v.versionNumber}? Only unused versions can be removed.`
    if (!window.confirm(label)) return
    setDeletingId(v.versionId)
    try {
      await deleteProfileVersion(profileId, v.versionId)
      notifyProfileWorkspace(profileId)
      await load()
    } catch (e: unknown) {
      const raw = e instanceof Error ? e.message : 'Delete failed.'
      window.alert(mapDeleteError(raw))
    } finally {
      setDeletingId(null)
    }
  }

  const handleFork = async (versionId: number) => {
    try {
      await forkDraftFromVersion(profileId, versionId)
      notifyProfileWorkspace(profileId)
      await load()
      const meta = await profileService.getProfileMeta(profileId)
      if (meta.draftVersionId) {
        setEditorTarget(meta.draftVersionId, false)
      }
    } catch (e: unknown) {
      window.alert(e instanceof Error ? e.message : 'Fork failed.')
    }
  }

  if (loading) return <Skeleton className="h-32 w-full" />

  return (
    <Card>
      <CardHeader>
        <CardTitle>Versions</CardTitle>
        <CardDescription>Published, draft, and historical profile versions.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-2">
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        {versions.length === 0 ? (
          <p className="text-sm text-muted-foreground">No versions yet.</p>
        ) : (
          versions.map((v) => {
            const isActivePublished =
              publishedVersionId != null &&
              v.versionId === publishedVersionId &&
              v.status === 'published'
            const canDelete = !isActivePublished
            return (
              <div
                key={v.versionId}
                className="flex flex-wrap items-center justify-between gap-2 rounded-md border px-3 py-2 text-sm"
              >
                <div>
                  <span className="font-medium">Version {v.versionNumber}</span>
                  <span className="ml-2 text-muted-foreground">{v.status}</span>
                  {v.publishedAt ? (
                    <span className="ml-2 text-xs text-muted-foreground">
                      {new Date(v.publishedAt).toLocaleDateString()}
                    </span>
                  ) : null}
                  {isActivePublished ? (
                    <span className="ml-2 text-xs text-primary">(current published)</span>
                  ) : null}
                </div>
                <div className="flex flex-wrap gap-2">
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() =>
                      setEditorTarget(v.versionId, v.status === 'published' || v.status === 'archived')
                    }
                  >
                    Open in editor
                  </Button>
                  {v.status === 'published' && !isActivePublished ? (
                    <Button type="button" size="sm" variant="outline" onClick={() => void handleFork(v.versionId)}>
                      Fork draft
                    </Button>
                  ) : null}
                  {canDelete ? (
                    <Button
                      type="button"
                      size="sm"
                      variant="destructive"
                      disabled={deletingId === v.versionId}
                      onClick={() => void handleDelete(v)}
                    >
                      {deletingId === v.versionId ? 'Deleting…' : 'Delete'}
                    </Button>
                  ) : null}
                </div>
              </div>
            )
          })
        )}
      </CardContent>
    </Card>
  )
}
