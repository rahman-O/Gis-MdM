/** Mirrors chart rows used inside `SummaryResponse.java`. */
export interface ChartItemLike {
  stringAttr?: string | null
  intAttr?: number | null
  number?: number | null
}

export interface DeviceSummaryPayload {
  statusSummary?: ChartItemLike[]
  installSummary?: ChartItemLike[]
  devicesTotal?: number | null
  devicesEnrolled?: number | null
  devicesEnrolledLastMonth?: number | null
  devicesEnrolledMonthly?: ChartItemLike[]
  topConfigs?: string[] | null
  statusOfflineByConfig?: number[] | null
  statusIdleByConfig?: number[] | null
  statusOnlineByConfig?: number[] | null
  appFailureByConfig?: number[] | null
  appMismatchByConfig?: number[] | null
  appSuccessByConfig?: number[] | null
}

export interface DashboardCounts {
  configurationCount: number
  applicationCount: number
}

export interface ParsedStatusCounts {
  offline: number
  idle: number
  online: number
}

export interface ParsedInstallCounts {
  failure: number
  mismatch: number
  success: number
}
