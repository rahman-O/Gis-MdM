import { Menu } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '@/shared/ui/button'
import { useAuth } from '@/shared/hooks/useAuth'
import { setUiLanguage } from '@/i18n/config'

interface HeaderProps {
  onMenuClick: () => void
}

export function Header({ onMenuClick }: HeaderProps) {
  const { t } = useTranslation()
  const { username, logout } = useAuth()

  return (
    <header className="flex h-12 items-center justify-between border-b bg-background px-3">
      <div className="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={onMenuClick}
          aria-label="Open navigation menu"
        >
          <Menu className="h-4 w-4" />
        </Button>
        <span className="hidden text-sm font-semibold md:block">{t('brand.title')}</span>
      </div>

      <div className="flex items-center gap-2">
        <div className="hidden items-center gap-1 md:flex">
          <Button variant="ghost" size="sm" type="button" onClick={() => setUiLanguage('en')}>
            EN
          </Button>
          <Button variant="ghost" size="sm" type="button" onClick={() => setUiLanguage('ar')}>
            AR
          </Button>
        </div>
        <Button variant="link" size="sm" className="text-muted-foreground h-auto px-2 py-0" asChild>
          <Link to="/profile">Profile</Link>
        </Button>
        {username && <span className="text-muted-foreground hidden text-sm sm:inline">{username}</span>}
        <Button variant="outline" size="sm" onClick={logout} aria-label="Sign out">
          Sign Out
        </Button>
      </div>
    </header>
  )
}
