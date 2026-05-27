/// Model representing the server sync response.
///
/// Matches the Go backend SyncResponse struct — contains the full
/// device configuration including apps, files, policies, and restrictions.
class SyncResponse {
  final String deviceId;
  final int configurationId;
  final List<SyncApplication> applications;
  final List<SyncFile> files;
  final bool? permissive;
  final bool kioskMode;
  final String? kioskApp;
  final String? restrictions;
  final String? password;
  final String? wifiSsid;
  final String? wifiPassword;
  final int? wifiSecurityType;
  final bool? gps;
  final bool? bluetooth;
  final bool? wifi;
  final bool? mobileData;
  final bool? usbStorage;
  final int? screenBrightness;
  final bool? manageVolume;
  final int? volumeLevel;
  final String? passwordMode;
  final String? timeZone;
  final bool? lockStatusBar;
  final bool? systemUpdateType;
  final String? systemUpdateFrom;
  final String? systemUpdateTo;
  final bool? factoryReset;
  final String? appPermissions;
  final bool? autoBrightness;
  final bool? managedUpdate;
  final int? scheduleType;

  SyncResponse({
    required this.deviceId,
    required this.configurationId,
    this.applications = const [],
    this.files = const [],
    this.permissive,
    this.kioskMode = false,
    this.kioskApp,
    this.restrictions,
    this.password,
    this.wifiSsid,
    this.wifiPassword,
    this.wifiSecurityType,
    this.gps,
    this.bluetooth,
    this.wifi,
    this.mobileData,
    this.usbStorage,
    this.screenBrightness,
    this.manageVolume,
    this.volumeLevel,
    this.passwordMode,
    this.timeZone,
    this.lockStatusBar,
    this.systemUpdateType,
    this.systemUpdateFrom,
    this.systemUpdateTo,
    this.factoryReset,
    this.appPermissions,
    this.autoBrightness,
    this.managedUpdate,
    this.scheduleType,
  });

