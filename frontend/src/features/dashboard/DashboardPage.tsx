import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { AlertCircle } from 'lucide-react'
import {
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Tooltip,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
} from 'recharts'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import type { DeviceSummaryPayload } from '@/features/dashboard/types'
import {
  getConfigurationApplicationCounts,
  getRecentDevices,
  getSummaryDevices,
  parseInstallSummary,
  parseStatusSummary,
} from '@/features/dashboard/dashboardService'
import type { DeviceView } from '@/features/devices/types'
import { StatusBadge } from '@/features/devices/StatusBadge'
import { formatLastSeen } from '@/features/devices/deviceFormat'

const COL_STATUS = ['#F7464A', '#FDB45C', '#46BFBD']
const COL_INSTALL = ['#F7464A', '#FDB45C', '#46BFBD']
const MONTH_COLOR = '#97BBCD'

function statSkeleton() {
  return <Skeleton className="h-8 w-24" />
}

export function DashboardPage() {
  const [summary, setSummary] = useState<DeviceSummaryPayload | null>(null)
  const [counts, setCounts] = useState<{ configurationCount: number; applicationCount: number } | null>(null)
  const [recent, setRecent] = useState<DeviceView[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  const initialLoadDone = useRef(false)

  const doInitialLoad = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [sum, rc, ctr] = await Promise.all([
        getSummaryDevices(),
        getRecentDevices(5),
        getConfigurationApplicationCounts(),
      ])
      setSummary(sum)
      setRecent(rc)
      setCounts(ctr)
      initialLoadDone.current = true
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load dashboard.')
    } finally {
      setLoading(false)
    }
  }, [])

  const refreshQuiet = useCallback(async () => {
    try {
      const [sum, rc, ctr] = await Promise.all([
        getSummaryDevices(),
        getRecentDevices(5),
        getConfigurationApplicationCounts(),
      ])
      setSummary(sum)
      setRecent(rc)
      setCounts(ctr)
    } catch {
      /* keep previous data — spec: ignore background failures */
    }
  }, [])

  useEffect(() => {
    void doInitialLoad()
  }, [doInitialLoad])

  useEffect(() => {
    const id = window.setInterval(() => void refreshQuiet(), 60_000)
    return () => window.clearInterval(id)
  }, [refreshQuiet])

  const statusPie = useMemo(() => {
    const s = parseStatusSummary(summary?.statusSummary ?? undefined)
    return [
      { name: 'Offline', value: s.offline },
      { name: 'Idle', value: s.idle },
      { name: 'Online', value: s.online },
    ].filter((d) => d.value > 0)
  }, [summary?.statusSummary])

  const installPie = useMemo(() => {
    const s = parseInstallSummary(summary?.installSummary ?? undefined)
    return [
      { name: 'Install failed', value: s.failure },
      { name: 'Version mismatch', value: s.mismatch },
      { name: 'Installed OK', value: s.success },
    ].filter((d) => d.value > 0)
  }, [summary?.installSummary])

  const monthlyBar = useMemo(() => {
    const rows =
      summary?.devicesEnrolledMonthly?.map((item) => ({
        label: item.stringAttr ?? '—',
        enrollments: typeof item.number === 'number' ? item.number : 0,
      })) ?? []
    return rows
  }, [summary?.devicesEnrolledMonthly])

  const enrollmentSplit = useMemo(() => {
    const enrolled = Number(summary?.devicesEnrolled ?? 0)
    const lastMo = Number(summary?.devicesEnrolledLastMonth ?? 0)
    const earlier = Math.max(0, enrolled - lastMo)
    return [
      { name: 'Earlier', value: earlier, fill: '#DCDCDC' },
      { name: 'Last 30 days', value: lastMo, fill: '#97BBCD' },
    ].filter((d) => d.value >= 0)
  }, [summary?.devicesEnrolled, summary?.devicesEnrolledLastMonth])

  const byConfigCombined = useMemo(() => {
    const labels = summary?.topConfigs ?? []
    const off = summary?.statusOfflineByConfig ?? []
    const idle = summary?.statusIdleByConfig ?? []
    const on = summary?.statusOnlineByConfig ?? []
    return labels.map((name, i) => ({
      name: name ?? `Cfg ${i + 1}`,
      offline: Number(off[i] ?? 0),
      idle: Number(idle[i] ?? 0),
      online: Number(on[i] ?? 0),
    }))
  }, [summary?.topConfigs, summary?.statusOfflineByConfig, summary?.statusIdleByConfig, summary?.statusOnlineByConfig])

  const enrolled = summary?.devicesEnrolled
  const lastMonth = summary?.devicesEnrolledLastMonth
  const totalDevices = summary?.devicesTotal

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Overview from{' '}
          <code className="bg-muted rounded px-1">GET /private/summary/devices</code> — refreshed every minute in the
          background.
        </p>
      </div>

      {loading ? <p className="text-muted-foreground text-sm">Loading dashboard…</p> : null}

      {error ? (
        <div className="border-destructive/50 bg-destructive/10 flex flex-wrap items-center gap-2 rounded-lg border px-3 py-2 text-sm">
          <AlertCircle className="text-destructive h-4 w-4 shrink-0" />
          <span className="flex-1">{error}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => void doInitialLoad()}>
            Retry
          </Button>
        </div>
      ) : null}

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
        <Card>
          <CardHeader className="py-3">
            <CardTitle className="text-muted-foreground text-xs font-medium uppercase">Devices total</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? statSkeleton() : <p className="text-2xl font-semibold tabular-nums">{String(totalDevices ?? '—')}</p>}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="py-3">
            <CardTitle className="text-muted-foreground text-xs font-medium uppercase">Enrolled</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? statSkeleton() : <p className="text-2xl font-semibold tabular-nums">{String(enrolled ?? '—')}</p>}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="py-3">
            <CardTitle className="text-muted-foreground text-xs font-medium uppercase">Online (recent)</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? (
              statSkeleton()
            ) : (
              <p className="text-2xl font-semibold tabular-nums">
                {String(parseStatusSummary(summary?.statusSummary).online)}
              </p>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="py-3">
            <CardTitle className="text-muted-foreground text-xs font-medium uppercase">Configurations</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? statSkeleton() : <p className="text-2xl font-semibold tabular-nums">{String(counts?.configurationCount ?? '—')}</p>}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="py-3">
            <CardTitle className="text-muted-foreground text-xs font-medium uppercase">Applications</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {loading ? statSkeleton() : <p className="text-2xl font-semibold tabular-nums">{String(counts?.applicationCount ?? '—')}</p>}
          </CardContent>
        </Card>
      </div>

      {!loading && summary ? (
        <div className="grid gap-4 lg:grid-cols-3">
          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Connectivity</CardTitle>
            </CardHeader>
            <CardContent className="h-52">
              {statusPie.length === 0 ? (
                <p className="text-muted-foreground text-sm">No status data.</p>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie data={statusPie} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label>
                      {statusPie.map((_, i) => (
                        <Cell key={i} fill={COL_STATUS[i % COL_STATUS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="text-sm">App installation</CardTitle>
            </CardHeader>
            <CardContent className="h-52">
              {installPie.length === 0 ? (
                <p className="text-muted-foreground text-sm">No install summary.</p>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie data={installPie} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label>
                      {installPie.map((_, i) => (
                        <Cell key={i} fill={COL_INSTALL[i % COL_INSTALL.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Enrollment window</CardTitle>
            </CardHeader>
            <CardContent className="h-52">
              {enrollmentSplit.every((x) => x.value === 0) ? (
                <p className="text-muted-foreground text-sm">No enrollment split.</p>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie data={enrollmentSplit} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label />
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </div>
      ) : null}

      {!loading && monthlyBar.length > 0 ? (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Enrollments by month</CardTitle>
          </CardHeader>
          <CardContent className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={monthlyBar} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                <XAxis dataKey="label" tick={{ fontSize: 11 }} />
                <YAxis allowDecimals={false} tick={{ fontSize: 11 }} />
                <Tooltip />
                <Bar dataKey="enrollments" name="Devices" fill={MONTH_COLOR} radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      ) : null}

      {!loading && byConfigCombined.length > 0 ? (
        <Card>
          <CardHeader>
            <CardTitle className="text-sm">Top configurations — connectivity mix</CardTitle>
          </CardHeader>
          <CardContent className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={byConfigCombined} margin={{ top: 8, right: 16, left: 0, bottom: 8 }}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                <XAxis dataKey="name" tick={{ fontSize: 10 }} interval={0} angle={-20} textAnchor="end" height={60} />
                <YAxis allowDecimals={false} tick={{ fontSize: 11 }} />
                <Tooltip />
                <Legend />
                <Bar dataKey="offline" stackId="a" fill="#F7464A" name="Offline" />
                <Bar dataKey="idle" stackId="a" fill="#FDB45C" name="Idle" />
                <Bar dataKey="online" stackId="a" fill="#46BFBD" name="Online" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle className="text-sm">Recently updated devices</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-2">
              {[1, 2, 3, 4, 5].map((i) => (
                <Skeleton key={i} className="h-8 w-full" />
              ))}
            </div>
          ) : recent.length === 0 ? (
            <p className="text-muted-foreground text-sm">No recent devices returned.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Last seen</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recent.map((d) => (
                  <TableRow key={d.id}>
                    <TableCell>{d.description?.trim() || d.number}</TableCell>
                    <TableCell>{formatLastSeen(d.lastUpdate)}</TableCell>
                    <TableCell>
                      <StatusBadge statusCode={d.statusCode} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {!loading ? (
        <p className="text-muted-foreground text-xs">
          Enrolled last 30 days:&nbsp;<span className="tabular-nums font-medium">{lastMonth ?? '—'}</span>
        </p>
      ) : null}
    </div>
  )
}
