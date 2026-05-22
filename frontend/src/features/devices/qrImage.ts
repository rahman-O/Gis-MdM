import apiClient from '@/services/apiClient'
import { buildEnrollmentQrImagePath, defaultCompactQrSize } from '@/features/devices/enrollmentQrQuery'

/** Jersey may use `APPLICATION_OCTET_STREAM`; some proxies strip MIME — detect common raster headers. */
export async function isLikelyQrImageBlob(blob: Blob): Promise<boolean> {
  if (!blob || blob.size === 0) return false
  const t = blob.type.toLowerCase()
  if (t.startsWith('image/')) return true

  let full: ArrayBuffer
  try {
    full =
      typeof blob.arrayBuffer === 'function'
        ? await blob.arrayBuffer()
        : await new Response(blob).arrayBuffer()
  } catch {
    return false
  }
  const head = new Uint8Array(full).subarray(0, Math.min(12, full.byteLength))
  if (head.length < 4) return false
  // PNG
  if (head[0] === 0x89 && head[1] === 0x50 && head[2] === 0x4e && head[3] === 0x47) return true
  // JPEG
  if (head[0] === 0xff && head[1] === 0xd8 && head[2] === 0xff) return true
  // GIF
  if (head[0] === 0x47 && head[1] === 0x49 && head[2] === 0x46) return true
  return false
}

export function buildPublicQrPaths(qrCodeKey: string, deviceNumber?: string): { primaryPath: string; fallbackPath: string } {
  const size = defaultCompactQrSize()
  const primaryPath = buildEnrollmentQrImagePath(qrCodeKey, {
    size,
    deviceId: deviceNumber?.trim() || undefined,
    deviceIdUseMode: deviceNumber?.trim() ? undefined : 'request',
  })
  const fallbackPath = buildEnrollmentQrImagePath(qrCodeKey, {
    size,
    deviceIdUseMode: 'request',
  })
  return { primaryPath, fallbackPath }
}

async function qrErrorFromResponse(response: { status: number; data: Blob }): Promise<string | null> {
  if (response.status < 400) return null
  try {
    const text = (await response.data.text()).trim()
    if (text) return text
  } catch {
    /* ignore */
  }
  if (response.status === 404) return 'Configuration not found for this QR key.'
  if (response.status === 400) {
    return 'Main application has no APK download URL. Upload an APK for the Main App version in Applications, then save the configuration.'
  }
  return `Server returned HTTP ${response.status} for the QR image.`
}

/** Fetch a single enrollment QR PNG path; optionally abort in-flight loads. */
export async function loadQrImageObjectUrl(
  imagePath: string,
  signal?: AbortSignal
): Promise<{ url: string | null; error: string | null }> {
  try {
    const response = await apiClient.get<Blob>(imagePath, {
      responseType: 'blob',
      signal,
      validateStatus: (status) => status < 500 || status === 500,
    })
    const blob = response.data
    if (await isLikelyQrImageBlob(blob)) {
      return { url: URL.createObjectURL(blob), error: null }
    }
    const serverMsg = await qrErrorFromResponse(response)
    return { url: null, error: serverMsg }
  } catch (e: unknown) {
    if (signal?.aborted) return { url: null, error: null }
    const msg = e instanceof Error ? e.message : null
    return { url: null, error: msg }
  }
}

/**
 * Loads QR PNG via Axios (same `baseURL`/proxy/`withCredentials` as the rest of the app).
 * Returns an object URL to revoke after use, or null.
 */
export async function loadDeviceQrObjectUrl(primaryPath: string, fallbackPath: string): Promise<string | null> {
  const primary = await loadQrImageObjectUrl(primaryPath)
  if (primary.url) return primary.url
  const fallback = await loadQrImageObjectUrl(fallbackPath)
  return fallback.url
}
