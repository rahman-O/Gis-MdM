import { useEffect, useId, useMemo, useState } from 'react'
import { Loader2, Upload } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as applicationService from '@/features/applications/services/applicationService'
import * as webUiFilesService from '@/features/applications/services/webUiFilesService'
import type { Application, ApplicationFormValues, ApplicationType } from '@/features/applications/model/types'
import { DuplicatePackageDialog } from '@/features/applications/components/DuplicatePackageDialog'
import { ANDROID_INTENT_SUGGESTIONS } from '@/features/applications/constants/androidIntentSuggestions'

interface Props {
  open: boolean
  initialData: Application | null
  onClose: () => void
  onSaved: () => Promise<void> | void
}

const DEFAULT_VALUES: ApplicationFormValues = {
  type: 'app',
  name: '',
  pkg: '',
  version: '',
  versionCode: null,
  url: '',
  urlArmeabi: '',
  urlArm64: '',
  split: false,
  arch: '',
  showIcon: true,
  useKiosk: false,
  system: false,
  runAfterInstall: false,
  runAtBoot: false,
  skipVersion: false,
  iconText: '',
  iconId: null,
  intent: '',
  filePath: '',
  autoUpdate: false,
}

function toValues(data: Application | null): ApplicationFormValues {
  if (!data) return DEFAULT_VALUES
  return {
    ...DEFAULT_VALUES,
    id: data.id ?? null,
    type: (data.type as ApplicationType) || 'app',
    name: String(data.name ?? ''),
    pkg: String(data.pkg ?? ''),
    version: String(data.version ?? ''),
    versionCode: data.versionCode ?? null,
    url: String(data.url ?? ''),
    urlArmeabi: String(data.urlArmeabi ?? ''),
    urlArm64: String(data.urlArm64 ?? ''),
    split: Boolean(data.split),
    arch: String(data.arch ?? ''),
    showIcon: data.showIcon !== false,
    useKiosk: Boolean(data.useKiosk),
    system: Boolean(data.system),
    runAfterInstall: Boolean(data.runAfterInstall),
    runAtBoot: Boolean(data.runAtBoot),
    skipVersion: Boolean(data.skipVersion),
    iconText: String(data.iconText ?? ''),
    iconId: data.iconId ?? null,
    intent: String(data.intent ?? ''),
    filePath: String(data.filePath ?? ''),
    autoUpdate: Boolean(data.autoUpdate),
  }
}

function toPayload(values: ApplicationFormValues): Application {
  return {
    id: values.id ?? null,
    type: values.type,
    name: values.name.trim(),
    pkg: values.pkg.trim() || null,
    version: values.version.trim() || null,
    versionCode: values.versionCode ?? null,
    url: values.url.trim() || null,
    urlArmeabi: values.urlArmeabi.trim() || null,
    urlArm64: values.urlArm64.trim() || null,
    split: values.split,
    arch: values.arch || null,
    showIcon: values.showIcon,
    useKiosk: values.useKiosk,
    system: values.system,
    runAfterInstall: values.runAfterInstall,
    runAtBoot: values.runAtBoot,
    skipVersion: values.skipVersion,
    iconText: values.iconText.trim() || null,
    iconId: values.iconId ?? null,
    intent: values.intent.trim() || null,
    filePath: values.filePath.trim() || null,
    autoUpdate: values.autoUpdate,
  }
}

function readErrorMessage(reason: unknown): string {
  if (reason instanceof Error && reason.message.trim()) return reason.message
  if (reason && typeof reason === 'object') {
    const rec = reason as Record<string, unknown>
    const response = rec.response as Record<string, unknown> | undefined
    const data = response?.data as Record<string, unknown> | undefined
    const msg = data?.message
    if (typeof msg === 'string' && msg.trim()) return msg
    const status = response?.status
    if (typeof status === 'number') return `Upload request failed (${status}).`
  }
  return 'Failed to upload file.'
}

/** Backend only runs APK analysis when the uploaded filename ends with `apk`. XAPK is still a valid ZIP. */
function fileForApkBackendParse(original: File): File {
  if (!original.name.toLowerCase().endsWith('.xapk')) return original
  const stem = original.name.replace(/\.xapk$/i, '') || 'upload'
  return new File([original], `${stem}.apk`, {
    type: original.type || 'application/vnd.android.package-archive',
    lastModified: original.lastModified,
  })
}

