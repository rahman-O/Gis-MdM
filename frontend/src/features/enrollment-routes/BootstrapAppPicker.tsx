import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { Check, Package } from 'lucide-react'
import type {
  BootstrapAppOption,
  BootstrapAppVersionOption,
  BootstrapIntent,
} from '@/features/enrollment-routes/enrollmentRouteService'
import { resolveBootstrapVersion } from '@/features/enrollment-routes/resolveBootstrapVersion'
import { Badge } from '@/shared/ui/badge'
import { Label } from '@/shared/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { cn } from '@/shared/utils/cn'

export interface BootstrapAppPickerProps {
  /** Available bootstrap applications loaded from the API */
  apps: BootstrapAppOption[]
  /** Currently selected application ID (controlled) */
  selectedAppId: number | ''
  /** Currently selected bootstrap intent (controlled) */
  intent: BootstrapIntent
  /** Currently selected version ID when intent is "specific" (controlled) */
  selectedVersionId: number | ''
  /** Called when the selected application changes */
  onAppChange: (appId: number) => void
  /** Called when the bootstrap intent changes */
  onIntentChange: (intent: BootstrapIntent) => void
  /** Called when the specific version changes */
  onVersionChange: (versionId: number) => void
  /** Whether the picker is in read-only mode */
  readOnly?: boolean
}

/**
 * BootstrapAppPicker — allows users to:
 * 1. Select a bootstrap application from a list
 * 2. Choose a bootstrap intent via radio buttons (stable / latest / specific)
 * 3. When "specific" intent is chosen, pick a specific version from a dropdown
 * 4. See the resolved version line below the picker
 */
