import { describe, expect, it, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'

const serviceMocks = vi.hoisted(() => ({
  getDevices: vi.fn(),
  deleteDevice: vi.fn(),
  getDevice: vi.fn(),
  updateDevice: vi.fn(),
  getGroups: vi.fn(),
  getConfigurations: vi.fn(),
  deleteBulk: vi.fn(),
  groupBulk: vi.fn(),
}))

vi.mock('@/features/devices/deviceService', () => serviceMocks)

import { DevicesPage } from '@/features/devices/DevicesPage'

const { getDevices, deleteDevice, getGroups, getConfigurations } = serviceMocks

describe('DevicesPage', () => {
  beforeEach(() => {
    getDevices.mockReset()
    deleteDevice.mockReset()
    getGroups.mockReset()
    getConfigurations.mockReset()
    getDevices.mockResolvedValue({
      devices: { items: [], totalItemsCount: 0 },
      configurations: {},
    })
    getGroups.mockResolvedValue([])
    getConfigurations.mockResolvedValue([])
  })

  it('shows loading skeleton then empty state', async () => {
    let resolveList: (v: unknown) => void = () => {}
    getDevices.mockImplementationOnce(
      () =>
        new Promise((r) => {
          resolveList = r
        })
    )
    render(<DevicesPage />)
    expect(document.querySelectorAll('.animate-pulse').length).toBeGreaterThan(0)
    resolveList({
      devices: { items: [], totalItemsCount: 0 },
      configurations: {},
    })
    await waitFor(() => expect(screen.getByText(/no devices yet/i)).toBeInTheDocument())
  })

  it('shows error banner with retry', async () => {
    getDevices.mockRejectedValueOnce(new Error('offline'))
    render(<DevicesPage />)
    await waitFor(() => expect(screen.getByText('offline')).toBeInTheDocument())
    getDevices.mockResolvedValueOnce({
      devices: { items: [], totalItemsCount: 0 },
      configurations: {},
    })
    await userEvent.click(screen.getByRole('button', { name: /retry/i }))
    await waitFor(() => expect(screen.queryByText('offline')).not.toBeInTheDocument())
  })

  it('shows no search results message when search returns empty', async () => {
    const user = userEvent.setup()
    getDevices.mockResolvedValue({
      devices: { items: [], totalItemsCount: 0 },
      configurations: {},
    })
    render(<DevicesPage />)
    await waitFor(() => expect(getDevices).toHaveBeenCalled())
    await user.type(screen.getByLabelText(/search devices/i), 'xyz')
    await waitFor(
      () => expect(screen.getByText(/no devices found for 'xyz'/i)).toBeInTheDocument(),
      { timeout: 4000 }
    )
  })
})
