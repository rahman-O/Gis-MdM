import { useCallback, useEffect, useMemo, useState } from 'react'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { AlertCircle, Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Skeleton } from '@/shared/ui/skeleton'
import { Switch } from '@/shared/ui/switch'
import { useToast } from '@/shared/hooks/use-toast'
import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { Configuration } from '@/features/configurations/types'
import { LANGUAGE_OPTIONS } from '@/features/settings/languageMaps'
import * as settingsService from '@/features/settings/settingsService'
import type { ConfigurationOption, Settings, SettingsPayload } from '@/features/settings/types'

/** Matches backend `Settings.passwordStrength` and legacy `password.service.js` (0–2). */
const PASSWORD_STRENGTH_OPTIONS = [
  { value: '0', label: 'Any (length only)' },
  { value: '1', label: 'Digits, upper & lower case' },
  { value: '2', label: 'Above + special characters' },
]

const NONE_CONFIG_VALUE = '__none__'

const settingsSchema = z
  .object({
    customerName: z.string().min(1, 'Customer name is required'),
    createNewDevices: z.boolean(),
    newDeviceConfigurationId: z.number().nullable(),
    language: z.string().min(1, 'Language is required'),
    passwordLength: z.coerce.number().int().min(1, 'Must be at least 1'),
    passwordStrength: z.coerce.number().int().min(0).max(2),
    sendDeviceInfoExpiryDays: z.coerce.number().int().min(1, 'Must be at least 1'),
    unsecureEnrollment: z.boolean(),
    deviceFastSearch: z.boolean(),
  })
  .refine((v) => !v.createNewDevices || v.newDeviceConfigurationId != null, {
    message:
      'Select a default configuration for new devices, or turn off "Create new devices on first access".',
    path: ['newDeviceConfigurationId'],
  })

type FormValues = z.infer<typeof settingsSchema>

function toFormValues(s: Settings): FormValues {
  const strength = s.passwordStrength ?? 0
  return {
    customerName: s.customerName ?? '',
    createNewDevices: s.createNewDevices ?? false,
    newDeviceConfigurationId: s.newDeviceConfigurationId ?? null,
    language: s.language ?? 'en',
    passwordLength: s.passwordLength ?? 0,
    passwordStrength: Math.min(2, Math.max(0, strength)),
    sendDeviceInfoExpiryDays: s.sendDeviceInfoExpiryDays ?? 0,
    unsecureEnrollment: s.unsecureEnrollment ?? false,
    deviceFastSearch: s.deviceFastSearch ?? false,
  }
}

function toPayload(v: FormValues): SettingsPayload {
  return {
    customerName: v.customerName.trim(),
    createNewDevices: v.createNewDevices,
    newDeviceConfigurationId: v.newDeviceConfigurationId,
    language: v.language,
    passwordLength: v.passwordLength,
    passwordStrength: v.passwordStrength,
    sendDeviceInfoExpiryDays: v.sendDeviceInfoExpiryDays,
    unsecureEnrollment: v.unsecureEnrollment,
    deviceFastSearch: v.deviceFastSearch,
  }
}

function mapConfigurations(list: Configuration[]): ConfigurationOption[] {
  return list
    .filter((c): c is Configuration & { id: number } => c.id != null && Number.isFinite(c.id))
    .map((c) => ({
      id: c.id as number,
      name: c.name?.trim() ? (c.name as string) : `Configuration #${c.id}`,
    }))
}

