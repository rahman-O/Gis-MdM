import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData, unwrapHmdmList } from '@/services/hmdmEnvelope'
import type { CreateProfilePayload, Profile, ProfileListItem, ProfileMeta } from '@/features/profiles/types'

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export async function listProfiles(): Promise<ProfileListItem[]> {
  const response = await apiClient.get<HmdmEnvelope<ProfileListItem[]>>('/private/profiles')
  return unwrapHmdmList(response.data, 'Failed to load profiles.')
}

export async function getProfileMeta(profileId: number): Promise<ProfileMeta> {
  const response = await apiClient.get<HmdmEnvelope<ProfileMeta>>(`/private/profiles/${profileId}`)
  return unwrap(response, 'Failed to load profile.')
}

export async function getProfileVersion(profileId: number, versionId: number): Promise<Profile> {
  const response = await apiClient.get<HmdmEnvelope<Profile>>(
    `/private/profiles/${profileId}/versions/${versionId}`
  )
  return unwrap(response, 'Failed to load profile version.')
}

export async function saveProfileVersion(
  profileId: number,
  versionId: number,
  body: Profile
): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>(
    `/private/profiles/${profileId}/versions/${versionId}`,
    { ...body, id: profileId, profileId }
  )
  assertHmdmOk(response.data, 'Failed to save profile draft.')
}

export async function createProfile(payload: CreateProfilePayload): Promise<ProfileMeta> {
  const response = await apiClient.post<HmdmEnvelope<ProfileMeta>>('/private/profiles', payload)
  return unwrap(response, 'Failed to create profile.')
}

export interface PublishImpactAssignment {
  assignmentId: number
  treeNodeId: number
  treeNodeName: string
  currentVersionNumber: number
  deviceCount: number
}

export interface ProfileImpact {
  deviceCount: number
  enrollmentRouteCount: number
  requiresConfirmDialog: boolean
  assignmentsToUpdate?: PublishImpactAssignment[]
}

export interface PublishProfileResult {
  publishedVersionId: number
  versionNumber: number
  artifactHash: string
  affectedDevices: number
  affectedRoutes: number
  assignmentsUpdated?: number
}

export async function getProfileImpact(profileId: number): Promise<ProfileImpact> {
  const response = await apiClient.get<HmdmEnvelope<ProfileImpact>>(`/private/profiles/${profileId}/impact`)
  return unwrap(response, 'Failed to load publish impact.')
}

export async function publishProfileVersion(
  profileId: number,
  versionId: number,
  confirmImpact: boolean
): Promise<PublishProfileResult> {
  const response = await apiClient.post<HmdmEnvelope<PublishProfileResult>>(
    `/private/profiles/${profileId}/versions/${versionId}/publish`,
    { confirmImpact }
  )
  return unwrap(response, 'Failed to publish profile.')
}

export {
  configurationKindToType as profileKindToType,
  typeToConfigurationKind as typeToProfileKind,
} from '@/features/configurations/configurationService'
