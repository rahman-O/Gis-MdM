import { describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { DeleteDialog } from '@/features/devices/DeleteDialog'
import type { DeviceView } from '@/features/devices/types'

const device: DeviceView = {
  id: 7,
  configurationId: 1,
  number: 'T-100',
  description: null,
  lastUpdate: null,
  imei: null,
  phone: null,
  model: null,
  batteryLevel: null,
  androidVersion: null,
  serial: null,
  statusCode: 'grey',
  groups: [],
  custom1: null,
  custom2: null,
  custom3: null,
  oldNumber: null,
  launcherVersion: null,
  info: null,
}

describe('DeleteDialog', () => {
  it('confirm calls onConfirm and closes on success', async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn().mockResolvedValue(undefined)
    const onCancel = vi.fn()
    render(<DeleteDialog device={device} onConfirm={onConfirm} onCancel={onCancel} />)
    await user.click(screen.getByRole('button', { name: /delete/i }))
    await waitFor(() => expect(onConfirm).toHaveBeenCalled())
    await waitFor(() => expect(onCancel).toHaveBeenCalled())
  })

  it('shows error and stays open on failure', async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn().mockRejectedValue(new Error('boom'))
    const onCancel = vi.fn()
    render(<DeleteDialog device={device} onConfirm={onConfirm} onCancel={onCancel} />)
    await user.click(screen.getByRole('button', { name: /delete/i }))
    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('boom'))
    expect(onCancel).not.toHaveBeenCalled()
  })

  it('cancel does not call onConfirm', async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn()
    const onCancel = vi.fn()
    render(<DeleteDialog device={device} onConfirm={onConfirm} onCancel={onCancel} />)
    await user.click(screen.getByRole('button', { name: /cancel/i }))
    expect(onConfirm).not.toHaveBeenCalled()
    expect(onCancel).toHaveBeenCalled()
  })
})
