import { Navigate, Route, Routes } from 'react-router-dom'
import { LoginPage } from '@/features/auth/LoginPage'
import { DashboardPage } from '@/features/dashboard/DashboardPage'
import { DevicesPage } from '@/features/devices/DevicesPage'
import { GroupsPage } from '@/features/groups/GroupsPage'
import { UsersPage } from '@/features/users/UsersPage'
import { RolesPage } from '@/features/roles/RolesPage'
import { SettingsPage } from '@/features/settings/SettingsPage'
import { ConfigurationsPage } from '@/features/configurations/ConfigurationsPage'
import { ConfigurationEditorPage } from '@/features/configurations/ConfigurationEditorPage'
import { EnrollmentQrPage } from '@/features/devices/EnrollmentQrPage'
import { ApplicationsPage } from '@/features/applications/ApplicationsPage'
import { AdminApplicationsPage } from '@/features/applications/AdminApplicationsPage'
import { ApplicationVersionsPage } from '@/features/applications/ApplicationVersionsPage'
import { AppLayout } from '@/features/layout/AppLayout'
import { AuthGuard } from '@/shared/components/AuthGuard'

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        element={
          <AuthGuard>
            <AppLayout />
          </AuthGuard>
        }
      >
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/devices" element={<DevicesPage />} />
        <Route path="/groups" element={<GroupsPage />} />
        <Route path="/applications" element={<ApplicationsPage />} />
        <Route path="/applications/admin" element={<AdminApplicationsPage />} />
        <Route path="/application/:id/versions" element={<ApplicationVersionsPage />} />
        <Route path="/configurations" element={<ConfigurationsPage />} />
        <Route path="/configurations/:id/edit" element={<ConfigurationEditorPage />} />
        <Route path="/qr/:qrCodeKey/:deviceId" element={<EnrollmentQrPage />} />
        <Route path="/qr/:qrCodeKey" element={<EnrollmentQrPage />} />
        <Route path="/users" element={<UsersPage />} />
        <Route path="/roles" element={<RolesPage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Route>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}
