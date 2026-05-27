import 'package:flutter/services.dart';

import '../../../core/utils/logger.dart';
import '../../sync/sync_response.dart';
import '../policy_engine.dart';

/// Applies hardware-related policies via platform channels.
///
/// Manages: GPS, Bluetooth, WiFi, mobile data, USB storage,
/// screenshots, brightness, volume, orientation, and screen timeout.
class HardwareEnforcer implements PolicyEnforcer {
  static const _channel = MethodChannel('com.mdm.agent/hardware');

  @override
  String get name => 'Hardware';

  @override
  Future<void> enforce(SyncResponse response) async {
    // GPS control
    if (response.gps != null) {
      await _setSetting('setGpsEnabled', response.gps!);
      Logger.debug('GPS set to: ${response.gps}', 'HardwareEnforcer');
    }

    // Bluetooth control
    if (response.bluetooth != null) {
      await _setSetting('setBluetoothEnabled', response.bluetooth!);
      Logger.debug('Bluetooth set to: ${response.bluetooth}', 'HardwareEnforcer');
    }

    // WiFi control
    if (response.wifi != null) {
      await _setSetting('setWifiEnabled', response.wifi!);
      Logger.debug('WiFi set to: ${response.wifi}', 'HardwareEnforcer');
    }

    // Mobile data control
    if (response.mobileData != null) {
      await _setSetting('setMobileDataEnabled', response.mobileData!);
      Logger.debug('Mobile data set to: ${response.mobileData}', 'HardwareEnforcer');
    }

    // USB storage control
    if (response.usbStorage != null) {
      await _setSetting('setUsbStorageEnabled', response.usbStorage!);
      Logger.debug('USB storage set to: ${response.usbStorage}', 'HardwareEnforcer');
    }

    // Screen brightness
    if (response.screenBrightness != null) {
      await _setIntSetting(
        'setScreenBrightness',
        response.screenBrightness!,
        autoBrightness: response.autoBrightness ?? false,
      );
      Logger.debug(
        'Brightness set to: ${response.screenBrightness} (auto: ${response.autoBrightness})',
        'HardwareEnforcer',
      );
    }

    // Volume control
    if (response.manageVolume == true && response.volumeLevel != null) {
      await _setIntSetting('setVolumeLevel', response.volumeLevel!);
      Logger.debug('Volume set to: ${response.volumeLevel}', 'HardwareEnforcer');
    }

    // Lock status bar
    if (response.lockStatusBar != null) {
      await _setSetting('setStatusBarLocked', response.lockStatusBar!);
      Logger.debug('Status bar locked: ${response.lockStatusBar}', 'HardwareEnforcer');
    }

    Logger.info('Hardware policies applied', 'HardwareEnforcer');
  }

  @override
  Future<void> clear() async {
    try {
      await _channel.invokeMethod('resetHardwareDefaults');
      Logger.info('Hardware policies reset to defaults', 'HardwareEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to reset hardware defaults', e, stack, 'HardwareEnforcer');
    }
  }

  /// Set a boolean hardware setting via platform channel.
  Future<void> _setSetting(String method, bool value) async {
    try {
      await _channel.invokeMethod(method, {'enabled': value});
    } catch (e, stack) {
      Logger.error('Failed to invoke $method', e, stack, 'HardwareEnforcer');
    }
  }

  /// Set an integer hardware setting via platform channel.
  Future<void> _setIntSetting(String method, int value, {bool? autoBrightness}) async {
    try {
      final args = <String, dynamic>{'value': value};
      if (autoBrightness != null) {
        args['autoBrightness'] = autoBrightness;
      }
      await _channel.invokeMethod(method, args);
    } catch (e, stack) {
      Logger.error('Failed to invoke $method', e, stack, 'HardwareEnforcer');
    }
  }
}
