import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ConfigurationForm } from '@/features/configurations/ConfigurationForm'
import {
  arbitraryConfiguration,
  arbitraryConfigurationPayload,
  arbitraryEmptyOrWhitespace,
} from '@/features/configurations/test/arbitraries'

const mocks = vi.hoisted(() => ({
  createConfiguration: vi.fn(),
  updateConfiguration: vi.fn(),
  typeToConfigurationKind: vi.fn((type?: number | null) => (type === 1 ? 'COMMON' : 'WORK')),
}))

vi.mock('@/features/configurations/configurationService', () => mocks)

describe('ConfigurationForm property tests', () => {
  // Property 2 (Req 2.3, 3.2): create/edit submit routes to correct service function.
  it('Property 2: form submit calls correct endpoint with form values', async () => {
    await fc.assert(
      fc.asyncProperty(
        arbitraryConfigurationPayload().map((p) => ({ ...p, type: 'WORK' as const })),
        fc.integer({ min: 1 }),
        async (payload, id) => {
        cleanup()
        const expectedPayload = {
          name: payload.name.trim(),
          type: payload.type,
          description: payload.description && payload.description.trim() !== '' ? payload.description.trim() : null,
        }
        mocks.createConfiguration.mockReset().mockResolvedValue({ id })
        mocks.updateConfiguration.mockReset().mockResolvedValue({ id })
        const onSuccess = vi.fn()
        const onOpenChange = vi.fn()
        const user = userEvent.setup()

        const { unmount } = render(
          <ConfigurationForm
            open
            mode="create"
            initialData={null}
            onSuccess={onSuccess}
            onOpenChange={onOpenChange}
          />
        )
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: payload.name } })
        fireEvent.input(screen.getByLabelText('Description'), { target: { value: payload.description ?? '' } })
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(mocks.createConfiguration).toHaveBeenCalled())
        expect(mocks.createConfiguration).toHaveBeenLastCalledWith(expectedPayload)
        unmount()

        const { unmount: unmountEdit } = render(
          <ConfigurationForm
            open
            mode="edit"
            initialData={{ id, name: 'init', type: 0, description: '', applications: [] }}
            onSuccess={onSuccess}
            onOpenChange={onOpenChange}
          />
        )
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: payload.name } })
        fireEvent.input(screen.getByLabelText('Description'), { target: { value: payload.description ?? '' } })
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(mocks.updateConfiguration).toHaveBeenCalledWith(id, expectedPayload))
        unmountEdit()
        }
      ),
      { numRuns: 100 }
    )
  }, 30000)

  // Property 3 (Req 2.7, 4.6): cancel does not call create/update service.
  it('Property 3: cancel makes no API call', async () => {
    await fc.assert(
      fc.asyncProperty(fc.boolean(), async (isEdit) => {
        cleanup()
        mocks.createConfiguration.mockReset()
        mocks.updateConfiguration.mockReset()
        const user = userEvent.setup()
        const onOpenChange = vi.fn()
        render(
          <ConfigurationForm
            open
            mode={isEdit ? 'edit' : 'create'}
            initialData={isEdit ? { id: 5, name: 'x', type: 1, description: '', applications: [] } : null}
            onSuccess={() => {}}
            onOpenChange={onOpenChange}
          />
        )
        await user.click(screen.getByRole('button', { name: 'Cancel' }))
        expect(mocks.createConfiguration).not.toHaveBeenCalled()
        expect(mocks.updateConfiguration).not.toHaveBeenCalled()
      }),
      { numRuns: 100 }
    )
  }, 15000)

  // Property 4 (Req 3.1): edit mode pre-populates name/description/type.
  it('Property 4: edit mode pre-populates fields for any configuration', () => {
    fc.assert(
      fc.property(arbitraryConfiguration(), (conf) => {
        cleanup()
        const { unmount } = render(
          <ConfigurationForm
            open
            mode="edit"
            initialData={conf}
            onSuccess={() => {}}
            onOpenChange={() => {}}
          />
        )
        expect(screen.getByLabelText('Name')).toHaveValue((conf.name ?? '').trim())
        expect(screen.getByLabelText('Description')).toHaveValue(conf.description ?? '')
        expect(screen.getByRole('combobox')).toHaveTextContent(conf.type === 1 ? 'Common (typical)' : 'Work (device)')
        unmount()
      }),
      { numRuns: 100 }
    )
  })

  // Property 7 (Req 5.1, 5.2): invalid empty/whitespace name rejected.
  it('Property 7: required field validation rejects empty or whitespace inputs', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryEmptyOrWhitespace(), async (badName) => {
        cleanup()
        mocks.createConfiguration.mockReset()
        const user = userEvent.setup()
        render(
          <ConfigurationForm
            open
            mode="create"
            initialData={null}
            onSuccess={() => {}}
            onOpenChange={() => {}}
          />
        )
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: badName } })
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(screen.getByText('Name is required')).toBeInTheDocument())
        expect(mocks.createConfiguration).not.toHaveBeenCalled()
      }),
      { numRuns: 100 }
    )
  }, 20000)

  // Property 8 (Req 5.3): optional description can be empty and still submit.
  it('Property 8: optional description allows submission', async () => {
    await fc.assert(
      fc.asyncProperty(fc.string({ minLength: 1, maxLength: 20 }), async (name) => {
        cleanup()
        mocks.createConfiguration.mockReset().mockResolvedValue({ id: 1 })
        const user = userEvent.setup()
        render(
          <ConfigurationForm
            open
            mode="create"
            initialData={null}
            onSuccess={() => {}}
            onOpenChange={() => {}}
          />
        )
        fireEvent.input(screen.getByLabelText('Name'), { target: { value: name } })
        fireEvent.input(screen.getByLabelText('Description'), { target: { value: '' } })
        await user.click(screen.getByRole('button', { name: 'Save' }))
        await waitFor(() => expect(mocks.createConfiguration).toHaveBeenCalled())
      }),
      { numRuns: 100 }
    )
  }, 20000)
})
