import 'package:battery_plus/battery_plus.dart';

import '../../../core/utils/logger.dart';
import '../telemetry_data.dart';

/// Collects battery level and charging state.
class BatteryCollector {
  final Battery _battery = Battery();

  /// Collect current battery information.
  Future<BatteryInfo> collect() async {
    try {
      final level = await _battery.batteryLevel;
      final state = await _battery.batteryState;

      return BatteryInfo(
        level: level,
        chargingState: _stateToString(state),
      );
    } catch (e) {
      Logger.warn('Battery collection failed: $e', 'BatteryCollector');
      return BatteryInfo(level: -1, chargingState: 'unknown');
    }
  }

  String _stateToString(BatteryState state) {
    switch (state) {
      case BatteryState.charging:
        return 'charging';
      case BatteryState.discharging:
        return 'discharging';
      case BatteryState.full:
        return 'full';
      case BatteryState.connectedNotCharging:
        return 'connected_not_charging';
      default:
        return 'unknown';
    }
  }
}
