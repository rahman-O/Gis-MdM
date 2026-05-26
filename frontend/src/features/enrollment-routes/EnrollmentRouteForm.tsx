import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AlertCircle, ChevronRight, Folder, Info, Wifi, Sliders, Package } from 'lucide-react'
import type { EnrollmentRouteFormValues } from '@/features/enrollment-routes/enrollmentRouteDialogState'
import type {
  BootstrapAppOption,
  BootstrapIntent,
  TreeNodeOption,
} from '@/features/enrollment-routes/enrollmentRouteService'
import {
  listBootstrapApps,
  listTreeNodeOptions,
} from '@/features/enrollment-routes/enrollmentRouteService'
import { BootstrapAppPicker } from '@/features/enrollment-routes/BootstrapAppPicker'
import { TargetNodePicker } from '@/features/enrollment-routes/TargetNodePicker'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Textarea } from '@/shared/ui/textarea'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/ui/tabs'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'

interface Props {
  values: EnrollmentRouteFormValues
  onChange: (next: EnrollmentRouteFormValues) => void
  readOnly?: boolean
  saveError?: string | null
  treeNodes?: TreeNodeOption[]
}

export function EnrollmentRouteForm({
  values,
  onChange,
  readOnly = false,
  saveError,
  treeNodes: propTreeNodes,
}: Props) {
  const { t } = useTranslation()
  const [internalTreeNodes, setInternalTreeNodes] = useState<TreeNodeOption[]>([])
  const [bootstrapApps, setBootstrapApps] = useState<BootstrapAppOption[]>([])
  const [optionsError, setOptionsError] = useState<string | null>(null)
  const [pickerOpen, setPickerOpen] = useState(false)


  useEffect(() => {
    if (propTreeNodes && propTreeNodes.length > 0) {
      setInternalTreeNodes(propTreeNodes)
      void listBootstrapApps()
        .then(setBootstrapApps)
        .catch((e: unknown) => {
          setOptionsError(e instanceof Error ? e.message : 'Failed to load bootstrap applications.')
        })
    } else {
      void Promise.all([listTreeNodeOptions(), listBootstrapApps()])
        .then(([nodes, apps]) => {
          setInternalTreeNodes(nodes)
          setBootstrapApps(apps)
        })
        .catch((e: unknown) => {
          setOptionsError(e instanceof Error ? e.message : 'Failed to load options.')
        })
    }
  }, [propTreeNodes])

  const treeNodes = propTreeNodes && propTreeNodes.length > 0 ? propTreeNodes : internalTreeNodes

  const selectedNode = treeNodes.find((n) => n.id === values.targetNodeId)
  const showContainerAck =
    selectedNode?.placementKind === 'inheritable' && !readOnly

  const patch = (partial: Partial<EnrollmentRouteFormValues>) => {
    onChange({ ...values, ...partial })
  }

  const handleNodeSelect = (nodeId: number, _node: TreeNodeOption) => {
    patch({ targetNodeId: nodeId })
    setPickerOpen(false)
  }

  const handlePickerCancel = () => {
    setPickerOpen(false)
  }

  /** Render the selected node path as a breadcrumb — always shows folder names */
  const renderSelectedPath = () => {
    if (!selectedNode) return null
    const segments = selectedNode.path.split('/').filter(Boolean)
    return (
      <div className="flex flex-wrap items-center gap-1 text-sm text-muted-foreground">
        <Folder className="h-3.5 w-3.5 shrink-0 text-primary" />
        {segments.map((segment, idx) => {
          const nodeId = Number(segment)
          const node = treeNodes.find((n) => n.id === nodeId)
          const displayName = node ? node.name : segment
          return (
            <span key={idx} className="flex items-center gap-1">
              {idx > 0 && <ChevronRight className="h-3 w-3 shrink-0 text-muted-foreground/60" />}
              <span className={idx === segments.length - 1 ? 'font-semibold text-foreground' : ''}>
                {displayName}
              </span>
            </span>
          )
        })}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {optionsError ? (
        <p className="text-sm text-destructive">{optionsError}</p>
      ) : null}
      {saveError ? (
        <div className="flex items-start gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-2.5 text-sm text-destructive shadow-sm animate-in fade-in slide-in-from-top-1 duration-200">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>{saveError}</span>
        </div>
      ) : null}

      <Tabs defaultValue="general" className="w-full">
        <TabsList className="grid w-full grid-cols-4 bg-muted/60 p-1 rounded-xl shadow-inner mb-4">
          <TabsTrigger value="general" className="gap-1.5 rounded-lg data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all duration-200">
            <Info className="h-3.5 w-3.5 shrink-0" />
            <span className="hidden sm:inline text-xs font-medium">{t('enrollmentRoute.tabs.general')}</span>
          </TabsTrigger>
          <TabsTrigger value="bootstrap" className="gap-1.5 rounded-lg data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all duration-200">
            <Package className="h-3.5 w-3.5 shrink-0" />
            <span className="hidden sm:inline text-xs font-medium">{t('enrollmentRoute.tabs.bootstrap')}</span>
          </TabsTrigger>
          <TabsTrigger value="wifi" className="gap-1.5 rounded-lg data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all duration-200">
            <Wifi className="h-3.5 w-3.5 shrink-0" />
            <span className="hidden sm:inline text-xs font-medium">{t('enrollmentRoute.tabs.wifi')}</span>
          </TabsTrigger>
          <TabsTrigger value="advanced" className="gap-1.5 rounded-lg data-[state=active]:bg-background data-[state=active]:shadow-sm transition-all duration-200">
            <Sliders className="h-3.5 w-3.5 shrink-0" />
            <span className="hidden sm:inline text-xs font-medium">{t('enrollmentRoute.tabs.advanced')}</span>
          </TabsTrigger>
        </TabsList>

        {/* --- GENERAL SETTINGS --- */}
        <TabsContent value="general" className="space-y-4 outline-none min-h-[340px]">
          <div className="rounded-xl border bg-card/40 backdrop-blur-sm p-4 shadow-sm space-y-4">
            <div className="space-y-2">
              <Label htmlFor="route-name" className="text-sm font-semibold">{t('enrollmentRoute.form.name')}</Label>
              <Input
                id="route-name"
                value={values.name}
                disabled={readOnly}
                onChange={(e) => patch({ name: e.target.value })}
                className="transition-all duration-200 focus:ring-2 focus:ring-primary/20"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="route-desc" className="text-sm font-semibold">{t('enrollmentRoute.form.description')}</Label>
              <Input
                id="route-desc"
                value={values.description}
                disabled={readOnly}
                onChange={(e) => patch({ description: e.target.value })}
                className="transition-all duration-200 focus:ring-2 focus:ring-primary/20"
              />
            </div>

            <div className="space-y-2 border-t pt-3">
              <div className="mb-2">
                <Label className="text-sm font-semibold">{t('enrollmentRoute.form.targetFolder')}</Label>
              </div>

              {!pickerOpen && (
                <div className="space-y-2">
                  {selectedNode ? (
                    <div className="flex items-center justify-between rounded-lg border bg-background px-3 py-2 shadow-sm transition-all hover:border-muted-foreground/30">
                      {renderSelectedPath()}
                      {!readOnly && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          className="ml-2 shrink-0 h-8 hover:bg-muted"
                          onClick={() => setPickerOpen(true)}
                        >
                          {t('enrollmentRoute.form.selectFolder')}
                        </Button>
                      )}
                    </div>
                  ) : (
                    <Button
                      type="button"
                      variant="outline"
                      className="w-full justify-start text-muted-foreground border-dashed h-10 hover:border-primary/50"
                      disabled={readOnly}
                      onClick={() => setPickerOpen(true)}
                    >
                      <Folder className="mr-2 h-4 w-4 text-muted-foreground" />
                      {t('enrollmentRoute.form.selectFolder')}
                    </Button>
                  )}
                  {selectedNode?.placementKind === 'inheritable' ? (
                    <p className="text-xs text-amber-700 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/20 p-2 rounded-md border border-amber-200/50 dark:border-amber-900/30">
                      {t('enrollmentRoute.form.containerWarning')}
                    </p>
                  ) : null}
                </div>
              )}

              {/* Scrollable picker container */}
              {pickerOpen && (
                <div className="max-h-[260px] overflow-y-auto rounded-lg border bg-background shadow-inner">
                  <TargetNodePicker
                    selectedNodeId={values.targetNodeId}
                    onSelect={handleNodeSelect}
                    onCancel={handlePickerCancel}
                    open={pickerOpen}
                  />
                </div>
              )}
              {!pickerOpen && (
                <TargetNodePicker
                  selectedNodeId={values.targetNodeId}
                  onSelect={handleNodeSelect}
                  onCancel={handlePickerCancel}
                  open={false}
                />
              )}
            </div>

            <div className="space-y-2 border-t pt-3">
              <Label className="text-sm font-semibold">{t('enrollmentRoute.form.deviceIdentity')}</Label>
              <Select
                disabled={readOnly}
                value={values.deviceIdentityMode}
                onValueChange={(v) => patch({ deviceIdentityMode: v })}
              >
                <SelectTrigger className="w-full transition-all focus:ring-2 focus:ring-primary/20">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="imei">IMEI</SelectItem>
                  <SelectItem value="serial">{t('enrollmentRoute.form.serial')}</SelectItem>
                  <SelectItem value="request">{t('enrollmentRoute.form.request')}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </TabsContent>

        {/* --- BOOTSTRAP APP --- */}
        <TabsContent value="bootstrap" className="space-y-4 outline-none min-h-[340px]">
          <div className="rounded-xl border bg-card/40 backdrop-blur-sm p-4 shadow-sm space-y-4 animate-in fade-in-50 duration-200">
            <BootstrapAppPicker
              apps={bootstrapApps}
              selectedAppId={values.bootstrapApplicationId || ''}
              intent={values.bootstrapIntent}
              selectedVersionId={values.bootstrapVersionId || ''}
              onAppChange={(appId) =>
                patch({ bootstrapApplicationId: appId, bootstrapVersionId: '' })
              }
              onIntentChange={(newIntent: BootstrapIntent) =>
                patch({
                  bootstrapIntent: newIntent,
                  bootstrapVersionId: newIntent === 'specific' ? values.bootstrapVersionId : '',
                })
              }
              onVersionChange={(versionId) => patch({ bootstrapVersionId: versionId })}
              readOnly={readOnly}
            />

            {showContainerAck ? (
              <div className="flex items-start gap-2.5 bg-amber-50 dark:bg-amber-950/10 p-3 rounded-lg border border-amber-200/50 dark:border-amber-900/20">
                <Checkbox
                  id="container-ack"
                  checked={values.acknowledgeContainerPlacement}
                  disabled={readOnly}
                  onCheckedChange={(c) => patch({ acknowledgeContainerPlacement: c === true })}
                  className="mt-0.5"
                />
                <Label htmlFor="container-ack" className="text-xs text-amber-900 dark:text-amber-300 font-normal leading-relaxed cursor-pointer select-none">
                  {t('enrollmentRoute.form.containerAck')}
                </Label>
              </div>
            ) : null}
          </div>
        </TabsContent>

        {/* --- WI-FI SETUP --- */}
        <TabsContent value="wifi" className="space-y-4 outline-none min-h-[340px]">
          <div className="rounded-xl border bg-card/40 backdrop-blur-sm p-4 shadow-sm space-y-4 animate-in fade-in-50 duration-200">
            <div className="space-y-1">
              <Label className="text-sm font-semibold text-foreground">{t('enrollmentRoute.provisioning.wifiTitle')}</Label>
              <p className="text-xs text-muted-foreground">Configure optional Wi-Fi so devices can connect automatically on boot.</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="wifi-ssid" className="text-xs font-medium">{t('enrollmentRoute.provisioning.wifiSsid')}</Label>
              <Input
                id="wifi-ssid"
                maxLength={32}
                value={values.wifiSsid}
                disabled={readOnly}
                onChange={(e) => patch({ wifiSsid: e.target.value })}
                className="transition-all duration-200 focus:ring-2 focus:ring-primary/20"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="wifi-password" className="text-xs font-medium">{t('enrollmentRoute.provisioning.wifiPassword')}</Label>
              <Input
                id="wifi-password"
                maxLength={63}
                value={values.wifiPassword}
                disabled={readOnly}
                onChange={(e) => patch({ wifiPassword: e.target.value })}
                className="transition-all duration-200 focus:ring-2 focus:ring-primary/20"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="wifi-security" className="text-xs font-medium">{t('enrollmentRoute.provisioning.wifiSecurityType')}</Label>
              <Select
                disabled={readOnly}
                value={values.wifiSecurityType}
                onValueChange={(v) => patch({ wifiSecurityType: v })}
              >
                <SelectTrigger id="wifi-security" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="NONE">NONE</SelectItem>
                  <SelectItem value="WPA">WPA</SelectItem>
                  <SelectItem value="WPA2">WPA2</SelectItem>
                  <SelectItem value="WEP">WEP</SelectItem>
                  <SelectItem value="WPA3">WPA3</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </TabsContent>

        {/* --- ADVANCED SETTINGS --- */}
        <TabsContent value="advanced" className="space-y-4 outline-none min-h-[340px]">
          <div className="rounded-xl border bg-card/40 backdrop-blur-sm p-4 shadow-sm space-y-4 animate-in fade-in-50 duration-200">
            {/* Parameters card */}
            <div className="space-y-3">
              <Label className="text-sm font-semibold">{t('enrollmentRoute.provisioning.advancedTitle')}</Label>
              
              <div className="space-y-2">
                <Label htmlFor="qr-parameters" className="text-xs font-medium">{t('enrollmentRoute.provisioning.qrParameters')}</Label>
                <Textarea
                  id="qr-parameters"
                  rows={3}
                  placeholder={t('enrollmentRoute.provisioning.qrParametersHint')}
                  value={values.qrParameters}
                  disabled={readOnly}
                  onChange={(e) => patch({ qrParameters: e.target.value })}
                  className="font-mono text-xs transition-all duration-200 focus:ring-2 focus:ring-primary/20"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="admin-extras" className="text-xs font-medium">{t('enrollmentRoute.provisioning.adminExtras')}</Label>
                <Textarea
                  id="admin-extras"
                  rows={3}
                  placeholder={t('enrollmentRoute.provisioning.adminExtrasHint')}
                  value={values.adminExtras}
                  disabled={readOnly}
                  onChange={(e) => patch({ adminExtras: e.target.value })}
                  className="font-mono text-xs transition-all duration-200 focus:ring-2 focus:ring-primary/20"
                />
                {values.adminExtras.trim() !== '' && (() => {
                  try { JSON.parse(values.adminExtras); return false } catch { return true }
                })() && (
                  <p className="text-xs text-destructive bg-destructive/10 p-2 rounded border border-destructive/20">{t('enrollmentRoute.provisioning.adminExtrasInvalid')}</p>
                )}
              </div>
            </div>

            {/* Flags */}
            <div className="border-t pt-3 space-y-3">
              <Label className="text-sm font-semibold">{t('enrollmentRoute.provisioning.flagsTitle')}</Label>
              
              <div className="flex items-start gap-2.5 transition-all hover:bg-muted/30 p-2 rounded-lg">
                <Checkbox
                  id="mobile-enrollment"
                  checked={values.mobileEnrollment}
                  disabled={readOnly}
                  onCheckedChange={(c) => patch({ mobileEnrollment: c === true })}
                  className="mt-0.5"
                />
                <Label htmlFor="mobile-enrollment" className="text-xs font-normal leading-relaxed cursor-pointer select-none">
                  {t('enrollmentRoute.provisioning.mobileEnrollment')}
                </Label>
              </div>

              <div className="flex items-start gap-2.5 transition-all hover:bg-muted/30 p-2 rounded-lg">
                <Checkbox
                  id="encrypt-device"
                  checked={values.encryptDevice}
                  disabled={readOnly}
                  onCheckedChange={(c) => patch({ encryptDevice: c === true })}
                  className="mt-0.5"
                />
                <Label htmlFor="encrypt-device" className="text-xs font-normal leading-relaxed cursor-pointer select-none">
                  {t('enrollmentRoute.provisioning.encryptDevice')}
                </Label>
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
