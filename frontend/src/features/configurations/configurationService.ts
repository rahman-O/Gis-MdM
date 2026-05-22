import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData, unwrapHmdmList } from '@/services/hmdmEnvelope'
import {
  buildCreateConfigurationBody,
  mergeConfigurationForUpdate,
} from '@/features/configurations/configurationNormalize'
import {
  mapApplicationCatalogRows,
  type ConfigurationAppCatalogItem,
} from '@/features/configurations/configurationCatalog'
import type {
  Configuration,
  ConfigurationAutocompleteRequest,
  ConfigurationLookupItem,
  ConfigurationPayload,
  CopyConfigurationPayload,
  UpgradeConfigurationApplicationPayload,
} from '@/features/configurations/types'

export function configurationKindToType(kind: ConfigurationPayload['type']): number {
  return kind === 'COMMON' ? 1 : 0
}

export function typeToConfigurationKind(type: number | null | undefined): ConfigurationPayload['type'] {
  return type === 1 ? 'COMMON' : 'WORK'
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export type { ConfigurationAppCatalogItem }

export async function getConfigurations(): Promise<Configuration[]> {
  const response = await apiClient.get<HmdmEnvelope<Configuration[]>>('/private/configurations/search')
  return unwrapHmdmList(response.data, 'Failed to load configurations.')
}

export async function searchConfigurations(value: string): Promise<Configuration[]> {
  const response = await apiClient.get<HmdmEnvelope<Configuration[]>>(
    `/private/configurations/search/${encodeURIComponent(value)}`
  )
  return unwrapHmdmList(response.data, 'Failed to search configurations.')
}

export async function listConfigurationNames(): Promise<ConfigurationLookupItem[]> {
  const response = await apiClient.get<HmdmEnvelope<ConfigurationLookupItem[]>>(
    '/private/configurations/list'
  )
  return unwrapHmdmList(response.data, 'Failed to load configuration names.')
}

export async function autocompleteConfigurations(
  request: ConfigurationAutocompleteRequest
): Promise<ConfigurationLookupItem[]> {
  const response = await apiClient.post<HmdmEnvelope<ConfigurationLookupItem[]>>(
    '/private/configurations/autocomplete',
    request.value
  )
  return unwrapHmdmList(response.data, 'Failed to autocomplete configurations.')
}

export async function getConfiguration(id: number): Promise<Configuration> {
  const response = await apiClient.get<HmdmEnvelope<Configuration>>(`/private/configurations/${id}`)
  return unwrap(response, 'Failed to load configuration.')
}

async function putConfiguration(body: Configuration): Promise<Configuration> {
  const response = await apiClient.put<HmdmEnvelope<Configuration>>('/private/configurations', body)
  return unwrap(response, 'Failed to save configuration.')
}

export async function saveConfiguration(body: Configuration): Promise<Configuration> {
  return putConfiguration(body)
}

export async function createConfiguration(data: ConfigurationPayload): Promise<Configuration> {
  return putConfiguration(buildCreateConfigurationBody(data))
}

export async function updateConfiguration(id: number, data: ConfigurationPayload): Promise<Configuration> {
  const existing = await getConfiguration(id)
  return putConfiguration(mergeConfigurationForUpdate(existing, id, data))
}

export async function copyConfiguration(payload: CopyConfigurationPayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/configurations/copy', payload)
  assertHmdmOk(response.data, 'Failed to copy configuration.')
}

export async function getConfigurationApplications(id: number): Promise<Configuration[]> {
  const response = await apiClient.get<HmdmEnvelope<Configuration[]>>(
    `/private/configurations/applications/${id}`
  )
  return unwrapHmdmList(response.data, 'Failed to load configuration applications.')
}

export async function getAllApplications(): Promise<ConfigurationAppCatalogItem[]> {
  try {
    const response = await apiClient.get<HmdmEnvelope<unknown>>('/private/applications/search')
    const list = mapApplicationCatalogRows(unwrapHmdmList(response.data, 'Failed to load applications.'))
    if (list.length > 0) return list
  } catch {
    // fallback below
  }
  const fallback = await apiClient.get<HmdmEnvelope<unknown>>('/private/configurations/applications')
  return mapApplicationCatalogRows(unwrapHmdmList(fallback.data, 'Failed to load applications.'))
}

export async function upgradeConfigurationApplication(
  payload: UpgradeConfigurationApplicationPayload
): Promise<Configuration> {
  const response = await apiClient.put<HmdmEnvelope<Configuration>>(
    '/private/configurations/application/upgrade',
    payload
  )
  return unwrap(response, 'Failed to upgrade configuration application.')
}

export async function deleteConfiguration(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown> | string>(
    `/private/configurations/${id}`
  )
  const data = response.data
  if (data && typeof data === 'object' && 'status' in data) {
    assertHmdmOk(data as HmdmEnvelope<unknown>, 'Failed to delete configuration.')
  }
}
