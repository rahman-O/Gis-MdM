/// Internal event bus for pub/sub within the agent.
class EventBus {
  static final _listeners =
      <String, List<void Function(Map<String, dynamic>)>>{};

  /// Subscribe to an event type.
  static void on(String eventType, void Function(Map<String, dynamic>) handler) {
    _listeners.putIfAbsent(eventType, () => []).add(handler);
  }

  /// Emit an event to all subscribers.
  static void emit(String eventType, Map<String, dynamic> data) {
    final handlers = _listeners[eventType];
    if (handlers != null) {
      for (final handler in handlers) {
        handler(data);
      }
    }
  }

  /// Remove a specific listener for an event type.
  static void off(
      String eventType, void Function(Map<String, dynamic>) handler) {
    _listeners[eventType]?.remove(handler);
  }

  /// Clear all listeners.
  static void clear() => _listeners.clear();
}
