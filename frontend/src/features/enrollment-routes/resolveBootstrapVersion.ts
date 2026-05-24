import type {
  BootstrapAppOption,
  BootstrapIntent,
} from '@/features/enrollment-routes/enrollmentRouteService'

export interface ResolvedBootstrapVersion {
  package: string
  version: string
  versionCode: number
}

export function resolveBootstrapVersion(
  apps: BootstrapAppOption[],
  applicationId: number | '',
  intent: BootstrapIntent,
  specificVersionId: number | ''
): ResolvedBootstrapVersion | null {
  if (!applicationId || applicationId <= 0) return null
  const app = apps.find((a) => a.applicationId === applicationId)
  if (!app) return null

  const versions = app.versions ?? []
  if (versions.length === 0) return null

  let picked: (typeof versions)[number] | undefined = versions[0]
  if (intent === 'specific') {
    if (!specificVersionId || specificVersionId <= 0) return null
    picked = versions.find((v) => v.versionId === specificVersionId)
    if (!picked) return null
  } else if (intent === 'latest') {
    picked = versions.find((v) => v.isLatest) ?? versions[0]
  } else {
    picked = versions.find((v) => v.isRecommended) ?? versions[0]
  }
  if (!picked) return null

  return {
    package: app.package,
    version: picked.version,
    versionCode: picked.versionCode,
  }
}
