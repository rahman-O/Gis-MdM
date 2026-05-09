import { describe, it, expect } from 'vitest'
import { existsSync } from 'fs'
import { resolve } from 'path'

// Resolve paths relative to the workspace root (three levels up from frontend/src/test → frontend/src → frontend → workspace root)
const root = resolve(__dirname, '../../..')

describe('Workspace structure', () => {
  it('backend/pom.xml exists (Requirements 1.1, 2.2)', () => {
    expect(existsSync(resolve(root, 'backend/pom.xml'))).toBe(true)
  })

  it('frontend/package.json exists (Requirement 1.2)', () => {
    expect(existsSync(resolve(root, 'frontend/package.json'))).toBe(true)
  })

  it('frontend/components.json exists (Requirement 3.7)', () => {
    expect(existsSync(resolve(root, 'frontend/components.json'))).toBe(true)
  })

  it('frontend/eslint.config.js exists (Requirement 3.3)', () => {
    expect(existsSync(resolve(root, 'frontend/eslint.config.js'))).toBe(true)
  })

  it('frontend/.prettierrc exists (Requirement 3.4)', () => {
    expect(existsSync(resolve(root, 'frontend/.prettierrc'))).toBe(true)
  })

  it('frontend/src/main.tsx exists (Requirement 4.11)', () => {
    expect(existsSync(resolve(root, 'frontend/src/main.tsx'))).toBe(true)
  })

  it('frontend/tailwind.config.ts exists (Requirement 3.8)', () => {
    expect(existsSync(resolve(root, 'frontend/tailwind.config.ts'))).toBe(true)
  })
})

describe('API client configuration', () => {
  it('apiClient baseURL is /rest (proxied to http://localhost:8080/rest via Vite)', async () => {
    const { default: apiClient } = await import('../services/apiClient')
    expect(apiClient.defaults.baseURL).toBe('/rest')
  })
})
