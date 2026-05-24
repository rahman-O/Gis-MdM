import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import * as profileService from '@/features/profiles/profileService'
import * as configurationService from '@/features/configurations/configurationService'
import { ConfigurationCommonTab } from '@/features/configurations/ConfigurationCommonTab'
import { ConfigurationDesignTab } from '@/features/configurations/ConfigurationDesignTab'
import { ConfigurationApplicationsTab } from '@/features/configurations/ConfigurationApplicationsTab'
import { ConfigurationAppSettingsTab } from '@/features/configurations/ConfigurationAppSettingsTab'
import { ConfigurationFilesTab } from '@/features/configurations/ConfigurationFilesTab'
import { ConfigurationMdmTab } from '@/features/configurations/ConfigurationMdmTab'
import { ConfigurationRestrictionsTab } from '@/features/configurations/ConfigurationRestrictionsTab'
import { hasPermission } from '@/features/auth/permissions'
import {
  ensureLinkedRowsForChosenVersions,
  normalizePolicyLocksForEditor,
  profileApplicationsForSaveFromApi,
} from '@/features/profiles/profileNormalize'
import { ProfileDisableBanner } from '@/features/profiles/ProfileDisableBanner'
import { ProfilePublishDialog } from '@/features/profiles/ProfilePublishDialog'
import { ProfileRolloutStatusPanel } from '@/features/profiles/ProfileRolloutStatusPanel'
import { ProfileTreeAssignmentPanel } from '@/features/profiles/ProfileTreeAssignmentPanel'
import { ProfileUsagePanel } from '@/features/profiles/ProfileUsagePanel'
import { ProfileVersionSelect } from '@/features/profiles/ProfileVersionSelect'
import type { Profile, ProfileMeta } from '@/features/profiles/types'
import {
  registerProfileEditorBridge,
  unregisterProfileEditorBridge,
} from '@/features/profiles/workspace/profileEditorBridge'
import { notifyProfileWorkspace } from '@/features/profiles/workspace/profileWorkspaceEvents'

interface AppOption {
  id: number
  name: string
  latestVersionId?: number | null
}
interface MdmAppOption {
  applicationId: number
  versionId: number
  action: number
  name: string
}

interface ProfileEditorPageProps {
  embedded?: boolean
  profileIdProp?: number
  versionIdOverride?: number | null
  readOnlyOverride?: boolean
  workspaceVersionSelect?: boolean
  workspaceChrome?: boolean
  hideWorkspacePublish?: boolean
  onWorkspaceVersionChange?: (versionId: number, readOnly: boolean) => void
  onDirtyChange?: (dirty: boolean) => void
  onLastSaved?: (at: Date) => void
}

