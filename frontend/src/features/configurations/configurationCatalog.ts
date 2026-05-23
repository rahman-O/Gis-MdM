/** Application row for configuration editor pickers (catalogue). */
export interface ConfigurationAppCatalogItem {
  id: number
  name: string
  latestVersionId: number | null
}

function asRecord(value: unknown): Record<string, unknown> | null {
  return value != null && typeof value === 'object' ? (value as Record<string, unknown>) : null
}

/** Maps `/private/applications/search` or `/private/configurations/applications` rows. */
export function mapApplicationCatalogRows(rows: unknown): ConfigurationAppCatalogItem[] {
  const list = Array.isArray(rows) ? rows : []
  return list
    .map((item) => {
      const rec = asRecord(item)
      if (!rec) return null
      const id = Number(rec.id ?? rec.applicationId ?? 0)
      const latest = Number(rec.latestVersion ?? rec.latestVersionId ?? 0)
      const name = String(rec.name ?? rec.pkg ?? '').trim()
      if (id <= 0) return null
      return {
        id,
        name: name || `Application #${id}`,
        latestVersionId: latest > 0 ? latest : null,
      }
    })
    .filter((x): x is ConfigurationAppCatalogItem => x != null)
}
