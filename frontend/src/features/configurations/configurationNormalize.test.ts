import { describe, expect, it } from 'vitest'
import {
  buildCreateConfigurationBody,
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
