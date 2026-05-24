import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { useSearchParams } from 'react-router-dom'

export type ProfileWorkspaceSection =
  | 'overview'
  | 'assignments'
  | 'rollout'
  | 'versions'
  | 'editor'
  | 'activity'

export type ProfileSecondaryPanel = 'publish-impact' | 'assignment-confirm' | null

export interface ProfileWorkspaceState {
  profileId: number | null
  section: ProfileWorkspaceSection
  editorVersionId: number | null
  editorReadOnly: boolean
  secondaryPanel: ProfileSecondaryPanel
  editorDirty: boolean
  open: (profileId: number, section?: ProfileWorkspaceSection) => void
  close: () => void
  setSection: (section: ProfileWorkspaceSection) => void
  setEditorTarget: (versionId: number | null, readOnly?: boolean) => void
  setSecondaryPanel: (panel: ProfileSecondaryPanel) => void
  setEditorDirty: (dirty: boolean) => void
}

const ProfileWorkspaceContext = createContext<ProfileWorkspaceState | null>(null)

const DEFAULT_SECTION: ProfileWorkspaceSection = 'overview'

function parseSection(raw: string | null): ProfileWorkspaceSection {
  const allowed: ProfileWorkspaceSection[] = [
    'overview',
    'assignments',
    'rollout',
    'versions',
    'editor',
    'activity',
  ]
  if (raw && allowed.includes(raw as ProfileWorkspaceSection)) {
    return raw as ProfileWorkspaceSection
  }
  return DEFAULT_SECTION
}

function parsePositiveInt(raw: string | null): number | null {
  if (!raw) return null
  const id = Number(raw)
  return Number.isFinite(id) && id > 0 ? id : null
}

export function ProfileWorkspaceProvider({ children }: { children: ReactNode }) {
  const [searchParams, setSearchParams] = useSearchParams()
  const openParam = searchParams.get('open')
  const sectionParam = searchParams.get('section')
  const versionIdParam = searchParams.get('versionId')
  const readOnlyParam = searchParams.get('readOnly')

  const profileId = useMemo(() => parsePositiveInt(openParam), [openParam])
  const section = parseSection(sectionParam)
  const editorVersionId = useMemo(() => parsePositiveInt(versionIdParam), [versionIdParam])
  const editorReadOnly = readOnlyParam === '1' || readOnlyParam === 'true'

  const [secondaryPanel, setSecondaryPanel] = useState<ProfileSecondaryPanel>(null)
  const [editorDirty, setEditorDirty] = useState(false)

  const syncUrl = useCallback(
    (
      id: number | null,
      nextSection: ProfileWorkspaceSection,
      versionId: number | null,
      readOnly: boolean
    ) => {
      setSearchParams(
        (prev) => {
          const next = new URLSearchParams(prev)
          if (id == null) {
            next.delete('open')
            next.delete('section')
            next.delete('versionId')
            next.delete('readOnly')
          } else {
            next.set('open', String(id))
            next.set('section', nextSection)
            if (nextSection === 'editor' && versionId != null) {
              next.set('versionId', String(versionId))
            } else {
              next.delete('versionId')
            }
            if (nextSection === 'editor' && readOnly) {
              next.set('readOnly', '1')
            } else {
              next.delete('readOnly')
            }
          }
          return next
        },
        { replace: true }
      )
    },
    [setSearchParams]
  )

  const open = useCallback(
    (id: number, nextSection: ProfileWorkspaceSection = DEFAULT_SECTION) => {
      setEditorDirty(false)
      setSecondaryPanel(null)
      syncUrl(id, nextSection, null, false)
    },
    [syncUrl]
  )

  const close = useCallback(() => {
    setEditorDirty(false)
    setSecondaryPanel(null)
    syncUrl(null, DEFAULT_SECTION, null, false)
  }, [syncUrl])

  const setSection = useCallback(
    (nextSection: ProfileWorkspaceSection) => {
      if (profileId == null) return
      if (editorDirty && section === 'editor' && nextSection !== 'editor') {
        if (!window.confirm('You have unsaved editor changes. Leave anyway?')) return
      }
      syncUrl(profileId, nextSection, editorVersionId, editorReadOnly && nextSection === 'editor')
    },
    [profileId, section, editorDirty, editorVersionId, editorReadOnly, syncUrl]
  )

  const setEditorTarget = useCallback(
    (versionId: number | null, readOnly = false) => {
      if (profileId == null) return
      syncUrl(profileId, 'editor', versionId, readOnly)
    },
    [profileId, syncUrl]
  )

  const value = useMemo<ProfileWorkspaceState>(
    () => ({
      profileId,
      section,
      editorVersionId,
      editorReadOnly,
      secondaryPanel,
      editorDirty,
      open,
      close,
      setSection,
      setEditorTarget,
      setSecondaryPanel,
      setEditorDirty,
    }),
    [
      profileId,
      section,
      editorVersionId,
      editorReadOnly,
      secondaryPanel,
      editorDirty,
      open,
      close,
      setSection,
      setEditorTarget,
    ]
  )

  return (
    <ProfileWorkspaceContext.Provider value={value}>{children}</ProfileWorkspaceContext.Provider>
  )
}

export function useProfileWorkspace(): ProfileWorkspaceState {
  const ctx = useContext(ProfileWorkspaceContext)
  if (!ctx) {
    throw new Error('useProfileWorkspace must be used within ProfileWorkspaceProvider')
  }
  return ctx
}

export function useIsMobileViewport(): boolean {
  const [mobile, setMobile] = useState(() =>
    typeof window !== 'undefined' ? window.matchMedia('(max-width: 767px)').matches : false
  )
  useEffect(() => {
    const mq = window.matchMedia('(max-width: 767px)')
    const onChange = () => setMobile(mq.matches)
    mq.addEventListener('change', onChange)
    onChange()
    return () => mq.removeEventListener('change', onChange)
  }, [])
  return mobile
}
