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
  batteryHealth: string | null
  batteryTemperature: number | null
  chargingState: string | null
  model: string | null
  manufacturer: string | null
  androidVersion: string | null
  serial: string | null
  imei: string | null
  phone: string | null
  location: string | { lat: number; lon: number; accuracy?: number; ts?: number } | null
  permissions: DevicePermissions | null
  applications: DeviceApplication[] | null
  files: DeviceFile[] | null
  defaultLauncher: string | null
  mdmMode: boolean | null
  kioskMode: boolean | null
  enrollTime: number | null
  lastSyncTime: number | null
  enrollmentState: string | null
  publicIp: string | null
  ipAddress: string | null
  wifiSsid: string | null
  networkType: string | null
  signalStrength: number | null
  ramUsed: number | null
  ramTotal: number | null
  storageUsed: number | null
  storageFree: number | null
  cpuUsage: number | null
  launcherVersion: string | null
  lastHeartbeat: number | null
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

// ---------------------------------------------------------------------------
// Device Data Optimization — Dashboard Types
// ---------------------------------------------------------------------------

/**
 * A single location data point with full metadata.
 * Used for displaying device positions on the interactive map,
 * history markers, and live tracking updates.
 */
export interface LocationPoint {
  /** Latitude in degrees (-90 to 90) */
  latitude: number
  /** Longitude in degrees (-180 to 180) */
  longitude: number
  /** GPS accuracy in meters */
  accuracy: number
  /** Device speed in m/s */
  speed: number
  /** Altitude in meters (optional, not always reported) */
  altitude?: number
  /** Battery level percentage (0–100) */
  batteryLevel?: number
  /** Network type at time of capture (e.g. "wifi", "cellular", "none") */
  networkType?: string
  /** Tracking mode active when this point was recorded */
  trackingMode?: string
  /** Unix timestamp in milliseconds */
  timestamp: number
}

/**
 * Hourly summary record for archived location data.
 * Used when querying location history older than the retention period.
 */
export interface LocationArchive {
  /** Start of the hour window (Unix timestamp in milliseconds) */
  hourStart: number
  /** Position at the beginning of the hour */
  startPosition: { latitude: number; longitude: number }
  /** Position at the end of the hour */
  endPosition: { latitude: number; longitude: number }
  /** Total distance traveled during the hour in meters */
  distanceTraveled: number
  /** Number of raw location points aggregated into this summary */
  pointCount: number
}

/**
 * An installed application entry for a device.
 * Displayed in the Apps tab of the Device Details Dialog.
 */
export interface DeviceApp {
  /** Android package name (e.g. "com.example.app") */
  packageName: string
  /** Human-readable application name */
  appName: string
  /** Application version string */
  version: string
  /** Current installation status */
  status: 'installed' | 'disabled' | 'uninstalled'
  /** Installation date as Unix timestamp in milliseconds */
  installDate: number
  /** Last update date as Unix timestamp in milliseconds (optional) */
  updateDate?: number
}

/**
 * A single log entry from one of the device log collections.
 * Displayed in the Logs tab of the Device Details Dialog.
 */
export interface DeviceLogEntry {
  /** Unique log entry identifier */
  id: string
  /** Log collection category */
  category: 'system_logs' | 'command_logs' | 'error_logs' | 'tracking_logs'
  /** Log severity level */
  severity: 'DEBUG' | 'INFO' | 'WARNING' | 'ERROR' | 'CRITICAL'
  /** Log message content (may be truncated to 500 characters for display) */
  message: string
  /** Unix timestamp in milliseconds */
  timestamp: number
  /** Optional structured metadata attached to the log entry */
  metadata?: Record<string, unknown>
}

/**
 * WebSocket connection state for live tracking and real-time updates.
 * Used by the connection status indicator in the Location tab.
 */
export type WsConnectionState = 'connected' | 'reconnecting' | 'disconnected'
