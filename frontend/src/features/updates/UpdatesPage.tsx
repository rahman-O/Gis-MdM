import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { useToast } from '@/shared/hooks/use-toast'
import * as updatesService from '@/features/updates/updatesService'

export function UpdatesPage() {
  const { toast } = useToast()
  const [items, setItems] = useState<updatesService.UpdateEntryRow[]>([])
  const [err, setErr] = useState<string | null>(null)
  const [applying, setApplying] = useState(false)

  useEffect(() => {
    updatesService.checkUpdates().then(setItems).catch((reason: unknown) => {
      setErr(reason instanceof Error ? reason.message : String(reason))
    })
  }, [])

  async function applyOutdated() {
    const outdated = items.filter((u) => u.outdated)
    if (!outdated.length) {
      toast({ title: 'No outdated packages to apply.' })
      return
    }
    setApplying(true)
    try {
      await updatesService.applyUpdates(outdated)
      toast({ title: 'Update request submitted' })
      const fresh = await updatesService.checkUpdates()
      setItems(fresh)
    } catch (reason: unknown) {
      toast({
        title: 'Apply failed',
        variant: 'destructive',
        description: reason instanceof Error ? reason.message : undefined,
      })
    } finally {
      setApplying(false)
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Updates</h1>
      <p className="text-muted-foreground text-sm">Super-admin or single-tenant check against Headwind manifests.</p>
      <Button type="button" variant="outline" disabled={applying} onClick={() => void applyOutdated()}>
        Apply outdated updates
      </Button>
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
