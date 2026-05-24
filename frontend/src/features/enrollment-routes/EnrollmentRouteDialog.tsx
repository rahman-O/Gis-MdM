import { useCallback, useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { buildEnrollmentContractPreview } from '@/features/enrollment-routes/buildEnrollmentContractPreview'
import { DeleteRouteConfirm } from '@/features/enrollment-routes/DeleteRouteConfirm'
import { EnrollmentRouteDialogHeader } from '@/features/enrollment-routes/EnrollmentRouteDialogHeader'
import { EnrollmentRouteForm } from '@/features/enrollment-routes/EnrollmentRouteForm'
import { EnrollmentRouteQrColumn } from '@/features/enrollment-routes/EnrollmentRouteQrColumn'
import {
  emptyFormValues,
  formValuesFromRoute,
  isEditDirty,
  usesPendingQr,
  type EnrollmentRouteDialogStateId,
  type EnrollmentRouteFormValues,
} from '@/features/enrollment-routes/enrollmentRouteDialogState'
import { resolveBootstrapVersion } from '@/features/enrollment-routes/resolveBootstrapVersion'
import type {
  BootstrapAppOption,
  EnrollmentRouteQrMeta,
  EnrollmentRouteView,
} from '@/features/enrollment-routes/enrollmentRouteService'
import {
  createEnrollmentRoute,
  deleteEnrollmentRoute,
  getEnrollmentRoute,
  getEnrollmentRouteImpact,
  getEnrollmentRouteQrMeta,
  listBootstrapApps,
  updateEnrollmentRoute,
} from '@/features/enrollment-routes/enrollmentRouteService'
import { Button } from '@/shared/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
} from '@/shared/ui/dialog'
import {
  Sheet,
  SheetContent,
  SheetHeader,
} from '@/shared/ui/sheet'
import { useMediaQuery } from '@/shared/hooks/useMediaQuery'

interface Props {
  open: boolean
  state: EnrollmentRouteDialogStateId
  routeId: number
  onStateChange: (state: EnrollmentRouteDialogStateId, routeId?: number) => void
  onClose: () => void
  onSaved: () => void
}

