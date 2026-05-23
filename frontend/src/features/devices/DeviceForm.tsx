import { useEffect, useMemo, useState } from 'react'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Loader2, ChevronsUpDown } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { Label } from '@/shared/ui/label'
import { Checkbox } from '@/shared/ui/checkbox'
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/ui/popover'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import * as deviceService from '@/features/devices/deviceService'
import type { ConfigurationOption, DevicePayload, DeviceView, LookupItem } from '@/features/devices/types'

const formSchema = z.object({
  number: z
    .string()
    .trim()
    .min(1, 'Device number is required')
    .refine((value) => !/[/?&]/.test(value), 'Device number cannot include /, ? or &'),
  description: z.string().nullable(),
  configurationId: z.number().nullable(),
  groups: z.array(z.object({ id: z.number(), name: z.string().nullable() })),
  imei: z
    .string()
    .nullable()
    .refine((value) => value == null || value === '' || /^\d{15}$/.test(value), 'IMEI must be exactly 15 digits'),
  phone: z.string().nullable(),
  custom1: z.string().nullable(),
  custom2: z.string().nullable(),
  custom3: z.string().nullable(),
})

type DeviceFormData = z.infer<typeof formSchema>

export interface DeviceFormProps {
  mode: 'create' | 'edit'
  initialData: DeviceView | null
  onSuccess: () => Promise<void> | void
  onClose: () => void
}

function mapInitialData(device: DeviceView | null): DeviceFormData {
  return {
    number: device?.number ?? '',
    description: device?.description ?? null,
    configurationId: device?.configurationId ?? null,
    groups: device?.groups ?? [],
    imei: device?.imei ?? null,
    phone: device?.phone ?? null,
    custom1: device?.custom1 ?? null,
    custom2: device?.custom2 ?? null,
    custom3: device?.custom3 ?? null,
  }
}

export function DeviceForm({ mode, initialData, onSuccess, onClose }: DeviceFormProps) {
  const [groups, setGroups] = useState<LookupItem[]>([])
  const [configurations, setConfigurations] = useState<ConfigurationOption[]>([])
  const [groupsLoading, setGroupsLoading] = useState(true)
  const [configurationsLoading, setConfigurationsLoading] = useState(true)
  const [groupsError, setGroupsError] = useState<string | null>(null)
  const [configurationsError, setConfigurationsError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [groupSelectorOpen, setGroupSelectorOpen] = useState(false)

  const form = useForm<DeviceFormData>({
    resolver: zodResolver(formSchema),
    defaultValues: mapInitialData(initialData),
  })

  useEffect(() => {
    form.reset(mapInitialData(initialData))
  }, [initialData, form])

  useEffect(() => {
    let cancelled = false
    setGroupsLoading(true)
    setConfigurationsLoading(true)
    void Promise.all([deviceService.getGroups(), deviceService.getConfigurations()])
      .then(([loadedGroups, loadedConfigurations]) => {
        if (cancelled) return
        setGroups(loadedGroups)
        setConfigurations(loadedConfigurations)
      })
      .catch((error: unknown) => {
        if (cancelled) return
        const message = error instanceof Error ? error.message : 'Failed to load form options.'
        setGroupsError(message)
        setConfigurationsError(message)
      })
      .finally(() => {
        if (cancelled) return
        setGroupsLoading(false)
        setConfigurationsLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const selectedGroups = form.watch('groups')
  const groupsLabel = useMemo(() => {
    if (!selectedGroups.length) return 'Select groups...'
    return selectedGroups.map((group) => group.name ?? `#${group.id}`).join(', ')
  }, [selectedGroups])

  const onSubmit = form.handleSubmit(async (values) => {
    setSubmitError(null)
    const payload: DevicePayload = {
      number: values.number.trim(),
      description: values.description?.trim() ? values.description.trim() : null,
      configurationId: values.configurationId,
      groups: values.groups,
      imei: values.imei?.trim() ? values.imei.trim() : null,
      phone: values.phone?.trim() ? values.phone.trim() : null,
      custom1: values.custom1?.trim() ? values.custom1.trim() : null,
      custom2: values.custom2?.trim() ? values.custom2.trim() : null,
      custom3: values.custom3?.trim() ? values.custom3.trim() : null,
      oldNumber: initialData?.oldNumber ?? null,
    }

    try {
      if (mode === 'create') {
        await deviceService.createDevice(payload)
      } else {
        if (!initialData?.id) {
          throw new Error('Invalid device ID for update.')
        }
        await deviceService.updateDevice({ ...payload, id: initialData.id })
      }
      await onSuccess()
      onClose()
    } catch (error: unknown) {
      setSubmitError(error instanceof Error ? error.message : 'Failed to save device.')
    }
  })

  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{mode === 'create' ? 'Add Device' : 'Edit Device'}</DialogTitle>
          <DialogDescription>Set the main enrollment and metadata fields for this device.</DialogDescription>
        </DialogHeader>

        <form onSubmit={(event) => void onSubmit(event)} className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="number">Device Number</Label>
              <Input id="number" readOnly={mode === 'edit'} {...form.register('number')} />
              {form.formState.errors.number ? (
                <p className="text-sm text-destructive">{form.formState.errors.number.message}</p>
              ) : null}
            </div>

            <div className="space-y-2">
              <Label>Configuration</Label>
              <Select
                value={form.watch('configurationId') == null ? 'none' : String(form.watch('configurationId'))}
                onValueChange={(value) => form.setValue('configurationId', value === 'none' ? null : Number(value))}
                disabled={configurationsLoading}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select configuration" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">None</SelectItem>
                  {configurations.map((item) => (
                    <SelectItem key={item.id} value={String(item.id)}>
                      {item.name ?? `Configuration #${item.id}`}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {configurationsError ? <p className="text-sm text-destructive">{configurationsError}</p> : null}
            </div>
          </div>

          <div className="space-y-2">
            <Label>Description</Label>
            <Textarea rows={3} {...form.register('description')} />
          </div>

          <div className="space-y-2">
            <Label>Groups</Label>
            <Popover open={groupSelectorOpen} onOpenChange={setGroupSelectorOpen}>
              <PopoverTrigger asChild>
                <Button type="button" variant="outline" className="w-full justify-between" disabled={groupsLoading}>
                  <span className="truncate text-left">{groupsLabel}</span>
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
            {groupsError ? <p className="text-sm text-destructive">{groupsError}</p> : null}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label>IMEI</Label>
              <Input {...form.register('imei')} />
              {form.formState.errors.imei ? <p className="text-sm text-destructive">{form.formState.errors.imei.message}</p> : null}
            </div>
            <div className="space-y-2">
              <Label>Phone</Label>
              <Input {...form.register('phone')} />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <div className="space-y-2">
              <Label>Custom 1</Label>
              <Input {...form.register('custom1')} />
            </div>
            <div className="space-y-2">
              <Label>Custom 2</Label>
              <Input {...form.register('custom2')} />
            </div>
            <div className="space-y-2">
              <Label>Custom 3</Label>
              <Input {...form.register('custom3')} />
            </div>
          </div>

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
      </DialogContent>
    </Dialog>
  )
}
