import { useEffect, useMemo, useState } from 'react'
import { ChevronDown, ChevronRight } from 'lucide-react'
import * as CheckboxPrimitive from '@radix-ui/react-checkbox'
import { Check, Minus } from 'lucide-react'
import { cn } from '@/shared/utils/cn'
import {
  buildChildrenMap,
  collectDescendantIds,
  getFolderCheckState,
  type FolderCheckState,
} from '@/features/device-tree/treeUtils'
import type { TreeNode } from '@/features/device-tree/deviceTreeService'

interface Props {
  nodes: TreeNode[]
  rootId: number
  selectedIds: number[]
  onSelectedIdsChange: (ids: number[]) => void
  disabled?: boolean
  className?: string
}

function TreeCheckbox({
  id,
  state,
  disabled,
  onToggle,
}: {
  id: string
  state: FolderCheckState
  disabled?: boolean
  onToggle: (nextChecked: boolean) => void
}) {
  const checked = state === true ? true : state === 'indeterminate' ? 'indeterminate' : false

  return (
    <CheckboxPrimitive.Root
      id={id}
      checked={checked}
      disabled={disabled}
      onCheckedChange={(v) => {
        if (v === 'indeterminate') {
          onToggle(true)
          return
        }
        onToggle(v === true)
      }}
      className={cn(
        'peer h-4 w-4 shrink-0 rounded-sm border border-primary ring-offset-background',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
        'disabled:cursor-not-allowed disabled:opacity-50',
        'data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground',
        'data-[state=indeterminate]:bg-primary data-[state=indeterminate]:text-primary-foreground'
      )}
    >
      <CheckboxPrimitive.Indicator className="flex items-center justify-center text-current">
        {state === 'indeterminate' ? (
          <Minus className="h-3.5 w-3.5" />
        ) : (
          <Check className="h-3.5 w-3.5" />
        )}
      </CheckboxPrimitive.Indicator>
    </CheckboxPrimitive.Root>
  )
}

export function DeviceTreeFolderChecklist({
  nodes,
  rootId,
  selectedIds,
  onSelectedIdsChange,
  disabled = false,
  className,
}: Props) {
  const [expanded, setExpanded] = useState<Set<number>>(() => new Set([rootId]))
  const childrenMap = useMemo(() => buildChildrenMap(nodes), [nodes])
  const selectedSet = useMemo(() => new Set(selectedIds), [selectedIds])

  useEffect(() => {
    setExpanded((prev) => {
      const next = new Set(prev)
      next.add(rootId)
      return next
    })
  }, [rootId])

  const applyChecked = (node: TreeNode, checked: boolean) => {
    const next = new Set(selectedIds)
    const descendants = collectDescendantIds(node.id, childrenMap)
    if (checked) {
      if (node.parentId != null) {
        next.add(node.id)
      }
      for (const id of descendants) {
        next.add(id)
      }
    } else {
      next.delete(node.id)
      for (const id of descendants) {
        next.delete(id)
      }
    }
    onSelectedIdsChange([...next])
  }

  const renderNode = (node: TreeNode, depth: number) => {
    const kids = childrenMap.get(node.id) ?? []
    const isExpanded = expanded.has(node.id)
    const isRoot = node.parentId == null
    const checkState = getFolderCheckState(node.id, childrenMap, selectedSet)
    const highlighted = checkState === true || checkState === 'indeterminate'

    return (
      <div key={node.id}>
        <div
          className={cn(
            'flex items-center gap-2 rounded-md py-1 pr-1 hover:bg-muted/60',
            highlighted && 'bg-muted/40'
          )}
          style={{ paddingLeft: `${depth * 16 + 4}px` }}
        >
          {kids.length > 0 ? (
            <button
              type="button"
              className="text-muted-foreground shrink-0 p-0.5"
              aria-label={isExpanded ? 'Collapse' : 'Expand'}
              disabled={disabled}
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
          <TreeCheckbox
            id={`tree-folder-${node.id}`}
            state={isRoot ? false : checkState}
            disabled={disabled || isRoot}
            onToggle={(nextChecked) => applyChecked(node, nextChecked)}
          />
          <label
            htmlFor={`tree-folder-${node.id}`}
            className={cn(
              'min-w-0 flex-1 cursor-pointer truncate text-sm',
              (disabled || isRoot) && 'cursor-not-allowed opacity-60'
            )}
            onClick={(e) => {
              if (disabled || isRoot) return
              e.preventDefault()
              applyChecked(node, checkState !== true)
            }}
          >
            {node.name}
            <span className="text-muted-foreground ml-1 text-xs">({node.deviceCount})</span>
          </label>
        </div>
        {isExpanded ? kids.map((child) => renderNode(child, depth + 1)) : null}
      </div>
    )
  }

  const roots = childrenMap.get(null) ?? []
  return (
    <div
      className={cn(
        'h-[min(22rem,48vh)] min-h-0 shrink-0 overflow-y-auto overflow-x-hidden',
        'rounded-md border bg-background p-2 scroll-smooth [scrollbar-gutter:stable]',
        className
      )}
      role="tree"
      aria-label="Device folder tree"
    >
      {roots.length === 0 ? (
        <p className="text-sm text-muted-foreground p-2">No folders in device tree.</p>
      ) : (
        roots.map((node) => renderNode(node, 0))
      )}
    </div>
  )
}