export function ProfileEditorPage({
  embedded = false,
  profileIdProp,
  versionIdOverride,
  readOnlyOverride = false,
  workspaceVersionSelect = false,
  workspaceChrome = false,
  hideWorkspacePublish = false,
  onWorkspaceVersionChange,
  onDirtyChange,
  onLastSaved,
}: ProfileEditorPageProps = {}) {
  const params = useParams<{ profileId: string; versionId?: string }>()
  const navigate = useNavigate()
  const profileId = profileIdProp ?? Number(params.profileId)
  const versionIdParam = params.versionId != null ? Number(params.versionId) : null
  const initialVersion =
    versionIdOverride ??
    (Number.isFinite(versionIdParam) && versionIdParam! > 0 ? versionIdParam : null)
  const [versionId, setVersionId] = useState<number | null>(initialVersion)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [saveSuccess, setSaveSuccess] = useState<string | null>(null)
  const [profile, setProfile] = useState<Profile | null>(null)
  const [applications, setApplications] = useState<AppOption[]>([])
  const [mdmApplications, setMdmApplications] = useState<MdmAppOption[]>([])
  const [activeTab, setActiveTab] = useState<
    'common' | 'mdm' | 'restrictions' | 'design' | 'applications' | 'appSettings' | 'files'
  >('common')
  const [pageSection, setPageSection] = useState<'editor' | 'assignment' | 'rollout'>('editor')
  const [meta, setMeta] = useState<ProfileMeta | null>(null)
  const [publishOpen, setPublishOpen] = useState(false)
  const [dirty, setDirty] = useState(false)
  const activeVersionId = versionIdOverride ?? versionId
  const loadGenRef = useRef(0)

  useEffect(() => {
    onDirtyChange?.(dirty)
  }, [dirty, onDirtyChange])

  useEffect(() => {
    if (versionIdOverride != null) {
      setVersionId(versionIdOverride)
    }
  }, [versionIdOverride])

  useEffect(() => {
    if (!Number.isFinite(profileId) || profileId <= 0) {
      setLoading(false)
      setError('Invalid profile id.')
      return
    }

    const loadGen = ++loadGenRef.current
    setLoading(true)
    setError(null)
    setProfile(null)

    void (async () => {
      try {
        let vid = versionIdOverride ?? versionId
        const profileMeta = await profileService.getProfileMeta(profileId)
        if (loadGen !== loadGenRef.current) return
        setMeta(profileMeta)
        if (vid == null) {
          if (readOnlyOverride && profileMeta.publishedVersionId) {
            vid = profileMeta.publishedVersionId
          } else if (profileMeta.draftVersionId == null) {
            setError('No draft version available for this profile.')
            return
          } else {
            vid = profileMeta.draftVersionId
          }
          setVersionId(vid)
        }
        const [cfg, allApps] = await Promise.all([
          profileService.getProfileVersion(profileId, vid),
          configurationService.getAllApplications(),
        ])
        if (loadGen !== loadGenRef.current) return
        const cfgApps = Array.isArray(cfg.applications) ? cfg.applications : []
        setProfile(
          normalizePolicyLocksForEditor({
            ...cfg,
            id: profileId,
            profileId,
            versionId: vid,
            versionNumber: cfg.versionNumber,
            versionStatus: cfg.versionStatus,
            applications: profileApplicationsForSaveFromApi(cfgApps),
          })
        )
        setApplications(Array.isArray(allApps) ? allApps : [])
        const mapped = cfgApps
          .map((item) => {
            const rec = item as Record<string, unknown>
            const applicationId = Number(rec.applicationId ?? rec.appId ?? rec.id ?? 0)
            const linkedVersionId = Number(rec.usedVersionId ?? rec.applicationVersionId ?? 0)
            const latestVersionId = Number(rec.latestVersion ?? 0)
            const versionIdResolved =
              linkedVersionId > 0 ? linkedVersionId : latestVersionId > 0 ? latestVersionId : 0
            const rawAct = rec.action
            const action = rawAct === undefined || rawAct === null ? 1 : Number(rawAct)
            return {
              applicationId,
              versionId: versionIdResolved,
              action,
              name: String(rec.name ?? rec.applicationName ?? rec.pkg ?? '').trim(),
            }
          })
          .filter((item) => item.applicationId > 0 && item.versionId > 0)
        setMdmApplications(mapped)
        setDirty(false)
        setSaveSuccess(null)
        setSaveError(null)
      } catch (reason: unknown) {
        if (loadGen !== loadGenRef.current) return
        setError(reason instanceof Error ? reason.message : 'Failed to load profile editor.')
      } finally {
        if (loadGen === loadGenRef.current) {
          setLoading(false)
        }
      }
    })()
  }, [profileId, versionIdOverride, readOnlyOverride])

  const readOnly =
    readOnlyOverride || profile?.versionStatus === 'published' || profile?.versionStatus === 'archived'
  const isDraftVersion =
    profile?.versionStatus === 'draft' ||
    (!profile?.versionStatus && meta?.draftVersionId === activeVersionId)
  const profileEnabled = meta?.enabled !== false

  const onProfileChange = (next: Profile) => {
    if (readOnly) return
    setDirty(true)
    setProfile(next)
  }

  const selectableMdmApps = useMemo(() => {
    if (mdmApplications.length === 0) {
      return applications
        .filter((app) => Number(app.latestVersionId ?? 0) > 0)
        .map((app) => ({
          applicationId: app.id,
          versionId: Number(app.latestVersionId),
          action: 1,
          name: app.name,
        }))
    }
    const installable = mdmApplications.filter((app) => app.action === 1)
    return installable.length > 0 ? installable : mdmApplications
  }, [mdmApplications, applications])

  const versionCatalogForSave = useMemo(() => {
    const out: { applicationId: number; versionId: number; name: string; action: number }[] = []
    const seen = new Set<number>()
    const add = (m: MdmAppOption) => {
      if (m.applicationId <= 0 || m.versionId <= 0) return
      if (seen.has(m.versionId)) return
      seen.add(m.versionId)
      out.push({
        applicationId: m.applicationId,
        versionId: m.versionId,
        name: m.name,
        action: m.action,
      })
    }
    mdmApplications.forEach(add)
    selectableMdmApps.forEach(add)
    return out
  }, [mdmApplications, selectableMdmApps])

  const validateBeforeSave = (cfg: Profile): string | null => {
    if (!String(cfg.name ?? '').trim()) return 'Common: name is required.'
    if (!String(cfg.password ?? '').trim()) return 'Common: admin password is required.'
    if (!String(cfg.pushOptions ?? '').trim()) return 'Common: push options are required.'
    if (cfg.kioskMode && !(Number(cfg.contentAppId ?? 0) > 0)) {
      return 'MDM: content app is required when kiosk mode is enabled.'
    }
    return null
  }

  const handleSave = useCallback(async (): Promise<boolean> => {
    const vid = versionIdOverride ?? versionId
    if (!profile || vid == null) return false
    if (readOnlyOverride || profile.versionStatus === 'published' || profile.versionStatus === 'archived') {
      return true
    }
    const validationError = validateBeforeSave(profile)
    if (validationError) {
      setSaveError(validationError)
      return false
    }
    setSaving(true)
    setSaveError(null)
    setSaveSuccess(null)
    try {
      const applicationsPayload = ensureLinkedRowsForChosenVersions(
        profile.applications,
        profile.mainAppId,
        profile.contentAppId,
        versionCatalogForSave
      )
      await profileService.saveProfileVersion(
        profileId,
        vid,
        normalizePolicyLocksForEditor({
          ...profile,
          applications: applicationsPayload,
        })
      )
      const fresh = await profileService.getProfileVersion(profileId, vid)
      const freshApps = Array.isArray(fresh.applications) ? fresh.applications : []
      setProfile(
        normalizePolicyLocksForEditor({
          ...fresh,
          id: profileId,
          profileId,
          versionId: vid,
          versionNumber: fresh.versionNumber,
          versionStatus: fresh.versionStatus,
          applications: profileApplicationsForSaveFromApi(freshApps),
        })
      )
      setSaveSuccess('Profile draft saved.')
      setDirty(false)
      onLastSaved?.(new Date())
      notifyProfileWorkspace(profileId)
      return true
    } catch (reason: unknown) {
      setSaveError(reason instanceof Error ? reason.message : 'Failed to save profile.')
      return false
    } finally {
      setSaving(false)
    }
  }, [
    profile,
    versionId,
    versionIdOverride,
    readOnlyOverride,
    profileId,
    versionCatalogForSave,
    onLastSaved,
  ])

  useEffect(() => {
    if (!embedded || !workspaceChrome) {
      return
    }
    registerProfileEditorBridge({
      saveIfDirty: async () => {
        if (!dirty) return true
        return handleSave()
      },
      getActiveVersionId: () => versionIdOverride ?? versionId,
    })
    return () => unregisterProfileEditorBridge()
  }, [embedded, workspaceChrome, dirty, handleSave, versionId, versionIdOverride])

  if (loading) {
    return <div className="text-sm text-muted-foreground">Loading profile editor...</div>
  }

  if (error) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4" />
          <span>{error}</span>
        </div>
        {!embedded ? (
          <Button variant="outline" onClick={() => navigate('/profiles')}>
            Back
          </Button>
        ) : null}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {!embedded ? (
        <>
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl font-semibold tracking-tight">Profile Editor</h1>
              <p className="text-sm text-muted-foreground">
                {readOnly
                  ? 'Viewing a published version (read-only). Fork a draft to edit.'
                  : 'Edit draft → publish → assign to a tree folder → monitor rollout status.'}
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              {activeVersionId != null ? (
                <ProfileVersionSelect
                  profileId={profileId}
                  currentVersionId={activeVersionId}
                  isDraft={isDraftVersion}
                  dirty={dirty}
                  workspaceMode={workspaceVersionSelect}
                  onForkComplete={() => {
                    void profileService.getProfileMeta(profileId).then((meta) => {
                      if (meta.draftVersionId) {
                        onWorkspaceVersionChange?.(meta.draftVersionId, false)
                      }
                    })
                  }}
                  onVersionChange={(vid, ro) => {
                    setDirty(false)
                    if (workspaceVersionSelect && onWorkspaceVersionChange) {
                      onWorkspaceVersionChange(vid, ro)
                    } else {
                      setVersionId(vid)
                    }
                  }}
                />
              ) : null}
              <Button variant="outline" onClick={() => navigate('/profiles')}>
                Cancel
              </Button>
              {!readOnly ? (
                <>
                  <Button disabled={saving || !profile} onClick={() => void handleSave()}>
                    {saving ? 'Saving...' : 'Save draft'}
                  </Button>
                  {activeVersionId != null ? (
                    <Button variant="secondary" disabled={saving} onClick={() => setPublishOpen(true)}>
                      Publish
                    </Button>
                  ) : null}
                </>
              ) : null}
            </div>
          </div>

          <ProfileDisableBanner
            profileId={profileId}
            enabled={profileEnabled}
            onChanged={(enabled) => setMeta((m) => (m ? { ...m, enabled } : m))}
          />

          <ProfileUsagePanel meta={meta} />

          <div className="flex flex-wrap gap-2 border-b pb-2">
            {(
              [
                { key: 'editor', label: 'Policy editor' },
                { key: 'assignment', label: 'Tree assignment' },
                { key: 'rollout', label: 'Rollout status' },
              ] as const
            ).map((s) => (
              <Button
                key={s.key}
                type="button"
                variant={pageSection === s.key ? 'default' : 'outline'}
                onClick={() => setPageSection(s.key)}
              >
                {s.label}
              </Button>
            ))}
          </div>

          {pageSection === 'assignment' ? (
            <ProfileTreeAssignmentPanel
              profileId={profileId}
              publishedVersionId={meta?.publishedVersionId}
            />
          ) : null}

          {pageSection === 'rollout' ? <ProfileRolloutStatusPanel profileId={profileId} /> : null}
        </>
      ) : workspaceChrome ? (
        <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
          {activeVersionId != null ? (
            <ProfileVersionSelect
              profileId={profileId}
              currentVersionId={activeVersionId}
              isDraft={isDraftVersion}
              dirty={dirty}
              workspaceMode={workspaceVersionSelect}
              onForkComplete={() => {
                void profileService.getProfileMeta(profileId).then((meta) => {
                  if (meta.draftVersionId) {
                    onWorkspaceVersionChange?.(meta.draftVersionId, false)
                  }
                })
              }}
              onVersionChange={(vid, ro) => {
                setDirty(false)
                if (workspaceVersionSelect && onWorkspaceVersionChange) {
                  onWorkspaceVersionChange(vid, ro)
                } else {
                  setVersionId(vid)
                }
              }}
            />
          ) : null}
          {!readOnly ? (
            <Button disabled={saving || !profile} size="sm" onClick={() => void handleSave()}>
              {saving ? 'Saving...' : 'Save draft'}
            </Button>
          ) : null}
        </div>
      ) : (
        <div className="flex flex-wrap items-center justify-between gap-2 border-b border-amber-500/40 bg-amber-50/50 px-3 py-2 dark:bg-amber-950/20">
          <p className="text-sm text-amber-900 dark:text-amber-100">
            Editing production policy draft. Save before leaving the editor section.
          </p>
          <div className="flex flex-wrap gap-2">
            {activeVersionId != null ? (
              <ProfileVersionSelect
                profileId={profileId}
                currentVersionId={activeVersionId}
                isDraft={isDraftVersion}
                dirty={dirty}
                workspaceMode={workspaceVersionSelect}
                onForkComplete={() => {
                  void profileService.getProfileMeta(profileId).then((meta) => {
                    if (meta.draftVersionId) {
                      onWorkspaceVersionChange?.(meta.draftVersionId, false)
                    }
                  })
                }}
                onVersionChange={(vid, ro) => {
                  setDirty(false)
                  if (workspaceVersionSelect && onWorkspaceVersionChange) {
                    onWorkspaceVersionChange(vid, ro)
                  } else {
                    setVersionId(vid)
                  }
                }}
              />
            ) : null}
            {!readOnly ? (
              <Button disabled={saving || !profile} size="sm" onClick={() => void handleSave()}>
                {saving ? 'Saving...' : 'Save draft'}
              </Button>
            ) : null}
          </div>
        </div>
      )}

      {(embedded || pageSection === 'editor') ? (
        <>
      <div className="flex flex-wrap gap-2">
        {[
          { key: 'common', label: 'Common' },
          { key: 'mdm', label: 'MDM' },
          { key: 'restrictions', label: 'Restrictions' },
          { key: 'design', label: 'Design' },
          { key: 'applications', label: 'Applications' },
          { key: 'appSettings', label: 'App Settings' },
          { key: 'files', label: 'Files' },
        ].map((tab) => (
          <Button
            key={tab.key}
            type="button"
            variant={activeTab === tab.key ? 'default' : 'outline'}
            onClick={() => setActiveTab(tab.key as typeof activeTab)}
          >
            {tab.label}
          </Button>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>
            {activeTab === 'common' && 'Common settings'}
            {activeTab === 'mdm' && 'MDM settings'}
            {activeTab === 'restrictions' && 'Restrictions'}
            {activeTab === 'design' && 'Design settings'}
            {activeTab === 'applications' && 'Applications'}
            {activeTab === 'appSettings' && 'Application settings'}
            {activeTab === 'files' && 'Files'}
          </CardTitle>
          <CardDescription>
            {activeTab === 'restrictions' && 'Device restrictions and connectivity.'}
            {activeTab === 'mdm' && 'Main app, receiver, and provisioning.'}
            {activeTab === 'common' && 'Name and general options.'}
            {activeTab === 'design' && 'Theme and display settings.'}
            {activeTab === 'applications' && 'Linked applications and versions.'}
            {activeTab === 'appSettings' && 'Profile-level app settings.'}
            {activeTab === 'files' && 'Linked file definitions.'}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {profile && activeTab === 'common' ? (
            <ConfigurationCommonTab
              configuration={profile}
              applications={applications}
              onChange={onProfileChange}
            />
          ) : null}
          {profile && activeTab === 'mdm' ? (
            <ConfigurationMdmTab
              configuration={profile}
              selectableMdmApps={selectableMdmApps}
              onChange={onProfileChange}
            />
          ) : null}
          {profile && activeTab === 'restrictions' ? (
            <ConfigurationRestrictionsTab configuration={profile} onChange={onProfileChange} />
          ) : null}
          {profile && activeTab === 'design' ? (
            <ConfigurationDesignTab configuration={profile} onChange={onProfileChange} />
          ) : null}
          {profile && activeTab === 'applications' ? (
            <ConfigurationApplicationsTab
              configuration={profile}
              applications={applications}
              upgrading={false}
              onChange={onProfileChange}
              onUpgrade={async () => {
                setSaveError('Application upgrade is not available in the profile editor yet.')
              }}
            />
          ) : null}
          {profile && activeTab === 'appSettings' ? (
            <ConfigurationAppSettingsTab
              configuration={profile}
              applications={applications}
              onChange={onProfileChange}
            />
          ) : null}
          {profile && activeTab === 'files' ? (
            <ConfigurationFilesTab configuration={profile} onChange={onProfileChange} />
          ) : null}
          {!hasPermission('copy_config') && activeTab === 'applications' ? (
            <p className="text-xs text-muted-foreground">Limited permissions: some actions may be restricted.</p>
          ) : null}
          {saveSuccess ? <p className="text-sm text-emerald-600">{saveSuccess}</p> : null}
          {saveError ? <p className="text-sm text-destructive">{saveError}</p> : null}
        </CardContent>
      </Card>

      {activeVersionId != null && !readOnly && !(embedded && hideWorkspacePublish) ? (
        <ProfilePublishDialog
          open={publishOpen}
          profileId={profileId}
          versionId={activeVersionId}
          onClose={() => setPublishOpen(false)}
          onPublished={() => {
            setSaveSuccess('Profile published.')
            void profileService.getProfileMeta(profileId).then(setMeta)
          }}
        />
      ) : null}
        </>
      ) : null}
    </div>
  )
}
