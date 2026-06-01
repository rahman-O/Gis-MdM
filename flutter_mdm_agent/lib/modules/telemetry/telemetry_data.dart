/// Complete telemetry snapshot collected from the device.
class TelemetryData {
  final String deviceId;
  final int timestamp;
  final BatteryInfo battery;
  final NetworkInfo network;
  final StorageInfo storage;
  final MemoryInfo memory;
  final LocationInfo? location;
  final ScreenInfo screen;
  final SystemInfo system;

  TelemetryData({
    required this.deviceId,
    required this.timestamp,
    required this.battery,
    required this.network,
    required this.storage,
    required this.memory,
    this.location,
    required this.screen,
    required this.system,
  });

  Map<String, dynamic> toJson() {
    return {
      'deviceId': deviceId,
      'timestamp': timestamp,
      'battery': battery.toJson(),
      'network': network.toJson(),
      'storage': storage.toJson(),
      'memory': memory.toJson(),
      if (location != null) 'location': location!.toJson(),
      'screen': screen.toJson(),
      'system': system.toJson(),
    };
  }

  factory TelemetryData.fromJson(Map<String, dynamic> json) {
    return TelemetryData(
      deviceId: json['deviceId'] as String,
      timestamp: json['timestamp'] as int,
      battery: BatteryInfo.fromJson(json['battery'] as Map<String, dynamic>),
      network: NetworkInfo.fromJson(json['network'] as Map<String, dynamic>),
      storage: StorageInfo.fromJson(json['storage'] as Map<String, dynamic>),
      memory: MemoryInfo.fromJson(json['memory'] as Map<String, dynamic>),
      location: json['location'] != null
          ? LocationInfo.fromJson(json['location'] as Map<String, dynamic>)
          : null,
      screen: ScreenInfo.fromJson(json['screen'] as Map<String, dynamic>),
      system: SystemInfo.fromJson(json['system'] as Map<String, dynamic>),
    );
  }
}

/// Battery status information.
class BatteryInfo {
  final int level;
  final String chargingState;

  BatteryInfo({required this.level, required this.chargingState});

  Map<String, dynamic> toJson() => {
        'level': level,
        'chargingState': chargingState,
      };

  factory BatteryInfo.fromJson(Map<String, dynamic> json) {
    return BatteryInfo(
      level: json['level'] as int? ?? 0,
      chargingState: json['chargingState'] as String? ?? 'unknown',
    );
  }
}

/// Network connectivity information.
class NetworkInfo {
  final String type;
  final bool connected;
  final String? wifiSsid;
  final String? ipAddress;

  NetworkInfo({
    required this.type,
    required this.connected,
    this.wifiSsid,
    this.ipAddress,
  });

  Map<String, dynamic> toJson() => {
        'type': type,
        'connected': connected,
        if (wifiSsid != null) 'wifiSsid': wifiSsid,
        if (ipAddress != null) 'ipAddress': ipAddress,
      };

  factory NetworkInfo.fromJson(Map<String, dynamic> json) {
    return NetworkInfo(
      type: json['type'] as String? ?? 'unknown',
      connected: json['connected'] as bool? ?? false,
      wifiSsid: json['wifiSsid'] as String?,
      ipAddress: json['ipAddress'] as String?,
    );
  }
}

/// Disk storage information.
class StorageInfo {
  final int totalBytes;
  final int freeBytes;

  StorageInfo({required this.totalBytes, required this.freeBytes});

  int get usedBytes => totalBytes - freeBytes;
  double get usedPercent =>
      totalBytes > 0 ? (usedBytes / totalBytes) * 100 : 0;

  Map<String, dynamic> toJson() => {
        'totalBytes': totalBytes,
        'freeBytes': freeBytes,
      };

  factory StorageInfo.fromJson(Map<String, dynamic> json) {
    return StorageInfo(
      totalBytes: json['totalBytes'] as int? ?? 0,
      freeBytes: json['freeBytes'] as int? ?? 0,
    );
  }
}

/// RAM memory information.
class MemoryInfo {
  final int totalBytes;
  final int freeBytes;

  MemoryInfo({required this.totalBytes, required this.freeBytes});

  Map<String, dynamic> toJson() => {
        'totalBytes': totalBytes,
        'freeBytes': freeBytes,
      };

  factory MemoryInfo.fromJson(Map<String, dynamic> json) {
    return MemoryInfo(
      totalBytes: json['totalBytes'] as int? ?? 0,
      freeBytes: json['freeBytes'] as int? ?? 0,
    );
  }
}

/// GPS location information.
class LocationInfo {
  final double latitude;
  final double longitude;
  final double accuracy;
  final int timestamp;

  LocationInfo({
    required this.latitude,
    required this.longitude,
    required this.accuracy,
    required this.timestamp,
  });

  Map<String, dynamic> toJson() => {
        'latitude': latitude,
        'longitude': longitude,
        'accuracy': accuracy,
        'timestamp': timestamp,
      };

  factory LocationInfo.fromJson(Map<String, dynamic> json) {
    return LocationInfo(
      latitude: (json['latitude'] as num?)?.toDouble() ?? 0.0,
      longitude: (json['longitude'] as num?)?.toDouble() ?? 0.0,
      accuracy: (json['accuracy'] as num?)?.toDouble() ?? 0.0,
      timestamp: json['timestamp'] as int? ?? 0,
    );
  }
}

/// Screen display information.
class ScreenInfo {
  final double brightness;
  final bool isOn;

  ScreenInfo({required this.brightness, required this.isOn});

  Map<String, dynamic> toJson() => {
        'brightness': brightness,
        'isOn': isOn,
      };

  factory ScreenInfo.fromJson(Map<String, dynamic> json) {
    return ScreenInfo(
      brightness: (json['brightness'] as num?)?.toDouble() ?? 0.0,
      isOn: json['isOn'] as bool? ?? true,
    );
  }
}

/// System-level device information.
class SystemInfo {
  final String model;
  final String manufacturer;
  final String androidVersion;
  final int sdkInt;
  final String serial;
  final int uptimeMillis;

  SystemInfo({
    required this.model,
    required this.manufacturer,
    required this.androidVersion,
    required this.sdkInt,
    required this.serial,
    required this.uptimeMillis,
  });

  Map<String, dynamic> toJson() => {
        'model': model,
        'manufacturer': manufacturer,
        'androidVersion': androidVersion,
        'sdkInt': sdkInt,
        'serial': serial,
        'uptimeMillis': uptimeMillis,
      };

  factory SystemInfo.fromJson(Map<String, dynamic> json) {
    return SystemInfo(
      model: json['model'] as String? ?? '',
      manufacturer: json['manufacturer'] as String? ?? '',
      androidVersion: json['androidVersion'] as String? ?? '',
      sdkInt: json['sdkInt'] as int? ?? 0,
      serial: json['serial'] as String? ?? '',
      uptimeMillis: json['uptimeMillis'] as int? ?? 0,
    );
  }
}
