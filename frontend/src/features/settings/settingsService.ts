import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { assertHmdmOk, unwrapHmdmData } from '@/services/hmdmEnvelope'
import { apiLanguageToForm, formLanguageToApi } from '@/features/settings/languageMaps'
import type { Settings, SettingsPayload } from '@/features/settings/types'

function unwrap<T>(response: { data: HmdmEnvelope<T> }, fallback: string): T {
  return unwrapHmdmData(response.data, fallback)
}

/** Backend uses 0–2 (see `Settings.java` / legacy `password.service.js`). Values above 2 are clamped. */
function normalizePasswordStrength(raw: unknown): number {
  const n = raw === null || raw === undefined ? 0 : Number(raw)
  if (Number.isNaN(n)) return 0
  return Math.min(2, Math.max(0, Math.trunc(n)))
}

/** Normalizes backend JSON (unknown keys tolerated) into the UI Settings model. */
export function normalizeSettings(raw: Record<string, unknown>): Settings {
  const id = Number(raw.id ?? 0)
  return {
    id,
    customerName: String(raw.customerName ?? raw.name ?? ''),
    createNewDevices: Boolean(raw.createNewDevices ?? false),
    newDeviceConfigurationId:
      raw.newDeviceConfigurationId === null || raw.newDeviceConfigurationId === undefined
        ? null
        : Number(raw.newDeviceConfigurationId),
    newDeviceGroupId:
      raw.newDeviceGroupId === null || raw.newDeviceGroupId === undefined
        ? null
        : Number(raw.newDeviceGroupId),
    phoneNumberFormat: String(raw.phoneNumberFormat ?? '+9 (999) 999-99-99'),
    customPropertyName1: String(raw.customPropertyName1 ?? ''),
    customPropertyName2: String(raw.customPropertyName2 ?? ''),
    customPropertyName3: String(raw.customPropertyName3 ?? ''),
    customMultiline1: Boolean(raw.customMultiline1 ?? false),
    customMultiline2: Boolean(raw.customMultiline2 ?? false),
    customMultiline3: Boolean(raw.customMultiline3 ?? false),
    customSend1: Boolean(raw.customSend1 ?? false),
    customSend2: Boolean(raw.customSend2 ?? false),
    customSend3: Boolean(raw.customSend3 ?? false),
    desktopHeaderTemplate: String(raw.desktopHeaderTemplate ?? ''),
    sendDescription: Boolean(raw.sendDescription ?? false),
    language: apiLanguageToForm(raw.language == null ? undefined : String(raw.language)),
    passwordLength: raw.passwordLength === null || raw.passwordLength === undefined ? 0 : Number(raw.passwordLength),
    passwordStrength: normalizePasswordStrength(raw.passwordStrength),
    sendDeviceInfoExpiryDays:
      raw.sendDeviceInfoExpiryDays === null || raw.sendDeviceInfoExpiryDays === undefined
        ? 0
        : Number(raw.sendDeviceInfoExpiryDays),
    unsecureEnrollment: Boolean(raw.unsecureEnrollment ?? false),
    deviceFastSearch: Boolean(raw.deviceFastSearch ?? false),
    idleLogout:
      raw.idleLogout === null || raw.idleLogout === undefined
        ? null
        : Math.max(0, Math.trunc(Number(raw.idleLogout))),
  }
}

export async function getSettings(): Promise<Settings> {
  const response = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>('/private/settings')
  const data = unwrap(response, 'Failed to load settings.')
  return normalizeSettings(data)
}

/**
 * Persists settings using the same split as legacy Angular: `misc` then `language`.
 * Returns fresh settings from the server.
 */
export async function updateSettings(data: SettingsPayload): Promise<Settings> {
  if (data.createNewDevices && data.newDeviceConfigurationId == null) {
    throw new Error(
      'When "Create new devices on first access" is enabled, you must select a default configuration.'
    )
  }

  const snapshot = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>('/private/settings')
  const base = unwrap(snapshot, 'Failed to load settings for save.')

  const miscBody: Record<string, unknown> = {
    ...base,
    createNewDevices: data.createNewDevices,
    newDeviceConfigurationId: data.newDeviceConfigurationId,
    newDeviceGroupId: data.newDeviceGroupId,
    phoneNumberFormat: data.phoneNumberFormat,
    customPropertyName1: data.customPropertyName1,
    customPropertyName2: data.customPropertyName2,
    customPropertyName3: data.customPropertyName3,
    customMultiline1: data.customMultiline1,
    customMultiline2: data.customMultiline2,
    customMultiline3: data.customMultiline3,
    customSend1: data.customSend1,
    customSend2: data.customSend2,
    customSend3: data.customSend3,
    desktopHeaderTemplate: data.desktopHeaderTemplate,
    sendDescription: data.sendDescription,
    passwordLength: data.passwordLength,
    passwordStrength: normalizePasswordStrength(data.passwordStrength),
    customerName: data.customerName,
    sendDeviceInfoExpiryDays: data.sendDeviceInfoExpiryDays,
    unsecureEnrollment: data.unsecureEnrollment,
    deviceFastSearch: data.deviceFastSearch,
    idleLogout: data.idleLogout == null || data.idleLogout === 0 ? null : data.idleLogout,
  }
  const miscRes = await apiClient.post<HmdmEnvelope<unknown>>('/private/settings/misc', miscBody)
  assertHmdmOk(miscRes.data, 'Failed to save settings.')

  const useDefaultLanguage =
    typeof base.useDefaultLanguage === 'boolean' ? base.useDefaultLanguage : true

  const langBody: Record<string, unknown> = {
    ...base,
    ...miscBody,
    language: formLanguageToApi(data.language),
    useDefaultLanguage,
  }
  const langRes = await apiClient.post<HmdmEnvelope<unknown>>('/private/settings/lang', langBody)
  assertHmdmOk(langRes.data, 'Failed to save language settings.')

  return getSettings()
}

export async function fetchRawSettings(): Promise<Record<string, unknown>> {
  const response = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>('/private/settings')
  const data = unwrap(response, 'Failed to load settings.')
  return data && typeof data === 'object' ? data : {}
}

export async function saveDefaultDesign(body: Record<string, unknown>): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/settings/design', body)
  assertHmdmOk(response.data, 'Failed to save design settings.')
}

export interface UserRoleListRow {
  id: number
  name?: string | null
}

export async function listAssignableUserRoles(): Promise<UserRoleListRow[]> {
  const response = await apiClient.get<HmdmEnvelope<UserRoleListRow[]>>('/private/users/roles')
  const data = unwrap(response, 'Failed to load user roles.')
  return Array.isArray(data) ? data.filter((r) => r?.id != null) : []
}

export async function getUserRoleColumns(roleId: number): Promise<Record<string, unknown>> {
  const response = await apiClient.get<HmdmEnvelope<Record<string, unknown>>>(
    `/private/settings/userRole/${roleId}`
  )
  const data = unwrap(response, 'Failed to load role preferences.')
  return data && typeof data === 'object' ? data : {}
}

export async function saveUserRolesCommon(rows: Record<string, unknown>[]): Promise<void> {
  const response = await apiClient.post<HmdmEnvelope<unknown>>('/private/settings/userRoles/common', rows)
  assertHmdmOk(response.data, 'Failed to save role column settings.')
}
