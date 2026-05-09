import { useEffect, useState } from 'react'
import { AlertCircle, Pencil, Trash2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Skeleton } from '@/shared/ui/skeleton'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
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
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import * as groupService from '@/features/groups/groupService'
import type { LookupItem } from '@/features/devices/types'

type FormMode = 'create' | 'edit' | null

export function GroupsPage() {
  const [groups, setGroups] = useState<LookupItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formMode, setFormMode] = useState<FormMode>(null)
  const [groupToEdit, setGroupToEdit] = useState<LookupItem | null>(null)
  const [groupToDelete, setGroupToDelete] = useState<LookupItem | null>(null)
  const [name, setName] = useState('')
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const loadGroups = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await groupService.getGroups()
      setGroups(response)
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load groups.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadGroups()
  }, [])

  const openCreate = () => {
    setFormMode('create')
    setGroupToEdit(null)
    setName('')
    setFormError(null)
  }

  const openEdit = (group: LookupItem) => {
    setFormMode('edit')
    setGroupToEdit(group)
    setName(group.name ?? '')
    setFormError(null)
  }

  const handleSave = async () => {
    if (!name.trim()) {
      setFormError('Group name is required.')
      return
    }
    setSubmitting(true)
    setFormError(null)
    try {
      if (formMode === 'create') {
        await groupService.createGroup(name.trim())
      } else if (groupToEdit) {
        await groupService.updateGroup({ id: groupToEdit.id, name: name.trim() })
      }
      setFormMode(null)
      await loadGroups()
    } catch (reason: unknown) {
      setFormError(reason instanceof Error ? reason.message : 'Failed to save group.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Groups</h1>
          <p className="text-sm text-muted-foreground">Manage device grouping used for filtering and policy assignment.</p>
        </div>
        <Button type="button" onClick={openCreate}>
          Add Group
        </Button>
      </div>

      {error ? (
        <div className="flex items-center gap-3 rounded border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button variant="outline" size="sm" onClick={() => void loadGroups()}>
            Retry
          </Button>
        </div>
      ) : null}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Group Name</TableHead>
              <TableHead className="w-[120px] text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              Array.from({ length: 4 }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell colSpan={2}>
                    <Skeleton className="h-8 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : groups.length === 0 ? (
              <TableRow>
                <TableCell colSpan={2} className="h-24 text-center text-muted-foreground">
                  No groups found
                </TableCell>
              </TableRow>
            ) : (
              groups.map((group) => (
                <TableRow key={group.id}>
                  <TableCell>{group.name ?? `Group #${group.id}`}</TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="icon" onClick={() => openEdit(group)}>
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button variant="ghost" size="icon" onClick={() => setGroupToDelete(group)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={formMode != null} onOpenChange={(open) => !open && setFormMode(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{formMode === 'create' ? 'Add Group' : 'Edit Group'}</DialogTitle>
          </DialogHeader>
          <div className="space-y-2">
            <Input value={name} onChange={(event) => setName(event.target.value)} placeholder="Group name" />
            {formError ? <p className="text-sm text-destructive">{formError}</p> : null}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setFormMode(null)} disabled={submitting}>
              Cancel
            </Button>
            <Button onClick={() => void handleSave()} disabled={submitting}>
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog open={groupToDelete != null} onOpenChange={(open) => !open && setGroupToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete group?</AlertDialogTitle>
            <AlertDialogDescription>
              This action will remove <span className="font-medium text-foreground">{groupToDelete?.name ?? ''}</span>.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (!groupToDelete) return
                void groupService.deleteGroup(groupToDelete.id).then(loadGroups)
                setGroupToDelete(null)
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
