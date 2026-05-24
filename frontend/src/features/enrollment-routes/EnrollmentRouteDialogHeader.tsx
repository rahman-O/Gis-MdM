import { useTranslation } from 'react-i18next'
import { Badge } from '@/shared/ui/badge'
import type { EnrollmentRouteDialogStateId } from '@/features/enrollment-routes/enrollmentRouteDialogState'
import {
  showActiveBadge,
  showDraftBadge,
  showUnsavedBadge,
} from '@/features/enrollment-routes/enrollmentRouteDialogState'

interface Props {
  title: string
  state: EnrollmentRouteDialogStateId
  routeId: number
  dirty: boolean
}

export function EnrollmentRouteDialogHeader({ title, state, routeId, dirty }: Props) {
  const { t } = useTranslation()

  return (
    <div className="flex flex-wrap items-center gap-2">
      <h2 className="text-lg font-semibold">{title}</h2>
      {showDraftBadge(state) ? (
        <Badge variant="secondary">{t('enrollmentRoute.status.draft')}</Badge>
      ) : null}
      {showActiveBadge(state, routeId) ? (
        <Badge variant="default">{t('enrollmentRoute.status.active')}</Badge>
      ) : null}
      {showUnsavedBadge(state, dirty) ? (
        <Badge variant="outline">{t('enrollmentRoute.status.unsaved')}</Badge>
      ) : null}
    </div>
  )
}
