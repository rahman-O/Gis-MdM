import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { encodePasswordForUserSave } from '@/features/users/userPasswordEncode'
import { arbitraryUserPayload } from '@/features/users/test/arbitraries'

const mocks = vi.hoisted(() => ({
  get: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}))

vi.mock('@/services/apiClient', () => ({
  default: {
    get: mocks.get,
    put: mocks.put,
    delete: mocks.del,
  },
}))

import * as userService from '@/features/users/userService'

function ok<T>(data: T): { data: HmdmEnvelope<T> } {
  return { data: { status: 'OK', data } }
}

async function didReject(promise: Promise<unknown>): Promise<boolean> {
  try {
    await promise
    return false
  } catch {
    return true
  }
}

describe('userService property tests', () => {
  // Feature: users-management, Property 9: Service routes to correct URL for any operation (Req 8.1, 8.2, 8.3, 8.4, 8.5)
  it('Property 9: service routes to correct URL for any operation', async () => {
    await fc.assert(
      fc.asyncProperty(fc.integer({ min: 1 }), arbitraryUserPayload(), async (id, payload) => {
        mocks.get.mockReset()
        mocks.put.mockReset()
        mocks.del.mockReset()

        mocks.get.mockResolvedValueOnce(ok([]))
        await userService.getUsers()
        expect(mocks.get).toHaveBeenCalledWith('/private/users/all')

        mocks.put.mockResolvedValueOnce(ok({}))
        await userService.createUser(payload)
        expect(mocks.put).toHaveBeenCalledWith(
          '/private/users',
          expect.objectContaining({
            login: payload.login,
            name: payload.name,
            email: payload.email,
            newPassword: encodePasswordForUserSave((payload.password ?? '').trim()),
            userRole: { id: payload.roleId },
            allDevicesAvailable: payload.allDevicesAvailable,
            allConfigAvailable: payload.allConfigAvailable,
            groups: payload.allDevicesAvailable ? null : payload.groups.map((g) => ({ id: g.id })),
            configurations: payload.allConfigAvailable ? null : payload.configurations.map((c) => ({ id: c.id })),
          })
        )

        mocks.put.mockResolvedValueOnce(ok({}))
        await userService.updateUser(id, payload)
        expect(mocks.put).toHaveBeenCalledWith(
          '/private/users',
          expect.objectContaining({
            id,
            login: payload.login,
            name: payload.name,
            email: payload.email,
            userRole: { id: payload.roleId },
            allDevicesAvailable: payload.allDevicesAvailable,
            allConfigAvailable: payload.allConfigAvailable,
            groups: payload.allDevicesAvailable ? null : payload.groups.map((g) => ({ id: g.id })),
            configurations: payload.allConfigAvailable ? null : payload.configurations.map((c) => ({ id: c.id })),
            newPassword: encodePasswordForUserSave((payload.password ?? '').trim()),
          })
        )

        mocks.del.mockResolvedValueOnce({ data: { status: 'OK' } })
        await userService.deleteUser(id)
        expect(mocks.del).toHaveBeenCalledWith(`/private/users/other/${id}`)

        mocks.get.mockResolvedValueOnce(ok([]))
        await userService.getRoles()
        expect(mocks.get).toHaveBeenCalledWith('/private/users/roles')
      }),
      { numRuns: 100 }
    )
  })

  // Feature: users-management, Property 10: Service error propagation (Req 8.7)
  it('Property 10: service error propagation', async () => {
    await fc.assert(
      fc.asyncProperty(fc.constantFrom('list', 'create', 'update', 'delete', 'roles'), fc.integer({ min: 1 }), arbitraryUserPayload(), async (op, id, payload) => {
        const err = new Error('boom')
        mocks.get.mockReset()
        mocks.put.mockReset()
        mocks.del.mockReset()

        if (op === 'list') {
          mocks.get.mockRejectedValueOnce(err)
          expect(await didReject(userService.getUsers())).toBe(true)
          return
        }

        if (op === 'create') {
          mocks.put.mockRejectedValueOnce(err)
          expect(await didReject(userService.createUser(payload))).toBe(true)
          return
        }

        if (op === 'update') {
          mocks.put.mockRejectedValueOnce(err)
          expect(await didReject(userService.updateUser(id, payload))).toBe(true)
          return
        }

        if (op === 'delete') {
          mocks.del.mockRejectedValueOnce(err)
          expect(await didReject(userService.deleteUser(id))).toBe(true)
          return
        }

        mocks.get.mockRejectedValueOnce(err)
        expect(await didReject(userService.getRoles())).toBe(true)
      }),
      { numRuns: 100 }
    )
  })
})
