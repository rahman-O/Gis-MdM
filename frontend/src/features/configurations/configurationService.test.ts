import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { AxiosResponse } from 'axios'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'

const mocks = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}))

vi.mock('@/services/apiClient', () => ({
  default: {
    get: mocks.get,
    post: mocks.post,
    put: mocks.put,
    delete: mocks.del,
  },
}))

import * as configurationService from '@/features/configurations/configurationService'

const { get, post, put, del } = mocks

function ok<T>(data: T): AxiosResponse<HmdmEnvelope<T>> {
  return { data: { status: 'OK', data } } as AxiosResponse<HmdmEnvelope<T>>
}

describe('configurationService', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    put.mockReset()
    del.mockReset()
  })

  it('getConfigurations GET /private/configurations/search', async () => {
    const list = [{ id: 1, name: 'A', type: 0 }]
    get.mockResolvedValueOnce(ok(list))
    await expect(configurationService.getConfigurations()).resolves.toEqual(list)
    expect(get).toHaveBeenCalledWith('/private/configurations/search')
  })

  it('getConfiguration GET by id', async () => {
    const c = { id: 2, name: 'B', type: 1 }
    get.mockResolvedValueOnce(ok(c))
    await expect(configurationService.getConfiguration(2)).resolves.toEqual(c)
    expect(get).toHaveBeenCalledWith('/private/configurations/2')
  })

  it('searchConfigurations GET by value', async () => {
    const list = [{ id: 4, name: 'A-1', type: 0 }]
    get.mockResolvedValueOnce(ok(list))
    await expect(configurationService.searchConfigurations('A 1')).resolves.toEqual(list)
    expect(get).toHaveBeenCalledWith('/private/configurations/search/A%201')
  })

  it('listConfigurationNames GET /private/configurations/list', async () => {
    const list = [{ id: 1, name: 'A' }]
    get.mockResolvedValueOnce(ok(list))
    await expect(configurationService.listConfigurationNames()).resolves.toEqual(list)
    expect(get).toHaveBeenCalledWith('/private/configurations/list')
  })

  it('autocompleteConfigurations POST /private/configurations/autocomplete', async () => {
    const list = [{ id: 1, name: 'A' }]
    post.mockResolvedValueOnce(ok(list))
    await expect(configurationService.autocompleteConfigurations({ value: 'A' })).resolves.toEqual(list)
    expect(post).toHaveBeenCalledWith('/private/configurations/autocomplete', 'A')
  })

  it('createConfiguration PUT body includes name and defaults', async () => {
    put.mockResolvedValueOnce(ok({ id: 9, name: 'New', type: 0 }))
    await configurationService.createConfiguration({
      name: 'New',
      description: null,
      type: 'WORK',
    })
    expect(put).toHaveBeenCalledWith(
      '/private/configurations',
      expect.objectContaining({
        name: 'New',
        type: 0,
        applications: [],
        iconSize: 'SMALL',
      })
    )
  })

  it('updateConfiguration merges payload onto existing', async () => {
    const existing = { id: 3, name: 'Old', type: 1, applications: [], kioskMode: false }
    get.mockResolvedValueOnce(ok(existing))
    put.mockResolvedValueOnce(ok({ ...existing, name: 'Renamed', type: 0 }))
    await configurationService.updateConfiguration(3, {
      name: 'Renamed',
      description: 'd',
      type: 'WORK',
    })
    expect(put).toHaveBeenCalledWith(
      '/private/configurations',
      expect.objectContaining({
        id: 3,
        name: 'Renamed',
        description: 'd',
        type: 0,
      })
    )
  })

  it('deleteConfiguration DELETE by id', async () => {
    del.mockResolvedValueOnce({ data: { status: 'OK' } })
    await configurationService.deleteConfiguration(5)
    expect(del).toHaveBeenCalledWith('/private/configurations/5')
  })

  it('copyConfiguration PUT /copy', async () => {
    put.mockResolvedValueOnce({ data: { status: 'OK' } })
    await configurationService.copyConfiguration({ id: 7, name: 'Copy' })
    expect(put).toHaveBeenCalledWith('/private/configurations/copy', { id: 7, name: 'Copy' })
  })

  it('getConfigurationApplications GET /applications/{id}', async () => {
    const apps = [{ id: 11, name: 'App A' }]
    get.mockResolvedValueOnce(ok(apps))
    await expect(configurationService.getConfigurationApplications(3)).resolves.toEqual(apps)
    expect(get).toHaveBeenCalledWith('/private/configurations/applications/3')
  })

  it('getAllApplications GET /applications', async () => {
    const apps = [{ id: 11, name: 'App A' }]
    get.mockResolvedValueOnce(ok(apps))
    await expect(configurationService.getAllApplications()).resolves.toEqual(apps)
    expect(get).toHaveBeenCalledWith('/private/configurations/applications')
  })

  it('upgradeConfigurationApplication PUT /application/upgrade', async () => {
    const payload = { configurationId: 5, applicationId: 9 }
    put.mockResolvedValueOnce(ok({ id: 5 }))
    await expect(configurationService.upgradeConfigurationApplication(payload)).resolves.toEqual({ id: 5 })
    expect(put).toHaveBeenCalledWith('/private/configurations/application/upgrade', payload)
  })
})
