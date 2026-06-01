import 'dart:convert';

import 'package:crypto/crypto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mdm_agent/modules/security/location_signer.dart';

void main() {
  group('LocationSigner', () {
    const testToken = 'test-enrollment-token-secret-key';
    const testDeviceId = 'device-123';

    late LocationSigner signer;

    setUp(() {
      signer = LocationSigner(
        enrollmentToken: testToken,
        deviceId: testDeviceId,
      );
    });

    group('constructor', () {
      test('throws ArgumentError for empty enrollment token', () {
        expect(
          () => LocationSigner(enrollmentToken: '', deviceId: testDeviceId),
          throwsA(isA<ArgumentError>()),
        );
      });

      test('throws ArgumentError for empty device ID', () {
        expect(
          () => LocationSigner(enrollmentToken: testToken, deviceId: ''),
          throwsA(isA<ArgumentError>()),
        );
      });

      test('creates successfully with valid parameters', () {
        final s = LocationSigner(
          enrollmentToken: testToken,
          deviceId: testDeviceId,
        );
        expect(s.deviceId, equals(testDeviceId));
      });
    });

    group('sign', () {
      test('produces a 64-character hex string (SHA-256 output)', () {
        final signature = signer.sign(
          timestamp: '2024-01-15T10:30:00.000Z',
          body: '{"latitude":33.31}',
        );
        expect(signature.length, equals(64));
        expect(RegExp(r'^[0-9a-f]{64}$').hasMatch(signature), isTrue);
      });

      test('matches backend computeSignature algorithm', () {
        // Backend algorithm: HMAC-SHA256(deviceId + timestamp + body, token)
        const timestamp = '2024-01-15T10:30:00.000Z';
        const body = '[{"latitude":33.312456,"longitude":44.366789}]';

        final signature = signer.sign(timestamp: timestamp, body: body);

        // Manually compute expected value using the same algorithm
        final signedContent = '$testDeviceId$timestamp$body';
        final key = utf8.encode(testToken);
        final bytes = utf8.encode(signedContent);
        final hmac = Hmac(sha256, key);
        final expected = hmac.convert(bytes).toString();

        expect(signature, equals(expected));
      });

      test('produces different signatures for different bodies', () {
        const timestamp = '2024-01-15T10:30:00.000Z';

        final sig1 = signer.sign(timestamp: timestamp, body: '{"a":1}');
        final sig2 = signer.sign(timestamp: timestamp, body: '{"a":2}');

        expect(sig1, isNot(equals(sig2)));
      });

      test('produces different signatures for different timestamps', () {
        const body = '{"latitude":33.31}';

        final sig1 = signer.sign(
          timestamp: '2024-01-15T10:30:00.000Z',
          body: body,
        );
        final sig2 = signer.sign(
          timestamp: '2024-01-15T10:31:00.000Z',
          body: body,
        );

        expect(sig1, isNot(equals(sig2)));
      });

      test('produces different signatures for different tokens', () {
        const timestamp = '2024-01-15T10:30:00.000Z';
        const body = '{"latitude":33.31}';

        final signer2 = LocationSigner(
          enrollmentToken: 'different-token',
          deviceId: testDeviceId,
        );

        final sig1 = signer.sign(timestamp: timestamp, body: body);
        final sig2 = signer2.sign(timestamp: timestamp, body: body);

        expect(sig1, isNot(equals(sig2)));
      });

      test('is deterministic — same inputs produce same output', () {
        const timestamp = '2024-01-15T10:30:00.000Z';
        const body = '{"latitude":33.31}';

        final sig1 = signer.sign(timestamp: timestamp, body: body);
        final sig2 = signer.sign(timestamp: timestamp, body: body);

        expect(sig1, equals(sig2));
      });
    });

    group('generateHeaders', () {
      test('returns all three required headers', () {
        final headers = signer.generateHeaders(
          timestamp: '2024-01-15T10:30:00.000Z',
          body: '{"test":true}',
        );

        expect(headers.containsKey('X-Device-Signature'), isTrue);
        expect(headers.containsKey('X-Device-Id'), isTrue);
        expect(headers.containsKey('X-Request-Timestamp'), isTrue);
      });

      test('X-Device-Id matches the configured device ID', () {
        final headers = signer.generateHeaders(
          timestamp: '2024-01-15T10:30:00.000Z',
          body: '{"test":true}',
        );

        expect(headers['X-Device-Id'], equals(testDeviceId));
      });

      test('X-Request-Timestamp matches the provided timestamp', () {
        const timestamp = '2024-01-15T10:30:00.000Z';
        final headers = signer.generateHeaders(
          timestamp: timestamp,
          body: '{"test":true}',
        );

        expect(headers['X-Request-Timestamp'], equals(timestamp));
      });

      test('X-Device-Signature matches sign() output', () {
        const timestamp = '2024-01-15T10:30:00.000Z';
        const body = '{"test":true}';

        final headers = signer.generateHeaders(
          timestamp: timestamp,
          body: body,
        );
        final expectedSig = signer.sign(timestamp: timestamp, body: body);

        expect(headers['X-Device-Signature'], equals(expectedSig));
      });
    });

    group('generateTimestamp', () {
      test('returns a valid ISO 8601 UTC timestamp', () {
        final timestamp = signer.generateTimestamp();

        // Should end with Z (UTC)
        expect(timestamp.endsWith('Z'), isTrue);

        // Should be parseable
        final parsed = DateTime.tryParse(timestamp);
        expect(parsed, isNotNull);
        expect(parsed!.isUtc, isTrue);
      });

      test('returns a timestamp close to now', () {
        final before = DateTime.now().toUtc();
        final timestamp = signer.generateTimestamp();
        final after = DateTime.now().toUtc();

        final parsed = DateTime.parse(timestamp);
        expect(parsed.isAfter(before.subtract(const Duration(seconds: 1))),
            isTrue);
        expect(
            parsed.isBefore(after.add(const Duration(seconds: 1))), isTrue);
      });
    });

    group('isTlsSecure', () {
      test('returns true for HTTPS URLs', () {
        expect(LocationSigner.isTlsSecure('https://api.example.com'), isTrue);
        expect(
          LocationSigner.isTlsSecure('https://api.example.com/locations'),
          isTrue,
        );
        expect(LocationSigner.isTlsSecure('HTTPS://API.EXAMPLE.COM'), isTrue);
      });

      test('returns false for HTTP URLs', () {
        expect(LocationSigner.isTlsSecure('http://api.example.com'), isFalse);
        expect(
          LocationSigner.isTlsSecure('http://api.example.com/locations'),
          isFalse,
        );
      });

      test('returns false for non-HTTP schemes', () {
        expect(LocationSigner.isTlsSecure('ftp://files.example.com'), isFalse);
        expect(LocationSigner.isTlsSecure('ws://socket.example.com'), isFalse);
      });

      test('returns false for invalid URLs', () {
        expect(LocationSigner.isTlsSecure(''), isFalse);
        expect(LocationSigner.isTlsSecure('not a url'), isFalse);
      });

      test('returns true for wss:// (WebSocket Secure) — only https', () {
        // wss is not https, so it should return false for strict TLS check
        expect(LocationSigner.isTlsSecure('wss://socket.example.com'), isFalse);
      });
    });
  });
}
