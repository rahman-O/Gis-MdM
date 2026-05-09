import type React from 'react'

export interface NavItem {
  label: string
  path: string
  icon: React.ComponentType<{ className?: string }>
  permission?: string
}
