import '../../../core/utils/logger.dart';
import '../../../platform/device_owner_channel.dart';
import '../../sync/sync_response.dart';
import '../policy_engine.dart';

/// Applies Android UserRestrictions via DeviceOwnerChannel.
///
/// Parses the comma-separated restrictions string from [SyncResponse]
/// and adds/removes restrictions to match the desired state.
class RestrictionsEnforcer implements PolicyEnforcer {
  Set<String> _currentRestrictions = {};

  @override
  String get name => 'Restrictions';

  @override
  Future<void> enforce(SyncResponse response) async {
    final desired = _parseRestrictions(response.restrictions);

    // Add new restrictions that aren't currently applied
    final toAdd = desired.difference(_currentRestrictions);
    for (final restriction in toAdd) {
      try {
        await DeviceOwnerChannel.addUserRestriction(restriction);
        Logger.debug('Added restriction: $restriction', 'RestrictionsEnforcer');
      } catch (e, stack) {
        Logger.error(
          'Failed to add restriction: $restriction',
          e,
          stack,
          'RestrictionsEnforcer',
        );
      }
    }

    // Remove restrictions that are no longer desired
    final toRemove = _currentRestrictions.difference(desired);
    for (final restriction in toRemove) {
      try {
        await DeviceOwnerChannel.clearUserRestriction(restriction);
        Logger.debug('Removed restriction: $restriction', 'RestrictionsEnforcer');
      } catch (e, stack) {
        Logger.error(
          'Failed to remove restriction: $restriction',
          e,
          stack,
          'RestrictionsEnforcer',
        );
      }
    }

    _currentRestrictions = desired;
    Logger.info(
      'Restrictions updated: ${_currentRestrictions.length} active',
      'RestrictionsEnforcer',
    );
  }

  @override
  Future<void> clear() async {
    for (final restriction in _currentRestrictions) {
      try {
        await DeviceOwnerChannel.clearUserRestriction(restriction);
      } catch (e, stack) {
        Logger.error(
          'Failed to clear restriction: $restriction',
          e,
          stack,
          'RestrictionsEnforcer',
        );
      }
    }
    _currentRestrictions = {};
    Logger.info('All restrictions cleared', 'RestrictionsEnforcer');
  }

  /// Get the currently applied restrictions.
  Set<String> get currentRestrictions => Set.unmodifiable(_currentRestrictions);

  /// Parse comma-separated restrictions string into a Set.
  Set<String> _parseRestrictions(String? restrictions) {
    if (restrictions == null || restrictions.trim().isEmpty) return {};
    return restrictions
        .split(',')
        .map((r) => r.trim())
        .where((r) => r.isNotEmpty)
        .toSet();
  }
}