  factory SyncResponse.fromJson(Map<String, dynamic> json) {
    return SyncResponse(
      deviceId: json['deviceId'] as String? ?? '',
      configurationId: json['configurationId'] as int? ?? 0,
      applications: (json['applications'] as List<dynamic>?)
              ?.map((e) => SyncApplication.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      files: (json['files'] as List<dynamic>?)
              ?.map((e) => SyncFile.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      permissive: json['permissive'] as bool?,
      kioskMode: json['kioskMode'] as bool? ?? false,
      kioskApp: json['kioskApp'] as String?,
      restrictions: json['restrictions'] as String?,
      password: json['password'] as String?,
      wifiSsid: json['wifiSsid'] as String?,
      wifiPassword: json['wifiPassword'] as String?,
      wifiSecurityType: json['wifiSecurityType'] as int?,
      gps: json['gps'] as bool?,
      bluetooth: json['bluetooth'] as bool?,
      wifi: json['wifi'] as bool?,
      mobileData: json['mobileData'] as bool?,
      usbStorage: json['usbStorage'] as bool?,
      screenBrightness: json['screenBrightness'] as int?,
      manageVolume: json['manageVolume'] as bool?,
      volumeLevel: json['volumeLevel'] as int?,
      passwordMode: json['passwordMode'] as String?,
      timeZone: json['timeZone'] as String?,
      lockStatusBar: json['lockStatusBar'] as bool?,
      systemUpdateType: json['systemUpdateType'] as bool?,
      systemUpdateFrom: json['systemUpdateFrom'] as String?,
      systemUpdateTo: json['systemUpdateTo'] as String?,
      factoryReset: json['factoryReset'] as bool?,
      appPermissions: json['appPermissions'] as String?,
      autoBrightness: json['autoBrightness'] as bool?,
      managedUpdate: json['managedUpdate'] as bool?,
      scheduleType: json['scheduleType'] as int?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'deviceId': deviceId,
      'configurationId': configurationId,
      'applications': applications.map((a) => a.toJson()).toList(),
      'files': files.map((f) => f.toJson()).toList(),
      if (permissive != null) 'permissive': permissive,
      'kioskMode': kioskMode,
      if (kioskApp != null) 'kioskApp': kioskApp,
      if (restrictions != null) 'restrictions': restrictions,
      if (password != null) 'password': password,
      if (wifiSsid != null) 'wifiSsid': wifiSsid,
      if (wifiPassword != null) 'wifiPassword': wifiPassword,
      if (wifiSecurityType != null) 'wifiSecurityType': wifiSecurityType,
      if (gps != null) 'gps': gps,
      if (bluetooth != null) 'bluetooth': bluetooth,
      if (wifi != null) 'wifi': wifi,
      if (mobileData != null) 'mobileData': mobileData,
      if (usbStorage != null) 'usbStorage': usbStorage,
      if (screenBrightness != null) 'screenBrightness': screenBrightness,
      if (manageVolume != null) 'manageVolume': manageVolume,
      if (volumeLevel != null) 'volumeLevel': volumeLevel,
      if (passwordMode != null) 'passwordMode': passwordMode,
      if (timeZone != null) 'timeZone': timeZone,
      if (lockStatusBar != null) 'lockStatusBar': lockStatusBar,
      if (systemUpdateType != null) 'systemUpdateType': systemUpdateType,
      if (systemUpdateFrom != null) 'systemUpdateFrom': systemUpdateFrom,
      if (systemUpdateTo != null) 'systemUpdateTo': systemUpdateTo,
      if (factoryReset != null) 'factoryReset': factoryReset,
      if (appPermissions != null) 'appPermissions': appPermissions,
      if (autoBrightness != null) 'autoBrightness': autoBrightness,
      if (managedUpdate != null) 'managedUpdate': managedUpdate,
      if (scheduleType != null) 'scheduleType': scheduleType,
    };
  }
}

/// Represents an application in the sync configuration.
class SyncApplication {
  final int id;
  final String name;
  final String pkg;
  final String version;
  final String url;
  final String type;
  final bool? showIcon;
  final bool? runAfterInstall;
  final bool? runAtBoot;
  final bool? system;
  final int? screenOrder;

  SyncApplication({
    required this.id,
    required this.name,
    required this.pkg,
    required this.version,
    required this.url,
    required this.type,
    this.showIcon,
    this.runAfterInstall,
    this.runAtBoot,
    this.system,
    this.screenOrder,
  });

  factory SyncApplication.fromJson(Map<String, dynamic> json) {
    return SyncApplication(
      id: json['id'] as int? ?? 0,
      name: json['name'] as String? ?? '',
      pkg: json['pkg'] as String? ?? '',
      version: json['version'] as String? ?? '',
      url: json['url'] as String? ?? '',
      type: json['type'] as String? ?? '',
      showIcon: json['showIcon'] as bool?,
      runAfterInstall: json['runAfterInstall'] as bool?,
      runAtBoot: json['runAtBoot'] as bool?,
      system: json['system'] as bool?,
      screenOrder: json['screenOrder'] as int?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'pkg': pkg,
      'version': version,
      'url': url,
      'type': type,
      if (showIcon != null) 'showIcon': showIcon,
      if (runAfterInstall != null) 'runAfterInstall': runAfterInstall,
      if (runAtBoot != null) 'runAtBoot': runAtBoot,
      if (system != null) 'system': system,
      if (screenOrder != null) 'screenOrder': screenOrder,
    };
  }
}

/// Represents a file to be pushed to the device.
class SyncFile {
  final String devicePath;
  final String url;
  final bool remove;

  SyncFile({
    required this.devicePath,
    required this.url,
    this.remove = false,
  });

  factory SyncFile.fromJson(Map<String, dynamic> json) {
    return SyncFile(
      devicePath: json['devicePath'] as String? ?? '',
      url: json['url'] as String? ?? '',
      remove: json['remove'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'devicePath': devicePath,
      'url': url,
      'remove': remove,
    };
  }
}
