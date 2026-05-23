# Auth stack — complete (Go)

The React/Angular shell can authenticate end-to-end against `serverBackendGo` without the Java WAR for auth flows.

## Endpoints

| Area | Method | Path | Module |
|------|--------|------|--------|
| Options | GET | `/rest/public/auth/options` | auth |
| Session login | POST | `/rest/public/auth/login` | auth |
| Logout | POST | `/rest/public/auth/logout` | auth |
| JWT login | POST | `/rest/public/jwt/login` | auth |
| Password reset settings | GET | `/rest/public/passwordReset/settings/:token` | passwordreset |
| Password reset | POST | `/rest/public/passwordReset/reset` | passwordreset |
| Recover password | GET | `/rest/public/passwordReset/recover/:username` | passwordreset |
| Signup verify email | POST | `/rest/public/signup/verifyEmail` | signup |
| Signup verify token | GET | `/rest/public/signup/verifyToken/:token` | signup |
| Signup complete | POST | `/rest/public/signup/complete` | signup |
| Current user | GET | `/rest/private/users/current` | users |
| 2FA QR | GET | `/rest/private/twofactor/qr/:userId` | twofactor |
| 2FA verify | GET | `/rest/private/twofactor/verify/:userId/:code` | twofactor |
| 2FA set | GET | `/rest/private/twofactor/set` | twofactor |
| 2FA reset | GET | `/rest/private/twofactor/reset` | twofactor |

## Configuration (`.env`)

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | Required for auth |
| `SESSION_SECRET` | Cookie session after `/auth/login` |
| `JWT_SECRET` | Bearer tokens for `/jwt/login` and private API |
| `EMAIL_CONFIGURED` | Enables recover/signup in options |
| `CUSTOMER_SIGNUP` | Enables signup in options + `/signup/*` |
| `TRANSMIT_PASSWORD` | RSA public key in options; decrypt login password |
| `FILES_DIRECTORY` | RSA key storage when `TRANSMIT_PASSWORD=true` |
| `MODULE_AUTH_ENABLED` | Auth + JWT (default true) |
| `MODULE_PASSWORDRESET_ENABLED` | Password reset routes (default true) |
| `MODULE_SIGNUP_ENABLED` | Signup routes (default false) |

## Run

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
```

Default user: `admin` / `admin` (or MD5 hex in JSON).

## Frontend flow

1. `GET /rest/public/auth/options`
2. `POST /rest/public/auth/login` → store `authToken`, session cookie
3. `GET /rest/private/users/current` → permissions / profile
4. Optional: `/twofactor` if `twoFactor` and not `twoFactorAccepted`
5. `POST /rest/public/auth/logout`

Proxy: point Vite to `http://localhost:8080` for `/rest`.

## Out of scope (next migration phases)

- Full `UserResource` (list/update users, password change on `/users/current` PUT)
- Mailchimp, signup copy-settings, customer limits from Java `SignupResource`
- Real SMTP (email is log/stub when `EMAIL_CONFIGURED=true`)

See [parity/auth.md](./parity/auth.md) for Java file mapping.
