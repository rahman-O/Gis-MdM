import { useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import * as configurationService from '@/features/configurations/configurationService'
import type { Configuration } from '@/features/configurations/types'
import { EnrollmentQrExperience } from '@/features/devices/EnrollmentQrExperience'
import * as deviceService from '@/features/devices/deviceService'
import type { LookupItem } from '@/features/devices/types'
import { canEnrollDevicesViaQr } from '@/features/auth/permissions'

export function EnrollmentQrPage() {
  const params = useParams<{ qrCodeKey?: string; deviceId?: string }>()
  const rawKeyParam = params.qrCodeKey
  let qrCodeKey = ''
  try {
    qrCodeKey = rawKeyParam != null ? decodeURIComponent(rawKeyParam) : ''
  } catch {
    qrCodeKey = rawKeyParam ?? ''
  }
  const routeDeviceId = params.deviceId != null ? decodeURIComponent(params.deviceId) : ''

  const [configurations, setConfigurations] = useState<Configuration[]>([])
  const [groups, setGroups] = useState<LookupItem[]>([])
  const [loadErr, setLoadErr] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoadErr(null)
      try {
        const list = await configurationService.getConfigurations()
        if (!cancelled) setConfigurations(Array.isArray(list) ? list : [])
      } catch (e: unknown) {
        if (!cancelled) setLoadErr(e instanceof Error ? e.message : 'Could not preload configurations.')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  useEffect(() => {
    void deviceService.getGroups().then(setGroups).catch(() => setGroups([]))
  }, [])

  const matchedConfig = useMemo(
    () => configurations.find((c) => String(c.qrCodeKey ?? '').trim() === qrCodeKey.trim()),
    [configurations, qrCodeKey]
  )

  const [eligibilityConfig, setEligibilityConfig] = useState<Configuration | null>(null)

  useEffect(() => {
    if (matchedConfig?.id == null) {
      setEligibilityConfig(null)
      return
    }
    let cancelled = false
    void (async () => {
      try {
        const full = await configurationService.getConfiguration(matchedConfig.id as number)
        if (!cancelled) setEligibilityConfig(full)
      } catch {
        if (!cancelled) setEligibilityConfig(matchedConfig)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [matchedConfig])

  if (!canEnrollDevicesViaQr()) {
    return (
      <div className="mx-auto max-w-lg space-y-4 p-6">
        <p className="text-sm text-muted-foreground">You do not have permission to enroll devices via QR.</p>
        <Button asChild variant="outline">
          <Link to="/dashboard">Back</Link>
        </Button>
      </div>
    )
  }

  if (!qrCodeKey.trim()) {
    return (
      <div className="mx-auto max-w-lg space-y-4 p-6">
        <p className="text-destructive text-sm">Missing QR code key in the URL.</p>
        <Button asChild variant="outline">
          <Link to="/configurations">
            <ArrowLeft className="mr-2 h-4 w-4" /> Configurations
          </Link>
        </Button>
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6 p-4 pb-16">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight">Device enrollment QR</h1>
          <p className="text-muted-foreground text-sm">
            Configuration QR key <span className="font-mono text-foreground">{qrCodeKey}</span>
          </p>
        </div>
        <Button asChild variant="ghost" size="sm">
          <Link to="/devices">
            <ArrowLeft className="mr-2 h-4 w-4" /> Devices
          </Link>
        </Button>
      </div>

      {loadErr ? <p className="text-sm text-amber-800 dark:text-amber-200">{loadErr}</p> : null}
      {!matchedConfig && !loadErr ? (
        <p className="text-sm text-muted-foreground">
          No matching configuration was found in search results for eligibility hints; QR generation still works if the
          key is valid.
        </p>
      ) : null}

      <EnrollmentQrExperience
        qrCodeKey={qrCodeKey.trim()}
        initialDeviceId={routeDeviceId}
        configuration={eligibilityConfig ?? matchedConfig ?? null}
        groups={groups}
      />
    </div>
  )
}
