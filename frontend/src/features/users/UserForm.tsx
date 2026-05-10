import { useEffect, useMemo, useState } from 'react'
import { Loader2, ChevronsUpDown } from 'lucide-react'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/shared/ui/dialog'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/ui/popover'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as groupService from '@/features/groups/groupService'
import * as deviceService from '@/features/devices/deviceService'
import * as userService from '@/features/users/userService'
import type { LookupItem } from '@/features/devices/types'
import type { Role, User, UserPayload } from '@/features/users/types'

const lookupSchema = z.object({
  id: z.number(),
  name: z.string().nullable(),
})

const baseSchema = z.object({
  login: z.string().trim().min(1, 'Login is required'),
  name: z.string().trim().min(1, 'Name is required'),
  email: z.string().regex(/[^@]+@[^.]+\..+/, 'Invalid email address'),
  roleId: z.number({ required_error: 'Role is required' }).int().positive('Role is required'),
  allDevicesAvailable: z.boolean().default(true),
  allConfigAvailable: z.boolean().default(true),
  groups: z.array(lookupSchema),
  configurations: z.array(lookupSchema),
  password: z.string().optional(),
  confirmPassword: z.string().optional(),
})

type UserFormData = z.infer<typeof baseSchema>

export interface UserFormProps {
  mode: 'create' | 'edit'
  initialData: User | null
  onSuccess: () => Promise<void> | void
  onClose: () => void
}

function mapInitialData(initialData: User | null): UserFormData {
  return {
    login: initialData?.login ?? '',
    name: initialData?.name ?? '',
    email: initialData?.email ?? '',
    password: '',
    confirmPassword: '',
    roleId: initialData?.role?.id ?? 0,
    allDevicesAvailable: initialData == null ? true : Boolean(initialData.allDevicesAvailable),
    allConfigAvailable: initialData == null ? true : Boolean(initialData.allConfigAvailable),
    groups: (initialData?.groups ?? []).map((g) => ({ id: g.id, name: g.name ?? null })),
    configurations: (initialData?.configurations ?? []).map((c) => ({ id: c.id, name: c.name ?? null })),
  }
}

