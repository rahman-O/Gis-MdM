type Listener = () => void

const listeners = new Map<number, Set<Listener>>()

function key(profileId: number): number {
  return profileId
}

/** Subscribe to workspace-wide refresh signals (save, publish, assignments, version delete). */
export function subscribeProfileWorkspace(profileId: number, listener: Listener): () => void {
  const k = key(profileId)
  let set = listeners.get(k)
  if (!set) {
    set = new Set()
    listeners.set(k, set)
  }
  set.add(listener)
  return () => {
    set?.delete(listener)
    if (set && set.size === 0) {
      listeners.delete(k)
    }
  }
}

export function notifyProfileWorkspace(profileId: number): void {
  const set = listeners.get(key(profileId))
  if (!set) return
  for (const fn of set) {
    fn()
  }
}
