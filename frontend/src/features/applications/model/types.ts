export type ApplicationType = 'app' | 'web' | 'intent'

export interface Application {
  id?: number | null
  name?: string | null
  pkg?: string | null
  version?: string | null
  versionCode?: number | null
  url?: string | null
  urlArmeabi?: string | null
  urlArm64?: string | null
  split?: boolean | null
  arch?: string | null
  type?: ApplicationType | string | null
  showIcon?: boolean | null
  useKiosk?: boolean | null
  system?: boolean | null
  runAfterInstall?: boolean | null
  runAtBoot?: boolean | null
  skipVersion?: boolean | null
  iconText?: string | null
  iconId?: number | null
  intent?: string | null
  latestVersion?: number | null
  latestVersionText?: string | null
  usedVersionId?: number | null
  customerId?: number | null
  customerName?: string | null
  common?: boolean | null
  commonApplication?: boolean | null
  deletionProhibited?: boolean | null
  outdated?: boolean | null
  action?: number | null
  filePath?: string | null
  autoUpdate?: boolean | null
  [key: string]: unknown
}

export interface ApplicationVersion {
  id?: number | null
  applicationId?: number | null
  version?: string | null
  versionCode?: number | null
  url?: string | null
  urlArmeabi?: string | null
  urlArm64?: string | null
  split?: boolean | null
  arch?: string | null
  action?: number | null
  showIcon?: boolean | null
  useKiosk?: boolean | null
  screenOrder?: number | null
  keyCode?: number | null
  bottom?: boolean | null
  longTap?: boolean | null
  intent?: string | null
  autoUpdate?: boolean | null
  filePath?: string | null
  [key: string]: unknown
}

export interface ApplicationConfigurationLink {
  configurationId: number
  name?: string | null
  action: number
  selected?: boolean
  notify?: boolean
  [key: string]: unknown
}

export interface ApplicationVersionConfigurationLink {
  configurationId: number
  name?: string | null
  action: number
  selected?: boolean
  notify?: boolean
  [key: string]: unknown
}

export interface LinkConfigurationsToAppRequest {
  applicationId: number
  configurations: ApplicationConfigurationLink[]
}

export interface LinkConfigurationsToAppVersionRequest {
  applicationVersionId: number
  configurations: ApplicationVersionConfigurationLink[]
}

export interface ApkFileDetails {
  pkg?: string | null
  version?: string | null
  versionCode?: number | null
  arch?: string | null
}

export interface FileUploadResult {
  serverPath?: string | null
  fileDetails?: ApkFileDetails | null
  application?: Application | null
  complete?: boolean
  exists?: boolean
  name?: string | null
  [key: string]: unknown
}

export interface ApplicationFormValues {
  id?: number | null
  type: ApplicationType
  name: string
  pkg: string
  version: string
  versionCode: number | null
  url: string
  urlArmeabi: string
  urlArm64: string
  split: boolean
  arch: string
  showIcon: boolean
  useKiosk: boolean
  system: boolean
  runAfterInstall: boolean
  runAtBoot: boolean
  skipVersion: boolean
  iconText: string
  iconId: number | null
  intent: string
  filePath: string
  autoUpdate: boolean
}
