import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { arbitrarySettingsPayload } from '@/features/settings/test/arbitraries'
import type { SettingsPayload } from '@/features/settings/types'

const mocks = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
}))

vi.mock('@/services/apiClient', () => ({
  default: {
    get: mocks.get,
    post: mocks.post,
  },
}))

import * as settingsService from '@/features/settings/settingsService'

function ok<T>(data: T): { data: HmdmEnvelope<T> } {
  return { data: { status: 'OK', data } }
}

const sampleBase: Record<string, unknown> = {
  id: 1,
  customerId: 1,
  customerName: 'Acme',
  createNewDevices: false,
  newDeviceConfigurationId: null,
  language: 'en_US',
  passwordLength: 8,
  passwordStrength: 0,
  useDefaultLanguage: true,
}

const minimalPayload: SettingsPayload = {
  customerName: 'Acme',
  createNewDevices: false,
  newDeviceConfigurationId: null,
  language: 'en',
  passwordLength: 8,
  passwordStrength: 0,
  sendDeviceInfoExpiryDays: 7,
  unsecureEnrollment: false,
  deviceFastSearch: false,
}

async function didReject(promise: Promise<unknown>): Promise<boolean> {
  try {
    await promise
    return false
  } catch {
    return true
  }
}

/** Matches legacy + service guard: cannot enable new devices without a default configuration id. */
function isValidSettingsPayload(p: SettingsPayload): boolean {
  return !p.createNewDevices || p.newDeviceConfigurationId != null
}

describe('settingsService property tests', () => {
  it('rejects update when createNewDevices is set but default configuration is missing', async () => {
    mocks.get.mockReset()
    mocks.post.mockReset()
    await expect(
      settingsService.updateSettings({
        ...minimalPayload,
        createNewDevices: true,
        newDeviceConfigurationId: null,
      })
    ).rejects.toThrow(/default configuration/i)
    expect(mocks.get).not.toHaveBeenCalled()
    expect(mocks.post).not.toHaveBeenCalled()
  })

  it('preserves useDefaultLanguage from GET snapshot on lang POST', async () => {
    mocks.get.mockReset()
    mocks.post.mockReset()
    mocks.get
      .mockResolvedValueOnce(ok({ ...sampleBase, useDefaultLanguage: false }))
      .mockResolvedValueOnce(ok({ ...sampleBase, useDefaultLanguage: false, language: 'en_US' }))
    mocks.post.mockResolvedValue({ data: { status: 'OK' } })

    await settingsService.updateSettings(minimalPayload)

    const langBody = mocks.post.mock.calls.find((c) => c[0] === '/private/settings/lang')?.[1] as Record<
      string,
      unknown
    >
    expect(langBody?.useDefaultLanguage).toBe(false)
  })

  // Feature: settings-management, Property 6: Service routes to correct URLs
  it('Property 6: service routes to correct URLs', async () => {
    await fc.assert(
      fc.asyncProperty(arbitrarySettingsPayload().filter(isValidSettingsPayload), async (payload) => {
        mocks.get.mockReset()
        mocks.post.mockReset()

        mocks.get.mockResolvedValueOnce(ok({ ...sampleBase }))
        await settingsService.getSettings()
        expect(mocks.get).toHaveBeenCalledWith('/private/settings')

        mocks.get.mockResolvedValueOnce(ok({ ...sampleBase }))
        mocks.post.mockResolvedValue(ok({ status: 'OK' }))
        mocks.get.mockResolvedValueOnce(ok({ ...sampleBase, ...payload, language: 'en_US' }))
        await settingsService.updateSettings(payload)

        expect(mocks.post).toHaveBeenCalledWith('/private/settings/misc', expect.anything())
        expect(mocks.post).toHaveBeenCalledWith('/private/settings/lang', expect.anything())
        expect(mocks.get.mock.calls.length).toBeGreaterThanOrEqual(2)
      }),
      { numRuns: 100 }
    )
  })

  // Feature: settings-management, Property 7: Service error propagation
  it('Property 7: service error propagation', async () => {
    await fc.assert(
      fc.asyncProperty(fc.constantFrom('get', 'update'), arbitrarySettingsPayload().filter(isValidSettingsPayload), async (op, payload) => {
        const err = new Error('boom')
        mocks.get.mockReset()
        mocks.post.mockReset()

        if (op === 'get') {
          mocks.get.mockRejectedValueOnce(err)
          expect(await didReject(settingsService.getSettings())).toBe(true)
          return
        }

        mocks.get.mockResolvedValue(ok({ ...sampleBase }))
        mocks.post.mockRejectedValueOnce(err)
        expect(await didReject(settingsService.updateSettings(payload))).toBe(true)
      }),
      { numRuns: 100 }
    )
  })
})
