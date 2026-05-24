import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { AlertCircle } from 'lucide-react'
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
import * as routeService from '@/features/enrollment-routes/enrollmentRouteService'
import { getOnboardingStatus } from '@/features/onboarding/onboardingService'
import { hasPermission } from '@/features/auth/permissions'

export function EnrollmentRouteListPage() {
  const navigate = useNavigate()
  const [routes, setRoutes] = useState<routeService.EnrollmentRouteListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const canAdd = hasPermission('add_config')

  const loadList = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const list = await routeService.listEnrollmentRoutes()
      setRoutes(Array.isArray(list) ? list : [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load enrollment routes.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadList()
  }, [loadList])

  const handleNewRoute = useCallback(async () => {
    try {
      const status = await getOnboardingStatus()
      if (!status.hasPublishedProfile) {
        navigate('/profiles', {
          state: { onboardingHint: 'Publish a profile before creating an enrollment route.' },
        })
        return
      }
    } catch {
      /* proceed if status check fails */
    }
    navigate('/enrollment-routes/new')
  }, [navigate])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Enrollment routes</h1>
          <p className="text-sm text-muted-foreground">
            مسار التسجيل — binds QR, default folder, and a published profile version.
          </p>
        </div>
        {canAdd ? (
          <Button onClick={() => void handleNewRoute()}>New route</Button>
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
              <TableHead>Name</TableHead>
              <TableHead>Profile</TableHead>
              <TableHead>Version</TableHead>
              <TableHead>Folder</TableHead>
              <TableHead className="w-[100px]" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {routes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-muted-foreground">
                  No enrollment routes yet.
                </TableCell>
              </TableRow>
            ) : (
              routes.map((r) => (
                <TableRow key={r.id}>
                  <TableCell className="font-medium">{r.name}</TableCell>
                  <TableCell>{r.profileId ?? '—'}</TableCell>
                  <TableCell>{r.profileVersionNumber ?? '—'}</TableCell>
                  <TableCell>{r.defaultTreeNodeName || '—'}</TableCell>
                  <TableCell>
                    <Button variant="outline" size="sm" onClick={() => navigate(`/enrollment-routes/${r.id}`)}>
                      Edit
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      )}
    </div>
  )
}
