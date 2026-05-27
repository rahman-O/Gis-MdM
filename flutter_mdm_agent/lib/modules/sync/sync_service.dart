import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/local_db.dart';
import '../../core/utils/logger.dart';
import 'sync_response.dart';

/// Service responsible for fetching device configuration from the server
/// and sending device info updates.
class SyncService {
  final ApiClient _api;

  SyncService(this._api);

  /// Fetch the full device configuration from the server.
  ///
  /// Parses the response into a [SyncResponse], caches it locally,
  /// and returns it. Returns `null` on failure.
  Future<SyncResponse?> fetchConfiguration(String deviceId) async {
    try {
      Logger.info('Fetching configuration for $deviceId', 'Sync');

      final response = await _api.get(
        Endpoints.syncConfiguration(deviceId),
      );

      if (response.statusCode == 200 && response.data != null) {
        final data = response.data as Map<String, dynamic>;
        final syncResponse = SyncResponse.fromJson(data);

        // Cache locally for offline access
        await LocalDb.saveConfig(syncResponse.toJson());

        Logger.info(
          'Configuration fetched: ${syncResponse.applications.length} apps, '
          '${syncResponse.files.length} files',
          'Sync',
        );
        return syncResponse;
      }

      Logger.warn(
        'Fetch configuration failed: status=${response.statusCode}',
        'Sync',
      );
      return null;
    } catch (e, stack) {
      Logger.error('Fetch configuration error', e, stack, 'Sync');
      return null;
    }
  }

  /// Send device info to the server.
  ///
  /// [info] should contain device metadata (model, OS version, battery, etc.)
  Future<void> sendDeviceInfo(
    String deviceId,
    Map<String, dynamic> info,
  ) async {
    try {
      Logger.debug('Sending device info for $deviceId', 'Sync');

      await _api.post(
        Endpoints.syncInfo,
        data: {
          'deviceId': deviceId,
          ...info,
        },
      );

      Logger.debug('Device info sent successfully', 'Sync');
    } catch (e, stack) {
      Logger.error('Send device info error', e, stack, 'Sync');
      rethrow;
    }
  }

  /// Load cached configuration from local storage.
  ///
  /// Useful when the device is offline.
  SyncResponse? getCachedConfiguration() {
    final cached = LocalDb.getConfig();
    if (cached != null) {
      return SyncResponse.fromJson(cached);
    }
    return null;
  }
}
