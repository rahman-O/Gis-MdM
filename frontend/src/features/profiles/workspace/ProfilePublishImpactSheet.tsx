import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/shared/ui/sheet'
import { useToast } from '@/shared/hooks/use-toast'
import * as profileService from '@/features/profiles/profileService'
import type { ProfileImpact } from '@/features/profiles/profileService'
import { notifyProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceEvents'
import {
  getProfileEditorActiveVersionId,
  saveProfileEditorIfDirty,
} from '@/features/profiles/workspace/profileEditorBridge'

interface Props {
  open: boolean
  profileId: number
  versionId: number
  draftVersionNumber?: number | null
  onClose: () => void
  onPublished?: () => void
}

export function ProfilePublishImpactSheet({
  open,
  profileId,
  versionId,
  draftVersionNumber,
  onClose,
  onPublished,
}: Props) {
  const { toast } = useToast()
  const [impact, setImpact] = useState<ProfileImpact | null>(null)
  const [loading, setLoading] = useState(false)
  const [publishing, setPublishing] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!open) return
    setLoading(true)
    setError(null)
    void profileService
      .getProfileImpact(profileId)
      .then(setImpact)
      .catch((e: unknown) => {
        setError(e instanceof Error ? e.message : 'Failed to load publish impact.')
      })
      .finally(() => setLoading(false))
  }, [open, profileId])

  const needsConfirm = impact?.requiresConfirmDialog ?? false
  const nextVersion = draftVersionNumber != null ? draftVersionNumber : 'new'

  const handlePublish = async () => {
    setPublishing(true)
    setError(null)
    try {
      const saved = await saveProfileEditorIfDirty()
      if (!saved) {
        setError('Save the draft (fix any validation errors) before publishing.')
        return
      }
      const publishVersionId = getProfileEditorActiveVersionId() ?? versionId
      const result = await profileService.publishProfileVersion(profileId, publishVersionId, true)
      notifyProfileWorkspace(profileId)
      toast({
        title: 'Profile published',
        description:
          result.assignmentsUpdated != null && result.assignmentsUpdated > 0
            ? `Published v${result.versionNumber}; ${result.assignmentsUpdated} folder assignment(s) updated.`
            : `Published v${result.versionNumber}.`,
      })
      onPublished?.()
      onClose()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Publish failed.')
    } finally {
      setPublishing(false)
    }
  }

  return (
    <Sheet open={open} onOpenChange={(v) => !v && onClose()}>
      <SheetContent side="right" className="flex w-full flex-col sm:max-w-md">
        <SheetHeader>
          <SheetTitle>Publish profile</SheetTitle>
          <SheetDescription>
            Review impact before publishing draft v{nextVersion} to production.
          </SheetDescription>
        </SheetHeader>
        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto py-4">
          {loading ? <p className="text-sm text-muted-foreground">Loading impact…</p> : null}
          {error ? <p className="text-sm text-destructive">{error}</p> : null}
          {impact?.assignmentsToUpdate && impact.assignmentsToUpdate.length > 0 ? (
            <div className="rounded-md border text-sm">
              <p className="border-b bg-muted/50 px-3 py-2 font-medium">
                Folder assignments will update to v{nextVersion}
              </p>
              <table className="w-full text-left text-xs">
                <thead>
                  <tr className="border-b text-muted-foreground">
                    <th className="px-3 py-2">Folder</th>
                    <th className="px-3 py-2">From</th>
                    <th className="px-3 py-2">Devices</th>
                  </tr>
                </thead>
                <tbody>
                  {impact.assignmentsToUpdate.map((a) => (
                    <tr key={a.assignmentId} className="border-b last:border-0">
                      <td className="px-3 py-2">{a.treeNodeName}</td>
                      <td className="px-3 py-2">v{a.currentVersionNumber}</td>
                      <td className="px-3 py-2">{a.deviceCount}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : null}
          {impact ? (
            <p className="text-sm text-muted-foreground">
              <strong>{impact.deviceCount}</strong> device(s) and{' '}
              <strong>{impact.enrollmentRouteCount}</strong> enrollment route(s) may be affected.
              {needsConfirm
                ? ' Confirm to publish and bump all folder assignments to the new version.'
                : null}
            </p>
          ) : null}
        </div>
        <div className="flex gap-2 border-t pt-4">
          <Button type="button" variant="outline" className="flex-1" onClick={onClose} disabled={publishing}>
            Cancel
          </Button>
          <Button
            type="button"
            className="flex-1"
            disabled={publishing || loading}
            onClick={() => void handlePublish()}
          >
            {publishing ? 'Publishing…' : needsConfirm ? 'Confirm & publish' : 'Publish'}
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  )
}
