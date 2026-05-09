export interface LoginRequest {
  login: string
  password: string
}

export interface LoginResponse {
  authToken: string
}
