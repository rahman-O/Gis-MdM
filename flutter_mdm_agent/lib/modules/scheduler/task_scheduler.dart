import '../../core/utils/logger.dart';

/// Schedules commands for execution at specific times.
class TaskScheduler {
  final List<ScheduledTask> _tasks = [];

  /// Add a scheduled task.
  void schedule(ScheduledTask task) {
    _tasks.add(task);
    _tasks.sort((a, b) => a.executeAt.compareTo(b.executeAt));
    Logger.info(
        'Task scheduled: ${task.commandType} at ${task.executeAt}', 'Scheduler');
  }

  /// Check and execute any due tasks.
  Future<void> tick() async {
    final now = DateTime.now();
    final due =
        _tasks.where((t) => t.executeAt.isBefore(now) && !t.executed).toList();
    for (final task in due) {
      task.executed = true;
      Logger.info(
          'Executing scheduled task: ${task.commandType}', 'Scheduler');
      // Dispatch to command handler
    }
    _tasks.removeWhere((t) => t.executed);
  }

  /// Get all pending (not yet executed) tasks.
  List<ScheduledTask> get pendingTasks =>
      _tasks.where((t) => !t.executed).toList();

  /// Cancel a scheduled task by ID.
  bool cancel(String taskId) {
    final removed = _tasks.length;
    _tasks.removeWhere((t) => t.id == taskId);
    return _tasks.length < removed;
  }
}

/// Represents a task scheduled for future execution.
class ScheduledTask {
  final String id;
  final String commandType;
  final Map<String, dynamic> payload;
  final DateTime executeAt;
  final bool recurring;
  final Duration? interval;
  bool executed = false;

  ScheduledTask({
    required this.id,
    required this.commandType,
    required this.payload,
    required this.executeAt,
    this.recurring = false,
    this.interval,
  });
}
