import 'package:device_info_plus/device_info_plus.dart';

import '../../../core/utils/logger.dart';

/// Collects information about installed applications.
class AppsCollector {
  final DeviceInfoPlugin _deviceInfo = DeviceInfoPlugin();

  /// Collect installed apps count and system info.
  ///
  /// Note: Full app list requires platform channel; this provides
  /// basic device info as a proxy for app-related telemetry.
  Future<AppsInfo> collect() async {
    try {
      final androidInfo = await _deviceInfo.androidInfo;

      return AppsInfo(
        systemFeatures: androidInfo.systemFeatures.length,
        supportedAbis: androidInfo.supportedAbis,
      );
    } catch (e) {
      Logger.warn('Apps collection failed: $e', 'AppsCollector');
      return AppsInfo(systemFeatures: 0, supportedAbis: []);
    }
  }
}

/// Basic apps/features information from the device.
class AppsInfo {
  final int systemFeatures;
  final List<String> supportedAbis;

  AppsInfo({required this.systemFeatures, required this.supportedAbis});

  Map<String, dynamic> toJson() => {
        'systemFeatures': systemFeatures,
        'supportedAbis': supportedAbis,
      };
}
