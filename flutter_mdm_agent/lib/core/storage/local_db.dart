import 'package:hive_flutter/hive_flutter.dart';
import '../utils/logger.dart';

/// Local database using Hive for offline data, events queue, and cached config.
class LocalDb {
  static const String _configBox = 'agent_config';
  static const String _eventsBox = 'events_queue';
  static const String _telemetryBox = 'telemetry_queue';
  static const String _commandsBox = 'commands_queue';
  static const String _locationsBox = 'locations_queue';

  static Future<void> initialize() async {
    await Hive.initFlutter();
    await Hive.openBox<Map>(_configBox);
    await Hive.openBox<Map>(_eventsBox);
    await Hive.openBox<Map>(_telemetryBox);
    await Hive.openBox<Map>(_commandsBox);
    await Hive.openBox<Map>(_locationsBox);
    Logger.info('Local DB initialized', 'LocalDb');
  }

  // Config cache
  static Box<Map> get configBox => Hive.box<Map>(_configBox);

  static Future<void> saveConfig(Map<String, dynamic> config) async {
    await configBox.put('current', config);
  }

  static Map<String, dynamic>? getConfig() {
    final raw = configBox.get('current');
    return raw?.cast<String, dynamic>();
  }

  // Events queue
  static Box<Map> get eventsBox => Hive.box<Map>(_eventsBox);

  static Future<void> addEvent(Map<String, dynamic> event) async {
    await eventsBox.add(event);
  }

  static List<Map<String, dynamic>> getPendingEvents({int limit = 100}) {
    return eventsBox.values.take(limit).map((e) => e.cast<String, dynamic>()).toList();
  }

  static Future<void> clearEvents(int count) async {
    final keys = eventsBox.keys.take(count).toList();
    await eventsBox.deleteAll(keys);
  }

  // Telemetry queue
  static Box<Map> get telemetryBox => Hive.box<Map>(_telemetryBox);

  static Future<void> addTelemetry(Map<String, dynamic> data) async {
    await telemetryBox.add(data);
  }

  static List<Map<String, dynamic>> getPendingTelemetry({int limit = 50}) {
    return telemetryBox.values.take(limit).map((e) => e.cast<String, dynamic>()).toList();
  }

  static Future<void> clearTelemetry(int count) async {
    final keys = telemetryBox.keys.take(count).toList();
    await telemetryBox.deleteAll(keys);
  }

  // Locations queue
  static Box<Map> get locationsBox => Hive.box<Map>(_locationsBox);

  static Future<void> addLocation(Map<String, dynamic> location) async {
    await locationsBox.add(location);
  }

  static List<Map<String, dynamic>> getPendingLocations({int limit = 100}) {
    return locationsBox.values.take(limit).map((e) => e.cast<String, dynamic>()).toList();
  }

  static Future<void> clearLocations(int count) async {
    final keys = locationsBox.keys.take(count).toList();
    await locationsBox.deleteAll(keys);
  }
}
