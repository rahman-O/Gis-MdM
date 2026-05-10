import { useEffect, useState } from 'react'
import { AlertCircle, Pencil, Trash2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Skeleton } from '@/shared/ui/skeleton'
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
import * as roleService from '@/features/roles/roleService'
import type { ManagedRole, RolePermission } from '@/features/roles/types'

type FormMode = 'create' | 'edit' | null

export function RolesPage() {
  const [roles, setRoles] = useState<ManagedRole[]>([])
  const [permissions, setPermissions] = useState<RolePermission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [formMode, setFormMode] = useState<FormMode>(null)
  const [roleToEdit, setRoleToEdit] = useState<ManagedRole | null>(null)
  const [roleToDelete, setRoleToDelete] = useState<ManagedRole | null>(null)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set())
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [deleteError, setDeleteError] = useState<string | null>(null)

  const loadData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [rolesList, perms] = await Promise.all([roleService.getManagedRoles(), roleService.getRolePermissions()])
      setRoles(rolesList)
      setPermissions(perms)
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load roles.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadData()
  }, [])

  const openCreate = () => {
    setFormMode('create')
    setRoleToEdit(null)
    setName('')
    setDescription('')
    setSelectedPermissionIds(new Set())
    setFormError(null)
  }

  const openEdit = (role: ManagedRole) => {
    setFormMode('edit')
    setRoleToEdit(role)
    setName(role.name ?? '')
    setDescription(role.description ?? '')
    setSelectedPermissionIds(new Set((role.permissions ?? []).map((p) => p.id)))
    setFormError(null)
  }

  const togglePermission = (id: number, checked: boolean) => {
    setSelectedPermissionIds((prev) => {
      const next = new Set(prev)
      if (checked) next.add(id)
      else next.delete(id)
      return next
    })
  }

  const handleSave = async () => {
    if (!name.trim()) {
      setFormError('Role name is required.')
      return
    }
    setSubmitting(true)
    setFormError(null)
    try {
      await roleService.saveRole({
        id: roleToEdit?.id ?? null,
        name: name.trim(),
        description: description.trim() || undefined,
        permissions: [...selectedPermissionIds].map((id) => ({ id })),
      })
      setFormMode(null)
      await loadData()
    } catch (reason: unknown) {
      setFormError(reason instanceof Error ? reason.message : 'Failed to save role.')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async () => {
    if (!roleToDelete?.id) return
    setDeleting(true)
    setDeleteError(null)
    try {
      await roleService.deleteRole(roleToDelete.id)
      setRoleToDelete(null)
      await loadData()
    } catch (reason: unknown) {
      setDeleteError(reason instanceof Error ? reason.message : 'Failed to delete role.')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Roles</h1>
          <p className="text-sm text-muted-foreground">
            Manage permission bundles assigned to administrators (parity with legacy Settings → Roles).
          </p>
        </div>
        <Button type="button" onClick={openCreate}>
          Add role
        </Button>
      </div>

      {error ? (
        <div className="flex items-center gap-3 rounded border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button variant="outline" size="sm" onClick={() => void loadData()}>
            Retry
          </Button>
        </div>
      ) : null}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Role name</TableHead>
              <TableHead className="w-[120px] text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell colSpan={2}>
                    <Skeleton className="h-8 w-full" />
                  </TableCell>
                </TableRow>
              ))
            ) : roles.length === 0 ? (
              <TableRow>
                <TableCell colSpan={2} className="h-24 text-center text-muted-foreground">
                  No roles found
                </TableCell>
              </TableRow>
            ) : (
              roles.map((role) => (
                <TableRow key={role.id != null ? String(role.id) : role.name}>
                  <TableCell>{role.name}</TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="ghost" size="icon" type="button" onClick={() => openEdit(role)} aria-label={`Edit role ${role.name}`}>
                        <Pencil className="h-4 w-4" />
                      </Button>
                      {!role.superAdmin ? (
                        <Button
                          variant="ghost"
                          size="icon"
                          type="button"
                          onClick={() => setRoleToDelete(role)}
                          aria-label={`Delete role ${role.name}`}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      ) : null}
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog open={formMode != null} onOpenChange={(open) => !open && setFormMode(null)}>
        <DialogContent className="max-h-[85vh] overflow-y-auto sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>{formMode === 'create' ? 'Add role' : 'Edit role'}</DialogTitle>
            <DialogDescription>Set the role name, optional description, and permission checkboxes matching the backend role model.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="role-name">Name</Label>
              <Input id="role-name" value={name} onChange={(e) => setName(e.target.value)} disabled={submitting || roleToEdit?.superAdmin} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="role-desc">Description</Label>
              <Input id="role-desc" value={description} onChange={(e) => setDescription(e.target.value)} disabled={submitting || roleToEdit?.superAdmin} />
            </div>
            <div className="space-y-2">
              <Label>Permissions</Label>
              <div className="max-h-52 space-y-2 overflow-y-auto rounded border p-3">
                {permissions.length === 0 ? (
                  <p className="text-sm text-muted-foreground">No permissions available.</p>
                ) : (
                  permissions.map((p) => (
                    <label key={p.id} className="flex cursor-pointer items-center gap-2 text-sm">
                      <Checkbox
                        checked={selectedPermissionIds.has(p.id)}
                        disabled={submitting || roleToEdit?.superAdmin}
                        onCheckedChange={(c) => togglePermission(p.id, Boolean(c))}
                      />
                      <span>{p.name ?? `Permission ${p.id}`}</span>
                    </label>
                  ))
                )}
              </div>
            </div>
            {formError ? <p className="text-sm text-destructive">{formError}</p> : null}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setFormMode(null)} disabled={submitting}>
              Cancel
            </Button>
            {!roleToEdit?.superAdmin ? (
              <Button type="button" onClick={() => void handleSave()} disabled={submitting}>
                Save
              </Button>
            ) : null}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={roleToDelete != null}
        onOpenChange={(open) => {
          if (!open) {
            setDeleteError(null)
            setRoleToDelete(null)
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete role?</AlertDialogTitle>
            <AlertDialogDescription>This will permanently remove the role "{roleToDelete?.name}".</AlertDialogDescription>
          </AlertDialogHeader>
          {deleteError ? <p className="text-sm text-destructive">{deleteError}</p> : null}
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleting}>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={() => void handleDelete()} disabled={deleting}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
