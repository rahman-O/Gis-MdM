import { useEffect, useState } from 'react'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import type { DeviceSummaryData } from '@/features/dashboard/summaryService'
import { getDeviceSummary } from '@/features/dashboard/summaryService'

function statCard(label: string, value: string) {
  return (
    <div className="rounded-lg border bg-card p-4 shadow-sm">
      <p className="text-muted-foreground text-xs font-medium uppercase tracking-wide">{label}</p>
      <p className="mt-2 text-3xl font-semibold tabular-nums">{value}</p>
    </div>
  )
}

export function DashboardPage() {
  const [summary, setSummary] = useState<DeviceSummaryData | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)
    void getDeviceSummary()
      .then((data) => {
        if (!cancelled) setSummary(data)
      })
      .catch((reason: unknown) => {
        if (!cancelled) setError(reason instanceof Error ? reason.message : 'Failed to load summary.')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const total = summary?.devicesTotal
  const enrolled = summary?.devicesEnrolled
  const lastMonth = summary?.devicesEnrolledLastMonth

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground text-sm mt-1">
          Welcome to Headwind MDM — quick enrollment stats from <code className="rounded bg-muted px-1">GET /private/summary/devices</code>.
        </p>
      </div>

      {loading ? <p className="text-sm text-muted-foreground">Loading statistics…</p> : null}

      {error ? (
        <div
          className="flex flex-wrap items-center gap-3 rounded-lg border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm"
          role="alert"
        >
          <AlertCircle className="h-4 w-4 shrink-0 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => window.location.reload()}>
            Reload
          </Button>
        </div>
      ) : null}

      {!loading && !error && summary ? (
        <div className="grid gap-4 sm:grid-cols-3">
          {statCard('Devices total', String(total ?? '—'))}
          {statCard('Enrolled devices', String(enrolled ?? '—'))}
          {statCard('Enrolled last 30 days', String(lastMonth ?? '—'))}
        </div>
      ) : null}
    </div>
  )
}
