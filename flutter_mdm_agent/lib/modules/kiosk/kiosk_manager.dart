import 'package:flutter/services.dart';

import '../../core/utils/logger.dart';
import '../../platform/device_owner_channel.dart';

/// Manages kiosk mode (Lock Task Mode) on the device.
///
/// Provides a high-level API to enable/disable kiosk mode with
/// configurable options for UI elements and exit behavior.
class KioskManager {
  static const _channel = MethodChannel('com.mdm.agent/kiosk');

  bool _active = false;
  String? _currentPackage;
  KioskOptions? _currentOptions;

  /// Enable kiosk mode for a specific package.
  ///
  /// Sets the package as the only allowed lock task package,
  /// applies [KioskOptions] for UI customization, and starts
  /// lock task mode.
  Future<void> enable(String packageName, KioskOptions options) async {
    Logger.info('Enabling kiosk mode for: $packageName', 'KioskManager');

    try {
      // Set allowed lock task packages
      await DeviceOwnerChannel.setLockTaskPackages([packageName]);

      // Apply kiosk UI options via platform channel
      await _applyKioskOptions(options);

      // Start lock task mode
      await DeviceOwnerChannel.startLockTask(packageName);

      _active = true;
      _currentPackage = packageName;
      _currentOptions = options;

      Logger.info('Kiosk mode enabled for: $packageName', 'KioskManager');
    } catch (e, stack) {
      Logger.error('Failed to enable kiosk mode', e, stack, 'KioskManager');
      rethrow;
    }
  }

  /// Disable kiosk mode.
  ///
  /// Stops lock task mode and restores default UI behavior.
  Future<void> disable() async {
    Logger.info('Disabling kiosk mode', 'KioskManager');

    try {
      await DeviceOwnerChannel.stopLockTask();

      // Reset kiosk UI options
      await _resetKioskOptions();

      _active = false;
      _currentPackage = null;
      _currentOptions = null;

      Logger.info('Kiosk mode disabled', 'KioskManager');
    } catch (e, stack) {
      Logger.error('Failed to disable kiosk mode', e, stack, 'KioskManager');
      rethrow;
    }
  }

  /// Whether kiosk mode is currently active.
  bool get isActive => _active;

  /// The package currently running in kiosk mode.
  String? get currentPackage => _currentPackage;

  /// The current kiosk options.
  KioskOptions? get currentOptions => _currentOptions;

  /// Apply kiosk UI options via platform channel.
  Future<void> _applyKioskOptions(KioskOptions options) async {
    try {
      await _channel.invokeMethod('applyKioskOptions', {
        'showHome': options.showHome,
        'showRecents': options.showRecents,
        'showNotifications': options.showNotifications,
        'showStatusBar': options.showStatusBar,
        'lockButtons': options.lockButtons,
        'keepScreenOn': options.keepScreenOn,
        'exitMethod': options.exitMethod,
      });
      Logger.debug('Kiosk options applied', 'KioskManager');
    } catch (e, stack) {
      Logger.error('Failed to apply kiosk options', e, stack, 'KioskManager');
    }
  }

  /// Reset kiosk UI options to defaults.
  Future<void> _resetKioskOptions() async {
    try {
      await _channel.invokeMethod('resetKioskOptions');
      Logger.debug('Kiosk options reset', 'KioskManager');
    } catch (e, stack) {
      Logger.error('Failed to reset kiosk options', e, stack, 'KioskManager');
    }
  }
}

/// Configuration options for kiosk mode.
///
/// Controls which UI elements are visible and how the user
/// can interact with the device while in kiosk mode.
class KioskOptions {
  /// Whether the home button is functional.
  final bool showHome;

  /// Whether the recents/overview button is functional.
  final bool showRecents;

  /// Whether notifications are shown.
  final bool showNotifications;

  /// Whether the status bar is visible.
  final bool showStatusBar;

  /// Whether hardware buttons are locked.
  final bool lockButtons;

  /// Whether the screen stays on indefinitely.
  final bool keepScreenOn;

  /// Exit method for kiosk mode: "password", "back", or "none".
  final String exitMethod;

  const KioskOptions({
    this.showHome = false,
    this.showRecents = false,
    this.showNotifications = false,
    this.showStatusBar = false,
    this.lockButtons = true,
    this.keepScreenOn = true,
    this.exitMethod = 'none',
  });

  /// Create KioskOptions from a JSON map.
  factory KioskOptions.fromJson(Map<String, dynamic> json) {
    return KioskOptions(
      showHome: json['showHome'] as bool? ?? false,
      showRecents: json['showRecents'] as bool? ?? false,
      showNotifications: json['showNotifications'] as bool? ?? false,
      showStatusBar: json['showStatusBar'] as bool? ?? false,
      lockButtons: json['lockButtons'] as bool? ?? true,
      keepScreenOn: json['keepScreenOn'] as bool? ?? true,
      exitMethod: json['exitMethod'] as String? ?? 'none',
    );
  }

  /// Convert to JSON map.
  Map<String, dynamic> toJson() {
    return {
      'showHome': showHome,
      'showRecents': showRecents,
      'showNotifications': showNotifications,
      'showStatusBar': showStatusBar,
      'lockButtons': lockButtons,
      'keepScreenOn': keepScreenOn,
      'exitMethod': exitMethod,
    };
  }

  @override
  String toString() =>
      'KioskOptions(home=$showHome, recents=$showRecents, notifications=$showNotifications, '
      'statusBar=$showStatusBar, lockButtons=$lockButtons, keepScreenOn=$keepScreenOn, '
      'exit=$exitMethod)';
}
