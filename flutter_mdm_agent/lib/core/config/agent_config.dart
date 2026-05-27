import 'constants.dart';

/// Dynamic configuration received from the server.
/// All intervals can be overridden remotely without app update.
class AgentConfig {
  final int heartbeatIntervalSec;
  final int telemetryIntervalSec;
  final int syncIntervalSec;
  final int commandPollIntervalSec;
  final bool fcmEnabled;
  final bool geofenceEnabled;
  final bool locationTrackingEnabled;
  final int locationIntervalSec;
  final List<String> enabledModules;
  final String logLevel;

  const AgentConfig({
    this.heartbeatIntervalSec = AgentConstants.defaultHeartbeatIntervalSec,
    this.telemetryIntervalSec = AgentConstants.defaultTelemetryIntervalSec,
    this.syncIntervalSec = AgentConstants.defaultSyncIntervalSec,
    this.commandPollIntervalSec = AgentConstants.defaultCommandPollIntervalSec,
    this.fcmEnabled = false,
    this.geofenceEnabled = false,
    this.locationTrackingEnabled = true,
    this.locationIntervalSec = 30,
    this.enabledModules = const [],
    this.logLevel = 'info',
  });

  factory AgentConfig.fromJson(Map<String, dynamic> json) {
    return AgentConfig(
      heartbeatIntervalSec: json['heartbeatIntervalSec'] as int? ?? AgentConstants.defaultHeartbeatIntervalSec,
      telemetryIntervalSec: json['telemetryIntervalSec'] as int? ?? AgentConstants.defaultTelemetryIntervalSec,
      syncIntervalSec: json['syncIntervalSec'] as int? ?? AgentConstants.defaultSyncIntervalSec,
      commandPollIntervalSec: json['commandPollIntervalSec'] as int? ?? AgentConstants.defaultCommandPollIntervalSec,
      fcmEnabled: json['fcmEnabled'] as bool? ?? false,
      geofenceEnabled: json['geofenceEnabled'] as bool? ?? false,
      locationTrackingEnabled: json['locationTrackingEnabled'] as bool? ?? true,
      locationIntervalSec: json['locationIntervalSec'] as int? ?? 30,
      enabledModules: (json['enabledModules'] as List<dynamic>?)?.cast<String>() ?? [],
      logLevel: json['logLevel'] as String? ?? 'info',
    );
  }

  Map<String, dynamic> toJson() => {
    'heartbeatIntervalSec': heartbeatIntervalSec,
    'telemetryIntervalSec': telemetryIntervalSec,
    'syncIntervalSec': syncIntervalSec,
    'commandPollIntervalSec': commandPollIntervalSec,
    'fcmEnabled': fcmEnabled,
    'geofenceEnabled': geofenceEnabled,
    'locationTrackingEnabled': locationTrackingEnabled,
    'locationIntervalSec': locationIntervalSec,
    'enabledModules': enabledModules,
    'logLevel': logLevel,
  };
}
