import { describe, expect, it } from 'vitest'
import {
  buildEnrollmentQrImagePath,
  buildEnrollmentQrJsonPath,
  buildEnrollmentQrJsonQueryString,
} from '@/features/devices/enrollmentQrQuery'

describe('enrollmentQrQuery', () => {
  const key = 'abc/ key'

  it('JSON query omits size and encodes key', () => {
    const q = buildEnrollmentQrJsonQueryString({
      size: 400,
      deviceId: ' dev1 ',
      create: true,
      groupIds: [1, 2],
    })
    expect(q).toContain('deviceId=dev1')
    expect(q).toContain('create=1')
    expect(q.match(/group=/g)?.length).toBe(2)
    expect(q.includes('size')).toBe(false)
  })

  it('when deviceId empty, maps useId for imei and serial', () => {
    expect(
      buildEnrollmentQrJsonQueryString({
        size: 300,
        deviceId: '',
        deviceIdUseMode: 'imei',
      })
    ).toContain('useId=imei')
    expect(
      buildEnrollmentQrJsonQueryString({
        size: 300,
        deviceId: undefined,
        deviceIdUseMode: 'serial',
      })
    ).toContain('useId=serial')
  })

  it('when deviceId non-empty, does not emit useId (request vs imei suppressed)', () => {
    const q = buildEnrollmentQrJsonQueryString({
      size: 300,
      deviceId: 'n1',
      deviceIdUseMode: 'imei',
    })
    expect(q).toContain('deviceId=n1')
    expect(q.includes('useId')).toBe(false)
  })

  it('PNG path includes size before other keys', () => {
    const path = buildEnrollmentQrImagePath(key, {
      size: 320,
      create: false,
      deviceIdUseMode: 'request',
    })
    expect(path.startsWith(`/public/qr/${encodeURIComponent(key)}?`)).toBe(true)
    expect(path).toMatch(/[?&]size=320(&|$)/)
  })

  it('JSON path strips size', () => {
    const path = buildEnrollmentQrJsonPath(key, { size: 999, deviceId: 'x' })
    expect(path.includes('size')).toBe(false)
  })
})
