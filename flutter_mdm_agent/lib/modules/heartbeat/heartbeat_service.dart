import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/utils/logger.dart';
import '../telemetry/collectors/battery_collector.dart';
import '../telemetry/collectors/network_collector.dart';
import '../telemetry/collectors/system_collector.dart';
import '../telemetry/collectors/location_collector.dart';
import '../telemetry/telemetry_data.dart';

/// Lightweight heartbeat service that periodically pings the server.
class HeartbeatService {
  final ApiClient _api;
  final String _deviceId;
  final BatteryCollector _batteryCollector = BatteryCollector();
  final NetworkCollector _networkCollector = NetworkCollector();
  final SystemCollector _systemCollector = SystemCollector();
  final LocationCollector _locationCollector = LocationCollector();

  HeartbeatService(this._api, this._deviceId);

  /// Send a heartbeat with full device info to the server.
  Future<void> sendHeartbeat() async {
    try {
      Logger.debug('Sending heartbeat', 'Heartbeat');

      final results = await Future.wait([
        _batteryCollector.collect(),
        _networkCollector.collect(),
        _systemCollector.collect(),
        _locationCollector.collect(),
      ]);

      final battery = results[0] as BatteryInfo;
      final network = results[1] as NetworkInfo;
      final system = results[2] as SystemInfo;
      final location = results[3];

      final body = <String, dynamic>{
        'deviceId': _deviceId,
        'batteryLevel': battery.level,
        // Extended fields — saved in device.info JSON
        'imei': _deviceId, // Device number is IMEI-based
        'model': _ensureNonNull(system.model, 'unknown'),
        'manufacturer': _ensureNonNull(system.manufacturer, 'unknown'),
        'androidVersion': _ensureNonNull(system.androidVersion, 'unknown'),
        'serial': _ensureNonNull(system.serial, 'unavailable'),
        'launcherVersion': '1.0.0',
        // Phone number — requires READ_PHONE_STATE + carrier support;
        // not reliably available on modern Android, send null (backend handles gracefully)
        'phone': null,
      };

      // Add network info fields
      if (network.type != 'unknown') {
        body['networkType'] = network.type;
      }
      if (network.wifiSsid != null) {
        body['wifiSsid'] = network.wifiSsid;
      }
      if (network.ipAddress != null) {
        body['ipAddress'] = network.ipAddress;
      }

      // Add location if available
      if (location != null) {
        final loc = location as LocationInfo;
        body['location'] = {
          'lat': loc.latitude,
          'lon': loc.longitude,
          'accuracy': loc.accuracy,
          'ts': loc.timestamp,
        };
      }

      Logger.debug(
        'Heartbeat payload: model=${body['model']}, manufacturer=${body['manufacturer']}, '
        'androidVersion=${body['androidVersion']}, serial=${body['serial']}, '
        'networkType=${network.type}, wifiSsid=${network.wifiSsid}, ipAddress=${network.ipAddress}',
        'Heartbeat',
      );

      await _api.post(Endpoints.syncInfo, data: body);
      Logger.debug('Heartbeat sent successfully', 'Heartbeat');
    } catch (e, stack) {
      Logger.error('Heartbeat failed', e, stack, 'Heartbeat');
      rethrow;
    }
  }

  /// Ensures a value is never null or empty before sending to the server.
  /// Returns [fallback] if [value] is null or empty.
  String _ensureNonNull(String? value, String fallback) {
    if (value == null || value.trim().isEmpty) {
      return fallback;
    }
    return value;
  }
}
