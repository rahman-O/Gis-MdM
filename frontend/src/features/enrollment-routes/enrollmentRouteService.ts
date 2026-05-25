import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData, unwrapHmdmList } from '@/services/hmdmEnvelope'

export type BootstrapIntent = 'stable' | 'specific' | 'latest'

export interface EnrollmentRouteView {
  id: number
  name: string
  description?: string
  qrcodekey?: string
  targetNodeId: number
  targetNodeName?: string
  targetNodePath?: string
  targetPlacementKind?: 'locked' | 'inheritable'
  containerPlacementAcknowledged?: boolean
  deviceIdentityMode: string
  bootstrapIntent: BootstrapIntent
  bootstrapApplicationId: number
  bootstrapApplicationName?: string
  bootstrapVersionId?: number | null
  resolvedMainAppVersionId?: number | null
  resolvedVersionLabel?: string
  resolvedPackage?: string
  status: 'draft' | 'active' | string
  type?: number
  wifiSsid?: string
  wifiPassword?: string
  wifiSecurityType?: string
  qrParameters?: string
  adminExtras?: string
  mobileEnrollment?: boolean
  encryptDevice?: boolean
}

export interface TreeNodeOption {
  id: number
  name: string
  path: string
  parentId?: number | null
  placementKind: 'locked' | 'inheritable'
  deviceCount: number
  heavilyLoaded: boolean
}

export interface BootstrapAppVersionOption {
  versionId: number
  version: string
  versionCode: number
  isRecommended: boolean
  isLatest: boolean
}

export interface BootstrapAppOption {
  applicationId: number
  name: string
  package: string
  versions: BootstrapAppVersionOption[]
}

export interface EnrollmentRouteQrMeta {
  qrcodekey: string
  defaultDeviceIdMode: string
  resolvedMainAppVersionId?: number | null
  mainAppPackage?: string
  mainAppVersion?: string
  mainAppVersionCode?: number
  targetNodeId?: number
  contract?: Record<string, unknown>
}

export interface EnrollmentDeleteImpact {
  enrollingNowCount: number
  historicalEnrolledCount: number
  activeQrScans7d: number
}

export interface CreateEnrollmentRoutePayload {
  name: string
  description?: string | null
  targetNodeId: number
  deviceIdentityMode?: string
  bootstrapIntent: BootstrapIntent
  bootstrapApplicationId: number
  bootstrapVersionId?: number | null
  acknowledgeContainerPlacement?: boolean
  wifiSsid?: string
  wifiPassword?: string
  wifiSecurityType?: string
  qrParameters?: string
  adminExtras?: string
  mobileEnrollment?: boolean
  encryptDevice?: boolean
}

export type UpdateEnrollmentRoutePayload = Partial<CreateEnrollmentRoutePayload>

function unwrap<T>(response: { data: HmdmEnvelope<T> }, msg: string): T {
  return unwrapHmdmData(response.data, msg)
}

export async function listEnrollmentRoutes(): Promise<EnrollmentRouteView[]> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentRouteView[]>>(
    '/private/enrollment-routes'
  )
  return unwrapHmdmList(response.data, 'Failed to load enrollment routes.')
}

export async function getEnrollmentRoute(id: number): Promise<EnrollmentRouteView> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentRouteView>>(
    `/private/enrollment-routes/${id}`
  )
  return unwrap(response, 'Failed to load enrollment route.')
}

export async function listTreeNodeOptions(): Promise<TreeNodeOption[]> {
  const response = await apiClient.get<HmdmEnvelope<TreeNodeOption[]>>(
    '/private/enrollment-routes/options/tree-nodes'
  )
  return unwrapHmdmList(response.data, 'Failed to load tree folders.')
}

export async function listBootstrapApps(): Promise<BootstrapAppOption[]> {
  const response = await apiClient.get<HmdmEnvelope<BootstrapAppOption[]>>(
    '/private/enrollment-routes/options/bootstrap-apps'
  )
  return unwrapHmdmList(response.data, 'Failed to load bootstrap applications.')
}

export async function createEnrollmentRoute(
  payload: CreateEnrollmentRoutePayload
): Promise<EnrollmentRouteView> {
  const response = await apiClient.post<HmdmEnvelope<EnrollmentRouteView>>(
    '/private/enrollment-routes',
    payload
  )
  return unwrap(response, 'Failed to create enrollment route.')
}

export async function updateEnrollmentRoute(
  id: number,
  payload: UpdateEnrollmentRoutePayload
): Promise<EnrollmentRouteView> {
  const response = await apiClient.put<HmdmEnvelope<EnrollmentRouteView>>(
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

export async function getEnrollmentRouteImpact(id: number): Promise<EnrollmentDeleteImpact> {
  const response = await apiClient.get<HmdmEnvelope<EnrollmentDeleteImpact>>(
    `/private/enrollment-routes/${id}/impact`
  )
  return unwrap(response, 'Failed to load delete impact.')
}

export async function deleteEnrollmentRoute(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(
    `/private/enrollment-routes/${id}`
  )
  assertHmdmOk(response.data, 'Failed to delete enrollment route.')
}
