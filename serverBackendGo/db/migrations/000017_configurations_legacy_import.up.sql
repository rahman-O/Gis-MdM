-- Merge legacy Java configuration columns into settingsjson when present (Java dump restore).
-- No-op on greenfield Go-only databases.

DO $$
DECLARE
    col_exists BOOLEAN;
    patch JSONB;
BEGIN
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'configurations' AND column_name = 'gps'
    ) INTO col_exists;
    IF NOT col_exists THEN
        RETURN;
    END IF;

    UPDATE configurations c SET settingsjson = COALESCE(c.settingsjson, '{}'::jsonb) || sub.patch
    FROM (
        SELECT id,
            jsonb_strip_nulls(jsonb_build_object(
                'gps', CASE WHEN gps IS NOT NULL THEN to_jsonb(gps) END,
                'bluetooth', CASE WHEN bluetooth IS NOT NULL THEN to_jsonb(bluetooth) END,
                'wifi', CASE WHEN wifi IS NOT NULL THEN to_jsonb(wifi) END,
                'mobileData', CASE WHEN mobiledata IS NOT NULL THEN to_jsonb(mobiledata) END,
                'kioskMode', CASE WHEN kioskmode IS NOT NULL THEN to_jsonb(kioskmode) END,
                'blockStatusBar', CASE WHEN blockstatusbar IS NOT NULL THEN to_jsonb(blockstatusbar) END,
                'systemUpdateType', CASE WHEN systemupdatetype IS NOT NULL THEN to_jsonb(systemupdatetype) END,
                'autoUpdate', CASE WHEN autoupdate IS NOT NULL THEN to_jsonb(autoupdate) END,
                'usbStorage', CASE WHEN usbstorage IS NOT NULL THEN to_jsonb(usbstorage) END,
                'requestUpdates', CASE WHEN requestupdates IS NOT NULL THEN to_jsonb(requestupdates) END,
                'pushOptions', CASE WHEN pushoptions IS NOT NULL THEN to_jsonb(pushoptions) END,
                'manageTimeout', CASE WHEN managetimeout IS NOT NULL THEN to_jsonb(managetimeout) END,
                'orientation', CASE WHEN orientation IS NOT NULL THEN to_jsonb(orientation) END,
                'timeZone', CASE WHEN timezone IS NOT NULL THEN to_jsonb(timezone) END,
                'disableScreenshots', CASE WHEN disablescreenshots IS NOT NULL THEN to_jsonb(disablescreenshots) END,
                'restrictions', CASE WHEN restrictions IS NOT NULL THEN to_jsonb(restrictions) END,
                'keepaliveTime', CASE WHEN keepalivetime IS NOT NULL THEN to_jsonb(keepalivetime) END,
                'manageVolume', CASE WHEN managevolume IS NOT NULL THEN to_jsonb(managevolume) END,
                'showWifi', CASE WHEN showwifi IS NOT NULL THEN to_jsonb(showwifi) END,
                'mobileEnrollment', CASE WHEN mobileenrollment IS NOT NULL THEN to_jsonb(mobileenrollment) END,
                'disableLocation', CASE WHEN disablelocation IS NOT NULL THEN to_jsonb(disablelocation) END,
                'appPermissions', CASE WHEN apppermissions IS NOT NULL THEN to_jsonb(apppermissions) END,
                'encryptDevice', CASE WHEN encryptdevice IS NOT NULL THEN to_jsonb(encryptdevice) END,
                'downloadUpdates', CASE WHEN downloadupdates IS NOT NULL THEN to_jsonb(downloadupdates) END
            )) AS patch
        FROM configurations
    ) sub
    WHERE c.id = sub.id AND sub.patch <> '{}'::jsonb;
END $$;
