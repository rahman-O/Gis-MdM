import type { Configuration } from '@/features/configurations/types'

export interface ConfigurationQrEligibility {
  eligible: boolean
  reason: string | null
}

export function getConfigurationQrEligibility(
  configuration: Pick<Configuration, 'qrCodeKey' | 'mainAppId' | 'eventReceivingComponent'> | null | undefined
): ConfigurationQrEligibility {
  if (!configuration) {
    return { eligible: false, reason: 'QR is unavailable: configuration data could not be loaded.' }
  }

  const qrCodeKey = String(configuration.qrCodeKey ?? '').trim()
  if (!qrCodeKey) {
    return { eligible: false, reason: 'QR is unavailable: this configuration has no QR key.' }
  }

  const mainAppId = Number(configuration.mainAppId ?? 0)
  if (!mainAppId || mainAppId <= 0) {
    return { eligible: false, reason: 'QR is unavailable: this configuration has no Main App assigned.' }
  }

  const eventReceivingComponent = String(configuration.eventReceivingComponent ?? '').trim()
  if (!eventReceivingComponent) {
    return { eligible: false, reason: 'QR is unavailable: Event Receiving Component is not configured.' }
  }

  return { eligible: true, reason: null }
}
