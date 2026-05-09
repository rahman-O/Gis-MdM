export type HmdmResponseStatus = 'OK' | 'ERROR'

export interface HmdmEnvelope<T = unknown> {
  status: HmdmResponseStatus | string
  message?: string | null
  data?: T | null
}

export function unwrapHmdmData<T>(envelope: HmdmEnvelope<T>, fallbackMessage: string): T {
  if (envelope.status !== 'OK') {
    throw new Error(envelope.message?.trim() || fallbackMessage)
  }
  if (envelope.data === undefined || envelope.data === null) {
    throw new Error(envelope.message?.trim() || fallbackMessage)
  }
  return envelope.data
}

/** For `Response.OK()` with no payload. */
export function assertHmdmOk(envelope: HmdmEnvelope<unknown>, fallbackMessage: string): void {
  if (envelope.status !== 'OK') {
    throw new Error(envelope.message?.trim() || fallbackMessage)
  }
}
