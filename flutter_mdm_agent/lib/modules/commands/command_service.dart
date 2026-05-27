import '../../core/network/api_client.dart';
import '../../core/network/endpoints.dart';
import '../../core/utils/logger.dart';
import 'models/command.dart';

/// Service for polling and acknowledging remote commands from the server.
class CommandService {
  final ApiClient _api;

  CommandService(this._api);

  /// Poll the server for pending commands for this device.
  ///
  /// Returns a list of [RemoteCommand] objects to be processed.
  Future<List<RemoteCommand>> pollCommands(String deviceId) async {
    try {
      Logger.debug('Polling commands for $deviceId', 'Commands');

      final response = await _api.get(
        Endpoints.notificationPolling(deviceId),
      );

      if (response.statusCode == 200 && response.data != null) {
        final List<RemoteCommand> commands = [];

        if (response.data is List) {
          for (final item in response.data as List) {
            if (item is Map<String, dynamic>) {
              commands.add(RemoteCommand.fromJson(item));
            }
          }
        } else if (response.data is Map<String, dynamic>) {
          // Single command response
          commands.add(
            RemoteCommand.fromJson(response.data as Map<String, dynamic>),
          );
        }

        Logger.info(
          'Polled ${commands.length} commands',
          'Commands',
        );
        return commands;
      }

      return [];
    } catch (e, stack) {
      Logger.error('Poll commands failed', e, stack, 'Commands');
      return [];
    }
  }

  /// Acknowledge a command as received/processed.
  Future<void> acknowledge(String deviceId, int commandId) async {
    try {
      await _api.post(
        Endpoints.syncInfo,
        data: {
          'deviceId': deviceId,
          'commandId': commandId,
          'status': 'delivered',
          'timestamp': DateTime.now().millisecondsSinceEpoch,
        },
      );
      Logger.debug('Command $commandId acknowledged', 'Commands');
    } catch (e, stack) {
      Logger.error(
        'Failed to acknowledge command $commandId',
        e,
        stack,
        'Commands',
      );
    }
  }
}
