import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { EnrollmentContractPreview } from '@/features/enrollment-routes/buildEnrollmentContractPreview'
import type { EnrollmentRouteQrMeta } from '@/features/enrollment-routes/enrollmentRouteService'
import {
  buildEnrollmentQrImagePath,
  defaultViewportQrSize,
} from '@/features/devices/enrollmentQrQuery'
import { loadQrImageObjectUrl } from '@/features/devices/qrImage'
import { Button } from '@/shared/ui/button'
import { Label } from '@/shared/ui/label'

interface PendingProps {
  mode: 'pending'
  preview: EnrollmentContractPreview | null
}

interface ActiveProps {
  mode: 'active'
  meta: EnrollmentRouteQrMeta
  dimmed?: boolean
}

type Props = PendingProps | ActiveProps

export function EnrollmentRouteQrColumn(props: Props) {
  const { t } = useTranslation()

  if (props.mode === 'pending') {
    return <PendingQrColumn preview={props.preview} />
  }

  return <ActiveQrColumn meta={props.meta} dimmed={props.dimmed} label={t('enrollmentRoute.qr.active')} />
}

function PendingQrColumn({ preview }: { preview: EnrollmentContractPreview | null }) {
  const { t } = useTranslation()
  const json = preview ? JSON.stringify(preview, null, 2) : ''

  return (
    <div className="flex h-full flex-col gap-3 rounded-md border bg-muted/30 p-3">
      <Label className="text-sm font-medium">{t('enrollmentRoute.qr.pending')}</Label>
      <div className="relative flex aspect-square max-h-56 w-full items-center justify-center rounded-md border border-dashed bg-background">
        <span className="px-4 text-center text-xs text-muted-foreground">
          {t('enrollmentRoute.qr.saveToActivate')}
        </span>
        <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
          <span className="rotate-[-18deg] text-xs font-semibold uppercase tracking-wide text-muted-foreground/70">
            {t('enrollmentRoute.qr.pendingWatermark')}
          </span>
        </div>
      </div>
      <pre className="max-h-48 overflow-auto rounded-md bg-background p-2 text-xs">
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
}: {
  meta: EnrollmentRouteQrMeta
  dimmed?: boolean
  label: string
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

  const contractJson = meta.contract
    ? JSON.stringify(meta.contract, null, 2)
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
      className={`flex h-full flex-col gap-3 rounded-md border p-3 ${dimmed ? 'opacity-60' : 'bg-muted/30'}`}
    >
      <Label className="text-sm font-medium">{label}</Label>
      <div className="flex aspect-square max-h-56 w-full items-center justify-center rounded-md border bg-background">
        {loading ? (
          <span className="text-sm text-muted-foreground">{t('enrollmentRoute.qr.loading')}</span>
        ) : error ? (
          <span className="px-2 text-center text-sm text-destructive">{error}</span>
        ) : qrUrl ? (
          <img src={qrUrl} alt="" className="h-full w-full object-contain" />
        ) : null}
      </div>
      {contractJson ? (
        <>
          <pre className="max-h-40 overflow-auto rounded-md bg-background p-2 text-xs">{contractJson}</pre>
          <Button type="button" variant="outline" size="sm" onClick={() => void copyJson()}>
            {t('enrollmentRoute.qr.copyContract')}
          </Button>
        </>
      ) : null}
      <p className="text-xs text-muted-foreground">{t('enrollmentRoute.qr.activeHint')}</p>
    </div>
  )
}
