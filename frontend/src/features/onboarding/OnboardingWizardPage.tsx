import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { CheckCircle2, Circle } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  getOnboardingStatus,
  type OnboardingStatus,
} from '@/features/onboarding/onboardingService'

const WIZARD_STEPS = [
  { id: 'tree', title: 'Device folder', description: 'Create at least one folder under the customer root in the device tree.' },
  { id: 'profile', title: 'Profile', description: 'Create a profile with restrictions, apps, and design settings.' },
  { id: 'publish', title: 'Publish profile', description: 'Publish the profile so enrollment routes can bind to it.' },
  { id: 'route', title: 'Enrollment route', description: 'Create a route that links a tree folder to the published profile.' },
  { id: 'qr', title: 'Test QR', description: 'Scan the route QR code and enroll a test device into the folder.' },
] as const

export function OnboardingWizardPage() {
  const [status, setStatus] = useState<OnboardingStatus | null>(null)
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      setStatus(await getOnboardingStatus())
    } catch {
      setStatus(null)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const stepDone = (id: string): boolean => {
    if (!status) return false
    switch (id) {
      case 'tree':
        return status.hasTreeBeyondRoot
      case 'profile':
        return status.steps.some((s) => s.id === 'profile' && s.done) || status.hasPublishedProfile
      case 'publish':
        return status.hasPublishedProfile
      case 'route':
        return status.hasEnrollmentRoute
      case 'qr':
        return status.steps.some((s) => s.id === 'qr' && s.done)
      default:
        return false
    }
  }

  const stepPath = (id: string): string => {
    switch (id) {
      case 'tree':
        return '/devices'
      case 'profile':
      case 'publish':
        return '/profiles'
      case 'route':
      case 'qr':
        return '/enrollment-routes'
      default:
        return '/dashboard'
    }
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Setup wizard</h1>
        <p className="text-muted-foreground text-sm mt-1">
          Follow these steps to configure the device control plane: tree, profile, publish, route, and QR test.
        </p>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[1, 2, 3, 4, 5].map((i) => (
            <Skeleton key={i} className="h-20 w-full" />
          ))}
        </div>
      ) : (
        <div className="space-y-3">
          {WIZARD_STEPS.map((step, index) => {
            const done = stepDone(step.id)
            return (
              <Card key={step.id} className={done ? 'opacity-80' : ''}>
                <CardHeader className="py-4">
                  <div className="flex items-start gap-3">
                    {done ? (
                      <CheckCircle2 className="h-5 w-5 text-emerald-600 shrink-0" />
                    ) : (
                      <Circle className="h-5 w-5 text-muted-foreground shrink-0" />
                    )}
                    <div className="flex-1 min-w-0">
                      <CardTitle className="text-base">
                        {index + 1}. {step.title}
                      </CardTitle>
                      <CardDescription className="mt-1">{step.description}</CardDescription>
                    </div>
                    {!done ? (
                      <Button variant="outline" size="sm" asChild>
                        <Link to={stepPath(step.id)}>Continue</Link>
                      </Button>
                    ) : null}
                  </div>
                </CardHeader>
              </Card>
            )
          })}
        </div>
      )}

      {status?.complete ? (
        <Card className="border-emerald-500/30 bg-emerald-500/5">
          <CardContent className="py-4 text-sm text-emerald-800 dark:text-emerald-200">
            Setup complete. You can manage devices, profiles, and enrollment routes from the sidebar.
          </CardContent>
        </Card>
      ) : null}

      <div className="flex gap-2">
        <Button variant="outline" asChild>
          <Link to="/dashboard">Back to dashboard</Link>
        </Button>
        <Button variant="ghost" onClick={() => void load()} disabled={loading}>
          Refresh status
        </Button>
      </div>
    </div>
  )
}
