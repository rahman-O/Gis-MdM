import 'dart:async';
import 'package:flutter/material.dart';
import 'core/config/constants.dart';
import 'core/network/api_client.dart';
import 'core/storage/local_db.dart';
import 'core/storage/secure_storage.dart';
import 'core/utils/logger.dart';
import 'modules/heartbeat/heartbeat_service.dart';
import 'modules/telemetry/telemetry_service.dart';
import 'modules/telemetry/telemetry_data.dart';

/// MDM Agent entry point.
void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await LocalDb.initialize();

  Logger.info('MDM Agent v${AgentConstants.appVersion} starting...', 'Main');

  runApp(const MdmAgentApp());
}

class MdmAgentApp extends StatefulWidget {
  const MdmAgentApp({super.key});

  @override
  State<MdmAgentApp> createState() => _MdmAgentAppState();
}

class _MdmAgentAppState extends State<MdmAgentApp> {
  final ApiClient _api = ApiClient();
  final TelemetryService _telemetry = TelemetryService();
  late HeartbeatService _heartbeat;

  Timer? _telemetryTimer;
  Timer? _heartbeatTimer;

  TelemetryData? _lastTelemetry;
  int _heartbeatCount = 0;
  int _telemetryCount = 0;
  String _status = 'Initializing...';
  final List<String> _logs = [];

  @override
  void initState() {
    super.initState();
    _heartbeat = HeartbeatService(_api, 'device-001');
    _startServices();
  }

  Future<void> _startServices() async {
    // Configure API (use stored server URL or default)
    final serverUrl = await SecureStorage.getServerUrl();
    if (serverUrl != null && serverUrl.isNotEmpty) {
      _api.configure(baseUrl: serverUrl);
    }

    _addLog('Services starting...');
    setState(() => _status = 'Running');

    // Start telemetry collection (every 30 seconds for demo)
    _telemetryTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) => _collectTelemetry(),
    );

    // Start heartbeat (every 15 seconds for demo — normally 60s)
    _heartbeatTimer = Timer.periodic(
      const Duration(seconds: 15),
      (_) => _sendHeartbeat(),
    );

    // Collect immediately on start
    await _collectTelemetry();
    _addLog('All services running ✓');
  }

  Future<void> _collectTelemetry() async {
    try {
      final data = await _telemetry.collect('device-001');
      _telemetryCount++;
      setState(() {
        _lastTelemetry = data;
      });
      _addLog('📊 Telemetry #$_telemetryCount collected');

      // Send to server if configured
      if (_api.isConfigured) {
        await _telemetry.sendToServer(_api, 'device-001', data);
        _addLog('📤 Telemetry sent to server');
      }
    } catch (e) {
      _addLog('❌ Telemetry error: $e');
    }
  }

  Future<void> _sendHeartbeat() async {
    _heartbeatCount++;
    _addLog('💓 Heartbeat #$_heartbeatCount');

    if (_api.isConfigured) {
      try {
        await _heartbeat.sendHeartbeat();
        _addLog('📤 Heartbeat sent');
      } catch (e) {
        _addLog('❌ Heartbeat error: $e');
      }
    }
  }

  void _addLog(String message) {
    final time = DateTime.now().toString().substring(11, 19);
    setState(() {
      _logs.insert(0, '[$time] $message');
      if (_logs.length > 50) _logs.removeLast();
    });
    Logger.info(message, 'Main');
  }

  @override
  void dispose() {
    _telemetryTimer?.cancel();
    _heartbeatTimer?.cancel();
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
              backgroundColor: _status == 'Running' ? Colors.green : Colors.orange,
            ),
            const SizedBox(width: 8),
          ],
        ),
        body: Column(
          children: [
            // Telemetry summary
            if (_lastTelemetry != null) _buildTelemetryCard(),
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

  Widget _buildTelemetryCard() {
    final t = _lastTelemetry!;
    return Card(
      margin: const EdgeInsets.all(8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('📱 Device Telemetry', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
            const SizedBox(height: 8),
            Wrap(
              spacing: 12,
              runSpacing: 4,
              children: [
                _infoItem('🔋', '${t.battery.level}% (${t.battery.chargingState})'),
                _infoItem('💾', '${(t.storage.freeBytes / 1024 / 1024 / 1024).toStringAsFixed(1)} GB free'),
                _infoItem('🧠', '${((t.memory.totalBytes - t.memory.freeBytes) / 1024 / 1024 / 1024).toStringAsFixed(1)}/${(t.memory.totalBytes / 1024 / 1024 / 1024).toStringAsFixed(1)} GB'),
                _infoItem('📶', t.network.type),
                _infoItem('🔗', t.network.connected ? 'Connected' : 'Offline'),
                _infoItem('🖥️', t.screen.isOn ? 'ON' : 'OFF'),
                _infoItem('📱', t.system.model),
                _infoItem('🤖', 'Android ${t.system.androidVersion}'),
                _infoItem('⏱️', '${(t.system.uptimeMillis / 3600000).toStringAsFixed(1)}h up'),
              ],
            ),
            if (t.location != null)
              Padding(
                padding: const EdgeInsets.only(top: 4),
                child: Text(
                  '📍 ${t.location!.latitude.toStringAsFixed(4)}, ${t.location!.longitude.toStringAsFixed(4)}',
                  style: const TextStyle(fontSize: 11, color: Colors.grey),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _infoItem(String icon, String value) {
    return Text('$icon $value', style: const TextStyle(fontSize: 11));
  }

  Widget _statChip(String label, String value) {
    return Chip(
      label: Text('$label: $value', style: const TextStyle(fontSize: 10)),
      padding: EdgeInsets.zero,
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    );
  }
}
