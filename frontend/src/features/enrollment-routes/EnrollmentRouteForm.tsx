import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AlertCircle, ChevronRight, Folder } from 'lucide-react'
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
}

export function EnrollmentRouteForm({ values, onChange, readOnly = false, saveError }: Props) {
  const { t } = useTranslation()
  const [treeNodes, setTreeNodes] = useState<TreeNodeOption[]>([])
  const [bootstrapApps, setBootstrapApps] = useState<BootstrapAppOption[]>([])
  const [optionsError, setOptionsError] = useState<string | null>(null)
  const [pickerOpen, setPickerOpen] = useState(false)

  useEffect(() => {
    void Promise.all([listTreeNodeOptions(), listBootstrapApps()])
      .then(([nodes, apps]) => {
        setTreeNodes(nodes)
        setBootstrapApps(apps)
      })
      .catch((e: unknown) => {
        setOptionsError(e instanceof Error ? e.message : 'Failed to load options.')
      })
  }, [])

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

  /** Render the selected node path as a breadcrumb */
  const renderSelectedPath = () => {
    if (!selectedNode) return null
    const segments = selectedNode.path.split('/').filter(Boolean)
    return (
      <div className="flex items-center gap-1 text-sm text-muted-foreground">
        <Folder className="h-3.5 w-3.5 shrink-0" />
        {segments.map((segment, idx) => (
          <span key={idx} className="flex items-center gap-1">
            {idx > 0 && <ChevronRight className="h-3 w-3 shrink-0 text-muted-foreground/60" />}
            <span className={idx === segments.length - 1 ? 'font-medium text-foreground' : ''}>
              {segment}
            </span>
          </span>
        ))}
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {optionsError ? (
        <p className="text-sm text-destructive">{optionsError}</p>
      ) : null}
      {saveError ? (
        <div className="flex items-start gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-2 text-sm text-destructive">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>{saveError}</span>
        </div>
      ) : null}

      <div className="space-y-2">
        <Label htmlFor="route-name">{t('enrollmentRoute.form.name')}</Label>
        <Input
          id="route-name"
          value={values.name}
          disabled={readOnly}
          onChange={(e) => patch({ name: e.target.value })}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="route-desc">{t('enrollmentRoute.form.description')}</Label>
        <Input
          id="route-desc"
          value={values.description}
          disabled={readOnly}
          onChange={(e) => patch({ description: e.target.value })}
        />
      </div>

      {/* Target folder section with TargetNodePicker */}
      <div className="space-y-2">
        <Label>{t('enrollmentRoute.form.targetFolder')}</Label>

        {/* Selected path display + trigger button */}
        {!pickerOpen && (
          <div className="space-y-2">
            {selectedNode ? (
              <div className="flex items-center justify-between rounded-md border px-3 py-2">
                {renderSelectedPath()}
                {!readOnly && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="ml-2 shrink-0"
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
                className="w-full justify-start text-muted-foreground"
                disabled={readOnly}
                onClick={() => setPickerOpen(true)}
              >
                <Folder className="mr-2 h-4 w-4" />
                {t('enrollmentRoute.form.selectFolder')}
              </Button>
            )}
            {selectedNode?.placementKind === 'inheritable' ? (
              <p className="text-xs text-amber-700 dark:text-amber-400">
                {t('enrollmentRoute.form.containerWarning')}
              </p>
            ) : null}
          </div>
        )}

        {/* Inline TargetNodePicker panel */}
        <TargetNodePicker
          selectedNodeId={values.targetNodeId}
          onSelect={handleNodeSelect}
          onCancel={handlePickerCancel}
          open={pickerOpen}
        />
      </div>

      <div className="space-y-2">
        <Label>{t('enrollmentRoute.form.deviceIdentity')}</Label>
        <Select
          disabled={readOnly}
          value={values.deviceIdentityMode}
          onValueChange={(v) => patch({ deviceIdentityMode: v })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="imei">IMEI</SelectItem>
            <SelectItem value="serial">{t('enrollmentRoute.form.serial')}</SelectItem>
            <SelectItem value="request">{t('enrollmentRoute.form.request')}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Bootstrap app + intent + version picker */}
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
        <div className="flex items-start gap-2">
          <Checkbox
            id="container-ack"
            checked={values.acknowledgeContainerPlacement}
            disabled={readOnly}
            onCheckedChange={(c) => patch({ acknowledgeContainerPlacement: c === true })}
          />
          <Label htmlFor="container-ack" className="text-sm font-normal leading-snug">
            {t('enrollmentRoute.form.containerAck')}
          </Label>
        </div>
      ) : null}

      {/* Wi-Fi Provisioning card */}
      <div className="rounded-md border p-4 space-y-3">
        <Label className="text-sm font-semibold">{t('enrollmentRoute.provisioning.wifiTitle')}</Label>
        <div className="space-y-2">
          <Label htmlFor="wifi-ssid">{t('enrollmentRoute.provisioning.wifiSsid')}</Label>
          <Input
            id="wifi-ssid"
            maxLength={32}
            value={values.wifiSsid}
            disabled={readOnly}
            onChange={(e) => patch({ wifiSsid: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="wifi-password">{t('enrollmentRoute.provisioning.wifiPassword')}</Label>
          <Input
            id="wifi-password"
            maxLength={63}
            value={values.wifiPassword}
            disabled={readOnly}
            onChange={(e) => patch({ wifiPassword: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="wifi-security">{t('enrollmentRoute.provisioning.wifiSecurityType')}</Label>
          <Select
            disabled={readOnly}
            value={values.wifiSecurityType}
            onValueChange={(v) => patch({ wifiSecurityType: v })}
          >
            <SelectTrigger id="wifi-security">
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

      {/* Advanced QR Parameters card */}
      <div className="rounded-md border p-4 space-y-3">
        <Label className="text-sm font-semibold">{t('enrollmentRoute.provisioning.advancedTitle')}</Label>
        <div className="space-y-2">
          <Label htmlFor="qr-parameters">{t('enrollmentRoute.provisioning.qrParameters')}</Label>
          <Textarea
            id="qr-parameters"
            rows={3}
            placeholder={t('enrollmentRoute.provisioning.qrParametersHint')}
            value={values.qrParameters}
            disabled={readOnly}
            onChange={(e) => patch({ qrParameters: e.target.value })}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="admin-extras">{t('enrollmentRoute.provisioning.adminExtras')}</Label>
          <Textarea
            id="admin-extras"
            rows={3}
            placeholder={t('enrollmentRoute.provisioning.adminExtrasHint')}
            value={values.adminExtras}
            disabled={readOnly}
            onChange={(e) => patch({ adminExtras: e.target.value })}
          />
          {values.adminExtras.trim() !== '' && (() => {
            try { JSON.parse(values.adminExtras); return false } catch { return true }
          })() && (
            <p className="text-sm text-destructive">{t('enrollmentRoute.provisioning.adminExtrasInvalid')}</p>
          )}
        </div>
      </div>

      {/* Enrollment Flags card */}
      <div className="rounded-md border p-4 space-y-3">
        <Label className="text-sm font-semibold">{t('enrollmentRoute.provisioning.flagsTitle')}</Label>
        <div className="flex items-start gap-2">
          <Checkbox
            id="mobile-enrollment"
            checked={values.mobileEnrollment}
            disabled={readOnly}
            onCheckedChange={(c) => patch({ mobileEnrollment: c === true })}
          />
          <Label htmlFor="mobile-enrollment" className="text-sm font-normal leading-snug">
            {t('enrollmentRoute.provisioning.mobileEnrollment')}
          </Label>
        </div>
        <div className="flex items-start gap-2">
          <Checkbox
            id="encrypt-device"
            checked={values.encryptDevice}
            disabled={readOnly}
            onCheckedChange={(c) => patch({ encryptDevice: c === true })}
          />
          <Label htmlFor="encrypt-device" className="text-sm font-normal leading-snug">
            {t('enrollmentRoute.provisioning.encryptDevice')}
          </Label>
        </div>
      </div>
    </div>
  )
}
