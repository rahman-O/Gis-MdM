package application

import "testing"

func TestBoolTrueOnly(t *testing.T) {
	if BoolTrueOnly(false) != nil {
		t.Fatal("expected nil for false")
	}
	if BoolTrueOnly(true) == nil || !*BoolTrueOnly(true) {
		t.Fatal("expected true pointer")
	}
}
