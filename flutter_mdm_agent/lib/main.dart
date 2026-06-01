import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_background_service/flutter_background_service.dart';

import 'background_service.dart';
import 'core/config/constants.dart';
import 'core/storage/local_db.dart';
import 'core/utils/logger.dart';

/// MDM Agent entry point.
void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await LocalDb.initialize();

  Logger.info('MDM Agent v${AgentConstants.appVersion} starting...', 'Main');

  // Initialize and start the background service
  await initializeBackgroundService();

  runApp(const MdmAgentApp());
}

class MdmAgentApp extends StatefulWidget {
  const MdmAgentApp({super.key});

  @override
  State<MdmAgentApp> createState() => _MdmAgentAppState();
}

class _MdmAgentAppState extends State<MdmAgentApp> {
  int _heartbeatCount = 0;
  int _telemetryCount = 0;
  bool _serviceRunning = false;
  String _status = 'Connecting to service...';
  final List<String> _logs = [];
  Timer? _statusPollTimer;

  @override
  void initState() {
    super.initState();
    _connectToService();
  }

  void _connectToService() {
    final service = FlutterBackgroundService();

    // Listen for status updates from the background service
    service.on('statusUpdate').listen((event) {
      if (event != null && mounted) {
        setState(() {
          _heartbeatCount = event['heartbeatCount'] as int? ?? _heartbeatCount;
          _telemetryCount = event['telemetryCount'] as int? ?? _telemetryCount;
          _serviceRunning = event['running'] as bool? ?? false;
          _status = _serviceRunning ? 'Running (Background)' : 'Stopped';
        });
        _addLog(
          '💓 HB: $_heartbeatCount | 📊 Tel: $_telemetryCount',
        );
      }
    });

    // Poll for status periodically
    _statusPollTimer = Timer.periodic(const Duration(seconds: 5), (_) {
      service.invoke('status');
    });

    // Request initial status
    service.invoke('status');

    // Check if service is running
    _checkServiceRunning();
  }

  Future<void> _checkServiceRunning() async {
    final service = FlutterBackgroundService();
    final running = await service.isRunning();
    if (mounted) {
      setState(() {
        _serviceRunning = running;
        _status = running ? 'Running (Background)' : 'Stopped';
      });
      if (running) {
        _addLog('✅ Background service is active');
      } else {
        _addLog('⚠️ Background service not running — starting...');
        // Auto-start if not running
        service.startService();
      }
    }
  }

  void _addLog(String message) {
    final time = DateTime.now().toString().substring(11, 19);
    setState(() {
      _logs.insert(0, '[$time] $message');
      if (_logs.length > 50) _logs.removeLast();
    });
  }

  @override
  void dispose() {
    _statusPollTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: ThemeData.dark(),
      home: Scaffold(
        appBar: AppBar(
          title: const Text('MDM Agent', style: TextStyle(fontSize: 14)),
          actions: [
            Chip(
              label: Text(_status, style: const TextStyle(fontSize: 10)),
              backgroundColor: _serviceRunning ? Colors.green : Colors.orange,
            ),
            const SizedBox(width: 8),
          ],
        ),
        body: Column(
          children: [
            // Service status card
            _buildServiceStatusCard(),
            // Stats
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Row(
                children: [
                  _statChip('💓 Heartbeats', '$_heartbeatCount'),
                  const SizedBox(width: 8),
                  _statChip('📊 Telemetry', '$_telemetryCount'),
                ],
              ),
            ),
            const Divider(height: 1),
            // Live logs
            Expanded(
              child: ListView.builder(
                padding: const EdgeInsets.all(8),
                itemCount: _logs.length,
                itemBuilder: (_, i) => Text(
                  _logs[i],
                  style: const TextStyle(fontSize: 11, fontFamily: 'monospace'),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildServiceStatusCard() {
    return Card(
      margin: const EdgeInsets.all(8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(
                  _serviceRunning ? Icons.check_circle : Icons.error,
                  color: _serviceRunning ? Colors.green : Colors.orange,
                  size: 20,
                ),
                const SizedBox(width: 8),
                const Text(
                  'Background Service',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              _serviceRunning
                  ? 'Service is running in background. Heartbeats and telemetry '
                    'will continue even when this UI is closed.'
                  : 'Service is not running. Tap restart to start it.',
              style: const TextStyle(fontSize: 11, color: Colors.grey),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                ElevatedButton.icon(
                  onPressed: _serviceRunning ? null : _startService,
                  icon: const Icon(Icons.play_arrow, size: 16),
                  label: const Text('Start', style: TextStyle(fontSize: 11)),
                ),
                const SizedBox(width: 8),
                ElevatedButton.icon(
                  onPressed: _serviceRunning ? _stopService : null,
                  icon: const Icon(Icons.stop, size: 16),
                  label: const Text('Stop', style: TextStyle(fontSize: 11)),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.red.shade700,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _startService() {
    final service = FlutterBackgroundService();
    service.startService();
    _addLog('▶️ Service start requested');
    Future.delayed(const Duration(seconds: 2), _checkServiceRunning);
  }

  void _stopService() {
    final service = FlutterBackgroundService();
    service.invoke('stopService');
    setState(() {
      _serviceRunning = false;
      _status = 'Stopped';
    });
    _addLog('⏹️ Service stop requested');
  }

  Widget _statChip(String label, String value) {
    return Chip(
      label: Text('$label: $value', style: const TextStyle(fontSize: 10)),
      padding: EdgeInsets.zero,
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    );
  }
}