async function readZipSignature(file: File): Promise<boolean> {
  const header = new Uint8Array(await file.slice(0, 4).arrayBuffer())
  if (header.length < 4) return false
  // ZIP local file header: PK\x03\x04, also allow empty archive/header variants.
  const isPk = header[0] === 0x50 && header[1] === 0x4b
  const validTail =
    (header[2] === 0x03 && header[3] === 0x04) ||
    (header[2] === 0x05 && header[3] === 0x06) ||
    (header[2] === 0x07 && header[3] === 0x08)
  return isPk && validTail
}

export function ApplicationFormDialog({ open, initialData, onClose, onSaved }: Props) {
  const intentListId = useId()
  const [values, setValues] = useState<ApplicationFormValues>(DEFAULT_VALUES)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [storageBanner, setStorageBanner] = useState<string | null>(null)
  const [uploading, setUploading] = useState(false)
  const [uploadPercent, setUploadPercent] = useState(0)
  const [uploadInfo, setUploadInfo] = useState<{
    fileName: string
    pkg: string
    version: string
    versionCode: number | null
    arch: string
  } | null>(null)
  const [dupOpen, setDupOpen] = useState(false)
  const [duplicates, setDuplicates] = useState<Application[]>([])

  useEffect(() => {
    if (!open) return
    setValues(toValues(initialData))
    setError(null)
    setUploadPercent(0)
    setUploadInfo(null)
  }, [open, initialData])

  useEffect(() => {
    if (!open) {
      setStorageBanner(null)
      return
    }
    let cancelled = false
    void webUiFilesService
      .getStorageLimit()
      .then((lim) => {
        if (cancelled) return
        if (lim.sizeLimit > 0) {
          const available = Math.max(0, lim.sizeLimit - lim.sizeUsed)
          if (available < 20) {
            setStorageBanner(`Available space: ${available} Mb`)
          } else {
            setStorageBanner(null)
          }
        } else {
          setStorageBanner(null)
        }
      })
      .catch(() => {
        if (!cancelled) setStorageBanner(null)
      })
    return () => {
      cancelled = true
    }
  }, [open])

  const isEdit = useMemo(() => values.id != null, [values.id])

  const setField = <K extends keyof ApplicationFormValues>(key: K, value: ApplicationFormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  const validate = (): string | null => {
    if (!values.name.trim()) return 'Name is required.'
    if (values.type === 'app' && !values.pkg.trim()) return 'Package is required.'
    if ((values.type === 'app' || values.type === 'web') && !values.url.trim() && !values.filePath.trim() && !values.split) {
      return 'URL or APK upload is required.'
    }
    const hasUploadedApk = Boolean(values.filePath.trim())
    const hasDownloadUrl =
      Boolean(values.url.trim()) || Boolean(values.urlArmeabi.trim()) || Boolean(values.urlArm64.trim())
    // DB: applicationVersions.version is NOT NULL. APK upload fills version server-side from the APK.
    if (
      values.type === 'app' &&
      hasDownloadUrl &&
      !hasUploadedApk &&
      !values.version.trim()
    ) {
      return 'Version name is required when you use download URL(s) without uploading an APK.'
    }
    if (values.showIcon && values.iconText.trim().length > 256) return 'Icon text exceeds 256 chars.'
    return null
  }

  const doSave = async (forceCreate = false, attachToAppId?: number) => {
    const validationError = validate()
    if (validationError) {
      setError(validationError)
      return
    }
    setSaving(true)
    setError(null)
    try {
      const payload = toPayload(values)
      if (attachToAppId != null) {
        await applicationService.createOrUpdateApplicationVersion({
          applicationId: attachToAppId,
          version: payload.version,
          versionCode: payload.versionCode,
          url: payload.url,
          urlArmeabi: payload.urlArmeabi,
          urlArm64: payload.urlArm64,
          split: payload.split,
          arch: payload.arch,
          filePath: payload.filePath,
          autoUpdate: payload.autoUpdate,
        })
      } else if (payload.type === 'app') {
        if (!isEdit && !forceCreate) {
          const list = await applicationService.validateApplicationPkg({
            id: payload.id ?? null,
            name: payload.name,
            pkg: payload.pkg,
          })
          if (Array.isArray(list) && list.length > 0) {
            setDuplicates(list)
            setDupOpen(true)
            setSaving(false)
            return
          }
        }
        await applicationService.createOrUpdateAndroidApplication(payload)
      } else {
        await applicationService.createOrUpdateWebApplication(payload)
      }
      await onSaved()
      onClose()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to save application.')
    } finally {
      setSaving(false)
    }
  }

  const onFilePicked = async (file: File | null) => {
    if (!file) return
    if (file.size <= 0) {
      setError('Selected APK/XAPK has size 0 bytes. Please choose a valid file.')
      return
    }
    const lower = file.name.toLowerCase()
    if (!lower.endsWith('.apk') && !lower.endsWith('.xapk')) {
      setError('Only APK/XAPK files are supported for application upload.')
      return
    }
    const looksLikeZip = await readZipSignature(file)
    if (!looksLikeZip) {
      setError(
        `Selected file does not look like a valid APK/XAPK archive (size: ${file.size} bytes).`
      )
      return
    }
    setUploading(true)
    setUploadPercent(0)
    setError(null)
    setUploadInfo(null)
    try {
      const uploadFile = fileForApkBackendParse(file)
      const result = await webUiFilesService.uploadApplicationFileWithProgress(uploadFile, setUploadPercent)
      const details = result.fileDetails
      const existing = result.application
      const nextPkg = String(details?.pkg ?? values.pkg)
      const nextVersion = String(details?.version ?? values.version)
      const nextVersionCode = details?.versionCode ?? values.versionCode
      const nextArch = String(details?.arch ?? values.arch)
      const inferredName =
        String(existing?.name ?? '').trim() ||
        (nextPkg ? nextPkg.split('.').filter(Boolean).slice(-1)[0] || nextPkg : '') ||
        file.name.replace(/\.(apk|xapk)$/i, '')
      setValues((prev) => ({
        ...prev,
        name: inferredName,
        filePath: String(result.serverPath ?? ''),
        pkg: nextPkg,
        version: nextVersion,
        versionCode: nextVersionCode,
        arch: nextArch,
      }))
      setUploadInfo({
        fileName: String(result.name ?? uploadFile.name),
        pkg: nextPkg,
        version: nextVersion,
        versionCode: nextVersionCode ?? null,
        arch: nextArch,
      })
      setUploadPercent(100)
    } catch (reason: unknown) {
      const message = readErrorMessage(reason)
      if (message.toLowerCase().includes('zip file is empty')) {
        setError('APK appears empty or corrupted. Please verify the file and try again.')
      } else if (message.includes('form.application.version.code.exists')) {
        setError('This APK versionCode already exists with a different version name.')
      } else if (message.includes('error.size.limit.exceeded')) {
        setError('Storage limit exceeded for this customer.')
      } else if (message.toLowerCase().includes('permission')) {
        setError('Permission denied while uploading file. Check your account permissions.')
      } else if (message.toLowerCase().includes('request failed with status code 500')) {
        setError('Upload failed on server while parsing APK. Please verify the APK file integrity.')
      } else {
        setError(message)
      }
    } finally {
      setUploading(false)
    }
  }

  return (
    <>
      <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>{isEdit ? 'Edit application' : 'Add application'}</DialogTitle>
            <DialogDescription>Configure Android, web, or intent application settings.</DialogDescription>
          </DialogHeader>
          {storageBanner ? (
            <div className="rounded border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm text-amber-950 dark:text-amber-100">
              {storageBanner}
            </div>
          ) : null}
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Type</Label>
              <Select value={values.type} onValueChange={(v) => setField('type', v as ApplicationType)} disabled={isEdit}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="app">Android</SelectItem>
                  <SelectItem value="web">Web</SelectItem>
                  <SelectItem value="intent">Intent</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Name</Label>
              <Input
                value={values.name}
                onChange={(e) => setField('name', e.target.value)}
                disabled={uploading || values.type === 'app'}
              />
            </div>
            {values.type === 'app' ? (
              <>
                <div className="space-y-2">
                  <Label>Package</Label>
                  <Input
                    value={values.pkg}
                    onChange={(e) => setField('pkg', e.target.value)}
                    disabled={uploading || !!uploadInfo}
                  />
                </div>
                <div className="space-y-2">
                  <Label>Version</Label>
                  <Input
                    value={values.version}
                    onChange={(e) => setField('version', e.target.value)}
                    disabled={uploading || !!uploadInfo}
                  />
                </div>
              </>
            ) : null}
            <div className="space-y-2 md:col-span-2">
              <Label>URL</Label>
              <Input value={values.url} onChange={(e) => setField('url', e.target.value)} />
            </div>
            {values.type === 'intent' ? (
              <div className="space-y-2 md:col-span-2">
                <Label>Intent</Label>
                <Input
                  value={values.intent}
                  list={intentListId}
                  onChange={(e) => setField('intent', e.target.value)}
                  placeholder="Choose an action or type a custom intent"
                />
                <datalist id={intentListId}>
                  {ANDROID_INTENT_SUGGESTIONS.map((s) => (
                    <option key={s} value={s} />
                  ))}
                </datalist>
              </div>
            ) : null}
            {values.type === 'app' ? (
              <>
                <div className="space-y-2">
                  <Label>Architecture</Label>
                  <Select
                    value={values.arch || 'none'}
                    onValueChange={(v) => setField('arch', v === 'none' ? '' : v)}
                    disabled={uploading || !!uploadInfo}
                  >
                    <SelectTrigger><SelectValue placeholder="Any" /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">Any</SelectItem>
                      <SelectItem value="armeabi">armeabi</SelectItem>
                      <SelectItem value="arm64">arm64</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Upload APK/XAPK</Label>
                  <Input
                    type="file"
                    accept=".apk,.xapk"
                    onChange={(e) => void onFilePicked(e.target.files?.[0] ?? null)}
                  />
                  {uploading ? (
                    <div className="space-y-1">
                      <p className="text-xs text-muted-foreground">
                        Uploading and analyzing APK... {uploadPercent}%
                      </p>
                      <div className="h-2 w-full overflow-hidden rounded bg-muted">
                        <div className="h-full bg-primary transition-all" style={{ width: `${uploadPercent}%` }} />
                      </div>
                    </div>
                  ) : null}
                  {uploadInfo ? (
                    <div className="rounded border bg-muted/30 p-2 text-xs">
                      <p className="font-medium">Analyzed successfully: {uploadInfo.fileName}</p>
                      <p>Package: {uploadInfo.pkg || '—'}</p>
                      <p>Version: {uploadInfo.version || '—'}</p>
                      <p>Version code: {uploadInfo.versionCode ?? '—'}</p>
                      <p>Arch: {uploadInfo.arch || '—'}</p>
                      <p className="text-emerald-700">
                        Name/package/version/architecture were auto-filled from APK analysis.
                      </p>
                    </div>
                  ) : null}
                </div>
              </>
            ) : null}
            <div className="flex flex-wrap items-center gap-4 md:col-span-2">
              <label className="flex items-center gap-2 text-sm"><Checkbox checked={values.showIcon} onCheckedChange={(v) => setField('showIcon', v === true)} />Show icon</label>
              <label className="flex items-center gap-2 text-sm"><Checkbox checked={values.system} onCheckedChange={(v) => setField('system', v === true)} />System app</label>
              <label className="flex items-center gap-2 text-sm"><Checkbox checked={values.runAfterInstall} onCheckedChange={(v) => setField('runAfterInstall', v === true)} />Run after install</label>
              <label className="flex items-center gap-2 text-sm"><Checkbox checked={values.runAtBoot} onCheckedChange={(v) => setField('runAtBoot', v === true)} />Run at boot</label>
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label>Icon text</Label>
              <Input value={values.iconText} onChange={(e) => setField('iconText', e.target.value)} />
            </div>
          </div>
          {error ? <p className="text-sm text-destructive">{error}</p> : null}
          <DialogFooter>
            <Button variant="outline" onClick={onClose} disabled={saving}>Cancel</Button>
            <Button onClick={() => void doSave()} disabled={saving || uploading}>
              {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Upload className="mr-2 h-4 w-4" />}
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <DuplicatePackageDialog
        open={dupOpen}
        duplicates={duplicates}
        onClose={() => setDupOpen(false)}
        onCreateNew={() => void doSave(true)}
        onAttachVersion={(appId) => void doSave(false, appId)}
      />
    </>
  )
}
