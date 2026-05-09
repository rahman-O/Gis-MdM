import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as applicationService from '@/features/applications/services/applicationService'
import * as configurationService from '@/features/configurations/configurationService'
import type { ApplicationVersionConfigurationLink } from '@/features/applications/model/types'

interface Props {
  open: boolean
  versionId: number | null
  onClose: () => void
}

export function ApplicationVersionConfigurationsDialog({ open, versionId, onClose }: Props) {
  const [rows, setRows] = useState<ApplicationVersionConfigurationLink[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!open || versionId == null) return
    setLoading(true)
    setError(null)
    void Promise.all([
      applicationService.getApplicationVersionConfigurations(versionId),
      configurationService.listConfigurationNames(),
    ])
      .then(([links, all]) => {
        const map = new Map(links.map((l) => [l.configurationId, l]))
        setRows(all.map((c) => ({
          configurationId: c.id,
          name: c.name,
          action: Number(map.get(c.id)?.action ?? 0),
          selected: map.has(c.id),
          notify: Boolean(map.get(c.id)?.notify),
        })))
      })
      .catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Failed to load links.'))
      .finally(() => setLoading(false))
  }, [open, versionId])

  const save = async () => {
    if (versionId == null) return
    setLoading(true)
    setError(null)
    try {
      await applicationService.updateApplicationVersionConfigurations({
        applicationVersionId: versionId,
        configurations: rows
          .filter((r) => Number(r.action ?? 0) !== 0 || r.selected)
          .map((r) => ({
            configurationId: r.configurationId,
            action: Number(r.action ?? 0),
            notify: Boolean(r.notify),
          })),
      })
      onClose()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to save version links.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="max-h-[80vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader><DialogTitle>Version configurations</DialogTitle></DialogHeader>
        {rows.map((r, i) => (
          <div key={r.configurationId} className="grid grid-cols-[1fr_auto_auto] items-center gap-3 rounded border p-2">
            <label className="flex items-center gap-2 text-sm">
              <Checkbox checked={Boolean(r.selected)} onCheckedChange={(v) => setRows((p) => p.map((x, idx) => idx === i ? { ...x, selected: v === true } : x))} />
              {r.name ?? `Configuration #${r.configurationId}`}
            </label>
            <Select value={String(r.action ?? 0)} onValueChange={(v) => setRows((p) => p.map((x, idx) => idx === i ? { ...x, action: Number(v) } : x))}>
              <SelectTrigger className="w-36"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="0">Hide</SelectItem>
                <SelectItem value="1">Install</SelectItem>
                <SelectItem value="2">Uninstall</SelectItem>
              </SelectContent>
            </Select>
            <label className="flex items-center gap-2 text-xs">
              <Checkbox checked={Boolean(r.notify)} onCheckedChange={(v) => setRows((p) => p.map((x, idx) => idx === i ? { ...x, notify: v === true } : x))} />
              Notify
            </label>
          </div>
        ))}
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cancel</Button>
          <Button onClick={() => void save()} disabled={loading}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
