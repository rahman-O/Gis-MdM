import '../../core/utils/logger.dart';
import 'models/command.dart';

/// Abstract handler interface for processing remote commands.
abstract class CommandHandler {
  /// Process a single command. Returns `true` if handled successfully.
  Future<bool> handle(RemoteCommand command);

  /// Whether this handler can process the given command type.
  bool canHandle(String messageType);
}

/// In-memory queue for processing remote commands in order.
///
/// Commands are dequeued and processed one at a time to prevent
/// conflicts (e.g., simultaneous lock + wipe).
class CommandQueue {
  final List<RemoteCommand> _queue = [];

  /// Number of commands currently in the queue.
  int get length => _queue.length;

  /// Whether the queue is empty.
  bool get isEmpty => _queue.isEmpty;

  /// Add a command to the end of the queue.
  void enqueue(RemoteCommand command) {
    _queue.add(command);
    Logger.debug('Command enqueued: $command', 'CommandQueue');
  }

  /// Remove and return the next command from the front of the queue.
  RemoteCommand? dequeue() {
    if (_queue.isEmpty) return null;
    return _queue.removeAt(0);
  }

  /// Process all queued commands using the provided handler.
  ///
  /// Commands are processed sequentially. Failed commands are marked
  /// but do not block subsequent commands.
  Future<void> processAll(CommandHandler handler) async {
    Logger.info('Processing ${_queue.length} queued commands', 'CommandQueue');

    while (_queue.isNotEmpty) {
      final command = _queue.removeAt(0);

      if (!handler.canHandle(command.messageType)) {
        Logger.warn(
          'No handler for command type: ${command.messageType}',
          'CommandQueue',
        );
        command.status = CommandStatus.failed;
        continue;
      }

      command.status = CommandStatus.executing;

      try {
        final success = await handler.handle(command);
        command.status =
            success ? CommandStatus.completed : CommandStatus.failed;
        Logger.info(
          'Command ${command.id} ${command.status.name}',
          'CommandQueue',
        );
      } catch (e, stack) {
        command.status = CommandStatus.failed;
        Logger.error(
          'Command ${command.id} failed',
          e,
          stack,
          'CommandQueue',
        );
      }
    }
  }
}
