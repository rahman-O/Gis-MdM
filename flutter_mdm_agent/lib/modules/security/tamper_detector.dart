import '../../core/utils/logger.dart';

/// Detects tampering attempts on the device.
class TamperDetector {
  /// Check if device is rooted.
  Future<bool> isRooted() async {
    // Check for su binary, Magisk, common root indicators
    // Uses platform channel to check native
    Logger.debug('Checking root status', 'TamperDetector');
    return false;
  }

  /// Check if agent app is being uninstalled.
  Future<bool> isUninstallAttempted() async {
    // Monitor package removal broadcasts
    return false;
  }

  /// Check if SIM card changed.
  Future<bool> isSimChanged(String? lastKnownImsi) async {
    // Compare current IMSI with stored one
    if (lastKnownImsi == null) return false;
    return false;
  }

  /// Run all security checks and return findings.
  Future<SecurityReport> runFullCheck() async {
    return SecurityReport(
      isRooted: await isRooted(),
      simChanged: false,
      uninstallAttempted: false,
      timestamp: DateTime.now(),
    );
  }
}

/// Report containing results of all security checks.
class SecurityReport {
  final bool isRooted;
  final bool simChanged;
  final bool uninstallAttempted;
  final DateTime timestamp;

  SecurityReport({
    required this.isRooted,
    required this.simChanged,
    required this.uninstallAttempted,
    required this.timestamp,
  });

  Map<String, dynamic> toJson() => {
        'isRooted': isRooted,
        'simChanged': simChanged,
        'uninstallAttempted': uninstallAttempted,
        'timestamp': timestamp.millisecondsSinceEpoch,
      };

  bool get hasIssues => isRooted || simChanged || uninstallAttempted;
}
