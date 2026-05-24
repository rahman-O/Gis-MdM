import type { ProfileMeta } from '@/features/profiles/types'

interface Props {
  meta: ProfileMeta | null
}

export function ProfileUsagePanel({ meta }: Props) {
  if (!meta) return null
  return (
    <div className="rounded-md border bg-muted/30 px-4 py-3 text-sm">
      <p className="font-medium">Usage</p>
      <ul className="mt-1 text-muted-foreground space-y-0.5">
        <li>{meta.deviceCount ?? 0} enrolled devices (via enrollment routes)</li>
        <li>{meta.enrollmentRouteCount ?? 0} enrollment routes</li>
        {meta.publishedVersion != null ? (
          <li>Published version: v{meta.publishedVersion}</li>
        ) : (
          <li>No published version yet</li>
        )}
      </ul>
    </div>
  )
}
