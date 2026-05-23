import * as fc from 'fast-check'
import type { Configuration, ConfigurationKind, ConfigurationPayload } from '@/features/configurations/types'

const maybeText = fc.option(fc.string({ minLength: 0, maxLength: 24 }), { nil: null })
const safeWord = fc.stringMatching(/^[A-Za-z0-9_-]{1,30}$/)

export function arbitraryConfigurationPayload(): fc.Arbitrary<ConfigurationPayload> {
  return fc.record({
    name: safeWord,
    description: fc.option(fc.stringMatching(/^[A-Za-z0-9 _-]{0,50}$/), { nil: null }),
    type: fc.constantFrom<ConfigurationKind>('WORK', 'COMMON'),
  })
}

export function arbitraryConfiguration(): fc.Arbitrary<Configuration> {
  return fc.record({
    id: fc.integer({ min: 1, max: 1_000_000 }),
    name: safeWord,
    description: maybeText,
    type: fc.integer({ min: 0, max: 1 }),
    deviceCount: fc.option(fc.integer({ min: 0, max: 2000 }), { nil: null }),
    applications: fc.constant([]),
    files: fc.constant([]),
  })
}

export function arbitraryEmptyOrWhitespace(): fc.Arbitrary<string> {
  return fc.oneof(
    fc.constant(''),
    fc.constant(' '),
    fc.constant('   '),
    fc.array(fc.constantFrom(' ', '\t'), { minLength: 1, maxLength: 10 }).map((chars) => chars.join(''))
  )
}

export function arbitraryConfigurationWithNulls(): fc.Arbitrary<Configuration> {
  return fc.record({
    id: fc.integer({ min: 1, max: 1_000_000 }),
    name: fc.option(fc.string({ minLength: 0, maxLength: 20 }), { nil: null }),
    description: fc.constant(null),
    type: fc.constant(0),
    deviceCount: fc.constant(null),
    applications: fc.constant([]),
    files: fc.constant([]),
  })
}
