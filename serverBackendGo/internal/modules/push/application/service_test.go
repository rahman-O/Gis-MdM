package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/push/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubTargets struct{}

func (stubTargets) DeviceIDsByNumbers(context.Context, int64, int64, bool, []string) ([]int64, error) {
	return []int64{1}, nil
}
func (stubTargets) DeviceIDsByGroupNames(context.Context, int64, int64, bool, []string) ([]int64, error) {
	return nil, nil
}
func (stubTargets) AllDeviceIDs(context.Context, int64, int64, bool) ([]int64, error) {
	return nil, nil
}

type stubQueue struct{ n int }

func (q *stubQueue) Enqueue(context.Context, int64, string, string) error {
	q.n++
	return nil
}

func TestSend_requiresPermission(t *testing.T) {
	q := &stubQueue{}
	svc := NewService(stubTargets{}, q)
	err := svc.Send(context.Background(), &platformauth.Principal{CustomerID: 1, ID: 1}, domain.PushRequest{MessageType: "x", DeviceNumbers: []string{"d1"}})
	if err != ErrPermissionDenied {
		t.Fatalf("got %v", err)
	}
}

func TestSend_enqueues(t *testing.T) {
	q := &stubQueue{}
	svc := NewService(stubTargets{}, q)
	p := &platformauth.Principal{CustomerID: 1, ID: 1, Permissions: []string{platformauth.PermPushAPI}}
	if err := svc.Send(context.Background(), p, domain.PushRequest{MessageType: "configUpdated", DeviceNumbers: []string{"hmdm-001"}}); err != nil {
		t.Fatal(err)
	}
	if q.n != 1 {
		t.Fatalf("want 1 enqueue, got %d", q.n)
	}
}
