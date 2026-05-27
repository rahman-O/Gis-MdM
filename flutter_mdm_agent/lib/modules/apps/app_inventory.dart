import 'package:flutter/services.dart';

import '../../core/utils/logger.dart';

/// Queries the device for installed application information.
///
/// Uses a platform channel to retrieve the list of installed packages
/// and their version details from the Android system.
class AppInventory {
  static const _channel = MethodChannel('com.mdm.agent/app_inventory');

  /// Get list of installed packages with versions.
  ///
  /// Returns a list of [InstalledApp] objects representing all
  /// user-installed applications on the device.
  Future<List<InstalledApp>> getInstalled() async {
    try {
      final result = await _channel.invokeMethod<List<dynamic>>('getInstalledApps');

      if (result == null) {
        Logger.warn('getInstalledApps returned null', 'AppInventory');
        return [];
      }

      final apps = result.map((item) {
        final map = Map<String, dynamic>.from(item as Map);
        return InstalledApp(
          packageName: map['packageName'] as String? ?? '',
          version: map['version'] as String? ?? '',
          versionCode: map['versionCode'] as int? ?? 0,
        );
      }).toList();

      Logger.debug('Found ${apps.length} installed apps', 'AppInventory');
      return apps;
    } catch (e, stack) {
      Logger.error('Failed to get installed apps', e, stack, 'AppInventory');
      return [];
    }
  }

  /// Check if a specific package is installed.
  ///
  /// Returns `true` if the package with [packageName] is installed.
  Future<bool> isInstalled(String packageName) async {
    try {
      final result = await _channel.invokeMethod<bool>(
        'isPackageInstalled',
        {'packageName': packageName},
      );
      return result ?? false;
    } catch (e, stack) {
      Logger.error('Failed to check if $packageName is installed', e, stack, 'AppInventory');
      return false;
    }
  }

  /// Get version info for a specific installed package.
  ///
  /// Returns the [InstalledApp] if found, or `null` if not installed.
  Future<InstalledApp?> getAppInfo(String packageName) async {
    try {
      final result = await _channel.invokeMethod<Map<dynamic, dynamic>>(
        'getAppInfo',
        {'packageName': packageName},
      );

      if (result == null) return null;

      final map = Map<String, dynamic>.from(result);
      return InstalledApp(
        packageName: map['packageName'] as String? ?? '',
        version: map['version'] as String? ?? '',
        versionCode: map['versionCode'] as int? ?? 0,
      );
    } catch (e, stack) {
      Logger.error('Failed to get app info for $packageName', e, stack, 'AppInventory');
      return null;
    }
  }
}

/// Represents an installed application on the device.
class InstalledApp {
  /// The application package name (e.g., "com.example.app").
  final String packageName;

  /// The human-readable version string (e.g., "1.2.3").
  final String version;

  /// The numeric version code for comparison.
  final int versionCode;

  InstalledApp({
    required this.packageName,
    required this.version,
    required this.versionCode,
  });

  @override
  String toString() => 'InstalledApp($packageName v$version [$versionCode])';
}
