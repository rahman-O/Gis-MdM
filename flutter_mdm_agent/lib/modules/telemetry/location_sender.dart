import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/utils/logger.dart';
import 'telemetry_data.dart';

/// Sends location data to the server's device-locations endpoint.
///
/// Posts a batch of location points to:
/// POST /rest/public/device-locations/{deviceId}
class LocationSender {
  /// Send a single location point as a batch of one.
  Future<void> send(ApiClient api, String deviceId, LocationInfo location) async {
    await sendBatch(api, deviceId, [location]);
  }

  /// Send a batch of location points to the server.
  Future<void> sendBatch(ApiClient api, String deviceId, List<LocationInfo> points) async {
    if (points.isEmpty) return;

    final body = points
        .map((p) => {
              'latitude': p.latitude,
              'longitude': p.longitude,
              'accuracy': p.accuracy,
              'speed': 0.0,
              'timestamp': p.timestamp,
            })
        .toList();

    try {
      await api.post(
        Endpoints.deviceLocations(deviceId),
        data: body,
      );
      Logger.debug('Sent ${points.length} location points to server', 'LocationSender');
    } catch (e, stack) {
      Logger.error('Failed to send locations', e, stack, 'LocationSender');
      rethrow;
    }
  }
}
