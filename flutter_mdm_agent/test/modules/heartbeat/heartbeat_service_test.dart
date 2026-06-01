import 'package:flutter_test/flutter_test.dart';
import 'package:mdm_agent/modules/telemetry/telemetry_data.dart';

/// Tests for the heartbeat payload construction logic.
///
/// Since HeartbeatService depends on platform plugins (device_info_plus,
/// battery_plus, etc.), we test the data validation logic that ensures
/// no null/empty values are sent to the server.
void main() {
  group('SystemInfo data integrity', () {
    test('SystemInfo fields are never null when constructed with values', () {
      final info = SystemInfo(
        model: 'SM-A525F',
        manufacturer: 'Samsung',
        androidVersion: '13',
        sdkInt: 33,
        serial: 'R58N12345',
        uptimeMillis: 123456789,
      );

      expect(info.model, isNotNull);
      expect(info.model, isNotEmpty);
      expect(info.manufacturer, isNotNull);
      expect(info.manufacturer, isNotEmpty);
      expect(info.androidVersion, isNotNull);
      expect(info.androidVersion, isNotEmpty);
      expect(info.serial, isNotNull);
    });

    test('SystemInfo.toJson() never produces null values for string fields', () {
      final info = SystemInfo(
        model: 'SM-A525F',
        manufacturer: 'Samsung',
        androidVersion: '13',
        sdkInt: 33,
        serial: 'R58N12345',
        uptimeMillis: 123456789,
      );

      final json = info.toJson();
      expect(json['model'], isNotNull);
      expect(json['manufacturer'], isNotNull);
      expect(json['androidVersion'], isNotNull);
      expect(json['serial'], isNotNull);
    });

    test('SystemInfo.fromJson() handles null values with defaults', () {
      final json = <String, dynamic>{
        'model': null,
        'manufacturer': null,
        'androidVersion': null,
        'sdkInt': null,
        'serial': null,
        'uptimeMillis': null,
      };

      final info = SystemInfo.fromJson(json);
      expect(info.model, equals(''));
      expect(info.manufacturer, equals(''));
      expect(info.androidVersion, equals(''));
      expect(info.sdkInt, equals(0));
      expect(info.serial, equals(''));
      expect(info.uptimeMillis, equals(0));
    });

    test('SystemInfo.fromJson() handles missing keys with defaults', () {
      final json = <String, dynamic>{};

      final info = SystemInfo.fromJson(json);
      expect(info.model, equals(''));
      expect(info.manufacturer, equals(''));
      expect(info.androidVersion, equals(''));
      expect(info.sdkInt, equals(0));
      expect(info.serial, equals(''));
      expect(info.uptimeMillis, equals(0));
    });

    test('SystemInfo with "unknown" serial (Android 10+ restriction) is valid', () {
      final info = SystemInfo(
        model: 'Pixel 6',
        manufacturer: 'Google',
        androidVersion: '14',
        sdkInt: 34,
        serial: 'unknown', // Android 10+ returns this for non-rooted devices
        uptimeMillis: 500000,
      );

      expect(info.serial, equals('unknown'));
      expect(info.model, isNotEmpty);
    });
  });

  group('Heartbeat payload null safety', () {
    /// Simulates the _ensureNonNull logic from HeartbeatService.
    String ensureNonNull(String? value, String fallback) {
      if (value == null || value.trim().isEmpty) {
        return fallback;
      }
      return value;
    }

    test('ensureNonNull returns value when non-null and non-empty', () {
      expect(ensureNonNull('Samsung', 'unknown'), equals('Samsung'));
      expect(ensureNonNull('SM-A525F', 'unknown'), equals('SM-A525F'));
    });

    test('ensureNonNull returns fallback for null', () {
      expect(ensureNonNull(null, 'unknown'), equals('unknown'));
    });

    test('ensureNonNull returns fallback for empty string', () {
      expect(ensureNonNull('', 'unknown'), equals('unknown'));
    });

    test('ensureNonNull returns fallback for whitespace-only string', () {
      expect(ensureNonNull('   ', 'unknown'), equals('unknown'));
      expect(ensureNonNull('\t', 'unknown'), equals('unknown'));
    });

    test('heartbeat body never contains null for device info fields', () {
      // Simulate what HeartbeatService.sendHeartbeat() does
      final system = SystemInfo(
        model: 'SM-A525F',
        manufacturer: 'Samsung',
        androidVersion: '13',
        sdkInt: 33,
        serial: 'unavailable',
        uptimeMillis: 123456789,
      );

      final body = <String, dynamic>{
        'deviceId': '351906200367061',
        'batteryLevel': 93,
        'imei': '351906200367061',
        'model': ensureNonNull(system.model, 'unknown'),
        'manufacturer': ensureNonNull(system.manufacturer, 'unknown'),
        'androidVersion': ensureNonNull(system.androidVersion, 'unknown'),
        'serial': ensureNonNull(system.serial, 'unavailable'),
        'launcherVersion': '1.0.0',
      };

      // Verify no null values in the payload
      for (final entry in body.entries) {
        expect(entry.value, isNotNull,
            reason: '${entry.key} should never be null');
      }

      // Verify string fields are non-empty
      expect(body['model'], isNotEmpty);
      expect(body['manufacturer'], isNotEmpty);
      expect(body['androidVersion'], isNotEmpty);
      expect(body['serial'], isNotEmpty);
    });

    test('heartbeat body uses fallback when SystemCollector returns empty', () {
      // Simulate a scenario where device_info_plus returns empty strings
      final system = SystemInfo(
        model: '',
        manufacturer: '',
        androidVersion: '',
        sdkInt: 0,
        serial: '',
        uptimeMillis: 0,
      );

      final body = <String, dynamic>{
        'model': ensureNonNull(system.model, 'unknown'),
        'manufacturer': ensureNonNull(system.manufacturer, 'unknown'),
        'androidVersion': ensureNonNull(system.androidVersion, 'unknown'),
        'serial': ensureNonNull(system.serial, 'unavailable'),
      };

      expect(body['model'], equals('unknown'));
      expect(body['manufacturer'], equals('unknown'));
      expect(body['androidVersion'], equals('unknown'));
      expect(body['serial'], equals('unavailable'));
    });
  });
}
