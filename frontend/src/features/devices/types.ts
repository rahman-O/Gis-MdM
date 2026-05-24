export interface LookupItem {
  id: number
  name: string | null
}

export interface DevicePermissions {
  permissionStatus: string | null
  details: string | null
}

export interface DeviceApplication {
  pkg: string | null
  version: string | null
  status: string | null
}

export interface DeviceFile {
  path: string | null
  status: string | null
}

export interface AppSetting {
  applicationPkg: string | null
  name: string | null
  type: string | null
  value: string | null
}

export interface ConfigurationView {
  id: number
  name: string | null
  description: string | null
  type: string | null
  qrCodeKey: string | null
  baseUrl: string | null
  permissiveMode: boolean | null
}

export interface ConfigurationOption {
  id: number
  name: string | null
}

export interface DeviceInfoView {
  batteryLevel: number | null
  model: string | null
  androidVersion: string | null
  serial: string | null
  imei: string | null
  phone: string | null
  location: string | null
  permissions: DevicePermissions | null
  applications: DeviceApplication[] | null
  files: DeviceFile[] | null
  defaultLauncher: string | null
  mdmMode: boolean | null
  kioskMode: boolean | null
  enrollTime: number | null
  publicIp: string | null
  launcherVersion: string | null
}

export interface DeviceView {
  id: number
  configurationId: number | null
  number: string
  description: string | null
  lastUpdate: number | null
  imei: string | null
  phone: string | null
  model: string | null
  batteryLevel: number | null
  androidVersion: string | null
  serial: string | null
  statusCode: string | null
  groups: LookupItem[] | null
  custom1: string | null
  custom2: string | null
  custom3: string | null
  oldNumber: string | null
  launcherVersion: string | null
  enrollmentState?: string | null
  treeNodeId?: number | null
  info: DeviceInfoView | null
}

export interface DevicePayload {
  id?: number
  number: string
  description: string | null
  configurationId: number | null
  imei: string | null
  phone: string | null
  custom1: string | null
  custom2: string | null
  custom3: string | null
  oldNumber: string | null
  groups: LookupItem[]
}

export interface DeviceSearchRequest {
  pageNum: number
  pageSize: number
  value?: string | null
  groupId?: number | null
  configurationId?: number | null
  status?: string | null
  androidVersion?: string | null
  sortBy?: string | null
  sortDir?: 'asc' | 'desc' | null
  dateFrom?: number | null
  dateTo?: number | null
  onlineEarlierMillis?: number | null
  onlineLaterMillis?: number | null
  enrollmentDateFrom?: number | null
  enrollmentDateTo?: number | null
  mdmMode?: boolean | null
  kioskMode?: boolean | null
  launcherVersion?: string | null
  installationStatus?: string | null
  imeiChanged?: boolean | null
  fastSearch?: boolean | null
  treeNodeId?: number | null
  includeDescendants?: boolean | null
}

export interface DeviceListResponse {
  devices: {
    items: DeviceView[]
    totalItemsCount: number
  }
  configurations: Record<number, ConfigurationView>
}

export interface DeviceFilters {
  groupId: number | null
  configurationId: number | null
  status: string | null
  androidVersion: string | null
  launcherVersion: string | null
  enrollmentDateFrom: number | null
  enrollmentDateTo: number | null
  onlineEarlierMillis: number | null
  onlineLaterMillis: number | null
  mdmMode: boolean | null
  kioskMode: boolean | null
  installationStatus: string | null
  imeiChanged: boolean | null
  fastSearch: boolean | null
  sortBy: string | null
  sortDir: 'asc' | 'desc' | null
}

export interface BulkDeletePayload {
  ids: number[]
}

export interface GroupBulkPayload {
  ids: number[]
  action: 'set' | 'unset'
  groups: LookupItem[]
}
