import { useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { useToast } from '@/shared/hooks/use-toast'
import * as pushService from '@/features/push/pushService'

export function PushPage() {
  const { toast } = useToast()
  const [messageType, setMessageType] = useState('config_updated')
  const [payload, setPayload] = useState('')
  const [broadcast, setBroadcast] = useState(false)
  const [deviceNumbers, setDeviceNumbers] = useState('')
  const [groups, setGroups] = useState('')
  const [busy, setBusy] = useState(false)

  async function send() {
    setBusy(true)
    try {
      await pushService.sendPush({
        messageType: messageType.trim(),
        payload: payload.trim(),
        broadcast,
        deviceNumbers:
          broadcast || !deviceNumbers.trim()
            ? undefined
            : deviceNumbers.split(/[\s,]+/).filter(Boolean),
        groups:
          broadcast || !groups.trim() ? undefined : groups.split(/[\s,]+/).filter(Boolean),
      })
      toast({ title: 'Push queued' })
    } catch {
      toast({ title: 'Push failed', variant: 'destructive' })
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="max-w-xl space-y-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Push</h1>
        <p className="text-muted-foreground text-sm">Requires `push_api` permission.</p>
      </div>
      <div>
        <label className="text-xs text-muted-foreground">Message type</label>
        <Input value={messageType} onChange={(e) => setMessageType(e.target.value)} />
      </div>
      <div>
        <label className="text-xs text-muted-foreground">Payload</label>
        <Textarea rows={4} value={payload} onChange={(e) => setPayload(e.target.value)} />
      </div>
      <label className="flex items-center gap-2">
        <Checkbox checked={broadcast} onCheckedChange={(c) => setBroadcast(c === true)} />
        Broadcast to all filtered devices for this tenant
      </label>
      <div>
        <label className="text-xs text-muted-foreground">Device numbers (comma or space separated)</label>
        <Input value={deviceNumbers} onChange={(e) => setDeviceNumbers(e.target.value)} disabled={broadcast} />
      </div>
      <div>
        <label className="text-xs text-muted-foreground">Group names (comma separated)</label>
        <Input value={groups} onChange={(e) => setGroups(e.target.value)} disabled={broadcast} />
      </div>
      <Button type="button" onClick={() => void send()} disabled={busy}>
        Send
      </Button>
    </div>
  )
}
