import type {
  Configuration,
  ConfigurationApplication,
  ConfigurationPayload,
} from '@/features/configurations/types'

/** Linked rows expose `configurationApplications.applicationVersionId` as `usedVersionId`; placeholders omit it (`NULL`). */
function hasLinkedApplicationVersion(rec: Record<string, unknown>): boolean {
  const v = Number(rec.usedVersionId ?? rec.applicationVersionId ?? 0)
  return Number.isFinite(v) && v > 0
}

function usedVersionIdFromPayload(app: ConfigurationApplication): number {
  const u = (app as Record<string, unknown>).usedVersionId
  return u != null && Number(u) > 0 ? Number(u) : 0
}

/** Row from the configuration editor MDM picker (version id = `applicationVersions.id`). */
export interface MdmVersionCatalogRow {
  applicationId: number
  versionId: number
  name: string
  action: number
}

/**
 * `configurations.mainAppId` / `contentAppId` reference **application version** ids.
 * If the chosen version is not yet in `configurationApplications`, `recheckConfigurationMainApplication`
 * returns no row and clears `mainAppId` → UI shows "None".
 */
export function ensureLinkedRowsForChosenVersions(
  applications: ConfigurationApplication[] | undefined,
  mainVersionId: number | null | undefined,
  contentVersionId: number | null | undefined,
  catalog: MdmVersionCatalogRow[]
): ConfigurationApplication[] {
  const next = Array.isArray(applications) ? applications.map((a) => ({ ...a })) : []
  const versionIds = new Set(
    next.map(usedVersionIdFromPayload).filter((v) => v > 0)
  )

  const addRow = (versionId: number) => {
    if (!Number.isFinite(versionId) || versionId <= 0 || versionIds.has(versionId)) return
    const pick = catalog.find((c) => c.versionId === versionId)
    if (!pick || pick.applicationId <= 0) return
    next.push({
      id: pick.applicationId,
      name: pick.name,
      action: 1,
      usedVersionId: versionId,
      showIcon: true,
    })
    versionIds.add(versionId)
  }

  addRow(Number(mainVersionId ?? 0))
  addRow(Number(contentVersionId ?? 0))
  return next
}

/**
 * Maps `GET /private/configurations/applications/:id` rows into payloads expected by `PUT /private/configurations`.
 * Linked rows carry `usedVersionId` (bound `configurationApplications.applicationVersionId`); unlinked UI rows omit it.
 * Do **not** use `applications.latestVersion` to decide linkage — every app row carries that field.
 *
 * Omitting linked apps cleared `configurationApplications` and allowed `recheckConfigurationMainApplication` to null main app.
 */
export function configurationApplicationsForSaveFromApi(rows: unknown): ConfigurationApplication[] {
  const list = Array.isArray(rows) ? rows : []
  return list
    .filter((item) => hasLinkedApplicationVersion(item as Record<string, unknown>))
    .map((item) => {
      const rec = item as Record<string, unknown>
      const id = Number(rec.id ?? 0)
      const usedVersionIdRaw = rec.usedVersionId ?? rec.applicationVersionId
      const latestVersionRaw = rec.latestVersion
      const rawAction = rec.action
      const actionNum = rawAction === undefined || rawAction === null ? 1 : Number(rawAction)

      const base: ConfigurationApplication = {
        id,
        name: rec.name != null ? String(rec.name) : null,
        pkg: rec.pkg != null ? String(rec.pkg) : null,
        type: rec.type != null ? String(rec.type) : null,
        action: actionNum,
        version: rec.version != null ? String(rec.version) : null,
        url: rec.url != null ? String(rec.url) : null,
        versionCode: rec.versionCode != null ? Number(rec.versionCode) : null,
      }
      const extra: Record<string, unknown> = {}
      if (usedVersionIdRaw != null && Number(usedVersionIdRaw) > 0) {
        extra.usedVersionId = Number(usedVersionIdRaw)
      }
      if (latestVersionRaw != null && Number(latestVersionRaw) > 0) {
        extra.latestVersion = Number(latestVersionRaw)
      }
      if (rec.remove !== undefined) base.remove = Boolean(rec.remove)
      if (rec.showIcon !== undefined) base.showIcon = Boolean(rec.showIcon)
      if (rec.screenOrder != null) extra.screenOrder = Number(rec.screenOrder)
      if (rec.keyCode != null) extra.keyCode = Number(rec.keyCode)
      if (rec.bottom !== undefined) extra.bottom = Boolean(rec.bottom)
      if (rec.longTap !== undefined) extra.longTap = Boolean(rec.longTap)
      if (rec.skipVersion === true) base.skipVersion = true
      return { ...base, ...extra }
    })
    .filter((row) => Number(row.id) > 0)
}

export function normalizeConfigurationPayload(payload: ConfigurationPayload): ConfigurationPayload {
  const description = payload.description?.trim()
  return {
    name: payload.name.trim(),
    type: payload.type,
    description: description ? description : null,
  }
}

/** Minimal defaults so inserts/updates remain backend-compatible. */
export function buildCreateConfigurationBody(payload: ConfigurationPayload): Configuration {
  const normalized = normalizeConfigurationPayload(payload)
  return {
    name: normalized.name,
    description: normalized.description,
    type: normalized.type === 'COMMON' ? 1 : 0,
    applications: [],
    iconSize: 'SMALL',
    desktopHeader: 'NO_HEADER',
    useDefaultDesignSettings: true,
    pushOptions: 'mqttWorker',
    defaultFilePath: '/',
    downloadUpdates: 'UNLIMITED',
    requestUpdates: 'DONOTTRACK',
    appPermissions: 'GRANTALL',
    kioskMode: false,
    encryptDevice: false,
    mobileEnrollment: false,
    displayStatus: false,
    blockStatusBar: false,
    systemUpdateType: 0,
    scheduleAppUpdate: false,
    disableLocation: false,
    permissive: false,
    selected: false,
  }
}

export function mergeConfigurationForUpdate(
  current: Configuration,
  id: number,
  payload: ConfigurationPayload
): Configuration {
  const normalized = normalizeConfigurationPayload(payload)
  return {
    ...current,
    id,
    name: normalized.name,
    description: normalized.description,
    type: normalized.type === 'COMMON' ? 1 : 0,
  }
}
