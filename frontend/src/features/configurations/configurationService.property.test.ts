import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { arbitraryConfigurationPayload } from '@/features/configurations/test/arbitraries'

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

function ok<T>(data: T): { data: HmdmEnvelope<T> } {
  return { data: { status: 'OK', data } }
}

async function didReject(promise: Promise<unknown>): Promise<boolean> {
  try {
    await promise
    return false
  } catch {
    return true
  }
}

describe('configurationService property tests', () => {
  // Property 10 (Req 8.1, 8.2, 8.3, 8.4, 8.5): correct method+URL routing for operations.
  it('Property 10: service routes to correct URL for any operation', async () => {
    await fc.assert(
      fc.asyncProperty(fc.integer({ min: 1 }), arbitraryConfigurationPayload(), async (id, payload) => {
        mocks.get.mockReset()
        mocks.post.mockReset()
        mocks.put.mockReset()
        mocks.del.mockReset()

        mocks.get.mockResolvedValue(ok([]))
        await configurationService.getConfigurations()
        expect(mocks.get).toHaveBeenCalledWith('/private/configurations/search')

        mocks.get.mockResolvedValue(ok({ id }))
        await configurationService.getConfiguration(id)
        expect(mocks.get).toHaveBeenCalledWith(`/private/configurations/${id}`)

        mocks.get.mockResolvedValue(ok([]))
        await configurationService.searchConfigurations(`q-${id}`)
        expect(mocks.get).toHaveBeenCalledWith(`/private/configurations/search/${encodeURIComponent(`q-${id}`)}`)

        mocks.get.mockResolvedValue(ok([]))
        await configurationService.listConfigurationNames()
        expect(mocks.get).toHaveBeenCalledWith('/private/configurations/list')

        mocks.post.mockResolvedValue(ok([]))
        await configurationService.autocompleteConfigurations({ value: `q-${id}` })
        expect(mocks.post).toHaveBeenCalledWith('/private/configurations/autocomplete', `q-${id}`)

        mocks.put.mockResolvedValue(ok({ id }))
        await configurationService.createConfiguration(payload)
        expect(mocks.put).toHaveBeenCalledWith(
          '/private/configurations',
          expect.objectContaining({ name: payload.name.trim() })
        )

        mocks.get.mockResolvedValue(ok({ id, name: 'existing', type: 0, applications: [] }))
        mocks.put.mockResolvedValue(ok({ id }))
        await configurationService.updateConfiguration(id, payload)
        expect(mocks.put).toHaveBeenCalledWith(
          '/private/configurations',
          expect.objectContaining({ id })
        )

        mocks.del.mockResolvedValue({ data: { status: 'OK' } })
        await configurationService.deleteConfiguration(id)
        expect(mocks.del).toHaveBeenCalledWith(`/private/configurations/${id}`)

        mocks.put.mockResolvedValue({ data: { status: 'OK' } })
        await configurationService.copyConfiguration({ id, name: `copy-${id}` })
        expect(mocks.put).toHaveBeenCalledWith('/private/configurations/copy', {
          id,
          name: `copy-${id}`,
        })

        mocks.get.mockResolvedValue(ok([]))
        await configurationService.getConfigurationApplications(id)
        expect(mocks.get).toHaveBeenCalledWith(`/private/configurations/applications/${id}`)

        mocks.get.mockResolvedValue(ok([]))
        await configurationService.getAllApplications()
        expect(mocks.get).toHaveBeenCalledWith('/private/configurations/applications')

        mocks.put.mockResolvedValue(ok({ id }))
        await configurationService.upgradeConfigurationApplication({
          configurationId: id,
          applicationId: id + 1,
        })
        expect(mocks.put).toHaveBeenCalledWith('/private/configurations/application/upgrade', {
          configurationId: id,
          applicationId: id + 1,
        })
      }),
      { numRuns: 100 }
    )
  })

  // Property 11 (Req 8.7): service propagates rejected requests.
  it('Property 11: service error propagation', async () => {
    await fc.assert(
      fc.asyncProperty(
        fc.constantFrom(
          'list',
          'getOne',
          'search',
          'listNames',
          'autocomplete',
          'create',
          'update',
          'delete',
          'copy',
          'getCfgApps',
          'getAllApps',
          'upgrade'
        ),
        fc.integer({ min: 1 }),
        arbitraryConfigurationPayload(),
        async (op, id, payload) => {
          const err = new Error('boom')
          mocks.get.mockReset()
          mocks.post.mockReset()
          mocks.put.mockReset()
          mocks.del.mockReset()

          if (op === 'list') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.getConfigurations())).toBe(true)
            return
          }

          if (op === 'getOne') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.getConfiguration(id))).toBe(true)
            return
          }

          if (op === 'create') {
            mocks.put.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.createConfiguration(payload))).toBe(true)
            return
          }

          if (op === 'search') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.searchConfigurations(payload.name))).toBe(true)
            return
          }

          if (op === 'listNames') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.listConfigurationNames())).toBe(true)
            return
          }

          if (op === 'autocomplete') {
            mocks.post.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.autocompleteConfigurations({ value: payload.name }))).toBe(true)
            return
          }

          if (op === 'update') {
            mocks.get.mockResolvedValueOnce(ok({ id, name: 'x', type: 0, applications: [] }))
            mocks.put.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.updateConfiguration(id, payload))).toBe(true)
            return
          }

          if (op === 'copy') {
            mocks.put.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.copyConfiguration({ id, name: payload.name }))).toBe(true)
            return
          }

          if (op === 'getCfgApps') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.getConfigurationApplications(id))).toBe(true)
            return
          }

          if (op === 'getAllApps') {
            mocks.get.mockRejectedValueOnce(err)
            expect(await didReject(configurationService.getAllApplications())).toBe(true)
            return
          }

          if (op === 'upgrade') {
            mocks.put.mockRejectedValueOnce(err)
            expect(
              await didReject(
                configurationService.upgradeConfigurationApplication({
                  configurationId: id,
                  applicationId: id + 1,
                })
              )
            ).toBe(true)
            return
          }

          mocks.del.mockRejectedValueOnce(err)
          expect(await didReject(configurationService.deleteConfiguration(id))).toBe(true)
        }
      ),
      { numRuns: 100 }
    )
  })
})
