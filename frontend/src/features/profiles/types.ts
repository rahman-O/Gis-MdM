import type {
  Configuration,
  ConfigurationApplication,
  ConfigurationApplicationSetting,
  ConfigurationFile,
  ConfigurationKind,
  ConfigurationPayload,
} from '@/features/configurations/types'

export type ProfileKind = ConfigurationKind
export type ProfilePayload = ConfigurationPayload
export type ProfileApplication = ConfigurationApplication
export type ProfileFile = ConfigurationFile
export type ProfileApplicationSetting = ConfigurationApplicationSetting

/** Full profile version editor payload (same shape as legacy Configuration). */
export type Profile = Configuration & {
  profileId?: number | null
  versionId?: number | null
  versionNumber?: number | null
  versionStatus?: string | null
}

export interface ProfileListItem {
  id: number
  name: string
  description: string
  enabled?: boolean
  publishedVersion?: number | null
  draftVersionId?: number | null
  deviceCount: number
  enrollmentRouteCount: number
  health?: string
  healthReasons?: string[]
  badges?: string[]
  assignmentCount?: number
  rolloutFailureCount?: number
}

export interface ProfileMeta {
  id: number
  name: string
  description: string
  enabled?: boolean
  draftVersionId?: number | null
  publishedVersionId?: number | null
  publishedVersion?: number | null
  deviceCount?: number
  enrollmentRouteCount?: number
}

export interface CreateProfilePayload {
  name: string
  description?: string | null
}
