import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import type { ApplicationVersion } from '@/features/applications/model/types'

interface Props {
  open: boolean
  initialData: ApplicationVersion | null
  onClose: () => void
  onSave: (payload: ApplicationVersion) => Promise<void>
}

export function ApplicationVersionDialog({ open, initialData, onClose, onSave }: Props) {
  const [form, setForm] = useState<ApplicationVersion>({})
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!open) return
    setForm(initialData ?? { split: false })
    setError(null)
  }, [open, initialData])

  const save = async () => {
    if (!String(form.version ?? '').trim() && !String(form.filePath ?? '').trim()) {
      setError('Version is required unless APK file was uploaded.')
      return
    }
    setSaving(true)
    setError(null)
    try {
      await onSave(form)
      onClose()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to save version.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{form.id ? 'Edit version' : 'Add version'}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-3 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Version</Label>
            <Input value={String(form.version ?? '')} onChange={(e) => setForm((p) => ({ ...p, version: e.target.value }))} />
          </div>
          <div className="space-y-2">
            <Label>Version code</Label>
            <Input value={String(form.versionCode ?? '')} onChange={(e) => setForm((p) => ({ ...p, versionCode: Number(e.target.value) || null }))} />
          </div>
          <div className="space-y-2 md:col-span-2">
            <Label>URL</Label>
            <Input value={String(form.url ?? '')} onChange={(e) => setForm((p) => ({ ...p, url: e.target.value }))} />
          </div>
          <div className="flex items-center gap-2 md:col-span-2">
            <Checkbox checked={Boolean(form.split)} onCheckedChange={(v) => setForm((p) => ({ ...p, split: v === true }))} />
            <Label>Split APK</Label>
          </div>
          {form.split ? (
            <>
              <div className="space-y-2">
                <Label>URL armeabi</Label>
                <Input value={String(form.urlArmeabi ?? '')} onChange={(e) => setForm((p) => ({ ...p, urlArmeabi: e.target.value }))} />
              </div>
              <div className="space-y-2">
                <Label>URL arm64</Label>
                <Input value={String(form.urlArm64 ?? '')} onChange={(e) => setForm((p) => ({ ...p, urlArm64: e.target.value }))} />
              </div>
            </>
          ) : null}
        </div>
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={saving}>Cancel</Button>
          <Button onClick={() => void save()} disabled={saving}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
