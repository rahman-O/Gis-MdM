import 'package:connectivity_plus/connectivity_plus.dart';

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects network connectivity type and status.
class NetworkCollector {
  final Connectivity _connectivity = Connectivity();

  /// Collect current network information.
  Future<NetworkInfo> collect() async {
    try {
      final results = await _connectivity.checkConnectivity();
      final result =
          results.isNotEmpty ? results.first : ConnectivityResult.none;

      return NetworkInfo(
        type: _typeToString(result),
        connected: result != ConnectivityResult.none,
      );
    } catch (e) {
      Logger.warn('Network collection failed: $e', 'NetworkCollector');
      return NetworkInfo(type: 'unknown', connected: false);
    }
  }

  String _typeToString(ConnectivityResult result) {
    switch (result) {
      case ConnectivityResult.wifi:
        return 'wifi';
      case ConnectivityResult.mobile:
        return 'mobile';
      case ConnectivityResult.ethernet:
        return 'ethernet';
      case ConnectivityResult.bluetooth:
        return 'bluetooth';
      case ConnectivityResult.vpn:
        return 'vpn';
      case ConnectivityResult.none:
        return 'none';
      default:
        return 'other';
    }
  }
}
