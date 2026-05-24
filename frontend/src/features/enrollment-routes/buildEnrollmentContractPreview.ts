export interface EnrollmentContractPreviewInput {
  routeId?: number
  targetNodeId: number
  mainAppPackage: string
  mainAppVersion: string
  mainAppVersionCode: number
  deviceIdentityMode: string
}

export interface EnrollmentContractPreview {
  routeId: number
  targetNodeId: number
  mainAppPackage: string
  mainAppVersion: string
  mainAppVersionCode: number
  deviceIdentityMode: string
  bootstrapFlags: { create: boolean }
  _preview: true
}

export function buildEnrollmentContractPreview(
  input: EnrollmentContractPreviewInput
): EnrollmentContractPreview | null {
  if (!input.targetNodeId || input.targetNodeId <= 0) return null
  if (!input.mainAppPackage?.trim()) return null
  if (!input.mainAppVersion?.trim()) return null

  return {
    routeId: input.routeId && input.routeId > 0 ? input.routeId : 0,
    targetNodeId: input.targetNodeId,
    mainAppPackage: input.mainAppPackage,
    mainAppVersion: input.mainAppVersion,
    mainAppVersionCode: input.mainAppVersionCode,
    deviceIdentityMode: input.deviceIdentityMode || 'imei',
    bootstrapFlags: { create: true },
    _preview: true,
  }
}
