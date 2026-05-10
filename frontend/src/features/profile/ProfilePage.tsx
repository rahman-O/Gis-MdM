import { useEffect, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { applySessionFromUserPayload } from '@/features/auth/session'
import { useToast } from '@/shared/hooks/use-toast'
import type { LoginUserPayload } from '@/features/auth/types'
import * as profileService from '@/features/profile/profileService'

export function ProfilePage() {
  const { toast } = useToast()
  const [user, setUser] = useState<(LoginUserPayload & Record<string, unknown>) | null>(null)
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [oldPw, setOldPw] = useState('')
  const [newPw, setNewPw] = useState('')
  const [confirmPw, setConfirmPw] = useState('')
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    profileService.fetchCurrentUser().then((u) => {
      setUser(u)
      setName(String(u?.name ?? ''))
      setEmail(String(u?.email ?? ''))
    })
  }, [])

  async function saveDetails() {
    const uid = Number(user?.id)
    if (!Number.isFinite(uid)) return
    setBusy(true)
    try {
      await profileService.updateUserDetails({ id: uid, name: name.trim(), email: email.trim() })
      toast({ title: 'Profile updated' })
      const refreshed = await profileService.fetchCurrentUser()
      setUser(refreshed)
      applySessionFromUserPayload(refreshed)
    } catch {
      toast({ variant: 'destructive', title: 'Save failed' })
    } finally {
      setBusy(false)
    }
  }

  async function savePassword() {
    const uid = Number(user?.id)
    if (!Number.isFinite(uid)) return
    if (!newPw || newPw !== confirmPw) {
      toast({ variant: 'destructive', title: 'Passwords do not match.' })
      return
    }
    setBusy(true)
    try {
      await profileService.updateCurrentPassword({
        id: uid,
        login: user?.login as string | undefined,
        oldPasswordPlain: oldPw,
        newPasswordPlain: newPw,
      })
      toast({ title: 'Password updated' })
      setOldPw('')
      setNewPw('')
      setConfirmPw('')
    } catch {
      toast({ variant: 'destructive', title: 'Password change failed' })
    } finally {
      setBusy(false)
    }
  }

  if (!user) {
    return (
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Loader2 className="h-4 w-4 animate-spin" />
        Loading profile…
      </div>
    )
  }

  return (
    <div className="max-w-lg space-y-8">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Profile</h1>
        <p className="text-muted-foreground mt-1 text-sm">User #{user.id ?? '—'}</p>
      </div>
      <div className="space-y-3">
        <h2 className="text-sm font-medium">Details</h2>
        <div>
          <label className="text-xs text-muted-foreground">Name</label>
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div>
          <label className="text-xs text-muted-foreground">Email</label>
          <Input value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <Button type="button" variant="outline" disabled={busy} onClick={() => void saveDetails()}>
          Save details
        </Button>
      </div>
      <div className="space-y-3">
        <h2 className="text-sm font-medium">Change password</h2>
        <p className="text-muted-foreground text-xs">Uses legacy MD5 client hashes like the Angular profile screen.</p>
        <Input type="password" placeholder="Current password" value={oldPw} onChange={(e) => setOldPw(e.target.value)} />
        <Input type="password" placeholder="New password" value={newPw} onChange={(e) => setNewPw(e.target.value)} />
        <Input
          type="password"
          placeholder="Confirm new password"
          value={confirmPw}
          onChange={(e) => setConfirmPw(e.target.value)}
        />
        <Button type="button" disabled={busy} onClick={() => void savePassword()}>
          Update password
        </Button>
      </div>
    </div>
  )
}
