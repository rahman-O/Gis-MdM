import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import { getDevices } from '@/features/devices/deviceService'
import type { DeviceView } from '@/features/devices/types'
import type {
  ChartItemLike,
  DashboardCounts,
  DeviceSummaryPayload,
  ParsedInstallCounts,
  ParsedStatusCounts,
} from '@/features/dashboard/types'

function unwrap<T>(response: { data: HmdmEnvelope<T> }, message: string): T {
  return unwrapHmdmData(response.data, message)
}

export type { DeviceSummaryPayload as DeviceSummaryData }

/** Full summary from the server device statistics endpoint (do not confuse with `/private/summary`). */
export async function getSummaryDevices(): Promise<DeviceSummaryPayload> {
  const response = await apiClient.get<HmdmEnvelope<DeviceSummaryPayload>>('/private/summary/devices')
  return unwrap(response, 'Failed to load device summary.')
}

export function parseStatusSummary(items: ChartItemLike[] | undefined): ParsedStatusCounts {
  let offline = 0
  let idle = 0
  let online = 0
  for (const item of items ?? []) {
    const v = typeof item.number === 'number' ? item.number : 0
    if (item.stringAttr === 'red') offline = v
    else if (item.stringAttr === 'yellow') idle = v
    else if (item.stringAttr === 'green') online = v
  }
  return { offline, idle, online }
}

export function parseInstallSummary(items: ChartItemLike[] | undefined): ParsedInstallCounts {
  let failure = 0
  let mismatch = 0
  let success = 0
  for (const item of items ?? []) {
    const v = typeof item.number === 'number' ? item.number : 0
    if (item.stringAttr === 'FAILURE') failure = v
    else if (item.stringAttr === 'VERSION_MISMATCH') mismatch = v
    else if (item.stringAttr === 'SUCCESS') success = v
  }
  return { failure, mismatch, success }
}

export async function getRecentDevices(limit = 5): Promise<DeviceView[]> {
  const res = await getDevices({
    pageNum: 1,
    pageSize: limit,
    sortBy: 'LAST_UPDATE',
    sortDir: 'desc',
  })
  return res.devices.items ?? []
}

export async function getConfigurationApplicationCounts(): Promise<DashboardCounts> {
  try {
    const [cRes, aRes] = await Promise.all([
      apiClient.get<HmdmEnvelope<unknown>>('/private/configurations/search'),
      apiClient.get<HmdmEnvelope<unknown>>('/private/applications/search'),
    ])
    const configs = unwrap(cRes, 'Failed to load configurations count.')
    const apps = unwrap(aRes, 'Failed to load applications count.')
    return {
      configurationCount: Array.isArray(configs) ? configs.length : 0,
      applicationCount: Array.isArray(apps) ? apps.length : 0,
    }
  } catch {
    return { configurationCount: 0, applicationCount: 0 }
  }
}

/** Prefer `getSummaryDevices`; kept as alias used by legacy imports/tests. */
export async function getDeviceSummary(): Promise<DeviceSummaryPayload> {
  return getSummaryDevices()
}
