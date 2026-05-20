# Auth API parity

**Status: complete** for React login shell (session, JWT, reset, signup, 2FA, `users/current`).

## Session auth (`AuthResource`)

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/public/auth/options` | `/public/auth/options` | Done |
| POST | `/rest/public/auth/login` | `/public/auth/login` | Done |
| POST | `/rest/public/auth/logout` | `/public/auth/logout` | Done (204) |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/AuthResource.java`

## JWT auth (`JWTAuthResource`)

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| POST | `/rest/public/jwt/login` | `/public/jwt/login` | Done |

**Java:** `backend/jwt/src/main/java/com/hmdm/security/jwt/rest/JWTAuthResource.java`

## Password reset (`PasswordResetResource`)

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/public/passwordReset/settings/:token` | `/public/passwordReset/settings/{token}` | Done |
| POST | `/rest/public/passwordReset/reset` | `/public/passwordReset/reset` | Done |
| GET | `/rest/public/passwordReset/recover/:username` | `/public/passwordReset/recover/{username}` | Done |
| GET | `/rest/public/passwordReset/canRecover` | `/public/passwordReset/canRecover` | Done (deprecated) |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/PasswordResetResource.java`

## Signup (`SignupResource`)

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| POST | `/rest/public/signup/verifyEmail` | `/public/signup/verifyEmail` | Done |
| GET | `/rest/public/signup/verifyToken/:token` | `/public/signup/verifyToken/{token}` | Done |
| POST | `/rest/public/signup/complete` | `/public/signup/complete` | Done |
| GET | `/rest/public/signup/canSignup` | `/public/signup/canSignup` | Done (deprecated) |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SignupResource.java`  
**Note:** Simplified signup (no Mailchimp / copy-settings / device limits).

## Users — current (`UserResource`)

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/users/current` | `/private/users/current` | Done |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/UserResource.java`  
**Note:** Other user CRUD endpoints not migrated yet.

## Two-factor (private)

| Method | Go path | Java (webapp) | Status |
|--------|---------|---------------|--------|
| GET | `/rest/private/twofactor/qr/:userId` | `rest/private/twofactor/qr/:id` | Done |
| GET | `/rest/private/twofactor/verify/:userId/:code` | `.../verify/:user/:code` | Done |
| GET | `/rest/private/twofactor/set` | `.../set` | Done |
| GET | `/rest/private/twofactor/reset` | `.../reset` | Done |

TOTP via `github.com/pquerna/otp`. Session blocks private routes until verify when customer `settings.twoFactor` is on.

## Password encoding

- Client sends **MD5 uppercase hex** of raw password (`frontend/src/features/auth/loginPasswordEncode.ts`).
- Server: `SHA1(md5 + "5YdSYHyg2U")` vs `users.password`.
- Go also accepts **raw password** in API tools (`NormalizeLoginPassword`).
- Optional `TRANSMIT_PASSWORD=true`: RSA decrypt then normalize.

## Swagger

- UI: `http://localhost:8080/swagger/index.html`
- `make swagger`

## Tests

```bash
cd serverBackendGo
go test ./internal/modules/auth/... ./internal/shared/crypto/... ./internal/platform/jwt/... -v
```
