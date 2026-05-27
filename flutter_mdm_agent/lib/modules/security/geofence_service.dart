import 'dart:math';

import '../../core/utils/logger.dart';

/// Manages geofencing zones and triggers events on enter/exit.
class GeofenceService {
  final List<GeofenceZone> _zones = [];
  bool _insideZone = false;

  /// Set geofence zones from server config.
  void setZones(List<GeofenceZone> zones) {
    _zones.clear();
    _zones.addAll(zones);
    Logger.info('Geofence zones updated: ${zones.length}', 'Geofence');
  }

  /// Check if a location is inside any defined zone.
  GeofenceEvent? checkLocation(double lat, double lon) {
    for (final zone in _zones) {
      final inside = _isInsideZone(lat, lon, zone);
      if (inside && !_insideZone) {
        _insideZone = true;
        return GeofenceEvent(type: 'enter', zone: zone, lat: lat, lon: lon);
      } else if (!inside && _insideZone) {
        _insideZone = false;
        return GeofenceEvent(type: 'exit', zone: zone, lat: lat, lon: lon);
      }
    }
    return null;
  }

  bool _isInsideZone(double lat, double lon, GeofenceZone zone) {
    final distance = _haversineDistance(lat, lon, zone.lat, zone.lon);
    return distance <= zone.radiusMeters;
  }

  /// Calculate distance between two coordinates using the Haversine formula.
  /// Returns distance in meters.
  double _haversineDistance(
      double lat1, double lon1, double lat2, double lon2) {
    const earthRadiusMeters = 6371000.0;

    final dLat = _degreesToRadians(lat2 - lat1);
    final dLon = _degreesToRadians(lon2 - lon1);

    final a = sin(dLat / 2) * sin(dLat / 2) +
        cos(_degreesToRadians(lat1)) *
            cos(_degreesToRadians(lat2)) *
            sin(dLon / 2) *
            sin(dLon / 2);

    final c = 2 * atan2(sqrt(a), sqrt(1 - a));

    return earthRadiusMeters * c;
  }

  double _degreesToRadians(double degrees) => degrees * pi / 180;
}

/// Represents a geofence zone with center point and radius.
class GeofenceZone {
  final String id;
  final String name;
  final double lat;
  final double lon;
  final double radiusMeters;

  GeofenceZone({
    required this.id,
    required this.name,
    required this.lat,
    required this.lon,
    required this.radiusMeters,
  });
}

/// Event emitted when a device enters or exits a geofence zone.
class GeofenceEvent {
  final String type; // "enter" or "exit"
  final GeofenceZone zone;
  final double lat;
  final double lon;

  GeofenceEvent({
    required this.type,
    required this.zone,
    required this.lat,
    required this.lon,
  });
}
