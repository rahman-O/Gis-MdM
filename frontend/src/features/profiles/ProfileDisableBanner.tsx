import { Button } from '@/shared/ui/button'
import { disableProfile, enableProfile } from '@/features/profiles/profileRolloutService'

interface Props {
  profileId: number
  enabled: boolean
  onChanged: (enabled: boolean) => void
}

export function ProfileDisableBanner({ profileId, enabled, onChanged }: Props) {
  const toggle = async () => {
    const next = !enabled
    const msg = next
      ? 'Enable this profile? Devices will be marked pending for rollout.'
      : 'Disable this profile? Sync will not push this policy until re-enabled.'
    if (!window.confirm(msg)) return
    try {
      if (next) {
        await enableProfile(profileId)
      } else {
        await disableProfile(profileId)
      }
      onChanged(next)
    } catch (e: unknown) {
      window.alert(e instanceof Error ? e.message : 'Update failed')
    }
  }

  return (
    <div
      className={
        enabled
          ? 'flex flex-wrap items-center justify-between gap-2 rounded-md border p-3 text-sm'
          : 'flex flex-wrap items-center justify-between gap-2 rounded-md border border-amber-500/50 bg-amber-500/10 p-3 text-sm'
      }
    >
      <span>
        {enabled
          ? 'Profile is active. Devices receive policy when assigned or enrolled via route.'
          : 'Profile is disabled. Policy is not pushed until you enable it again.'}
      </span>
      <Button type="button" variant={enabled ? 'outline' : 'default'} size="sm" onClick={() => void toggle()}>
        {enabled ? 'Disable profile' : 'Enable profile'}
      </Button>
    </div>
  )
}
