-- 031: Move provisioning settings from profile to enrollment routes

ALTER TABLE enrollment_routes
  ADD COLUMN IF NOT EXISTS wifi_ssid VARCHAR(32),
  ADD COLUMN IF NOT EXISTS wifi_password VARCHAR(63),
  ADD COLUMN IF NOT EXISTS wifi_security_type VARCHAR(10),
  ADD COLUMN IF NOT EXISTS qr_parameters TEXT,
  ADD COLUMN IF NOT EXISTS admin_extras TEXT,
  ADD COLUMN IF NOT EXISTS mobile_enrollment BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS encrypt_device BOOLEAN NOT NULL DEFAULT false;

-- Backfill from linked profile_versions.settingsjson (idempotent: only where columns are still NULL/default)
UPDATE enrollment_routes er
SET
  wifi_ssid = COALESCE(NULLIF(TRIM(pv.settingsjson->>'wifiSSID'), ''), er.wifi_ssid),
  wifi_password = COALESCE(NULLIF(TRIM(pv.settingsjson->>'wifiPassword'), ''), er.wifi_password),
  wifi_security_type = COALESCE(NULLIF(TRIM(pv.settingsjson->>'wifiSecurityType'), ''), er.wifi_security_type),
  qr_parameters = COALESCE(NULLIF(TRIM(pv.settingsjson->>'qrParameters'), ''), er.qr_parameters),
  admin_extras = COALESCE(NULLIF(TRIM(pv.settingsjson->>'adminExtras'), ''), er.admin_extras),
  mobile_enrollment = COALESCE((pv.settingsjson->>'mobileEnrollment')::boolean, er.mobile_enrollment),
  encrypt_device = COALESCE((pv.settingsjson->>'encryptDevice')::boolean, er.encrypt_device)
FROM profile_versions pv
WHERE er.profile_version_id = pv.id
  AND er.profile_version_id IS NOT NULL
  AND er.wifi_ssid IS NULL
  AND er.mobile_enrollment = false
  AND er.encrypt_device = false;
