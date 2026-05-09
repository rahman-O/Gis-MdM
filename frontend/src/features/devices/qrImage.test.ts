/** @vitest-environment node */
import { describe, expect, it } from 'vitest'
import { isLikelyQrImageBlob } from '@/features/devices/qrImage'

describe('isLikelyQrImageBlob', () => {
  it('rejects empty blob', async () => {
    await expect(isLikelyQrImageBlob(new Blob([]))).resolves.toBe(false)
  })

  it('accepts image/* MIME', async () => {
    const b = new Blob([new Uint8Array([1, 2, 3])], { type: 'image/png' })
    await expect(isLikelyQrImageBlob(b)).resolves.toBe(true)
  })

  it('accepts PNG by magic bytes when type is octet-stream', async () => {
    const bytes = Uint8Array.of(0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a)
    const b = new Blob([bytes], { type: 'application/octet-stream' })
    await expect(isLikelyQrImageBlob(b)).resolves.toBe(true)
  })

  it('rejects random bytes without image signature', async () => {
    const b = new Blob([new Uint8Array([1, 2, 3, 4])], { type: 'application/octet-stream' })
    await expect(isLikelyQrImageBlob(b)).resolves.toBe(false)
  })
})
