import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface IconRow {
  id?: number | null
  name?: string | null
  fileId?: number | null
  fileName?: string | null
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function listIcons(searchValue?: string): Promise<IconRow[]> {
  const v = searchValue?.trim()
  const url = v ? `/private/icons/search/${encodeURIComponent(v)}` : '/private/icons/search'
  const response = await apiClient.get<HmdmEnvelope<IconRow[]>>(url)
  const data = unwrap(response, 'Failed to load icons.')
  return Array.isArray(data) ? data : []
}

export async function saveIcon(icon: IconRow): Promise<IconRow> {
  const response = await apiClient.put<HmdmEnvelope<IconRow>>('/private/icons', icon)
  return unwrap(response, 'Failed to save icon.')
}

export async function deleteIcon(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/icons/${id}`)
  assertHmdmOk(response.data, 'Failed to delete icon.')
}
