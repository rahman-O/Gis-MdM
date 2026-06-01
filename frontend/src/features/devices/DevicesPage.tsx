import { useCallback, useEffect, useMemo, useState } from 'react'
import { AlertCircle, Search, Settings2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuTrigger } from '@/shared/ui/dropdown-menu'
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from '@/shared/ui/pagination'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { useDebounce } from '@/shared/hooks/useDebounce'
import * as deviceService from '@/features/devices/deviceService'
import type { ConfigurationOption, ConfigurationView, DeviceFilters, DeviceSearchRequest, DeviceView, LookupItem } from '@/features/devices/types'
import { formatLastSeen } from '@/features/devices/deviceFormat'
import { StatusBadge } from '@/features/devices/StatusBadge'
import { DeleteDialog } from '@/features/devices/DeleteDialog'
import { DeviceDetailsDialog } from '@/features/devices/DeviceDetailsDialog'
import { DeviceForm } from '@/features/devices/DeviceForm'
import { FilterPanel } from '@/features/devices/FilterPanel'
import { BulkActionBar } from '@/features/devices/BulkActionBar'
import { DeviceTreeSidebar } from '@/features/device-tree/DeviceTreeSidebar'

const PAGE_SIZE = 20
const EMPTY_FILTERS: DeviceFilters = {
  groupId: null,
  configurationId: null,
  status: null,
  androidVersion: null,
  launcherVersion: null,
  enrollmentDateFrom: null,
  enrollmentDateTo: null,
  onlineEarlierMillis: null,
  onlineLaterMillis: null,
  mdmMode: null,
  kioskMode: null,
  installationStatus: null,
  imeiChanged: null,
  fastSearch: null,
  sortBy: null,
  sortDir: null,
}

function renderGroups(groups: DeviceView['groups']) {
  if (!groups?.length) return <span className="text-muted-foreground/60">—</span>
  return (
    <div className="flex flex-wrap gap-1 max-w-[240px]">
      {groups.map((g) => (
        <span key={g.id} className="inline-flex items-center px-1.5 py-0.5 rounded-md text-[10px] font-medium bg-secondary/70 text-secondary-foreground border border-border/30">
          {g.name ?? `#${g.id}`}
        </span>
      ))}
    </div>
  )
}

function areFiltersEqual(a: DeviceFilters, b: DeviceFilters): boolean {
  const keys = Object.keys(EMPTY_FILTERS) as Array<keyof DeviceFilters>
  return keys.every((key) => a[key] === b[key])
}