export function SettingsPage() {
  const { toast } = useToast()
  const [configurations, setConfigurations] = useState<ConfigurationOption[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const form = useForm<FormValues>({
    resolver: zodResolver(settingsSchema),
    defaultValues: {
      customerName: '',
      createNewDevices: false,
      newDeviceConfigurationId: null,
      language: 'en',
      passwordLength: 1,
      passwordStrength: 0,
      sendDeviceInfoExpiryDays: 1,
      unsecureEnrollment: false,
      deviceFastSearch: false,
    },
  })

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const settled = await Promise.allSettled([
        settingsService.getSettings(),
        apiClient.get<HmdmEnvelope<Configuration[]>>('/private/configurations/search'),
      ])

      if (settled[0].status === 'rejected') {
        const reason = settled[0].reason
        throw reason instanceof Error ? reason : new Error('Failed to load settings.')
      }
      const nextSettings = settled[0].value as Settings

      let configs: ConfigurationOption[] = []
      if (settled[1].status === 'fulfilled') {
        try {
          const data = unwrapHmdmData(settled[1].value.data, 'Failed to load configurations.')
          configs = mapConfigurations(Array.isArray(data) ? data : [])
        } catch {
          configs = []
        }
      }
      setConfigurations(configs)
      form.reset(toFormValues(nextSettings))
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load settings.')
    } finally {
      setLoading(false)
    }
  }, [form])

  useEffect(() => {
    void load()
  }, [load])

  const handleSave = form.handleSubmit(async (values) => {
    setSubmitting(true)
    try {
      const payload = toPayload(values)
      const updated = await settingsService.updateSettings(payload)
      form.reset(toFormValues(updated))
      toast({ title: 'Settings saved' })
    } catch (reason: unknown) {
      toast({
        title: 'Failed to save settings',
        variant: 'destructive',
        description: reason instanceof Error ? reason.message : undefined,
      })
    } finally {
      setSubmitting(false)
    }
  })

  const skeletonBlock = useMemo(
    () => (
      <div className="max-w-2xl space-y-4">
        <Skeleton className="h-9 w-full" />
        <Skeleton className="h-9 w-full" />
        <Skeleton className="h-9 w-full" />
        <Skeleton className="h-9 w-2/3" />
      </div>
    ),
    []
  )

  if (loading) {
    return (
      <div className="space-y-5">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Settings</h1>
          <p className="text-sm text-muted-foreground">Instance-wide preferences and security defaults.</p>
        </div>
        {skeletonBlock}
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-5">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Settings</h1>
          <p className="text-sm text-muted-foreground">Instance-wide preferences and security defaults.</p>
        </div>
        <div className="flex items-center gap-3 rounded border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm">
          <AlertCircle className="h-4 w-4 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button variant="outline" size="sm" onClick={() => void load()}>
            Retry
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Settings</h1>
        <p className="text-sm text-muted-foreground">Instance-wide preferences and security defaults.</p>
      </div>

      <Form {...form}>
        <form onSubmit={(e) => void handleSave(e)} className="max-w-2xl space-y-5">
          <FormField
            control={form.control}
            name="customerName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Customer name</FormLabel>
                <FormControl>
                  <Input autoComplete="organization" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="createNewDevices"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded border p-3">
                <div>
                  <FormLabel>Create new devices on first access</FormLabel>
                </div>
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="newDeviceConfigurationId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Default configuration for new devices</FormLabel>
                <Select
                  value={field.value == null ? NONE_CONFIG_VALUE : String(field.value)}
                  onValueChange={(v) => field.onChange(v === NONE_CONFIG_VALUE ? null : Number(v))}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select configuration" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={NONE_CONFIG_VALUE}>None</SelectItem>
                    {configurations.map((c) => (
                      <SelectItem key={c.id} value={String(c.id)}>
                        {c.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="language"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Language</FormLabel>
                <Select value={field.value} onValueChange={field.onChange}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Language" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {LANGUAGE_OPTIONS.map((opt) => (
                      <SelectItem key={opt.value} value={opt.value}>
                        {opt.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="passwordLength"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Minimum password length</FormLabel>
                <FormControl>
                  <Input type="number" min={1} {...field} value={field.value ?? ''} onChange={field.onChange} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="passwordStrength"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password strength</FormLabel>
                <Select
                  value={String(field.value)}
                  onValueChange={(v) => field.onChange(Number(v))}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {PASSWORD_STRENGTH_OPTIONS.map((opt) => (
                      <SelectItem key={opt.value} value={opt.value}>
                        {opt.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="sendDeviceInfoExpiryDays"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Send device info expiry (days)</FormLabel>
                <FormControl>
                  <Input type="number" min={1} {...field} value={field.value ?? ''} onChange={field.onChange} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="unsecureEnrollment"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded border p-3">
                <div>
                  <FormLabel>Unsecure enrollment</FormLabel>
                </div>
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="deviceFastSearch"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded border p-3">
                <div>
                  <FormLabel>Device fast search</FormLabel>
                </div>
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
              </FormItem>
            )}
          />

          <Button type="submit" disabled={submitting}>
            {submitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
            Save Settings
          </Button>
        </form>
      </Form>
    </div>
  )
}
