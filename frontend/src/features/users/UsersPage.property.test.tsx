import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { UsersPage } from '@/features/users/UsersPage'
import { arbitraryUser, arbitraryUserWithNulls } from '@/features/users/test/arbitraries'

const mocks = vi.hoisted(() => ({
  getUsers: vi.fn(),
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

function renderUsersPage() {
  render(
    <MemoryRouter>
      <UsersPage />
    </MemoryRouter>
  )
}

describe('UsersPage property tests', () => {
  // Feature: users-management, Property 1: Table columns rendered for any user list (Req 1.3)
  it('Property 1: table columns rendered for any user list', async () => {
    await fc.assert(
      fc.asyncProperty(fc.array(arbitraryUser(), { minLength: 1, maxLength: 10 }), async (users) => {
        cleanup()
        mocks.getUsers.mockReset().mockResolvedValue(users)
        renderUsersPage()
        await waitFor(() => expect(mocks.getUsers).toHaveBeenCalled())
        expect(screen.getByRole('columnheader', { name: 'Login' })).toBeInTheDocument()
        expect(screen.getByRole('columnheader', { name: 'Name' })).toBeInTheDocument()
        expect(screen.getByRole('columnheader', { name: 'Email' })).toBeInTheDocument()
        expect(screen.getByRole('columnheader', { name: 'Role' })).toBeInTheDocument()
        expect(screen.getByRole('columnheader', { name: 'Status' })).toBeInTheDocument()
      }),
      { numRuns: 100 }
    )
  }, 30000)

  // Feature: users-management, Property 5: Delete dialog shows user login and name for any user (Req 4.1)
  it('Property 5: delete dialog shows user login and name for any user', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryUser(), async (userData) => {
        cleanup()
        mocks.getUsers.mockReset().mockResolvedValue([userData])
        const user = userEvent.setup()
        renderUsersPage()
        await waitFor(() => expect(screen.getByText(userData.email)).toBeInTheDocument())
        await user.click(screen.getAllByRole('button', { name: /Actions for/ })[0])
        await user.click(screen.getByText('Delete'))
        const dialog = await screen.findByRole('alertdialog')
        expect(dialog).toHaveTextContent(userData.login)
        expect(dialog).toHaveTextContent(userData.name)
      }),
      { numRuns: 100 }
    )
  }, 45000)

  // Feature: users-management, Property 6: Delete confirm calls DELETE with correct id (Req 4.2)
  it('Property 6: delete confirm calls DELETE with correct id', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryUser(), async (userData) => {
        cleanup()
        mocks.getUsers.mockReset().mockResolvedValue([userData])
        mocks.deleteUser.mockReset().mockResolvedValue(undefined)
        const user = userEvent.setup()
        renderUsersPage()
        await waitFor(() => expect(screen.getByText(userData.email)).toBeInTheDocument())
        await user.click(screen.getAllByRole('button', { name: /Actions for/ })[0])
        await user.click(screen.getByText('Delete'))
        await user.click(screen.getByRole('button', { name: 'Delete' }))
        await waitFor(() => expect(mocks.deleteUser).toHaveBeenCalledWith(userData.id))
      }),
      { numRuns: 100 }
    )
  }, 45000)

  // Feature: users-management, Property 11: Null-safe rendering for any user with null optional fields (Req 9.4, 9.5)
  it('Property 11: null-safe rendering for any user with null optional fields', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryUserWithNulls(), async (userData) => {
        cleanup()
        mocks.getUsers.mockReset().mockResolvedValue([userData])
        mocks.getRoles.mockReset().mockResolvedValue([{ id: 1, name: 'Role A' }])
        optionMocks.getGroups.mockReset().mockResolvedValue([])
        optionMocks.getConfigurations.mockReset().mockResolvedValue([])
        const user = userEvent.setup()
        renderUsersPage()
        await waitFor(() => expect(screen.getByText(userData.email)).toBeInTheDocument())
        await user.click(screen.getAllByRole('button', { name: /Actions for/ })[0])
        await user.click(screen.getByText('Edit'))
        await waitFor(() => expect(screen.getByRole('dialog')).toBeInTheDocument())
      }),
      { numRuns: 100 }
    )
  }, 45000)
})
