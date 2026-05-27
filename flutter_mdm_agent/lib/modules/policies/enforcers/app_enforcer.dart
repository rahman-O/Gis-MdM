import 'package:flutter/services.dart';

import '../../../core/utils/logger.dart';
import '../../sync/sync_response.dart';
import '../policy_engine.dart';

/// Manages application whitelist/blacklist policies.
///
/// Reads the [SyncResponse.permissive] flag and [SyncResponse.applications]
/// list to determine which apps should be allowed or blocked on the device.
class AppEnforcer implements PolicyEnforcer {
  static const _channel = MethodChannel('com.mdm.agent/app_policy');

  bool? _currentPermissive;
  Set<String> _allowedPackages = {};

  @override
  String get name => 'Applications';

  @override
  Future<void> enforce(SyncResponse response) async {
    final permissive = response.permissive;
    final applications = response.applications;

    // Extract package names from the application list
    final desiredPackages = applications
        .map((app) => app.pkg)
        .where((pkg) => pkg.isNotEmpty)
        .toSet();

    // Determine if policy mode changed
    final modeChanged = permissive != _currentPermissive;
    final packagesChanged = !_setsEqual(desiredPackages, _allowedPackages);

    if (!modeChanged && !packagesChanged) {
      Logger.debug('App policy unchanged', 'AppEnforcer');
      return;
    }

    try {
      if (permissive == true) {
        // Permissive mode: blacklist approach (block specific apps)
        await _applyPermissiveMode(desiredPackages);
      } else if (permissive == false) {
        // Non-permissive mode: whitelist approach (only allow listed apps)
        await _applyRestrictiveMode(desiredPackages);
      } else {
        // No policy specified, clear any existing restrictions
        await _clearAppRestrictions();
      }

      _currentPermissive = permissive;
      _allowedPackages = desiredPackages;
      Logger.info(
        'App policy applied: permissive=$permissive, packages=${desiredPackages.length}',
        'AppEnforcer',
      );
    } catch (e, stack) {
      Logger.error('Failed to apply app policy', e, stack, 'AppEnforcer');
    }
  }

  @override
  Future<void> clear() async {
    try {
      await _clearAppRestrictions();
      _currentPermissive = null;
      _allowedPackages = {};
      Logger.info('App policy cleared', 'AppEnforcer');
    } catch (e, stack) {
      Logger.error('Failed to clear app policy', e, stack, 'AppEnforcer');
    }
  }

  /// Get the current permissive mode.
  bool? get currentPermissive => _currentPermissive;

  /// Get the currently allowed packages.
  Set<String> get allowedPackages => Set.unmodifiable(_allowedPackages);

  /// Apply permissive mode — all apps allowed except those not in the list.
  Future<void> _applyPermissiveMode(Set<String> packages) async {
    await _channel.invokeMethod('setAppPolicy', {
      'mode': 'permissive',
      'packages': packages.toList(),
    });
    Logger.debug('Permissive mode applied with ${packages.length} packages', 'AppEnforcer');
  }

  /// Apply restrictive mode — only listed apps are allowed.
  Future<void> _applyRestrictiveMode(Set<String> packages) async {
    await _channel.invokeMethod('setAppPolicy', {
      'mode': 'restrictive',
      'packages': packages.toList(),
    });
    Logger.debug('Restrictive mode applied with ${packages.length} packages', 'AppEnforcer');
  }

  /// Clear all app restrictions.
  Future<void> _clearAppRestrictions() async {
    await _channel.invokeMethod('clearAppPolicy');
    Logger.debug('App restrictions cleared', 'AppEnforcer');
  }

  /// Compare two sets for equality.
  bool _setsEqual(Set<String> a, Set<String> b) {
    if (a.length != b.length) return false;
    return a.containsAll(b);
  }
}
