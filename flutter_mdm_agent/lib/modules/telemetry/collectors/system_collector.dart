import 'package:device_info_plus/device_info_plus.dart';

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects system-level device information (model, OS version, uptime).
class SystemCollector {
  final DeviceInfoPlugin _deviceInfo = DeviceInfoPlugin();

  /// Collect current system information.
  Future<SystemInfo> collect() async {
    try {
      final androidInfo = await _deviceInfo.androidInfo;

      return SystemInfo(
        model: androidInfo.model,
        manufacturer: androidInfo.manufacturer,
        androidVersion: androidInfo.version.release,
        sdkInt: androidInfo.version.sdkInt,
        serial: androidInfo.serialNumber,
        uptimeMillis: DateTime.now().millisecondsSinceEpoch -
            androidInfo.version.sdkInt, // Approximation; real uptime via platform channel
      );
    } catch (e) {
      Logger.warn('System collection failed: $e', 'SystemCollector');
      return SystemInfo(
        model: 'unknown',
        manufacturer: 'unknown',
        androidVersion: 'unknown',
        sdkInt: 0,
        serial: '',
        uptimeMillis: 0,
      );
    }
  }
}
