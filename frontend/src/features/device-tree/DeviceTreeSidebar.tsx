import { useCallback, useEffect, useRef, useState } from 'react'
import { ChevronDown, ChevronRight, FolderPlus, Trash2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { cn } from '@/shared/utils/cn'
import * as deviceTreeService from '@/features/device-tree/deviceTreeService'
import type { TreeNode } from '@/features/device-tree/deviceTreeService'
import { DeleteTreeNodeDialog } from '@/features/device-tree/DeleteTreeNodeDialog'

type Props = {
  selectedNodeId: number | null
  onSelectNode: (nodeId: number | null) => void
  onTreeChanged?: () => void
}

function buildChildrenMap(nodes: TreeNode[]): Map<number | null, TreeNode[]> {
  const map = new Map<number | null, TreeNode[]>()
  for (const node of nodes) {
    const key = node.parentId
    const list = map.get(key) ?? []
    list.push(node)
    map.set(key, list)
  }
  for (const list of map.values()) {
    list.sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name))
  }
  return map
}

export function DeviceTreeSidebar({ selectedNodeId, onSelectNode, onTreeChanged }: Props) {
  const [nodes, setNodes] = useState<TreeNode[]>([])
  const [rootId, setRootId] = useState<number | null>(null)
  const [expanded, setExpanded] = useState<Set<number>>(new Set())
  const [loading, setLoading] = useState(true)
  const [newFolderParentId, setNewFolderParentId] = useState<number | null>(null)
  const [newFolderName, setNewFolderName] = useState('')
  const [deleteNode, setDeleteNode] = useState<TreeNode | null>(null)
  const [error, setError] = useState<string | null>(null)
  const didAutoSelect = useRef(false)

  const loadTree = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await deviceTreeService.getDeviceTree()
      setNodes(data.nodes ?? [])
      setRootId(data.rootId)
      setExpanded((prev) => {
        const next = new Set(prev)
        next.add(data.rootId)
        return next
      })
      if (!didAutoSelect.current) {
        didAutoSelect.current = true
        onSelectNode(data.rootId)
      }
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load tree.')
    } finally {
      setLoading(false)
    }
  }, [onSelectNode])

  useEffect(() => {
    void loadTree()
  }, [loadTree])

  const childrenMap = buildChildrenMap(nodes)

  const renderNode = (node: TreeNode, depth: number) => {
    const kids = childrenMap.get(node.id) ?? []
    const isExpanded = expanded.has(node.id)
    const isRoot = node.parentId == null
    const isSelected = selectedNodeId === node.id

    return (
      <div key={node.id}>
        <div
          className={cn(
            'flex items-center gap-1 rounded-md px-1 py-1 text-sm hover:bg-muted/80',
            isSelected && 'bg-muted font-medium'
          )}
          style={{ paddingLeft: `${depth * 12 + 4}px` }}
        >
          {kids.length > 0 ? (
            <button
              type="button"
              className="text-muted-foreground shrink-0 p-0.5"
              aria-label={isExpanded ? 'Collapse' : 'Expand'}
              onClick={() =>
                setExpanded((prev) => {
                  const next = new Set(prev)
                  if (next.has(node.id)) next.delete(node.id)
                  else next.add(node.id)
                  return next
                })
              }
            >
              {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            </button>
          ) : (
            <span className="w-5 shrink-0" />
          )}
          <button type="button" className="min-w-0 flex-1 truncate text-left" onClick={() => onSelectNode(node.id)}>
            {node.name}
            <span className="text-muted-foreground ml-1">({node.deviceCount})</span>
          </button>
          {!isRoot ? (
            <button
              type="button"
              className="text-muted-foreground hover:text-destructive shrink-0 p-1"
              aria-label="Delete folder"
              onClick={() => setDeleteNode(node)}
            >
              <Trash2 className="h-3.5 w-3.5" />
            </button>
          ) : null}
        </div>
        {isExpanded ? kids.map((child) => renderNode(child, depth + 1)) : null}
      </div>
    )
  }

  const handleCreateFolder = async () => {
    const parentId = newFolderParentId ?? rootId
    const name = newFolderName.trim()
    if (!parentId || !name) return
    try {
      await deviceTreeService.createTreeNode(parentId, name)
      setNewFolderName('')
      setNewFolderParentId(null)
      await loadTree()
      onTreeChanged?.()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to create folder.')
    }
  }

  const roots = childrenMap.get(null) ?? []

  return (
    <div className="flex h-full min-h-[320px] w-56 shrink-0 flex-col gap-2 border-r pr-3">
      <div className="flex items-center justify-between gap-1">
        <span className="text-sm font-medium">Folders</span>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          disabled={rootId == null}
          onClick={() => setNewFolderParentId(selectedNodeId ?? rootId)}
          aria-label="New folder"
        >
          <FolderPlus className="h-4 w-4" />
        </Button>
      </div>
      {error ? <p className="text-destructive text-xs">{error}</p> : null}
      {loading ? <p className="text-muted-foreground text-xs">Loading…</p> : roots.map((node) => renderNode(node, 0))}
      {newFolderParentId != null ? (
        <div className="space-y-2 border-t pt-2">
          <Input
            placeholder="Folder name"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') void handleCreateFolder()
            }}
          />
          <div className="flex gap-1">
            <Button type="button" size="sm" className="flex-1" onClick={() => void handleCreateFolder()}>
              Add
            </Button>
            <Button type="button" size="sm" variant="outline" onClick={() => setNewFolderParentId(null)}>
              Cancel
            </Button>
          </div>
        </div>
      ) : null}
      {deleteNode ? (
        <DeleteTreeNodeDialog
          node={deleteNode}
          nodes={nodes}
          open
          onOpenChange={(open) => {
            if (!open) setDeleteNode(null)
          }}
          onDeleted={async () => {
            setDeleteNode(null)
            await loadTree()
            if (selectedNodeId === deleteNode.id && rootId != null) {
              onSelectNode(rootId)
            }
            onTreeChanged?.()
          }}
        />
      ) : null}
    </div>
  )
}
