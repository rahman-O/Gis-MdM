import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { AlertCircle, Plus } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
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
import * as applicationService from '@/features/applications/services/applicationService'
import { ApplicationVersionDialog } from '@/features/applications/components/ApplicationVersionDialog'
import { ApplicationVersionConfigurationsDialog } from '@/features/applications/components/ApplicationVersionConfigurationsDialog'
import type { Application, ApplicationVersion } from '@/features/applications/model/types'

export function ApplicationVersionsPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const appId = Number(id)
  const [application, setApplication] = useState<Application | null>(null)
  const [versions, setVersions] = useState<ApplicationVersion[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [selected, setSelected] = useState<ApplicationVersion | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<ApplicationVersion | null>(null)
  const [cfgVersionId, setCfgVersionId] = useState<number | null>(null)

  const load = async () => {
    if (!Number.isFinite(appId) || appId <= 0) return
    setLoading(true)
    setError(null)
    try {
      const [app, list] = await Promise.all([
        applicationService.getApplication(appId),
        applicationService.getApplicationVersions(appId),
      ])
      setApplication(app)
      setVersions(list ?? [])
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load versions.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { void load() }, [appId])

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <Button variant="ghost" onClick={() => navigate('/applications')}>Back</Button>
          <h1 className="text-2xl font-semibold tracking-tight">
            Versions {application?.name ? `— ${application.name}` : ''}
          </h1>
        </div>
        <Button onClick={() => { setSelected(null); setDialogOpen(true) }}>
          <Plus className="mr-2 h-4 w-4" />
          Add version
        </Button>
      </div>

      {error ? (
        <div className="flex items-center gap-2 rounded border border-destructive/50 bg-destructive/10 p-3 text-sm">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button variant="outline" size="sm" onClick={() => void load()}>Retry</Button>
        </div>
      ) : null}

      {loading ? (
        <div className="space-y-2">{Array.from({ length: 4 }).map((_, i) => <Skeleton key={i} className="h-9 w-full" />)}</div>
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Version</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>URL</TableHead>
                <TableHead>Split</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {versions.length === 0 ? (
                <TableRow><TableCell colSpan={5} className="h-24 text-center text-muted-foreground">No versions found.</TableCell></TableRow>
              ) : versions.map((v) => (
                <TableRow key={v.id ?? `${v.version}-${v.versionCode}`}>
                  <TableCell>{v.version ?? '—'}</TableCell>
                  <TableCell>{v.versionCode ?? '—'}</TableCell>
                  <TableCell className="max-w-[260px] truncate">{v.url ?? '—'}</TableCell>
                  <TableCell>{v.split ? 'Yes' : 'No'}</TableCell>
                  <TableCell className="text-right">
                    <div className="inline-flex gap-2">
                      <Button size="sm" variant="outline" onClick={() => setCfgVersionId(v.id ?? null)}>Configurations</Button>
                      <Button size="sm" variant="outline" onClick={() => { setSelected(v); setDialogOpen(true) }}>Edit</Button>
                      <Button size="sm" variant="destructive" onClick={() => setDeleteTarget(v)}>Delete</Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      <ApplicationVersionDialog
        open={dialogOpen}
        initialData={selected}
        onClose={() => setDialogOpen(false)}
        onSave={async (payload) => {
          await applicationService.createOrUpdateApplicationVersion({
            ...payload,
            applicationId: appId,
          })
          await load()
        }}
      />

      <ApplicationVersionConfigurationsDialog
        open={cfgVersionId != null}
        versionId={cfgVersionId}
        onClose={() => setCfgVersionId(null)}
      />

      <AlertDialog open={deleteTarget != null} onOpenChange={(v) => !v && setDeleteTarget(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete version?</AlertDialogTitle>
            <AlertDialogDescription>
              This action can fail if the version is referenced by configurations.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={async () => {
                if (deleteTarget?.id == null) return
                await applicationService.deleteApplicationVersion(deleteTarget.id)
                setDeleteTarget(null)
                await load()
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
