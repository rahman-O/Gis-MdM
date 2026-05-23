import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { FileUploadResult } from '@/features/applications/model/types'

export interface StorageLimitInfo {
  sizeLimit: number
  sizeUsed: number
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

function toFormData(file: File): FormData {
  const fd = new FormData()
  fd.append('file', file)
  return fd
}

export async function uploadApplicationFile(file: File): Promise<FileUploadResult> {
  const response = await apiClient.post<HmdmEnvelope<FileUploadResult>>(
    '/private/web-ui-files',
    toFormData(file)
  )
  return unwrap(response, 'Failed to upload file.')
}

export async function uploadApplicationFileWithProgress(
  file: File,
  onProgress: (percent: number) => void
): Promise<FileUploadResult> {
  const response = await apiClient.post<HmdmEnvelope<FileUploadResult>>(
    '/private/web-ui-files',
    toFormData(file),
    {
      onUploadProgress: (event) => {
        if (!event.total || event.total <= 0) return
        const percent = Math.max(0, Math.min(100, Math.round((event.loaded / event.total) * 100)))
        onProgress(percent)
      },
    }
  )
  return unwrap(response, 'Failed to upload file.')
}

export async function uploadRawFile(file: File): Promise<FileUploadResult> {
  const response = await apiClient.post<HmdmEnvelope<FileUploadResult>>(
    '/private/web-ui-files/raw',
    toFormData(file)
  )
  return unwrap(response, 'Failed to upload raw file.')
}

export async function commitUploadedFile(payload: Record<string, unknown>): Promise<FileUploadResult | void> {
  const response = await apiClient.post<HmdmEnvelope<FileUploadResult | undefined>>(
    '/private/web-ui-files/update',
    payload
  )
  return unwrap(response, 'Failed to commit uploaded file.')
}

export async function getStorageLimit(): Promise<StorageLimitInfo> {
  const response = await apiClient.get<HmdmEnvelope<StorageLimitInfo>>('/private/web-ui-files/limit')
  return unwrap(response, 'Failed to load storage limit.')
}

export async function getApplicationsByFileUrl(url: string): Promise<Array<Record<string, unknown>>> {
  const response = await apiClient.get<HmdmEnvelope<Array<Record<string, unknown>>>>(
    `/private/web-ui-files/apps/${encodeURIComponent(url)}`
  )
  return unwrap(response, 'Failed to load file applications.')
}

export async function assertUploadOk(envelope: HmdmEnvelope<unknown>, msg: string): Promise<void> {
  assertHmdmOk(envelope, msg)
}
