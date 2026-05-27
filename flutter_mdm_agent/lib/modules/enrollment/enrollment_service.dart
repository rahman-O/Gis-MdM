import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/storage/secure_storage.dart';
import '../../core/utils/logger.dart';

/// Handles device enrollment with the MDM server.
///
/// Enrollment is the first step in the device lifecycle — it registers
/// the device with the server using a configuration key and device ID.
class EnrollmentService {
  final ApiClient _api;

  EnrollmentService(this._api);

  /// Enroll the device with the MDM server.
  ///
  /// [serverUrl] — base URL of the MDM server (e.g. https://mdm.example.com)
  /// [configKey] — configuration key provided by the admin
  /// [deviceId] — unique device identifier
  ///
  /// Returns `true` if enrollment succeeded, `false` otherwise.
  Future<bool> enrollWithToken(
    String serverUrl,
    String configKey,
    String deviceId,
  ) async {
    try {
      Logger.info(
        'Starting enrollment: server=$serverUrl, deviceId=$deviceId',
        'Enrollment',
      );

      // Configure the API client with the server URL
      _api.configure(baseUrl: serverUrl);

      // POST enrollment request
      final response = await _api.post(
        Endpoints.syncConfiguration(deviceId),
        data: {'configuration': configKey},
      );

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Persist enrollment data securely
        await SecureStorage.setServerUrl(serverUrl);
        await SecureStorage.setDeviceId(deviceId);
        await SecureStorage.setEnrolledAt(
          DateTime.now().toIso8601String(),
        );

        Logger.info('Enrollment successful', 'Enrollment');
        return true;
      }

      Logger.warn(
        'Enrollment failed: status=${response.statusCode}',
        'Enrollment',
      );
      return false;
    } catch (e, stack) {
      Logger.error('Enrollment error', e, stack, 'Enrollment');
      return false;
    }
  }

  /// Check if the device is currently enrolled.
  Future<bool> isEnrolled() => SecureStorage.isEnrolled();
}
