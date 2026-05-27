import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/utils/logger.dart';
import '../telemetry/collectors/battery_collector.dart';
import '../telemetry/collectors/network_collector.dart';

/// Lightweight heartbeat service that periodically pings the server.
///
/// Sends minimal device status (battery, network, timestamp) to let
/// the server know the device is alive and reachable.
class HeartbeatService {
  final ApiClient _api;
  final String _deviceId;
  final BatteryCollector _batteryCollector = BatteryCollector();
  final NetworkCollector _networkCollector = NetworkCollector();

  HeartbeatService(this._api, this._deviceId);

  /// Send a heartbeat to the server with minimal device status.
  Future<void> sendHeartbeat() async {
    try {
      Logger.debug('Sending heartbeat', 'Heartbeat');

      // Collect minimal info in parallel
      final results = await Future.wait([
        _batteryCollector.collect(),
        _networkCollector.collect(),
      ]);

      final battery = results[0] as dynamic;
      final network = results[1] as dynamic;

      await _api.post(
        Endpoints.syncInfo,
        data: {
          'deviceId': _deviceId,
          'batteryLevel': battery.level,
          'networkType': network.type,
          'timestamp': DateTime.now().millisecondsSinceEpoch,
        },
      );

      Logger.debug('Heartbeat sent successfully', 'Heartbeat');
    } catch (e, stack) {
      Logger.error('Heartbeat failed', e, stack, 'Heartbeat');
      rethrow;
    }
  }
}
