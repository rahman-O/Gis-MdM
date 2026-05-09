import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as applicationService from '@/features/applications/services/applicationService'
import * as configurationService from '@/features/configurations/configurationService'
import type { ApplicationConfigurationLink } from '@/features/applications/model/types'

interface Props {
  open: boolean
  applicationId: number | null
  onClose: () => void
}

export function ApplicationConfigurationsDialog({ open, applicationId, onClose }: Props) {
  const [rows, setRows] = useState<ApplicationConfigurationLink[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!open || applicationId == null) return
    setLoading(true)
    setError(null)
    void Promise.all([
      applicationService.getApplicationConfigurations(applicationId),
      configurationService.listConfigurationNames(),
    ])
      .then(([links, all]) => {
        const map = new Map(links.map((l) => [l.configurationId, l]))
        const combined = all.map((c) => {
          const current = map.get(c.id)
          return {
            configurationId: c.id,
            name: c.name,
            selected: current != null,
            action: Number(current?.action ?? 0),
            notify: Boolean(current?.notify),
          }
        })
        setRows(combined)
      })
      .catch((reason: unknown) => setError(reason instanceof Error ? reason.message : 'Failed to load configurations.'))
      .finally(() => setLoading(false))
  }, [open, applicationId])

  const save = async () => {
    if (applicationId == null) return
    setLoading(true)
    setError(null)
    try {
      await applicationService.updateApplicationConfigurations({
        applicationId,
        configurations: rows.map((r) => ({
          configurationId: r.configurationId,
          action: Number(r.action ?? 0),
          notify: Boolean(r.notify),
        })),
      })
      onClose()
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to update links.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="max-h-[80vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader><DialogTitle>Application configurations</DialogTitle></DialogHeader>
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
        {loading ? <p className="text-sm text-muted-foreground">Saving...</p> : null}
        {error ? <p className="text-sm text-destructive">{error}</p> : null}
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cancel</Button>
          <Button onClick={() => void save()} disabled={loading}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
