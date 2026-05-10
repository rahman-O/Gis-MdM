import { useEffect, useState } from 'react'
import * as updatesService from '@/features/updates/updatesService'

export function UpdatesPage() {
  const [items, setItems] = useState<updatesService.UpdateEntryRow[]>([])
  const [err, setErr] = useState<string | null>(null)

  useEffect(() => {
    updatesService.checkUpdates().then(setItems).catch((reason: unknown) => {
      setErr(reason instanceof Error ? reason.message : String(reason))
    })
  }, [])

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Updates</h1>
      <p className="text-muted-foreground text-sm">Super-admin or single-tenant check against Headwind manifests.</p>
      {err ? <p className="text-destructive text-sm">{err}</p> : null}
      <ul className="space-y-2 text-sm">
        {items.map((u, i) => (
          <li key={i} className="rounded border p-3">
            <strong>{u.pkg ?? '—'}</strong> v{u.version ?? '?'}{u.outdated ? ' (outdated)' : ''}
            {u.description ? <p className="text-muted-foreground mt-1">{u.description}</p> : null}
          </li>
        ))}
      </ul>
      {items.length === 0 && !err ? <p className="text-muted-foreground text-sm">No manifest entries returned.</p> : null}
    </div>
  )
}
