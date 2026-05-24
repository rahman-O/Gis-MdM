import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  listProfileVersions,
  forkDraftFromVersion,
  type ProfileVersionListItem,
} from '@/features/profiles/profileRolloutService'
import { Button } from '@/shared/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'

interface Props {
  profileId: number
  currentVersionId: number | null
  isDraft: boolean
  dirty: boolean
  onVersionChange: (versionId: number, readOnly: boolean) => void
  /** When true, stay on workspace URL (no legacy /edit navigation). */
  workspaceMode?: boolean
  onForkComplete?: () => void
}

export function ProfileVersionSelect({
  profileId,
  currentVersionId,
  isDraft,
  dirty,
  onVersionChange,
  workspaceMode = false,
  onForkComplete,
}: Props) {
  const navigate = useNavigate()
  const [versions, setVersions] = useState<ProfileVersionListItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void listProfileVersions(profileId)
      .then(setVersions)
      .catch(() => setVersions([]))
      .finally(() => setLoading(false))
  }, [profileId])

  const handleSelect = (value: string) => {
    const versionId = Number(value)
    const item = versions.find((v) => v.versionId === versionId)
    if (!item) return
    if (dirty && isDraft) {
      const ok = window.confirm('You have unsaved changes. Discard and switch version?')
      if (!ok) return
    }
    onVersionChange(versionId, item.status === 'published')
    if (!workspaceMode) {
      navigate(`/profiles/${profileId}/versions/${versionId}/edit`, { replace: true })
    }
  }

  const handleFork = async () => {
    if (!currentVersionId) return
    try {
      await forkDraftFromVersion(profileId, currentVersionId)
      if (workspaceMode) {
        onForkComplete?.()
        return
      }
      navigate(`/profiles/${profileId}/edit`, { replace: true })
      window.location.reload()
    } catch (e: unknown) {
      window.alert(e instanceof Error ? e.message : 'Fork failed')
    }
  }

  if (loading) return null

  return (
    <div className="flex flex-wrap items-center gap-2">
      <Select value={currentVersionId ? String(currentVersionId) : undefined} onValueChange={handleSelect}>
        <SelectTrigger className="w-[220px]">
          <SelectValue placeholder="Version" />
        </SelectTrigger>
        <SelectContent>
          {versions.map((v) => (
            <SelectItem key={v.versionId} value={String(v.versionId)}>
              v{v.versionNumber} · {v.status}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {!isDraft && currentVersionId ? (
        <Button type="button" variant="outline" size="sm" onClick={() => void handleFork()}>
          Fork draft from this version
        </Button>
      ) : null}
    </div>
  )
}
