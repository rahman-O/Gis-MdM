import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { AlertCircle } from 'lucide-react'
import { EnrollmentRouteDialog } from '@/features/enrollment-routes/EnrollmentRouteDialog'
import type { EnrollmentRouteDialogStateId } from '@/features/enrollment-routes/enrollmentRouteDialogState'
import * as routeService from '@/features/enrollment-routes/enrollmentRouteService'
import { hasPermission } from '@/features/auth/permissions'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/ui/table'

export function EnrollmentRouteListPage() {
  const { t } = useTranslation()
  const [routes, setRoutes] = useState<routeService.EnrollmentRouteView[]>([])
  const [treeNodes, setTreeNodes] = useState<routeService.TreeNodeOption[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [dialogState, setDialogState] = useState<EnrollmentRouteDialogStateId>('LIST')
  const [dialogRouteId, setDialogRouteId] = useState(0)
  const canAdd = hasPermission('add_config')

  const loadList = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [list, nodes] = await Promise.all([
        routeService.listEnrollmentRoutes(),
        routeService.listTreeNodeOptions(),
      ])
      setRoutes(Array.isArray(list) ? list : [])
      setTreeNodes(Array.isArray(nodes) ? nodes : [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load enrollment routes.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadList()
  }, [loadList])

  const openDialog = (state: EnrollmentRouteDialogStateId, routeId = 0) => {
    setDialogState(state)
    setDialogRouteId(routeId)
    setDialogOpen(true)
  }

  const closeDialog = () => {
    setDialogOpen(false)
    setDialogState('LIST')
    setDialogRouteId(0)
  }

  const handleDialogStateChange = (state: EnrollmentRouteDialogStateId, routeId?: number) => {
    setDialogState(state)
    if (routeId !== undefined) setDialogRouteId(routeId)
  }

  const formatNodePath = (path?: string) => {
    if (!path) return '—'
    const segments = path.split('/').filter(Boolean)
    const resolvedSegments = segments.map((segment) => {
      const id = Number(segment)
      if (!isNaN(id)) {
        const node = treeNodes.find((n) => n.id === id)
        return node ? node.name : segment
      }
      return segment
    })
    return resolvedSegments.join(' / ')
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            {t('nav.enrollmentRoutes')}
          </h1>
          <p className="text-sm text-muted-foreground">{t('enrollmentRoute.list.subtitle')}</p>
        </div>
        {canAdd ? (
          <Button onClick={() => openDialog('DIALOG_CREATE')}>{t('enrollmentRoute.actions.new')}</Button>
        ) : null}
      </div>

      {error ? (
        <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4" />
          <span>{error}</span>
        </div>
      ) : null}

      {loading ? (
        <Skeleton className="h-48 w-full" />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('enrollmentRoute.list.name')}</TableHead>
              <TableHead>{t('enrollmentRoute.list.folder')}</TableHead>
              <TableHead>{t('enrollmentRoute.list.bootstrap')}</TableHead>
              <TableHead>{t('enrollmentRoute.list.status')}</TableHead>
              <TableHead className="w-[100px]" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {routes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-muted-foreground">
                  {t('enrollmentRoute.list.empty')}
                </TableCell>
              </TableRow>
            ) : (
              routes.map((r) => (
                <TableRow key={r.id}>
                  <TableCell className="font-medium">{r.name}</TableCell>
                  <TableCell>{formatNodePath(r.targetNodePath || r.targetNodeName)}</TableCell>
                  <TableCell>
                    {r.bootstrapApplicationName || '—'}
                    {r.resolvedVersionLabel ? ` (${r.resolvedVersionLabel})` : ''}
                  </TableCell>
                  <TableCell>{r.status}</TableCell>
                  <TableCell>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => openDialog('DIALOG_OVERVIEW', r.id)}
                    >
                      {t('enrollmentRoute.actions.open')}
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      )}

      <EnrollmentRouteDialog
        open={dialogOpen}
        state={dialogState}
        routeId={dialogRouteId}
        onStateChange={handleDialogStateChange}
        onClose={closeDialog}
        onSaved={() => void loadList()}
      />
    </div>
  )
}
