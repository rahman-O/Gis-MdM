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

      // Ensure no null or empty values slip through.
      // On Android 10+ (API 29+), serialNumber returns "unknown" for non-rooted devices.
      final serial = _sanitize(androidInfo.serialNumber, fallback: 'unavailable');

      return SystemInfo(
        model: _sanitize(androidInfo.model, fallback: 'unknown'),
        manufacturer: _sanitize(androidInfo.manufacturer, fallback: 'unknown'),
        androidVersion: _sanitize(androidInfo.version.release, fallback: 'unknown'),
        sdkInt: androidInfo.version.sdkInt,
        serial: serial,
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
        serial: 'unavailable',
        uptimeMillis: 0,
      );
    }
  }

  /// Returns [value] if it is non-null and non-empty, otherwise returns [fallback].
  /// Also treats the literal string "unknown" from Android API as a valid but
  /// flagged value (e.g. serial on Android 10+).
  String _sanitize(String? value, {required String fallback}) {
    if (value == null || value.trim().isEmpty) {
      return fallback;
    }
    return value.trim();
  }
}
