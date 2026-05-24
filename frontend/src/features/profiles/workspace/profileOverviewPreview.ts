import * as profileService from '@/features/profiles/profileService'
import type { Profile } from '@/features/profiles/types'
import type { ConfigurationApplication } from '@/features/configurations/types'

export interface OverviewPolicyPreview {
  versionId: number
  versionNumber: number | null
  status: string
  kioskMode: boolean
  appCount: number
  mainAppName?: string
}

function appRecord(app: ConfigurationApplication): Record<string, unknown> {
  return app as Record<string, unknown>
}

function countInstallApps(apps: ConfigurationApplication[]): number {
  const install = apps.filter((a) => {
    const action = a.action
    const n = action === undefined || action === null ? 1 : Number(action)
    return n === 1
  })
  return install.length > 0 ? install.length : apps.length
}

function resolveMainAppName(profile: Profile, apps: ConfigurationApplication[]): string | undefined {
  const mainVersionId = Number(profile.mainAppId ?? 0)
  const contentVersionId = Number(profile.contentAppId ?? 0)
  const targetVersionId = mainVersionId > 0 ? mainVersionId : contentVersionId
  if (targetVersionId <= 0) return undefined

  for (const app of apps) {
    const rec = appRecord(app)
    const usedVersionId = Number(rec.usedVersionId ?? rec.applicationVersionId ?? 0)
    const appId = Number(rec.applicationId ?? rec.id ?? app.id ?? 0)
    if (usedVersionId === targetVersionId || appId === targetVersionId) {
      const name = String(rec.name ?? rec.applicationName ?? '').trim()
      if (name) return name
    }
  }

  const byId = apps.find((a) => {
    const rec = appRecord(a)
    return Number(rec.id ?? rec.applicationId ?? 0) === targetVersionId
  })
  if (byId) {
    const rec = appRecord(byId)
    const name = String(rec.name ?? '').trim()
    if (name) return name
  }

  return undefined
}

export function previewFromProfile(
  profile: Profile,
  versionId: number,
  status: string
): OverviewPolicyPreview {
  const apps = Array.isArray(profile.applications) ? profile.applications : []
  return {
    versionId,
    versionNumber: profile.versionNumber ?? null,
    status,
    kioskMode: Boolean(profile.kioskMode),
    appCount: countInstallApps(apps),
    mainAppName: resolveMainAppName(profile, apps),
  }
}

export async function loadDraftOverviewPreview(
  profileId: number,
  draftVersionId: number
): Promise<OverviewPolicyPreview | null> {
  try {
    const cfg = await profileService.getProfileVersion(profileId, draftVersionId)
    return previewFromProfile(cfg, draftVersionId, cfg.versionStatus ?? 'draft')
  } catch {
    return null
  }
}

/** Enrich published summary pinned fields when hub app count is missing. */
export async function enrichPublishedPreviewFromVersion(
  profileId: number,
  publishedVersionId: number,
  pinned: { kioskMode: boolean; appCount: number; mainAppName?: string; lastPublishedAt?: string | null }
): Promise<OverviewPolicyPreview> {
  if (pinned.appCount > 0 && pinned.mainAppName) {
    return {
      versionId: publishedVersionId,
      versionNumber: null,
      status: 'published',
      kioskMode: pinned.kioskMode,
      appCount: pinned.appCount,
      mainAppName: pinned.mainAppName,
    }
  }
  try {
    const cfg = await profileService.getProfileVersion(profileId, publishedVersionId)
    const fromCfg = previewFromProfile(cfg, publishedVersionId, 'published')
    return {
      ...fromCfg,
      versionNumber: cfg.versionNumber ?? fromCfg.versionNumber,
      kioskMode: fromCfg.kioskMode || pinned.kioskMode,
      appCount: fromCfg.appCount > 0 ? fromCfg.appCount : pinned.appCount,
      mainAppName: fromCfg.mainAppName ?? pinned.mainAppName,
    }
  } catch {
    return {
      versionId: publishedVersionId,
      versionNumber: null,
      status: 'published',
      kioskMode: pinned.kioskMode,
      appCount: pinned.appCount,
      mainAppName: pinned.mainAppName,
    }
  }
}