export function BootstrapAppPicker({
  apps,
  selectedAppId,
  intent,
  selectedVersionId,
  onAppChange,
  onIntentChange,
  onVersionChange,
  readOnly = false,
}: BootstrapAppPickerProps) {
  const { t } = useTranslation()

  const selectedApp = useMemo(
    () => apps.find((a) => a.applicationId === selectedAppId) ?? null,
    [apps, selectedAppId]
  )

  const versions: BootstrapAppVersionOption[] = selectedApp?.versions ?? []

  const resolved = useMemo(
    () => resolveBootstrapVersion(apps, selectedAppId, intent, selectedVersionId),
    [apps, selectedAppId, intent, selectedVersionId]
  )

  const handleIntentChange = (newIntent: BootstrapIntent) => {
    if (readOnly) return
    onIntentChange(newIntent)
  }

  return (
    <div className="space-y-4">
      {/* Application selector */}
      <div className="space-y-2">
        <Label>{t('enrollmentRoute.form.bootstrapApp')}</Label>
        <Select
          disabled={readOnly}
          value={selectedAppId ? String(selectedAppId) : ''}
          onValueChange={(v) => onAppChange(Number(v))}
        >
          <SelectTrigger>
            <SelectValue placeholder={t('enrollmentRoute.form.selectApp')} />
          </SelectTrigger>
          <SelectContent>
            {apps.map((app) => {
              const recommended = app.versions.some((v) => v.isRecommended)
              return (
                <SelectItem key={app.applicationId} value={String(app.applicationId)}>
                  <span className="flex items-center gap-2">
                    <Package className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
                    <span>{app.name}</span>
                    <span className="text-muted-foreground">({app.package})</span>
                    {recommended && (
                      <Badge variant="secondary" className="ml-1 text-[10px] px-1.5 py-0">
                        {t('enrollmentRoute.form.intentStable')}
                      </Badge>
                    )}
                  </span>
                </SelectItem>
              )
            })}
          </SelectContent>
        </Select>
      </div>

      {/* Intent radio group */}
      {selectedAppId ? (
        <div className="space-y-2">
          <Label>{t('enrollmentRoute.form.bootstrapIntent')}</Label>
          <div
            role="radiogroup"
            aria-label={t('enrollmentRoute.form.bootstrapIntent')}
            className="space-y-1.5"
          >
            <IntentRadioOption
              value="stable"
              checked={intent === 'stable'}
              disabled={readOnly}
              onChange={() => handleIntentChange('stable')}
              label={t('enrollmentRoute.form.intentStable')}
              description={t('enrollmentRoute.bootstrap.stableDesc', 'Uses the recommended version from the app catalog')}
            />
            <IntentRadioOption
              value="latest"
              checked={intent === 'latest'}
              disabled={readOnly}
              onChange={() => handleIntentChange('latest')}
              label={t('enrollmentRoute.form.intentLatest')}
              description={t('enrollmentRoute.bootstrap.latestDesc', 'Uses the newest available version')}
            />
            <IntentRadioOption
              value="specific"
              checked={intent === 'specific'}
              disabled={readOnly}
              onChange={() => handleIntentChange('specific')}
              label={t('enrollmentRoute.form.intentSpecific')}
              description={t('enrollmentRoute.bootstrap.specificDesc', 'Pin to an exact version you choose')}
            />
          </div>
        </div>
      ) : null}

      {/* Version dropdown (only when intent is "specific") */}
      {intent === 'specific' && selectedAppId ? (
        <div className="space-y-2">
          <Label>{t('enrollmentRoute.form.version')}</Label>
          <Select
            disabled={readOnly}
            value={selectedVersionId ? String(selectedVersionId) : ''}
            onValueChange={(v) => onVersionChange(Number(v))}
          >
            <SelectTrigger>
              <SelectValue placeholder={t('enrollmentRoute.form.selectVersion')} />
            </SelectTrigger>
            <SelectContent>
              {versions.map((ver) => (
                <SelectItem key={ver.versionId} value={String(ver.versionId)}>
                  <span className="flex items-center gap-2">
                    <span>
                      {ver.version} ({ver.versionCode})
                    </span>
                    {ver.isRecommended && (
                      <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                        {t('enrollmentRoute.form.intentStable')}
                      </Badge>
                    )}
                    {ver.isLatest && !ver.isRecommended && (
                      <Badge variant="outline" className="text-[10px] px-1.5 py-0">
                        {t('enrollmentRoute.form.intentLatest')}
                      </Badge>
                    )}
                  </span>
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      ) : null}

      {/* Resolved version line */}
      {selectedAppId && resolved ? (
        <div className="flex items-center gap-2 rounded-md border bg-muted/40 px-3 py-2 text-sm">
          <Check className="h-3.5 w-3.5 shrink-0 text-green-600 dark:text-green-400" />
          <span className="text-muted-foreground">
            {resolved.package}{' '}
            <span className="font-medium text-foreground">
              v{resolved.version}
            </span>{' '}
            ({resolved.versionCode})
          </span>
        </div>
      ) : null}
    </div>
  )
}

/* ─── Internal: Intent radio option ─── */

interface IntentRadioOptionProps {
  value: BootstrapIntent
  checked: boolean
  disabled: boolean
  onChange: () => void
  label: string
  description: string
}

function IntentRadioOption({
  value,
  checked,
  disabled,
  onChange,
  label,
  description,
}: IntentRadioOptionProps) {
  return (
    <label
      className={cn(
        'flex cursor-pointer items-start gap-3 rounded-md border px-3 py-2.5 transition-colors',
        checked && 'border-primary bg-primary/5',
        !checked && 'border-border hover:bg-muted/50',
        disabled && 'cursor-not-allowed opacity-60'
      )}
    >
      <input
        type="radio"
        name="bootstrap-intent"
        value={value}
        checked={checked}
        disabled={disabled}
        onChange={onChange}
        className="mt-0.5 h-4 w-4 shrink-0 accent-primary"
      />
      <div className="min-w-0 flex-1">
        <span className="text-sm font-medium">{label}</span>
        <p className="mt-0.5 text-xs text-muted-foreground">{description}</p>
      </div>
    </label>
  )
}
