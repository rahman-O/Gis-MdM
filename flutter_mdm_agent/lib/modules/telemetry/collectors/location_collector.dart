import 'package:geolocator/geolocator.dart';

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects GPS location data (latitude, longitude, accuracy).
class LocationCollector {
  /// Collect current location.
  ///
  /// Returns `null` if location services are disabled or permission denied.
  Future<LocationInfo?> collect() async {
    try {
      final serviceEnabled = await Geolocator.isLocationServiceEnabled();
      if (!serviceEnabled) {
        Logger.debug('Location services disabled', 'LocationCollector');
        return null;
      }

      var permission = await Geolocator.checkPermission();
      if (permission == LocationPermission.denied) {
        permission = await Geolocator.requestPermission();
        if (permission == LocationPermission.denied) {
          Logger.debug('Location permission denied', 'LocationCollector');
          return null;
        }
      }

      if (permission == LocationPermission.deniedForever) {
        Logger.debug('Location permission permanently denied', 'LocationCollector');
        return null;
      }

      final position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );

      return LocationInfo(
        latitude: position.latitude,
        longitude: position.longitude,
        accuracy: position.accuracy,
        timestamp: position.timestamp.millisecondsSinceEpoch,
      );
    } catch (e) {
      Logger.warn('Location collection failed: $e', 'LocationCollector');
      return null;
    }
  }
}
