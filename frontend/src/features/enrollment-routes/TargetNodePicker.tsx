import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AlertTriangle, Check, ChevronDown, ChevronRight, Folder, Search } from 'lucide-react'
import type { TreeNodeOption } from '@/features/enrollment-routes/enrollmentRouteService'
import { listTreeNodeOptions } from '@/features/enrollment-routes/enrollmentRouteService'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Badge } from '@/shared/ui/badge'
import { cn } from '@/shared/utils/cn'

interface Props {
  /** Currently selected node ID (controlled) */
  selectedNodeId: number | ''
  /** Called when user confirms a node selection */
  onSelect: (nodeId: number, node: TreeNodeOption) => void
  /** Called when user closes the picker without selecting */
  onCancel: () => void
  /** Whether the picker panel is visible */
  open: boolean
}

/**
 * Builds a parent→children map from the flat tree-node options list.
 * Nodes with parentId=null are roots.
 */
function buildChildrenMap(nodes: TreeNodeOption[]): Map<number | null, TreeNodeOption[]> {
  const map = new Map<number | null, TreeNodeOption[]>()
  for (const node of nodes) {
    const key = node.parentId ?? null
    const list = map.get(key) ?? []
    list.push(node)
    map.set(key, list)
  }
  return map
}

/**
 * Collects all ancestor IDs for a given node (for auto-expanding the tree to show the selected node).
 */
function getAncestorIds(nodeId: number, nodesById: Map<number, TreeNodeOption>): number[] {
  const ancestors: number[] = []
  let current = nodesById.get(nodeId)
  while (current?.parentId != null) {
    ancestors.push(current.parentId)
    current = nodesById.get(current.parentId)
  }
  return ancestors
}

/**
 * Filters nodes that match the search query (by name or path).
 * Returns a Set of node IDs that match or have a descendant that matches.
 */
function getVisibleNodeIds(
  nodes: TreeNodeOption[],
  childrenMap: Map<number | null, TreeNodeOption[]>,
  query: string
): Set<number> | null {
  if (!query.trim()) return null // show all
  const lowerQuery = query.toLowerCase()
  const matchingIds = new Set<number>()

  // Find direct matches
  for (const node of nodes) {
    if (
      node.name.toLowerCase().includes(lowerQuery) ||
      node.path.toLowerCase().includes(lowerQuery)
    ) {
      matchingIds.add(node.id)
    }
  }

  // Include all ancestors of matching nodes
  const nodesById = new Map(nodes.map((n) => [n.id, n]))
  const visible = new Set(matchingIds)
  for (const id of matchingIds) {
    let current = nodesById.get(id)
    while (current?.parentId != null) {
      visible.add(current.parentId)
      current = nodesById.get(current.parentId)
    }
  }

  // Include all descendants of matching nodes (so subtrees are visible)
  const addDescendants = (parentId: number) => {
    const children = childrenMap.get(parentId) ?? []
    for (const child of children) {
      visible.add(child.id)
      addDescendants(child.id)
    }
  }
  for (const id of matchingIds) {
    addDescendants(id)
  }

  return visible
}

