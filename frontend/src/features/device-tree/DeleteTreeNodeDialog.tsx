import { useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as deviceTreeService from '@/features/device-tree/deviceTreeService'
import type { TreeNode } from '@/features/device-tree/deviceTreeService'

type Props = {
  node: TreeNode
  nodes: TreeNode[]
  open: boolean
  onOpenChange: (open: boolean) => void
  onDeleted: () => void | Promise<void>
}

export function DeleteTreeNodeDialog({ node, nodes, open, onOpenChange, onDeleted }: Props) {
  const [targetNodeId, setTargetNodeId] = useState<string>('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const targets = nodes.filter((n) => n.id !== node.id && !n.path.startsWith(node.path))

  const handleDelete = async () => {
    const target = Number(targetNodeId)
    if (!target) {
      setError('Select a target folder.')
      return
    }
    setSubmitting(true)
    setError(null)
    try {
      await deviceTreeService.deleteTreeNode(node.id, target)
      await onDeleted()
      onOpenChange(false)
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to delete folder.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete folder</DialogTitle>
          <DialogDescription>
            Devices in «{node.name}» and its subfolders will be moved to the folder you select, then this folder will be
            removed.
          </DialogDescription>
        </DialogHeader>
        <Select value={targetNodeId} onValueChange={setTargetNodeId}>
          <SelectTrigger>
            <SelectValue placeholder="Move devices to…" />
          </SelectTrigger>
          <SelectContent>
            {targets.map((t) => (
              <SelectItem key={t.id} value={String(t.id)}>
                {t.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {error ? <p className="text-destructive text-sm">{error}</p> : null}
        <DialogFooter>
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button type="button" variant="destructive" disabled={submitting} onClick={() => void handleDelete()}>
            Delete folder
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
