import { beforeEach, describe, expect, it, vi } from 'vitest'
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

import * as applicationService from '@/features/applications/services/applicationService'

function ok<T>(data: T): { data: HmdmEnvelope<T> } {
  return { data: { status: 'OK', data } }
}

describe('applicationService', () => {
  beforeEach(() => {
    mocks.get.mockReset()
    mocks.post.mockReset()
    mocks.put.mockReset()
    mocks.del.mockReset()
  })

  it('loads all applications', async () => {
    mocks.get.mockResolvedValueOnce(ok([{ id: 1, name: 'A' }]))
    await expect(applicationService.getAllApplications()).resolves.toHaveLength(1)
    expect(mocks.get).toHaveBeenCalledWith('/private/applications/search')
  })

  it('saves android app', async () => {
    mocks.put.mockResolvedValueOnce(ok({ id: 1 }))
    await applicationService.createOrUpdateAndroidApplication({ name: 'A', pkg: 'x.y', type: 'app' })
    expect(mocks.put).toHaveBeenCalledWith('/private/applications/android', expect.any(Object))
  })

  it('updates app configurations', async () => {
    mocks.post.mockResolvedValueOnce({ data: { status: 'OK' } })
    await applicationService.updateApplicationConfigurations({
      applicationId: 1,
      configurations: [{ configurationId: 2, action: 1 }],
    })
    expect(mocks.post).toHaveBeenCalledWith('/private/applications/configurations', expect.any(Object))
  })

  it('loads admin applications (super-admin catalog)', async () => {
    mocks.get.mockResolvedValueOnce(ok([{ id: 9, name: 'AdminApp' }]))
    await expect(applicationService.getAllAdminApplications()).resolves.toHaveLength(1)
    expect(mocks.get).toHaveBeenCalledWith('/private/applications/admin/search')
  })

  it('searches admin applications with path param', async () => {
    mocks.get.mockResolvedValueOnce(ok([]))
    await applicationService.searchAdminApplications('com.example')
    expect(mocks.get).toHaveBeenCalledWith('/private/applications/admin/search/com.example')
  })

  it('delegates turn-into-common via GET', async () => {
    mocks.get.mockResolvedValueOnce({ data: { status: 'OK' } })
    await applicationService.turnApplicationIntoCommon(42)
    expect(mocks.get).toHaveBeenCalledWith('/private/applications/admin/common/42')
  })
})
