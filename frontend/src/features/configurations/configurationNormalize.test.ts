import { describe, expect, it } from 'vitest'
import {
  buildCreateConfigurationBody,
  configurationApplicationsForSaveFromApi,
  ensureLinkedRowsForChosenVersions,
  mergeConfigurationForUpdate,
  normalizeConfigurationPayload,
} from '@/features/configurations/configurationNormalize'

describe('configurationNormalize', () => {
  it('normalizeConfigurationPayload trims and nullifies description', () => {
    expect(
      normalizeConfigurationPayload({
        name: '  Test  ',
        type: 'WORK',
        description: '   ',
      })
    ).toEqual({
      name: 'Test',
      type: 'WORK',
      description: null,
    })
  })

  it('buildCreateConfigurationBody adds backend-safe defaults', () => {
    const body = buildCreateConfigurationBody({
      name: 'A',
      type: 'COMMON',
      description: null,
    })
    expect(body).toEqual(
      expect.objectContaining({
        name: 'A',
        type: 1,
        applications: [],
        pushOptions: 'mqttWorker',
        defaultFilePath: '/',
      })
    )
  })

  it('configurationApplicationsForSaveFromApi keeps only linked rows (non-null usedVersionId)', () => {
    const rows = [
      {
        id: 10,
        name: 'Placeholder (latest shown in UI only)',
        action: 0,
        selected: false,
        latestVersion: 999,
      },
      {
        id: 20,
        usedVersionId: 200,
        action: 1,
        latestVersion: 555,
        name: 'Actually linked',
        showIcon: true,
      },
    ]
    expect(configurationApplicationsForSaveFromApi(rows)).toEqual([
      expect.objectContaining({
        id: 20,
        usedVersionId: 200,
        latestVersion: 555,
        name: 'Actually linked',
        showIcon: true,
      }),
    ])
  })

  it('configurationApplicationsForSaveFromApi ignores invalid ids', () => {
    expect(
      configurationApplicationsForSaveFromApi([{ id: 0, usedVersionId: 9, action: 1 }])
    ).toHaveLength(0)
  })

  it('configurationApplicationsForSaveFromApi defaults missing action on linked rows to 1', () => {
    const out = configurationApplicationsForSaveFromApi([{ id: 7, usedVersionId: 99, name: 'X' }])
    expect(out).toHaveLength(1)
    expect(out[0].action).toBe(1)
  })

  it('ensureLinkedRowsForChosenVersions inserts catalog rows when version not listed', () => {
    const catalog = [
      { applicationId: 5, versionId: 100, name: 'Launcher', action: 1 },
      { applicationId: 6, versionId: 200, name: 'Content', action: 1 },
    ]
    const merged = ensureLinkedRowsForChosenVersions([], 100, 200, catalog)
    expect(merged.map((x) => (x as { usedVersionId?: number }).usedVersionId)).toEqual(expect.arrayContaining([100, 200]))
    expect(merged.find((x) => (x as { usedVersionId?: number }).usedVersionId === 100)).toEqual(
      expect.objectContaining({
        id: 5,
        usedVersionId: 100,
        action: 1,
      })
    )
  })

  it('mergeConfigurationForUpdate keeps existing fields and replaces base metadata', () => {
    const merged = mergeConfigurationForUpdate(
      { id: 3, kioskMode: true, name: 'Old', type: 0, description: null },
      3,
      {
        name: ' New ',
        type: 'COMMON',
        description: ' Desc ',
      }
    )
    expect(merged).toEqual(
      expect.objectContaining({
        id: 3,
        kioskMode: true,
        name: 'New',
        type: 1,
        description: 'Desc',
      })
    )
  })
})
