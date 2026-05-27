import '../../core/network/api_client.dart';
import '../../core/utils/logger.dart';
import '../../platform/device_owner_channel.dart';
import '../sync/sync_response.dart';

/// Manages self-update of the MDM Agent app.
class OtaUpdater {
  final ApiClient _api;
  final String _currentVersion;

  OtaUpdater(this._api, this._currentVersion);

  /// Check if an update is available from the sync response.
  bool isUpdateAvailable(List<SyncApplication> apps) {
    final self =
        apps.where((a) => a.pkg == 'com.gismdm.mdm_agent').firstOrNull;
    if (self == null) return false;
    return self.version != _currentVersion;
  }

  /// Download and install the update.
  Future<bool> performUpdate(SyncApplication app) async {
    Logger.info('OTA update available: ${app.version}', 'OtaUpdater');
    // 1. Download APK
    // 2. Verify checksum (if provided)
    // 3. Install via DeviceOwnerChannel
    const tempPath = '/data/local/tmp/mdm_agent_update.apk';
    try {
      await _api.download(app.url, tempPath);
      final success = await DeviceOwnerChannel.installPackage(tempPath);
      if (success) {
        Logger.info('OTA update installed: ${app.version}', 'OtaUpdater');
      }
      return success;
    } catch (e) {
      Logger.error('OTA update failed', e, null, 'OtaUpdater');
      return false;
    }
  }
}
