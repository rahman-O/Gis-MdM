/** Global instance settings (UI model aligned with `.kiro/specs/settings-management/design.md`). */
export interface Settings {
  id: number
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  newDeviceGroupId: number | null
  phoneNumberFormat: string
  customPropertyName1: string
  customPropertyName2: string
  customPropertyName3: string
  customMultiline1: boolean
  customMultiline2: boolean
  customMultiline3: boolean
  customSend1: boolean
  customSend2: boolean
  customSend3: boolean
  desktopHeaderTemplate: string
  sendDescription: boolean
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
  /** Seconds before idle logout (`Settings.idleLogout`; `null`/0 disables in UI). */
  idleLogout: number | null
}

/** PUT-style payload (no id) — all fields required for submit. */
export interface SettingsPayload {
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  newDeviceGroupId: number | null
  phoneNumberFormat: string
  customPropertyName1: string
  customPropertyName2: string
  customPropertyName3: string
  customMultiline1: boolean
  customMultiline2: boolean
  customMultiline3: boolean
  customSend1: boolean
  customSend2: boolean
  customSend3: boolean
  desktopHeaderTemplate: string
  sendDescription: boolean
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
  idleLogout: number | null
}

export interface ConfigurationOption {
  id: number
  name: string
}

export interface GroupOption {
  id: number
  name: string
}
