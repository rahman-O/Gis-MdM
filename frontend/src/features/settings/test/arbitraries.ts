import * as fc from 'fast-check'
import type { ConfigurationOption, Settings, SettingsPayload } from '@/features/settings/types'

const tokenArb = fc.stringOf(fc.constantFrom(...'abcdefghijklmnopqrstuvwxyz0123456789'), {
  minLength: 1,
  maxLength: 12,
})

const langArb = fc.constantFrom('en', 'ru', 'de', 'fr', 'es', 'pt', 'zh')

export const arbitraryConfigurationOption = (): fc.Arbitrary<ConfigurationOption> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    name: tokenArb,
  })

export const arbitrarySettings = (): fc.Arbitrary<Settings> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    customerName: tokenArb,
    createNewDevices: fc.boolean(),
    newDeviceConfigurationId: fc.oneof(fc.integer({ min: 1 }), fc.constant(null)),
    language: langArb,
    passwordLength: fc.integer({ min: 1, max: 64 }),
    passwordStrength: fc.integer({ min: 0, max: 2 }),
    sendDeviceInfoExpiryDays: fc.integer({ min: 1, max: 3650 }),
    unsecureEnrollment: fc.boolean(),
    deviceFastSearch: fc.boolean(),
    idleLogout: fc.oneof(fc.constant(null), fc.integer({ min: 1, max: 86400 })),
  })

export const arbitrarySettingsPayload = (): fc.Arbitrary<SettingsPayload> =>
  fc.record({
    customerName: tokenArb,
    createNewDevices: fc.boolean(),
    newDeviceConfigurationId: fc.oneof(fc.integer({ min: 1 }), fc.constant(null)),
    language: langArb,
    passwordLength: fc.integer({ min: 1, max: 64 }),
    passwordStrength: fc.integer({ min: 0, max: 2 }),
    sendDeviceInfoExpiryDays: fc.integer({ min: 1, max: 3650 }),
    unsecureEnrollment: fc.boolean(),
    deviceFastSearch: fc.boolean(),
    idleLogout: fc.oneof(fc.constant(null), fc.integer({ min: 1, max: 86400 })),
  })

/** Raw-like settings with nulls coerced by `normalizeSettings` (simulates API). */
export const arbitrarySettingsWithNulls = (): fc.Arbitrary<Record<string, unknown>> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    customerName: tokenArb,
    createNewDevices: fc.constant(null),
    newDeviceConfigurationId: fc.constant(null),
    language: langArb,
    passwordLength: fc.constant(null),
    passwordStrength: fc.constant(null),
    sendDeviceInfoExpiryDays: fc.constant(null),
    unsecureEnrollment: fc.constant(null),
    deviceFastSearch: fc.constant(null),
    idleLogout: fc.constant(null),
  })

export const arbitraryInvalidCustomerName = (): fc.Arbitrary<string> =>
  fc.stringOf(fc.constantFrom(' ', '\t', '\n'), { minLength: 0, maxLength: 8 })

export const arbitraryInvalidPositiveInt = (): fc.Arbitrary<number> =>
  fc.oneof(fc.constant(0), fc.integer({ max: 0 }))
