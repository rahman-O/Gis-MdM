import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Encrypted storage for sensitive data (enrollment token, server URL, device ID).
class SecureStorage {
  static const _storage = FlutterSecureStorage(
    aOptions: AndroidOptions(encryptedSharedPreferences: true),
  );

  static const _keyServerUrl = 'server_url';
  static const _keyDeviceId = 'device_id';
  static const _keyEnrollmentToken = 'enrollment_token';
  static const _keyEnrolledAt = 'enrolled_at';
  static const _keyFcmToken = 'fcm_token';

  // Server URL
  static Future<String?> getServerUrl() => _storage.read(key: _keyServerUrl);
  static Future<void> setServerUrl(String url) => _storage.write(key: _keyServerUrl, value: url);

  // Device ID
  static Future<String?> getDeviceId() => _storage.read(key: _keyDeviceId);
  static Future<void> setDeviceId(String id) => _storage.write(key: _keyDeviceId, value: id);

  // Enrollment Token
  static Future<String?> getEnrollmentToken() => _storage.read(key: _keyEnrollmentToken);
  static Future<void> setEnrollmentToken(String token) => _storage.write(key: _keyEnrollmentToken, value: token);

  // Enrolled At
  static Future<String?> getEnrolledAt() => _storage.read(key: _keyEnrolledAt);
  static Future<void> setEnrolledAt(String timestamp) => _storage.write(key: _keyEnrolledAt, value: timestamp);

  // FCM Token
  static Future<String?> getFcmToken() => _storage.read(key: _keyFcmToken);
  static Future<void> setFcmToken(String token) => _storage.write(key: _keyFcmToken, value: token);

  // Check if enrolled
  static Future<bool> isEnrolled() async {
    final url = await getServerUrl();
    final deviceId = await getDeviceId();
    return url != null && url.isNotEmpty && deviceId != null && deviceId.isNotEmpty;
  }

  // Clear all (for unenrollment)
  static Future<void> clearAll() => _storage.deleteAll();
}
