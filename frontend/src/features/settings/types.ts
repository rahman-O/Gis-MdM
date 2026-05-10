/** Global instance settings (UI model aligned with `.kiro/specs/settings-management/design.md`). */
export interface Settings {
  id: number
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
}

/** PUT-style payload (no id) — all fields required for submit. */
export interface SettingsPayload {
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
}

export interface ConfigurationOption {
  id: number
  name: string
}
