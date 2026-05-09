import { Menu } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { useAuth } from '@/shared/hooks/useAuth'

interface HeaderProps {
  onMenuClick: () => void
}

export function Header({ onMenuClick }: HeaderProps) {
  const { username, logout } = useAuth()

  return (
    <header className="flex h-14 items-center justify-between border-b bg-background px-4">
      <div className="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={onMenuClick}
          aria-label="Open navigation menu"
        >
          <Menu className="h-5 w-5" />
        </Button>
        <span className="font-semibold text-sm hidden md:block">Headwind MDM</span>
      </div>

      <div className="flex items-center gap-3">
        {username && (
          <span className="text-sm text-muted-foreground">{username}</span>
        )}
        <Button variant="outline" size="sm" onClick={logout} aria-label="Sign out">
          Sign Out
        </Button>
      </div>
    </header>
  )
}
