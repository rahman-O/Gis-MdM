import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { UserForm } from '@/features/users/UserForm'
import {
  arbitraryEmptyOrWhitespace,
  arbitraryRole,
  arbitraryUser,
  arbitraryUserPayload,
} from '@/features/users/test/arbitraries'

const mocks = vi.hoisted(() => ({
  getRoles: vi.fn(),
  createUser: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
}))

const optionMocks = vi.hoisted(() => ({
  getGroups: vi.fn(),
  getConfigurations: vi.fn(),
}))

vi.mock('@/features/users/userService', () => mocks)
vi.mock('@/features/groups/groupService', () => ({
  getGroups: optionMocks.getGroups,
}))
vi.mock('@/features/devices/deviceService', () => ({
  getConfigurations: optionMocks.getConfigurations,
}))

describe('UserForm property tests', () => {
  // Feature: users-management, Property 2: Form submit calls correct endpoint with form values (Req 2.4, 3.3)
  it('Property 2: form submit calls correct endpoint with form values', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryUserPayload(), fc.integer({ min: 1 }), async (payload, id) => {
        cleanup()
        mocks.getRoles.mockReset().mockResolvedValue([{ id: payload.roleId, name: 'Role A' }])
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        mocks.createUser.mockReset().mockResolvedValue({ id, ...payload })
        mocks.updateUser.mockReset().mockResolvedValue({ id, ...payload })

        const onSuccess = vi.fn()
        const onClose = vi.fn()
        const user = userEvent.setup()

        const { unmount } = render(<UserForm mode="create" initialData={null} onSuccess={onSuccess} onClose={onClose} />)
        fireEvent.input(screen.getByLabelText('Login'), { target: { value: payload.login } })
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: payload.name } })
        fireEvent.input(screen.getByLabelText('Email'), { target: { value: payload.email } })
        const plainPw = payload.password ?? 'secret'
        fireEvent.input(screen.getByLabelText('Password'), { target: { value: plainPw } })
        fireEvent.input(screen.getByLabelText('Confirm password'), { target: { value: plainPw } })
        await waitFor(() => expect(screen.queryByText('Loading roles...')).not.toBeInTheDocument())
        await user.click(screen.getByRole('combobox'))
        await waitFor(() => expect(screen.getByRole('option', { name: 'Role A' })).toBeInTheDocument())
        await user.click(screen.getByRole('option', { name: 'Role A' }))
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(mocks.createUser).toHaveBeenCalled())
        expect(mocks.createUser).toHaveBeenLastCalledWith({
          login: payload.login.trim(),
          name: payload.name.trim(),
          email: payload.email.trim(),
          password: plainPw.trim(),
          roleId: payload.roleId,
          allDevicesAvailable: true,
          allConfigAvailable: true,
          groups: [],
          configurations: [],
        })
        unmount()

        render(
          <UserForm
            mode="edit"
            initialData={{
              id,
              login: 'x',
              name: 'y',
              email: 'x@y.com',
              role: { id: payload.roleId, name: 'Role A' },
              allDevicesAvailable: false,
              allConfigAvailable: true,
              groups: [],
              configurations: [],
            }}
            onSuccess={onSuccess}
            onClose={onClose}
          />
        )
        fireEvent.input(screen.getByLabelText('Login'), { target: { value: payload.login } })
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: payload.name } })
        fireEvent.input(screen.getByLabelText('Email'), { target: { value: payload.email } })
        await waitFor(() => expect(screen.queryByText('Loading roles...')).not.toBeInTheDocument())
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() =>
          expect(mocks.updateUser).toHaveBeenCalledWith(
            id,
            expect.objectContaining({
              login: payload.login.trim(),
              name: payload.name.trim(),
              email: payload.email.trim(),
              roleId: payload.roleId,
              allDevicesAvailable: false,
              allConfigAvailable: true,
              groups: [],
              configurations: [],
            })
          )
        )
      }),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: users-management, Property 3: Cancel makes no API call (Req 2.8, 4.6)
  it('Property 3: cancel makes no API call', async () => {
    await fc.assert(
      fc.asyncProperty(fc.boolean(), async (isEdit) => {
        cleanup()
        mocks.getRoles.mockReset().mockResolvedValue([{ id: 1, name: 'Role A' }])
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        mocks.createUser.mockReset()
        mocks.updateUser.mockReset()
        mocks.deleteUser.mockReset()
        const user = userEvent.setup()
        render(
          <UserForm
            mode={isEdit ? 'edit' : 'create'}
            initialData={
              isEdit
                ? {
                    id: 9,
                    login: 'login',
                    name: 'Name',
                    email: 'mail@test.com',
                    role: { id: 1, name: 'Role A' },
                    allDevicesAvailable: false,
                    allConfigAvailable: true,
                    groups: [],
                    configurations: [],
                  }
                : null
            }
            onSuccess={() => {}}
            onClose={() => {}}
          />
        )
        await user.click(screen.getByRole('button', { name: 'Cancel' }))
        expect(mocks.createUser).not.toHaveBeenCalled()
        expect(mocks.updateUser).not.toHaveBeenCalled()
        expect(mocks.deleteUser).not.toHaveBeenCalled()
      }),
      { numRuns: 100 }
    )
  }, 20000)

  // Feature: users-management, Property 4: Edit mode pre-populates fields for any user (Req 3.1)
  it('Property 4: edit mode pre-populates fields for any user', () => {
    fc.assert(
      fc.property(arbitraryUser(), (userData) => {
        cleanup()
        const role = userData.role ?? { id: 1, name: 'Default Role' }
        mocks.getRoles.mockReset().mockResolvedValue([role])
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        render(
          <UserForm
            mode="edit"
            initialData={{ ...userData, role }}
            onSuccess={() => {}}
            onClose={() => {}}
          />
        )
        expect(screen.getByLabelText('Login')).toHaveValue(userData.login)
        expect(screen.getByLabelText('Name')).toHaveValue(userData.name)
        expect(screen.getByLabelText('Email')).toHaveValue(userData.email)
      }),
      { numRuns: 100 }
    )
  })

  // Feature: users-management, Property 7: Role select populates and pre-selects correctly (Req 5.3, 5.4)
  it('Property 7: role select populates and pre-selects correctly', async () => {
    await fc.assert(
      fc.asyncProperty(fc.array(arbitraryRole(), { minLength: 1, maxLength: 5 }), arbitraryUser(), async (roles, u) => {
        cleanup()
        const selectedRole = roles[0]
        mocks.getRoles.mockReset().mockResolvedValue(roles)
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        render(
          <UserForm
            mode="edit"
            initialData={{ ...u, role: selectedRole }}
            onSuccess={() => {}}
            onClose={() => {}}
          />
        )
        await waitFor(() => expect(screen.queryByText('Loading roles...')).not.toBeInTheDocument())
        await waitFor(() => {
          expect(screen.getByRole('combobox')).toHaveTextContent(selectedRole.name)
        })
      }),
      { numRuns: 100 }
    )
  }, 120000)

  // Feature: users-management, Property 8: Required field validation rejects invalid inputs (Req 5.5, 6.1, 6.2, 6.3, 6.4)
  it('Property 8: required field validation rejects invalid inputs', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryEmptyOrWhitespace(), fc.string(), async (badText, invalidEmail) => {
        cleanup()
        mocks.getRoles.mockReset().mockResolvedValue([{ id: 1, name: 'Role A' }])
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        mocks.createUser.mockReset()
        const user = userEvent.setup()
        render(<UserForm mode="create" initialData={null} onSuccess={() => {}} onClose={() => {}} />)
        fireEvent.input(screen.getByLabelText('Login'), { target: { value: badText } })
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: badText } })
        fireEvent.input(screen.getByLabelText('Email'), { target: { value: invalidEmail.replace('@', '') || 'bad' } })
        fireEvent.input(screen.getByLabelText('Password'), { target: { value: badText } })
        fireEvent.input(screen.getByLabelText('Confirm password'), { target: { value: badText } })
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(mocks.createUser).not.toHaveBeenCalled())
      }),
      { numRuns: 100 }
    )
  }, 30000)
})
