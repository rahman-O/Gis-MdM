import type { LookupItem } from '@/features/devices/types'

export interface Role {
  id: number
  name: string
}

export interface User {
  id: number
  login: string
  name: string
  email: string
  role: Role | null
  allDevicesAvailable: boolean
  allConfigAvailable: boolean
  groups: LookupItem[]
  configurations: LookupItem[]
}

export interface UserPayload {
  login: string
  name: string
  email: string
  /** Plaintext from the form; encoded in `userService` before PUT */
  password?: string
  roleId: number
  allDevicesAvailable: boolean
  allConfigAvailable: boolean
  groups: LookupItem[]
  configurations: LookupItem[]
}
