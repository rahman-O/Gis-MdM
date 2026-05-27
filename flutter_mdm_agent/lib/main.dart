import 'package:flutter/material.dart';
import 'core/storage/local_db.dart';
import 'core/storage/secure_storage.dart';
import 'core/utils/logger.dart';

/// MDM Agent entry point.
/// This app runs as a background service with minimal UI.
/// It is installed silently via the MDM profile and starts on boot.
void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize local database
  await LocalDb.initialize();

  // Check enrollment status
  final enrolled = await SecureStorage.isEnrolled();
  Logger.info('Agent started. Enrolled: $enrolled', 'Main');

  // Run minimal app (required for Flutter engine, but no visible UI)
  runApp(const MdmAgentApp());
}

/// Minimal Material app — the agent has no visible UI.
/// It exists only to keep the Flutter engine alive.
class MdmAgentApp extends StatelessWidget {
  const MdmAgentApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      home: const Scaffold(
        body: Center(
          child: Text(
            'MDM Agent Running',
            style: TextStyle(color: Colors.grey),
          ),
        ),
      ),
    );
  }
}
