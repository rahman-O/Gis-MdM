import { describe, expect, it, vi, afterEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useMediaQuery } from '@/shared/hooks/useMediaQuery'

describe('useMediaQuery', () => {
  let listeners: Array<(e: MediaQueryListEvent) => void> = []
  let currentMatches = false

  const mockMatchMedia = vi.fn((query: string) => ({
    matches: currentMatches,
    media: query,
    addEventListener: (_event: string, handler: (e: MediaQueryListEvent) => void) => {
      listeners.push(handler)
    },
    removeEventListener: (_event: string, handler: (e: MediaQueryListEvent) => void) => {
      listeners = listeners.filter((l) => l !== handler)
    },
  }))

  afterEach(() => {
    listeners = []
    currentMatches = false
    vi.restoreAllMocks()
  })

  it('returns true when media query matches initially', () => {
    currentMatches = true
    vi.stubGlobal('matchMedia', mockMatchMedia)

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'))
    expect(result.current).toBe(true)
  })

  it('returns false when media query does not match initially', () => {
    currentMatches = false
    vi.stubGlobal('matchMedia', mockMatchMedia)

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'))
    expect(result.current).toBe(false)
  })

  it('updates when media query match changes', () => {
    currentMatches = false
    vi.stubGlobal('matchMedia', mockMatchMedia)

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'))
    expect(result.current).toBe(false)

    act(() => {
      listeners.forEach((l) => l({ matches: true } as MediaQueryListEvent))
    })
    expect(result.current).toBe(true)
  })

  it('cleans up listener on unmount', () => {
    currentMatches = false
    vi.stubGlobal('matchMedia', mockMatchMedia)

    const { unmount } = renderHook(() => useMediaQuery('(min-width: 768px)'))
    expect(listeners.length).toBe(1)

    unmount()
    expect(listeners.length).toBe(0)
  })
})
