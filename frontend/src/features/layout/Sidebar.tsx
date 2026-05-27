import { NavLink, useLocation } from 'react-router-dom'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/shared/ui/sheet'
import { cn } from '@/shared/utils/cn'
import { canManageRoles, hasPermission, isSuperAdmin } from '@/features/auth/permissions'
import { NAV_ITEMS } from './navItems'

interface SidebarProps {
  mobileOpen: boolean
  collapsed: boolean
  onMobileClose: () => void
  onToggleCollapse: () => void
}

function SidebarNav({ collapsed, onNavigate }: { collapsed: boolean; onNavigate?: () => void }) {
  const { t } = useTranslation()
  return (
    <nav aria-label="Main navigation" className="flex flex-col gap-0.5 p-2">
      {NAV_ITEMS.filter((item) => {
        if (item.requiresSuperAdmin && !isSuperAdmin()) return false
        if (item.requiresCanManageRoles && !canManageRoles()) return false
        if (item.anyPermission?.length) {
          if (!item.anyPermission.some((p) => hasPermission(p))) return false
        } else if (item.permission && !hasPermission(item.permission)) return false
        return true
      }).map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          onClick={onNavigate}
          title={collapsed ? (item.labelKey ? t(item.labelKey) : item.label) : undefined}
          className={({ isActive }) =>
            cn(
              'flex items-center rounded-md text-sm font-medium transition-all duration-200',
              collapsed ? 'h-9 w-9 justify-center' : 'h-9 gap-2.5 px-2.5',
              isActive
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
            )
          }
          aria-current={undefined}
        >
          {({ isActive }) => (
            <>
              <item.icon className="h-4 w-4 shrink-0" />
              {!collapsed && <span className="truncate">{item.labelKey ? t(item.labelKey) : item.label}</span>}
              {isActive && <span className="sr-only">(current page)</span>}
            </>
          )}
        </NavLink>
      ))}
    </nav>
  )
}

export function Sidebar({ mobileOpen, collapsed, onMobileClose, onToggleCollapse }: SidebarProps) {
  const location = useLocation()

  // Auto-collapse sidebar when navigating to a new page
  useEffect(() => {
    if (!collapsed) {
      onToggleCollapse()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  return (
    <>
      {/* Desktop sidebar */}
      <aside
        className={cn(
          'hidden md:flex md:flex-col md:border-r md:bg-background transition-all duration-300 ease-in-out',
          collapsed ? 'md:w-[52px]' : 'md:w-52'
        )}
      >
        <div className={cn(
          'flex h-12 items-center border-b transition-all duration-300',
          collapsed ? 'justify-center px-1' : 'justify-between px-3'
        )}>
          {!collapsed && <span className="font-bold text-sm truncate">Headwind MDM</span>}
          <button
            type="button"
            onClick={onToggleCollapse}
            className="flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
            aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
          >
            {collapsed ? <PanelLeftOpen className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
          </button>
        </div>
        <SidebarNav collapsed={collapsed} />
      </aside>

      {/* Mobile drawer */}
      <Sheet open={mobileOpen} onOpenChange={(open) => !open && onMobileClose()}>
        <SheetContent side="left" className="w-52 p-0">
          <SheetHeader className="border-b px-3 py-2.5">
            <SheetTitle className="text-sm font-bold">Headwind MDM</SheetTitle>
          </SheetHeader>
          <SidebarNav collapsed={false} onNavigate={onMobileClose} />
        </SheetContent>
      </Sheet>
    </>
  )
}
