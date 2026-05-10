import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { AxiosResponse } from 'axios'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'

const mocks = vi.hoisted(() => ({
  get: vi.fn(),
  getDevices: vi.fn(),
}))

vi.mock('@/services/apiClient', () => ({
  default: { get: mocks.get, post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/features/devices/deviceService', () => ({
  getDevices: mocks.getDevices,
}))

import * as dashboardService from '@/features/dashboard/dashboardService'

function ok<T>(data: T): AxiosResponse<HmdmEnvelope<T>> {
  return { data: { status: 'OK', data } } as AxiosResponse<HmdmEnvelope<T>>
}

describe('dashboardService', () => {
  beforeEach(() => {
    mocks.get.mockReset()
    mocks.getDevices.mockReset()
  })

  it('getSummaryDevices GETs /private/summary/devices', async () => {
    const summary = {
      devicesTotal: 12,
      statusSummary: [],
      installSummary: [],
    }
    mocks.get.mockResolvedValueOnce(ok(summary))
    const result = await dashboardService.getSummaryDevices()
    expect(mocks.get).toHaveBeenCalledWith('/private/summary/devices')
    expect(result).toEqual(summary)
  })

  it('throws when envelope is not OK for summary devices', async () => {
    mocks.get.mockResolvedValueOnce({ data: { status: 'ERROR', message: 'x' } } as AxiosResponse<HmdmEnvelope<unknown>>)
    await expect(dashboardService.getSummaryDevices()).rejects.toThrow()
  })

  it('getRecentDevices delegates to getDevices with paging and sort', async () => {
    mocks.getDevices.mockResolvedValueOnce({
      devices: {
        items: [
          {
            id: 1,
            configurationId: null,
            number: 'A',
            groups: [],
            description: null,
            statusCode: 'green',
            lastUpdate: null,
            imei: null,
            phone: null,
            model: null,
            batteryLevel: null,
            androidVersion: null,
            serial: null,
            custom1: null,
            custom2: null,
            custom3: null,
            oldNumber: null,
          },
        ],
      },
      configurations: {},
    })
    const rows = await dashboardService.getRecentDevices(5)
    expect(mocks.getDevices).toHaveBeenCalledWith({
      pageNum: 1,
      pageSize: 5,
      sortBy: 'LAST_UPDATE',
      sortDir: 'desc',
    })
    expect(rows).toHaveLength(1)
  })

  it('getConfigurationApplicationCounts counts list endpoints leniently', async () => {
    mocks.get
      .mockResolvedValueOnce(ok([{ id: 1 }, { id: 2 }]))
      .mockResolvedValueOnce(ok([{ id: 10 }]))
    const c = await dashboardService.getConfigurationApplicationCounts()
    expect(c).toEqual({ configurationCount: 2, applicationCount: 1 })
  })
})
