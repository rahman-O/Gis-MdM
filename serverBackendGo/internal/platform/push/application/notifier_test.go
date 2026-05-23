package application

import (
	"context"
	"testing"
)

type stubQueue struct {
	calls []struct {
		deviceID    int64
		messageType string
	}
}

func (q *stubQueue) Enqueue(_ context.Context, deviceID int64, messageType, _ string) error {
	q.calls = append(q.calls, struct {
		deviceID    int64
		messageType string
	}{deviceID, messageType})
	return nil
}

type stubLookup struct {
	ids []int64
}

func (l *stubLookup) DeviceIDsByConfiguration(context.Context, int) ([]int64, error) {
	return l.ids, nil
}

func TestNotifier_NotifyConfigurationChanged(t *testing.T) {
	q := &stubQueue{}
	n := NewNotifier(q, &stubLookup{ids: []int64{10, 20}}, nil)
	if err := n.NotifyConfigurationChanged(5); err != nil {
		t.Fatal(err)
	}
	if len(q.calls) != 2 {
		t.Fatalf("calls=%d", len(q.calls))
	}
	if q.calls[0].messageType != TypeConfigUpdated || q.calls[0].deviceID != 10 {
		t.Fatalf("first call=%+v", q.calls[0])
	}
}

func TestNotifier_NotifyAppSettings(t *testing.T) {
	q := &stubQueue{}
	n := NewNotifier(q, nil, nil)
	if err := n.NotifyAppSettings(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if len(q.calls) != 1 || q.calls[0].messageType != TypeAppConfigUpdated {
		t.Fatalf("calls=%+v", q.calls)
	}
}
