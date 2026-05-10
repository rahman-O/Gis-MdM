import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface PluginRow {
  id: number
  identifier?: string
  nameLocalizationKey?: string
}

export async function fetchActivePlugins(): Promise<PluginRow[]> {
  const response = await apiClient.get<HmdmEnvelope<PluginRow[]>>('/plugin/main/private/active')
  const data = unwrapHmdmData(response.data, 'Failed to load active plugins.')
  return Array.isArray(data) ? data : []
}

export async function fetchAvailablePlugins(): Promise<PluginRow[]> {
  const response = await apiClient.get<HmdmEnvelope<PluginRow[]>>('/plugin/main/private/available')
  const data = unwrapHmdmData(response.data, 'Failed to load available plugins.')
  return Array.isArray(data) ? data : []
}

/** Marks plugins whose IDs appear in `disabledIds` as disabled for this tenant (matches legacy tab logic). */
export async function saveDisabledPlugins(disabledIds: number[]): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/plugin/main/private/disabled', disabledIds)
  assertHmdmOk(response.data, 'Failed to save plugins.')
}
