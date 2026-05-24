import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData, unwrapHmdmList } from '@/services/hmdmEnvelope'

export interface EnrollmentRouteListItem {
  id: number
  name: string
  description?: string
  qrcodekey?: string
  profileId?: number
  profileVersionId?: number
  profileVersionNumber?: number | null
  defaultTreeNodeId?: number
  defaultTreeNodeName?: string
  defaultDeviceIdMode: string
  mainAppId?: number | null
}

export interface EnrollmentRouteDetail extends EnrollmentRouteListItem {
  type?: number
}

export interface PublishedProfileVersionOption {
  profileVersionId: number
  profileId: number
  profileName: string
  versionNumber: number
  profileEnabled?: boolean
  mainAppId?: number | null
}

export interface CreateEnrollmentRoutePayload {
  name: string
  description?: string | null
  profileVersionId?: number | null
  defaultTreeNodeId: number
  defaultDeviceIdMode?: string
  mainAppId?: number | null
}

export interface UpdateEnrollmentRoutePayload {
  name?: string
  description?: string | null
  profileVersionId?: number
  defaultTreeNodeId?: number
  defaultDeviceIdMode?: string
  mainAppId?: number | null
}

export interface EnrollmentRouteQrMeta {
  qrcodekey: string
  defaultDeviceIdMode: string
  mainAppId?: number | null
}

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export async function listEnrollmentRoutes(): Promise<EnrollmentRouteListItem[]> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentRouteListItem[]>>(
    '/private/enrollment-routes'
  )
  return unwrapHmdmList(response.data, 'Failed to load enrollment routes.')
}

export async function getEnrollmentRoute(id: number): Promise<EnrollmentRouteDetail> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentRouteDetail>>(
    `/private/enrollment-routes/${id}`
  )
  return unwrap(response, 'Failed to load enrollment route.')
}

export async function listPublishedProfileVersions(): Promise<PublishedProfileVersionOption[]> {
  const response = await apiClient.get<HmdmEnvelope<PublishedProfileVersionOption[]>>(
    '/private/enrollment-routes/options/published-profile-versions'
  )
  return unwrapHmdmList(response.data, 'Failed to load published profile versions.')
}

export async function createEnrollmentRoute(
  payload: CreateEnrollmentRoutePayload
): Promise<EnrollmentRouteDetail> {
  const response = await apiClient.post<HmdmEnvelope<EnrollmentRouteDetail>>(
    '/private/enrollment-routes',
    payload
  )
  return unwrap(response, 'Failed to create enrollment route.')
}

export async function updateEnrollmentRoute(
  id: number,
  payload: UpdateEnrollmentRoutePayload
): Promise<EnrollmentRouteDetail> {
  const response = await apiClient.put<HmdmEnvelope<EnrollmentRouteDetail>>(
    `/private/enrollment-routes/${id}`,
    payload
  )
  return unwrap(response, 'Failed to save enrollment route.')
}

export async function getEnrollmentRouteQrMeta(id: number): Promise<EnrollmentRouteQrMeta> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentRouteQrMeta>>(
    `/private/enrollment-routes/${id}/qr`
  )
  return unwrap(response, 'Failed to load QR metadata.')
}

export async function deleteEnrollmentRoute(_id: number): Promise<void> {
  assertHmdmOk({ status: 'OK' } as HmdmEnvelope<unknown>, '')
}
