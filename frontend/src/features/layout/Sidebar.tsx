import { NavLink } from 'react-router-dom'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/shared/ui/sheet'
import { cn } from '@/shared/utils/cn'
import { canManageRoles, hasPermission, isSuperAdmin } from '@/features/auth/permissions'
import { NAV_ITEMS } from './navItems'

interface SidebarProps {
  mobileOpen: boolean
  onMobileClose: () => void
}

function SidebarNav({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <nav aria-label="Main navigation" className="flex flex-col gap-1 p-2.5">
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
          className={({ isActive }) =>
            cn(
              'flex min-h-[36px] items-center gap-2 rounded-md px-2.5 py-1.5 text-sm font-medium transition-colors',
              isActive
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
            )
          }
          aria-current={undefined}
        >
          {({ isActive }) => (
            <>
              <item.icon className="h-3.5 w-3.5 shrink-0" />
              <span>{item.label}</span>
              {isActive && <span className="sr-only">(current page)</span>}
            </>
          )}
        </NavLink>
      ))}
    </nav>
  )
}

export function Sidebar({ mobileOpen, onMobileClose }: SidebarProps) {
  return (
    <>
      {/* Desktop sidebar */}
      <aside className="hidden md:flex md:w-56 md:flex-col md:border-r md:bg-background">
        <div className="flex h-12 items-center border-b px-3">
          <span className="font-bold text-sm">Headwind MDM</span>
        </div>
        <SidebarNav />
      </aside>

      {/* Mobile drawer */}
      <Sheet open={mobileOpen} onOpenChange={(open) => !open && onMobileClose()}>
        <SheetContent side="left" className="w-56 p-0">
          <SheetHeader className="border-b px-3 py-2.5">
            <SheetTitle className="text-sm font-bold">Headwind MDM</SheetTitle>
          </SheetHeader>
          <SidebarNav onNavigate={onMobileClose} />
        </SheetContent>
      </Sheet>
    </>
  )
}
