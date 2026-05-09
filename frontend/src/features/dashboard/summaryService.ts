import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface ChartItem {
  stringAttr?: string | null
  intAttr?: number
  number?: number
}

/** Subset of {@code SummaryResponse.java} fields used by the dashboard. */
export interface DeviceSummaryData {
  statusSummary?: ChartItem[]
  installSummary?: ChartItem[]
  devicesTotal?: number
  devicesEnrolled?: number
  devicesEnrolledLastMonth?: number
  devicesEnrolledMonthly?: ChartItem[]
  topConfigs?: string[]
}

export async function getDeviceSummary(): Promise<DeviceSummaryData> {
  const response = await apiClient.get<HmdmEnvelope<DeviceSummaryData>>('/private/summary/devices')
  return unwrapHmdmData(response.data, 'Failed to load device summary.')
}
