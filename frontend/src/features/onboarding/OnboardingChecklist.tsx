import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { CheckCircle2, Circle, ListChecks } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  getOnboardingStatus,
  type OnboardingStatus,
} from '@/features/onboarding/onboardingService'

export function OnboardingChecklist() {
  const [status, setStatus] = useState<OnboardingStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [hidden, setHidden] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await getOnboardingStatus()
      setStatus(data)
      if (data.complete) {
        setHidden(true)
      }
    } catch {
      setStatus(null)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  if (hidden || (!loading && status?.complete)) {
    return null
  }

  if (loading) {
    return (
      <Card>
        <CardHeader className="py-3">
          <Skeleton className="h-5 w-48" />
        </CardHeader>
        <CardContent className="space-y-2">
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-3/4" />
        </CardContent>
      </Card>
    )
  }

  if (!status) {
    return null
  }

  return (
    <Card className="border-primary/20 bg-primary/5">
      <CardHeader className="flex flex-row items-center justify-between py-3">
        <CardTitle className="text-sm font-semibold flex items-center gap-2">
          <ListChecks className="h-4 w-4 text-primary" />
          Complete your MDM setup
        </CardTitle>
        <Button variant="outline" size="sm" asChild>
          <Link to="/onboarding">Open wizard</Link>
        </Button>
      </CardHeader>
      <CardContent className="pt-0 pb-4 space-y-2">
        {status.steps.map((step) => (
          <div key={step.id} className="flex items-start gap-2 text-sm">
            {step.done ? (
              <CheckCircle2 className="h-4 w-4 text-emerald-600 shrink-0 mt-0.5" />
            ) : (
              <Circle className="h-4 w-4 text-muted-foreground shrink-0 mt-0.5" />
            )}
            <span className={step.done ? 'text-muted-foreground line-through' : ''}>{step.label}</span>
            {!step.done && step.path ? (
              <Link to={step.path} className="text-primary text-xs ml-auto shrink-0 hover:underline">
                Go
              </Link>
            ) : null}
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
