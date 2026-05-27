/// Static constants for the MDM Agent.
class AgentConstants {
  AgentConstants._();

  static const String appName = 'MDM Agent';
  static const String appVersion = '1.0.0';
  static const int appVersionCode = 1;

  // Default intervals (seconds) — overridden by server config
  static const int defaultHeartbeatIntervalSec = 60;
  static const int defaultTelemetryIntervalSec = 300; // 5 minutes
  static const int defaultSyncIntervalSec = 300; // 5 minutes
  static const int defaultCommandPollIntervalSec = 30;

  // Retry
  static const int maxRetryAttempts = 3;
  static const int retryDelayMs = 5000;

  // Offline queue limits
  static const int maxOfflineEvents = 10000;
  static const int maxOfflineLocations = 10000;

  // Notification channel
  static const String foregroundChannelId = 'mdm_agent_foreground';
  static const String foregroundChannelName = 'MDM Agent Service';

  // Method channel
  static const String methodChannelName = 'com.gismdm.mdm_agent/device_owner';
}
