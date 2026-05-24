import { useCallback, useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Label } from '@/shared/ui/label'
import { DeviceTreeFolderChecklist } from '@/features/device-tree/DeviceTreeFolderChecklist'
import { getDeviceTree, type TreeNode } from '@/features/device-tree/deviceTreeService'
import {
  getAssignmentImpact,
  listAssignments,
  putAssignment,
  type ProfileTreeAssignment,
} from '@/features/profiles/profileRolloutService'

interface Props {
  profileId: number
  publishedVersionId?: number | null
  publishedVersionNumber?: number | null
  hasUnpublishedDraft?: boolean
  /** Bump after publish/save so folder checkboxes re-sync from server assignments. */
  refreshKey?: number
  onAssignmentChange?: () => void
}

function folderIdsFromAssignments(
  assignments: ProfileTreeAssignment[],
  publishedVersionId: number | null
): number[] {
  if (publishedVersionId == null) return []
  return assignments
    .filter((a) => a.profileVersionId === publishedVersionId)
    .map((a) => a.treeNodeId)
}

export function ProfileTreeAssignmentPanel({
  profileId,
  publishedVersionId,
  publishedVersionNumber,
  hasUnpublishedDraft = false,
  refreshKey = 0,
  onAssignmentChange,
}: Props) {
  const [nodes, setNodes] = useState<TreeNode[]>([])
  const [rootId, setRootId] = useState(0)
  const [selectedFolderIds, setSelectedFolderIds] = useState<number[]>([])
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const assignVersionId =
    publishedVersionId != null && publishedVersionId > 0 ? publishedVersionId : null

  const load = useCallback(async () => {
    try {
      const [a, tree] = await Promise.all([listAssignments(profileId), getDeviceTree()])
      setNodes(tree.nodes ?? [])
      setRootId(tree.rootId)
      setSelectedFolderIds(folderIdsFromAssignments(a, assignVersionId))
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load assignment data.')
    }
  }, [profileId, assignVersionId])

  useEffect(() => {
    void load()
  }, [load, refreshKey])

  const handleAssign = async () => {
    if (!assignVersionId || selectedFolderIds.length === 0) return
    setSaving(true)
    setError(null)
    try {
      for (const treeNodeId of selectedFolderIds) {
        let confirmImpact = false
        const impact = await getAssignmentImpact(profileId, treeNodeId)
        if (impact.requiresConfirmDialog) {
          const ok = window.confirm(
            `Assign to "${impact.folderName}"? This affects ${impact.deviceCount} devices.`
          )
          if (!ok) {
            continue
          }
          confirmImpact = true
        }
        await putAssignment(profileId, {
          treeNodeId,
          profileVersionId: assignVersionId,
          confirmImpact,
        })
      }
      await load()
      onAssignmentChange?.()
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Assignment failed.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="space-y-4 rounded-lg border p-4">
      <p className="text-sm text-muted-foreground">
        Select one or more folders (parent and/or branches). Devices in each folder and its subfolders
        receive the <strong>published</strong> policy on sync.
      </p>

      {assignVersionId == null ? (
        <p className="rounded-md border border-amber-500/40 bg-amber-50/80 px-3 py-2 text-sm text-amber-900 dark:bg-amber-950/30 dark:text-amber-100">
          Publish a version before assigning folders. Draft-only changes are not pushed until you publish.
        </p>
      ) : (
        <div className="rounded-md border bg-muted/30 px-3 py-2 text-sm">
          <span className="font-medium">
            Assigning published policy
            {publishedVersionNumber != null ? ` · v${publishedVersionNumber}` : ''}
          </span>
          {hasUnpublishedDraft ? (
            <p className="mt-1 text-xs text-amber-800 dark:text-amber-200">
              You have a saved draft with unpublished changes. Publish to apply those changes to
              assignments and devices.
            </p>
          ) : null}
        </div>
      )}

      <div className="flex min-h-0 flex-col gap-2">
        <Label>Folders (tree)</Label>
        {nodes.length > 0 && rootId > 0 ? (
          <div className="min-h-0 shrink-0 overflow-hidden">
            <DeviceTreeFolderChecklist
              nodes={nodes}
              rootId={rootId}
              selectedIds={selectedFolderIds}
              onSelectedIdsChange={setSelectedFolderIds}
              disabled={assignVersionId == null || saving}
            />
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">Loading device tree…</p>
        )}
        <p className="text-xs text-muted-foreground">
          Check one or more folders. Root cannot be assigned directly.
        </p>
      </div>

      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      <Button
        type="button"
        onClick={() => void handleAssign()}
        disabled={saving || assignVersionId == null || selectedFolderIds.length === 0}
      >
        {saving
          ? 'Assigning…'
          : selectedFolderIds.length > 1
            ? `Assign to ${selectedFolderIds.length} folders`
            : 'Assign to selected folder(s)'}
      </Button>
    </div>
  )
}
