import { describe, expect, it, vi } from 'vitest'
import { MemoryRouter } from 'react-router-dom'
import { render, screen, waitFor } from '@testing-library/react'
import { ApplicationsPage } from '@/features/applications/ApplicationsPage'

const serviceMocks = vi.hoisted(() => ({
  getAllApplications: vi.fn(),
  searchApplications: vi.fn(),
  deleteApplication: vi.fn(),
}))

vi.mock('@/features/applications/services/applicationService', () => ({
  ...serviceMocks,
}))

vi.mock('@/features/applications/services/webUiFilesService', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/features/applications/services/webUiFilesService')>()
  return {
    ...actual,
    getStorageLimit: vi.fn().mockResolvedValue({ sizeLimit: 0, sizeUsed: 0 }),
  }
})

describe('ApplicationsPage', () => {
  it('renders loaded applications', async () => {
    serviceMocks.getAllApplications.mockResolvedValueOnce([
      { id: 1, name: 'Launcher', pkg: 'com.hmdm.launcher', type: 'app' },
    ])
    render(
      <MemoryRouter>
        <ApplicationsPage />
      </MemoryRouter>
    )
    await waitFor(() => expect(serviceMocks.getAllApplications).toHaveBeenCalled())
    expect(screen.getByText('Launcher')).toBeInTheDocument()
  })
})
