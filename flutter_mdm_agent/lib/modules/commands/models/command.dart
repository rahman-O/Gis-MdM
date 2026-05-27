/// Represents a remote command received from the MDM server.
class RemoteCommand {
  final int id;
  final String messageType;
  final String payload;
  final DateTime receivedAt;
  CommandStatus status;

  RemoteCommand({
    required this.id,
    required this.messageType,
    required this.payload,
    required this.receivedAt,
    this.status = CommandStatus.pending,
  });

  factory RemoteCommand.fromJson(Map<String, dynamic> json) {
    return RemoteCommand(
      id: json['id'] as int? ?? 0,
      messageType: json['messageType'] as String? ?? json['type'] as String? ?? '',
      payload: json['payload'] as String? ?? json['message'] as String? ?? '',
      receivedAt: json['receivedAt'] != null
          ? DateTime.parse(json['receivedAt'] as String)
          : DateTime.now(),
      status: CommandStatus.pending,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'messageType': messageType,
      'payload': payload,
      'receivedAt': receivedAt.toIso8601String(),
      'status': status.name,
    };
  }

  @override
  String toString() =>
      'RemoteCommand(id=$id, type=$messageType, status=${status.name})';
}

/// Status of a remote command in the processing pipeline.
enum CommandStatus {
  /// Command received but not yet processed.
  pending,

  /// Command is currently being executed.
  executing,

  /// Command completed successfully.
  completed,

  /// Command execution failed.
  failed,
}
