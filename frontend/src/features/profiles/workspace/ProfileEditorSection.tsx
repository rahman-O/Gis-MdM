import { useState } from 'react'
import { ProfileEditorCore } from '@/features/profiles/ProfileEditorCore'
import { useProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceState'

interface Props {
  profileId: number
}

export function ProfileEditorSection({ profileId }: Props) {
  const { editorVersionId, editorReadOnly, setEditorDirty, setEditorTarget } = useProfileWorkspace()
  const [lastSaved, setLastSaved] = useState<Date | null>(null)

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="min-h-0 flex-1 overflow-y-auto p-4 pb-24">
        {editorReadOnly ? (
          <p className="mb-3 rounded-md border bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
            Viewing published policy (read-only). Use the version switcher or fork a draft to edit.
          </p>
        ) : (
          <p className="mb-3 text-sm text-muted-foreground">
            Edit the draft, save, then publish from the workspace header.
          </p>
        )}
        <ProfileEditorCore
          embedded
          profileIdProp={profileId}
          versionIdOverride={editorVersionId}
          readOnlyOverride={editorReadOnly}
          workspaceVersionSelect
          workspaceChrome
          hideWorkspacePublish
          onWorkspaceVersionChange={(vid, ro) => setEditorTarget(vid, ro)}
          onDirtyChange={setEditorDirty}
          onLastSaved={setLastSaved}
        />
      </div>
      <div className="sticky bottom-0 z-10 border-t bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/80">
        <p className="text-xs text-muted-foreground">
          {lastSaved
            ? `Last saved ${lastSaved.toLocaleString()}`
            : editorReadOnly
              ? 'Read-only — no changes saved'
              : 'Save your draft before switching sections'}
        </p>
      </div>
    </div>
  )
}
