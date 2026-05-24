import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { AlertCircle, Smartphone, UserCheck, Activity, Settings, AppWindow, RefreshCw } from 'lucide-react'
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
import { OnboardingChecklist } from '@/features/onboarding/OnboardingChecklist'

interface CustomTooltipProps {
  active?: boolean
  payload?: any[]
  label?: string
}

function CustomTooltip({ active, payload, label }: CustomTooltipProps) {
  if (active && payload && payload.length) {
    return (
      <div className="bg-background/95 border-border/80 backdrop-blur-md rounded-lg border p-2.5 shadow-lg max-w-[220px]">
        {label && <p className="text-[10px] font-bold text-muted-foreground mb-1.5 uppercase tracking-wider">{label}</p>}
        <div className="space-y-1">
          {payload.map((item, idx) => (
            <div key={idx} className="flex items-center gap-2 text-xs">
              <span className="h-2 w-2 rounded-full shrink-0" style={{ backgroundColor: item.color || item.fill }} />
              <span className="font-medium text-muted-foreground truncate">{item.name}:</span>
              <span className="font-semibold tabular-nums text-foreground ml-auto">{item.value}</span>
            </div>
          ))}
        </div>
      </div>
    )
  }
  return null
}

function statSkeleton() {
  return <Skeleton className="h-8 w-24" />
}

