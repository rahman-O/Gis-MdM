-- 031 rollback: remove provisioning columns from enrollment_routes

ALTER TABLE enrollment_routes
  DROP COLUMN IF EXISTS wifi_ssid,
  DROP COLUMN IF EXISTS wifi_password,
  DROP COLUMN IF EXISTS wifi_security_type,
  DROP COLUMN IF EXISTS qr_parameters,
  DROP COLUMN IF EXISTS admin_extras,
  DROP COLUMN IF EXISTS mobile_enrollment,
  DROP COLUMN IF EXISTS encrypt_device;
