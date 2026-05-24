import { EnrollmentQrExperience } from '@/features/devices/EnrollmentQrExperience'
import type { EnrollmentRouteQrMeta } from '@/features/enrollment-routes/enrollmentRouteService'

interface Props {
  meta: EnrollmentRouteQrMeta
}

/** QR panel for enrollment routes (policy lives on Profile). */
export function EnrollmentRouteQrPanel({ meta }: Props) {
  return (
    <EnrollmentQrExperience
      qrCodeKey={meta.qrcodekey}
      configuration={{
        qrCodeKey: meta.qrcodekey,
        mainAppId: meta.mainAppId ?? undefined,
      }}
    />
  )
}
