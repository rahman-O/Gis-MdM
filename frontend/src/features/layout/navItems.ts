import { FolderOpen, LayoutDashboard, LayoutList, Smartphone, Users, Settings, AppWindow } from 'lucide-react'
import type { NavItem } from './types'

export const NAV_ITEMS: NavItem[] = [
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { label: 'Devices',         path: '/devices',         icon: Smartphone },
  { label: 'Applications', path: '/applications', icon: AppWindow, permission: 'applications' },
  { label: 'Groups', path: '/groups', icon: FolderOpen },
  { label: 'Configurations', path: '/configurations', icon: LayoutList, permission: 'configurations' },
  { label: 'Users',           path: '/users',           icon: Users },
  { label: 'Settings',  path: '/settings',  icon: Settings },
]
