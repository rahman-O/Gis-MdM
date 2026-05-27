import 'package:flutter/services.dart';
import '../core/config/constants.dart';
import '../core/utils/logger.dart';

/// Method channel bridge to Kotlin Device Owner APIs.
class DeviceOwnerChannel {
  static const _channel = MethodChannel(AgentConstants.methodChannelName);

  /// Check if this app is Device Owner.
  static Future<bool> isDeviceOwner() async {
    try {
      return await _channel.invokeMethod<bool>('isDeviceOwner') ?? false;
    } catch (e) {
      Logger.error('isDeviceOwner failed', e, null, 'DeviceOwner');
      return false;
    }
  }

  /// Add a UserRestriction (e.g. "no_camera").
  static Future<void> addUserRestriction(String restriction) async {
    await _channel.invokeMethod('addUserRestriction', {'restriction': restriction});
  }

  /// Clear a UserRestriction.
  static Future<void> clearUserRestriction(String restriction) async {
    await _channel.invokeMethod('clearUserRestriction', {'restriction': restriction});
  }

  /// Grant a runtime permission to a package.
  static Future<void> grantPermission(String packageName, String permission) async {
    await _channel.invokeMethod('grantPermission', {
      'packageName': packageName,
      'permission': permission,
    });
  }

  /// Install APK silently (Device Owner only).
  static Future<bool> installPackage(String apkPath) async {
    try {
      return await _channel.invokeMethod<bool>('installPackage', {'apkPath': apkPath}) ?? false;
    } catch (e) {
      Logger.error('installPackage failed', e, null, 'DeviceOwner');
      return false;
    }
  }

  /// Uninstall package silently.
  static Future<bool> uninstallPackage(String packageName) async {
    try {
      return await _channel.invokeMethod<bool>('uninstallPackage', {'packageName': packageName}) ?? false;
    } catch (e) {
      Logger.error('uninstallPackage failed', e, null, 'DeviceOwner');
      return false;
    }
  }

  /// Lock device screen immediately.
  static Future<void> lockNow() async {
    await _channel.invokeMethod('lockNow');
  }

  /// Factory reset (wipe data).
  static Future<void> wipeData() async {
    await _channel.invokeMethod('wipeData');
  }

  /// Reboot device.
  static Future<void> reboot() async {
    await _channel.invokeMethod('reboot');
  }

  /// Start Lock Task Mode (kiosk) for a package.
  static Future<void> startLockTask(String packageName) async {
    await _channel.invokeMethod('startLockTask', {'packageName': packageName});
  }

  /// Stop Lock Task Mode.
  static Future<void> stopLockTask() async {
    await _channel.invokeMethod('stopLockTask');
  }

  /// Set Lock Task packages (whitelist for kiosk).
  static Future<void> setLockTaskPackages(List<String> packages) async {
    await _channel.invokeMethod('setLockTaskPackages', {'packages': packages});
  }
}
