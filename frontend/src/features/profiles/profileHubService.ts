import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'

export type ProfileHealth = 'healthy' | 'warning' | 'error' | 'draft_only'

export interface ProfilePinnedSettings {
  kioskMode: boolean
  mainAppName?: string
  appCount: number
  lastPublishedAt?: string | null
}

export interface PublishedContext {
  versionId: number
  versionNumber: number
  status: string
  pinnedSettings: ProfilePinnedSettings
}

export interface ProfileSummary {
  id: number
  name: string
  description: string
  enabled: boolean
  health: ProfileHealth
  healthReasons?: string[]
  lifecycle: string
  publishedVersionId?: number | null
  publishedVersionNumber?: number | null
  draftVersionId?: number | null
  hasUnpublishedDraft: boolean
  canPublish: boolean
  assignmentCount: number
  assignedFolders: string[]
  rollout: {
    pending: number
    installed: number
    partial: number
    failed: number
    total: number
  }
  pinnedSettings: ProfilePinnedSettings
  publishedContext?: PublishedContext | null
}

export interface ProfileActivityEvent {
  id: number
  eventType: string
  summaryKey: string
  summaryParams?: Record<string, unknown>
  occurredAt: string
  actorUserId?: number | null
}

export interface ProfileActivityPage {
  items: ProfileActivityEvent[]
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export async function getProfileSummary(profileId: number): Promise<ProfileSummary> {
  const response = await apiClient.get<HmdmEnvelope<ProfileSummary>>(
    `/private/profiles/${profileId}/summary`
  )
  return unwrap(response, 'Failed to load profile summary.')
}

export async function getProfileActivity(
  profileId: number,
  limit = 50
): Promise<ProfileActivityPage> {
  const response = await apiClient.get<HmdmEnvelope<ProfileActivityPage>>(
    `/private/profiles/${profileId}/activity`,
    { params: { limit } }
  )
  return unwrap(response, 'Failed to load profile activity.')
}

export async function deleteProfileVersion(profileId: number, versionId: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(
    `/private/profiles/${profileId}/versions/${versionId}`
  )
  unwrap(response, 'Failed to delete profile version.')
}
