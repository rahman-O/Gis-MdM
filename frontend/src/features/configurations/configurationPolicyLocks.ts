import type { Configuration } from '@/features/configurations/types'

export function isPolicyLocked(configuration: Configuration | null, fieldKey: string): boolean {
  return Boolean(configuration?.policyLocks?.[fieldKey])
}

export function togglePolicyLock(
  configuration: Configuration,
  fieldKey: string,
  locked: boolean
): Configuration {
  const policyLocks = { ...(configuration.policyLocks ?? {}) }
  if (locked) {
    policyLocks[fieldKey] = true
  } else {
    delete policyLocks[fieldKey]
  }
  return { ...configuration, policyLocks }
}
