import { useCallback, useEffect, useRef, useState } from 'react'
import { ChevronDown, ChevronRight, FolderPlus, MoreHorizontal, Pencil, Trash2, Check, X } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
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
  const [editingNodeId, setEditingNodeId] = useState<number | null>(null)
  const [editingName, setEditingName] = useState('')
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

  const handleRename = async (nodeId: number) => {
    const name = editingName.trim()
    if (!name) return
    try {
      await deviceTreeService.renameTreeNode(nodeId, name)
      setEditingNodeId(null)
      setEditingName('')
      await loadTree()
      onTreeChanged?.()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to rename folder.')
    }
  }

  const renderNode = (node: TreeNode, depth: number) => {
    const kids = childrenMap.get(node.id) ?? []
    const isExpanded = expanded.has(node.id)
    const isRoot = node.parentId == null
    const isSelected = selectedNodeId === node.id
    const isEditing = editingNodeId === node.id
    const hasChildren = kids.length > 0

    return (
      <div key={node.id}>
        <div
          className={cn(
            'group flex items-center gap-1 rounded px-1 py-[3px] hover:bg-muted/60 transition-colors',
            isSelected && 'bg-muted font-medium'
          )}
          style={{ paddingLeft: `${depth * 12 + 4}px` }}
        >
          {/* Expand/Collapse */}
          {hasChildren ? (
            <button
              type="button"
              className="text-muted-foreground shrink-0 p-0.5"
              onClick={() =>
                setExpanded((prev) => {
                  const next = new Set(prev)
                  if (next.has(node.id)) next.delete(node.id)
                  else next.add(node.id)
                  return next
                })
              }
            >
              {isExpanded ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronRight className="h-3.5 w-3.5" />}
            </button>
          ) : (
            <span className="w-[18px] shrink-0" />
          )}

          {/* Name or Edit Input */}
          {isEditing ? (
            <div className="flex flex-1 items-center gap-1">
              <Input
                className="h-6 text-sm px-1.5 py-0"
                value={editingName}
                onChange={(e) => setEditingName(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') void handleRename(node.id)
                  if (e.key === 'Escape') setEditingNodeId(null)
                }}
                autoFocus
              />
              <button type="button" className="text-muted-foreground hover:text-foreground p-0.5" onClick={() => void handleRename(node.id)}>
                <Check className="h-3.5 w-3.5" />
              </button>
              <button type="button" className="text-muted-foreground hover:text-foreground p-0.5" onClick={() => setEditingNodeId(null)}>
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          ) : (
            <>
              <button
                type="button"
                className="min-w-0 flex-1 truncate text-left text-sm"
                onClick={() => onSelectNode(node.id)}
              >
                {node.name}
              </button>

              {/* Device count — on every node that has children */}
              {hasChildren ? (
                <span className="text-[10px] text-muted-foreground tabular-nums shrink-0">
                  {node.deviceCount}
                </span>
              ) : null}

              {/* Context menu — Rename + Delete inside dropdown */}
              {!isRoot ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button
                      type="button"
                      className="shrink-0 p-0.5 opacity-0 group-hover:opacity-100 transition-opacity text-muted-foreground hover:text-foreground"
                      aria-label="Folder actions"
                    >
                      <MoreHorizontal className="h-3.5 w-3.5" />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-32">
                    <DropdownMenuItem
                      className="text-xs gap-2"
                      onClick={() => {
                        setEditingNodeId(node.id)
                        setEditingName(node.name)
                      }}
                    >
                      <Pencil className="h-3 w-3" />
                      Rename
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-xs gap-2 text-destructive focus:text-destructive"
                      onClick={() => setDeleteNode(node)}
                    >
                      <Trash2 className="h-3 w-3" />
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : null}
            </>
          )}
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
    <div className="flex h-full min-h-[320px] w-52 shrink-0 flex-col gap-1.5 border-r pr-2">
      {/* Header */}
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Folders</span>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-6 w-6"
          disabled={rootId == null}
          onClick={() => setNewFolderParentId(selectedNodeId ?? rootId)}
          aria-label="New folder"
        >
          <FolderPlus className="h-3.5 w-3.5" />
        </Button>
      </div>

      {/* Error */}
      {error ? <p className="text-destructive text-[10px]">{error}</p> : null}

      {/* Tree */}
      {loading ? (
        <p className="text-muted-foreground text-xs">Loading…</p>
      ) : (
        <div className="flex-1 overflow-y-auto">
          {roots.map((node) => renderNode(node, 0))}
        </div>
      )}

      {/* New folder input — compact */}
      {newFolderParentId != null ? (
        <div className="space-y-1.5 border-t pt-2">
          <Input
            className="h-7 text-xs"
            placeholder="New folder name"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') void handleCreateFolder()
              if (e.key === 'Escape') setNewFolderParentId(null)
            }}
            autoFocus
          />
          <div className="flex gap-1">
            <Button type="button" size="sm" className="h-6 flex-1 text-xs" onClick={() => void handleCreateFolder()}>
              Add
            </Button>
            <Button type="button" size="sm" variant="outline" className="h-6 text-xs" onClick={() => setNewFolderParentId(null)}>
              Cancel
            </Button>
          </div>
        </div>
      ) : null}

      {/* Delete dialog */}
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
