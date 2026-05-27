import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/local_db.dart';
import '../../core/utils/logger.dart';
import 'collectors/battery_collector.dart';
import 'collectors/location_collector.dart';
import 'collectors/network_collector.dart';
import 'collectors/storage_collector.dart';
import 'collectors/system_collector.dart';
import 'telemetry_data.dart';

/// Service that orchestrates telemetry collection and transmission.
///
/// Runs all collectors in parallel, combines results into a [TelemetryData]
/// snapshot, stores it locally, and sends it to the server.
class TelemetryService {
  final BatteryCollector _batteryCollector = BatteryCollector();
  final NetworkCollector _networkCollector = NetworkCollector();
  final StorageCollector _storageCollector = StorageCollector();
  final LocationCollector _locationCollector = LocationCollector();
  final SystemCollector _systemCollector = SystemCollector();

  /// Collect telemetry from all device sensors/sources.
  ///
  /// Runs collectors in parallel for efficiency. Stores the result
  /// in the local queue for offline resilience.
  Future<TelemetryData> collect(String deviceId) async {
    Logger.debug('Starting telemetry collection', 'Telemetry');

    // Run all collectors in parallel
    final results = await Future.wait([
      _batteryCollector.collect(),
      _networkCollector.collect(),
      _storageCollector.collect(),
      _locationCollector.collect(),
      _systemCollector.collect(),
    ]);

    final battery = results[0] as BatteryInfo;
    final network = results[1] as NetworkInfo;
    final storage = results[2] as StorageInfo;
    final location = results[3] as LocationInfo?;
    final system = results[4] as SystemInfo;

    final data = TelemetryData(
      deviceId: deviceId,
      timestamp: DateTime.now().millisecondsSinceEpoch,
      battery: battery,
      network: network,
      storage: storage,
      memory: MemoryInfo(totalBytes: 0, freeBytes: 0), // Requires platform channel
      location: location,
      screen: ScreenInfo(brightness: 0.5, isOn: true), // Requires platform channel
      system: system,
    );

    // Store in local queue for offline resilience
    await LocalDb.addTelemetry(data.toJson());

    Logger.info(
      'Telemetry collected: battery=${battery.level}%, '
      'network=${network.type}, location=${location != null}',
      'Telemetry',
    );

    return data;
  }

  /// Send telemetry data to the server.
  Future<void> sendToServer(
    ApiClient api,
    String deviceId,
    TelemetryData data,
  ) async {
    try {
      await api.put(
        Endpoints.deviceInfo(deviceId),
        data: data.toJson(),
      );
      Logger.debug('Telemetry sent to server', 'Telemetry');
    } catch (e, stack) {
      Logger.error('Failed to send telemetry', e, stack, 'Telemetry');
      rethrow;
    }
  }

  /// Flush any pending telemetry from the local queue.
  Future<void> flushPending(ApiClient api, String deviceId) async {
    final pending = LocalDb.getPendingTelemetry();
    if (pending.isEmpty) return;

    Logger.info('Flushing ${pending.length} pending telemetry items', 'Telemetry');

    int sent = 0;
    for (final item in pending) {
      try {
        await api.post(Endpoints.deviceInfo(deviceId), data: item);
        sent++;
      } catch (e) {
        Logger.warn('Failed to flush telemetry item: $e', 'Telemetry');
        break;
      }
    }

    if (sent > 0) {
      await LocalDb.clearTelemetry(sent);
    }
  }
}
