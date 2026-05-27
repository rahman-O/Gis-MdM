import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/local_db.dart';
import '../../core/utils/logger.dart';

/// Batches events and sends them to the server periodically.
class EventReporter {
  final ApiClient _api;

  EventReporter(this._api);

  /// Record an event locally.
  Future<void> record(String eventType, Map<String, dynamic> data) async {
    final event = {
      'type': eventType,
      'data': data,
      'timestamp': DateTime.now().millisecondsSinceEpoch,
    };
    await LocalDb.addEvent(event);
  }

  /// Flush pending events to server.
  Future<void> flush(String deviceId) async {
    final pending = LocalDb.getPendingEvents(limit: 100);
    if (pending.isEmpty) return;
    try {
      await _api.post(Endpoints.deviceLog(deviceId), data: pending);
      await LocalDb.clearEvents(pending.length);
      Logger.debug('Flushed ${pending.length} events', 'EventReporter');
    } catch (e) {
      Logger.warn('Event flush failed (will retry)', 'EventReporter');
    }
  }
}
