/**
 * Query parameters for Headwind {@code QRCodeResource}: PNG {@code GET /public/qr/{key}}
 * and JSON {@code GET /public/qr/json/{key}} (same params except {@code size} is PNG-only).
 */

export type QrDeviceIdUseMode = 'request' | 'imei' | 'serial'

export interface EnrollmentQrFields {
  /** Pixel edge length for generated PNG. */
  size: number
  /** Explicit device number / id string embedded in provisioning; empty => server auto-assign semantics. */
  deviceId?: string
  /** When true, adds {@code create=1} (create device on demand on first enrollment). */
  create?: boolean
  /** How device id is interpreted when {@code deviceId} is empty. Ignored when {@code deviceId} is non-empty (legacy Angular). */
  deviceIdUseMode?: QrDeviceIdUseMode
  /** Group ids repeated as {@code group} query keys when {@code create} is true. */
  groupIds?: number[]
}

function appendCommonQrParams(search: URLSearchParams, fields: Omit<EnrollmentQrFields, 'size'>): void {
  const deviceTrimmed = (fields.deviceId ?? '').trim()
  if (deviceTrimmed.length > 0) {
    search.set('deviceId', deviceTrimmed)
  } else {
    const mode = fields.deviceIdUseMode ?? 'imei'
    if (mode === 'imei') search.set('useId', 'imei')
    else if (mode === 'serial') search.set('useId', 'serial')
  }
  if (fields.create) {
    search.set('create', '1')
  }
  for (const gid of fields.groupIds ?? []) {
    if (Number.isFinite(gid) && gid > 0) {
      search.append('group', String(gid))
    }
  }
}

/** Builds search string for provisioning JSON endpoint (no {@code size}). */
export function buildEnrollmentQrJsonQueryString(fields: EnrollmentQrFields): string {
  const search = new URLSearchParams()
  appendCommonQrParams(search, fields)
  return search.toString()
}

/** PNG path includes {@code size}. */
export function buildEnrollmentQrImagePath(qrCodeKey: string, fields: EnrollmentQrFields): string {
  const key = encodeURIComponent(qrCodeKey)
  const search = new URLSearchParams()
  search.set('size', String(Math.max(64, Math.round(fields.size))))
  appendCommonQrParams(search, fields)
  return `/public/qr/${key}?${search.toString()}`
}

export function buildEnrollmentQrJsonPath(qrCodeKey: string, fields: EnrollmentQrFields): string {
  const key = encodeURIComponent(qrCodeKey)
  const q = buildEnrollmentQrJsonQueryString(fields)
  return q.length > 0 ? `/public/qr/json/${key}?${q}` : `/public/qr/json/${key}`
}

/** Default QR size for simple device-row actions (legacy used larger dynamic size in full screen). */
export function defaultCompactQrSize(): number {
  return 360
}

/** Match Angular QR screen: ~80% of smaller viewport edge. */
export function defaultViewportQrSize(): number {
  if (typeof window === 'undefined') return 480
  return Math.max(200, Math.round(Math.min(window.innerWidth, window.innerHeight) * 0.8))
}
