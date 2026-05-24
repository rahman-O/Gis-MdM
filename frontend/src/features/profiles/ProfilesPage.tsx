import { useCallback, useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom'
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
import * as profileService from '@/features/profiles/profileService'
import { ProfileForm } from '@/features/profiles/ProfileForm'
import { ProfileHealthBadge } from '@/features/profiles/ProfileHealthBadge'
import { ProfileListBadges } from '@/features/profiles/ProfileListBadges'
import { ProfileWorkspace } from '@/features/profiles/workspace/ProfileWorkspace'
import {
  ProfileWorkspaceProvider,
  useProfileWorkspace,
} from '@/features/profiles/workspace/profileWorkspaceState'
import { hasPermission } from '@/features/auth/permissions'
import type { ProfileListItem } from '@/features/profiles/types'

function ProfilesPageInner() {
  const { open } = useProfileWorkspace()
  const location = useLocation()
  const [profiles, setProfiles] = useState<ProfileListItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formOpen, setFormOpen] = useState(false)
  const canAdd = hasPermission('add_config')

  const loadList = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const list = await profileService.listProfiles()
      setProfiles(Array.isArray(list) ? list : [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load profiles.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadList()
  }, [loadList])

  const hint = (location.state as { onboardingHint?: string } | null)?.onboardingHint

  const openWorkspace = (p: ProfileListItem) => {
    open(p.id, 'overview')
  }

  const handleCreate = async (payload: { name: string; description?: string | null }) => {
    const meta = await profileService.createProfile(payload)
    open(meta.id, 'assignments')
    void loadList()
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Profiles</h1>
          <p className="text-sm text-muted-foreground">
            Versioned device policy. Assign published versions to tree folders; enrollment routes no
            longer bind profiles.
          </p>
        </div>
        {canAdd ? (
          <Button onClick={() => setFormOpen(true)}>New profile</Button>
        ) : null}
      </div>

      {hint ? (
        <p className="rounded-md border border-amber-500/40 bg-amber-50/80 px-3 py-2 text-sm text-amber-900 dark:bg-amber-950/30 dark:text-amber-100">
          {hint}
        </p>
      ) : null}

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
              <TableHead>Health</TableHead>
              <TableHead>Signals</TableHead>
              <TableHead>Published</TableHead>
              <TableHead>Assignments</TableHead>
              <TableHead className="w-[100px]" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {profiles.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-muted-foreground">
                  No profiles yet.
                </TableCell>
              </TableRow>
            ) : (
              profiles.map((p) => (
                <TableRow
                  key={p.id}
                  className="cursor-pointer"
                  onClick={() => openWorkspace(p)}
                >
                  <TableCell className="font-medium">{p.name}</TableCell>
                  <TableCell>
                    {p.health ? (
                      <ProfileHealthBadge health={p.health} />
                    ) : (
                      <span className="text-muted-foreground">—</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <ProfileListBadges badges={p.badges} />
                  </TableCell>
                  <TableCell>{p.publishedVersion ?? '—'}</TableCell>
                  <TableCell>{p.assignmentCount ?? 0}</TableCell>
                  <TableCell onClick={(e) => e.stopPropagation()}>
                    <Button variant="outline" size="sm" onClick={() => openWorkspace(p)}>
                      Open
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      )}

      <ProfileForm open={formOpen} onOpenChange={setFormOpen} onSubmit={handleCreate} />
      <ProfileWorkspace />
    </div>
  )
}

export function ProfilesPage() {
  return (
    <ProfileWorkspaceProvider>
      <ProfilesPageInner />
    </ProfileWorkspaceProvider>
  )
}
