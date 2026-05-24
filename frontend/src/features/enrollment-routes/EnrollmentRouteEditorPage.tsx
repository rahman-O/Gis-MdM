import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { AlertCircle, Info } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { getDeviceTree, type TreeNode } from '@/features/device-tree/deviceTreeService'
import { EnrollmentRouteQrPanel } from '@/features/enrollment-routes/EnrollmentRouteQrPanel'
import * as routeService from '@/features/enrollment-routes/enrollmentRouteService'
import type { EnrollmentRouteQrMeta } from '@/features/enrollment-routes/enrollmentRouteService'
import { getOnboardingStatus } from '@/features/onboarding/onboardingService'

export function EnrollmentRouteEditorPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const params = useParams<{ id: string }>()
  const isNew = params.id === 'new'
  const routeId = isNew ? 0 : Number(params.id)

  const [loading, setLoading] = useState(!isNew)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saveError, setSaveError] = useState<string | null>(null)

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [defaultTreeNodeId, setDefaultTreeNodeId] = useState<number | ''>('')
  const [defaultDeviceIdMode, setDefaultDeviceIdMode] = useState('imei')
  const [mainAppId, setMainAppId] = useState<number | ''>('')

  const [treeNodes, setTreeNodes] = useState<TreeNode[]>([])
  const [qrMeta, setQrMeta] = useState<EnrollmentRouteQrMeta | null>(null)

  useEffect(() => {
    if (!isNew) return
    void getOnboardingStatus()
      .then((s) => {
        if (!s.hasPublishedProfile) {
          navigate('/profiles', {
            replace: true,
            state: { onboardingHint: 'Publish a profile before creating an enrollment route.' },
          })
        }
      })
      .catch(() => {
        /* allow editor if status unavailable */
      })
  }, [isNew, navigate])

  useEffect(() => {
    void getDeviceTree()
      .then((tree) => {
        setTreeNodes(Array.isArray(tree.nodes) ? tree.nodes : [])
      })
      .catch(() => {
        setSaveError('Failed to load tree folders.')
      })
  }, [])

  useEffect(() => {
    if (isNew || !Number.isFinite(routeId) || routeId <= 0) {
      setLoading(false)
      return
    }
    setLoading(true)
    setError(null)
    void Promise.all([routeService.getEnrollmentRoute(routeId), routeService.getEnrollmentRouteQrMeta(routeId)])
      .then(([route, meta]) => {
        setName(route.name)
        setDescription(route.description ?? '')
        setDefaultTreeNodeId(route.defaultTreeNodeId ?? '')
        setDefaultDeviceIdMode(route.defaultDeviceIdMode || 'imei')
        setMainAppId(route.mainAppId ?? '')
        setQrMeta(meta)
      })
      .catch((e: unknown) => {
        setError(e instanceof Error ? e.message : 'Failed to load enrollment route.')
      })
      .finally(() => setLoading(false))
  }, [isNew, routeId])

  const validate = (): string | null => {
    if (!name.trim()) return 'Name is required.'
    if (!defaultTreeNodeId || defaultTreeNodeId <= 0) {
      return 'Default tree folder is required.'
    }
    if (!mainAppId || mainAppId <= 0) return 'Main app is required for QR enrollment.'
    return null
  }

  const handleSave = async () => {
    const validationError = validate()
    if (validationError) {
      setSaveError(validationError)
      return
    }
    setSaving(true)
    setSaveError(null)
    try {
      const payload = {
        name: name.trim(),
        description: description.trim() || null,
        defaultTreeNodeId: Number(defaultTreeNodeId),
        defaultDeviceIdMode,
        mainAppId: Number(mainAppId),
      }
      if (isNew) {
        const created = await routeService.createEnrollmentRoute(payload)
        navigate(`/enrollment-routes/${created.id}`, { replace: true })
      } else {
        await routeService.updateEnrollmentRoute(routeId, payload)
        const meta = await routeService.getEnrollmentRouteQrMeta(routeId)
        setQrMeta(meta)
      }
    } catch (e: unknown) {
      setSaveError(e instanceof Error ? e.message : 'Failed to save enrollment route.')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="text-sm text-muted-foreground">Loading enrollment route…</div>
  }

  if (error) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
          <AlertCircle className="h-4 w-4" />
          <span>{error}</span>
        </div>
        <Button variant="outline" onClick={() => navigate('/enrollment-routes')}>
          Back
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">
            {isNew ? 'New enrollment route' : 'Edit enrollment route'}
          </h1>
          <p className="text-sm text-muted-foreground">{t('enrollmentRoute.help.profileOnly')}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate('/enrollment-routes')}>
            Cancel
          </Button>
          <Button disabled={saving} onClick={() => void handleSave()}>
            {saving ? 'Saving…' : 'Save'}
          </Button>
        </div>
      </div>

      <div className="flex gap-2 rounded-md border border-blue-500/30 bg-blue-50/80 p-3 text-sm text-blue-950 dark:bg-blue-950/30 dark:text-blue-100">
        <Info className="mt-0.5 h-4 w-4 shrink-0" />
        <p>
          Device policy is delivered from <strong>tree folder assignments</strong> in Profiles, not
          from this enrollment route. Open{' '}
          <Link className="font-medium underline" to="/profiles">
            Profiles
          </Link>{' '}
          and use the <strong>Assignments</strong> section to bind a published version to a folder.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Enrollment settings</CardTitle>
          <CardDescription>QR provisioning, default folder placement, and main app.</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="route-name">Name</Label>
            <Input id="route-name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="route-desc">Description</Label>
            <Input id="route-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label>Default tree folder</Label>
            <Select
              value={defaultTreeNodeId === '' ? '' : String(defaultTreeNodeId)}
              onValueChange={(v) => setDefaultTreeNodeId(Number(v))}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select folder" />
              </SelectTrigger>
              <SelectContent>
                {treeNodes.map((n) => (
                  <SelectItem key={n.id} value={String(n.id)}>
                    {n.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Default device id mode</Label>
            <Select value={defaultDeviceIdMode} onValueChange={setDefaultDeviceIdMode}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="imei">IMEI (default)</SelectItem>
                <SelectItem value="serial">Serial</SelectItem>
                <SelectItem value="request">Request</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label>Main app version id</Label>
            <Input
              type="number"
              value={mainAppId === '' ? '' : String(mainAppId)}
              onChange={(e) => setMainAppId(e.target.value === '' ? '' : Number(e.target.value))}
            />
          </div>
          {saveError ? <p className="text-sm text-destructive md:col-span-2">{saveError}</p> : null}
        </CardContent>
      </Card>

      {qrMeta?.qrcodekey ? (
        <Card>
          <CardHeader>
            <CardTitle>Enrollment QR</CardTitle>
            <CardDescription>Provisioning QR for this route.</CardDescription>
          </CardHeader>
          <CardContent>
            <EnrollmentRouteQrPanel meta={qrMeta} />
          </CardContent>
        </Card>
      ) : null}
    </div>
  )
}
