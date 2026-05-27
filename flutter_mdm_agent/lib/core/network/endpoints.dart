/// API endpoint definitions matching serverBackendGo routes.
class Endpoints {
  Endpoints._();

  // Sync (public — no auth)
  static String syncConfiguration(String deviceId) => '/rest/public/sync/configuration/$deviceId';
  static const String syncInfo = '/rest/public/sync/info';
  static String syncAppSettings(String deviceId) => '/rest/public/sync/applicationSettings/$deviceId';

  // Notifications (public — signature auth)
  static String notificationPolling(String deviceId) => '/rest/notification/polling/$deviceId';

  // Device Info Plugin (public)
  static String deviceInfo(String deviceNumber) => '/rest/plugins/deviceinfo/deviceinfo/public/$deviceNumber';

  // Device Log Plugin (public)
  static String deviceLog(String deviceNumber) => '/rest/plugins/devicelog/log/list/$deviceNumber';

  // Device Location (public — agent sends batch locations)
  static String deviceLocations(String deviceId) => '/rest/public/device-locations/$deviceId';

  // Public
  static const String publicName = '/rest/public/name';
}
