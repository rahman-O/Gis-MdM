import 'dart:io';

import 'package:path_provider/path_provider.dart';

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects disk storage information (free/total space).
class StorageCollector {
  /// Collect current storage information.
  Future<StorageInfo> collect() async {
    try {
      final dir = await getExternalStorageDirectory() ??
          await getApplicationDocumentsDirectory();

      final stat = await _getStorageStat(dir.path);
      return stat;
    } catch (e) {
      Logger.warn('Storage collection failed: $e', 'StorageCollector');
      return StorageInfo(totalBytes: 0, freeBytes: 0);
    }
  }

  Future<StorageInfo> _getStorageStat(String path) async {
    try {
      // Use the statfs-equivalent via ProcessResult on Android
      final result = await Process.run('df', [path]);
      if (result.exitCode == 0) {
        final lines = (result.stdout as String).trim().split('\n');
        if (lines.length >= 2) {
          final parts = lines[1].split(RegExp(r'\s+'));
          if (parts.length >= 4) {
            // df outputs in 1K blocks
            final total = int.tryParse(parts[1]) ?? 0;
            final free = int.tryParse(parts[3]) ?? 0;
            return StorageInfo(
              totalBytes: total * 1024,
              freeBytes: free * 1024,
            );
          }
        }
      }
    } catch (_) {
      // Fallback: cannot determine storage
    }
    return StorageInfo(totalBytes: 0, freeBytes: 0);
  }
}
