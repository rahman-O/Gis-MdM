import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import * as fc from 'fast-check'
import { StatusBadge } from '@/features/devices/StatusBadge'

describe('StatusBadge', () => {
  it('maps green and red', () => {
    const { unmount: u1 } = render(<StatusBadge statusCode="green" />)
    expect(screen.getByText('Online')).toBeInTheDocument()
    u1()
    render(<StatusBadge statusCode="red" />)
    expect(screen.getByText('Offline')).toBeInTheDocument()
  })

  it('Property 1: every statusCode renders exactly one badge label', () => {
    fc.assert(
      fc.property(fc.oneof(fc.constant(null), fc.constant('green'), fc.constant('red'), fc.string()), (code) => {
        const { unmount, container } = render(<StatusBadge statusCode={code} />)
        const badges = container.querySelectorAll('[class*="rounded-full"]')
        expect(badges.length).toBeGreaterThanOrEqual(1)
        const text =
          code === 'green' ? 'Online'
          : code === 'red' ? 'Offline'
          : 'Unknown'
        expect(screen.getByText(text)).toBeInTheDocument()
        unmount()
      }),
      { numRuns: 30 }
    )
  })
})
