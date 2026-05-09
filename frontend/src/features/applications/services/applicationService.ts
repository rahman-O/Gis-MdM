import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type {
  Application,
  ApplicationConfigurationLink,
  ApplicationVersion,
  ApplicationVersionConfigurationLink,
  LinkConfigurationsToAppRequest,
  LinkConfigurationsToAppVersionRequest,
} from '@/features/applications/model/types'

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function getAllApplications(): Promise<Application[]> {
  const response = await apiClient.get<HmdmEnvelope<Application[]>>('/private/applications/search')
  return unwrap(response, 'Failed to load applications.')
}

export async function searchApplications(value: string): Promise<Application[]> {
  const normalized = value.trim()
  if (!normalized) return getAllApplications()
  const response = await apiClient.get<HmdmEnvelope<Application[]>>(
    `/private/applications/search/${encodeURIComponent(normalized)}`
  )
  return unwrap(response, 'Failed to search applications.')
}

/** Super-admin: all tenants’ applications (control panel). */
export async function getAllAdminApplications(): Promise<Application[]> {
  const response = await apiClient.get<HmdmEnvelope<Application[]>>('/private/applications/admin/search')
  return unwrap(response, 'Failed to load shared applications catalog.')
}

export async function searchAdminApplications(value: string): Promise<Application[]> {
  const normalized = value.trim()
  if (!normalized) return getAllAdminApplications()
  const response = await apiClient.get<HmdmEnvelope<Application[]>>(
    `/private/applications/admin/search/${encodeURIComponent(normalized)}`
  )
  return unwrap(response, 'Failed to search shared applications catalog.')
}

/** Merges duplicates by package into one shared app (server GET, legacy API). */
export async function turnApplicationIntoCommon(id: number): Promise<void> {
  const response = await apiClient.get<HmdmEnvelope<unknown>>(`/private/applications/admin/common/${id}`)
  assertHmdmOk(response.data, 'Failed to turn application into shared application.')
}

export async function getApplicationsForAutocomplete(filter: string): Promise<Array<{ id: number; name: string }>> {
  const response = await apiClient.post<HmdmEnvelope<Array<{ id: number; name: string }>>>(
    '/private/applications/autocomplete',
    filter
  )
  return unwrap(response, 'Failed to autocomplete applications.')
}

export async function getApplication(id: number): Promise<Application> {
  const response = await apiClient.get<HmdmEnvelope<Application>>(`/private/applications/${id}`)
  return unwrap(response, 'Failed to load application.')
}

export async function getApplicationVersions(id: number): Promise<ApplicationVersion[]> {
  const response = await apiClient.get<HmdmEnvelope<ApplicationVersion[]>>(`/private/applications/${id}/versions`)
  return unwrap(response, 'Failed to load application versions.')
}

export async function validateApplicationPkg(payload: Pick<Application, 'id' | 'name' | 'pkg'>): Promise<Application[]> {
  const response = await apiClient.put<HmdmEnvelope<Application[]>>('/private/applications/validatePkg', payload)
  return unwrap(response, 'Failed to validate package.')
}

export async function createOrUpdateAndroidApplication(payload: Application): Promise<Application | void> {
  const response = await apiClient.put<HmdmEnvelope<Application | undefined>>('/private/applications/android', payload)
  return unwrap(response, 'Failed to save Android application.')
}

export async function createOrUpdateWebApplication(payload: Application): Promise<Application | void> {
  const response = await apiClient.put<HmdmEnvelope<Application | undefined>>('/private/applications/web', payload)
  return unwrap(response, 'Failed to save web application.')
}

export async function createOrUpdateApplicationVersion(payload: ApplicationVersion): Promise<ApplicationVersion | void> {
  const response = await apiClient.put<HmdmEnvelope<ApplicationVersion | undefined>>(
    '/private/applications/versions',
    payload
  )
  return unwrap(response, 'Failed to save application version.')
}

export async function deleteApplication(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/applications/${id}`)
  assertHmdmOk(response.data, 'Failed to delete application.')
}

export async function deleteApplicationVersion(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/applications/versions/${id}`)
  assertHmdmOk(response.data, 'Failed to delete application version.')
}

export async function getApplicationConfigurations(id: number): Promise<ApplicationConfigurationLink[]> {
  const response = await apiClient.get<HmdmEnvelope<ApplicationConfigurationLink[]>>(
    `/private/applications/configurations/${id}`
  )
  return unwrap(response, 'Failed to load application configurations.')
}

export async function updateApplicationConfigurations(
  payload: LinkConfigurationsToAppRequest
): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/applications/configurations', payload)
  assertHmdmOk(response.data, 'Failed to update application configurations.')
}

export async function getApplicationVersionConfigurations(
  versionId: number
): Promise<ApplicationVersionConfigurationLink[]> {
  const response = await apiClient.get<HmdmEnvelope<ApplicationVersionConfigurationLink[]>>(
    `/private/applications/version/${versionId}/configurations`
  )
  return unwrap(response, 'Failed to load application version configurations.')
}

export async function updateApplicationVersionConfigurations(
  payload: LinkConfigurationsToAppVersionRequest
): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(
    '/private/applications/version/configurations',
    payload
  )
  assertHmdmOk(response.data, 'Failed to update application version configurations.')
}
