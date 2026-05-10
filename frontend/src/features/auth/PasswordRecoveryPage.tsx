import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import * as passwordResetPublicService from '@/features/auth/passwordResetPublicService'

export function PasswordRecoveryPage() {
  const [username, setUsername] = useState('')
  const [busy, setBusy] = useState(false)
  const [done, setDone] = useState(false)

  async function submit() {
    setBusy(true)
    try {
      await passwordResetPublicService.requestPasswordRecovery(username)
      setDone(true)
    } catch {
      setDone(true)
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="bg-background flex min-h-screen items-center justify-center px-4">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Password recovery</CardTitle>
          <CardDescription>
            If the account exists and email is configured on the server, a reset message will be sent. (The API always
            returns success to avoid leaking which logins exist.)
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {!done ? (
            <>
              <div>
                <label className="text-muted-foreground text-xs">Username or email</label>
                <Input value={username} onChange={(e) => setUsername(e.target.value)} autoComplete="username" />
              </div>
              <Button className="w-full" disabled={busy || !username.trim()} type="button" onClick={() => void submit()}>
                {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Send recovery link'}
              </Button>
            </>
          ) : (
            <p className="text-sm">If your credentials match a user with email on file, check your inbox.</p>
          )}
          <Button variant="link" className="px-0" asChild>
            <Link to="/login">Back to sign in</Link>
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
