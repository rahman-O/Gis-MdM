import 'dart:developer' as dev;

enum LogLevel { debug, info, warn, error }

/// Structured logger for the MDM Agent.
class Logger {
  static LogLevel _level = LogLevel.info;

  static void setLevel(String level) {
    switch (level.toLowerCase()) {
      case 'debug':
        _level = LogLevel.debug;
      case 'warn':
        _level = LogLevel.warn;
      case 'error':
        _level = LogLevel.error;
      default:
        _level = LogLevel.info;
    }
  }

  static void debug(String message, [String? module]) {
    if (_level.index <= LogLevel.debug.index) {
      _log('DEBUG', message, module);
    }
  }

  static void info(String message, [String? module]) {
    if (_level.index <= LogLevel.info.index) {
      _log('INFO', message, module);
    }
  }

  static void warn(String message, [String? module]) {
    if (_level.index <= LogLevel.warn.index) {
      _log('WARN', message, module);
    }
  }

  static void error(String message, [Object? error, StackTrace? stack, String? module]) {
    _log('ERROR', message, module);
    if (error != null) {
      dev.log('  Error: $error', name: 'MDM');
    }
    if (stack != null) {
      dev.log('  Stack: $stack', name: 'MDM');
    }
  }

  static void _log(String level, String message, String? module) {
    final prefix = module != null ? '[$module]' : '';
    dev.log('$level $prefix $message', name: 'MDM');
  }
}
