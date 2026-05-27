import 'dart:io';

import 'package:path_provider/path_provider.dart';

import '../../../core/network/api_client.dart';
import '../../../core/utils/logger.dart';
import '../../../platform/device_owner_channel.dart';
import '../command_queue.dart';
import '../models/command.dart';

/// Handles "appInstall" commands by downloading and installing APKs.
///
/// Downloads the APK from the URL in the command payload, saves it
/// to a temporary directory, and uses Device Owner APIs to install silently.
class AppInstallHandler implements CommandHandler {
  final ApiClient _api;

  AppInstallHandler(this._api);

  @override
  bool canHandle(String messageType) =>
      messageType == 'appInstall' ||
      messageType == 'APP_INSTALL' ||
      messageType == 'installApp';

  @override
  Future<bool> handle(RemoteCommand command) async {
    Logger.info('Executing app install command', 'AppInstallHandler');

    try {
      // Parse payload — expects JSON with "url" and optionally "pkg"
      final url = _extractUrl(command.payload);
      if (url == null || url.isEmpty) {
        Logger.warn('App install: no URL in payload', 'AppInstallHandler');
        return false;
      }

      // Download APK to temp directory
      final dir = await getTemporaryDirectory();
      final fileName = 'install_${command.id}.apk';
      final savePath = '${dir.path}/$fileName';

      Logger.info('Downloading APK: $url', 'AppInstallHandler');
      await _api.download(url, savePath);

      // Verify file exists
      final file = File(savePath);
      if (!await file.exists()) {
        Logger.warn('APK download failed: file not found', 'AppInstallHandler');
        return false;
      }

      // Install silently via Device Owner
      Logger.info('Installing APK: $savePath', 'AppInstallHandler');
      final success = await DeviceOwnerChannel.installPackage(savePath);

      // Clean up temp file
      try {
        await file.delete();
      } catch (_) {
        // Non-critical cleanup failure
      }

      if (success) {
        Logger.info('App installed successfully', 'AppInstallHandler');
      } else {
        Logger.warn('App installation returned false', 'AppInstallHandler');
      }

      return success;
    } catch (e, stack) {
      Logger.error('App install failed', e, stack, 'AppInstallHandler');
      return false;
    }
  }

  /// Extract the download URL from the command payload.
  ///
  /// Payload can be a plain URL string or a JSON string with a "url" field.
  String? _extractUrl(String payload) {
    final trimmed = payload.trim();

    // If it looks like a URL directly
    if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
      return trimmed;
    }

    // Try to parse as simple key=value or JSON-like
    // Simple approach: look for url field
    final urlMatch = RegExp(r'"url"\s*:\s*"([^"]+)"').firstMatch(trimmed);
    if (urlMatch != null) {
      return urlMatch.group(1);
    }

    return null;
  }
}
