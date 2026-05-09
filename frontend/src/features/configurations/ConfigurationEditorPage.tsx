import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Textarea } from '@/shared/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Checkbox } from '@/shared/ui/checkbox'
import * as configurationService from '@/features/configurations/configurationService'
import { getConfigurationQrEligibility } from '@/features/configurations/configurationQr'
import { ConfigurationCommonTab } from '@/features/configurations/ConfigurationCommonTab'
import { ConfigurationDesignTab } from '@/features/configurations/ConfigurationDesignTab'
import { ConfigurationApplicationsTab } from '@/features/configurations/ConfigurationApplicationsTab'
import { ConfigurationAppSettingsTab } from '@/features/configurations/ConfigurationAppSettingsTab'
import { ConfigurationFilesTab } from '@/features/configurations/ConfigurationFilesTab'
import { hasPermission } from '@/features/auth/permissions'
import type { Configuration } from '@/features/configurations/types'

interface AppOption {
  id: number
  name: string
}
interface MdmAppOption {
  id: number
  appId: number
  usedVersionId: number
  action: number
  name: string
}

function toText(value: unknown): string {
  return value == null ? '' : String(value)
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
  const [diagnosticCounts, setDiagnosticCounts] = useState<{
    allAppsRaw: number
    allAppsMapped: number
    configAppsRaw: number
    configAppsMapped: number
  } | null>(null)
  const [activeTab, setActiveTab] = useState<'common' | 'mdm' | 'design' | 'applications' | 'appSettings' | 'files'>('common')

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
        setConfiguration(cfg)
        const allAppsRaw = Array.isArray(allApps) ? allApps.length : 0
        const cfgAppsRaw = Array.isArray(cfgApps) ? cfgApps.length : 0
        setApplications(
          (Array.isArray(allApps) ? allApps : [])
            .map((item) => ({
              id: Number((item as { id?: unknown }).id ?? 0),
              name: String((item as { name?: unknown }).name ?? '').trim(),
            }))
            .filter((item) => item.id > 0)
        )
        const mapped = (Array.isArray(cfgApps) ? cfgApps : [])
          .map((item) => {
            const rec = item as Record<string, unknown>
            const appId = Number(
              rec.applicationId ??
                rec.appId ??
                rec.id ??
                0
            )
            const usedVersionId = Number(
              rec.usedVersionId ??
                rec.applicationVersionId ??
                rec.versionId ??
                rec.latestVersion ??
                0
            )
            const action = Number(rec.action ?? 0)
            return {
              id: usedVersionId > 0 ? usedVersionId : appId,
              appId,
              usedVersionId,
              action,
              name: String(rec.name ?? rec.applicationName ?? rec.pkg ?? '').trim(),
            }
          })
          .filter((item) => item.id > 0)
        setMdmApplications(mapped)
        setDiagnosticCounts({
          allAppsRaw,
          allAppsMapped: (Array.isArray(allApps) ? allApps : []).filter(
            (item) => Number((item as { id?: unknown }).id ?? 0) > 0
          ).length,
          configAppsRaw: cfgAppsRaw,
          configAppsMapped: mapped.length,
        })
      })
      .catch((reason: unknown) => {
        setError(reason instanceof Error ? reason.message : 'Failed to load configuration editor.')
      })
      .finally(() => setLoading(false))
  }, [configId])

  const selectableMdmApps = useMemo(() => {
    if (mdmApplications.length === 0) {
      const mapped = applications.map((app) => ({
        id: app.id,
        appId: app.id,
        usedVersionId: app.id,
        action: 1,
        name: app.name,
      }))
      if (mapped.length === 0 && configuration?.mainAppId && configuration.mainAppId > 0) {
        mapped.push({
          id: configuration.mainAppId,
          appId: configuration.mainAppId,
          usedVersionId: configuration.mainAppId,
          action: 1,
          name: `Current Main App (Version #${configuration.mainAppId})`,
        })
      }
      if (mapped.length === 0 && configuration?.contentAppId && configuration.contentAppId > 0) {
        mapped.push({
          id: configuration.contentAppId,
          appId: configuration.contentAppId,
          usedVersionId: configuration.contentAppId,
          action: 1,
          name: `Current Content App (Version #${configuration.contentAppId})`,
        })
      }
      return mapped
    }
    const installable = mdmApplications.filter((app) => app.action === 1)
    if (installable.length > 0) return installable
    return mdmApplications
  }, [mdmApplications, applications, configuration?.mainAppId, configuration?.contentAppId])

  const qrEligibility = useMemo(
    () => getConfigurationQrEligibility(configuration),
    [configuration]
  )

  const validateBeforeSave = (cfg: Configuration): string | null => {
    if (!String(cfg.name ?? '').trim()) return 'Name is required.'
    if (!String(cfg.password ?? '').trim()) return 'Admin password is required.'
    if (!String(cfg.pushOptions ?? '').trim()) return 'Push options are required.'
    if (cfg.kioskMode && !(Number(cfg.contentAppId ?? 0) > 0)) {
      return 'Content app is required when kiosk mode is enabled.'
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
      await configurationService.saveConfiguration(configuration)
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
      setConfiguration(next)
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
          <p className="text-sm text-muted-foreground">MDM block (phase 1): QR-critical fields.</p>
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

      <Card>
        <CardHeader>
          <CardTitle>QR readiness</CardTitle>
          <CardDescription>
            {qrEligibility.eligible ? 'Configuration is eligible for QR generation.' : qrEligibility.reason}
          </CardDescription>
        </CardHeader>
      </Card>
      {diagnosticCounts ? (
        <p className="text-xs text-muted-foreground">
          Apps diagnostic — all(raw/mapped): {diagnosticCounts.allAppsRaw}/{diagnosticCounts.allAppsMapped}, config(raw/mapped): {diagnosticCounts.configAppsRaw}/{diagnosticCounts.configAppsMapped}
        </p>
      ) : null}

      <div className="flex flex-wrap gap-2">
        {[
          { key: 'common', label: 'Common' },
          { key: 'mdm', label: 'MDM' },
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
            {activeTab === 'design' && 'Design settings'}
            {activeTab === 'applications' && 'Applications'}
            {activeTab === 'appSettings' && 'Application settings'}
            {activeTab === 'files' && 'Files'}
          </CardTitle>
          <CardDescription>
            {activeTab === 'mdm' && 'Main app, receiver component, content app, and provisioning options.'}
            {activeTab === 'common' && 'Core validation fields and behavior flags.'}
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
            <>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label>Main app</Label>
                  <Select
                    value={configuration.mainAppId != null && configuration.mainAppId > 0 ? String(configuration.mainAppId) : 'none'}
                    onValueChange={(value) =>
                      setConfiguration((current) =>
                        current ? { ...current, mainAppId: value === 'none' ? null : Number(value) } : current
                      )
                    }
                  >
                    <SelectTrigger><SelectValue placeholder="Select main app" /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">None</SelectItem>
                      {selectableMdmApps.map((app) => (
                        <SelectItem key={app.id} value={String(app.id)}>
                          {app.name || `Application #${app.id}`}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {selectableMdmApps.length === 0 ? (
                    <p className="text-xs text-muted-foreground">
                      No applications were returned by backend for this customer/session.
                    </p>
                  ) : null}
                  {selectableMdmApps.length > 0 && selectableMdmApps.every((app) => app.usedVersionId === app.appId) ? (
                    <p className="text-xs text-muted-foreground">
                      Version mapping is missing in backend response; using application ids as fallback.
                    </p>
                  ) : null}
                </div>
                <div className="space-y-2">
                  <Label>Content app</Label>
                  <Select
                    value={configuration.contentAppId != null && configuration.contentAppId > 0 ? String(configuration.contentAppId) : 'none'}
                    onValueChange={(value) =>
                      setConfiguration((current) =>
                        current ? { ...current, contentAppId: value === 'none' ? null : Number(value) } : current
                      )
                    }
                  >
                    <SelectTrigger><SelectValue placeholder="Select content app" /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">None</SelectItem>
                      {selectableMdmApps.map((app) => (
                        <SelectItem key={app.id} value={String(app.id)}>
                          {app.name || `Application #${app.id}`}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="space-y-2">
                <Label>Event receiving component</Label>
                <Input
                  placeholder="com.example/.AdminReceiver"
                  value={toText(configuration.eventReceivingComponent)}
                  onChange={(event) =>
                    setConfiguration((current) =>
                      current ? { ...current, eventReceivingComponent: event.target.value } : current
                    )
                  }
                />
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label>Launcher URL override</Label>
                  <Input
                    placeholder="https://..."
                    value={toText(configuration.launcherUrl)}
                    onChange={(event) =>
                      setConfiguration((current) =>
                        current ? { ...current, launcherUrl: event.target.value } : current
                      )
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label>Wi-Fi SSID</Label>
                  <Input
                    value={toText(configuration.wifiSSID)}
                    onChange={(event) =>
                      setConfiguration((current) =>
                        current ? { ...current, wifiSSID: event.target.value } : current
                      )
                    }
                  />
                </div>
              </div>
              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label>Wi-Fi password</Label>
                  <Input
                    value={toText(configuration.wifiPassword)}
                    onChange={(event) =>
                      setConfiguration((current) =>
                        current ? { ...current, wifiPassword: event.target.value } : current
                      )
                    }
                  />
                </div>
                <div className="space-y-2">
                  <Label>Wi-Fi security type</Label>
                  <Input
                    value={toText(configuration.wifiSecurityType)}
                    onChange={(event) =>
                      setConfiguration((current) =>
                        current ? { ...current, wifiSecurityType: event.target.value } : current
                      )
                    }
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label>QR parameters</Label>
                <Textarea
                  rows={3}
                  value={toText(configuration.qrParameters)}
                  onChange={(event) =>
                    setConfiguration((current) =>
                      current ? { ...current, qrParameters: event.target.value } : current
                    )
                  }
                />
              </div>
              <div className="space-y-2">
                <Label>Admin extras</Label>
                <Textarea
                  rows={3}
                  value={toText(configuration.adminExtras)}
                  onChange={(event) =>
                    setConfiguration((current) =>
                      current ? { ...current, adminExtras: event.target.value } : current
                    )
                  }
                />
              </div>
              <div className="flex items-center gap-6">
                <div className="flex items-center gap-2">
                  <Checkbox
                    checked={Boolean(configuration.mobileEnrollment)}
                    onCheckedChange={(checked) =>
                      setConfiguration((current) =>
                        current ? { ...current, mobileEnrollment: checked === true } : current
                      )
                    }
                  />
                  <Label>Mobile enrollment</Label>
                </div>
                <div className="flex items-center gap-2">
                  <Checkbox
                    checked={Boolean(configuration.encryptDevice)}
                    onCheckedChange={(checked) =>
                      setConfiguration((current) =>
                        current ? { ...current, encryptDevice: checked === true } : current
                      )
                    }
                  />
                  <Label>Encrypt device</Label>
                </div>
              </div>
            </>
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