export function DevicesPage() {
  const [devices, setDevices] = useState<DeviceView[]>([])
  const [configurationsMap, setConfigurationsMap] = useState<Record<number, ConfigurationView>>({})
  const [groups, setGroups] = useState<LookupItem[]>([])
  const [configurations, setConfigurations] = useState<ConfigurationOption[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [filters, setFilters] = useState<DeviceFilters>(EMPTY_FILTERS)
  const debouncedSearch = useDebounce(search, 300)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedDevice, setSelectedDevice] = useState<DeviceView | null>(null)
  const [deviceToDelete, setDeviceToDelete] = useState<DeviceView | null>(null)
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set())
  const [formMode, setFormMode] = useState<'create' | 'edit' | null>(null)
  const [deviceToEdit, _setDeviceToEdit] = useState<DeviceView | null>(null)
  const [bulkAction, setBulkAction] = useState<'delete' | 'configuration' | 'group' | null>(null)
  const [bulkConfigurationId, setBulkConfigurationId] = useState<number | null>(null)
  const [bulkGroupId, setBulkGroupId] = useState<number | null>(null)
  const [bulkError, setBulkError] = useState<string | null>(null)
  const [selectedTreeNodeId, setSelectedTreeNodeId] = useState<number | null>(null)
  const [columnVisibility, setColumnVisibility] = useState({
    imei: false,
    phone: false,
    model: false,
    battery: false,
    android: false,
    serial: false,
    description: false,
  })

  const handleFiltersChange = useCallback((nextFilters: DeviceFilters) => {
    setFilters((prev) => (areFiltersEqual(prev, nextFilters) ? prev : nextFilters))
    setPage(1)
  }, [])

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))
  const showPagination = total > PAGE_SIZE
  const selectedCount = selectedIds.size
  const isAllSelected = devices.length > 0 && devices.every((device) => selectedIds.has(device.id))
  const selectedDevices = useMemo(() => devices.filter((device) => selectedIds.has(device.id)), [devices, selectedIds])

  const buildRequest = (): DeviceSearchRequest => ({
    pageNum: page,
    pageSize: PAGE_SIZE,
    value: debouncedSearch.trim() || null,
    ...filters,
    treeNodeId: selectedTreeNodeId,
    includeDescendants: selectedTreeNodeId != null ? true : null,
  })

  const fetchDevices = async (withLoading = true) => {
    if (withLoading) setLoading(true)
    setError(null)
    try {
      const response = await deviceService.getDevices(buildRequest())
      setDevices(response.devices.items ?? [])
      setTotal(response.devices.totalItemsCount ?? 0)
      setConfigurationsMap(response.configurations ?? {})
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load devices.')
    } finally {
      if (withLoading) setLoading(false)
    }
  }

  useEffect(() => {
    let cancelled = false
    void Promise.allSettled([deviceService.getGroups(), deviceService.getConfigurations()]).then(([groupsResult, configurationsResult]) => {
      if (cancelled) return
      if (groupsResult.status === 'fulfilled') setGroups(groupsResult.value)
      if (configurationsResult.status === 'fulfilled') setConfigurations(configurationsResult.value)
    })
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    void fetchDevices(true)
    setSelectedIds(new Set())
  }, [page, debouncedSearch, filters, selectedTreeNodeId])

  useEffect(() => {
    const intervalId = window.setInterval(() => void fetchDevices(false), 60_000)
    return () => window.clearInterval(intervalId)
  }, [page, debouncedSearch, filters, selectedTreeNodeId])

  const allColumns = [
    { key: 'imei' as const, label: 'IMEI' },
    { key: 'phone' as const, label: 'Phone' },
    { key: 'model' as const, label: 'Model' },
    { key: 'battery' as const, label: 'Battery' },
    { key: 'android' as const, label: 'Android' },
    { key: 'serial' as const, label: 'Serial' },
    { key: 'description' as const, label: 'Description' },
  ]

  const executeBulkAction = async () => {
    setBulkError(null)
    try {
      if (bulkAction === 'delete') {
        await deviceService.deleteBulk(Array.from(selectedIds))
      } else if (bulkAction === 'configuration') {
        await Promise.allSettled(
          selectedDevices.map((device) =>
            deviceService.updateDevice({
              id: device.id,
              number: device.number,
              description: device.description,
              configurationId: bulkConfigurationId,
              imei: device.imei,
              phone: device.phone,
              custom1: device.custom1,
              custom2: device.custom2,
              custom3: device.custom3,
              oldNumber: device.oldNumber,
              groups: device.groups ?? [],
            })
          )
        )
      } else if (bulkAction === 'group') {
        const selectedGroup = groups.find((group) => group.id === bulkGroupId)
        if (!selectedGroup) throw new Error('Select a group first.')
        await deviceService.groupBulk({ ids: Array.from(selectedIds), action: 'set', groups: [selectedGroup] })
      }
      setBulkAction(null)
      setSelectedIds(new Set())
      await fetchDevices(true)
    } catch (reason: unknown) {
      setBulkError(reason instanceof Error ? reason.message : 'Bulk action failed.')
    }
  }

  return (
    <div className="flex gap-6 p-2 max-w-[1600px] mx-auto">
      <DeviceTreeSidebar
        selectedNodeId={selectedTreeNodeId}
        onSelectNode={setSelectedTreeNodeId}
        onTreeChanged={() => void fetchDevices(true)}
      />
      <div className="min-w-0 flex-1 space-y-5">
        {/* Title Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border/40 pb-4">
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground/90">Devices</h1>
            <p className="text-muted-foreground text-sm mt-1">Monitor, filter, and manage all your MDM enrolled devices in one place.</p>
          </div>
          <div className="flex items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="sm" className="h-9 rounded-lg border-border/80 shadow-xs hover:bg-accent/60">
                  <Settings2 className="mr-2 h-4 w-4 text-muted-foreground" />
                  Columns
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48 p-1 border border-border/80 shadow-md rounded-lg">
                {allColumns.map((column) => (
                  <DropdownMenuCheckboxItem
                    key={column.key}
                    checked={columnVisibility[column.key]}
                    onCheckedChange={(checked) => setColumnVisibility((current) => ({ ...current, [column.key]: checked === true }))}
                    className="text-xs"
                  >
                    {column.label}
                  </DropdownMenuCheckboxItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {/* Search & Filters Container */}
        <div className="bg-card border border-border/50 rounded-xl p-4 space-y-4 shadow-xs">
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <div className="relative flex-1 w-full">
              <Search className="text-muted-foreground/70 absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2" />
              <Input 
                className="pl-9 h-9 text-sm rounded-lg border-border/80 focus-visible:ring-1 focus-visible:ring-primary/45 focus-visible:ring-offset-0" 
                placeholder="Search devices by number, configuration, description..." 
                aria-label="Search devices" 
                value={search} 
                onChange={(event) => setSearch(event.target.value)} 
              />
            </div>
            <Button variant="default" className="h-9 rounded-lg px-4 shadow-sm w-full sm:w-auto" onClick={() => { setPage(1); void fetchDevices(true) }}>
              Search
            </Button>
          </div>

          <FilterPanel filters={filters} groups={groups} configurations={configurations} onChange={handleFiltersChange} />
        </div>

        {selectedCount > 0 ? (
          <BulkActionBar selectedCount={selectedCount} onDeleteSelected={() => setBulkAction('delete')} onSetConfiguration={() => setBulkAction('configuration')} onSetGroup={() => setBulkAction('group')} />
        ) : null}

        {error ? (
          <div className="flex flex-wrap items-center gap-2 rounded-lg border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm" role="alert">
            <AlertCircle className="h-4 w-4 shrink-0 text-destructive" />
            <span className="flex-1">{error}</span>
            <Button type="button" variant="outline" size="sm" onClick={() => void fetchDevices(true)}>Retry</Button>
          </div>
        ) : null}

        {/* Table Wrapper */}
        <div className="rounded-xl border border-border/60 bg-card overflow-hidden shadow-xs">
          <Table>
            <TableHeader className="bg-muted/40">
              <TableRow className="hover:bg-transparent border-b border-border/55">
                <TableHead className="w-10">
                  <Checkbox checked={isAllSelected} onCheckedChange={(checked) => setSelectedIds(checked ? new Set(devices.map((d) => d.id)) : new Set())} className="rounded" />
                </TableHead>
                <TableHead className="font-semibold text-foreground/80">Status</TableHead>
                <TableHead className="font-semibold text-foreground/80">Last Seen</TableHead>
                <TableHead className="font-semibold text-foreground/80">Number</TableHead>
                <TableHead className="font-semibold text-foreground/80">Configuration</TableHead>
                <TableHead className="font-semibold text-foreground/80">Enrollment</TableHead>
                <TableHead className="font-semibold text-foreground/80">Groups</TableHead>
                {columnVisibility.imei ? <TableHead className="font-semibold text-foreground/80">IMEI</TableHead> : null}
                {columnVisibility.phone ? <TableHead className="font-semibold text-foreground/80">Phone</TableHead> : null}
                {columnVisibility.model ? <TableHead className="font-semibold text-foreground/80">Model</TableHead> : null}
                {columnVisibility.battery ? <TableHead className="font-semibold text-foreground/80">Battery</TableHead> : null}
                {columnVisibility.android ? <TableHead className="font-semibold text-foreground/80">Android</TableHead> : null}
                {columnVisibility.serial ? <TableHead className="font-semibold text-foreground/80">Serial</TableHead> : null}
                {columnVisibility.description ? <TableHead className="font-semibold text-foreground/80">Description</TableHead> : null}
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={index}><TableCell colSpan={15} className="py-4"><Skeleton className="h-9 w-full rounded-md" /></TableCell></TableRow>
              )) : devices.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={15} className="h-32 text-center text-muted-foreground">
                    {debouncedSearch.trim() ? `No devices found for '${debouncedSearch.trim()}'.` : 'No devices yet.'}
                  </TableCell>
                </TableRow>
              ) : devices.map((device) => (
                <TableRow key={device.id} className="hover:bg-secondary/25 transition-colors duration-150 cursor-pointer border-b border-border/40" onClick={() => setSelectedDevice(device)}>
                  <TableCell onClick={(event) => event.stopPropagation()} className="py-3">
                    <Checkbox
                      checked={selectedIds.has(device.id)}
                      onCheckedChange={(checked) => {
                        const next = new Set(selectedIds)
                        if (checked) next.add(device.id)
                        else next.delete(device.id)
                        setSelectedIds(next)
                      }}
                      className="rounded"
                    />
                  </TableCell>
                  <TableCell className="py-3">
                    <StatusBadge statusCode={device.statusCode ?? (Date.now() - (device.lastUpdate ?? 0) < 5 * 60 * 1000 ? 'green' : 'red')} lastUpdate={device.lastUpdate} />
                  </TableCell>
                  <TableCell className="py-3 text-xs text-muted-foreground">{formatLastSeen(device.lastUpdate)}</TableCell>
                  <TableCell className="py-3 font-mono text-[13px] font-semibold text-foreground/85">{device.number || '—'}</TableCell>
                  <TableCell className="py-3 text-sm">{device.configurationId != null ? configurationsMap[device.configurationId]?.name?.trim() || `Configuration #${device.configurationId}` : '—'}</TableCell>
                  <TableCell className="py-3">
                    <span className="capitalize px-2 py-0.5 rounded-full text-[10px] font-medium bg-muted text-muted-foreground/85">
                      {device.enrollmentState ?? '—'}
                    </span>
                  </TableCell>
                  <TableCell className="py-3 max-w-[240px] truncate">{renderGroups(device.groups)}</TableCell>
                  {columnVisibility.imei ? <TableCell className="py-3 text-xs font-mono">{device.imei || device.info?.imei || '—'}</TableCell> : null}
                  {columnVisibility.phone ? <TableCell className="py-3 text-xs font-mono">{device.phone || device.info?.phone || '—'}</TableCell> : null}
                  {columnVisibility.model ? <TableCell className="py-3 text-xs">{device.model || device.info?.model || '—'}</TableCell> : null}
                  {columnVisibility.battery ? (
                    <TableCell className="py-3">
                      {(device.batteryLevel ?? device.info?.batteryLevel) != null ? (
                        <span className="inline-flex items-center gap-1 text-xs font-semibold">
                          {(device.batteryLevel ?? device.info?.batteryLevel)}%
                        </span>
                      ) : '—'}
                    </TableCell>
                  ) : null}
                  {columnVisibility.android ? <TableCell className="py-3 text-xs">{device.androidVersion || device.info?.androidVersion || '—'}</TableCell> : null}
                  {columnVisibility.serial ? <TableCell className="py-3 text-xs font-mono">{device.serial || device.info?.serial || '—'}</TableCell> : null}
                  {columnVisibility.description ? <TableCell className="py-3 text-xs text-muted-foreground max-w-[200px] truncate">{device.description || '—'}</TableCell> : null}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>

      {showPagination ? (
        <Pagination className="mx-0 w-full justify-end">
          <PaginationContent>
            <PaginationItem><PaginationPrevious disabled={page <= 1 || loading} onClick={() => setPage((p) => Math.max(1, p - 1))} /></PaginationItem>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((pnum) => (
              <PaginationItem key={pnum}><PaginationLink isActive={pnum === page} onClick={() => setPage(pnum)} disabled={loading}>{pnum}</PaginationLink></PaginationItem>
            ))}
            <PaginationItem><PaginationNext disabled={page >= totalPages || loading} onClick={() => setPage((p) => Math.min(totalPages, p + 1))} /></PaginationItem>
          </PaginationContent>
        </Pagination>
      ) : null}

      <DeviceDetailsDialog
        device={selectedDevice}
        configurationName={selectedDevice?.configurationId != null ? configurationsMap[selectedDevice.configurationId]?.name?.trim() || null : null}
        onClose={() => setSelectedDevice(null)}
      />
      <DeleteDialog
        device={deviceToDelete}
        onConfirm={async () => {
          if (!deviceToDelete?.id) throw new Error('Invalid device.')
          await deviceService.deleteDevice(deviceToDelete.id)
          await fetchDevices(true)
        }}
        onCancel={() => setDeviceToDelete(null)}
      />
      {formMode ? <DeviceForm mode={formMode} initialData={deviceToEdit} onSuccess={async () => { await fetchDevices(true) }} onClose={() => setFormMode(null)} /> : null}



      <Dialog open={bulkAction != null} onOpenChange={(open) => !open && setBulkAction(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>{bulkAction === 'delete' ? 'Delete selected devices' : bulkAction === 'configuration' ? 'Set configuration' : 'Set group'}</DialogTitle></DialogHeader>
          {bulkAction === 'configuration' ? (
            <Select value={bulkConfigurationId == null ? 'none' : String(bulkConfigurationId)} onValueChange={(value) => setBulkConfigurationId(value === 'none' ? null : Number(value))}>
              <SelectTrigger><SelectValue placeholder="Select configuration" /></SelectTrigger>
              <SelectContent>
                <SelectItem value="none">None</SelectItem>
                {configurations.map((item) => <SelectItem key={item.id} value={String(item.id)}>{item.name ?? `Configuration #${item.id}`}</SelectItem>)}
              </SelectContent>
            </Select>
          ) : null}
          {bulkAction === 'group' ? (
            <Select value={bulkGroupId == null ? 'none' : String(bulkGroupId)} onValueChange={(value) => setBulkGroupId(value === 'none' ? null : Number(value))}>
              <SelectTrigger><SelectValue placeholder="Select group" /></SelectTrigger>
              <SelectContent>
                <SelectItem value="none">None</SelectItem>
                {groups.map((item) => <SelectItem key={item.id} value={String(item.id)}>{item.name ?? `Group #${item.id}`}</SelectItem>)}
              </SelectContent>
            </Select>
          ) : null}
          {bulkError ? <p className="text-sm text-destructive">{bulkError}</p> : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setBulkAction(null)}>Cancel</Button>
            <Button onClick={() => void executeBulkAction()}>Apply</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      </div>
    </div>
  )
}
