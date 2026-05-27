import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface TreeNode {
  id: number
  parentId: number | null
  name: string
  sortOrder: number
  path: string
  depth: number
  deviceCount: number
}

export interface TreeListResponse {
  nodes: TreeNode[]
  rootId: number
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function getDeviceTree(): Promise<TreeListResponse> {
  const response = await apiClient.get<HmdmEnvelope<TreeListResponse>>('/private/device-tree')
  return unwrap(response, 'Failed to load device tree.')
}

export async function createTreeNode(parentId: number, name: string): Promise<TreeNode> {
  const response = await apiClient.post<HmdmEnvelope<TreeNode>>('/private/device-tree/nodes', {
    parentId,
    name,
    sortOrder: 0,
  })
  return unwrap(response, 'Failed to create folder.')
}

export async function deleteTreeNode(id: number, targetNodeId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/device-tree/nodes/${id}/delete`, {
    targetNodeId,
  })
  assertHmdmOk(response.data, 'Failed to delete folder.')
}

export async function renameTreeNode(id: number, name: string): Promise<TreeNode> {
  const response = await apiClient.put<HmdmEnvelope<TreeNode>>(`/private/device-tree/nodes/${id}`, {
    name,
  })
  return unwrap(response, 'Failed to rename folder.')
}

export async function moveDeviceToTree(deviceId: number, treeNodeId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/devices/${deviceId}/move-tree`, {
    treeNodeId,
  })
  assertHmdmOk(response.data, 'Failed to move device.')
}
