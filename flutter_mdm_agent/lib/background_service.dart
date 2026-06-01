import 'dart:async';
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_background_service/flutter_background_service.dart';

import 'core/config/constants.dart';
import 'core/network/api_client.dart';
import 'core/storage/local_db.dart';
import 'core/storage/secure_storage.dart';
import 'core/utils/logger.dart';
import 'modules/heartbeat/heartbeat_service.dart';
import 'modules/telemetry/telemetry_service.dart';
import 'modules/telemetry/location_sender.dart';

/// Initialize and configure the background service.
///
/// Call this once from main() before runApp().
Future<void> initializeBackgroundService() async {
  final service = FlutterBackgroundService();

  await service.configure(
    androidConfiguration: AndroidConfiguration(
      onStart: onStart,
      autoStart: true,
      isForegroundMode: true,
      autoStartOnBoot: true,
      notificationChannelId: AgentConstants.foregroundChannelId,
      initialNotificationTitle: 'MDM Agent',
      initialNotificationContent: 'Device management active',
      foregroundServiceNotificationId: 1001,
      foregroundServiceTypes: [
        AndroidForegroundType.location,
        AndroidForegroundType.dataSync,
      ],
    ),
    iosConfiguration: IosConfiguration(
      autoStart: false,
      onForeground: onStart,
      onBackground: onIosBackground,
    ),
  );
}

/// iOS background handler (no-op for this MDM agent, Android only).
@pragma('vm:entry-point')
Future<bool> onIosBackground(ServiceInstance service) async {
  return true;
}

/// Background service entry point.
///
/// This runs in a separate isolate from the UI. It handles:
/// - Heartbeat sending (every 15 seconds for demo)
/// - Telemetry collection (every 30 seconds)
/// - Location sending
/// - Self-restart via AlarmManager scheduling
@pragma('vm:entry-point')
Future<void> onStart(ServiceInstance service) async {
  // Ensure Flutter bindings are initialized in the background isolate
  DartPluginRegistrant.ensureInitialized();
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize local storage
  await LocalDb.initialize();

  Logger.info('Background service started', 'BackgroundService');

  // Device ID (hardcoded for now — in production from secure storage)
  const deviceId = '351906200367061';

  // Set up API client
  final api = ApiClient();
  final serverUrl = await SecureStorage.getServerUrl();
  if (serverUrl != null && serverUrl.isNotEmpty) {
    api.configure(baseUrl: serverUrl);
  } else {
    api.configure(baseUrl: 'https://mdm.studhub.app');
  }

  // Initialize services
  final heartbeat = HeartbeatService(api, deviceId);
  final telemetry = TelemetryService();
  final locationSender = LocationSender();

  int heartbeatCount = 0;
  int telemetryCount = 0;

  // Listen for stop command from UI
  service.on('stopService').listen((event) {
    service.stopSelf();
    Logger.info('Background service stopped by user', 'BackgroundService');
  });

  // Listen for status request from UI
  service.on('status').listen((event) {
    service.invoke('statusUpdate', {
      'heartbeatCount': heartbeatCount,
      'telemetryCount': telemetryCount,
      'running': true,
    });
  });

  // Heartbeat timer — every 15 seconds (demo) / 60 seconds (production)
  Timer.periodic(const Duration(seconds: 15), (timer) async {
    heartbeatCount++;
    Logger.info('💓 Background heartbeat #$heartbeatCount', 'BackgroundService');

    if (api.isConfigured) {
      try {
        await heartbeat.sendHeartbeat();
        Logger.info('📤 Background heartbeat sent', 'BackgroundService');
      } catch (e) {
        Logger.warn('❌ Background heartbeat error: $e', 'BackgroundService');
      }
    }

    // Update notification with latest status
    if (service is AndroidServiceInstance) {
      service.setForegroundNotificationInfo(
        title: 'MDM Agent',
        content: 'Heartbeat #$heartbeatCount • Telemetry #$telemetryCount',
      );
    }

    // Notify UI if it's listening
    service.invoke('statusUpdate', {
      'heartbeatCount': heartbeatCount,
      'telemetryCount': telemetryCount,
      'running': true,
    });
  });

  // Telemetry timer — every 30 seconds
  Timer.periodic(const Duration(seconds: 30), (timer) async {
    telemetryCount++;
    Logger.info('📊 Background telemetry #$telemetryCount', 'BackgroundService');

    try {
      final data = await telemetry.collect(deviceId);

      if (api.isConfigured) {
        await telemetry.sendToServer(api, deviceId, data);
        Logger.info('📤 Background telemetry sent', 'BackgroundService');

        // Send location to dedicated endpoint
        if (data.location != null) {
          try {
            await locationSender.send(api, deviceId, data.location!);
            Logger.info('📍 Background location sent', 'BackgroundService');
          } catch (e) {
            Logger.warn('⚠️ Background location send failed: $e', 'BackgroundService');
          }
        }
      }
    } catch (e) {
      Logger.warn('❌ Background telemetry error: $e', 'BackgroundService');
    }
  });

  // Collect telemetry immediately on start
  try {
    final data = await telemetry.collect(deviceId);
    telemetryCount++;
    if (api.isConfigured) {
      await telemetry.sendToServer(api, deviceId, data);
      if (data.location != null) {
        await locationSender.send(api, deviceId, data.location!);
      }
    }
    Logger.info('Initial background telemetry collected and sent', 'BackgroundService');
  } catch (e) {
    Logger.warn('Initial background telemetry failed: $e', 'BackgroundService');
  }

  // Send initial heartbeat
  try {
    await heartbeat.sendHeartbeat();
    heartbeatCount++;
    Logger.info('Initial background heartbeat sent', 'BackgroundService');
  } catch (e) {
    Logger.warn('Initial background heartbeat failed: $e', 'BackgroundService');
  }
}
