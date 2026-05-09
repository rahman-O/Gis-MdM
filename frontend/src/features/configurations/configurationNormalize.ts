import type { Configuration, ConfigurationPayload } from '@/features/configurations/types'

export function normalizeConfigurationPayload(payload: ConfigurationPayload): ConfigurationPayload {
  const description = payload.description?.trim()
  return {
    name: payload.name.trim(),
    type: payload.type,
    description: description ? description : null,
  }
}

/** Minimal defaults so inserts/updates remain backend-compatible. */
export function buildCreateConfigurationBody(payload: ConfigurationPayload): Configuration {
  const normalized = normalizeConfigurationPayload(payload)
  return {
    name: normalized.name,
    description: normalized.description,
    type: normalized.type === 'COMMON' ? 1 : 0,
    applications: [],
    iconSize: 'SMALL',
    desktopHeader: 'NO_HEADER',
    useDefaultDesignSettings: true,
    pushOptions: 'mqttWorker',
    defaultFilePath: '/',
    downloadUpdates: 'UNLIMITED',
    requestUpdates: 'DONOTTRACK',
    appPermissions: 'GRANTALL',
    kioskMode: false,
    encryptDevice: false,
    mobileEnrollment: false,
    displayStatus: false,
    blockStatusBar: false,
    systemUpdateType: 0,
    scheduleAppUpdate: false,
    disableLocation: false,
    permissive: false,
    selected: false,
  }
}

export function mergeConfigurationForUpdate(
  current: Configuration,
  id: number,
  payload: ConfigurationPayload
): Configuration {
  const normalized = normalizeConfigurationPayload(payload)
  return {
    ...current,
    id,
    name: normalized.name,
    description: normalized.description,
    type: normalized.type === 'COMMON' ? 1 : 0,
  }
}
