import '../../../core/utils/logger.dart';
import '../../../platform/device_owner_channel.dart';
import '../command_queue.dart';
import '../models/command.dart';

/// Handles "wipe" commands by performing a factory reset.
///
/// WARNING: This is a destructive, irreversible operation.
/// The device will be completely wiped and reset to factory state.
class WipeHandler implements CommandHandler {
  @override
  bool canHandle(String messageType) =>
      messageType == 'wipe' || messageType == 'WIPE_DATA';

  @override
  Future<bool> handle(RemoteCommand command) async {
    Logger.info('Executing WIPE command — factory reset', 'WipeHandler');

    try {
      // This is irreversible — log before executing
      Logger.warn(
        'Factory reset initiated by server command id=${command.id}',
        'WipeHandler',
      );

      await DeviceOwnerChannel.wipeData();

      // If we reach here, wipe may not have executed immediately
      return true;
    } catch (e, stack) {
      Logger.error('Wipe command failed', e, stack, 'WipeHandler');
      return false;
    }
  }
}
