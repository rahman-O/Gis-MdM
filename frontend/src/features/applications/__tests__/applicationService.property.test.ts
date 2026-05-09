import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
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

describe('applicationService property tests', () => {
  it('routes searches correctly for any value', async () => {
    await fc.assert(
      fc.asyncProperty(fc.string({ minLength: 1 }), async (value) => {
        mocks.get.mockReset()
        mocks.get.mockResolvedValueOnce(ok([]))
        await applicationService.searchApplications(value)
        expect(mocks.get).toHaveBeenCalledWith(
          `/private/applications/search/${encodeURIComponent(value.trim())}`
        )
      }),
      { numRuns: 100 }
    )
  })
})
