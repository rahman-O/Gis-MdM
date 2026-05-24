import {
  FolderOpen,
  LayoutDashboard,
  LayoutList,
  Smartphone,
  Users,
  Settings,
  AppWindow,
  Shield,
  FileArchive,
  ImageIcon,
  Lightbulb,
  RefreshCw,
  Radio,
  Route,
} from 'lucide-react'
import type { NavItem } from './types'

/**
 * Mirrors legacy tab visibility (`content.html`): SUMMARY/DEVICES always;
 * FILES needs `files`; settings submenu items live under `settings`.
 */
export const NAV_ITEMS: NavItem[] = [
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { label: 'Devices', path: '/devices', icon: Smartphone },
  { label: 'Applications', path: '/applications', icon: AppWindow, permission: 'applications' },
  { labelKey: 'nav.profiles', label: 'Profiles', path: '/profiles', icon: LayoutList, permission: 'configurations' },
  {
    labelKey: 'nav.enrollmentRoutes',
    label: 'Enrollment routes',
    path: '/enrollment-routes',
    icon: Route,
    permission: 'configurations',
  },
  { label: 'Files', path: '/files', icon: FileArchive, permission: 'files' },
  { label: 'Groups', path: '/groups', icon: FolderOpen, permission: 'settings' },
  { label: 'Users', path: '/users', icon: Users, permission: 'settings' },
  {
    label: 'Roles',
    path: '/roles',
    icon: Shield,
    permission: 'settings',
    requiresCanManageRoles: true,
  },
  { label: 'Settings', path: '/settings', icon: Settings, permission: 'settings' },
  { label: 'Icons', path: '/icons', icon: ImageIcon, permission: 'settings' },
  { label: 'Hints', path: '/hints', icon: Lightbulb, permission: 'settings' },
  { label: 'Updates', path: '/updates', icon: RefreshCw },
  { label: 'Push', path: '/push', icon: Radio, permission: 'push_api' },
  {
    label: 'Plugins',
    path: '/plugin-settings',
    icon: Settings,
    permission: 'plugins_customer_access_management',
  },
]
