import '../../../core/utils/logger.dart';
import '../../../platform/device_owner_channel.dart';
import '../../sync/sync_response.dart';
import '../policy_engine.dart';

/// Applies kiosk mode policies via DeviceOwnerChannel.
///
/// Manages Lock Task Mode (Android kiosk) by starting/stopping
/// lock task based on the sync response configuration.
class KioskEnforcer implements PolicyEnforcer {
  bool _kioskActive = false;
  String? _currentKioskApp;

  @override
  String get name => 'Kiosk';

  @override
  Future<void> enforce(SyncResponse response) async {
    final shouldBeActive = response.kioskMode;
    final desiredApp = response.kioskApp;

    if (shouldBeActive && desiredApp != null && desiredApp.isNotEmpty) {
      // Enable or update kiosk mode
      if (!_kioskActive || _currentKioskApp != desiredApp) {
        await _enableKiosk(desiredApp, response);
      }
    } else if (_kioskActive) {
      // Disable kiosk mode
      await _disableKiosk();
    }
  }

  @override
  Future<void> clear() async {
    if (_kioskActive) {
      await _disableKiosk();
    }
  }

  /// Whether kiosk mode is currently active.
  bool get isActive => _kioskActive;

  /// The package currently in kiosk mode.
  String? get currentKioskApp => _currentKioskApp;

  /// Enable kiosk mode for the specified package.
  Future<void> _enableKiosk(String packageName, SyncResponse response) async {
    try {
      // If already in kiosk with a different app, stop first
      if (_kioskActive) {
        await DeviceOwnerChannel.stopLockTask();
      }

      // Set allowed packages for lock task
      await DeviceOwnerChannel.setLockTaskPackages([packageName]);

      // Apply lock status bar setting if specified
      if (response.lockStatusBar == true) {
        Logger.debug('Status bar will be locked in kiosk mode', 'KioskEnforcer');
      }

      // Start lock task mode
      await DeviceOwnerChannel.startLockTask(packageName);

      _kioskActive = true;
      _currentKioskApp = packageName;
      Logger.info('Kiosk mode enabled for: $packageName', 'KioskEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to enable kiosk mode', e, stack, 'KioskEnforcer');
    }
  }

  /// Disable kiosk mode.
  Future<void> _disableKiosk() async {
    try {
      await DeviceOwnerChannel.stopLockTask();
      _kioskActive = false;
      _currentKioskApp = null;
      Logger.info('Kiosk mode disabled', 'KioskEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to disable kiosk mode', e, stack, 'KioskEnforcer');
    }
  }
}
