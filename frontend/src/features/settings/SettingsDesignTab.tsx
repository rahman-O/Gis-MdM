import { useCallback, useEffect, useRef, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Textarea } from '@/shared/ui/textarea'
import { useToast } from '@/shared/hooks/use-toast'
import * as settingsService from '@/features/settings/settingsService'

const ICON_SIZES = ['SMALL', 'MEDIUM', 'LARGE'] as const

const DESKTOP_HEADERS = [
  'NO_HEADER',
  'DEVICE_ID',
  'DESCRIPTION',
  'TEMPLATE',
  'CUSTOM1',
  'CUSTOM2',
  'CUSTOM3',
] as const

export function SettingsDesignTab() {
  const { toast } = useToast()
  const snapshot = useRef<Record<string, unknown>>({})
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [backgroundColor, setBackgroundColor] = useState('')
  const [textColor, setTextColor] = useState('')
  const [backgroundImageUrl, setBackgroundImageUrl] = useState('')
  const [iconSize, setIconSize] = useState<string>('SMALL')
  const [desktopHeader, setDesktopHeader] = useState<string>('NO_HEADER')
  const [desktopHeaderTemplate, setDesktopHeaderTemplate] = useState('')

  const [custom1, setCustom1] = useState('')
  const [custom2, setCustom2] = useState('')
  const [custom3, setCustom3] = useState('')

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const raw = await settingsService.fetchRawSettings()
      snapshot.current = raw
      setBackgroundColor(String(raw.backgroundColor ?? ''))
      setTextColor(String(raw.textColor ?? ''))
      setBackgroundImageUrl(String(raw.backgroundImageUrl ?? ''))
      setIconSize(String(raw.iconSize ?? 'SMALL'))
      setDesktopHeader(String(raw.desktopHeader ?? 'NO_HEADER'))
      setDesktopHeaderTemplate(String(raw.desktopHeaderTemplate ?? ''))
      setCustom1(String(raw.customPropertyName1 ?? ''))
      setCustom2(String(raw.customPropertyName2 ?? ''))
      setCustom3(String(raw.customPropertyName3 ?? ''))
    } catch {
      toast({ variant: 'destructive', title: 'Failed to load design settings' })
    } finally {
      setLoading(false)
    }
  }, [toast])

  useEffect(() => {
    void load()
  }, [load])

  async function save() {
    setBusy(true)
    try {
      const merged: Record<string, unknown> = {
        ...snapshot.current,
        backgroundColor,
        textColor,
        backgroundImageUrl,
        iconSize,
        desktopHeader,
        desktopHeaderTemplate,
      }
      await settingsService.saveDefaultDesign(merged)
      snapshot.current = merged
      toast({ title: 'Design settings saved' })
    } catch (e) {
      toast({
        variant: 'destructive',
        title: 'Save failed',
        description: e instanceof Error ? e.message : undefined,
      })
    } finally {
      setBusy(false)
    }
  }

  if (loading) {
    return (
      <div className="text-muted-foreground flex items-center gap-2 text-sm">
        <Loader2 className="h-4 w-4 animate-spin" />
        Loading design…
      </div>
    )
  }

  const headerChoices = [...DESKTOP_HEADERS].filter((h) => {
    if (h === 'CUSTOM1') return Boolean(custom1)
    if (h === 'CUSTOM2') return Boolean(custom2)
    if (h === 'CUSTOM3') return Boolean(custom3)
    return true
  })

  return (
    <div className="max-w-2xl space-y-4">
      <p className="text-muted-foreground text-sm">
        Default launcher appearance mirrored from the legacy Settings → Design tab (POST `/private/settings/design`).
      </p>
      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label className="text-muted-foreground text-xs">Background color</label>
          <Input value={backgroundColor} onChange={(e) => setBackgroundColor(e.target.value)} placeholder="#FFFFFF" />
        </div>
        <div>
          <label className="text-muted-foreground text-xs">Launcher text color</label>
          <Input value={textColor} onChange={(e) => setTextColor(e.target.value)} placeholder="#000000" />
        </div>
      </div>
      <div>
        <label className="text-muted-foreground text-xs">Background image URL</label>
        <Input
          value={backgroundImageUrl}
          onChange={(e) => setBackgroundImageUrl(e.target.value)}
          placeholder="https://…"
        />
      </div>
      <div>
        <label className="text-muted-foreground text-xs">Icon size</label>
        <Select value={iconSize} onValueChange={setIconSize}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {ICON_SIZES.map((s) => (
              <SelectItem key={s} value={s}>
                {s}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div>
        <label className="text-muted-foreground text-xs">Desktop header</label>
        <Select value={desktopHeader} onValueChange={setDesktopHeader}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {headerChoices.map((h) => (
              <SelectItem key={h} value={h}>
                {h === 'CUSTOM1' ? custom1 : h === 'CUSTOM2' ? custom2 : h === 'CUSTOM3' ? custom3 : h}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      {desktopHeader === 'TEMPLATE' ? (
        <div>
          <label className="text-muted-foreground text-xs">Header template tokens</label>
          <Textarea
            value={desktopHeaderTemplate}
            onChange={(e) => setDesktopHeaderTemplate(e.target.value)}
            placeholder="deviceId, description, custom1, custom2, custom3 …"
            rows={3}
          />
        </div>
      ) : null}
      <Button type="button" variant="outline" disabled={busy} onClick={() => void save()}>
        {busy ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
        Save design
      </Button>
    </div>
  )
}
