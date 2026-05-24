export type HmdmResponseStatus = 'OK' | 'ERROR'

export interface HmdmEnvelope<T = unknown> {
  status: HmdmResponseStatus | string
  message?: string | null
  data?: T | null
}

/** Backend often returns localization keys from `Response`; map common ones for the React UI. */
const KNOWN_SERVER_MESSAGE_KEYS: Record<string, string> = {
  'error.internal.server': 'Internal server error. If this persists, check server logs for the underlying exception.',
  'error.profile.version.delete.activePublished':
    'Cannot delete the version currently published for this profile.',
  'error.profile.version.delete.assigned':
    'This version is still assigned to a folder. Remove the assignment first.',
  'error.profile.version.delete.devicesTarget':
    'Devices are still targeting this version.',
}

function resolveEnvelopeMessage(raw: string | null | undefined, fallback: string): string {
  const key = raw?.trim()
  if (!key) return fallback
  return KNOWN_SERVER_MESSAGE_KEYS[key] ?? key
}

export function unwrapHmdmData<T>(envelope: HmdmEnvelope<T>, fallbackMessage: string): T {
  if (envelope.status !== 'OK') {
    throw new Error(resolveEnvelopeMessage(envelope.message, fallbackMessage))
  }
  if (envelope.data === undefined || envelope.data === null) {
    throw new Error(resolveEnvelopeMessage(envelope.message, fallbackMessage))
  }
  return envelope.data
}

/** Java often returns an empty list; Go may serialize a nil slice as JSON `null`. */
export function unwrapHmdmList<T>(
  envelope: HmdmEnvelope<T[] | null | undefined>,
  fallbackMessage: string
): T[] {
  if (envelope.status !== 'OK') {
    throw new Error(resolveEnvelopeMessage(envelope.message, fallbackMessage))
  }
  const data = envelope.data
  if (data == null) return []
  return Array.isArray(data) ? data : []
}

/** For `Response.OK()` with no payload. */
export function assertHmdmOk(envelope: HmdmEnvelope<unknown>, fallbackMessage: string): void {
  if (envelope.status !== 'OK') {
    throw new Error(resolveEnvelopeMessage(envelope.message, fallbackMessage))
  }
}
