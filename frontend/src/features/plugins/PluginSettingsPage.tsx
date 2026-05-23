import { useEffect, useMemo, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import * as pluginService from '@/features/plugins/pluginService'

export function PluginSettingsPage() {
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [plugins, setPlugins] = useState<pluginService.PluginRow[]>([])
  const [selection, setSelection] = useState<Record<number, boolean>>({})

  useEffect(() => {
    let cancelled = false
    async function run() {
      setLoading(true)
      try {
        const active = await pluginService.fetchActivePlugins()
        const avail = await pluginService.fetchAvailablePlugins()
        if (cancelled) return
        active.sort((a, b) => String(a.nameLocalizationKey ?? a.identifier ?? '').localeCompare(String(b.nameLocalizationKey ?? b.identifier ?? '')))
        setPlugins(active)
        const init: Record<number, boolean> = {}
        for (const p of active) init[p.id] = false
        for (const p of avail) init[p.id] = true
        setSelection(init)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    void run()
    return () => {
      cancelled = true
    }
  }, [])

  const toggle = (id: number, checked: boolean) => {
    setSelection((prev) => ({ ...prev, [id]: checked }))
  }

  const disabledPayload = useMemo(() => {
    const ids: number[] = []
    for (const p of plugins) {
      if (selection[p.id] === false) ids.push(p.id)
    }
    return ids
  }, [plugins, selection])

  async function save() {
    setBusy(true)
    try {
      await pluginService.saveDisabledPlugins(disabledPayload)
    } finally {
      setBusy(false)
    }
  }

  if (loading) {
    return <p className="text-muted-foreground text-sm">Loading plugins…</p>
  }

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Plugins</h1>
      <p className="text-muted-foreground max-w-xl text-sm">
        Checked plugins stay enabled for your tenant; unchecked IDs are posted as disabled (same semantics as legacy Angular tab).
      </p>
      <div className="max-w-xl space-y-3">
        {plugins.map((p) => (
          <label key={p.id} className="flex items-center gap-2 rounded border px-3 py-2">
            <Checkbox checked={selection[p.id] === true} onCheckedChange={(c) => toggle(p.id, c === true)} />
            <span className="text-sm">
              {p.identifier?.trim() || p.nameLocalizationKey || `Plugin #${p.id}`}
            </span>
          </label>
        ))}
      </div>
      <Button type="button" onClick={() => void save()} disabled={busy}>
        {busy ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
        Save
      </Button>
    </div>
  )
}
