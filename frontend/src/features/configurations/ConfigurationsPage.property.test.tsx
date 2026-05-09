import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { ConfigurationsPage } from '@/features/configurations/ConfigurationsPage'
import {
  arbitraryConfiguration,
  arbitraryConfigurationWithNulls,
} from '@/features/configurations/test/arbitraries'

const serviceMocks = vi.hoisted(() => ({
  getConfigurations: vi.fn(),
  getConfiguration: vi.fn(),
  createConfiguration: vi.fn(),
  updateConfiguration: vi.fn(),
  deleteConfiguration: vi.fn(),
  typeToConfigurationKind: vi.fn((type?: number | null) => (type === 1 ? 'COMMON' : 'WORK')),
}))

vi.mock('@/features/configurations/configurationService', () => serviceMocks)

function renderPage() {
  return render(
    <MemoryRouter>
      <ConfigurationsPage />
    </MemoryRouter>
  )
}

describe('ConfigurationsPage property tests', () => {
  // Property 12 (Req 9.5): null-safe field access while rendering row values.
  it('Property 12: null-safe field access does not throw when row fields are null', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryConfigurationWithNulls(), async (conf) => {
        serviceMocks.getConfigurations.mockReset().mockResolvedValue([conf])
        const { unmount } = renderPage()
        await waitFor(() => expect(serviceMocks.getConfigurations).toHaveBeenCalled())
        expect(screen.queryAllByText('—').length).toBeGreaterThan(0)
        unmount()
      }),
      { numRuns: 100 }
    )
  })

  // Property 1 (Req 1.3, 6.1): table columns are always rendered.
  it('Property 1: table columns rendered for any configuration list', async () => {
    await fc.assert(
      fc.asyncProperty(fc.array(arbitraryConfiguration(), { minLength: 1, maxLength: 8 }), async (list) => {
        serviceMocks.getConfigurations.mockReset().mockResolvedValue(list)
        const { unmount } = renderPage()
        await waitFor(() => expect(serviceMocks.getConfigurations).toHaveBeenCalled())
        expect(screen.getByText('Name')).toBeInTheDocument()
        expect(screen.getByText('Type')).toBeInTheDocument()
        expect(screen.getByText('Description')).toBeInTheDocument()
        expect(screen.getByText('Device count')).toBeInTheDocument()
        unmount()
      }),
      { numRuns: 100 }
    )
  })

  // Property 9 (Req 6.2, 6.3): displayed deviceCount equals returned value.
  it('Property 9: device count displayed correctly for any configuration', async () => {
    await fc.assert(
      fc.asyncProperty(
        fc.array(
          arbitraryConfiguration().map((c) => ({ ...c, deviceCount: Math.abs(c.deviceCount ?? 0) })),
          { minLength: 1, maxLength: 6 }
        ),
        async (list) => {
          serviceMocks.getConfigurations.mockReset().mockResolvedValue(list)
          const { unmount } = renderPage()
          await waitFor(() => expect(serviceMocks.getConfigurations).toHaveBeenCalled())
          list.forEach((c) => {
            expect(screen.getAllByText(String(c.deviceCount ?? 0)).length).toBeGreaterThan(0)
          })
          unmount()
        }
      ),
      { numRuns: 100 }
    )
  })

  // Property 12 (Req 9.5): null-safe rendering and form opening with null optional fields.
  it('Property 12: null-safe rendering with form open does not throw', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryConfigurationWithNulls(), async (conf) => {
        serviceMocks.getConfigurations.mockReset().mockResolvedValue([conf])
        const { unmount } = renderPage()
        await waitFor(() => expect(serviceMocks.getConfigurations).toHaveBeenCalled())
        expect(screen.queryAllByText('—').length).toBeGreaterThan(0)
        fireEvent.click(screen.getByRole('button', { name: 'New configuration' }))
        expect(screen.getAllByText('New configuration').length).toBeGreaterThan(0)
        unmount()
      }),
      { numRuns: 100 }
    )
  })
})
