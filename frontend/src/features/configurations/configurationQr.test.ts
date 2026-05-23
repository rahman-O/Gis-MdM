import { describe, expect, it } from 'vitest'
import { getConfigurationQrEligibility } from '@/features/configurations/configurationQr'

describe('configurationQr', () => {
  it('returns ineligible when qrCodeKey is missing', () => {
    const result = getConfigurationQrEligibility({
      qrCodeKey: '',
      mainAppId: 10,
      eventReceivingComponent: 'a/b',
    })
    expect(result.eligible).toBe(false)
    expect(result.reason).toContain('no QR key')
  })

  it('returns ineligible when mainAppId is missing', () => {
    const result = getConfigurationQrEligibility({
      qrCodeKey: 'k',
      mainAppId: null,
      eventReceivingComponent: 'a/b',
    })
    expect(result.eligible).toBe(false)
    expect(result.reason).toContain('no Main App')
  })

  it('returns eligible when all prerequisites are present', () => {
    const result = getConfigurationQrEligibility({
      qrCodeKey: 'k',
      mainAppId: 3,
      eventReceivingComponent: 'pkg/.Receiver',
    })
    expect(result).toEqual({ eligible: true, reason: null })
  })
})
