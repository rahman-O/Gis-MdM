import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import type {
  AppSetting,
  BulkDeletePayload,
  ConfigurationOption,
  DeviceListResponse,
  DevicePayload,
  DeviceSearchRequest,
  DeviceView,
  GroupBulkPayload,
  LookupItem,
} from '@/features/devices/types'

function unwrapEnvelope<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export async function getDevices(params: DeviceSearchRequest): Promise<DeviceListResponse> {
  const body: Record<string, unknown> = {
    pageNum: params.pageNum,
    pageSize: params.pageSize,
  }
  const candidates: Record<string, unknown> = {
    value: params.value?.trim() ? params.value.trim() : null,
    groupId: params.groupId,
    configurationId: params.configurationId,
    status: params.status?.trim() ? params.status.trim() : null,
    androidVersion: params.androidVersion?.trim() ? params.androidVersion.trim() : null,
    sortBy: params.sortBy?.trim() ? params.sortBy.trim() : null,
    sortDir: params.sortDir,
    dateFrom: params.dateFrom,
    dateTo: params.dateTo,
    onlineEarlierMillis: params.onlineEarlierMillis,
    onlineLaterMillis: params.onlineLaterMillis,
    enrollmentDateFrom: params.enrollmentDateFrom,
    enrollmentDateTo: params.enrollmentDateTo,
    mdmMode: params.mdmMode,
    kioskMode: params.kioskMode,
    launcherVersion: params.launcherVersion?.trim() ? params.launcherVersion.trim() : null,
    installationStatus: params.installationStatus?.trim() ? params.installationStatus.trim() : null,
    imeiChanged: params.imeiChanged,
    fastSearch: params.fastSearch,
  }
  Object.entries(candidates).forEach(([key, value]) => {
    if (value !== null && value !== undefined && value !== '') {
      body[key] = value
    }
  })
  const response = await apiClient.post<HmdmEnvelope<DeviceListResponse>>('/private/devices/search', body)
  return unwrapEnvelope(response, 'Failed to load devices.')
}

export async function getDevice(number: string): Promise<DeviceView> {
  const response = await apiClient.get<HmdmEnvelope<DeviceView>>(
    `/private/devices/number/${encodeURIComponent(number)}`
  )
  return unwrapEnvelope(response, 'Device not found.')
}

export async function createDevice(payload: DevicePayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/devices', payload)
  assertHmdmOk(response.data, 'Failed to create device.')
}

export async function updateDevice(payload: DevicePayload): Promise<void> {
  const response = await apiClient.put<HmdmEnvelope<unknown>>('/private/devices', payload)
  assertHmdmOk(response.data, 'Failed to update device.')
}

export async function deleteDevice(id: number): Promise<void> {
  const response = await apiClient.delete<HmdmEnvelope<unknown>>(`/private/devices/${id}`)
  assertHmdmOk(response.data, 'Failed to delete device.')
}

export async function getGroups(): Promise<LookupItem[]> {
  const response = await apiClient.get<HmdmEnvelope<LookupItem[]>>('/private/groups/search')
  return unwrapEnvelope(response, 'Failed to load groups.')
}

export async function getConfigurations(): Promise<ConfigurationOption[]> {
  const response = await apiClient.get<HmdmEnvelope<ConfigurationOption[]>>('/private/configurations/list')
  return unwrapEnvelope(response, 'Failed to load configurations.')
}

export async function deleteBulk(ids: number[]): Promise<void> {
  const payload: BulkDeletePayload = { ids }
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/devices/deleteBulk', payload)
  assertHmdmOk(response.data, 'Failed to delete selected devices.')
}

export async function groupBulk(payload: GroupBulkPayload): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/devices/groupBulk', payload)
  assertHmdmOk(response.data, 'Failed to update selected groups.')
}

export async function getAppSettings(deviceId: number): Promise<AppSetting[]> {
  const response = await apiClient.get<HmdmEnvelope<AppSetting[]>>(`/private/devices/${deviceId}/applicationSettings`)
  return unwrapEnvelope(response, 'Failed to load app settings.')
}

export async function saveAppSettings(deviceId: number, settings: AppSetting[]): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/devices/${deviceId}/applicationSettings`, settings)
  assertHmdmOk(response.data, 'Failed to save app settings.')
}

export async function notifyAppSettings(deviceId: number): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(
    `/private/devices/${deviceId}/applicationSettings/notify`
  )
  assertHmdmOk(response.data, 'Failed to notify device.')
}

export async function updateDescription(deviceId: number, description: string): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>(`/private/devices/${deviceId}/description`, description)
  assertHmdmOk(response.data, 'Failed to update description.')
}

export async function autocomplete(value: string): Promise<string[]> {
  const response = await apiClient.post<HmdmEnvelope<Array<{ value?: string | null; number?: string | null }>>>(
    '/private/devices/autocomplete',
    value
  )
  const items = unwrapEnvelope(response, 'Failed to fetch autocomplete.')
  return items.map((item) => item.value ?? item.number ?? '').filter(Boolean)
}
