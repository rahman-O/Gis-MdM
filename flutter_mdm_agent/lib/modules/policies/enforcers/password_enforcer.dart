import 'package:flutter/services.dart';

import '../../../core/utils/logger.dart';
import '../../sync/sync_response.dart';
import '../policy_engine.dart';

/// Applies password quality policies via DeviceOwnerChannel.
///
/// Reads the [SyncResponse.passwordMode] field and sets the
/// appropriate password quality level on the device.
class PasswordEnforcer implements PolicyEnforcer {
  static const _channel = MethodChannel('com.mdm.agent/password');

  String? _currentMode;

  @override
  String get name => 'Password';

  @override
  Future<void> enforce(SyncResponse response) async {
    final desiredMode = response.passwordMode;

    if (desiredMode == null || desiredMode.isEmpty) {
      // No password policy specified, skip
      return;
    }

    if (desiredMode == _currentMode) {
      // Already applied, no change needed
      Logger.debug('Password mode unchanged: $desiredMode', 'PasswordEnforcer');
      return;
    }

    try {
      final quality = _mapModeToQuality(desiredMode);
      await _channel.invokeMethod('setPasswordQuality', {
        'quality': quality,
        'mode': desiredMode,
      });

      _currentMode = desiredMode;
      Logger.info('Password quality set to: $desiredMode (quality=$quality)', 'PasswordEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to set password quality', e, stack, 'PasswordEnforcer');
    }
  }

  @override
  Future<void> clear() async {
    try {
      await _channel.invokeMethod('clearPasswordPolicy');
      _currentMode = null;
      Logger.info('Password policy cleared', 'PasswordEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to clear password policy', e, stack, 'PasswordEnforcer');
    }
  }

  /// Get the currently applied password mode.
  String? get currentMode => _currentMode;

  /// Map password mode string to Android password quality constant.
  ///
  /// Android DevicePolicyManager quality constants:
  /// - PASSWORD_QUALITY_UNSPECIFIED = 0
  /// - PASSWORD_QUALITY_SOMETHING = 65536
  /// - PASSWORD_QUALITY_NUMERIC = 131072
  /// - PASSWORD_QUALITY_NUMERIC_COMPLEX = 196608
  /// - PASSWORD_QUALITY_ALPHABETIC = 262144
  /// - PASSWORD_QUALITY_ALPHANUMERIC = 327680
  /// - PASSWORD_QUALITY_COMPLEX = 393216
  int _mapModeToQuality(String mode) {
    switch (mode.toLowerCase()) {
      case 'none':
        return 0; // PASSWORD_QUALITY_UNSPECIFIED
      case 'something':
      case 'pattern':
        return 65536; // PASSWORD_QUALITY_SOMETHING
      case 'numeric':
      case 'pin':
        return 131072; // PASSWORD_QUALITY_NUMERIC
      case 'numeric_complex':
        return 196608; // PASSWORD_QUALITY_NUMERIC_COMPLEX
      case 'alphabetic':
        return 262144; // PASSWORD_QUALITY_ALPHABETIC
      case 'alphanumeric':
        return 327680; // PASSWORD_QUALITY_ALPHANUMERIC
      case 'complex':
        return 393216; // PASSWORD_QUALITY_COMPLEX
      default:
        Logger.warn('Unknown password mode: $mode, defaulting to unspecified', 'PasswordEnforcer');
        return 0;
    }
  }
}
