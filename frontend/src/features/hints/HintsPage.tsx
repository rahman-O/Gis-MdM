import { useCallback, useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import * as hintsService from '@/features/hints/hintsService'

export function HintsPage() {
  const [keys, setKeys] = useState<string[]>([])
  const [err, setErr] = useState<string | null>(null)

  const load = useCallback(async () => {
    setErr(null)
    try {
      setKeys(await hintsService.fetchHintHistory())
    } catch (reason: unknown) {
      setErr(reason instanceof Error ? reason.message : 'Failed to load.')
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Hints</h1>
      <p className="text-muted-foreground text-sm">Tutorial hints tracked per logged-in account.</p>
      <div className="flex flex-wrap gap-2">
        <Button type="button" variant="outline" onClick={() => hintsService.enableHints().then(load)}>
          Enable hints
        </Button>
        <Button type="button" variant="outline" onClick={() => hintsService.disableHints().then(load)}>
          Disable hints
        </Button>
      </div>
      {err ? <p className="text-destructive text-sm">{err}</p> : null}
      <div>
        <h2 className="text-sm font-medium">Recorded hint keys ({keys.length})</h2>
        <ul className="mt-2 list-inside list-disc text-sm">
          {keys.map((k) => (
            <li key={k}>{k}</li>
          ))}
        </ul>
        {keys.length === 0 && !err ? <p className="text-muted-foreground text-sm">None yet.</p> : null}
      </div>
    </div>
  )
}
