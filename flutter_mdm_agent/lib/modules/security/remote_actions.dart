import '../../core/utils/logger.dart';
import '../../platform/device_owner_channel.dart';

/// High-level remote security actions.
class RemoteActions {
  /// Lock device immediately.
  static Future<void> lockDevice() async {
    await DeviceOwnerChannel.lockNow();
    Logger.info('Device locked remotely', 'RemoteActions');
  }

  /// Factory reset (wipe all data).
  static Future<void> wipeDevice() async {
    Logger.warn('WIPE COMMAND RECEIVED — wiping device', 'RemoteActions');
    await DeviceOwnerChannel.wipeData();
  }

  /// Reboot device.
  static Future<void> rebootDevice() async {
    await DeviceOwnerChannel.reboot();
    Logger.info('Device rebooted remotely', 'RemoteActions');
  }
}
