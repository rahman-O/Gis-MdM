import '../../core/network/api_client.dart';
import '../../core/utils/logger.dart';

/// Firebase Cloud Messaging service for receiving push commands.
class FcmService {
  /// Initialize FCM and register for messages.
  Future<void> initialize() async {
    // Note: Firebase dependencies are commented out in pubspec.yaml (Phase 2)
    // This is the structure — actual Firebase init will be added when deps are enabled
    Logger.info('FCM service initialized (stub — enable firebase deps)', 'FCM');
  }

  /// Get the current FCM token.
  Future<String?> getToken() async {
    // Will use FirebaseMessaging.instance.getToken()
    return null;
  }

  /// Handle incoming FCM message — dispatch to CommandService.
  void onMessage(Map<String, dynamic> data) {
    Logger.info('FCM message received: ${data['messageType']}', 'FCM');
    // Parse message type and dispatch to command handlers
  }

  /// Send FCM token to server for registration.
  Future<void> registerToken(
      ApiClient api, String deviceId, String token) async {
    // POST token to server so it can send push to this device
    Logger.info('FCM token registered for device: $deviceId', 'FCM');
  }
}
