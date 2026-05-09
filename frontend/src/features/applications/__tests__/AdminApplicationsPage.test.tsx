import { describe, expect, it, vi, beforeEach } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { render, screen, waitFor } from '@testing-library/react'
import { AdminApplicationsPage } from '@/features/applications/AdminApplicationsPage'

const permMocks = vi.hoisted(() => ({
  isSuperAdmin: vi.fn(() => false),
  hasPermission: vi.fn(() => true),
}))

vi.mock('@/features/auth/permissions', () => ({
  isSuperAdmin: permMocks.isSuperAdmin,
  hasPermission: permMocks.hasPermission,
}))

const serviceMocks = vi.hoisted(() => ({
  getAllAdminApplications: vi.fn(),
  searchAdminApplications: vi.fn(),
  turnApplicationIntoCommon: vi.fn(),
  deleteApplication: vi.fn(),
}))

vi.mock('@/features/applications/services/applicationService', () => serviceMocks)

vi.mock('@/features/applications/components/ApplicationFormDialog', () => ({
  ApplicationFormDialog: () => null,
}))

describe('AdminApplicationsPage', () => {
  beforeEach(() => {
    permMocks.isSuperAdmin.mockReset()
    permMocks.hasPermission.mockReset()
    permMocks.hasPermission.mockImplementation(() => true)
    serviceMocks.getAllAdminApplications.mockReset()
    serviceMocks.searchAdminApplications.mockReset()
  })

  it('blocks non-super-admin users', () => {
    permMocks.isSuperAdmin.mockReturnValue(false)
    render(
      <MemoryRouter>
        <AdminApplicationsPage />
      </MemoryRouter>
    )
    expect(screen.getByText(/super administrators/i)).toBeInTheDocument()
    expect(serviceMocks.getAllAdminApplications).not.toHaveBeenCalled()
  })

  it('loads and lists admin applications for super-admin', async () => {
    permMocks.isSuperAdmin.mockReturnValue(true)
    serviceMocks.getAllAdminApplications.mockResolvedValueOnce([{ id: 1, name: 'Z', pkg: 'a.b', commonApplication: false, customerName: 'Org' }])
    render(
      <MemoryRouter>
        <AdminApplicationsPage />
      </MemoryRouter>
    )
    await waitFor(() => expect(serviceMocks.getAllAdminApplications).toHaveBeenCalled())
    expect(await screen.findByText('Z')).toBeInTheDocument()
    expect(screen.getByText('Org')).toBeInTheDocument()
  })
})
