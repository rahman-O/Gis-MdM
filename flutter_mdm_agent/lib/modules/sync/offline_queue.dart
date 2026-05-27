import 'package:uuid/uuid.dart';

import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/local_db.dart';
import '../../core/utils/logger.dart';

/// Represents a single queued item waiting to be sent to the server.
class QueueItem {
  final String id;
  final String type;
  final Map<String, dynamic> data;
  final int createdAt;

  QueueItem({
    required this.id,
    required this.type,
    required this.data,
    required this.createdAt,
  });

  factory QueueItem.fromJson(Map<String, dynamic> json) {
    return QueueItem(
      id: json['id'] as String,
      type: json['type'] as String,
      data: (json['data'] as Map).cast<String, dynamic>(),
      createdAt: json['createdAt'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type,
      'data': data,
      'createdAt': createdAt,
    };
  }
}

/// Offline queue for storing events/telemetry when the device is offline.
///
/// Items are persisted in LocalDb and flushed to the server when
/// connectivity is restored.
class OfflineQueue {
  static const _uuid = Uuid();

  /// Add an item to the offline queue.
  Future<void> enqueue(String type, Map<String, dynamic> data) async {
    final item = QueueItem(
      id: _uuid.v4(),
      type: type,
      data: data,
      createdAt: DateTime.now().millisecondsSinceEpoch,
    );

    await LocalDb.addEvent(item.toJson());
    Logger.debug('Enqueued offline item: type=$type, id=${item.id}', 'OfflineQueue');
  }

  /// Get all pending items from the queue.
  Future<List<QueueItem>> getPending({int limit = 100}) async {
    final events = LocalDb.getPendingEvents(limit: limit);
    return events.map((e) => QueueItem.fromJson(e)).toList();
  }

  /// Mark items as sent and remove them from the queue.
  Future<void> markSent(List<String> ids) async {
    if (ids.isEmpty) return;
    await LocalDb.clearEvents(ids.length);
    Logger.debug('Cleared ${ids.length} sent items from queue', 'OfflineQueue');
  }

  /// Flush all pending items to the server.
  ///
  /// Sends each item based on its type and removes successfully sent items.
  Future<void> flushAll(ApiClient api, String deviceId) async {
    final pending = await getPending();
    if (pending.isEmpty) {
      Logger.debug('No pending items to flush', 'OfflineQueue');
      return;
    }

    Logger.info('Flushing ${pending.length} offline items', 'OfflineQueue');
    final sentIds = <String>[];

    for (final item in pending) {
      try {
        await _sendItem(api, deviceId, item);
        sentIds.add(item.id);
      } catch (e) {
        Logger.warn(
          'Failed to flush item ${item.id}: $e',
          'OfflineQueue',
        );
        // Stop on first failure to preserve order
        break;
      }
    }

    if (sentIds.isNotEmpty) {
      await markSent(sentIds);
      Logger.info('Flushed ${sentIds.length}/${pending.length} items', 'OfflineQueue');
    }
  }

  /// Send a single queue item to the appropriate endpoint.
  Future<void> _sendItem(ApiClient api, String deviceId, QueueItem item) async {
    switch (item.type) {
      case 'device_info':
        await api.post(
          Endpoints.syncInfo,
          data: {'deviceId': deviceId, ...item.data},
        );
      case 'telemetry':
        await api.post(
          Endpoints.deviceInfo(deviceId),
          data: item.data,
        );
      default:
        await api.post(
          Endpoints.syncInfo,
          data: {'deviceId': deviceId, 'type': item.type, ...item.data},
        );
    }
  }
}
