import { Navigate, Route, Routes } from 'react-router-dom'
import { LoginPage } from '@/features/auth/LoginPage'
import { PasswordRecoveryPage } from '@/features/auth/PasswordRecoveryPage'
import { PasswordResetPage } from '@/features/auth/PasswordResetPage'
import { SignupPage } from '@/features/auth/SignupPage'
import { SignupCompletePage } from '@/features/auth/SignupCompletePage'
import { DashboardPage } from '@/features/dashboard/DashboardPage'
import { DevicesPage } from '@/features/devices/DevicesPage'
import { GroupsPage } from '@/features/groups/GroupsPage'
import { UsersPage } from '@/features/users/UsersPage'
import { RolesPage } from '@/features/roles/RolesPage'
import { SettingsPage } from '@/features/settings/SettingsPage'
import { ProfilesPage } from '@/features/profiles/ProfilesPage'
import { ProfileEditRedirect } from '@/features/profiles/ProfileEditRedirect'
import { EnrollmentRouteListPage } from '@/features/enrollment-routes/EnrollmentRouteListPage'
import { OnboardingWizardPage } from '@/features/onboarding/OnboardingWizardPage'
import { ConfigurationEditorPage } from '@/features/configurations/ConfigurationEditorPage'
import { EnrollmentQrPage } from '@/features/devices/EnrollmentQrPage'
import { ApplicationsPage } from '@/features/applications/ApplicationsPage'
import { AdminApplicationsPage } from '@/features/applications/AdminApplicationsPage'
import { ApplicationVersionsPage } from '@/features/applications/ApplicationVersionsPage'
import { AppLayout } from '@/features/layout/AppLayout'
import { FilesPage } from '@/features/files/FilesPage'
import { IconsPage } from '@/features/icons/IconsPage'
import { HintsPage } from '@/features/hints/HintsPage'
import { UpdatesPage } from '@/features/updates/UpdatesPage'
import { PushPage } from '@/features/push/PushPage'
import { ControlPanelPage } from '@/features/customers/ControlPanelPage'
import { PluginSettingsPage } from '@/features/plugins/PluginSettingsPage'
import { ProfilePage } from '@/features/profile/ProfilePage'
import { TwofactorPage } from '@/features/auth/TwofactorPage'
import { MapsInfoPage } from '@/features/maps/MapsInfoPage'
import { AuthGuard } from '@/shared/components/AuthGuard'

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/password-recovery" element={<PasswordRecoveryPage />} />
      <Route path="/passwordRecovery" element={<PasswordRecoveryPage />} />
      <Route path="/signup" element={<SignupPage />} />
      <Route path="/signup-complete/:token" element={<SignupCompletePage />} />
      <Route path="/signupComplete/:token" element={<SignupCompletePage />} />

      <Route
        path="/password-reset/:token"
        element={
          <AuthGuard>
            <PasswordResetPage />
          </AuthGuard>
        }
      />
      <Route
        path="/passwordReset/:token"
        element={
          <AuthGuard>
            <PasswordResetPage />
          </AuthGuard>
        }
      />

      <Route
        element={
          <AuthGuard>
            <AppLayout />
          </AuthGuard>
        }
      >
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/onboarding" element={<OnboardingWizardPage />} />
        <Route path="/devices" element={<DevicesPage />} />
        <Route path="/groups" element={<GroupsPage />} />
        <Route path="/applications" element={<ApplicationsPage />} />
        <Route path="/applications/admin" element={<AdminApplicationsPage />} />
        <Route path="/application/:id/versions" element={<ApplicationVersionsPage />} />
        <Route path="/profiles" element={<ProfilesPage />} />
        <Route path="/profiles/:profileId/edit" element={<ProfileEditRedirect />} />
        <Route path="/profiles/:profileId/versions/:versionId/edit" element={<ProfileEditRedirect />} />
        <Route path="/enrollment-routes" element={<EnrollmentRouteListPage />} />
        <Route path="/enrollment-routes/*" element={<EnrollmentRouteListPage />} />
        <Route path="/configurations" element={<Navigate to="/profiles" replace />} />
        <Route path="/configurations/:id/edit" element={<ConfigurationEditorPage />} />
        <Route path="/qr/:qrCodeKey/:deviceId" element={<EnrollmentQrPage />} />
        <Route path="/qr/:qrCodeKey" element={<EnrollmentQrPage />} />
        <Route path="/users" element={<UsersPage />} />
        <Route path="/roles" element={<RolesPage />} />
        <Route path="/settings" element={<SettingsPage />} />
        <Route path="/files" element={<FilesPage />} />
        <Route path="/icons" element={<IconsPage />} />
        <Route path="/hints" element={<HintsPage />} />
        <Route path="/updates" element={<UpdatesPage />} />
        <Route path="/push" element={<PushPage />} />
        <Route path="/control-panel" element={<ControlPanelPage />} />
        <Route path="/plugin-settings" element={<PluginSettingsPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/twofactor" element={<TwofactorPage />} />
        <Route path="/maps" element={<MapsInfoPage />} />
      </Route>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}
