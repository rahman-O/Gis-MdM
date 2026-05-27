import 'package:flutter_test/flutter_test.dart';
import 'package:mdm_agent/main.dart';

void main() {
  testWidgets('Agent app starts', (WidgetTester tester) async {
    await tester.pumpWidget(const MdmAgentApp());
    expect(find.text('MDM Agent Running'), findsOneWidget);
  });
}
