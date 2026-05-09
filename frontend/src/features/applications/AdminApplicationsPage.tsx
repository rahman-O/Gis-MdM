import { useCallback, useEffect, useMemo, useState } from 'react'
import { AlertCircle, Plus, Search } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
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
import { hasPermission, isSuperAdmin } from '@/features/auth/permissions'
import * as applicationService from '@/features/applications/services/applicationService'
import { AdminApplicationsTable } from '@/features/applications/components/AdminApplicationsTable'
import { ApplicationFormDialog } from '@/features/applications/components/ApplicationFormDialog'
import type { Application } from '@/features/applications/model/types'

const PAGE_SIZE = 50

export function AdminApplicationsPage() {
  const [list, setList] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [formOpen, setFormOpen] = useState(false)
  const [selected, setSelected] = useState<Application | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Application | null>(null)
  const [turnTarget, setTurnTarget] = useState<Application | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const debouncedSearch = useDebounce(search, 400)

  const superAdmin = isSuperAdmin()
  const canEditApps = hasPermission('edit_applications')

  const load = useCallback(async () => {
    if (!isSuperAdmin()) return
    setLoading(true)
    setError(null)
    try {
      const data = debouncedSearch.trim()
        ? await applicationService.searchAdminApplications(debouncedSearch)
        : await applicationService.getAllAdminApplications()
      setList(Array.isArray(data) ? data : [])
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load shared applications.')
    } finally {
      setLoading(false)
    }
  }, [debouncedSearch])

  useEffect(() => {
    setCurrentPage(1)
  }, [debouncedSearch])

  useEffect(() => {
    if (!superAdmin) {
      setLoading(false)
      return
    }
    void load()
  }, [superAdmin, load])

  const totalPages = Math.max(1, Math.ceil(list.length / PAGE_SIZE))
  useEffect(() => {
    if (currentPage > totalPages) setCurrentPage(totalPages)
  }, [currentPage, totalPages])

  const paginatedApps = useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE
    return list.slice(start, start + PAGE_SIZE)
  }, [list, currentPage])

  if (!superAdmin) {
    return (
      <p className="text-sm text-muted-foreground">
        Managing shared applications is only allowed for super administrators.
      </p>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Shared applications</h1>
          <p className="text-sm text-muted-foreground">Super-admin catalog across all organizations.</p>
        </div>
        <Button onClick={() => { setSelected(null); setFormOpen(true) }} disabled={!canEditApps}>
          <Plus className="mr-2 h-4 w-4" />
          Add application
        </Button>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <div className="relative min-w-[200px] max-w-md flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder="Search by name or package…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

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
          <AdminApplicationsTable
            applications={paginatedApps}
            canEdit={canEditApps}
            onEdit={(app) => { setSelected(app); setFormOpen(true) }}
            onDelete={(app) => setDeleteTarget(app)}
            onTurnCommon={(app) => setTurnTarget(app)}
          />
          {list.length > PAGE_SIZE ? (
            <div className="flex flex-wrap items-center justify-between gap-3">
              <p className="text-sm text-muted-foreground">
                Showing {(currentPage - 1) * PAGE_SIZE + 1}–
                {Math.min(currentPage * PAGE_SIZE, list.length)} of {list.length}
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

      <AlertDialog open={deleteTarget != null} onOpenChange={(v) => !v && setDeleteTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete application?</AlertDialogTitle>
            <AlertDialogDescription>
              This may fail if the application is still referenced. This removes the shared catalog entry after merge workflows.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={async () => {
                if (deleteTarget?.id == null) return
                try {
                  await applicationService.deleteApplication(deleteTarget.id)
                  setDeleteTarget(null)
                  await load()
                } catch (reason: unknown) {
                  setError(reason instanceof Error ? reason.message : 'Delete failed.')
                }
              }}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={turnTarget != null} onOpenChange={(v) => !v && setTurnTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Share application across organizations?</AlertDialogTitle>
            <AlertDialogDescription>
              Merge all copies of package &quot;{turnTarget?.pkg ?? ''}&quot; (&quot;{turnTarget?.name ?? ''}&quot;) into one shared application. This operation cannot be undone from the UI.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={async () => {
                if (turnTarget?.id == null) return
                try {
                  await applicationService.turnApplicationIntoCommon(turnTarget.id)
                  setTurnTarget(null)
                  await load()
                } catch (reason: unknown) {
                  setError(reason instanceof Error ? reason.message : 'Failed to convert application.')
                }
              }}
            >
              Confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
