import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { EnrollmentDeleteImpact } from '@/features/enrollment-routes/enrollmentRouteService'
import { getEnrollmentRouteImpact } from '@/features/enrollment-routes/enrollmentRouteService'
import type { EnrollmentRouteDialogStateId } from '@/features/enrollment-routes/enrollmentRouteDialogState'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'

interface Props {
  routeId: number
  routeName: string
  step: Extract<EnrollmentRouteDialogStateId, 'DELETE_STEP1' | 'DELETE_STEP2' | 'DELETE_CONFIRM_ZERO'>
  confirmName: string
  onConfirmNameChange: (v: string) => void
}

export function DeleteRouteConfirm({
  routeId,
  routeName,
  step,
  confirmName,
  onConfirmNameChange,
}: Props) {
  const { t } = useTranslation()
  const [impact, setImpact] = useState<EnrollmentDeleteImpact | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    void getEnrollmentRouteImpact(routeId)
      .then(setImpact)
      .catch((e: unknown) => {
        setError(e instanceof Error ? e.message : 'Failed to load impact.')
      })
      .finally(() => setLoading(false))
  }, [routeId])

  const anyImpact =
    (impact?.enrollingNowCount ?? 0) > 0 ||
    (impact?.historicalEnrolledCount ?? 0) > 0 ||
    (impact?.activeQrScans7d ?? 0) > 0

  if (loading) {
    return <p className="text-sm text-muted-foreground">{t('enrollmentRoute.delete.loading')}</p>
  }
  if (error) {
    return <p className="text-sm text-destructive">{error}</p>
  }

  return (
    <div className="space-y-4">
      <p className="text-sm">{t('enrollmentRoute.delete.summary', { name: routeName })}</p>
      <ul className="list-inside list-disc text-sm text-muted-foreground">
        <li>
          {t('enrollmentRoute.delete.metricEnrolling', {
            count: impact?.enrollingNowCount ?? 0,
          })}
        </li>
        <li>
          {t('enrollmentRoute.delete.metricHistorical', {
            count: impact?.historicalEnrolledCount ?? 0,
          })}
        </li>
        <li>
          {t('enrollmentRoute.delete.metricScans', {
            count: impact?.activeQrScans7d ?? 0,
          })}
        </li>
      </ul>
      {step === 'DELETE_STEP2' && anyImpact ? (
        <div className="space-y-2">
          <Label htmlFor="confirm-name">{t('enrollmentRoute.delete.typeName')}</Label>
          <Input
            id="confirm-name"
            value={confirmName}
            onChange={(e) => onConfirmNameChange(e.target.value)}
            placeholder={routeName}
          />
        </div>
      ) : null}
    </div>
  )
}

export function deleteRequiresTypedName(
  step: EnrollmentRouteDialogStateId,
  impact: EnrollmentDeleteImpact | null
): boolean {
  if (step !== 'DELETE_STEP2') return false
  const any =
    (impact?.enrollingNowCount ?? 0) > 0 ||
    (impact?.historicalEnrolledCount ?? 0) > 0 ||
    (impact?.activeQrScans7d ?? 0) > 0
  return any
}
