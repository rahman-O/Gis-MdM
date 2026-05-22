import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { AlertCircle, Plus } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from '@/shared/ui/pagination'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/shared/ui/alert-dialog'
import { useDebounce } from '@/shared/hooks/useDebounce'
import { hasPermission } from '@/features/auth/permissions'
import * as applicationService from '@/features/applications/services/applicationService'
import { ApplicationFilters } from '@/features/applications/components/ApplicationFilters'
import { ApplicationsTable } from '@/features/applications/components/ApplicationsTable'
import { ApplicationFormDialog } from '@/features/applications/components/ApplicationFormDialog'
import { ApplicationConfigurationsDialog } from '@/features/applications/components/ApplicationConfigurationsDialog'
import type { Application } from '@/features/applications/model/types'
import * as webUiFilesService from '@/features/applications/services/webUiFilesService'

const MY_APPS_KEY = 'HMDM_showMyAppsOnly'
const SYSTEM_APPS_KEY = 'HMDM_showSystemApps'
const PAGE_SIZE = 50

export function ApplicationsPage() {
  const navigate = useNavigate()
  const [list, setList] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [showMyAppsOnly, setShowMyAppsOnly] = useState(
    () => window.localStorage.getItem(MY_APPS_KEY) === 'true'
  )
  const [showSystemApps, setShowSystemApps] = useState(
    () => window.localStorage.getItem(SYSTEM_APPS_KEY) === 'true'
  )
  const [formOpen, setFormOpen] = useState(false)
  const [selected, setSelected] = useState<Application | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Application | null>(null)
  const [cfgAppId, setCfgAppId] = useState<number | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [storageBanner, setStorageBanner] = useState<string | null>(null)
  const debouncedSearch = useDebounce(search, 400)

  const canEditApps = hasPermission('edit_applications')
  const canViewApps = hasPermission('applications')

  const load = async () => {
    setLoading(true)
    setError(null)
    try {
      const data = debouncedSearch.trim()
        ? await applicationService.searchApplications(debouncedSearch)
        : await applicationService.getAllApplications()
      setList(Array.isArray(data) ? data : [])
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load applications.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    window.localStorage.setItem(MY_APPS_KEY, showMyAppsOnly ? 'true' : 'false')
  }, [showMyAppsOnly])
  useEffect(() => {
    window.localStorage.setItem(SYSTEM_APPS_KEY, showSystemApps ? 'true' : 'false')
  }, [showSystemApps])

  useEffect(() => {
    if (!canViewApps) return
    void load()
  }, [debouncedSearch, canViewApps])

  useEffect(() => {
    setCurrentPage(1)
  }, [debouncedSearch, showMyAppsOnly, showSystemApps])

  useEffect(() => {
    if (!canViewApps) {
      setStorageBanner(null)
      return
    }
    let cancelled = false
    void webUiFilesService
      .getStorageLimit()
      .then((lim) => {
        if (cancelled) return
        if (lim.sizeLimit > 0) {
          const available = Math.max(0, lim.sizeLimit - lim.sizeUsed)
          if (available < 20) {
            setStorageBanner(`Available space: ${available} Mb`)
          } else {
            setStorageBanner(null)
          }
        } else {
          setStorageBanner(null)
        }
      })
      .catch(() => {
        if (!cancelled) setStorageBanner(null)
      })
    return () => {
      cancelled = true
    }
  }, [canViewApps, debouncedSearch])

  const filtered = useMemo(() => {
    return list.filter((app) => {
      if (!showSystemApps && app.system) return false
      if (showMyAppsOnly && app.commonApplication && app.customerId == null) return false
      return true
    })
  }, [list, showSystemApps, showMyAppsOnly])

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE))
  useEffect(() => {
    if (currentPage > totalPages) {
      setCurrentPage(totalPages)
    }
  }, [currentPage, totalPages])

  const paginatedApps = useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE
    return filtered.slice(start, start + PAGE_SIZE)
  }, [filtered, currentPage])

  if (!canViewApps) {
    return <p className="text-sm text-muted-foreground">You do not have permission to view applications.</p>
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Applications</h1>
          <p className="text-sm text-muted-foreground">Manage application catalog and deployment behavior.</p>
        </div>
        <Button onClick={() => { setSelected(null); setFormOpen(true) }} disabled={!canEditApps}>
          <Plus className="mr-2 h-4 w-4" />
          Add application
        </Button>
      </div>

      <ApplicationFilters
        search={search}
        onSearchChange={setSearch}
        showSystemApps={showSystemApps}
        onShowSystemAppsChange={setShowSystemApps}
        showMyAppsOnly={showMyAppsOnly}
        onShowMyAppsOnlyChange={setShowMyAppsOnly}
      />

      {storageBanner ? (
        <div className="rounded border border-amber-500/40 bg-amber-500/10 px-3 py-2 text-sm text-amber-950 dark:text-amber-100">
          {storageBanner}
        </div>
      ) : null}

      {error ? (
        <div className="flex items-center gap-2 rounded border border-destructive/50 bg-destructive/10 p-3 text-sm">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button variant="outline" size="sm" onClick={() => void load()}>Retry</Button>
        </div>
      ) : null}

      {loading ? (
        <div className="space-y-2">
          {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-9 w-full" />)}
        </div>
      ) : (
        <>
          <ApplicationsTable
            applications={paginatedApps}
            onEdit={(app) => { setSelected(app); setFormOpen(true) }}
            onDelete={(app) => setDeleteTarget(app)}
            onVersions={(app) => app.id != null && navigate(`/application/${app.id}/versions`)}
            onConfigurations={(app) => setCfgAppId(app.id ?? null)}
          />
          {filtered.length > PAGE_SIZE ? (
            <div className="flex flex-wrap items-center justify-between gap-3">
              <p className="text-sm text-muted-foreground">
                Showing {(currentPage - 1) * PAGE_SIZE + 1}–
                {Math.min(currentPage * PAGE_SIZE, filtered.length)} of {filtered.length}
              </p>
              <Pagination>
                <PaginationContent className="flex-wrap justify-end gap-2">
                  <PaginationItem>
                    <PaginationPrevious
                      disabled={currentPage <= 1}
                      onClick={() => currentPage > 1 && setCurrentPage((p) => p - 1)}
                    />
                  </PaginationItem>
                  <PaginationItem>
                    <span className="px-2 text-sm text-muted-foreground">
                      Page {currentPage} / {totalPages}
                    </span>
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationNext
                      disabled={currentPage >= totalPages}
                      onClick={() => currentPage < totalPages && setCurrentPage((p) => p + 1)}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          ) : null}
        </>
      )}

      <ApplicationFormDialog
        open={formOpen}
        initialData={selected}
        onClose={() => setFormOpen(false)}
        onSaved={load}
      />

      <ApplicationConfigurationsDialog
        open={cfgAppId != null}
        applicationId={cfgAppId}
        onClose={() => setCfgAppId(null)}
      />

      <AlertDialog open={deleteTarget != null} onOpenChange={(v) => !v && setDeleteTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete application?</AlertDialogTitle>
            <AlertDialogDescription>
              This action can fail if the application is referenced by configurations.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                void (async () => {
                  if (deleteTarget?.id == null) return
                  const id = deleteTarget.id
                  try {
                    await applicationService.deleteApplication(id)
                    setDeleteTarget(null)
                    setList((prev) => prev.filter((a) => a.id !== id))
                    setError(null)
                    await load()
                  } catch (reason: unknown) {
                    setError(reason instanceof Error ? reason.message : 'Failed to delete application.')
                  }
                })()
              }}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
