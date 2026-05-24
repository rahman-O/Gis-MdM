import { useCallback, useEffect, useMemo, useState } from 'react'
import { AlertCircle, MoreHorizontal, Plus, QrCode, Search, Settings2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/shared/ui/dropdown-menu'
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from '@/shared/ui/pagination'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { useDebounce } from '@/shared/hooks/useDebounce'
import * as deviceService from '@/features/devices/deviceService'
import * as configurationService from '@/features/configurations/configurationService'
import type { Configuration } from '@/features/configurations/types'
import { getConfigurationQrEligibility } from '@/features/configurations/configurationQr'
import { canEnrollDevicesViaQr } from '@/features/auth/permissions'
import { EnrollmentQrExperience } from '@/features/devices/EnrollmentQrExperience'
import type { ConfigurationOption, ConfigurationView, DeviceFilters, DeviceSearchRequest, DeviceView, LookupItem } from '@/features/devices/types'
import { StatusBadge } from '@/features/devices/StatusBadge'
import { formatLastSeen } from '@/features/devices/deviceFormat'
import { DeleteDialog } from '@/features/devices/DeleteDialog'
import { DeviceDetailPanel } from '@/features/devices/DeviceDetailPanel'
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

function groupsCell(groups: DeviceView['groups']): string {
  if (!groups?.length) return '—'
  return groups.map((g) => g.name ?? `#${g.id}`).join(', ')
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
  const [deviceToEdit, setDeviceToEdit] = useState<DeviceView | null>(null)
  const [bulkAction, setBulkAction] = useState<'delete' | 'configuration' | 'group' | null>(null)
  const [bulkConfigurationId, setBulkConfigurationId] = useState<number | null>(null)
  const [bulkGroupId, setBulkGroupId] = useState<number | null>(null)
  const [bulkError, setBulkError] = useState<string | null>(null)
  const [qrLoadingId, setQrLoadingId] = useState<number | null>(null)
  const [qrEnrollmentContext, setQrEnrollmentContext] = useState<{
    qrCodeKey: string
    deviceNumber: string
    configuration: Configuration | null
  } | null>(null)
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

  const getQrCodeKey = (device: DeviceView): string | null => {
    if (device.configurationId == null) return null
    const config = configurationsMap[device.configurationId]
    const qrCodeKey = config?.qrCodeKey?.trim()
    return qrCodeKey || null
  }

  const openDeviceQr = async (device: DeviceView) => {
    if (!canEnrollDevicesViaQr()) {
      setError('You do not have permission to enroll devices via QR.')
      return
    }

    const deviceNumber = device.number?.trim()
    if (!deviceNumber) {
      setError('Device number is missing.')
      return
    }
    if (device.configurationId == null) {
      setError('This device has no configuration assigned.')
      return
    }

    setQrLoadingId(device.id)
    let qrCodeKey: string | null = getQrCodeKey(device)
    let fullConfiguration: Configuration | null = null

    try {
      fullConfiguration = await configurationService.getConfiguration(device.configurationId)
      const resolvedKey = String(fullConfiguration.qrCodeKey ?? '').trim()
      if (resolvedKey) qrCodeKey = resolvedKey
    } catch {
      // ignore load error; qrCodeKey may still come from list map
    } finally {
      setQrLoadingId(null)
    }

    if (!qrCodeKey?.trim()) {
      const eligibility = getConfigurationQrEligibility(fullConfiguration)
      setError(eligibility.reason ?? 'QR is not available for this device configuration.')
      return
    }

    setQrEnrollmentContext({
      qrCodeKey: qrCodeKey.trim(),
      deviceNumber,
      configuration: fullConfiguration,
    })
  }

  return (
    <div className="flex gap-4">
      <DeviceTreeSidebar
        selectedNodeId={selectedTreeNodeId}
        onSelectNode={setSelectedTreeNodeId}
        onTreeChanged={() => void fetchDevices(true)}
      />
      <div className="min-w-0 flex-1 space-y-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Devices</h1>
          <p className="text-muted-foreground text-sm">Manage enrolled devices.</p>
        </div>
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm">
                <Settings2 className="mr-2 h-4 w-4" />
                Columns
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              {allColumns.map((column) => (
                <DropdownMenuCheckboxItem
                  key={column.key}
                  checked={columnVisibility[column.key]}
                  onCheckedChange={(checked) => setColumnVisibility((current) => ({ ...current, [column.key]: checked === true }))}
                >
                  {column.label}
                </DropdownMenuCheckboxItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <Button onClick={() => { setFormMode('create'); setDeviceToEdit(null) }}>
            <Plus className="mr-2 h-4 w-4" />
            Add Device
          </Button>
        </div>
      </div>

      <div className="flex max-w-xl items-center gap-2">
        <div className="relative flex-1">
          <Search className="text-muted-foreground absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2" />
          <Input className="pl-9" placeholder="Search devices..." aria-label="Search devices" value={search} onChange={(event) => setSearch(event.target.value)} />
        </div>
        <Button variant="outline" onClick={() => { setPage(1); void fetchDevices(true) }}>Search</Button>
      </div>

      <FilterPanel filters={filters} groups={groups} configurations={configurations} onChange={handleFiltersChange} />

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

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-8">
                <Checkbox checked={isAllSelected} onCheckedChange={(checked) => setSelectedIds(checked ? new Set(devices.map((d) => d.id)) : new Set())} />
              </TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Seen</TableHead>
              <TableHead>Number</TableHead>
              <TableHead>Configuration</TableHead>
              <TableHead>Enrollment</TableHead>
              <TableHead>Groups</TableHead>
              {columnVisibility.imei ? <TableHead>IMEI</TableHead> : null}
              {columnVisibility.phone ? <TableHead>Phone</TableHead> : null}
              {columnVisibility.model ? <TableHead>Model</TableHead> : null}
              {columnVisibility.battery ? <TableHead>Battery</TableHead> : null}
              {columnVisibility.android ? <TableHead>Android</TableHead> : null}
              {columnVisibility.serial ? <TableHead>Serial</TableHead> : null}
              {columnVisibility.description ? <TableHead>Description</TableHead> : null}
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? Array.from({ length: 5 }).map((_, index) => (
              <TableRow key={index}><TableCell colSpan={15}><Skeleton className="h-9 w-full" /></TableCell></TableRow>
            )) : devices.length === 0 ? (
              <TableRow>
                <TableCell colSpan={15} className="h-24 text-center text-muted-foreground">
                  {debouncedSearch.trim() ? `No devices found for '${debouncedSearch.trim()}'.` : 'No devices yet.'}
                </TableCell>
              </TableRow>
            ) : devices.map((device) => (
              <TableRow key={device.id} className="cursor-pointer" onClick={() => setSelectedDevice(device)}>
                <TableCell onClick={(event) => event.stopPropagation()}>
                  <Checkbox
                    checked={selectedIds.has(device.id)}
                    onCheckedChange={(checked) => {
                      const next = new Set(selectedIds)
                      if (checked) next.add(device.id)
                      else next.delete(device.id)
                      setSelectedIds(next)
                    }}
                  />
                </TableCell>
                <TableCell><StatusBadge statusCode={device.statusCode} lastUpdate={device.lastUpdate} /></TableCell>
                <TableCell>{formatLastSeen(device.lastUpdate)}</TableCell>
                <TableCell className="font-medium">{device.number || '—'}</TableCell>
                <TableCell>{device.configurationId != null ? configurationsMap[device.configurationId]?.name?.trim() || `Configuration #${device.configurationId}` : '—'}</TableCell>
                <TableCell className="capitalize text-muted-foreground text-xs">{device.enrollmentState ?? '—'}</TableCell>
                <TableCell className="max-w-[240px] truncate">{groupsCell(device.groups)}</TableCell>
                {columnVisibility.imei ? <TableCell>{device.imei || '—'}</TableCell> : null}
                {columnVisibility.phone ? <TableCell>{device.phone || '—'}</TableCell> : null}
                {columnVisibility.model ? <TableCell>{device.model || '—'}</TableCell> : null}
                {columnVisibility.battery ? <TableCell>{device.batteryLevel != null ? `${device.batteryLevel}%` : '—'}</TableCell> : null}
                {columnVisibility.android ? <TableCell>{device.androidVersion || '—'}</TableCell> : null}
                {columnVisibility.serial ? <TableCell>{device.serial || '—'}</TableCell> : null}
                {columnVisibility.description ? <TableCell>{device.description || '—'}</TableCell> : null}
                <TableCell className="text-right" onClick={(event) => event.stopPropagation()}>
                  <div className="inline-flex items-center gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      disabled={!canEnrollDevicesViaQr() || qrLoadingId === device.id}
                      title="Open enrollment QR"
                      onClick={() => void openDeviceQr(device)}
                    >
                      <QrCode className="h-4 w-4" />
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8"><MoreHorizontal className="h-4 w-4" /></Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onSelect={() => setSelectedDevice(device)}>View Details</DropdownMenuItem>
                        <DropdownMenuItem onSelect={() => { setDeviceToEdit(device); setFormMode('edit') }}>Edit</DropdownMenuItem>
                        <DropdownMenuItem className="text-destructive focus:text-destructive" onSelect={() => setDeviceToDelete(device)}>Delete</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </TableCell>
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

      <DeviceDetailPanel deviceNumber={selectedDevice?.number ?? null} onClose={() => setSelectedDevice(null)} />
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

      <Dialog
        open={qrEnrollmentContext != null}
        onOpenChange={(open) => {
          if (!open) setQrEnrollmentContext(null)
        }}
      >
        <DialogContent className="flex max-h-[90vh] max-w-3xl flex-col overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Device enrollment QR</DialogTitle>
            <DialogDescription>Adjust device id / provisioning options; the QR updates automatically (same as legacy UI).</DialogDescription>
          </DialogHeader>
          {qrEnrollmentContext ? (
            <EnrollmentQrExperience
              key={`${qrEnrollmentContext.qrCodeKey}-${qrEnrollmentContext.deviceNumber}`}
              qrCodeKey={qrEnrollmentContext.qrCodeKey}
              initialDeviceId={qrEnrollmentContext.deviceNumber}
              configuration={qrEnrollmentContext.configuration}
              groups={groups}
              footer={
                <Button type="button" variant="outline" onClick={() => setQrEnrollmentContext(null)}>
                  Close
                </Button>
              }
            />
          ) : null}
        </DialogContent>
      </Dialog>

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
