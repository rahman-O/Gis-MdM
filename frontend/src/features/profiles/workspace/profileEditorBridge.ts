/** Lets workspace chrome (publish) flush unsaved editor state before server actions. */

export type ProfileEditorBridge = {
  saveIfDirty: () => Promise<boolean>
  getActiveVersionId: () => number | null
}

let bridge: ProfileEditorBridge | null = null

export function registerProfileEditorBridge(next: ProfileEditorBridge): void {
  bridge = next
}

export function unregisterProfileEditorBridge(): void {
  bridge = null
}

/** Saves draft when dirty; returns false if validation/save failed. */
export async function saveProfileEditorIfDirty(): Promise<boolean> {
  if (!bridge) return true
  return bridge.saveIfDirty()
}

export function getProfileEditorActiveVersionId(): number | null {
  return bridge?.getActiveVersionId() ?? null
}
