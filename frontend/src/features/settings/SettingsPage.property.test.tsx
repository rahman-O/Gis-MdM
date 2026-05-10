import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import {
  arbitraryConfigurationOption,
  arbitraryInvalidCustomerName,
  arbitraryInvalidPositiveInt,
  arbitrarySettings,
  arbitrarySettingsPayload,
  arbitrarySettingsWithNulls,
} from '@/features/settings/test/arbitraries'
import { SettingsPage } from '@/features/settings/SettingsPage'
import type { SettingsPayload } from '@/features/settings/types'

const PASSWORD_STRENGTH_LABELS = ['Any (length only)', 'Digits, upper & lower case', 'Above + special characters']

const toastMock = vi.fn()

vi.mock('@/shared/hooks/use-toast', () => ({
  useToast: () => ({ toast: toastMock }),
}))

vi.mock('@/services/apiClient', () => ({
  default: {
    get: vi.fn(),
  },
}))

vi.mock('@/features/settings/settingsService', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/features/settings/settingsService')>()
  return {
    ...actual,
    getSettings: vi.fn(),
    updateSettings: vi.fn(),
  }
})

import apiClient from '@/services/apiClient'
import * as settingsService from '@/features/settings/settingsService'

function ok<T>(data: T): { data: HmdmEnvelope<T> } {
  return { data: { status: 'OK', data } }
}