export function EnrollmentRouteDialog({
  open,
  state,
  routeId,
  onStateChange,
  onClose,
  onSaved,
}: Props) {
  const { t } = useTranslation()
  const isDesktop = useMediaQuery('(min-width: 768px)')
  const [route, setRoute] = useState<EnrollmentRouteView | null>(null)
  const [qrMeta, setQrMeta] = useState<EnrollmentRouteQrMeta | null>(null)
  const [form, setForm] = useState<EnrollmentRouteFormValues>(emptyFormValues())
  const [savedSnapshot, setSavedSnapshot] = useState<EnrollmentRouteFormValues | null>(null)
  const [bootstrapApps, setBootstrapApps] = useState<BootstrapAppOption[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [saveError, setSaveError] = useState<string | null>(null)
  const [deleteConfirmName, setDeleteConfirmName] = useState('')

  const isDeleteFlow =
    state === 'DELETE_STEP1' || state === 'DELETE_STEP2' || state === 'DELETE_CONFIRM_ZERO'

  useEffect(() => {
    void listBootstrapApps().then(setBootstrapApps).catch(() => setBootstrapApps([]))
  }, [])

  const loadRoute = useCallback(async (id: number) => {
    setLoading(true)
    setSaveError(null)
    try {
      const [detail, meta] = await Promise.all([
        getEnrollmentRoute(id),
        getEnrollmentRouteQrMeta(id),
      ])
      setRoute(detail)
      setQrMeta(meta)
      const snap = formValuesFromRoute(detail)
      setForm(snap)
      setSavedSnapshot(snap)
    } catch (e: unknown) {
      setSaveError(e instanceof Error ? e.message : 'Failed to load route.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!open) return
    if (state === 'DIALOG_CREATE') {
      setRoute(null)
      setQrMeta(null)
      setForm(emptyFormValues())
      setSavedSnapshot(null)
      setSaveError(null)
      return
    }
    if (routeId > 0 && (state === 'DIALOG_OVERVIEW' || state === 'DIALOG_EDIT' || isDeleteFlow)) {
      void loadRoute(routeId)
    }
  }, [open, state, routeId, loadRoute, isDeleteFlow])

  const dirty = isEditDirty(form, savedSnapshot)

  const resolvedBootstrap = useMemo(
    () =>
      resolveBootstrapVersion(
        bootstrapApps,
        form.bootstrapApplicationId,
        form.bootstrapIntent,
        form.bootstrapVersionId
      ),
    [bootstrapApps, form]
  )

  const pendingPreview = useMemo(() => {
    if (!usesPendingQr(state) || !resolvedBootstrap) return null
    if (!form.targetNodeId || form.targetNodeId <= 0) return null
    return buildEnrollmentContractPreview({
      routeId: routeId > 0 ? routeId : 0,
      targetNodeId: Number(form.targetNodeId),
      mainAppPackage: resolvedBootstrap.package,
      mainAppVersion: resolvedBootstrap.version,
      mainAppVersionCode: resolvedBootstrap.versionCode,
      deviceIdentityMode: form.deviceIdentityMode,
    })
  }, [state, resolvedBootstrap, form, routeId])

  const validate = (): string | null => {
    if (!form.name.trim()) return t('enrollmentRoute.validation.nameRequired')
    if (!form.targetNodeId || form.targetNodeId <= 0) {
      return t('enrollmentRoute.validation.folderRequired')
    }
    if (!form.bootstrapApplicationId || form.bootstrapApplicationId <= 0) {
      return t('enrollmentRoute.validation.appRequired')
    }
    if (form.bootstrapIntent === 'specific' && (!form.bootstrapVersionId || form.bootstrapVersionId <= 0)) {
      return t('enrollmentRoute.validation.versionRequired')
    }
    return null
  }

  const buildPayload = () => ({
    name: form.name.trim(),
    description: form.description.trim() || null,
    targetNodeId: Number(form.targetNodeId),
    deviceIdentityMode: form.deviceIdentityMode,
    bootstrapIntent: form.bootstrapIntent,
    bootstrapApplicationId: Number(form.bootstrapApplicationId),
    bootstrapVersionId:
      form.bootstrapIntent === 'specific' ? Number(form.bootstrapVersionId) : null,
    acknowledgeContainerPlacement: form.acknowledgeContainerPlacement,
  })

  const handleSave = async () => {
    const err = validate()
    if (err) {
      setSaveError(err)
      return
    }
    setSaving(true)
    setSaveError(null)
    try {
      const payload = buildPayload()
      if (state === 'DIALOG_CREATE') {
        const created = await createEnrollmentRoute(payload)
        onSaved()
        onStateChange('DIALOG_OVERVIEW', created.id)
      } else if (state === 'DIALOG_EDIT' && routeId > 0) {
        const updated = await updateEnrollmentRoute(routeId, payload)
        const snap = formValuesFromRoute(updated)
        setForm(snap)
        setSavedSnapshot(snap)
        setRoute(updated)
        const meta = await getEnrollmentRouteQrMeta(routeId)
        setQrMeta(meta)
        onSaved()
        onStateChange('DIALOG_OVERVIEW', routeId)
      }
    } catch (e: unknown) {
      setSaveError(e instanceof Error ? e.message : 'Save failed.')
    } finally {
      setSaving(false)
    }
  }

  const handleDeleteContinue = () => {
    void getEnrollmentRouteImpact(routeId).then((impact) => {
      const any =
        impact.enrollingNowCount > 0 ||
        impact.historicalEnrolledCount > 0 ||
        impact.activeQrScans7d > 0
      onStateChange(any ? 'DELETE_STEP2' : 'DELETE_CONFIRM_ZERO', routeId)
    })
  }

  const handleDeleteConfirm = async () => {
    if (!route) return
    if (state === 'DELETE_STEP2' && deleteConfirmName.trim() !== route.name.trim()) {
      setSaveError(t('enrollmentRoute.delete.nameMismatch'))
      return
    }
    setSaving(true)
    try {
      await deleteEnrollmentRoute(routeId)
      onSaved()
      onClose()
    } catch (e: unknown) {
      setSaveError(e instanceof Error ? e.message : 'Delete failed.')
    } finally {
      setSaving(false)
    }
  }

  const title =
    state === 'DIALOG_CREATE'
      ? t('enrollmentRoute.dialog.createTitle')
      : route?.name ?? t('enrollmentRoute.dialog.routeTitle')

  const readOnlyForm = state === 'DIALOG_OVERVIEW' || isDeleteFlow

  // --- Shared content rendered in both Dialog and Sheet ---

  const headerContent = (
    <EnrollmentRouteDialogHeader
      title={title}
      state={state}
      routeId={routeId}
      dirty={dirty}
    />
  )

  const bodyContent = loading ? (
    <p className="text-sm text-muted-foreground">{t('enrollmentRoute.dialog.loading')}</p>
  ) : (
    <div className="grid gap-4 grid-cols-1 md:grid-cols-2">
      <div>
        {isDeleteFlow ? (
          <DeleteRouteConfirm
            routeId={routeId}
            routeName={route?.name ?? ''}
            step={state as 'DELETE_STEP1' | 'DELETE_STEP2' | 'DELETE_CONFIRM_ZERO'}
            confirmName={deleteConfirmName}
            onConfirmNameChange={setDeleteConfirmName}
          />
        ) : (
          <EnrollmentRouteForm
            values={form}
            onChange={setForm}
            readOnly={readOnlyForm}
            saveError={saveError}
          />
        )}
      </div>
      <div>
        {usesPendingQr(state) ? (
          <EnrollmentRouteQrColumn mode="pending" preview={pendingPreview} />
        ) : qrMeta ? (
          <EnrollmentRouteQrColumn mode="active" meta={qrMeta} dimmed={isDeleteFlow} />
        ) : null}
      </div>
    </div>
  )

  const hintContent =
    state === 'DIALOG_EDIT' && route?.qrcodekey ? (
      <p className="text-xs text-muted-foreground">{t('enrollmentRoute.qr.lastActiveHint')}</p>
    ) : null

  const footerContent = (
    <>
      {state === 'DIALOG_OVERVIEW' ? (
        <>
          <Button type="button" variant="outline" onClick={onClose}>
            {t('enrollmentRoute.actions.close')}
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={() => onStateChange('DELETE_STEP1', routeId)}
          >
            {t('enrollmentRoute.actions.delete')}
          </Button>
          <Button type="button" onClick={() => onStateChange('DIALOG_EDIT', routeId)}>
            {t('enrollmentRoute.actions.edit')}
          </Button>
        </>
      ) : null}

      {state === 'DIALOG_CREATE' || state === 'DIALOG_EDIT' ? (
        <>
          <Button type="button" variant="outline" onClick={onClose} disabled={saving}>
            {t('enrollmentRoute.actions.cancel')}
          </Button>
          <Button type="button" onClick={() => void handleSave()} disabled={saving}>
            {t('enrollmentRoute.actions.save')}
          </Button>
        </>
      ) : null}

      {state === 'DELETE_STEP1' ? (
        <>
          <Button type="button" variant="outline" onClick={() => onStateChange('DIALOG_OVERVIEW', routeId)}>
            {t('enrollmentRoute.actions.back')}
          </Button>
          <Button type="button" variant="destructive" onClick={handleDeleteContinue}>
            {t('enrollmentRoute.actions.continueDelete')}
          </Button>
        </>
      ) : null}

      {(state === 'DELETE_STEP2' || state === 'DELETE_CONFIRM_ZERO') && (
        <>
          <Button type="button" variant="outline" onClick={() => onStateChange('DELETE_STEP1', routeId)}>
            {t('enrollmentRoute.actions.back')}
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={saving}
            onClick={() => void handleDeleteConfirm()}
          >
            {t('enrollmentRoute.actions.confirmDelete')}
          </Button>
        </>
      )}
    </>
  )

  // --- Desktop: Dialog (≥md) ---
  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={(o) => !o && onClose()}>
        <DialogContent className="max-h-[90vh] max-w-4xl overflow-y-auto">
          <DialogHeader>{headerContent}</DialogHeader>
          {bodyContent}
          {hintContent}
          <DialogFooter className="flex flex-wrap gap-2">{footerContent}</DialogFooter>
        </DialogContent>
      </Dialog>
    )
  }

  // --- Mobile: Sheet (<md) — stacks config then QR vertically ---
  return (
    <Sheet open={open} onOpenChange={(o) => !o && onClose()}>
      <SheetContent side="bottom" className="flex h-[90vh] flex-col overflow-y-auto">
        <SheetHeader>{headerContent}</SheetHeader>
        <div className="flex-1 space-y-4 overflow-y-auto py-4">
          {bodyContent}
          {hintContent}
        </div>
        <div className="flex flex-col-reverse gap-2 border-t pt-4 sm:flex-row sm:justify-end">
          {footerContent}
        </div>
      </SheetContent>
    </Sheet>
  )
}
