import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '@/shared/ui/dialog'
import { Sheet, SheetContent } from '@/shared/ui/sheet'
import { ProfileCockpitHeader } from '@/features/profiles/workspace/ProfileCockpitHeader'
import { ProfilePublishImpactSheet } from '@/features/profiles/workspace/ProfilePublishImpactSheet'
import { ProfileWorkspaceContent } from '@/features/profiles/workspace/ProfileWorkspaceContent'
import { ProfileWorkspaceSidebar } from '@/features/profiles/workspace/ProfileWorkspaceSidebar'
import {
  useIsMobileViewport,
  useProfileWorkspace,
} from '@/features/profiles/workspace/profileWorkspaceState'
import { subscribeProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceEvents'
import { saveProfileEditorIfDirty } from '@/features/profiles/workspace/profileEditorBridge'
import { getProfileSummary, type ProfileSummary } from '@/features/profiles/profileHubService'
import { useToast } from '@/shared/hooks/use-toast'
import { useEffect, useState } from 'react'

function WorkspaceBody({
  summary,
  summaryLoading,
  onClose,
  onSummaryRefresh,
}: {
  summary: ProfileSummary | null
  summaryLoading: boolean
  onClose: () => void
  onSummaryRefresh: () => void
}) {
  const { setSection, setSecondaryPanel, secondaryPanel, profileId } = useProfileWorkspace()
  const { toast } = useToast()
  const [publishOpen, setPublishOpen] = useState(false)

  useEffect(() => {
    if (secondaryPanel === 'publish-impact') {
      setPublishOpen(true)
      setSecondaryPanel(null)
    }
  }, [secondaryPanel, setSecondaryPanel])

  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <ProfileCockpitHeader
        summary={summary}
        loading={summaryLoading}
        onClose={onClose}
        onEdit={() => setSection('editor')}
        onPublish={
          summary?.canPublish && summary.draftVersionId
            ? () => {
                void (async () => {
                  const saved = await saveProfileEditorIfDirty()
                  if (!saved) {
                    toast({
                      title: 'Save required',
                      description: 'Save your draft changes before publishing.',
                      variant: 'destructive',
                    })
                    return
                  }
                  setSecondaryPanel('publish-impact')
                })()
              }
            : undefined
        }
      />
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden md:flex-row">
        <ProfileWorkspaceSidebar />
        <div className="min-h-0 min-w-0 flex-1 overflow-y-auto">
          <ProfileWorkspaceContent />
        </div>
      </div>
      {profileId != null && summary?.draftVersionId != null ? (
        <ProfilePublishImpactSheet
          open={publishOpen}
          profileId={profileId}
          versionId={summary.draftVersionId}
          draftVersionNumber={summary.publishedVersionNumber != null ? summary.publishedVersionNumber + 1 : null}
          onClose={() => setPublishOpen(false)}
          onPublished={() => {
            setPublishOpen(false)
            onSummaryRefresh()
          }}
        />
      ) : null}
    </div>
  )
}

export function ProfileWorkspace() {
  const mobile = useIsMobileViewport()
  const { profileId, close, editorDirty, section, setSection } = useProfileWorkspace()
  const open = profileId != null
  const [summary, setSummary] = useState<ProfileSummary | null>(null)
  const [summaryLoading, setSummaryLoading] = useState(false)

  const refreshSummary = () => {
    if (profileId == null) return
    setSummaryLoading(true)
    void getProfileSummary(profileId)
      .then(setSummary)
      .catch(() => setSummary(null))
      .finally(() => setSummaryLoading(false))
  }

  useEffect(() => {
    if (profileId == null) {
      setSummary(null)
      return
    }
    refreshSummary()
  }, [profileId, section])

  useEffect(() => {
    if (profileId == null) return
    return subscribeProfileWorkspace(profileId, refreshSummary)
    // eslint-disable-next-line react-hooks/exhaustive-deps -- stable refresh for event bus
  }, [profileId])

  const handleClose = () => {
    if (editorDirty && section === 'editor') {
      if (!window.confirm('You have unsaved editor changes. Close anyway?')) return
    }
    close()
  }

  const body = (
    <WorkspaceBody
      summary={summary}
      summaryLoading={summaryLoading}
      onClose={handleClose}
      onSummaryRefresh={() => {
        refreshSummary()
        if (profileId != null) {
          setSection('overview')
        }
      }}
    />
  )

  if (mobile) {
    return (
      <Sheet open={open} onOpenChange={(v) => !v && handleClose()}>
        <SheetContent side="bottom" className="!flex h-[100dvh] max-h-[100dvh] w-full flex-col p-0 sm:max-w-full">
          <DialogTitle className="sr-only">Profile workspace</DialogTitle>
          <DialogDescription className="sr-only">Profile management workspace</DialogDescription>
          {body}
        </SheetContent>
      </Sheet>
    )
  }

  return (
    <Dialog open={open} onOpenChange={(v) => !v && handleClose()}>
      <DialogContent
        className="!flex h-[94vh] max-h-[94vh] w-[96vw] max-w-[96vw] flex-col gap-0 overflow-hidden p-0 sm:max-w-[96vw]"
        onInteractOutside={(e) => e.preventDefault()}
      >
        <DialogTitle className="sr-only">Profile workspace</DialogTitle>
        <DialogDescription className="sr-only">Profile management workspace</DialogDescription>
        {body}
      </DialogContent>
    </Dialog>
  )
}
