import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { useToast } from '@/shared/hooks/use-toast'
import * as twoFactorAuthService from '@/features/auth/twoFactorAuthService'
import * as profileService from '@/features/profile/profileService'

export function TwofactorPage() {
  const { toast } = useToast()
  const navigate = useNavigate()
  const [userId, setUserId] = useState<number | null>(null)
  const [qrUrl, setQrUrl] = useState<string | null>(null)
  const qrObjectUrl = useRef<string | null>(null)
  const [code, setCode] = useState('')
  const [busy, setBusy] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let cancelled = false
    void (async () => {
      try {
        const u = await profileService.fetchCurrentUser()
        const id = Number(u.id)
        if (!Number.isFinite(id)) return
        setUserId(id)
        const blob = await twoFactorAuthService.fetchTwoFactorQrPngBlob(id)
        if (cancelled) return
        if (qrObjectUrl.current) URL.revokeObjectURL(qrObjectUrl.current)
        const next = URL.createObjectURL(blob)
        if (cancelled) {
          URL.revokeObjectURL(next)
          return
        }
        qrObjectUrl.current = next
        setQrUrl(next)
      } catch {
        toast({
          variant: 'destructive',
          title: 'Two-factor QR could not load',
          description: 'The server must expose `/rest/private/twofactor/qr/:id` for this tenant.',
        })
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
      if (qrObjectUrl.current) {
        URL.revokeObjectURL(qrObjectUrl.current)
        qrObjectUrl.current = null
      }
    }
  }, [toast])

  async function verify() {
    if (userId == null) return
    if (code.length !== 6 || !/^\d+$/.test(code)) {
      toast({ variant: 'destructive', title: 'Enter the 6-digit code from your authenticator app.' })
      return
    }
    setBusy(true)
    try {
      await twoFactorAuthService.verifyTwoFactor(userId, code)
      await twoFactorAuthService.enableTwoFactorAfterVerify()
      toast({ title: 'Two-factor authentication enabled' })
      navigate('/dashboard')
    } catch (e) {
      toast({
        variant: 'destructive',
        title: 'Verification failed',
        description: e instanceof Error ? e.message : undefined,
      })
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="max-w-lg space-y-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Two-factor authentication</h1>
        <p className="text-muted-foreground text-sm">
          Scan the QR code, then enter the verification code once to finish enrolling (legacy Headwind flow).
        </p>
      </div>
      {loading ? (
        <div className="text-muted-foreground flex items-center gap-2 text-sm">
          <Loader2 className="h-4 w-4 animate-spin" />
          Loading…
        </div>
      ) : qrUrl ? (
        <img src={qrUrl} alt="Two-factor QR" className="max-w-xs rounded border p-2" />
      ) : (
        <p className="text-muted-foreground text-sm">QR not available.</p>
      )}
      <div className="flex max-w-xs flex-col gap-2">
        <Input inputMode="numeric" maxLength={6} placeholder="6-digit code" value={code} onChange={(e) => setCode(e.target.value)} />
        <Button type="button" disabled={busy} onClick={() => void verify()}>
          {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Verify and continue'}
        </Button>
      </div>
    </div>
  )
}
