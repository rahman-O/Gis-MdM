import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import * as profileService from '@/features/profiles/profileService'
import type { ProfileImpact } from '@/features/profiles/profileService'

interface Props {
  open: boolean
  profileId: number
  versionId: number
  onClose: () => void
  onPublished: () => void
}

export function ProfilePublishDialog({ open, profileId, versionId, onClose, onPublished }: Props) {
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
        setError(e instanceof Error ? e.message : 'Failed to load impact.')
      })
      .finally(() => setLoading(false))
  }, [open, profileId])

  if (!open) return null

  const needsConfirm = impact?.requiresConfirmDialog ?? false

  const handlePublish = async (confirmImpact: boolean) => {
    setPublishing(true)
    setError(null)
    try {
      await profileService.publishProfileVersion(profileId, versionId, confirmImpact)
      onPublished()
      onClose()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Publish failed.')
    } finally {
      setPublishing(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div className="w-full max-w-md space-y-4 rounded-lg border bg-background p-6 shadow-lg">
        <h2 className="text-lg font-semibold">Publish profile</h2>
        {loading ? <p className="text-sm text-muted-foreground">Loading impact…</p> : null}
        {impact?.assignmentsToUpdate && impact.assignmentsToUpdate.length > 0 ? (
          <div className="max-h-40 overflow-y-auto rounded border text-sm">
            <p className="border-b bg-muted/50 px-3 py-2 font-medium">
              Folder assignments will move to the new published version
            </p>
            <ul className="divide-y">
              {impact.assignmentsToUpdate.map((a) => (
                <li key={a.assignmentId} className="px-3 py-2">
                  {a.treeNodeName} · v{a.currentVersionNumber} · {a.deviceCount} device(s)
                </li>
              ))}
            </ul>
          </div>
        ) : null}
        {impact ? (
          <p className="text-sm text-muted-foreground">
            This will affect <strong>{impact.deviceCount}</strong> devices across{' '}
            <strong>{impact.enrollmentRouteCount}</strong> enrollment routes.
            {needsConfirm ? (
              <span className="mt-2 block text-amber-700">
                At least 50 devices are affected — confirm to continue.
              </span>
            ) : null}
          </p>
        ) : null}
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline" onClick={onClose} disabled={publishing}>
            Cancel
          </Button>
          {needsConfirm ? (
            <Button type="button" disabled={publishing || loading} onClick={() => void handlePublish(true)}>
              {publishing ? 'Publishing…' : 'Confirm publish'}
            </Button>
          ) : (
            <Button type="button" disabled={publishing || loading} onClick={() => void handlePublish(false)}>
              {publishing ? 'Publishing…' : 'Publish'}
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
