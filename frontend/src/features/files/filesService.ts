import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface FileRecord {
  id?: number
  url?: string | null
  filePath?: string | null
  description?: string | null
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

/** Admin file library (`FilesResource`): list all uploads. */
export async function searchAllFiles(): Promise<FileRecord[]> {
  const response = await apiClient.get<HmdmEnvelope<FileRecord[]>>('/private/web-ui-files/search')
  const data = unwrap(response, 'Failed to load files.')
  return Array.isArray(data) ? data : []
}

export async function removeFile(payload: Partial<FileRecord> & Pick<FileRecord, 'id'>): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/web-ui-files/remove', payload)
  assertHmdmOk(response.data, 'Failed to delete file.')
}
