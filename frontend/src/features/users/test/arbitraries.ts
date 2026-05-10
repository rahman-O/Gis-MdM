import * as fc from 'fast-check'
import type { LookupItem } from '@/features/devices/types'
import type { Role, User, UserPayload } from '@/features/users/types'

const tokenArb = fc.stringOf(fc.constantFrom(...'abcdefghijklmnopqrstuvwxyz0123456789'), {
  minLength: 1,
  maxLength: 16,
})

const lookupArb: fc.Arbitrary<LookupItem> = fc.record({
  id: fc.integer({ min: 1 }),
  name: fc.oneof(fc.string({ minLength: 0, maxLength: 12 }), fc.constant(null)),
})

export const arbitraryRole = (): fc.Arbitrary<Role> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    name: tokenArb,
  })

/** Matches UserForm zod email rule `[^@]+@[^.]+\\..+` (fast-check emails may omit a TLD segment). */
const formValidEmailArb = fc.emailAddress().filter((e) => /[^@]+@[^.]+\..+/.test(e))

export const arbitraryUser = (): fc.Arbitrary<User> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    login: tokenArb,
    name: tokenArb,
    email: formValidEmailArb,
    role: fc.oneof(arbitraryRole(), fc.constant(null)),
    allDevicesAvailable: fc.boolean(),
    allConfigAvailable: fc.boolean(),
    groups: fc.array(lookupArb, { maxLength: 3 }),
    configurations: fc.array(lookupArb, { maxLength: 3 }),
  })

export const arbitraryUserPayload = (): fc.Arbitrary<UserPayload> =>
  fc.record({
    login: tokenArb,
    name: tokenArb,
    email: formValidEmailArb,
    password: tokenArb,
    roleId: fc.integer({ min: 1 }),
    allDevicesAvailable: fc.boolean(),
    allConfigAvailable: fc.boolean(),
    groups: fc.array(lookupArb, { maxLength: 3 }),
    configurations: fc.array(lookupArb, { maxLength: 3 }),
  })

export const arbitraryEmptyOrWhitespace = (): fc.Arbitrary<string> =>
  fc.stringOf(fc.constantFrom(' ', '\t', '\n'), { minLength: 0, maxLength: 8 })

/** Relaxed shape for null-safe UI tests (mock API response). */
export const arbitraryUserWithNulls = (): fc.Arbitrary<
  Omit<User, 'role' | 'allDevicesAvailable' | 'allConfigAvailable' | 'groups' | 'configurations'> & {
    role: null
    allDevicesAvailable: boolean | null
    allConfigAvailable: boolean | null
    groups: LookupItem[] | null
    configurations: LookupItem[] | null
  }
> =>
  fc.record({
    id: fc.integer({ min: 1 }),
    login: tokenArb,
    name: tokenArb,
    email: fc.emailAddress(),
    role: fc.constant(null),
    allDevicesAvailable: fc.constant(null),
    allConfigAvailable: fc.constant(null),
    groups: fc.constant(null),
    configurations: fc.constant(null),
  })
