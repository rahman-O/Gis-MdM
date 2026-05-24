import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData, unwrapHmdmList } from '@/services/hmdmEnvelope'

export interface ProfileVersionListItem {
  versionId: number
  versionNumber: number
  status: string
  publishedAt?: string | null
  createdAt: string
}

export interface ProfileTreeAssignment {
  assignmentId: number
  treeNodeId: number
  treeNodeName: string
  treePath?: string
  profileVersionId: number
  versionNumber: number
  deviceCount: number
  createdAt: string
}

export interface AssignmentImpact {
  deviceCount: number
  requiresConfirmDialog: boolean
  folderName: string
}

export interface DeviceRolloutRow {
  deviceId: number
  deviceName: string
  treeNodeId?: number
  treeNodeName?: string
  targetVersionId?: number
  targetVersionNumber?: number
  appliedVersionId?: number
  appliedVersionNumber?: number
  status: string
  reason?: string
  lastUpdate?: number
}

export interface RolloutDevicesPage {
  items: DeviceRolloutRow[]
  totalCount: number
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export async function listProfileVersions(profileId: number): Promise<ProfileVersionListItem[]> {
  const response = await apiClient.get<HmdmEnvelope<ProfileVersionListItem[]>>(
    `/private/profiles/${profileId}/versions`
  )
  return unwrapHmdmList(response.data, 'Failed to load profile versions.')
}

export async function forkDraftFromVersion(profileId: number, versionId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(
    `/private/profiles/${profileId}/versions/${versionId}/fork-draft`,
    {}
  )
  unwrap(response, 'Failed to create draft from version.')
}

export async function listAssignments(profileId: number): Promise<ProfileTreeAssignment[]> {
  const response = await apiClient.get<HmdmEnvelope<ProfileTreeAssignment[]>>(
    `/private/profiles/${profileId}/assignments`
  )
  return unwrapHmdmList(response.data, 'Failed to load assignments.')
}

export async function getAssignmentImpact(
  profileId: number,
  treeNodeId: number
): Promise<AssignmentImpact> {
  const response = await apiClient.get<HmdmEnvelope<AssignmentImpact>>(
    `/private/profiles/${profileId}/assignments/impact`,
    { params: { treeNodeId } }
  )
  return unwrap(response, 'Failed to load assignment impact.')
}

export async function putAssignment(
  profileId: number,
  body: { treeNodeId: number; profileVersionId: number; confirmImpact: boolean }
): Promise<{ affectedDevices: number }> {
  const response = await apiClient.put<HmdmEnvelope<{ affectedDevices: number }>>(
    `/private/profiles/${profileId}/assignments`,
    body
  )
  return unwrap(response, 'Failed to save assignment.')
}

export async function deleteAssignment(profileId: number, assignmentId: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(
    `/private/profiles/${profileId}/assignments/${assignmentId}`
  )
  unwrap(response, 'Failed to remove assignment.')
}

export { deleteProfileVersion } from '@/features/profiles/profileHubService'

export async function listRolloutDevices(
  profileId: number,
  params?: { treeNodeId?: number; status?: string; page?: number; pageSize?: number }
): Promise<RolloutDevicesPage> {
  const response = await apiClient.get<HmdmEnvelope<RolloutDevicesPage>>(
    `/private/profiles/${profileId}/rollout/devices`,
    { params }
  )
  return unwrap(response, 'Failed to load rollout status.')
}

export async function recomputeRollout(profileId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(
    `/private/profiles/${profileId}/rollout/recompute`,
    {}
  )
  unwrap(response, 'Failed to refresh rollout status.')
}

export async function disableProfile(profileId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/profiles/${profileId}/disable`, {})
  unwrap(response, 'Failed to disable profile.')
}

export async function enableProfile(profileId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/profiles/${profileId}/enable`, {})
  unwrap(response, 'Failed to enable profile.')
}
