import 'package:flutter_test/flutter_test.dart';
import 'package:mdm_agent/modules/telemetry/models/location_models.dart';
import 'package:mdm_agent/modules/telemetry/smart_location_engine.dart';
import 'package:mdm_agent/modules/telemetry/telemetry_data.dart';

void main() {
  group('SmartLocationEngine', () {
    late SmartLocationEngine engine;

    setUp(() {
      engine = SmartLocationEngine();
    });

    group('haversineDistance', () {
      test('returns 0 for identical points', () {
        final distance = SmartLocationEngine.haversineDistance(
          33.3152, 44.3661,
          33.3152, 44.3661,
        );
        expect(distance, equals(0.0));
      });

      test('calculates ~100m for nearby points', () {
        // ~0.0009 degrees latitude ≈ 100m
        final distance = SmartLocationEngine.haversineDistance(
          33.3152, 44.3661,
          33.3161, 44.3661,
        );
        expect(distance, closeTo(100.0, 5.0));
      });

      test('calculates ~111km for 1 degree latitude at equator', () {
        final distance = SmartLocationEngine.haversineDistance(
          0.0, 0.0,
          1.0, 0.0,
        );
        expect(distance, closeTo(111195.0, 100.0));
      });

      test('calculates ~111km for 1 degree longitude at equator', () {
        final distance = SmartLocationEngine.haversineDistance(
          0.0, 0.0,
          0.0, 1.0,
        );
        expect(distance, closeTo(111195.0, 100.0));
      });

      test('handles antipodal points (max distance ~20,000km)', () {
        final distance = SmartLocationEngine.haversineDistance(
          0.0, 0.0,
          0.0, 180.0,
        );
        // Half circumference of Earth ≈ 20,015 km
        expect(distance, closeTo(20015087.0, 1000.0));
      });

      test('handles negative coordinates', () {
        final distance = SmartLocationEngine.haversineDistance(
          -33.8688, 151.2093, // Sydney
          51.5074, -0.1278, // London
        );
        // Sydney to London ≈ 16,993 km
        expect(distance, closeTo(16993000.0, 50000.0));
      });

      test('is symmetric (distance A→B == distance B→A)', () {
        final ab = SmartLocationEngine.haversineDistance(
          33.3152, 44.3661,
          30.5085, 47.7834,
        );
        final ba = SmartLocationEngine.haversineDistance(
          30.5085, 47.7834,
          33.3152, 44.3661,
        );
        expect(ab, equals(ba));
      });
    });

    group('validateCoordinates', () {
      test('accepts valid coordinates with good accuracy', () {
        expect(
          SmartLocationEngine.validateCoordinates(33.3152, 44.3661, 10.0),
          isTrue,
        );
      });

      test('accepts edge latitude values (-90 and 90)', () {
        expect(SmartLocationEngine.validateCoordinates(90.0, 0.0, 10.0), isTrue);
        expect(SmartLocationEngine.validateCoordinates(-90.0, 0.0, 10.0), isTrue);
      });

      test('accepts edge longitude values (-180 and 180)', () {
        expect(SmartLocationEngine.validateCoordinates(0.0, 180.0, 10.0), isTrue);
        expect(SmartLocationEngine.validateCoordinates(0.0, -180.0, 10.0), isTrue);
      });

      test('accepts accuracy exactly at 500m', () {
        expect(
          SmartLocationEngine.validateCoordinates(33.3, 44.3, 500.0),
          isTrue,
        );
      });

      test('rejects latitude > 90', () {
        expect(
          SmartLocationEngine.validateCoordinates(90.001, 44.3, 10.0),
          isFalse,
        );
      });

      test('rejects latitude < -90', () {
        expect(
          SmartLocationEngine.validateCoordinates(-90.001, 44.3, 10.0),
          isFalse,
        );
      });

      test('rejects longitude > 180', () {
        expect(
          SmartLocationEngine.validateCoordinates(33.3, 180.001, 10.0),
          isFalse,
        );
      });

      test('rejects longitude < -180', () {
        expect(
          SmartLocationEngine.validateCoordinates(33.3, -180.001, 10.0),
          isFalse,
        );
      });

      test('rejects accuracy > 500m', () {
        expect(
          SmartLocationEngine.validateCoordinates(33.3, 44.3, 500.1),
          isFalse,
        );
      });

      test('accepts zero coordinates', () {
        expect(
          SmartLocationEngine.validateCoordinates(0.0, 0.0, 0.0),
          isTrue,
        );
      });
    });

    group('shouldStore', () {
      test('stores unconditionally when lastStored is null (first reading)', () {
        final newReading = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );

        expect(engine.shouldStore(null, newReading, 50), isTrue);
      });

      test('stores when distance >= threshold', () {
        final lastStored = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        // ~100m away
        final newReading = LocationInfo(
          latitude: 33.3161,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        expect(engine.shouldStore(lastStored, newReading, 50), isTrue);
      });

      test('discards when distance < threshold', () {
        final lastStored = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        // ~10m away (very close)
        final newReading = LocationInfo(
          latitude: 33.31529,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        expect(engine.shouldStore(lastStored, newReading, 50), isFalse);
      });

      test('stores when distance equals threshold', () {
        final lastStored = LocationInfo(
          latitude: 0.0,
          longitude: 0.0,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        // ~100m away
        final newReading = LocationInfo(
          latitude: 0.0009,
          longitude: 0.0,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        // Calculate actual distance and use it as threshold
        final distance = SmartLocationEngine.haversineDistance(
          0.0, 0.0, 0.0009, 0.0,
        );
        // Use floor of distance as threshold — distance >= threshold should be true
        expect(engine.shouldStore(lastStored, newReading, distance.floor()), isTrue);
      });

      test('discards when same location is reported again', () {
        final lastStored = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        final newReading = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        expect(engine.shouldStore(lastStored, newReading, 50), isFalse);
      });
    });

    group('setThreshold', () {
      test('accepts value within valid range', () {
        expect(engine.setThreshold(100), isTrue);
        expect(engine.state.distanceThreshold, equals(100));
      });

      test('accepts minimum threshold (10m)', () {
        expect(engine.setThreshold(10), isTrue);
        expect(engine.state.distanceThreshold, equals(10));
      });

      test('accepts maximum threshold (10000m)', () {
        expect(engine.setThreshold(10000), isTrue);
        expect(engine.state.distanceThreshold, equals(10000));
      });

      test('rejects value below minimum and retains current', () {
        engine.setThreshold(100); // Set to 100 first
        expect(engine.setThreshold(9), isFalse);
        expect(engine.state.distanceThreshold, equals(100));
      });

      test('rejects value above maximum and retains current', () {
        engine.setThreshold(100); // Set to 100 first
        expect(engine.setThreshold(10001), isFalse);
        expect(engine.state.distanceThreshold, equals(100));
      });

      test('default threshold is 50m', () {
        expect(engine.state.distanceThreshold, equals(50));
      });
    });

    group('processReading', () {
      test('stores valid first reading', () {
        final reading = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );

        expect(engine.processReading(reading), isTrue);
        expect(engine.state.lastStoredLocation, equals(reading));
      });

      test('rejects reading with invalid coordinates', () {
        final reading = LocationInfo(
          latitude: 91.0,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );

        expect(engine.processReading(reading), isFalse);
        expect(engine.state.lastStoredLocation, isNull);
      });

      test('rejects reading with poor accuracy', () {
        final reading = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 501.0,
          timestamp: 1000000,
        );

        expect(engine.processReading(reading), isFalse);
        expect(engine.state.lastStoredLocation, isNull);
      });

      test('discards reading within threshold distance', () {
        // Store first reading
        final first = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        engine.processReading(first);

        // Try to store a very close reading (~1m away)
        final second = LocationInfo(
          latitude: 33.31521,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        expect(engine.processReading(second), isFalse);
        expect(engine.state.lastStoredLocation, equals(first));
      });

      test('stores reading beyond threshold distance', () {
        // Store first reading
        final first = LocationInfo(
          latitude: 33.3152,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1000000,
        );
        engine.processReading(first);

        // Store a reading ~100m away
        final second = LocationInfo(
          latitude: 33.3161,
          longitude: 44.3661,
          accuracy: 10.0,
          timestamp: 1001000,
        );

        expect(engine.processReading(second), isTrue);
        expect(engine.state.lastStoredLocation, equals(second));
      });
    });

    group('updateMovementState', () {
      test('initial state is stationary', () {
        expect(engine.state.movementState, equals(DeviceMovementState.stationary));
      });

      test('transitions to moving when speed > 2 m/s', () {
        final now = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(3.0, now);

        expect(result, equals(DeviceMovementState.moving));
        expect(engine.state.movementState, equals(DeviceMovementState.moving));
      });

      test('transitions to moving immediately on first speed > 2 m/s', () {
        final now = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(2.1, now);

        expect(result, equals(DeviceMovementState.moving));
        expect(engine.state.movementState, equals(DeviceMovementState.moving));
      });

      test('remains stationary when speed <= 2 m/s and not yet 60s', () {
        // Start as moving first
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);
        expect(engine.state.movementState, equals(DeviceMovementState.moving));

        // Speed drops below threshold
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        final result = engine.updateMovementState(1.5, t1);

        // Still moving — hasn't been 60s yet
        expect(result, equals(DeviceMovementState.moving));
        expect(engine.state.movementState, equals(DeviceMovementState.moving));
      });

      test('transitions to stationary after speed <= 2 m/s for 60 consecutive seconds', () {
        // Start as moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        // Speed drops below threshold
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        engine.updateMovementState(1.0, t1);

        // 60 seconds later, still below threshold
        final t2 = DateTime(2024, 1, 1, 12, 1, 10);
        final result = engine.updateMovementState(0.5, t2);

        expect(result, equals(DeviceMovementState.stationary));
        expect(engine.state.movementState, equals(DeviceMovementState.stationary));
      });

      test('resets stationary timer when speed goes above threshold during countdown', () {
        // Start as moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        // Speed drops below threshold
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        engine.updateMovementState(1.0, t1);

        // 30 seconds later, speed goes back up (resets the timer)
        final t2 = DateTime(2024, 1, 1, 12, 0, 40);
        engine.updateMovementState(3.0, t2);
        expect(engine.state.movementState, equals(DeviceMovementState.moving));

        // Speed drops again
        final t3 = DateTime(2024, 1, 1, 12, 0, 50);
        engine.updateMovementState(1.0, t3);

        // 50 seconds later (not 60 from t3)
        final t4 = DateTime(2024, 1, 1, 12, 1, 40);
        final result = engine.updateMovementState(0.5, t4);

        // Should still be moving — only 50s since t3
        expect(result, equals(DeviceMovementState.moving));
      });

      test('null speed preserves moving state', () {
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);
        expect(engine.state.movementState, equals(DeviceMovementState.moving));

        final t1 = DateTime(2024, 1, 1, 12, 0, 30);
        final result = engine.updateMovementState(null, t1);

        expect(result, equals(DeviceMovementState.moving));
        expect(engine.state.movementState, equals(DeviceMovementState.moving));
      });

      test('null speed preserves stationary state', () {
        // Engine starts as stationary
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(null, t0);

        expect(result, equals(DeviceMovementState.stationary));
        expect(engine.state.movementState, equals(DeviceMovementState.stationary));
      });

      test('speed exactly at 2 m/s does not trigger moving', () {
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(2.0, t0);

        // 2.0 m/s is NOT > 2 m/s, so stays stationary
        expect(result, equals(DeviceMovementState.stationary));
        expect(engine.state.movementState, equals(DeviceMovementState.stationary));
      });

      test('speed just above threshold (2.01) triggers moving', () {
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(2.01, t0);

        expect(result, equals(DeviceMovementState.moving));
        expect(engine.state.movementState, equals(DeviceMovementState.moving));
      });

      test('transitions stationary exactly at 60 seconds', () {
        // Start as moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        // Speed drops below threshold
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        engine.updateMovementState(1.0, t1);

        // Exactly 60 seconds after t1
        final t2 = DateTime(2024, 1, 1, 12, 1, 10);
        final result = engine.updateMovementState(1.0, t2);

        expect(result, equals(DeviceMovementState.stationary));
      });

      test('does not transition stationary at 59 seconds', () {
        // Start as moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        // Speed drops below threshold
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        engine.updateMovementState(1.0, t1);

        // 59 seconds after t1 — not enough
        final t2 = DateTime(2024, 1, 1, 12, 1, 9);
        final result = engine.updateMovementState(1.0, t2);

        expect(result, equals(DeviceMovementState.moving));
      });

      test('multiple speed readings below threshold accumulate correctly', () {
        // Start as moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        // Series of low-speed readings over 60+ seconds
        final t1 = DateTime(2024, 1, 1, 12, 0, 10);
        engine.updateMovementState(1.5, t1);

        final t2 = DateTime(2024, 1, 1, 12, 0, 30);
        expect(engine.updateMovementState(0.5, t2), equals(DeviceMovementState.moving));

        final t3 = DateTime(2024, 1, 1, 12, 0, 50);
        expect(engine.updateMovementState(1.0, t3), equals(DeviceMovementState.moving));

        // 60 seconds from t1 (first low-speed reading)
        final t4 = DateTime(2024, 1, 1, 12, 1, 10);
        expect(engine.updateMovementState(0.8, t4), equals(DeviceMovementState.stationary));
      });

      test('records lastSpeedAboveThreshold when moving', () {
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        expect(engine.state.lastSpeedAboveThreshold, equals(t0));
      });

      test('stationary to moving transition clears stationarySince', () {
        // Start stationary, then go moving
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(1.0, t0);
        expect(engine.state.stationarySince, isNotNull);

        final t1 = DateTime(2024, 1, 1, 12, 0, 30);
        engine.updateMovementState(5.0, t1);

        expect(engine.state.stationarySince, isNull);
      });

      test('zero speed is treated as stationary (below threshold)', () {
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        final result = engine.updateMovementState(0.0, t0);

        expect(result, equals(DeviceMovementState.stationary));
      });
    });

    group('stationary rate cap (Requirement 13.5)', () {
      /// Helper: creates a LocationInfo at a given latitude offset and timestamp.
      /// Each reading is far enough apart (>50m) to pass the distance filter.
      LocationInfo _readingAt(int index, int timestampMs) {
        // Each index adds ~111m of latitude offset (0.001 degrees ≈ 111m)
        return LocationInfo(
          latitude: 33.0 + (index * 0.001),
          longitude: 44.0,
          accuracy: 10.0,
          timestamp: timestampMs,
        );
      }

      test('allows up to 12 records per hour while stationary', () {
        // Engine starts in stationary state by default.
        // Use a base timestamp at the start of an hour.
        const baseTs = 1704067200000; // 2024-01-01 00:00:00 UTC

        // Store 12 readings within the same hour — all should be accepted.
        for (int i = 0; i < 12; i++) {
          final reading = _readingAt(i, baseTs + (i * 60000)); // 1 min apart
          expect(
            engine.processReading(reading),
            isTrue,
            reason: 'Reading $i should be stored (within 12/hour cap)',
          );
        }

        expect(engine.state.stationaryRecordsThisHour, equals(12));
      });

      test('discards the 13th record in the same hour while stationary', () {
        const baseTs = 1704067200000; // 2024-01-01 00:00:00 UTC

        // Store 12 readings.
        for (int i = 0; i < 12; i++) {
          engine.processReading(_readingAt(i, baseTs + (i * 60000)));
        }

        // 13th reading in the same hour — should be discarded.
        final thirteenth = _readingAt(12, baseTs + (12 * 60000));
        expect(engine.processReading(thirteenth), isFalse);
      });

      test('resets counter when a new hour begins', () {
        const baseTs = 1704067200000; // 2024-01-01 00:00:00 UTC
        const nextHourTs = baseTs + 3600000; // 2024-01-01 01:00:00 UTC

        // Fill up the first hour.
        for (int i = 0; i < 12; i++) {
          engine.processReading(_readingAt(i, baseTs + (i * 60000)));
        }

        // Verify cap is reached.
        expect(
          engine.processReading(_readingAt(12, baseTs + (12 * 60000))),
          isFalse,
        );

        // New hour — should reset and allow storage.
        final newHourReading = _readingAt(13, nextHourTs);
        expect(engine.processReading(newHourReading), isTrue);
        expect(engine.state.stationaryRecordsThisHour, equals(1));
      });

      test('does not apply rate cap when device is moving', () {
        // Transition to moving state.
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);
        expect(engine.state.movementState, equals(DeviceMovementState.moving));

        const baseTs = 1704110400000; // 2024-01-01 12:00:00 UTC

        // Store more than 12 readings — all should be accepted since moving.
        for (int i = 0; i < 15; i++) {
          final reading = _readingAt(i, baseTs + (i * 60000));
          expect(
            engine.processReading(reading),
            isTrue,
            reason: 'Reading $i should be stored (moving — no rate cap)',
          );
        }
      });

      test('rate cap applies immediately when device becomes stationary', () {
        // Start moving, store some readings.
        final t0 = DateTime(2024, 1, 1, 12, 0, 0);
        engine.updateMovementState(5.0, t0);

        const baseTs = 1704110400000; // 2024-01-01 12:00:00 UTC

        // Store 5 readings while moving.
        for (int i = 0; i < 5; i++) {
          engine.processReading(_readingAt(i, baseTs + (i * 60000)));
        }

        // Transition to stationary.
        final t1 = DateTime(2024, 1, 1, 12, 5, 0);
        engine.updateMovementState(1.0, t1);
        final t2 = DateTime(2024, 1, 1, 12, 6, 0);
        engine.updateMovementState(0.5, t2); // 60s elapsed → stationary
        expect(engine.state.movementState, equals(DeviceMovementState.stationary));

        // Now store up to 12 readings while stationary in the same hour.
        for (int i = 5; i < 17; i++) {
          engine.processReading(_readingAt(i, baseTs + (i * 60000)));
        }

        // The 13th stationary reading should be discarded.
        final extra = _readingAt(17, baseTs + (17 * 60000));
        expect(engine.processReading(extra), isFalse);
      });

      test('first reading in stationary state always counts as record 1', () {
        const baseTs = 1704067200000;

        // First reading ever — should be stored and count as 1.
        final first = _readingAt(0, baseTs);
        expect(engine.processReading(first), isTrue);
        expect(engine.state.stationaryRecordsThisHour, equals(1));
        expect(engine.state.currentHourTimestamp, equals(baseTs));
      });

      test('readings spanning multiple hours each get their own counter', () {
        const hour1 = 1704067200000; // 00:00 UTC
        const hour2 = hour1 + 3600000; // 01:00 UTC
        const hour3 = hour2 + 3600000; // 02:00 UTC

        // Fill hour 1.
        for (int i = 0; i < 12; i++) {
          engine.processReading(_readingAt(i, hour1 + (i * 60000)));
        }
        expect(engine.processReading(_readingAt(12, hour1 + 720000)), isFalse);

        // Hour 2 — fresh counter.
        for (int i = 13; i < 25; i++) {
          final reading = _readingAt(i, hour2 + ((i - 13) * 60000));
          expect(engine.processReading(reading), isTrue);
        }
        expect(
          engine.processReading(_readingAt(25, hour2 + 720000)),
          isFalse,
        );

        // Hour 3 — fresh counter again.
        final hour3Reading = _readingAt(26, hour3);
        expect(engine.processReading(hour3Reading), isTrue);
        expect(engine.state.stationaryRecordsThisHour, equals(1));
      });
    });
  });
}
