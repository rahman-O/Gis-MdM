import type { EnrollmentRouteView } from '@/features/enrollment-routes/enrollmentRouteService'

/** Dialog state IDs from specs/021-enrollment-routes-ux/contracts/enrollment-route-dialog-ux.md */
export type EnrollmentRouteDialogStateId =
  | 'LIST'
  | 'DIALOG_CREATE'
  | 'DIALOG_OVERVIEW'
  | 'DIALOG_EDIT'
  | 'DELETE_STEP1'
  | 'DELETE_STEP2'
  | 'DELETE_CONFIRM_ZERO'

export interface EnrollmentRouteFormValues {
  name: string
  description: string
  targetNodeId: number | ''
  deviceIdentityMode: string
  bootstrapIntent: 'stable' | 'specific' | 'latest'
  bootstrapApplicationId: number | ''
  bootstrapVersionId: number | ''
  acknowledgeContainerPlacement: boolean
  wifiSsid: string
  wifiPassword: string
  wifiSecurityType: string
  qrParameters: string
  adminExtras: string
  mobileEnrollment: boolean
  encryptDevice: boolean
}

export function emptyFormValues(): EnrollmentRouteFormValues {
  return {
    name: '',
    description: '',
    targetNodeId: '',
    deviceIdentityMode: 'imei',
    bootstrapIntent: 'stable',
    bootstrapApplicationId: '',
    bootstrapVersionId: '',
    acknowledgeContainerPlacement: false,
    wifiSsid: '',
    wifiPassword: '',
    wifiSecurityType: '',
    qrParameters: '',
    adminExtras: '',
    mobileEnrollment: false,
    encryptDevice: false,
  }
}

export function formValuesFromRoute(route: EnrollmentRouteView): EnrollmentRouteFormValues {
  return {
    name: route.name,
    description: route.description ?? '',
    targetNodeId: route.targetNodeId,
    deviceIdentityMode: route.deviceIdentityMode || 'imei',
    bootstrapIntent: route.bootstrapIntent || 'stable',
    bootstrapApplicationId: route.bootstrapApplicationId,
    bootstrapVersionId: route.bootstrapVersionId ?? '',
    acknowledgeContainerPlacement: route.containerPlacementAcknowledged ?? false,
    wifiSsid: route.wifiSsid ?? '',
    wifiPassword: route.wifiPassword ?? '',
    wifiSecurityType: route.wifiSecurityType ?? '',
    qrParameters: route.qrParameters ?? '',
    adminExtras: route.adminExtras ?? '',
    mobileEnrollment: route.mobileEnrollment ?? false,
    encryptDevice: route.encryptDevice ?? false,
  }
}

export function formValuesEqual(a: EnrollmentRouteFormValues, b: EnrollmentRouteFormValues): boolean {
  return JSON.stringify(a) === JSON.stringify(b)
}

export function isEditDirty(
  form: EnrollmentRouteFormValues,
  saved: EnrollmentRouteFormValues | null
): boolean {
  if (!saved) return true
  return !formValuesEqual(form, saved)
}

export function showDraftBadge(state: EnrollmentRouteDialogStateId): boolean {
  return state === 'DIALOG_CREATE'
}

export function showActiveBadge(state: EnrollmentRouteDialogStateId, routeId: number): boolean {
  return routeId > 0 && state !== 'LIST'
}

export function showUnsavedBadge(
  state: EnrollmentRouteDialogStateId,
  dirty: boolean
): boolean {
  return state === 'DIALOG_EDIT' && dirty
}

export function usesPendingQr(state: EnrollmentRouteDialogStateId): boolean {
  return state === 'DIALOG_CREATE' || state === 'DIALOG_EDIT'
}