export function TargetNodePicker({ selectedNodeId, onSelect, onCancel, open }: Props) {
  const { t } = useTranslation()
  const [nodes, setNodes] = useState<TreeNodeOption[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [expanded, setExpanded] = useState<Set<number>>(new Set())
  const [highlightedId, setHighlightedId] = useState<number | ''>(selectedNodeId)
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    if (!open) return
    setLoading(true)
    setError(null)
    void listTreeNodeOptions()
      .then((data) => {
        setNodes(data)
        // Auto-expand root nodes and path to selected node
        const nodesById = new Map(data.map((n) => [n.id, n]))
        const roots = data.filter((n) => n.parentId == null)
        const initialExpanded = new Set(roots.map((r) => r.id))
        if (selectedNodeId && selectedNodeId > 0) {
          const ancestors = getAncestorIds(selectedNodeId, nodesById)
          for (const id of ancestors) {
            initialExpanded.add(id)
          }
        }
        setExpanded(initialExpanded)
        setHighlightedId(selectedNodeId)
      })
      .catch((e: unknown) => {
        setError(e instanceof Error ? e.message : 'Failed to load tree folders.')
      })
      .finally(() => setLoading(false))
  }, [open, selectedNodeId])

  const childrenMap = useMemo(() => buildChildrenMap(nodes), [nodes])
  const nodesById = useMemo(() => new Map(nodes.map((n) => [n.id, n])), [nodes])

  const visibleIds = useMemo(
    () => getVisibleNodeIds(nodes, childrenMap, searchQuery),
    [nodes, childrenMap, searchQuery]
  )

  const highlightedNode = highlightedId ? nodesById.get(highlightedId) ?? null : null

  const handleConfirm = () => {
    if (!highlightedNode) return
    onSelect(highlightedNode.id, highlightedNode)
  }

  const toggleExpand = (nodeId: number) => {
    setExpanded((prev) => {
      const next = new Set(prev)
      if (next.has(nodeId)) next.delete(nodeId)
      else next.add(nodeId)
      return next
    })
  }

  const renderNode = (node: TreeNodeOption, depth: number) => {
    // If filtering, skip nodes not in visible set
    if (visibleIds && !visibleIds.has(node.id)) return null

    const children = childrenMap.get(node.id) ?? []
    const visibleChildren = visibleIds
      ? children.filter((c) => visibleIds.has(c.id))
      : children
    const hasChildren = visibleChildren.length > 0
    const isExpanded = expanded.has(node.id)
    const isHighlighted = highlightedId === node.id
    const isCurrentlySelected = selectedNodeId === node.id

    return (
      <div key={node.id}>
        <div
          role="treeitem"
          aria-selected={isHighlighted}
          aria-expanded={hasChildren ? isExpanded : undefined}
          tabIndex={0}
          className={cn(
            'flex cursor-pointer items-center gap-1 rounded-md px-1 py-1.5 text-sm transition-colors hover:bg-muted/80',
            isHighlighted && 'bg-accent text-accent-foreground',
            isCurrentlySelected && !isHighlighted && 'bg-muted/50'
          )}
          style={{ paddingLeft: `${depth * 16 + 4}px` }}
          onClick={() => setHighlightedId(node.id)}
          onDoubleClick={() => {
            setHighlightedId(node.id)
            onSelect(node.id, node)
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault()
              setHighlightedId(node.id)
            }
          }}
        >
          {hasChildren ? (
            <button
              type="button"
              className="shrink-0 p-0.5 text-muted-foreground"
              aria-label={isExpanded ? 'Collapse' : 'Expand'}
              onClick={(e) => {
                e.stopPropagation()
                toggleExpand(node.id)
              }}
            >
              {isExpanded ? (
                <ChevronDown className="h-3.5 w-3.5" />
              ) : (
                <ChevronRight className="h-3.5 w-3.5" />
              )}
            </button>
          ) : (
            <span className="w-5 shrink-0" />
          )}

          <Folder className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />

          <span className="min-w-0 flex-1 truncate">{node.name}</span>

          <span className="shrink-0 text-xs text-muted-foreground">
            {node.deviceCount}
          </span>

          {node.heavilyLoaded ? (
            <AlertTriangle className="h-3.5 w-3.5 shrink-0 text-amber-500" aria-label={t('enrollmentRoute.tree.heavilyLoaded')} />
          ) : null}

          {isCurrentlySelected ? (
            <Check className="h-3.5 w-3.5 shrink-0 text-primary" />
          ) : null}
        </div>

        {isExpanded
          ? visibleChildren.map((child) => renderNode(child, depth + 1))
          : null}
      </div>
    )
  }

  if (!open) return null

  const roots = childrenMap.get(null) ?? []
  const visibleRoots = visibleIds ? roots.filter((r) => visibleIds.has(r.id)) : roots

  return (
    <div className="flex flex-col gap-3 rounded-md border bg-background p-3">
      {/* Header */}
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">{t('enrollmentRoute.form.targetFolder')}</span>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder={t('enrollmentRoute.form.selectFolder')}
          value={searchQuery}
          onChange={(e) => {
            setSearchQuery(e.target.value)
            // Auto-expand all when searching
            if (e.target.value.trim()) {
              setExpanded(new Set(nodes.map((n) => n.id)))
            }
          }}
          className="pl-9"
        />
      </div>

      {/* Error state */}
      {error ? (
        <p className="text-sm text-destructive">{error}</p>
      ) : null}

      {/* Loading state */}
      {loading ? (
        <p className="text-sm text-muted-foreground">{t('enrollmentRoute.dialog.loading')}</p>
      ) : null}

      {/* Tree list */}
      {!loading && !error ? (
        <div
          role="tree"
          aria-label={t('enrollmentRoute.form.targetFolder')}
          className="max-h-[240px] overflow-y-auto"
        >
          {visibleRoots.length === 0 ? (
            <p className="py-4 text-center text-sm text-muted-foreground">
              {searchQuery ? t('enrollmentRoute.tree.noResults') : t('enrollmentRoute.tree.empty')}
            </p>
          ) : (
            visibleRoots.map((node) => renderNode(node, 0))
          )}
        </div>
      ) : null}

      {/* Context preview: selected node path + warnings */}
      {highlightedNode ? (
        <div className="space-y-1 border-t pt-2">
          <p className="text-xs text-muted-foreground">{highlightedNode.path}</p>
          <div className="flex flex-wrap items-center gap-1.5">
            {highlightedNode.placementKind === 'inheritable' ? (
              <Badge variant="outline" className="text-amber-700 border-amber-300 dark:text-amber-400 dark:border-amber-700">
                {t('enrollmentRoute.tree.containerWarning')}
              </Badge>
            ) : null}
            {highlightedNode.heavilyLoaded ? (
              <Badge variant="outline" className="text-amber-700 border-amber-300 dark:text-amber-400 dark:border-amber-700">
                {t('enrollmentRoute.tree.heavilyLoaded')}
              </Badge>
            ) : null}
            <span className="text-xs text-muted-foreground">
              {highlightedNode.deviceCount} {highlightedNode.deviceCount === 1 ? 'device' : 'devices'}
            </span>
          </div>
        </div>
      ) : null}

      {/* Actions */}
      <div className="flex justify-end gap-2 border-t pt-2">
        <Button type="button" variant="outline" size="sm" onClick={onCancel}>
          {t('enrollmentRoute.actions.cancel')}
        </Button>
        <Button
          type="button"
          size="sm"
          disabled={!highlightedNode}
          onClick={handleConfirm}
        >
          {t('enrollmentRoute.actions.confirm', 'Confirm')}
        </Button>
      </div>
    </div>
  )
}
