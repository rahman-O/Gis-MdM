import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { md5UpperHex } from '@/features/auth/loginPasswordEncode'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import * as signupPublicService from '@/features/auth/signupPublicService'

const ID_RE = /^[a-zA-Z0-9][a-zA-Z0-9.\-_]{4,48}[a-zA-Z0-9]$/

export function SignupCompletePage() {
  const { token: tokenParam } = useParams()
  const token = tokenParam ? decodeURIComponent(tokenParam) : ''

  const [valid, setValid] = useState<boolean | null>(null)
  const [name, setName] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [company, setCompany] = useState('')
  const [description, setDescription] = useState('')
  const [p1, setP1] = useState('')
  const [p2, setP2] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [complete, setComplete] = useState(false)

  useEffect(() => {
    if (!token) {
      setValid(false)
      return
    }
    void signupPublicService.signupVerifyTokenOk(token).then(setValid)
  }, [token])

  async function submit() {
    setError(null)
    if (!ID_RE.test(name)) {
      setError('Customer ID must be 6–50 chars, start/end with alphanumeric, and may include . - _')
      return
    }
    if (!firstName.trim() || !lastName.trim()) {
      setError('First and last name are required.')
      return
    }
    if (p1.length < 6 || p1 !== p2) {
      setError('Password must be at least 6 characters and match confirmation.')
      return
    }
    setBusy(true)
    try {
      await signupPublicService.signupComplete({
        token,
        name: name.trim(),
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        company: company.trim() || undefined,
        description: description.trim() || undefined,
        passwd: md5UpperHex(p1),
      })
      setComplete(true)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Signup failed.')
    } finally {
      setBusy(false)
    }
  }

  if (!token) {
    return <p className="text-muted-foreground p-6 text-sm">Missing token.</p>
  }

  return (
    <div className="bg-background flex min-h-screen items-center justify-center px-4 py-8">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Complete signup</CardTitle>
          <CardDescription>Create your tenant admin user and workspace.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {valid === null ? (
            <div className="text-muted-foreground flex items-center gap-2 text-sm">
              <Loader2 className="h-4 w-4 animate-spin" />
              Validating link…
            </div>
          ) : !valid ? (
            <p className="text-destructive text-sm">This signup link is invalid or expired.</p>
          ) : complete ? (
            <p className="text-sm">Registration finished. You can sign in with the password you chose.</p>
          ) : (
            <>
              <div>
                <label className="text-muted-foreground text-xs">Customer ID</label>
                <Input value={name} onChange={(e) => setName(e.target.value)} autoComplete="off" />
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-muted-foreground text-xs">First name</label>
                  <Input value={firstName} onChange={(e) => setFirstName(e.target.value)} />
                </div>
                <div>
                  <label className="text-muted-foreground text-xs">Last name</label>
                  <Input value={lastName} onChange={(e) => setLastName(e.target.value)} />
                </div>
              </div>
              <div>
                <label className="text-muted-foreground text-xs">Company (optional)</label>
                <Input value={company} onChange={(e) => setCompany(e.target.value)} />
              </div>
              <div>
                <label className="text-muted-foreground text-xs">Description (optional)</label>
                <Textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={3} />
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-muted-foreground text-xs">Password</label>
                  <Input type="password" value={p1} onChange={(e) => setP1(e.target.value)} />
                </div>
                <div>
                  <label className="text-muted-foreground text-xs">Repeat password</label>
                  <Input type="password" value={p2} onChange={(e) => setP2(e.target.value)} />
                </div>
              </div>
              {error ? (
                <p role="alert" className="text-destructive text-sm">
                  {error}
                </p>
              ) : null}
              <Button type="button" disabled={busy} onClick={() => void submit()}>
                {busy ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Finish signup
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
