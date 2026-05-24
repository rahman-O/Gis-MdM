import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'

export interface OnboardingStep {
  id: string
  label: string
  done: boolean
  path?: string
}

export interface OnboardingStatus {
  complete: boolean
  hasTreeBeyondRoot: boolean
  hasPublishedProfile: boolean
  hasEnrollmentRoute: boolean
  steps: OnboardingStep[]
}

export async function getOnboardingStatus(): Promise<OnboardingStatus> {
  const response = await apiClient.get<HmdmEnvelope<OnboardingStatus>>('/private/onboarding/status')
  return unwrapHmdmData(response.data, 'Failed to load setup status.')
}
