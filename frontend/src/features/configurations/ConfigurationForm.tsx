import { useEffect, useState } from 'react'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { Loader2 } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import type { Configuration, ConfigurationPayload } from '@/features/configurations/types'
import * as configurationService from '@/features/configurations/configurationService'

const formSchema = z.object({
  name: z.string().refine((s) => s.trim().length > 0, 'Name is required'),
  type: z.enum(['WORK', 'COMMON']),
  description: z.string(),
})

export type ConfigurationFormValues = z.infer<typeof formSchema>

export interface ConfigurationFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  mode: 'create' | 'edit'
  initialData: Configuration | null
  onSuccess: () => void
}

export function ConfigurationForm({
  open,
  onOpenChange,
  mode,
  initialData,
  onSuccess,
}: ConfigurationFormProps) {
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  const form = useForm<ConfigurationFormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      type: 'WORK',
      description: '',
    },
  })

  useEffect(() => {
    if (!open) return
    setSubmitError(null)
    if (mode === 'edit' && initialData?.id != null) {
      form.reset({
        name: initialData.name?.trim() ?? '',
        type: configurationService.typeToConfigurationKind(initialData.type),
        description: initialData.description ?? '',
      })
    } else {
      form.reset({
        name: '',
        type: 'WORK',
        description: '',
      })
    }
  }, [open, mode, initialData, form])

  const handleClose = () => {
    if (!submitting) {
      setSubmitError(null)
      onOpenChange(false)
    }
  }

  const onSubmit = async (values: ConfigurationFormValues) => {
    const descTrim = values.description.trim()
    const payload: ConfigurationPayload = {
      name: values.name.trim(),
      type: values.type,
      description: descTrim === '' ? null : descTrim,
    }
    setSubmitError(null)
    setSubmitting(true)
    try {
      if (mode === 'create') {
        await configurationService.createConfiguration(payload)
      } else {
        const id = initialData?.id
        if (id == null) throw new Error('Missing configuration id.')
        await configurationService.updateConfiguration(id, payload)
      }
      onSuccess()
      onOpenChange(false)
    } catch (e) {
      setSubmitError(e instanceof Error ? e.message : 'Save failed.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(o) => !o && handleClose()}>
      <DialogContent className="sm:max-w-md" onPointerDownOutside={(e) => submitting && e.preventDefault()}>
        <DialogHeader>
          <DialogTitle>
            {mode === 'create' ? 'New configuration' : 'Edit configuration'}
          </DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Configuration name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Type</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="WORK">Work (device)</SelectItem>
                      <SelectItem value="COMMON">Common (typical)</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea placeholder="Optional" rows={3} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {submitError ?
              <p className="text-destructive text-sm" role="alert">
                {submitError}
              </p>
            : null}

            <DialogFooter className="gap-2 sm:gap-0">
              <Button type="button" variant="outline" disabled={submitting} onClick={handleClose}>
                Cancel
              </Button>
              <Button type="submit" disabled={submitting}>
                {submitting ?
                  <Loader2 className="h-4 w-4 animate-spin" aria-hidden />
                : null}
                Save
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
