import { cn } from '@/shared/utils/cn'
import { Button } from '@/shared/ui/button'
import {
  type ProfileWorkspaceSection,
  useProfileWorkspace,
} from '@/features/profiles/workspace/profileWorkspaceState'

const SECTIONS: { key: ProfileWorkspaceSection; label: string }[] = [
  { key: 'overview', label: 'Overview' },
  { key: 'assignments', label: 'Assignments' },
  { key: 'rollout', label: 'Rollout' },
  { key: 'versions', label: 'Versions' },
  { key: 'editor', label: 'Editor' },
  { key: 'activity', label: 'Activity' },
]

export function ProfileWorkspaceSidebar() {
  const { section, setSection } = useProfileWorkspace()

  return (
    <nav className="flex w-full shrink-0 flex-row gap-1 overflow-x-auto border-b p-2 md:w-48 md:flex-col md:border-b-0 md:border-r md:overflow-x-visible">
      {SECTIONS.map((s) => (
        <Button
          key={s.key}
          type="button"
          variant={section === s.key ? 'secondary' : 'ghost'}
          className={cn('justify-start whitespace-nowrap', section === s.key && 'font-medium')}
          onClick={() => setSection(s.key)}
        >
          {s.label}
        </Button>
      ))}
    </nav>
  )
}
