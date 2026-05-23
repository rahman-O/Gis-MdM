export interface LoginRequest {
  login: string
  password: string
}

/** Server returns `UserView` on login; only fields we persist are modeled here. */
export interface LoginUserPayload {
  id?: number
  login?: string | null
  name?: string | null
  email?: string | null
  authToken?: string | null
  superAdmin?: boolean
  passwordReset?: boolean
  passwordResetToken?: string | null
  twoFactor?: boolean | null
  twoFactorAccepted?: boolean | null
  /** Mirrors `UserView.isSingleCustomer` from the Headwind API. */
  singleCustomer?: boolean
  userRole?: {
    id?: number
    superAdmin?: boolean
    permissions?: Array<{ name?: string | null } | null>
  } | null
}

/** Outcome consumed by `AuthProvider` for post-login routing. */
export interface LoginOutcome {
  authToken: string
  redirectPath: string
}

export interface AuthLandingOptions {
  signup: boolean
  recover: boolean
  publicKey?: string | null
}
