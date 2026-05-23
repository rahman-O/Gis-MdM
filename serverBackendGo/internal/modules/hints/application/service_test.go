package application

import (
	"context"
	"testing"

	"github.com/gis-mdm/server-backend-go/internal/modules/hints/domain"
	platformauth "github.com/gis-mdm/server-backend-go/internal/platform/auth"
)

type stubRepo struct {
	history map[int64][]string
	catalog []string
}

func (s *stubRepo) GetHistory(_ context.Context, userID int64) ([]string, error) {
	if s.history == nil {
		return []string{}, nil
	}
	return append([]string(nil), s.history[userID]...), nil
}

func (s *stubRepo) MarkShown(_ context.Context, userID int64, hintKey string) error {
	if s.history == nil {
		s.history = map[int64][]string{}
	}
	for _, k := range s.history[userID] {
		if k == hintKey {
			return nil
		}
	}
	s.history[userID] = append(s.history[userID], hintKey)
	return nil
}

func (s *stubRepo) Enable(_ context.Context, userID int64) error {
	if s.history == nil {
		s.history = map[int64][]string{}
	}
	s.history[userID] = nil
	return nil
}

func (s *stubRepo) Disable(_ context.Context, userID int64) error {
	if s.history == nil {
		s.history = map[int64][]string{}
	}
	s.history[userID] = append([]string(nil), s.catalog...)
	return nil
}

func testPrincipal() *platformauth.Principal {
	return &platformauth.Principal{ID: 1, AuthLoaded: true}
}

func TestGetHistory_empty(t *testing.T) {
	svc := NewService(&stubRepo{})
	keys, err := svc.GetHistory(context.Background(), testPrincipal())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Fatalf("want empty, got %v", keys)
	}
}

func TestGetHistory_populated(t *testing.T) {
	svc := NewService(&stubRepo{history: map[int64][]string{1: {"hint.step.1"}}})
	keys, err := svc.GetHistory(context.Background(), testPrincipal())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0] != "hint.step.1" {
		t.Fatalf("got %v", keys)
	}
}

func TestMarkShown_duplicateIdempotent(t *testing.T) {
	repo := &stubRepo{history: map[int64][]string{}}
	svc := NewService(repo)
	p := testPrincipal()
	if err := svc.MarkShown(context.Background(), p, "hint.step.1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.MarkShown(context.Background(), p, "hint.step.1"); err != nil {
		t.Fatal(err)
	}
	if len(repo.history[1]) != 1 {
		t.Fatalf("want one key, got %v", repo.history[1])
	}
}

func TestMarkShown_emptyKey(t *testing.T) {
	svc := NewService(&stubRepo{})
	err := svc.MarkShown(context.Background(), testPrincipal(), "  ")
	if err != domain.ErrEmptyHintKey {
		t.Fatalf("want empty key err, got %v", err)
	}
}

func TestEnable_clearsHistory(t *testing.T) {
	repo := &stubRepo{history: map[int64][]string{1: {"hint.step.1"}}}
	svc := NewService(repo)
	if err := svc.Enable(context.Background(), testPrincipal()); err != nil {
		t.Fatal(err)
	}
	keys, _ := svc.GetHistory(context.Background(), testPrincipal())
	if len(keys) != 0 {
		t.Fatalf("want cleared, got %v", keys)
	}
}

func TestDisable_insertsCatalog(t *testing.T) {
	catalog := []string{"hint.step.1", "hint.step.2", "hint.step.3", "hint.step.4"}
	svc := NewService(&stubRepo{catalog: catalog})
	if err := svc.Disable(context.Background(), testPrincipal()); err != nil {
		t.Fatal(err)
	}
	keys, err := svc.GetHistory(context.Background(), testPrincipal())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 4 {
		t.Fatalf("want 4 keys, got %v", keys)
	}
}

func TestGetHistory_unauthenticated(t *testing.T) {
	svc := NewService(&stubRepo{})
	_, err := svc.GetHistory(context.Background(), nil)
	if err != ErrUnauthenticated {
		t.Fatalf("got %v", err)
	}
}
