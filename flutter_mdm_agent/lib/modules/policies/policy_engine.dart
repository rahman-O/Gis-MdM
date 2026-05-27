import '../../core/utils/logger.dart';
import '../sync/sync_response.dart';

/// Core policy engine that reads SyncResponse and dispatches to enforcers.
///
/// The engine iterates through all registered [PolicyEnforcer] instances
/// and applies the policies defined in the sync response. Each enforcer
/// handles its own domain (restrictions, hardware, kiosk, etc.).
class PolicyEngine {
  final List<PolicyEnforcer> _enforcers;

  PolicyEngine(this._enforcers);

  /// Apply all policies from a sync response.
  ///
  /// Iterates through each enforcer and calls [PolicyEnforcer.enforce].
  /// Failures in one enforcer do not prevent others from executing.
  Future<void> applyPolicies(SyncResponse response) async {
    Logger.info(
      'Applying policies from config ${response.configurationId}',
      'PolicyEngine',
    );

    for (final enforcer in _enforcers) {
      try {
        await enforcer.enforce(response);
        Logger.debug('Enforcer "${enforcer.name}" applied successfully', 'PolicyEngine');
      } catch (e, stack) {
        Logger.error(
          'Policy enforcer "${enforcer.name}" failed',
          e,
          stack,
          'PolicyEngine',
        );
      }
    }

    Logger.info('All policy enforcers processed', 'PolicyEngine');
  }

  /// Remove all policies (reset to default).
  ///
  /// Calls [PolicyEnforcer.clear] on each enforcer to revert
  /// any applied policies back to their default state.
  Future<void> clearAll() async {
    Logger.info('Clearing all policies', 'PolicyEngine');

    for (final enforcer in _enforcers) {
      try {
        await enforcer.clear();
        Logger.debug('Enforcer "${enforcer.name}" cleared', 'PolicyEngine');
      } catch (e, stack) {
        Logger.error(
          'Failed to clear enforcer "${enforcer.name}"',
          e,
          stack,
          'PolicyEngine',
        );
      }
    }
  }

  /// Get the list of registered enforcers.
  List<PolicyEnforcer> get enforcers => List.unmodifiable(_enforcers);
}

/// Abstract base class for policy enforcers.
///
/// Each enforcer is responsible for a specific domain of device policy
/// (e.g., restrictions, hardware settings, kiosk mode).
abstract class PolicyEnforcer {
  /// Human-readable name for logging and identification.
  String get name;

  /// Apply policies from the given [SyncResponse].
  Future<void> enforce(SyncResponse response);

  /// Clear/revert all policies managed by this enforcer.
  Future<void> clear();
}