export function UserForm({ mode, initialData, onSuccess, onClose }: UserFormProps) {
  const [roles, setRoles] = useState<Role[]>([])
  const [rolesLoading, setRolesLoading] = useState(true)
  const [rolesError, setRolesError] = useState<string | null>(null)
  const [groups, setGroups] = useState<LookupItem[]>([])
  const [configurations, setConfigurations] = useState<{ id: number; name: string | null }[]>([])
  const [optionsLoading, setOptionsLoading] = useState(true)
  const [optionsError, setOptionsError] = useState<string | null>(null)
  const [groupPopoverOpen, setGroupPopoverOpen] = useState(false)
  const [configPopoverOpen, setConfigPopoverOpen] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  const resolver = useMemo(
    () =>
      zodResolver(
        baseSchema.superRefine((value, context) => {
          const pw = value.password?.trim() ?? ''
          const cf = value.confirmPassword?.trim() ?? ''
          if (mode === 'create') {
            if (!pw) {
              context.addIssue({ code: z.ZodIssueCode.custom, path: ['password'], message: 'Password is required' })
            }
            if (pw !== cf) {
              context.addIssue({ code: z.ZodIssueCode.custom, path: ['confirmPassword'], message: 'Passwords do not match' })
            }
            return
          }
          if (pw.length > 0 || cf.length > 0) {
            if (pw !== cf) {
              context.addIssue({ code: z.ZodIssueCode.custom, path: ['confirmPassword'], message: 'Passwords do not match' })
            }
          }
        })
      ),
    [mode]
  )

  const form = useForm<UserFormData>({
    resolver,
    defaultValues: mapInitialData(initialData),
  })

  useEffect(() => {
    form.reset(mapInitialData(initialData))
    setSubmitError(null)
  }, [mode, initialData, form])

  useEffect(() => {
    let cancelled = false
    setRolesLoading(true)
    setRolesError(null)
    void userService
      .getRoles()
      .then((data) => {
        if (cancelled) return
        setRoles(data)
      })
      .catch((error: unknown) => {
        if (cancelled) return
        setRolesError(error instanceof Error ? error.message : 'Failed to load roles.')
      })
      .finally(() => {
        if (cancelled) return
        setRolesLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    let cancelled = false
    setOptionsLoading(true)
    setOptionsError(null)
    void Promise.all([groupService.getGroups(), deviceService.getConfigurations()])
      .then(([g, c]) => {
        if (cancelled) return
        setGroups(g)
        setConfigurations(c.map((row) => ({ id: row.id, name: row.name })))
      })
      .catch((error: unknown) => {
        if (cancelled) return
        setOptionsError(error instanceof Error ? error.message : 'Failed to load groups or configurations.')
      })
      .finally(() => {
        if (cancelled) return
        setOptionsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const selectedGroups = form.watch('groups')
  const selectedConfigurations = form.watch('configurations')
  const allDevicesAvailable = form.watch('allDevicesAvailable')
  const allConfigAvailable = form.watch('allConfigAvailable')

  const groupsButtonLabel = useMemo(() => {
    if (!selectedGroups.length) return 'Select groups…'
    return selectedGroups.map((g) => g.name ?? `#${g.id}`).join(', ')
  }, [selectedGroups])

  const configsButtonLabel = useMemo(() => {
    if (!selectedConfigurations.length) return 'Select configurations…'
    return selectedConfigurations.map((c) => c.name ?? `#${c.id}`).join(', ')
  }, [selectedConfigurations])

  const handleSubmit = form.handleSubmit(async (values) => {
    setSubmitError(null)
    const passwordTrimmed = values.password?.trim() ?? ''
    const basePayload: UserPayload = {
      login: values.login.trim(),
      name: values.name.trim(),
      email: values.email.trim(),
      roleId: values.roleId,
      allDevicesAvailable: Boolean(values.allDevicesAvailable),
      allConfigAvailable: Boolean(values.allConfigAvailable),
      groups: values.groups as LookupItem[],
      configurations: values.configurations as LookupItem[],
    }
    if (passwordTrimmed.length > 0) {
      basePayload.password = passwordTrimmed
    }

    try {
      if (mode === 'create') {
        await userService.createUser({ ...basePayload, password: passwordTrimmed })
      } else {
        if (!initialData) {
          throw new Error('Missing user data for edit mode.')
        }
        await userService.updateUser(initialData.id, basePayload)
      }
      await onSuccess()
      onClose()
    } catch (error: unknown) {
      setSubmitError(error instanceof Error ? error.message : 'Failed to save user.')
    }
  })

  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{mode === 'create' ? 'Add User' : 'Edit User'}</DialogTitle>
          <DialogDescription>
            Matches legacy admin user form: role, device/configuration scope, and password rules aligned with the server.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={(event) => void handleSubmit(event)} className="space-y-4">
            <FormField
              control={form.control}
              name="login"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Login</FormLabel>
                  <FormControl>
                    <Input {...field} readOnly={mode === 'edit'} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input type="email" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="roleId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Role</FormLabel>
                  <FormControl>
                    <Select
                      disabled={rolesLoading || !!rolesError}
                      value={field.value > 0 ? String(field.value) : ''}
                      onValueChange={(value) => field.onChange(Number(value))}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder={rolesLoading ? 'Loading roles...' : 'Select role'} />
                      </SelectTrigger>
                      <SelectContent>
                        {roles.map((role) => (
                          <SelectItem key={role.id} value={String(role.id)}>
                            {role.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                  {rolesError ? <p className="text-sm text-destructive">{rolesError}</p> : null}
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="allDevicesAvailable"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded border p-3">
                  <div>
                    <FormLabel>All devices available</FormLabel>
                    <p className="text-xs text-muted-foreground">If off, restrict this user to selected groups.</p>
                  </div>
                  <FormControl>
                    <Checkbox checked={field.value} onCheckedChange={(checked) => field.onChange(Boolean(checked))} />
                  </FormControl>
                </FormItem>
              )}
            />

            {!allDevicesAvailable ? (
              <div className="space-y-2">
                <Label>Available groups</Label>
                <Popover open={groupPopoverOpen} onOpenChange={setGroupPopoverOpen}>
                  <PopoverTrigger asChild>
                    <Button type="button" variant="outline" className="w-full justify-between" disabled={optionsLoading}>
                      <span className="truncate text-left">{groupsButtonLabel}</span>
                      <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="max-h-64 w-[var(--radix-popover-trigger-width)] overflow-y-auto p-2">
                    <div className="space-y-2">
                      {groups.map((group) => {
                        const checked = selectedGroups.some((selected) => selected.id === group.id)
                        return (
                          <label key={group.id} className="flex cursor-pointer items-center gap-2 rounded px-2 py-1 hover:bg-muted">
                            <Checkbox
                              checked={checked}
                              onCheckedChange={(nextChecked) => {
                                if (nextChecked) {
                                  form.setValue('groups', [...selectedGroups, group], { shouldDirty: true })
                                } else {
                                  form.setValue(
                                    'groups',
                                    selectedGroups.filter((selected) => selected.id !== group.id),
                                    { shouldDirty: true }
                                  )
                                }
                              }}
                            />
                            <span className="text-sm">{group.name ?? `Group #${group.id}`}</span>
                          </label>
                        )
                      })}
                    </div>
                  </PopoverContent>
                </Popover>
              </div>
            ) : null}

            <FormField
              control={form.control}
              name="allConfigAvailable"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded border p-3">
                  <div>
                    <FormLabel>All configurations available</FormLabel>
                    <p className="text-xs text-muted-foreground">If off, restrict this user to selected configurations.</p>
                  </div>
                  <FormControl>
                    <Checkbox checked={field.value} onCheckedChange={(checked) => field.onChange(Boolean(checked))} />
                  </FormControl>
                </FormItem>
              )}
            />

            {!allConfigAvailable ? (
              <div className="space-y-2">
                <Label>Available configurations</Label>
                <Popover open={configPopoverOpen} onOpenChange={setConfigPopoverOpen}>
                  <PopoverTrigger asChild>
                    <Button type="button" variant="outline" className="w-full justify-between" disabled={optionsLoading}>
                      <span className="truncate text-left">{configsButtonLabel}</span>
                      <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="max-h-64 w-[var(--radix-popover-trigger-width)] overflow-y-auto p-2">
                    <div className="space-y-2">
                      {configurations.map((configuration) => {
                        const checked = selectedConfigurations.some((selected) => selected.id === configuration.id)
                        return (
                          <label
                            key={configuration.id}
                            className="flex cursor-pointer items-center gap-2 rounded px-2 py-1 hover:bg-muted"
                          >
                            <Checkbox
                              checked={checked}
                              onCheckedChange={(nextChecked) => {
                                const row: LookupItem = { id: configuration.id, name: configuration.name }
                                if (nextChecked) {
                                  form.setValue('configurations', [...selectedConfigurations, row], { shouldDirty: true })
                                } else {
                                  form.setValue(
                                    'configurations',
                                    selectedConfigurations.filter((selected) => selected.id !== configuration.id),
                                    { shouldDirty: true }
                                  )
                                }
                              }}
                            />
                            <span className="text-sm">{configuration.name ?? `Configuration #${configuration.id}`}</span>
                          </label>
                        )
                      })}
                    </div>
                  </PopoverContent>
                </Popover>
              </div>
            ) : null}

            {optionsError ? <p className="text-sm text-destructive">{optionsError}</p> : null}

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{mode === 'create' ? 'Password' : 'New password (optional)'}</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} value={field.value ?? ''} autoComplete="new-password" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm password</FormLabel>
                  <FormControl>
                    <Input type="password" {...field} value={field.value ?? ''} autoComplete="new-password" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {submitError ? <p className="text-sm text-destructive">{submitError}</p> : null}

            <DialogFooter>
              <Button type="button" variant="outline" onClick={onClose} disabled={form.formState.isSubmitting}>
                Cancel
              </Button>
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Save
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
