import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { md5UpperHex } from '@/features/auth/loginPasswordEncode'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { checkPasswordQuality, passwordPolicyHint } from '@/features/auth/passwordPolicy'
import * as passwordResetPublicService from '@/features/auth/passwordResetPublicService'

export function PasswordResetPage() {
  const { token: tokenParam } = useParams()
  const navigate = useNavigate()
  const token = tokenParam ? decodeURIComponent(tokenParam) : ''

  const [loading, setLoading] = useState(true)
  const [valid, setValid] = useState(false)
  const [policy, setPolicy] = useState<{ length: number; strength: number }>({ length: 6, strength: 0 })
  const [pwd1, setPwd1] = useState('')
  const [pwd2, setPwd2] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    if (!token) {
      setLoading(false)
      return
    }
    let cancelled = false
    void (async () => {
      try {
        const s = await passwordResetPublicService.fetchResetSettings(token)
        if (cancelled) return
        setValid(true)
        setPolicy({
          length: Math.max(1, Number(s.passwordLength ?? 6)),
          strength: Math.min(2, Math.max(0, Number(s.passwordStrength ?? 0))),
        })
      } catch {
        if (!cancelled) setValid(false)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [token])

  async function save() {
    setError(null)
    if (!pwd1 || pwd1 !== pwd2) {
      setError('Passwords must match.')
      return
    }
    if (!checkPasswordQuality(pwd1, policy.length, policy.strength)) {
      setError(passwordPolicyHint(policy.length, policy.strength))
      return
    }
    setBusy(true)
    try {
      await passwordResetPublicService.submitPasswordReset(token, md5UpperHex(pwd1))
      navigate('/dashboard')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Reset failed.')
    } finally {
      setBusy(false)
    }
  }

  if (!token) {
    return (
      <div className="bg-background flex min-h-screen items-center justify-center px-4">
        <p className="text-muted-foreground text-sm">Missing reset token.</p>
      </div>
    )
  }

  return (
    <div className="bg-background flex min-h-screen items-center justify-center px-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Choose new password</CardTitle>
          <CardDescription>{passwordPolicyHint(policy.length, policy.strength)}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {loading ? (
            <div className="text-muted-foreground flex items-center gap-2 text-sm">
              <Loader2 className="h-4 w-4 animate-spin" />
              Checking link…
            </div>
          ) : !valid ? (
            <p className="text-destructive text-sm">This reset link is invalid or expired.</p>
          ) : (
            <>
              <Input type="password" placeholder="New password" value={pwd1} onChange={(e) => setPwd1(e.target.value)} />
              <Input
                type="password"
                placeholder="Confirm password"
                value={pwd2}
                onChange={(e) => setPwd2(e.target.value)}
              />
              {error ? (
                <p role="alert" className="text-destructive text-sm">
                  {error}
                </p>
              ) : null}
              <Button className="w-full" disabled={busy} type="button" onClick={() => void save()}>
                {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Save password'}
              </Button>
            </>
          )}
          <Button variant="link" className="px-0" asChild>
            <Link to="/login">Back to sign in</Link>
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
