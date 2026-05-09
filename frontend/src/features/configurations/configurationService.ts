import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import {
  buildCreateConfigurationBody,
  mergeConfigurationForUpdate,
} from '@/features/configurations/configurationNormalize'
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

function asObject(value: unknown): Record<string, unknown> | null {
  return value != null && typeof value === 'object' ? (value as Record<string, unknown>) : null
}

function extractApplications(payload: unknown): Configuration[] {
  if (Array.isArray(payload)) return payload as Configuration[]
  const obj = asObject(payload)
  if (!obj) return []
  const itemsCandidate = obj.items
  if (Array.isArray(itemsCandidate)) return itemsCandidate as Configuration[]
  return []
}

export async function getConfigurations(): Promise<Configuration[]> {
  const response = await apiClient.get<HmdmEnvelope<Configuration[]>>('/private/configurations/search')
  return unwrap(response, 'Failed to load configurations.')
}

export async function searchConfigurations(value: string): Promise<Configuration[]> {
  const response = await apiClient.get<HmdmEnvelope<Configuration[]>>(
    `/private/configurations/search/${encodeURIComponent(value)}`
  )
  return unwrap(response, 'Failed to search configurations.')
}

export async function listConfigurationNames(): Promise<ConfigurationLookupItem[]> {
  const response = await apiClient.get<HmdmEnvelope<ConfigurationLookupItem[]>>(
    '/private/configurations/list'
  )
  return unwrap(response, 'Failed to load configuration names.')
}

export async function autocompleteConfigurations(
  request: ConfigurationAutocompleteRequest
): Promise<ConfigurationLookupItem[]> {
  const response = await apiClient.post<HmdmEnvelope<ConfigurationLookupItem[]>>(
    '/private/configurations/autocomplete',
    request.value
  )
  return unwrap(response, 'Failed to autocomplete configurations.')
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
  return unwrap(response, 'Failed to load configuration applications.')
}

export async function getAllApplications(): Promise<Configuration[]> {
  try {
    const response = await apiClient.get<HmdmEnvelope<unknown>>('/private/configurations/applications')
    const data = unwrap(response, 'Failed to load applications.')
    const list = extractApplications(data)
    if (list.length > 0) return list
  } catch {
    // fallback below
  }
  const fallback = await apiClient.get<HmdmEnvelope<unknown>>('/private/applications/search')
  const data = unwrap(fallback, 'Failed to load applications.')
  return extractApplications(data)
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
