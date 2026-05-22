import { describe, expect, it } from 'vitest'
import { mapApplicationCatalogRows } from '@/features/configurations/configurationCatalog'

describe('mapApplicationCatalogRows', () => {
  it('maps applications search rows with latestVersion', () => {
    expect(
      mapApplicationCatalogRows([
        { id: 3, name: 'My App', latestVersion: 99 },
        { applicationId: 4, pkg: 'com.example', latestVersionId: 10 },
      ])
    ).toEqual([
      { id: 3, name: 'My App', latestVersionId: 99 },
      { id: 4, name: 'com.example', latestVersionId: 10 },
    ])
  })

  it('ignores rows without id', () => {
    expect(mapApplicationCatalogRows([{ name: 'orphan' }])).toEqual([])
  })
})
