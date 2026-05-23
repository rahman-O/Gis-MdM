import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
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
  configurationApplicationsForSaveFromApi,
  ensureLinkedRowsForChosenVersions,
  normalizePolicyLocksForEditor,
} from '@/features/configurations/configurationNormalize'
import type { Configuration } from '@/features/configurations/types'

interface AppOption {
  id: number
  name: string
  /** Backend `applications.latestVersion` (id of newest `applicationVersions` row). */
  latestVersionId?: number | null
}
interface MdmAppOption {
  applicationId: number
  /** `applicationVersions.id` — same meaning as stored `configuration.mainAppId` / `contentAppId`. */
  versionId: number
  action: number
  name: string
}

export function ConfigurationEditorPage() {
  const params = useParams<{ id: string }>()
  const navigate = useNavigate()
  const configId = Number(params.id)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [saveSuccess, setSaveSuccess] = useState<string | null>(null)
  const [upgradingApp, setUpgradingApp] = useState(false)
  const [configuration, setConfiguration] = useState<Configuration | null>(null)
  const [applications, setApplications] = useState<AppOption[]>([])
  const [mdmApplications, setMdmApplications] = useState<MdmAppOption[]>([])
  const [activeTab, setActiveTab] = useState<
    'common' | 'mdm' | 'restrictions' | 'design' | 'applications' | 'appSettings' | 'files'
  >('common')

  useEffect(() => {
    if (!Number.isFinite(configId) || configId <= 0) {
      setLoading(false)
      setError('Invalid configuration id.')
      return
    }

    setLoading(true)
    setError(null)
    void Promise.all([
      configurationService.getConfiguration(configId),
      configurationService.getAllApplications(),
      configurationService.getConfigurationApplications(configId),
    ])
      .then(([cfg, allApps, cfgApps]) => {
        setConfiguration(
          normalizePolicyLocksForEditor({
            ...cfg,
            applications: configurationApplicationsForSaveFromApi(cfgApps),
          })
        )
        setApplications(Array.isArray(allApps) ? allApps : [])
        const mapped = (Array.isArray(cfgApps) ? cfgApps : [])
          .map((item) => {
            const rec = item as Record<string, unknown>
            const applicationId = Number(rec.applicationId ?? rec.appId ?? rec.id ?? 0)
            const linkedVersionId = Number(rec.usedVersionId ?? rec.applicationVersionId ?? 0)
            const latestVersionId = Number(rec.latestVersion ?? 0)
            const versionId =
              linkedVersionId > 0 ? linkedVersionId : latestVersionId > 0 ? latestVersionId : 0
            const rawAct = rec.action
            const action = rawAct === undefined || rawAct === null ? 1 : Number(rawAct)
            return {
              applicationId,
              versionId,
              action,
              name: String(rec.name ?? rec.applicationName ?? rec.pkg ?? '').trim(),
            }
          })
          .filter((item) => item.applicationId > 0 && item.versionId > 0)
        setMdmApplications(mapped)
      })
      .catch((reason: unknown) => {
        setError(reason instanceof Error ? reason.message : 'Failed to load configuration editor.')
      })
      .finally(() => setLoading(false))
  }, [configId])

  const selectableMdmApps = useMemo(() => {
    if (mdmApplications.length === 0) {
      const mapped: MdmAppOption[] = applications
        .filter((app) => Number(app.latestVersionId ?? 0) > 0)
        .map((app) => ({
          applicationId: app.id,
          versionId: Number(app.latestVersionId),
          action: 1,
          name: app.name,
        }))
      const addSynthetic = (vid: number | null | undefined, kind: 'main' | 'content') => {
        if (!vid || vid <= 0) return
        if (mapped.some((x) => x.versionId === vid)) return
        mapped.push({
          applicationId: 0,
          versionId: vid,
          action: 1,
          name:
            kind === 'main'
              ? `Current main app (version #${vid})`
              : `Current content app (version #${vid})`,
        })
      }
      addSynthetic(configuration?.mainAppId, 'main')
      addSynthetic(configuration?.contentAppId, 'content')
      return mapped
    }
    const installable = mdmApplications.filter((app) => app.action === 1)
    if (installable.length > 0) return installable
    return mdmApplications
  }, [mdmApplications, applications, configuration?.mainAppId, configuration?.contentAppId])

  /** Catalog rows with real `applicationId` for injecting `configurationApplications` on save */
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

  const validateBeforeSave = (cfg: Configuration): string | null => {
    if (!String(cfg.name ?? '').trim()) return 'Common: name is required.'
    if (!String(cfg.password ?? '').trim()) return 'Common: admin password is required.'
    if (!String(cfg.pushOptions ?? '').trim()) return 'Common: push options are required.'
    if (cfg.kioskMode && !(Number(cfg.contentAppId ?? 0) > 0)) {
      return 'MDM: content app is required when kiosk mode is enabled.'
    }
    return null
  }

  const handleSave = async () => {
    if (!configuration || configuration.id == null) return
    const validationError = validateBeforeSave(configuration)
    if (validationError) {
      setSaveError(validationError)
      return
    }
    setSaving(true)
    setSaveError(null)
    setSaveSuccess(null)
    try {
      const applicationsPayload = ensureLinkedRowsForChosenVersions(
        configuration.applications,
        configuration.mainAppId,
        configuration.contentAppId,
        versionCatalogForSave
      )
      const versionIdsInPayload = new Set(
        applicationsPayload
          .map((a) => Number((a as Record<string, unknown>).usedVersionId ?? 0))
          .filter((v) => v > 0)
      )
      const mainV = Number(configuration.mainAppId ?? 0)
      const contentV = Number(configuration.contentAppId ?? 0)
      if (mainV > 0 && !versionIdsInPayload.has(mainV)) {
        setSaveError(
          'Cannot persist main app: no application catalogue entry for this version id. Reload the editor and pick main app again from the list.'
        )
        return
      }
      if (contentV > 0 && !versionIdsInPayload.has(contentV)) {
        setSaveError(
          'Cannot persist content app: no application catalogue entry for this version id. Reload and pick content app again from the list.'
        )
        return
      }

      const savedCfg = await configurationService.saveConfiguration(
        normalizePolicyLocksForEditor({
          ...configuration,
          applications: applicationsPayload,
        })
      )
      const freshCfgApps = await configurationService.getConfigurationApplications(configId)
      setConfiguration(
        normalizePolicyLocksForEditor({
          ...savedCfg,
          applications: configurationApplicationsForSaveFromApi(freshCfgApps),
        })
      )
      setSaveSuccess('Configuration saved successfully.')
    } catch (reason: unknown) {
      setSaveError(reason instanceof Error ? reason.message : 'Failed to save configuration.')
    } finally {
      setSaving(false)
    }
  }

  const handleUpgrade = async (applicationId: number) => {
    if (!configuration?.id) return
    setUpgradingApp(true)
    setSaveError(null)
    try {
      const next = await configurationService.upgradeConfigurationApplication({
        configurationId: configuration.id,
        applicationId,
      })
      const cfgApps = await configurationService.getConfigurationApplications(configuration.id)
      setConfiguration({
        ...next,
        applications: configurationApplicationsForSaveFromApi(cfgApps),
      })
    } catch (reason: unknown) {
      setSaveError(reason instanceof Error ? reason.message : 'Failed to upgrade application.')
    } finally {
      setUpgradingApp(false)
    }
  }

  if (loading) {
    return <div className="text-sm text-muted-foreground">Loading configuration editor...</div>
  }

  if (error) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4" />
          <span>{error}</span>
        </div>
        <Button variant="outline" onClick={() => navigate('/configurations')}>
          Back
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Configuration Editor</h1>
          <p className="text-sm text-muted-foreground">Device policy profile — changes sync to enrolled devices after save.</p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => navigate('/configurations')}>
            Cancel
          </Button>
          <Button disabled={saving || !configuration} onClick={() => void handleSave()}>
            {saving ? 'Saving...' : 'Save'}
          </Button>
        </div>
      </div>

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
            {activeTab === 'mdm' && 'Main app, receiver, and provisioning.'}
            {activeTab === 'restrictions' && 'Device restrictions and connectivity.'}
            {activeTab === 'common' && 'Name, QR, and general options.'}
            {activeTab === 'design' && 'Theme, header, orientation, and display settings.'}
            {activeTab === 'applications' && 'Linked applications, actions, versions, and upgrade flow.'}
            {activeTab === 'appSettings' && 'Configuration-level app settings.'}
            {activeTab === 'files' && 'Default path and linked file definitions.'}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {configuration && activeTab === 'common' ? (
            <ConfigurationCommonTab
              configuration={configuration}
              applications={applications}
              onChange={setConfiguration}
            />
          ) : null}

          {configuration && activeTab === 'mdm' ? (
            <ConfigurationMdmTab
              configuration={configuration}
              selectableMdmApps={selectableMdmApps}
              onChange={setConfiguration}
            />
          ) : null}

          {configuration && activeTab === 'restrictions' ? (
            <ConfigurationRestrictionsTab configuration={configuration} onChange={setConfiguration} />
          ) : null}

          {configuration && activeTab === 'design' ? (
            <ConfigurationDesignTab configuration={configuration} onChange={setConfiguration} />
          ) : null}

          {configuration && activeTab === 'applications' ? (
            <ConfigurationApplicationsTab
              configuration={configuration}
              applications={applications}
              upgrading={upgradingApp}
              onChange={setConfiguration}
              onUpgrade={handleUpgrade}
            />
          ) : null}

          {configuration && activeTab === 'appSettings' ? (
            <ConfigurationAppSettingsTab
              configuration={configuration}
              applications={applications}
              onChange={setConfiguration}
            />
          ) : null}

          {configuration && activeTab === 'files' ? (
            <ConfigurationFilesTab configuration={configuration} onChange={setConfiguration} />
          ) : null}

          {!hasPermission('copy_config') && activeTab === 'applications' ? (
            <p className="text-xs text-muted-foreground">Limited permissions detected: upgrade/copy actions may be restricted.</p>
          ) : null}
          {saveSuccess ? <p className="text-sm text-emerald-600">{saveSuccess}</p> : null}
          {saveError ? <p className="text-sm text-destructive">{saveError}</p> : null}
        </CardContent>
      </Card>
    </div>
  )
}
