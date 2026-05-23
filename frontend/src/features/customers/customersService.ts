import apiClient from '@/services/apiClient'
import type { HmdmEnvelope } from '@/services/hmdmEnvelope'
import { unwrapHmdmData } from '@/services/hmdmEnvelope'
import type { LoginUserPayload } from '@/features/auth/types'
import { applySessionFromUserPayload } from '@/features/auth/session'

export interface CustomerRow {
  id?: number | null
  name?: string | null
  email?: string | null
}

export interface PaginatedCustomers {
  records?: CustomerRow[]
  rows?: CustomerRow[]
  items?: CustomerRow[]
  data?: CustomerRow[]
  totalRecords?: number
  /** Headwind paginated envelope */
  totalItemsCount?: number
}

export async function searchCustomers(payload: Record<string, unknown>): Promise<PaginatedCustomers | CustomerRow[]> {
  const response = await apiClient.post<HmdmEnvelope<PaginatedCustomers | CustomerRow[]>>(
    '/private/customers/search',
    payload
  )
  const data = unwrapHmdmData(response.data, 'Failed to search customers.')
  return Array.isArray(data) ? ({ records: data } as PaginatedCustomers) : data
}

export function unwrapCustomerRows(parsed: PaginatedCustomers | CustomerRow[]): CustomerRow[] {
  if (Array.isArray(parsed)) return parsed
  return parsed.items ?? parsed.records ?? parsed.rows ?? parsed.data ?? []
}

/** Returns target user/token when server supports bearer tokens for impersonated admin. */
export async function impersonateCustomer(customerId: number): Promise<LoginUserPayload> {
  const response = await apiClient.get<HmdmEnvelope<LoginUserPayload>>(`/private/customers/impersonate/${customerId}`)
  return unwrapHmdmData(response.data, 'Impersonation failed.')
}

export function hydrateSessionAfterImpersonation(user: LoginUserPayload): void {
  const tok = user.authToken ?? undefined
  applySessionFromUserPayload(
    {
      superAdmin: user.superAdmin,
      singleCustomer: user.singleCustomer,
      userRole: user.userRole ?? undefined,
    },
    typeof tok === 'string' ? tok : undefined
  )
}
