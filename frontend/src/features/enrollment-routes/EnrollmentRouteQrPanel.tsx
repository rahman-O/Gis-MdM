import { EnrollmentQrExperience } from '@/features/devices/EnrollmentQrExperience'
import type { EnrollmentRouteQrMeta } from '@/features/enrollment-routes/enrollmentRouteService'

interface Props {
  meta: EnrollmentRouteQrMeta
}

/** QR panel for enrollment routes — renders scannable QR code. */
export function EnrollmentRouteQrPanel({ meta }: Props) {
  return (
    <EnrollmentQrExperience
      qrCodeKey={meta.qrcodekey}
      configuration={{
        qrCodeKey: meta.qrcodekey,
        mainAppId: meta.resolvedMainAppVersionId ?? undefined,
      }}
    />
  )
}
