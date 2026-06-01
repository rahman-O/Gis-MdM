import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:mdm_agent/modules/telemetry/location_queue.dart';
import 'package:mdm_agent/modules/telemetry/models/location_models.dart';

/// Unit tests for [LocationQueue], focusing on batch dequeue behavior.
///
/// These tests verify:
/// - `dequeueBatch(count)` returns records in ascending timestamp order
///   (oldest first) — **Requirement 5.2**
/// - Records remain in the queue until explicitly removed via `removeBatch()`
///   — **Requirement 5.6**
/// - `removeBatch()` removes the oldest records matching the dequeued batch
void main() {
  late LocationQueue queue;
  late Box<Map> box;

  setUp(() async {
    // Initialize Hive with a unique temporary path per test.
    final tempPath = '/tmp/hive_test_${DateTime.now().microsecondsSinceEpoch}';
    Hive.init(tempPath);

    // Open an unencrypted box for testing (bypasses secure storage).
    box = await Hive.openBox<Map>('test_location_queue');

    // Use the test-friendly constructor.
    queue = LocationQueue.withBox(box);
  });

  tearDown(() async {
    await queue.dispose();
    await Hive.deleteFromDisk();
  });

  /// Helper to create a [QueuedLocationRecord] with a specific timestamp.
  QueuedLocationRecord createRecord({
    required int timestamp,
    double latitude = 33.3152,
    double longitude = 44.3661,
  }) {
    return QueuedLocationRecord(
      latitude: latitude,
      longitude: longitude,
      accuracy: 10.0,
      speed: 1.5,
      batteryLevel: 80,
      networkType: 'wifi',
      trackingMode: 'normal',
      timestamp: timestamp,
      deviceId: 'test-device-001',
    );
  }

  group('dequeueBatch - chronological ordering (Requirement 5.2)', () {
    test('returns empty list when queue is empty', () {
      final batch = queue.dequeueBatch(10);
      expect(batch, isEmpty);
    });

    test('returns records in ascending timestamp order (oldest first)', () async {
      // Insert records out of chronological order.
      await queue.enqueue(createRecord(timestamp: 3000)); // newest
      await queue.enqueue(createRecord(timestamp: 1000)); // oldest
      await queue.enqueue(createRecord(timestamp: 2000)); // middle

      final batch = queue.dequeueBatch(3);

      expect(batch.length, equals(3));
      expect(batch[0].timestamp, equals(1000));
      expect(batch[1].timestamp, equals(2000));
      expect(batch[2].timestamp, equals(3000));
    });

    test('returns at most count records (oldest subset)', () async {
      for (int i = 1; i <= 5; i++) {
        await queue.enqueue(createRecord(timestamp: i * 1000));
      }

      final batch = queue.dequeueBatch(3);

      expect(batch.length, equals(3));
      expect(batch[0].timestamp, equals(1000));
      expect(batch[1].timestamp, equals(2000));
      expect(batch[2].timestamp, equals(3000));
    });

    test('returns all records when count exceeds queue size', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));

      final batch = queue.dequeueBatch(100);

      expect(batch.length, equals(2));
      expect(batch[0].timestamp, equals(1000));
      expect(batch[1].timestamp, equals(2000));
    });

    test('maintains order regardless of insertion order', () async {
      // Insert in reverse order.
      await queue.enqueue(createRecord(timestamp: 5000));
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 3000));
      await queue.enqueue(createRecord(timestamp: 2000));
      await queue.enqueue(createRecord(timestamp: 4000));

      final batch = queue.dequeueBatch(5);

      for (int i = 0; i < batch.length - 1; i++) {
        expect(
          batch[i].timestamp <= batch[i + 1].timestamp,
          isTrue,
          reason: 'Record at index $i (ts=${batch[i].timestamp}) should be '
              '<= record at index ${i + 1} (ts=${batch[i + 1].timestamp})',
        );
      }
    });

    test('handles count of 0 gracefully', () async {
      await queue.enqueue(createRecord(timestamp: 1000));

      final batch = queue.dequeueBatch(0);
      expect(batch, isEmpty);
    });

    test('handles negative count gracefully', () async {
      await queue.enqueue(createRecord(timestamp: 1000));

      final batch = queue.dequeueBatch(-5);
      expect(batch, isEmpty);
    });

    test('preserves all record fields correctly', () async {
      final record = QueuedLocationRecord(
        latitude: 33.312456,
        longitude: 44.366789,
        accuracy: 12.5,
        speed: 3.2,
        altitude: 45.0,
        batteryLevel: 72,
        networkType: 'cellular',
        trackingMode: 'lowPower',
        timestamp: 1716595200000,
        deviceId: 'device-xyz',
      );

      await queue.enqueue(record);
      final batch = queue.dequeueBatch(1);

      expect(batch.length, equals(1));
      final dequeued = batch[0];
      expect(dequeued.latitude, equals(33.312456));
      expect(dequeued.longitude, equals(44.366789));
      expect(dequeued.accuracy, equals(12.5));
      expect(dequeued.speed, equals(3.2));
      expect(dequeued.altitude, equals(45.0));
      expect(dequeued.batteryLevel, equals(72));
      expect(dequeued.networkType, equals('cellular'));
      expect(dequeued.trackingMode, equals('lowPower'));
      expect(dequeued.timestamp, equals(1716595200000));
      expect(dequeued.deviceId, equals('device-xyz'));
    });
  });

  group('dequeueBatch - non-destructive (Requirement 5.6)', () {
    test('does NOT remove records from queue', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));
      await queue.enqueue(createRecord(timestamp: 3000));

      queue.dequeueBatch(3);

      // Queue should still have all records.
      expect(queue.length, equals(3));
    });

    test('repeated dequeueBatch returns same records', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));

      final batch1 = queue.dequeueBatch(2);
      final batch2 = queue.dequeueBatch(2);
      final batch3 = queue.dequeueBatch(2);

      expect(batch1[0].timestamp, equals(batch2[0].timestamp));
      expect(batch1[1].timestamp, equals(batch2[1].timestamp));
      expect(batch2[0].timestamp, equals(batch3[0].timestamp));
      expect(batch2[1].timestamp, equals(batch3[1].timestamp));
    });

    test('records persist until removeBatch is called', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));

      // Multiple dequeue calls — records persist.
      queue.dequeueBatch(2);
      queue.dequeueBatch(2);
      queue.dequeueBatch(2);

      expect(queue.length, equals(2));

      // Only after explicit removal.
      await queue.removeBatch(1);
      expect(queue.length, equals(1));
    });
  });

  group('removeBatch - acknowledgment-based removal (Requirement 5.6)', () {
    test('removes oldest records after acknowledgment', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));
      await queue.enqueue(createRecord(timestamp: 3000));

      // Acknowledge 2 oldest.
      final removed = await queue.removeBatch(2);
      expect(removed, equals(2));
      expect(queue.length, equals(1));

      // Remaining record should be the newest.
      final remaining = queue.dequeueBatch(1);
      expect(remaining[0].timestamp, equals(3000));
    });

    test('removeBatch with count exceeding queue size removes all', () async {
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 2000));

      final removed = await queue.removeBatch(100);
      expect(removed, equals(2));
      expect(queue.isEmpty, isTrue);
    });

    test('removeBatch with 0 removes nothing', () async {
      await queue.enqueue(createRecord(timestamp: 1000));

      final removed = await queue.removeBatch(0);
      expect(removed, equals(0));
      expect(queue.length, equals(1));
    });

    test('removeBatch removes in chronological order (oldest first)', () async {
      // Insert out of order.
      await queue.enqueue(createRecord(timestamp: 5000));
      await queue.enqueue(createRecord(timestamp: 1000));
      await queue.enqueue(createRecord(timestamp: 3000));

      // Remove 1 — should remove the oldest (1000).
      await queue.removeBatch(1);
      expect(queue.length, equals(2));

      final remaining = queue.dequeueBatch(2);
      expect(remaining[0].timestamp, equals(3000));
      expect(remaining[1].timestamp, equals(5000));
    });
  });

  group('Full batch upload workflow simulation', () {
    test('dequeue → upload → acknowledge → dequeue next batch', () async {
      // Simulate collecting 5 location records over time.
      for (int i = 1; i <= 5; i++) {
        await queue.enqueue(createRecord(timestamp: i * 1000));
      }
      expect(queue.length, equals(5));

      // Step 1: Dequeue batch of 3 for upload.
      final batch1 = queue.dequeueBatch(3);
      expect(batch1.length, equals(3));
      expect(batch1[0].timestamp, equals(1000));
      expect(batch1[1].timestamp, equals(2000));
      expect(batch1[2].timestamp, equals(3000));

      // Records still in queue (not yet acknowledged).
      expect(queue.length, equals(5));

      // Step 2: Server acknowledges — remove batch.
      await queue.removeBatch(3);
      expect(queue.length, equals(2));

      // Step 3: Dequeue next batch.
      final batch2 = queue.dequeueBatch(3);
      expect(batch2.length, equals(2));
      expect(batch2[0].timestamp, equals(4000));
      expect(batch2[1].timestamp, equals(5000));

      // Step 4: Server acknowledges — remove.
      await queue.removeBatch(2);
      expect(queue.isEmpty, isTrue);
    });

    test('connectivity loss mid-upload preserves all records', () async {
      for (int i = 1; i <= 5; i++) {
        await queue.enqueue(createRecord(timestamp: i * 1000));
      }

      // Dequeue batch for upload.
      final batch = queue.dequeueBatch(5);
      expect(batch.length, equals(5));

      // Simulate connectivity loss — do NOT call removeBatch.
      expect(queue.length, equals(5));

      // On retry, dequeue again — same records available.
      final retryBatch = queue.dequeueBatch(5);
      expect(retryBatch.length, equals(5));
      expect(retryBatch[0].timestamp, equals(1000));
      expect(retryBatch[4].timestamp, equals(5000));
    });

    test('partial acknowledgment removes only acknowledged records', () async {
      for (int i = 1; i <= 10; i++) {
        await queue.enqueue(createRecord(timestamp: i * 1000));
      }

      // Dequeue 5 for upload.
      final batch = queue.dequeueBatch(5);
      expect(batch.length, equals(5));

      // Server acknowledges only 3 of the 5.
      await queue.removeBatch(3);
      expect(queue.length, equals(7));

      // Next dequeue should start from the 4th oldest.
      final nextBatch = queue.dequeueBatch(5);
      expect(nextBatch[0].timestamp, equals(4000));
      expect(nextBatch[1].timestamp, equals(5000));
      expect(nextBatch[2].timestamp, equals(6000));
      expect(nextBatch[3].timestamp, equals(7000));
      expect(nextBatch[4].timestamp, equals(8000));
    });
  });
}
