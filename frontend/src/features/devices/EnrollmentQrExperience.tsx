import type { ReactNode } from 'react'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { getConfigurationQrEligibility } from '@/features/configurations/configurationQr'
import type { Configuration } from '@/features/configurations/types'
import {
  buildEnrollmentQrImagePath,
  buildEnrollmentQrJsonPath,
  defaultViewportQrSize,
  type EnrollmentQrFields,
  type QrDeviceIdUseMode,
} from '@/features/devices/enrollmentQrQuery'
import { loadQrImageObjectUrl } from '@/features/devices/qrImage'
import type { LookupItem } from '@/features/devices/types'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import apiClient from '@/services/apiClient'

export interface EnrollmentQrExperienceProps {
  qrCodeKey: string
  initialDeviceId?: string
  configuration?: Pick<Configuration, 'qrCodeKey' | 'mainAppId' | 'eventReceivingComponent'> | null
  groups?: LookupItem[]
  /** Footer actions (e.g. Close in a dialog) */
  footer?: ReactNode
}

export function EnrollmentQrExperience({
  qrCodeKey,
  initialDeviceId = '',
  configuration = null,
  groups = [],
  footer,
}: EnrollmentQrExperienceProps) {
  const [deviceIdInput, setDeviceIdInput] = useState(initialDeviceId)
  const [useMode, setUseMode] = useState<QrDeviceIdUseMode>('request')
  const [createOnDemand, setCreateOnDemand] = useState(false)
  const [groupSelection, setGroupSelection] = useState<Set<number>>(new Set())
  const [pixelSize, setPixelSize] = useState(() => defaultViewportQrSize())
  const [showHelp, setShowHelp] = useState(false)

  const [qrObjectUrl, setQrObjectUrl] = useState<string | null>(null)
  const [qrLoading, setQrLoading] = useState(false)
  const [qrError, setQrError] = useState<string | null>(null)

  const [jsonText, setJsonText] = useState<string | null>(null)
  const [jsonLoading, setJsonLoading] = useState(false)
  const [jsonError, setJsonError] = useState<string | null>(null)
  const [jsonCopied, setJsonCopied] = useState(false)

  const latestQrBlobUrl = useRef<string | null>(null)

  const eligibility = useMemo(() => getConfigurationQrEligibility(configuration), [configuration])

  useEffect(() => {
    return () => {
      const url = latestQrBlobUrl.current
      if (url?.startsWith('blob:')) URL.revokeObjectURL(url)
      latestQrBlobUrl.current = null
    }
  }, [])

  useEffect(() => {
    setDeviceIdInput(initialDeviceId)
  }, [initialDeviceId, qrCodeKey])

  useEffect(() => {
    const onResize = () => setPixelSize(defaultViewportQrSize())
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  }, [])

  const qrFields = useMemo((): EnrollmentQrFields => {
    const trimmed = deviceIdInput.trim()
    return {
      size: pixelSize,
      deviceId: trimmed.length > 0 ? trimmed : undefined,
      create: createOnDemand,
      deviceIdUseMode: trimmed.length > 0 ? 'request' : useMode,
      groupIds: createOnDemand ? [...groupSelection].filter((id) => id > 0) : [],
    }
  }, [deviceIdInput, createOnDemand, useMode, groupSelection, pixelSize])

  useEffect(() => {
    if (!qrCodeKey.trim()) return undefined

    const ac = new AbortController()
    let disposed = false

    const run = async () => {
      setQrLoading(true)
      setQrError(null)
      setQrObjectUrl((prev) => {
        if (prev?.startsWith('blob:')) URL.revokeObjectURL(prev)
        return null
      })
      latestQrBlobUrl.current = null

      const path = buildEnrollmentQrImagePath(qrCodeKey.trim(), qrFields)
      const { url: next, error: loadErr } = await loadQrImageObjectUrl(path, ac.signal)
      if (disposed) {
        if (next?.startsWith('blob:')) URL.revokeObjectURL(next)
        return
      }
      if (next) {
        latestQrBlobUrl.current = next
        setQrObjectUrl(next)
      } else {
        setQrError(
          loadErr ??
            'The server did not return a QR image. Check Main App APK URL, launcher URL override, event receiving component, or configuration QR key.'
        )
      }
      setQrLoading(false)
    }

    void run()
    return () => {
      disposed = true
      ac.abort()
    }
  }, [qrCodeKey, qrFields])

  const fetchJson = useCallback(async () => {
    if (!qrCodeKey.trim()) return
    setJsonLoading(true)
    setJsonError(null)
    setJsonCopied(false)
    try {
      const path = buildEnrollmentQrJsonPath(qrCodeKey.trim(), qrFields)
      const res = await apiClient.get<string>(path, { responseType: 'text', transformResponse: (r) => r })
      const raw = res.data
      setJsonText(typeof raw === 'string' ? raw : raw != null ? JSON.stringify(raw, null, 2) : '')
    } catch (e: unknown) {
      setJsonText(null)
      const msg = e instanceof Error ? e.message : 'Failed to load provisioning JSON.'
      setJsonError(
        `${msg} Ensure the configuration has a Main App with APK URL (or launcher URL), and the Go server BASE_URL is reachable from the device.`
      )
    } finally {
      setJsonLoading(false)
    }
  }, [qrCodeKey, qrFields])

  const copyJson = useCallback(async () => {
    if (!jsonText) return
    try {
      await navigator.clipboard.writeText(jsonText)
      setJsonCopied(true)
      window.setTimeout(() => setJsonCopied(false), 2000)
    } catch {
      setJsonCopied(false)
    }
  }, [jsonText])

  const useIdDisabled = deviceIdInput.trim().length > 0

  const displaySize = Math.min(pixelSize, 720)

  return (
    <div className="space-y-4">
      {!eligibility.eligible && eligibility.reason ? (
        <div className="rounded-md border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm text-amber-950 dark:text-amber-100">
          {eligibility.reason}
        </div>
      ) : null}

      <div className="flex flex-col items-center justify-center rounded-md border bg-muted/20 p-4">
        {qrLoading ? <p className="text-sm text-muted-foreground">Loading QR…</p> : null}
        {qrError ? <p className="text-center text-sm text-destructive">{qrError}</p> : null}
        {qrObjectUrl ? (
          <img
            src={qrObjectUrl}
            alt="Enrollment QR code"
            width={displaySize}
            height={displaySize}
            className="max-h-[80vh] max-w-full rounded border"
          />
        ) : null}
        {!qrLoading && !qrObjectUrl && !qrError ? (
          <p className="text-sm text-muted-foreground">QR preview will appear here.</p>
        ) : null}
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2 sm:col-span-2">
          <Label htmlFor="enroll-device-id">Device number</Label>
          <Input
            id="enroll-device-id"
            maxLength={100}
            placeholder="Optional — clear to let the server assign an id"
            value={deviceIdInput}
            onChange={(e) => setDeviceIdInput(e.target.value)}
          />
          <p className="text-xs text-muted-foreground">
            Clearing the device number enables automatic device number assignment in the QR payload (legacy behavior).
          </p>
        </div>

        <div className="space-y-2">
          <Label>Device id mode (when number is empty)</Label>
          <Select
            value={useMode}
            disabled={useIdDisabled}
            onValueChange={(v) => setUseMode(v as QrDeviceIdUseMode)}
          >
            <SelectTrigger id="enroll-use-id">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="request">User provides value on device</SelectItem>
              <SelectItem value="imei">Use IMEI</SelectItem>
              <SelectItem value="serial">Use serial number</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="flex flex-col gap-3">
          <div className="flex items-center gap-2">
            <Checkbox
              id="enroll-create"
              checked={createOnDemand}
              onCheckedChange={(c) => setCreateOnDemand(c === true)}
            />
            <Label htmlFor="enroll-create" className="font-normal">
              Add to device list if missing ({'create=1'})
            </Label>
          </div>
          {createOnDemand && groups.length > 0 ? (
            <div className="space-y-2 rounded-md border p-3">
              <p className="text-xs font-medium text-muted-foreground">Assign to groups on create</p>
              <div className="max-h-40 space-y-2 overflow-y-auto pr-1">
                {groups.map((g) => {
                  const checked = groupSelection.has(g.id)
                  return (
                    <label key={g.id} className="flex cursor-pointer items-center gap-2 text-sm">
                      <Checkbox
                        checked={checked}
                        onCheckedChange={(c) => {
                          setGroupSelection((prev) => {
                            const next = new Set(prev)
                            if (c === true) next.add(g.id)
                            else next.delete(g.id)
                            return next
                          })
                        }}
                      />
                      <span>{g.name ?? `Group #${g.id}`}</span>
                    </label>
                  )
                })}
              </div>
            </div>
          ) : null}
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <Button type="button" variant="secondary" onClick={() => void fetchJson()} disabled={jsonLoading}>
          {jsonLoading ? 'Loading JSON…' : 'Get provisioning JSON'}
        </Button>
        <Button type="button" variant="outline" onClick={() => void copyJson()} disabled={!jsonText}>
          {jsonCopied ? 'Copied' : 'Copy JSON'}
        </Button>
        <Button type="button" variant="outline" onClick={() => setShowHelp((h) => !h)}>
          {showHelp ? 'Hide help' : 'Help'}
        </Button>
      </div>

      {jsonError ? <p className="text-sm text-destructive">{jsonError}</p> : null}
      {jsonText ? (
        <textarea
          readOnly
          className="min-h-[150px] w-full rounded-md border bg-muted/30 p-3 font-mono text-xs"
          value={jsonText}
        />
      ) : null}

      {showHelp ? (
        <div className="rounded-md border bg-muted/30 p-4 text-sm leading-relaxed">
          <p className="font-medium">Android Enterprise — managed device (QR), Android 7+</p>
          <ol className="mt-2 list-decimal space-y-1 pl-5">
            <li>Factory-reset the device.</li>
            <li>On the welcome screen, tap anywhere on the setup wizard in the same spot 6–7 times to open the QR scanner.</li>
            <li>Accept terms, connect to Wi‑Fi if prompted, then scan this QR code.</li>
          </ol>
          <p className="mt-3">
            The QR may embed a device id so the user does not type it. Enter a device number above and the code
            refreshes automatically, or leave it empty for automatic assignment / IMEI / serial per the mode.
          </p>
        </div>
      ) : null}

      {footer ? <div className="flex flex-wrap justify-end gap-2 pt-2">{footer}</div> : null}
    </div>
  )
}