describe('SettingsPage property tests', () => {
  // Feature: settings-management, Property 1: Form populated with any settings object
  it('Property 1: form populated with any settings object', async () => {
    await fc.assert(
      fc.asyncProperty(arbitrarySettings(), async (s) => {
        cleanup()
        toastMock.mockReset()
        vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue(s)
        vi.mocked(apiClient.get).mockReset().mockResolvedValue(ok([]))

        render(<SettingsPage />)
        await waitFor(() =>
          expect(screen.queryByText('Instance-wide preferences and security defaults.')).toBeInTheDocument()
        )

        expect(screen.getByLabelText(/customer name/i)).toHaveValue(s.customerName)
        expect(screen.getByLabelText(/minimum password length/i)).toHaveValue(s.passwordLength)
        expect(screen.getByLabelText(/send device info expiry/i)).toHaveValue(s.sendDeviceInfoExpiryDays)
      }),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: settings-management, Property 2: Configurations populate select options
  it('Property 2: configurations populate select options', async () => {
    await fc.assert(
      fc.asyncProperty(
        fc.array(arbitraryConfigurationOption(), { minLength: 1, maxLength: 6 }),
        async (configs) => {
          cleanup()
          vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue({
            id: 1,
            customerName: 'Test',
            createNewDevices: false,
            newDeviceConfigurationId: null,
            language: 'en',
            passwordLength: 8,
            passwordStrength: 0,
            sendDeviceInfoExpiryDays: 7,
            unsecureEnrollment: false,
            deviceFastSearch: false,
          })
          vi.mocked(apiClient.get).mockReset().mockImplementation((url: string) => {
            if (url === '/private/configurations/search') {
              return Promise.resolve(
                ok(
                  configs.map((c) => ({
                    id: c.id,
                    name: c.name,
                  }))
                )
              )
            }
            return Promise.reject(new Error(`Unexpected GET ${url}`))
          })

          const user = userEvent.setup()
          render(<SettingsPage />)
          await waitFor(() => expect(screen.getByText(/default configuration for new devices/i)).toBeInTheDocument())

          const combos = screen.getAllByRole('combobox')
          const configCombo = combos[0]
          await user.click(configCombo)

          for (const c of configs) {
            await waitFor(() => expect(screen.getByRole('option', { name: c.name })).toBeInTheDocument())
            const opt = screen.getByRole('option', { name: c.name })
            const dv = opt.getAttribute('data-value')
            if (dv != null) {
              expect(dv).toBe(String(c.id))
            }
          }
        }
      ),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: settings-management, Property 3: Submit calls service with form values
  it('Property 3: submit calls service with form values', async () => {
    const payloadArb = arbitrarySettingsPayload().filter(
      (p) =>
        !p.createNewDevices &&
        p.newDeviceConfigurationId === null &&
        p.language === 'en' &&
        p.passwordStrength === 0
    )
    await fc.assert(
      fc.asyncProperty(payloadArb, async (payload) => {
        cleanup()
        toastMock.mockReset()
        vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue({
          id: 1,
          customerName: '',
          createNewDevices: false,
          newDeviceConfigurationId: null,
          language: 'en',
          passwordLength: 1,
          passwordStrength: 0,
          sendDeviceInfoExpiryDays: 1,
          unsecureEnrollment: false,
          deviceFastSearch: false,
        })
        vi.mocked(settingsService.updateSettings).mockReset().mockResolvedValue({
          id: 1,
          ...payload,
        })
        vi.mocked(apiClient.get).mockReset().mockResolvedValue(ok([]))

        const user = userEvent.setup()
        render(<SettingsPage />)
        await waitFor(() => expect(screen.getByLabelText(/customer name/i)).toBeInTheDocument())

        await applySettingsPayload(user, payload)

        await user.click(screen.getByRole('button', { name: /save settings/i }))
        await waitFor(() => expect(settingsService.updateSettings).toHaveBeenCalled())
        expect(settingsService.updateSettings).toHaveBeenCalledWith(payload)
      }),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: settings-management, Property 4: Save button re-enabled after any PUT outcome
  it('Property 4: save button re-enabled after any PUT outcome', async () => {
    await fc.assert(
      fc.asyncProperty(fc.boolean(), async (shouldFail) => {
        cleanup()
        vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue({
          id: 1,
          customerName: 'Acme',
          createNewDevices: false,
          newDeviceConfigurationId: null,
          language: 'en',
          passwordLength: 8,
          passwordStrength: 0,
          sendDeviceInfoExpiryDays: 7,
          unsecureEnrollment: false,
          deviceFastSearch: false,
        })
        vi.mocked(settingsService.updateSettings).mockReset().mockImplementation(() =>
          shouldFail ? Promise.reject(new Error('x')) : Promise.resolve({ id: 1, customerName: 'Acme', createNewDevices: false, newDeviceConfigurationId: null, language: 'en', passwordLength: 8, passwordStrength: 0, sendDeviceInfoExpiryDays: 7, unsecureEnrollment: false, deviceFastSearch: false })
        )
        vi.mocked(apiClient.get).mockReset().mockResolvedValue(ok([]))

        const user = userEvent.setup()
        render(<SettingsPage />)
        await waitFor(() => expect(screen.getByRole('button', { name: /save settings/i })).toBeEnabled())
        await user.click(screen.getByRole('button', { name: /save settings/i }))
        await waitFor(() => expect(screen.getByRole('button', { name: /save settings/i })).toBeEnabled())
      }),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: settings-management, Property 5: Validation rejects invalid required fields
  it('Property 5: validation rejects invalid required fields', async () => {
    await fc.assert(
      fc.asyncProperty(
        arbitraryInvalidCustomerName(),
        arbitraryInvalidPositiveInt(),
        arbitraryInvalidPositiveInt(),
        async (badName, badPwLen, badExpiry) => {
          cleanup()
          vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue({
            id: 1,
            customerName: 'ok',
            createNewDevices: false,
            newDeviceConfigurationId: null,
            language: 'en',
            passwordLength: 8,
            passwordStrength: 0,
            sendDeviceInfoExpiryDays: 7,
            unsecureEnrollment: false,
            deviceFastSearch: false,
          })
          vi.mocked(settingsService.updateSettings).mockReset()
          vi.mocked(apiClient.get).mockReset().mockResolvedValue(ok([]))

          const user = userEvent.setup()
          render(<SettingsPage />)
          await waitFor(() => expect(screen.getByLabelText(/customer name/i)).toBeInTheDocument())

          fireEvent.input(screen.getByLabelText(/customer name/i), { target: { value: badName } })
          fireEvent.input(screen.getByLabelText(/minimum password length/i), { target: { value: String(badPwLen) } })
          fireEvent.input(screen.getByLabelText(/send device info expiry/i), { target: { value: String(badExpiry) } })

          await user.click(screen.getByRole('button', { name: /save settings/i }))
          await waitFor(() => expect(settingsService.updateSettings).not.toHaveBeenCalled())
        }
      ),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: settings-management, Property 8: Null-safe rendering
  it('Property 8: null-safe rendering for any settings with null fields', async () => {
    await fc.assert(
      fc.asyncProperty(arbitrarySettingsWithNulls(), async (raw) => {
        cleanup()
        const normalized = settingsService.normalizeSettings(raw)
        vi.mocked(settingsService.getSettings).mockReset().mockResolvedValue(normalized)
        vi.mocked(apiClient.get).mockReset().mockResolvedValue(ok([]))

        expect(() => render(<SettingsPage />)).not.toThrow()
        await waitFor(() => expect(screen.getByLabelText(/minimum password length/i)).toHaveValue(0))
      }),
      { numRuns: 100 }
    )
  }, 120000)
})

async function applySettingsPayload(user: ReturnType<typeof userEvent.setup>, p: SettingsPayload) {
  await user.clear(screen.getByLabelText(/customer name/i))
  await user.type(screen.getByLabelText(/customer name/i), p.customerName)

  const switches = screen.getAllByRole('switch')
  if (p.createNewDevices) await user.click(switches[0])
  if (p.unsecureEnrollment) await user.click(switches[1])
  if (p.deviceFastSearch) await user.click(switches[2])

  await user.clear(screen.getByLabelText(/minimum password length/i))
  await user.type(screen.getByLabelText(/minimum password length/i), String(p.passwordLength))

  await user.clear(screen.getByLabelText(/send device info expiry/i))
  await user.type(screen.getByLabelText(/send device info expiry/i), String(p.sendDeviceInfoExpiryDays))

  if (p.language !== 'en') {
    await user.click(screen.getAllByRole('combobox')[1])
    const label = { en: 'English', ru: 'Russian', de: 'German', fr: 'French', es: 'Spanish', pt: 'Portuguese', zh: 'Chinese' }[p.language]
    await user.click(screen.getByRole('option', { name: label ?? p.language }))
  }

  if (p.passwordStrength !== 0) {
    await user.click(screen.getAllByRole('combobox')[2])
    await user.click(screen.getByRole('option', { name: PASSWORD_STRENGTH_LABELS[p.passwordStrength] }))
  }
}
