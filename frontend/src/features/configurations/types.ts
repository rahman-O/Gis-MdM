/** Device configuration kind shown in UI — maps to backend `type` int (0 = device/work, 1 = typical/common). */
export type ConfigurationKind = 'WORK' | 'COMMON'

/** Payload for create/update dialogs (subset of full `Configuration`). */
export interface ConfigurationPayload {
  name: string
  description: string | null
  type: ConfigurationKind
}

export interface ConfigurationApplication {
  id?: number
  name?: string | null
  pkg?: string | null
  type?: string | null
  action?: number | null
  version?: string | null
  url?: string | null
  versionCode?: number | null
  oldVersion?: string | null
  oldVersionCode?: number | null
  selected?: boolean
  [key: string]: unknown
}

export interface ConfigurationFile {
  id?: number
  fileId?: number
  url?: string | null
  externalUrl?: string | null
  path?: string | null
  checksum?: string | null
  remove?: boolean
  variables?: string[] | null
  [key: string]: unknown
}

export interface ConfigurationApplicationSetting {
  id?: number | null
  applicationId?: number | null
  applicationName?: string | null
  name?: string | null
  type?: string | null
  value?: string | null
}

export interface ConfigurationLookupItem {
  id: number
  name: string
}

export interface ConfigurationAutocompleteRequest {
  value: string
}

export interface CopyConfigurationPayload {
  id: number
  name: string
  description?: string | null
}

export interface UpgradeConfigurationApplicationPayload {
  configurationId: number
  applicationId: number
}

/**
 * Full configuration object from Headwind REST (`PUT /private/configurations`).
 * Many optional fields exist server-side; UI list mainly uses id, name, description, type, deviceCount.
 */
export interface Configuration {
  id?: number | null
  name?: string | null
  description?: string | null
  /** Backend: 0 = regular (WORK), 1 = typical template (COMMON). */
  type?: number
  deviceCount?: number | null
  password?: string | null
  pushOptions?: string | null
  requestUpdates?: string | null
  appPermissions?: string | null
  downloadUpdates?: string | null
  passwordMode?: string | null
  iconSize?: string | null
  desktopHeader?: string | null
  desktopHeaderText?: string | null
  desktopHeaderTemplate?: string | null
  orientation?: string | null
  useDefaultDesignSettings?: boolean | null
  backgroundColor?: string | null
  textColor?: string | null
  backgroundImageUrl?: string | null
  displayStatus?: boolean | null
  defaultFilePath?: string | null
  qrCodeKey?: string | null
  baseUrl?: string | null
  mainAppId?: number | null
  contentAppId?: number | null
  eventReceivingComponent?: string | null
  qrParameters?: string | null
  adminExtras?: string | null
  launcherUrl?: string | null
  wifiSSID?: string | null
  wifiPassword?: string | null
  wifiSecurityType?: string | null
  mobileEnrollment?: boolean | null
  encryptDevice?: boolean | null
  kioskMode?: boolean | null
  gps?: boolean | null
  bluetooth?: boolean | null
  wifi?: boolean | null
  mobileData?: boolean | null
  usbStorage?: boolean | null
  autoBrightness?: boolean | null
  brightness?: number | null
  manageTimeout?: boolean | null
  timeout?: number | null
  lockVolume?: boolean | null
  manageVolume?: boolean | null
  volume?: number | null
  keepaliveTime?: number | null
  timeZone?: string | null
  timeZoneMode?: string | null
  systemUpdateType?: number | null
  scheduleAppUpdate?: boolean | null
  runDefaultLauncher?: boolean | null
  showWifi?: boolean | null
  disableScreenshots?: boolean | null
  autostartForeground?: boolean | null
  systemUpdateFrom?: string | null
  systemUpdateTo?: string | null
  appUpdateFrom?: string | null
  appUpdateTo?: string | null
  kioskHome?: boolean | null
  kioskRecents?: boolean | null
  kioskNotifications?: boolean | null
  kioskSystemInfo?: boolean | null
  kioskKeyguard?: boolean | null
  kioskLockButtons?: boolean | null
  kioskScreenOn?: boolean | null
  kioskExit?: boolean | null
  permissive?: boolean | null
  lockSafeSettings?: boolean | null
  allowedClasses?: string | null
  restrictions?: string | null
  newServerUrl?: string | null
  policyLocks?: Record<string, boolean> | null
  applications?: ConfigurationApplication[]
  applicationSettings?: ConfigurationApplicationSetting[]
  files?: ConfigurationFile[]
  [key: string]: unknown
}
