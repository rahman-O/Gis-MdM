import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { AxiosResponse } from 'axios'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'

const mocks = vi.hoisted(() => ({
  post: vi.fn(),
  get: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}))

vi.mock('@/services/apiClient', () => ({
  default: {
    post: mocks.post,
    get: mocks.get,
    put: mocks.put,
    delete: mocks.del,
  },
}))

import * as deviceService from '@/features/devices/deviceService'

const { post, get, put, del } = mocks

function ok<T>(data: T): AxiosResponse<HmdmEnvelope<T>> {
  return { data: { status: 'OK', data } } as AxiosResponse<HmdmEnvelope<T>>
}

describe('deviceService', () => {
  beforeEach(() => {
    post.mockReset()
    get.mockReset()
    put.mockReset()
    del.mockReset()
  })

  it('getDevices posts search body with pageNum, pageSize, value', async () => {
    const payload = {
      devices: { items: [], totalItemsCount: 0 },
      configurations: {},
    }
    post.mockResolvedValueOnce(ok(payload))
    const result = await deviceService.getDevices({ pageNum: 2, pageSize: 20, value: 'ab' })
    expect(post).toHaveBeenCalledWith('/private/devices/search', {
      pageNum: 2,
      pageSize: 20,
      value: 'ab',
    })
    expect(result).toEqual(payload)
  })

  it('getDevice GETs by number in path', async () => {
    const d = { id: 1, configurationId: 1, number: 'DEV001', groups: [] }
    get.mockResolvedValueOnce(ok(d))
    const result = await deviceService.getDevice('DEV001')
    expect(get).toHaveBeenCalledWith('/private/devices/number/DEV001')
    expect(result).toEqual(d)
  })

  it('deleteDevice DELETEs with numeric id', async () => {
    del.mockResolvedValueOnce({ data: { status: 'OK' } })
    await deviceService.deleteDevice(42)
    expect(del).toHaveBeenCalledWith('/private/devices/42')
  })

  it('propagates axios errors from getDevices', async () => {
    post.mockRejectedValueOnce(new Error('network'))
    await expect(deviceService.getDevices({ pageNum: 1, pageSize: 20 })).rejects.toThrow('network')
  })
})
