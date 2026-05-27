import '../../../core/utils/logger.dart';
import '../../sync/sync_service.dart';
import '../command_queue.dart';
import '../models/command.dart';

/// Handles "configUpdated" commands by triggering a full sync.
///
/// When the server pushes a config update notification, this handler
/// fetches the latest configuration from the server.
class ConfigUpdatedHandler implements CommandHandler {
  final SyncService _syncService;
  final String _deviceId;

  ConfigUpdatedHandler(this._syncService, this._deviceId);

  @override
  bool canHandle(String messageType) =>
      messageType == 'configUpdated' || messageType == 'CONFIG_UPDATED';

  @override
  Future<bool> handle(RemoteCommand command) async {
    Logger.info('Config update received, triggering sync', 'ConfigUpdatedHandler');

    try {
      final config = await _syncService.fetchConfiguration(_deviceId);
      if (config != null) {
        Logger.info(
          'Config sync successful: ${config.applications.length} apps',
          'ConfigUpdatedHandler',
        );
        return true;
      }

      Logger.warn('Config sync returned null', 'ConfigUpdatedHandler');
      return false;
    } catch (e, stack) {
      Logger.error('Config update handler failed', e, stack, 'ConfigUpdatedHandler');
      return false;
    }
  }
}
