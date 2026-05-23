import '@testing-library/jest-dom'

class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

if (!('ResizeObserver' in globalThis)) {
  ;(globalThis as unknown as { ResizeObserver: typeof ResizeObserverMock }).ResizeObserver = ResizeObserverMock
}

/* eslint-disable @typescript-eslint/no-explicit-any -- test env polyfills for Radix/jsdom */
if (typeof Element !== 'undefined') {
  const elProto = Element.prototype as any
  if (typeof elProto.hasPointerCapture !== 'function') {
    elProto.hasPointerCapture = () => false
  }
  if (typeof elProto.setPointerCapture !== 'function') {
    elProto.setPointerCapture = () => {}
  }
  if (typeof elProto.releasePointerCapture !== 'function') {
    elProto.releasePointerCapture = () => {}
  }
  if (typeof elProto.scrollIntoView !== 'function') {
    elProto.scrollIntoView = () => {}
  }
}
