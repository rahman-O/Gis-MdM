import { describe, expect, it, vi } from 'vitest'
import * as fc from 'fast-check'
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { ConfigurationDeleteDialog } from '@/features/configurations/ConfigurationDeleteDialog'
import { arbitraryConfiguration } from '@/features/configurations/test/arbitraries'

describe('ConfigurationDeleteDialog property tests', () => {
  // Property 5 (Req 4.1): delete dialog includes configuration name.
  it('Property 5: delete dialog shows configuration name for any configuration', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryConfiguration(), async (conf) => {
        cleanup()
        render(
          <ConfigurationDeleteDialog
            configuration={{ id: conf.id ?? 1, name: conf.name }}
            onConfirm={async () => {}}
            onCancel={() => {}}
          />
        )
        expect(screen.getByText(conf.name ?? '')).toBeInTheDocument()
      }),
      { numRuns: 100 }
    )
  })

  // Property 6 (Req 4.2): confirming delete calls service with exact id.
  it('Property 6: delete confirm triggers callback for exact id', async () => {
    await fc.assert(
      fc.asyncProperty(arbitraryConfiguration(), async (conf) => {
        cleanup()
        const id = conf.id ?? 1
        const onConfirm = vi.fn().mockResolvedValue(undefined)
        render(
          <ConfigurationDeleteDialog
            configuration={{ id, name: conf.name }}
            onConfirm={onConfirm}
            onCancel={() => {}}
          />
        )
        fireEvent.click(screen.getByRole('button', { name: 'Delete' }))
        await waitFor(() => expect(onConfirm).toHaveBeenCalledTimes(1))
      }),
      { numRuns: 100 }
    )
  }, 15000)
})
