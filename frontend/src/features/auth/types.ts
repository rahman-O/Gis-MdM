export interface LoginRequest {
  login: string
  password: string
}

/** Server returns `UserView` on login; only fields we persist are modeled here. */
export interface LoginUserPayload {
  authToken?: string | null
  superAdmin?: boolean
  userRole?: {
    superAdmin?: boolean
    permissions?: Array<{ name?: string | null } | null>
  } | null
}

export interface LoginResponse {
  authToken: string
}
