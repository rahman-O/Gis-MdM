import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { EnrollmentContractPreview } from '@/features/enrollment-routes/buildEnrollmentContractPreview'
import type {
  EnrollmentRouteQrMeta,
  TreeNodeOption,
  BootstrapAppOption,
} from '@/features/enrollment-routes/enrollmentRouteService'
import {
  buildEnrollmentQrImagePath,
  defaultViewportQrSize,
} from '@/features/devices/enrollmentQrQuery'
import { loadQrImageObjectUrl } from '@/features/devices/qrImage'
import { Button } from '@/shared/ui/button'
import { Label } from '@/shared/ui/label'
import { Switch } from '@/shared/ui/switch'

interface PendingProps {
  mode: 'pending'
  preview: EnrollmentContractPreview | null
  treeNodes?: TreeNodeOption[]
  bootstrapApps?: BootstrapAppOption[]
}

interface ActiveProps {
  mode: 'active'
  meta: EnrollmentRouteQrMeta
  dimmed?: boolean
  treeNodes?: TreeNodeOption[]
  bootstrapApps?: BootstrapAppOption[]
}

type Props = PendingProps | ActiveProps

export function EnrollmentRouteQrColumn(props: Props) {
  const { t } = useTranslation()
  const [showNames, setShowNames] = useState(true)

  if (props.mode === 'pending') {
    return (
      <PendingQrColumn
        preview={props.preview}
        showNames={showNames}
        setShowNames={setShowNames}
        treeNodes={props.treeNodes}
        bootstrapApps={props.bootstrapApps}
      />
    )
  }

  return (
    <ActiveQrColumn
      meta={props.meta}
      dimmed={props.dimmed}
      showNames={showNames}
      setShowNames={setShowNames}
      treeNodes={props.treeNodes}
      bootstrapApps={props.bootstrapApps}
      label={t('enrollmentRoute.qr.active')}
    />
  )
}

function PendingQrColumn({
  preview,
  showNames,
  setShowNames,
  treeNodes = [],
  bootstrapApps = [],
}: {
  preview: EnrollmentContractPreview | null
  showNames: boolean
  setShowNames: (val: boolean) => void
  treeNodes?: TreeNodeOption[]
  bootstrapApps?: BootstrapAppOption[]
}) {
  const { t } = useTranslation()

  const resolvedPreview = useMemo(() => {
    if (!preview) return null
    if (!showNames) return preview

    const node = treeNodes.find((n) => n.id === preview.targetNodeId)
    const app = bootstrapApps.find((a) => a.package === preview.mainAppPackage)

    return {
      ...preview,
      targetNodeId: node ? `${node.name} (${preview.targetNodeId})` : preview.targetNodeId,
      mainAppPackage: app ? `${app.name} (${preview.mainAppPackage})` : preview.mainAppPackage,
    } as Record<string, unknown>
  }, [preview, showNames, treeNodes, bootstrapApps])

  const json = resolvedPreview ? JSON.stringify(resolvedPreview, null, 2) : ''

  return (
    <div className="flex h-full flex-col gap-3 rounded-xl border bg-muted/30 p-4 shadow-sm">
      <div className="flex items-center justify-between border-b pb-2 flex-wrap gap-2">
        <Label className="text-sm font-semibold">{t('enrollmentRoute.qr.pending')}</Label>
        {preview && (
          <div className="flex items-center gap-2">
            <Switch
              id="show-names-pending"
              checked={showNames}
              onCheckedChange={setShowNames}
            />
            <Label htmlFor="show-names-pending" className="text-xs text-muted-foreground cursor-pointer select-none">
              {t('enrollmentRoute.qr.showNamesInPreview')}
            </Label>
          </div>
        )}
      </div>

      <div className="relative flex aspect-square max-h-56 w-full items-center justify-center rounded-lg border border-dashed bg-background shadow-inner">
        <span className="px-4 text-center text-xs text-muted-foreground">
          {t('enrollmentRoute.qr.saveToActivate')}
        </span>
        <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
          <span className="rotate-[-18deg] text-xs font-semibold uppercase tracking-widest text-muted-foreground/30">
            {t('enrollmentRoute.qr.pendingWatermark')}
          </span>
        </div>
      </div>
      <pre className="max-h-48 overflow-auto rounded-lg border bg-background p-3 text-xs font-mono text-foreground/80 shadow-inner">
        {json || t('enrollmentRoute.qr.previewEmpty')}
      </pre>
      <p className="text-xs text-muted-foreground">{t('enrollmentRoute.qr.pendingHint')}</p>
    </div>
  )
}

