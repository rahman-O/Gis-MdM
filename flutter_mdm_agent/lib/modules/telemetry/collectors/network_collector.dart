import 'dart:io';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:network_info_plus/network_info_plus.dart' as net_info;

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects network connectivity type, WiFi SSID, and local IP address.
class NetworkCollector {
  final Connectivity _connectivity = Connectivity();
  final net_info.NetworkInfo _networkInfo = net_info.NetworkInfo();

  /// Collect current network information including WiFi SSID and IP.
  Future<NetworkInfo> collect() async {
    try {
      final results = await _connectivity.checkConnectivity();
      final result =
          results.isNotEmpty ? results.first : ConnectivityResult.none;

      final networkType = _typeToString(result);
      final connected = result != ConnectivityResult.none;

      String? wifiSsid;
      String? ipAddress;

      // Collect WiFi SSID if connected to WiFi
      if (result == ConnectivityResult.wifi) {
        try {
          wifiSsid = await _networkInfo.getWifiName();
          // Remove surrounding quotes if present (Android returns "SSID")
          if (wifiSsid != null) {
            wifiSsid = wifiSsid.replaceAll('"', '');
            if (wifiSsid.isEmpty || wifiSsid == '<unknown ssid>') {
              wifiSsid = null;
            }
          }
        } catch (e) {
          Logger.warn('WiFi SSID collection failed: $e', 'NetworkCollector');
        }
      }

      // Collect local IP address
      try {
        ipAddress = await _getLocalIpAddress();
      } catch (e) {
        Logger.warn('IP address collection failed: $e', 'NetworkCollector');
      }

      return NetworkInfo(
        type: networkType,
        connected: connected,
        wifiSsid: wifiSsid,
        ipAddress: ipAddress,
      );
    } catch (e) {
      Logger.warn('Network collection failed: $e', 'NetworkCollector');
      return NetworkInfo(type: 'unknown', connected: false);
    }
  }

  /// Get the device's local IP address from network interfaces.
  Future<String?> _getLocalIpAddress() async {
    try {
      final interfaces = await NetworkInterface.list(
        type: InternetAddressType.IPv4,
        includeLinkLocal: false,
      );
      for (final iface in interfaces) {
        for (final addr in iface.addresses) {
          // Skip loopback addresses
          if (!addr.isLoopback) {
            return addr.address;
          }
        }
      }
    } catch (e) {
      Logger.warn('Failed to get local IP: $e', 'NetworkCollector');
    }
    return null;
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
