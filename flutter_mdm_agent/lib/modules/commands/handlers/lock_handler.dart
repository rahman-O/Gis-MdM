import '../../../core/utils/logger.dart';
import '../../../platform/device_owner_channel.dart';
import '../command_queue.dart';
import '../models/command.dart';

/// Handles "lock" commands by locking the device screen immediately.
class LockHandler implements CommandHandler {
  @override
  bool canHandle(String messageType) =>
      messageType == 'lock' || messageType == 'LOCK_DEVICE';

  @override
  Future<bool> handle(RemoteCommand command) async {
    Logger.info('Executing device lock command', 'LockHandler');

    try {
      await DeviceOwnerChannel.lockNow();
      Logger.info('Device locked successfully', 'LockHandler');
      return true;
    } catch (e, stack) {
      Logger.error('Device lock failed', e, stack, 'LockHandler');
      return false;
    }
  }
}
