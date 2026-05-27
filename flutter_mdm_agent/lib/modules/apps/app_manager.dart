import 'dart:io';

import '../../core/network/api_client.dart';
import '../../core/utils/logger.dart';
import '../../platform/device_owner_channel.dart';
import '../sync/sync_response.dart';
import 'app_inventory.dart';

/// Manages application installation, uninstallation, and updates.
///
/// Uses [ApiClient] to download APKs and [DeviceOwnerChannel] to
/// perform silent install/uninstall operations as Device Owner.
class AppManager {
  final ApiClient _api;
  final AppInventory _inventory;

  AppManager({
    required ApiClient api,
    AppInventory? inventory,
  })  : _api = api,
        _inventory = inventory ?? AppInventory();

  /// Install APK from URL silently.
  ///
  /// Downloads the APK to a temporary directory, installs it via
  /// DeviceOwnerChannel, and cleans up the temp file afterwards.
  /// Returns `true` if installation succeeded.
  Future<bool> installFromUrl(String url, String packageName) async {
    Logger.info('Installing $packageName from: $url', 'AppManager');

    final tempDir = Directory.systemTemp;
    final tempFile = File('${tempDir.path}/$packageName.apk');

    try {
      // 1. Download APK to temp directory
      await _api.download(url, tempFile.path);
      Logger.debug('APK downloaded to: ${tempFile.path}', 'AppManager');

      // 2. Verify file exists and has content
      if (!await tempFile.exists() || await tempFile.length() == 0) {
        Logger.error('Downloaded APK is empty or missing', null, null, 'AppManager');
        return false;
      }

      // 3. Call DeviceOwnerChannel to install
      final success = await DeviceOwnerChannel.installPackage(tempFile.path);

      if (success) {
        Logger.info('Successfully installed: $packageName', 'AppManager');
      } else {
        Logger.error('Installation failed for: $packageName', null, null, 'AppManager');
      }

      return success;
    } catch (e, stack) {
      Logger.error('Failed to install $packageName', e, stack, 'AppManager');
      return false;
    } finally {
      // 4. Clean up temp file
      try {
        if (await tempFile.exists()) {
          await tempFile.delete();
          Logger.debug('Temp APK cleaned up: ${tempFile.path}', 'AppManager');
        }
      } catch (e) {
        Logger.warn('Failed to clean up temp file: ${tempFile.path}', 'AppManager');
      }
    }
  }

  /// Uninstall package silently.
  ///
  /// Returns `true` if uninstallation succeeded.
  Future<bool> uninstall(String packageName) async {
    Logger.info('Uninstalling: $packageName', 'AppManager');

    try {
      final success = await DeviceOwnerChannel.uninstallPackage(packageName);

      if (success) {
        Logger.info('Successfully uninstalled: $packageName', 'AppManager');
      } else {
        Logger.error('Uninstall failed for: $packageName', null, null, 'AppManager');
      }

      return success;
    } catch (e, stack) {
      Logger.error('Failed to uninstall $packageName', e, stack, 'AppManager');
      return false;
    }
  }

  /// Check and install/update apps from SyncResponse.
  ///
  /// Compares the desired application list against currently installed
  /// apps and performs install, update, or removal as needed.
  Future<void> syncApps(List<SyncApplication> desired) async {
    Logger.info('Syncing ${desired.length} applications', 'AppManager');

    final installed = await _inventory.getInstalled();
    final installedMap = {for (final app in installed) app.packageName: app};

    for (final app in desired) {
      if (app.pkg.isEmpty || app.url.isEmpty) {
        Logger.warn('Skipping app with missing pkg or url: ${app.name}', 'AppManager');
        continue;
      }

      final existing = installedMap[app.pkg];

      if (existing == null) {
        // App not installed — install it
        Logger.debug('App not installed, installing: ${app.pkg}', 'AppManager');
        await installFromUrl(app.url, app.pkg);
      } else if (_needsUpdate(existing, app)) {
        // App installed but outdated — update it
        Logger.debug('App needs update: ${app.pkg} (${existing.version} -> ${app.version})', 'AppManager');
        await installFromUrl(app.url, app.pkg);
      } else {
        Logger.debug('App up to date: ${app.pkg}', 'AppManager');
      }
    }

    Logger.info('App sync complete', 'AppManager');
  }

  /// Determine if an installed app needs to be updated.
  bool _needsUpdate(InstalledApp installed, SyncApplication desired) {
    if (desired.version.isEmpty) return false;
    return installed.version != desired.version;
  }
}
