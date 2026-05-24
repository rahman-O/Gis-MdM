import type { TreeNode } from '@/features/device-tree/deviceTreeService'

export type FolderCheckState = boolean | 'indeterminate'

export function buildChildrenMap(nodes: TreeNode[]): Map<number | null, TreeNode[]> {
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

/** All descendant folder ids (not including the node itself). */
export function collectDescendantIds(
  nodeId: number,
  childrenMap: Map<number | null, TreeNode[]>
): number[] {
  const kids = childrenMap.get(nodeId) ?? []
  const out: number[] = []
  for (const child of kids) {
    out.push(child.id)
    out.push(...collectDescendantIds(child.id, childrenMap))
  }
  return out
}

/** Checked, unchecked, or partial (some descendants selected). */
export function getFolderCheckState(
  nodeId: number,
  childrenMap: Map<number | null, TreeNode[]>,
  selected: Set<number>
): FolderCheckState {
  const descendants = collectDescendantIds(nodeId, childrenMap)
  if (descendants.length === 0) {
    return selected.has(nodeId)
  }
  const selectedDesc = descendants.filter((id) => selected.has(id)).length
  const selfSelected = selected.has(nodeId)
  if (selectedDesc === descendants.length && (selfSelected || selectedDesc > 0)) {
    return true
  }
  if (selectedDesc === 0 && !selfSelected) {
    return false
  }
  return 'indeterminate'
}