export function DashboardPage() {
  const [summary, setSummary] = useState<DeviceSummaryPayload | null>(null)
  const [counts, setCounts] = useState<{ configurationCount: number; applicationCount: number } | null>(null)
  const [recent, setRecent] = useState<DeviceView[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [lastRefreshedAt, setLastRefreshedAt] = useState<string>('')

  const initialLoadDone = useRef(false)

  const formatTime = () => new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })

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
      setLastRefreshedAt(formatTime())
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
      setLastRefreshedAt(formatTime())
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

  const handleManualRefresh = useCallback(async () => {
    await doInitialLoad()
  }, [doInitialLoad])

  const statusPie = useMemo(() => {
    const s = parseStatusSummary(summary?.statusSummary ?? undefined)
    return [
      { name: 'Offline', value: s.offline, fill: '#F43F5E' },
      { name: 'Idle', value: s.idle, fill: '#F59E0B' },
      { name: 'Online', value: s.online, fill: '#10B981' },
    ].filter((d) => d.value > 0)
  }, [summary?.statusSummary])

  const installPie = useMemo(() => {
    const s = parseInstallSummary(summary?.installSummary ?? undefined)
    return [
      { name: 'Install failed', value: s.failure, fill: '#F43F5E' },
      { name: 'Version mismatch', value: s.mismatch, fill: '#F59E0B' },
      { name: 'Installed OK', value: s.success, fill: '#10B981' },
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
      { name: 'Earlier', value: earlier, fill: '#94A3B8' },
      { name: 'Last 30 days', value: lastMo, fill: '#6366F1' },
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

  const totalStatus = useMemo(() => statusPie.reduce((acc, curr) => acc + curr.value, 0), [statusPie])
  const totalInstall = useMemo(() => installPie.reduce((acc, curr) => acc + curr.value, 0), [installPie])
  const totalEnrollment = useMemo(() => enrollmentSplit.reduce((acc, curr) => acc + curr.value, 0), [enrollmentSplit])

  const enrolled = summary?.devicesEnrolled
  const lastMonth = summary?.devicesEnrolledLastMonth
  const totalDevices = summary?.devicesTotal

  const totalDevicesNum = summary ? Number(summary.devicesTotal || 0) : 0
  const enrolledNum = summary ? Number(summary.devicesEnrolled || 0) : 0
  const enrollmentRate = totalDevicesNum ? Math.round((enrolledNum / totalDevicesNum) * 100) : 0

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b pb-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-1 text-sm">
            Real-time overview of your managed devices, configurations, and applications.
          </p>
        </div>
        <div className="flex items-center gap-3">
          {lastRefreshedAt && (
            <span className="text-xs text-muted-foreground tabular-nums">
              Last updated: {lastRefreshedAt}
            </span>
          )}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleManualRefresh}
            disabled={loading}
            className="flex items-center gap-2 hover:bg-accent hover:text-accent-foreground active:scale-95 transition-transform"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      <OnboardingChecklist />

      {error ? (
        <div className="border-destructive/50 bg-destructive/10 flex flex-wrap items-center gap-2 rounded-lg border px-3 py-2 text-sm">
          <AlertCircle className="text-destructive h-4 w-4 shrink-0" />
          <span className="flex-1">{error}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => void doInitialLoad()}>
            Retry
          </Button>
        </div>
      ) : null}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
        {/* Total Devices */}
        <Card className="hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 transition-all duration-200 overflow-hidden relative group">
          <div className="absolute top-4 right-4 text-muted-foreground group-hover:text-primary transition-colors duration-200">
            <div className="bg-primary/5 p-2 rounded-lg group-hover:bg-primary/10 transition-colors">
              <Smartphone className="h-4 w-4" />
            </div>
          </div>
          <CardHeader className="py-4">
            <CardTitle className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Devices Total</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 pb-4">
            {loading ? (
              statSkeleton()
            ) : (
              <div>
                <p className="text-2xl font-bold tabular-nums tracking-tight">{String(totalDevices ?? '—')}</p>
                <p className="text-muted-foreground mt-1 text-[11px] font-normal leading-none">
                  {enrollmentRate ? `${enrollmentRate}% enrollment rate` : 'All registered in system'}
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Enrolled */}
        <Card className="hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 transition-all duration-200 overflow-hidden relative group">
          <div className="absolute top-4 right-4 text-muted-foreground group-hover:text-primary transition-colors duration-200">
            <div className="bg-primary/5 p-2 rounded-lg group-hover:bg-primary/10 transition-colors">
              <UserCheck className="h-4 w-4" />
            </div>
          </div>
          <CardHeader className="py-4">
            <CardTitle className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Enrolled</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 pb-4">
            {loading ? (
              statSkeleton()
            ) : (
              <div>
                <p className="text-2xl font-bold tabular-nums tracking-tight">{String(enrolled ?? '—')}</p>
                <p className="text-muted-foreground mt-1 text-[11px] font-normal leading-none">
                  {lastMonth !== undefined ? `+${lastMonth} in last 30 days` : 'Active enrolments'}
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Online (recent) */}
        <Card className="hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 transition-all duration-200 overflow-hidden relative group">
          <div className="absolute top-4 right-4 text-muted-foreground group-hover:text-primary transition-colors duration-200">
            <div className="bg-primary/5 p-2 rounded-lg group-hover:bg-primary/10 transition-colors">
              <Activity className="h-4 w-4" />
            </div>
          </div>
          <CardHeader className="py-4">
            <CardTitle className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Online (Recent)</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 pb-4">
            {loading ? (
              statSkeleton()
            ) : (
              <div>
                <p className="text-2xl font-bold tabular-nums tracking-tight">
                  {String(parseStatusSummary(summary?.statusSummary).online)}
                </p>
                <p className="text-muted-foreground mt-1 text-[11px] font-normal leading-none">
                  {parseStatusSummary(summary?.statusSummary).offline} offline devices
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Configurations */}
        <Card className="hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 transition-all duration-200 overflow-hidden relative group">
          <div className="absolute top-4 right-4 text-muted-foreground group-hover:text-primary transition-colors duration-200">
            <div className="bg-primary/5 p-2 rounded-lg group-hover:bg-primary/10 transition-colors">
              <Settings className="h-4 w-4" />
            </div>
          </div>
          <CardHeader className="py-4">
            <CardTitle className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Configurations</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 pb-4">
            {loading ? (
              statSkeleton()
            ) : (
              <div>
                <p className="text-2xl font-bold tabular-nums tracking-tight">{String(counts?.configurationCount ?? '—')}</p>
                <p className="text-muted-foreground mt-1 text-[11px] font-normal leading-none">
                  Active profiles in use
                </p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Applications */}
        <Card className="hover:shadow-md hover:border-primary/20 hover:-translate-y-0.5 transition-all duration-200 overflow-hidden relative group">
          <div className="absolute top-4 right-4 text-muted-foreground group-hover:text-primary transition-colors duration-200">
            <div className="bg-primary/5 p-2 rounded-lg group-hover:bg-primary/10 transition-colors">
              <AppWindow className="h-4 w-4" />
            </div>
          </div>
          <CardHeader className="py-4">
            <CardTitle className="text-muted-foreground text-xs font-semibold uppercase tracking-wider">Applications</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 pb-4">
            {loading ? (
              statSkeleton()
            ) : (
              <div>
                <p className="text-2xl font-bold tabular-nums tracking-tight">{String(counts?.applicationCount ?? '—')}</p>
                <p className="text-muted-foreground mt-1 text-[11px] font-normal leading-none">
                  Managed app catalog
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {!loading && summary ? (
        <div className="grid gap-6 md:grid-cols-3">
          {/* Connectivity */}
          <Card className="overflow-hidden hover:shadow-md transition-shadow duration-200">
            <CardHeader className="border-b bg-muted/20 py-3">
              <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Connectivity Status</CardTitle>
            </CardHeader>
            <CardContent className="h-56 pt-6">
              {statusPie.length === 0 ? (
                <div className="flex h-full items-center justify-center">
                  <p className="text-muted-foreground text-sm">No status data.</p>
                </div>
              ) : (
                <div className="relative flex items-center justify-center h-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={statusPie}
                        dataKey="value"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        innerRadius={52}
                        outerRadius={72}
                        paddingAngle={3}
                      >
                        {statusPie.map((entry, i) => (
                          <Cell key={i} fill={entry.fill} />
                        ))}
                      </Pie>
                      <Tooltip content={<CustomTooltip />} />
                    </PieChart>
                  </ResponsiveContainer>
                  <div className="absolute flex flex-col items-center justify-center text-center pointer-events-none">
                    <span className="text-2xl font-bold tracking-tight leading-none tabular-nums">{totalStatus}</span>
                    <span className="text-[10px] uppercase font-semibold text-muted-foreground mt-0.5 tracking-wider">Devices</span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* App installation */}
          <Card className="overflow-hidden hover:shadow-md transition-shadow duration-200">
            <CardHeader className="border-b bg-muted/20 py-3">
              <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">App Installation</CardTitle>
            </CardHeader>
            <CardContent className="h-56 pt-6">
              {installPie.length === 0 ? (
                <div className="flex h-full items-center justify-center">
                  <p className="text-muted-foreground text-sm">No install summary.</p>
                </div>
              ) : (
                <div className="relative flex items-center justify-center h-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={installPie}
                        dataKey="value"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        innerRadius={52}
                        outerRadius={72}
                        paddingAngle={3}
                      >
                        {installPie.map((entry, i) => (
                          <Cell key={i} fill={entry.fill} />
                        ))}
                      </Pie>
                      <Tooltip content={<CustomTooltip />} />
                    </PieChart>
                  </ResponsiveContainer>
                  <div className="absolute flex flex-col items-center justify-center text-center pointer-events-none">
                    <span className="text-2xl font-bold tracking-tight leading-none tabular-nums">{totalInstall}</span>
                    <span className="text-[10px] uppercase font-semibold text-muted-foreground mt-0.5 tracking-wider">Apps</span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Enrollment window */}
          <Card className="overflow-hidden hover:shadow-md transition-shadow duration-200">
            <CardHeader className="border-b bg-muted/20 py-3">
              <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Enrollment Window</CardTitle>
            </CardHeader>
            <CardContent className="h-56 pt-6">
              {totalEnrollment === 0 ? (
                <div className="flex h-full items-center justify-center">
                  <p className="text-muted-foreground text-sm">No enrollment split.</p>
                </div>
              ) : (
                <div className="relative flex items-center justify-center h-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={enrollmentSplit}
                        dataKey="value"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        innerRadius={52}
                        outerRadius={72}
                        paddingAngle={3}
                      >
                        {enrollmentSplit.map((entry, i) => (
                          <Cell key={i} fill={entry.fill} />
                        ))}
                      </Pie>
                      <Tooltip content={<CustomTooltip />} />
                    </PieChart>
                  </ResponsiveContainer>
                  <div className="absolute flex flex-col items-center justify-center text-center pointer-events-none">
                    <span className="text-2xl font-bold tracking-tight leading-none tabular-nums">{totalEnrollment}</span>
                    <span className="text-[10px] uppercase font-semibold text-muted-foreground mt-0.5 tracking-wider">Enrolled</span>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      ) : null}

      {!loading && monthlyBar.length > 0 ? (
        <Card className="hover:shadow-md transition-shadow duration-200">
          <CardHeader className="border-b bg-muted/20 py-3 flex flex-row items-center justify-between">
            <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Enrollments by month</CardTitle>
          </CardHeader>
          <CardContent className="h-64 pt-6">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={monthlyBar} margin={{ top: 8, right: 8, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" vertical={false} />
                <XAxis dataKey="label" tick={{ fontSize: 11, fill: '#64748B' }} tickLine={false} axisLine={false} />
                <YAxis allowDecimals={false} tick={{ fontSize: 11, fill: '#64748B' }} tickLine={false} axisLine={false} />
                <Tooltip content={<CustomTooltip />} />
                <Bar dataKey="enrollments" name="Devices" fill="#6366F1" radius={[4, 4, 0, 0]} maxBarSize={50} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      ) : null}

      {!loading && byConfigCombined.length > 0 ? (
        <Card className="hover:shadow-md transition-shadow duration-200">
          <CardHeader className="border-b bg-muted/20 py-3">
            <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Top configurations — connectivity mix</CardTitle>
          </CardHeader>
          <CardContent className="h-72 pt-6">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={byConfigCombined} margin={{ top: 8, right: 16, left: -20, bottom: 8 }}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" vertical={false} />
                <XAxis dataKey="name" tick={{ fontSize: 10, fill: '#64748B' }} tickLine={false} axisLine={false} interval={0} angle={-20} textAnchor="end" height={60} />
                <YAxis allowDecimals={false} tick={{ fontSize: 11, fill: '#64748B' }} tickLine={false} axisLine={false} />
                <Tooltip content={<CustomTooltip />} />
                <Legend iconType="circle" iconSize={8} wrapperStyle={{ fontSize: 11, paddingTop: 10 }} />
                <Bar dataKey="offline" stackId="a" fill="#F43F5E" name="Offline" maxBarSize={40} />
                <Bar dataKey="idle" stackId="a" fill="#F59E0B" name="Idle" maxBarSize={40} />
                <Bar dataKey="online" stackId="a" fill="#10B981" name="Online" maxBarSize={40} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      ) : null}

      <Card className="hover:shadow-md transition-shadow duration-200">
        <CardHeader className="border-b bg-muted/20 py-3">
          <CardTitle className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Recently updated devices</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-4 space-y-3">
              {[1, 2, 3, 4, 5].map((i) => (
                <Skeleton key={i} className="h-8 w-full" />
              ))}
            </div>
          ) : recent.length === 0 ? (
            <div className="p-4 text-center">
              <p className="text-muted-foreground text-sm">No recent devices returned.</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table>
                <TableHeader className="bg-muted/50">
                  <TableRow>
                    <TableHead className="font-semibold text-xs text-muted-foreground h-10 px-4">Name</TableHead>
                    <TableHead className="font-semibold text-xs text-muted-foreground h-10 px-4">Last seen</TableHead>
                    <TableHead className="font-semibold text-xs text-muted-foreground h-10 px-4 text-right">Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recent.map((d) => (
                    <TableRow key={d.id} className="hover:bg-muted/30 transition-colors border-b last:border-0">
                      <TableCell className="font-medium text-sm px-4 py-3">{d.description?.trim() || d.number}</TableCell>
                      <TableCell className="text-sm text-muted-foreground px-4 py-3">{formatLastSeen(d.lastUpdate)}</TableCell>
                      <TableCell className="px-4 py-3 text-right">
                        <StatusBadge statusCode={d.statusCode} lastUpdate={d.lastUpdate} />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>

      {!loading ? (
        <div className="flex items-center justify-between text-muted-foreground text-xs border-t pt-4">
          <p>
            Enrolled last 30 days:&nbsp;<span className="tabular-nums font-semibold text-foreground">{lastMonth ?? '—'}</span>
          </p>
          <p>
            Refreshed automatically every minute
          </p>
        </div>
      ) : null}
    </div>
  )
}

