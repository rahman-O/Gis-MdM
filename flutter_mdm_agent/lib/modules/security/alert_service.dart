import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/local_db.dart';
import '../../core/utils/logger.dart';

/// Sends security alerts to the server.
class AlertService {
  final ApiClient _api;

  AlertService(this._api);

  /// Send an alert event to the server.
  Future<void> sendAlert(
      String deviceId, String alertType, Map<String, dynamic> details) async {
    final payload = {
      'deviceId': deviceId,
      'alertType': alertType,
      'details': details,
      'timestamp': DateTime.now().millisecondsSinceEpoch,
    };
    try {
      await _api.post(Endpoints.deviceLog(deviceId), data: [payload]);
      Logger.info('Alert sent: $alertType', 'AlertService');
    } catch (e) {
      // Store locally if offline
      await LocalDb.addEvent({'type': 'alert', ...payload});
      Logger.warn('Alert queued (offline): $alertType', 'AlertService');
    }
  }
}