function ActiveQrColumn({
  meta,
  dimmed,
  label,
  showNames,
  setShowNames,
  treeNodes = [],
  bootstrapApps = [],
}: {
  meta: EnrollmentRouteQrMeta
  dimmed?: boolean
  label: string
  showNames: boolean
  setShowNames: (val: boolean) => void
  treeNodes?: TreeNodeOption[]
  bootstrapApps?: BootstrapAppOption[]
}) {
  const { t } = useTranslation()
  const [qrUrl, setQrUrl] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const blobRef = useRef<string | null>(null)

  const fields = useMemo(
    () => ({
      size: defaultViewportQrSize(),
      create: true,
      deviceIdUseMode: (meta.defaultDeviceIdMode === 'serial'
        ? 'serial'
        : meta.defaultDeviceIdMode === 'request'
          ? 'request'
          : 'imei') as 'imei' | 'serial' | 'request',
    }),
    [meta.defaultDeviceIdMode]
  )

  const loadQr = useCallback(async () => {
    if (!meta.qrcodekey) return
    setLoading(true)
    setError(null)
    const path = buildEnrollmentQrImagePath(meta.qrcodekey, fields)
    const { url, error: loadErr } = await loadQrImageObjectUrl(path)
    if (blobRef.current?.startsWith('blob:')) URL.revokeObjectURL(blobRef.current)
    blobRef.current = url
    setQrUrl(url)
    setError(loadErr)
    setLoading(false)
  }, [meta.qrcodekey, fields])

  useEffect(() => {
    void loadQr()
    return () => {
      if (blobRef.current?.startsWith('blob:')) URL.revokeObjectURL(blobRef.current)
      blobRef.current = null
    }
  }, [loadQr])

  const resolvedContract = useMemo(() => {
    if (!meta.contract) return null
    if (!showNames) return meta.contract

    const contract = { ...meta.contract }
    if (typeof contract.targetNodeId === 'number' || typeof contract.targetNodeId === 'string') {
      const nodeId = Number(contract.targetNodeId)
      const node = treeNodes.find((n) => n.id === nodeId)
      if (node) {
        contract.targetNodeId = `${node.name} (${nodeId})`
      }
    }

    if (typeof contract.mainAppPackage === 'string') {
      const app = bootstrapApps.find((a) => a.package === contract.mainAppPackage)
      if (app) {
        contract.mainAppPackage = `${app.name} (${contract.mainAppPackage})`
      }
    }

    return contract
  }, [meta.contract, showNames, treeNodes, bootstrapApps])

  const contractJson = resolvedContract
    ? JSON.stringify(resolvedContract, null, 2)
    : ''

  const copyJson = async () => {
    if (!contractJson) return
    try {
      await navigator.clipboard.writeText(contractJson)
    } catch {
      /* ignore */
    }
  }

  return (
    <div
      className={`flex h-full flex-col gap-3 rounded-xl border p-4 shadow-sm transition-all duration-200 ${
        dimmed ? 'opacity-60' : 'bg-muted/30'
      }`}
    >
      <div className="flex items-center justify-between border-b pb-2 flex-wrap gap-2">
        <Label className="text-sm font-semibold">{label}</Label>
        {meta.contract && (
          <div className="flex items-center gap-2">
            <Switch
              id="show-names-active"
              checked={showNames}
              onCheckedChange={setShowNames}
            />
            <Label htmlFor="show-names-active" className="text-xs text-muted-foreground cursor-pointer select-none">
              {t('enrollmentRoute.qr.showNamesInPreview')}
            </Label>
          </div>
        )}
      </div>

      <div className="flex aspect-square max-h-56 w-full items-center justify-center rounded-lg border bg-background shadow-inner">
        {loading ? (
          <span className="text-sm text-muted-foreground animate-pulse">{t('enrollmentRoute.qr.loading')}</span>
        ) : error ? (
          <span className="px-2 text-center text-sm text-destructive">{error}</span>
        ) : qrUrl ? (
          <img src={qrUrl} alt="" className="h-full w-full object-contain p-2" />
        ) : null}
      </div>
      {contractJson ? (
        <>
          <pre className="max-h-40 overflow-auto rounded-lg border bg-background p-3 text-xs font-mono text-foreground/80 shadow-inner">{contractJson}</pre>
          <Button type="button" variant="outline" size="sm" className="w-full h-9 hover:bg-muted" onClick={() => void copyJson()}>
            {t('enrollmentRoute.qr.copyContract')}
          </Button>
        </>
      ) : null}
      <p className="text-xs text-muted-foreground">{t('enrollmentRoute.qr.activeHint')}</p>
    </div>
  )
}
